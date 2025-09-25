package mga

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"github.com/mohae/deepcopy"
	"gitlab.dev.wopta.it/goworkspace/document"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/compare"
	"gitlab.dev.wopta.it/goworkspace/models"
	plc "gitlab.dev.wopta.it/goworkspace/policy"
	usr "gitlab.dev.wopta.it/goworkspace/user"
)

func modifyPolicyFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err                                     error
		inputPolicy, modifiedPolicy, diffPolicy models.Policy
		hasDiff                                 bool
	)

	log.AddPrefix("ModifyPolicyFx")
	defer func() {
		r.Body.Close()
		if err != nil {
			log.ErrorF("error: %s", err.Error())
		}
		log.Println("Handler end ---------------------------------------------")
		log.PopPrefix()
	}()
	log.Println("Handler start -----------------------------------------------")

	log.Println("loading authToken from idToken...")
	token := r.Header.Get("Authorization")
	authToken, err := lib.GetAuthTokenFromIdToken(token)
	if err != nil {
		log.ErrorF("error getting authToken")
		return "", nil, err
	}
	log.Printf(
		"authToken - type: '%s' role: '%s' uid: '%s' email: '%s'",
		authToken.Type,
		authToken.Role,
		authToken.UserID,
		authToken.Email,
	)

	if err = json.NewDecoder(r.Body).Decode(&inputPolicy); err != nil {
		log.ErrorF("error decoding request body: %s", err.Error())
		return "", nil, err
	}
	inputPolicy.Normalize()

	log.Printf("fetching policy %s from Firestore...", inputPolicy.Uid)
	originalPolicy, err := plc.GetPolicy(inputPolicy.Uid)
	if err != nil {
		log.ErrorF("error fetching policy from Firestore: %s", err.Error())
		return "", nil, err
	}
	rawPolicy, err := json.Marshal(originalPolicy)
	if err != nil {
		log.ErrorF("error marshaling db policy: %s", err.Error())
		return "", nil, err
	}
	log.Printf("original policy: %s", string(rawPolicy))

	log.Printf("modifying policy...")
	modifiedPolicy, modifiedUser, err := modifyController(originalPolicy, inputPolicy)
	if err != nil {
		log.ErrorF("error during policy modification: %s", err.Error())
		return "", nil, err
	}
	log.Printf("policy %s modified successfully", modifiedPolicy.Uid)

	if diffPolicy, hasDiff = generateDiffPolicy(originalPolicy, modifiedPolicy); hasDiff {
		log.Println("generating addendum document for chages...")
		if addendumResp, err := document.Addendum(&diffPolicy); err == nil {
			res, err := addendumResp.Save()
			if err != nil {
				return "", nil, err
			}
			os.WriteFile("add.pdf", addendumResp.Bytes, 0777)
			addendumAtt := models.Attachment{
				Name:        "Appendice - Modifica dati di polizza",
				FileName:    addendumResp.FileName,
				ContentType: lib.GetContentType("pdf"),
				Link:        res.LinkGcs,
				IsPrivate:   false,
				Section:     models.DocumentSectionContracts,
				Note:        "",
			}
			if modifiedPolicy.Attachments == nil {
				modifiedPolicy.Attachments = new([]models.Attachment)
			}
			*modifiedPolicy.Attachments = append(*modifiedPolicy.Attachments, addendumAtt)
		} else if !errors.Is(err, document.ErrNotImplemented) {
			log.ErrorF("error generating addendum for policy %s: %s", inputPolicy.Uid, err.Error())
			return "", nil, err
		}
	}
	if err = writePolicyToDb(modifiedPolicy); err != nil {
		return "", nil, err
	}
	log.Printf("policy %s successfully saved", modifiedPolicy.Uid)
	if err = writeUserToDB(modifiedUser); err != nil {
		return "", nil, err
	}
	log.Printf("user %s modified successfully", modifiedUser.Uid)

	rawPolicy, err = json.Marshal(modifiedPolicy)

	return string(rawPolicy), modifiedPolicy, err
}

