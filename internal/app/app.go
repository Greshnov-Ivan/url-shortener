package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/speps/go-hashids/v2"
	httpSwagger "github.com/swaggo/http-swagger"
	"log/slog"
	"net/http"
	"url-shortener/internal/config"
	"url-shortener/internal/http/handlers/error/expired"
	"url-shortener/internal/http/handlers/links/shorten"
	"url-shortener/internal/http/handlers/redirect"
	mwCors "url-shortener/internal/http/middleware/cors"
	mwLogger "url-shortener/internal/http/middleware/logger"
	"url-shortener/internal/repository/postgres"
	"url-shortener/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
	_ "url-shortener/docs"
)

type Repository interface {
	CloseDB() error
}

type App struct {
	server *http.Server
	log    *slog.Logger
	rep    Repository
}

func NewApp(cfg config.Config, log *slog.Logger) (*App, error) {
	// open connect to db
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.ConnectionStrings.User,
		cfg.ConnectionStrings.Password,
		cfg.ConnectionStrings.HOST,
		cfg.ConnectionStrings.PORT,
		cfg.ConnectionStrings.DB,
		cfg.ConnectionStrings.SSLMode,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Error("failed to connect to database", slog.Any("error", err))
		return nil, err
	}
	if err = db.Ping(); err != nil {
		log.Error("failed to ping database", slog.Any("error", err))
		return nil, err
	}

	// init repository: postgres
	rep := postgres.NewRepository(db)

	// init hashIds
	hd := hashids.NewData()
	hd.Salt = cfg.HashIDConfiguration.Salt
	hd.MinLength = cfg.HashIDConfiguration.MinHashLength
	hd.Alphabet = cfg.HashIDConfiguration.Alphabet
	h, err := hashids.NewWithData(hd)
	if err != nil {
		log.Error("failed to create new hashID", slog.Any("error", err))
		return nil, err
	}
	// init service url shortener
	urlShortenerService := service.NewUrlShortenerService(log, h, rep)

	// configuring the router
	router := chi.NewRouter()
	router.Use(mwCors.CORSMiddleware)
	router.Use(middleware.RequestID)
	// router.Use(middleware.Logger)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Get("/swagger/*", httpSwagger.WrapHandler)
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("OK")); err != nil {
			log.Error("failed to write health response", slog.Any("error", err))
		}
	})
	router.Get(cfg.ExpiredURLPath, expired.ExpiredHandler(log))
	router.Post("/links", shorten.Handler(log, urlShortenerService))
	router.Get("/{code}", redirect.Handler(log, urlShortenerService, cfg.ExpiredURLPath))

	// init http-server
	server := &http.Server{
		Addr:              cfg.Address,
		Handler:           router,
		ReadTimeout:       cfg.HTTPServer.Timeout,
		WriteTimeout:      cfg.HTTPServer.Timeout,
		IdleTimeout:       cfg.HTTPServer.IdleTimeout,
		ReadHeaderTimeout: cfg.HTTPServer.ReadHeaderTimeout,
	}

	return &App{
		server: server,
		log:    log,
		rep:    rep,
	}, nil
}

// run server
func (a *App) Run() error {
	a.log.Info("server running", slog.String("address", a.server.Addr))

	err := a.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		a.log.Error("server crashed", slog.Any("error", err))
	}
	return err
}

// stopping the server and freeing up resources
func (a *App) Shutdown(ctx context.Context) error {
	a.log.Info("closing repository...")
	if err := a.rep.CloseDB(); err != nil {
		a.log.Error("failed to close repository", slog.Any("error", err))
		return err
	}

	a.log.Info("shutting down HTTP server...")
	return a.server.Shutdown(ctx)
}
