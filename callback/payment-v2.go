package callback

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/wopta/goworkspace/document"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/mail"
	"github.com/wopta/goworkspace/models"
	plc "github.com/wopta/goworkspace/policy"
)

const fabrickBillPaid string = "PAID"

var operations = make(map[string]map[string]interface{})

func PaymentV2Fx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println("[PaymentV2Fx] Handler start -----------------------------------")

	var (
		response        string
		err             error
		fabrickCallback FabrickCallback
	)

	policyUid := r.URL.Query().Get("uid")
	schedule := r.URL.Query().Get("schedule")
	origin := r.URL.Query().Get("origin")
	log.Printf("[PaymentV2Fx] uid %s, schedule %s", policyUid, schedule)

	request := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	log.Printf("[PaymentV2Fx] request payload: %s", string(request))
	err = json.Unmarshal([]byte(request), &fabrickCallback)
	if err != nil {
		log.Printf("[PaymentV2Fx] ERROR unmarshaling request: %s", err.Error())
		return response, nil, err
	}

	if policyUid == "" || origin == "" {
		ext := strings.Split(fabrickCallback.ExternalID, "_")
		policyUid = ext[0]
		schedule = ext[1]
		origin = ext[2]
	}

	switch fabrickCallback.Bill.Status {
	case fabrickBillPaid:
		err = fabrickPayment(origin, policyUid, schedule)
	default:
	}

	if err != nil {
		log.Printf("[PaymentV2Fx] ERROR: %s", err.Error())
		return response, nil, err
	}

	response = `{
		"result": true,
		"requestPayload": ` + string(request) + `,
		"locale": "it"
	}`
	log.Printf("[PaymentV2Fx] response: %s", response)

	return response, nil, nil
}

func fabrickPayment(origin, policyUid, schedule string) error {
	log.Printf("[fabrickPayment] Policy %s", policyUid)

	policy := plc.GetPolicyByUid(policyUid, origin)

	if !policy.IsPay && policy.Status == models.PolicyStatusToPay {
		fireTransactions := lib.GetDatasetByEnv(origin, models.TransactionsCollection)
		firePolicy := lib.GetDatasetByEnv(origin, models.PolicyCollection)

		// update policy and create/update user
		err := updatePolicyPaidAndSaveUser(&policy, origin)
		if err != nil {
			log.Printf("[fabrickPayment] ERROR modifying policy/user")
			return errors.New("policy/user error")
		}

		// Update the first transaction in policy as paid
		transactionOperation := setPolicyFirstTransactionPaid(policyUid, schedule, origin)
		if transactionOperation != nil {
			operations[fireTransactions] = transactionOperation
		} else {
			log.Printf("[fabrickPayment] ERROR modifying transaction")
			return errors.New("transaction error")
		}

		// Update agency if present
		err = models.UpdateAgencyPortfolio(&policy, origin)
		// check handling of error

		// Update agent if present
		err = models.UpdateAgentPortfolio(&policy, origin)
		// check habdling of error

		// Do batch operations
		err = lib.SetBatchFirestoreErr(operations)
		if err != nil {
			log.Printf("[fabrickPayment] ERROR %s", err.Error())
			return err
		}
		for _, tr := range operations[fireTransactions] {
			tr.(*models.Transaction).BigQuerySave(origin)
		}
		for _, pl := range operations[firePolicy] {
			pl.(*models.Policy).BigquerySave(origin)
		}

		// Send mail with the contract to the user
		mail.SendMailContract(policy, nil)

		return nil
	}

	log.Printf("[fabrickPayment] ERROR Policy %s with status %s and isPay %t cannot be paid", policyUid, policy.Status, policy.IsPay)
	return errors.New("cannot pay policy")
}

func updatePolicyPaidAndSaveUser(policy *models.Policy, origin string) error {
	log.Printf("[updatePolicyPaidAndSaveUser] Policy %s", policy.Uid)

	// promove documents from temp bucket to user
	err := updateIdentityDocument(policy)
	if err != nil {
		log.Printf("[updatePolicyPaidAndSaveUser] ERROR %s", err.Error())
		return err
	}

	// Create/Update document on user collection based on contractor fiscalCode
	err = setUserIntoPolicyContractor(policy, origin)
	if err != nil {
		log.Printf("[updatePolicyPaidAndSaveUser] ERROR %s", err.Error())
		return err
	}

	// Get Policy contract
	addPolicyContract(policy)

	// Update Policy as paid
	setPolicyPaid(policy, origin)

	firePolicy := lib.GetDatasetByEnv(origin, models.PolicyCollection)
	fireUser := lib.GetDatasetByEnv(origin, models.UserCollection)
	policyOperation := make(map[string]interface{})
	userOperation := make(map[string]interface{})

	policyOperation[policy.Uid] = policy
	userOperation[policy.Contractor.Uid] = &policy.Contractor
	operations[firePolicy] = policyOperation
	operations[fireUser] = userOperation

	return nil
}

