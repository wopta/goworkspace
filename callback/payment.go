package callback

import (
	"cloud.google.com/go/civil"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/wopta/goworkspace/document"

	"cloud.google.com/go/firestore"
	lib "github.com/wopta/goworkspace/lib"
	mail "github.com/wopta/goworkspace/mail"
	"github.com/wopta/goworkspace/models"
)

func Payment(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println("Payment")
	var response string
	var e error
	var query *firestore.DocumentIterator
	var fabrickCallback FabrickCallback
	uid := r.URL.Query().Get("uid")
	schedule := r.URL.Query().Get("schedule")
	request := lib.ErrorByte(io.ReadAll(r.Body))
	origin := r.URL.Query().Get("origin")

	log.Println(string(request))
	log.Println(string(r.RequestURI))
	json.Unmarshal([]byte(request), &fabrickCallback)

	now := time.Now().UTC()
	// Unmarshal or Decode the JSON to the interface.
	if fabrickCallback.Bill.Status == "PAID" {
		if uid == "" || origin == "" {
			ext := strings.Split(fabrickCallback.ExternalID, "_")
			uid = ext[0]
			schedule = ext[1]
			origin = ext[2]
		}
		firePolicy := lib.GetDatasetByEnv(origin, "policy")
		log.Println(uid)
		log.Println(schedule)
		policyF := lib.GetFirestore(firePolicy, uid)
		var policy models.Policy
		policyF.DataTo(&policy)
		policyM, _ := policy.Marshal()
		log.Println(uid+" payment ", string(policyM))
		if !policy.IsPay && policy.Status == models.PolicyStatusToPay {
			// Get User UID by fiscal code
			userUID, newUser, err := models.GetUserUIDByFiscalCode(r.Header.Get("origin"), policy.Contractor.FiscalCode)
			lib.CheckError(err)
			policy.Contractor.Uid = userUID
			log.Println("Contractor UID: ", userUID)
			log.Println("Policy Contractor UID: ", policy.Contractor.Uid)

			gsLink := <-document.GetFileV6(policy, uid)
			log.Println("contractGsLink: ", gsLink)
			timestamp := strconv.FormatInt(now.Unix(), 10)
			*policy.Attachments = append(*policy.Attachments, models.Attachment{
				Name: "Contratto",
				Link: gsLink,
				FileName: "Contratto_" + strings.ReplaceAll(policy.NameDesc, " ", "_") +
					"_" + timestamp + ".pdf",
			})

			// Move user identity documents to user folder on Google Storage
			for _, identityDocument := range policy.Contractor.IdentityDocuments {
				frontMediaBytes, e := lib.GetFromGoogleStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"),
					"temp/"+policy.Uid+"/"+identityDocument.FrontMedia.FileName)
				lib.CheckError(e)
				frontGsLink, e := lib.PutToGoogleStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "assets/users/"+
					userUID+"/"+identityDocument.FrontMedia.FileName, frontMediaBytes)
				log.Println("frontGsLink: ", frontGsLink)
				identityDocument.FrontMedia.Link = frontGsLink

				if identityDocument.BackMedia != nil {
					backMediaBytes, e := lib.GetFromGoogleStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"),
						"temp/"+policy.Uid+"/"+identityDocument.BackMedia.FileName)
					lib.CheckError(e)
					backGsLink, e := lib.PutToGoogleStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "assets/users/"+
						userUID+"/"+identityDocument.FrontMedia.FileName, backMediaBytes)
					log.Println("backGsLink: ", backGsLink)
					identityDocument.BackMedia.Link = backGsLink
				}
			}

			if newUser {
				policy.Contractor.CreationDate = now
				fireUsers := lib.GetDatasetByEnv(r.Header.Get("origin"), "users")
				lib.SetFirestore(fireUsers, userUID, policy.Contractor)
			} else {
				_, err = models.UpdateUserByFiscalCode(r.Header.Get("origin"), policy.Contractor)
				lib.CheckError(err)
			}

			policy.IsPay = true
			policy.Updated = now
			policy.Status = models.PolicyStatusPay
			policy.StatusHistory = append(policy.StatusHistory, models.PolicyStatusPay)
			//policy.StatusHistory = append(policy.StatusHistory, models.PolicyStatusToPay)
			lib.SetFirestore(firePolicy, uid, policy)
			policy.BigquerySave(r.Header.Get("origin"))
			q := lib.Firequeries{
				Queries: []lib.Firequery{
					{
						Field:      "policyUid",
						Operator:   "==",
						QueryValue: uid,
					},
					{
						Field:      "scheduleDate",
						Operator:   "==",
						QueryValue: schedule,
					},
				},
			}
			fireTransactions := lib.GetDatasetByEnv(origin, "transactions")
			query, e = q.FirestoreWherefields(fireTransactions)
			transactions := models.TransactionToListData(query)
			transaction := transactions[0]
			tr, _ := json.Marshal(transaction)
			log.Println(uid+" payment ", string(tr))
			transaction.IsPay = true
			transaction.Status = models.TransactionStatusPay
			transaction.StatusHistory = append(transaction.StatusHistory, models.TransactionStatusPay)
			transaction.PayDate = now
			transaction.BigPayDate = civil.DateTimeOf(transaction.PayDate)
			transaction.BigCreationDate = civil.DateTimeOf(transaction.CreationDate)
			lib.SetFirestore(fireTransactions, transaction.Uid, transaction)
			e = lib.InsertRowsBigQuery("wopta", fireTransactions, transaction)
			log.Println(uid + " payment sendMail ")
			var contractbyte []byte
			name := policy.Uid + ".pdf"
			contractbyte, e = lib.GetFromGoogleStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "assets/users/"+
				policy.Contractor.Uid+"/contract_"+name)

			mail.SendMailContract(policy, &[]mail.Attachment{{
				Byte:        base64.StdEncoding.EncodeToString(contractbyte),
				ContentType: "application/pdf",
				Name: policy.Contractor.Name + "_" + policy.Contractor.Surname + "_" +
					strings.ReplaceAll(policy.NameDesc, " ", "_") + "_contratto.pdf",
			}})

			response = `{
			"result": true,
			"requestPayload": ` + string(request) + `,
			"locale": "it"
		}`
			log.Println(response)
		}
	}
	return response, nil, e
}
func UnmarshalFabrickCallback(data []byte) (FabrickCallback, error) {
	var r FabrickCallback
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *FabrickCallback) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type FabrickCallback struct {
	ExternalID string  `json:"externalId,omitempty"`
	PaymentID  *string `json:"paymentId,omitempty"`
	Bill       *Bill   `json:"bill,omitempty"`
}

