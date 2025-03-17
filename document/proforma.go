package document

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/wopta/goworkspace/document/internal/engine"
	"github.com/wopta/goworkspace/document/pkg/proforma"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
	prd "github.com/wopta/goworkspace/product"
)

type ProformaResponse struct {
	FileName string
	LinkGcs  string `json:"linkGcs"`
}

func ProformaFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println("[Proforma]")
	origin := r.Header.Get("Origin")
	req := lib.ErrorByte(io.ReadAll(r.Body))
	var data models.Policy
	defer r.Body.Close()
	err := json.Unmarshal([]byte(req), &data)
	lib.CheckError(err)

	var warrant *models.Warrant
	networkNode := network.GetNetworkNodeByUid(data.ProducerUid)
	if networkNode != nil {
		warrant = networkNode.GetWarrant()
	}

	product := prd.GetProductV2(data.Name, data.ProductVersion, models.MgaChannel, networkNode, warrant)

	respObj := <-ProformaObj(origin, data, networkNode, product) // TODO review product nil
	resp, err := json.Marshal(respObj)

	lib.CheckError(err)
	return string(resp), respObj, nil
}

func ProformaObj(origin string, data models.Policy, networkNode *models.NetworkNode, product *models.Product) <-chan ProformaResponse {
	r := make(chan ProformaResponse)

	log.Println("[ProformaObj] function start -------------------------------")

	rawPolicy, _ := data.Marshal()
	log.Printf("[ProformaObj] policy: %s", string(rawPolicy))

	go func() {
		var (
			err      error
			filename string
			out      []byte
		)

		generator := proforma.NewProformaGenerator(engine.NewFpdf(), &data, networkNode, *product)
		out, err = generator.Contract()
		if err != nil {
			log.Printf("error generating proforma: %v", err)
			return
		}
		filename, err = generator.Save(out)
		if err != nil {
			log.Printf("error saving contract: %v", err)
			return
		}

		log.Println(data.Uid + " ContractObj end")
		r <- ProformaResponse{
			LinkGcs:  filename,
			FileName: filename,
		}
	}()

	log.Println("[ProformaObj] function end -------------------------------..")

	return r
}
