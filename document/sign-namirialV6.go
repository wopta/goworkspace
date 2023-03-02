package document

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"strings"

	"os"
	"path/filepath"
	"time"

	lib "github.com/wopta/goworkspace/lib"
	model "github.com/wopta/goworkspace/models"
)

func SignNamirialV6(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {

	req := lib.ErrorByte(ioutil.ReadAll(r.Body))
	var data model.Policy
	defer r.Body.Close()
	err := json.Unmarshal([]byte(req), &data)
	lib.CheckError(err)
	//file, _ := os.Open("document/billing.pdf")

	return NamirialOtpV6(data)

}

func NamirialOtpV6(data model.Policy) (string, NamirialOtpResponse, error) {
	var file []byte
	if os.Getenv("env") == "local" {
		file = lib.ErrorByte(ioutil.ReadFile("document/contract.pdf"))

	} else {
		file = lib.GetFromStorage("function-data", data.DocumentName, "")
	}

	SspFileId := <-postDataV6(file)
	log.Println(data.Uid+"postData:", SspFileId)
	//prepareEnvelop(SspFileId)
	unassigned := <-prepareEnvelopV6(SspFileId)
	log.Println("prepare body:", unassigned)
	id := <-sendEnvelopV6(SspFileId, data, unassigned)

	log.Println(data.Uid+"sendEnvelop:", id)
	url := <-GetEnvelopV6(id)
	resp := NamirialOtpResponse{
		EnvelopeId: id,
		Url:        url,
		FileId:     SspFileId,
	}
	return "{}", resp, nil
}

func prepareEnvelopV6(id string) <-chan string {
	r := make(chan string)
	go func() {
		log.Println("prepare")
		//var b bytes.Buffer
		//fileReader := bytes.NewReader([]byte())
		var urlstring = os.Getenv("ESIGN_BASEURL") + "v6/file/prepare"
		client := &http.Client{
			Timeout: time.Second * 10,
		}
		req, _ := http.NewRequest(http.MethodPost, urlstring, strings.NewReader(getPrepareV6(id)))
		req.Header.Set("apiToken", os.Getenv("ESIGN_TOKEN_API"))
		req.Header.Set("Content-Type", "application/json")
		//header('Content-Length: ' . filesize($pdf));
		log.Println("url parse:", req.Header)
		res, err := client.Do(req)
		lib.CheckError(err)

		if res != nil {
			body, err := ioutil.ReadAll(res.Body)
			lib.CheckError(err)
			var result map[string]string
			json.Unmarshal([]byte(body), &result)
			res.Body.Close()
			//res, e := json.Marshal(result)
			//lib.CheckError(e)
			r <- string(body)

			log.Println("body prepareEnvelopV6:", string(body))
		}
	}()
	return r
}
func sendEnvelopV6(id string, data model.Policy, unassigned string) <-chan string {
	r := make(chan string)

	go func() {
		defer close(r)
		log.Println("Send")
		var urlstring = os.Getenv("ESIGN_BASEURL") + "/v6/envelope/send"
		client := &http.Client{
			Timeout: time.Second * 10,
		}
		log.Println(data.Uid+" body:", string(getSendV6(id, data, unassigned)))
		req, _ := http.NewRequest(http.MethodPost, urlstring, strings.NewReader(getSendV6(id, data, unassigned)))
		req.Header.Set("apiToken", os.Getenv("ESIGN_TOKEN_API"))
		req.Header.Set("Content-Type", "application/json")
		//header('Content-Length: ' . filesize($pdf));

		res, err := client.Do(req)
		lib.CheckError(err)

		if res != nil {
			body, err := ioutil.ReadAll(res.Body)
			lib.CheckError(err)
			var result map[string]string
			json.Unmarshal([]byte(body), &result)
			res.Body.Close()

			log.Println(data.Uid+" body sendEnvelopV6: ", string(body))
			r <- result["EnvelopeId"]

		}
	}()
	return r
}
func GetEnvelopV6(id string) <-chan string {
	r := make(chan string)

	go func() {
		defer close(r)
		log.Println("GetEnvelopV6")
		//var b bytes.Buffer
		//fileReader := bytes.NewReader([]byte())GET /v6/envelope/{envelopeId}/viewerlinks
		var urlstring = os.Getenv("ESIGN_BASEURL") + "/v6/envelope/" + id + "/viewerlinks"
		client := &http.Client{
			Timeout: time.Second * 10,
		}
		req, _ := http.NewRequest(http.MethodGet, urlstring, nil)
		req.Header.Set("apiToken", os.Getenv("ESIGN_TOKEN_API"))
		req.Header.Set("Content-Type", "application/json")
		//header('Content-Length: ' . filesize($pdf));

		res, err := client.Do(req)
		lib.CheckError(err)

		if res != nil {
			body, err := ioutil.ReadAll(res.Body)
			lib.CheckError(err)
			var result GetEvelopViewerLinkResponse
			json.Unmarshal([]byte(body), &result)
			res.Body.Close()

			log.Println("body getEnvelop:", string(body))

			r <- result.ViewerLinks[0].ViewerLink

		}
	}()
	return r
}
func postDataV6(data []byte) <-chan string {
	r := make(chan string)
	var urlstring = os.Getenv("ESIGN_BASEURL") + "v6/file/upload"
	go func() {
		defer close(r)
		var b bytes.Buffer
		w := multipart.NewWriter(&b)
		// Add the field
		fw, err := w.CreateFormFile("file", filepath.Base("contract.pdf"))
		lib.CheckError(err)
		fw.Write((data)[:])
		w.Close()
		log.Println("postDataV6")
		req, err := http.NewRequest("POST", urlstring, &b)
		lib.CheckError(err)
		req.Header.Set("apiToken", os.Getenv("ESIGN_TOKEN_API"))
		req.Header.Set("Content-Type", w.FormDataContentType())
		client := &http.Client{
			Timeout: time.Second * 10,
		}
		res, err := client.Do(req)
		var result map[string]string
		if res != nil {
			resByte, err := ioutil.ReadAll(res.Body)
			lib.CheckError(err)
			json.Unmarshal(resByte, &result)
			res.Body.Close()
			log.Println(result["FileId"])
			r <- result["FileId"]
			fmt.Println(res.StatusCode)
		}
	}()

	return r
}
func GetFileV6(id string) string {
	r := make(chan string)

	go func() {
		defer close(r)
		var urlstring = os.Getenv("ESIGN_BASEURL") + "v4/authorization"
		client := &http.Client{
			Timeout: time.Second * 10,
		}
		req, _ := http.NewRequest(http.MethodGet, urlstring, nil)
		req.Header.Set("apiToken", os.Getenv("ESIGN_TOKEN_API"))
		log.Println("url parse:", req.Header)
		res, err := client.Do(req)
		lib.CheckError(err)

		if res != nil {
			body, err := ioutil.ReadAll(res.Body)
			lib.CheckError(err)
			res.Body.Close()
			lib.PutToFireStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "document/contracts"+id, body)
			log.Println("body:", string(body))
		}

	}()
	return ``
}
func getPrepareV6(id string) string {
	return `{
		"FileIds": [
			"` + id + `"
		],
		"ClearFieldMarkupString": true,
		"SigStringConfigurations": [
		  {
			"StartPattern": "string",
			"EndPattern": "string",
			"ClearSigString": true,
			"SearchEntireWordOnly": true
		  }
		]
	  }`
}
func getSendV6(id string, data model.Policy, prepare string) string {
	var preparePointer PrepareResponse
	json.Unmarshal([]byte(prepare), &preparePointer)
	unassignedElements := preparePointer.UnassignedElements
	elements := preparePointer.Activities[0].Action.Sign.Elements
	for i, element := range elements.Signatures {
		element.DocumentNumber = 0
		elements.Signatures[i].FieldDefinition.Position.X = elements.Signatures[i].FieldDefinition.Position.X + 20
		elements.Signatures[i].TaskConfiguration.OrderDefinition.OrderIndex = int64(i + 1)
	}
	unassignedJson, e := json.Marshal(unassignedElements)
	elementsJson, e := json.Marshal(elements)
	log.Println(string(unassignedJson))
	lib.CheckError(e)
	calbackurl := `"https://europe-west1-` + os.Getenv("GOOGLE_PROJECT_ID") + `.cloudfunctions.net/callback/v1/sign?envelope=##EnvelopeId##&action=##Action##&uid=` + data.Uid + `&token=` + os.Getenv("WOPTA_TOKEN_API") + `"`
	var testPin string
	if os.Getenv("env") == "local" || os.Getenv("env") == "dev" {
		testPin = `
			"AccessCode": {
			  "Code": "test"
			}, `
	}

	return `{
		"Documents": [
			{
				"FileId": "` + id + `",
				"DocumentNumber": 1
			}
		],
		"Name": "Test",
		"MetaData": "string",
		"AddDocumentTimestamp": true,
		"ShareWithTeam": true,
		"LockFormFieldsOnFinish": true,
		"UnassignedElements": ` + string(unassignedJson) + `,
		"Activities": [
			 
			{
			 
				"Action": {
					"Sign": {
						"SendEmails": false,
						"AllowAccessAfterFinish": false,
						"AllowDelegation": false,
                        "RequireViewContentBeforeFormFilling": false,
						"RecipientConfiguration": {
							"ContactInformation": {
								"Email": "` + data.Contractor.Mail + `",
								"GivenName": "` + data.Contractor.Name + `",
								"Surname": "` + data.Contractor.Surname + `",
								"PhoneNumber": "` + data.Contractor.Phone + `",
								"LanguageCode": "IT"
							},
							"PersonalMessage": "FIRMA LA TUA POLIZZA",
							"AuthenticationConfiguration": {
								` + testPin + `
							  "SmsOneTimePassword": {
								"PhoneNumber": "` + data.Contractor.Phone + `"
							  }
						}},
						"Elements": ` + string(elementsJson) + `,
						"SigningGroup": "CONTRAENTE"
					}
				}
			}
		
		],
		"AgentRedirectConfiguration": {
			"Policy": "None",
			"Allow": true,
			"IframeWhitelisting": [
			  "dev.wopta.it"
			]
		  },
		"ReminderConfiguration": {
			"Enabled": true,
			"FirstReminderInDays": 2,
			"ReminderResendIntervalInDays": 0,
			"BeforeExpirationInDays": 0
		  },
		"CallbackConfiguration": {
			"CallbackUrl": ` + calbackurl + `,
			"StatusUpdateCallbackUrl":` + calbackurl + ` ,
			"ActivityActionCallbackConfiguration": {
			  "Url": ` + calbackurl + `,
			  "ActionCallbackSelection": {
				"ConfirmTransactionCode": true,
				"AgreementAccepted": true,
				"AgreementRejected": true,
				"PrepareAuthenticationSuccess": true,
				"AuthenticationFailed": true,
				"AuthenticationSuccess": true,
				"AuditTrailRequested": true,
				"AuditTrailXmlRequested": true,
				"CalledPage": true,
				"DocumentDownloaded": true,
				"FlattenedDocumentDownloaded": true,
				"AddedAnnotation": true,
				"AddedAttachment": true,
				"AppendedDocument": true,
				"FormsFilled": true,
				"ConfirmReading": true,
				"SendTransactionCode": true,
				"PrepareSignWorkstepDocument": true,
				"SignWorkstepDocument": true,
				"UndoAction": true,
				"WorkstepCreated": true,
				"WorkstepFinished": true,
				"WorkstepRejected": true,
				"DisablePolicyAndValidityChecks": true,
				"EnablePolicyAndValidityChecks": true,
				"AppendFileToWorkstep": true,
				"AppendTasksToWorkstep": true,
				"SetOptionalDocumentState": true,
				"PreparePayloadForBatch": true
			  }
			}
		  }
	}`
}

