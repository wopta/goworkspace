package companydata

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	lib "github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func LifeIn(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		guarantees    []models.Guarante
		sumPriseGross float64
	)
	ricAteco := lib.GetFromStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "track/in/life/life.csv", "")

	df := lib.CsvToDataframe(ricAteco)
	log.Println("LifeIn  row", df.Nrow())
	log.Println("LifeIn  col", df.Ncol())
	group := df.GroupBy("N° adesione individuale univoco")

	for _, d := range group.GetGroups() {
		sumPriseGross = 0
		log.Println("LifeIn  row", d.Nrow())
		log.Println("LifeIn  col", d.Ncol())
		_, _, _, version := LifeMapCodecCompanyAxaRevert(d.Elem(0, 1).String())
		policy := models.Policy{
			Name:           "life",
			CodeCompany:    "",
			Company:        "axa",
			ProductVersion: version,
			IsPay:          true,
			IsSign:         true,
			PaymentSplit:   "",
			StartDate:      ParseDateDDMMYYYY(d.Elem(0, 4).String()),
			EndDate:        ParseDateDDMMYYYY(d.Elem(0, 5).String()),

			Contractor: models.User{
				Type:       d.Elem(0, 22).String(),
				Name:       d.Elem(0, 23).String(),
				Surname:    d.Elem(0, 24).String(),
				FiscalCode: d.Elem(0, 27).String(),
				Gender:     d.Elem(0, 25).String(),
				BirthDate:  d.Elem(0, 26).String(),
				IdentityDocuments: []*models.IdentityDocument{{
					Code:             d.Elem(0, 57).String(),
					Type:             d.Elem(0, 56).String(),
					DateOfIssue:      ParseDateDDMMYYYY(d.Elem(0, 58).String()),
					IssuingAuthority: d.Elem(0, 59).String(),
				}},
				Residence: &models.Address{
					StreetName: d.Elem(0, 28).String(),

					City:       d.Elem(0, 31).String(),
					CityCode:   d.Elem(0, 31).String(),
					PostalCode: d.Elem(0, 29).String(),
					Locality:   d.Elem(0, 30).String(),
				},
			},
			Assets: []models.Asset{{
				Name: "person",

				Person: &models.User{
					Type:       d.Elem(0, 22).String(),
					Name:       d.Elem(0, 35).String(),
					Surname:    d.Elem(0, 34).String(),
					FiscalCode: d.Elem(0, 38).String(),
					Gender:     d.Elem(0, 36).String(),
					BirthDate:  d.Elem(0, 37).String(),
					Mail:       d.Elem(0, 71).String(),
					Phone:      d.Elem(0, 72).String(),
					IdentityDocuments: []*models.IdentityDocument{{
						Code:             d.Elem(0, 77).String(),
						Type:             d.Elem(0, 76).String(),
						DateOfIssue:      ParseDateDDMMYYYY(d.Elem(0, 78).String()),
						IssuingAuthority: d.Elem(0, 79).String(),
					}},
					BirthCity: d.Elem(0, 37).String(),

					BirthProvince: d.Elem(0, 37).String(),
					Residence: &models.Address{
						StreetName: d.Elem(0, 63).String(),
						City:       d.Elem(0, 66).String(),
						CityCode:   d.Elem(0, 66).String(),
						PostalCode: d.Elem(0, 64).String(),
						Locality:   d.Elem(0, 65).String(),
					},
					Domicile: &models.Address{
						StreetName: d.Elem(0, 67).String(),
						City:       d.Elem(0, 70).String(),
						CityCode:   d.Elem(0, 70).String(),
						PostalCode: d.Elem(0, 68).String(),
						Locality:   d.Elem(0, 69).String(),
					},
				},
			},
			},
		}

		for _, r := range d.Records() {
			result, _, slug, _ := LifeMapCodecCompanyAxaRevert(r[1])
			var (
				beneficiaries []models.Beneficiary
				benef1        models.Beneficiary
			)
			benef1 = ParseAxaBeneficiary(r, 0)
			benef2 := ParseAxaBeneficiary(r, 1)
			benef3 := ParseAxaBeneficiary(r, 2)
			beneficiaries = append(beneficiaries, benef1)
			beneficiaries = append(beneficiaries, benef2)
			beneficiaries = append(beneficiaries, benef3)
			dur, _ := strconv.Atoi(r[7])
			priceGross := ParseAxaFloat(r[8])
			sumPriseGross = priceGross + sumPriseGross
			var guarante models.Guarante = models.Guarante{
				Slug:                       slug,
				CompanyCodec:               result,
				SumInsuredLimitOfIndemnity: 0,

				Beneficiaries: &beneficiaries,
				Value: &models.GuaranteValue{
					SumInsuredLimitOfIndemnity: ParseAxaFloat(r[9]),
					PremiumGrossYearly:         priceGross,
					Duration: &models.Duration{
						Year: dur / 12,
					},
				},
			}

			guarantees = append(guarantees, guarante)
		}

		policy.Assets[0].Guarantees = guarantees
		policy.PriceGross = sumPriseGross

		log.Println("LifeIn policy:", policy)
		lib.PutFirestoreErr("policy", policy)
		//user,e:=models.UpdateUserByFiscalCode("", policy.Contractor)
		//tr := transaction.PutByPolicy(policy, "", "", "", "", sumPriseGross, 0, "", "BO", true)
		//accounting.CreateNetworkTransaction(tr, "uat")

	}

	return "", nil, e
}

func LifeMapCodecCompanyAxaRevert(g string) (string, string, string, string) {
	var result, pay, slug, version string
	version = g[0:0]
	code := g[2:2]

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
	return result, slug, version, pay
}
func ParseDateDDMMYYYY(date string) time.Time {
	if len(date) < 8 {
		date = "0" + date
	}
	d, e := strconv.Atoi(date[0:1])
	m, e := strconv.Atoi(date[2:3])
	y, e := strconv.Atoi(date[4:7])
	log.Println("LifeIn LifeMapCodecCompanyAxaRevert d:", d)
	res := time.Date(y, time.Month(m),
		d, 0, 0, 0, 0, time.UTC)
	log.Println(e)
	return res

}
func ParseAxaFloat(price string) float64 {
	d := price[len(price)-2 : len(price)-1]
	i := price[0 : len(price)-3]
	f64string := d + "." + i
	res, e := strconv.ParseFloat(f64string, 64)
	log.Println("LifeIn ParseAxaFloat d:", res)
	log.Println(e)
	return res

}
func ParseAxaBeneficiary(r []string, base int) models.Beneficiary {
	var (
		benef models.Beneficiary
	)
	rangeCell := 11 * base

	if r[82] == "GE" {
		benef = models.Beneficiary{

			User: models.User{Name: "",
				Surname:    "",
				FiscalCode: ""},
			IsLegitimateSuccessors: true,
		}

	}

	if r[82] == "NM" {
		benef = models.Beneficiary{

			User: models.User{
				Name:       r[84+rangeCell],
				Surname:    r[83+rangeCell],
				FiscalCode: r[85+rangeCell],
				Mail:       r[91+rangeCell],

				Residence: &models.Address{
					StreetName: r[87+rangeCell],
					City:       r[90+rangeCell],
					CityCode:   r[90+rangeCell],
					PostalCode: r[89+rangeCell],
					Locality:   r[88+rangeCell],
				},
			},
			IsLegitimateSuccessors: false,
		}

	} else {
		benef = models.Beneficiary{

			User: models.User{Name: "",
				Surname:    "",
				FiscalCode: ""},
			IsLegitimateSuccessors: false,
		}
	}
	return benef

}
