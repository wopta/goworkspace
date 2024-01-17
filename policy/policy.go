package policy

import (
	"log"
	"net/http"
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
				Handler: DeletePolicyFx,
				Method:  http.MethodDelete,
				Roles:   []string{models.UserRoleAdmin, models.UserRoleManager},
			},
			{
				Route:   "attachment/v1/:uid",
				Handler: GetPolicyAttachmentsFx,
				Method:  http.MethodGet,
				Roles:   []string{models.UserRoleAll},
			},
			{
				Route:   "/portfolio/v1",
				Handler: GetPortfolioPoliciesFx,
				Method:  http.MethodPost,
				Roles:   []string{models.UserRoleAgent, models.UserRoleAgency},
			},
			{
				Route:   "/media/v1",
				Handler: GetPolicyMediaFx,
				Method:  http.MethodPost,
				Roles:   []string{models.UserRoleAdmin, models.UserRoleManager, models.UserRoleAgent, models.UserRoleAgency},
			},
			{
				Route:   "/v1",
				Handler: GetPoliciesByQueryFx,
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

func SetPolicyPaid(policy *models.Policy, origin string) {
	firePolicy := lib.GetDatasetByEnv(origin, models.PolicyCollection)

	// Update payment fields
	policy.IsPay = true
	policy.Updated = time.Now().UTC()
	policy.Status = models.PolicyStatusPay
	policy.StatusHistory = append(policy.StatusHistory, models.PolicyStatusPay)
	lib.SetFirestore(firePolicy, policy.Uid, policy)
	policy.BigquerySave(origin)
}
