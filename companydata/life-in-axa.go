package companydata

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-gota/gota/dataframe"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
)

func LifeIn(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	const (
		slide       int = -1
		headervalue     = "N° adesione individuale univoco"
	)
	var skippedPolicies = make([]string, 0)

	data := lib.GetFromStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "track/in/life/life.csv", "")
	df := lib.CsvToDataframe(data)
	//log.Println("LifeIn  df.Describe: ", df.Describe())
	log.Println("LifeIn  row", df.Nrow())
	log.Println("LifeIn  col", df.Ncol())
	//group := df.GroupBy("N\xb0 adesione individuale univoco")
	group := GroupBy(df, 2)

	for v, d := range group {
		var (
			guarantees    []models.Guarante
			sumPriceGross float64
			maxDuration   int
		)
		if v != headervalue {
			if d[0][13] == "W1" || d[0][22] == "PG" {
				continue
			}

			for i, r := range d {
				log.Println("LifeIn  i: ", i)
				log.Println("LifeIn  d: ", r)
				companyCodec, slug, _, _ := LifeMapCodecCompanyAxaRevert(r[1])
				var (
					beneficiaries []models.Beneficiary
				)

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
			log.Println("LifeIn  row", len(d))
			log.Println("LifeIn  col", len(d[0]))
			log.Println("LifeIn  d: ", d)
			log.Println("LifeIn  elemets (0-0 ): ", d[0][0])
			log.Println("LifeIn  elemets (0-1 ): ", d[0][1])
			log.Println("LifeIn  elemets (0-2 ): ", d[0][2])
			log.Println("LifeIn  elemets (0-3 ): ", d[0][3])
			//1998-09-27T00:00:00Z RFC3339
			_, _, version, paymentSplit := LifeMapCodecCompanyAxaRevert(d[0][1])
			networkNode := network.GetNetworkNodeByCode(d[0][13])
			if networkNode == nil {
				log.Println("node not found!")
				skippedPolicies = append(skippedPolicies, d[0][0])
				continue
			}

			// create contractor

			contractor := models.User{
				Uid:        lib.NewDoc(models.UserCollection),
				Type:       d[0][22],
				Name:       d[0][23],
				Surname:    d[0][24],
				FiscalCode: strings.ToUpper(d[0][27]),
				Gender:     d[0][25],
				BirthDate:  ParseDateDDMMYYYY(d[0][26]).Format(time.RFC3339),
				Phone:      d[0][33],
				IdentityDocuments: []*models.IdentityDocument{{
					Code:             d[0][56],
					Type:             identityDocumentMap[d[0][56]],
					Number:           d[0][57],
					DateOfIssue:      ParseDateDDMMYYYY(d[0][58]),
					IssuingAuthority: d[0][59],
					PlaceOfIssue:     d[0][59],
				}},
				Residence: &models.Address{
					StreetName: d[0][28],
					City:       d[0][31],
					CityCode:   d[0][31],
					PostalCode: d[0][29],
					Locality:   d[0][30],
				},
			}

			// create insured

			insured := &models.User{
				Type:       d[0][22],
				Name:       d[0][35],
				Surname:    d[0][34],
				FiscalCode: strings.ToUpper(d[0][38]),
				Gender:     d[0][36],
				BirthDate:  ParseDateDDMMYYYY(d[0][37]).Format(time.RFC3339),
				Mail:       d[0][71],
				Phone:      d[0][72],
				IdentityDocuments: []*models.IdentityDocument{{
					Code:             d[0][77],
					Type:             d[0][76],
					DateOfIssue:      ParseDateDDMMYYYY(d[0][78]),
					IssuingAuthority: d[0][79],
				}},
				BirthCity:     d[0][37],
				BirthProvince: d[0][37],
				Residence: &models.Address{
					StreetName: d[0][63],
					City:       d[0][66],
					CityCode:   d[0][66],
					PostalCode: d[0][64],
					Locality:   d[0][65],
				},
				Domicile: &models.Address{
					StreetName: d[0][67],
					City:       d[0][70],
					CityCode:   d[0][70],
					PostalCode: d[0][68],
					Locality:   d[0][69],
				},
			}

			policy := models.Policy{
				Uid:            lib.NewDoc(models.PolicyCollection),
				Status:         models.PolicyStatusPay,
				StatusHistory:  []string{"Imported", models.PolicyStatusInitLead, models.PolicyStatusContact, models.PolicyStatusToSign, models.PolicyStatusSign, models.NetworkTransactionStatusToPay, models.PolicyStatusPay},
				Name:           models.LifeProduct,
				NameDesc:       "Wopta per te Vita",
				CodeCompany:    fmt.Sprintf("%07s", d[0][2]),
				Company:        models.AxaCompany,
				ProductVersion: "v" + version,
				IsPay:          true,
				IsSign:         true,
				CompanyEmit:    true,
				CompanyEmitted: true,
				Channel:        models.NetworkChannel,
				PaymentSplit:   paymentSplit,
				CreationDate:   ParseDateDDMMYYYY(d[0][4]),
				EmitDate:       ParseDateDDMMYYYY(d[0][4]),
				StartDate:      ParseDateDDMMYYYY(d[0][4]),
				EndDate:        ParseDateDDMMYYYY(d[0][4]).AddDate(maxDuration, 0, 0),
				Updated:        time.Now().UTC(),
				PriceGross:     sumPriceGross,
				PriceNett:      0,
				Payment:        models.ManualPaymentProvider,
				FundsOrigin:    "Proprie risorse economiche",
				ProducerCode:   networkNode.Code,
				ProducerUid:    networkNode.Uid,
				ProducerType:   networkNode.Type,
				Contractor:     contractor,
				Assets: []models.Asset{{
					Guarantees: guarantees,
					Person:     insured,
				}},
			}

			// create transactions

			// create network transactions

			// update node portfolio

			// save policy firestore

			// save policy bigquery

			// save contractor firestore

			// save contractor bigquery

			//log.Println("LifeIn policy:", policy)
			b, e := json.Marshal(policy)
			log.Println("LifeIn policy:", e)
			log.Println("LifeIn policy:", string(b))
			// docref, _, _ := lib.PutFirestoreErr("test-policy", policy)
			// log.Println("LifeIn doc id: ", docref.ID)

			//_, e = models.UpdateUserByFiscalCode("uat", policy.Contractor)
			//log.Println("LifeIn policy:", policy)
			//tr := transaction.PutByPolicy(policy, "", "uat", "", "", sumPriseGross, 0, "", "manual", true)
			//	log.Println("LifeIn transactionpolicy:",tr)
			//accounting.CreateNetworkTransaction(tr, "uat")

		}

	}

	log.Printf("Skipped %d policies: %v", len(skippedPolicies), skippedPolicies)

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
