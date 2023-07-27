package policy

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func init() {
	log.Println("INIT Policy")
	functions.HTTP("Policy", Policy)
}

func Policy(w http.ResponseWriter, r *http.Request) {
	log.Println("Policy")
	lib.EnableCors(&w, r)
	route := lib.RouteData{
		Routes: []lib.Route{
			{
				Route:   "fiscalcode/v1/:fiscalcode",
				Handler: GetPolicyByFiscalCodeFx, // Broker.PolicyFiscalcode
				Method:  "GET",
				Roles:   []string{models.UserRoleAll},
			},
			{
				Route:   "v1/:uid",
				Handler: GetPolicyFx,
				Method:  "GET",
				Roles:   []string{models.UserRoleAll},
			},
			{
				Route:   "v1/:uid",
				Handler: PatchPolicyFx, // Broker.UpdatePolicy
				Method:  http.MethodPatch,
				Roles:   []string{models.UserRoleAll},
			},
			{
				Route:   "v1/:uid",
				Handler: DeletePolicyFx, // Broker.DeletePolicy
				Method:  http.MethodDelete,
				Roles:   []string{models.UserRoleAdmin, models.UserRoleManager},
			},
			{
				Route:   "attachment/v1/:uid",
				Handler: GetPolicyAttachmentFx,
				Method:  http.MethodGet,
				Roles:   []string{models.UserRoleAll},
			},
			{
				Route:   "v1",
				Handler: GetPoliciesByQueryFx, // Broker.GetPoliciesFx,
				Method:  http.MethodPost,
				Roles:   []string{models.UserRoleAdmin, models.UserRoleManager},
			},
		},
	}
	route.Router(w, r)
}

func GetPolicyByUid(policyUid string, origin string) models.Policy {
	firePolicy := lib.GetDatasetByEnv(origin, "policy")
	policyF := lib.GetFirestore(firePolicy, policyUid)
	var policy models.Policy
	policyF.DataTo(&policy)
	policyM, _ := policy.Marshal()
	log.Println("GetPolicyByUid: Policy "+policyUid+" found: ", string(policyM))

	return policy
}

func SetPolicyPaid(policy *models.Policy, contractLink string, origin string) {
	firePolicy := lib.GetDatasetByEnv(origin, "policy")
	now := time.Now().UTC()
	// Add Contract
	timestamp := strconv.FormatInt(now.Unix(), 10)
	*policy.Attachments = append(*policy.Attachments, models.Attachment{
		Name: "Contratto",
		Link: contractLink,
		FileName: "Contratto_" + strings.ReplaceAll(policy.NameDesc, " ", "_") +
			"_" + timestamp + ".pdf",
	})
	// Update payment fields
	policy.IsPay = true
	policy.Updated = now
	policy.Status = models.PolicyStatusPay
	policy.StatusHistory = append(policy.StatusHistory, models.PolicyStatusPay)
	lib.SetFirestore(firePolicy, policy.Uid, policy)
	policy.BigquerySave(origin)
}