func generateDiffPolicy(originalPolicy, inputPolicy models.Policy) (models.Policy, bool) {
	var (
		diff            models.Policy
		diffContractor  models.Contractor
		diffAssets      []models.Asset
		hasPolicyVaried bool
	)
	if !reflect.DeepEqual(originalPolicy.Contractor, inputPolicy.Contractor) {
		hasPolicyVaried = true
		diffContractor = inputPolicy.Contractor
	}

	diffAssets = make([]models.Asset, len(inputPolicy.Assets))

	// TODO: missing asset uid for life - Enhance how assets are recognized
	for idx, asset := range inputPolicy.Assets {
		diffAssets[idx].Guarantees = make([]models.Guarante, 0, len(inputPolicy.Assets[idx].Guarantees))
		var modifiedAsset models.Asset
		if asset.Person != nil {
			if hasVaried := diffForUser(originalPolicy.Assets[idx].Person, inputPolicy.Assets[idx].Person); hasVaried {
				hasPolicyVaried = true
				modifiedAsset.Person = asset.Person
			}
		}
		for gIndex, guarantee := range asset.Guarantees {
			var (
				modifiedGuarantee            models.Guarante
				modifiedBeneficiary          *models.User
				modifiedBeneficiaryReference *models.User
				modifiedBeneficiaries        *[]models.Beneficiary
			)
			if hasVaried := diffForUser(originalPolicy.Assets[idx].Guarantees[gIndex].Beneficiary, inputPolicy.Assets[idx].Guarantees[gIndex].Beneficiary); hasVaried {
				hasPolicyVaried = true
				modifiedBeneficiary = guarantee.Beneficiary
			}
			if hasVaried := diffForUser(originalPolicy.Assets[idx].Guarantees[gIndex].BeneficiaryReference, inputPolicy.Assets[idx].Guarantees[gIndex].BeneficiaryReference); hasVaried {
				hasPolicyVaried = true
				if guarantee.BeneficiaryReference == nil {
					modifiedBeneficiaryReference = &models.User{}
				} else {
					modifiedBeneficiaryReference = guarantee.BeneficiaryReference
				}
			}
			if hasVaried := diffForBeneficiaries(originalPolicy.Assets[idx].Guarantees[gIndex].Beneficiaries, inputPolicy.Assets[idx].Guarantees[gIndex].Beneficiaries); hasVaried {
				hasPolicyVaried = true
				modifiedBeneficiaries = guarantee.Beneficiaries
			}

			modifiedGuarantee.Beneficiary = modifiedBeneficiary
			modifiedGuarantee.BeneficiaryReference = modifiedBeneficiaryReference
			modifiedGuarantee.Beneficiaries = modifiedBeneficiaries

			if !reflect.DeepEqual(models.Guarante{}, modifiedGuarantee) {
				modifiedGuarantee = guarantee
				modifiedGuarantee.Beneficiary = modifiedBeneficiary
				modifiedGuarantee.BeneficiaryReference = modifiedBeneficiaryReference
				modifiedGuarantee.Beneficiaries = modifiedBeneficiaries
				modifiedAsset.Guarantees = append(modifiedAsset.Guarantees, modifiedGuarantee)
			}
		}
		if !reflect.DeepEqual(models.Asset{}, modifiedAsset) {
			diffAssets[idx] = modifiedAsset
		}
	}
	var diffContractors []models.User
	if inputPolicy.Contractors != nil {
		if len(*inputPolicy.Contractors) != len(*originalPolicy.Contractors) {
			hasPolicyVaried = true
			diffContractors = *inputPolicy.Contractors
		} else {
			for idx, contr := range *inputPolicy.Contractors {
				if diffForUser(&contr, &(*originalPolicy.Contractors)[idx]) {
					diffContractors = append(diffContractors, contr)
					hasPolicyVaried = true
				}
			}
		}
	}

	if hasPolicyVaried {
		diff = originalPolicy
		diff.Contractor = diffContractor
		diff.Contractors = &diffContractors
		diff.Assets = diffAssets
	}

	return diff, hasPolicyVaried
}

func diffForUser(orig, input *models.User) bool {
	compareFunc := func(a, b *models.User) bool {
		if a.FiscalCode != b.FiscalCode {
			return false
		}

		if a.Mail != b.Mail {
			return false
		}

		if a.Phone != b.Phone {
			return false
		}

		if a.Name != b.Name {
			return false
		}

		if a.Surname != b.Surname {
			return false
		}
		if !compare.AreEqual(a.Residence, b.Residence) {
			return false
		}

		if !compare.AreEqual(a.Domicile, b.Domicile) {
			return false
		}
		return true

	}
	return !compare.AreEqualFunc(orig, input, compareFunc)
}

