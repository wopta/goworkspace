package _script

import (
	// "log"

	"github.com/go-jose/go-jose/v4"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/network"
)

var nameDesc string = "Wopta per te Vita"

func CreatePartnerhipNodes() {
	// var err error

	// err = createBeprofNode()
	// if err != nil {
	// 	log.Println(err.Error())
	// }

	// err = createFacileNode()
	// if err != nil {
	// 	log.Println(err.Error())
	// }

	// err = createFpinsuranceNode()
	// if err != nil {
	// 	log.Println(err.Error())
	// }

	// err = createMultiTestNode()
	// if err != nil {
	// 	log.Println(err.Error())
	// }

	// err = createELeadsNode()
	// if err != nil {
	// 	log.Println(err.Error())
	// }

	// err := createSegugioNode()
	// if err != nil {
	//	log.Println(err.Error())
	// }

	// err := createSwitchoNode()
	// if err != nil {
	// 	log.Println(err.Error())
	// }
}

func createBeprofNode() error {
	partnershipModel := models.NetworkNode{
		Uid:      "beprof",
		Code:     "beprof",
		Type:     "partnership",
		IsActive: true,
		Partnership: &models.PartnershipNode{
			Name: "beprof",
			Skin: &models.Skin{
				PrimaryColor:   "#30415f",
				SecondaryColor: "#f1bd12",
				LogoUrl:        "assets/images/logo-beprof-trasp-72dpi.png",
			},
		},
		JwtConfig: lib.JwtConfig{
			KeyName:            "BEPROF_SIGNING_KEY",
			SignatureAlgorithm: jose.HS256,
		},
		Products: []models.Product{{
			Name:         "life",
			NameDesc:     &nameDesc,
			Version:      "v2",
			NameTitle:    "Wopta per te",
			NameSubtitle: "Vita",
			Companies: []models.Company{{
				Name: "axa",
			}},
			Steps: []models.Step{{
				Children: []models.Child{{
					Attributes: map[string]interface{}{
						"consens": "Il sottoscritto, letta e compresa l'informativa sul trattamento dei dati personali, ACCONSENTE al trattamento dei propri dati personali da parte di Wopta Assicurazioni per l'invio di comunicazioni e proposte commerciali e di marketing, incluso l'invio di newsletter e ricerche di mercato, attraverso strumenti automatizzati (sms, mms, e-mail, ecc.) e non (posta cartacea e telefono con operatore).",
						"key":     2,
						"title":   "Privacy",
					},
					Widget: "privacyConsent",
				}},
				Widget: "guaranteeconfigurationstep",
			}, {
				Attributes: map[string]interface{}{
					"companyPrivacy":      "PRESTO IL CONSENSO al trattamento dei miei dati personali ad AXA France VIE S.A. – Rappresentanza Generale per l’Italia, ivi inclusi quelli eventualmente da me conferiti in riferimento al mio stato di salute, per le finalità indicate nell’informativa, consultabile all’interno dei documenti precontrattuali (ricevuti via mail o consultabili al link nella pagina che precede), nonché alla loro comunicazione, per successivo trattamento, da parte dei soggetti indicati nella informativa predetta.",
					"companyPrivacyTitle": "Privacy Assicurativa",
					"statementsEndpoint":  "question/v1/surveys",
				},
				Widget: "quotersurvey",
			}, {
				Attributes: map[string]interface{}{
					"statementsEndpoint": "question/v1/statements",
				},
				Widget: "quoterstatements",
			}, {
				Attributes: map[string]interface{}{
					"beneficiaryText":              "Per procedere indicare chi sono i beneficiari della polizza in caso di decesso dell'assicurato.\\nPuoi indicare genericamente i tuoi eredi legittimi e/o testamentari. Oppure inserire in maniera puntuale i nomi dei Beneficiari (Massimo due).",
					"guaranteeSlug":                "death",
					"maximumNumberOfBeneficiaries": 2,
					"thirdPartyReferenceText":      "In caso di specifiche esigenze di riservatezza, potrai indicare il nominativo ed i dati di recapito (inluso email e/o telefono) di un soggetto terno (diverso dal Beneficiario) a cui l'impresa di Assicurazione potrà rivolgersi in caso di decesso dell'Assicurato al fine di contattare il Beneficiario.",
				},
				Widget: "quoterbeneficiary",
			}, {
				Widget: "quotercontractordata",
			}, {
				Attributes: map[string]interface{}{
					"guaranteeSlug": "death",
				},
				Widget: "quoteruploaddocuments",
			}, {
				Attributes: map[string]interface{}{
					"showDuration":   false,
					"showEndDate":    true,
					"showGuarantees": true,
				},
				Widget: "quoterrecap",
			}, {
				Widget: "quotersignpay",
			}, {
				Attributes: map[string]interface{}{
					"productLogo": "assets/images/wopta-logo-vita-magenta.png",
				},
				Widget: "quoterthankyou",
			}},
		}},
	}

	nn, err := network.CreateNode(partnershipModel)
	if err != nil {
		return err
	}

	return nn.SaveBigQuery("")
}

