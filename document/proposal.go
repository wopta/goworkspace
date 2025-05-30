package document

import (
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

	result, err := Proposal(origin, policy, networkNode, product)
	response, err := result.Save()
	if err != nil {
		return "", nil, err
	}
	respJson, err := json.Marshal(response)

	log.Println("handler end ---------------------------------")

	return string(respJson), result, err
}

func Proposal(origin string, policy *models.Policy, networkNode *models.NetworkNode, product *models.Product) (DocumentGenerated, error) {
	var (
		pdf      *fpdf.Fpdf
		err      error
		document DocumentGenerated
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
		document, err = lifeProposal(pdf, origin, policy, networkNode, product)
	case models.GapProduct:
		log.Println("call gapProposal...")
		document, err = gapSogessurProposalV1(pdf, origin, policy, networkNode)
	case models.PersonaProduct:
		log.Println("call personaProposal...")
		pdf := engine.NewFpdf()
		generator := contract.NewPersonaGenerator(pdf, policy, networkNode, *product, true)
		personaGlobalProposalV1(pdf.GetPdf(), policy, networkNode, product)
		generator.AddMup()
		document, err = generateProposalDocument(pdf.GetPdf(), policy)
	case models.CatNatProduct:
		//to change
		pdf := initFpdf()
		document, err = generateProposalDocument(pdf, policy)
	}

	log.Printf("proposal document generated for proposal n. %d", policy.ProposalNumber)

	log.Println("function end ----------------------------------")

	return document, err
}

func lifeProposal(pdf *fpdf.Fpdf, origin string, policy *models.Policy, networkNode *models.NetworkNode, product *models.Product) (DocumentGenerated, error) {
	var (
		document DocumentGenerated
		err      error
	)
	log.AddPrefix("lifeProposal")
	defer log.PopPrefix()
	log.Println("function start ------------------------------")

	switch policy.ProductVersion {
	case models.ProductV1:
		log.Println("life v1")
		document, err = lifeAxaProposalV1(pdf, origin, policy, networkNode, product)
	case models.ProductV2:
		log.Println("life v2")
		pdf := engine.NewFpdf()
		gen := contract.NewLifeGenerator(pdf, policy, networkNode, *product, true)
		gen.Generate()
		lifeAxaProposalV2(pdf.GetPdf(), origin, policy, networkNode, product)
		gen.AddMup()
		document, err = generateProposalDocument(pdf.GetPdf(), policy)
	}

	log.Println("function end --------------------------------")

	return document, err
}
