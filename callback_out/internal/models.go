package internal

import "net/http"

type CallbackInfo struct {
	Request     *http.Request
	RequestBody []byte
	Response    *http.Response
	Error       error
}
