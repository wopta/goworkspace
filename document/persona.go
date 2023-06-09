package document

import (
	"github.com/go-pdf/fpdf"
	"github.com/wopta/goworkspace/models"
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

}
