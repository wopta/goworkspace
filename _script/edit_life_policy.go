package _script

import (
	"log"
	"os"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func EditLifePolicy(policyUid string) {
	rawData, err := os.ReadFile("./_script/policy_80.txt")
	if err != nil {
		log.Fatal(err)
	}

	df, err := lib.CsvToDataframeV2(rawData, ';', false)
	if err != nil {
		log.Fatal(err)
	}

	groups := groupBy(df, 2)
	delete(groups, "X2")

	for _, rawPolicy := range groups {
		log.Printf("%v", rawPolicy)

		// TODO: implementare estrazione dati contraente persona giuridica
		// TODO: implementare estrazione dati 3 titolari effettivi
		// TODO: implementare estrazione dati assicurato
	}

}

func extractContractorData(rawPolicy []string) models.Contractor {
	contractor := models.Contractor{}

	return contractor
}
