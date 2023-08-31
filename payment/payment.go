package payment

import (
	"fmt"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	prd "github.com/wopta/goworkspace/product"
	"log"
	"net/http"
)

func init() {
	log.Println("INIT Payment")
	functions.HTTP("Payment", Payment)
}

func Payment(w http.ResponseWriter, r *http.Request) {
	log.Println("Payment")
	lib.EnableCors(&w, r)
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	route := lib.RouteData{
		Routes: []lib.Route{

			{
				Route:   "/v1/fabrick",
				Handler: FabrickPay,
				Method:  http.MethodPost,
				Roles:   []string{models.UserRoleAll},
			},
			{
				Route:   "/v1/fabrick/montly",
				Handler: FabrickPayMontly,
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

func PaymentController(origin string, policy models.Policy) (string, error) {
	var (
		payUrl, paymentProvider string
		paymentMethods          []string
	)

	log.Printf("[PaymentController] init")

	paymentProvider = policy.Payment
	paymentMethods = getPaymentMethods(policy)

	log.Printf("[PaymentController] genereting payment URL")
	switch paymentProvider {
	case models.FabrickPaymentProvider:
		var payRes FabrickPaymentResponse

		if policy.PaymentSplit == string(models.PaySplitYear) {
			payRes = FabbrickYearPay(policy, origin, paymentMethods)
		}
		if policy.PaymentSplit == string(models.PaySplitMonthly) {
			payRes = FabbrickMontlyPay(policy, origin, paymentMethods)
		}
		if payRes.Payload.PaymentPageURL == nil {
			return "", fmt.Errorf("fabrick error: %v", payRes.Errors)
		}
		payUrl = *payRes.Payload.PaymentPageURL
	default:
		return "", fmt.Errorf("payment provider %s not supported", policy.Payment)
	}
	return payUrl, nil
}

func getPaymentMethods(policy models.Policy) []string {
	paymentMethods := make([]string, 0)

	log.Printf("[GetPaymentMethods] loading available payment methods for %s payment provider", policy.Payment)

	product, err := prd.GetProduct(policy.Name, policy.ProductVersion, models.UserRoleAdmin)
	lib.CheckError(err)

	for _, provider := range product.PaymentProviders {
		if provider.Name == policy.Payment {
			for _, method := range provider.Methods {
				if lib.SliceContains(method.Rates, policy.PaymentSplit) {
					paymentMethods = append(paymentMethods, method.Name)
				}
			}
		}
	}
	return paymentMethods
}
