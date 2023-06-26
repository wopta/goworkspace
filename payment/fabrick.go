package payment

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/civil"
	"github.com/google/uuid"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func getFabrickClient(urlstring string, req *http.Request) (*http.Response, error) {
	client := &http.Client{
		Timeout: time.Second * 15,
	}

	req.Header.Set("api-key", os.Getenv("FABRICK_TOKEN_BACK_API"))
	req.Header.Set("Auth-Schema", "S2S")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-auth-token", os.Getenv("FABRICK_TOKEN_BACK_API"))
	req.Header.Set("Accept", "application/json")
	res, err := client.Do(req)

	return res, err
}
func FabrickPayObj(data models.Policy, firstSchedule bool, scheduleDate string, expireDate string, customerId string, amount float64, amountNet float64, origin string) <-chan FabrickPaymentResponse {
	r := make(chan FabrickPaymentResponse)

	go func() {
		defer close(r)
		log.Println("FabrickPay")

		var (
			urlstring        = os.Getenv("FABRICK_BASEURL") + "api/fabrick/pace/v4.0/mods/back/v1.0/payments"
			commission       float64
			commissionAgent  float64
			commissionAgency float64
			netCommission    map[string]float64
		)
		client := &http.Client{
			Timeout: time.Second * 15,
		}

		marshal := getfabbricPayments(data, firstSchedule, scheduleDate, expireDate, customerId, amount, origin)
		log.Printf(data.Uid + " " + string(marshal))
		//log.Println(getFabrickPay(data))
		req, _ := http.NewRequest(http.MethodPost, urlstring, strings.NewReader(string(marshal)))
		req.Header.Set("api-key", os.Getenv("FABRICK_TOKEN_BACK_API"))
		req.Header.Set("Auth-Schema", "S2S")
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("x-auth-token", os.Getenv("FABRICK_TOKEN_BACK_API"))
		req.Header.Set("Accept", "application/json")
		log.Printf(data.Uid+" ", req.Header)
		res, err := client.Do(req)
		lib.CheckError(err)

		if res != nil {
			log.Println("header:", res.Header)
			body, err := io.ReadAll(res.Body)
			log.Println(data.Uid+"pay response body:", string(body))
			lib.CheckError(err)
			var result FabrickPaymentResponse
			json.Unmarshal([]byte(body), &result)
			res.Body.Close()
			//prod, err := product.GetName(data.Name, "v"+fmt.Sprint(data.ProductVersion))
			prodb, err := lib.GetFromGoogleStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "products/"+data.Name+"-"+data.ProductVersion+".json")
			//var prod models.Product
			prod, err := models.UnmarshalProduct(prodb)
			log.Println(data.Uid+" pay error marsh product:", err)

			commission = getCommissionProduct(data, prod)

			if data.AgentUid != "" {
				var agent models.Agent
				dn := lib.GetFirestore(models.AgentCollection, data.AgentUid)
				dn.DataTo(&agent)
				commissionAgent = getCommissionProducts(data, agent.Products)

			}
			if data.AgencyUid != "" {
				var agency models.Agency
				dn := lib.GetFirestore(models.AgencyCollection, data.AgentUid)
				dn.DataTo(&agency)
				commissionAgent = getCommissionProducts(data, agency.Products)
			}
			log.Println(data.Uid+"pay commission: ", commission)
			layout2 := "2006-01-02"
			var sd string
			if scheduleDate == "" {
				sd = time.Now().UTC().Format(layout2)
			} else {
				sd = scheduleDate
			}
			//tr := models.SetTransactionPolicy(data, data.Uid+"_"+scheduleDate, amount, scheduleDate, data.PriceNett * commission)
			transactionsFire := lib.GetDatasetByEnv(origin, "transactions")
			transactionUid := lib.NewDoc(transactionsFire)

			tr := models.Transaction{
				Amount:             amount,
				AmountNet:          amountNet,
				Id:                 "",
				Uid:                transactionUid,
				PolicyName:         data.Name,
				PolicyUid:          data.Uid,
				CreationDate:       time.Now().UTC(),
				Status:             models.TransactionStatusToPay,
				StatusHistory:      []string{models.TransactionStatusToPay},
				ScheduleDate:       sd,
				ExpirationDate:     expireDate,
				NumberCompany:      data.CodeCompany,
				Commissions:        amountNet * commission,
				IsPay:              false,
				Name:               data.Contractor.Name + " " + data.Contractor.Surname,
				Company:            data.Company,
				CommissionsCompany: commission,
				IsDelete:           false,
				ProviderId:         *result.Payload.PaymentID,
				UserToken:          customerId,
				ProviderName:       "fabrick",
				AgentUid:           data.AgencyUid,
				AgencyUid:          data.AgencyUid,
				CommissionsAgent:   amountNet * commissionAgent,
				CommissionsAgency:  amountNet * commissionAgency,
				NetworkCommissions: netCommission,
			}

			lib.SetFirestore(transactionsFire, transactionUid, tr)
			tr.BigPayDate = bigquery.NullDateTime{}
			tr.BigTransactionDate = bigquery.NullDateTime{}
			tr.BigCreationDate = civil.DateTimeOf(time.Now().UTC())
			tr.BigStatusHistory = strings.Join(tr.StatusHistory, ",")
			err = lib.InsertRowsBigQuery("wopta", transactionsFire, tr)
			lib.CheckError(err)
			r <- result

		}
	}()
	return r
}
func FabbrickMontlyPay(data models.Policy, origin string) FabrickPaymentResponse {
	customerId := uuid.New().String()
	log.Println(data.Uid + " FabbrickMontlyPay")
	layout := "2006-01-02"
	firstres := <-FabrickPayObj(data, true, "", "", customerId, data.PriceGrossMonthly, data.PriceNettMonthly, origin)
	time.Sleep(100)
	for i := 1; i <= 11; i++ {
		date := data.StartDate.AddDate(0, i, 0)
		expireDate := date.AddDate(0, 0, 4)
		res := <-FabrickPayObj(data, false, date.Format(layout), expireDate.Format(layout), customerId, data.PriceGrossMonthly, data.PriceNettMonthly, origin)
		log.Println(data.Uid+" FabbrickMontlyPay res:", res)
		time.Sleep(100)
	}
	return firstres
}
func FabbrickYearPay(data models.Policy, origin string) FabrickPaymentResponse {

	customerId := uuid.New().String()
	log.Println(data.Uid + " FabbrickYearPay")
	res := <-FabrickPayObj(data, false, "", "", customerId, data.PriceGross, data.PriceNett, origin)

	return res
}

