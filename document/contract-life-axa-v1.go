package document

import (
	"github.com/go-pdf/fpdf"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"strings"
)

func lifeAxaContractV1(pdf *fpdf.Fpdf, origin string, policy *models.Policy, networkNode *models.NetworkNode, product *models.Product) (string, []byte) {
	signatureID = 0

	mainHeader(pdf, policy, false)

	mainFooter(pdf, policy.Name)

	pdf.AddPage()

	insuredInfoSection(pdf, policy)

	guaranteesMap, slugs := loadLifeGuarantees(policy, product)

	lifeGuaranteesTable(pdf, guaranteesMap, slugs)

	avvertenzeBeneficiariSection(pdf)

	beneficiaries, legitimateSuccessorsChoice, designatedSuccessorsChoice := loadLifeBeneficiariesInfo(policy)

	beneficiariesSection(pdf, beneficiaries, legitimateSuccessorsChoice, designatedSuccessorsChoice)

	beneficiaryReferenceSection(pdf, policy)

	surveysSection(pdf, policy, false)

	pdf.AddPage()

	statementsSection(pdf, policy, false)

	offerResumeSection(pdf, policy)

	paymentResumeSection(pdf, policy)

	contractWithdrawlSection(pdf, false)

	pdf.AddPage()

	paymentMethodSection(pdf)

	emitResumeSection(pdf, policy)

	companiesDescriptionSection(pdf, policy.Company)

	axaHeader(pdf)

	pdf.AddPage()

	axaFooter(pdf)

	axaDeclarationsConsentSection(pdf, policy, false)

	pdf.AddPage()

	axaTableSection(pdf, policy)

	pdf.AddPage()

	axaTablePart2Section(pdf, policy, false)

	pdf.Ln(15)

	axaTablePart3Section(pdf)

	woptaHeader(pdf)

	pdf.AddPage()

	woptaFooter(pdf)

	producerInfo := loadProducerInfo(origin, networkNode)

	allegato3Section(pdf, producerInfo)

	pdf.AddPage()

	allegato4Section(pdf, producerInfo)

	pdf.AddPage()

	allegato4TerSection(pdf, producerInfo)

	pdf.AddPage()

	woptaPrivacySection(pdf)

	personalDataHandlingSection(pdf, policy, false)

	filename, out := saveContract(pdf, policy)
	return filename, out
}

func insuredInfoSection(pdf *fpdf.Fpdf, policy *models.Policy) {
	getParagraphTitle(pdf, "La tua assicurazione è operante per il seguente Assicurato e Garanzie")
	insuredInfoTable(pdf, policy.Assets[0].Person)
}

func insuredInfoTable(pdf *fpdf.Fpdf, insured *models.User) {
	residenceAddress := strings.ToUpper(insured.Residence.StreetName + ", " + insured.Residence.StreetNumber +
		" - " + insured.Residence.PostalCode + " " + insured.Residence.City + " (" + insured.Residence.CityCode + ")")

	drawPinkHorizontalLine(pdf, thickLineWidth)
	pdf.Ln(2)
	setBlackBoldFont(pdf, standardTextSize)
	pdf.Cell(20, 2, "Cognome e Nome")
	setBlackRegularFont(pdf, standardTextSize)
	pdf.SetX(pdf.GetX() + 24)
	pdf.Cell(20, 2, strings.ToUpper(insured.Surname+" "+insured.Name))
	pdf.SetX(pdf.GetX() + 60)
	setBlackBoldFont(pdf, standardTextSize)
	pdf.Cell(10, 2, "Codice fiscale:")
	pdf.SetX(pdf.GetX() + 20)
	setBlackRegularFont(pdf, standardTextSize)
	pdf.Cell(20, 2, strings.ToUpper(insured.FiscalCode))
	pdf.Ln(2.5)
	drawPinkHorizontalLine(pdf, thinLineWidth)
	pdf.Ln(2)
	setBlackBoldFont(pdf, standardTextSize)
	pdf.Cell(20, 2, "Residente in")
	setBlackRegularFont(pdf, standardTextSize)
	pdf.SetX(pdf.GetX() + 24)
	pdf.Cell(20, 2, residenceAddress)
	pdf.Ln(2.5)
	drawPinkHorizontalLine(pdf, thinLineWidth)
	pdf.Ln(2)
	setBlackBoldFont(pdf, standardTextSize)
	pdf.Cell(20, 2, "Mail")
	setBlackRegularFont(pdf, standardTextSize)
	pdf.SetX(pdf.GetX() + 24)
	pdf.Cell(20, 2, insured.Mail)
	pdf.SetX(pdf.GetX() + 60)
	setBlackBoldFont(pdf, standardTextSize)
	pdf.Cell(10, 2, "Telefono:")
	pdf.SetX(pdf.GetX() + 20)
	setBlackRegularFont(pdf, 9)
	pdf.Cell(20, 2, insured.Phone)
	pdf.Ln(3)
	drawPinkHorizontalLine(pdf, thickLineWidth)
	pdf.Ln(1)
}

