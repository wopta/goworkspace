package mga

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

type CreateWarrantResponse struct {
	Success bool `json:"success"`
}

func CreateWarrantFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println("[CreateWarrantFx] Handler start ------------")

	var (
		response CreateWarrantResponse
		warrant  models.Warrant
	)

	bodyBytes, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Printf("[CreateWarrantFx] error reading request body: %s", err.Error())
		return "", "", err
	}

	err = json.Unmarshal(bodyBytes, &warrant)
	if err != nil {
		log.Printf("[CreateWarrantFx] error marshaling request: %s", err.Error())
		return "", "", err
	}

	err = CreateWarrant(warrant)
	if err != nil {
		log.Printf("[CreateWarrantFx] error creating warrant: %s", err.Error())
		return "", response, err
	}

	response.Success = err == nil
	if err == nil {
		log.Printf("[CreateWarrantFx] created warrant %s", warrant.Name)
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		log.Printf("[CreateWarrantFx] error marshaling response: %s", err.Error())
		return "", response, err
	}

	models.CreateAuditLog(r, string(bodyBytes))

	return string(responseBytes), response, err
}

func CreateWarrant(warrant models.Warrant) error {
	fileName := models.WarrantsFolder + warrant.Name + ".json"

	bytesToWrite, err := json.Marshal(warrant)
	if err != nil {
		log.Printf("[CreateWarrant] error marshaling warrant: %s", err.Error())
		return err
	}

	_, err = lib.PutToStorageIfNotExists(os.Getenv("GOOGLE_STORAGE_BUCKET"), fileName, bytesToWrite)

	if err != nil {
		log.Printf("[CreateWarrant] error writing warrant: %s", err.Error())
	}

	return err
}
