package document

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"gitlab.dev.wopta.it/goworkspace/document/namirial"
	lib "gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
)

func SignNamirial(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {

	req := lib.ErrorByte(ioutil.ReadAll(r.Body))
	var data models.Policy
	defer r.Body.Close()
	err := json.Unmarshal([]byte(req), &data)
	lib.CheckError(err)

	input := namirial.NamirialInput{
		DocumentsFullPath: []string{data.DocumentName},
		Policy:            data,
	}
	output, err := namirial.Sign(input)
	return fmt.Sprintf("%+v", output), output, err

}
