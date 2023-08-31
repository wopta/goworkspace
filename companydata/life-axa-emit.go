package companydata

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	lib "github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func LifeAxaEmit(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		from          time.Time
		to            time.Time
		filenamesplit string
		cabCsv        []byte
		result        [][]string
		refMontly     time.Time
		upload        bool
	)
	var lifeAxaEmitQuery = lib.Firequeries{
		Queries: []lib.Firequery{

			{
				Field:      "isDelete", //
				Operator:   "==",       //
				QueryValue: false,
			},
			{
				Field:      "isPay", //
				Operator:   "==",    //
				QueryValue: true,
			},
			{
				Field:      "company", //
				Operator:   "==",      //
				QueryValue: "axa",
			},
			{
				Field:      "policyName", //
				Operator:   "==",         //
				QueryValue: "life",
			}, {
				Field:      "payDate", //
				Operator:   ">",       //
				QueryValue: from,
			},
			{
				Field:      "payDate", //
				Operator:   "<",       //
				QueryValue: to,
			},
		},
	}
	log.Println("----------------LifeAxalEmit-----------------")
	req := lib.ErrorByte(ioutil.ReadAll(r.Body))
	defer r.Body.Close()
	log.Println("LifeAxalEmit: ", r.Header)
	log.Println("LifeAxalEmit: ", string(req))
	now, upload := getRequestData(req)
	from, to, refMontly, filenamesplit = AxaPartnersSchedule(now)
	cabCsv = lib.GetFilesByEnv("data/cab-cap-istat.csv")

	log.Println("LifeAxalEmit now: " + now.String())
	log.Println("LifeAxalEmit now.Day: ", now.Day())
	log.Println("LifeAxalEmit from: " + from.String())
	log.Println("LifeAxalEmit to: " + to.String())
	log.Println("LifeAxalEmit: " + filenamesplit)
	lifeAxaEmitQuery = lib.Firequeries{
		Queries: []lib.Firequery{

			{
				Field:      "isDelete", //
				Operator:   "==",       //
				QueryValue: false,
			},
			{
				Field:      "isPay", //
				Operator:   "==",    //
				QueryValue: true,
			},
			{
				Field:      "company", //
				Operator:   "==",      //
				QueryValue: "axa",
			},
			{
				Field:      "policyName", //
				Operator:   "==",         //
				QueryValue: "life",
			}, {
				Field:      "payDate", //
				Operator:   ">",       //
				QueryValue: from,
			},
			{
				Field:      "payDate", //
				Operator:   "<",       //
				QueryValue: to,
			},
		},
	}
	df := lib.CsvToDataframe(cabCsv)
	query, e := lifeAxaEmitQuery.FirestoreWherefields("transactions")
	log.Println("LifeAxalEmit: ", e)
	transactions := models.TransactionToListData(query)
	log.Println("LifeAxalEmit: transaction len: ", len(transactions))
	//result = append(result, getHeader())
	for _, transaction := range transactions {
		var (
			policy models.Policy
		)
		docsnap := lib.GetFirestore("policy", transaction.PolicyUid)
		docsnap.DataTo(&policy)
		result = append(result, setRowLifeEmit(policy, df, transaction)...)

		transaction.IsEmit = true

		lib.SetFirestore("transactions", transaction.Uid, transaction)

	}

	filepath := "WOPTAKEYweb_NB" + filenamesplit + "_" + strconv.Itoa(refMontly.Year()) + fmt.Sprintf("%02d", int(refMontly.Month())) + "_" + fmt.Sprintf("%02d", now.Day()) + fmt.Sprintf("%02d", int(now.Month())) + ".txt"
	lib.WriteCsv("../tmp/"+filepath, result, ';')
	source, _ := ioutil.ReadFile("../tmp/" + filepath)
	lib.PutToStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "track/axa/life/"+strconv.Itoa(refMontly.Year())+"/"+filepath, source)
	if upload {
		AxaPartnersSftpUpload(filepath)
	}
	return "", nil, e
}

