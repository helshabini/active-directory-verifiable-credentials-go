package auth

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/Azure-Samples/active-directory-verifiable-credentials-go/1-go-api-idtokenhint/config"
	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
)

var (
	cacheAccessor = &TokenCache{"serialized_cache.json"}
)

func AcquireTokenClientSecret(config config.Config) (string, error) {
	cred, err := confidential.NewCredFromSecret(config.AzClientSecret)
	if err != nil {
		log.Fatal(err)
	}
	app, err := confidential.New(config.AzClientId, cred, confidential.WithAuthority("https://login.microsoftonline.com/"+config.AzTenantId), confidential.WithAccessor(cacheAccessor))
	if err != nil {
		log.Fatal(err)
	}
	result, err := app.AcquireTokenSilent(context.Background(), []string{config.VerifiableCredentialsDefaultScope})
	if err != nil {
		result, err = app.AcquireTokenByCredential(context.Background(), []string{config.VerifiableCredentialsDefaultScope})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Access Token Is " + result.AccessToken)
		return result.AccessToken, nil
	}
	fmt.Println("Silently acquired token " + result.AccessToken)
	return result.AccessToken, nil
}

func AcquireTokenClientCertificate(config config.Config, password string) (string, error) {
	pemData, err := ioutil.ReadFile(config.AzPEMFileLocation)
	if err != nil {
		log.Fatal(err)
	}

	// This extracts our public certificates and private key from the PEM file.
	// The private key must be in PKCS8 format. If it is encrypted, the second argument
	// must be password to decode.
	certs, privateKey, err := confidential.CertFromPEM(pemData, password)
	if err != nil {
		log.Fatal(err)
	}

	// PEM files can have multiple certs. This is usually for certificate chaining where roots
	// sign to leafs. Useful for TLS, not for this use case.
	if len(certs) > 1 {
		log.Fatal("too many certificates in PEM file")
	}

	cred := confidential.NewCredFromCert(certs[0], privateKey)
	if err != nil {
		log.Fatal(err)
	}
	app, err := confidential.New(config.AzClientId, cred, confidential.WithAuthority("https://login.microsoftonline.com/"+config.AzTenantId), confidential.WithAccessor(cacheAccessor))
	if err != nil {
		log.Fatal(err)
	}
	result, err := app.AcquireTokenSilent(context.Background(), []string{config.VerifiableCredentialsDefaultScope})
	if err != nil {
		result, err = app.AcquireTokenByCredential(context.Background(), []string{config.VerifiableCredentialsDefaultScope})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Access Token is " + result.AccessToken)
		return result.AccessToken, nil
	}
	fmt.Println("Silently acquired token " + result.AccessToken)
	return result.AccessToken, nil
}
