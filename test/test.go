package test

/*

 */
import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
)

var (
	signatureID = 0
)

func init() {
	log.Println("INIT Test")
	functions.HTTP("Test", Test)
}

func Test(w http.ResponseWriter, r *http.Request) {
	log.Println("Test")
	lib.EnableCors(&w, r)
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	route := lib.RouteData{
		Routes: []lib.Route{},
	}
	route.Router(w, r)

}
