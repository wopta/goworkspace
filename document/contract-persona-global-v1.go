package document

import (
	"fmt"
	"github.com/go-pdf/fpdf"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"strings"
	"time"
)

type keyValue struct {
	key   string
	value string
}

func personaGlobalContractV1(pdf *fpdf.Fpdf, policy *models.Policy, networkNode *models.NetworkNode, product *models.Product) (string, []byte) {
	signatureID = 0

	personaMainHeaderV1(pdf, policy, networkNode, false)

	mainFooter(pdf, policy.Name)

	pdf.AddPage()

	personaInsuredInfoSection(pdf, policy)

	guaranteesMap, slugs := loadPersonaGuarantees(policy, product)

	personaGuaranteesTable(pdf, guaranteesMap, slugs)

	pdf.Ln(5)

	personaSurveySection(pdf, policy, false)

	personaStatementsSection(pdf, policy, false)

	if policy.HasGuarantee("IPM") {
		pdf.AddPage()
	}

	personaOfferResumeSection(pdf, policy)

	paymentMethodSection(pdf)

	emitResumeSection(pdf, policy)

	companiesDescriptionSection(pdf, policy.Company)

	personalDataHandlingSection(pdf, policy, false)

	filename, out := saveContract(pdf, policy)
	return filename, out
}

func personaMainHeaderV1(pdf *fpdf.Fpdf, policy *models.Policy, networkNode *models.NetworkNode, isProposal bool) {
	var (
		opt                                                       fpdf.ImageOptions
		logoPath, cfpi, policyInfoHeader, policyInfo, productName string
	)

	location, err := time.LoadLocation("Europe/Rome")
	lib.CheckError(err)

	policyStartDate := policy.StartDate.In(location)
	policyEndDate := policy.EndDate.In(location)

	if isProposal {
		policyInfoHeader = "I dati della tua proposta"
		policyInfo = fmt.Sprintf("Numero: %d\n", policy.ProposalNumber)
	} else {
		policyInfoHeader = "I dati della tua polizza"
		policyInfo = fmt.Sprintf("Numero: %s\n", policy.CodeCompany)
	}

	policyInfo += "Decorre dal: " + policyStartDate.Format(dateLayout) + " ore 24:00\n" +
		"Scade il: " + policyEndDate.In(location).Format(dateLayout) + " ore 24:00\n"

	logoPath = lib.GetAssetPathByEnvV2() + "logo_persona.png"
	productName = "Persona"
	policyInfo += "Si rinnova a scadenza salvo disdetta da inviare 30 giorni prima\n" + "Prossimo pagamento "
	if policy.PaymentSplit == string(models.PaySplitMonthly) {
		policyInfo += policyStartDate.In(location).AddDate(0, 1, 0).Format(dateLayout) + "\n"
	} else if policy.PaymentSplit == string(models.PaySplitYear) {
		policyInfo += policyStartDate.In(location).AddDate(1, 0, 0).Format(dateLayout) + "\n"
	}
	policyInfo += "Sostituisce la polizza ========\n"

	if networkNode != nil {
		policyInfo += "Produttore: " + getProducerName(networkNode)
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
		pdf.ImageOptions(logoPath, 10, 6, 13, 13, false, opt, 0, "")
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
		pdf.SetXY(-75, 20)
		pdf.Cell(0, 3, "I tuoi dati")
		setBlackRegularFont(pdf, standardTextSize)
		pdf.SetXY(-75, pdf.GetY()+3)
		pdf.MultiCell(0, 3.5, contractorInfo, "", "", false)
		pdf.Ln(5)
	})
}

