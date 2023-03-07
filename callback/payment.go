package callback

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	lib "github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func Payment(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println("Payment")
	var response string
	var e error
	var fabrickCallback FabrickCallback
	request := lib.ErrorByte(ioutil.ReadAll(r.Body))
	log.Println(request)
	json.Unmarshal([]byte(request), &fabrickCallback)
	// Unmarshal or Decode the JSON to the interface.
	if fabrickCallback.Bill.Status == "PAID" {

		uid := r.URL.Query().Get("uid")
		schedule := r.URL.Query().Get("schedule")

		policyF := lib.GetFirestore("policy", uid)
		var policy models.Policy
		policyF.DataTo(policy)
		policy.IsPay = true
		policy.Updated = time.Now()
		policy.Status = models.PolicyStatusPay
		policy.StatusHistory = append(policy.StatusHistory, models.PolicyStatusPay)
		//policy.StatusHistory = append(policy.StatusHistory, models.PolicyStatusToPay)
		lib.SetFirestore("policy", uid, policy)
		q := lib.Firequeries{
			Queries: []lib.Firequery{{
				Field:      "uid",
				Operator:   "==",
				QueryValue: uid,
			},
				{
					Field:      "schedule",
					Operator:   "==",
					QueryValue: schedule,
				},
			},
		}
		query := q.FirestoreWherefields("transactions")

		transactions := models.TransactionToListData(query)
		transaction := transactions[0]
		transaction.IsPay = true
		transaction.Status = models.TransactionStatusPay
		transaction.StatusHistory = append(transaction.StatusHistory, models.TransactionStatusPay)
		lib.SetFirestore("transactions", transaction.Uid, transaction)
		e = lib.InsertRowsBigQuery("wopta", "transactions-day", transaction)
		log.Println(schedule)
		//log.Println(token)
		log.Println(q)
		response = `{
			"result": true,
			"requestPayload": ` + string(request) + `,
			"locale": "it"
		}`

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
	ExternalID *string `json:"externalId,omitempty"`
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
