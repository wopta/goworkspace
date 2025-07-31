package document

import (
	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"github.com/go-pdf/fpdf"
	"gitlab.dev.wopta.it/goworkspace/document/internal/engine"
	"gitlab.dev.wopta.it/goworkspace/document/pkg/contract"
	"gitlab.dev.wopta.it/goworkspace/models"
)

func Proposal(policy *models.Policy, networkNode *models.NetworkNode, product *models.Product) (DocumentGenerated, error) {
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
		document, err = lifeProposal(pdf, policy, networkNode, product)
	case models.GapProduct:
		log.Println("call gapProposal...")
		document, err = gapSogessurProposalV1(pdf, policy, networkNode)
	case models.PersonaProduct:
		log.Println("call personaProposal...")
		pdf := engine.NewFpdf()
		generator := contract.NewPersonaGenerator(pdf, policy, networkNode, *product, true)
		personaGlobalProposalV1(pdf.GetPdf(), policy, networkNode, product)
		generator.AddMup()
		document, err = generateProposalDocument(pdf.GetPdf(), policy)
	case models.CatNatProduct:
		pdf := engine.NewFpdf()
		generator := contract.NewCatnatGenerator(pdf, policy, networkNode, *product, true)
		generator.Generate()
		document, err = generateContractDocument(pdf.GetPdf(), policy)
	}

	log.Printf("proposal document generated for proposal n. %d", policy.ProposalNumber)

	log.Println("function end ----------------------------------")

	return document, err
}

func lifeProposal(pdf *fpdf.Fpdf, policy *models.Policy, networkNode *models.NetworkNode, product *models.Product) (DocumentGenerated, error) {
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
		document, err = lifeAxaProposalV1(pdf, policy, networkNode, product)
	case models.ProductV2:
		log.Println("life v2")
		pdf := engine.NewFpdf()
		gen := contract.NewLifeGenerator(pdf, policy, networkNode, *product, true)
		gen.Generate()
		lifeAxaProposalV2(pdf.GetPdf(), policy, networkNode, product)
		gen.AddMup()
		document, err = generateProposalDocument(pdf.GetPdf(), policy)
	}

	log.Println("function end --------------------------------")

	return document, err
}
