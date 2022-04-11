package handler

import (
	"crypto/rand"
	"encoding/json"
	"io"
	"io/ioutil"
	"time"

	"github.com/Azure-Samples/active-directory-verifiable-credentials-go/1-go-api-idtokenhint/auth"
	"github.com/Azure-Samples/active-directory-verifiable-credentials-go/1-go-api-idtokenhint/config"
	"github.com/google/uuid"

	"github.com/patrickmn/go-cache"
)

var (
	Config                  = config.Config{}
	IssuanceRequestFile     = "issuance_request_config.json"
	PresentationRequestFile = "presentation_request_config.json"
	PEMFilePassword         = ""
	apiKey                  = uuid.New().String()
	localcache              = cache.New(5*time.Minute, 10*time.Minute)
)

func getJson(issuancefile string) (map[string]interface{}, error) {
	bytes, err := ioutil.ReadFile(issuancefile)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	json.Unmarshal(bytes, &result)
	return result, nil
}

func getAccessToken() (string, error) {
	//Issuing access token using Client Certificate if it exists
	if Config.AzPEMFileLocation != "" {
		//Second parameter should be the password for PEM file
		accessToken, err := auth.AcquireTokenClientCertificate(Config, PEMFilePassword)
		if err != nil {
			return "", err
		}
		return accessToken, err
	}

	//Issuing access token using Client Secret
	accessToken, err := auth.AcquireTokenClientSecret(Config)
	if err != nil {
		return "", err
	}
	return accessToken, err
}

func generatePin(max int) string {
	table := [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}
