package document

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"os"
	"time"

	lib "gitlab.dev.wopta.it/goworkspace/lib"
	env "gitlab.dev.wopta.it/goworkspace/lib/environment"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/network"
)

func SignNamirialV6(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	req := lib.ErrorByte(ioutil.ReadAll(r.Body))
	var data models.Policy
	defer r.Body.Close()
	err := json.Unmarshal([]byte(req), &data)
	lib.CheckError(err)
	//file, _ := os.Open("document/billing.pdf")

	sendEmail := true // CHECK

	return NamirialOtpV6(data, r.Header.Get("origin"), sendEmail)
}

func NamirialOtpV6(data models.Policy, origin string, sendEmail bool) (string, NamirialOtpResponse, error) {
	var file []byte

	if env.IsLocal() {
		file = lib.ErrorByte(os.ReadFile("document/contract.pdf"))
	} else {
		file = lib.GetFromStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), data.DocumentName, "")
	}

	SspFileId := <-postDataV6(file, data.NameDesc)
	log.Println(data.Uid+"postData:", SspFileId)
	//prepareEnvelop(SspFileId)
	unassigned := <-prepareEnvelopV6(SspFileId)

	id := <-sendEnvelopV6(SspFileId, data, unassigned, origin, sendEmail)

	log.Println(data.Uid+" sendEnvelop:", id)
	url := <-GetEnvelopV6(id)
	resp := NamirialOtpResponse{
		EnvelopeId: id,
		Url:        url,
		FileId:     SspFileId,
	}
	return "{}", resp, nil
}

func GetClient(method string, urlstring string, payload io.Reader) ([]byte, error) {
	var (
		r []byte
		e error
	)
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, _ := http.NewRequest(method, urlstring, payload)
	req.Header.Set("apiToken", os.Getenv("ESIGN_TOKEN_API"))
	res, e := client.Do(req)

	if res != nil {
		r, e = ioutil.ReadAll(res.Body)

	}

	return r, e
}

func prepareEnvelopV6(id string) <-chan string {
	r := make(chan string)
	go func() {
		log.Println("prepare")
		//var b bytes.Buffer
		//fileReader := bytes.NewReader([]byte())
		var urlstring = os.Getenv("ESIGN_BASEURL") + "v6/file/prepare"

		req, _ := http.NewRequest(http.MethodPost, urlstring, strings.NewReader(getPrepareV6(id)))
		req.Header.Set("apiToken", os.Getenv("ESIGN_TOKEN_API"))
		req.Header.Set("Content-Type", "application/json")

		res, err := lib.RetryDo(req, 10, 30)
		lib.CheckError(err)

		if res != nil {
			body, err := ioutil.ReadAll(res.Body)
			lib.CheckError(err)
			var result map[string]string
			json.Unmarshal([]byte(body), &result)
			res.Body.Close()
			//res, e := json.Marshal(result)
			//lib.CheckError(e)
			log.Println("resp  prepareEnvelopV6:", string(body))
			r <- string(body)

		}
	}()
	return r
}

