package product

import (
	"io/ioutil"
	"log"
	"net/http"

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
				Route:   "/v1/name/:name",
				Handler: GetNameFx,
				Method:  "GET",
			},
			{
				Route:   "/v1/:uid",
				Handler: GetFx,
				Method:  "GET",
			},

			{
				Route:   "/v1",
				Handler: GetNameFx,
				Method:  "POST",
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

func GetFx(resp http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println(r.Header.Get("uid"))
	p, e := Get(r.Header.Get("uid"))
	jsonString, e := p.Marshal()
	return string(jsonString), p, e
}
func Get(uid string) (models.Product, error) {
	log.Println(uid)
	productFire := lib.GetFirestore("products", uid)
	var product models.Product
	e := productFire.DataTo(product)
	return product, e
}

func GetNameFx(resp http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	name := r.Header.Get("name")
	log.Println(name)
	product, e := GetName(name)
	jsonString, e := product.Marshal()
	return string(jsonString), product, e

}
func GetName(name string) (models.Product, error) {

	productFire, e := lib.QueryWhereFirestore("products", "name", "==", name)

	products := models.ProductToListData(productFire)

	return products[0], e

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
