package productFlow

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"

	"gitlab.dev.wopta.it/goworkspace/bpmn/bpmnEngine"
	"gitlab.dev.wopta.it/goworkspace/bpmn/bpmnEngine/flow"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models/catnat"
)

const nameCatnatDocument = "Contratto NetInsurance"

func AddCatnatHandlers(builder *bpmnEngine.BpnmBuilder) error {
	return bpmnEngine.IsError(
		builder.AddHandler("catnatIntegration", catnatIntegration),
		builder.AddHandler("catnatdownloadPolicy", catnatDownloadCertification),
		builder.AddHandler("catnatUploadDocument", catnatUpload),
	)
}

func catnatUpload(store *bpmnEngine.StorageBpnm) error {
	policy, err := bpmnEngine.GetData[*flow.Policy]("policy", store)
	if err != nil {
		return err
	}
	client := catnat.NewNetClient()
	attachmentName := fmt.Sprint(nameCatnatDocument)
	var document string
	for i := range *policy.Attachments {
		if (*policy.Attachments)[i].FileName == attachmentName {
			bytes, err := lib.ReadFileFromGoogleStorageEitherGsOrNot((*policy.Attachments)[i].Link)
			if err != nil {
				return err
			}
			document = base64.StdEncoding.EncodeToString(bytes)
			break
		}
	}
	if document == "" {
		return errors.New("Document not found")
	}
	return client.UploadDocument(*policy.Policy, document)
}
func catnatIntegration(store *bpmnEngine.StorageBpnm) error {
	policy, err := bpmnEngine.GetData[*flow.Policy]("policy", store)
	if err != nil {
		return err
	}
	client := catnat.NewNetClient()
	requestDto := catnat.QuoteRequest{}
	err = requestDto.FromPolicyForEmit(policy.Policy)
	if err != nil {
		return err
	}
	res, err := client.Emit(requestDto, policy.Policy)
	if err != nil {
		return err
	}
	policy.CodeCompany = res.PolicyNumber
	store.AddLocal("numeroPolizza", &flow.String{String: res.PolicyNumber})
	return nil
}

func catnatDownloadCertification(store *bpmnEngine.StorageBpnm) error {
	var policy *flow.Policy
	var numeroPoliza *flow.String
	err := bpmnEngine.IsError(
		bpmnEngine.GetDataRef("policy", &policy, store),
		bpmnEngine.GetDataRef("numeroPolizza", &numeroPoliza, store),
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
	filePath := fmt.Sprint("temp/", policy.Uid, "/namirial/", policy.NameDesc, " ", nameCatnatDocument, " ", policy.CodeCompany)
	_, err = lib.PutToStorageErr(os.Getenv("GOOGLE_STORAGE_BUCKET"), filePath, bytes)
	if err != nil {
		return err
	}

	return err
}
