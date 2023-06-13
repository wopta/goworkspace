package document

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"

	lib "github.com/wopta/goworkspace/lib"
	model "github.com/wopta/goworkspace/models"
)

func ContractFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println("Contract")
	//lib.Files("./serverless_function_source_code")
	req := lib.ErrorByte(io.ReadAll(r.Body))
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
		var (
			filename string
			out      []byte
		)

		switch data.Name {
		case "pmi":
			skin := getVar()
			m := skin.initDefault()
			skin.GlobalContract(m, data)
			//-----------Save file
			filename, out = Save(m, data)
		case "life":
			pdf := initFpdf()
			filename, out = LifeContract(pdf, &data)
		case "persona":
			pdf := initFpdf()
			filename, out = PersonaContract(pdf, &data)
		}

		data.DocumentName = filename
		log.Println(data.Uid + " ContractObj end")
		r <- DocumentResponse{
			LinkGcs: filename,
			Bytes:   base64.StdEncoding.EncodeToString(out),
		}
	}()
	return r
}
