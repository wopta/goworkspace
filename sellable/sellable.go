package sellable

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
)

const (
	rulesFilename = "sellable"
)

var sellableRoutes []lib.Route = []lib.Route{
	{
		Route:   "/v1/sales/life",
		Handler: lib.ResponseLoggerWrapper(LifeFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/v1/risk/person",
		Handler: lib.ResponseLoggerWrapper(PersonaFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/v1/commercial-combined",
		Handler: lib.ResponseLoggerWrapper(CommercialCombinedFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/v1/catnat/{step}",
		Handler: lib.ResponseLoggerWrapper(CatnatFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAll},
	},
}

func init() {
	log.Println("INIT Sellable")
	functions.HTTP("Sellable", Sellable)
}

func Sellable(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix)

	router := lib.GetRouter("sellable", sellableRoutes)
	router.ServeHTTP(w, r)
}
