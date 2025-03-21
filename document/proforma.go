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

	prof, err := ProformaObj(origin, data, networkNode, product) // TODO review product nil
	if err != nil {
		log.Printf("unable to generate proforma: %s", err.Error())
		return "", nil, err
	}
	resp, err := json.Marshal(prof)

	lib.CheckError(err)
	return string(resp), prof, nil
}

func ProformaObj(origin string, data models.Policy, networkNode *models.NetworkNode, product *models.Product) (ProformaResponse, error) {
	var (
		err      error
		filename string
		out      []byte
	)

	log.Println("[ProformaObj] function start -------------------------------")

	generator := proforma.NewProformaGenerator(engine.NewFpdf(), &data, networkNode, *product)
	out, err = generator.Generate()
	if err != nil {
		log.Printf("error generating proforma: %v", err)
		return ProformaResponse{}, err
	}
	filename, err = generator.Save(out)
	if err != nil {
		log.Printf("error saving proforma: %v", err)
		return ProformaResponse{}, err
	}

	log.Println(data.Uid + " ProformaObj end")
	res := ProformaResponse{
		LinkGcs:  filename,
		FileName: filename,
	}

	log.Println("[ProformaObj] function end -------------------------------..")

	return res, nil
}
