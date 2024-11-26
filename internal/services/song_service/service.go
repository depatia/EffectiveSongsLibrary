package songservice

import (
	"SongsLibrary/internal/models"
	externalservice "SongsLibrary/internal/services/external_service"
	"SongsLibrary/internal/storage/sterrors"
	"SongsLibrary/pkg/logger/sl"
	"SongsLibrary/pkg/validator"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type SongStore struct {
	log             *slog.Logger
	songCreater     SongCreater
	songGetter      SongGetter
	songProvider    SongProvider
	externalService *externalservice.ExternalStore
}

func New(
	log *slog.Logger,
	songCreater SongCreater,
	songGetter SongGetter,
	songProvider SongProvider,
	externalService *externalservice.ExternalStore) *SongStore {
	return &SongStore{log: log, songCreater: songCreater, songGetter: songGetter, songProvider: songProvider}
}

type SongCreater interface {
	CreateSong(ctx context.Context, song *models.Song) error
}

type SongGetter interface {
	GetSongs(ctx context.Context, limit, offset int, filter map[string]string) ([]*models.Song, error)
	GetText(ctx context.Context, ID int64) (string, error)
}

type SongProvider interface {
	UpdateSong(ctx context.Context, song *models.Song) error
	DelSong(ctx context.Context, ID int64) error
}

func (s *SongStore) CreateSong(ctx context.Context, song *models.CreateSongReq) (int, error) {
	const op = "services.song_service.CreateSong"
	log := s.log.With(
		slog.String("Operation:", op),
	)

	log.Info("creating song")

	if err := validator.ValidateStruct(song); err != nil {
		log.Error("failed to create song", sl.Err(err))

		return http.StatusBadRequest, fmt.Errorf("%s: %w", op, err)
	}

	songDetail, status, err := s.externalService.GetSongDetail(ctx, song.Group, song.Song)
	if err != nil {
		log.Error("failed to get response from external service", sl.Err(err))

		return status, err
	}

	err = s.songCreater.CreateSong(ctx, songDetail)
	if err != nil {
		log.Error("failed to create song", sl.Err(err))

		if errors.Is(err, sterrors.ErrSongAlreadyExists) {
			return http.StatusConflict, fmt.Errorf("%s: %w", op, err)
		}

		return http.StatusInternalServerError, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("song created")

	return http.StatusCreated, nil
}

func (s *SongStore) UpdateSong(ctx context.Context, song *models.Song) (int, error) {
	const op = "services.song_service.UpdateSong"
	log := s.log.With(
		slog.String("Operation:", op),
	)

	log.Info("updating song")
	if song.ReleaseDate != "" {

		date, err := time.Parse(time.DateOnly, song.ReleaseDate)
		if err != nil {
			log.Error("failed to update song", sl.Err(err))

			return http.StatusBadRequest, fmt.Errorf("%s: %w", op, err)
		}
		song.ReleaseDate = date.Format(time.DateOnly)
	}

	err := s.songProvider.UpdateSong(ctx, song)
	if err != nil {
		log.Error("failed to update song", sl.Err(err))

		if errors.Is(err, sterrors.ErrSongAlreadyExists) {
			return http.StatusConflict, fmt.Errorf("%s: %w", op, err)
		}

		return http.StatusInternalServerError, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("song updated")

	return http.StatusNoContent, nil
}

func (s *SongStore) GetSongs(ctx context.Context, limit, offset string, filter map[string]string) ([]*models.Song, int, error) {
	const op = "services.song_service.GetSongs"
	log := s.log.With(
		slog.String("Operation:", op),
	)

	lim, _ := strconv.Atoi(limit)
	ofs, _ := strconv.Atoi(offset)

	log.Info("getting songs")

	if filter["release_date"] != "" {
		date, err := time.Parse(time.DateOnly, filter["release_date"])
		if err != nil {
			return nil, http.StatusBadRequest, fmt.Errorf("%s: %w", op, err)
		}
		filter["release_date"] = date.Format(time.DateOnly)
	}

	songs, err := s.songGetter.GetSongs(ctx, lim, ofs, filter)
	if err != nil {
		log.Error("failed to get songs", sl.Err(err))

		return nil, http.StatusInternalServerError, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("songs gave")

	return songs, http.StatusOK, nil
}

func (s *SongStore) GetText(ctx context.Context, id int64, limit, offset string) ([]string, int, error) {
	const op = "services.song_service.GetText"
	log := s.log.With(
		slog.String("Operation:", op),
	)

	log.Info("getting song's text")

	lim, _ := strconv.Atoi(limit)
	ofs, _ := strconv.Atoi(offset)

	if ofs < 0 || ofs+lim < ofs {
		return nil, http.StatusBadRequest, nil
	}

	text, err := s.songGetter.GetText(ctx, id)
	if err != nil {
		log.Error("failed to get song's text", sl.Err(err))

		if errors.Is(err, sterrors.ErrSongNotFound) {
			return nil, http.StatusNotFound, fmt.Errorf("%s: %w", op, err)
		}

		return nil, http.StatusInternalServerError, fmt.Errorf("%s: %w", op, err)
	}

	verses := strings.Split(text, "\\n\\n")

	if lim == 0 {
		lim = len(verses)
	}

	if ofs+ofs+lim > len(verses)+1 || lim == len(verses)+1 {
		return nil, http.StatusBadRequest, nil
	}

	log.Info("song's text gave")

	return verses[ofs : ofs+lim], http.StatusOK, nil
}

func (s *SongStore) DelSong(ctx context.Context, ID int64) (int, error) {
	const op = "services.song_service.DelSong"
	log := s.log.With(
		slog.String("Operation:", op),
	)

	log.Info("deleting song")

	err := s.songProvider.DelSong(ctx, ID)
	if err != nil {
		log.Error("failed to delete song", sl.Err(err))

		return http.StatusInternalServerError, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("song deleted")

	return http.StatusNoContent, nil
}
