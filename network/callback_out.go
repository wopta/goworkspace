package network

import (
	"log"

	"github.com/wopta/goworkspace/callback_out"
	"github.com/wopta/goworkspace/models"
)

func ExecuteCallback(node *models.NetworkNode, policy models.Policy) {
	if node != nil && node.CallbackConfig != nil {
		log.Println("executing node callback...")
		if err := callback_out.Handler(node.CallbackConfig.FxName, policy); err != nil {
			log.Printf("error executing node callback: %s", err.Error())
		}
	}
}
