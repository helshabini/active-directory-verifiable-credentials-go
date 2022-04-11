package handler

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/uuid"
)

//This method is called from the UI to initiate the presentation of the verifiable credential
func PresentationRequest(w http.ResponseWriter, r *http.Request) {
	//Loading issuance request template
	payload, err := getJson(PresentationRequestFile)
	if err != nil {
		log.Fatal(err)
	}

	//Ref to callback element
	callback := payload["callback"].(map[string]interface{})

	//Setting callback url
	callbackUrl := "https://" + r.Host + "/api/verifier/presentation-request-callback"
	callback["url"] = callbackUrl

	//Setting callback state
	requestId := uuid.New()
	callback["state"] = requestId

	//Ref to headers element
	headers := callback["headers"].(map[string]interface{})

	//Setting Api-Key
	headers["api-key"] = apiKey

	//Setting Verifier Authority
	payload["authority"] = Config.VerifierAuthority

	//Ref to presentation element
	presentation := payload["presentation"].(map[string]interface{})

	//Setting accepted issuer ["presentation"]["requestedCredentials"][0]["acceptedIssuers"][0]
	requestedCredentials := presentation["requestedCredentials"].([]interface{})
	credentials := requestedCredentials[0].(map[string]interface{})
	acceptedIssuers := credentials["acceptedIssuers"].([]interface{})
	acceptedIssuers[0] = Config.IssuerAuthority

	//Generating presentation request payload
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

	//Posting presentation request
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
	body, _ = json.Marshal(respBody)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(body)
}

//This function is called from the UI polling for a response from the AAD VC Service.
//When a callback is recieved at the presentationCallback service the session will be updated
//This method will respond with the status so the UI can reflect if the QR code was scanned and with the result of the presentation
func PresentationResponse(w http.ResponseWriter, r *http.Request) {
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
