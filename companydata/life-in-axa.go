package companydata

import (
	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/civil"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/wopta/goworkspace/network"
	"github.com/wopta/goworkspace/product"
	"github.com/wopta/goworkspace/user"
	"io"
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

type ResultStruct struct {
	Policy       models.Policy                 `json:"policy"`
	Transactions map[string]TransactionsOutput `json:"transactions"`
}

type TransactionsOutput struct {
	Transaction         models.Transaction           `json:"transaction"`
	NetworkTransactions []*models.NetworkTransaction `json:"networkTransactions"`
}

var (
	skippedPolicies                    = make([]string, 0)
	missingContractorBirthCityPolicies = make([]string, 0)
	missingInsuredBirthCityPolicies    = make([]string, 0)
	missingProducerPolicies            = make([]string, 0)
	missingProducers                   = make([]string, 0)
	wrongContractorFiscalCodePolicies  = make([]string, 0)
	wrongInsuredFiscalCodePolicies     = make([]string, 0)
	monthlyPolicies                    = make(map[string]map[string][][]string, 0)
)

type LifeInReq struct {
	DryRun           *bool  `json:"dryRun"`
	CollectionPrefix string `json:"collectionPrefix"`
}

func LifeInFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	const (
		slide            int = -1
		headervalue          = "N° adesione individuale univoco"
		titleHeaderValue     = "DATI DEL CONTRATTO DI ASSICURAZIONE"
	)
	var (
		policies                 = make([]models.Policy, 0)
		result                   = make(map[string]ResultStruct, 0)
		codes                    map[string]map[string]string
		startDateJob, endDateJob time.Time
		contractorEqualInsured   bool
		insured                  *models.User
		req                      LifeInReq
	)

	startDateJob = time.Now().UTC()

	taxesByGuarantee := map[string]float64{
		"death":                0,
		"permanent-disability": 0.025,
		"serious-ill":          0.025,
		"temporary-disability": 0.025,
	}

	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()
	err := json.Unmarshal(body, &req)
	if err != nil {
		log.Printf("error unmrashalling request body")
		return "", nil, err
	}

	b, err := os.ReadFile(lib.GetAssetPathByEnv("companyData") + "/reverse-codes.json")
	err = json.Unmarshal(b, &codes)
	if err != nil {
		return "", nil, err
	}

	data, _ := os.ReadFile("./companydata/track_in_life.csv")
	//data := lib.GetFromStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "track/in/life/life.csv", "")
	df := lib.CsvToDataframe(data)
	//log.Println("LifeInFx  df.Describe: ", df.Describe())
	log.Println("LifeInFx  row", df.Nrow())
	log.Println("LifeInFx  col", df.Ncol())
	//group := df.GroupBy("N\xb0 adesione individuale univoco")
	group := GroupBy(df, 2)

	mgaProducts := map[string]*models.Product{
		models.ProductV1: product.GetProductV2(models.LifeProduct, models.ProductV1, models.MgaChannel, nil,
			nil),
		models.ProductV2: product.GetProductV2(models.LifeProduct, models.ProductV2, models.MgaChannel, nil,
			nil),
	}

	networkProducts := map[string]*models.Product{
		models.ProductV1: product.GetProductV2(models.LifeProduct, models.ProductV1, models.NetworkChannel, nil,
			nil),
		models.ProductV2: product.GetProductV2(models.LifeProduct, models.ProductV2, models.NetworkChannel, nil,
			nil),
	}

	for v, pol := range group {
		var (
			row                                                                 []string
			guarantees                                                          []models.Guarante
			sumPriceGross, sumPriceTaxAmount, sumPriceNett                      float64
			sumPriceGrossMonthly, sumPriceTaxAmountMonthly, sumPriceNettMonthly float64
			maxDuration                                                         int
		)

		if pol[0][2] == headervalue || pol[0][1] == titleHeaderValue || pol[0][1] == "1" {
			continue
		}

		row = pol[0]

		codeCompany := fmt.Sprintf("%07s", strings.TrimSpace(row[2]))
		payDate := fmt.Sprintf("%08s", strings.TrimSpace(row[5]))

		for i, r := range pol {
			var (
				beneficiaries []models.Beneficiary
			)

			log.Println("LifeInFx  i: ", i)
			log.Println("LifeInFx  pol: ", r)

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

			companyCodec, slug, version, _ := LifeMapCodecCompanyAxaRevert(r[1])
			if slug == "death" {
				for i := 0; i < 2; i++ {
					benef := ParseAxaBeneficiary(r, i)
					if benef != nil {
						beneficiaries = append(beneficiaries, *benef)
						if lib.SliceContains([]string{models.BeneficiarySelfLegalEntity, models.BeneficiaryLegalAndWillSuccessors}, benef.BeneficiaryType) {
							break
						}
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
			sumPriceGrossMonthly += priceGross

			productVersion := fmt.Sprintf("v%s", version)

			guarante := models.Guarante{
				Slug:                       slug,
				CompanyCodec:               companyCodec,
				Description:                mgaProducts[productVersion].Companies[0].GuaranteesMap[slug].Description,
				Group:                      mgaProducts[productVersion].Companies[0].GuaranteesMap[slug].Group,
				Type:                       mgaProducts[productVersion].Companies[0].GuaranteesMap[slug].Type,
				Name:                       mgaProducts[productVersion].Companies[0].GuaranteesMap[slug].Name,
				CompanyName:                mgaProducts[productVersion].Companies[0].GuaranteesMap[slug].CompanyName,
				SumInsuredLimitOfIndemnity: 0,
				Beneficiaries:              &beneficiaries,
				Value: &models.GuaranteValue{
					SumInsuredLimitOfIndemnity: lib.RoundFloat(ParseAxaFloat(r[9]), 0),
					PremiumGrossYearly:         lib.RoundFloat(priceGross, 2),
					PremiumTaxAmountYearly:     lib.RoundFloat(priceGross*taxesByGuarantee[slug], 2),
					PremiumNetYearly:           lib.RoundFloat(priceGross-(priceGross*taxesByGuarantee[slug]), 2),
					PremiumGrossMonthly:        lib.RoundFloat(priceGross, 2),
					PremiumTaxAmountMonthly:    lib.RoundFloat(priceGross*taxesByGuarantee[slug], 2),
					PremiumNetMonthly:          lib.RoundFloat(priceGross-(priceGross*taxesByGuarantee[slug]), 2),
					Duration: &models.Duration{
						Year: guaranteeYearDuration,
					},
				},
				Config:         networkProducts[productVersion].Companies[0].GuaranteesMap[slug].Config,
				IsSellable:     true,
				IsSelected:     true,
				IsConfigurable: true,
				Order:          mgaProducts[productVersion].Companies[0].GuaranteesMap[slug].Order,
			}

			if guarante.Slug == "temporary-disability" {
				guarante.Value.SumInsuredLimitOfIndemnity = lib.RoundFloat(ParseAxaFloat(r[10]), 0)
			} else if guarante.Slug == "death" {
				guarante.BeneficiaryOptions = mgaProducts[productVersion].Companies[0].GuaranteesMap["death"].BeneficiaryOptions
			}

			sumPriceTaxAmount += guarante.Value.PremiumTaxAmountYearly
			sumPriceNett += guarante.Value.PremiumNetYearly
			sumPriceTaxAmountMonthly += guarante.Value.PremiumTaxAmountMonthly
			sumPriceNettMonthly += guarante.Value.PremiumNetMonthly

			guarantees = append(guarantees, guarante)
		}

		log.Println("LifeInFx  value", v)
		log.Println("LifeInFx  row", len(row))
		//log.Println("LifeInFx  col", len(row))
		//log.Println("LifeInFx  pol: ", pol)
		log.Println("LifeInFx  elemets (0-0 ): ", row[0])
		log.Println("LifeInFx  elemets (0-1 ): ", row[1])
		log.Println("LifeInFx  elemets (0-2 ): ", row[2])
		log.Println("LifeInFx  elemets (0-3 ): ", row[3])
		//1998-09-27T00:00:00Z RFC3339

		_, _, version, paymentSplit := LifeMapCodecCompanyAxaRevert(row[1])
		nodeCode := strings.TrimSpace(strings.ToUpper(row[13]))
		if nodeCode == "W1" {
			nodeCode = "W1.DIRAgent"
		}
		networkNode := network.GetNetworkNodeByCode(nodeCode)
		if networkNode == nil {
			log.Println("node not found!")
			missingProducerPolicies = append(missingProducerPolicies, codeCompany)
			skippedPolicies = append(skippedPolicies, fmt.Sprintf("%07s", strings.TrimSpace(row[2])))
			if !lib.SliceContains(missingProducers, strings.TrimSpace(strings.ToUpper(row[13]))) {
				missingProducers = append(missingProducers, strings.TrimSpace(strings.ToUpper(row[13])))
			}
			continue
		}

		// check if contractor equal insured

		contractorEqualInsured = strings.TrimSpace(strings.ToUpper(row[45])) == "SI"

		// create contractor

		isLegalEntity := row[22] == "PG"

		var contractor *models.Contractor
		contractors := new([]models.User)
		if isLegalEntity {
			// parsing contractor

			contractor = parseEnterpriseContractor(row)
			if contractor == nil {
				skippedPolicies = append(skippedPolicies, fmt.Sprintf("%07s", strings.TrimSpace(row[2])))
				continue
			}

			// parsing esecutore info

			esecutore := parseEsecutore(row)
			*contractors = append(*contractors, esecutore)

			// parsing titolare effettivo info

			offset := 26
			titolariEffettivi := make([]models.User, 0)
			for i := 0; i < 3; i++ {
				if strings.TrimSpace(strings.ToUpper(row[116+(offset*i)])) == "" || strings.TrimSpace(strings.ToUpper(row[116+(offset*i)])) == "NO" {
					break
				}
				titolareEffettivo := parsingTitolareEffettivo(row, offset, i)
				titolariEffettivi = append(titolariEffettivi, titolareEffettivo)
			}
			*contractors = append(*contractors, titolariEffettivi...)
		} else {
			contractor = parseIndividualContractor(codeCompany, row, codes)
			if contractor == nil {
				skippedPolicies = append(skippedPolicies, fmt.Sprintf("%07s", strings.TrimSpace(row[2])))
				continue
			}
		}

		if !contractorEqualInsured {
			// create insured
			insured = parseInsured(codeCompany, row, codes)
			if insured == nil {
				skippedPolicies = append(skippedPolicies, fmt.Sprintf("%07s", strings.TrimSpace(row[2])))
				continue
			}
		} else {
			insured = contractor.ToUser()
		}

		policy := models.Policy{
			Uid:               lib.NewDoc(models.PolicyCollection),
			Status:            models.PolicyStatusPay,
			StatusHistory:     []string{"Imported", models.PolicyStatusInitLead, models.PolicyStatusContact, models.PolicyStatusToSign, models.PolicyStatusSign, models.NetworkTransactionStatusToPay, models.PolicyStatusPay},
			Name:              models.LifeProduct,
			NameDesc:          "Wopta per te Vita",
			CodeCompany:       codeCompany,
			Company:           models.AxaCompany,
			ProductVersion:    "v" + version,
			IsPay:             true,
			IsSign:            true,
			CompanyEmit:       true,
			CompanyEmitted:    true,
			Channel:           models.NetworkChannel,
			PaymentSplit:      paymentSplit,
			CreationDate:      ParseDateDDMMYYYY(row[4]),
			EmitDate:          ParseDateDDMMYYYY(row[4]),
			StartDate:         ParseDateDDMMYYYY(row[4]),
			EndDate:           ParseDateDDMMYYYY(row[4]).AddDate(maxDuration, 0, 0),
			Updated:           time.Now().UTC(),
			PriceGross:        lib.RoundFloat(sumPriceGross, 2),
			PriceNett:         lib.RoundFloat(sumPriceNett, 2),
			TaxAmount:         lib.RoundFloat(sumPriceTaxAmount, 2),
			PriceGrossMonthly: lib.RoundFloat(sumPriceGrossMonthly, 2),
			PriceNettMonthly:  lib.RoundFloat(sumPriceNettMonthly, 2),
			TaxAmountMonthly:  lib.RoundFloat(sumPriceTaxAmountMonthly, 2),
			Payment:           models.ManualPaymentProvider,
			FundsOrigin:       "Proprie risorse economiche",
			ProducerCode:      networkNode.Code,
			ProducerUid:       networkNode.Uid,
			ProducerType:      networkNode.Type,
			Contractors:       contractors,
			Assets: []models.Asset{{
				Guarantees: guarantees,
			}},
			OffersPrices: map[string]map[string]*models.Price{
				"default": {
					string(models.PaySplitMonthly): &models.Price{
						Net:   lib.RoundFloat(sumPriceNettMonthly, 2),
						Tax:   lib.RoundFloat(sumPriceTaxAmountMonthly, 2),
						Gross: lib.RoundFloat(sumPriceGrossMonthly, 2),
					},
					string(models.PaySplitYearly): &models.Price{
						Net:   lib.RoundFloat(sumPriceNett, 2),
						Tax:   lib.RoundFloat(sumPriceTaxAmount, 2),
						Gross: lib.RoundFloat(sumPriceGross, 2),
					},
				},
			},
		}

		// get monthly prices

		calculateMonthlyPrices(&policy)

		// enrich contractor data if guarantee death present

		if policy.HasGuarantee("death") {

			// TODO: need to split case individual contractor and legalEntity contractor

			if !isLegalEntity {
				// setting contractor identity document

				rawCode, _ := strconv.Atoi(strings.TrimSpace(row[56]))
				identityDocumentCode := fmt.Sprintf("%02d", rawCode)
				contractor.IdentityDocuments = []*models.IdentityDocument{{
					Number:           strings.TrimSpace(strings.ToUpper(row[57])),
					Code:             identityDocumentCode,
					Type:             identityDocumentMap[identityDocumentCode],
					DateOfIssue:      ParseDateDDMMYYYY(row[58]),
					ExpiryDate:       ParseDateDDMMYYYY(row[58]).AddDate(10, 0, 0),
					IssuingAuthority: strings.TrimSpace(lib.Capitalize(row[59])),
					PlaceOfIssue:     strings.TrimSpace(lib.Capitalize(row[59])),
					LastUpdate:       policy.EmitDate,
				}}

				// setting contractor domicile

				contractor.Domicile = &models.Address{
					StreetName: strings.TrimSpace(lib.Capitalize(row[46])),
					City:       strings.TrimSpace(lib.Capitalize(row[48])),
					CityCode:   strings.TrimSpace(strings.ToUpper(row[49])),
					PostalCode: row[47],
					Locality:   strings.TrimSpace(lib.Capitalize(row[48])),
				}

				if !contractorEqualInsured {
					// setting insured identity documents

					rawCode, _ = strconv.Atoi(strings.TrimSpace(row[76]))
					identityDocumentCode = fmt.Sprintf("%02d", rawCode)
					insured.IdentityDocuments = []*models.IdentityDocument{{
						Number:           strings.TrimSpace(strings.ToUpper(row[77])),
						Code:             identityDocumentCode,
						Type:             identityDocumentMap[identityDocumentCode],
						DateOfIssue:      ParseDateDDMMYYYY(row[78]),
						ExpiryDate:       ParseDateDDMMYYYY(row[78]).AddDate(10, 0, 0),
						IssuingAuthority: strings.TrimSpace(lib.Capitalize(row[79])),
						PlaceOfIssue:     strings.TrimSpace(lib.Capitalize(row[79])),
						LastUpdate:       policy.EmitDate,
					}}

					// setting insured domicile

					insured.Domicile = &models.Address{
						StreetName: strings.TrimSpace(lib.Capitalize(row[67])),
						City:       strings.TrimSpace(lib.Capitalize(row[69])),
						CityCode:   strings.TrimSpace(strings.ToUpper(row[70])),
						PostalCode: row[68],
						Locality:   strings.TrimSpace(lib.Capitalize(row[69])),
					}
				} else {
					insured.IdentityDocuments = contractor.IdentityDocuments
					insured.Domicile = contractor.Domicile
				}
			} else {
				// setting insured identity documents

				rawCode, _ := strconv.Atoi(strings.TrimSpace(row[76]))
				identityDocumentCode := fmt.Sprintf("%02d", rawCode)
				insured.IdentityDocuments = []*models.IdentityDocument{{
					Number:           strings.TrimSpace(strings.ToUpper(row[77])),
					Code:             identityDocumentCode,
					Type:             identityDocumentMap[identityDocumentCode],
					DateOfIssue:      ParseDateDDMMYYYY(row[78]),
					ExpiryDate:       ParseDateDDMMYYYY(row[78]).AddDate(10, 0, 0),
					IssuingAuthority: strings.TrimSpace(lib.Capitalize(row[79])),
					PlaceOfIssue:     strings.TrimSpace(lib.Capitalize(row[79])),
					LastUpdate:       policy.EmitDate,
				}}

				// setting insured domicile

				insured.Domicile = &models.Address{
					StreetName: strings.TrimSpace(lib.Capitalize(row[67])),
					City:       strings.TrimSpace(lib.Capitalize(row[69])),
					CityCode:   strings.TrimSpace(strings.ToUpper(row[70])),
					PostalCode: row[68],
					Locality:   strings.TrimSpace(lib.Capitalize(row[69])),
				}
			}
		}

		policy.Assets[0].Person = insured
		policy.Contractor = *contractor
		policy.Contractor.Uid = lib.NewDoc(models.UserCollection)

		// check if user is already present

		/*query := fmt.Sprintf(
			"SELECT * FROM `%s.%s` WHERE isDelete = false AND JSON_VALUE(data, '$.contractor.fiscalCode') = '%s'",
			models.WoptaDataset,
			models.PoliciesViewCollection,
			insured.FiscalCode,
		)
		retrievedPolicies, err := lib.QueryRowsBigQuery[models.Policy](query)
		if err != nil {
			log.Printf("error retrieving policies bigquery: %s", err.Error())
			continue
		}
		for _, rp := range retrievedPolicies {
			if rp.Name == models.LifeProduct {
				log.Printf("error user already has a life policy")
				return "", nil, nil
			}
		}

		if len(retrievedPolicies) > 0 {
			policy.Contractor.Uid = retrievedPolicies[0].Contractor.Uid
		} else {
			policy.Contractor.Uid = lib.NewDoc(models.UserCollection)
		}*/

		transactionsOutput := make(map[string]TransactionsOutput, 0)

		// create transactions and network node transactions

		scheduleDate := policy.StartDate
		transactionPayDate := policy.StartDate

		// if monthly create remaining transactions and network transactions if transaction is paid

		if monthlyPolicies[policy.CodeCompany] != nil {
			tr := createTransaction(policy, mgaProducts[policy.ProductVersion], "", scheduleDate, transactionPayDate, policy.PriceGrossMonthly, policy.PriceNettMonthly, true)

			transactionsOutput = map[string]TransactionsOutput{
				scheduleDate.Format(models.TimeDateOnly): {
					Transaction:         tr,
					NetworkTransactions: createNetworkTransactions(&policy, &tr, networkNode, mgaProducts[policy.ProductVersion]),
				},
			}

			for i := 1; i < 12; i++ {
				transactionPayDate = time.Time{}
				scheduleDate = scheduleDate.AddDate(0, 1, 0)
				isPay := false
				payDateString := scheduleDate.Format("02012006")
				if monthlyPolicies[policy.CodeCompany][payDateString] != nil {
					isPay = true
					transactionPayDate = scheduleDate
				}
				tr = createTransaction(policy, mgaProducts[policy.ProductVersion], "", scheduleDate, transactionPayDate, policy.PriceGrossMonthly, policy.PriceNettMonthly, isPay)
				sc := scheduleDate.Format(models.TimeDateOnly)

				if isPay {
					transactionsOutput[sc] = TransactionsOutput{
						Transaction:         tr,
						NetworkTransactions: createNetworkTransactions(&policy, &tr, networkNode, mgaProducts[policy.ProductVersion]),
					}
				} else {
					transactionsOutput[sc] = TransactionsOutput{
						Transaction:         tr,
						NetworkTransactions: []*models.NetworkTransaction{},
					}
				}
			}
		} else {
			tr := createTransaction(policy, mgaProducts[policy.ProductVersion], "", scheduleDate, transactionPayDate, policy.PriceGross, policy.PriceNett, true)

			transactionsOutput = map[string]TransactionsOutput{
				scheduleDate.Format(models.TimeDateOnly): {
					Transaction:         tr,
					NetworkTransactions: createNetworkTransactions(&policy, &tr, networkNode, mgaProducts[policy.ProductVersion]),
				},
			}
		}

		result[codeCompany] = ResultStruct{
			Policy:       policy,
			Transactions: transactionsOutput,
		}

		// update node portfolio

		networkNode.Policies = append(networkNode.Policies, policy.Uid)
		networkNode.Users = append(networkNode.Users, policy.Contractor.Uid)

		dryRun := true
		if req.DryRun != nil {
			dryRun = *req.DryRun
		}
		log.Printf("dryRun: %v", dryRun)
		if !dryRun {
			collectionPrefix := req.CollectionPrefix

			// save policy firestore

			err := lib.SetFirestoreErr(fmt.Sprintf("%s%s", collectionPrefix, models.PolicyCollection), policy.Uid, policy)
			if err != nil {
				log.Printf("error saving policy firestore: %s", err.Error())
				continue
			}

			// save policy bigquery

			policyBigquerySave(policy, collectionPrefix)

			// save transactions firestore

			for _, res := range transactionsOutput {
				err := lib.SetFirestoreErr(fmt.Sprintf("%s%s", collectionPrefix, models.TransactionsCollection), res.Transaction.Uid, res.Transaction)
				if err != nil {
					log.Printf("error saving transaction firestore: %s", err.Error())
					continue
				}

				// save transactions bigquery

				transactionBigQuerySave(res.Transaction, collectionPrefix)

				for _, nt := range res.NetworkTransactions {
					// save network transactions bigquery
					networkTransactionBigQuerySave(*nt, collectionPrefix)
				}
			}

			// save user firestore

			err = lib.SetFirestoreErr(fmt.Sprintf("%s%s", collectionPrefix, models.UserCollection), policy.Contractor.Uid, policy.Contractor)
			if err != nil {
				log.Printf("error saving contractor firestore: %s", err.Error())
				continue
			}

			// save user bigquery

			contractorBigQuerySave(&policy.Contractor, collectionPrefix)

			// save network node firestore

			err = lib.SetFirestoreErr(fmt.Sprintf("%s%s", collectionPrefix, models.NetworkNodesCollection), networkNode.Uid, networkNode)
			if err != nil {
				log.Printf("error saving network node firestore: %s", err.Error())
				continue
			}

			// save network node bigquery

			networkNodeBigQuerySave(*networkNode, collectionPrefix)

			// save single guarantees into bigquery
			models.SetGuaranteBigquery(policy, "emit", fmt.Sprintf("%s%s", collectionPrefix, models.GuaranteeCollection))
		}

		//log.Println("LifeInFx policy:", policy)
		b, e := json.Marshal(policy)
		log.Println("LifeInFx policy:", e)
		log.Println("LifeInFx policy:", string(b))
		policies = append(policies, policy)
		// docref, _, _ := lib.PutFirestoreErr("test-policy", policy)
		// log.Println("LifeInFx doc id: ", docref.ID)

		//_, e = models.UpdateUserByFiscalCode("uat", policy.Contractor)
		//log.Println("LifeInFx policy:", policy)
		//tr := transaction.PutByPolicy(policy, "", "uat", "", "", sumPriseGross, 0, "", "manual", true)
		//	log.Println("LifeInFx transactionpolicy:",tr)
		//accounting.CreateNetworkTransaction(tr, "uat")

	}

	log.Printf("Skipped %d policies: %v\n", len(skippedPolicies), skippedPolicies)
	log.Printf("Missing %d producers: %v\n", len(missingProducers), missingProducers)
	log.Printf("Wrong contractor fiscal code %d policies: %v\n", len(wrongContractorFiscalCodePolicies), wrongContractorFiscalCodePolicies)
	log.Printf("Wrong insured fiscal code %d policies: %v\n", len(wrongInsuredFiscalCodePolicies), wrongInsuredFiscalCodePolicies)
	log.Printf("Missing contractor Birth City %d policies: %v\n", len(missingContractorBirthCityPolicies), missingContractorBirthCityPolicies)
	log.Printf("Missing insured Birth City %d policies: %v\n", len(missingInsuredBirthCityPolicies), missingInsuredBirthCityPolicies)
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

	endDateJob = time.Now().UTC()

	log.Printf("Script started at %s", startDateJob.String())
	log.Printf("Script ended at %s", endDateJob.String())

	return "", nil, e
}

func parsingTitolareEffettivo(row []string, offset int, i int) models.User {
	phone := strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(row[132], " ", ""), " ", ""))
	if phone != "" {
		phone = fmt.Sprintf("+39%s", phone)
	}
	rawDocumentCode, _ := strconv.Atoi(strings.TrimSpace(row[136+(offset*i)]))
	identityDocumentCode := fmt.Sprintf("%02d", rawDocumentCode)
	titolareEffettivo := models.User{
		Type:          models.UserLegalEntity,
		Name:          strings.TrimSpace(lib.Capitalize(row[118+(offset*i)])),
		Surname:       strings.TrimSpace(lib.Capitalize(row[117+(offset*i)])),
		FiscalCode:    strings.TrimSpace(strings.ToUpper(row[121+(offset*i)])),
		Gender:        strings.TrimSpace(strings.ToUpper(row[119+(offset*i)])),
		BirthDate:     ParseDateDDMMYYYY(row[120+(offset*i)]).Format(time.RFC3339),
		Mail:          strings.TrimSpace(strings.ToLower(row[131+(offset*i)])),
		Phone:         phone,
		BirthCity:     strings.TrimSpace(lib.Capitalize(row[133+(offset*i)])),
		BirthProvince: strings.TrimSpace(strings.ToUpper(row[134+(offset*i)])),
		Residence: &models.Address{
			StreetName: strings.TrimSpace(lib.Capitalize(row[122+(offset*i)])),
			City:       strings.TrimSpace(lib.Capitalize(row[124+(offset*i)])),
			CityCode:   strings.TrimSpace(strings.ToUpper(row[125+(offset*i)])),
			PostalCode: row[123+(offset*i)],
			Locality:   strings.TrimSpace(lib.Capitalize(row[124+(offset*i)])),
		},
		Domicile: &models.Address{
			StreetName: strings.TrimSpace(lib.Capitalize(row[126+(offset*i)])),
			City:       strings.TrimSpace(lib.Capitalize(row[128+(offset*i)])),
			CityCode:   strings.TrimSpace(strings.ToUpper(row[129+(offset*i)])),
			PostalCode: row[127+(offset*i)],
			Locality:   strings.TrimSpace(lib.Capitalize(row[128+(offset*i)])),
		},
		IdentityDocuments: []*models.IdentityDocument{{
			Number:           strings.TrimSpace(strings.ToUpper(row[137+(offset*i)])),
			Code:             identityDocumentCode,
			Type:             identityDocumentMap[identityDocumentCode],
			DateOfIssue:      ParseDateDDMMYYYY(row[138+(offset*i)]),
			ExpiryDate:       ParseDateDDMMYYYY(row[138+(offset*i)]).AddDate(10, 0, 0),
			IssuingAuthority: strings.TrimSpace(lib.Capitalize(row[139+(offset*i)])),
			PlaceOfIssue:     strings.TrimSpace(lib.Capitalize(row[139+(offset*i)])),
		}},
		Work:            strings.TrimSpace(lib.Capitalize(row[130+(offset*i)])),
		LegalEntityType: models.TitolareEffettivo,
	}
	return titolareEffettivo
}

func parseEsecutore(row []string) models.User {
	rawDocumentCode, _ := strconv.Atoi(strings.TrimSpace(row[238]))
	identityDocumentCode := fmt.Sprintf("%02d", rawDocumentCode)

	phone := strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(row[234], " ", ""), " ", ""))
	if phone != "" {
		phone = fmt.Sprintf("+39%s", phone)
	}

	esecutore := models.User{
		Type:          models.UserLegalEntity,
		Name:          strings.TrimSpace(lib.Capitalize(row[221])),
		Surname:       strings.TrimSpace(lib.Capitalize(row[220])),
		FiscalCode:    strings.TrimSpace(strings.ToUpper(row[224])),
		Gender:        strings.TrimSpace(strings.ToUpper(row[222])),
		Mail:          strings.TrimSpace(strings.ToLower(row[233])),
		Phone:         phone,
		BirthDate:     ParseDateDDMMYYYY(row[223]).Format(time.RFC3339),
		BirthCity:     strings.TrimSpace(lib.Capitalize(row[235])),
		BirthProvince: strings.TrimSpace(strings.ToUpper(row[236])),
		Residence: &models.Address{
			StreetName: strings.TrimSpace(lib.Capitalize(row[225])),
			City:       strings.TrimSpace(lib.Capitalize(row[227])),
			CityCode:   strings.TrimSpace(strings.ToUpper(row[228])),
			PostalCode: row[226],
			Locality:   strings.TrimSpace(lib.Capitalize(row[227])),
		},
		Domicile: &models.Address{
			StreetName: strings.TrimSpace(lib.Capitalize(row[229])),
			City:       strings.TrimSpace(lib.Capitalize(row[231])),
			CityCode:   strings.TrimSpace(strings.ToUpper(row[232])),
			PostalCode: row[230],
			Locality:   strings.TrimSpace(lib.Capitalize(row[231])),
		},
		IdentityDocuments: []*models.IdentityDocument{{
			Number:           strings.TrimSpace(strings.ToUpper(row[239])),
			Code:             identityDocumentCode,
			Type:             identityDocumentMap[identityDocumentCode],
			DateOfIssue:      ParseDateDDMMYYYY(row[240]),
			ExpiryDate:       ParseDateDDMMYYYY(row[240]).AddDate(10, 0, 0),
			IssuingAuthority: strings.TrimSpace(lib.Capitalize(row[241])),
			PlaceOfIssue:     strings.TrimSpace(lib.Capitalize(row[241])),
		}},
		CompanyRole:     strings.TrimSpace(lib.Capitalize(row[244])),
		LegalEntityType: models.Esecutore,
	}
	return esecutore
}

