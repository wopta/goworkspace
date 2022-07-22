package enrichVatCode

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

// GOOGLE_CLOUD_PROJECT is a user-set environment variable.
/*{
	"BaseUrl":"https://api-devexternal.munichre.com/flowin/dev/api/V1",
	"ApimKey":"59c92bc0095d4b8c803656a207150c32",
	"TokenEndPoint":"https://login.microsoftonline.com/9f2c9c2d-da50-4f33-8dfb-a780f38b50dd/oauth2/v2.0/token",
	"Scope":"46e8daaf-f894-464a-942a-e06852ed4526/.default",
	"ClientId":"194d46f8-0779-4e17-a96d-62c7bdd81901",
	"ClientSecret":"nrDLtgtLiVhvaChj1sU7JiCBUZbztRXw2ROMBYxZ",
	"GrantType":"client_credentials",
	"UWRole":"Agent",
	"SubProductId_PMIW":"35"
}*/
var projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")

func init() {
	// err is pre-declared to avoid shadowing client.
}

type publishRequest struct {
	vat string `json:"vat"`
}

// PublishMessage publishes a message to Pub/Sub. PublishMessage only works
// with topics that already exist.
func enrichVatCode(w http.ResponseWriter, r *http.Request) {
	// Parse the request body to get the topic name and message.
	p := publishRequest{}
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		log.Printf("json.NewDecoder: %v", err)
		http.Error(w, "Error parsing request", http.StatusBadRequest)
		return
	}

	var url = ""
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Ocp-Apim-Subscription-Key", "value")
	res, err := client.Do(req)
	if err != nil {
		//Handle Error
	}
	if res != nil {
	}

	if p.vat == "" {
		s := "missing 'topic' or 'message' parameter"
		log.Println(s)
		http.Error(w, s, http.StatusBadRequest)
		return
	}

	// Publish and Get use r.Context() because they are only needed for this
	// function invocation. If this were a background function, they would use
	// the ctx passed as an argument.

	fmt.Fprintf(w, "Message published: %v")
}
