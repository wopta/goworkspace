package lib

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-jose/go-jose/v4"
	oldJose "gopkg.in/square/go-jose.v2"
)

func DecryptJwt[T interface{}](jwtData, key string, claims *T) error {
	object, err := oldJose.ParseEncrypted(jwtData)
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

func ParseJwtClaims[T any](jwt, key string, isEncrypted bool, claims *T) error {
	bytes, err := ParseJwt(jwt, key, isEncrypted)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, &claims)
}

func ParseJwt(jwt, key string, isEncrypted bool) (bytes []byte, err error) {
	if isEncrypted {
		return decryptJwt(jwt, key)
	}
	return parseSigned(jwt, key)
}

func decryptJwt(jwt, key string) ([]byte, error) {
	object, err := jose.ParseEncrypted(
		jwt,
		[]jose.KeyAlgorithm{jose.DIRECT},
		[]jose.ContentEncryption{jose.A128CBC_HS256},
	)
	if err != nil {
		log.Printf("[decryptJwt] could not parse jwt - %s", err.Error())
		return nil, fmt.Errorf("could not parse jwt")
	}

	decryptionKey, err := b64.StdEncoding.DecodeString(key)
	if err != nil {
		log.Printf("[decryptJwt] could not decode signing key - %s", err.Error())
		return nil, fmt.Errorf("could not decode jwt key")
	}

	return object.Decrypt(decryptionKey)
}

func parseSigned(jwt, key string) ([]byte, error) {
	object, err := jose.ParseSigned(jwt, []jose.SignatureAlgorithm{jose.HS256})
	if err != nil {
		log.Printf("[DecryptJwt] could not parse jwt - %s", err.Error())
		return nil, fmt.Errorf("could not parse jwt")
	}

	decodedKey, err := b64.StdEncoding.DecodeString(key)
	if err != nil {
		log.Printf("[decryptJwt] could not decode signing key - %s", err.Error())
		return nil, fmt.Errorf("could not decode jwt key")
	}

	return object.Verify(decodedKey)
}
