package companydata

import (
	"encoding/json"
	"fmt"
	"github.com/wopta/goworkspace/network"
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
		policies        = make([]models.Policy, 0)
		skippedPolicies = make([]string, 0)
		monthlyPolicies = make(map[string]map[string][][]string, 0)
	)

	log.Println(os.Getwd())

	data, _ := os.ReadFile("./companydata/track_in_life.csv")
	//data := lib.GetFromStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "track/in/life/life.csv", "")
	df := lib.CsvToDataframe(data)
	//log.Println("LifeIn  df.Describe: ", df.Describe())
	log.Println("LifeIn  row", df.Nrow())
	log.Println("LifeIn  col", df.Ncol())
	//group := df.GroupBy("N\xb0 adesione individuale univoco")
	group := GroupBy(df, 2)

	for v, pol := range group {
		var (
			row           []string
			guarantees    []models.Guarante
			sumPriceGross float64
			maxDuration   int
		)

		if pol[0][3] == headervalue || pol[0][13] == "W1" || pol[0][22] == "PG" {
			continue
		}

		row = pol[0]

		for i, r := range pol {
			var (
				beneficiaries []models.Beneficiary
			)

			log.Println("LifeIn  i: ", i)
			log.Println("LifeIn  pol: ", r)

			if r[3] == "R" {
				if monthlyPolicies[r[2]] == nil {
					monthlyPolicies[r[2]] = make(map[string][][]string, 0)
				}
				if monthlyPolicies[r[2]][r[5]] == nil {
					monthlyPolicies[r[2]][r[5]] = make([][]string, 0)
				}
				monthlyPolicies[r[2]][r[5]] = append(monthlyPolicies[r[2]][r[5]], r)
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
					beneficiaries = append(beneficiaries, benef1)
					beneficiaries = append(beneficiaries, benef2)
				}
			}
			dur, _ := strconv.Atoi(r[7])
			guaranteeYearDuration := dur / 12

			if guaranteeYearDuration > maxDuration {
				maxDuration = guaranteeYearDuration
			}

			priceGross := ParseAxaFloat(r[8])
			sumPriceGross += priceGross
			var guarante models.Guarante = models.Guarante{
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
		networkNode := network.GetNetworkNodeByCode(row[13])
		if networkNode == nil {
			log.Println("node not found!")
			skippedPolicies = append(skippedPolicies, row[0])
			continue
		}

		// create insured

		insured := &models.User{
			Type:       row[22],
			Name:       strings.TrimSpace(lib.Capitalize(row[35])),
			Surname:    strings.TrimSpace(lib.Capitalize(row[34])),
			FiscalCode: strings.TrimSpace(strings.ToUpper(row[38])),
			Gender:     strings.TrimSpace(strings.ToUpper(row[36])),
			BirthDate:  ParseDateDDMMYYYY(row[37]).Format(time.RFC3339),
			Mail:       strings.TrimSpace(strings.ToLower(row[71])),
			Phone:      row[72],
			IdentityDocuments: []*models.IdentityDocument{{
				Code:             row[77],
				Type:             identityDocumentMap[row[76]],
				DateOfIssue:      ParseDateDDMMYYYY(row[78]),
				IssuingAuthority: strings.TrimSpace(lib.Capitalize(row[79])),
			}},
			BirthCity:     strings.TrimSpace(lib.Capitalize(row[73])),
			BirthProvince: strings.TrimSpace(strings.ToUpper(row[74])),
			Residence: &models.Address{
				StreetName: strings.TrimSpace(lib.Capitalize(row[63])),
				City:       strings.TrimSpace(lib.Capitalize(row[65])),
				CityCode:   strings.TrimSpace(strings.ToUpper(row[66])),
				PostalCode: row[64],
				Locality:   strings.TrimSpace(lib.Capitalize(row[65])),
			},
			Domicile: &models.Address{
				StreetName: strings.TrimSpace(lib.Capitalize(row[67])),
				City:       strings.TrimSpace(lib.Capitalize(row[70])),
				CityCode:   strings.TrimSpace(strings.ToUpper(row[70])),
				PostalCode: row[68],
				Locality:   strings.TrimSpace(lib.Capitalize(row[69])),
			},
		}

		policy := models.Policy{
			Uid:            lib.NewDoc(models.PolicyCollection),
			Status:         models.PolicyStatusPay,
			StatusHistory:  []string{"Imported", models.PolicyStatusInitLead, models.PolicyStatusContact, models.PolicyStatusToSign, models.PolicyStatusSign, models.NetworkTransactionStatusToPay, models.PolicyStatusPay},
			Name:           models.LifeProduct,
			NameDesc:       "Wopta per te Vita",
			CodeCompany:    fmt.Sprintf("%07s", strings.TrimSpace(row[2])),
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
			Contractor:     *insured,
			Assets: []models.Asset{{
				Guarantees: guarantees,
				Person:     insured,
			}},
		}
		policy.Contractor.Uid = lib.NewDoc(models.UserCollection)

		// create transactions

		// create network transactions

		// update node portfolio

		networkNode.Policies = append(networkNode.Policies, policy.Uid)
		networkNode.Users = append(networkNode.Users, policy.Contractor.Uid)

		// save policy firestore

		// save policy bigquery

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

	log.Printf("Skipped %d policies: %v", len(skippedPolicies), skippedPolicies)
	log.Printf("Created %d policies ", len(policies))

	out, err := json.Marshal(policies)
	if err != nil {
		log.Printf("error: %s", err.Error())
	}
	err = os.WriteFile("./companydata/result.json", out, 0777)
	if err != nil {
		log.Printf("error: %s", err.Error())
	}

	return "", nil, e
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
	// //pricelen:=len("0000001500000")
	// if len(price) > 3 {
	// 	log.Println("LifeIn ParseAxaFloat price:", price)

	// 	d := price[len(price)-2:]
	// 	i := price[:len(price)-3]
	// 	f64string := i + "." + d
	// 	res, e := strconv.ParseFloat(f64string, 64)
	// 	log.Println("LifeIn ParseAxaFloat d:", res)
	// 	log.Println(e)
	// 	return res
	// }
	// return 0

	princeInCents, _ := strconv.ParseFloat(price, 64)
	return princeInCents / 100.0
}

func ParseAxaBeneficiary(r []string, base int) models.Beneficiary {
	var (
		benef models.Beneficiary
	)
	rangeCell := 11 * base

	if r[82] == "GE" {
		benef = models.Beneficiary{
			BeneficiaryType: "legalAndWillSuccessor",
		}
	}
	if r[82] == "NM" {
		benef = models.Beneficiary{
			BeneficiaryType: "chosenBeneficiary",
			User: models.User{
				Name:       r[84+rangeCell],
				Surname:    r[83+rangeCell],
				FiscalCode: strings.ToUpper(r[85+rangeCell]),
				Mail:       r[91+rangeCell],
				Residence: &models.Address{
					StreetName: r[87+rangeCell],
					City:       r[90+rangeCell],
					CityCode:   r[90+rangeCell],
					PostalCode: r[89+rangeCell],
					Locality:   r[88+rangeCell],
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
