package network

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

// Use network node method
func GetWarrant(filename string) *models.Warrant {
	if filename == "" {
		log.Println("[GetWarrant] no filename specified")
		return nil
	}

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
