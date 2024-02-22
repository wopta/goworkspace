package mga

import (
	"encoding/json"
	"errors"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	plc "github.com/wopta/goworkspace/policy"
	usr "github.com/wopta/goworkspace/user"
	"io"
	"log"
	"net/http"
	"strings"
)

func ModifyPolicyFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err                         error
		inputPolicy, modifiedPolicy models.Policy
	)

	log.SetPrefix("[ModifyPolicyFx] ")
	defer log.SetPrefix("")
	log.Println("Handler start -----------------------------------------------")

	body := lib.ErrorByte(io.ReadAll(r.Body))
	log.Printf("request body: %s", string(body))
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
	modifiedPolicy, err = modifyController(originalPolicy, modifiedPolicy)
	if err != nil {
		log.Printf("error during policy modification: %s", err.Error())
		return "{}", nil, err
	}
	log.Printf("policy %s modified successfully", modifiedPolicy.Uid)

	log.Printf("writing modified policy to Firestore...")
	err = lib.SetFirestoreErr(models.PolicyCollection, modifiedPolicy.Uid, modifiedPolicy)
	if err != nil {
		log.Printf("error writing modified policy to Firestore: %s", err.Error())
		return "{}", nil, err
	}
	log.Printf("policy %s written to Firestore", modifiedPolicy.Uid)

	log.Printf("writing modified policy to BigQuery...")
	modifiedPolicy.BigquerySave("")
	log.Printf("policy %s written to BigQuery", modifiedPolicy.Uid)

	rawPolicy, err = json.Marshal(modifiedPolicy)
	log.Printf("modified policy: %s", string(rawPolicy))

	log.Println("Handler end -------------------------------------------------")

	return string(rawPolicy), modifiedPolicy, err
}

func modifyController(originalPolicy, inputPolicy models.Policy) (models.Policy, error) {
	var (
		err            error
		modifiedPolicy models.Policy
	)

	switch originalPolicy.Name {
	case models.LifeProduct:
		modifiedPolicy, err = lifeModifier(inputPolicy, originalPolicy)
	default:
		return models.Policy{}, errors.New("product not supported")
	}

	return modifiedPolicy, err
}

func lifeModifier(inputPolicy, originalPolicy models.Policy) (models.Policy, error) {
	var (
		err                error
		modifiedPolicy     models.Policy
		modifiedContractor models.Contractor
		modifiedInsured    *models.User
	)

	modifiedContractor, err = updateContractorInfo(inputPolicy.Contractor, originalPolicy.Contractor)
	if err != nil {
		log.Printf("error modifying contractor %s info: %s", originalPolicy.Contractor.Uid, err.Error())
		return models.Policy{}, err
	}
	if strings.EqualFold(inputPolicy.Contractor.FiscalCode, inputPolicy.Assets[0].Person.FiscalCode) {
		modifiedInsured = modifiedContractor.ToUser()
	}

	modifiedPolicy.Contractor = modifiedContractor
	modifiedPolicy.Assets[0].Person = modifiedInsured
	return modifiedPolicy, err
}

func updateContractorInfo(inputContractor, originalContractor models.Contractor) (models.Contractor, error) {
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

	// TODO: handle omocodia
	user := modifiedContractor.ToUser()
	computedFiscalCode, _, err := usr.CalculateFiscalCode(*user)
	if err != nil {
		log.Printf("error computing fiscalCode for contractor %s: %s", modifiedContractor.Uid, err.Error())
		return models.Contractor{}, err
	}
	if !strings.EqualFold(modifiedContractor.FiscalCode, computedFiscalCode) {
		log.Printf("computed fiscalCode %s not matching inputted fiscalCode %s", computedFiscalCode, modifiedContractor.FiscalCode)
		return models.Contractor{}, errors.New("invalid fiscalCode")
	}

	return *modifiedContractor, err
}
