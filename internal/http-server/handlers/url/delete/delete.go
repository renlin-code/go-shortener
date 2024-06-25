package delete

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/renlin-code/go-shortener/internal/lib/api/response"
	chiResponse "github.com/renlin-code/go-shortener/internal/lib/api/response/chi"
	sl "github.com/renlin-code/go-shortener/internal/lib/logger/slog"
	"github.com/renlin-code/go-shortener/internal/storage"
)

type Response struct {
	response.Response
	Alias string `json:"alias,omitempty"`
}

type URLDeleter interface {
	DeleteURL(alias string) error
}

func NewHandler(log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.NewHandler"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")

			chiResponse.ResponseWithStatusCode(w, r, http.StatusBadRequest, response.Error("invalid request"))

			return
		}

		err := urlDeleter.DeleteURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", "alias", alias)

			chiResponse.ResponseWithStatusCode(w, r, http.StatusNotFound, response.Error("not found"))

			return
		}
		if err != nil {
			log.Error("failed to delete url", sl.Err(err))

			chiResponse.ResponseWithStatusCode(w, r, http.StatusInternalServerError, response.Error("internal error"))

			return
		}

		log.Info("deleted url", slog.String("alias", alias))

		chiResponse.ResponseWithStatusCode(w, r, http.StatusOK, Response{
			Response: response.OK(),
			Alias:    alias,
		})
	}
}
