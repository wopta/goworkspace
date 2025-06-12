package utility

import (
	"strings"
	"time"

	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/models"
	plc "gitlab.dev.wopta.it/goworkspace/policy"
	"gitlab.dev.wopta.it/goworkspace/reserved"
)

func SetRequestApprovalData(policy *models.Policy, networkNode *models.NetworkNode, mgaProduct *models.Product, origin string) {
	log.AddPrefix("setRequestApprovalData")
	defer log.PopPrefix()
	log.Printf("policy uid %s: reserved flow", policy.Uid)

	SetProposalNumber(policy)

	if policy.Status == models.PolicyStatusInitLead {
		plc.AddProposalDoc(origin, policy, networkNode, mgaProduct)
		for _, reason := range policy.ReservedInfo.Reasons {
			// TODO: add key/id for reasons so we do not have to check string equallity
			if !strings.HasPrefix(reason, "Cliente già assicurato") {
				reserved.SetReservedInfo(policy)
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
