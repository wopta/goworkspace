package lib

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"log"

	"gopkg.in/square/go-jose.v2"
)

func DecryptJwt[T interface{}](jwtData, key string, claims *T) error {
	object, err := jose.ParseEncrypted(jwtData)
	if err != nil {
		log.Printf("[DecryptJwt] could not parse jwt - %s", err.Error())
		return fmt.Errorf("could not parse jwt")
	}

	decryptionKey, err := b64.StdEncoding.DecodeString(key)
	if err != nil {
		log.Printf("[DecryptJwt] could not decode signing key - %s", err.Error())
		return fmt.Errorf("could not decode jwt key")
	}
	decrypted, err := object.Decrypt(decryptionKey)
	if err != nil {
		log.Printf("[DecryptJwt] could not decrypt jwt - %s", err.Error())
		return fmt.Errorf("could not decrypt jwt")
	}

	return json.Unmarshal(decrypted, &claims)
}
