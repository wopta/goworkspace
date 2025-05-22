package constants

import "gitlab.dev.wopta.it/goworkspace/models"

const WoptaLogo string = "logo_wopta.png"

var ProductLogoMap map[string]string = map[string]string{
	models.LifeProduct:    "logo_vita.png",
	models.PmiProduct:     "logo_pmi.png",
	models.PersonaProduct: "logo_persona.png",
	models.GapProduct:     "logo_gap.png",
}

var CompanyLogoMap map[string]string = map[string]string{
	models.AxaCompany:      "logo_axa.png",
	models.GlobalCompany:   "logo_global.png",
	models.SogessurCompany: "logo_sogessur.png",
}
