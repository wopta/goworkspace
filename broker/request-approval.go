package broker

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/wopta/goworkspace/callback_out"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/mail"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
	plc "github.com/wopta/goworkspace/policy"
	prd "github.com/wopta/goworkspace/product"
	"github.com/wopta/goworkspace/reserved"
)

type RequestApprovalReq = BrokerBaseRequest

var errNotAllowed = errors.New("operation not allowed")

func RequestApprovalFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err    error
		req    RequestApprovalReq
		policy models.Policy
	)

	log.SetPrefix("[RequestApprovalFx] ")
	defer func() {
		r.Body.Close()
		if err != nil {
			log.Printf("error: %s", err)
		}
		log.Println("Handler end ---------------------------------------------")
		log.SetPrefix("")
	}()
	log.Println("Handler start -----------------------------------------------")

	token := r.Header.Get("Authorization")
	authToken, err := lib.GetAuthTokenFromIdToken(token)
	if err != nil {
		log.Printf("error getting authToken")
		return "", nil, err
	}
	log.Printf(
		"authToken - type: '%s' role: '%s' uid: '%s' email: '%s'",
		authToken.Type,
		authToken.Role,
		authToken.UserID,
		authToken.Email,
	)

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("error decoding request body: %s", err)
		return "", nil, err
	}

	if policy, err = plc.GetPolicy(req.PolicyUid, origin); err != nil {
		log.Printf("error fetching policy %s from Firestore...", req.PolicyUid)
		return "", nil, err
	}

	if policy.ProducerUid != authToken.UserID {
		log.Printf("user %s cannot request approval for policy %s because producer not equal to request user",
			authToken.UserID, policy.Uid)
		err = errNotAllowed
		return "", nil, err
	}

	allowedStatus := []string{models.PolicyStatusInitLead, models.PolicyStatusNeedsApproval}
	if !policy.IsReserved || !lib.SliceContains(allowedStatus, policy.Status) {
		log.Printf("cannot request approval for policy with status %s and isReserved %t", policy.Status, policy.IsReserved)
		err = errNotAllowed
		return "", nil, err
	}

	brokerUpdatePolicy(&policy, req)

	var currentFlow string
	if currentFlow, err = requestApproval(&policy); err != nil {
		log.Println("error request approval")
		return "", nil, err
	}
	sendRequestApprovalMail(&policy, currentFlow)

	jsonOut, err := policy.Marshal()

	return string(jsonOut), policy, err
}

// Fallback to old reserved structure
func oldRequestApproval(policy *models.Policy) {
	setProposalNumber(policy)

	if policy.Status == models.PolicyStatusInitLead {
		plc.AddProposalDoc(origin, policy, networkNode, mgaProduct)

		for _, reason := range policy.ReservedInfo.Reasons {
			// TODO: add key/id for reasons so we do not have to check string equallity
			if !strings.HasPrefix(reason, "Cliente già assicurato") {
				reserved.SetReservedInfo(policy, mgaProduct)
				break
			}
		}
	}

	policy.Status = models.PolicyStatusWaitForApproval
	for _, reason := range policy.ReservedInfo.Reasons {
		// TODO: add key/id for reasons so we do not have to check string equallity
		if strings.HasPrefix(reason, "Cliente già assicurato") {
			policy.Status = models.PolicyStatusWaitForApprovalMga
			break
		}
	}

	policy.StatusHistory = append(policy.StatusHistory, policy.Status)
	policy.Updated = time.Now().UTC()
}

func requestCustomerApproval(policy *models.Policy) {
	policy.Status = models.PolicyStatusWaitForApprovalCustomer
	policy.ReservedInfo.CustomerApproval.Status = models.WaitingApproval
	policy.ReservedInfo.CustomerApproval.UpdateDate = time.Now().UTC()
}

func requestMgaApproval(policy *models.Policy) {
	policy.Status = models.PolicyStatusWaitForApproval
	policy.ReservedInfo.MgaApproval.Status = models.WaitingApproval
	policy.ReservedInfo.MgaApproval.UpdateDate = time.Now().UTC()
}

