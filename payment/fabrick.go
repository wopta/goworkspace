package payment

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/civil"
	"github.com/google/uuid"
	lib "github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	model "github.com/wopta/goworkspace/models"
)

func FabrickPayObj(data model.Policy, firstSchedule bool, scheduleDate string, customerId string, amount float64) <-chan FabrickPaymentResponse {
	r := make(chan FabrickPaymentResponse)

	go func() {
		defer close(r)
		log.Println("FabrickPay")
		//var b bytes.Buffer
		//fileReader := bytes.NewReader([]byte())

		var urlstring = os.Getenv("FABRICK_BASEURL") + "api/fabrick/pace/v4.0/mods/back/v1.0/payments"
		client := &http.Client{
			Timeout: time.Second * 15,
		}

		marshal := getfabbricPayments(data, firstSchedule, scheduleDate, customerId, amount)
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
			body, err := ioutil.ReadAll(res.Body)
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
			var commission float64
			for _, x := range prod.Companies {
				log.Println(data.Uid+" pay product name:", x.Name)
				log.Println(data.Uid+" pay product name:", data.Company)
				if x.Name == data.Company {
					if data.IsRenew {
						commission = x.CommissionRenew
					} else {
						commission = x.Commission
					}
				}

			}
			log.Println(data.Uid+"pay commission: ", commission)
			layout2 := "2006-01-02"
			var sd string
			if scheduleDate == "" {
				sd = time.Now().Format(layout2)
			} else {
				sd = scheduleDate
			}
			//tr := models.SetTransactionPolicy(data, data.Uid+"_"+scheduleDate, amount, scheduleDate, data.PriceNett * commission)
			tr := models.Transaction{
				Amount:             amount,
				Id:                 "",
				PolicyName:         data.Name,
				PolicyUid:          data.Uid,
				CreationDate:       time.Now(),
				Status:             models.TransactionStatusToPay,
				StatusHistory:      []string{models.TransactionStatusToPay},
				ScheduleDate:       sd,
				NumberCompany:      data.CodeCompany,
				Commissions:        data.PriceNett * commission,
				IsPay:              false,
				Name:               data.Contractor.Name + " " + data.Contractor.Surname,
				Company:            data.Company,
				CommissionsCompany: commission,
			}
			transactionsFire := lib.GetDatasetByContractorName(data.Contractor.Name, "transactions")

			ref, _ := lib.PutFirestore(transactionsFire, tr)
			tr.Uid = ref.ID
			tr.BigPayDate = civil.DateTimeOf(time.Now())
			tr.BigCreationDate = civil.DateTimeOf(time.Now())
			tr.BigStatusHistory = strings.Join(tr.StatusHistory, ",")
			err = lib.InsertRowsBigQuery("wopta", transactionsFire, tr)
			lib.CheckError(err)
			r <- result

		}
	}()
	return r
}
func FabbrickMontlyPay(data model.Policy) FabrickPaymentResponse {

	installment := data.PriceGross / 12
	customerId := uuid.New().String()
	log.Println(data.Uid + " FabbrickMontlyPay")
	layout := "2006-01-02"
	firstres := <-FabrickPayObj(data, true, data.StartDate.Format(layout), customerId, installment)
	time.Sleep(100)
	for i := 1; i <= 11; i++ {
		date := data.StartDate.AddDate(0, i, 0)
		res := <-FabrickPayObj(data, false, date.Format(layout), customerId, installment)
		log.Println(data.Uid+" FabbrickMontlyPay res:", res)
		time.Sleep(500)
	}
	return firstres
}
func FabbrickYearPay(data model.Policy) FabrickPaymentResponse {

	customerId := uuid.New().String()
	log.Println(data.Uid + " FabbrickYearPay")
	res := <-FabrickPayObj(data, false, "", customerId, data.PriceGross)

	return res
}

func getfabbricPayments(data model.Policy, firstSchedule bool, scheduleDate string, customerId string, amount float64) string {
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
	next := now.AddDate(0, 0, 4)
	layout := "2006-01-02T15:04:05.000Z"
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
	log.Println(next.Format(layout))
	externalId := data.Uid + "_" + scheduleDate
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
	calbackurl := "https://europe-west1-" + os.Getenv("GOOGLE_PROJECT_ID") + ".cloudfunctions.net/callback/v1/payment?uid=" + data.Uid + `&schedule=` + scheduleDate
	calbackurl = strings.Replace(calbackurl, `\u0026`, `&`, 1)
	log.Println(calbackurl)
	pay.PaymentConfiguration = PaymentConfiguration{

		//ExpirationDate: next.Format(layout),
		PaymentPageRedirectUrls: PaymentPageRedirectUrls{
			OnSuccess: "https://www.wopta.it",
			OnFailure: "https://www.wopta.it",
			//OnInterruption: "https://www.wopta.it",
		},

		AllowedPaymentMethods: &[]AllowedPaymentMethod{{Role: "payer", PaymentMethods: paymentMethods}},
		CallbackURL:           calbackurl,
		//PayByLink:             []PayByLink{{Type: "EMAIL", Recipients: data.Contractor.Mail, Template: "pay-by-link"}},
	}
	pay.Bill = bill

	res, _ := pay.Marshal()
	result := strings.Replace(string(res), `\u0026`, `&`, -1)
	return result
}
