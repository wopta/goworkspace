package payment

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
	prd "github.com/wopta/goworkspace/product"
)

func FabrickPayFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
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

	resultPay := <-FabrickPayObj(data, false, "", data.StartDate.AddDate(10, 0, 0).Format(models.TimeDateOnly), "", data.PriceGross,
		data.PriceNett, getOrigin(r.Header.Get("origin")), paymentMethods, mgaProduct, data.StartDate)

	log.Println(resultPay)
	return "", nil, err
}

func FabrickYearPay(data models.Policy, origin string, paymentMethods []string, mgaProduct *models.Product) FabrickPaymentResponse {
	log.Printf("[FabrickYearPay] Policy %s", data.Uid)

	customerId := uuid.New().String()
	res := <-FabrickPayObj(data, false, "", data.StartDate.AddDate(10, 0, 0).Format(models.TimeDateOnly), customerId, data.PriceGross, data.PriceNett, origin, paymentMethods, mgaProduct, data.StartDate)

	return res
}
