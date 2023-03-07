package broker

import (
	"net/http"
)

func PolicyFiscalcode(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {

	r.Header.Get("fiscalcode")

	return "", nil, nil
}