func parseIndividualContractor(codeCompany string, row []string, codes map[string]map[string]string) *models.Contractor {
	phone := strings.TrimSpace(strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(row[33], " ", ""), " ", "")))
	if phone != "" {
		phone = fmt.Sprintf("+39%s", phone)
	}

	contractor := &models.Contractor{
		Type:          models.UserIndividual,
		Name:          strings.TrimSpace(lib.Capitalize(row[24])),
		Surname:       strings.TrimSpace(lib.Capitalize(row[23])),
		FiscalCode:    strings.TrimSpace(strings.ToUpper(row[27])),
		Gender:        strings.TrimSpace(strings.ToUpper(row[25])),
		BirthDate:     ParseDateDDMMYYYY(row[26]).Format(time.RFC3339),
		Mail:          strings.TrimSpace(strings.ToLower(row[32])),
		Phone:         phone,
		BirthCity:     strings.TrimSpace(lib.Capitalize(row[50])),
		BirthProvince: strings.TrimSpace(strings.ToUpper(row[51])),
		Residence: &models.Address{
			StreetName: strings.TrimSpace(lib.Capitalize(row[28])),
			City:       strings.TrimSpace(lib.Capitalize(row[30])),
			CityCode:   strings.TrimSpace(strings.ToUpper(row[31])),
			PostalCode: row[29],
			Locality:   strings.TrimSpace(lib.Capitalize(row[30])),
		},
		CreationDate: ParseDateDDMMYYYY(row[4]),
		UpdatedDate:  time.Now().UTC(),
		Consens: &[]models.Consens{
			{
				Title:        "Privacy",
				Consens:      "Il sottoscritto, letta e compresa l'informativa sul trattamento dei dati personali, ACCONSENTE al trattamento dei propri dati personali da parte di Wopta Assicurazioni per l'invio di comunicazioni e proposte commerciali e di marketing, incluso l'invio di newsletter e ricerche di mercato, attraverso strumenti automatizzati (sms, mms, e-mail, ecc.) e non (posta cartacea e telefono con operatore).",
				Key:          2,
				Answer:       false,
				CreationDate: ParseDateDDMMYYYY(row[4]),
			},
		},
	}

	var usr models.User
	_, usr, err := user.CalculateFiscalCode(*contractor.ToUser())
	if err != nil {
		if strings.ToLower(err.Error()) == "invalid birth city" {
			_, extractedUser, _ := ExtractUserDataFromFiscalCode(contractor.FiscalCode, codes)

			contractor.BirthCity = extractedUser.BirthCity
			contractor.BirthProvince = extractedUser.BirthProvince

			_, usr, err = user.CalculateFiscalCode(*contractor.ToUser())

			missingContractorBirthCityPolicies = append(missingContractorBirthCityPolicies, codeCompany)
		} else {
			log.Printf("error: %s", err.Error())
			skippedPolicies = append(skippedPolicies, codeCompany)
			return nil
		}
	}

	if strings.ToUpper(usr.FiscalCode) != strings.ToUpper(contractor.FiscalCode) {
		contractor.FiscalCode = usr.FiscalCode
		wrongContractorFiscalCodePolicies = append(wrongContractorFiscalCodePolicies, codeCompany)
	}

	return contractor
}

