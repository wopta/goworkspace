package broker

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	bpmn "github.com/wopta/goworkspace/broker/draftBpmn"
	"github.com/wopta/goworkspace/broker/draftBpmn/flow"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/quote/catnat"
)

func getProductFlow() (*bpmn.BpnmBuilder, error) {
	store := bpmn.NewStorageBpnm()
	builder, e := bpmn.NewBpnmBuilder("flows/draft/product_flows.json")
	if e != nil {
		return nil, e
	}
	builder.SetStorage(store)
	err := bpmn.IsError(
		builder.AddHandler("catnatIntegration", catnatIntegration),
		builder.AddHandler("catnatdownloadPolicy", catnatDownloadCertification),
	)
	if err != nil {
		return nil, err
	}
	return builder, nil
}

func catnatIntegration(store bpmn.StorageData) error {
	policy, err := bpmn.GetData[*flow.PolicyDraft]("policy", store)
	if err != nil {
		return err
	}
	client := catnat.NewNetClient()
	requestDto := catnat.RequestDTO{}
	err = requestDto.FromPolicy(policy.Policy, true)
	if err != nil {
		return err
	}
	res, e := client.Emit(requestDto)
	if e != nil {
		return e
	}
	store.AddLocal("numeroPoliza", &flow.StringBpmn{String: res.PolicyNumber})
	return nil
}

func catnatDownloadCertification(store bpmn.StorageData) error {
	var policy *flow.PolicyDraft
	var numeroPoliza *flow.StringBpmn
	err := bpmn.IsError(
		bpmn.GetDataRef("policy", &policy, store),
		bpmn.GetDataRef("numeroPoliza", &numeroPoliza, store),
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
	filename := fmt.Sprintf(models.NetInsuranceDocument, policy.NameDesc)
	filePath := strings.ReplaceAll(fmt.Sprintf("%s/%s/%s", "temp", policy.Uid, filename), " ", "_")
	link, err := lib.PutToStorageErr(os.Getenv("GOOGLE_STORAGE_BUCKET"), filePath, bytes)
	if err != nil {
		return err
	}
	if policy.Attachments == nil {
		policy.Attachments = new([]models.Attachment)
	}
	*policy.Attachments = append(*policy.Attachments, models.Attachment{
		Name:     models.ContractAttachmentName + "NetInsurance",
		FileName: filename + ".pdf",
		Link:     link,
		MimeType: "application/pdf",
	})
	return err
}
