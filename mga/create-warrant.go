package mga

import (
	"encoding/json"
	"github.com/wopta/goworkspace/lib/log"
	"io"
	"net/http"
	"os"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func CreateWarrantFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var warrant models.Warrant

	log.AddPrefix("CreateWarrantFx")
	defer log.PopPrefix()

	log.Println("Handler start -----------------------------------------------")

	bodyBytes := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	err := json.Unmarshal(bodyBytes, &warrant)
	if err != nil {
		log.ErrorF("error unmarshaling request: %s", err.Error())
		return "", nil, err
	}

	log.Println("creating warrant...")
	err = CreateWarrant(warrant)
	if err != nil {
		log.ErrorF("error creating warrant: %s", err.Error())
		return "", nil, err
	}

	log.Println("warrant created successfully!")
	log.Println("Handler end -------------------------------------------------")

	return "{}", nil, err
}

func CreateWarrant(warrant models.Warrant) error {
	log.AddPrefix("CreateWarrant")
	defer log.PopPrefix()
	fileName := models.WarrantsFolder + warrant.Name + ".json"

	bytesToWrite, err := json.Marshal(warrant)
	if err != nil {
		log.ErrorF("error marshaling warrant: %s", err.Error())
		return err
	}

	_, err = lib.PutToStorageIfNotExists(os.Getenv("GOOGLE_STORAGE_BUCKET"), fileName, bytesToWrite)

	if err != nil {
		log.ErrorF("error writing warrant: %s", err.Error())
	}

	return err
}
