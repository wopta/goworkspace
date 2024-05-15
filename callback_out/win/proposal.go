package win

import (
	"log"
	"net/http"

	"github.com/wopta/goworkspace/models"
)

func proposalCallback(policy models.Policy) (*http.Request, *http.Response, error) {
	log.Println("win proposal calback...")
	return emitCallback(policy)
}
