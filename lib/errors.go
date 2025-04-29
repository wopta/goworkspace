package lib

import (
	"encoding/json"
	"errors"
	"github.com/wopta/goworkspace/lib/log"
	"net/http"
)

func CheckError(e error) {
	if e != nil {
		log.Error(e)
		panic(e)

	}
}

func ErrorByte(b []byte, e error) []byte {
	CheckError(e)
	return b
}

func CheckErrorResp(w http.ResponseWriter, e error) {

	if e != nil {
		log.Error(e)
		http.Error(w, e.Error(), 500)
		panic(e)

	}
}

type ErrorResponse struct {
	Code    int    `firestore:"-" json:"code,omitempty" bigquery:"name"`          //h-Nome
	Type    string `firestore:"-" json:"type,omitempty" bigquery:"surname"`       //Cognome
	Message string `firestore:"-" json:"message,omitempty" bigquery:"fiscalCode"` //Codice fiscale

}

func GetErrorJson(code int, typeE string, message string) error {
	var (
		e     error
		eResp ErrorResponse
		b     []byte
	)
	eResp = ErrorResponse{Code: code, Type: typeE, Message: message}
	b, e = json.Marshal(eResp)
	e = errors.New(string(b))
	return e
}
