package policy

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"gitlab.dev.wopta.it/goworkspace/document"
	"gitlab.dev.wopta.it/goworkspace/document/namirial"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/models"
)

func FillAttachments(policy *models.Policy, origin string) error {
	firePolicy := lib.PolicyCollection

	if policy.Attachments == nil {
		policy.Attachments = new([]models.Attachment)
	}
	policy.Updated = time.Now().UTC()

	return lib.SetFirestoreErr(firePolicy, policy.Uid, policy)
}

func Sign(policy *models.Policy, origin string) error {
	if !lib.SliceContains(policy.StatusHistory, models.PolicyStatusToSign) {
		return errors.New("policy has not been set to be signed")
	}

	firePolicy := lib.PolicyCollection

	policy.IsSign = true
	policy.Status = models.PolicyStatusSign
	policy.StatusHistory = append(policy.StatusHistory, models.PolicyStatusSign)
	policy.Updated = time.Now().UTC()

	return lib.SetFirestoreErr(firePolicy, policy.Uid, policy)
}

func SetToPay(policy *models.Policy, origin string) error {
	firePolicy := lib.PolicyCollection

	policy.Status = models.PolicyStatusToPay
	policy.StatusHistory = append(policy.StatusHistory, models.PolicyStatusToPay)
	policy.Updated = time.Now().UTC()

	return lib.SetFirestoreErr(firePolicy, policy.Uid, policy)
}

func promoteContractorDocumentsToUser(policy *models.Policy, origin string) error {
	var (
		tempPathFormat = "temp/%s/%s"
		userPathFormat = "assets/users/%s/%s"
	)
	log.AddPrefix("UpdateIdentityDocument")
	defer log.PopPrefix()

	for _, identityDocument := range policy.Contractor.IdentityDocuments {
		frontGsLink, err := lib.PromoteFile(
			fmt.Sprintf(tempPathFormat, policy.Uid, identityDocument.FrontMedia.FileName),
			fmt.Sprintf(userPathFormat, policy.Contractor.Uid, identityDocument.FrontMedia.FileName),
		)
		if err != nil {
			log.ErrorF("error saving front file: %s", err.Error())
			return err
		}
		identityDocument.FrontMedia.Link = frontGsLink

		if identityDocument.BackMedia != nil {
			backGsLink, err := lib.PromoteFile(
				fmt.Sprintf(tempPathFormat, policy.Uid, identityDocument.BackMedia.FileName),
				fmt.Sprintf(userPathFormat, policy.Contractor.Uid, identityDocument.BackMedia.FileName),
			)
			if err != nil {
				log.ErrorF("error saving back file: %s", err.Error())
				return err
			}
			identityDocument.BackMedia.Link = backGsLink
		}
	}
	policy.Updated = time.Now().UTC()

	firePolicy := lib.PolicyCollection

	return lib.SetFirestoreErr(firePolicy, policy.Uid, policy)
}

func SetUserIntoPolicyContractor(policy *models.Policy, origin string) error {
	log.AddPrefix("setUserIntoPolicyContractor")
	defer log.PopPrefix()
	log.Printf("Policy %s", policy.Uid)
	userUid, newUser, err := models.GetUserUIDByFiscalCode(origin, policy.Contractor.FiscalCode)
	if err != nil {
		log.ErrorF("error finding user: %s", err.Error())
		return err
	}

	policy.Contractor.Uid = userUid
	err = promoteContractorDocumentsToUser(policy, origin)
	if err != nil {
		log.ErrorF("error updating documents: %s", err.Error())
		return err
	}

	err = promotePolicyAttachments(policy, origin)
	if err != nil {
		log.ErrorF("error updating attachments: %s", err.Error())
		return err
	}

	if newUser {
		policy.Contractor.CreationDate = time.Now().UTC()
		policy.Contractor.UpdatedDate = policy.Contractor.CreationDate
		fireUsers := lib.UserCollection
		err = lib.SetFirestoreErr(fireUsers, userUid, policy.Contractor)
		if err != nil {
			log.ErrorF("error creating/updating user %s: %s", policy.Contractor.Uid, err.Error())
			return err
		}
		return policy.Contractor.BigquerySave(origin)
	}

	user := policy.Contractor.ToUser()
	if user == nil {
		return fmt.Errorf("invalid user")
	}
	_, err = models.UpdateUserByFiscalCode(origin, *user)
	return err
}

// Not sure if this is the right place
// because it creates a dependency with document
// DEPRECATED
func AddContract(policy *models.Policy, origin string) error {
	if slices.Contains(policy.StatusHistory, models.PolicyStatusManualSigned) {
		return nil
	}
	gsLink := <-document.GetFileV6(*policy, policy.Uid)
	filename := strings.ReplaceAll(fmt.Sprintf(models.ContractDocumentFormat, policy.NameDesc, policy.CodeCompany), " ", "_")
	*policy.Attachments = append(*policy.Attachments, models.Attachment{
		Name:     models.ContractAttachmentName,
		Link:     gsLink,
		FileName: filename,
		Section:  models.DocumentSectionContracts,
	})
	policy.Updated = time.Now().UTC()

	firePolicy := lib.PolicyCollection

	return lib.SetFirestoreErr(firePolicy, policy.Uid, policy)
}