func setRowLifeEmit(policy models.Policy, df dataframe.DataFrame, trans models.Transaction) [][]string {
	var (
		result       [][]string
		residenceCab string
	)
	log.Println("LifeAxalEmit:  policy.Uid: ", policy.Uid)
	fil := df.Filter(
		dataframe.F{Colidx: 4, Colname: "CAP", Comparator: series.Eq, Comparando: policy.Contractor.Residence.PostalCode},
	)
	if fil.Nrow() > 0 {
		residenceCab = fil.Records()[1][5]
	}

	log.Println("LifeAxalEmit:  residenceCab:", residenceCab)
	log.Println("LifeAxalEmit:  fil.Records()[0]:", fil.Records()[0])
	log.Println("LifeAxalEmit:  filtered col", fil.Ncol())
	log.Println("LifeAxalEmit: filtered row", fil.Nrow())
	for _, asset := range policy.Assets {

		for _, g := range asset.Guarantees {
			var (
				price                                     float64
				beneficiary1, beneficiary2, beneficiary1T string
				beneficiary1S, beneficiary2S              models.Beneficiary
			)
			beneficiary1, beneficiary1S, beneficiary1T = mapBeneficiary(g, 0) //Codice Fiscale Beneficiario
			beneficiary2, beneficiary2S, _ = mapBeneficiary(g, 0)
			if policy.PaymentSplit == string(models.PaySplitMonthly) {
				price = g.Value.PremiumGrossMonthly
			} else {
				price = g.Value.PremiumGrossYearly
			}

			log.Println("LifeAxalEmit: ", price)
			var intNum = int(price * 100)
			priceGrossFormat := fmt.Sprintf("%012d", intNum) // 000000001220
			log.Println("LifeAxalEmit: ", priceGrossFormat)
			row := []string{
				mapCodecCompany(policy, g.CompanyCodec),                 //Codice schema
				policy.CodeCompany,                                      //N° adesione individuale univoco
				getRenew(policy),                                        //Tipo di Transazione
				getFormatdate(policy.StartDate),                         //Data di decorrenza
				getFormatdate(getRenewDate(policy, trans)),              //"Data di rinnovo"
				mapCoverageDuration(policy),                             //"Durata copertura assicurativa"
				fmt.Sprint(g.Value.Duration.Year * 12),                  //"Durata complessiva"
				priceGrossFormat,                                        //"Premio assicurativo lordo"
				fmt.Sprintf("%.0f", g.Value.SumInsuredLimitOfIndemnity), //"Importo Assicurato"
				"0",                       //indennizzo mensile
				"",                        //campo disponibile
				"",                        //% di sovrappremio da applicare alla garanzia
				"W1",                      //Codice Concessionario /dipendenti (iscr.E)
				"",                        //Codice Campagna
				"T",                       //Copertura Assicurativa: Totale o Pro quota
				"",                        //% assicurata dell'assicurato
				"",                        //campo disponibile
				"",                        //Maxi rata finale/Valore riscatto
				"",                        //Stato occupazionale dell'Assicurato
				"1",                       //Tipo aderente
				"WEB",                     //Canale di vendita
				"PF",                      //Tipo contraente / Contraente
				policy.Contractor.Surname, //Denominazione Sociale o Cognome contraente
				policy.Contractor.Name,    //campo vuoto o nome
				policy.Contractor.Gender,  //Sesso
				getFormatBithdate(policy.Contractor.BirthDate),       //Data di nascita
				policy.Contractor.FiscalCode,                         //Codice Fiscale
				policy.Contractor.Residence.StreetName,               //Indirizzo di residenza
				policy.Contractor.Residence.PostalCode,               //C.A.P. Di residenza
				policy.Contractor.Residence.Locality,                 //Comune di residenza
				policy.Contractor.Residence.City,                     //Provincia di residenza
				policy.Contractor.Mail,                               //Indirizzo e-mail
				policy.Contractor.Phone,                              //Numero di Cellulare
				policy.Assets[0].Person.Surname,                      //Cognome Assicurato
				policy.Assets[0].Person.Name,                         //Nome
				policy.Assets[0].Person.Gender,                       //Sesso
				getFormatBithdate(policy.Assets[0].Person.BirthDate), //Data di nascita
				policy.Assets[0].Person.FiscalCode,                   //Codice Fiscale
				beneficiary1,                                         //Codice Fiscale Beneficiario
				beneficiary2,                                         //Codice Fiscale Beneficiario 2
				"",                                                   //Codice Fiscale Beneficiario 3
				"VIT",                                                //Natura del rapporto
				"PAS",                                                //Scopo del rapporto
				"BO",                                                 //Modalità di pagamento del premio assicurativo (all'intermediario)
				"SI",                                                 //contraente = Assicurato?
				ChekDomicilie(policy.Contractor).StreetName,          //Indirizzo di domicilio contraente
				ChekDomicilie(policy.Contractor).PostalCode,          //C.A.P. Di domicilio
				ChekDomicilie(policy.Contractor).Locality,            //Comune di domicilio
				ChekDomicilie(policy.Contractor).CityCode,            //Provincia di domicilio
				policy.Contractor.BirthCity,                          //Luogo di nascita dell’contraente persona fisica
				policy.Contractor.BirthProvince,                      //Provincia di nascita dell’contraente persona fisica
				"086",                                                //Stato di residenza dell’contraente
				residenceCab,                                         //Cab della città di residenza dell’contraente
				"600",                                                //Sottogruppo attività economica
				"600",                                                //Ramo gruppo attività economica
				ExistIdentityDocument(policy.Contractor.IdentityDocuments).Code,                       //Tipo documento dell'contraente persona fisica
				ExistIdentityDocument(policy.Contractor.IdentityDocuments).Number,                     //Numero documento dell'contraente persona fisica
				getFormatdate(ExistIdentityDocument(policy.Contractor.IdentityDocuments).DateOfIssue), //Data rilascio documento dell'contraente persona fisica
				ExistIdentityDocument(policy.Contractor.IdentityDocuments).IssuingAuthority,           //Ente rilascio documento dell'contraente persona fisica
				"NO", //PEP - Persona Politicamente Esposta
				"",   //Tipologia di PEP
				"E",  //Modalità di comunicazione prescelta tra Compagnia ed contraente
				policy.Assets[0].Person.Residence.StreetName,                                                //Indirizzo di residenza Assicurato
				policy.Assets[0].Person.Residence.PostalCode,                                                //C.A.P. Residenza
				policy.Assets[0].Person.Residence.Locality,                                                  //Comune Residenza
				policy.Assets[0].Person.Residence.CityCode,                                                  //Provincia Residenza
				ChekDomicilie(*policy.Assets[0].Person).StreetName,                                          //Indirizzo di domicilio
				ChekDomicilie(*policy.Assets[0].Person).PostalCode,                                          //C.A.P. Domicilio
				ChekDomicilie(*policy.Assets[0].Person).Locality,                                            //Comune Domicilio
				ChekDomicilie(*policy.Assets[0].Person).CityCode,                                            //Provincia Domicilio
				policy.Assets[0].Person.Mail,                                                                //Indirizzo e-mail
				policy.Assets[0].Person.Phone,                                                               //Numero di cellulare
				policy.Assets[0].Person.BirthCity,                                                           //Luogo di nascita
				policy.Assets[0].Person.BirthCity,                                                           //Provincia di nascita
				"ITA",                                                                                       //Stato di residenza
				ExistIdentityDocument(policy.Assets[0].Person.IdentityDocuments).Code,                       //Tipo documento
				ExistIdentityDocument(policy.Assets[0].Person.IdentityDocuments).Number,                     //Numero documento
				getFormatdate(ExistIdentityDocument(policy.Assets[0].Person.IdentityDocuments).DateOfIssue), //Data rilascio documento
				ExistIdentityDocument(policy.Assets[0].Person.IdentityDocuments).IssuingAuthority,           //Ente rilascio documento
				"NO",                                  //PEP - Persona Politicamente Esposta
				"",                                    //Tipologia di PEP
				beneficiary1T,                         //Eredi designati nominativamente o genericamente?
				beneficiary1S.Surname,                 //Cognome Beneficiario 1
				beneficiary1S.Name,                    //Nome
				beneficiary1S.FiscalCode,              //Codice Fiscale
				beneficiary1S.Phone,                   //Numero di Telefono del Beneficiario
				beneficiary1S.Residence.StreetName,    //Indirizzo di residenza
				beneficiary1S.Residence.City,          //Città /Comune di Residenza
				beneficiary1S.Residence.PostalCode,    //CAP
				beneficiary1S.Residence.CityCode,      //Provincia
				beneficiary1S.Mail,                    //Email
				MapBool(beneficiary1S.IsFamilyMember), //Legame del Cliente col Beneficiario

				MapBool(beneficiary1S.IsContactable),  //Lcontraente ha escluso l invio di comunicazioni da parte dell Impresa al Beneficiario?
				beneficiary2S.Surname,                 //Cognome Beneficiario 2
				beneficiary2S.Name,                    //Nome
				beneficiary2S.FiscalCode,              //Codice Fiscale
				beneficiary2S.Phone,                   //Numero di Telefono del Beneficiario
				beneficiary2S.Residence.StreetName,    //Indirizzo di residenza
				beneficiary2S.Residence.City,          //Città /Comune di Residenza
				beneficiary2S.Residence.PostalCode,    //CAP
				beneficiary2S.Residence.CityCode,      //Provincia
				beneficiary2S.Mail,                    //Email
				MapBool(beneficiary2S.IsFamilyMember), //Legame del Cliente col Beneficiario
				MapBool(beneficiary2S.IsContactable),  //NUCLEO FAMILIARE
				"",                                    //L contraente ha escluso l invio di comunicazioni da parte dell Impresa al Beneficiario?
				"",                                    //Cognome Beneficiario 3
				"",                                    //Nome
				"",                                    //Codice Fiscale
				"",                                    //Numero di Telefono del Beneficiario
				"",                                    //Indirizzo di residenza
				"",                                    //Città /Comune di Residenza
				"",                                    //CAP
				"",                                    //Provincia
				"",                                    //Email
				"",                                    //Legame del Cliente col Beneficiario
				"",                                    //L'contraente ha escluso l'invio di comunicazioni da parte dell Impresa al Beneficiario?
				"NO",                                  //Esistenza Titolare effettiv
				"",                                    //Cognome
				"",                                    //Nome
				"",                                    //Sesso
				"",                                    //Data di nascita
				"",                                    //Codice Fiscale
				"",                                    //Indirizzo di residenza
				"",                                    //C.A.P. Residenza
				"",                                    //Comune Residenza
				"",                                    //Provincia Residenza
				"",                                    //Indirizzo di domicilio
				"",                                    //C.A.P. Domicilio
				"",                                    //Comune Domicilio
				"",                                    //Provincia Domicilio
				"",                                    //Stato occupazionale
				"",                                    //Indirizzo e-mail
				"",                                    //Numero di cellulare
				"",                                    //Luogo di nascita
				"",                                    //Provincia di nascita
				"",                                    //Stato di residenza
				"",                                    //Tipo documento
				"",                                    //Numero documento
				"",                                    //Data rilascio documento
				"",                                    //Ente rilascio documento
				"",                                    //PEP - Persona Politicamente Esposta
				"",                                    //Tipologia di PEP
				"NO",                                  //Esistenza Titolare effettivo
				"",                                    //Cognome
				"",                                    //Nome
				"",                                    //Sesso
				"",                                    //Data di nascita
				"",                                    //Codice Fiscale
				"",                                    //Indirizzo di residenza
				"",                                    //C.A.P. Residenza
				"",                                    //Comune Residenza
				"",                                    //Provincia Residenza
				"",                                    //Indirizzo di domicilio
				"",                                    //C.A.P. Domicilio
				"",                                    //Comune Domicilio
				"",                                    //Provincia Domicilio
				"",                                    //Stato occupazionale
				"",                                    //Indirizzo e-mail
				"",                                    //Numero di cellulare
				"",                                    //Luogo di nascita
				"",                                    //Provincia di nascita
				"",                                    //Stato di residenza
				"",                                    //Tipo documento
				"",                                    //Numero documento
				"",                                    //Data rilascio documento
				"",                                    //Ente rilascio documento
				"",                                    //PEP - Persona Politicamente Esposta
				"",                                    //Tipologia di PEP
				"NO",                                  //Esistenza Titolare effettivo
				"",                                    //Cognome
				"",                                    //Nome
				"",                                    //Sesso
				"",                                    //Data di nascita
				"",                                    //Codice Fiscale
				"",                                    //Indirizzo di residenza
				"",                                    //C.A.P. Residenza
				"",                                    //Comune Residenza
				"",                                    //Provincia Residenza
				"",                                    //Indirizzo di domicilio
				"",                                    //C.A.P. Domicilio
				"",                                    //Comune Domicilio
				"",                                    //Provincia Domicilio
				"",                                    //Stato occupazionale
				"",                                    //Indirizzo e-mail
				"",                                    //Numero di cellulare
				"",                                    //Luogo di nascita
				"",                                    //Provincia di nascita
				"",                                    //Stato di residenza
				"",                                    //Tipo documento
				"",                                    //Numero documento
				"",                                    //Data rilascio documento
				"",                                    //Ente rilascio documento
				"",                                    //PEP - Persona Politicamente Esposta
				"",                                    //Tipologia di PEP
				"NO",                                  //Esistenza Titolare effettivo
				"",                                    //Cognome
				"",                                    //Nome
				"",                                    //Sesso
				"",                                    //Data di nascita
				"",                                    //Codice Fiscale
				"",                                    //Indirizzo di residenza
				"",                                    //C.A.P. Residenza
				"",                                    //Comune Residenza
				"",                                    //Provincia Residenza
				"",                                    //Indirizzo di domicilio
				"",                                    //C.A.P. Domicilio
				"",                                    //Comune Domicilio
				"",                                    //Provincia Domicilio
				"",                                    //Stato occupazionale
				"",                                    //Indirizzo e-mail
				"",                                    //umero di cellulare
				"",                                    //Luogo di nascita
				"",                                    //Provincia di nascita
				"",                                    //Stato di residenza
				"",                                    //Tipo documento"
				"",                                    //Numero documento
				"",                                    //Data rilascio documento
				"",                                    //Ente rilascio documento
				"",                                    //PEP - Persona Politicamente Esposta
				"",                                    //Tipologia di PEP
				"",                                    //Cognome
				"",                                    //Nome
				"",                                    //Sesso
				"",                                    //Data di nascita
				"",                                    //Codice Fiscale
				"",                                    //Indirizzo di residenza
				"",                                    //C.A.P. Di residenza
				"",                                    //Comune di residenza
				"",                                    //Provincia di residenza
				"",                                    //Indirizzo di domicilio
				"",                                    //C.A.P. Di domicilio
				"",                                    //Comune di domicilio
				"",                                    //Provincia di domicilio
				"",                                    //Indirizzo e-mail
				"",                                    //Numero di Cellulare
				"",                                    //Luogo di nascita dell’esecutore
				"",                                    //Provincia di nascita dell’esecutore
				"",                                    //Stato di residenza dell’esecutore
				"",                                    //Tipo documento
				"",                                    //Numero documento
				"",                                    //Data rilascio documento"
				"",                                    //Ente rilascio documento
				"",                                    //PEP - Persona Politicamente Esposta
				"",                                    //Tipologia di PEP
				"",                                    //Carica ricoperta dell'esecutore
				"",                                    //Cognome
				"",                                    //Nome
				"",                                    //Indirizzo di residenza
				"",                                    //Città /Comune di Residenza
				"",                                    //CAP
				"",                                    //Codice Fiscale
				"",                                    //Numero di Telefono
				"",                                    //Email
			}

			result = append(result, row)

		}

	}
	log.Println("----------------End LifeAxalEmit-----------------")
	return result
}
func getFormatdate(d time.Time) string {
	var res string
	res = fmt.Sprintf("%02d", d.Day()) + fmt.Sprintf("%02d", int(d.Month())) + strconv.Itoa(d.Year())
	return res

}

