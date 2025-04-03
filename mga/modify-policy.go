package mga

import (
	"encoding/json"
	"errors"
	"github.com/wopta/goworkspace/lib/log"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/mohae/deepcopy"
	"github.com/wopta/goworkspace/document"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	plc "github.com/wopta/goworkspace/policy"
	usr "github.com/wopta/goworkspace/user"
)

func ModifyPolicyFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err                                     error
		inputPolicy, modifiedPolicy, diffPolicy models.Policy
		hasDiff                                 bool
	)

	log.AddPrefix("[ModifyPolicyFx] ")
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
	originalPolicy, err := plc.GetPolicy(inputPolicy.Uid, "")
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
			addendumAtt := models.Attachment{
				Name:      "Appendice - Modifica dati di polizza",
				FileName:  addendumResp.Filename,
				MimeType:  "application/pdf",
				Link:      addendumResp.LinkGcs,
				IsPrivate: false,
				Section:   models.DocumentSectionContracts,
				Note:      "",
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
		hasPolicyVaried bool
	)

	if hasVaried := diffForUser(originalPolicy.Contractor.ToUser(), inputPolicy.Contractor.ToUser()); hasVaried {
		hasPolicyVaried = true
		diff.Contractor = inputPolicy.Contractor
	}

	diff.Assets = make([]models.Asset, len(inputPolicy.Assets), len(inputPolicy.Assets))

	// TODO: missing asset uid for life - Enhance how assets are recognized
	for idx, asset := range inputPolicy.Assets {
		diff.Assets[idx].Guarantees = make([]models.Guarante, 0, len(inputPolicy.Assets[idx].Guarantees))
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
			diff.Assets[idx] = modifiedAsset
		}
	}

	if hasPolicyVaried {
		diff.Uid = originalPolicy.Uid
		diff.Name = originalPolicy.Name
		diff.NameDesc = originalPolicy.NameDesc
		diff.ProposalNumber = originalPolicy.ProposalNumber
		diff.StartDate = originalPolicy.StartDate
		diff.EndDate = originalPolicy.EndDate
		diff.Company = originalPolicy.Company
	}

	return diff, hasPolicyVaried
}

func diffForUser(orig, input *models.User) bool {
	if orig == nil && input == nil {
		return false
	}

	if (orig == nil && input != nil) || (orig != nil && input == nil) {
		return true
	}

	if orig.FiscalCode != input.FiscalCode {
		return true
	}

	if orig.Mail != input.Mail {
		return true
	}

	if orig.Phone != input.Phone {
		return true
	}

	if orig.Name != input.Name {
		return true
	}

	if orig.Surname != input.Surname {
		return true
	}

	if input.Residence == nil || orig.Residence == nil {
		return input.Residence != orig.Residence
	}

	if input.Domicile == nil || orig.Domicile == nil {
		return input.Domicile != orig.Domicile
	}

	return false
}

func diffForBeneficiaries(orig, input *[]models.Beneficiary) bool {
	if orig == nil && input == nil {
		return false
	}
	if (orig == nil && input != nil) || (orig != nil && input == nil) {
		return true
	}
	if len(*orig) != len(*input) {
		return true
	}

	for idx := range *orig {
		if hasChanged := diffForBeneficiary((*orig)[idx], (*input)[idx]); hasChanged {
			return true
		}
	}

	return false
}

func diffForBeneficiary(orig, input models.Beneficiary) bool {
	if orig.Name != input.Name {
		return true
	}
	if orig.Surname != input.Surname {
		return true
	}
	if orig.Mail != input.Mail {
		return true
	}
	if orig.Phone != input.Phone {
		return true
	}
	if orig.FiscalCode != input.FiscalCode {
		return true
	}
	if orig.VatCode != input.VatCode {
		return true
	}
	if orig.Residence == nil || input.Residence == nil {
		return orig.Residence != input.Residence
	}
	if orig.CompanyAddress == nil || input.CompanyAddress == nil {
		return orig.CompanyAddress != input.CompanyAddress
	}
	if orig.IsFamilyMember != input.IsFamilyMember {
		return true
	}
	if orig.IsContactable != input.IsContactable {
		return true
	}
	if orig.IsLegitimateSuccessors != input.IsLegitimateSuccessors {
		return true
	}
	if orig.BeneficiaryType != input.BeneficiaryType {
		return true
	}
	return false
}