func parseEnterpriseContractor(row []string) *models.Contractor {
	phone := strings.TrimSpace(strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(row[33], " ", ""), " ", "")))
	if phone != "" {
		phone = fmt.Sprintf("+39%s", phone)
	}

	contractor := &models.Contractor{
		Type:         models.UserLegalEntity,
		Name:         strings.TrimSpace(lib.Capitalize(row[23])),
		VatCode:      fmt.Sprintf("%011s", strings.TrimSpace(row[27])),
		Mail:         strings.TrimSpace(strings.ToLower(row[32])),
		Phone:        phone,
		CreationDate: ParseDateDDMMYYYY(row[4]),
		UpdatedDate:  time.Now().UTC(),
		Consens: &[]models.Consens{
			{
				Title:        "Privacy",
				Consens:      "Il sottoscritto, letta e compresa l'informativa sul trattamento dei dati personali, ACCONSENTE al trattamento dei propri dati personali da parte di Wopta Assicurazioni per l'invio di comunicazioni e proposte commerciali e di marketing, incluso l'invio di newsletter e ricerche di mercato, attraverso strumenti automatizzati (sms, mms, e-mail, ecc.) e non (posta cartacea e telefono con operatore).",
				Key:          2,
				Answer:       false,
				CreationDate: ParseDateDDMMYYYY(row[4]),
			},
		},
		CompanyAddress: &models.Address{
			StreetName: strings.TrimSpace(lib.Capitalize(row[28])),
			City:       strings.TrimSpace(lib.Capitalize(row[30])),
			CityCode:   strings.TrimSpace(strings.ToUpper(row[31])),
			PostalCode: row[29],
			Locality:   strings.TrimSpace(lib.Capitalize(row[30])),
		},
	}

	return contractor
}

