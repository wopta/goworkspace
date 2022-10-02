package appcheck-proxy

import (

	"fmt"
	import "google.golang.org/api/firebaseappcheck/v1"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
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
	w.Header().Set("Content-Type", "application/pdf")
	w.WriteHeader(http.StatusOK)

	fmt.Fprintf(w, string(pdfg.Bytes()))

}

type PdfData struct {
	name string
}
