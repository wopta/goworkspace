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

	lifeGuaranteesTable(pdf, guaranteesMap, slugs)

	avvertenzeBeneficiariSection(pdf)

	beneficiaries, legitimateSuccessorsChoice, designatedSuccessorsChoice := loadLifeBeneficiariesInfo(policy)

	beneficiariesSection(pdf, beneficiaries, legitimateSuccessorsChoice, designatedSuccessorsChoice)

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
