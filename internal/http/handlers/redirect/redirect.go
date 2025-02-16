package redirect

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"url-shortener/internal/lib/httphelper"
	"url-shortener/internal/service/serverrors"
)

type Shortener interface {
	Resolve(ctx context.Context, code string) (string, error)
}

func Handler(log *slog.Logger, shortener Shortener, expiredURLPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.handlers.redirect.redirect.Handler"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		code := chi.URLParam(r, "code")
		if code == "" {
			log.Info("code is empty")
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		sourceURL, err := shortener.Resolve(r.Context(), code)
		if err != nil {
			handleResolveError(w, r, err, log, expiredURLPath)
			return
		}

		log.Info("redirecting", slog.String("url", sourceURL))
		http.Redirect(w, r, sourceURL, http.StatusFound)
	}
}

func handleResolveError(w http.ResponseWriter, r *http.Request, err error, log *slog.Logger, expiredURLPath string) {
	switch {
	case errors.Is(err, serverrors.ErrURLNotFound):
		log.Warn("url not found", slog.Any("error", err))
		http.Error(w, "url not found", http.StatusNotFound)

	case errors.Is(err, serverrors.ErrURLExpired):
		expiredURL, err := httphelper.CreateURL(r, expiredURLPath)
		if err != nil {
			log.Error("failed to create expired url", slog.Any("error", err))
			http.Error(w, "failed to create expired url", http.StatusInternalServerError)
			return
		}
		log.Warn("redirecting to expired page", slog.String("url", expiredURL))
		http.Redirect(w, r, expiredURL, http.StatusTemporaryRedirect)

	default:
		log.Error("failed to resolve code", slog.Any("error", err))
		http.Error(w, "failed to resolve code", http.StatusInternalServerError)
	}
}