type PrepareResponse struct {
	UnassignedElements Elements   `json:"UnassignedElements"`
	Activities         []Activity `json:"Activities"`
}

type Activity struct {
	Action Action `json:"Action"`
}

type Action struct {
	Sign Sign `json:"Sign"`
}

type Sign struct {
	Elements Elements `json:"Elements"`
}

type Elements struct {
	TextBoxes         []interface{}     `json:"TextBoxes"`
	CheckBoxes        []interface{}     `json:"CheckBoxes"`
	ComboBoxes        []interface{}     `json:"ComboBoxes"`
	RadioButtons      []interface{}     `json:"RadioButtons"`
	ListBoxes         []interface{}     `json:"ListBoxes"`
	Signatures        []Signature       `json:"Signatures"`
	Attachments       []interface{}     `json:"Attachments"`
	LinkConfiguration LinkConfiguration `json:"LinkConfiguration"`
}

type LinkConfiguration struct {
	HyperLinks []interface{} `json:"HyperLinks"`
}

type Signature struct {
	ElementID             string                `json:"ElementId"`
	Required              bool                  `json:"Required"`
	DocumentNumber        int64                 `json:"DocumentNumber"`
	DisplayName           string                `json:"DisplayName"`
	AllowedSignatureTypes AllowedSignatureTypes `json:"AllowedSignatureTypes"`
	FieldDefinition       FieldDefinition       `json:"FieldDefinition"`
	TaskConfiguration     TaskConfiguration     `json:"TaskConfiguration"`
}

type AllowedSignatureTypes struct {
	ClickToSign      ClickToSign   `json:"ClickToSign"`
	SignaturePlugins []interface{} `json:"SignaturePlugins"`
}

type ClickToSign struct {
	UseExternalSignatureImage string `json:"UseExternalSignatureImage"`
}

type FieldDefinition struct {
	Position Position `json:"Position"`
	Size     Size     `json:"Size"`
}

type Position struct {
	PageNumber int64   `json:"PageNumber"`
	X          float64 `json:"X"`
	Y          float64 `json:"Y"`
}

type Size struct {
	Width  float64 `json:"Width"`
	Height float64 `json:"Height"`
}

type TaskConfiguration struct {
	BatchGroup      string          `json:"BatchGroup"`
	OrderDefinition OrderDefinition `json:"OrderDefinition"`
}

type OrderDefinition struct {
	OrderIndex int64 `json:"OrderIndex"`
}

type GetEvelopViewerLinkResponse struct {
	ViewerLinks []ViewerLink `json:"ViewerLinks"`
}
type ViewerLink struct {
	ActivityID string `json:"ActivityId"`
	Email      string `json:"Email"`
	ViewerLink string `json:"ViewerLink"`
}