func personaInsuredInfoSection(pdf *fpdf.Fpdf, policy *models.Policy) {
	coverageTypeMap := map[string]string{
		"24h":   "Professionale ed Extraprofessionale",
		"prof":  "Professionale",
		"extra": "Extraprofessionale",
	}

	getParagraphTitle(pdf, "La tua assicurazione è operante per il seguente Assicurato e Garanzie")
	drawPinkHorizontalLine(pdf, thickLineWidth)
	pdf.Ln(2)
	contractorInfo := []keyValue{
		{key: "Assicurato: ", value: "1"},
		{key: "Cognome e Nome: ", value: policy.Contractor.Surname + " " + policy.Contractor.Name},
		{key: "Codice Fiscale: ", value: policy.Contractor.FiscalCode},
		{key: "Professione: ", value: policy.Contractor.Work},
		{key: "Tipo Professione: ", value: strings.ToUpper(policy.Contractor.WorkType[:1]) + policy.Contractor.WorkType[1:]},
		{key: "Classe rischio: ", value: "Classe " + policy.Contractor.RiskClass},
		{key: "Forma di copertura: ", value: coverageTypeMap[policy.Assets[0].Guarantees[0].Type]},
	}

	maxLength := 0
	for _, info := range contractorInfo {
		if len(info.key) > maxLength {
			maxLength = len(info.key)
		}
	}

	for _, info := range contractorInfo {
		setBlackBoldFont(pdf, standardTextSize)
		pdf.CellFormat(40, 4, info.key, "B", 0, fpdf.AlignRight, false, 0, "")
		setBlackRegularFont(pdf, standardTextSize)
		pdf.CellFormat(2.5, 4, "", "", 0, "", false, 0, "")
		pdf.CellFormat(0, 4, info.value, "", 2, fpdf.AlignLeft, false, 0, "")
		pdf.Ln(1)
	}
}

func personaGuaranteesTable(pdf *fpdf.Fpdf, guaranteesMap map[string]map[string]string,
	slugs []slugStruct) {
	var table [][]string

	for _, slug := range slugs {
		r := []string{guaranteesMap[slug.name]["name"], guaranteesMap[slug.name]["sumInsuredLimitOfIndemnity"],
			guaranteesMap[slug.name]["details"], guaranteesMap[slug.name]["price"]}
		table = append(table, r)
	}

	setBlackBoldFont(pdf, titleTextSize)
	pdf.CellFormat(80, titleTextSize, "Garanzie", "B", 0, fpdf.AlignCenter, false, 0, "")
	pdf.CellFormat(30, titleTextSize, "Somma Assicurata", "B", 0, fpdf.AlignCenter, false, 0, "")
	pdf.CellFormat(5, titleTextSize, "", "B", 0, fpdf.AlignCenter, false, 0, "")
	pdf.CellFormat(60, titleTextSize, "Opzioni/Dettagli", "B", 0, fpdf.AlignLeft, false, 0, "")
	pdf.CellFormat(15, titleTextSize, "Premio", "B", 1, fpdf.AlignRight, false, 0, "")

	for _, slug := range slugs {
		numLines := float64(len(pdf.SplitText(guaranteesMap[slug.name]["details"], 60)))

		setBlackBoldFont(pdf, standardTextSize)
		pdf.CellFormat(80, 6*numLines, guaranteesMap[slug.name]["name"], "", 0, fpdf.AlignMiddle+fpdf.AlignLeft, false, 0, "")
		setBlackRegularFont(pdf, standardTextSize)
		pdf.CellFormat(30, 6*numLines, guaranteesMap[slug.name]["sumInsuredLimitOfIndemnity"], "", 0, fpdf.AlignMiddle+fpdf.AlignRight, false, 0, "")
		pdf.CellFormat(5, 6*numLines, "", "", 0, fpdf.AlignMiddle+fpdf.AlignRight, false, 0, "")
		if numLines > 1 {
			pdf.MultiCell(60, 6, guaranteesMap[slug.name]["details"], "", fpdf.AlignMiddle+fpdf.AlignLeft, false)
			pdf.SetXY(pdf.GetX()+175, pdf.GetY()-6*numLines)
			pdf.CellFormat(15, 6*numLines, guaranteesMap[slug.name]["price"], "", 1, fpdf.AlignMiddle+fpdf.AlignRight, false, 0, "")
		} else {
			pdf.CellFormat(60, 6, guaranteesMap[slug.name]["details"], "", 0, fpdf.AlignMiddle+fpdf.AlignLeft, false, 0, "")
			pdf.CellFormat(15, 6, guaranteesMap[slug.name]["price"], "", 1, fpdf.AlignMiddle+fpdf.AlignRight, false, 0, "")
		}
		drawPinkHorizontalLine(pdf, thinLineWidth)
	}
}

