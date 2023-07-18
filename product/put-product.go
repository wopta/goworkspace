package product

import (
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"io"
	"log"
	"net/http"
)

func PutFx(resp http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println(models.ProductsCollection)
	origin := r.Header.Get("origin")
	request := lib.ErrorByte(io.ReadAll(r.Body))
	pr, e := models.UnmarshalProduct([]byte(request))
	p, e := Put(origin, pr)
	return "{}", p, e
}

func Put(origin string, p models.Product) (models.Product, error) {
	fireProducts := lib.GetDatasetByEnv(origin, "products")
	r, _, e := lib.PutFirestoreErr(fireProducts, p)
	log.Println(r.ID)

	return p, e
}
