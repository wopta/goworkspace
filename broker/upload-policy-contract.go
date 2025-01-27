package broker

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
	plc "github.com/wopta/goworkspace/policy"
)

func UploadPolicyContractFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err    error
		policy models.Policy
	)

	log.SetPrefix("[UploadPolicyContractFx] ")
	defer func() {
		r.Body.Close()
		if err != nil {
			log.Printf("error: %s", err.Error())
		}
		log.Println("Handler end ----------------------------------------------")
		log.SetPrefix("")
	}()

	log.Println("Handler start -----------------------------------------------")

	policyUid := chi.URLParam(r, "uid")

	policy, err = plc.GetPolicy(policyUid, "")
	if err != nil {
		return "", nil, err
	}

	if !policy.CompanyEmit || policy.IsSign {
		err = fmt.Errorf("cannot upload policy contract, policy %s companyEmit: %v isSign: %v", policyUid,
			policy.CompanyEmit, policy.IsSign)
		return "", nil, err
	}

	// maximum attachment size is 32MB
	err = r.ParseMultipartForm(32 << 20)
	if err != nil {
		return "", nil, fmt.Errorf("error parsing multipart form: %v", err)
	}

	mimeType := r.PostFormValue("mimeType")
	file, _, err := r.FormFile("bytes")
	if err != nil {
		err = fmt.Errorf("error getting file from request: %v", err)
		return "", nil, err
	}
	defer file.Close()

	log.Printf("policyUid: %s, mimeType: %s", policyUid, mimeType)

	rawDoc, err := io.ReadAll(file)
	if err != nil {
		err = fmt.Errorf("error reading document from request: %v", err)
		return "", nil, err
	}

	filename := fmt.Sprintf("%s/%s/"+models.ContractDocumentFormat, "temp", policyUid, policy.NameDesc, policy.CodeCompany)
	gsLink, err := lib.PutToGoogleStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), filename, rawDoc)
	if err != nil {
		err = fmt.Errorf("error uploading document to GoogleBucket: %v", err)
		return "", nil, err
	}

	att := models.Attachment{
		Name:      models.ContractNonDigitalAttachmentName,
		Link:      gsLink,
		FileName:  filename,
		MimeType:  mimeType,
		IsPrivate: false,
		Section:   models.DocumentSectionContracts,
	}
	if policy.Attachments == nil {
		policy.Attachments = new([]models.Attachment)
	}
	*policy.Attachments = append(*policy.Attachments, att)

	flow := models.ProviderMgaFlow
	newStatus := models.PolicyStatusToPay
	newStatusHistory := []string{models.PolicyStatusManualSigned, models.PolicyStatusSign, models.PolicyStatusToPay}

	if policy.ProducerUid != "" {
		node := network.GetNetworkNodeByUid(policy.ProducerUid)
		if node != nil {
			flow = node.GetWarrant().GetFlowName(policy.Name)
		}
	}

	if flow == models.RemittanceMgaFlow {
		newStatus = models.PolicyStatusSign
		newStatusHistory = newStatusHistory[:len(newStatusHistory)-1]
	}

	policy.IsSign = true
	policy.Status = newStatus
	policy.StatusHistory = append(policy.StatusHistory, newStatusHistory...)
	policy.Updated = time.Now().UTC()

	// TODO: expire link namirial for signature

	err = lib.SetFirestoreErr(lib.PolicyCollection, policyUid, policy)
	if err != nil {
		return "", nil, err
	}

	policy.BigquerySave("")

	return "{}", "", nil
}