func modifyController(originalPolicy, inputPolicy models.Policy) (models.Policy, models.User, error) {
	var (
		err            error
		modifiedPolicy models.Policy
		modifiedUser   models.User
	)

	modifiedPolicy = deepcopy.Copy(originalPolicy).(models.Policy)

	err = checkEmailUniqueness(originalPolicy.Contractor, inputPolicy.Contractor)
	if err != nil {
		return models.Policy{}, models.User{}, err
	}

	modifiedPolicy.Contractor, err = modifyContractorInfo(inputPolicy.Contractor, originalPolicy.Contractor)
	if err != nil {
		return models.Policy{}, models.User{}, err
	}

	if modifiedPolicy.Contractor.Uid != "" {
		tmpUser := modifiedPolicy.Contractor.ToUser()
		modifiedUser, err = modifyUserInfo(*tmpUser)
		if err != nil {
			return models.Policy{}, models.User{}, err
		}
	}

	modifiedPolicy.Assets, err = modifyAssets(modifiedPolicy, inputPolicy)
	if err != nil {
		return models.Policy{}, models.User{}, err
	}

	return modifiedPolicy, modifiedUser, err
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
	err = usr.CheckFiscalCode(*user)
	if err != nil {
		return models.Contractor{}, err
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
				modifiedAsset.Person, err = modifyInsuredInfo(*inputPolicy.Assets[0].Person, *modifiedPolicy.Assets[0].Person)
				if err != nil {
					return nil, err
				}
			}

			if asset.Guarantees != nil {
				modifiedAsset.Guarantees = modifyBeneficiaryInfo(inputPolicy.Assets[0].Guarantees, modifiedPolicy.Assets[0].Guarantees)
			}
		}

		assets = append(assets, modifiedAsset)
	}
	return assets, nil
}

func modifyBeneficiaryInfo(inputGuarantees, originalGuarantees []models.Guarante) []models.Guarante {
	var (
		modifiedGuarantees = new([]models.Guarante)
	)

	log.Println("modifying beneficiary info...")
	modifiedGuarantees = &originalGuarantees

	for i, g := range inputGuarantees {
		(*modifiedGuarantees)[i].BeneficiaryReference = g.BeneficiaryReference
		(*modifiedGuarantees)[i].Beneficiaries = g.Beneficiaries
	}

	return *modifiedGuarantees
}

func modifyInsuredInfo(inputInsured, originalInsured models.User) (*models.User, error) {
	var (
		err             error
		modifiedInsured = new(models.User)
	)

	log.Println("modifying insured info...")
	*modifiedInsured = originalInsured

	modifiedInsured.Name = inputInsured.Name
	modifiedInsured.Surname = inputInsured.Surname
	modifiedInsured.BirthDate = inputInsured.BirthDate
	modifiedInsured.BirthCity = inputInsured.BirthCity
	modifiedInsured.BirthProvince = inputInsured.BirthProvince
	modifiedInsured.Gender = inputInsured.Gender
	modifiedInsured.FiscalCode = inputInsured.FiscalCode
	modifiedInsured.Mail = inputInsured.Mail
	modifiedInsured.Residence = inputInsured.Residence
	if inputInsured.Consens != nil {
		if modifiedInsured.Consens == nil {
			modifiedInsured.Consens = inputInsured.Consens
		} else {
			for index, consensus := range *modifiedInsured.Consens {
				for _, inputConsensus := range *inputInsured.Consens {
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

	log.Printf("insured %s modified", originalInsured.Uid)

	log.Println("insured modified successfully")
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

	modifiedPolicy.BigquerySave("")

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

	err = user.BigquerySave("")
	if err != nil {
		log.ErrorF("error writing modified user to BigQuery: %s", err.Error())
		return err
	}

	log.Printf("user %s written into DBs", user.Uid)

	return err
}
