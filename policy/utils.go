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

func FillAttachments(policy *models.Policy) error {
	firePolicy := lib.PolicyCollection

	if policy.Attachments == nil {
		policy.Attachments = new([]models.Attachment)
	}
	policy.Updated = time.Now().UTC()

	return lib.SetFirestoreErr(firePolicy, policy.Uid, policy)
}

func Sign(policy *models.Policy) error {
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

func SetToPay(policy *models.Policy) error {
	firePolicy := lib.PolicyCollection

	policy.Status = models.PolicyStatusToPay
	policy.StatusHistory = append(policy.StatusHistory, models.PolicyStatusToPay)
	policy.Updated = time.Now().UTC()

	return lib.SetFirestoreErr(firePolicy, policy.Uid, policy)
}

func promoteContractorDocumentsToUser(policy *models.Policy) error {
	log.AddPrefix("UpdateIdentityDocument")
	defer log.PopPrefix()

	for _, identityDocument := range policy.Contractor.IdentityDocuments {
		if err := promoteIdentityDocument(policy, identityDocument); err != nil {
			return err
		}
	}
	policy.Updated = time.Now().UTC()

	firePolicy := lib.PolicyCollection

	return lib.SetFirestoreErr(firePolicy, policy.Uid, policy)
}

func PromotePolicy(policy *models.Policy) error {
	log.AddPrefix("promotePolicy")
	defer log.PopPrefix()
	log.Printf("Policy %s", policy.Uid)
	userUid, newUser, err := models.GetUserUIDByFiscalCode(policy.Contractor.FiscalCode)
	if err != nil {
		log.ErrorF("error finding user: %s", err.Error())
		return err
	}

	policy.Contractor.Uid = userUid
	err = promoteContractorDocumentsToUser(policy)
	if err != nil {
		log.ErrorF("error updating documents: %s", err.Error())
		return err
	}

	err = promotePolicyAttachments(policy)
	if err != nil {
		log.ErrorF("error updating attachments: %s", err.Error())
		return err
	}

	err = promoteNamirialDirectory(policy)
	if err != nil {
		log.ErrorF("error promoting namirial documents: %s", err.Error())
		return err
	}
	err = promoteAssets(policy)
	if err != nil {
		log.ErrorF("error promoting assent documents: %s", err.Error())
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
		return policy.Contractor.BigquerySave()
	}

	user := policy.Contractor.ToUser()
	if user == nil {
		return fmt.Errorf("invalid user")
	}
	_, err = models.UpdateUserByFiscalCode(*user)
	return err
}
func promoteAssets(policy *models.Policy) error {
	for _, assent := range policy.Assets {
		//TODO: to change when there will be more type of assents that has document
		if assent.Person == nil {
			continue
		}
		for _, identityDocument := range assent.Person.IdentityDocuments {
			if err := promoteIdentityDocument(policy, identityDocument); err != nil {
				return err
			}
		}
	}
	return nil
}
func RemoveTempPolicy(policy *models.Policy) error {
	tempPathFormat := "temp/%s"
	log.AddPrefix("RemoveTempPolicy")
	defer log.PopPrefix()
	log.Printf("Policy %s,path %v", policy.Uid, fmt.Sprintf(tempPathFormat, policy.Uid))

	return lib.RemoveFromGoogleStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), fmt.Sprintf(tempPathFormat, policy.Uid))
}

// Download the signed file from the envelope, add them inside the policy's attachments and save the policy.
// The name of the attachment is given by file's name sent to namirial(with the extension removed)
func AddSignedDocumentsInPolicy(policy *models.Policy) error {
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

		filePath := strings.ReplaceAll(fmt.Sprint("temp/", policy.Uid, "/", documents.Documents[i].FileName), " ", "_")
		log.Println("path file path:", filePath)

		gsLink, err := lib.PutToGoogleStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), filePath, body)

		fileName, _ = strings.CutPrefix(fileName, policy.NameDesc+" ")
		fileName, _ = strings.CutSuffix(fileName, " "+policy.CodeCompany)

		*policy.Attachments = append(*policy.Attachments, models.Attachment{
			Name:        fileName,
			Link:        gsLink,
			FileName:    documents.Documents[i].FileName,
			Section:     models.DocumentSectionContracts,
			ContentType: lib.GetContentType(typeFile),
			MimeType:    lib.GetContentType(typeFile),
		})

	}
	policy.Updated = time.Now().UTC()

	firePolicy := lib.PolicyCollection

	return lib.SetFirestoreErr(firePolicy, policy.Uid, policy)
}

func Pay(policy *models.Policy) error {
	firePolicy := lib.PolicyCollection

	policy.IsPay = true
	policy.Status = models.PolicyStatusPay
	policy.StatusHistory = append(policy.StatusHistory, models.PolicyStatusPay)
	policy.Updated = time.Now().UTC()

	return lib.SetFirestoreErr(firePolicy, policy.Uid, policy)
}

func promotePolicyAttachments(policy *models.Policy) error {
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
		path := attachment.Link
		path, _ = strings.CutPrefix(path, "gs://"+os.Getenv("GOOGLE_STORAGE_BUCKET")+"/")
		if !(strings.HasPrefix(path, "temp")) {
			continue
		}
		log.Printf("promoting %s", path)
		gsLink, err := lib.PromoteFile(
			path,
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

func promoteNamirialDirectory(policy *models.Policy) error {
	const (
		tempPathFormat string = "temp/%s/namirial"
		userPathFormat string = "assets/users/%s/namirial/%v"
	)
	log.AddPrefix("promotoNamirialDirectory")

	defer log.PopPrefix()
	listFilesPath, err := lib.ListGoogleStorageFolderContent(fmt.Sprintf(tempPathFormat, policy.Uid))
	if err != nil {
		return err
	}
	for _, path := range listFilesPath {
		split := strings.Split(path, "/")
		finalFullPath := fmt.Sprintf(userPathFormat, policy.Contractor.Uid, split[len(split)-1])
		_, err = lib.PromoteFile(path, finalFullPath)
		if err != nil {
			return err
		}
		log.Printf("Promoted %v", finalFullPath)
	}
	policy.DocumentName = fmt.Sprintf(userPathFormat, policy.Contractor.Uid, "")
	return err
}

func promoteIdentityDocument(policy *models.Policy, identityDocument *models.IdentityDocument) error {
	tempPathFormat := "temp/%s/%s"
	userPathFormat := "assets/users/%s/%s"

	if identityDocument.FrontMedia != nil {
		frontGsLink, err := lib.PromoteFile(
			fmt.Sprintf(tempPathFormat, policy.Uid, identityDocument.FrontMedia.FileName),
			fmt.Sprintf(userPathFormat, policy.Contractor.Uid, identityDocument.FrontMedia.FileName),
		)
		if err != nil {
			log.ErrorF("error saving front file: %s", err.Error())
			return err
		}
		identityDocument.FrontMedia.Link = frontGsLink
	}

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
	return nil
}

func AddProposalDoc(policy *models.Policy, networkNode *models.NetworkNode, mgaProduct *models.Product) error {
	fileGenerated, err := document.Proposal(policy, networkNode, mgaProduct)
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
		Name:        models.ProposalAttachmentName,
		Link:        response.LinkGcs,
		FileName:    filename,
		Section:     models.DocumentSectionContracts,
		MimeType:    lib.GetContentType("pdf"),
		ContentType: lib.GetContentType("pdf"),
	})
	return nil
}