func sendEnvelopV6(id string, data models.Policy, unassigned string, origin string, sendEmail bool) <-chan string {
	r := make(chan string)

	go func() {
		defer close(r)
		log.Println("Send")
		var urlstring = os.Getenv("ESIGN_BASEURL") + "/v6/envelope/send"
		req, _ := http.NewRequest(http.MethodPost, urlstring, strings.NewReader(getSendV6(id, data, unassigned, origin, sendEmail)))
		req.Header.Set("apiToken", os.Getenv("ESIGN_TOKEN_API"))
		req.Header.Set("Content-Type", "application/json")

		res, err := lib.RetryDo(req, 5, 30)
		lib.CheckError(err)

		if res != nil {
			body, err := io.ReadAll(res.Body)
			lib.CheckError(err)
			var result map[string]string
			json.Unmarshal([]byte(body), &result)
			res.Body.Close()

			log.Println(data.Uid+" body response sendEnvelopV6: ", string(body))
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

		req, _ := http.NewRequest(http.MethodGet, urlstring, nil)
		req.Header.Set("apiToken", os.Getenv("ESIGN_TOKEN_API"))
		req.Header.Set("Content-Type", "application/json")

		res, err := lib.RetryDo(req, 5, 30)
		lib.CheckError(err)

		if res != nil {
			body, err := ioutil.ReadAll(res.Body)
			lib.CheckError(err)
			var result GetEvelopViewerLinkResponse
			json.Unmarshal([]byte(body), &result)
			res.Body.Close()

			log.Println("body response getEnvelop:", string(body))

			r <- result.ViewerLinks[0].ViewerLink

		}
	}()
	return r
}

func postDataV6(data []byte, productNameDesc string) <-chan string {
	r := make(chan string)
	var urlstring = os.Getenv("ESIGN_BASEURL") + "v6/file/upload"
	go func() {
		defer close(r)
		var b bytes.Buffer
		w := multipart.NewWriter(&b)
		// Add the field
		fw, err := w.CreateFormFile("file", productNameDesc+" Polizza.pdf")
		lib.CheckError(err)
		fw.Write((data)[:])
		w.Close()
		log.Println("postDataV6")
		req, err := http.NewRequest("POST", urlstring, &b)
		lib.CheckError(err)
		req.Header.Set("apiToken", os.Getenv("ESIGN_TOKEN_API"))
		req.Header.Set("Content-Type", w.FormDataContentType())

		res, err := lib.RetryDo(req, 5, 30)
		var result map[string]string
		if res != nil {
			resByte, err := io.ReadAll(res.Body)
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

func GetFileV6(policy models.Policy, uid string) chan string {
	r := make(chan string)
	log.Println("Get file: ", policy.IdSign)
	contractPath := "assets/users/%s/" + models.ContractDocumentFormat
	go func() {

		defer close(r)
		files := <-GetFilesV6(policy.IdSign)

		var urlstring = os.Getenv("ESIGN_BASEURL") + "v6/file/" + files.Documents[0].FileID

		req, _ := http.NewRequest(http.MethodGet, urlstring, nil)
		req.Header.Set("apiToken", os.Getenv("ESIGN_TOKEN_API"))
		log.Println("url parse:", req.Header)

		res, err := lib.RetryDo(req, 5, 30)
		lib.CheckError(err)
		if res != nil {
			body, _ := io.ReadAll(res.Body)
			defer res.Body.Close()
			//log.Println("Get body: ", string(body))
			log.Println("Document Policy Contractor UID: ", policy.Contractor.Uid)
			gsLink, err := lib.PutToGoogleStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"),
				fmt.Sprintf(contractPath, policy.Contractor.Uid, policy.NameDesc, policy.CodeCompany), body)
			lib.CheckError(err)
			r <- gsLink

		}

	}()
	return r
}

func GetFilesV6(envelopeId string) chan NamirialFiles {
	r := make(chan NamirialFiles)

	go func() {
		defer close(r)
		var urlstring = os.Getenv("ESIGN_BASEURL") + "v6/envelope/" + envelopeId + "/files"

		req, _ := http.NewRequest(http.MethodGet, urlstring, nil)
		req.Header.Set("apiToken", os.Getenv("ESIGN_TOKEN_API"))

		res, err := lib.RetryDo(req, 5, 30)
		lib.CheckError(err)

		if res != nil {
			body, _ := io.ReadAll(res.Body)
			resp, _ := UnmarshalNamirialFiles(body)
			res.Body.Close()

			log.Println("body:", string(body))
			r <- resp
		}
	}()
	return r
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

func getSendV6(id string, data models.Policy, prepare string, origin string, sendEmail bool) string {
	var (
		preparePointer PrepareResponse
		redirectUrl    string
	)

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
	lib.CheckError(e)
	calbackurl := `"https://europe-west1-` + os.Getenv("GOOGLE_PROJECT_ID") + `.cloudfunctions.net/callback/v1/sign?envelope=##EnvelopeId##&action=##Action##&uid=` + data.Uid + `&token=` + os.Getenv("WOPTA_TOKEN_API") + `&origin=` + origin + `&sendEmail=` + strconv.FormatBool(sendEmail) + `"`
	var testPin string
	if !env.IsProduction() {
		testPin = `
		"AccessCode": {
			"Code": "test"
		}, `
	}

	nn := network.GetNetworkNodeByUid(data.ProducerUid)
	if data.Channel == models.NetworkChannel && nn != nil {
		warrant := nn.GetWarrant()
		if warrant != nil {
			flow := warrant.GetFlowName(data.Name)
			if flow == models.RemittanceMgaFlow {
				var baseUrl string = "https://www.wopta.it"
				if os.Getenv("env") != "prod" {
					baseUrl = "https://dev.wopta.it"
				}
				if lib.SliceContains([]string{models.LifeProduct, models.GapProduct}, data.Name) {
					redirectUrl = `
					"FinishActionConfiguration": {
						"SignAnyWhereViewer": {
							"RedirectUri": "` + baseUrl + `/it/quote/` + data.Name + `/thank-you"
						}
					},
					`
				}
			}
		}
	}

	return `{
		"Documents": [
		{
			"FileId": "` + id + `",
			"DocumentNumber": 1
		}
		],
		"Name": "` + data.CodeCompany + `",
		"MetaData": "string",
		"AddDocumentTimestamp": true,
		"ShareWithTeam": true,
		"LockFormFieldsOnFinish": true,
		"UnassignedElements": ` + string(unassignedJson) + `,
		"Activities": [
		{
			"Action": {
				"Sign": {
					"RequireViewContentBeforeFormFilling": false,
					"RecipientConfiguration": {
						"SendEmails": false,
						"AllowAccessAfterFinish": false,
						"AllowDelegation": false,
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
						}
					},
					"Elements": ` + string(elementsJson) + `,
					` + redirectUrl + `
					"SigningGroup": "CONTRAENTE"
				}
			}
		}

		],
		"AgentRedirectConfiguration": {
			"Policy": "None",
			"Allow": true,
			"IframeWhitelisting": [
			"dev.wopta.it", "wopta.it"
			]
		},
		"ReminderConfiguration": {
			"Enabled": true,
			"FirstReminderInDays": 2,
			"ReminderResendIntervalInDays": 1,
			"BeforeExpirationInDays": 1
		},
		"CallbackConfiguration": {
			"CallbackUrl": ` + calbackurl + `,
			"StatusUpdateCallbackUrl":` + calbackurl + ` ,
			"ActivityActionCallbackConfiguration": {
				"Url": ` + calbackurl + `,
				"ActionCallbackSelection": {
					"ConfirmTransactionCode": false,
					"AgreementAccepted": false,
					"AgreementRejected": false,
					"PrepareAuthenticationSuccess": false,
					"AuthenticationFailed":false,
					"AuthenticationSuccess":false,
					"AuditTrailRequested": false,
					"AuditTrailXmlRequested": false,
					"CalledPage": false,
					"DocumentDownloaded": false,
					"FlattenedDocumentDownloaded": false,
					"AddedAnnotation": false,
					"AddedAttachment": false,
					"AppendedDocument": false,
					"FormsFilled": false,
					"ConfirmReading": false,
					"SendTransactionCode": false,
					"PrepareSignWorkstepDocument": false,
					"SignWorkstepDocument": false,
					"UndoAction": false,
					"WorkstepCreated": false,
					"WorkstepFinished": true,
					"WorkstepRejected": true,
					"DisablePolicyAndValidityChecks": true,
					"EnablePolicyAndValidityChecks": true,
					"AppendFileToWorkstep": false,
					"AppendTasksToWorkstep": false,
					"SetOptionalDocumentState": false,
					"PreparePayloadForBatch": false
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
	Elements               Elements `json:"Elements"`
	RecipientConfiguration RecipientConfiguration
}

type RecipientConfiguration struct {
	SendEmails                  bool                        `json:"SendEmails"`
	AllowAccessAfterFinish      bool                        `json:"AllowAccessAfterFinish"`
	AllowDelegation             bool                        `json:"AllowDelegation"`
	ContactInformation          ContactInformation          `json:"ContactInformation"`
	PersonalMessage             string                      `json:"PersonalMessage"`
	AuthenticationConfiguration AuthenticationConfiguration `json:"AuthenticationConfiguration"`
}

type ContactInformation struct {
	Email        string `json:"Email"`
	GivenName    string `json:"GivenName"`
	Surname      string `json:"Surname"`
	PhoneNumber  string `json:"PhoneNumber"`
	LanguageCode string `json:"LanguageCode"`
}

type AuthenticationConfiguration struct {
	SmsOneTimePassword SmsOneTimePassword `json:"SmsOneTimePassword"`
	AccessCode         AccessCode
}

type AccessCode struct {
	Code string
}

type SmsOneTimePassword struct {
	PhoneNumber string `json:"PhoneNumber"`
}
type Elements struct {
	//TO implement if needed
	TextBoxes    []interface{} `json:"TextBoxes"`
	CheckBoxes   []interface{} `json:"CheckBoxes"`
	ComboBoxes   []interface{} `json:"ComboBoxes"`
	RadioButtons []interface{} `json:"RadioButtons"`
	ListBoxes    []interface{} `json:"ListBoxes"`
	Attachments  []interface{} `json:"Attachments"`

	Signatures []Signature `json:"Signatures"`
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

func UnmarshalNamirialFiles(data []byte) (NamirialFiles, error) {
	var r NamirialFiles
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *NamirialFiles) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type NamirialFiles struct {
	Documents      []Documents     `json:"Documents"`
	AuditTrail     AuditTrail      `json:"AuditTrail"`
	LegalDocuments []LegalDocument `json:"LegalDocuments"`
}

type AuditTrail struct {
	FileID    string `json:"FileId"`
	XMLFileID string `json:"XmlFileId"`
}

type Documents struct {
	FileID           string       `json:"FileId"`
	FileName         string       `json:"FileName"`
	AuditTrailFileID string       `json:"AuditTrailFileId"`
	Attachments      []Attachment `json:"Attachments"`
	PageCount        int64        `json:"PageCount"`
	DocumentNumber   int64        `json:"DocumentNumber"`
}

type Attachment struct {
	FileID   string `json:"FileId"`
	FileName string `json:"FileName"`
}

type LegalDocument struct {
	FileID     string `json:"FileId"`
	FileName   string `json:"FileName"`
	ActivityID string `json:"ActivityId"`
	Email      string `json:"Email"`
}
