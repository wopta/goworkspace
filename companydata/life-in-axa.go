package companydata

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/wopta/goworkspace/network"
	"github.com/wopta/goworkspace/product"
	"github.com/wopta/goworkspace/user"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-gota/gota/dataframe"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func LifeIn(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	const (
		slide       int = -1
		headervalue     = "N° adesione individuale univoco"
	)
	var (
		policies                 = make([]models.Policy, 0)
		skippedPolicies          = make([]string, 0)
		missingBirthCityPolicies = make([]string, 0)
		missingProducerPolicies  = make([]string, 0)
		missingProducers         = make([]string, 0)
		wrongFiscalCodePolicies  = make([]string, 0)
		monthlyPolicies          = make(map[string]map[string][][]string, 0)
		result                   = make(map[string]map[string]interface{}, 0)
		codes                    map[string]map[string]string
	)

	taxesByGuarantee := map[string]float64{
		"death":                0,
		"permanent-disability": 0.025,
		"serious-ill":          0.025,
		"temporary-disability": 0.025,
	}

	b, err := os.ReadFile(lib.GetAssetPathByEnv("companyData") + "/reverse-codes.json")
	err = json.Unmarshal(b, &codes)
	if err != nil {
		return "", nil, err
	}

	data, _ := os.ReadFile("./companydata/track_in_life.csv")
	//data := lib.GetFromStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "track/in/life/life.csv", "")
	df := lib.CsvToDataframe(data)
	//log.Println("LifeIn  df.Describe: ", df.Describe())
	log.Println("LifeIn  row", df.Nrow())
	log.Println("LifeIn  col", df.Ncol())
	//group := df.GroupBy("N\xb0 adesione individuale univoco")
	group := GroupBy(df, 2)

	mgaProducts := map[string]*models.Product{
		models.ProductV1: product.GetProductV2(models.LifeProduct, models.ProductV1, models.MgaChannel, nil,
			nil),
		models.ProductV2: product.GetProductV2(models.LifeProduct, models.ProductV2, models.MgaChannel, nil,
			nil),
	}

	for v, pol := range group {
		var (
			row           []string
			guarantees    []models.Guarante
			sumPriceGross float64
			maxDuration   int
		)

		if pol[0][2] == headervalue {
			continue
		}
		if strings.TrimSpace(pol[0][13]) == "W1" || strings.TrimSpace(pol[0][22]) == "PG" {
			skippedPolicies = append(skippedPolicies, fmt.Sprintf("%07s", strings.TrimSpace(pol[0][2])))
			continue
		}

		row = pol[0]

		codeCompany := fmt.Sprintf("%07s", strings.TrimSpace(row[2]))
		payDate := fmt.Sprintf("%08s", strings.TrimSpace(row[5]))

		for i, r := range pol {
			var (
				beneficiaries []models.Beneficiary
			)

			log.Println("LifeIn  i: ", i)
			log.Println("LifeIn  pol: ", r)

			payDate = fmt.Sprintf("%08s", strings.TrimSpace(r[5]))

			if strings.TrimSpace(r[3]) == "R" {
				if monthlyPolicies[codeCompany] == nil {
					monthlyPolicies[codeCompany] = make(map[string][][]string, 0)
				}
				if monthlyPolicies[codeCompany][payDate] == nil {
					monthlyPolicies[codeCompany][payDate] = make([][]string, 0)
				}
				monthlyPolicies[codeCompany][payDate] = append(monthlyPolicies[codeCompany][payDate], r)
				continue
			}

			companyCodec, slug, _, _ := LifeMapCodecCompanyAxaRevert(r[1])
			if slug == "death" {
				if r[82] == "GE" {
					beneficiaries = append(beneficiaries, models.Beneficiary{
						BeneficiaryType: "legalAndWillSuccessor",
					})
				} else {
					benef1 := ParseAxaBeneficiary(r, 0)
					benef2 := ParseAxaBeneficiary(r, 1)
					if benef1 != nil {
						beneficiaries = append(beneficiaries, *benef1)
					}
					if benef2 != nil {
						beneficiaries = append(beneficiaries, *benef2)
					}
				}
			}
			dur, _ := strconv.Atoi(r[7])
			guaranteeYearDuration := dur / 12

			if guaranteeYearDuration > maxDuration {
				maxDuration = guaranteeYearDuration
			}

			priceGross := ParseAxaFloat(r[8])
			sumPriceGross += priceGross

			guarante := models.Guarante{
				Slug:                       slug,
				CompanyCodec:               companyCodec,
				SumInsuredLimitOfIndemnity: 0,
				Beneficiaries:              &beneficiaries,
				Value: &models.GuaranteValue{
					SumInsuredLimitOfIndemnity: lib.RoundFloat(ParseAxaFloat(r[9]), 0),
					PremiumGrossYearly:         lib.RoundFloat(priceGross, 2),
					Duration: &models.Duration{
						Year: guaranteeYearDuration,
					},
				},
			}

			guarantees = append(guarantees, guarante)
		}

		log.Println("LifeIn  value", v)
		log.Println("LifeIn  row", len(row))
		//log.Println("LifeIn  col", len(row))
		//log.Println("LifeIn  pol: ", pol)
		log.Println("LifeIn  elemets (0-0 ): ", row[0])
		log.Println("LifeIn  elemets (0-1 ): ", row[1])
		log.Println("LifeIn  elemets (0-2 ): ", row[2])
		log.Println("LifeIn  elemets (0-3 ): ", row[3])
		//1998-09-27T00:00:00Z RFC3339

		_, _, version, paymentSplit := LifeMapCodecCompanyAxaRevert(row[1])
		networkNode := network.GetNetworkNodeByCode(strings.TrimSpace(strings.ToUpper(row[13])))
		if networkNode == nil {
			log.Println("node not found!")
			missingProducerPolicies = append(missingProducerPolicies, codeCompany)
			skippedPolicies = append(skippedPolicies, fmt.Sprintf("%07s", strings.TrimSpace(row[2])))
			if !lib.SliceContains(missingProducers, strings.TrimSpace(strings.ToUpper(row[13]))) {
				missingProducers = append(missingProducers, strings.TrimSpace(strings.ToUpper(row[13])))
			}
			continue
		}

		// create insured

		insured := &models.User{
			Type:          row[22],
			Name:          strings.TrimSpace(lib.Capitalize(row[24])),
			Surname:       strings.TrimSpace(lib.Capitalize(row[23])),
			FiscalCode:    strings.TrimSpace(strings.ToUpper(row[27])),
			Gender:        strings.TrimSpace(strings.ToUpper(row[25])),
			BirthDate:     ParseDateDDMMYYYY(row[26]).Format(time.RFC3339),
			Phone:         row[72],
			BirthCity:     strings.TrimSpace(lib.Capitalize(row[73])),
			BirthProvince: strings.TrimSpace(strings.ToUpper(row[74])),
			Residence: &models.Address{
				StreetName: strings.TrimSpace(lib.Capitalize(row[28])),
				City:       strings.TrimSpace(lib.Capitalize(row[30])),
				CityCode:   strings.TrimSpace(strings.ToUpper(row[31])),
				PostalCode: row[29],
				Locality:   strings.TrimSpace(lib.Capitalize(row[30])),
			},
			CreationDate: ParseDateDDMMYYYY(row[4]),
			UpdatedDate:  time.Now().UTC(),
		}

		policy := models.Policy{
			Uid:            lib.NewDoc(models.PolicyCollection),
			Status:         models.PolicyStatusPay,
			StatusHistory:  []string{"Imported", models.PolicyStatusInitLead, models.PolicyStatusContact, models.PolicyStatusToSign, models.PolicyStatusSign, models.NetworkTransactionStatusToPay, models.PolicyStatusPay},
			Name:           models.LifeProduct,
			NameDesc:       "Wopta per te Vita",
			CodeCompany:    codeCompany,
			Company:        models.AxaCompany,
			ProductVersion: "v" + version,
			IsPay:          true,
			IsSign:         true,
			CompanyEmit:    true,
			CompanyEmitted: true,
			Channel:        models.NetworkChannel,
			PaymentSplit:   paymentSplit,
			CreationDate:   ParseDateDDMMYYYY(row[4]),
			EmitDate:       ParseDateDDMMYYYY(row[4]),
			StartDate:      ParseDateDDMMYYYY(row[4]),
			EndDate:        ParseDateDDMMYYYY(row[4]).AddDate(maxDuration, 0, 0),
			Updated:        time.Now().UTC(),
			PriceGross:     sumPriceGross,
			PriceNett:      0,
			Payment:        models.ManualPaymentProvider,
			FundsOrigin:    "Proprie risorse economiche",
			ProducerCode:   networkNode.Code,
			ProducerUid:    networkNode.Uid,
			ProducerType:   networkNode.Type,
			Assets: []models.Asset{{
				Guarantees: guarantees,
			}},
		}

		if policy.PaymentSplit == string(models.PaySplitMonthly) {
			calculateMonthlyPrices(&policy, taxesByGuarantee)
		}

		if policy.HasGuarantee("death") {

			// setting identity documents

			tmpCode, _ := strconv.Atoi(strings.TrimSpace(row[76]))
			identityDocumentCode := fmt.Sprintf("%02d", tmpCode)
			insured.IdentityDocuments = []*models.IdentityDocument{{
				Number:           strings.TrimSpace(strings.ToUpper(row[77])),
				Code:             identityDocumentCode,
				Type:             identityDocumentMap[identityDocumentCode],
				DateOfIssue:      ParseDateDDMMYYYY(row[78]),
				IssuingAuthority: strings.TrimSpace(lib.Capitalize(row[79])),
			}}

			for index, _ := range insured.IdentityDocuments {
				insured.IdentityDocuments[index].ExpiryDate = insured.IdentityDocuments[0].DateOfIssue.AddDate(10, 0, 0)
				insured.IdentityDocuments[index].LastUpdate = policy.EmitDate
			}

			// setting email

			insured.Mail = strings.TrimSpace(strings.ToLower(row[71]))

			// setting domicile

			insured.Domicile = &models.Address{
				StreetName: strings.TrimSpace(lib.Capitalize(row[67])),
				City:       strings.TrimSpace(lib.Capitalize(row[69])),
				CityCode:   strings.TrimSpace(strings.ToUpper(row[70])),
				PostalCode: row[68],
				Locality:   strings.TrimSpace(lib.Capitalize(row[69])),
			}

		}

		policy.Assets[0].Person = insured
		policy.Contractor = *insured
		policy.Contractor.Uid = lib.NewDoc(models.UserCollection)

		// check fiscalcode
		var usr models.User
		_, usr, err = user.CalculateFiscalCode(*insured)
		if err != nil {
			if strings.ToLower(err.Error()) == "invalid birth city" {
				_, extractedUser, _ := ExtractUserDataFromFiscalCode(insured.FiscalCode, codes)
				insured.BirthCity = extractedUser.BirthCity
				insured.BirthProvince = extractedUser.BirthProvince

				_, usr, err = user.CalculateFiscalCode(*insured)

				missingBirthCityPolicies = append(missingBirthCityPolicies, policy.CodeCompany)
			} else {
				log.Printf("error: %s", err.Error())
				continue
			}

		}

		if strings.ToUpper(usr.FiscalCode) != strings.ToUpper(insured.FiscalCode) {
			wrongFiscalCodePolicies = append(wrongFiscalCodePolicies, policy.CodeCompany)
			skippedPolicies = append(skippedPolicies, policy.CodeCompany)
			continue
		}

		result[codeCompany] = map[string]interface{}{
			"policy":              policy,
			"transactions":        []models.Transaction{},
			"networkTransactions": []*models.NetworkTransaction{},
		}

		// create transactions and network node transactions

		transactions := make([]models.Transaction, 0)
		networkTransactions := make([]*models.NetworkTransaction, 0)

		if monthlyPolicies[policy.CodeCompany] != nil {
			scheduleDate := policy.EmitDate
			transactionPayDate := policy.StartDate
			tr := createTransaction(policy, mgaProducts[policy.ProductVersion], "", scheduleDate, transactionPayDate, lib.RoundFloat(policy.PriceGross/12, 2), true)
			transactions = append(transactions, tr)
			networkTransactions = append(networkTransactions, createNetworkTransactions(&policy, &tr, networkNode, mgaProducts[policy.ProductVersion])...)
			for i := 1; i < 12; i++ {
				transactionPayDate = time.Time{}
				scheduleDate = scheduleDate.AddDate(0, 1, 0)
				isPay := false
				payDateString := scheduleDate.Format("02012006")
				if monthlyPolicies[policy.CodeCompany][payDateString] != nil {
					isPay = true
					transactionPayDate = scheduleDate
				}
				tr = createTransaction(policy, mgaProducts[policy.ProductVersion], "", scheduleDate, transactionPayDate, lib.RoundFloat(policy.PriceGross/12, 2), isPay)
				transactions = append(transactions, tr)
				if isPay {
					networkTransactions = append(networkTransactions, createNetworkTransactions(&policy, &tr, networkNode, mgaProducts[policy.ProductVersion])...)
				}
			}
		} else {
			tr := createTransaction(policy, mgaProducts[policy.ProductVersion], "", policy.EmitDate, policy.EmitDate, lib.RoundFloat(policy.PriceGross, 2), true)
			transactions = append(transactions, tr)
			networkTransactions = append(networkTransactions, createNetworkTransactions(&policy, &tr, networkNode, mgaProducts[policy.ProductVersion])...)

		}

		result[codeCompany]["transactions"] = transactions
		result[codeCompany]["networkTransactions"] = networkTransactions

		// update node portfolio

		networkNode.Policies = append(networkNode.Policies, policy.Uid)
		networkNode.Users = append(networkNode.Users, policy.Contractor.Uid)

		// save policy firestore

		// save policy bigquery

		// save single guarantees into bigquery

		// save contractor firestore

		// save contractor bigquery

		//log.Println("LifeIn policy:", policy)
		b, e := json.Marshal(policy)
		log.Println("LifeIn policy:", e)
		log.Println("LifeIn policy:", string(b))
		policies = append(policies, policy)
		// docref, _, _ := lib.PutFirestoreErr("test-policy", policy)
		// log.Println("LifeIn doc id: ", docref.ID)

		//_, e = models.UpdateUserByFiscalCode("uat", policy.Contractor)
		//log.Println("LifeIn policy:", policy)
		//tr := transaction.PutByPolicy(policy, "", "uat", "", "", sumPriseGross, 0, "", "manual", true)
		//	log.Println("LifeIn transactionpolicy:",tr)
		//accounting.CreateNetworkTransaction(tr, "uat")

	}

	log.Printf("Skipped %d policies: %v\n", len(skippedPolicies), skippedPolicies)
	log.Printf("Missing %d producers: %v\n", len(missingProducers), missingProducers)
	log.Printf("Wrong fiscal code %d policies: %v\n", len(wrongFiscalCodePolicies), wrongFiscalCodePolicies)
	log.Printf("Missing Birth City %d policies: %v\n", len(missingBirthCityPolicies), missingBirthCityPolicies)
	log.Printf("Missing Producer %d policies: %v\n", len(missingProducerPolicies), missingProducerPolicies)
	log.Printf("Created %d policies ", len(policies))

	out, err := json.Marshal(result)
	if err != nil {
		log.Printf("error: %s", err.Error())
	}
	err = os.WriteFile("./companydata/result.json", out, 0777)
	if err != nil {
		log.Printf("error: %s", err.Error())
	}

	return "", nil, e
}

