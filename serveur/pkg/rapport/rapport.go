package rapport

import (
	"net/http"
	"os"
)

func GetRapport(w http.ResponseWriter, r *http.Request, filePath string) {

	// Lire le contenu du fichier
	data, err := os.ReadFile(filePath)
	if err != nil {
		http.Error(w, "Erreur lecture fichier", http.StatusInternalServerError)
		return
	}

	// Renvoyer le contenu dans le body
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