func lifeGuaranteesTable(pdf *fpdf.Fpdf, guaranteesMap map[string]map[string]string, slugs []slugStruct) {
	setBlackBoldFont(pdf, standardTextSize)
	pdf.MultiCell(90, 3, "Garanzie", "", "CM", false)
	pdf.SetXY(pdf.GetX()+90, pdf.GetY()-3)
	pdf.MultiCell(30, 3, "Somma\nassicurata €", "", "CM", false)
	pdf.SetXY(pdf.GetX()+117, pdf.GetY()-6)
	pdf.MultiCell(20, 3, "Durata\nanni", "", "CM", false)
	pdf.SetXY(pdf.GetX()+140, pdf.GetY()-6)
	pdf.MultiCell(25, 3, "Scade il", "", "CM", false)
	pdf.SetXY(pdf.GetX()+169, pdf.GetY()-3)
	pdf.MultiCell(0, 3, "Premio annuale €", "", "CM", false)
	pdf.Ln(1)
	drawPinkHorizontalLine(pdf, thinLineWidth)

	for _, slug := range slugs {
		setBlackBoldFont(pdf, standardTextSize)
		pdf.CellFormat(90, 6, guaranteesMap[slug.name]["name"], "", 0, "", false, 0, "")
		setBlackRegularFont(pdf, standardTextSize)
		pdf.CellFormat(25, 6, guaranteesMap[slug.name]["sumInsuredLimitOfIndemnity"],
			"", 0, "RM", false, 0, "")
		pdf.CellFormat(25, 6, guaranteesMap[slug.name]["duration"], "", 0, "CM",
			false, 0, "")
		pdf.CellFormat(25, 6, guaranteesMap[slug.name]["endDate"], "", 0, "CM", false, 0, "")
		pdf.CellFormat(0, 6, guaranteesMap[slug.name]["price"], "",
			0, "RM", false, 0, "")
		pdf.Ln(5)
		drawPinkHorizontalLine(pdf, thinLineWidth)
	}
	pdf.Ln(0.5)
	setBlackRegularFont(pdf, smallTextSize)
	pdf.Cell(80, 3, "(*) imposte assicurative di legge incluse nella misura del 2,50% del premio imponibile")
	pdf.Ln(5)
}

func avvertenzeBeneficiariSection(pdf *fpdf.Fpdf) {
	getParagraphTitle(pdf, "Nomina dei Beneficiari e Referente terzo, per il caso di garanzia Decesso "+
		"(qualora sottoscritta)")
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, "AVVERTENZE: Può scegliere se designare nominativamente i beneficiari o se "+
		"designare genericamente come beneficiari i suoi eredi legittimi e/o testamentari. In caso di mancata "+
		"designazione nominativa, la Compagnia potrà incontrare, al decesso dell’Assicurato, maggiori difficoltà "+
		"nell’identificazione e nella ricerca dei beneficiari. La modifica o revoca del/i beneficiario/i deve essere "+
		"comunicata alla Compagnia in forma scritta.\nIn caso di specifiche esigenze di riservatezza, la Compagnia "+
		"potrà rivolgersi ad un soggetto terzo (diverso dal Beneficiario) in caso di Decesso al fine di contattare "+
		"il Beneficiario designato.", "", "", false)
	pdf.Ln(3)
}

