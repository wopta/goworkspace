package namirial

import (
	"bytes"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	"gitlab.dev.wopta.it/goworkspace/document"
	"gitlab.dev.wopta.it/goworkspace/lib"
	env "gitlab.dev.wopta.it/goworkspace/lib/environment"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/models"
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

	log.Println("Start uploading files")
	var file []byte
	var buffer bytes.Buffer
	var idsFile []string
	for i := range files {
		if env.IsLocal() {
			if i%2 == 0 {
				file, err = os.ReadFile("document/net.pdf")
			} else {
				file, err = os.ReadFile("document/contract.pdf")
			}
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
		log.Printf("files uploaded, idFiles %v", res)
	}
	log.Println("End uploading files")
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
		}, {
			StartPattern:         "Rappresentante_Legale",
			ClearSigString:       true,
			SearchEntireWordOnly: true,
		}}}

	req, err := doNamirialRequest(http.MethodPost, url, request)
	if err != nil {
		return resp, err
	}

	resp, err = handleResponse[document.PrepareResponse](lib.RetryDo(req, 5, 30))
	if err != nil {
		return resp, err
	}
	if len(resp.Activities) == 0 {
		resp.Activities = append(resp.Activities, document.Activity{})
	}
	//The signatures that dont use the default place holder are put inside Unassigned,so you need to iterate them and fix their size and position
	for i := range resp.UnassignedElements.Signatures {
		sign := &resp.UnassignedElements.Signatures[i]
		if sign.FieldDefinition.Size.Height < sign.FieldDefinition.Size.Width {
			continue
		}
		sign.FieldDefinition.Size.Height = 50
		sign.FieldDefinition.Position.X -= 25
		sign.FieldDefinition.Position.Y -= 10
		sign.FieldDefinition.Size.Width = 150
	}
	resp.Activities[0].Action.Sign.Elements.Signatures = append(resp.Activities[0].Action.Sign.Elements.Signatures, resp.UnassignedElements.Signatures...)
	for i := range resp.Activities[0].Action.Sign.Elements.Signatures {
		sign := &resp.Activities[0].Action.Sign.Elements.Signatures[i]
		sign.DisplayName = "firma qui"
	}
	log.Println("End preparing files")
	return resp, nil
}

func sendDocuments(preSendBody document.PrepareResponse, idFiles []string, policy models.Policy, callbackUrl string) (idEnvelope string, err error) {
	var url = os.Getenv("ESIGN_BASEURL") + "v6/envelope/send"
	var body sendNamirialRequest
	log.Println("Sending documents")

	body.Activities = preSendBody.Activities

	if env.IsLocal() || env.IsDevelopment() {
		for i := range body.Activities {
			body.Activities[i].Action.Sign.RecipientConfiguration.AuthenticationConfiguration.AccessCode.Code = "test"
		}
	}
	body.Documents = make([]documentDescription, len(idFiles))
	for i := range idFiles {
		body.Documents[i] = documentDescription{FileId: idFiles[i], DocumentNumber: i + 1} //the document number has to start from 1
	}
	body.CallbackConfiguration.CallbackUrl = callbackUrl
	body.CallbackConfiguration.StatusUpdateCallbackUrl = callbackUrl
	body.CallbackConfiguration.ActivityActionCallbackConfig = activityActionCallbackConfiguration{
		Url: callbackUrl,
	}
	err = setContractorDataInSendBody(&body, policy)
	if err != nil {
		return idEnvelope, err
	}
	log.PrintStruct("request send", body)
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
func setContractorDataInSendBody(bodySend *sendNamirialRequest, policy models.Policy) error {
	var signer *models.User
	if policy.Contractor.Type == "legalEntity" { //for legalentity who pay is between contractors
		for _, contractor := range *policy.Contractors {
			if contractor.IsSignatory {
				signer = &contractor
				break
			}
		}
	} else { //otherwise i use contractor
		signer = policy.Contractor.ToUser()
	}

	if signer == nil {
		return errors.New("You need to populate contractors to sign")
	}
	for i := range bodySend.Activities {
		contactInfo := &bodySend.Activities[i].Action.Sign.RecipientConfiguration.ContactInformation
		contactInfo.LanguageCode = "IT"
		contactInfo.Surname = signer.Surname
		contactInfo.GivenName = signer.Name
		contactInfo.Email = signer.Mail
		contactInfo.PhoneNumber = signer.Phone
		contactInfo.PhoneNumber = signer.Phone
	}
	bodySend.Name = policy.Name
	return nil
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
