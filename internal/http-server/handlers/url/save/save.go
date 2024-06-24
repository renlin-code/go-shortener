package save

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/renlin-code/go-shortener/internal/lib/api/response"
	chi "github.com/renlin-code/go-shortener/internal/lib/api/response/chi"
	sl "github.com/renlin-code/go-shortener/internal/lib/logger/slog"
	"github.com/renlin-code/go-shortener/internal/lib/random"
	"github.com/renlin-code/go-shortener/internal/storage"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	response.Response
	Alias string `json:"alias,omitempty"`
}

const aliasLength = 4

type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
}

func NewHandler(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.NewHandler"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)

		if err != nil {
			log.Error("failed to decode request", sl.Err(err))

			chi.ResponseWithStatusCode(w, r, http.StatusBadRequest, response.Error("failed to decode request"))

			return
		}

		log.Info("request", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			chi.ResponseWithStatusCode(w, r, http.StatusBadRequest, response.ValidationError(validateErr))

			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}

		id, err := urlSaver.SaveURL(req.URL, alias)

		if errors.Is(err, storage.ErrURLAlreadyExists) {
			log.Info("url already exists", slog.String("url", req.URL))

			chi.ResponseWithStatusCode(w, r, http.StatusConflict, response.Error("url already exists"))

			return
		}

		if err != nil {
			log.Error("failed to save url", sl.Err(err))

			chi.ResponseWithStatusCode(w, r, http.StatusInternalServerError, response.Error("failed to save url"))

			return
		}

		log.Info("url saved", slog.Int64("id", id))

		chi.ResponseWithStatusCode(w, r, http.StatusCreated, Response{
			Response: response.OK(),
			Alias:    alias,
		})
	}
}
