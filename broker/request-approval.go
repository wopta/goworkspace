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
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
	plc "github.com/wopta/goworkspace/policy"
	"github.com/wopta/goworkspace/reserved"
)

type RequestApprovalReq = BrokerBaseRequest

const AlreadyInsured int = 9999

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

	if err = requestApproval(&policy); err != nil {
		log.Println("error request approval")
		return "", nil, err
	}

	jsonOut, err := policy.Marshal()

	return string(jsonOut), policy, err
}

func requestApproval(policy *models.Policy) error {
	var (
		err error
	)

	log.Println("[requestApproval] start -------------------------------------")

	networkNode = network.GetNetworkNodeByUid(policy.ProducerUid)
	if networkNode != nil {
		warrant = networkNode.GetWarrant()
	}

	log.Println("[requestApproval] starting bpmn flow...")

	state := runBrokerBpmn(policy, requestApprovalFlowKey)
	if state == nil || state.Data == nil {
		log.Println("[requestApproval] error bpmn - state not set")
		return errors.New("error on bpmn - no data present")
	}
	if state.IsFailed {
		log.Println("[requestApproval] error bpmn - state failed")
		return nil
	}

	*policy = *state.Data

	log.Printf("[requestApproval] saving policy with uid %s to bigquery...", policy.Uid)
	policy.BigquerySave(origin)

	callback_out.Execute(networkNode, *policy, callback_out.RequestApproval)

	log.Println("[requestApproval] end ---------------------------------------")

	return err
}

func setRequestApprovalData(policy *models.Policy) error {
	const (
		mgaApprovalFlow     = "MgaApprovalFlow"
		companyApprovalFlow = "CompanyApprovalFlow"
	)
	var (
		flow              string
	)
	log.Printf("[setRequestApprovalData] policy uid %s: reserved flow", policy.Uid)

	isOldReserved := len(policy.ReservedInfo.Reasons) > 0

	if isOldReserved {
		oldRequestApproval(policy)
		return nil
	}

	setProposalNumber(policy)
	createProposalDoc := policy.Status == models.PolicyStatusInitLead

	if policy.ReservedInfo.CompanyApproval.Mandatory {
		flow = companyApprovalFlow
	}
	if policy.ReservedInfo.MgaApproval.Mandatory {
		flow = mgaApprovalFlow
	}

	switch flow {
	case mgaApprovalFlow:
		requestMgaApproval(policy)
	case companyApprovalFlow:
		requestCompanyApproval(policy)
	default:
		log.Println("flow not set")
		return errNotAllowed
	}

	for _, reason := range policy.ReservedInfo.ReservedReasons {
		if reason.Id == AlreadyInsured {
			policy.Status = models.PolicyStatusWaitForApprovalMga
			break
		}
	}

	reserved.SetReservedInfo(policy, mgaProduct)

	if createProposalDoc {
		plc.AddProposalDoc(origin, policy, networkNode, mgaProduct)
	}

	policy.StatusHistory = append(policy.StatusHistory, policy.Status)
	policy.Updated = time.Now().UTC()

	return nil
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