func calculateMonthlyPrices(policy *models.Policy, taxesByGuarantee map[string]float64) {
	policy.PriceGrossMonthly = policy.PriceGross
	policy.PriceGross *= 12

	for index, guarantee := range policy.Assets[0].Guarantees {
		policy.Assets[0].Guarantees[index].Value.PremiumGrossMonthly = lib.RoundFloat(guarantee.Value.PremiumGrossYearly, 2)
		policy.Assets[0].Guarantees[index].Value.PremiumTaxAmountMonthly = lib.RoundFloat(policy.Assets[0].Guarantees[index].Value.PremiumGrossMonthly*taxesByGuarantee[guarantee.Slug], 2)
		policy.Assets[0].Guarantees[index].Value.PremiumNetMonthly = lib.RoundFloat(policy.Assets[0].Guarantees[index].Value.PremiumGrossMonthly-policy.Assets[0].Guarantees[index].Value.PremiumTaxAmountMonthly, 2)
		policy.TaxAmountMonthly += lib.RoundFloat(policy.Assets[0].Guarantees[index].Value.PremiumTaxAmountMonthly, 2)
		policy.PriceNettMonthly += lib.RoundFloat(policy.Assets[0].Guarantees[index].Value.PremiumNetMonthly, 2)

		policy.Assets[0].Guarantees[index].Value.PremiumGrossYearly = lib.RoundFloat(guarantee.Value.PremiumGrossYearly*12, 2)
		policy.Assets[0].Guarantees[index].Value.PremiumTaxAmountYearly = lib.RoundFloat(policy.Assets[0].Guarantees[index].Value.PremiumGrossYearly*taxesByGuarantee[guarantee.Slug], 2)
		policy.Assets[0].Guarantees[index].Value.PremiumNetYearly = lib.RoundFloat(policy.Assets[0].Guarantees[index].Value.PremiumGrossYearly-policy.Assets[0].Guarantees[index].Value.PremiumTaxAmountYearly, 2)
		policy.TaxAmount += lib.RoundFloat(policy.Assets[0].Guarantees[index].Value.PremiumTaxAmountYearly, 2)
		policy.PriceNett += lib.RoundFloat(policy.Assets[0].Guarantees[index].Value.PremiumNetYearly, 2)
	}
}

