package postgre

import (
	"SongsLibrary/internal/models"

	"SongsLibrary/internal/storage/sterrors"
	"SongsLibrary/pkg/tools"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/jackc/pgx/v5/pgconn"
)

type StDb struct {
	db *sql.DB
}

func (s *StDb) GetDB() *sql.DB {
	return s.db
}

func New(path string) (*StDb, error) {
	db, err := sql.Open("pgx", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open db due to error: %w", err)
	}
	return &StDb{db: db}, nil
}

func (s *StDb) CreateSong(ctx context.Context, song *models.Song) error {
	const op = "storage.postgre.CreateSong"

	_, err := s.db.ExecContext(ctx,
		"INSERT INTO service.songs(group, song, release_date, text, link) VALUES ($1, $2, $3, $4, $5)",
		song.Group,
		song.Song,
		song.ReleaseDate,
		song.Text,
		song.Link,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return fmt.Errorf("%s: %w", op, sterrors.ErrSongAlreadyExists)
			}
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *StDb) UpdateSong(ctx context.Context, song *models.Song) error {
	const op = "storage.postgre.UpdateSong"

	builder := &tools.UpdateQueryBuilder{}
	builder.Update("service.songs")
	builder.Set("text", song.Text)
	builder.Set("song", song.Song)
	builder.Set("group", song.Group)
	builder.Set("link", song.Link)
	builder.Set("release_date", song.ReleaseDate)
	builder.Where("id", song.ID)

	fmt.Println(builder.String())

	res, err := s.db.ExecContext(ctx,
		builder.String(),
		builder.Args()...,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return fmt.Errorf("%s: %w", op, sterrors.ErrSongAlreadyExists)
			}
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	if count, _ := res.RowsAffected(); count == 0 {
		return fmt.Errorf("%s: %w", op, sterrors.ErrSongNotFound)
	}

	return nil
}

func (s *StDb) DelSong(ctx context.Context, ID int64) error {
	const op = "storage.postgre.DelSong"

	res, err := s.db.ExecContext(
		ctx,
		"DELETE FROM service.songs WHERE id = $1",
		ID,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if count, _ := res.RowsAffected(); count == 0 {
		return fmt.Errorf("%s: %w", op, sterrors.ErrSongNotFound)
	}

	return nil
}

func (s *StDb) GetSongs(ctx context.Context, limit, offset int, filter map[string]string) ([]*models.Song, error) {
	const op = "storage.postgre.GetSongs"

	builder := &tools.SelectQueryBuilder{}
	builder.Select("service.songs")

	for k, v := range filter {
		if v != "" {
			builder.Where(k, v)
		}
	}

	builder.Limit(limit)
	builder.Offset(offset)

	rows, err := s.db.QueryContext(ctx,
		builder.String(),
		builder.Args()...,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer rows.Close()

	var songs []*models.Song
	for rows.Next() {
		song := new(models.Song)
		err = rows.Scan(
			&song.ID,
			&song.Group,
			&song.Song,
			&song.ReleaseDate,
			&song.Link,
			&song.Text,
		)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, sterrors.ErrSongNotFound
			}

			return nil, fmt.Errorf("%s: %w", op, err)
		}

		songs = append(songs, song)
	}

	return songs, nil
}

func (s *StDb) GetText(ctx context.Context, ID int64) (string, error) {
	const op = "storage.postgre.GetText"

	row := s.db.QueryRowContext(
		ctx,
		"SELECT text FROM service.songs WHERE id = $1",
		ID,
	)

	var text string

	err := row.Scan(&text)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", sterrors.ErrSongNotFound
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return text, nil
}
