package renew

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/transaction"
)

func PromoteFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err      error
		response struct{} // TODO
	)

	log.SetPrefix("[PromoteFx] ")
	defer log.SetPrefix("")

	log.Println("Handler start -----------------------------------------------")

	// select all policies from renew draft collection where startDate + 1 year == now
	now := time.Now().UTC()
	query, params := buildQuery(now)

	policies, err := lib.QueryParametrizedRowsBigQuery[models.Policy](query, params)
	if err != nil {
		return "", nil, err
	}

	for _, p := range policies {
		if p.IsPay {
			promoteRenewPolicyDraftIntoPolicy(p)
			promoteRenewDraftTransactionIntoTransaction(p)
		} else {
			setPolicyNotPaid(p)
		}
	}

	responseJson, err := json.Marshal(response)

	log.Println("Handler end -------------------------------------------------")

	return string(responseJson), response, err
}

func buildQuery(date time.Time) (string, map[string]interface{}) {
	var (
		query  bytes.Buffer
		params = make(map[string]interface{})
	)

	// SELECT * FROM `wopta.renewPolicyDraft` WHERE EXTRACT(MONTH FROM startDate) = @date.Month() AND EXTRACT(DAY FROM startDate) = @date.Day()
	params["month"] = date.Month()
	params["day"] = date.Day()

	query.WriteString(fmt.Sprintf("SELECT * FROM `wopta.renewPolicyDraft` WHERE "+
		"EXTRACT(MONTH FROM startDate) = @%d AND "+
		"EXTRACT(DAY FROM startDate) = @%d",
		params["month"],
		params["day"]))

	return query.String(), params
}

func promoteRenewPolicyDraftIntoPolicy(policy models.Policy) error {
	policy.Updated = time.Now().UTC()
	// // // save draft in policy
	err := lib.SetFirestoreErr(models.PolicyCollection, policy.Uid, policy)
	policy.BigquerySave("")

	// // // delete draft
	_, err = lib.DeleteFirestoreErr("renewPolicyDraft", policy.Uid)
	// TODO: DELETE from bigquery

	return err
}

func promoteRenewDraftTransactionIntoTransaction(policy models.Policy) error {
	// // // save draft-transactions in
	var err error
	trs := transaction.GetPolicyTransactions("", policy.Uid) // TODO: get only next annuity - create function

	for _, tr := range trs {
		err = lib.SetFirestoreErr(models.TransactionsCollection, tr.Uid, tr)
		tr.BigQuerySave("")

		_, err = lib.DeleteFirestoreErr("renewTransactionDraft", tr.Uid)
		// TODO: DELETE from bigquery
	}

	return err
}

func setPolicyNotPaid(policy models.Policy) error {
	policy.Status = "INSOLUTO"
	policy.StatusHistory = append(policy.StatusHistory, policy.Status)
	policy.Updated = time.Now().UTC()

	return nil
}

func POCRoutine2Fx() {
	var (
		policies      = make([]models.Policy, 0)
		paidWg        sync.WaitGroup
		paidChannel   = make(chan models.Policy)
		unpaidWg      sync.WaitGroup
		unpaidChannel = make(chan models.Policy)
	)

	docIterator := lib.OrderFirestore(lib.PolicyCollection, "creationDate", firestore.Asc)
	snapshots, err := docIterator.GetAll()
	if err != nil {
		log.Println(err.Error())
		return
	}
	for _, sp := range snapshots {
		var policy models.Policy
		sp.DataTo(&policy)
		policies = append(policies, policy)

	}
	log.Printf("Got %d policies", len(policies))

	for _, p := range policies {
		if p.IsPay {
			paidWg.Add(1)
			go func(uid string) {
				paidChannel <- payPolicy(p)
			}(p.Uid)
			paidWg.Done()
		} else {
			unpaidWg.Add(1)
			go func(uid string) {
				unpaidChannel <- demotedPolicy(p)
			}(p.Uid)
			unpaidWg.Done()
		}
	}

	go func() {
		paidWg.Wait()
		close(paidChannel)
	}()

	go func() {
		unpaidWg.Wait()
		close(unpaidChannel)
	}()

	for p := range paidChannel {
		log.Printf("=== paid policy %s", p.Uid)
	}

	for p := range unpaidChannel {
		log.Printf("--- demoted policy %s", p.Uid)
	}
}

func demotedPolicy(p models.Policy) models.Policy {
	time.Sleep(time.Millisecond * 3000)
	return p
}

func payPolicy(p models.Policy) models.Policy {
	time.Sleep(time.Millisecond * 2000)
	payTransactions(p.Uid)
	return p
}

func payTransactions(policyUid string) {
	var trWg sync.WaitGroup
	promotedTransactions := make(chan models.Transaction)
	trs := transaction.GetPolicyTransactions("", policyUid) // TODO: get only next annuity - create function

	for _, tr := range trs {
		trWg.Add(1)
		go func(t models.Transaction) {
			time.Sleep(time.Millisecond * 1500)
			promotedTransactions <- t
		}(tr)
		trWg.Done()
	}

	go func() {
		trWg.Wait()
		close(promotedTransactions)
	}()

	for tr := range promotedTransactions {
		log.Printf("promoted transaction %s", tr.Uid)
	}
}