// 1989-03-13T00:00:00Z
func getFormatBithdate(d string) string {
	var res string
	if d != "" {
		splitD := strings.Split(d, "-")
		split2 := strings.Split(splitD[2], "T")
		day, _ := strconv.Atoi(split2[0])
		month, _ := strconv.Atoi(splitD[1])
		res = fmt.Sprintf("%02d", day) + fmt.Sprintf("%02d", month) + splitD[0]
	}

	return res

}
func getRenew(p models.Policy) string {
	var result string
	now := time.Now()
	addMonth := p.StartDate.AddDate(0, 1, 0)
	if p.PaymentSplit == "year" {
		result = "A"
	}
	if p.PaymentSplit == string(models.PaySplitMonthly) {
		if addMonth.Before(now) {
			result = "A"
		} else {
			result = "R"
		}
	}
	return result
}
func getRenewDate(p models.Policy, trans models.Transaction) time.Time {
	var result time.Time
	now := time.Now()
	addMonth := p.StartDate.AddDate(0, 1, 0)
	if p.PaymentSplit == "year" {
		result = p.StartDate
	}
	if p.PaymentSplit == string(models.PaySplitMonthly) {

		if addMonth.Before(now) {
			result = p.StartDate
		} else {
			result = trans.PayDate
		}
	}
	return result
}
func mapCoverageDuration(p models.Policy) string {
	var result string
	if p.PaymentSplit == "year" {
		result = "012"
	}
	if p.PaymentSplit == string(models.PaySplitMonthly) {
		result = "001"
	}
	return result
}
func mapCodecCompany(p models.Policy, g string) string {
	var result, pay string

	if p.PaymentSplit == "year" {
		pay = "W"
	}
	if p.PaymentSplit == string(models.PaySplitMonthly) {
		pay = "M"
	}
	if g == "D" {
		result = "1" + pay + "5"
	}
	if g == "PTD" {
		result = "1" + pay + "6"
	}
	if g == "TTD" {
		result = "1" + pay + "7"
	}
	if g == "CI" {
		result = "1" + pay + "8"
	}
	return result
}
func ChekDomicilie(u models.User) models.Address {
	var res models.Address
	//log.Println(reflect.ValueOf(u.Domicile))
	if reflect.ValueOf(u.Domicile).IsNil() {
		res = *u.Residence
	} else {
		res = *u.Domicile
	}
	return res
}

