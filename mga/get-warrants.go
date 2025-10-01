package mga

import (
	"encoding/json"
	"net/http"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/network"
)

type warrantProduct struct {
	Name string `json:"name"`
	Flow string `json:"flow"`
}

type warrant struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Products    []warrantProduct `json:"products"`
}

func (w *warrant) fromDomain(input models.Warrant) {
	w.Name = input.Name
	w.Description = input.Description
	products := make([]warrantProduct, 0, len(input.Products))
	for _, p := range input.Products {
		products = append(products, warrantProduct{
			Name: p.Name,
			Flow: p.Flow,
		})
	}
	w.Products = products
}

type getWarrantsResponse struct {
	Warrants []warrant `json:"warrants"`
}

func getWarrantsFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err  error
		resp getWarrantsResponse
	)

	log.AddPrefix("GetWarrantsFx")
	log.Println("Handler start -----------------------------------------------")

	defer func() {
		if err != nil {
			log.ErrorF("error: %s", err.Error())
		}
		log.Println("Handler end -------------------------------------------------")
		log.PopPrefix()
	}()

	retrievedWarrant, err := network.GetWarrants()
	if err != nil {
		return "", "", err
	}

	warrants := make([]warrant, 0, len(retrievedWarrant))
	for _, war := range retrievedWarrant {
		dto := new(warrant)
		dto.fromDomain(war)
		warrants = append(warrants, *dto)
	}

	resp.Warrants = warrants

	rawResp, err := json.Marshal(resp)
	if err != nil {
		return "", "", err
	}

	return string(rawResp), resp, nil
}