var identityDocumentMap map[string]string = map[string]string{
	"01": "Carta di Identità",
	"02": "Patente di Guida",
	"03": "Passaporto",
}

func LifeMapCodecCompanyAxaRevert(g string) (string, string, string, string) {
	log.Println("LifeIn LifeMapCodecCompanyAxaRevert:", g)
	var result, pay, slug, version string
	version = g[:1]
	code := g[2:3]
	payCode := g[1:2]

	switch payCode {
	case "W":
		pay = string(models.PaySplitYearly)
	case "M":
		pay = string(models.PaySplitMonthly)
	}

	if code == "5" {
		result = "D"
		slug = "death"
	}
	if code == "6" {
		result = "PTD"
		slug = "permanent-disability"
	}
	if code == "7" {
		result = "TTD"
		slug = "temporary-disability"
	}
	if code == "8" {
		result = "CI"
		slug = "serious-ill"
	}
	log.Println("LifeIn LifeMapCodecCompanyAxaRevert:", version)
	log.Println("LifeIn LifeMapCodecCompanyAxaRevert:", code)
	return result, slug, version, pay
}

func ParseDateDDMMYYYY(date string) time.Time {
	var (
		res time.Time
	)
	log.Println("LifeIn ParseDateDDMMYYYY date:", date)
	log.Println("LifeIn ParseDateDDMMYYYY len(date):", len(date))
	if len(date) == 7 {
		date = "0" + date
	}
	if len(date) == 8 {
		d, e := strconv.Atoi(date[:2])
		m, e := strconv.Atoi(date[2:4])
		y, e := strconv.Atoi(date[4:8])

		res = time.Date(y, time.Month(m),
			d, 0, 0, 0, 0, time.UTC)
		log.Println(e)
		log.Println("LifeIn ParseDateDDMMYYYY d:", d)
		log.Println("LifeIn ParseDateDDMMYYYY m:", m)
		log.Println("LifeIn ParseDateDDMMYYYY y:", y)
		log.Println("LifeIn ParseDateDDMMYYYY res:", res)
	}
	return res

}

