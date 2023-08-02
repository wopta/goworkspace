package policy

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/wopta/goworkspace/document"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func FillAttachments(policy *models.Policy, origin string) error {
	firePolicy := lib.GetDatasetByEnv(origin, models.PolicyCollection)

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

	firePolicy := lib.GetDatasetByEnv(origin, models.PolicyCollection)

	policy.IsSign = true
	policy.Status = models.PolicyStatusSign
	policy.StatusHistory = append(policy.StatusHistory, models.PolicyStatusSign)
	policy.Updated = time.Now().UTC()

	return lib.SetFirestoreErr(firePolicy, policy.Uid, policy)
}

func SetToPay(policy *models.Policy, origin string) error {
	firePolicy := lib.GetDatasetByEnv(origin, models.PolicyCollection)

	policy.Status = models.PolicyStatusToPay
	policy.StatusHistory = append(policy.StatusHistory, models.PolicyStatusToPay)
	policy.Updated = time.Now().UTC()

	return lib.SetFirestoreErr(firePolicy, policy.Uid, policy)
}

func promoteContractorDocumentsToUser(policy *models.Policy, origin string) error {
	var (
		tempPathFormat string = "temp/%s/%s"
		userPathFormat string = "assets/users/%s/%s"
	)

	for _, identityDocument := range policy.Contractor.IdentityDocuments {
		frontGsLink, err := lib.PromoteFile(
			fmt.Sprintf(tempPathFormat, policy.Uid, identityDocument.FrontMedia.FileName),
			fmt.Sprintf(userPathFormat, policy.Contractor.Uid, identityDocument.FrontMedia.FileName),
		)
		if err != nil {
			log.Printf("[updateIdentityDocument] ERROR saving front file: %s", err.Error())
			return err
		}
		identityDocument.FrontMedia.Link = frontGsLink

		if identityDocument.BackMedia != nil {
			backGsLink, err := lib.PromoteFile(
				fmt.Sprintf(tempPathFormat, policy.Uid, identityDocument.BackMedia.FileName),
				fmt.Sprintf(userPathFormat, policy.Contractor.Uid, identityDocument.BackMedia.FileName),
			)
			if err != nil {
				log.Printf("[updateIdentityDocument] ERROR saving back file: %s", err.Error())
				return err
			}
			identityDocument.BackMedia.Link = backGsLink
		}
	}
	policy.Updated = time.Now().UTC()

	firePolicy := lib.GetDatasetByEnv(origin, models.PolicyCollection)

	return lib.SetFirestoreErr(firePolicy, policy.Uid, policy)
}

func SetUserIntoPolicyContractor(policy *models.Policy, origin string) error {
	log.Printf("[setUserIntoPolicyContractor] Policy %s", policy.Uid)
	userUid, newUser, err := models.GetUserUIDByFiscalCode(origin, policy.Contractor.FiscalCode)
	if err != nil {
		log.Printf("[setUserIntoPolicyContractor] ERROR finding user: %s", err.Error())
		return err
	}

	policy.Contractor.Uid = userUid
	err = promoteContractorDocumentsToUser(policy, origin)
	if err != nil {
		log.Printf("[setUserIntoPolicyContractor] ERROR updating documents: %s", err.Error())
		return err
	}

	if newUser {
		policy.Contractor.CreationDate = time.Now().UTC()
		fireUsers := lib.GetDatasetByEnv(origin, "users")
		err = lib.SetFirestoreErr(fireUsers, userUid, policy.Contractor)
	} else {
		_, err = models.UpdateUserByFiscalCode(origin, policy.Contractor)
	}

	if err != nil {
		log.Printf("[setUserIntoPolicyContractor] ERROR creating/updating user %s: %s", policy.Contractor.Uid,
			err.Error())
		return err
	}

	return policy.Contractor.BigquerySave(origin)
}

// Not sure if this is the right place
// because it creates a dependency with document
func AddContract(policy *models.Policy, origin string) error {
	// Get Policy contract
	gsLink := <-document.GetFileV6(*policy, policy.Uid)
	// Add Contract
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	filenameParts := []string{"Contratto", policy.NameDesc, timestamp, ".pdf"}
	filename := strings.Join(filenameParts, "_")
	filename = strings.ReplaceAll(filename, " ", "_")
	*policy.Attachments = append(*policy.Attachments, models.Attachment{
		Name:     "Contratto",
		Link:     gsLink,
		FileName: filename,
	})
	policy.Updated = time.Now().UTC()

	firePolicy := lib.GetDatasetByEnv(origin, models.PolicyCollection)

	return lib.SetFirestoreErr(firePolicy, policy.Uid, policy)
}

func Pay(policy *models.Policy, origin string) error {
	firePolicy := lib.GetDatasetByEnv(origin, models.PolicyCollection)

	policy.IsPay = true
	policy.Status = models.PolicyStatusPay
	policy.StatusHistory = append(policy.StatusHistory, models.PolicyStatusPay)
	policy.Updated = time.Now().UTC()

	return lib.SetFirestoreErr(firePolicy, policy.Uid, policy)
}
