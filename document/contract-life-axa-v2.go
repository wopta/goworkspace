package document

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/go-pdf/fpdf"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"strings"
	"time"
)

func lifeAxaContractV2(pdf *fpdf.Fpdf, origin string, policy *models.Policy, networkNode *models.NetworkNode, product *models.Product) (string, []byte) {
	signatureID = 0

	lifeMainHeaderV2(pdf, policy, networkNode, false)

	mainFooter(pdf, policy.Name)

	pdf.AddPage()

	lifeInsuredInfoSectionV2(pdf, policy, false)

	guaranteesMap, slugs := loadLifeGuarantees(policy, product)

	lifeGuaranteesTableV2(pdf, guaranteesMap, slugs)

	lifeAvvertenzeBeneficiariSectionV2(pdf)

	beneficiaries, legitimateSuccessorsChoice, designatedSuccessorsChoice := loadLifeBeneficiariesInfo(policy)

	lifeBeneficiariesSectionV2(pdf, beneficiaries, legitimateSuccessorsChoice, designatedSuccessorsChoice)

	lifeBeneficiaryReferenceSectionV2(pdf, policy)

	//pdf.AddPage()

	lifeSurveysSectionV2(pdf, policy, false)

	pdf.AddPage()

	lifeStatementsSectionV2(pdf, policy, false)

	lifeOfferResumeSectionV2(pdf, policy)

	lifePaymentResumeSectionV2(pdf, policy)

	lifeContractWithdrawlSectionV2(pdf, false)

	pdf.AddPage()

	lifePaymentMethodSectionV2(pdf)

	lifeEmitResumeSectionV2(pdf, policy)

	companiesDescriptionSection(pdf, policy.Company)

	axaHeader(pdf, false)

	pdf.AddPage()

	axaFooter(pdf)

	axaDeclarationsConsentSection(pdf, policy, false)

	pdf.AddPage()

	axaTableSection(pdf, policy)

	pdf.AddPage()

	axaTablePart2Section(pdf, policy, false)

	pdf.Ln(15)

	axaTablePart3Section(pdf)

	woptaHeader(pdf, false)

	//woptaFooter(pdf)

	generatePolicyAnnex(pdf, origin, networkNode)

	pdf.AddPage()

	woptaFooter(pdf)

	woptaPrivacySection(pdf)

	personalDataHandlingSection(pdf, policy, false)

	filename, out := saveContract(pdf, policy)
	return filename, out
}

func lifeMainHeaderV2(pdf *fpdf.Fpdf, policy *models.Policy, networkNode *models.NetworkNode, isProposal bool) {
	var (
		opt                                                             fpdf.ImageOptions
		logoPath, policyInfoHeader, policyInfo, expiryInfo, productName string
	)

	location, err := time.LoadLocation("Europe/Rome")
	lib.CheckError(err)

	policyStartDate := policy.StartDate.In(location)
	policyEndDate := policy.EndDate.In(location)

	if policy.PaymentSplit == string(models.PaySplitMonthly) {
		expiryInfo = "Prima scandenza mensile il: " +
			policyStartDate.AddDate(0, 1, 0).Format(dateLayout) + "\n"
	} else if policy.PaymentSplit == string(models.PaySplitYear) || policy.PaymentSplit == string(models.PaySplitYearly) {
		expiryInfo = "Prima scadenza annuale il: " +
			policyStartDate.AddDate(1, 0, 0).Format(dateLayout) + "\n"
	}

	if isProposal {
		policyInfoHeader = "I dati della tua proposta"
		policyInfo = fmt.Sprintf("Numero: %d\n", policy.ProposalNumber)
	} else {
		policyInfoHeader = "I dati della tua polizza"
		policyInfo = fmt.Sprintf("Numero: %s\n", policy.CodeCompany)
	}

	policyInfo += "Decorre dal: " + policyStartDate.Format(dateLayout) + " ore 24:00\n" +
		"Scade il: " + policyEndDate.In(location).Format(dateLayout) + " ore 24:00\n" +
		expiryInfo +
		"Non si rinnova a scadenza.\n"

	if networkNode != nil {
		policyInfo += "Produttore: " + getProducerName(networkNode)
	}

	logoPath = lib.GetAssetPathByEnvV2() + "logo_vita.png"
	productName = "Vita"

	contractor := policy.Contractor
	address := strings.ToUpper(contractor.Residence.StreetName + ", " + contractor.Residence.StreetNumber + "\n" +
		contractor.Residence.PostalCode + " " + contractor.Residence.City + " (" + contractor.Residence.CityCode + ")\n")

	contractorInfo := "Contraente: " + strings.ToUpper(contractor.Surname+" "+contractor.Name+"\n"+
		"C.F./P.IVA: "+contractor.FiscalCode) + "\n" +
		"Indirizzo: " + strings.ToUpper(address) + "Mail: " + contractor.Mail + "\n" +
		"Telefono: " + contractor.Phone

	pdf.SetHeaderFunc(func() {
		opt.ImageType = "png"
		pdf.ImageOptions(logoPath, 10.5, 6, 13, 13, false, opt, 0, "")
		pdf.SetXY(23.5, 7)
		setPinkBoldFont(pdf, 18)
		pdf.Cell(10, 6, "Wopta per te")
		setPinkItalicFont(pdf, 18)
		pdf.SetXY(23.5, 13)
		pdf.SetTextColor(92, 89, 92)
		pdf.Cell(10, 6, productName)
		pdf.ImageOptions(lib.GetAssetPathByEnvV2()+"logo_wopta.png", 170, 6, 0, 8, false, opt, 0, "")

		setBlackBoldFont(pdf, standardTextSize)
		pdf.SetXY(10, 20)
		pdf.Cell(0, 3, policyInfoHeader)
		setBlackRegularFont(pdf, standardTextSize)
		pdf.SetXY(10, pdf.GetY()+3)
		pdf.MultiCell(0, 3.5, policyInfo, "", "", false)

		setBlackBoldFont(pdf, standardTextSize)
		pdf.SetXY(-95, 20)
		pdf.Cell(0, 3, "I tuoi dati")
		setBlackRegularFont(pdf, standardTextSize)
		pdf.SetXY(-95, pdf.GetY()+3)
		pdf.MultiCell(0, 3.5, contractorInfo, "", "", false)
		pdf.Ln(5)

		if isProposal {
			insertWatermark(pdf, proposal)
		}
	})
}