func beneficiariesSection(pdf *fpdf.Fpdf, beneficiaries []map[string]string, legitimateSuccessorsChoice,
	designatedSuccessorsChoice string) {
	getParagraphTitle(pdf, "Beneficiario")
	setBlackRegularFont(pdf, standardTextSize)
	pdf.CellFormat(0, 3, "Io sottoscritto Assicurato, con la sottoscrizione della presente polizza, in "+
		"riferimento alla garanzia Decesso:", "", 0, "", false, 0, "")
	pdf.Ln(4)
	setBlackDrawColor(pdf)
	pdf.SetX(11)
	pdf.CellFormat(3, 3, legitimateSuccessorsChoice, "1", 0, "CM", false, 0, "")
	pdf.CellFormat(0, 3, "Designo genericamente quali beneficiari della prestazione i miei eredi "+
		"(legittimi e/o testamentari)", "", 0, "", false, 0, "")
	pdf.Ln(4)
	pdf.SetX(11)
	pdf.CellFormat(3, 3, designatedSuccessorsChoice, "1", 0, "CM", false, 0, "")
	pdf.CellFormat(0, 3, "Designo nominativamente il/i seguente/i soggetto/i quale beneficiario/i della "+
		"prestazione", "", 0, "", false, 0, "")
	pdf.Ln(5)
	beneficiariesTable(pdf, beneficiaries)
	pdf.Ln(1)
}

func beneficiariesTable(pdf *fpdf.Fpdf, beneficiaries []map[string]string) {
	for _, beneficiary := range beneficiaries {
		drawPinkHorizontalLine(pdf, thickLineWidth)
		pdf.Ln(1.5)
		setBlackBoldFont(pdf, standardTextSize)
		pdf.Cell(50, 2, "Cognome e nome")
		setBlackRegularFont(pdf, standardTextSize)
		pdf.Cell(20, 2, beneficiary["name"])
		pdf.SetX(pdf.GetX() + 60)
		setBlackBoldFont(pdf, standardTextSize)
		pdf.Cell(20, 2, "Cod. Fisc.: ")
		setBlackRegularFont(pdf, standardTextSize)
		pdf.Cell(20, 2, beneficiary["fiscCode"])
		pdf.Ln(3)
		drawPinkHorizontalLine(pdf, thinLineWidth)
		pdf.Ln(2)
		setBlackBoldFont(pdf, standardTextSize)
		pdf.Cell(50, 2, "Indirizzo")
		setBlackRegularFont(pdf, standardTextSize)
		pdf.Cell(20, 2, beneficiary["address"])
		pdf.SetX(pdf.GetX() + 60)
		setBlackBoldFont(pdf, standardTextSize)
		pdf.Ln(3)
		drawPinkHorizontalLine(pdf, thinLineWidth)
		pdf.Ln(2)
		setBlackBoldFont(pdf, standardTextSize)
		pdf.Cell(50, 2, "Mail")
		setBlackRegularFont(pdf, standardTextSize)
		pdf.Cell(20, 2, beneficiary["mail"])
		pdf.SetX(pdf.GetX() + 60)
		setBlackBoldFont(pdf, standardTextSize)
		pdf.Cell(20, 2, "Telefono: ")
		setBlackRegularFont(pdf, standardTextSize)
		pdf.Cell(20, 2, beneficiary["phone"])
		pdf.Ln(3)
		drawPinkHorizontalLine(pdf, thinLineWidth)
		pdf.Ln(2)
		setBlackBoldFont(pdf, standardTextSize)
		pdf.Cell(50, 2, "Relazione con Assicurato")
		setBlackRegularFont(pdf, standardTextSize)
		pdf.Cell(20, 2, beneficiary["relation"])
		pdf.Ln(3)
		drawPinkHorizontalLine(pdf, thinLineWidth)
		pdf.Ln(2)
		pdf.Cell(165, 2, "Consenso ad invio comunicazioni da parte della Compagnia al beneficiario, prima "+
			"dell'evento Decesso:")
		pdf.Cell(20, 2, beneficiary["contactConsent"])
		pdf.Ln(3)
		drawPinkHorizontalLine(pdf, thinLineWidth)
		pdf.Ln(2)
	}
}

