package quote

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
	"github.com/wopta/goworkspace/sellable"
)

func LifeModFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		policy  models.Policy
		warrant *models.Warrant
	)

	log.Println("[LifeFx] handler start ----------------------")

	req := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	log.Println("[LifeFx] body: ", string(req))

	err := json.Unmarshal(req, &policy)
	if err != nil {
		log.Printf("[LifeFx] error unmarshaling body: %s", err.Error())
		return "", nil, err
	}

	policy.Normalize()

	authToken, err := models.GetAuthTokenFromIdToken(r.Header.Get("Authorization"))
	if err != nil {
		log.Printf("[LifeFx] error getting authToken from idToken: %s", err.Error())
		return "", nil, err
	}

	flow := authToken.GetChannelByRoleV2()

	log.Println("[LifeFx] loading network node")
	networkNode := network.GetNetworkNodeByUid(authToken.UserID)
	if networkNode != nil {
		warrant = networkNode.GetWarrant()
		if warrant != nil {
			flow = warrant.GetFlowName(policy.Name)
		}
	}

	// Extract file bytes for quoting
	channel := authToken.GetChannelByRoleV2()
	taxesBytes := lib.GetFilesByEnv(fmt.Sprintf("products-v2/%s/%s/taxes.csv", policy.Name, policy.ProductVersion))

	quoter := Quoter{
		taxesBytes: taxesBytes,
		sellable: func() (*models.Product, error) {
			return sellable.Life(&policy, channel, networkNode, warrant)
		},
		channel: channel,
		flow:    flow,
		policy:  policy,
	}

	log.Println("[LifeFx] start quoting")

	result, err := LifeMod(quoter)
	jsonOut, err := json.Marshal(result)

	log.Printf("[LifeFx] response: %s", string(jsonOut))

	log.Println("[LifeFx] handler end ---------------------------------------")

	return string(jsonOut), result, err

}

type Quoter struct {
	taxesBytes []byte
	sellable   func() (*models.Product, error)
	channel    string
	flow       string
	policy     models.Policy
}

func LifeMod(quoter Quoter) (models.Policy, error) {
	var err error

	log.Println("[Life] function start --------------------------------------")

	contractorAge, err := quoter.policy.CalculateContractorAge()

	log.Printf("[Life] contractor age: %d", contractorAge)

	df := lib.CsvToDataframe(quoter.taxesBytes)
	var selectRow []string

	rulesProduct, err := quoter.sellable()
	if err != nil {
		log.Printf("[LifeFx] error in sellable: %s", err.Error())
		return models.Policy{}, err
	}

	log.Printf("[Life] loading available rates for flow %s", quoter.flow)

	availableRates := getAvailableRates(rulesProduct, quoter.flow)

	log.Printf("[Life] available rates: %s", availableRates)

	log.Printf("[Life] add default guarantees")

	addDefaultGuarantees(quoter.policy, *rulesProduct)

	switch quoter.policy.ProductVersion {
	case models.ProductV1:
		death, err := quoter.policy.ExtractGuarantee(deathGuarantee)
		lib.CheckError(err)

		if quoter.channel == models.ECommerceChannel {
			log.Println("[Life] e-commerce flow")
			log.Println("[Life] setting sumInsuredLimitOfIndeminity")
			calculateSumInsuredLimitOfIndemnity(quoter.policy.Assets, death.Value.SumInsuredLimitOfIndemnity)
			log.Println("[Life] setting guarantees duration")
			calculateGuaranteeDuration(quoter.policy.Assets, death.Value.Duration.Year)
		}
	case models.ProductV2:
		if quoter.channel == models.ECommerceChannel {
			death, err := quoter.policy.ExtractGuarantee(deathGuarantee)
			lib.CheckError(err)
			log.Println("[Life] e-commerce flow")
			log.Println("[Life] setting sumInsuredLimitOfIndeminity")
			calculateSumInsuredLimitOfIndemnity(quoter.policy.Assets, death.Value.SumInsuredLimitOfIndemnity)
			log.Println("[Life] setting guarantees duration")
			calculateGuaranteeDuration(quoter.policy.Assets, death.Value.Duration.Year)
		} else {
			log.Println("[Life] mga, network flow")
			log.Println("[Life] setting sumInsuredLimitOfIndeminity")
			calculateSumInsuredLimitOfIndemnityV2(&quoter.policy)
		}
	}

	log.Println("[Life] updating policy start and end date")

	updatePolicyStartEndDate(&quoter.policy)

	log.Println("[Life] set guarantees subtitle")

	getGuaranteeSubtitle(quoter.policy.Assets)

	for _, row := range df.Records() {
		if row[0] == strconv.Itoa(contractorAge) {
			selectRow = row
			break
		}
	}

	quoter.policy.OffersPrices = map[string]map[string]*models.Price{
		"default": {
			"yearly":  &models.Price{},
			"monthly": &models.Price{},
		},
	}

	log.Println("[Life] calculate guarantees and offers prices")

	for assetIndex, asset := range quoter.policy.Assets {
		for guaranteeIndex, _ := range asset.Guarantees {
			guarantee := quoter.policy.Assets[assetIndex].Guarantees[guaranteeIndex]
			base, baseTax := getMultipliersIndex(guarantee.Slug)

			offset := getOffset(guarantee.Value.Duration.Year)

			baseFloat, taxFloat := getMultipliers(selectRow, offset, base, baseTax)

			calculateGuaranteePrices(&guarantee, baseFloat, taxFloat, *rulesProduct)

			if guarantee.IsSelected && guarantee.IsSellable {
				calculateOfferPrices(quoter.policy, guarantee)
			}
		}

	}

	log.Println("[Life] check monthly limit")

	monthlyToBeRemoved := !rulesProduct.Companies[0].IsMonthlyPaymentAvailable ||
		quoter.policy.OffersPrices["default"]["monthly"].Gross < rulesProduct.Companies[0].MinimumMonthlyPrice
	if monthlyToBeRemoved {
		log.Println("[Life] monthly payment disabled")
		delete(quoter.policy.OffersPrices["default"], "monthly")
	}

	log.Println("[Life] filtering available rates")

	removeOfferRate(&quoter.policy, availableRates)

	log.Println("[Life] round offers prices")

	roundOfferPrices(quoter.policy.OffersPrices)

	log.Println("[Life] sort guarantees list")

	sort.Slice(quoter.policy.Assets[0].Guarantees, func(i, j int) bool {
		return quoter.policy.Assets[0].Guarantees[i].Order < quoter.policy.Assets[0].Guarantees[j].Order
	})

	log.Println("[Life] function end --------------------------------------")

	return quoter.policy, err
}
