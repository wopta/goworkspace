package document

import (
	"fmt"
	"github.com/go-pdf/fpdf"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"strings"
	"time"
)

func lifeAxaV2(pdf *fpdf.Fpdf, origin string, policy *models.Policy, networkNode *models.NetworkNode, product *models.Product) (string, []byte) {
	signatureID = 0

	lifeMainHeaderV2(pdf, policy, networkNode)

	mainFooter(pdf, policy.Name)

	pdf.AddPage()

	lifeInsuredInfoSectionV2(pdf, policy)

	guaranteesMap, slugs := loadLifeGuarantees(policy, product)

	lifeGuaranteesTableV2(pdf, guaranteesMap, slugs)

	lifeAvvertenzeBeneficiariSectionV2(pdf)

	beneficiaries, legitimateSuccessorsChoice, designatedSuccessorsChoice := loadLifeBeneficiariesInfo(policy)

	lifeBeneficiariesSectionV2(pdf, beneficiaries, legitimateSuccessorsChoice, designatedSuccessorsChoice)

	beneficiaryReferenceSection(pdf, policy)

	surveysSection(pdf, policy)

	pdf.AddPage()

	statementsSection(pdf, policy)

	offerResumeSection(pdf, policy)

	paymentResumeSection(pdf, policy)

	contractWithdrawlSection(pdf)

	pdf.AddPage()

	paymentMethodSection(pdf)

	emitResumeSection(pdf, policy)

	companiesDescriptionSection(pdf, policy.Company)

	axaHeader(pdf)

	pdf.AddPage()

	axaFooter(pdf)

	axaDeclarationsConsentSection(pdf, policy)

	pdf.AddPage()

	axaTableSection(pdf, policy)

	pdf.AddPage()

	axaTablePart2Section(pdf, policy)

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

	personalDataHandlingSection(pdf, policy)

	filename, out := saveContract(pdf, policy)
	return filename, out
}

func lifeMainHeaderV2(pdf *fpdf.Fpdf, policy *models.Policy, networkNode *models.NetworkNode) {
	var (
		opt                                     fpdf.ImageOptions
		logoPath, cfpi, expiryInfo, productName string
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

	policyInfo := "Numero: " + policy.CodeCompany + "\n" +
		"Decorre dal: " + policyStartDate.Format(dateLayout) + " ore 24:00\n" +
		"Scade il: " + policyEndDate.In(location).Format(dateLayout) + " ore 24:00\n" +
		expiryInfo +
		"Non si rinnova a scadenza.\n"

	logoPath = lib.GetAssetPathByEnvV2() + "logo_vita.png"
	productName = "Vita"

	if networkNode != nil {
		policyInfo += "Produttore: "
		switch networkNode.Type {
		case models.AgentNetworkNodeType:
			policyInfo += strings.ToUpper(fmt.Sprintf("%s %s\n", networkNode.Agent.Surname, networkNode.Agent.Name))
		case models.AgencyNetworkNodeType:
			policyInfo += strings.ToUpper(fmt.Sprintf("%s\n", networkNode.Agency.Name))
		}
	}

	contractor := policy.Contractor
	address := strings.ToUpper(contractor.Residence.StreetName + ", " + contractor.Residence.StreetNumber + "\n" +
		contractor.Residence.PostalCode + " " + contractor.Residence.City + " (" + contractor.Residence.CityCode + ")\n")

	if contractor.VatCode == "" {
		cfpi = contractor.FiscalCode
	} else {
		cfpi = contractor.VatCode
	}

	contractorInfo := "Contraente: " + strings.ToUpper(contractor.Surname+" "+contractor.Name+"\n"+
		"C.F./P.IVA: "+cfpi) + "\n" +
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
		pdf.Cell(0, 3, "I dati della tua polizza")
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
	})
}

func lifeInsuredInfoSectionV2(pdf *fpdf.Fpdf, policy *models.Policy) {
	getParagraphTitle(pdf, "La tua assicurazione è operante per il seguente Assicurato e Garanzie")
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
	if err != nil {
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
			{"Cognome e nome", beneficiary["name"], "Cod. Fisc.:", beneficiary["fiscCode"]},
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
