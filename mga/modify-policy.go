package mga

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/wopta/goworkspace/document"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	plc "github.com/wopta/goworkspace/policy"
	usr "github.com/wopta/goworkspace/user"
)

func ModifyPolicyFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err                         error
		inputPolicy, modifiedPolicy models.Policy
	)

	log.SetPrefix("[ModifyPolicyFx] ")
	defer log.SetPrefix("")

	log.Println("Handler start -----------------------------------------------")

	log.Println("loading authToken from idToken...")

	token := r.Header.Get("Authorization")
	authToken, err := lib.GetAuthTokenFromIdToken(token)
	if err != nil {
		log.Printf("error getting authToken")
		return "{}", nil, err
	}
	log.Printf(
		"authToken - type: '%s' role: '%s' uid: '%s' email: '%s'",
		authToken.Type,
		authToken.Role,
		authToken.UserID,
		authToken.Email,
	)

	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	err = json.Unmarshal(body, &inputPolicy)
	if err != nil {
		log.Printf("error unmarshaling request body: %s", err.Error())
		return "{}", nil, err
	}

	inputPolicy.Normalize()

	log.Printf("fetching policy %s from Firestore...", inputPolicy.Uid)
	originalPolicy, err := plc.GetPolicy(inputPolicy.Uid, "")
	if err != nil {
		log.Printf("error fetching policy from Firestore: %s", err.Error())
		return "{}", nil, err
	}
	rawPolicy, err := json.Marshal(originalPolicy)
	if err != nil {
		log.Printf("error marshaling db policy: %s", err.Error())
		return "{}", nil, err
	}
	log.Printf("original policy: %s", string(rawPolicy))

	log.Printf("modifying policy...")
	modifiedPolicy, modifiedUser, err := modifyController(originalPolicy, inputPolicy)
	if err != nil {
		log.Printf("error during policy modification: %s", err.Error())
		return "{}", nil, err
	}
	log.Printf("policy %s modified successfully", modifiedPolicy.Uid)
	err = writePolicyToDb(modifiedPolicy)
	if err != nil {
		return "{}", nil, err
	}
	log.Printf("policy %s successfully saved", modifiedPolicy.Uid)
	err = writeUserToDB(modifiedUser)
	if err != nil {
		return "{}", nil, err
	}
	log.Printf("user %s modified successfully", modifiedUser.Uid)

	diffPolicy, err := generateDiffPolicy(originalPolicy, inputPolicy)
	if err != nil {
		return "", models.Policy{}, err
	}
	p := <-document.AddendumObj("", diffPolicy, nil, nil)
	log.Println(p.LinkGcs)

	rawPolicy, err = json.Marshal(modifiedPolicy)

	log.Println("Handler end -------------------------------------------------")

	return string(rawPolicy), modifiedPolicy, err
}

func generateDiffPolicy(originalPolicy, inputPolicy models.Policy) (models.Policy, error) {
	var diff models.Policy
	diff.Uid = originalPolicy.Uid
	diff.Name = originalPolicy.Name
	diff.NameDesc = originalPolicy.NameDesc
	diff.ProposalNumber = originalPolicy.ProposalNumber
	diff.StartDate = originalPolicy.StartDate
	diff.EndDate = originalPolicy.EndDate
	diff.Company = originalPolicy.Company

	diff.Contractor = diffForContractor(originalPolicy.Contractor, inputPolicy.Contractor)
	assets := make([]models.Asset, 0)
	ass := models.Asset{}

	gar := models.Guarante{}
	if inputPolicy.Assets[0].Person != nil {
		ass.Person = diffForUser(originalPolicy.Assets[0].Person, inputPolicy.Assets[0].Person)
	}

	if inputPolicy.Assets[0].Guarantees[0].Beneficiary != nil {
		gar.Beneficiary = diffForUser(originalPolicy.Assets[0].Guarantees[0].Beneficiary, inputPolicy.Assets[0].Guarantees[0].Beneficiary)
	}
	if inputPolicy.Assets[0].Guarantees[0].BeneficiaryReference != nil {
		gar.BeneficiaryReference = diffForUser(originalPolicy.Assets[0].Guarantees[0].BeneficiaryReference, inputPolicy.Assets[0].Guarantees[0].BeneficiaryReference)
	}
	if inputPolicy.Assets[0].Guarantees[0].Beneficiaries != nil {
		if originalPolicy.Assets[0].Guarantees[0].Beneficiaries != nil {
			gar.Beneficiaries = diffForBeneficiaries(*(originalPolicy.Assets[0].Guarantees[0]).Beneficiaries, *(inputPolicy.Assets[0].Guarantees[0]).Beneficiaries)
		} else {
			gar.Beneficiaries = inputPolicy.Assets[0].Guarantees[0].Beneficiaries
		}
	} else {
		gar.Beneficiaries = nil
	}
	garS := make([]models.Guarante, 0)
	garS = append(garS, gar)
	ass.Guarantees = garS
	assets = append(assets, ass)
	diff.Assets = assets
	return diff, nil
}

