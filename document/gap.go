package document

import (
	"github.com/go-pdf/fpdf"
	"github.com/wopta/goworkspace/models"
)

func GapContract(pdf *fpdf.Fpdf, origin string, policy *models.Policy) (string, []byte) {
	var (
		filename string
		out      []byte
	)

	filename, out = GapSogessur(pdf, origin, policy)

	return filename, out
}

func GapSogessur(pdf *fpdf.Fpdf, origin string, policy *models.Policy) (string, []byte) {
	signatureID = 0

	mainMotorHeader(pdf, policy)

	mainFooter(pdf, policy.Name)

	pdf.AddPage()

	getParagraphTitle(pdf, "La tua assicurazione è operante sui dati sotto riportati, verifica la loro correttezza"+
		" e segnala eventuali inesattezze")

	filename, out := save(pdf, policy)
	return filename, out
}