func createFacileNode() error {
	partnershipModel := models.NetworkNode{
		Uid:      "facile",
		Code:     "facile",
		Type:     "partnership",
		IsActive: true,
		Partnership: &models.PartnershipNode{
			Name: "facile",
			Skin: &models.Skin{
				PrimaryColor:   "",
				SecondaryColor: "",
				LogoUrl:        "https://upload.wikimedia.org/wikipedia/commons/7/78/Logo_facile_%28azienda%29.png",
			},
		},
		JwtConfig: lib.JwtConfig{
			KeyName:           "FACILE_SIGNING_KEY",
			KeyAlgorithm:      jose.DIRECT,
			ContentEncryption: jose.A128CBC_HS256,
		},
		Products: []models.Product{{
			Name:         "life",
			NameDesc:     &nameDesc,
			Version:      "v2",
			NameTitle:    "Wopta per te",
			NameSubtitle: "Vita",
			Companies: []models.Company{{
				Name: "axa",
			}},
			Steps: []models.Step{{
				Children: []models.Child{
					{
						Attributes: map[string]interface{}{
							"consens": "Il sottoscritto, letta e compresa l'informativa sul trattamento dei dati personali, ACCONSENTE al trattamento dei propri dati personali da parte di Wopta Assicurazioni per l'invio di comunicazioni e proposte commerciali e di marketing, incluso l'invio di newsletter e ricerche di mercato, attraverso strumenti automatizzati (sms, mms, e-mail, ecc.) e non (posta cartacea e telefono con operatore).",
							"key":     2,
							"title":   "Privacy",
						},
						Widget: "privacyConsent",
					}},
				Widget: "guaranteeconfigurationstep",
			}, {
				Attributes: map[string]interface{}{
					"companyPrivacy":      "PRESTO IL CONSENSO al trattamento dei miei dati personali ad AXA France VIE S.A. – Rappresentanza Generale per l’Italia, ivi inclusi quelli eventualmente da me conferiti in riferimento al mio stato di salute, per le finalità indicate nell’informativa, consultabile all’interno dei documenti precontrattuali (ricevuti via mail o consultabili al link nella pagina che precede), nonché alla loro comunicazione, per successivo trattamento, da parte dei soggetti indicati nella informativa predetta.",
					"companyPrivacyTitle": "Privacy Assicurativa",
					"statementsEndpoint":  "question/v1/surveys",
				},
				Widget: "quotersurvey",
			}, {
				Attributes: map[string]interface{}{
					"statementsEndpoint": "question/v1/statements",
				},
				Widget: "quoterstatements",
			}, {
				Attributes: map[string]interface{}{
					"beneficiaryText":              "Per procedere indicare chi sono i beneficiari della polizza in caso di decesso dell'assicurato.\\nPuoi indicare genericamente i tuoi eredi legittimi e/o testamentari. Oppure inserire in maniera puntuale i nomi dei Beneficiari (Massimo due).",
					"guaranteeSlug":                "death",
					"maximumNumberOfBeneficiaries": 2,
					"thirdPartyReferenceText":      "In caso di specifiche esigenze di riservatezza, potrai indicare il nominativo ed i dati di recapito (inluso email e/o telefono) di un soggetto terno (diverso dal Beneficiario) a cui l'impresa di Assicurazione potrà rivolgersi in caso di decesso dell'Assicurato al fine di contattare il Beneficiario.",
				},
				Widget: "quoterbeneficiary",
			}, {
				Widget: "quotercontractordata",
			}, {
				Attributes: map[string]interface{}{
					"guaranteeSlug": "death",
				},
				Widget: "quoteruploaddocuments",
			}, {
				Attributes: map[string]interface{}{
					"showDuration":   false,
					"showEndDate":    true,
					"showGuarantees": true,
				},
				Widget: "quoterrecap",
			}, {
				Widget: "quotersignpay",
			}, {
				Attributes: map[string]interface{}{
					"productLogo": "assets/images/wopta-logo-vita-magenta.png",
				},
				Widget: "quoterthankyou",
			}},
		}},
	}

	nn, err := network.CreateNode(partnershipModel)
	if err != nil {
		return err
	}

	return nn.SaveBigQuery("")
}

