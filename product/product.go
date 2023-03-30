package product

import (
	"encoding/json"
	"fmt"
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
	product, err := GetName(name, version)
	if err != nil {

		return "", nil, err
	}
	jsonOut, err := product.Marshal()
	if err != nil {
		return "", nil, err
	}

	jsonString := string(jsonOut)
	switch name {
	case "persona":
		jsonString, product, err = ReplaceDatesInProduct(product, 75)
	case "life":
		jsonString, product, err = ReplaceDatesInProduct(product, 69)
	}

	return jsonString, product, err
}

func ReplaceDatesInProduct(product models.Product, minYear int) (string, models.Product, error) {
	jsonOut, err := product.Marshal()
	if err != nil {
		return "", models.Product{}, err
	}

	productJson := string(jsonOut)

	initialDate := time.Now().AddDate(-18, 0, 0).Format("2006-01-02")
	minDate := time.Now().AddDate(-minYear, 0, 1).Format("2006-01-02")

	regexInitialDate := regexp.MustCompile("{{INITIAL_DATE}}")
	regexMinDate := regexp.MustCompile("{{MIN_DATE}}")

	productJson = regexInitialDate.ReplaceAllString(productJson, initialDate)
	productJson = regexMinDate.ReplaceAllString(productJson, minDate)

	err = json.Unmarshal(jsonOut, &product)
	if err != nil {
		return "", models.Product{}, err
	}
	return productJson, models.Product{}, err
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
	if len(products) == 0 {
		return models.Product{}, fmt.Errorf("no product json file found for %s %s", name, version)
	}

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
