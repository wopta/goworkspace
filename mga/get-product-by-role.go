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
		resp    models.Product
		request GetProductByRoleRequest
		err     error
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
	log.Printf("GetProductByRoleFx response: %s", string(jsonResp))

	return string(jsonResp), resp, err
}
