package test

import (
	"encoding/json"
	"errors"
)

func GetErrorJson(code int, typeEr string, message string) error {
	var (
		e     error
		eResp map[string]interface{} = make(map[string]interface{})
		b     []byte
	)
	eResp["code"] = code
	eResp["type"] = typeEr
	eResp["message"] = message
	b, e = json.Marshal(eResp)
	e = errors.New(string(b))
	return e
}