func mapBeneficiary(g models.Guarante, b int) (string, models.Beneficiary, string) {
	var (
		result      string
		result2     string
		resulStruct models.Beneficiary
	)
	resulStructDefault := models.Beneficiary{User: models.User{Residence: &models.Address{}}}
	resulStruct = resulStructDefault
	if g.Beneficiaries != nil {

		if len(*g.Beneficiaries) > 0 && len(*g.Beneficiaries) > b {
			result = ""
			if (*g.Beneficiaries)[b].IsLegitimateSuccessors || (*g.Beneficiaries)[b].IsFamilyMember || (*g.Beneficiaries)[b].BeneficiaryType == models.BeneficiaryLegalAndWillSuccessors {
				result = "GE"
				result2 = "GE"
				resulStruct = resulStructDefault

			} else {
				result = (*g.Beneficiaries)[b].FiscalCode
				result2 = "NM"
				resulStruct = (*g.Beneficiaries)[b]
			}

		}
	}
	return result, resulStruct, result2
}
func MapBool(s bool) string {
	var res string
	res = "NO"
	if s {
		res = "SI"
	}
	return res
}
func ExistIdentityDocument(docs []*models.IdentityDocument) *models.IdentityDocument {
	var (
		result *models.IdentityDocument
	)
	result = &models.IdentityDocument{}
	if len(docs) > 0 {
		for _, doc := range docs {
			log.Println("LifeAxalEmit: ", doc)
			result = doc

		}

	}
	return result
}

