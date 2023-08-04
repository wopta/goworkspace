package transaction

import (
	"log"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/product"
)

func PutByPolicy(policy models.Policy, scheduleDate string, origin string, expireDate string, customerId string, amount float64, amountNet float64, providerId string, isPay bool) {
	log.Printf("[PutByPolicy] Policy %s", policy.Uid)
	var (
		commissionMga    float64
		commissionAgent  float64
		commissionAgency float64
		netCommission    map[string]float64
		sd               string
		status           string
		statusHistory    = make([]string, 0)
	)

	prod, err := product.GetProduct(policy.Name, policy.ProductVersion, models.UserRoleAdmin)
	if err != nil {
		log.Printf("[PutByPolicy] ERROR getting mga product: %s", err.Error())
		return
	}

	commissionMga = product.GetCommissionProduct(policy, *prod)
	log.Printf("[PutByPolicy] commissionMga: %g", commissionMga)

	if policy.AgentUid != "" {
		commissionAgent = getAgentCommission(policy)
		log.Printf("[PutByPolicy] commissionAgent: %g", commissionAgent)
	}

	if policy.AgencyUid != "" {
		commissionAgency = getAgencyCommission(policy)
		log.Printf("[PutByPolicy] commissionAgency: %g", commissionAgency)
	}

	if scheduleDate == "" {
		sd = time.Now().UTC().Format(models.TimeDateOnly)
	} else {
		sd = scheduleDate
	}

	if isPay {
		status = models.TransactionStatusPay
		statusHistory = append(statusHistory, models.TransactionStatusToPay, models.TransactionStatusPay)
	} else {
		status = models.TransactionStatusToPay
		statusHistory = append(statusHistory, models.TransactionStatusToPay)
	}

	fireTransactions := lib.GetDatasetByEnv(origin, models.TransactionsCollection)
	transactionUid := lib.NewDoc(fireTransactions)

	tr := models.Transaction{
		Amount:             amount,
		AmountNet:          amountNet,
		Id:                 "",
		Uid:                transactionUid,
		PolicyName:         policy.Name,
		PolicyUid:          policy.Uid,
		CreationDate:       time.Now().UTC(),
		Status:             status,
		StatusHistory:      statusHistory,
		ScheduleDate:       sd,
		ExpirationDate:     expireDate,
		NumberCompany:      policy.CodeCompany,
		Commissions:        commissionMga,
		IsPay:              isPay,
		Name:               policy.Contractor.Name + " " + policy.Contractor.Surname,
		Company:            policy.Company,
		IsDelete:           false,
		ProviderId:         providerId,
		UserToken:          customerId,
		ProviderName:       policy.Payment,
		AgentUid:           policy.AgencyUid,
		AgencyUid:          policy.AgencyUid,
		CommissionsAgent:   commissionAgent,
		CommissionsAgency:  commissionAgency,
		NetworkCommissions: netCommission,
	}

	err = lib.SetFirestoreErr(fireTransactions, transactionUid, tr)
	lib.CheckError(err)

	tr.BigQuerySave(origin)
}

func getAgentCommission(policy models.Policy) float64 {
	var agent models.Agent
	dn := lib.GetFirestore(models.AgentCollection, policy.AgentUid)
	dn.DataTo(&agent)

	return product.GetCommissionProducts(policy, agent.Products)
}

func getAgencyCommission(policy models.Policy) float64 {
	var agency models.Agency
	dn := lib.GetFirestore(models.AgencyCollection, policy.AgencyUid)
	dn.DataTo(&agency)
	return product.GetCommissionProducts(policy, agency.Products)
}