func ParseAxaFloat(price string) float64 {
	princeInCents, _ := strconv.ParseFloat(price, 64)
	return princeInCents / 100.0
}

func ParseAxaBeneficiary(r []string, base int) *models.Beneficiary {
	var (
		benef *models.Beneficiary
	)
	rangeCell := 11 * base

	if r[82] == "GE" {
		benef = &models.Beneficiary{
			BeneficiaryType: "legalAndWillSuccessor",
		}
	}
	if r[82] == "NM" {
		if strings.TrimSpace(strings.ToUpper(r[85+rangeCell])) == "" || strings.TrimSpace(strings.ToUpper(r[85+rangeCell])) == "0" {
			return nil
		}

		benef = &models.Beneficiary{
			BeneficiaryType: "chosenBeneficiary",
			User: models.User{
				Name:       strings.TrimSpace(lib.Capitalize(r[84+rangeCell])),
				Surname:    strings.TrimSpace(lib.Capitalize(r[83+rangeCell])),
				FiscalCode: strings.TrimSpace(strings.ToUpper(r[85+rangeCell])),
				Mail:       strings.TrimSpace(strings.ToLower(r[91+rangeCell])),
				Residence: &models.Address{
					StreetName: strings.TrimSpace(lib.Capitalize(r[87+rangeCell])),
					City:       strings.TrimSpace(lib.Capitalize(r[88+rangeCell])),
					CityCode:   strings.TrimSpace(strings.ToUpper(r[90+rangeCell])),
					PostalCode: strings.TrimSpace(r[89+rangeCell]),
					Locality:   strings.TrimSpace(lib.Capitalize(r[88+rangeCell])),
				},
			},
		}
	}
	return benef

}
func GroupBy(df dataframe.DataFrame, col int) map[string][][]string {
	log.Println("GroupBy")
	res := make(map[string][][]string)
	for _, k := range df.Records() {
		if _, found := res[k[col]]; found {
			res[k[col]] = append(res[k[col]], k)
		} else {
			res[k[col]] = [][]string{k}
		}
	}
	return res
}

