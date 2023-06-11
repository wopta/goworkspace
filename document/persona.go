package document

import (
	"github.com/dustin/go-humanize"
	"github.com/go-pdf/fpdf"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/product"
	"sort"
	"strings"
)

type keyValue struct {
	key   string
	value string
}

func PersonaContract(pdf *fpdf.Fpdf, policy *models.Policy) (string, []byte) {
	var (
		filename string
		out      []byte
	)

	filename, out = Persona(pdf, policy)

	return filename, out
}

func Persona(pdf *fpdf.Fpdf, policy *models.Policy) (string, []byte) {
	signatureID = 0

	mainHeader(pdf, policy)

	mainFooter(pdf, policy.Name)

	pdf.AddPage()

	personaInsuredInfoSection(pdf, policy)

	personaGuaranteesTable(pdf, policy)

	personaSurveySection(pdf, policy)

	personaStatementsSection(pdf, policy)

	personaOfferResumeSection(pdf, policy)

	paymentMethodSection(pdf)

	emitResumeSection(pdf, policy)

	companiesDescriptionSection(pdf, policy.Company)

	personalDataHandlingSection(pdf, policy)

	globalHeader(pdf)

	pdf.AddPage()

	globalFooter(pdf)

	globalStamentsAndConsens(pdf)

	filename, out := save(pdf, policy)
	return filename, out
}

func personaInsuredInfoSection(pdf *fpdf.Fpdf, policy *models.Policy) {
	coverageTypeMap := map[string]string{
		"24h":   "Professionale ed Extraprofessionale",
		"prof":  "Professionale",
		"extra": "Extraprofessionale",
	}

	getParagraphTitle(pdf, "La tua assicurazione per il seguente Assicurato e Garanzie")
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

func personaGuaranteesTable(pdf *fpdf.Fpdf, policy *models.Policy) {
	type slugStruct struct {
		name  string
		order int64
	}

	var table [][]string
	offerName := policy.OfferlName
	prod, err := product.GetProduct("persona", "v1")
	lib.CheckError(err)

	guaranteesMap := map[string]map[string]string{}
	var slugs []slugStruct

	for guaranteeSlug, guarantee := range prod.Companies[0].GuaranteesMap {
		guaranteesMap[guaranteeSlug] = make(map[string]string, 0)

		guaranteesMap[guaranteeSlug]["name"] = guarantee.CompanyName
		guaranteesMap[guaranteeSlug]["sumInsuredLimitOfIndemnity"] = "====="
		guaranteesMap[guaranteeSlug]["details"] = "====="
		guaranteesMap[guaranteeSlug]["price"] = "====="
		slugs = append(slugs, slugStruct{name: guaranteeSlug, order: guarantee.Order})
	}

	sort.Slice(slugs, func(i, j int) bool {
		return slugs[i].order < slugs[j].order
	})

	for _, asset := range policy.Assets {
		for _, guarantee := range asset.Guarantees {
			var price float64
			var details string

			guaranteesMap[guarantee.Slug]["sumInsuredLimitOfIndemnity"] = humanize.FormatFloat("#.###,", guarantee.Offer[offerName].SumInsuredLimitOfIndemnity) + " €"
			if policy.PaymentSplit == string(models.PaySplitMonthly) {
				price = guarantee.Value.PremiumGrossMonthly * 12
			} else {
				price = guarantee.Value.PremiumGrossYearly
			}
			guaranteesMap[guarantee.Slug]["price"] = humanize.FormatFloat("#.###,##", price) + " €"

			switch guarantee.Slug {
			case "IPI":
				details = "Franchigia " + guarantee.Value.Deductible + guarantee.Value.DeductibleUnit
				if guarantee.Value.DeductibleType == "absolute" {
					details += " Assoluta"
				} else {
					details += " Assorbibile"
				}
			case "D":
				if guarantee.Beneficiaries != nil {
					details = "Beneficiari\n"
					for _, beneficiary := range *guarantee.Beneficiaries {
						details += beneficiary.Name + " " + beneficiary.Surname + "\n"
					}
				} else {
					details = "===="
				}
			case "ITI":
				details = "Franchigia " + guarantee.Value.Deductible + " " + guarantee.Offer[offerName].DeductibleUnit
			default:
				details = "====="
			}
			guaranteesMap[guarantee.Slug]["details"] = details
		}
	}

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
		setBlackBoldFont(pdf, standardTextSize)
		pdf.CellFormat(80, 6, guaranteesMap[slug.name]["name"], "B", 0, fpdf.AlignLeft, false, 0, "")
		setBlackRegularFont(pdf, standardTextSize)
		pdf.CellFormat(30, 6, guaranteesMap[slug.name]["sumInsuredLimitOfIndemnity"], "B", 0, fpdf.AlignRight, false, 0, "")
		pdf.CellFormat(5, 6, "", "B", 0, fpdf.AlignRight, false, 0, "")
		pdf.CellFormat(60, 6, guaranteesMap[slug.name]["details"], "B", 0, fpdf.AlignLeft, false, 0, "")
		pdf.CellFormat(15, 6, guaranteesMap[slug.name]["price"], "B", 1, fpdf.AlignRight, false, 0, "")
	}
}

