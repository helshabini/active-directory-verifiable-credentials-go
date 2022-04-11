package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	AzTenantId                        string `json:"azTenantId"`
	AzClientId                        string `json:"azClientId"`
	AzClientSecret                    string `json:"azClientSecret"`
	AzPEMFileLocation                 string `json:"azPEMFileLocation"`
	VerifiableCredentialsDefaultScope string `json:"VerifiableCredentialsDefaultScope"`
	MsIdentityHostName                string `json:"msIdentityHostName"`
	CredentialManifest                string `json:"CredentialManifest"`
	IssuerAuthority                   string `json:"IssuerAuthority"`
	VerifierAuthority                 string `json:"VerifierAuthority"`
}

func LoadConfig(filename string) (Config, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
