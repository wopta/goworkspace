package companydata

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	lib "github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func LifeAxalEmit(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	layout := "20060102"
	var (
		cabCsv []byte
		result [][]string
	)
	q := lib.Firequeries{
		Queries: []lib.Firequery{{
			Field:      "companyEmit", //
			Operator:   "==",          //
			QueryValue: true,
		},
			{
				Field:      "companyEmitted", //
				Operator:   "==",             //
				QueryValue: false,
			},
			{
				Field:      "company", //
				Operator:   "==",      //
				QueryValue: "axa",
			},
			{
				Field:      "name", //
				Operator:   "==",   //
				QueryValue: "life",
			},
		},
	}
	query, e := q.FirestoreWherefields("policy")
	policies := models.PolicyToListData(query)
	switch os.Getenv("env") {
	case "local":
		cabCsv = lib.ErrorByte(ioutil.ReadFile("function-data/data/rules/Riclassificazione_Ateco.csv"))

	default:
		cabCsv = lib.GetFromStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "data/cab-cap-istat.csv", "")
	}
	df := lib.CsvToDataframe(cabCsv)
	log.Println("df.Describe(): ", df.Describe())
	result = append(result, getHeader())
	for _, policy := range policies {
		fil := df.Filter(
			dataframe.F{Colidx: 4, Colname: "CAP", Comparator: series.Eq, Comparando: policy.Contractor.Residence.PostalCode},
		)
		residenceCab := fil.Records()[0][5]
		log.Println("policy.Uid: ", policy.Uid)
		log.Println("residenceCab", residenceCab)
		log.Println("filtered col", fil.Ncol())
		log.Println("filtered row", fil.Nrow())
		for _, asset := range policy.Assets {

			for _, g := range asset.Guarantees {
				var (
					beneficiary1, beneficiary2, beneficiary1T string
					beneficiary1S, beneficiary2S              models.Beneficiary
				)
				beneficiary1, beneficiary1S, beneficiary1T = mapBeneficiary(g, 0) //Codice Fiscale Beneficiario
				beneficiary2, beneficiary2S, _ = mapBeneficiary(g, 1)
				fmt.Println(g)
				row := []string{
					mapCodecCompany(policy, g.CompanyCodec), //Codice schema
					policy.CodeCompany,                      //N° adesione individuale univoco
					"A",                                     //Tipo di Transazione
					policy.StartDate.Format(layout),         //Data di decorrenza
					policy.EndDate.Format(layout),           //"Data di rinnovo"
					"012",                                   //"Durata copertura assicurativa"
					fmt.Sprint(g.Value.Duration.Year * 12),  //"Durata complessiva"
					fmt.Sprintf("%.2f", g.PriceGross),       //"Premio assicurativo lordo"
					fmt.Sprintf("%.0f", g.SumInsuredLimitOfIndemnity), //"Importo Assicurato"
					"0",                                //indennizzo mensile
					"",                                 //campo disponibile
					"",                                 //% di sovrappremio da applicare alla garanzia
					"W1",                               //Codice Concessionario /dipendenti (iscr.E)
					"",                                 //Codice Banca
					"",                                 //Codice Campagna
					"T",                                //Copertura Assicurativa: Totale o Pro quota
					"",                                 //% assicurata dell'assicurato
					"",                                 //campo disponibile
					"",                                 //Maxi rata finale/Valore riscatto
					"",                                 //Stato occupazionale dell'Assicurato
					"2",                                //Tipo aderente
					"WEB",                              //Canale di vendita
					"PF",                               //Tipo contraente / Contraente
					policy.Contractor.Surname,          //Denominazione Sociale o Cognome contraente
					policy.Contractor.Name,             //campo vuoto o nome
					policy.Contractor.Gender,           //Sesso
					policy.Contractor.BirthDate,        //Data di nascita
					policy.Contractor.FiscalCode,       //Codice Fiscale
					policy.Contractor.Address,          //Indirizzo di residenza
					policy.Contractor.PostalCode,       //C.A.P. Di residenza
					policy.Contractor.Locality,         //Comune di residenza
					policy.Contractor.City,             //Provincia di residenza
					policy.Contractor.Mail,             //Indirizzo e-mail
					policy.Contractor.Phone,            //Numero di Cellulare
					policy.Assets[0].Person.Surname,    //Cognome Assicurato
					policy.Assets[0].Person.Name,       //Nome
					policy.Assets[0].Person.Gender,     //Sesso
					policy.Assets[0].Person.BirthDate,  //Data di nascita
					policy.Assets[0].Person.FiscalCode, //Codice Fiscale
					beneficiary1,                       //Codice Fiscale Beneficiario
					beneficiary2,                       //Codice Fiscale Beneficiario 2
					"",                                 //Codice Fiscale Beneficiario 3
					"VIT",                              //Natura del rapporto
					"PAS ",                             //Scopo del rapporto
					"BO",                               //Modalità di pagamento del premio assicurativo (all'intermediario)
					"SI",                               //contraente = Assicurato?
					ChekDomicilie(policy.Contractor).StreetName, //Indirizzo di domicilio contraente
					ChekDomicilie(policy.Contractor).PostalCode, //C.A.P. Di domicilio
					ChekDomicilie(policy.Contractor).Locality,   //Comune di domicilio
					ChekDomicilie(policy.Contractor).CityCode,   //Provincia di domicilio
					policy.Contractor.BirthCity,                 //Luogo di nascita dell’contraente persona fisica
					policy.Contractor.BirthCity,                 //Provincia di nascita dell’contraente persona fisica
					"086",                                       //Stato di residenza dell’contraente
					residenceCab,                                //Cab della città di residenza dell’contraente
					"600",                                       //Sottogruppo attività economica
					"600",                                       //Ramo gruppo attività economica
					ExistIdentityDocument(policy.Contractor.IdentityDocuments).Code,                       //Tipo documento dell'contraente persona fisica
					ExistIdentityDocument(policy.Contractor.IdentityDocuments).Number,                     //Numero documento dell'contraente persona fisica
					ExistIdentityDocument(policy.Contractor.IdentityDocuments).DateOfIssue.Format(layout), //Data rilascio documento dell'contraente persona fisica
					ExistIdentityDocument(policy.Contractor.IdentityDocuments).IssuingAuthority,           //Ente rilascio documento dell'contraente persona fisica
					"NO", //PEP - Persona Politicamente Esposta
					"",   //Tipologia di PEP
					"E",  //Modalità di comunicazione prescelta tra Compagnia ed contraente
					policy.Assets[0].Person.Residence.StreetName,                            //Indirizzo di residenza Assicurato
					policy.Assets[0].Person.Residence.PostalCode,                            //C.A.P. Residenza
					policy.Assets[0].Person.Residence.Locality,                              //Comune Residenza
					policy.Assets[0].Person.Residence.CityCode,                              //Provincia Residenza
					policy.Assets[0].Person.Domicile.StreetName,                             //Indirizzo di domicilio
					policy.Assets[0].Person.Domicile.PostalCode,                             //C.A.P. Domicilio
					policy.Assets[0].Person.Domicile.Locality,                               //Comune Domicilio
					policy.Assets[0].Person.Domicile.CityCode,                               //Provincia Domicilio
					policy.Assets[0].Person.Mail,                                            //Indirizzo e-mail
					policy.Assets[0].Person.Phone,                                           //Numero di cellulare
					policy.Assets[0].Person.BirthCity,                                       //Luogo di nascita
					policy.Assets[0].Person.BirthCity,                                       //Provincia di nascita
					"ITA",                                                                   //Stato di residenza
					ExistIdentityDocument(policy.Assets[0].Person.IdentityDocuments).Code,   //Tipo documento
					ExistIdentityDocument(policy.Assets[0].Person.IdentityDocuments).Number, //Numero documento
					ExistIdentityDocument(policy.Assets[0].Person.IdentityDocuments).DateOfIssue.Format(layout), //Data rilascio documento
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

		now := time.Now()
		refMontly := now.AddDate(0, -1, 0)
		//year, month, day := time.Now().Date()
		//year2, month2, day2 := time.Now().AddDate(0, -1, 0).Date()
		filepath := "WOPTAKEY_NBM_" + strconv.Itoa(refMontly.Year()) + strconv.Itoa(int(refMontly.Month())) + "_" + strconv.Itoa(now.Day()) + strconv.Itoa(int(now.Month())) + ".txt"
		lib.WriteCsv("../tmp/"+filepath, result)
		source, _ := ioutil.ReadFile("../tmp/" + filepath)
		lib.PutToStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "axa/life/"+filepath, source)
		SftpUpload(filepath)
	}
	return "", nil, e
}
func mapCodecCompany(p models.Policy, g string) string {
	var result, pay string

	if p.PaymentSplit == "year" {
		pay = "W"
	}
	if p.PaymentSplit == "montly" {
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
	log.Println(reflect.ValueOf(u.Domicile))
	if reflect.ValueOf(u.Domicile).IsNil() {
		res = *u.Residence
	}
	return res
}
func CheckStructNil[T interface{}](s interface{}) T {
	var result T
	result1 := new(T)
	result = *result1
	log.Println(reflect.TypeOf(s))
	if reflect.TypeOf(s) != nil {
		log.Println("is not nill")
		result = s.(T)
	}
	log.Println(s)
	log.Println(result)
	return result
}
func mapBeneficiary(g models.Guarante, b int) (string, models.Beneficiary, string) {
	var (
		result      string
		result2     string
		resulStruct models.Beneficiary
	)
	resulStruct = models.Beneficiary{}
	if g.Beneficiaries != nil {
		if len(*g.Beneficiaries) > 0 && len(*g.Beneficiaries) > b {
			result = ""
			if (*g.Beneficiaries)[b].IsLegitimateSuccessors || (*g.Beneficiaries)[b].IsFamilyMember {
				result = "GE"
				result2 = "GE"
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
func ExistIdentityDocument(docs []*models.IdentityDocument) models.IdentityDocument {
	var (
		result models.IdentityDocument
	)
	result = models.IdentityDocument{}
	if len(docs) > 0 {
		for _, doc := range docs {
			log.Println(doc)
			//doc.DateOfIssue

		}

	}
	return result
}

func getHeader() []string {
	return []string{"Codice schema",
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
		"Tasso di Interesse",
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
func SftpUpload(filePath string) {

	pk := lib.GetFromStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "env/axa-life.ppk", "")
	config := lib.SftpConfig{
		Username:     os.Getenv("AXA_LIFE_SFTP_USER"),
		Password:     "",                                                                                                          // required only if password authentication is to be used
		PrivateKey:   string(pk),                                                                                                  //                           // required only if private key authentication is to be used
		Server:       os.Getenv("AXA_LIFE_SFTP_HOST") + ":10026",                                                                  //
		KeyExchanges: []string{"diffie-hellman-group-exchange-sha1", "diffie-hellman-group1-sha1", "diffie-hellman-group14-sha1"}, // optional
		Timeout:      time.Second * 30,
		KeyPsw:       "", // 0 for not timeout
	}
	client, e := lib.NewSftpclient(config)
	lib.CheckError(e)
	defer client.Close()
	log.Println("Open local file for reading.:")
	source, e := os.Open("../tmp/" + filePath)
	lib.CheckError(e)
	//defer source.Close()
	log.Println("Create remote file for writing:")
	// Create remote file for writing.
	lib.Files("../tmp")
	destination, e := client.Create("To_CLP/" + filePath)
	lib.CheckError(e)
	defer destination.Close()
	log.Println("Upload local file to a remote location as in 1MB (byte) chunks.")
	info, e := source.Stat()
	log.Println(info.Size())
	// Upload local file to a remote location as in 1MB (byte) chunks.
	e = client.Upload(source, destination, int(info.Size()))
	lib.CheckError(e)
	/*
		// Download remote file.
		file, err := client.Download("tmp/file.txt")
		if err != nil {
			log.Fatalln(err)
		}
		defer file.Close()

		// Read downloaded file.
		data, err := ioutil.ReadAll(file)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(string(data))

		// Get remote file stats.
		info, err := client.Info("tmp/file.txt")
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("%+v\n", info)*/
}
