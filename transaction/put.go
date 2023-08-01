package transaction

import (
	"log"
	"strings"
	"time"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/civil"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/product"
	pr "github.com/wopta/goworkspace/product"
)

func PutByPolicy(data models.Policy, scheduleDate string, origin string, expireDate string, customerId string, amount float64, amountNet float64, providerId string, isPay bool, role string) {
	var (
		commission       float64
		commissionAgent  float64
		commissionAgency float64
		netCommission    map[string]float64
	)

	//var prod models.Product
	prod, err := product.GetProduct(data.Name, data.ProductVersion, role)
	log.Println(data.Uid+" pay error marsh product:", err)
	commission = pr.GetCommissionProduct(data, *prod)

	if data.AgentUid != "" {
		var agent models.Agent
		dn := lib.GetFirestore(models.AgentCollection, data.AgentUid)
		dn.DataTo(&agent)
		commissionAgent = pr.GetCommissionProducts(data, agent.Products)

	}
	if data.AgencyUid != "" {
		var agency models.Agency
		dn := lib.GetFirestore(models.AgencyCollection, data.AgencyUid)
		dn.DataTo(&agency)
		commissionAgent = pr.GetCommissionProducts(data, agency.Products)
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
		ProviderId:         providerId,
		UserToken:          customerId,
		ProviderName:       data.Payment,
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
}