func createFpinsuranceNode() error {
	var (
		err error
		nn  *models.NetworkNode
	)
	partnershipModel := models.NetworkNode{
		Uid:         models.PartnershipFpinsurance,
		Code:        models.PartnershipFpinsurance,
		Type:        models.PartnershipNetworkNodeType,
		IsActive:    true,
		Partnership: &models.PartnershipNode{Name: models.PartnershipFpinsurance},
		Products: []models.Product{{
			Name:         models.LifeProduct,
			NameDesc:     &nameDesc,
			Version:      models.ProductV2,
			NameTitle:    "Wopta per te",
			NameSubtitle: "Vita",
			Companies:    []models.Company{{Name: "axa"}},
		}},
	}

	nn, err = network.CreateNode(partnershipModel)
	if err != nil {
		return err
	}

	return nn.SaveBigQuery("")
}

func createSegugioNode() error {
	var (
		err error
		nn  *models.NetworkNode
	)
	partnershipModel := models.NetworkNode{
		Uid:         models.PartnershipSegugio,
		Code:        models.PartnershipSegugio,
		Type:        models.PartnershipNetworkNodeType,
		IsActive:    true,
		Partnership: &models.PartnershipNode{Name: models.PartnershipSegugio},
		Products: []models.Product{{
			Name:         models.LifeProduct,
			NameDesc:     &nameDesc,
			Version:      models.ProductV2,
			NameTitle:    "Wopta per te",
			NameSubtitle: "Vita",
			Companies:    []models.Company{{Name: "axa"}},
		}},
	}

	nn, err = network.CreateNode(partnershipModel)
	if err != nil {
		return err
	}

	return nn.SaveBigQuery("")
}

func createMultiTestNode() error {
	var (
		err             error
		partnershipName string = "multi-test"
		nameDesc2       string = "Wopta per te Persona"
	)

	partnershipModel := models.NetworkNode{
		Uid:         partnershipName,
		Code:        partnershipName,
		Type:        models.PartnershipNetworkNodeType,
		IsActive:    true,
		Partnership: &models.PartnershipNode{Name: partnershipName},
		Products: []models.Product{{
			Name:         models.LifeProduct,
			NameDesc:     &nameDesc,
			Version:      models.ProductV2,
			NameTitle:    "Wopta per te",
			NameSubtitle: "Vita",
			Companies:    []models.Company{{Name: "axa"}},
		}, {
			Name:         models.PersonaProduct,
			NameDesc:     &nameDesc2,
			Version:      models.ProductV1,
			NameTitle:    "Wopta per te",
			NameSubtitle: "Persona",
			Companies:    []models.Company{{Name: "global"}},
		}},
	}

	nn, err := network.CreateNode(partnershipModel)
	if err != nil {
		return err
	}

	return nn.SaveBigQuery("")
}

