package feedback

import (
	"io"
	"log"
	"net/http"
	"os"
)

// messageHandler est la fonction qui gère les requêtes POST sur /api/V0/messages
func CreateFeedback(w http.ResponseWriter, r *http.Request, filePath string) {
	// 1. Lire le corps de la requête (le message)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		log.Printf("Error reading body: %v", err)
		return
	}
	defer r.Body.Close() // Assurer la fermeture du corps de la requête

	// 3. Stocker dans un fichier (messages.log)
	if err := storeMessageInFile(filePath, string(body)); err != nil {
		http.Error(w, "Error storing message", http.StatusInternalServerError)
		log.Printf("Error writing to file: %v", err)
		return
	}

	// 4. Répondre au client
	w.WriteHeader(http.StatusCreated) // 201 Created est approprié pour une création de ressource
	w.Write([]byte("Message received and stored successfully."))
}

// storeMessageInFile ouvre le fichier en mode APPEND et écrit la ligne.
func storeMessageInFile(filePath string, line string) error {
	// os.O_TRUNC : Remplacer le contenu précédent
	// os.O_CREATE : Créer le fichier s'il n'existe pas
	// 0644 : Permissions pour le nouveau fichier (-rw-r--r--)
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close() // Toujours fermer le fichier

	if _, err := f.Write([]byte(line)); err != nil {
		return err
	}

	return nil
}
