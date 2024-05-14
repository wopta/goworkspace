package win

import (
	"log"

	"github.com/wopta/goworkspace/models"
)

func contractCallback(policy models.Policy) error {
	log.Println("win contract calback...")
	return nil
}
