package lib

import (
	b64 "encoding/base64"
	"fmt"
	"os"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"github.com/go-jose/go-jose/v4"
)

type JwtConfig struct {
	KeyName            string                  `json:"keyName,omitempty" firestore:"keyName,omitempty" bigquery:"-"`
	KeyAlgorithm       jose.KeyAlgorithm       `json:"keyAlgorithm,omitempty" firestore:"keyAlgorithm,omitempty" bigquery:"-"`
	ContentEncryption  jose.ContentEncryption  `json:"contentEncryption,omitempty" firestore:"contentEncryption,omitempty" bigquery:"-"`
	SignatureAlgorithm jose.SignatureAlgorithm `json:"signatureAlgorithm,omitempty" firestore:"signatureAlgorithm,omitempty" bigquery:"-"`
}

func ParseJwt(jwt string, jwtConfig JwtConfig) (bytes []byte, err error) {
	if jwtConfig.KeyAlgorithm != "" {
		return decryptJwt(jwt, os.Getenv(jwtConfig.KeyName), jwtConfig.KeyAlgorithm, jwtConfig.ContentEncryption)
	}
	return parseSigned(jwt, os.Getenv(jwtConfig.KeyName), jwtConfig.SignatureAlgorithm)
}

func decryptJwt(jwt, key string, keyAlgorithm jose.KeyAlgorithm, contentEncription jose.ContentEncryption) ([]byte, error) {
	log.AddPrefix("DecryptJwt")
	defer log.PopPrefix()
	object, err := jose.ParseEncrypted(
		jwt,
		[]jose.KeyAlgorithm{keyAlgorithm},
		[]jose.ContentEncryption{contentEncription},
	)

	if err != nil {
		log.Printf("could not parse jwt - %s", err.Error())
		return nil, fmt.Errorf("could not parse jwt")
	}

	decryptionKey, err := b64.StdEncoding.DecodeString(key)
	if err != nil {
		log.Printf("could not decode signing key - %s", err.Error())
		return nil, fmt.Errorf("could not decode jwt key")
	}

	return object.Decrypt(decryptionKey)
}

func parseSigned(jwt, key string, algorithm jose.SignatureAlgorithm) ([]byte, error) {
	log.AddPrefix("DecryptJwt")
	log.PopPrefix()
	object, err := jose.ParseSigned(jwt, []jose.SignatureAlgorithm{algorithm})
	if err != nil {
		log.Printf("could not parse jwt - %s", err.Error())
		return nil, fmt.Errorf("could not parse jwt")
	}

	decodedKey, err := b64.StdEncoding.DecodeString(key)
	if err != nil {
		log.Printf("could not decode signing key - %s", err.Error())
		return nil, fmt.Errorf("could not decode jwt key")
	}

	return object.Verify(decodedKey)
}
