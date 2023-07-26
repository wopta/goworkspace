package policy

import (
	"errors"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func FillAttachments(policy *models.Policy, origin string) error {
	firePolicy := lib.GetDatasetByEnv(origin, models.PolicyCollection)

	if policy.Attachments == nil {
		policy.Attachments = new([]models.Attachment)
	}
	policy.Updated = time.Now().UTC()

	return lib.SetFirestoreErr(firePolicy, policy.Uid, policy)
}

func Sign(policy *models.Policy, origin string) error {
	if policy.Status != models.PolicyStatusToSign {
		return errors.New("policy wrong status")
	}

	firePolicy := lib.GetDatasetByEnv(origin, models.PolicyCollection)

	policy.IsSign = true
	policy.Status = models.PolicyStatusSign
	policy.StatusHistory = append(policy.StatusHistory, models.PolicyStatusSign)
	policy.Updated = time.Now().UTC()

	return lib.SetFirestoreErr(firePolicy, policy.Uid, policy)
}

func SetToPay(policy *models.Policy, origin string) error {
	firePolicy := lib.GetDatasetByEnv(origin, models.PolicyCollection)

	policy.Status = models.PolicyStatusToPay
	policy.StatusHistory = append(policy.StatusHistory, models.PolicyStatusToPay)
	policy.Updated = time.Now().UTC()

	return lib.SetFirestoreErr(firePolicy, policy.Uid, policy)
}