func lifeInsuredInfoSectionV2(pdf *fpdf.Fpdf, policy *models.Policy, isProposal bool) {
	title := "La tua assicurazione è operante per il seguente Assicurato e Garanzie"

	if isProposal {
		title = "La tua assicurazione sarà operante per il seguente Assicurato e Garanzie"
	}
	getParagraphTitle(pdf, title)
	lifeInsuredInfoTableV2(pdf, policy.Assets[0].Person)
}

func lifeInsuredInfoTableV2(pdf *fpdf.Fpdf, insured *models.User) {
	var (
		residenceAddress, domicileAddress, birthDate string
	)

	residenceAddress = strings.ToUpper(insured.Residence.StreetName + ", " + insured.Residence.StreetNumber +
		" - " + insured.Residence.PostalCode + " " + insured.Residence.City + " (" + insured.Residence.CityCode + ")")

	if insured.Domicile != nil {
		domicileAddress = strings.ToUpper(insured.Domicile.StreetName + ", " + insured.Domicile.StreetNumber +
			" - " + insured.Domicile.PostalCode + " " + insured.Domicile.City + " (" + insured.Domicile.CityCode + ")")
	} else {
		domicileAddress = residenceAddress
	}

	tmpBirthDate, err := time.Parse(time.RFC3339, insured.BirthDate)
	if err == nil {
		birthDate = tmpBirthDate.Format(dateLayout)
	}

	tableRows := [][]string{
		{"Cognome e Nome", strings.ToUpper(insured.Surname + " " + insured.Name), "Codice fiscale:", insured.FiscalCode},
		{"Residente in", residenceAddress, "Data nascita:", birthDate},
		{"Domicilio", domicileAddress, "", ""},
		{"Mail", insured.Mail, "Telefono:", formatPhoneNumber(insured.Phone)},
	}

	drawPinkHorizontalLine(pdf, thickLineWidth)
	for index, row := range tableRows {
		setBlackBoldFont(pdf, standardTextSize)
		pdf.CellFormat(32, 5, row[0], "", 0, fpdf.AlignLeft+fpdf.AlignMiddle, false, 0, "")
		setBlackRegularFont(pdf, standardTextSize)
		pdf.CellFormat(98, 5, row[1], "", 0, fpdf.AlignLeft+fpdf.AlignMiddle, false, 0, "")
		setBlackBoldFont(pdf, standardTextSize)
		pdf.CellFormat(26, 5, row[2], "", 0, fpdf.AlignLeft+fpdf.AlignMiddle, false, 0, "")
		setBlackRegularFont(pdf, standardTextSize)
		pdf.CellFormat(30, 5, row[3], "", 1, fpdf.AlignLeft+fpdf.AlignMiddle, false, 0, "")

		if index == len(tableRows)-1 {
			drawPinkHorizontalLine(pdf, thickLineWidth)
			pdf.Ln(2)
			break
		}
		drawPinkHorizontalLine(pdf, thinLineWidth)
	}
}

