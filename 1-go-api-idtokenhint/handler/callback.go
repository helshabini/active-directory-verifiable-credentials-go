package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

//Struct used to temporary hold request status until UI querys for it
type CacheData struct {
	Status               string `json:"status"`
	Message              string `json:"message"`
	Payload              string `json:"payload"`
	Subject              string `json:"subject"`
	FirstName            string `json:"firstName"`
	LastName             string `json:"lastName"`
	PresentationResponse string `json:"presentationResponse"`
}

//This method is called by the VC Request API when the user scans a QR code and presents a Verifiable Credential to the service
func IssuanceRequestCallback(w http.ResponseWriter, r *http.Request) {
	var issuanceResponse map[string]interface{}
	json.NewDecoder(r.Body).Decode(&issuanceResponse)
	log.Print(issuanceResponse)
	if r.Header.Get("api-key") != apiKey {
		log.Print("api-key wrong or missing")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{'error':'api-key wrong or missing'}"))
		return
	}

	//Setting the VC response into local cache, to be retrieved later by the UI
	switch issuanceResponse["code"] {
	case "request_retrieved":
		cacheData := CacheData{
			Status:  issuanceResponse["code"].(string),
			Message: "QR Code is scanned. Waiting for issuance to complete...",
		}
		cacheJson, err := json.Marshal(cacheData)
		if err != nil {
			log.Print(err)
			return
		}
		localcache.Set(issuanceResponse["state"].(string), string(cacheJson), 0)
		return
	case "issuance_successful":
		cacheData := CacheData{
			Status:  issuanceResponse["code"].(string),
			Message: "Credential successfully issued",
		}
		cacheJson, err := json.Marshal(cacheData)
		if err != nil {
			log.Print(err)
			return
		}
		localcache.Set(issuanceResponse["state"].(string), string(cacheJson), 0)
		return
	case "issuance_error":
		errorResponse := issuanceResponse["error"].(map[string]interface{})
		cacheData := CacheData{
			Status:  issuanceResponse["code"].(string),
			Message: errorResponse["message"].(string),
		}
		cacheJson, err := json.Marshal(cacheData)
		if err != nil {
			log.Print(err)
			return
		}
		localcache.Set(issuanceResponse["state"].(string), string(cacheJson), 0)
		return
	default:
		return
	}
}

//This method is called by the VC Request API when the user scans a QR code and presents a Verifiable Credential to the service
func PresentationRequestCallback(w http.ResponseWriter, r *http.Request) {
	var presentationResponse map[string]interface{}
	json.NewDecoder(r.Body).Decode(&presentationResponse)
	log.Print(presentationResponse)
	if r.Header.Get("api-key") != apiKey {
		log.Print("api-key wrong or missing")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{'error':'api-key wrong or missing'}"))
		return
	}

	//Setting the VC response into local cache, to be retrieved later by the UI
	switch presentationResponse["code"] {
	case "request_retrieved":
		cacheData := CacheData{
			Status:  presentationResponse["code"].(string),
			Message: "QR Code is scanned. Waiting for validation...",
		}
		cacheJson, err := json.Marshal(cacheData)
		if err != nil {
			log.Print(err)
			return
		}
		localcache.Set(presentationResponse["state"].(string), string(cacheJson), 0)
		return
	case "presentation_verified":
		issuers := presentationResponse["issuers"].([]interface{})
		issuer := issuers[0].(map[string]interface{})
		claims := issuer["claims"].(map[string]interface{})
		responseBody, _ := io.ReadAll(r.Body)
		cacheData := CacheData{
			Status:               presentationResponse["code"].(string),
			Message:              "Presentation received",
			Payload:              fmt.Sprint(issuers),
			Subject:              presentationResponse["subject"].(string),
			FirstName:            claims["firstName"].(string),
			LastName:             claims["lastName"].(string),
			PresentationResponse: string(responseBody),
		}
		cacheJson, err := json.Marshal(cacheData)
		if err != nil {
			log.Print(err)
			return
		}
		localcache.Set(presentationResponse["state"].(string), string(cacheJson), 0)
		return
	default:
		return
	}
}
