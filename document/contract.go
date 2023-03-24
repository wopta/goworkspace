package document

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	lib "github.com/wopta/goworkspace/lib"
	model "github.com/wopta/goworkspace/models"
)

func ContractFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println("Contract")
	//lib.Files("./serverless_function_source_code")
	req := lib.ErrorByte(ioutil.ReadAll(r.Body))
	var data model.Policy
	defer r.Body.Close()
	err := json.Unmarshal([]byte(req), &data)
	lib.CheckError(err)
	respObj := <-ContractObj(data)
	resp, err := json.Marshal(respObj)

	lib.CheckError(err)
	return string(resp), respObj, nil
}

func ContractObj(data model.Policy) <-chan DocumentResponse {
	r := make(chan DocumentResponse)

	//now := time.Now()
	//next := now.AddDate(0, 0, 4)
	//layout := "2006-01-02T15:04:05.000Z"

	go func() {
		skin := getVar()
		m := skin.initDefault()
		switch data.Company {
		case "global":
			skin.GlobalContract(m, data)
		}
		//-----------Save file
		Save(m, data)

		log.Println(data.Uid + " ContractObj end")
	}()
	return r
}
