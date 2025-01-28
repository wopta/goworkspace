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
	const pdfMimeType = "application/pdf"
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
	if mimeType != pdfMimeType {
		err = fmt.Errorf("cannot upload policy contract, invalid mime type: %s", mimeType)
		return "", nil, err
	}
	file, _, err := r.FormFile("bytes")
	if err != nil {
		err = fmt.Errorf("error getting file from request: %v", err)
		return "", nil, err
	}
	defer file.Close()

	rawDoc, err := io.ReadAll(file)
	if err != nil {
		err = fmt.Errorf("error reading document from request: %v", err)
		return "", nil, err
	}

	flow := models.ProviderMgaFlow
	pathPrefix := fmt.Sprintf("temp/%s/", policy.Uid)
	filename := fmt.Sprintf(models.ContractDocumentFormat, policy.NameDesc, policy.CodeCompany)
	newStatus := models.PolicyStatusToPay
	newStatusHistory := []string{models.PolicyStatusManualSigned, models.PolicyStatusSign, models.PolicyStatusToPay}

	if policy.ProducerUid != "" {
		node := network.GetNetworkNodeByUid(policy.ProducerUid)
		if node != nil {
			flow = node.GetWarrant().GetFlowName(policy.Name)
		}
	}

	if flow == models.RemittanceMgaFlow {
		pathPrefix = fmt.Sprintf("assets/users/%s/", policy.Contractor.Uid)
		newStatus = models.PolicyStatusSign
		newStatusHistory = newStatusHistory[:len(newStatusHistory)-1]
	}

	gsLink, err := lib.PutToGoogleStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), pathPrefix+filename, rawDoc)
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
