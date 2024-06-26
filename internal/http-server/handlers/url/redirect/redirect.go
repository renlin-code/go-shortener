package redirect

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

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=URLGetter
type URLGetter interface {
	GetURL(alias string) (string, error)
}

func NewHandler(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.redirect.NewHandler"

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

		resURL, err := urlGetter.GetURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", "alias", alias)

			chiResponse.ResponseWithStatusCode(w, r, http.StatusNotFound, response.Error("not found"))

			return
		}
		if err != nil {
			log.Error("failed to get url", sl.Err(err))

			chiResponse.ResponseWithStatusCode(w, r, http.StatusInternalServerError, response.Error("internal error"))

			return
		}

		log.Info("got url", slog.String("url", resURL))

		http.Redirect(w, r, resURL, http.StatusFound)
	}
}
