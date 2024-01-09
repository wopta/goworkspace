package payment

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/mail"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
	plc "github.com/wopta/goworkspace/policy"
	prd "github.com/wopta/goworkspace/product"
	"github.com/wopta/goworkspace/transaction"
)

type FabrickRefreshPayByLinkRequest struct {
	PolicyUid string `json:"policyUid"`
}

func FabrickRefreshPayByLinkFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.SetPrefix("[FabrickRefreshPayByLinkFx] ")
	log.Println("Handler start -----------------------------------------------")

	var (
		request FabrickRefreshPayByLinkRequest
		err     error
	)

	origin := r.Header.Get("Origin")
	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()
	log.Printf("request body: %s", string(body))
	err = json.Unmarshal(body, &request)
	if err != nil {
		log.Printf("error unmarshaling body: %s", err.Error())
		return "", nil, err
	}

	policy := plc.GetPolicyByUid(request.PolicyUid, origin)

	err = fabrickRefreshPayByLink(&policy, origin)
	if err != nil {
		log.Printf("error refreshing payment link: %s", err.Error())
		return "", nil, err
	}

	models.CreateAuditLog(r, string(body))

	err = sendPayByLinkEmail(policy)
	if err != nil {
		log.Printf("error sending payment email: %s", err.Error())
		return "", nil, err
	}

	log.Println("Handler end -------------------------------------------------")

	return `{"success":true}`, `{"success":true}`, nil
}

func fabrickRefreshPayByLink(policy *models.Policy, origin string) error {
	var (
		paymentMethods          []string
		toBeDeletedTransactions []models.Transaction = make([]models.Transaction, 0)
		fabrickResponse         *FabrickPaymentResponse
		err                     error
	)

	if policy.Payment != models.FabrickPaymentProvider {
		return fmt.Errorf("payment provider '%s' not supported", policy.Payment)
	}

	mgaProduct := prd.GetProductV2(policy.Name, policy.ProductVersion, models.MgaChannel, nil, nil)
	product := prd.GetProductV2(policy.Name, policy.ProductVersion, policy.Channel, nil, nil)

	// extract payment methods
	paymentMethods = getPaymentMethods(*policy, product)

	// check policy payment split
	switch policy.PaymentSplit {
	// yearly/singleInstallment
	case string(models.PaySplitYear), string(models.PaySplitYearly), string(models.PaySplitSingleInstallment):
		fabrickResponse, toBeDeletedTransactions, err =
			fabrickRefreshPayByLinkSingleRate(policy, origin, paymentMethods, mgaProduct)
	case string(models.PaySplitMonthly):
		fabrickResponse, toBeDeletedTransactions, err =
			fabrickRefreshPayByLinkMultiRate(policy, origin, paymentMethods, mgaProduct, 12)
	default:
		err = fmt.Errorf("unhandle payment split '%s'", policy.PaymentSplit)
	}

	if err != nil {
		log.Printf("error refreshing pay by link: %s", err.Error())
		return err
	}

	if fabrickResponse == nil || fabrickResponse.Payload.PaymentPageURL == nil {
		return fmt.Errorf("fabrickResponse does not contain payment url: %v", fabrickResponse)
	}

	// update policy payurl and update date
	policy.PayUrl = *fabrickResponse.Payload.PaymentPageURL
	policy.Updated = time.Now().UTC()

	// loop all to be deleted transactions
	for _, tr := range toBeDeletedTransactions {
		// delete
		err = transaction.DeleteTransaction(&tr, origin, "Cancellata per ricreazione link di pagamento")
		if err != nil {
			log.Printf("error deleting transaction '%s': %s", tr.Uid, err.Error())
			return err
		}
		// TODO: delete on Fabrick
	}

	// save policy firestore
	firePolicy := lib.GetDatasetByEnv(origin, models.PolicyCollection)

	log.Println("saving policy to firestore...")
	err = lib.SetFirestoreErr(firePolicy, policy.Uid, policy)
	if err != nil {
		log.Printf("error saving policy to firestore: %s", err.Error())
		return err
	}

	// save policy bigquery
	log.Println("saving policy to bigquery...")
	policy.BigquerySave(origin)

	return nil
}