func updateIdentityDocument(policy *models.Policy) error {
	// Move user identity documents to user folder on Google Storage
	for _, identityDocument := range policy.Contractor.IdentityDocuments {
		frontMediaBytes, err := lib.GetFromGoogleStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"),
			"temp/"+policy.Uid+"/"+identityDocument.FrontMedia.FileName)
		if err != nil {
			log.Printf("[updateIdentityDocument] ERROR getting front file: %s", err.Error())
			return err
		}
		frontGsLink, err := lib.PutToGoogleStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "assets/users/"+
			policy.Contractor.Uid+"/"+identityDocument.FrontMedia.FileName, frontMediaBytes)
		if err != nil {
			log.Printf("[updateIdentityDocument] ERROR saving front file: %s", err.Error())
			return err
		}
		identityDocument.FrontMedia.Link = frontGsLink

		if identityDocument.BackMedia != nil {
			backMediaBytes, err := lib.GetFromGoogleStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"),
				"temp/"+policy.Uid+"/"+identityDocument.BackMedia.FileName)
			if err != nil {
				log.Printf("[updateIdentityDocument] ERROR getting back file: %s", err.Error())
				return err
			}
			backGsLink, err := lib.PutToGoogleStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "assets/users/"+
				policy.Contractor.Uid+"/"+identityDocument.FrontMedia.FileName, backMediaBytes)
			if err != nil {
				log.Printf("[updateIdentityDocument] ERROR saving back file: %s", err.Error())
				return err
			}
			identityDocument.BackMedia.Link = backGsLink
		}
	}
	policy.Updated = time.Now().UTC()
	log.Println("[updateIdentityDocument] file(s) saved")

	return nil
}

func setUserIntoPolicyContractor(policy *models.Policy, origin string) error {
	log.Printf("[setUserIntoPolicyContractor] Policy %s", policy.Uid)
	userUID, newUser, err := models.GetUserUIDByFiscalCode(origin, policy.Contractor.FiscalCode)
	if err != nil {
		log.Printf("[setUserIntoPolicyContractor] ERROR finding user: %s", err.Error())
		return err
	}

	policy.Contractor.Uid = userUID

	if newUser {
		policy.Contractor.CreationDate = time.Now().UTC()
	} else {
		log.Printf("[setUserIntoPolicyContractor] Found user uid: %s", userUID)
		contractor := &policy.Contractor
		usersFire := lib.GetDatasetByEnv(origin, models.UserCollection)
		docSnap := lib.WhereFirestore(usersFire, "fiscalCode", "==", contractor.FiscalCode)
		retrievedUser, err := models.FirestoreDocumentToUser(docSnap)
		if err != nil {
			log.Printf("[setUserIntoPolicyContractor] ERROR getting user: %s", err.Error())
			return err
		}

		if retrievedUser.Uid != "" {
			retrievedUser.IdentityDocuments = append(retrievedUser.IdentityDocuments, contractor.IdentityDocuments...)
			retrievedUser.Consens = models.UpdateUserConsens(retrievedUser.Consens, contractor.Consens)
			retrievedUser.Address = contractor.Address
			retrievedUser.PostalCode = contractor.PostalCode
			retrievedUser.City = contractor.City
			retrievedUser.Locality = contractor.Locality
			retrievedUser.CityCode = contractor.CityCode
			retrievedUser.StreetNumber = contractor.StreetNumber
			retrievedUser.Location = contractor.Location
			retrievedUser.Residence = contractor.Residence
			retrievedUser.Domicile = contractor.Domicile
			retrievedUser.UpdatedDate = time.Now().UTC()
			if contractor.Height != 0 {
				retrievedUser.Height = contractor.Height
			}
			if contractor.Weight != 0 {
				retrievedUser.Weight = contractor.Weight
			}
			policy.Contractor = retrievedUser
		}
	}

	policy.Updated = time.Now().UTC()

	return nil
}

func addPolicyContract(policy *models.Policy) {
	// Get Policy contract
	gsLink := <-document.GetFileV6(*policy, policy.Uid)
	// Add Contract
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	filename := buildFilename([]string{"Contratto", policy.NameDesc, timestamp, ".pdf"})
	*policy.Attachments = append(*policy.Attachments, models.Attachment{
		Name:     "Contratto",
		Link:     gsLink,
		FileName: filename,
	})
	policy.Updated = time.Now().UTC()
}

func setPolicyPaid(policy *models.Policy, origin string) {
	policy.IsPay = true
	policy.Status = models.PolicyStatusPay
	policy.StatusHistory = append(policy.StatusHistory, models.PolicyStatusPay)
	policy.Updated = time.Now().UTC()
}

func setPolicyFirstTransactionPaid(policyUid, scheduleDate, origin string) map[string]interface{} {
	log.Printf("[setPolicyFirstTransactionPaid] Policy %s", policyUid)

	q := lib.Firequeries{
		Queries: []lib.Firequery{
			{
				Field:      "policyUid",
				Operator:   "==",
				QueryValue: policyUid,
			},
			{
				Field:      "scheduleDate",
				Operator:   "==",
				QueryValue: scheduleDate,
			},
		},
	}
	fireTransactions := lib.GetDatasetByEnv(origin, models.TransactionsCollection)
	query, err := q.FirestoreWherefields(fireTransactions)
	if err != nil {
		log.Printf("[setPolicyFirstTransactionPaid] ERROR getting from firestore %s", err.Error())
		return nil
	}

	transactions := models.TransactionToListData(query)
	transaction := transactions[0]
	tr, err := json.Marshal(transaction)
	if err != nil {
		log.Printf("[setPolicyFirstTransactionPaid] ERROR marshaling %s", err.Error())
		return nil
	}
	log.Printf("[setPolicyFirstTransactionPaid] transaction %s", string(tr))

	transaction.IsPay = true
	transaction.Status = models.TransactionStatusPay
	transaction.StatusHistory = append(transaction.StatusHistory, models.TransactionStatusPay)
	transaction.PayDate = time.Now().UTC()
	transaction.TransactionDate = time.Now().UTC()

	transactionOperation := make(map[string]interface{})
	transactionOperation[transaction.Uid] = &transaction

	return transactionOperation
}

func buildFilename(parts []string) string {
	return strings.ReplaceAll(
		strings.Join(parts, "_"),
		" ",
		"_",
	)
}
