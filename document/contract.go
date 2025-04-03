package document

import (
	"encoding/base64"
	"encoding/json"
	"github.com/wopta/goworkspace/lib/log"
	"io"
	"net/http"

	"github.com/go-pdf/fpdf"
	"github.com/wopta/goworkspace/document/internal/engine"
	"github.com/wopta/goworkspace/document/pkg/contract"
	"github.com/wopta/goworkspace/network"
	prd "github.com/wopta/goworkspace/product"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

var (
	signatureID int
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
	log.AddPrefix("ContractObj")
	defer log.PopPrefix()

	log.Println("function start -------------------------------")

	rawPolicy, _ := data.Marshal()
	log.Printf("policy: %s", string(rawPolicy))

	go func() {
		var (
			err      error
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
		case models.CommercialCombinedProduct:
			generator := contract.NewCommercialCombinedGenerator(engine.NewFpdf(), &data, networkNode, *product, false)
			out, err = generator.Contract()
			if err != nil {
				log.ErrorF("error generating contract: %v", err)
				return
			}
			filename, err = generator.Save(out)
			if err != nil {
				log.ErrorF("error generating contract: %v", err)
				return
			}
		}

		data.DocumentName = filename
		log.Println(data.Uid + " ContractObj end")
		r <- DocumentResponse{
			LinkGcs: filename,
			Bytes:   base64.StdEncoding.EncodeToString(out),
		}
	}()

	log.Println("function end -------------------------------..")

	return r
}

func lifeContract(pdf *fpdf.Fpdf, origin string, policy *models.Policy, networkNode *models.NetworkNode, product *models.Product) (string, []byte) {
	var (
		filename string
		out      []byte
	)
	log.AddPrefix("LifeContract")
	defer log.PopPrefix()
	log.Println("function start ------------------------------")

	switch policy.ProductVersion {
	case models.ProductV1:
		log.Println("life v1")
		filename, out = lifeAxaContractV1(pdf, origin, policy, networkNode, product)
	case models.ProductV2:
		log.Println("life v2")
		filename, out = lifeAxaContractV2(pdf, origin, policy, networkNode, product)
	}

	log.Println("function end --------------------------------")

	return filename, out
}

func gapContract(pdf *fpdf.Fpdf, origin string, policy *models.Policy, networkNode *models.NetworkNode) (string, []byte) {
	var (
		filename string
		out      []byte
	)
	log.AddPrefix("GapContract")
	defer log.PopPrefix()
	log.Println("function start -------------------------------")

	filename, out = gapSogessurContractV1(pdf, origin, policy, networkNode)

	log.Println("function end ---------------------------------")

	return filename, out
}

func personaContract(pdf *fpdf.Fpdf, policy *models.Policy, networkNode *models.NetworkNode, product *models.Product) (string, []byte) {
	var (
		filename string
		out      []byte
	)
	log.AddPrefix("personaContract")
	defer log.PopPrefix()
	log.Println("function start ---------------------------")

	filename, out = personaGlobalContractV1(pdf, policy, networkNode, product)

	log.Println("function end -----------------------------")

	return filename, out
}