func diffForBeneficiaries(orig, input *[]models.Beneficiary) bool {
	compareFunc := func(a, b models.Beneficiary) bool {
		if a.Name != b.Name {
			return false
		}
		if a.Surname != b.Surname {
			return false
		}
		if a.Mail != b.Mail {
			return false
		}
		if a.Phone != b.Phone {
			return false
		}
		if a.FiscalCode != b.FiscalCode {
			return false
		}
		if a.VatCode != b.VatCode {
			return false
		}
		if !compare.AreEqual(a.Residence, b.Residence) {
			return false
		}
		if !compare.AreEqual(a.CompanyAddress, b.CompanyAddress) {
			return false
		}
		if a.IsFamilyMember != b.IsFamilyMember {
			return false
		}
		if a.IsContactable != b.IsContactable {
			return false
		}
		if a.IsLegitimateSuccessors != b.IsLegitimateSuccessors {
			return false
		}
		if a.BeneficiaryType != b.BeneficiaryType {
			return false
		}
		return true
	}
	return !compare.AreSlicesEqualFunc(orig, input, compareFunc)
}

func modifyController(originalPolicy, inputPolicy models.Policy) (models.Policy, models.User, error) {
	var (
		err            error
		policyToModify models.Policy
		modifiedUser   models.User
	)

	policyToModify = deepcopy.Copy(originalPolicy).(models.Policy)

	err = checkEmailUniqueness(originalPolicy.Contractor, inputPolicy.Contractor)
	if err != nil {
		return models.Policy{}, models.User{}, err
	}

	policyToModify.Contractor, err = modifyContractorInfo(inputPolicy.Contractor, originalPolicy.Contractor)
	if err != nil {
		return models.Policy{}, models.User{}, err
	}

	policyToModify.Assets, err = modifyAssets(policyToModify, inputPolicy)
	if err != nil {
		return models.Policy{}, models.User{}, err
	}
	policyToModify.Contractors, err = modifyContractorsInfo(inputPolicy, policyToModify)
	if err != nil {
		return models.Policy{}, models.User{}, err
	}
	if policyToModify.Contractor.Uid != "" {
		tmpUser := policyToModify.Contractor.ToUser()
		modifiedUser, err = modifyUserInfo(*tmpUser)
		if err != nil {
			return models.Policy{}, models.User{}, err
		}
	}

	return policyToModify, modifiedUser, err
}

func checkEmailUniqueness(originalContractor, inputContractor models.Contractor) error {
	if strings.EqualFold(originalContractor.Mail, inputContractor.Mail) {
		return nil
	}

	iterator := lib.WhereFirestore(models.UserCollection, "mail", "==", inputContractor.Mail)
	users := models.UsersToListData(iterator)

	for _, usr := range users {
		if !strings.EqualFold(usr.FiscalCode, inputContractor.FiscalCode) {
			return errors.New("mail duplicated")
		}
	}

	return nil
}

func modifyContractorsInfo(input, original models.Policy) (*[]models.User, error) {
	var res []models.User
	if original.Contractors == nil {
		return nil, nil
	}
	for i := range *original.Contractors {
		//if contractor is ditta individuale the signer is the same person
		if input.Contractor.Type == models.UserIndividual && (*input.Contractors)[i].IsSignatory {
			newContractor, err := modifyContractorInfo(input.Contractor, original.Contractor)
			if err != nil {
				return nil, err
			}
			user := *newContractor.ToUser()
			user.IsSignatory = true
			res = append(res, user)
			continue
		}
		newContractor, err := modifyPersonInfo((*input.Contractors)[i], (*original.Contractors)[i])
		if err != nil {
			return nil, err
		}
		res = append(res, *newContractor)
	}
	return &res, nil
}

