package product

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	lib "github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	//"google.golang.org/api/firebaseappcheck/v1"
)

func init() {
	log.Println("INIT AppcheckProxy")
	functions.HTTP("Product", Product)
}

func Product(w http.ResponseWriter, r *http.Request) {

	log.Println("Product")
	lib.EnableCors(&w, r)
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	route := lib.RouteData{
		Routes: []lib.Route{
			{
				Route:   "/v1/:name",
				Handler: GetNameFx,
				Method:  "GET",
			},
			{
				Route:   "/v1",
				Handler: PutFx,
				Method:  "PUT",
			},
		},
	}
	route.Router(w, r)

}

const (
	productCollection = "products"
)

func GetNameFx(resp http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	name := r.Header.Get("name")
	v := strings.Split(r.RequestURI, "/")
	version := v[0]
	log.Println(v[0])
	product, e := GetName(name, version)
	jsonString, e := product.Marshal()
	return string(jsonString), product, e

}
func GetName(name string, version string) (models.Product, error) {
	q := lib.Firequeries{
		Queries: []lib.Firequery{{
			Field:      "name",
			Operator:   "==",
			QueryValue: name,
		},
			{
				Field:      "version",
				Operator:   "==",
				QueryValue: version,
			},
		},
	}
	query := q.FirestoreWherefields("product")
	products := models.ProductToListData(query)

	return products[0], nil

}

func PutFx(resp http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println(productCollection)
	request := lib.ErrorByte(ioutil.ReadAll(r.Body))
	pr, e := models.UnmarshalProduct([]byte(request))
	p, e := Put(pr)
	return "{}", p, e
}
func Put(p models.Product) (models.Product, error) {

	r, _, e := lib.PutFirestoreErr("products", p)
	log.Println(r.ID)

	return p, e
}
