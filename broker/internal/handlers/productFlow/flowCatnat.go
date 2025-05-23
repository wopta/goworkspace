package productFlow

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	bpmn "gitlab.dev.wopta.it/goworkspace/broker/draftBpmn"
	"gitlab.dev.wopta.it/goworkspace/broker/draftBpmn/flow"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/quote/catnat"
)

func CatnatIntegration(store bpmn.StorageData) error {
	policy, err := bpmn.GetData[*flow.PolicyDraft]("policy", store)
	if err != nil {
		return err
	}
	client := catnat.NewNetClient()
	requestDto := catnat.QuoteRequest{}
	err = requestDto.FromPolicy(policy.Policy, true)
	if err != nil {
		return err
	}
	res, e := client.Emit(requestDto)
	if e != nil {
		return e
	}
	catnat.MappingQuoteResponseToPolicy(res, policy.Policy)
	store.AddLocal("numeroPolizza", &flow.StringBpmn{String: res.PolicyNumber})
	return nil
}

func CatnatDownloadCertification(store bpmn.StorageData) error {
	var policy *flow.PolicyDraft
	var numeroPoliza *flow.StringBpmn
	err := bpmn.IsError(
		bpmn.GetDataRef("policy", &policy, store),
		bpmn.GetDataRef("numeroPolizza", &numeroPoliza, store),
	)
	if err != nil {
		return err
	}

	client := catnat.NewNetClient()
	resp, err := client.Download(numeroPoliza.String)
	if err != nil {
		return err
	}
	bytes, err := base64.StdEncoding.DecodeString(resp.Documento[0].DatiDocumento)
	if err != nil {
		return err
	}
	filePath := strings.ReplaceAll(fmt.Sprintf("%s/%s/%s", "temp", policy.Uid, fmt.Sprintf(models.NetInsuranceDocument, policy.NameDesc)), " ", "_")
	link, err := lib.PutToStorageErr(os.Getenv("GOOGLE_STORAGE_BUCKET"), filePath, bytes)
	if err != nil {
		return err
	}
	if policy.Attachments == nil {
		policy.Attachments = new([]models.Attachment)
	}
	*policy.Attachments = append(*policy.Attachments, models.Attachment{
		Name:     models.ContractAttachmentName + " NetInsurance",
		FileName: filePath,
		Link:     link,
		MimeType: "application/pdf",
	})
	return err
}