func modifyContractorInfo(inputContractor, originalContractor models.Contractor) (models.Contractor, error) {
	var (
		err                error
		modifiedContractor = new(models.Contractor)
	)

	log.Printf("modifying contractor %s info...", originalContractor.Uid)
	*modifiedContractor = originalContractor

	modifiedContractor.Name = inputContractor.Name
	modifiedContractor.Surname = inputContractor.Surname
	modifiedContractor.BirthDate = inputContractor.BirthDate
	modifiedContractor.BirthCity = inputContractor.BirthCity
	modifiedContractor.BirthProvince = inputContractor.BirthProvince
	modifiedContractor.Gender = inputContractor.Gender
	modifiedContractor.FiscalCode = inputContractor.FiscalCode
	modifiedContractor.Mail = inputContractor.Mail
	modifiedContractor.Residence = inputContractor.Residence
	modifiedContractor.CompanyAddress = inputContractor.CompanyAddress
	modifiedContractor.CompanyName = inputContractor.CompanyName
	modifiedContractor.VatCode = inputContractor.VatCode
	//TODO: check if it is correct
	modifiedContractor.Ateco = inputContractor.Ateco
	if inputContractor.Consens != nil {
		if modifiedContractor.Consens == nil {
			modifiedContractor.Consens = inputContractor.Consens
		} else {
			for index, consensus := range *modifiedContractor.Consens {
				for _, inputConsensus := range *inputContractor.Consens {
					if inputConsensus.Key == consensus.Key {
						(*modifiedContractor.Consens)[index].Answer = inputConsensus.Answer
					}
				}
			}
		}
	}

	user := modifiedContractor.ToUser()
	if originalContractor.Type == models.UserIndividual {
		return *modifiedContractor, err
	}
	if user.FiscalCode != "" {
		err = usr.CheckFiscalCode(*user)
		if err != nil {
			return models.Contractor{}, err
		}
	}

	log.Printf("contractor %s modified", originalContractor.Uid)

	return *modifiedContractor, err
}

func modifyAssets(modifiedPolicy models.Policy, inputPolicy models.Policy) ([]models.Asset, error) {
	assets := make([]models.Asset, 0, len(inputPolicy.Assets))

	for index, asset := range inputPolicy.Assets {
		var (
			err           error
			modifiedAsset models.Asset
		)

		modifiedAsset = modifiedPolicy.Assets[index]

		if asset.Person != nil {
			if asset.Person.FiscalCode == modifiedPolicy.Contractor.FiscalCode {
				modifiedAsset.Person = modifiedPolicy.Contractor.ToUser()
			} else {
				modifiedAsset.Person, err = modifyPersonInfo(*inputPolicy.Assets[0].Person, *modifiedPolicy.Assets[0].Person)
				if err != nil {
					return nil, err
				}
			}

			if asset.Guarantees != nil {
				modifiedAsset.Guarantees, err = modifyBeneficiaryInfo(inputPolicy.Assets[0].Guarantees, modifiedPolicy.Assets[0].Guarantees)
				if err != nil {
					return nil, err
				}
			}
		}

		assets = append(assets, modifiedAsset)
	}
	return assets, nil
}

func modifyBeneficiaryInfo(inputGuarantees, originalGuarantees []models.Guarante) ([]models.Guarante, error) {
	log.Println("modifying beneficiary info...")
	modifiedGuarantees := deepcopy.Copy(&originalGuarantees).(*[]models.Guarante)

	for i, g := range inputGuarantees {
		if originalGuarantees[i].Beneficiaries != nil && (g.Beneficiaries == nil || (len(*g.Beneficiaries) == 0 && len(*originalGuarantees[i].Beneficiaries) > 0)) {
			return nil, fmt.Errorf("must have at least one beneficiary")
		}
		(*modifiedGuarantees)[i].BeneficiaryReference = g.BeneficiaryReference
		(*modifiedGuarantees)[i].Beneficiaries = g.Beneficiaries
	}

	return *modifiedGuarantees, nil
}

