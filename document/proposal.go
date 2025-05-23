package document

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"github.com/go-pdf/fpdf"
	"gitlab.dev.wopta.it/goworkspace/document/internal/engine"
	"gitlab.dev.wopta.it/goworkspace/document/pkg/contract"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/network"
	prd "gitlab.dev.wopta.it/goworkspace/product"
)

func ProposalFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err         error
		policy      *models.Policy
		warrant     *models.Warrant
		product     *models.Product
		networkNode *models.NetworkNode
	)
	log.AddPrefix("ProposalFx")
	defer log.PopPrefix()
	log.Println("handler start ---------------------------------")

	origin := r.Header.Get("Origin")
	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	log.Printf("body: %s", string(body))

	err = json.Unmarshal(body, &policy)
	if err != nil {
		log.ErrorF("error unmarshaling request body: %s", err.Error())
		return "", nil, err
	}

	networkNode = network.GetNetworkNodeByUid(policy.ProducerUid)
	if networkNode != nil {
		warrant = networkNode.GetWarrant()
	}

	product = prd.GetProductV2(policy.Name, policy.ProductVersion, models.MgaChannel, networkNode, warrant)

	result := Proposal(origin, policy, networkNode, product)

	respJson, err := json.Marshal(result)

	log.Println("handler end ---------------------------------")

	return string(respJson), result, err
}

func Proposal(origin string, policy *models.Policy, networkNode *models.NetworkNode, product *models.Product) *DocumentResponse {
	var (
		pdf      *fpdf.Fpdf
		err      error
		rawDoc   []byte
		filename string
	)
	log.AddPrefix("Proposal")
	defer log.PopPrefix()
	log.Println("function start ----------------------------------")

	rawPolicy, _ := policy.Marshal()
	log.Printf("policy: %s", string(rawPolicy))

	log.Printf("generating proposal document for proposal n. %d", policy.ProposalNumber)

	pdf = initFpdf()

	switch policy.Name {
	case models.LifeProduct:
		log.Println("call lifeProposal...")
		filename, rawDoc = lifeProposal(pdf, origin, policy, networkNode, product)
	case models.GapProduct:
		log.Println("call gapProposal...")
		filename, rawDoc = gapProposal(pdf, origin, policy, networkNode)
	case models.PersonaProduct:
		log.Println("call personaProposal...")
		filename, rawDoc = personaProposal(pdf, policy, networkNode, product)
	case models.CatNatProduct:
		//to change
		pdf := initFpdf()
		filename, rawDoc = saveProposal(pdf, policy)
	case models.CommercialCombinedProduct:
		generator := contract.NewCommercialCombinedGenerator(engine.NewFpdf(), policy, networkNode, *product, true)
		rawDoc, err = generator.Contract()
		if err != nil {
			log.ErrorF("error generating contract: %v", err)
			return nil
		}
		filename, err = generator.Save(rawDoc)
		if err != nil {
			log.ErrorF("error generating contract: %v", err)
			return nil
		}
	}

	log.Printf("proposal document generated for proposal n. %d", policy.ProposalNumber)

	log.Println("function end ----------------------------------")

	return &DocumentResponse{
		LinkGcs: filename,
		Bytes:   base64.StdEncoding.EncodeToString(rawDoc),
	}
}

func lifeProposal(pdf *fpdf.Fpdf, origin string, policy *models.Policy, networkNode *models.NetworkNode, product *models.Product) (string, []byte) {
	var (
		rawDoc   []byte
		filename string
	)
	log.AddPrefix("lifeProposal")
	defer log.PopPrefix()
	log.Println("function start ------------------------------")

	switch policy.ProductVersion {
	case models.ProductV1:
		log.Println("life v1")
		filename, rawDoc = lifeAxaProposalV1(pdf, origin, policy, networkNode, product)
	case models.ProductV2:
		log.Println("life v2")
		pdf := engine.NewFpdf()
		gen := contract.NewLifeGenerator(pdf, policy, networkNode, *product, true)
		gen.Generate()
		filename, rawDoc = lifeAxaProposalV2(pdf.GetPdf(), origin, policy, networkNode, product)
	}

	log.Println("function end --------------------------------")

	return filename, rawDoc
}

func gapProposal(pdf *fpdf.Fpdf, origin string, policy *models.Policy, networkNode *models.NetworkNode) (string, []byte) {
	var (
		filename string
		out      []byte
	)
	log.AddPrefix("gapProposal")
	log.Println("function start -------------------------------")

	filename, out = gapSogessurProposalV1(pdf, origin, policy, networkNode)

	log.Println("function end ---------------------------------")

	return filename, out
}

func personaProposal(pdf *fpdf.Fpdf, policy *models.Policy, networkNode *models.NetworkNode, product *models.Product) (string, []byte) {
	var (
		filename string
		out      []byte
	)
	log.AddPrefix("personaProposal")
	defer log.PopPrefix()
	log.Println("function start ---------------------------")

	filename, out = personaGlobalProposalV1(pdf, policy, networkNode, product)

	log.Println("function end -----------------------------")

	return filename, out
}
