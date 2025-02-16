package expired

import (
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
)

const expPath = "web/views/expired.html"

func ExpiredHandler(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.handlers.expired.expired.ExpiredHandler"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		htmlContent, err := os.ReadFile(expPath)
		if err != nil {
			log.Error("couldn't load the page", slog.Any("error", err))
			http.Error(w, "couldn't load the page", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, err = w.Write(htmlContent)
		if err != nil {
			log.Error("couldn't write the page", slog.Any("error", err))
			http.Error(w, "couldn't write the page", http.StatusInternalServerError)
			return
		}

		log.Info("expiring")
	}
}