func requestCompanyApproval(policy *models.Policy) {
	policy.Status = models.PolicyStatusWaitForApprovalCompany
	policy.ReservedInfo.CompanyApproval.Status = models.WaitingApproval
	policy.ReservedInfo.CompanyApproval.UpdateDate = time.Now().UTC()
}

const (
	oldFlow              = "oldFlow"
	customerApprovalFlow = "CustomerApprovalFlow"
	mgaApprovalFlow      = "MgaApprovalFlow"
	companyApprovalFlow  = "CompanyApprovalFlow"
)

func getCurrentApprovalFlow(policy *models.Policy) string {
	flow := ""
	if policy.ReservedInfo.CompanyApproval.Mandatory {
		flow = companyApprovalFlow
	}
	if policy.ReservedInfo.MgaApproval.Mandatory {
		flow = mgaApprovalFlow
	}
	if policy.ReservedInfo.CustomerApproval.Mandatory {
		flow = customerApprovalFlow
	}
	if len(policy.ReservedInfo.Reasons) > 0 {
		flow = oldFlow
	}
	return flow
}

func requestApproval(policy *models.Policy) (string, error) {
	var err error

	flow := getCurrentApprovalFlow(policy)

	if flow == oldFlow {
		oldRequestApproval(policy)
		return flow, nil
	}

	setProposalNumber(policy)
	createProposalDoc := policy.Status == models.PolicyStatusInitLead

	switch flow {
	case customerApprovalFlow:
		requestCustomerApproval(policy)
	case mgaApprovalFlow:
		requestMgaApproval(policy)
	case companyApprovalFlow:
		requestCompanyApproval(policy)
	default:
		return "", errors.New("approval flow not set")
	}

	for _, reason := range policy.ReservedInfo.ReservedReasons {
		if reason.Id == reserved.AlreadyInsured {
			policy.Status = models.PolicyStatusWaitForApprovalMga
			break
		}
	}

	product = prd.GetProductV2(policy.Name, policy.ProductVersion, models.MgaChannel, nil, nil)

	reserved.SetReservedInfo(policy, product)

	if createProposalDoc {
		plc.AddProposalDoc(origin, policy, networkNode, product)
	}

	policy.StatusHistory = append(policy.StatusHistory, policy.Status)
	policy.Updated = time.Now().UTC()

	if err := lib.SetFirestoreErr(lib.PolicyCollection, policy.Uid, policy); err != nil {
		return "", err
	}

	policy.BigquerySave(origin)

	callback_out.Execute(networkNode, *policy, callback_out.RequestApproval)

	return flow, err
}

func sendRequestApprovalMail(policy *models.Policy, flow string) {
	switch flow {
	case oldFlow:
		sendOldFlowRequestApprovalMail(policy)
	case customerApprovalFlow:
		sendOldFlowRequestApprovalMail(policy)
	case mgaApprovalFlow:
		sendMgaRequestApprovalMail(policy)
	case companyApprovalFlow:
		sendCompanyRequestApprovalMail(policy)
	}
}

func sendOldFlowRequestApprovalMail(policy *models.Policy) {
	var (
		nn  *models.NetworkNode
		wrt *models.Warrant
	)

	if policy.Status == models.PolicyStatusWaitForApprovalMga {
		return
	}

	if nn = network.GetNetworkNodeByUid(policy.ProducerUid); nn != nil {
		wrt = nn.GetWarrant()
	}
	flowName, _ := policy.GetFlow(nn, wrt)

	from := mail.AddressAnna
	to := mail.GetContractorEmail(policy)
	cc := mail.Address{}
	switch flowName {
	case models.ProviderMgaFlow, models.RemittanceMgaFlow:
		cc = mail.GetNetworkNodeEmail(nn)
	case models.ECommerceChannel:
		to = mail.Address{} // fail safe for not sending email on ecommerce reserved
	}

	mail.SendMailReserved(*policy, from, to, cc, flowName,
		[]string{models.ProposalAttachmentName})
}

func sendMgaRequestApprovalMail(policy *models.Policy) {
	from := mail.AddressAnna
	to := mail.AddressAssunzioneTest
	cc := mail.Address{}

	mail.SendMailMgaRequestApproval(*policy, from, to, cc)
}

func sendCompanyRequestApprovalMail(policy *models.Policy) {
	// TODO: implement
}
