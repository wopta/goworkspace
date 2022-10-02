package appcheck-proxy

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	import "google.golang.org/api/firebaseappcheck/v1"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	pdf "github.com/SebastiaanKlippert/go-wkhtmltopdf"
	lib "github.com/wopta/goworkspace/lib"
)

func init() {
	log.Println("INIT AppcheckProxy")

	functions.HTTP("AppcheckProxy", AppcheckProxy)
}

func AppcheckProxy(w http.ResponseWriter, r *http.Request) {

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