func getfabbricPayments(data models.Policy, firstSchedule bool, scheduleDate string, expireDate string, customerId string, amount float64, origin string) string {
	var mandate string

	if firstSchedule {
		mandate = "true"
	} else {
		mandate = "false"
	}
	if customerId == "" {
		customerId = uuid.New().String()
	}
	now := time.Now()
	layout2 := "2006-01-02"

	paymentMethods := []string{
		"CREDITCARD",
		//"FBKR2P",
		//"SDD",
		//"SMARTPOS",
	}
	var scheduleTransaction ScheduleTransaction

	var bill Bill
	if scheduleDate != "" {
		scheduleTransaction = ScheduleTransaction{DueDate: scheduleDate, PaymentInstrumentResolutionStrategy: "BY_PAYER"}
		bill.ScheduleTransaction = &scheduleTransaction
	} else {
		scheduleDate = now.Format(layout2)
	}
	externalId := data.Uid + "_" + scheduleDate + "_" + data.CodeCompany + "_" + strings.ReplaceAll(origin, "https://", "")
	var pay FabrickPaymentsRequest
	pay.MerchantID = "wop134b31-5926-4b26-1411-726bc9f0b111"
	pay.ExternalID = externalId

	bill.ExternalID = externalId
	bill.Amount = amount
	bill.Currency = "EUR"
	bill.Description = "Pagamento polizza n° " + data.CodeCompany

	bill.MandateCreation = mandate

	bill.Items = []Item{{ExternalID: externalId, Amount: amount, Currency: "EUR"}}
	bill.Subjects = &[]Subject{{ExternalID: customerId, Role: "customer", Email: data.Contractor.Mail, Name: data.Contractor.Name + ` ` + data.Contractor.Surname}}
	calbackurl := "https://europe-west1-" + os.Getenv("GOOGLE_PROJECT_ID") + ".cloudfunctions.net/callback/v1/payment?uid=" + data.Uid + `&schedule=` + scheduleDate + `&token=` + os.Getenv("WOPTA_TOKEN_API") + `&origin=` + origin
	calbackurl = strings.Replace(calbackurl, `\u0026`, `&`, 1)
	log.Println(calbackurl)
	pay.PaymentConfiguration = PaymentConfiguration{

		PaymentPageRedirectUrls: PaymentPageRedirectUrls{
			OnSuccess: "https://www.wopta.it",
			OnFailure: "https://www.wopta.it",
			//OnInterruption: "https://www.wopta.it",
		},

		AllowedPaymentMethods: &[]AllowedPaymentMethod{{Role: "payer", PaymentMethods: paymentMethods}},
		CallbackURL:           calbackurl,
		//PayByLink:             []PayByLink{{Type: "EMAIL", Recipients: data.Contractor.Mail, Template: "pay-by-link"}},
	}
	/*if expireDate != "" {
		pay.PaymentConfiguration.ExpirationDate = expireDate
	}*/

	pay.Bill = bill

	res, _ := pay.Marshal()
	result := strings.Replace(string(res), `\u0026`, `&`, -1)
	return result
}
func getCommissionProduct(data models.Policy, prod models.Product) float64 {
	var commission float64
	for _, x := range prod.Companies {
		if x.Name == data.Company {
			if data.IsRenew {
				return x.CommissionRenew
			} else {
				return x.Commission
			}
		}

	}
	return commission
}
func getCommissionProducts(data models.Policy, products []models.Product) float64 {
	var commission float64
	for _, prod := range products {
		if prod.Name == data.Name {
			return getCommissionProduct(data, prod)
		}

	}
	return commission
}