func beneficiaryReferenceSection(pdf *fpdf.Fpdf, policy *models.Policy) {
	beneficiaryReference := map[string]string{
		"name":     "=====",
		"fiscCode": "=====",
		"address":  "=====",
		"mail":     "=====",
		"phone":    "=====",
	}

	deathGuarantee, err := policy.ExtractGuarantee("death")
	lib.CheckError(err)

	if deathGuarantee.BeneficiaryReference != nil {
		beneficiary := deathGuarantee.BeneficiaryReference
		address := strings.ToUpper(beneficiary.Residence.StreetName + ", " + beneficiary.Residence.StreetNumber +
			" - " + beneficiary.Residence.PostalCode + " " + beneficiary.Residence.City +
			" (" + beneficiary.Residence.CityCode + ")")
		beneficiaryReference["name"] = strings.ToUpper(beneficiary.Surname + " " + beneficiary.Name)
		beneficiaryReference["fiscCode"] = strings.ToUpper(beneficiary.FiscalCode)
		beneficiaryReference["address"] = address
		beneficiaryReference["mail"] = beneficiary.Mail
		beneficiaryReference["phone"] = beneficiary.Phone
	}

	getParagraphTitle(pdf, "Referente terzo")
	beneficiaryReferenceTable(pdf, beneficiaryReference)
	pdf.Ln(2)
}

func beneficiaryReferenceTable(pdf *fpdf.Fpdf, beneficiaryReference map[string]string) {
	drawPinkHorizontalLine(pdf, thickLineWidth)
	pdf.Ln(1.5)
	setBlackBoldFont(pdf, standardTextSize)
	pdf.Cell(50, 2, "Cognome e nome")
	setBlackRegularFont(pdf, standardTextSize)
	pdf.Cell(20, 2, beneficiaryReference["name"])
	pdf.SetX(pdf.GetX() + 60)
	setBlackBoldFont(pdf, standardTextSize)
	pdf.Cell(20, 2, "Cod. Fisc.: ")
	setBlackRegularFont(pdf, standardTextSize)
	pdf.Cell(20, 2, beneficiaryReference["fiscCode"])
	pdf.Ln(3)
	drawPinkHorizontalLine(pdf, thinLineWidth)
	pdf.Ln(2)
	setBlackBoldFont(pdf, standardTextSize)
	pdf.Cell(50, 2, "Indirizzo")
	setBlackRegularFont(pdf, standardTextSize)
	pdf.Cell(20, 2, beneficiaryReference["address"])
	pdf.SetX(pdf.GetX() + 60)
	setBlackBoldFont(pdf, standardTextSize)
	pdf.Ln(3)
	drawPinkHorizontalLine(pdf, thinLineWidth)
	pdf.Ln(2)
	setBlackBoldFont(pdf, standardTextSize)
	pdf.Cell(50, 2, "Mail")
	setBlackRegularFont(pdf, standardTextSize)
	pdf.Cell(20, 2, beneficiaryReference["mail"])
	pdf.SetX(pdf.GetX() + 60)
	setBlackBoldFont(pdf, standardTextSize)
	pdf.Cell(20, 2, "Telefono: ")
	setBlackRegularFont(pdf, standardTextSize)
	pdf.Cell(20, 2, beneficiaryReference["phone"])
	pdf.Ln(3)
	drawPinkHorizontalLine(pdf, thinLineWidth)
}

func surveysSection(pdf *fpdf.Fpdf, policy *models.Policy, isProposal bool) {
	surveys := *policy.Surveys

	if policy.PartnershipName == models.PartnershipBeProf {
		surveys[0].Questions[len(surveys[0].Questions)-1].Question += " In caso anche di una sola risposta positiva, ovvero in caso di somme" +
			" assicurate per le garanzie Decesso e/o Invalidità Totale Permanente da Infortunio o Malattia superiori" +
			" a 200.000 €, è richiesto che l’Assicurato si sottoponga a visita medica come indicato al punto c)" +
			" che precede."
	}

	getParagraphTitle(pdf, "Dichiarazioni da leggere con attenzione prima di firmare")
	err := printSurvey(pdf, surveys[0], policy.Company, isProposal)
	lib.CheckError(err)

	pdf.AddPage()

	getParagraphTitle(pdf, "Questionario Medico")
	for _, survey := range surveys[1:] {
		err := printSurvey(pdf, survey, policy.Company, isProposal)
		lib.CheckError(err)
	}

	if policy.IsReserved {
		pdf.Ln(3)
		setBlackRegularFont(pdf, standardTextSize)
		pdf.MultiCell(0, 3, "Nota  bene:  Ai  fini  della  valutazione  del  rischio,  l’Assicurato  ha  "+
			"inviato  alla  compagnia  un  Rapporto  di  Visita  Medica,  sottoscritto  dal medico curante, che  "+
			"costituisce  parte  integrante della  presente  Polizza  e  la Compagnia, valutato il  rischio,  ha  "+
			"accettato  il rischio  alle condizioni indicate nella presente Polizza", "", fpdf.AlignLeft, false)

		setBlackBoldFont(pdf, standardTextSize)
	}
}

