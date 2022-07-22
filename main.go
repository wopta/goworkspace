package main

import (
	"fmt"
	"log"
	"net/http"
	enr "github.com/wopta/goworkspace/enrichVatCode"
)

func main() {

	http.HandleFunc("/", enr.enrichVatCode)
	fmt.Println("Listening on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