func diffForContractor(orig, input models.Contractor) models.Contractor {
	c := false
	if orig.FiscalCode != input.FiscalCode {
		return input
	}

	if input.Residence != nil && orig.Residence != nil {
		if orig.Residence.StreetName != input.Residence.StreetName {
			c = true
		}
		if orig.Residence.StreetNumber != input.Residence.StreetNumber {
			c = true
		}
		if orig.Residence.City != input.Residence.City {
			c = true
		}
		if orig.Residence.CityCode != input.Residence.CityCode {
			c = true
		}
	}
	if input.Domicile != nil && orig.Domicile != nil {
		if orig.Domicile.StreetName != input.Domicile.StreetName {
			c = true
		}
		if orig.Domicile.StreetNumber != input.Domicile.StreetNumber {
			c = true
		}
		if orig.Domicile.City != input.Domicile.City {
			c = true
		}
		if orig.Domicile.CityCode != input.Domicile.CityCode {
			c = true
		}
	}
	if orig.Mail != input.Mail {
		c = true
	}
	if orig.Phone != input.Phone {
		c = true
	}

	if c {
		return input
	}
	return models.Contractor{}
}

func diffForBeneficiaries(orig, input []models.Beneficiary) *[]models.Beneficiary {
	if len(input) != len(orig) {
		return &input
	}
	if reflect.DeepEqual(orig, input) {
		diffS := make([]models.Beneficiary, 0)
		return &diffS
	}
	return &input
}

func diffForUser(orig, input *models.User) *models.User {
	c := false
	if orig.FiscalCode != input.FiscalCode {
		return input
	}

	if input.Residence != nil && orig.Residence != nil {
		if orig.Residence.StreetName != input.Residence.StreetName {
			c = true
		}
		if orig.Residence.StreetNumber != input.Residence.StreetNumber {
			c = true
		}
		if orig.Residence.City != input.Residence.City {
			c = true
		}
		if orig.Residence.CityCode != input.Residence.CityCode {
			c = true
		}
	}
	if input.Domicile != nil && orig.Domicile != nil {
		if orig.Domicile.StreetName != input.Domicile.StreetName {
			c = true
		}
		if orig.Domicile.StreetNumber != input.Domicile.StreetNumber {
			c = true
		}
		if orig.Domicile.City != input.Domicile.City {
			c = true
		}
		if orig.Domicile.CityCode != input.Domicile.CityCode {
			c = true
		}
	}
	if orig.Mail != input.Mail {
		c = true
	}
	if orig.Phone != input.Phone {
		c = true
	}

	if c {
		return input
	}
	return &models.User{}
}

func modifyController(originalPolicy, inputPolicy models.Policy) (models.Policy, models.User, error) {
	var (
		err            error
		modifiedPolicy models.Policy
		modifiedUser   models.User
	)

	modifiedPolicy = originalPolicy

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
		if g.BeneficiaryReference != nil {
			(*modifiedGuarantees)[i].BeneficiaryReference = g.BeneficiaryReference
		}
		if g.Beneficiaries != nil {
			(*modifiedGuarantees)[i].Beneficiaries = g.Beneficiaries
		}
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
		log.Printf("error retrieving user %s from Firestore: %s", inputUser.Uid, err.Error())
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
			log.Printf("error modifying authentication email: %s", err.Error())
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
		log.Printf("error writing modified policy to Firestore: %s", err.Error())
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
		log.Printf("error writing modified user to Firestore: %s", err.Error())
		return err
	}

	err = user.BigquerySave("")
	if err != nil {
		log.Printf("error writing modified user to BigQuery: %s", err.Error())
		return err
	}

	log.Printf("user %s written into DBs", user.Uid)

	return err
}