func lifeGuaranteesTableV2(pdf *fpdf.Fpdf, guaranteesMap map[string]map[string]string, slugs []slugStruct) {
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
	pdf.Ln(1)

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
			1, "RM", false, 0, "")
		//pdf.Ln(5)
		drawPinkHorizontalLine(pdf, thinLineWidth)
	}
	pdf.Ln(0.5)
	setBlackRegularFont(pdf, smallTextSize)
	pdf.Cell(80, 3, "(*) imposte assicurative di legge incluse nella misura del 2,50% del premio imponibile")
	pdf.Ln(5)
}

func lifeAvvertenzeBeneficiariSectionV2(pdf *fpdf.Fpdf) {
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
	pdf.Ln(2)
}

func lifeBeneficiariesSectionV2(pdf *fpdf.Fpdf, beneficiaries []map[string]string, legitimateSuccessorsChoice,
	designatedSuccessorsChoice string) {
	getParagraphTitle(pdf, "Beneficiario")
	setBlackRegularFont(pdf, standardTextSize)
	pdf.CellFormat(0, 4, "Io sottoscritto Assicurato, con la sottoscrizione della presente polizza, in "+
		"riferimento alla garanzia Decesso:", "", 1, "", false, 0, "")

	rows := [][]string{
		{legitimateSuccessorsChoice, "Designo genericamente quali beneficiari della prestazione i miei eredi " +
			"(legittimi e/o testamentari)"},
		{designatedSuccessorsChoice, "Designo nominativamente il/i seguente/i soggetto/i quale beneficiario/i della " +
			"prestazione"},
	}

	setBlackDrawColor(pdf)
	for _, row := range rows {
		pdf.SetX(11.4)
		pdf.CellFormat(3, 3, row[0], "1", 0, "CM", false, 0, "")
		pdf.CellFormat(0, 3.5, row[1], "", 1, "LM", false, 0, "")
	}
	pdf.Ln(1)

	lifeBeneficiariesTableV2(pdf, beneficiaries)
}

func lifeBeneficiariesTableV2(pdf *fpdf.Fpdf, beneficiaries []map[string]string) {
	tables := make([][][]string, 0)

	for _, beneficiary := range beneficiaries {
		tableRows := [][]string{
			{"Cognome e nome", beneficiary["name"], "Cod. Fisc.:", beneficiary["fiscalCode"]},
			{"Indirizzo", beneficiary["address"], "", ""},
			{"Mail", beneficiary["mail"], "Telefono:", beneficiary["phone"]},
			{"Relazione con assicurato", beneficiary["relation"], "", ""},
		}
		tables = append(tables, tableRows)
	}

	for tableIndex, table := range tables {
		drawPinkHorizontalLine(pdf, thickLineWidth)

		for _, row := range table {
			setBlackBoldFont(pdf, standardTextSize)
			pdf.CellFormat(45, 5, row[0], "", 0, fpdf.AlignLeft, false, 0, "")
			setBlackRegularFont(pdf, standardTextSize)
			pdf.CellFormat(80, 5, row[1], "", 0, fpdf.AlignLeft, false, 0, "")
			setBlackBoldFont(pdf, standardTextSize)
			pdf.CellFormat(20, 5, row[2], "", 0, fpdf.AlignLeft, false, 0, "")
			setBlackRegularFont(pdf, standardTextSize)
			pdf.CellFormat(45, 5, row[3], "", 1, fpdf.AlignLeft, false, 0, "")
			drawPinkHorizontalLine(pdf, thinLineWidth)
		}

		pdf.CellFormat(165, 5, "Consenso ad invio comunicazioni da parte della Compagnia al beneficiario, prima "+
			"dell'evento Decesso:", "", 0, fpdf.AlignLeft, false, 0, "")
		pdf.CellFormat(50, 5, beneficiaries[tableIndex]["contactConsent"], "", 1, fpdf.AlignLeft, false, 0, "")
		drawPinkHorizontalLine(pdf, thinLineWidth)
		pdf.Ln(2)
	}
}

