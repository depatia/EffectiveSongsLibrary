package app

import (
	_ "SongsLibrary/docs"
	"SongsLibrary/internal/config"
	"SongsLibrary/internal/delivery/http/song"
	externalservice "SongsLibrary/internal/services/external_service"
	songservice "SongsLibrary/internal/services/song_service"
	"SongsLibrary/internal/storage/postgre"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/mattes/migrate/source/file"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Server struct {
	server *http.Server
	Router *httprouter.Router
	log    *slog.Logger
	cfg    *config.Config
}

func New(
	ctx context.Context,
	log *slog.Logger,
	cfg *config.Config,
) *Server {
	router := httprouter.New()

	storage, err := postgre.New(cfg.DBPath)
	if err != nil {
		panic(err)
	}
	db := storage.GetDB()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Error(err.Error())
	}

	m, err := migrate.NewWithDatabaseInstance(
		cfg.MigrationsPath,
		"postgres", driver)
	if err != nil {
		log.Error(err.Error())
	}

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Error(err.Error())
	}

	router.HandlerFunc(http.MethodGet, "/swagger/:any", httpSwagger.WrapHandler)

	externalService := externalservice.New(log, cfg)
	songService := songservice.New(log, storage, storage, storage, externalService)

	songHandler := song.Handler{Service: songService}
	songHandler.Register(router)

	return &Server{Router: router, cfg: cfg, log: log}

}

func (s *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.cfg.IP, s.cfg.Port))
	if err != nil {
		s.log.Error(err.Error())
	}

	s.server = &http.Server{
		Handler:      s.Router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	s.log.Info("application started")

	if err := s.server.Serve(listener); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			s.log.Info("server shutdown")
		default:
			s.log.Error(err.Error())
		}
	}
}

func (s *Server) Stop() {
	const op = "app.Stop"

	s.log.With(slog.String("op", op)).
		Info("stopping http server", slog.Int("port", s.cfg.Port))

	s.server.Close()
}
