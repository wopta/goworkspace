package document

import (
	"encoding/base64"
	"encoding/json"
	"github.com/go-pdf/fpdf"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
	prd "github.com/wopta/goworkspace/product"
	"io"
	"log"
	"net/http"
)

func ProposalFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err         error
		policy      *models.Policy
		warrant     *models.Warrant
		product     *models.Product
		networkNode *models.NetworkNode
	)

	log.Println("[ProposalFx] handler start ---------------------------------")

	origin := r.Header.Get("Origin")
	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	log.Printf("[ProposalFx] body: %s", string(body))

	err = json.Unmarshal(body, &policy)
	if err != nil {
		log.Printf("[ProposalFx] error unmarshaling request body: %s", err.Error())
		return "", nil, err
	}

	networkNode = network.GetNetworkNodeByUid(policy.ProducerUid)
	if networkNode != nil {
		warrant = networkNode.GetWarrant()
	}

	product = prd.GetProductV2(policy.Name, policy.ProductVersion, models.MgaChannel, networkNode, warrant)

	result := Proposal(origin, policy, networkNode, product)

	respJson, err := json.Marshal(result)

	log.Println("[ProposalFx] handler end ---------------------------------")

	return string(respJson), result, err
}

func Proposal(origin string, policy *models.Policy, networkNode *models.NetworkNode, product *models.Product) *DocumentResponse {
	var (
		pdf      *fpdf.Fpdf
		rawDoc   []byte
		filename string
	)

	log.Println("[Proposal] function start ----------------------------------")

	log.Printf("[Proposal] generating proposal document for proposal n. %d", policy.ProposalNumber)

	pdf = initFpdf()

	switch policy.Name {
	case models.LifeProduct:
		log.Println("[Proposal] call LifeProposal...")
		filename, rawDoc = LifeProposal(pdf, origin, policy, networkNode, product)
	}

	log.Printf("[Proposal] proposal document generated for proposal n. %d", policy.ProposalNumber)

	log.Println("[Proposal] function end ----------------------------------")

	return &DocumentResponse{
		LinkGcs: filename,
		Bytes:   base64.StdEncoding.EncodeToString(rawDoc),
	}
}
