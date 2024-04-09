package product

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
	//"google.golang.org/api/firebaseappcheck/v1"
)

var productRoutes []lib.ChiRoute = []lib.ChiRoute{
	{
		Route:   "/v1/{name}",
		Handler: lib.ResponseLoggerWrapper(GetNameFx),
		Method:  http.MethodGet,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/v1/name/{name}",
		Handler: lib.ResponseLoggerWrapper(GetNameFx),
		Method:  http.MethodGet,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/v1",
		Handler: lib.ResponseLoggerWrapper(PutFx),
		Method:  http.MethodPut,
		Roles:   []string{lib.UserRoleAll},
	},
}

func init() {
	log.Println("INIT Product")
	functions.HTTP("Product", Product)
}

func Product(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix)

	router := lib.GetChiRouter("product", productRoutes)
	router.ServeHTTP(w, r)
}
