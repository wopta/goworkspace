package mga

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/product"
)

type GetProductByRoleRequest struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Company string `json:"company"`
}

func GetProductByRoleFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println("GetProductByRoleFx")
	var (
		resp       models.Product
		respString string
		request    GetProductByRoleRequest
		err        error
	)

	body := lib.ErrorByte(io.ReadAll(r.Body))
	err = json.Unmarshal(body, &request)
	lib.CheckError(err)
	log.Printf("GetProductByRoleFx body: %s", string(body))

	authToken, err := models.GetAuthTokenFromIdToken(r.Header.Get("Authorization"))
	lib.CheckError(err)
	log.Printf("GetProductByRoleFx authToken: %s", authToken)

	resp, err = product.GetProductByRole(request.Name, request.Version, request.Company, authToken)
	if err != nil {
		return "", resp, err
	}
	jsonResp, err := json.Marshal(resp)

	respString = string(jsonResp)
	switch request.Name {
	case "persona":
		respString, resp, err = product.ReplaceDatesInProduct(resp, 70, 0)
	case "life":
		respString, resp, err = product.ReplaceDatesInProduct(resp, 70, 55)
	}

	log.Printf("GetProductByRoleFx response: %s", respString)
	return respString, resp, err
}
