package mga

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
)

type GetWarrantsResponse struct {
	Warrants []models.Warrant `json:"warrants"`
}

func GetWarrantsFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var response GetWarrantsResponse

	log.SetPrefix("[GetWarrantsFx] ")
	defer log.SetPrefix("")

	log.Println("Handler start -----------------------------------------------")

	warrants, err := network.GetWarrants()
	if err != nil {
		log.Printf("error getting warrants: %s", err.Error())
		return "", "", err
	}

	response.Warrants = warrants

	responseBytes, err := json.Marshal(response)
	if err != nil {
		log.Printf("error marshaling response: %s", err.Error())
		return "", "", err
	}

	log.Println("Handler end -------------------------------------------------")

	return string(responseBytes), response, nil
}
