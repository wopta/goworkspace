package callback

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/wopta/goworkspace/document"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/mail"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/policy"
	"github.com/wopta/goworkspace/transaction"
	"github.com/wopta/goworkspace/user"
)

func Payment(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println("Payment")
	var response string
	var e error
	var fabrickCallback FabrickCallback
	uid := r.URL.Query().Get("uid")
	schedule := r.URL.Query().Get("schedule")
	request := lib.ErrorByte(io.ReadAll(r.Body))
	origin := r.URL.Query().Get("origin")

	log.Println(string(request))
	log.Println(string(r.RequestURI))
	json.Unmarshal([]byte(request), &fabrickCallback)

	// Unmarshal or Decode the JSON to the interface.
	if fabrickCallback.Bill.Status == "PAID" {
		if uid == "" || origin == "" {
			ext := strings.Split(fabrickCallback.ExternalID, "_")
			uid = ext[0]
			schedule = ext[1]
			origin = ext[2]
		}
		log.Println("Payment::uid: " + uid)
		log.Println("Payment::schedule: " + schedule)

		p := policy.GetPolicyByUid(uid, origin)
		if !p.IsPay && p.Status == models.PolicyStatusToPay {
			// Create/Update document on user collection based on contractor fiscalCode
			user.SetUserIntoPolicyContractor(&p, origin)

			// Get Policy contract
			gsLink := <-document.GetFileV6(p, uid)
			log.Println("Payment::contractGsLink: ", gsLink)

			// Update Policy as paid
			policy.SetPolicyPaid(&p, gsLink, origin)

			// Update the first transaction in policy as paid
			transaction.SetPolicyFirstTransactionPaid(uid, schedule, origin)

			// Update agency if present
			if p.AgencyUid != "" {
				var agency models.Agency
				fireAgency := lib.GetDatasetByEnv(origin, models.AgencyCollection)
				docsnap, err := lib.GetFirestoreErr(fireAgency, p.AgentUid)
				lib.CheckError(err)
				docsnap.DataTo(&agency)
				agency.Policies = append(agency.Policies, p.Uid)
				found := false
				for _, contractorUid := range agency.Users {
					if contractorUid == p.Contractor.Uid {
						found = true
						break
					}
				}
				if !found {
					agency.Users = append(agency.Users, p.Contractor.Uid)
				}
				err = lib.SetFirestoreErr(fireAgency, agency.Uid, agency)
				lib.CheckError(err)
			}

			// Update agent if present
			if p.AgentUid != "" {
				var agent models.Agent
				fireAgent := lib.GetDatasetByEnv(origin, models.AgentCollection)
				docsnap, err := lib.GetFirestoreErr(fireAgent, p.AgentUid)
				lib.CheckError(err)
				docsnap.DataTo(&agent)
				agent.Policies = append(agent.Policies, p.Uid)
				found := false
				for _, contractorUid := range agent.Users {
					if contractorUid == p.Contractor.Uid {
						found = true
						break
					}
				}
				if !found {
					agent.Users = append(agent.Users, p.Contractor.Uid)
				}
				err = lib.SetFirestoreErr(fireAgent, agent.Uid, agent)
				lib.CheckError(err)
			}

			// Send mail with the contract to the user
			log.Println("Payment: " + uid + " sendMail ")
			var contractbyte []byte
			name := p.Uid + ".pdf"
			contractbyte, e = lib.GetFromGoogleStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "assets/users/"+
				p.Contractor.Uid+"/contract_"+name)

			mail.SendMailContract(p, &[]mail.Attachment{{
				Byte:        base64.StdEncoding.EncodeToString(contractbyte),
				ContentType: "application/pdf",
				Name: strings.ReplaceAll(p.Contractor.Name+"_"+p.Contractor.Surname+"_"+p.NameDesc, " ",
					"_") + "_contratto.pdf",
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