func personaSurveySection(pdf *fpdf.Fpdf, policy *models.Policy, isProposal bool) {
	surveys := *policy.Surveys

	getParagraphTitle(pdf, "Dichiarazioni da leggere con attenzione prima di firmare")
	err := printSurvey(pdf, surveys[0], policy.Company, isProposal)
	lib.CheckError(err)

	for _, survey := range surveys[1:] {
		err := printSurvey(pdf, survey, policy.Company, isProposal)
		lib.CheckError(err)
	}
	pdf.Ln(3)
}

func personaStatementsSection(pdf *fpdf.Fpdf, policy *models.Policy, isProposal bool) {
	statements := *policy.Statements

	for _, statement := range statements {
		printStatement(pdf, statement, policy.Company, isProposal)
	}
	pdf.Ln(3)
}

func personaOfferResumeSection(pdf *fpdf.Fpdf, policy *models.Policy) {
	var (
		tableInfo [][]string
	)

	switch policy.PaymentSplit {
	case string(models.PaySplitMonthly):
		tableInfo = [][]string{
			{
				"Annuale",
				lib.HumanaizePriceEuro(policy.PriceNettMonthly * 12),
				lib.HumanaizePriceEuro(policy.PriceGrossMonthly - policy.PriceNettMonthly*12),
				lib.HumanaizePriceEuro(policy.PriceGrossMonthly * 12),
			},
			{
				"Rata firma della polizza",
				lib.HumanaizePriceEuro(policy.PriceNettMonthly),
				lib.HumanaizePriceEuro(policy.PriceGrossMonthly - policy.PriceNettMonthly),
				lib.HumanaizePriceEuro(policy.PriceGrossMonthly),
			},
			{
				"Rata mensile",
				lib.HumanaizePriceEuro(policy.PriceNettMonthly),
				lib.HumanaizePriceEuro(policy.PriceGrossMonthly - policy.PriceNettMonthly),
				lib.HumanaizePriceEuro(policy.PriceGrossMonthly),
			},
		}
	case string(models.PaySplitYear):
		tableInfo = [][]string{
			{
				"Annuale firma della polizza",
				lib.HumanaizePriceEuro(policy.PriceNett),
				lib.HumanaizePriceEuro(policy.PriceGross - policy.PriceNett),
				lib.HumanaizePriceEuro(policy.PriceGross),
			},
		}
	}

	getParagraphTitle(pdf, "Il premio per tutte le coperture assicurative attivate sulla polizza")
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
	setBlackRegularFont(pdf, smallTextSize)
	pdf.MultiCell(0, 3, "In caso di sostituzione, il premio alla firma è al netto dell’eventuale rimborso"+
		" dei premi non goduti sulla polizza sostituita e tiene conto dell’eventuale diversa durata rispetto alle"+
		" rate successive.", "", fpdf.AlignLeft, false)
	pdf.Ln(3)
}

func globalStamentsAndConsens(pdf *fpdf.Fpdf) {
	setBlackBoldFont(pdf, titleTextSize)
	pdf.MultiCell(0, 3, "DICHIARAZIONI E CONSENSI", "", fpdf.AlignLeft, false)
	pdf.Ln(3)
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, "Io Sottoscritto, dichiaro di avere perso visione dell’Informativa Privacy ai "+
		"sensi dell’art. 13 del GDPR (informativa resa all’interno del set documentale contenente anche la "+
		"Documentazione Informativa Precontrattuale, il Glossario e le Condizioni di Assicurazione) e di averne "+
		"compreso i contenuti", "", fpdf.AlignLeft, false)
	pdf.Ln(3)
	drawSignatureForm(pdf)
	pdf.Ln(9)

	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, "Qui di seguito esprimo il mio consenso al trattamento dei dati personali "+
		"particolari per le finalità sopra indicate, in conformità con quanto previsto all’interno "+
		"dell’informativa", "", fpdf.AlignLeft, false)
	pdf.Ln(1)
	setBlackBoldFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, "1. Consenso al trattamento dei miei dati al fine di perfezionamento "+
		"dell’offerta assicurativa e riassicurativa di cui alle lettere b) ed f) della presente "+
		"informativa", "", fpdf.AlignLeft, false)
	pdf.Ln(3)
	drawSignatureForm(pdf)
}