func createTransaction(policy models.Policy, mgaProduct *models.Product, customerId string, scheduleDate time.Time, payDate time.Time, priceGross float64, isPay bool) models.Transaction {
	status := models.TransactionStatusToPay
	statusHistory := []string{models.TransactionStatusToPay}
	paymentMethod := ""

	if isPay {
		status = models.TransactionStatusPay
		statusHistory = append(statusHistory, models.TransactionStatusPay)
		paymentMethod = models.PayMethodTransfer
	}

	expireDate := scheduleDate.AddDate(10, 0, 0)

	return models.Transaction{
		Amount:          priceGross,
		Uid:             lib.NewDoc(models.TransactionsCollection),
		PolicyName:      policy.Name,
		PolicyUid:       policy.Uid,
		CreationDate:    policy.EmitDate,
		UpdateDate:      time.Now().UTC(),
		Status:          status,
		StatusHistory:   statusHistory,
		ScheduleDate:    scheduleDate.Format(models.TimeDateOnly),
		ExpirationDate:  expireDate.Format(models.TimeDateOnly),
		NumberCompany:   policy.CodeCompany,
		IsPay:           isPay,
		PayDate:         payDate,
		TransactionDate: payDate,
		Name:            policy.Contractor.Name + " " + policy.Contractor.Surname,
		Company:         policy.Company,
		IsDelete:        false,
		UserToken:       customerId,
		ProviderName:    policy.Payment,
		PaymentMethod:   paymentMethod,
		Commissions:     lib.RoundFloat(product.GetCommissionByProduct(&policy, mgaProduct, false), 2),
	}
}

