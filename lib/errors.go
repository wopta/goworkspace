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