func modifyPersonInfo(inputPerson, originalPerson models.User) (*models.User, error) {
	var (
		err             error
		modifiedInsured = new(models.User)
	)

	*modifiedInsured = originalPerson

	modifiedInsured.Name = inputPerson.Name
	modifiedInsured.Surname = inputPerson.Surname
	modifiedInsured.BirthDate = inputPerson.BirthDate
	modifiedInsured.BirthCity = inputPerson.BirthCity
	modifiedInsured.BirthProvince = inputPerson.BirthProvince
	modifiedInsured.Gender = inputPerson.Gender
	modifiedInsured.FiscalCode = inputPerson.FiscalCode
	modifiedInsured.Mail = inputPerson.Mail
	modifiedInsured.Residence = inputPerson.Residence
	if inputPerson.Consens != nil {
		if modifiedInsured.Consens == nil {
			modifiedInsured.Consens = inputPerson.Consens
		} else {
			for index, consensus := range *modifiedInsured.Consens {
				for _, inputConsensus := range *inputPerson.Consens {
					if inputConsensus.Key == consensus.Key {
						(*modifiedInsured.Consens)[index].Answer = inputConsensus.Answer
					}
				}
			}
		}
	}
	err = usr.CheckFiscalCode(*modifiedInsured)
	if err != nil {
		return nil, err
	}

	log.Printf("person %s modified", originalPerson.Uid)

	return modifiedInsured, err
}

func modifyUserInfo(inputUser models.User) (models.User, error) {
	var (
		err                  error
		dbUser, modifiedUser models.User
	)

	log.Println("modifying user info...")

	docsnap, err := lib.GetFirestoreErr(models.UserCollection, inputUser.Uid)
	if err != nil {
		log.ErrorF("error retrieving user %s from Firestore: %s", inputUser.Uid, err.Error())
		return models.User{}, err
	}
	docsnap.DataTo(&dbUser)

	modifiedUser = dbUser

	modifiedUser.Name = inputUser.Name
	modifiedUser.Surname = inputUser.Surname
	modifiedUser.BirthDate = inputUser.BirthDate
	modifiedUser.BirthCity = inputUser.BirthCity
	modifiedUser.BirthProvince = inputUser.BirthProvince
	modifiedUser.Gender = inputUser.Gender
	modifiedUser.FiscalCode = inputUser.FiscalCode
	modifiedUser.Mail = inputUser.Mail
	modifiedUser.Residence = inputUser.Residence
	if inputUser.Consens != nil {
		if modifiedUser.Consens == nil {
			modifiedUser.Consens = inputUser.Consens
		} else {
			for index, consensus := range *modifiedUser.Consens {
				for _, inputConsensus := range *inputUser.Consens {
					if inputConsensus.Key == consensus.Key {
						(*modifiedUser.Consens)[index].Answer = inputConsensus.Answer
					}
				}
			}
		}
	}
	if !strings.EqualFold(modifiedUser.Mail, dbUser.Mail) && dbUser.AuthId != "" {
		log.Printf("modifying user %s email from %s to %s...", modifiedUser.Uid, modifiedUser.Mail, dbUser.Mail)
		_, err = lib.UpdateUserEmail(modifiedUser.Uid, modifiedUser.Mail)
		if err != nil {
			log.ErrorF("error modifying authentication email: %s", err.Error())
			return models.User{}, err
		}
		log.Printf("mail modified successfully")
	}

	log.Printf("user %s modified successfully", modifiedUser.Uid)

	return modifiedUser, err
}

func writePolicyToDb(modifiedPolicy models.Policy) error {
	var err error

	modifiedPolicy.Updated = time.Now().UTC()

	log.Printf("writing policy %s to DBs...", modifiedPolicy.Uid)

	err = lib.SetFirestoreErr(models.PolicyCollection, modifiedPolicy.Uid, modifiedPolicy)
	if err != nil {
		log.ErrorF("error writing modified policy to Firestore: %s", err.Error())
		return err
	}

	modifiedPolicy.BigquerySave()

	log.Printf("policy %s written into DBs", modifiedPolicy.Uid)

	return err
}

func writeUserToDB(user models.User) error {
	var err error

	if user.Uid == "" {
		return nil
	}

	user.UpdatedDate = time.Now().UTC()

	log.Printf("writing user %s to DBs...", user.Uid)

	err = lib.SetFirestoreErr(models.UserCollection, user.Uid, user)
	if err != nil {
		log.ErrorF("error writing modified user to Firestore: %s", err.Error())
		return err
	}

	err = user.BigquerySave()
	if err != nil {
		log.ErrorF("error writing modified user to BigQuery: %s", err.Error())
		return err
	}

	log.Printf("user %s written into DBs", user.Uid)

	return err
}
