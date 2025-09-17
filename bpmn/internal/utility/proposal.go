package utility

import (
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/models"
	plc "gitlab.dev.wopta.it/goworkspace/policy"
	"gitlab.dev.wopta.it/goworkspace/question"
	"gitlab.dev.wopta.it/goworkspace/reserved"
)

func SetProposalData(policy *models.Policy, networkNode *models.NetworkNode, mgaProduct *models.Product) {
	log.AddPrefix("setProposalData")
	defer log.PopPrefix()
	SetProposalNumber(policy)
	policy.Status = models.PolicyStatusProposal

	if policy.IsReserved {
		log.Println("setting NeedsApproval status")
		policy.Status = models.PolicyStatusNeedsApproval

		for _, reason := range policy.ReservedInfo.Reasons {
			// TODO: add key/id for reasons so we do not have to cjeck string equallity
			if !strings.HasPrefix(reason, "Clientele già assicurato") {
				reserved.SetReservedInfo(policy)
				break
			}
		}
	}

	if policy.Statements == nil || len(*policy.Statements) == 0 {
		var err error
		policy.Statements = new([]models.Statement)

		log.Println("setting policy statements")

		*policy.Statements, err = question.GetStatements(policy, false)
		if err != nil {
			log.ErrorF("error setting policy statements: %s", err.Error())
			return
		}

	}

	plc.AddProposalDoc(policy, networkNode, mgaProduct)

	log.Printf("policy status %s", policy.Status)

	policy.StatusHistory = append(policy.StatusHistory, policy.Status)
	policy.Updated = time.Now().UTC()
}

func SetProposalNumber(policy *models.Policy) {
	log.AddPrefix("SetProposalNumber")
	defer log.PopPrefix()
	log.Println("set proposal number start ---------------")

	if policy.ProposalNumber != 0 {
		log.Printf("proposal number already set %d", policy.ProposalNumber)
		return
	}

	log.Println("setting proposal number...")
	firePolicy := lib.PolicyCollection
	policy.ProposalNumber = GetSequenceProposal("", firePolicy)
	log.Printf("proposal number %d", policy.ProposalNumber)
}

func GetSequenceProposal(name string, firePolicy string) int {
	var number int
	r, e := lib.OrderLimitFirestoreErr(firePolicy, "proposalNumber", firestore.Desc, 1)
	lib.CheckError(e)
	policyCompany := models.PolicyToListData(r)
	if len(policyCompany) == 0 {
		number = 1
	} else {

		number = policyCompany[0].ProposalNumber + 1
	}
	log.Println("GetSequenceProposal: ", number)
	return number
}