func lifeBeneficiaryReferenceSectionV2(pdf *fpdf.Fpdf, policy *models.Policy) {
	beneficiaryReference := map[string]string{
		"name":       "=====",
		"fiscalCode": "=====",
		"address":    "=====",
		"mail":       "=====",
		"phone":      "=====",
	}

	deathGuarantee, err := policy.ExtractGuarantee("death")

	if err == nil && deathGuarantee.BeneficiaryReference != nil {
		beneficiary := deathGuarantee.BeneficiaryReference
		address := strings.ToUpper(beneficiary.Residence.StreetName + ", " + beneficiary.Residence.StreetNumber +
			" - " + beneficiary.Residence.PostalCode + " " + beneficiary.Residence.City +
			" (" + beneficiary.Residence.CityCode + ")")
		beneficiaryReference["name"] = strings.ToUpper(beneficiary.Surname + " " + beneficiary.Name)
		beneficiaryReference["fiscalCode"] = strings.ToUpper(beneficiary.FiscalCode)
		beneficiaryReference["address"] = address
		beneficiaryReference["mail"] = beneficiary.Mail
		beneficiaryReference["phone"] = beneficiary.Phone
	}

	getParagraphTitle(pdf, "Referente terzo")
	lifeBeneficiaryReferenceTableV2(pdf, beneficiaryReference)
	pdf.Ln(2)
}

func lifeBeneficiaryReferenceTableV2(pdf *fpdf.Fpdf, beneficiaryReference map[string]string) {
	tableRows := [][]string{
		{"Cognome e nome", beneficiaryReference["name"], "Cod, Fisc.:", beneficiaryReference["fiscalCode"]},
		{"Indirizzo", beneficiaryReference["address"], "", ""},
		{"Mail", beneficiaryReference["mail"], "Telefono:", beneficiaryReference["phone"]},
	}

	drawPinkHorizontalLine(pdf, thickLineWidth)
	for _, row := range tableRows {
		setBlackBoldFont(pdf, standardTextSize)
		pdf.CellFormat(45, 5, row[0], "", 0, fpdf.AlignLeft, false, 0, "")
		setBlackRegularFont(pdf, standardTextSize)
		pdf.CellFormat(80, 5, row[1], "", 0, fpdf.AlignLeft, false, 0, "")
		setBlackBoldFont(pdf, standardTextSize)
		pdf.CellFormat(20, 5, row[2], "", 0, fpdf.AlignLeft, false, 0, "")
		setBlackRegularFont(pdf, standardTextSize)
		pdf.CellFormat(45, 5, row[3], "", 1, fpdf.AlignLeft, false, 0, "")
		drawPinkHorizontalLine(pdf, thinLineWidth)
	}
}

func lifeSurveysSectionV2(pdf *fpdf.Fpdf, policy *models.Policy, isProposal bool) {
	surveys := *policy.Surveys

	getParagraphTitle(pdf, "Dichiarazioni da leggere con attenzione prima di firmare")
	err := printSurvey(pdf, surveys[0], policy.Company, isProposal)
	lib.CheckError(err)

	pdf.AddPage()

	getParagraphTitle(pdf, "Questionario Medico")
	for _, survey := range surveys[1:] {
		err = printSurvey(pdf, survey, policy.Company, isProposal)
		lib.CheckError(err)
	}
}

func lifeStatementsSectionV2(pdf *fpdf.Fpdf, policy *models.Policy, isProposal bool) {
	statements := *policy.Statements
	for _, statement := range statements {
		printStatement(pdf, statement, policy.Company, isProposal)
	}
}

