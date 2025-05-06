package quote

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	netclient "github.com/wopta/goworkspace/models/client"
	"github.com/wopta/goworkspace/models/dto/net"
	"github.com/wopta/goworkspace/network"
	prd "github.com/wopta/goworkspace/product"
)

func CatNatFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err       error
		reqPolicy *models.Policy
	)

	log.SetPrefix("[CatNatFx] ")
	defer func() {
		r.Body.Close()
		if err != nil {
			log.Printf("error: %s", err.Error())
		}
		log.Println("Handler end ---------------------------------------------")
		log.SetPrefix("")
	}()
	log.Println("Handler start -----------------------------------------------")

	_, err = lib.GetAuthTokenFromIdToken(r.Header.Get("Authorization"))
	if err != nil {
		log.Printf("error getting authToken")
		return "", nil, err
	}

	if err = json.NewDecoder(r.Body).Decode(&reqPolicy); err != nil {
		log.Println("error decoding request body")
		return "", nil, err
	}

	var cnReq net.RequestDTO
	err = cnReq.FromPolicy(reqPolicy, false)
	if err != nil {
		log.Printf("error building NetInsurance DTO: %s", err.Error())
		return "", nil, err
	}

	netClient := netclient.NewNetClient()
	netClient.Authenticate()

	resp, errResp, err := netClient.Quote(cnReq)
	if err != nil {
		log.Printf("error calling NetInsurance api: %s", err.Error())
		return "", nil, err
	}
	var out []byte
	if errResp != nil {
		out, err = json.Marshal(errResp)
		if err != nil {
			log.Println("error encoding response %w", err.Error())
			return "", nil, err
		}

		return string(out), out, err
	}

	if resp.Result != "OK" {
		out, err = json.Marshal(resp)
		if err != nil {
			log.Println("error encoding response %w", err.Error())
			return "", nil, err
		}
		return string(out), out, err
	}

	_ = resp.ToPolicy(reqPolicy)

	networkNode := network.GetNetworkNodeByUid(reqPolicy.ProducerUid)
	warrant := networkNode.GetWarrant()
	product := prd.GetProductV2(reqPolicy.Name, reqPolicy.ProductVersion, reqPolicy.Channel, networkNode, warrant)
	addConsultacyPrice(reqPolicy, product)

	out, err = json.Marshal(reqPolicy)
	if err != nil {
		log.Println("error encoding response %w", err.Error())
		return "", nil, err
	}

	return string(out), out, err
}
