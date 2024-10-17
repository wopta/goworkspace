package broker

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func GetSequenceByCompany(companyName, productName string, firePolicy string) (string, int, int) {
	var (
		codeCompany         string
		companyDefault      int
		companyPrefix       string
		companyPrefixLenght string
		numberCompany       int
		number              int
	)
	switch companyName {
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
	}

	ctx := context.Background()
	client, err := firestore.NewClient(ctx, os.Getenv("GOOGLE_PROJECT_ID"))
	lib.CheckError(err)
	docSnap := client.Collection(lib.PolicyCollection).
		Where("company", "==", companyName).
		Where("name", "==", productName).
		OrderBy("numberCompany", firestore.Desc).
		Limit(1).
		Documents(ctx)

	policy := models.PolicyToListData(docSnap)
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
func GetSequenceProposal(name string, firePolicy string) int {
	var number int
	r, e := lib.OrderLimitFirestoreErr(firePolicy, "proposalNumber", firestore.Desc, 1)
	lib.CheckError(e)
	policyCompany := models.PolicyToListData(r)
	if len(policyCompany) == 0 {
		number = 1
	} else {

		number = policyCompany[0].ProposalNumber + 1
	}
	log.Println("GetSequenceProposal: ", number)
	return number
}
