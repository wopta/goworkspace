package WiseBroker

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	lib "github.com/wopta/goworkspace/lib"
	//"google.golang.org/api/firebaseappcheck/v1"
)

func init() {
	log.Println("INIT WiseBroker")
	functions.HTTP("WiseBroker", WiseBroker)
}

func WiseBroker(w http.ResponseWriter, r *http.Request) {
	var (
		//idToken string
		service = "http://test-wopta.northeurope.cloudapp.azure.com/WiseWebServicePtf/service.asmx"
	)
	lib.EnableCors(&w, r)
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	log.Println("WiseBroker")
	req := lib.ErrorByte(ioutil.ReadAll(r.Body))
	log.Println(string(req))
	var send Request
	// Unmarshal or Decode the JSON to the interface.
	//json.NewDecoder(req).Decode(&send)
	defer r.Body.Close()

	json.Unmarshal([]byte(req), &send)
	tmplt := template.New("action")
	var tpl bytes.Buffer
	tmplt, _ = template.ParseFiles("wise-broker/request-epoli.html")
	tmplt.Execute(&tpl, send)
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	response, err := client.Post(service, "text/xml", bytes.NewBufferString(tpl.String()))
	lib.CheckError(err)
	defer response.Body.Close()

	content, _ := ioutil.ReadAll(response.Body)
	s := strings.TrimSpace(string(content))
	log.Println(s)
}

type Request struct {
	From         string   `json:"from"`
	To           []string `json:"to"`
	Message      string   `json:"message"`
	Subject      string   `json:"subject"`
	IsHtml       bool     `json:"isHtml,omitempty"`
	IsAttachment bool     `json:"isAttachment,omitempty"`
	Cc           string   `json:"cc,omitempty"`
	TemplateName string   `json:"templateName,omitempty"`
}
