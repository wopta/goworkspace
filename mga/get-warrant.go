package mga

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

type GetWarrantsResponse struct {
	Warrants []models.Warrant `json:"warrants"`
}

func GetWarrantsFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var response GetWarrantsResponse

	log.SetPrefix("[GetWarrantsFx] ")

	log.Println("Handler start -----------------------------------------------")

	warrants, err := GetWarrants()
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

	log.Printf("found warrants: %s", string(responseBytes))

	log.Println("Handler end -------------------------------------------------")

	return string(responseBytes), response, nil
}

func GetWarrants() ([]models.Warrant, error) {
	var (
		err      error
		warrants []models.Warrant
	)

	warrantsBytes := lib.GetFolderContentByEnv(models.WarrantsFolder)

	for _, warrantBytes := range warrantsBytes {
		var warrant models.Warrant
		err = json.Unmarshal(warrantBytes, &warrant)
		if err != nil {
			log.Printf("[GetWarrants] error unmarshaling warrant: %s", err.Error())
			return warrants, err
		}

		warrants = append(warrants, warrant)
	}
	return warrants, nil
}