// Download the signed file from the envelope and add them inside the policy's attachments, and save the policy.
// The name of the attachment is given by file's name sended to namirial(with the extension removed)
func AddSignedDocumentsInPolicy(policy *models.Policy, origin string) error {
	log.AddPrefix("AddDocumentsInPolicy")
	defer log.PopPrefix()
	if slices.Contains(policy.StatusHistory, models.PolicyStatusManualSigned) {
		return nil
	}
	if policy.IdSign == "" {
		return fmt.Errorf("No IdSign for policy with uid '%v'", policy.Uid)
	}
	documents, err := namirial.GetFiles(policy.IdSign)
	if err != nil {
		return err
	}
	if policy.Attachments == nil {
		policy.Attachments = &[]models.Attachment{}
	}
	if len(documents.Documents) == 0 {
		return fmt.Errorf("No document in envelope '%s'", policy.IdSign)
	}
	for i := range documents.Documents {
		body, err := namirial.GetFile(documents.Documents[i].FileID)
		if err != nil {
			return err
		}
		var typeFile string
		fileName := documents.Documents[i].FileName
		fileName, typeFile, _ = strings.Cut(fileName, ".")
		fileName = strings.ReplaceAll(fileName, "_", " ")

		filePath := strings.ReplaceAll(fmt.Sprintf("temp/%s/%v", policy.Uid, documents.Documents[i].FileName), " ", "_")
		log.Println("path file path:", filePath)

		gsLink, err := lib.PutToGoogleStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), filePath, body)
		//TODO: to remove eventually
		//With the new implementation of namirial we use the file's name to extract the label that will be showed in FE
		//Instead in the old one, the label was hardcode independently of 'NameDesc' (that happened to be the filename that we used for namirial)
		//olfImplementation of namirial: fw, err := w.CreateFormFile("file", NameDesc+" Polizza.pdf")'
		//So to allow retrocompatibility we use this, old file sended with old implementation
		if strings.Contains(fileName, policy.NameDesc) {
			fileName = models.ContractAttachmentName
		}

		*policy.Attachments = append(*policy.Attachments, models.Attachment{
			Name:        fileName,
			Link:        gsLink,
			FileName:    fileName,
			Section:     models.DocumentSectionContracts,
			ContentType: lib.GetContentType(typeFile),
			MimeType:    lib.GetContentType(typeFile),
		})

	}
	policy.Updated = time.Now().UTC()

	firePolicy := lib.PolicyCollection

	return lib.SetFirestoreErr(firePolicy, policy.Uid, policy)
}

func Pay(policy *models.Policy, origin string) error {
	firePolicy := lib.PolicyCollection

	policy.IsPay = true
	policy.Status = models.PolicyStatusPay
	policy.StatusHistory = append(policy.StatusHistory, models.PolicyStatusPay)
	policy.Updated = time.Now().UTC()

	return lib.SetFirestoreErr(firePolicy, policy.Uid, policy)
}

func promotePolicyAttachments(policy *models.Policy, origin string) error {
	const (
		tempPathFormat string = "temp/%s/%s"
		userPathFormat string = "assets/users/%s/%s"
	)
	log.AddPrefix("promotoPolicyAttachments")
	defer log.PopPrefix()
	if policy.Attachments == nil {
		return nil
	}

	for index, attachment := range *policy.Attachments {
		if !strings.HasPrefix(attachment.Link, "temp") {
			continue
		}
		gsLink, err := lib.PromoteFile(
			fmt.Sprintf(tempPathFormat, policy.Uid, attachment.FileName),
			fmt.Sprintf(userPathFormat, policy.Contractor.Uid, attachment.FileName),
		)
		if err != nil {
			log.Error(err)
			return err
		}
		(*policy.Attachments)[index].Link = gsLink
	}
	return nil
}

func AddProposalDoc(origin string, policy *models.Policy, networkNode *models.NetworkNode, mgaProduct *models.Product) error {
	fileGenerated, err := document.Proposal(origin, policy, networkNode, mgaProduct)
	if err != nil {
		log.Error(err)
		return err
	}
	if policy.Attachments == nil {
		policy.Attachments = new([]models.Attachment)
	}

	response, err := fileGenerated.Save()
	if err != nil {
		log.Error(err)
		return err
	}
	filename := strings.ReplaceAll(fmt.Sprintf(models.ProposalDocumentFormat, policy.NameDesc, policy.ProposalNumber), " ", "_")
	*policy.Attachments = append(*policy.Attachments, models.Attachment{
		Name:     models.ProposalAttachmentName,
		Link:     response.LinkGcs,
		FileName: filename,
		Section:  models.DocumentSectionContracts,
	})
	return nil
}
