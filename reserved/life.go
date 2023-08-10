package reserved

import (
	"encoding/json"
	"fmt"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"io"
	"log"
	"net/http"
	"strings"
)

type Contact struct {
	ContactType string `json:"contactType"`
	Address     string `json:"address"`
	Object      string `json:"object,omitempty"`
}

type LifeReservedResp struct {
	Documents []string  `json:"documents"`
	Contacts  []Contact `json:"contacts"`
}

func LifeReservedFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	const (
		rulesFileName = "life-reserved.json"
	)
	var (
		policy models.Policy
		err    error
	)

	log.Println("[LifeReserved]")

	err = json.Unmarshal(lib.ErrorByte(io.ReadAll(r.Body)), &policy)
	lib.CheckError(err)

	resp := &LifeReservedResp{
		Documents: make([]string, 0),
		Contacts:  getContactsDetails(policy),
	}

	fx := new(models.Fx)
	rulesFile := lib.GetRulesFile(rulesFileName)

	_, ruleOutput := lib.RulesFromJsonV2(fx, rulesFile, resp, getInputData(policy), nil)

	resp = ruleOutput.(*LifeReservedResp)
	jsonOut, err := json.Marshal(resp)

	return string(jsonOut), resp, err
}

func getContactsDetails(policy models.Policy) []Contact {
	return []Contact{
		{
			ContactType: "e-mail",
			Address:     "clp.it.sinistri@partners.axa",
			Object: fmt.Sprintf("%s proposta %d - UNDERWRITING MEDICO - %s", policy.NameDesc, policy.ProposalNumber,
				strings.ToUpper(policy.Contractor.Surname+" "+policy.Contractor.Name)),
		},
		{
			ContactType: "posta",
			Address:     "AXA PARTNERS Ufficio Underwriting Medico – Corso Como n. 17 – 20154 MILANO",
		},
	}
}

func getInputData(policy models.Policy) []byte {
	var err error

	in := make(map[string]interface{})
	in["gender"] = policy.Contractor.Gender
	in["age"], err = policy.CalculateContractorAge()
	lib.CheckError(err)
	maxSumInsuredLimitOfIndemnity := 0.0
	for _, guarantee := range policy.Assets[0].Guarantees {
		if guarantee.Value.SumInsuredLimitOfIndemnity > maxSumInsuredLimitOfIndemnity {
			maxSumInsuredLimitOfIndemnity = guarantee.Value.SumInsuredLimitOfIndemnity
		}
	}
	in["sumInsuredLimitOfIndemnity"] = maxSumInsuredLimitOfIndemnity

	out, err := json.Marshal(in)
	lib.CheckError(err)

	return out
}
