package companydata

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-gota/gota/dataframe"
	lib "github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func LifeIn(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	const (
		slide int = -1
	)
	var (
		guarantees    []models.Guarante
		sumPriseGross float64
	)
	data := lib.GetFromStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "track/in/life/life.csv", "")
	df := lib.CsvToDataframe(data)
	//log.Println("LifeIn  df.Describe: ", df.Describe())
	log.Println("LifeIn  row", df.Nrow())
	log.Println("LifeIn  col", df.Ncol())
	//group := df.GroupBy("N\xb0 adesione individuale univoco")
	group := GroupBy(df, 2)

	for v, d := range group {

		log.Println("LifeIn  value", v)
		sumPriseGross = 0
		log.Println("LifeIn  row", len(d))
		log.Println("LifeIn  col", len(d[0]))
		log.Println("LifeIn  d: ", d)
		log.Println("LifeIn  elemets (0-0 ): ", d[0][0])
		log.Println("LifeIn  elemets (0-1 ): ", d[0][1])
		log.Println("LifeIn  elemets (0-2 ): ", d[0][2])
		log.Println("LifeIn  elemets (0-3 ): ", d[0][3])
		_, _, _, version := LifeMapCodecCompanyAxaRevert(d[0][1])
		policy := models.Policy{
			Name:           "life",
			CodeCompany:    "",
			Company:        "axa",
			ProductVersion: version,
			IsPay:          true,
			IsSign:         true,
			PaymentSplit:   "",
			StartDate:      ParseDateDDMMYYYY(d[0][4]),
			EndDate:        ParseDateDDMMYYYY(d[0][5]),

			Contractor: models.User{
				Type:       d[0][22],
				Name:       d[0][23],
				Surname:    d[0][24],
				FiscalCode: d[0][27],
				Gender:     d[0][25],
				BirthDate:  d[0][26],
				IdentityDocuments: []*models.IdentityDocument{{
					Code:             d[0][57],
					Type:             d[0][56],
					DateOfIssue:      ParseDateDDMMYYYY(d[0][58]),
					IssuingAuthority: d[0][59],
				}},
				Residence: &models.Address{
					StreetName: d[0][28],

					City:       d[0][31],
					CityCode:   d[0][31],
					PostalCode: d[0][29],
					Locality:   d[0][30],
				},
			},
			Assets: []models.Asset{{
				Name: "person",

				Person: &models.User{
					Type:       d[0][22],
					Name:       d[0][35],
					Surname:    d[0][34],
					FiscalCode: d[0][38],
					Gender:     d[0][36],
					BirthDate:  d[0][37],
					Mail:       d[0][71],
					Phone:      d[0][72],
					IdentityDocuments: []*models.IdentityDocument{{
						Code:             d[0][77],
						Type:             d[0][76],
						DateOfIssue:      ParseDateDDMMYYYY(d[0][78]),
						IssuingAuthority: d[0][79],
					}},
					BirthCity: d[0][37],

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
				},
			},
			},
		}

		for _, r := range d {
			log.Println("LifeIn  d: ", r)
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
	log.Println("LifeIn LifeMapCodecCompanyAxaRevert:", g)
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
	log.Println("LifeIn ParseDateDDMMYYYY date:", date)
	log.Println("LifeIn ParseDateDDMMYYYY len(date):", date)
	if len(date) < 8 {
		date = "0" + date
	}
	d, e := strconv.Atoi(date[0:1])
	m, e := strconv.Atoi(date[2:3])
	y, e := strconv.Atoi(date[4:7])

	res := time.Date(y, time.Month(m),
		d, 0, 0, 0, 0, time.UTC)
	log.Println(e)
	return res

}
func ParseAxaFloat(price string) float64 {
	log.Println("LifeIn ParseAxaFloat price:", price)
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
func GroupBy(df dataframe.DataFrame, col int) map[string][][]string {
	log.Println("GroupBy")
	res := make(map[string][][]string)
	for _, k := range df.Records() {
		if resFound, found := res[k[col]]; found {
			resFound = append(resFound, k)
		} else {
			res[k[col]] = [][]string{k}
		}

	}
	return res
}