func createNetworkTransaction(
	policy *models.Policy,
	transaction *models.Transaction,
	node *models.NetworkNode,
	commission float64, // Amount
	mgaAccountType, paymentType, name string,
) (*models.NetworkTransaction, error) {
	var amount float64

	switch paymentType {
	case models.PaymentTypeRemittanceCompany, models.PaymentTypeCommission:
		amount = lib.RoundFloat(commission, 2)
	case models.PaymentTypeRemittanceMga:
		amount = lib.RoundFloat(transaction.Amount-commission, 2)
	}

	netTransaction := models.NetworkTransaction{
		Uid:              uuid.New().String(),
		PolicyUid:        policy.Uid,
		TransactionUid:   transaction.Uid,
		NetworkNodeUid:   node.Uid,
		NetworkNodeType:  node.Type,
		AccountType:      mgaAccountType,
		PaymentType:      paymentType,
		Amount:           amount,
		AmountNet:        amount, // TBD
		Name:             name,
		Status:           models.NetworkTransactionStatusPaid,
		StatusHistory:    []string{models.NetworkTransactionStatusCreated, models.NetworkTransactionStatusToPay, models.NetworkTransactionStatusPaid},
		IsPay:            true,
		IsConfirmed:      false,
		IsDelete:         false,
		CreationDate:     lib.GetBigQueryNullDateTime(transaction.PayDate),
		PayDate:          lib.GetBigQueryNullDateTime(transaction.PayDate),
		TransactionDate:  lib.GetBigQueryNullDateTime(transaction.PayDate),
		ConfirmationDate: lib.GetBigQueryNullDateTime(time.Time{}),
		DeletionDate:     lib.GetBigQueryNullDateTime(time.Time{}),
	}

	//jsonLog, _ := json.Marshal(&netTransaction)

	/*err := netTransaction.SaveBigQuery()
	if err != nil {
		log.Printf("[createNetworkTransaction] error saving network transaction to bigquery: %s", err.Error())
		return nil, err
	}*/

	//log.Printf("[createNetworkTransaction] network transaction created! %s", string(jsonLog))

	return &netTransaction, nil
}

