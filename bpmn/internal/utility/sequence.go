package utility

import (
	"fmt"

	"cloud.google.com/go/firestore"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/models"
)

func GetSequenceByCompany(name string, firePolicy string) (string, int, int) {
	var (
		codeCompany         string
		companyDefault      int
		companyPrefix       string
		companyPrefixLenght string
		numberCompany       int
		number              int
	)
	switch name {
	case models.GlobalCompany:
		companyDefault = 1
		companyPrefix = "WB"
		companyPrefixLenght = `%07d`
	case models.AxaCompany:
		companyDefault = 100001
		companyPrefixLenght = `%06d`
	case models.SogessurCompany:
		companyDefault = 1
		companyPrefixLenght = `%07d`
		companyPrefix = "G"
	case "qbe": // TODO
		companyDefault = 1
		companyPrefixLenght = `%07d`
		companyPrefix = "Q"
	}

	rn, e := lib.OrderWhereLimitFirestoreErr(firePolicy, "company", "numberCompany", "==", name, firestore.Desc, 1)
	log.Println(e)

	policy := models.PolicyToListData(rn)
	log.Println("len(policy):", len(policy))
	if len(policy) == 0 {
		//WE0000001
		numberCompany = companyDefault
		codeCompany = companyPrefix + fmt.Sprintf(companyPrefixLenght, numberCompany)
		number = 1
	} else {
		numberCompany = policy[0].NumberCompany + 1
		log.Println("numberCompany:", numberCompany)
		codeCompany = companyPrefix + fmt.Sprintf(companyPrefixLenght, numberCompany)
		number = policy[0].Number + 1
	}
	r, e := lib.OrderLimitFirestoreErr(firePolicy, "number", firestore.Desc, 1)
	log.Println(e)
	policyCompany := models.PolicyToListData(r)
	if len(policyCompany) == 0 {
		number = 1
	} else {

		number = policyCompany[0].Number + 1
	}
	log.Println("GetSequenceByCompany: ", codeCompany)

	return codeCompany, numberCompany, number
}
