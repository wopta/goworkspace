package document

import (
	"github.com/go-pdf/fpdf"
	"log"
	"net/http"
)

func FpdfHandler(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {

	pdf := initFpdf()

	GetHeader(pdf, "life")

	pdf.AddPage()
	pdf.AddPage()

	err := pdf.OutputFileAndClose("document/test.pdf")
	log.Println(err)
	return "", nil, err
}

func initFpdf() *fpdf.Fpdf {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(10, 15, 10)
	loadCustomFonts(pdf)
	return pdf
}

func loadCustomFonts(pdf *fpdf.Fpdf) {
	pdf.AddUTF8Font("Montserrat", "", "document/assets/montserrat_regular.ttf")
	pdf.AddUTF8Font("Montserrat", "B", "document/assets/montserrat_bold.ttf")
	pdf.AddUTF8Font("Montserrat", "I", "document/assets/montserrat_italic.ttf")
}

func GetHeader(pdf *fpdf.Fpdf, name string) {
	var opt fpdf.ImageOptions
	var product, logoPath string
	pathPrefix := "document/assets/logo_"

	switch name {
	case "life":
		product = "Vita"
		logoPath = pathPrefix + "vita.png"
	}

	pdf.SetHeaderFunc(func() {
		opt.ImageType = "png"
		pdf.ImageOptions(logoPath, 10, 6, 13, 13, false, opt, 0, "")
		pdf.SetXY(23, 6)
		pdf.SetTextColor(229, 0, 117)
		pdf.SetFont("Montserrat", "B", 18)
		pdf.Text(24, 12, "Wopta per te")
		pdf.SetFont("Montserrat", "I", 18)
		pdf.SetTextColor(92, 89, 92)
		pdf.Text(24, 18, product)
		pdf.ImageOptions("document/assets/ARTW_LOGO_RGB_400px.png", 158, 6, 0, 10, false, opt, 0, "")

		pdf.SetTextColor(0, 0, 0)
		pdf.SetFont("Montserrat", "B", 8)
		pdf.Text(12, 23, "I dati della tua polizza")
		pdf.SetFont("Montserrat", "", 8)
		pdf.SetXY(11, pdf.GetY()+18)
		pdf.MultiCell(0, 3, "Numero: 12345\nDecorre dal: 03/04/2023 ore 24:00\nScade il: 03/04/2043 ore 24:00\nPrima scadenza annuale il: 04/04/2024\nNon si rinnova a scadenza.", "", "", false)

		pdf.SetFont("Montserrat", "B", 8)
		pdf.Text(120, 23, "I tuoi dati")
		pdf.SetFont("Montserrat", "", 8)
		pdf.SetXY(119, 24)
		pdf.MultiCell(0, 3, "Contraente: HAMMAR YOUSEF\nC.F./P.IVA: HMMYSF94R07D912M\nIndirizzo: Via Unicef, 4\n20033 SOLARO (MI)\nMail: yousef.hammar@wopta.it\nTelefono: +393451031004", "", "", false)
	})
}
