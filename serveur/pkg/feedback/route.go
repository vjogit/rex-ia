package feedback

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RouteFeedback(r chi.Router, filePath string) {
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		CreateFeedback(w, r, filePath)
	})
}
