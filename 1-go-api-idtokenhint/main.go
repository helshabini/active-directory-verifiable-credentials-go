package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Azure-Samples/active-directory-verifiable-credentials-go/1-go-api-idtokenhint/config"
	"github.com/Azure-Samples/active-directory-verifiable-credentials-go/1-go-api-idtokenhint/handler"
)

func main() {
	args := os.Args[1:]
	if len(args) > 0 && len(args) < 3 {
		log.Fatal("Invalid arguments.")
	}

	//Loading configuration from config.json
	configfile := "config.json"
	if args[0] != "" {
		configfile = args[0]
	}
	config, err := config.LoadConfig(configfile)
	if err != nil {
		log.Fatal(err)
	}

	//Injecting configuration into Handler package
	handler.Config = config
	if args[1] != "" {
		handler.IssuanceRequestFile = args[1]
	}
	if args[2] != "" {
		handler.PresentationRequestFile = args[2]
	}
	if len(args) == 4 {
		handler.PEMFilePassword = args[3]
	}

	//Setting up routes
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/api/issuer/issuance-request", handler.IssuanceRequest)
	http.HandleFunc("/api/issuer/issuance-request-callback", handler.IssuanceRequestCallback)
	http.HandleFunc("/api/issuer/issuance-response", handler.IssuanceResponse)
	http.HandleFunc("/api/verifier/presentation-request", handler.PresentationRequest)
	http.HandleFunc("/api/verifier/presentation-request-callback", handler.PresentationRequestCallback)
	http.HandleFunc("/api/verifier/presentation-response", handler.PresentationResponse)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
