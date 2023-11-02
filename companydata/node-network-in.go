package companydata

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	lib "github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func NodeNetworkIn(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	const ()
	var (
		agent       *models.AgentNode
		agency      *models.AgencyNode
		broker      *models.AgencyNode
		areaManager *models.AgentNode
	)
	data := lib.GetFromStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "track/in/life/life.csv", "")
	df := lib.CsvToDataframe(data)
	//log.Println("LifeIn  df.Describe: ", df.Describe())
	log.Println("LifeIn  row", df.Nrow())
	log.Println("LifeIn  col", df.Ncol())
	//group := df.GroupBy("N\xb0 adesione individuale univoco")

	for _, d := range df.Records() {
		if d[0] == "agent" {
			agent = &models.AgentNode{
				Name:            d[0],
				Surname:         d[0],
				FiscalCode:      d[0],
				VatCode:         d[0],
				RuiCode:         d[0],
				RuiSection:      d[0],
				RuiRegistration: ParseDateDDMMYYYY(d[0]),
			}
		}
		if d[0] == "agency" {
			agency = &models.AgencyNode{}
		}
		if d[0] == "area" {
			areaManager = &models.AgentNode{}
		}
		if d[0] == "broker" {
			broker = &models.AgencyNode{}
		}

		node := models.NetworkNode{
			Type:d[0] ,
			Agent:       agent,
			Agency:      agency,
			Broker:      broker,
			AreaManager: areaManager,
			Code:  d[0],
			NetworkCode:  d[0],
			Mail:  d[0],
			Warrant:  d[0],
			ParentUid: "",
			NetworkUid:"",
			ManagerUid: "",
			
		}

		b, e := json.Marshal(node)
		log.Println("LifeIn policy:", e)
		log.Println("LifeIn policy:", string(b))
		docref, _, _ := lib.PutFirestoreErr("test-network-node", node)
		log.Println("LifeIn doc id: ", docref.ID)
	

	}

	return "", nil, e
}
