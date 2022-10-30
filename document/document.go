package document

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	lib "github.com/wopta/goworkspace/lib"
	//model "github.com/wopta/goworkspace/models"
)

func init() {
	log.Println("INIT Document")
	functions.HTTP("Document", Document)
}

func Document(w http.ResponseWriter, r *http.Request) {
	log.Println("Document")
	lib.EnableCors(&w, r)
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	lib.Router(w, r, map[string]func(http.ResponseWriter, *http.Request){
		"/v1/proposal": Proposal,
		"/v1/contract": Contract,
		"/v1/sign":     SignNamirial,
	})

}
func getFilesByEnv(file string) []byte {
	var res1 []byte
	switch os.Getenv("env") {

	case "local":
		res1 = lib.ErrorByte(ioutil.ReadFile("function-data/" + file))

	case "dev":
		res1 = lib.GetFromStorage("function-data", file, "")

	case "prod":
		res1 = lib.GetFromStorage("core-350507-function-data", file, "")

	default:

	}
	return res1
}

type Kv struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
type PdfData struct {
	Class        string `json:"class"`
	CoverageType string `json:"coverageType"`
	FiscalCode   string `json:"fiscalCode"`
	Name         string `json:"name"`
	Price        int64  `json:"price"`
	PriceNett    int64  `json:"priceNett"`
	Surname      string `json:"surname"`
	Work         string `json:"work"`
	WorkType     string `json:"workType"`
	Coverages    []struct {
		Deductible                 string `json:"deductible"`
		Name                       string `json:"name"`
		Price                      int64  `json:"price"`
		PriceNett                  int64  `json:"priceNett"`
		SelfInsurance              string `json:"selfInsurance"`
		SumInsuredLimitOfIndemnity int64  `json:"sumInsuredLimitOfIndemnity"`
	} `json:"coverages"`
}
