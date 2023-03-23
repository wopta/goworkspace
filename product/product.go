package product

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	lib "github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	//"google.golang.org/api/firebaseappcheck/v1"
)

func init() {
	log.Println("INIT Product")
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
				Route:   "/v1/name/:name",
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

func GetNameFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	name := r.Header.Get("name")
	v := strings.Split(r.RequestURI, "/")
	version := v[1]
	log.Println(r.RequestURI)
	log.Println(v)
	log.Println(v[1])
	product, e := GetName(name, version)
	jsonString, e := product.Marshal()
	out := string(jsonString)
	if name == "persona" {
		initialDate := time.Now().AddDate(-18, 0, 0)
		minDate := time.Now().AddDate(-74, 0, 1)
		out = strings.Replace(out, "{{INITIAL_DATE}}", initialDate.Format(time.DateOnly), 2)
		out = strings.Replace(out, "{{MIN_DATE}}", minDate.Format(time.DateOnly), 1)

		err := json.Unmarshal([]byte(out), &product)
		if err != nil {
			return "", nil, err
		}
	}
	return out, product, e

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
	query := q.FirestoreWherefields("products")
	products := models.ProductToListData(query)

	return products[0], nil

}

func PutFx(resp http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println(productCollection)
	request := lib.ErrorByte(io.ReadAll(r.Body))
	pr, e := models.UnmarshalProduct([]byte(request))
	p, e := Put(pr)
	return "{}", p, e
}
func Put(p models.Product) (models.Product, error) {

	r, _, e := lib.PutFirestoreErr("products", p)
	log.Println(r.ID)

	return p, e
}
