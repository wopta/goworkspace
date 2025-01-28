package sellable

import (
	"fmt"
	"net/http"
)

func CommercialCombinedFx(_ http.ResponseWriter, r *http.Request) (string, any, error) {
	return "", nil, fmt.Errorf("policy not sellable by: %s", "random reason")
}