func POCRoutineFx() {

	pChannel := policyChannel()
	defer close(pChannel)

	paidChannel := make(chan string)
	unpaidChannel := make(chan string)

	paidCount := 0
	unpaidCount := 0

	for p := range pChannel {
		log.Println("-------------------------- " + p.Uid)
		if p.IsPay {
			go func(uid string) {
				paidChannel <- branch1(uid)
			}(p.Uid)
		} else {
			go func(uid string) {
				unpaidChannel <- branch2(uid)
			}(p.Uid)
		}
	}

	for p := range paidChannel {
		paidCount++
		log.Printf("%s - %d", p, paidCount)
	}

	for p := range unpaidChannel {
		unpaidCount++
		log.Printf("%s - %d", p, unpaidCount)
	}
}

func policyChannel() chan models.Policy {
	pChannel := make(chan models.Policy)

	log.Println("query start....")
	docIterator := lib.OrderFirestore(lib.PolicyCollection, "creationDate", firestore.Asc)
	log.Println("iterator start....")
	snapshots, err := docIterator.GetAll()
	if err != nil {
		log.Println(err.Error())
		return pChannel
	}
	log.Printf("looping %d snapshots....", len(snapshots))
	for idx, sp := range snapshots {
		var policy models.Policy
		sp.DataTo(&policy)
		go func(isLast bool) {
			pChannel <- policy
			if isLast {
				close(pChannel)
			}
		}(idx+1 == len(snapshots))
	}

	return pChannel
}

func branch1(uid string) string {
	msg := fmt.Sprintf("branch1! promoting policy %s at %d", uid, time.Now().UnixMilli())
	time.Sleep(time.Millisecond * 2000)
	nestedBranch(uid)
	return msg
}

func branch2(uid string) string {
	msg := fmt.Sprintf("branch2! %s INSOLUTO at %d", uid, time.Now().UnixMilli())

	time.Sleep(time.Millisecond * 2500)
	return msg
}

func nestedBranch(uid string) {
	nestedChannel := make(chan string)
	trs := transaction.GetPolicyTransactions("", uid)

	for _, tr := range trs {
		go func(uid, policyUid string) {
			nestedChannel <- branch3(uid, policyUid)
		}(tr.Uid, uid)
	}

	for n := range nestedChannel {
		log.Println(n)
	}
}

func branch3(uid, policyUid string) string {
	msg := fmt.Sprintf("branch3! promoting transaction %s - policy %s", uid, policyUid)
	time.Sleep(time.Millisecond * 1500)

	return msg
}

var (
	wg sync.WaitGroup
)

// ///////////
func POC() {
	var (
		policies       []models.Policy
		resultsChannel chan string = make(chan string)
		startTime                  = time.Now().UTC().UnixMilli()
	)

	log.Printf("Start at %d", startTime)

	policies = getAllPolicies()

	log.Printf("found %d policies", len(policies))

	wg.Add(len(policies))
	for _, p := range policies {
		if p.IsPay {
			// log.Printf("promoting %s ....", p.Uid)
			go promoPolicy(p, resultsChannel)
		} else {
			// log.Printf("demoting %s ....", p.Uid)
			go demotePolicy(p, resultsChannel)
		}
	}

	go func() {
		wg.Wait()
		log.Println("closing channel ...........")
		close(resultsChannel)
	}()

	for res := range resultsChannel {
		log.Println(res)
	}

	endTime := time.Now().UTC().UnixMilli()

	log.Printf("End at %d - duration %d", endTime, endTime-startTime)
}

func getAllPolicies() (policies []models.Policy) {
	docIterator := lib.OrderFirestore(lib.PolicyCollection, "creationDate", firestore.Asc)
	snapshots, err := docIterator.GetAll()
	if err != nil {
		log.Println(err.Error())
		return
	}
	for _, sp := range snapshots {
		var policy models.Policy
		sp.DataTo(&policy)
		policies = append(policies, policy)
	}
	return policies
}

func demotePolicy(p models.Policy, c chan<- string) {
	// log.Printf("demotion start %s", p.Uid)
	defer wg.Done()

	// do business logic
	time.Sleep(3000)

	c <- fmt.Sprintf("------ demoted policy %s", p.Uid)
}

func promoPolicy(p models.Policy, c chan<- string) {
	// log.Printf("promotion start %s", p.Uid)
	defer wg.Done()

	trs := transaction.GetPolicyTransactions("", p.Uid) // TODO: get only next annuity - create function

	wg.Add(len(trs))
	for _, tr := range trs {
		go promoTr(tr, c)
	}

	// do business logic
	time.Sleep(2000)

	c <- fmt.Sprintf("++++++ promoted policy %s", p.Uid)
}

func promoTr(tr models.Transaction, c chan<- string) {
	defer wg.Done()

	// do business logic
	time.Sleep(1500)

	c <- fmt.Sprintf("******* promoted transaction %s", tr.Uid)
}
