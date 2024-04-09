package form

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
)

var formRoutes []lib.Route = []lib.Route{
	{
		Route:   "/axafleet",
		Handler: lib.ResponseLoggerWrapper(AxaFleetTway),
		Method:  http.MethodGet,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/v1/{uid}",
		Handler: lib.ResponseLoggerWrapper(GetFx),
		Method:  http.MethodGet,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/fleet/assistance/v1",
		Handler: lib.ResponseLoggerWrapper(GetFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAll},
	},
}

func init() {
	log.Println("INIT Form")
	functions.HTTP("Form", Form)
}

func Form(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix)

	router := lib.GetRouter("form", formRoutes)
	router.ServeHTTP(w, r)
}
