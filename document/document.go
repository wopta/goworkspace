package document

/*

 */
import (
	"log"
	"net/http"

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
	route := lib.RouteData{
		Routes: []lib.Route{{
			Route:   "/v1/proposal",
			Hendler: Contract,
		},
			{
				Route:   "/v1/contract",
				Hendler: Contract,
			},
			{
				Route:   "/v1/sign",
				Hendler: SignNamirial,
			},
		},
	}
	route.Router(w, r)

}

type Kv struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type DodumentResponse struct {
	EnvelopSignId string `json:"envelopSignId"`
	LinkGcs       string `json:"linkGcs"`
	Bytes         string `json:"bytes"`
}

type DodumentData struct {
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
