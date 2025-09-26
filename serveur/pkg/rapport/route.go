package rapport

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RapportFeedback(r chi.Router, filePath string) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		GetRapport(w, r, filePath)
	})
}
