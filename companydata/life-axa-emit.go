package companydata

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"cloud.google.com/go/firestore"
	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	lib "gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
	"google.golang.org/api/iterator"
)

func LifeAxaEmit(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		from           time.Time
		to             time.Time
		filenamesplit  string
		cabCsv         []byte
		result         [][]string
		refMontly      time.Time
		upload         bool
		queryListdate  []string
		queryListdate2 []string
	)

	log.Println("----------------LifeAxalEmit-----------------")
	req := lib.ErrorByte(io.ReadAll(r.Body))
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

	i := from
	count := 0
	for next := true; next; next = i.Before(to) {
		if count > 15 {
			queryListdate2 = append(queryListdate2, i.Format("2006-01-02"))
		} else {
			queryListdate = append(queryListdate, i.Format("2006-01-02"))
		}
		fmt.Println("LifeAxalEmit: ", i.Format("2006-01-02"))
		//2023-09-26
		i = i.AddDate(0, 0, 1)
		count++

	}
	fmt.Println("LifeAxalEmit: ", i.Format("2006-01-02"))
	queryListdate = append(queryListdate, i.Format("2006-01-02"))
	fmt.Println("LifeAxalEmit:queryListdate ", queryListdate)
	fmt.Println("LifeAxalEmit:queryListdate ", queryListdate2)

	lifeAxaEmitQuery := lib.Firequeries{
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
			},
			{
				Field:      "scheduleDate", //
				Operator:   "in",           //
				QueryValue: queryListdate,
			},
		},
	}
	lifeAxaEmitQuery2 := lib.Firequeries{
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
			},
			{
				Field:      "scheduleDate", //
				Operator:   "in",           //
				QueryValue: queryListdate2,
			},
		},
	}
	df := lib.CsvToDataframe(cabCsv)
	query, e := lifeAxaEmitQuery.FirestoreWherefields("transactions")
	log.Println("LifeAxalEmit error: ", e)
	query2, e := lifeAxaEmitQuery2.FirestoreWherefields("transactions")
	log.Println("LifeAxalEmit error: ", e)

	transactions := TransactionToListData(query)
	transactions2 := TransactionToListData(query2)
	transactionstot := append(transactions, transactions2...)
	log.Println("LifeAxalEmit: transaction len: ", len(transactions))
	//result = append(result, getHeader())
	for _, transaction := range transactionstot {
		var (
			policy models.Policy
		)

		docsnap := lib.GetFirestore("policy", transaction.PolicyUid)
		docsnap.DataTo(&policy)
		if policy.Contractor.Type != "legalEntity" {
			result = append(result, setRowLifeEmit(policy, df, transaction, now)...)
			transaction.IsEmit = true
			lib.SetFirestore("transactions", transaction.Uid, transaction)
		}

	}

	filepath := "WOPTAKEYweb_NB" + filenamesplit + "_" + strconv.Itoa(refMontly.Year()) + fmt.Sprintf("%02d", int(refMontly.Month())) + "_" + fmt.Sprintf("%02d", now.Day()) + fmt.Sprintf("%02d", int(now.Month())) + ".txt"
	lib.WriteCsv("../tmp/"+filepath, result, ';')
	source, _ := os.ReadFile("../tmp/" + filepath)
	lib.PutToStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "track/axa/life/"+strconv.Itoa(refMontly.Year())+"/"+filepath, source)
	if upload {
		AxaPartnersSftpUpload(filepath)
	}
	return "", nil, e
}
func mapContractorTypeAxaPLife(policy models.Policy) (models.Contractor, string) {
	contractor := policy.Contractor
	typeContractorAxa := ""
	if policy.Contractor.Type == "legalEntity" {
		typeContractorAxa = "PG"
	} else {
		typeContractorAxa = "PF"
	}
	return contractor, typeContractorAxa
}
func fixImportedId(policy models.Policy) string {
	result := policy.CodeCompany
	if lib.SliceContains[string](policy.StatusHistory, "Imported") {
		result = "0000" + policy.CodeCompany
	}
	return result
}
func setRowLifeEmit(policy models.Policy, df dataframe.DataFrame, trans models.Transaction, now time.Time) [][]string {
	var (
		result               [][]string
		residenceCab         string
		networkCode, channel string
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
	if policy.ProducerUid != "" && policy.Channel == "network" {
		var node *models.NetworkNode
		snap, _ := lib.GetFirestoreErr("networkNodes", policy.ProducerUid)
		snap.DataTo(&node)
		log.Println("LifeAxalEmit: node.Code", node.Code)
		networkCode = node.Code
		channel = "POS"
	} else {
		channel = "WEB"
		networkCode = "W1"
	}
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
			priceGrossFormat := fmt.Sprintf("%013d", intNum) // 000000001220
			log.Println("LifeAxalEmit: ", priceGrossFormat)
			user := policy.Contractor.ToUser()
			if user == nil {
				return nil
			}

			r, m := getIndennity(g)
			if e != nil {
				log.Error(e)
			}
			row := []string{
				mapCodecCompany(policy, g.CompanyCodec),         //Codice schema
				fixImportedId(policy),                           //N° adesione individuale univoco
				getTypeTransactionAR(policy, trans),             //Tipo di Transazione
				getFormatdate(policy.StartDate),                 //Data di decorrenza
				getFormatdate(getRenewDate(policy, trans, now)), //"Data di rinnovo"
				mapCoverageDuration(policy),                     //"Durata copertura assicurativa"
				fmt.Sprint(g.Value.Duration.Year * 12),          //"Durata complessiva"
				priceGrossFormat,                                //"Premio assicurativo lordo"
				r,                                               //"Importo Assicurato"
				m,                                               //indennizzo mensile
				"",                                              //campo disponibile
				"",                                              //% di sovrappremio da applicare alla garanzia
				networkCode,                                     //Codice Concessionario /dipendenti (iscr.E)
				"",                                              //Codice Campagna
				"T",                                             //Copertura Assicurativa: Totale o Pro quota
				"",                                              //% assicurata dell'assicurato
				"",                                              //campo disponibile
				"",                                              //Maxi rata finale/Valore riscatto
				"",                                              //Stato occupazionale dell'Assicurato
				"1",                                             //Tipo aderente
				channel,                                         //Canale di vendita
				"PF",                                            //Tipo contraente / Contraente
				policy.Contractor.Surname,                       //Denominazione Sociale o Cognome contraente
				policy.Contractor.Name,                          //campo vuoto o nome
				policy.Contractor.Gender,                        //Sesso
				getFormatBithdate(policy.Contractor.BirthDate),                                           //Data di nascita
				policy.Contractor.FiscalCode,                                                             //Codice Fiscale
				policy.Contractor.Residence.StreetName + ", " + policy.Contractor.Residence.StreetNumber, //Indirizzo di residenza
				policy.Contractor.Residence.PostalCode,                                                   //C.A.P. Di residenza
				policy.Contractor.Residence.Locality,                                                     //Comune di residenza
				policy.Contractor.Residence.CityCode,                                                     //Provincia di residenza
				policy.Contractor.Mail,                                                                   //Indirizzo e-mail
				policy.Contractor.Phone,                                                                  //Numero di Cellulare
				policy.Assets[0].Person.Surname,                                                          //Cognome Assicurato
				policy.Assets[0].Person.Name,                                                             //Nome
				policy.Assets[0].Person.Gender,                                                           //Sesso
				getFormatBithdate(policy.Assets[0].Person.BirthDate),                                     //Data di nascita
				policy.Assets[0].Person.FiscalCode,                                                       //Codice Fiscale
				beneficiary1,                                                                             //Codice Fiscale Beneficiario
				beneficiary2,                                                                             //Codice Fiscale Beneficiario 2
				"",                                                                                       //Codice Fiscale Beneficiario 3
				"VIT",                                                                                    //Natura del rapporto
				"PAS",                                                                                    //Scopo del rapporto
				"BO",                                                                                     //Modalità di pagamento del premio assicurativo (all'intermediario)
				"SI",                                                                                     //contraente = Assicurato?
				ChekDomicilie(*user).StreetName + ", " + ChekDomicilie(*user).StreetNumber, //Indirizzo di domicilio contraente
				ChekDomicilie(*user).PostalCode,                                            //C.A.P. Di domicilio
				ChekDomicilie(*user).Locality,                                              //Comune di domicilio
				ChekDomicilie(*user).CityCode,                                              //Provincia di domicilio
				policy.Contractor.BirthCity,                                                //Luogo di nascita dell’contraente persona fisica
				policy.Contractor.BirthProvince,                                            //Provincia di nascita dell’contraente persona fisica
				"086",                                                                      //Stato di residenza dell’contraente
				residenceCab,                                                               //Cab della città di residenza dell’contraente
				"600",                                                                      //Sottogruppo attività economica
				"600",                                                                      //Ramo gruppo attività economica
				ExistIdentityDocument(policy.Contractor.IdentityDocuments).Code,                                                                                                //Tipo documento dell'contraente persona fisica
				ExistIdentityDocument(policy.Contractor.IdentityDocuments).Number,                                                                                              //Numero documento dell'contraente persona fisica
				getFormatdate(ExistIdentityDocument(policy.Contractor.IdentityDocuments).DateOfIssue),                                                                          //Data rilascio documento dell'contraente persona fisica
				ExistIdentityDocument(policy.Contractor.IdentityDocuments).IssuingAuthority + " di " + ExistIdentityDocument(policy.Contractor.IdentityDocuments).PlaceOfIssue, //Ente rilascio documento dell'contraente persona fisica
				"NO", //PEP - Persona Politicamente Esposta
				"",   //Tipologia di PEP
				"E",  //Modalità di comunicazione prescelta tra Compagnia ed contraente
				policy.Assets[0].Person.Residence.StreetName + ", " + policy.Assets[0].Person.Residence.StreetNumber, //Indirizzo di residenza Assicurato
				policy.Assets[0].Person.Residence.PostalCode,                                                         //C.A.P. Residenza
				policy.Assets[0].Person.Residence.Locality,                                                           //Comune Residenza
				policy.Assets[0].Person.Residence.CityCode,                                                           //Provincia Residenza
				ChekDomicilie(*policy.Assets[0].Person).StreetName,                                                   //Indirizzo di domicilio
				ChekDomicilie(*policy.Assets[0].Person).PostalCode,                                                   //C.A.P. Domicilio
				ChekDomicilie(*policy.Assets[0].Person).Locality,                                                     //Comune Domicilio
				ChekDomicilie(*policy.Assets[0].Person).CityCode,                                                     //Provincia Domicilio
				policy.Assets[0].Person.Mail,                                                                         //Indirizzo e-mail
				policy.Assets[0].Person.Phone,                                                                        //Numero di cellulare
				policy.Assets[0].Person.BirthCity,                                                                    //Luogo di nascita
				policy.Assets[0].Person.BirthCity,                                                                    //Provincia di nascita
				"ITA",                                                                                                //Stato di residenza
				ExistIdentityDocument(policy.Assets[0].Person.IdentityDocuments).Code,                                //Tipo documento
				ExistIdentityDocument(policy.Assets[0].Person.IdentityDocuments).Number,                              //Numero documento
				getFormatdate(ExistIdentityDocument(policy.Assets[0].Person.IdentityDocuments).DateOfIssue),          //Data rilascio documento
				ExistIdentityDocument(policy.Assets[0].Person.IdentityDocuments).IssuingAuthority + " di " + ExistIdentityDocument(policy.Assets[0].Person.IdentityDocuments).PlaceOfIssue, //Ente rilascio documento
				"NO",                     //PEP - Persona Politicamente Esposta
				"",                       //Tipologia di PEP
				beneficiary1T,            //Eredi designati nominativamente o genericamente?
				beneficiary1S.Surname,    //Cognome Beneficiario 1
				beneficiary1S.Name,       //Nome
				beneficiary1S.FiscalCode, //Codice Fiscale
				beneficiary1S.Phone,      //Numero di Telefono del Beneficiario
				beneficiary1S.Residence.StreetName + ", " + beneficiary1S.Residence.StreetNumber, //Indirizzo di residenza
				beneficiary1S.Residence.City,          //Città /Comune di Residenza
				beneficiary1S.Residence.PostalCode,    //CAP
				beneficiary1S.Residence.CityCode,      //Provincia
				beneficiary1S.Mail,                    //Email
				MapBool(beneficiary1S.IsFamilyMember), //Legame del Cliente col Beneficiario

				MapBool(beneficiary1S.IsContactable), //Lcontraente ha escluso l invio di comunicazioni da parte dell Impresa al Beneficiario?
				beneficiary2S.Surname,                //Cognome Beneficiario 2
				beneficiary2S.Name,                   //Nome
				beneficiary2S.FiscalCode,             //Codice Fiscale
				beneficiary2S.Phone,                  //Numero di Telefono del Beneficiario
				beneficiary2S.Residence.StreetName + ", " + beneficiary2S.Residence.StreetNumber, //Indirizzo di residenza
				beneficiary2S.Residence.City,          //Città /Comune di Residenza
				beneficiary2S.Residence.PostalCode,    //CAP
				beneficiary2S.Residence.CityCode,      //Provincia
				beneficiary2S.Mail,                    //Email
				MapBool(beneficiary2S.IsFamilyMember), //Legame del Cliente col Beneficiario
				MapBool(beneficiary2S.IsContactable),  //L contraente ha escluso l invio di comunicazioni da parte dell Impresa al Beneficiario?
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
func getFormatIndennity(i float64) string {
	var res string
	sumInsuredLimitRound := fmt.Sprintf("%.0f", i)
	sumInsuredLimit, _ := strconv.Atoi(sumInsuredLimitRound)
	res = fmt.Sprintf("%013d", sumInsuredLimit)
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

func getRenew(p models.Policy, now time.Time) string {
	var result string
	addMonth := p.StartDate.AddDate(0, 1, 0)
	if p.PaymentSplit == string(models.PaySplitYear) || p.PaymentSplit == string(models.PaySplitYearly) {

		result = "A"
	}

	if p.PaymentSplit == string(models.PaySplitMonthly) {
		log.Println("LifeAxalEmit : getRenew CodeCompany", p.CodeCompany)
		log.Println("LifeAxalEmit : getRenew addMonth", addMonth)
		log.Println("LifeAxalEmit : getRenew now", now)

		if now.Before(addMonth) {

			result = "A"

		} else {

			result = "R"
		}
	}
	log.Println("LifeAxalEmit : getRenew result", result)
	return result
}
func getTypeTransactionAR(p models.Policy, tr models.Transaction) string {
	var result string
	if p.PaymentSplit == string(models.PaySplitYear) || p.PaymentSplit == string(models.PaySplitYearly) {

		result = "A"
	}
	if lib.SliceContains(p.StatusHistory, "Renewed") {

		result = "R"
	}

	if p.PaymentSplit == string(models.PaySplitMonthly) {
		trdate, e := time.Parse("2006-01-02", tr.ScheduleDate)
		trMounthInt := int(trdate.Month())
		trYearInt := int(trdate.Year())
		log.Println("LifeAxalEmit : getTypeTransaction error", e)
		policyStartMounth := int(p.StartDate.Month())
		policyStartYear := int(p.StartDate.Year())
		log.Println("LifeAxalEmit : getTypeTransaction policyStartMounth", policyStartMounth)
		log.Println("LifeAxalEmit : getTypeTransaction trMounthInt", trMounthInt)
		log.Println("LifeAxalEmit : getTypeTransaction policyStartYear == trYearInt", policyStartYear, " - ", trYearInt)
		log.Println("LifeAxalEmit : getTypeTransaction lib.SliceContains", !lib.SliceContains(p.StatusHistory, "Renewed"))

		if policyStartMounth == trMounthInt && policyStartYear == trYearInt && !lib.SliceContains(p.StatusHistory, "Renewed") {

			result = "A"

		} else {

			result = "R"
		}
	}
	log.Println("LifeAxalEmit : getTypeTransactionresult", result)
	return result
}
func getRenewDate(p models.Policy, trans models.Transaction, now time.Time) time.Time {
	var result time.Time

	if p.PaymentSplit == "year" || p.PaymentSplit == string(models.PaySplitYearly) {
		result = p.StartDate
	}
	if p.PaymentSplit == string(models.PaySplitMonthly) {
		trdate, e := time.Parse("2006-01-02", trans.ScheduleDate)
		trMounthInt := int(trdate.Month())
		log.Println("LifeAxalEmit : getTypeTransaction error", e)
		policyStartMounth := int(p.StartDate.Month())
		log.Println("LifeAxalEmit : getTypeTransaction policyStartMounth", policyStartMounth)
		log.Println("LifeAxalEmit : getTypeTransaction trMounthInt", trMounthInt)
		if policyStartMounth == trMounthInt {

			result = p.StartDate
		} else {

			result = trdate
		}
	}
	return result
}
func mapCoverageDuration(p models.Policy) string {
	var result string
	if p.PaymentSplit == "year" || p.PaymentSplit == string(models.PaySplitYearly) {
		result = "012"
	}
	if p.PaymentSplit == string(models.PaySplitMonthly) {
		result = "001"
	}
	return result
}
func mapCodecCompany(p models.Policy, g string) string {
	var result, pay string

	if p.PaymentSplit == "year" || p.PaymentSplit == string(models.PaySplitYearly) {
		pay = "W"
	}
	if p.PaymentSplit == string(models.PaySplitMonthly) {
		pay = "M"
	}
	if g == "D" {
		result = p.ProductVersion[1:] + pay + "5"
	}
	if g == "PTD" {
		result = p.ProductVersion[1:] + pay + "6"
	}
	if g == "TTD" {
		result = p.ProductVersion[1:] + pay + "7"
	}
	if g == "CI" {
		result = p.ProductVersion[1:] + pay + "8"
	}
	return result
}
func getIndennity(g models.Guarante) (string, string) {
	var result, monthly string
	sumInsuredLimitRound := g.Value.SumInsuredLimitOfIndemnity * 100
	sumInsuredLimit := int(sumInsuredLimitRound)

	if g.CompanyCodec == "TTD" {
		monthly = fmt.Sprintf("%013d", sumInsuredLimit)
	} else {
		result = fmt.Sprintf("%013d", sumInsuredLimit)
	}

	return result, monthly
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
	resulStructDefault := models.Beneficiary{Residence: &models.Address{}}
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

			result = doc

		}

	}
	return result
}

func TransactionToListData(query *firestore.DocumentIterator) []models.Transaction {
	result := make([]models.Transaction, 0)
	log.Println("TransactionToListDatam start")
	for {
		d, err := query.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}
		}
		var value models.Transaction
		log.Println("TransactionToListData ref id:", d.Ref.ID)
		e := d.DataTo(&value)
		value.Uid = d.Ref.ID
		lib.CheckError(e)
		result = append(result, value)

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
