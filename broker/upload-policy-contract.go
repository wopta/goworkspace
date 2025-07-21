package broker

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"github.com/go-chi/chi/v5"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/network"
	plc "gitlab.dev.wopta.it/goworkspace/policy"
)

func UploadPolicyContractFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	const pdfMimeType = "application/pdf"
	var (
		err    error
		policy models.Policy
	)

	log.AddPrefix("UploadPolicyContractFx")
	defer func() {
		r.Body.Close()
		if err != nil {
			log.ErrorF("error: %s", err.Error())
		}
		log.Println("Handler end ----------------------------------------------")
		log.PopPrefix()
	}()

	log.Println("Handler start -----------------------------------------------")

	policyUid := chi.URLParam(r, "uid")

	policy, err = plc.GetPolicy(policyUid)
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

	note := r.PostFormValue("note")

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

	flow := models.ECommerceFlow
	pathPrefix := fmt.Sprintf("temp/%s/", policy.Uid)
	filename := fmt.Sprintf(models.ContractDocumentFormat, policy.NameDesc, policy.CodeCompany)
	newStatus := models.PolicyStatusToPay
	newStatusHistory := []string{models.PolicyStatusManualSigned, models.PolicyStatusSign, models.PolicyStatusToPay}

	if policy.Channel == models.NetworkChannel {
		var node *models.NetworkNode
		if node = network.GetNetworkNodeByUid(policy.ProducerUid); node == nil {
			err = fmt.Errorf("error getting node %s", policy.ProductUid)
			return "", nil, err
		}
		flow = node.GetWarrant().GetFlowName(policy.Name)
	}

	if flow == models.RemittanceMgaFlow {
		pathPrefix = fmt.Sprintf("assets/users/%s/", policy.Contractor.Uid)
		newStatus = models.PolicyStatusSign
		newStatusHistory = newStatusHistory[:len(newStatusHistory)-1]
	}

	filePath := pathPrefix + filename

	if _, err = lib.PutToGoogleStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), filePath, rawDoc); err != nil {
		err = fmt.Errorf("error uploading document to GoogleBucket: %v", err)
		return "", nil, err
	}

	att := models.Attachment{
		Name:      models.ContractNonDigitalAttachmentName,
		Link:      filePath,
		FileName:  filename,
		MimeType:  mimeType,
		IsPrivate: false,
		Section:   models.DocumentSectionContracts,
		Note:      note,
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

	if err = lib.SetFirestoreErr(lib.PolicyCollection, policyUid, policy); err != nil {
		return "", nil, err
	}

	policy.BigquerySave()

	return "{}", "", nil
}
