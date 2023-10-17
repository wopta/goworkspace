package transaction

import (
	"log"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/product"
)

func PutByPolicy(policy models.Policy, scheduleDate, origin, expireDate, customerId string, amount, amountNet float64, providerId, paymentMethod string, isPay bool) *models.Transaction {
	log.Printf("[PutByPolicy] Policy %s", policy.Uid)
	var (
		commissionMga    float64
		commissionAgent  float64
		commissionAgency float64
		netCommission    map[string]float64
		sd               string
		status           string
		statusHistory    = make([]string, 0)
		payDate          time.Time
		transactionDate  time.Time
	)

	prod := product.GetProductV2(policy.Name, policy.ProductVersion, models.MgaChannel, nil, nil)
	if prod == nil {
		log.Printf("[PutByPolicy] error getting mga product")
		return nil
	}

	// TODO fix me - workaround for Gap camper mga commissions
	if isGapCamper(&policy) {
		log.Println("[PutByPolicy] overrinding product commissions for Gap camper")
		prod.Offers["base"].Commissions.NewBusiness = 0.22
		prod.Offers["base"].Commissions.Renew = 0.22
		prod.Offers["complete"].Commissions.NewBusiness = 0.37
		prod.Offers["complete"].Commissions.Renew = 0.37
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

	now := time.Now().UTC()

	if isPay {
		status = models.TransactionStatusPay
		statusHistory = append(statusHistory, models.TransactionStatusToPay, models.TransactionStatusPay)
		payDate = now
		transactionDate = now
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
		CreationDate:       now,
		Status:             status,
		StatusHistory:      statusHistory,
		ScheduleDate:       sd,
		ExpirationDate:     expireDate,
		NumberCompany:      policy.CodeCompany,
		Commissions:        commissionMga,
		IsPay:              isPay,
		PayDate:            payDate,
		TransactionDate:    transactionDate,
		Name:               policy.Contractor.Name + " " + policy.Contractor.Surname,
		Company:            policy.Company,
		IsDelete:           false,
		ProviderId:         providerId,
		UserToken:          customerId,
		ProviderName:       policy.Payment,
		AgentUid:           policy.AgentUid,
		AgencyUid:          policy.AgencyUid,
		CommissionsAgent:   commissionAgent,
		CommissionsAgency:  commissionAgency,
		NetworkCommissions: netCommission,
		PaymentMethod:      paymentMethod,
	}

	err := lib.SetFirestoreErr(fireTransactions, transactionUid, tr)
	if err != nil {
		log.Printf("[PutByPolicy] error saving transaction to firestore: %s", err.Error())
		return nil
	}

	tr.BigQuerySave(origin)

	return &tr
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

func isGapCamper(policy *models.Policy) bool {
	return policy.Name == models.GapProduct &&
		len(policy.Assets) > 0 &&
		policy.Assets[0].Vehicle != nil &&
		policy.Assets[0].Vehicle.VehicleType == "camper"
}
