package document

import (
	"github.com/go-pdf/fpdf"
	"log"
	"net/http"
)

func FpdfHandler(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {

	pdf := initFpdf()

	GetHeader(pdf, "life")

	err := pdf.OutputFileAndClose("document/test.pdf")
	log.Println(err)
	return "", nil, err
}

func initFpdf() *fpdf.Fpdf {
	pdf := fpdf.New("P", "mm", "A4", "")
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
		pdf.Cell(40, 10, "Wopta per te")
		pdf.SetXY(23, 12)
		pdf.SetFont("Montserrat", "I", 18)
		pdf.SetTextColor(92, 89, 92)
		pdf.Cell(40, 10, product)
		pdf.ImageOptions("document/assets/ARTW_LOGO_RGB_400px.png", 158, 6, 0, 10, false, opt, 0, "")

	})
}
