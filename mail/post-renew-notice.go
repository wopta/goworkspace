package mail

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
)

type RenewNoticeReq struct {
	Date string `json:"date"`
}

type NodePortfolio struct {
	Node     models.NetworkNode
	Policies []models.Policy
}

func RenewNoticeFx(w http.ResponseWriter, r *http.Request) (string, any, error) {
	var (
		err           error
		wg            = new(sync.WaitGroup)
		today         = time.Now().UTC()
		request       RenewNoticeReq
		policies      []models.Policy
		nodePolicyMap = make(map[string]NodePortfolio)
	)

	log.SetPrefix("[RenewNoticeFx] ")
	defer func() {
		if err != nil {
			log.Printf("error: %s", err.Error())
		}
		log.Println("Handler end -------------------------------------------------")
	}()
	log.Println("Handler start -----------------------------------------------")

	if err = json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Println("error decoding request")
		return "", nil, err
	}

	if request.Date != "" {
		tmpDate, err := time.Parse(time.DateOnly, request.Date)
		if err != nil {
			log.Println("error parsing request date")
			return "", nil, err
		}
		today = tmpDate
	}

	targetDate := today.AddDate(0, 1, 0)
	log.Printf("executing query for date: %v", targetDate)
	if policies, err = getRenewingPolicies(targetDate); err != nil {
		log.Println("error getting policies")
		return "", nil, err
	}

	for _, p := range policies {
		if _, ok := nodePolicyMap[p.ProducerUid]; !ok {
			var node *models.NetworkNode
			if node = network.GetNetworkNodeByUid(p.ProducerUid); node == nil {
				log.Println("error getting node")
				return "", nil, fmt.Errorf("node '%s' not found", p.ProducerUid)
			}
			nodePolicyMap[p.ProducerUid] = NodePortfolio{
				Node:     *node,
				Policies: make([]models.Policy, 0),
			}
		}
		current := nodePolicyMap[p.ProducerUid]
		current.Policies = append(current.Policies, p)
		nodePolicyMap[p.ProducerUid] = current
	}

	for _, value := range nodePolicyMap {
		wg.Add(1)
		go SendMailRenewNotice(value.Node, AddressAnna, GetNetworkNodeEmail(&value.Node), Address{})
	}

	go func() {
		wg.Wait()
	}()

	return "", nil, nil
}

func getRenewingPolicies(date time.Time) (policies []models.Policy, err error) {
	var (
		query  bytes.Buffer
		params = make(map[string]interface{})
	)

	params["month"] = int64(date.Month())
	params["year"] = int64(date.Year())
	params["isDeleted"] = false
	params["isRenewable"] = true
	params["channel"] = models.NetworkChannel

	query.WriteString(fmt.Sprintf("SELECT * FROM `%s.%s` WHERE "+
		"EXTRACT(MONTH FROM startDate) = @month AND "+
		"EXTRACT(YEAR FROM startDate) = @year AND "+
		"channel = @channel AND "+
		"isDeleted = @isDeleted AND "+
		"isRenewable = @isRenewable",
		models.WoptaDataset,
		lib.RenewPolicyViewCollection))

	log.Printf("query: %s", query.String())
	log.Printf("params: %v", params)

	if policies, err = lib.QueryParametrizedRowsBigQuery[models.Policy](query.String(), params); err != nil {
		return
	}

	for index, policy := range policies {
		var temp models.Policy
		if err = json.Unmarshal([]byte(policy.Data), &temp); err != nil {
			policies = nil
			return
		}
		policies[index] = temp
	}

	return
}
