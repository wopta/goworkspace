package namirial

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/wopta/goworkspace/document"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
	"github.com/wopta/goworkspace/policy"
	"github.com/wopta/goworkspace/product"
)

type Namirial struct {
	idEnvelope      string
	prepareDocument document.PrepareResponse
	filesIds        FilesIdsResponse
	status          StatusNamirial
	dtos            []dataForDocument
}

func NewNamirial(origin string, uids ...string) (*Namirial, error) {
	dtos := make([]dataForDocument, len(uids))
	var (
		policyModel  models.Policy
		networkModel *models.NetworkNode
		productModel *models.Product
		warrant      *models.Warrant
		err          error
	)
	if len(uids) == 0 {
		return nil, fmt.Errorf("it is necessary have atleast 1 uid")
	}

	for i, uid := range uids {
		policyModel, err = policy.GetPolicy(uid, origin)
		if err != nil {
			return nil, err
		}
		networkModel = network.GetNetworkNodeByUid(policyModel.ProducerUid)
		if networkModel != nil {
			warrant = networkModel.GetWarrant()
		}
		productModel = product.GetProductV2(policyModel.Name, policyModel.ProductVersion, policyModel.Channel, networkModel, warrant)
		dtos[i] = dataForDocument{
			policy:  &policyModel,
			product: productModel,
			warrant: warrant,
		}
		if os.Getenv("env") == "local" {
			<-document.ContractObj("", policyModel, networkModel, productModel)
		}
	}
	return &Namirial{
		dtos:   dtos,
		status: Idle,
	}, nil
}

// upload the files for each policy passed throught NewNamirial
func (n *Namirial) UploadFiles() error {
	if n.status != Idle {
		return fmt.Errorf("Error: cant upload file, status has to be Idle, instead is %v", n.status)
	}
	var buffer bytes.Buffer
	log.SetPrefix("[UploadFiles]")
	defer log.SetPrefix("")

	log.Println("Start uploading files")

	var urlstring = os.Getenv("ESIGN_BASEURL") + "v6/file/upload"
	for i, dto := range n.dtos {
		var file []byte
		var err error
		if os.Getenv("env") == "local" {
			file, err = os.ReadFile("document/proposal.pdf")
		} else {
			file = lib.GetFromStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), dto.policy.DocumentName, "")
		}
		if err != nil {
			return err
		}
		if file == nil || len(file) == 0 {
			return fmt.Errorf("Error getting the file %v", dto.policy.DocumentName)
		}

		//Create form body
		w := multipart.NewWriter(&buffer)
		fw, err := w.CreateFormFile("file", dto.policy.DocumentName+" Polizza.pdf")
		if err != nil {
			return err
		}
		nWrite, err := fw.Write(file)
		if err != nil || nWrite == 0 {
			return err
		}
		w.Close()

		req, err := http.NewRequest("POST", urlstring, &buffer)
		if err != nil {
			return err
		}
		req.Header.Set("apiToken", os.Getenv("ESIGN_TOKEN_API"))
		req.Header.Set("Content-Type", w.FormDataContentType())

		res, err := handleResponse[struct{ FileId string }](lib.RetryDo(req, 5, 30))
		if err != nil {
			return err
		}
		if res.FileId == "" {
			return fmt.Errorf("Error: no fileId found")
		}
		n.dtos[i].idDocument = res.FileId

		n.status = Upload

		log.Printf("End uploading files, idFiles %v", res)
	}
	return nil
}