func lifeOfferResumeSectionV2(pdf *fpdf.Fpdf, policy *models.Policy) {
	var (
		paymentSplit string
		paymentInfo  [][]string
		tableRows    = [][]string{
			{"Premio", "Imponibile", "Imposte Assicurative", "Totale"},
		}
	)

	switch policy.PaymentSplit {
	case string(models.PaySplitMonthly):
		paymentSplit = "MENSILE"
		paymentInfo = [][]string{
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
		paymentInfo = [][]string{
			{
				"Annuale firma del contratto",
				lib.HumanaizePriceEuro(policy.OffersPrices["default"]["yearly"].Net),
				lib.HumanaizePriceEuro(policy.OffersPrices["default"]["yearly"].Tax),
				lib.HumanaizePriceEuro(policy.OffersPrices["default"]["yearly"].Gross),
			},
		}
	}
	tableRows = append(tableRows, paymentInfo...)

	getParagraphTitle(pdf, "Il premio per tutte le coperture assicurative attivate sulla polizza – Frazionamento: "+paymentSplit)

	for _, row := range tableRows {
		setBlackRegularFont(pdf, standardTextSize)
		pdf.CellFormat(70, 5, row[0], "", 0, fpdf.AlignLeft, false, 0, "")
		pdf.CellFormat(50, 5, row[1], "", 0, fpdf.AlignLeft, false, 0, "")
		pdf.CellFormat(40, 5, row[2], "", 0, fpdf.AlignLeft, false, 0, "")
		pdf.CellFormat(30, 5, row[3], "", 1, fpdf.AlignRight, false, 0, "")
		drawPinkHorizontalLine(pdf, thinLineWidth)
	}
	pdf.Ln(3)
}

func lifePaymentResumeSectionV2(pdf *fpdf.Fpdf, policy *models.Policy) {
	var (
		paymentSplit string
		payments     = make([]float64, 20)
	)

	policyStartDate := policy.StartDate

	cellWidth := pdf.GetStringWidth("00/00/0000:") + pdf.GetStringWidth("€ ###.###,##")

	switch policy.PaymentSplit {
	case string(models.PaySplitYear), string(models.PaySplitYearly):
		paymentSplit = "ANNUALE"
		for _, guarantee := range policy.Assets[0].Guarantees {
			for i := 0; i < guarantee.Value.Duration.Year; i++ {
				payments[i] += guarantee.Value.PremiumGrossYearly
			}
		}
	case string(models.PaySplitMonthly):
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
	drawPinkHorizontalLine(pdf, thickLineWidth)
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

func lifeContractWithdrawlSectionV2(pdf *fpdf.Fpdf, isProposal bool) {
	getParagraphTitle(pdf, "Informativa sul diritto di recesso")

	paragraphs := [][]string{
		{"Diritto di recesso entro i primi 30 giorni dalla stipula (diritto di ripensamento)", "Il Contraente può recedere dal contratto entro il termine di 30 giorni dalla " +
			"decorrenza dell’assicurazione (diritto di ripensamento). In tal caso, l’assicurazione si intende come mai " +
			"entrata in vigore e la Compagnia, per il tramite dell’intermediario, provvederà a rimborsare al Contraente " +
			"l’importo di Premio già versato (al netto delle imposte)."},
		{"Diritto di recesso annuale (disdetta alla annualità)", "Il Contraente può recedere dal contratto annualmente, entro il termine di 30 " +
			"giorni dalla scadenza annuale della polizza (disdetta alla annualità). In tal caso, l’assicurazione cessa alle " +
			"ore 24:00 dell’ultimo giorno della annualità in corso. È possibile disdettare singolarmente una o più delle " +
			"coperture attivate in fase di sottoscrizione."},
		{"Modalità per l’esercizio del diritto di recesso", "Il Contraente è tenuto ad esercitare il diritto di recesso mediante invio di una " +
			"lettera raccomandata a.r. al seguente indirizzo: Wopta Assicurazioni srl – Gestione Portafoglio – Galleria del " +
			"Corso, 1 – 201212 Milano (MI) oppure via posta elettronica certificata (PEC) all’indirizzo " +
			"email: woptaassicurazioni@legalmail.it"},
	}

	for _, paragraph := range paragraphs {
		setBlackBoldFont(pdf, standardTextSize)
		pdf.MultiCell(0, 3, paragraph[0], "", "", false)
		setBlackRegularFont(pdf, standardTextSize)
		pdf.MultiCell(0, 3, paragraph[1], "", "", false)
	}

	if !isProposal {
		pdf.Ln(5)
		drawSignatureForm(pdf)
		pdf.Ln(5)
	}
}

func lifePaymentMethodSectionV2(pdf *fpdf.Fpdf) {
	getParagraphTitle(pdf, "Come puoi pagare il premio")
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, "I mezzi di pagamento consentiti, nei confronti di Wopta, sono esclusivamente "+
		"bonifico e strumenti di pagamento elettronico, quali ad esempio, carte di credito e/o carte di debito, "+
		"incluse le carte prepagate. Oppure può essere pagato direttamente alla Compagnia alla "+
		"stipula del contratto, via bonifico o carta di credito.", "", "", false)
	pdf.Ln(3)
}

func lifeEmitResumeSectionV2(pdf *fpdf.Fpdf, policy *models.Policy) {
	var offerPrice string
	emitDate := time.Now().UTC().Format(dateLayout)
	if policy.PaymentSplit == string(models.PaySplitMonthly) {
		offerPrice = humanize.FormatFloat("#.###,##", policy.PriceGrossMonthly)
	} else {
		offerPrice = humanize.FormatFloat("#.###,##", policy.PriceGross)
	}
	text := "Polizza emessa a Milano il " + emitDate + " per un importo di € " + offerPrice + " quale " +
		"prima rata alla firma, il cui pagamento a saldo è da effettuarsi con i metodi di pagamento sopra indicati."

	getParagraphTitle(pdf, "Emissione polizza e pagamento della prima rata")
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, text, "", "", false)
	pdf.Ln(3)
}
