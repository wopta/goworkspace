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
			log.ErrorF("error creating/updating user %s: %s", policy.Contractor.Uid,
				err.Error())
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

func AddFilesInPolicy(policy *models.Policy, origin string) error {
	if slices.Contains(policy.StatusHistory, models.PolicyStatusManualSigned) {
		return nil
	}
	documents, err := namirial.GetFiles(policy.IdSign)
	if err != nil {
		return err
	}
	if policy.Attachments == nil {
		policy.Attachments = &[]models.Attachment{}
	}
	for i := range documents.Documents {
		body, err := namirial.GetFile(documents.Documents[i].FileID)
		if err != nil {
			return err
		}
		fileName := documents.Documents[i].FileName
		fileName, _, _ = strings.Cut(fileName, ".")

		filePath := strings.ReplaceAll(fmt.Sprintf("temp/%s/%v", policy.Uid, documents.Documents[i].FileName), " ", "_")
		log.Println("path file path:", filePath)

		gsLink, err := lib.PutToGoogleStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), filePath, body)
		*policy.Attachments = append(*policy.Attachments, models.Attachment{
			Name:     fileName,
			Link:     gsLink,
			FileName: fileName,
			Section:  models.DocumentSectionContracts,
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