func parseInsured(codeCompany string, row []string, codes map[string]map[string]string) *models.User {
	insured := &models.User{
		Type:          models.UserIndividual,
		Name:          strings.TrimSpace(lib.Capitalize(row[35])),
		Surname:       strings.TrimSpace(lib.Capitalize(row[34])),
		FiscalCode:    strings.TrimSpace(strings.ToUpper(row[38])),
		Gender:        strings.TrimSpace(strings.ToUpper(row[36])),
		BirthDate:     ParseDateDDMMYYYY(row[37]).Format(time.RFC3339),
		Mail:          strings.TrimSpace(strings.ToLower(row[71])),
		Phone:         fmt.Sprintf("+39%s", strings.TrimSpace(strings.ReplaceAll(row[72], " ", ""))),
		BirthCity:     strings.TrimSpace(lib.Capitalize(row[73])),
		BirthProvince: strings.TrimSpace(strings.ToUpper(row[74])),
		Residence: &models.Address{
			StreetName: strings.TrimSpace(lib.Capitalize(row[63])),
			City:       strings.TrimSpace(lib.Capitalize(row[65])),
			CityCode:   strings.TrimSpace(strings.ToUpper(row[66])),
			PostalCode: row[64],
			Locality:   strings.TrimSpace(lib.Capitalize(row[65])),
		},
		CreationDate: ParseDateDDMMYYYY(row[4]),
		UpdatedDate:  time.Now().UTC(),
		Consens: &[]models.Consens{
			{
				Title:        "Privacy",
				Consens:      "Il sottoscritto, letta e compresa l'informativa sul trattamento dei dati personali, ACCONSENTE al trattamento dei propri dati personali da parte di Wopta Assicurazioni per l'invio di comunicazioni e proposte commerciali e di marketing, incluso l'invio di newsletter e ricerche di mercato, attraverso strumenti automatizzati (sms, mms, e-mail, ecc.) e non (posta cartacea e telefono con operatore).",
				Key:          2,
				Answer:       false,
				CreationDate: ParseDateDDMMYYYY(row[4]),
			},
		},
	}

	// check fiscalcode

	_, usr, err := user.CalculateFiscalCode(*insured)
	if err != nil {
		if strings.ToLower(err.Error()) == "invalid birth city" {
			_, extractedUser, _ := ExtractUserDataFromFiscalCode(insured.FiscalCode, codes)

			insured.BirthCity = extractedUser.BirthCity
			insured.BirthProvince = extractedUser.BirthProvince

			_, usr, err = user.CalculateFiscalCode(*insured)

			missingInsuredBirthCityPolicies = append(missingInsuredBirthCityPolicies, codeCompany)
		} else {
			log.Printf("error: %s", err.Error())
			skippedPolicies = append(skippedPolicies, codeCompany)
			return nil
		}
	}

	if strings.ToUpper(usr.FiscalCode) != strings.ToUpper(insured.FiscalCode) {
		insured.FiscalCode = usr.FiscalCode
		wrongInsuredFiscalCodePolicies = append(wrongInsuredFiscalCodePolicies, codeCompany)
		//skippedPolicies = append(skippedPolicies, policy.CodeCompany)
		//continue
	}
	return insured
}

