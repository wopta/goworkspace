package payment

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
	prd "github.com/wopta/goworkspace/product"
)

func FabrickPayMonthlyFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	req := lib.ErrorByte(io.ReadAll(r.Body))

	var (
		data    models.Policy
		warrant *models.Warrant
	)

	defer r.Body.Close()
	err := json.Unmarshal([]byte(req), &data)
	log.Println(data.PriceGross)
	lib.CheckError(err)

	networkNode := network.GetNetworkNodeByUid(data.ProducerUid)
	if networkNode != nil {
		warrant = networkNode.GetWarrant()
	}
	product := prd.GetProductV2(data.Name, data.ProductVersion, data.Channel, networkNode, warrant)
	mgaProduct := prd.GetProductV2(data.Name, data.ProductVersion, models.MgaChannel, nil, nil)

	paymentMethods := getPaymentMethods(data, product)

	resultPay := FabrickMonthlyPay(data, getOrigin(r.Header.Get("origin")), paymentMethods, mgaProduct)
	b, err := json.Marshal(resultPay)
	log.Println(resultPay)
	return string(b), resultPay, err
}

func FabrickMonthlyPay(data models.Policy, origin string, paymentMethods []string, mgaProduct *models.Product) FabrickPaymentResponse {
	log.Printf("[FabrickMonthlyPay] Policy %s", data.Uid)

	customerId := uuid.New().String()
	firstres := <-FabrickPayObj(data, true, "", "", customerId, data.PriceGrossMonthly, data.PriceNettMonthly, origin, paymentMethods, mgaProduct, data.StartDate)
	time.Sleep(100)

	for i := 1; i <= 11; i++ {
		date := data.StartDate.AddDate(0, i, 0)
		expireDate := date.AddDate(10, 0, 0)
		res := <-FabrickPayObj(data, false, date.Format(models.TimeDateOnly), expireDate.Format(models.TimeDateOnly), customerId, data.PriceGrossMonthly, data.PriceNettMonthly, origin, paymentMethods, mgaProduct, date)
		log.Printf("[FabrickMonthlyPay] Policy %s - Index %d - response: %v", data.Uid, i, res)
		time.Sleep(100)
	}

	return firstres
}
