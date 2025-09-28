package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"ia-serveur/pkg/feedback"
	"ia-serveur/pkg/rapport"
	"ia-serveur/pkg/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Ces variables seront injectées au moment de la compilation
var (
	buildTime string
	version   string
)

const path = "/home/vjo/projets/rex/rex-ia/testdata/x509/"

const (
	server_cert    = path + "server_cert.pem"
	server_key     = path + "server_key.pem"     // La clé privée du client
	client_ca_cert = path + "client_ca_cert.pem" // Le certificat de l'autorité de certification (CA) du client
)

func main() {

	// Affiche les informations de compilation
	log.Printf("Application version: %s", version)
	log.Printf("Compilation time: %s", buildTime)

	r := chi.NewRouter()
	r.Use(middleware.Logger) // Log HTTP requests

	configPath := "./config.yaml"
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	cfg, err := services.LoadConfigYaml(configPath)
	if err != nil {
		log.Fatal("Erreur chargement config YAML :", err)
	}

	// version api1
	r.Route("/api/v0", func(r chi.Router) {
		r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
			log.Println("Ping reçu api 1!")
			w.Write([]byte("pong"))
		})

		r.Route("/feedback", func(r chi.Router) {
			feedback.RouteFeedback(r, cfg.Location.Feedback)
		})

		r.Route("/rapport", func(r chi.Router) {
			rapport.RapportFeedback(r, cfg.Location.Rapport)
		})
	})

	log.Printf("Serveur démarré sur le port %d (HTTP)", cfg.Server.Port)
	log.Fatal(http.ListenAndServe(
		fmt.Sprintf(":%d", cfg.Server.Port),
		r,
	))
}