// prepare and set the documents uploaded, fix the sign,set the position ecc
func (n *Namirial) PrepareDocument() error {
	if n.status != Upload {
		return fmt.Errorf("Error: cant prepare files, status has to be upload, instead is %v", n.status)
	}
	log.SetPrefix("[PrepareFiles]")
	defer log.SetPrefix("")

	log.Println("Start preparing files")

	var urlstring = os.Getenv("ESIGN_BASEURL") + "v6/file/prepare"
	var ids = make([]string, len(n.dtos))
	for i, dto := range n.dtos {
		ids[i] = dto.idDocument
	}

	request := prepareNamirialDocumentRequest{
		FileIds:                   ids,
		ClearAdvancedDocumentTags: true,
		SigStringConfigurations: []sigStringConfiguration{{
			StartPattern:         "string",
			EndPattern:           "string",
			ClearSigString:       true,
			SearchEntireWordOnly: true,
		},
	}}

	req, err := doRequestNamirial(http.MethodPost, urlstring, request)
	if err != nil {
		return err
	}

	resPrepare, err := handleResponse[document.PrepareResponse](lib.RetryDo(req, 5, 30))
	n.prepareDocument = resPrepare
	n.status = Prepared
	if os.Getenv("env") == "local" || os.Getenv("env") == "dev" {
		for i := range n.prepareDocument.Activities {
			n.prepareDocument.Activities[i].Action.Sign.RecipientConfiguration.AuthenticationConfiguration.AccessCode.Code = "test"
		}

	}
	log.Println("End preparing files")
	return nil

}

// send each documents, previusly prepared,return an envelope id
func (n *Namirial) SendDocuments() (string, error) {
	if n.status != Prepared {
		return "", fmt.Errorf("Error: cant send files, status has to be preparared, instead is %v", n.status)
	}
	log.SetPrefix("[SendEnvelop]")
	defer log.SetPrefix("")

	log.Println("Start Sending files")

	var urlstring = os.Getenv("ESIGN_BASEURL") + "v6/envelope/send"

	var request sendNamirialRequest
	request.Documents = make([]documentDescription, len(n.dtos))
	request.Activities = n.prepareDocument.Activities
	n.adjectSendBody(&request)

	for i := range n.dtos {
		request.Documents[i] = documentDescription{FileId: n.dtos[i].idDocument, DocumentNumber: i + 1}
	}

	req, err := doRequestNamirial("POST", urlstring, request)
	if err != nil {
		return "", err
	}

	body, err := handleResponse[responseSendDocuments](lib.RetryDo(req, 5, 30))
	if err != nil {
		return "", err
	}

	if body.EnvelopeId == "" {
		return "", fmt.Errorf("Error: no envelopId found")
	}

	n.idEnvelope = body.EnvelopeId
	n.status = Sended
	log.Printf("End sending files,idEnvelope %v", n.idEnvelope)

	return n.idEnvelope, nil
}

// adjust the request to insert information regard the contractor
func (n *Namirial) adjectSendBody(d *sendNamirialRequest) {
	for i := range d.Activities {
		for _, el := range d.Activities[i].Action.Sign.Elements.Signatures {
			contactInfo := &d.Activities[i].Action.Sign.RecipientConfiguration.ContactInformation

			contactInfo.LanguageCode = "IT"
			contactInfo.Surname = n.dtos[el.DocumentNumber-1].policy.Contractor.Surname
			contactInfo.GivenName = n.dtos[el.DocumentNumber-1].policy.Contractor.Name
			contactInfo.Email = n.dtos[el.DocumentNumber-1].policy.Contractor.Mail
			contactInfo.PhoneNumber = n.dtos[el.DocumentNumber-1].policy.Contractor.Phone
			contactInfo.PhoneNumber = n.dtos[el.DocumentNumber-1].policy.Contractor.Phone
		}
	}
	//TODO: i dont know if it is correct
	d.Name = fmt.Sprint(d.Name, ",", n.dtos[0].policy.CodeCompany)
}

