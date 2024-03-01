package _script

import (
	"log"

	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
)

var nameDesc string = "Wopta per te Vita"

func CreatePartnerhipNodes() {
	err := createBeprofNode()
	if err != nil {
		log.Println(err.Error())
	}

	err = createFacileNode()
	if err != nil {
		log.Println(err.Error())
	}
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
		Products: []models.Product{
			models.Product{
				Name:         "life",
				NameDesc:     &nameDesc,
				Version:      "v2",
				NameTitle:    "Wopta per te",
				NameSubtitle: "Vita",
				Companies: []models.Company{
					models.Company{
						Name: "axa",
					},
				},
				Steps: []models.Step{
					models.Step{
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
					},
					models.Step{
						Attributes: map[string]interface{}{
							"companyPrivacy":      "PRESTO IL CONSENSO al trattamento dei miei dati personali ad AXA France VIE S.A. – Rappresentanza Generale per l’Italia, ivi inclusi quelli eventualmente da me conferiti in riferimento al mio stato di salute, per le finalità indicate nell’informativa, consultabile all’interno dei documenti precontrattuali (ricevuti via mail o consultabili al link nella pagina che precede), nonché alla loro comunicazione, per successivo trattamento, da parte dei soggetti indicati nella informativa predetta.",
							"companyPrivacyTitle": "Privacy Assicurativa",
							"statementsEndpoint":  "question/v1/surveys",
						},
						Widget: "quotersurvey",
					},
					models.Step{
						Attributes: map[string]interface{}{
							"statementsEndpoint": "question/v1/statements",
						},
						Widget: "quoterstatements",
					},
					models.Step{
						Attributes: map[string]interface{}{
							"beneficiaryText":              "Per procedere indicare chi sono i beneficiari della polizza in caso di decesso dell'assicurato.\\nPuoi indicare genericamente i tuoi eredi legittimi e/o testamentari. Oppure inserire in maniera puntuale i nomi dei Beneficiari (Massimo due).",
							"guaranteeSlug":                "death",
							"maximumNumberOfBeneficiaries": 2,
							"thirdPartyReferenceText":      "In caso di specifiche esigenze di riservatezza, potrai indicare il nominativo ed i dati di recapito (inluso email e/o telefono) di un soggetto terno (diverso dal Beneficiario) a cui l'impresa di Assicurazione potrà rivolgersi in caso di decesso dell'Assicurato al fine di contattare il Beneficiario.",
						},
						Widget: "quoterbeneficiary",
					},
					models.Step{
						Widget: "quotercontractordata",
					},
					models.Step{
						Attributes: map[string]interface{}{
							"guaranteeSlug": "death",
						},
						Widget: "quoteruploaddocuments",
					},
					models.Step{
						Attributes: map[string]interface{}{
							"showDuration":   false,
							"showEndDate":    true,
							"showGuarantees": true,
						},
						Widget: "quoterrecap",
					},
					models.Step{
						Widget: "quotersignpay",
					},
					models.Step{
						Attributes: map[string]interface{}{
							"productLogo": "assets/images/wopta-logo-vita-magenta.png",
						},
						Widget: "quoterthankyou",
					},
				},
			},
		},
	}

	_, err := network.CreateNode(partnershipModel)

	return err
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
		Products: []models.Product{
			models.Product{
				Name:         "life",
				NameDesc:     &nameDesc,
				Version:      "v2",
				NameTitle:    "Wopta per te",
				NameSubtitle: "Vita",
				Companies: []models.Company{
					models.Company{
						Name: "axa",
					},
				},
				Steps: []models.Step{
					models.Step{
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
					},
					models.Step{
						Attributes: map[string]interface{}{
							"companyPrivacy":      "PRESTO IL CONSENSO al trattamento dei miei dati personali ad AXA France VIE S.A. – Rappresentanza Generale per l’Italia, ivi inclusi quelli eventualmente da me conferiti in riferimento al mio stato di salute, per le finalità indicate nell’informativa, consultabile all’interno dei documenti precontrattuali (ricevuti via mail o consultabili al link nella pagina che precede), nonché alla loro comunicazione, per successivo trattamento, da parte dei soggetti indicati nella informativa predetta.",
							"companyPrivacyTitle": "Privacy Assicurativa",
							"statementsEndpoint":  "question/v1/surveys",
						},
						Widget: "quotersurvey",
					},
					models.Step{
						Attributes: map[string]interface{}{
							"statementsEndpoint": "question/v1/statements",
						},
						Widget: "quoterstatements",
					},
					models.Step{
						Attributes: map[string]interface{}{
							"beneficiaryText":              "Per procedere indicare chi sono i beneficiari della polizza in caso di decesso dell'assicurato.\\nPuoi indicare genericamente i tuoi eredi legittimi e/o testamentari. Oppure inserire in maniera puntuale i nomi dei Beneficiari (Massimo due).",
							"guaranteeSlug":                "death",
							"maximumNumberOfBeneficiaries": 2,
							"thirdPartyReferenceText":      "In caso di specifiche esigenze di riservatezza, potrai indicare il nominativo ed i dati di recapito (inluso email e/o telefono) di un soggetto terno (diverso dal Beneficiario) a cui l'impresa di Assicurazione potrà rivolgersi in caso di decesso dell'Assicurato al fine di contattare il Beneficiario.",
						},
						Widget: "quoterbeneficiary",
					},
					models.Step{
						Widget: "quotercontractordata",
					},
					models.Step{
						Attributes: map[string]interface{}{
							"guaranteeSlug": "death",
						},
						Widget: "quoteruploaddocuments",
					},
					models.Step{
						Attributes: map[string]interface{}{
							"showDuration":   false,
							"showEndDate":    true,
							"showGuarantees": true,
						},
						Widget: "quoterrecap",
					},
					models.Step{
						Widget: "quotersignpay",
					},
					models.Step{
						Attributes: map[string]interface{}{
							"productLogo": "assets/images/wopta-logo-vita-magenta.png",
						},
						Widget: "quoterthankyou",
					},
				},
			},
		},
	}

	_, err := network.CreateNode(partnershipModel)

	return err
}