func createELeadsNode() error {
	partnershipModel := models.NetworkNode{
		Uid:      "eleads",
		Code:     "eleads",
		Type:     "partnership",
		IsActive: true,
		Partnership: &models.PartnershipNode{
			Name: "eleads",
		},
		JwtConfig: lib.JwtConfig{
			KeyName:           "ELEADS_SIGNING_KEY",
			KeyAlgorithm:      jose.DIRECT,
			ContentEncryption: jose.A128CBC_HS256,
		},
		Products: []models.Product{{
			Name:         "life",
			NameDesc:     &nameDesc,
			Version:      "v2",
			NameTitle:    "Wopta per te",
			NameSubtitle: "Vita",
			Companies: []models.Company{{
				Name: "axa",
			}},
			Steps: []models.Step{{
				Children: []models.Child{
					{
						Attributes: map[string]interface{}{
							"consens": "Il sottoscritto, letta e compresa l'informativa sul trattamento dei dati personali, ACCONSENTE al trattamento dei propri dati personali da parte di Wopta Assicurazioni per l'invio di comunicazioni e proposte commerciali e di marketing, incluso l'invio di newsletter e ricerche di mercato, attraverso strumenti automatizzati (sms, mms, e-mail, ecc.) e non (posta cartacea e telefono con operatore).",
							"key":     2,
							"title":   "Privacy",
						},
						Widget: "privacyConsent",
					}},
				Widget: "guaranteeconfigurationstep",
			}, {
				Attributes: map[string]interface{}{
					"companyPrivacy":      "PRESTO IL CONSENSO al trattamento dei miei dati personali ad AXA France VIE S.A. – Rappresentanza Generale per l’Italia, ivi inclusi quelli eventualmente da me conferiti in riferimento al mio stato di salute, per le finalità indicate nell’informativa, consultabile all’interno dei documenti precontrattuali (ricevuti via mail o consultabili al link nella pagina che precede), nonché alla loro comunicazione, per successivo trattamento, da parte dei soggetti indicati nella informativa predetta.",
					"companyPrivacyTitle": "Privacy Assicurativa",
					"statementsEndpoint":  "question/v1/surveys",
				},
				Widget: "quotersurvey",
			}, {
				Attributes: map[string]interface{}{
					"statementsEndpoint": "question/v1/statements",
				},
				Widget: "quoterstatements",
			}, {
				Attributes: map[string]interface{}{
					"beneficiaryText":              "Per procedere indicare chi sono i beneficiari della polizza in caso di decesso dell'assicurato.\\nPuoi indicare genericamente i tuoi eredi legittimi e/o testamentari. Oppure inserire in maniera puntuale i nomi dei Beneficiari (Massimo due).",
					"guaranteeSlug":                "death",
					"maximumNumberOfBeneficiaries": 2,
					"thirdPartyReferenceText":      "In caso di specifiche esigenze di riservatezza, potrai indicare il nominativo ed i dati di recapito (inluso email e/o telefono) di un soggetto terno (diverso dal Beneficiario) a cui l'impresa di Assicurazione potrà rivolgersi in caso di decesso dell'Assicurato al fine di contattare il Beneficiario.",
				},
				Widget: "quoterbeneficiary",
			}, {
				Widget: "quotercontractordata",
			}, {
				Attributes: map[string]interface{}{
					"guaranteeSlug": "death",
				},
				Widget: "quoteruploaddocuments",
			}, {
				Attributes: map[string]interface{}{
					"showDuration":   false,
					"showEndDate":    true,
					"showGuarantees": true,
				},
				Widget: "quoterrecap",
			}, {
				Widget: "quotersignpay",
			}, {
				Attributes: map[string]interface{}{
					"productLogo": "assets/images/wopta-logo-vita-magenta.png",
				},
				Widget: "quoterthankyou",
			}},
		}},
	}

	nn, err := network.CreateNode(partnershipModel)
	if err != nil {
		return err
	}

	return nn.SaveBigQuery("")
}

func createSwitchoNode() error {
	var (
		err error
		nn  *models.NetworkNode
	)
	partnershipModel := models.NetworkNode{
		Uid:      models.PartnershipSwitcho,
		Code:     models.PartnershipSwitcho,
		Type:     models.PartnershipNetworkNodeType,
		IsActive: true,
		Partnership: &models.PartnershipNode{
			Name: models.PartnershipSwitcho,
			Skin: &models.Skin{
				LogoUrl: "https://storage.googleapis.com/wopta-public-assets/logo_switcho.png",
			},
		},
		Products: []models.Product{{
			Name:         models.LifeProduct,
			NameDesc:     &nameDesc,
			Version:      models.ProductV2,
			NameTitle:    "Wopta per te",
			NameSubtitle: "Vita",
			Companies:    []models.Company{{Name: "axa"}},
		}},
	}

	nn, err = network.CreateNode(partnershipModel)
	if err != nil {
		return err
	}

	return nn.SaveBigQuery("")
}

func CreateAdvTestNode() error {
	var (
		err error
		nn  *models.NetworkNode
	)
	partnershipModel := models.NetworkNode{
		Uid:      "adv1",
		Code:     "adv1",
		Type:     models.PartnershipNetworkNodeType,
		IsActive: true,
		Partnership: &models.PartnershipNode{
			Name: "adv1",
		},
		Products: []models.Product{{
			Name:         models.LifeProduct,
			NameDesc:     &nameDesc,
			Version:      models.ProductV2,
			NameTitle:    "Wopta per te",
			NameSubtitle: "Vita",
			Companies:    []models.Company{{Name: "axa"}},
			ConsultancyConfig: &models.ConsultancyConfig{
				Min:            0.05,
				Max:            0.05,
				Step:           0,
				DefaultValue:   0.05,
				IsActive:       true,
				IsConfigurable: false,
			},
		}},
	}

	nn, err = network.CreateNode(partnershipModel)
	if err != nil {
		return err
	}

	return nn.SaveBigQuery("")
}
