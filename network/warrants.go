package network

import (
	"encoding/json"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/lib/log"
	"github.com/wopta/goworkspace/models"
)

func GetWarrants() ([]models.Warrant, error) {
	var (
		err      error
		warrants []models.Warrant
	)
	log.AddPrefix("GetWarrants")
	defer log.PopPrefix()
	warrantsBytes := lib.GetFolderContentByEnv(models.WarrantsFolder)

	for _, warrantBytes := range warrantsBytes {
		var warrant models.Warrant
		err = json.Unmarshal(warrantBytes, &warrant)
		if err != nil {
			log.ErrorF("error unmarshaling warrant: %s", err.Error())
			return warrants, err
		}

		warrants = append(warrants, warrant)
	}
	return warrants, nil
}
