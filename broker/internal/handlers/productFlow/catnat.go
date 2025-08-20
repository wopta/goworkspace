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
	"gitlab.dev.wopta.it/goworkspace/models/catnat"
)

func CatnatIntegration(store bpmn.StorageData) error {
	policy, err := bpmn.GetData[*flow.Policy]("policy", store)
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

func CatnatDownloadCertification(store bpmn.StorageData) error {
	var policy *flow.Policy
	var numeroPoliza *flow.String
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
	filePath := strings.ReplaceAll(fmt.Sprintf("%s/%s/namirial/%s %s", "temp", policy.Uid, models.ContractAttachmentName, "NetInsurance"), " ", "_")
	_, err = lib.PutToStorageErr(os.Getenv("GOOGLE_STORAGE_BUCKET"), filePath, bytes)
	if err != nil {
		return err
	}

	return err
}
