package main

import (
	"log"
	"net/http"

	"github.com/AndySchubert/k8s-admission-webhook/internal/admission"
)

func main() {
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
		"/tls/tls.crt",
		"/tls/tls.key",
		mux,
	)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
