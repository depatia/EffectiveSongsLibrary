package externalservice

import (
	"SongsLibrary/internal/config"
	"SongsLibrary/internal/models"
	"SongsLibrary/pkg/logger/sl"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
)

type ExternalStore struct {
	log *slog.Logger
	cfg *config.Config
}

func New(log *slog.Logger, cfg *config.Config) *ExternalStore {
	return &ExternalStore{log: log, cfg: cfg}
}

func (s *ExternalStore) GetSongDetail(ctx context.Context, group, songName string) (*models.Song, int, error) {
	const op = "services.external_service.GetSongDetail"
	log := s.log.With(
		slog.String("Operation:", op),
	)

	log.Info("getting song info")

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/info", s.cfg.ExternalAPIUrl), nil)
	if err != nil {
		log.Error("failed to get song detail from external API ", sl.Err(err))

		return nil, http.StatusInternalServerError, err
	}
	q := url.Values{}

	q.Add("group", group)
	q.Add("song", songName)

	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error("failed to get song detail from external API. status code: %d", sl.Err(err), resp.StatusCode)

		return nil, resp.StatusCode, err
	}

	defer resp.Body.Close()

	song := &models.Song{}

	if err := json.NewDecoder(resp.Body).Decode(song); err != nil {
		log.Error("failed to get song detail from external API", sl.Err(err))

		return nil, http.StatusInternalServerError, err
	}

	song.Group = group
	song.Song = songName

	log.Info("song info gave")

	return song, resp.StatusCode, nil
}
