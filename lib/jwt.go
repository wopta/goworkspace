package lib

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-jose/go-jose/v4"
)

type JwtConfig struct {
	KeyAlgorithm       jose.KeyAlgorithm       `json:"keyAlgorithm,omitempty" firestore:"keyAlgorithm,omitempty" bigquery:"-"`
	ContentEncryption  jose.ContentEncryption  `json:"contentEncryption,omitempty" firestore:"contentEncryption,omitempty" bigquery:"-"`
	SignatureAlgorithm jose.SignatureAlgorithm `json:"signatureAlgorithm,omitempty" firestore:"signatureAlgorithm,omitempty" bigquery:"-"`
}

func ParseJwtClaims[T any](jwt, key string, jwtConfig JwtConfig, claims *T) error {
	bytes, err := ParseJwt(jwt, key, jwtConfig)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, &claims)
}

func ParseJwt(jwt, key string, jwtConfig JwtConfig) (bytes []byte, err error) {
	if jwtConfig.KeyAlgorithm != "" {
		return decryptJwt(jwt, key, jwtConfig.KeyAlgorithm, jwtConfig.ContentEncryption)
	}
	return parseSigned(jwt, key, jwtConfig.SignatureAlgorithm)
}

func decryptJwt(jwt, key string, keyAlgorithm jose.KeyAlgorithm, contentEncription jose.ContentEncryption) ([]byte, error) {
	object, err := jose.ParseEncrypted(
		jwt,
		[]jose.KeyAlgorithm{keyAlgorithm},
		[]jose.ContentEncryption{contentEncription},
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

func parseSigned(jwt, key string, algorithm jose.SignatureAlgorithm) ([]byte, error) {
	object, err := jose.ParseSigned(jwt, []jose.SignatureAlgorithm{algorithm})
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
