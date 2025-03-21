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

type StatusNamirial string

const (
	Idle        StatusNamirial = "Idle"
	Upload      StatusNamirial = "Uploaded Files"
	Prepared    StatusNamirial = "Prepared Files"
	Sended      StatusNamirial = "Sended Files"
	GetEnvelope StatusNamirial = "Get Envelope"
)

type dataForDocument struct {
	policy     *models.Policy
	product    *models.Product
	warrant    *models.Warrant
	idDocument string
}

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
	}
	return &Namirial{
		dtos:   dtos,
		status: Idle,
	}, nil
}

// upload the files for each policy passed throught NewNamirial
func (n *Namirial) UploadFiles() error {
	if n.status != Idle {
		return fmt.Errorf("Error: cant upload file, status has to be Idle instead is %v", n.status)
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
			file, err = os.ReadFile("document/contract.pdf")
		} else {
			file = lib.GetFromStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), dto.policy.DocumentName, "")
		}
		if file == nil || len(file) == 0 {
			log.Printf("Error: getting file %v", dto.policy.DocumentName)
			return fmt.Errorf("Error getting the file %v", dto.policy.DocumentName)
		}
		if err != nil {
			log.Printf("Error:%v", err.Error())
			return err
		}

		//Create form body
		w := multipart.NewWriter(&buffer)
		fw, err := w.CreateFormFile("file", dto.policy.DocumentName+" Polizza.pdf")
		if err != nil {
			log.Printf("Error: creating form")
			return err
		}
		nWrite, err := fw.Write(file)
		if err != nil || nWrite == 0 {
			log.Printf("Error: writing into form")
			return err
		}
		w.Close()

		req, err := http.NewRequest("POST", urlstring, &buffer)
		if err != nil {
			log.Printf("Error: creating request")
			return err
		}
		req.Header.Set("apiToken", os.Getenv("ESIGN_TOKEN_API"))
		req.Header.Set("Content-Type", w.FormDataContentType())
		resResp, err := lib.RetryDo(req, 5, 30)
		if err != nil {
			log.Printf("Error: error request")
			return err
		}

		if resResp.StatusCode != http.StatusOK {
			log.Printf("Error: request %v ", resResp.Status)
			return err
		}

		res, err := getBodyResponse[struct{ FileId string }](resResp)
		if err != nil {
			log.Printf("Error: %v", err.Error())
			return err
		}
		if res.FileId == "" {
			log.Printf("Error: no fileId found")
			return err
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
		return fmt.Errorf("Error: cant prepare files, status has to be upload instead is %v", n.status)
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

	reqJson, err := json.Marshal(request)
	if err != nil {
		log.Printf("Error: %v", err.Error())
		return err
	}

	req, err := http.NewRequest(http.MethodPost, urlstring, bytes.NewReader(reqJson))
	if err != nil {
		log.Printf("Error: %v", err.Error())
		return err
	}

	req.Header.Set("apiToken", os.Getenv("ESIGN_TOKEN_API"))
	req.Header.Set("Content-Type", "application/json")

	resResp, err := lib.RetryDo(req, 5, 30)
	if err != nil {
		log.Printf("Error: error request %v", err.Error())
		return err
	}

	if resResp.StatusCode != http.StatusOK {
		log.Printf("Error: request %v ", resResp.Status)
		return err
	}
	resPrepare, err := getBodyResponse[document.PrepareResponse](resResp)
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
		return "", fmt.Errorf("Error: cant send files, status has to be preparared instead is %v", n.status)
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
	requestJson, err := json.Marshal(request)
	if err != nil {
		log.Printf("Error:%v", err.Error())
		return "", err
	}

	req, err := http.NewRequest("POST", urlstring, bytes.NewReader(requestJson))
	if err != nil {
		log.Printf("Error:%v", err.Error())
		return "", err
	}
	req.Header.Set("apiToken", os.Getenv("ESIGN_TOKEN_API"))
	req.Header.Set("Content-Type", "application/json")
	res, err := lib.RetryDo(req, 5, 30)

	if err != nil {
		log.Printf("Error: error request")
		return "", err
	}
	if res.StatusCode != http.StatusOK {
		log.Printf("Error: request %v", res.Status)
		return "", err
	}

	body, err := getBodyResponse[responseSendDocuments](res)
	if err != nil {
		log.Printf("Error: %v", err.Error())
		return "", err
	}
	if body.EnvelopeId == "" {
		log.Printf("Error: no envelopId found")
		return "", err
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
		return resp, fmt.Errorf("Error: cant open the envelope, status has to be sended instead is %v", n.status)
	}
	log.SetPrefix("[GetEnvelope]")
	defer log.SetPrefix("")

	log.Println("Start Getting envelop")

	if n.idEnvelope == "" {
		return resp, fmt.Errorf("Error:no envelope id founded")
	}
	var urlstring = fmt.Sprint(os.Getenv("ESIGN_BASEURL")+"v6/envelope/"+n.idEnvelope, "/viewerlinks")
	req, err := http.NewRequest("GET", urlstring, nil)

	req.Header.Set("apiToken", os.Getenv("ESIGN_TOKEN_API"))
	req.Header.Set("Content-Type", "none")
	resResp, err := lib.RetryDo(req, 5, 30)

	if err != nil {
		log.Printf("Error: error request %v", err.Error())
		return resp, err
	}
	if resResp.StatusCode != http.StatusOK {
		log.Printf("Error: request %v ", resResp.Status)
		return resp, err
	}

	resp, err = getBodyResponse[ResponeGetEvelop](resResp)

	if err != nil {
		log.Printf("Error: %v", err.Error())
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

	log.Println("Start Getting IdsFiles")

	if n.idEnvelope == "" {
		return resp, fmt.Errorf("Error:no envelope id founded")
	}
	var urlstring = os.Getenv("ESIGN_BASEURL") + "v6/envelope/" + n.idEnvelope + "/files"

	req, err := http.NewRequest("GET", urlstring, nil)
	if err != nil {
		log.Printf("Error:%v", err.Error())
		return resp, err
	}
	req.Header.Set("apiToken", os.Getenv("ESIGN_TOKEN_API"))
	req.Header.Set("Content-Type", "application/json")
	res, err := lib.RetryDo(req, 5, 30)
	if err != nil {
		log.Printf("Error: error request")
		return resp, err
	}

	if res.StatusCode != http.StatusOK {
		log.Printf("Error: request %v", res.Status)
		return resp, err
	}
	body, err := getBodyResponse[FilesIdsResponse](res)
	if err != nil {
		log.Printf("Error: %v", err.Error())
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

func getBodyResponse[T any](r *http.Response) (T, error) {
	var (
		req T
	)
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