func createCompanyNetworkTransaction(
	policy *models.Policy,
	transaction *models.Transaction,
	producerNode *models.NetworkNode,
	mgaProduct *models.Product,
) (*models.NetworkTransaction, error) {
	var code string

	commissionMga := product.GetCommissionByProduct(policy, mgaProduct, false)
	commissionCompany := lib.RoundFloat(transaction.Amount-commissionMga, 2)
	code = producerNode.Code

	name := strings.ToUpper(strings.Join([]string{code, policy.Company}, "-"))

	return createNetworkTransaction(
		policy,
		transaction,
		&models.NetworkNode{},
		commissionCompany,
		models.AccountTypePassive,
		models.PaymentTypeRemittanceCompany,
		name,
	)
}

func createNetworkTransactions(
	policy *models.Policy,
	transaction *models.Transaction,
	producerNode *models.NetworkNode,
	mgaProduct *models.Product,
) []*models.NetworkTransaction {
	var err error

	networkTransactions := make([]*models.NetworkTransaction, 0)

	if policy.CodeCompany == "0000071" {
		log.Printf("hello")
	}

	nt, err := createCompanyNetworkTransaction(policy, transaction, producerNode, mgaProduct)
	if err != nil {
		log.Printf("[CreateNetworkTransactions] error creating company network-transaction: %s", err.Error())
		return nil
	}

	networkTransactions = append(networkTransactions, nt)

	if policy.ProducerUid != "" && policy.ProducerType != models.PartnershipNetworkNodeType {
		network.TraverseWithCallbackNetworkByNodeUid(producerNode, "", func(currentNode *models.NetworkNode, currentName string) string {
			var (
				accountType, paymentType string
				baseName                 string
			)

			warrant := currentNode.GetWarrant()
			if warrant == nil {
				log.Printf("[CreateNetworkTransactions] error getting warrant for node: %s", currentNode.Uid)
				return baseName
			}
			prod := warrant.GetProduct(policy.Name)
			if warrant == nil {
				log.Printf("[CreateNetworkTransactions] error getting product for warrant: %s", warrant.Name)
				return baseName
			}

			accountType = getAccountType(transaction)
			paymentType = getPaymentType(transaction, policy, currentNode)
			commission := product.GetCommissionByProduct(policy, prod, policy.ProducerUid == currentNode.Uid)

			if currentName != "" {
				baseName = strings.ToUpper(strings.Join([]string{currentName, currentNode.Code}, "__"))
			} else {
				baseName = strings.ToUpper(currentNode.Code)
			}
			nodeName := strings.ToUpper(strings.Join([]string{
				baseName,
				strings.ReplaceAll(currentNode.GetName(), " ", "-"),
			}, "-"))

			nt, err = createNetworkTransaction(policy, transaction, currentNode, commission, accountType, paymentType, nodeName)
			if err != nil {
				log.Printf("[CreateNetworkTransactions] error creating network-transaction: %s", err.Error())
			} else {
				log.Printf("[CreateNetworkTransactions] created network-transaction for node: %s", currentNode.Uid)
			}

			networkTransactions = append(networkTransactions, nt)
			return baseName
		})
	}

	return networkTransactions
}

func getAccountType(transaction *models.Transaction) string {
	if transaction.PaymentMethod == models.PayMethodRemittance {
		return models.AccountTypeActive
	}
	return models.AccountTypePassive
}

func getPaymentType(transaction *models.Transaction, policy *models.Policy, producerNode *models.NetworkNode) string {
	if policy.ProducerUid == producerNode.Uid && transaction.PaymentMethod == models.PayMethodRemittance {
		return models.PaymentTypeRemittanceMga
	}
	return models.PaymentTypeCommission
}
