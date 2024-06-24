package chi

import (
	"net/http"

	"github.com/go-chi/render"
)

func ResponseWithStatusCode(w http.ResponseWriter, r *http.Request, statusCode int, responseBody interface{}) {
	render.Status(r, statusCode)
	render.JSON(w, r, responseBody)
}