func calculateMonthlyPrices(policy *models.Policy) {
	if policy.PaymentSplit == string(models.PaySplitMonthly) {
		policy.PriceGross = lib.RoundFloat(policy.PriceGross*12, 2)
		policy.PriceNett = lib.RoundFloat(policy.PriceNett*12, 2)
		policy.TaxAmount = lib.RoundFloat(policy.TaxAmount*12, 2)

		for index, guarantee := range policy.Assets[0].Guarantees {
			policy.Assets[0].Guarantees[index].Value.PremiumGrossYearly = lib.RoundFloat(guarantee.Value.PremiumGrossYearly*12, 2)
			policy.Assets[0].Guarantees[index].Value.PremiumNetYearly = lib.RoundFloat(guarantee.Value.PremiumNetYearly*12, 2)
			policy.Assets[0].Guarantees[index].Value.PremiumTaxAmountYearly = lib.RoundFloat(policy.Assets[0].Guarantees[index].Value.PremiumGrossYearly-policy.Assets[0].Guarantees[index].Value.PremiumNetYearly, 2)
		}

		policy.OffersPrices["default"][string(models.PaySplitYearly)].Gross = policy.PriceGross
		policy.OffersPrices["default"][string(models.PaySplitYearly)].Net = policy.PriceNett
		policy.OffersPrices["default"][string(models.PaySplitYearly)].Tax = policy.TaxAmount
	} else {
		policy.PriceGrossMonthly = lib.RoundFloat(policy.PriceGrossMonthly/12, 2)
		policy.PriceNettMonthly = lib.RoundFloat(policy.PriceNettMonthly/12, 2)
		policy.TaxAmountMonthly = lib.RoundFloat(policy.TaxAmountMonthly/12, 2)

		for index, guarantee := range policy.Assets[0].Guarantees {
			policy.Assets[0].Guarantees[index].Value.PremiumGrossMonthly = lib.RoundFloat(guarantee.Value.PremiumGrossMonthly/12, 2)
			policy.Assets[0].Guarantees[index].Value.PremiumNetMonthly = lib.RoundFloat(guarantee.Value.PremiumNetMonthly/12, 2)
			policy.Assets[0].Guarantees[index].Value.PremiumTaxAmountMonthly = lib.RoundFloat(policy.Assets[0].Guarantees[index].Value.PremiumGrossMonthly-policy.Assets[0].Guarantees[index].Value.PremiumNetMonthly, 2)
		}

		policy.OffersPrices["default"][string(models.PaySplitMonthly)].Gross = policy.PriceGrossMonthly
		policy.OffersPrices["default"][string(models.PaySplitMonthly)].Net = policy.PriceNettMonthly
		policy.OffersPrices["default"][string(models.PaySplitMonthly)].Tax = policy.TaxAmountMonthly
	}
}

