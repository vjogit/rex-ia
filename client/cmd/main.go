package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"ia-client/pkg/feedback"
	"ia-client/pkg/rapport"
	"ia-client/pkg/services"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {

	configPath := "./config.yaml"
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	cfg, err := services.LoadConfigYaml(configPath)
	if err != nil {
		log.Fatal("Erreur chargement config YAML :", err)
	}

	pg, err := services.NewPG(&cfg.Database)
	if err != nil {
		log.Fatal("Erreur connexion base de données :", err)
	}

	client, err := newclient(cfg.Server)
	if err != nil {
		log.Fatal(err.Error())
	}

	urlServeur := cfg.Server.Host + ":" + strconv.Itoa(cfg.Server.Port)

	for {
		err = feedback.SendFeedback(pg, client, urlServeur)
		if err != nil {
			log.Println("Erreur SendFeedback:", err)
		}

		err = rapport.GetRapport(pg, client, urlServeur)
		if err != nil {
			log.Println("Erreur GetRapport:", err)
		}
		time.Sleep(2 * time.Second)
	}

}

func newclient(cfg services.ServerConfig) (*http.Client, error) {

	client_cert := cfg.CertPath + "/client_cert.pem" // Le certificat du client (public)
	client_key := cfg.CertPath + "/client_key.pem"   // La clé privée du client
	ca_cert := cfg.CertPath + "/ca_cert.pem"         // Le certificat de l'autorité de certification (CA) du serveur

	cert, err := tls.LoadX509KeyPair(client_cert, client_key)
	if err != nil {
		return nil, fmt.Errorf("failed to load client cert: %v", err)
	}

	ca := x509.NewCertPool()

	caBytes, err := os.ReadFile(ca_cert)
	if err != nil {
		return nil, fmt.Errorf("failed to read ca cert %q: %v", ca_cert, err)
	}
	if ok := ca.AppendCertsFromPEM(caBytes); !ok {
		return nil, fmt.Errorf("failed to parse %q", ca_cert)
	}

	tlsConfig := &tls.Config{
		ServerName:   "x.test.example.com",
		Certificates: []tls.Certificate{cert},
		RootCAs:      ca,
	}

	// 4. Créer un transport HTTP personnalisé
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
		// Définir des délais d'attente pour le transport est une bonne pratique
		// DialContext: (&net.Dialer{
		// 	Timeout: 5 * time.Second,
		// }).DialContext,
		TLSHandshakeTimeout: 5 * time.Second,
	}

	// 5. Créer le client HTTP
	client := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}

	return client, nil
}
