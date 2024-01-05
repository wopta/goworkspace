package payment

import (
	"fmt"
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func init() {
	log.Println("INIT Payment")
	functions.HTTP("Payment", Payment)
}

func Payment(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix)

	log.Println("Payment")
	lib.EnableCors(&w, r)
	w.Header().Set("Access-Control-Allow-Methods", "POST")

	route := lib.RouteData{
		Routes: []lib.Route{
			{
				Route:   "/v1/fabrick/recreate",
				Handler: FabrickRefreshPayByLinkFx,
				Method:  http.MethodPost,
				Roles:   []string{models.UserRoleAdmin, models.UserRoleManager},
			},
			{
				Route:   "/v1/fabrick",
				Handler: FabrickPayFx,
				Method:  http.MethodPost,
				Roles:   []string{models.UserRoleAll},
			},
			{
				Route:   "/v1/fabrick/montly",
				Handler: FabrickPayMonthlyFx,
				Method:  http.MethodPost,
				Roles:   []string{models.UserRoleAll},
			},
			{
				Route:   "/v1/cripto",
				Handler: CriptoPay,
				Method:  http.MethodPost,
				Roles:   []string{models.UserRoleAll},
			},
			{
				Route:   "/v1/:uid",
				Handler: FabrickExpireBill,
				Method:  http.MethodDelete,
				Roles:   []string{models.UserRoleAdmin, models.UserRoleManager},
			},
			{
				Route:   "/manual/v1/:transactionUid",
				Handler: ManualPaymentFx,
				Method:  http.MethodPost,
				Roles:   []string{models.UserRoleAdmin, models.UserRoleManager},
			},
		},
	}
	route.Router(w, r)

}

func CriptoPay(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	return "", nil, nil
}

func PaymentController(origin string, policy *models.Policy, product, mgaProduct *models.Product) (string, error) {
	var (
		payUrl         string
		paymentMethods []string
	)

	log.Printf("[PaymentController] init")

	// TODO: fix me
	if policy.Payment == "" || policy.Payment == "fabrik" {
		policy.Payment = models.FabrickPaymentProvider
	}
	paymentMethods = getPaymentMethods(*policy, product)

	log.Printf("[PaymentController] generating payment URL")
	switch policy.Payment {
	case models.FabrickPaymentProvider:
		var payRes FabrickPaymentResponse

		switch policy.PaymentSplit {
		case string(models.PaySplitYear), string(models.PaySplitYearly), string(models.PaySplitSingleInstallment):
			log.Printf("[PaymentController] fabrick yearly pay")
			payRes = FabrickYearPay(*policy, origin, paymentMethods, mgaProduct)
		case string(models.PaySplitMonthly):
			log.Printf("[PaymentController] fabrick monthly pay")
			payRes = FabrickMonthlyPay(*policy, origin, paymentMethods, mgaProduct)
		}
		if payRes.Payload == nil || payRes.Payload.PaymentPageURL == nil {
			log.Println("[PaymentController] fabrick error payload or paymentUrl empty")
			return "", fmt.Errorf("fabrick error: %v", payRes.Errors)
		}
		payUrl = *payRes.Payload.PaymentPageURL
	default:
		return "", fmt.Errorf("payment provider %s not supported", policy.Payment)
	}

	log.Printf("[PaymentController] payUrl: %s", payUrl)

	return payUrl, nil
}

func getPaymentMethods(policy models.Policy, product *models.Product) []string {
	var paymentMethods = make([]string, 0)

	log.Printf("[GetPaymentMethods] loading available payment methods for %s payment provider", policy.Payment)

	// TODO: remove me once established standard
	if policy.PaymentSplit == string(models.PaySplitYear) {
		policy.PaymentSplit = string(models.PaySplitYearly)
	}

	for _, provider := range product.PaymentProviders {
		if provider.Name == policy.Payment {
			for _, method := range provider.Methods {
				if lib.SliceContains(method.Rates, policy.PaymentSplit) {
					paymentMethods = append(paymentMethods, method.Name)
				}
			}
		}
	}

	log.Printf("[GetPaymentMethods] found %v", paymentMethods)
	return paymentMethods
}
