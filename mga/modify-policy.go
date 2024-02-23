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
		modifiedInsured    models.User
	)

	modifiedContractor, err = editContractorInfo(inputPolicy.Contractor, originalPolicy.Contractor)
	if err != nil {
		log.Printf("error modifying contractor %s info: %s", originalPolicy.Contractor.Uid, err.Error())
		return models.Policy{}, err
	}
	modifiedInsured, err = editInsuredInfo(*inputPolicy.Assets[0].Person, *originalPolicy.Assets[0].Person)
	if err != nil {
		log.Printf("error editing insured for policy %s info: %s", originalPolicy.Uid, err.Error())
		return models.Policy{}, err
	}

	modifiedPolicy.Contractor = modifiedContractor
	modifiedPolicy.Assets[0].Person = &modifiedInsured
	return modifiedPolicy, err
}

func editContractorInfo(inputContractor, originalContractor models.Contractor) (models.Contractor, error) {
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
	err = usr.CheckFiscalCode(*user)
	if err != nil {
		return models.Contractor{}, err
	}
	/*computedFiscalCode, _, err := usr.CalculateFiscalCode(*user)
	if err != nil {
		log.Printf("error computing fiscalCode for contractor %s: %s", modifiedContractor.Uid, err.Error())
		return models.Contractor{}, err
	}
	if !strings.EqualFold(modifiedContractor.FiscalCode, computedFiscalCode) {
		log.Printf("computed fiscalCode %s not matching inputted fiscalCode %s", computedFiscalCode, modifiedContractor.FiscalCode)
		return models.Contractor{}, errors.New("invalid fiscalCode")
	}*/

	log.Printf("contractor %s modified", originalContractor.Uid)

	return *modifiedContractor, err
}

func editInsuredInfo(inputInsured, originalInsured models.User) (models.User, error) {
	var (
		err             error
		modifiedInsured = new(models.User)
	)

	log.Println("editing insured info...")
	*modifiedInsured = originalInsured

	modifiedInsured.Name = inputInsured.Name
	modifiedInsured.Surname = inputInsured.Surname
	modifiedInsured.BirthDate = inputInsured.BirthDate
	modifiedInsured.BirthCity = inputInsured.BirthCity
	modifiedInsured.BirthProvince = inputInsured.BirthProvince
	modifiedInsured.Gender = inputInsured.Gender
	modifiedInsured.FiscalCode = inputInsured.FiscalCode

	// TODO: handle omocodia
	err = usr.CheckFiscalCode(*modifiedInsured)
	if err != nil {
		return models.User{}, err
	}
	/*computedFiscalCode, _, err := usr.CalculateFiscalCode(*modifiedInsured)
	if err != nil {
		log.Printf("error computing fiscalCode for contractor %s: %s", modifiedInsured.Uid, err.Error())
		return models.User{}, err
	}
	if !strings.EqualFold(modifiedInsured.FiscalCode, computedFiscalCode) {
		log.Printf("computed fiscalCode %s not matching inputted fiscalCode %s", computedFiscalCode, modifiedInsured.FiscalCode)
		return models.User{}, errors.New("invalid fiscalCode")
	}*/

	log.Printf("contractor %s modified", originalInsured.Uid)

	log.Println("insured edited successfully")
	return *modifiedInsured, err
}
