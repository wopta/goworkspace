package win

import (
	"log"

	"github.com/wopta/goworkspace/models"
)

func proposalCallback(policy models.Policy) error {
	log.Println("win proposal calback...")
	return emitCallback(policy)
}
