package document

import (
	"encoding/base64"
	"encoding/json"
	"github.com/wopta/goworkspace/network"
	prd "github.com/wopta/goworkspace/product"
	"io"
	"log"
	"net/http"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func ContractFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println("[Contract]")
	//lib.Files("./serverless_function_source_code")
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

	respObj := <-ContractObj(origin, data, networkNode, product) // TODO review product nil
	resp, err := json.Marshal(respObj)

	lib.CheckError(err)
	return string(resp), respObj, nil
}

func ContractObj(origin string, data models.Policy, networkNode *models.NetworkNode, product *models.Product) <-chan DocumentResponse {
	r := make(chan DocumentResponse)

	go func() {
		var (
			filename string
			out      []byte
		)

		switch data.Name {
		case models.PmiProduct:
			skin := getVar()
			m := skin.initDefault()
			skin.GlobalContract(m, data)
			//-----------Save file
			filename, out = Save(m, data)
		case models.LifeProduct:
			pdf := initFpdf()
			filename, out = lifeContract(pdf, origin, &data, networkNode, product)
		case models.PersonaProduct:
			pdf := initFpdf()
			filename, out = personaContract(pdf, &data, networkNode, product)
		case models.GapProduct:
			pdf := initFpdf()
			filename, out = gapContract(pdf, origin, &data, networkNode)
		}

		data.DocumentName = filename
		log.Println(data.Uid + " ContractObj end")
		r <- DocumentResponse{
			LinkGcs: filename,
			Bytes:   base64.StdEncoding.EncodeToString(out),
		}
	}()
	return r
}
