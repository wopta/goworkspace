package partnership

import "github.com/wopta/goworkspace/models"

var encryptedPartnerships map[string]bool = map[string]bool{
	models.PartnershipBeProf:      false,
	models.PartnershipFacile:      true,
	models.PartnershipFpinsurance: false,
}
