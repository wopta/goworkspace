package test

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/test/namirial"
)

type requestEnvelop struct {
	Uids []string
}

func HandlerEnvelop(w http.ResponseWriter, r *http.Request) (string, any, error) {
	req, _ := getBodyRequest[requestEnvelop](r)
	nav,_:=namirial.NewNamirial("",req.Uids...)
	//is it neccessary use the channel?
	var envelopId string
	var envInfo namirial.ResponeGetEvelop
	var filesIds namirial.FilesIdsResponse
	var err error

	if err:=nav.UploadFiles();err!=nil{
		return "","",err
	}
	if err:=nav.PrepareDocument();err!=nil{
		return "","",err
	}

	if envelopId,err=nav.SendDocuments();err!=nil{
		return "","",err
	}

	if envInfo,err=nav.GetEnvelope();err!=nil{
		return "","",err
	}

	if filesIds,err=nav.GetIdsFiles();err!=nil{
		return "","",err
	}

	return fmt.Sprintf("Envelope id: %v\n, EnvelopeInfo %v\n, Files info %v\n",envelopId,envInfo,filesIds),nil,nil
}

func getBodyRequest[T any](r *http.Request) (T, error) {
	var (
		req T
	)
	log.SetPrefix("[Unmarshal] ")
	defer log.SetPrefix("")

	body := lib.ErrorByte(io.ReadAll(r.Body))

	defer r.Body.Close()
	if err := json.Unmarshal(body, &req); err != nil {
		return *new(T), err
	}
	return req, nil
}
