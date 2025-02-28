package product

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

var productRoutes []lib.Route = []lib.Route{
	{
		Route:       "/v1/{name}",
		Fn:          GetNameFx,
		Method:      http.MethodGet,
		Roles:       []string{lib.UserRoleAll},
		Entitlement: "product.get.product",
	},
	{
		Route:       "/v1/name/{name}",
		Fn:          GetNameFx,
		Method:      http.MethodGet,
		Roles:       []string{lib.UserRoleAll},
		Entitlement: "product.get.product",
	},
	{
		Route:       "/v1",
		Fn:          PutFx,
		Method:      http.MethodPut,
		Roles:       []string{lib.UserRoleAll},
		Entitlement: "product.update.product",
	},
}

func init() {
	log.Println("INIT Product")
	functions.HTTP("Product", Product)
}

func Product(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix)

	router := models.GetExtendedRouter("product", productRoutes)
	router.ServeHTTP(w, r)
}