// return an object that contains a link to open and sign the documents
func (n *Namirial) GetEnvelope() (ResponeGetEvelop, error) {
	var resp ResponeGetEvelop
	if n.status != Sended {
		return resp, fmt.Errorf("Error: cant open the envelope, status has to be sended, instead is %v", n.status)
	}
	log.SetPrefix("[GetEnvelope]")
	defer log.SetPrefix("")

	log.Println("Start Getting envelop")

	if n.idEnvelope == "" {
		return resp, fmt.Errorf("Error:no envelope id founded")
	}
	var urlstring = os.Getenv("ESIGN_BASEURL") + "v6/envelope/" + n.idEnvelope + "/viewerlinks"

	req, err := doRequestNamirial(http.MethodGet, urlstring, nil)
	if err != nil {
		return resp, err
	}

	resp, err = handleResponse[ResponeGetEvelop](lib.RetryDo(req, 5, 30))

	if err != nil {
		return resp, err
	}

	n.status = GetEnvelope
	log.Println("End getting evenlop")

	return resp, err
}

// get ids of each files to download them eventually, they have to be signead
func (n *Namirial) GetIdsFiles() (FilesIdsResponse, error) {
	var resp FilesIdsResponse
	log.SetPrefix("[GetIdsFiles]")
	defer log.SetPrefix("")

	if n.status != GetEnvelope {
		return resp, fmt.Errorf("Error: cant have the id files if at least you have status Open Envelope, instead is %v", n.status)
	}

	log.Println("Start Getting IdsFiles")

	if n.idEnvelope == "" {
		return resp, fmt.Errorf("Error:no envelope id founded")
	}
	var urlstring = os.Getenv("ESIGN_BASEURL") + "v6/envelope/" + n.idEnvelope + "/files"

	req, err := doRequestNamirial(http.MethodGet, urlstring, nil)
	if err != nil {
		return resp, err
	}

	body, err := handleResponse[FilesIdsResponse](lib.RetryDo(req, 5, 30))
	if err != nil {
		return resp, err
	}

	resp = body
	n.filesIds = body

	log.Println("End getting ids files")
	return resp, nil
}

//i dont know how to manage the download of file for the front end
//TODO
//func (n *Namirial) DowloadFiles(idfile string) ([]byte,error){
//	log.SetPrefix("[DownloadFiles]")
//	defer log.SetPrefix("")
//	if idfile==""{
//		return nil,fmt.Errorf("Error:no idFile")
//	}
//	log.Println("Start Getting IdsFiles")
//
//	var urlstring = os.Getenv("ESIGN_BASEURL") + "v6/file/" + idfile
//
//	req, err := http.NewRequest("GET", urlstring, nil)
//	if err != nil {
//		log.Printf("Error:%v", err.Error())
//		return nil,err
//	}
//	req.Header.Set("apiToken", os.Getenv("ESIGN_TOKEN_API"))
//	req.Header.Set("Content-Type", "application/json")
//	res, err := lib.RetryDo(req, 5, 30)
//	if err != nil {
//		log.Printf("Error: error request")
//		return nil,err
//	}
//
//	if res.StatusCode != http.StatusOK {
//		log.Printf("Error: request %v", res.Status)
//		return nil,err
//	}
//	body, err := getBodyResponse[[]byte](res)
//	if err != nil {
//		log.Printf("Error: %v", err.Error())
//		return nil,err
//	}
//	return body,err
//}

// check http response(status code) and unmarshal the body
func handleResponse[T any](r *http.Response, err error) (T, error) {
	var (
		req T
	)
	if err != nil {
		return req, err
	}

	if r.StatusCode != http.StatusOK {
		var body []byte
		r.Body.Read(body)
		return req, err
	}

	log.SetPrefix("[Unmarshal] ")
	defer log.SetPrefix("")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return *new(T), err
	}
	defer r.Body.Close()
	if err := json.Unmarshal(body, &req); err != nil {
		return *new(T), err
	}
	return req, nil
}

// create the request passing the body(struct)
func doRequestNamirial(method string, url string, body any) (*http.Request, error) {
	var (
		err error
		req *http.Request
	)

	if body == nil {
		req, err = http.NewRequest(method, url, nil)
	} else {
		requestJson, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		requestReader := bytes.NewReader(requestJson)
		req, err = http.NewRequest(method, url, requestReader)
	}

	if err != nil {
		return nil, err
	}

	req.Header.Set("apiToken", os.Getenv("ESIGN_TOKEN_API"))
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}
