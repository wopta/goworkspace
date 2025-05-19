package namirial

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	"github.com/wopta/goworkspace/document"
	"github.com/wopta/goworkspace/lib"
	env "github.com/wopta/goworkspace/lib/environment"
	"github.com/wopta/goworkspace/lib/log"
	"github.com/wopta/goworkspace/models"
)

func Sign(input NamirialInput) (response NamirialOutput, err error) {
	log.AddPrefix("Namirial")
	defer log.PopPrefix()

	fileIds, err := uploadFiles(input.FilesName...)
	if err != nil {
		return response, err
	}
	resp, err := prepareDocuments(fileIds...)
	if err != nil {
		return response, err
	}
	callbackurl := `"https://europe-west1-` + os.Getenv("GOOGLE_PROJECT_ID") + `.cloudfunctions.net/callback/v1/sign?envelope=##EnvelopeId##&action=##Action##&uid=` + input.Policy.Uid + `&token=` + os.Getenv("WOPTA_TOKEN_API") + `&origin=` + input.Origin + `&sendEmail=` + strconv.FormatBool(input.SendEmail) + `"`
	idEnvelope, err := sendDocuments(resp, fileIds, input.Policy, callbackurl)
	if err != nil {
		return response, err
	}
	envelope, err := getEnvelope(idEnvelope)
	if err != nil {
		return response, err
	}
	return NamirialOutput{
		Url:        envelope.ViewerLinks[0].ViewerLink,
		IdEnvelope: idEnvelope,
		FileIds:    fileIds,
	}, err
}

// upload the files thought namirial
// TODO: could be done in parallel ?
func uploadFiles(files ...string) (fileIds []string, err error) {
	var url = os.Getenv("ESIGN_BASEURL") + "v6/file/upload"

	var file []byte
	var buffer bytes.Buffer
	var idsFile []string
	for i := range files {
		if env.IsLocal() {
			file, err = os.ReadFile("document/contract.pdf")
		} else {
			file, err = lib.GetFromStorageErr(os.Getenv("GOOGLE_STORAGE_BUCKET"), files[i], "")
		}
		if err != nil {
			return fileIds, err
		}
		if file == nil || len(file) == 0 {
			return fileIds, fmt.Errorf("Error getting the file %v", files[i])
		}
		w := multipart.NewWriter(&buffer)
		fw, err := w.CreateFormFile("file", files[i]+".pdf")
		if err != nil {
			return fileIds, err
		}
		nWrite, err := fw.Write(file)
		if err != nil || nWrite == 0 {
			return fileIds, err
		}
		w.Close()

		req, err := http.NewRequest(http.MethodPost, url, &buffer)
		if err != nil {
			return fileIds, err
		}
		req.Header.Set("apiToken", os.Getenv("ESIGN_TOKEN_API"))
		req.Header.Set("Content-Type", w.FormDataContentType())
		res, err := handleResponse[struct{ FileId string }](lib.RetryDo(req, 5, 30))
		if err != nil {
			return fileIds, err
		}
		if res.FileId == "" {
			return fileIds, fmt.Errorf("Error: no fileId found")
		}
		idsFile = append(idsFile, res.FileId)
		log.Printf("End uploading files, idFiles %v", res)
	}
	return idsFile, nil
}

func prepareDocuments(idsDocument ...string) (resp document.PrepareResponse, err error) {
	var url = os.Getenv("ESIGN_BASEURL") + "v6/file/prepare"

	log.Println("Start preparing files")

	request := prepareNamirialDocumentRequest{
		FileIds:                   idsDocument,
		ClearAdvancedDocumentTags: true,
		SigStringConfigurations: []sigStringConfiguration{{
			StartPattern:         "string",
			EndPattern:           "string",
			ClearSigString:       true,
			SearchEntireWordOnly: true,
		},
		}}

	req, err := doNamirialRequest(http.MethodPost, url, request)
	if err != nil {
		return resp, err
	}

	resp, err = handleResponse[document.PrepareResponse](lib.RetryDo(req, 5, 30))
	if err != nil {
		return resp, err
	}
	if env.IsLocal() || env.IsDevelopment() {
		for i := range resp.Activities {
			resp.Activities[i].Action.Sign.RecipientConfiguration.AuthenticationConfiguration.AccessCode.Code = "test"
		}
	}

	log.Println("End preparing files")
	return resp, nil
}

func sendDocuments(preSendBody document.PrepareResponse, idFiles []string, policy models.Policy, callbackUrl string) (idEnvelope string, err error) {
	var url = os.Getenv("ESIGN_BASEURL") + "v6/envelope/send"
	var body sendNamirialRequest
	log.Println("Sending documents")

	body.Activities = preSendBody.Activities
	body.Documents = make([]documentDescription, len(idFiles))
	for i := range idFiles {
		body.Documents[i] = documentDescription{FileId: idFiles[i], DocumentNumber: i + 1} //the document number has to start from 1
	}
	body.CallbackConfiguration.CallbackUrl = callbackUrl
	body.CallbackConfiguration.StatusUpdateCallbackUrl = callbackUrl
	body.CallbackConfiguration.ActivityActionCallbackConfig = activityActionCallbackConfiguration{
		Url: callbackUrl,
	}
	setContractorDataInSendBody(&body, policy)
	req, err := doNamirialRequest("POST", url, body)
	if err != nil {
		return idEnvelope, err
	}
	resp, err := handleResponse[responseSendDocuments](lib.RetryDo(req, 5, 30))
	if err != nil {
		return idEnvelope, err
	}
	if resp.EnvelopeId == "" {
		return idEnvelope, fmt.Errorf("Error: no envelopId found")
	}
	idEnvelope = resp.EnvelopeId
	log.Println("End sending documents")
	return idEnvelope, err
}

// adjust the request to insert information regard the contractor
func setContractorDataInSendBody(bodySend *sendNamirialRequest, policy models.Policy) {
	contractor := policy.Contractor
	for i := range bodySend.Activities {
		for range bodySend.Activities[i].Action.Sign.Elements.Signatures {
			contactInfo := &bodySend.Activities[i].Action.Sign.RecipientConfiguration.ContactInformation
			contactInfo.LanguageCode = "IT"
			contactInfo.Surname = contractor.Surname
			contactInfo.GivenName = contractor.Name
			contactInfo.Email = contractor.Mail
			contactInfo.PhoneNumber = contractor.Phone
			contactInfo.PhoneNumber = contractor.Phone
		}
	}
	//TODO: i dont know if it is correct
	bodySend.Name = fmt.Sprint(bodySend.Name, ",", policy.CodeCompany)
}

// return an object that contains a link to open and sign the documents
func getEnvelope(idEvenelope string) (ResponseGetEvelop, error) {
	var resp ResponseGetEvelop
	var url = os.Getenv("ESIGN_BASEURL") + "v6/envelope/" + idEvenelope + "/viewerlinks"
	log.Println("Start Getting envelop")

	if idEvenelope == "" {
		return resp, fmt.Errorf("Error:no envelope id founded")
	}

	req, err := doNamirialRequest(http.MethodGet, url, nil)
	if err != nil {
		return resp, err
	}

	resp, err = handleResponse[ResponseGetEvelop](lib.RetryDo(req, 5, 30))

	if err != nil {
		return resp, err
	}
	log.Println("End getting evenlop")

	return resp, err
}
