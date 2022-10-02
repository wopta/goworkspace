package document

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	pdf "github.com/SebastiaanKlippert/go-wkhtmltopdf"
	lib "github.com/wopta/goworkspace/lib"
)

func init() {
	log.Println("INIT Document")

	functions.HTTP("Document", Document)
}

func Document(w http.ResponseWriter, r *http.Request) {
	data := PdfData{name: ""}
	var templ *template.Template
	var err error
	// use Go's default HTML template generation tools to generate your HTML
	templ, err = template.ParseFiles("template.html")
	// apply the parsed HTML template data and keep the result in a Buffer
	var body bytes.Buffer
	err = templ.Execute(&body, data)
	lib.CheckError(err)
	//SETTING
	pdfg, err := pdf.NewPDFGenerator()
	pdfg.Dpi.Set(600)
	pdfg.NoCollate.Set(false)
	pdfg.PageSize.Set(pdf.PageSizeA4)
	pdfg.MarginBottom.Set(40)
	lib.CheckError(err)
	page := pdf.NewPageReader(bytes.NewReader(body.Bytes()))
	page.EnableLocalFileAccess.Set(true)
	// add the page to your generator
	pdfg.AddPage(page)
	// manipulate page attributes as needed
	pdfg.MarginLeft.Set(0)
	pdfg.MarginRight.Set(0)
	pdfg.Dpi.Set(300)
	pdfg.PageSize.Set(pdf.PageSizeA4)
	pdfg.Orientation.Set(pdf.OrientationLandscape)

	// magic
	err = pdfg.Create()
	lib.CheckError(err)

	w.Header().Set("Content-Disposition", "attachment; filename=wopta document test.pdf")
	w.Header().Set("Content-Type", "application/pdf")
	w.WriteHeader(http.StatusOK)

	fmt.Fprintf(w, string(pdfg.Bytes()))

}

type PdfData struct {
	name string
}
