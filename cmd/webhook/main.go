package main

import (
	"log"
	"net/http"
	"os"

	"github.com/AndySchubert/k8s-admission-webhook/internal/admission"
)

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	certFile := envOrDefault("TLS_CERT", "/tls/tls.crt")
	keyFile := envOrDefault("TLS_KEY", "/tls/tls.key")

	mux := http.NewServeMux()

	mux.HandleFunc("/validate", admission.HandleValidate)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	log.Println("Starting webhook server on :8443")

	err := http.ListenAndServeTLS(
		":8443",
		certFile,
		keyFile,
		mux,
	)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
