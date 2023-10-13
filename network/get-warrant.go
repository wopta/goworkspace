package network

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func GetWarrant(filename string) *models.Warrant {
	log.Printf("[GetWarrant] requesting warrant %s", filename)

	var (
		warrant       models.Warrant
		warrantFormat string = "warrants/%s.json"
	)

	warrantBytes := lib.GetFilesByEnv(fmt.Sprintf(warrantFormat, filename))
	err := json.Unmarshal(warrantBytes, &warrant)
	if err != nil {
		log.Printf("[GetWarrant] error unmarshaling warrant %s: %s", filename, err.Error())
		return nil
	}

	return &warrant
}
