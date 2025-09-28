package main

import (
	"crypto/tls"
	"crypto/x509"
	"ia-client/pkg/feedback"
	"ia-client/pkg/rapport"
	"ia-client/pkg/services"
	"log"
	"net"
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

	clientCertPath := cfg.CertPath + "/client.crt" // Le certificat du client (public)
	clientKeyPath := cfg.CertPath + "/client.key"  // La clé privée du client
	caCertPath := cfg.CertPath + "/ca.crt"         // Le certificat de l'autorité de certification (CA) du serveur

	// 1. Charger la clé et le certificat du client
	cert, err := tls.LoadX509KeyPair(clientCertPath, clientKeyPath)
	if err != nil {
		log.Fatalf("❌ Échec du chargement du certificat client: %v", err)
	}
	log.Println("✅ Certificat client chargé avec succès.")

	// 2. Préparer le pool de certificats CA du serveur (RootCAs)
	ca := x509.NewCertPool()
	caBytes, err := os.ReadFile(caCertPath)
	if err != nil {
		log.Fatalf("❌ Échec de la lecture du certificat CA serveur %q: %v", caCertPath, err)
	}
	if ok := ca.AppendCertsFromPEM(caBytes); !ok {
		log.Fatalf("❌ Échec du parsing du certificat CA serveur %q", caCertPath)
	}
	log.Println("✅ CA du serveur chargée pour la vérification.")

	// 3. Créer la configuration TLS
	tlsConfig := &tls.Config{
		ServerName:   "localhost",             // Doit correspondre au CN/SAN du 'server_cert.pem' du serveur de test
		Certificates: []tls.Certificate{cert}, // Certificat du client pour le mTLS
		RootCAs:      ca,                      // CA pour vérifier le certificat du serveur
		MinVersion:   tls.VersionTLS12,
	}

	// 4. Créer un transport HTTP personnalisé avec le DialContext corrigé
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
		// Utilisation de net.Dialer pour la correction de l'erreur "http.Dialer est inconnu"
		DialContext: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout: 5 * time.Second,
	}

	// 5. Créer le client HTTP
	client := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}

	return client, nil
}
