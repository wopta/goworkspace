package partnership

import "github.com/wopta/goworkspace/models"

type PartnershipNode struct {
	Name string       `json:"name"`
	Skin *models.Skin `json:"skin,omitempty"`
}

type PartnershipResponse struct {
	Policy      models.Policy   `json:"policy"`
	Partnership PartnershipNode `json:"partnership"`
	Product     models.Product  `json:"product"`
}
