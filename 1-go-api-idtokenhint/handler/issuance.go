package handler

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func IssuanceRequest(w http.ResponseWriter, r *http.Request) {
	//Loading issuance request template
	payload, err := getJson(IssuanceRequestFile)
	if err != nil {
		log.Fatal(err)
	}

	//Ref to callback element
	callback := payload["callback"].(map[string]interface{})

	//Setting callback url
	callbackUrl := "https://" + r.Host + "/api/issuer/issuance-request-callback"
	callback["url"] = callbackUrl

	//Setting callback state
	requestId := uuid.New()
	callback["state"] = requestId

	//Setting callback apiKey
	headers := callback["headers"].(map[string]interface{})
	headers["api-key"] = apiKey

	//Ref to issuance element
	issuance := payload["issuance"].(map[string]interface{})

	//Setting manifest
	issuance["manifest"] = Config.CredentialManifest

	//Setting Custom Claims (optional)
	claims := issuance["claims"].(map[string]interface{})
	if claims != nil {
		claims["given_name"] = "Megan"
		claims["family_name"] = "Bowen"
	}

	//Setting PIN code (optional except if setting custom claims)
	pin := issuance["pin"].(map[string]interface{})
	if pin != nil {
		pinLength := int(pin["length"].(float64))
		if pinLength != 0 && claims != nil {
			pin["value"] = generatePin(pinLength)
		} else {
			issuance["pin"] = nil
		}
	}

	//Setting issuance authority
	payload["authority"] = Config.IssuerAuthority

	//Generating issuance request payload
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(string(payloadJson))

	//Generating access token
	accessToken, err := getAccessToken()
	if err != nil {
		log.Fatal(err)
	}

	//Posting issuance request
	req, _ := http.NewRequest("POST", Config.MsIdentityHostName+Config.AzTenantId+"/verifiablecredentials/request", bytes.NewBuffer(payloadJson))
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Print("response Body:", string(body))

	var respBody map[string]interface{}
	json.Unmarshal(body, &respBody)
	respBody["id"] = requestId
	if pin != nil {
		respBody["pin"] = pin["value"]
	}
	body, _ = json.Marshal(respBody)
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

//This function is called from the UI polling for a response from the AAD VC Service.
//When a callback is recieved at the presentationCallback service the session will be updated
func IssuanceResponse(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("id")
	cachedata, found := localcache.Get(code)
	if !found {
		log.Print("id not found")
		w.Write([]byte("id not found"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(cachedata.(string)))
}