func statementsSection(pdf *fpdf.Fpdf, policy *models.Policy, isProposal bool) {
	statements := *policy.Statements
	for _, statement := range statements {
		printStatement(pdf, statement, policy.Company, isProposal)
	}
}

func offerResumeSection(pdf *fpdf.Fpdf, policy *models.Policy) {
	var (
		paymentSplit string
		tableInfo    [][]string
	)

	switch policy.PaymentSplit {
	case string(models.PaySplitMonthly):
		paymentSplit = "MENSILE"
		tableInfo = [][]string{
			{
				"Mensile firma del contratto",
				lib.HumanaizePriceEuro(policy.OffersPrices["default"]["monthly"].Net),
				lib.HumanaizePriceEuro(policy.OffersPrices["default"]["monthly"].Tax),
				lib.HumanaizePriceEuro(policy.OffersPrices["default"]["monthly"].Gross),
			},
			{
				"Pari ad un premio Annuale",
				lib.HumanaizePriceEuro(policy.OffersPrices["default"]["monthly"].Net * 12),
				lib.HumanaizePriceEuro(policy.OffersPrices["default"]["monthly"].Tax * 12),
				lib.HumanaizePriceEuro(policy.OffersPrices["default"]["monthly"].Gross * 12),
			},
		}
	case string(models.PaySplitYear), string(models.PaySplitYearly):
		paymentSplit = "ANNUALE"
		tableInfo = [][]string{
			{
				"Annuale firma del contratto",
				lib.HumanaizePriceEuro(policy.OffersPrices["default"]["yearly"].Net),
				lib.HumanaizePriceEuro(policy.OffersPrices["default"]["yearly"].Tax),
				lib.HumanaizePriceEuro(policy.OffersPrices["default"]["yearly"].Gross),
			},
		}
	}

	getParagraphTitle(pdf, "Il premio per tutte le coperture assicurative attivate sulla polizza – Frazionamento: "+paymentSplit)
	setBlackRegularFont(pdf, standardTextSize)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(40, 2, "Premio", "", 0, "", false, 0, "")
	pdf.SetX(pdf.GetX() + 20)
	pdf.CellFormat(40, 2, "Imponibile", "", 0, "", false, 0, "")
	pdf.SetX(pdf.GetX() + 15)
	pdf.CellFormat(40, 2, "Imposte Assicurative", "", 0, "", false, 0, "")
	pdf.SetX(pdf.GetX() + 15)
	pdf.CellFormat(40, 2, "Totale", "", 0, "", false, 0, "")
	pdf.Ln(3)
	drawPinkHorizontalLine(pdf, thinLineWidth)
	pdf.Ln(1)
	for _, info := range tableInfo {
		pdf.CellFormat(40, 2, info[0], "", 0, "", false, 0,
			"")
		pdf.SetX(pdf.GetX() + 8)
		pdf.CellFormat(40, 2, info[1], "", 0,
			"CM", false, 0, "")
		pdf.SetX(pdf.GetX() + 20)
		pdf.CellFormat(40, 2, info[2], "", 0,
			"CM", false, 0, "")
		pdf.SetX(pdf.GetX() + 18)
		pdf.CellFormat(20, 2, info[3], "",
			0, "CM", false, 0, "")
		pdf.Ln(3)
		drawPinkHorizontalLine(pdf, thinLineWidth)
		pdf.Ln(1)
	}
	pdf.Ln(3)
}