func fabrickRefreshPayByLinkSingleRate(
	policy *models.Policy,
	origin string,
	paymentMethods []string,
	mgaProduct *models.Product,
) (fabrickResponse *FabrickPaymentResponse, toBeDeletedTransactions []models.Transaction, err error) {
	// paid
	if policy.IsPay {
		// exit with error - cannot recreated paid annual policy
		return nil, nil, fmt.Errorf(
			"cannot refresh pay by link for policy with isPay '%T' and paymentSplit '%s'",
			policy.IsPay,
			policy.PaymentSplit,
		)
	}
	// not paid
	// // add old transaction to to be deleted array
	toBeDeletedTransactions = transaction.GetPolicyTransactions(origin, policy.Uid)
	// // create new one with current date - fabrickYearPay
	payRes := FabrickYearPay(*policy, origin, paymentMethods, mgaProduct)

	return &payRes, toBeDeletedTransactions, nil
}

func fabrickRefreshPayByLinkMultiRate(
	policy *models.Policy,
	origin string,
	paymentMethods []string,
	mgaProduct *models.Product,
	totalRates int,
) (fabrickResponse *FabrickPaymentResponse, toBeDeletedTransactions []models.Transaction, err error) {
	var refreshScheduleDates []time.Time = make([]time.Time, 0)

	// get all transactions
	allTransactions := transaction.GetPolicyTransactions(origin, policy.Uid)
	// add all unpaid transactions to to be deleted array
	for _, tr := range allTransactions {
		// TODO: check control flow
		if (tr.IsPay && !tr.IsDelete) || (!tr.IsPay && tr.IsDelete) {
			continue
		}
		scheduleDate, err := time.Parse(models.TimeDateOnly, tr.ScheduleDate)
		if err != nil {
			log.Printf("error parsing schedule date '%s' of transaction '%s': %s", tr.ScheduleDate, tr.Uid, err.Error())
			return nil, nil, err
		}
		log.Printf("adding transcation with schedule date '%s' to be recreated", tr.ScheduleDate)
		refreshScheduleDates = append(refreshScheduleDates, scheduleDate)
		if !tr.IsDelete {
			log.Printf("adding transcation with schedule date '%s' to be deleted", tr.ScheduleDate)
			toBeDeletedTransactions = append(toBeDeletedTransactions, tr)
		}
	}
	// create (totalRates - n paid transactions)
	payRes := fabrickMultiRatePayment(*policy, origin, paymentMethods, mgaProduct, refreshScheduleDates)

	return payRes, toBeDeletedTransactions, nil
}

func fabrickMultiRatePayment(
	policy models.Policy,
	origin string,
	paymentMethods []string,
	mgaProduct *models.Product,
	rateScheduleDates []time.Time,
) *FabrickPaymentResponse {
	log.Printf("creating payments for %d transactions", len(rateScheduleDates))

	customerId := uuid.New().String()
	// first transaction schedule date == now
	firstres := <-FabrickPayObj(policy, true, time.Now().UTC().Format(models.TimeDateOnly), "", customerId, policy.PriceGrossMonthly, policy.PriceNettMonthly, origin, paymentMethods, mgaProduct)
	time.Sleep(100)

	// we skip the first element since it is always created with todays date
	for _, sd := range rateScheduleDates[1:] {
		// other transactions same schedule date as before to respect policy.StartDate/EndDate
		expireDate := sd.AddDate(10, 0, 0)

		res := <-FabrickPayObj(policy, false, sd.Format(models.TimeDateOnly), expireDate.Format(models.TimeDateOnly), customerId, policy.PriceGrossMonthly, policy.PriceNettMonthly, origin, paymentMethods, mgaProduct)
		log.Printf("ScheduleDate: '%s' - response: %v", sd, res)
		time.Sleep(100)
	}

	return &firstres
}

func sendPayByLinkEmail(policy models.Policy) error {
	var (
		flowName    string
		networkNode *models.NetworkNode
		warrant     *models.Warrant
		toAddress   mail.Address
	)

	flowName = models.ECommerceFlow
	if policy.Channel == models.NetworkChannel {
		networkNode = network.GetNetworkNodeByUid(policy.ProducerUid)
		if networkNode == nil {
			return fmt.Errorf("networkNode not found")
		}
		toAddress = mail.GetNetworkNodeEmail(networkNode)
		warrant = networkNode.GetWarrant()
		if warrant == nil {
			return fmt.Errorf("warrant not found")
		}
		flowName = warrant.GetFlowName(policy.Name)
	} else {
		toAddress = mail.GetContractorEmail(&policy)
	}

	log.Printf("flowName '%s'", flowName)
	log.Printf("send pay mail to '%s'...", toAddress.String())

	mail.SendMailPay(
		policy,
		mail.AddressAnna,
		toAddress,
		mail.Address{},
		flowName,
	)

	return nil
}
