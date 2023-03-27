package product

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"regexp"
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
		e = replaceDatesForPersonaProduct(&out, &product)
	}
	return out, product, e
}

func replaceDatesForPersonaProduct(productJson *string, product *models.Product) error {
	initialDate := time.Now().AddDate(-18, 0, 0).Format("2006-01-02")
	minDate := time.Now().AddDate(-75, 0, 1).Format("2006-01-02")

	regexInitialDate := regexp.MustCompile("{{INITIAL_DATE}}")
	regexMinDate := regexp.MustCompile("{{MIN_DATE}}")

	*productJson = regexInitialDate.ReplaceAllString(*productJson, initialDate)
	*productJson = regexMinDate.ReplaceAllString(*productJson, minDate)

	err := json.Unmarshal([]byte(*productJson), product)
	if err != nil {
		return err
	}
	return nil
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