type Bill struct {
	ExternalID     *string       `json:"externalId,omitempty"`
	BillID         *string       `json:"billId,omitempty"`
	Amount         *float64      `json:"amount,omitempty"`
	Currency       *string       `json:"currency,omitempty"`
	Description    *string       `json:"description,omitempty"`
	ReservedAmount *int64        `json:"reservedAmount,omitempty"`
	ResidualAmount *int64        `json:"residualAmount,omitempty"`
	RefundedAmount *int64        `json:"refundedAmount,omitempty"`
	PaidAmout      *float64      `json:"paidAmout,omitempty"`
	Items          []Item        `json:"items,omitempty"`
	Status         string        `json:"status,omitempty"`
	Transactions   []Transaction `json:"transactions,omitempty"`
}

type Item struct {
	ExternalID  *string     `json:"externalId,omitempty"`
	ItemID      *string     `json:"itemId,omitempty"`
	Amount      *float64    `json:"amount,omitempty"`
	Currency    *string     `json:"currency,omitempty"`
	Description *string     `json:"description,omitempty"`
	XInfo       interface{} `json:"xInfo"`
	Status      *string     `json:"status,omitempty"`
	Xinfo       interface{} `json:"xinfo"`
}

type Transaction struct {
	TransactionID       *string     `json:"transactionId,omitempty"`
	TransactionDateTime interface{} `json:"transactionDateTime"`
	Amount              *float64    `json:"amount,omitempty"`
	Currency            *string     `json:"currency,omitempty"`
	GatewayID           interface{} `json:"gatewayId"`
	AcquirerID          interface{} `json:"acquirerId"`
	Status              *string     `json:"status,omitempty"`
	PaymentMethod       *string     `json:"paymentMethod,omitempty"`
}