var identityDocumentMap map[string]string = map[string]string{
	"01": "Carta di Identità",
	"02": "Patente di Guida",
	"03": "Passaporto",
}

func LifeMapCodecCompanyAxaRevert(g string) (string, string, string, string) {
	log.Println("LifeInFx LifeMapCodecCompanyAxaRevert:", g)
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
	log.Println("LifeInFx LifeMapCodecCompanyAxaRevert:", version)
	log.Println("LifeInFx LifeMapCodecCompanyAxaRevert:", code)
	return result, slug, version, pay
}

func ParseDateDDMMYYYY(date string) time.Time {
	var (
		res time.Time
	)
	log.Println("LifeInFx ParseDateDDMMYYYY date:", date)
	log.Println("LifeInFx ParseDateDDMMYYYY len(date):", len(date))
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
		log.Println("LifeInFx ParseDateDDMMYYYY d:", d)
		log.Println("LifeInFx ParseDateDDMMYYYY m:", m)
		log.Println("LifeInFx ParseDateDDMMYYYY y:", y)
		log.Println("LifeInFx ParseDateDDMMYYYY res:", res)
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
			IsLegitimateSuccessors: true,
			BeneficiaryType:        models.BeneficiaryLegalAndWillSuccessors,
		}
	} else if r[82] == "NM" {
		if strings.TrimSpace(strings.ToUpper(r[85+rangeCell])) == "" || strings.TrimSpace(strings.ToUpper(r[85+rangeCell])) == "0" {
			return nil
		}

		isFamilyMember := false
		isContactable := true
		if strings.TrimSpace(strings.ToUpper(r[92+rangeCell])) == "SI" {
			isFamilyMember = true
		}

		if strings.TrimSpace(strings.ToUpper(r[93+rangeCell])) == "SI" {
			isContactable = false
		}

		benef = &models.Beneficiary{
			BeneficiaryType: models.BeneficiaryChosenBeneficiary,
			Name:            strings.TrimSpace(lib.Capitalize(r[84+rangeCell])),
			Surname:         strings.TrimSpace(lib.Capitalize(r[83+rangeCell])),
			FiscalCode:      strings.TrimSpace(strings.ToUpper(r[85+rangeCell])),
			Mail:            strings.TrimSpace(strings.ToLower(r[91+rangeCell])),
			Phone:           strings.TrimSpace(strings.ReplaceAll(r[86+rangeCell], " ", "")),
			Residence: &models.Address{
				StreetName: strings.TrimSpace(lib.Capitalize(r[87+rangeCell])),
				City:       strings.TrimSpace(lib.Capitalize(r[88+rangeCell])),
				CityCode:   strings.TrimSpace(strings.ToUpper(r[90+rangeCell])),
				PostalCode: strings.TrimSpace(r[89+rangeCell]),
				Locality:   strings.TrimSpace(lib.Capitalize(r[88+rangeCell])),
			},
			IsContactable:  isContactable,
			IsFamilyMember: isFamilyMember,
		}
	} else if r[82] == "PG" {
		benef = &models.Beneficiary{
			BeneficiaryType: models.BeneficiarySelfLegalEntity,
			Name:            strings.TrimSpace(lib.Capitalize(r[23])),
			VatCode:         fmt.Sprintf("%011s", strings.TrimSpace(r[27])),
			Mail:            strings.TrimSpace(strings.ToLower(r[32])),
			Phone:           fmt.Sprintf("+39%s", strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(r[33], " ", ""), " ", ""))),
			CompanyAddress: &models.Address{
				StreetName: strings.TrimSpace(lib.Capitalize(r[28])),
				City:       strings.TrimSpace(lib.Capitalize(r[30])),
				CityCode:   strings.TrimSpace(strings.ToUpper(r[31])),
				PostalCode: r[29],
				Locality:   strings.TrimSpace(lib.Capitalize(r[30])),
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

func createTransaction(policy models.Policy, mgaProduct *models.Product, customerId string, scheduleDate time.Time, payDate time.Time, priceGross, priceNett float64, isPay bool) models.Transaction {
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
		AmountNet:       priceNett,
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
		IsEmit:          isPay,
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

func policyBigquerySave(policy models.Policy, collectionPrefix string) {
	log.Printf("[policyBigquerySave] parsing data for policy %s", policy.Uid)

	policyBig := lib.GetDatasetByEnv("", fmt.Sprintf("%s%s", collectionPrefix, models.PolicyCollection))
	policyJson, err := policy.Marshal()
	if err != nil {
		log.Printf("[policy.BigquerySave] error marshaling policy: %s", err.Error())
	}

	policy.Data = string(policyJson)
	policy.BigStartDate = civil.DateTimeOf(policy.StartDate)
	policy.BigRenewDate = civil.DateTimeOf(policy.RenewDate)
	policy.BigEndDate = civil.DateTimeOf(policy.EndDate)
	policy.BigEmitDate = civil.DateTimeOf(policy.EmitDate)
	policy.BigStatusHistory = strings.Join(policy.StatusHistory, ",")
	if policy.ReservedInfo != nil {
		policy.BigReasons = strings.Join(policy.ReservedInfo.Reasons, ",")
		policy.BigAcceptanceNote = policy.ReservedInfo.AcceptanceNote
		policy.BigAcceptanceDate = lib.GetBigQueryNullDateTime(policy.ReservedInfo.AcceptanceDate)
	}

	log.Println("[policyBigquerySave] saving to bigquery...")
	err = lib.InsertRowsBigQuery(models.WoptaDataset, policyBig, policy)
	if err != nil {
		log.Println("[policyBigquerySave] error saving policy to bigquery: ", err.Error())
		return
	}
	log.Println("[policyBigquerySave] bigquery saved!")
}

func transactionBigQuerySave(transaction models.Transaction, collectionPrefix string) {
	fireTransactions := lib.GetDatasetByEnv("", fmt.Sprintf("%s%s", collectionPrefix, models.TransactionsCollection))

	transaction.BigPayDate = lib.GetBigQueryNullDateTime(transaction.PayDate)
	transaction.BigTransactionDate = lib.GetBigQueryNullDateTime(transaction.TransactionDate)
	transaction.BigCreationDate = civil.DateTimeOf(transaction.CreationDate)
	transaction.BigStatusHistory = strings.Join(transaction.StatusHistory, ",")
	transaction.BigUpdateDate = lib.GetBigQueryNullDateTime(transaction.UpdateDate)
	log.Println("Transaction save BigQuery: " + transaction.Uid)

	err := lib.InsertRowsBigQuery(models.WoptaDataset, fireTransactions, transaction)
	if err != nil {
		log.Println("ERROR Transaction "+transaction.Uid+" save BigQuery: ", err)
		return
	}
	log.Println("Transaction BigQuery saved!")
}

func networkTransactionBigQuerySave(nt models.NetworkTransaction, collectionPrefix string) error {
	log.Println("[NetworkTransaction.SaveBigQuery]")

	var (
		err       error
		datasetId = models.WoptaDataset
		tableId   = fmt.Sprintf("%s%s", collectionPrefix, models.NetworkTransactionCollection)
	)

	baseQuery := fmt.Sprintf("SELECT * FROM `%s.%s` WHERE ", datasetId, tableId)
	whereClause := fmt.Sprintf("uid = '%s'", nt.Uid)
	query := fmt.Sprintf("%s %s", baseQuery, whereClause)

	result, err := lib.QueryRowsBigQuery[models.NetworkTransaction](query)
	if err != nil {
		log.Printf("[NetworkTransaction.SaveBigQuery] error querying db with query %s: %s", query, err.Error())
		return err
	}

	if len(result) == 0 {
		log.Printf("[NetworkTransaction.SaveBigQuery] creating new NetworkTransaction %s", nt.Uid)
		err = lib.InsertRowsBigQuery(datasetId, tableId, nt)
	} else {
		log.Printf("[NetworkTransaction.SaveBigQuery] updating NetworkTransaction %s", nt.Uid)
		updatedFields := make(map[string]interface{})
		updatedFields["status"] = nt.Status
		updatedFields["statusHistory"] = nt.StatusHistory
		updatedFields["isPay"] = nt.IsPay
		updatedFields["isConfirmed"] = nt.IsConfirmed
		updatedFields["isDelete"] = nt.IsDelete
		if nt.PayDate.Valid {
			updatedFields["payDate"] = nt.PayDate
		}
		if nt.TransactionDate.Valid {
			updatedFields["transactionDate"] = nt.TransactionDate
		}
		if nt.ConfirmationDate.Valid {
			updatedFields["confirmationDate"] = nt.ConfirmationDate
		}
		if nt.DeletionDate.Valid {
			updatedFields["deletionDate"] = nt.DeletionDate
		}

		err = lib.UpdateRowBigQueryV2(datasetId, tableId, updatedFields, "WHERE "+whereClause)
	}

	if err != nil {
		log.Printf("[NetworkTransaction.SaveBigQuery] error saving to db: %s", err.Error())
		return err
	}

	log.Println("[NetworkTransaction.SaveBigQuery] NetworkTransaction saved!")
	return nil
}

func networkNodeBigQuerySave(nn models.NetworkNode, collectionPrefix string) error {
	log.Println("[networkNodeSaveBigQuery]")

	nnJson, _ := json.Marshal(nn)

	nn.Data = string(nnJson)
	nn.BigCreationDate = lib.GetBigQueryNullDateTime(nn.CreationDate)
	nn.BigUpdatedDate = lib.GetBigQueryNullDateTime(nn.UpdatedDate)
	nn.Agent = parseBigQueryAgentNode(nn.Agent)
	nn.AreaManager = parseBigQueryAgentNode(nn.AreaManager)
	nn.Agency = parseBigQueryAgencyNode(nn.Agency)
	nn.Broker = parseBigQueryAgencyNode(nn.Broker)

	for _, p := range nn.Products {
		companies := make([]models.NodeCompany, 0)
		for _, c := range p.Companies {
			companies = append(companies, models.NodeCompany{
				Name:         c.Name,
				ProducerCode: c.ProducerCode,
			})
		}
		nn.BigProducts = append(nn.BigProducts, models.NodeProduct{
			Name:      p.Name,
			Companies: companies,
		})
	}

	err := lib.InsertRowsBigQuery(models.WoptaDataset, fmt.Sprintf("%s%s", collectionPrefix, models.NetworkNodesCollection), nn)
	return err
}

func parseBigQueryAgentNode(agent *models.AgentNode) *models.AgentNode {
	if agent == nil {
		return nil
	}

	if agent.BirthDate != "" {
		birthDate, _ := time.Parse(time.RFC3339, agent.BirthDate)
		agent.BigBirthDate = lib.GetBigQueryNullDateTime(birthDate)
	}
	if agent.Residence != nil {
		agent.Residence.BigLocation = lib.GetBigQueryNullGeography(
			agent.Residence.Location.Lng,
			agent.Residence.Location.Lat,
		)
	}
	if agent.Domicile != nil {
		agent.Domicile.BigLocation = lib.GetBigQueryNullGeography(
			agent.Domicile.Location.Lng,
			agent.Domicile.Location.Lat,
		)
	}
	agent.BigRuiRegistration = lib.GetBigQueryNullDateTime(agent.RuiRegistration)

	return agent
}

func parseBigQueryAgencyNode(agency *models.AgencyNode) *models.AgencyNode {
	if agency == nil {
		return nil
	}

	if agency.Address != nil {
		agency.Address.BigLocation = lib.GetBigQueryNullGeography(
			agency.Address.Location.Lng,
			agency.Address.Location.Lat,
		)
	}
	agency.Manager = parseBigQueryAgentNode(agency.Manager)
	agency.BigRuiRegistration = lib.GetBigQueryNullDateTime(agency.RuiRegistration)

	return agency
}

func contractorBigQuerySave(contractor *models.Contractor, collectionPrefix string) error {
	table := lib.GetDatasetByEnv("", fmt.Sprintf("%s%s", collectionPrefix, models.UserCollection))

	result, err := initBigqueryData(contractor)
	if err != nil {
		return err
	}

	log.Println("user save big query: " + result.Uid)

	return lib.InsertRowsBigQuery(models.WoptaDataset, table, contractor)
}

func initBigqueryData(contractor *models.Contractor) (*models.Contractor, error) {
	rawContractor, err := json.Marshal(contractor)
	if err != nil {
		return nil, err
	}
	contractor.Data = string(rawContractor)

	if contractor.BirthDate != "" {
		birthDate, err := time.Parse(time.RFC3339, contractor.BirthDate)
		if err != nil {
			return nil, err
		}
		contractor.BigBirthDate = lib.GetBigQueryNullDateTime(birthDate)
	}

	if contractor.Residence != nil {
		contractor.BigResidenceStreetName = contractor.Residence.StreetName
		contractor.BigResidenceStreetNumber = contractor.Residence.StreetNumber
		contractor.BigResidenceCity = contractor.Residence.City
		contractor.BigResidencePostalCode = contractor.Residence.PostalCode
		contractor.BigResidenceLocality = contractor.Residence.Locality
		contractor.BigResidenceCityCode = contractor.Residence.CityCode
	}

	if contractor.Domicile != nil {
		contractor.BigDomicileStreetName = contractor.Domicile.StreetName
		contractor.BigDomicileStreetNumber = contractor.Domicile.StreetNumber
		contractor.BigDomicileCity = contractor.Domicile.City
		contractor.BigDomicilePostalCode = contractor.Domicile.PostalCode
		contractor.BigDomicileLocality = contractor.Domicile.Locality
		contractor.BigDomicileCityCode = contractor.Domicile.CityCode
	}

	contractor.BigLocation = bigquery.NullGeography{
		// TODO: Check if correct: Geography type uses the WKT format for geometry
		GeographyVal: fmt.Sprintf("POINT (%f %f)", contractor.Location.Lng, contractor.Location.Lat),
		Valid:        true,
	}
	contractor.BigCreationDate = lib.GetBigQueryNullDateTime(contractor.CreationDate)
	contractor.BigUpdatedDate = lib.GetBigQueryNullDateTime(contractor.UpdatedDate)

	return contractor, err
}