func paymentResumeSection(pdf *fpdf.Fpdf, policy *models.Policy) {
	payments := make([]float64, 20)
	var paymentSplit string
	policyStartDate := policy.StartDate

	cellWidth := pdf.GetStringWidth("00/00/0000:") + pdf.GetStringWidth("€ ###.###,##")

	if policy.PaymentSplit == string(models.PaySplitYear) || policy.PaymentSplit == string(models.PaySplitYearly) {
		paymentSplit = "ANNUALE"
		for _, guarantee := range policy.Assets[0].Guarantees {
			for i := 0; i < guarantee.Value.Duration.Year; i++ {
				payments[i] += guarantee.Value.PremiumGrossYearly
			}
		}
	} else if policy.PaymentSplit == string(models.PaySplitMonthly) {
		paymentSplit = "MENSILE"
		for _, guarantee := range policy.Assets[0].Guarantees {
			for i := 0; i < guarantee.Value.Duration.Year; i++ {
				for y := 0; y < 12; y++ {
					payments[i] += guarantee.Value.PremiumGrossMonthly
				}
			}
		}
	}

	getParagraphTitle(pdf, "Pagamento dei premi successivi al primo")
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, "Il Contraente è tenuto a pagare i Premi entro 30 giorni dalle relative scadenze. "+
		"In caso di mancato pagamento del premio entro 30 giorni dalla scadenza (c.d. termine di tolleranza) "+
		"l’assicurazione è sospesa. Il contratto è risolto automaticamente in caso di mancato pagamento "+
		"del Premio entro 90 giorni dalla scadenza.", "", "", false)
	pdf.Ln(3)
	drawPinkHorizontalLine(pdf, 0.4)
	pdf.Ln(1)
	setBlackBoldFont(pdf, standardTextSize)
	pdf.CellFormat(pdf.GetStringWidth("Tipologia di premio"), 3, "Tipologia di premio:", "",
		0, "", false, 0, "")
	pdf.SetX(pdf.GetX() + 5)
	setBlackRegularFont(pdf, standardTextSize)
	setBlackDrawColor(pdf)
	pdf.CellFormat(3, 3, "", "1", 0, "", false, 0, "")
	pdf.SetX(pdf.GetX() + 1)
	pdf.CellFormat(pdf.GetStringWidth("naturale variabile annualmente"), 3, "naturale variabile annualmente",
		"", 0, "", false, 0, "")
	pdf.SetX(pdf.GetX() + 5)
	pdf.CellFormat(3, 3, "X", "1", 0, "CM", false, 0, "")
	pdf.SetX(pdf.GetX() + 1)
	pdf.CellFormat(pdf.GetStringWidth("fisso"), 3, "fisso", "", 0, "", false, 0,
		"")
	pdf.SetX(pdf.GetX() + 40)
	setBlackBoldFont(pdf, standardTextSize)
	pdf.CellFormat(pdf.GetStringWidth("Frazionamento"), 3, "Frazionamento:", "", 0, "",
		false, 0, "")
	pdf.SetX(pdf.GetX() + 3)
	pdf.CellFormat(pdf.GetStringWidth(paymentSplit), 3, paymentSplit, "", 0, "", false,
		0, "")
	pdf.Ln(4)
	drawPinkHorizontalLine(pdf, 0.1)
	pdf.Ln(1)
	setBlackRegularFont(pdf, standardTextSize)
	pdf.Cell(0, 3, "Il Premio è dovuto alle diverse annualità di Polizza, alle date qui sotto indicate:")
	pdf.Ln(4)
	drawPinkHorizontalLine(pdf, 0.1)

	for x := 0; x < len(payments)/4; x++ {
		pdf.Ln(1)
		for y := 0; y < 4; y++ {
			pdf.SetX(pdf.GetX() + 4)

			if payments[x+(5*y)] != 0 {
				var date string
				if x == 0 && y == 0 {
					date = "Alla firma:"
				} else {
					date = policyStartDate.AddDate(x+(5*y), 0, 0).Format(dateLayout) + ":"
				}
				price := lib.HumanaizePriceEuro(payments[x+(5*y)])

				pdf.CellFormat(cellWidth, 3, date+" "+price, "", 0, fpdf.AlignRight, false,
					0, "")

			} else {
				pdf.CellFormat(cellWidth, 3, "===========", "", 0, fpdf.AlignRight, false,
					0, "")
			}

		}
		pdf.Ln(4)
		drawPinkHorizontalLine(pdf, 0.1)
	}

	pdf.Ln(1)
	setBlackRegularFont(pdf, smallTextSize)
	pdf.MultiCell(0, 3, "In caso di frazionamento mensile i Premi sopra riportati sono dovuti, alle date "+
		"indicate e con successiva frequenza mensile, in misura di 1/12 per ogni mensilità. Non sono previsti oneri "+
		"o interessi di frazionamento.", "", "", false)
	pdf.Ln(3)
}
