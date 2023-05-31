package companydata

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	lib "github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func LifeAxaDelete(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {

	now := time.Now()
	fromM := time.Now().AddDate(0, -1, 0)
	var (
		result [][]string
	)
	q := lib.Firequeries{
		Queries: []lib.Firequery{{
			Field:      "companyEmit", //
			Operator:   "==",          //
			QueryValue: true,
		},
			{
				Field:      "companyEmitted", //
				Operator:   "==",             //
				QueryValue: false,
			},
			{
				Field:      "company", //
				Operator:   "==",      //
				QueryValue: "axa",
			},
			{
				Field:      "name", //
				Operator:   "==",   //
				QueryValue: "life",
			}, {
				Field:      "startSate", //
				Operator:   ">",         //
				QueryValue: strconv.Itoa(int(fromM.Unix())),
			},
			{
				Field:      "startSate", //
				Operator:   "<",         //
				QueryValue: strconv.Itoa(int(now.Unix())),
			},
		},
	}
	query, e := q.FirestoreWherefields("policy")
	policies := models.PolicyToListData(query)
	result = append(result, getHeader())
	for _, policy := range policies {

		for _, asset := range policy.Assets {

			for _, g := range asset.Guarantees {

				fmt.Println(g)
				row := []string{
					mapCodecCompany(policy, g.CompanyCodec), //Codice schema
					policy.CodeCompany,                      //Tipo Rimborso
					"A",                                     //Motivo Transazione
					policy.CodeCompany,                      //N° adesione
					getFormatdate(policy.StartDate),         //"Inizio copertura"
					"012",                                   //"Data estinzione"
					fmt.Sprintf("%.2f", g.PriceGross),       //"Importo di rimborso"
					policy.Contractor.Surname,               //"Cognome "
					policy.Contractor.Name,                  //"Nome"

				}

				result = append(result, row)

			}

		}

	}

	refMontly := now.AddDate(0, -1, 0)
	//year, month, day := time.Now().Date()
	//year2, month2, day2 := time.Now().AddDate(0, -1, 0).Date()
	filepath := "WOPTAKEY_CANM_" + strconv.Itoa(refMontly.Year()) + fmt.Sprintf("%02d", int(refMontly.Month())) + "_" + fmt.Sprintf("%02d", now.Day()) + fmt.Sprintf("%02d", int(now.Month())) + ".txt"
	lib.WriteCsv("../tmp/"+filepath, result)
	source, _ := ioutil.ReadFile("../tmp/" + filepath)
	lib.PutToStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "axa/life/"+filepath, source)
	SftpUpload(filepath)
	return "", nil, e
}

//Codice schema

func getHeaderDelete() []string {
	return []string{
		"Tipo Rimborso",
		"Motivo Transazione",
		"N° adesione",
		"Inizio copertura",
		"Data estinzione",
		"Importo di rimborso",
		"Cognome",
		"Nome",
	}
}
