package shorten

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"io"
	"log/slog"
	"net/http"
	"time"
	"url-shortener/internal/lib/httphelper"
	"url-shortener/internal/service/serverrors"
)

type Shortener interface {
	Shorten(ctx context.Context, sourceUrl string, expiresAt *time.Time) (string, error)
}

type Request struct {
	SourceURL string     `json:"sourceUrl" validate:"required,url"`
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`
}

type Response struct {
	ShortURL string `json:"shortUrl"`
}

// Handler сокращает URL
// @Summary Создание короткой ссылки
// @Description Принимает длинный URL и возвращает сокращенный URL
// @Tags links
// @Accept  json
// @Produce  json
// @Param request body Request true "Данные запроса"
// @Success 200 {object} Response "Успешный ответ с сокращенной ссылкой"
// @Failure 400 {string} string "Некорректный запрос"
// @Failure 422 {string} string "Невалидные параметры запроса"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /links [post]
func Handler(log *slog.Logger, shortener Shortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.handlers.shorten.Handler"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")
			http.Error(w, "empty request", http.StatusBadRequest)
			return
		}
		if err != nil {
			log.Error("failed to decode request body", slog.Any("error", err))
			http.Error(w, "failed to decode request", http.StatusBadRequest)
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request", slog.Any("error", err))
			http.Error(w, "invalid request: "+validateErr.Error(), http.StatusUnprocessableEntity)
			return
		}

		code, err := shortener.Shorten(r.Context(), req.SourceURL, req.ExpiresAt)
		if err != nil {
			if errors.Is(err, serverrors.ErrExpired) {
				//TODO: add in customValidator
				log.Error("invalid request", slog.Any("error", err))
				http.Error(w, "invalid request: "+err.Error(), http.StatusUnprocessableEntity)
			}

			log.Error("failed to shorten url", slog.Any("error", err))
			http.Error(w, "failed to shorten url", http.StatusInternalServerError)
			return
		}

		shortURL, err := httphelper.CreateURL(r, code)
		if err != nil {
			log.Error("failed to create url", slog.Any("error", err))
			http.Error(w, "failed to create url", http.StatusInternalServerError)
			return
		}

		log.Info("short url created", slog.String("shortUrl", shortURL))

		render.JSON(w, r, Response{ShortURL: shortURL})
	}
}
