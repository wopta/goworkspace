package utility

import (
	"slices"
	"time"

	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/models"
)

func SetLeadData(policy *models.Policy, product models.Product, networkNode *models.NetworkNode) error {
	log.AddPrefix("SetLeadData")
	defer log.PopPrefix()
	log.Println("start -----------------------------------------")
	now := time.Now().UTC()

	if policy.CreationDate.IsZero() {
		policy.CreationDate = now
	}
	if policy.Status != models.PolicyStatusInitLead {
		policy.Status = models.PolicyStatusInitLead
		policy.StatusHistory = append(policy.StatusHistory, policy.Status)
	}
	log.Printf("policy status %s", policy.Status)

	policy.IsSign = false
	policy.IsPay = false
	policy.Updated = now

	if networkNode != nil {
		setPolicyProducerNode(policy, networkNode)
	}

	// TODO delete me when PMI is fixed
	if policy.Name == models.PmiProduct {
		policy.NameDesc = "Wopta per te Artigiani & Imprese"
	}
	if policy.ProductVersion == "" {
		policy.ProductVersion = "v1"
	}

	setRenewInfo(policy, product)
	linkDocument, err := lib.GetLastVersionSetInformativo(policy.Name, policy.ProductVersion)
	if err != nil {
		return err
	}
	log.Println("add information set")
	informationSet := models.Attachment{
		Name:        "Set Informativo Precontrattuale",
		FileName:    "Set Informativo Precontrattuale.pdf",
		MimeType:    lib.GetContentType("pdf"),
		ContentType: lib.GetContentType("pdf"),
		Link:        "gs://" + linkDocument,
	}
	if policy.Attachments == nil {
		policy.Attachments = new([]models.Attachment)
	}
	attIdx := slices.IndexFunc(*policy.Attachments, func(a models.Attachment) bool {
		return a.Name == informationSet.Name
	})
	if attIdx == -1 {
		*policy.Attachments = append(*policy.Attachments, informationSet)
	}

	log.Println("end -------------------------------------------")
	return nil
}

func setPolicyProducerNode(policy *models.Policy, node *models.NetworkNode) {
	policy.ProducerUid = node.Uid
	policy.ProducerCode = node.Code
	policy.ProducerType = node.Type
	policy.NetworkUid = node.NetworkUid
}

func setRenewInfo(policy *models.Policy, product models.Product) {
	policy.Annuity = 0
	policy.IsRenewable = product.IsRenewable
	policy.IsAutoRenew = product.IsAutoRenew
	policy.PolicyType = product.PolicyType
	policy.QuoteType = product.QuoteType
}
