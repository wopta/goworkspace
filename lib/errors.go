package lib

import (
	"log"
	"net/http"
)

func CheckError(e error) {
	if e != nil {
		log.Fatal(e)
		panic(e)

	}
}

func ErrorByte(b []byte, e error) []byte {
	CheckError(e)
	return b
}

func CheckErrorResp(w http.ResponseWriter, e error) {

	if e != nil {
		log.Fatal(e)
		http.Error(w, e.Error(), 500)
		panic(e)

	}
}
type ErrorResponse struct {
	Code    int    `firestore:"-" json:"code,omitempty" bigquery:"name"`          //h-Nome
	Type    string `firestore:"-" json:"type,omitempty" bigquery:"surname"`       //Cognome
	Message string `firestore:"-" json:"message,omitempty" bigquery:"fiscalCode"` //Codice fiscale

}
func GetErrorJson(code int , type string, message string)error{
	var (
		e     error
		eResp ErrorResponse
		b     []byte
	)
	eResp = ErrorResponse{Code: code, Type: type, Message: message}
	b, e = json.Marshal(eResp)
	e = errors.New(string(b))
	return e
}