func personaSurveySection(pdf *fpdf.Fpdf, policy *models.Policy) {
	surveys := *policy.Surveys

	getParagraphTitle(pdf, "Dichiarazioni da leggere con attenzione prima di firmare")
	err := printSurvey(pdf, surveys[0])
	lib.CheckError(err)

	getParagraphTitle(pdf, "Questionario Medico")
	if len(surveys) == 3 {
		for _, survey := range surveys[1:2] {
			err = printSurvey(pdf, survey)
			lib.CheckError(err)
		}
		pdf.AddPage()
	} else {
		for _, survey := range surveys[1:3] {
			err = printSurvey(pdf, survey)
			lib.CheckError(err)
		}
	}

	surveys[len(surveys)-1].Title = ""
	getParagraphTitle(pdf, "Tutela Privacy")
	err = printSurvey(pdf, surveys[len(surveys)-1])
	lib.CheckError(err)

	pdf.Ln(5)
	drawSignatureForm(pdf)
	pdf.Ln(10)
}

func personaStatementsSection(pdf *fpdf.Fpdf, policy *models.Policy) {
	statements := *policy.Statements

	for _, statement := range statements {
		printStatement(pdf, statement)
	}
	pdf.SetY(pdf.GetY() - 28)
	setBlackBoldFont(pdf, standardTextSize)
	pdf.MultiCell(70, 3, "Global Assistance", "",
		fpdf.AlignCenter, false)
	var opt fpdf.ImageOptions
	opt.ImageType = "png"
	pdf.ImageOptions(lib.GetAssetPathByEnv(basePath)+"/firma_global.png", 30, pdf.GetY()+3, 40, 12,
		false, opt, 0, "")
	pdf.Ln(20)
}

func personaOfferResumeSection(pdf *fpdf.Fpdf, policy *models.Policy) {
	var (
		tableInfo [][]string
	)

	switch policy.PaymentSplit {
	case string(models.PaySplitMonthly):
		tableInfo = [][]string{
			{
				"Mensile firma del contratto",
				lib.HumanaizePriceEuro(policy.PriceNett),
				lib.HumanaizePriceEuro(policy.PriceGross - policy.PriceNett),
				lib.HumanaizePriceEuro(policy.PriceGross),
			},
			{
				"Pari ad un premio Annuale",
				lib.HumanaizePriceEuro(policy.PriceNett * 12),
				lib.HumanaizePriceEuro(policy.PriceGross - policy.PriceNett*12),
				lib.HumanaizePriceEuro(policy.PriceGross * 12),
			},
		}
	case string(models.PaySplitYear):
		tableInfo = [][]string{
			{
				"Annuale firma del contratto",
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
