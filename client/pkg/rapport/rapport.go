package rapport

import (
	"fmt"
	"ia-client/pkg/services"
	"log"
	"net/http"
)

func GetRapport(pg *services.Postgres, client *http.Client, urlServeur string) error {
	// 6. Construire l'URL complète et faire la requête
	url := urlServeur + "/api/v0/rapport"
	log.Printf("Tentative de requête vers: %s", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("echec de la création de la requête: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("echec de la requête TLS: %v", err)
	}
	defer resp.Body.Close()

	body, err := services.ProcessResponse(resp)
	if err != nil {
		return err
	}
	log.Printf("Rapport reçu: %s", *body)

	return nil
}