func getHeader() []string {
	return []string{
		"Codice schema",
		"N° adesione individuale univoco",
		"Tipo di Transazione",
		"Data di decorrenza",
		"Data di rinnovo",
		"Durata copertura assicurativa",
		"Durata complessiva",
		"Premio assicurativo lordo",
		"Importo Assicurato",
		"indennizzo mensile",
		"campo disponibile",
		"% di sovrappremio da applicare alla garanzia",
		"Codice Concessionario /dipendenti (iscr.E)",
		"Codice Banca",
		"Codice Campagna",
		"Copertura Assicurativa: Totale o Pro quota",
		"% assicurata dell assicurato ",
		"campo disponibile",
		"Maxi rata finale/Valore riscatto",
		"Stato occupazionale dell Assicurato",
		"Tipo aderente",
		"Canale di vendita ",
		"Tipo contraente / Contraente",
		"Denominazione Sociale o Cognome",
		"campo vuoto o nome",
		"Sesso",
		"Data di nascita",
		"Codice Fiscale ",
		"Indirizzo di residenza",
		"C.A.P. Di residenza",
		"Comune di residenza",
		"Provincia di residenza",
		"Indirizzo e-mail",
		"Numero di Cellulare",
		"Cognome Assicurato ",
		"Nome",
		"Sesso ",
		"Data di nascita ",
		"Codice Fiscale ",
		"Codice Fiscale Beneficiario 1",
		"Codice Fiscale Beneficiario 2",
		"Codice Fiscale Beneficiario 3",
		"Natura del rapporto ",
		"Scopo del rapporto ",
		"Modalità di pagamento del premio assicurativo (all intermediario) ",
		"contraente = Assicurato?",
		"Indirizzo di domicilio",
		"C.A.P. Di domicilio",
		"Comune di domicilio",
		"Provincia di domicilio",
		"Luogo di nascita dell’contraente persona fisica",
		"Provincia di nascita dell’contraente persona fisica",
		"Stato di residenza dell’contraente ",
		"Cab della città di residenza dell’contraente",
		"Sottogruppo attività economica",
		"Ramo gruppo attività economica",
		"Tipo documento dell contraente persona fisica ",
		"Numero documento dell contraente persona fisica ",
		"Data rilascio documento dell contraente persona fisica ",
		"Ente rilascio documento dell contraente persona fisica ",
		"PEP - Persona Politicamente Esposta",
		"Tipologia di PEP",
		"Modalità di comunicazione prescelta tra Compagnia ed contraente",
		"Indirizzo di residenza ",
		"C.A.P. Residenza ",
		"Comune Residenza ",
		"Provincia Residenza ",
		"Indirizzo di domicilio",
		"C.A.P. Domicilio",
		"Comune Domicilio",
		"Provincia Domicilio",
		"Indirizzo e-mail ",
		"Numero di cellulare ",
		"Luogo di nascita ",
		"Provincia di nascita ",
		"Stato di residenza",
		"Tipo documento",
		"Numero documento",
		"Data rilascio documento",
		"Ente rilascio documento",
		"PEP - Persona Politicamente Esposta",
		"Tipologia di PEP",
		"Eredi designati nominativamente o genericamente?",
		"Cognome Beneficiario 1",
		"Nome",
		"Codice Fiscale ",
		"Numero di Telefono del Beneficiario",
		"Indirizzo di residenza ",
		"Città /Comune di Residenza",
		"CAP",
		"Provincia",
		"Email ",
		"Legame del Cliente col Beneficiario",
		"L contraente ha escluso l invio di comunicazioni da parte dell Impresa al Beneficiario?",
		"Cognome Beneficiario 2",
		"Nome",
		"Codice Fiscale ",
		"Numero di Telefono del Beneficiario",
		"Indirizzo di residenza ",
		"Città /Comune di Residenza",
		"CAP",
		"Provincia",
		"Email ",
		"Legame del Cliente col Beneficiario",
		"L contraente ha escluso l invio di comunicazioni da parte dell Impresa al Beneficiario?",
		"Cognome Beneficiario 3",
		"Nome",
		"Codice Fiscale ",
		"Numero di Telefono del Beneficiario",
		"Indirizzo di residenza ",
		"Città /Comune di Residenza",
		"CAP",
		"Provincia",
		"Email ",
		"Legame del Cliente col Beneficiario",
		"L contraente ha escluso l invio di comunicazioni da parte dell Impresa al Beneficiario?",
		"Esistenza Titolare effettivo",
		"Cognome ",
		"Nome",
		"Sesso ",
		"Data di nascita ",
		"Codice Fiscale ",
		"Indirizzo di residenza ",
		"C.A.P. Residenza ",
		"Comune Residenza ",
		"Provincia Residenza ",
		"Indirizzo di domicilio",
		"C.A.P. Domicilio",
		"Comune Domicilio",
		"Provincia Domicilio",
		"Stato occupazionale ",
		"Indirizzo e-mail ",
		"Numero di cellulare ",
		"Luogo di nascita ",
		"Provincia di nascita ",
		"Stato di residenza",
		"Tipo documento",
		"Numero documento",
		"Data rilascio documento",
		"Ente rilascio documento",
		"PEP - Persona Politicamente Esposta",
		"Tipologia di PEP",
		"Esistenza Titolare effettivo",
		"Cognome ",
		"Nome",
		"Sesso ",
		"Data di nascita ",
		"Codice Fiscale ",
		"Indirizzo di residenza ",
		"C.A.P. Residenza ",
		"Comune Residenza ",
		"Provincia Residenza ",
		"Indirizzo di domicilio",
		"C.A.P. Domicilio",
		"Comune Domicilio",
		"Provincia Domicilio",
		"Stato occupazionale ",
		"Indirizzo e-mail ",
		"Numero di cellulare ",
		"Luogo di nascita ",
		"Provincia di nascita ",
		"Stato di residenza",
		"Tipo documento",
		"Numero documento",
		"Data rilascio documento",
		"Ente rilascio documento",
		"PEP - Persona Politicamente Esposta",
		"Tipologia di PEP",
		"Esistenza Titolare effettivo",
		"Cognome ",
		"Nome",
		"Sesso ",
		"Data di nascita ",
		"Codice Fiscale ",
		"Indirizzo di residenza ",
		"C.A.P. Residenza ",
		"Comune Residenza ",
		"Provincia Residenza ",
		"Indirizzo di domicilio",
		"C.A.P. Domicilio",
		"Comune Domicilio",
		"Provincia Domicilio",
		"Stato occupazionale ",
		"Indirizzo e-mail ",
		"Numero di cellulare ",
		"Luogo di nascita ",
		"Provincia di nascita ",
		"Stato di residenza",
		"Tipo documento",
		"Numero documento",
		"Data rilascio documento",
		"Ente rilascio documento",
		"PEP - Persona Politicamente Esposta",
		"Tipologia di PEP",
		"Esistenza Titolare effettivo",
		"Cognome ",
		"Nome",
		"Sesso ",
		"Data di nascita ",
		"Codice Fiscale ",
		"Indirizzo di residenza ",
		"C.A.P. Residenza ",
		"Comune Residenza ",
		"Provincia Residenza ",
		"Indirizzo di domicilio",
		"C.A.P. Domicilio",
		"Comune Domicilio",
		"Provincia Domicilio",
		"Stato occupazionale ",
		"Indirizzo e-mail ",
		"Numero di cellulare ",
		"Luogo di nascita ",
		"Provincia di nascita ",
		"Stato di residenza",
		"Tipo documento",
		"Numero documento",
		"Data rilascio documento",
		"Ente rilascio documento",
		"PEP - Persona Politicamente Esposta",
		"Tipologia di PEP",
		"Cognome ",
		"Nome",
		"Sesso",
		"Data di nascita",
		"Codice Fiscale ",
		"Indirizzo di residenza",
		"C.A.P. Di residenza",
		"Comune di residenza",
		"Provincia di residenza",
		"Indirizzo di domicilio",
		"C.A.P. Di domicilio",
		"Comune di domicilio",
		"Provincia di domicilio",
		"Indirizzo e-mail",
		"Numero di Cellulare",
		"Luogo di nascita dell’esecutore",
		"Provincia di nascita dell’esecutore",
		"Stato di residenza dell’esecutore",
		"Tipo documento",
		"Numero documento",
		"Data rilascio documento",
		"Ente rilascio documento",
		"PEP - Persona Politicamente Esposta",
		"Tipologia di PEP",
		"Carica ricoperta dell esecutore ",
		"Cognome ",
		"Nome",
		"Indirizzo di residenza ",
		"Città /Comune di Residenza",
		"CAP",
		"Codice Fiscale",
		"Numero di Telefono",
		"Email"}
}
