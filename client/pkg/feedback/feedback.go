package feedback

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"ia-client/pkg/services"
	"log"
	"net/http"
)

func SendFeedback(pg *services.Postgres, client *http.Client, urlServeur string) error {
	// 6. Construire l'URL complète et faire la requête
	url := urlServeur + "/api/v0/feedback"
	log.Printf("Tentative de requête vers: %s", url)

	feedbacksReader, err := getFeedbackFomBD(pg)
	if err != nil {
		return fmt.Errorf("echec de la récupération des feedbacks: %v", err)
	}

	req, err := http.NewRequest("POST", url, feedbacksReader)
	if err != nil {
		return fmt.Errorf("echec de la création de la requête: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("echec de la requête TLS: %v", err)
	}
	defer resp.Body.Close()

	if _, err := services.ProcessResponse(resp); err != nil {
		return err
	}

	return nil
}

func getFeedbackFomBD(pg *services.Postgres) (*bytes.Reader, error) {

	query := New(pg.Db)

	feedbacks, err := query.ListFeedbacks(context.Background())
	if err != nil {
		return nil, err
	}
	data, err := json.Marshal(feedbacks)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(data), nil
}
