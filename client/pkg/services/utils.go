package services

import (
	"fmt"
	"io"
	"net/http"
)

func ProcessResponse(resp *http.Response) (*string, error) {
	// 7. Traiter la réponse
	fmt.Printf("Réponse du serveur:\n")
	fmt.Printf("Statut HTTP: %s\n", resp.Status)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("echec de la lecture du corps de la réponse: %v", err)
	}

	payload := string(body)
	fmt.Printf("Corps de la réponse:\n%s\n", payload)
	return &payload, nil
}
