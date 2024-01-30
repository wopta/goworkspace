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

func CreateWarrantFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var warrant models.Warrant

	log.SetPrefix("[CreateWarrantFx] ")
	defer log.SetPrefix("")

	log.Println("Handler start -----------------------------------------------")

	bodyBytes, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Printf("error reading request body: %s", err.Error())
		return "", nil, err
	}
	log.Printf("request: %s", string(bodyBytes))

	err = json.Unmarshal(bodyBytes, &warrant)
	if err != nil {
		log.Printf("error unmarshaling request: %s", err.Error())
		return "", nil, err
	}

	log.Println("creating warrant...")
	err = CreateWarrant(warrant)
	if err != nil {
		log.Printf("error creating warrant: %s", err.Error())
		return "", nil, err
	}

	log.Println("warrant created successfully!")
	log.Println("Handler end -------------------------------------------------")

	return "{}", nil, err
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
