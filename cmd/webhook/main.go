package main

import (
	"log"
	"net/http"

	"github.com/AndySchubert/k8s-admission-webhook/internal/admission"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/validate", admission.HandleValidate)

	log.Println("Starting webhook server on :8443")

	err := http.ListenAndServe(
		":8443",
		// "/tls/tls.crt",
		// "/tls/tls.key",
		mux,
	)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
