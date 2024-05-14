package win

import (
	"log"

	"github.com/wopta/goworkspace/models"
)

func approvalCallback(policy models.Policy) error {
	log.Println("win wait for approval calback...")
	return nil
}
