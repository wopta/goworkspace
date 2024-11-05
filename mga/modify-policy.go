package mga

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

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
	modifiedPolicy, err = modifyController(originalPolicy, inputPolicy)
	if err != nil {
		log.Printf("error during policy modification: %s", err.Error())
		return "{}", nil, err
	}
	log.Printf("policy %s modified successfully", modifiedPolicy.Uid)

	rawPolicy, err = json.Marshal(modifiedPolicy)

	log.Println("Handler end -------------------------------------------------")

	return string(rawPolicy), modifiedPolicy, err
}

func modifyController(originalPolicy, inputPolicy models.Policy) (models.Policy, error) {
	var (
		err            error
		modifiedPolicy models.Policy
		modifiedUser   models.User
	)

	err = checkEmailUniqueness(originalPolicy.Contractor, inputPolicy.Contractor)
	if err != nil {
		return models.Policy{}, err
	}

	switch originalPolicy.Name {
	case models.LifeProduct, models.PersonaProduct:
		modifiedPolicy, modifiedUser, err = lifeModifier(originalPolicy, inputPolicy)
	case models.GapProduct:
		modifiedPolicy, modifiedUser, err = gapModifier(originalPolicy, inputPolicy)
	default:
		return models.Policy{}, errors.New("product not supported")
	}

	err = writePolicyToDb(modifiedPolicy)
	if err != nil {
		return models.Policy{}, err
	}

	err = writeUserToDB(modifiedUser)
	if err != nil {
		return models.Policy{}, err
	}

	return modifiedPolicy, err
}

func lifeModifier(originalPolicy, inputPolicy models.Policy) (models.Policy, models.User, error) {
	var (
		err                           error
		modifiedPolicy                models.Policy
		modifiedContractor            models.Contractor
		modifiedInsured, modifiedUser models.User
	)

	modifiedPolicy = originalPolicy

	modifiedContractor, err = modifyContractorInfo(inputPolicy.Contractor, originalPolicy.Contractor)
	if err != nil {
		log.Printf("error modifying contractor %s info: %s", originalPolicy.Contractor.Uid, err.Error())
		return models.Policy{}, models.User{}, err
	}

	modifiedInsured, err = modifyInsuredInfo(*inputPolicy.Assets[0].Person, *originalPolicy.Assets[0].Person)
	if err != nil {
		log.Printf("error modifying insured for policy %s info: %s", originalPolicy.Uid, err.Error())
		return models.Policy{}, models.User{}, err
	}

	modifiedPolicy.Contractor = modifiedContractor
	modifiedPolicy.Assets[0].Person = &modifiedInsured

	if modifiedContractor.Uid != "" {
		tmpUser := modifiedContractor.ToUser()
		modifiedUser, err = modifyUserInfo(*tmpUser)
		if err != nil {
			log.Printf("error modifying user %s info: %s", tmpUser.Uid, err.Error())
			return models.Policy{}, models.User{}, err
		}
	}

	return modifiedPolicy, modifiedUser, err
}

func gapModifier(originalPolicy, inputPolicy models.Policy) (models.Policy, models.User, error) {
	var (
		err                error
		modifiedPolicy     models.Policy
		modifiedContractor models.Contractor
		modifiedUser       models.User
	)

	modifiedPolicy = originalPolicy

	modifiedContractor, err = modifyContractorInfo(inputPolicy.Contractor, originalPolicy.Contractor)
	if err != nil {
		log.Printf("error modifying contractor %s info: %s", originalPolicy.Contractor.Uid, err.Error())
		return models.Policy{}, models.User{}, err
	}

	modifiedPolicy.Contractor = modifiedContractor

	if modifiedContractor.Uid != "" {
		tmpUser := modifiedContractor.ToUser()
		modifiedUser, err = modifyUserInfo(*tmpUser)
		if err != nil {
			log.Printf("error modifying user %s info: %s", tmpUser.Uid, err.Error())
			return models.Policy{}, models.User{}, err
		}
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

func modifyInsuredInfo(inputInsured, originalInsured models.User) (models.User, error) {
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
		return models.User{}, err
	}

	log.Printf("insured %s modified", originalInsured.Uid)

	log.Println("insured modified successfully")
	return *modifiedInsured, err
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
