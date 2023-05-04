package companydata

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	lib "github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/xuri/excelize/v2"
)

func LifeAxalEmit(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	layout := "20060102"

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
		},
	}
	query, e := q.FirestoreWherefields("policy")
	policies := models.PolicyToListData(query)

	for _, policy := range policies {
		var (
			result [][]string
			//enterpriseName string
			//employer       string
			//revenue        string
		)
		for _, asset := range policy.Assets {
			if asset.Building != nil {

			}
			if asset.Enterprise != nil {
				//enterpriseName = asset.Enterprise.Name
				//employer = fmt.Sprint(asset.Enterprise.Employer)
				//revenue = fmt.Sprint(asset.Enterprise.Revenue)
			}
		}
		for _, asset := range policy.Assets {

			if asset.Building != nil {
				for _, g := range asset.Guarantees {
					var (
						beneficiary1, beneficiary2 string
					)
					beneficiary1, _ = mapBeneficiary(g, 0) //Codice Fiscale Beneficiario
					beneficiary2, _ = mapBeneficiary(g, 1)
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
						"0",                                   //indennizzo mensile
						"",                                    //campo disponibile
						"",                                    //% di sovrappremio da applicare alla garanzia
						"W1.42",                               //Codice Concessionario /dipendenti (iscr.E)
						"",                                    //Codice Banca
						"",                                    //Codice Campagna
						"T",                                   //Copertura Assicurativa: Totale o Pro quota
						"% assicurata dell'assicurato ",       //% assicurata dell'assicurato
						"campo disponibile",                   //campo disponibile
						"Maxi rata finale/Valore riscatto",    //Maxi rata finale/Valore riscatto
						"Stato occupazionale dell'Assicurato", //Stato occupazionale dell'Assicurato
						"Tasso di Interesse",                  //Tasso di Interesse
						"ONLINE",                              //Canale di vendita
						"PF",                                  //Tipo contraente / Contraente
						policy.Contractor.Surname,             //Denominazione Sociale o Cognome contraente
						policy.Contractor.Name,                //campo vuoto o nome
						policy.Contractor.Gender,              //Sesso
						policy.Contractor.BirthDate,           //Data di nascita
						policy.Contractor.FiscalCode,          //Codice Fiscale
						policy.Contractor.Address,             //Indirizzo di residenza
						policy.Contractor.PostalCode,          //C.A.P. Di residenza
						policy.Contractor.Locality,            //Comune di residenza
						policy.Contractor.City,                //Provincia di residenza
						policy.Contractor.Mail,                //Indirizzo e-mail
						policy.Contractor.Phone,               //Numero di Cellulare
						"",                                    //Cognome Assicurato
						"",                                    //Nome
						"",                                    //Sesso
						"",                                    //Data di nascita
						"",                                    //Codice Fiscale
						beneficiary1,                          //Codice Fiscale Beneficiario
						beneficiary2,                          //Codice Fiscale Beneficiario 2
						"",                                    //Codice Fiscale Beneficiario 3
						"VIT",                                 //Natura del rapporto
						"PAS ",                                //Scopo del rapporto
						"BO",                                  //Modalità di pagamento del premio assicurativo (all'intermediario)
						"SI",                                  //contraente = Assicurato?
						"",                                    //Indirizzo di domicilio
						"",                                    //C.A.P. Di domicilio
						"",                                    //Comune di domicilio
						"",                                    //Provincia di domicilio
						policy.Contractor.BirthCity,           //Luogo di nascita dell’contraente persona fisica
						policy.Contractor.BirthCity,           //Provincia di nascita dell’contraente persona fisica
						"086",                                 //Stato di residenza dell’contraente
						"Cab della città di residenza dell’contraente",            //Cab della città di residenza dell’contraente
						"Sottogruppo attività economica",                          //Sottogruppo attività economica
						"Ramo gruppo attività economica",                          //Ramo gruppo attività economica
						"Tipo documento dell'contraente persona fisica ",          //Tipo documento dell'contraente persona fisica
						"Numero documento dell'contraente persona fisica ",        //Numero documento dell'contraente persona fisica
						"Data rilascio documento dell'contraente persona fisica ", //Data rilascio documento dell'contraente persona fisica
						"Ente rilascio documento dell'contraente persona fisica ", //Ente rilascio documento dell'contraente persona fisica
						"NO",                                  //PEP - Persona Politicamente Esposta
						"",                                    //Tipologia di PEP
						"E",                                   //Modalità di comunicazione prescelta tra Compagnia ed contraente
						"Indirizzo di residenza ",             //Indirizzo di residenza
						"C.A.P. Residenza ",                   //C.A.P. Residenza
						"Comune Residenza ",                   //Comune Residenza
						"Provincia Residenza ",                //Provincia Residenza
						"Indirizzo di domicilio",              //Indirizzo di domicilio
						"C.A.P. Domicilio",                    //C.A.P. Domicilio
						"Comune Domicilio",                    //Comune Domicilio
						"Provincia Domicilio",                 //Provincia Domicilio
						"Indirizzo e-mail ",                   //Indirizzo e-mail
						"Numero di cellulare ",                //Numero di cellulare
						"Luogo di nascita ",                   //Luogo di nascita
						"Provincia di nascita ",               //Provincia di nascita
						"Stato di residenza",                  //Stato di residenza
						"Tipo documento",                      //Tipo documento
						"Numero documento",                    //Numero documento
						"Data rilascio documento",             //Data rilascio documento
						"Ente rilascio documento",             //Ente rilascio documento
						"PEP - Persona Politicamente Esposta", //PEP - Persona Politicamente Esposta
						"Tipologia di PEP",                    //Tipologia di PEP
						"Eredi designati nominativamente o genericamente?", //Eredi designati nominativamente o genericamente?
						"Cognome Beneficiario 1",                           //Cognome Beneficiario 1
						"Nome",                                             //Nome
						"Codice Fiscale ",                                  //Codice Fiscale
						"Numero di Telefono del Beneficiario",              //Numero di Telefono del Beneficiario
						"Indirizzo di residenza ",                          //Indirizzo di residenza
						"Città /Comune di Residenza",                       //Città /Comune di Residenza
						"CAP",                                              //CAP
						"Provincia",                                        //CAP
						"Email ",                                           //Email
						"Legame del Cliente col Beneficiario",              //Legame del Cliente col Beneficiario
						"NUCLEO FAMILIARE",                                 //NUCLEO FAMILIARE
						"Lcontraente ha escluso l invio di comunicazioni da parte dell Impresa al Beneficiario?", //Lcontraente ha escluso l invio di comunicazioni da parte dell Impresa al Beneficiario?
						"Cognome Beneficiario 2",              //Cognome Beneficiario 2
						"Nome",                                //Nome
						"Codice Fiscale ",                     //Codice Fiscale
						"Numero di Telefono del Beneficiario", //Numero di Telefono del Beneficiario
						"Indirizzo di residenza ",             //Indirizzo di residenza
						"Città /Comune di Residenza",          //Città /Comune di Residenza
						"CAP",                                 //CAP
						"Provincia",                           //Provincia
						"Email ",                              //Email
						"Legame del Cliente col Beneficiario", //Legame del Cliente col Beneficiario
						"NUCLEO FAMILIARE",                    //NUCLEO FAMILIARE
						"L contraente ha escluso l invio di comunicazioni da parte dell Impresa al Beneficiario?", //L contraente ha escluso l invio di comunicazioni da parte dell Impresa al Beneficiario?
						"Cognome Beneficiario 3",              //Cognome Beneficiario 3
						"Nome",                                //Nome
						"Codice Fiscale ",                     //Codice Fiscale
						"Numero di Telefono del Beneficiario", //Numero di Telefono del Beneficiario
						"Indirizzo di residenza ",             //Indirizzo di residenza
						"Città /Comune di Residenza",          //Città /Comune di Residenza
						"CAP",                                 //CAP
						"Provincia",                           //Provincia
						"Email ",                              //Email
						"Legame del Cliente col Beneficiario", //Legame del Cliente col Beneficiario
						"NUCLEO FAMILIARE",                    //NUCLEO FAMILIARE
						"L'contraente ha escluso l'invio di comunicazioni da parte dell Impresa al Beneficiario?", //L'contraente ha escluso l'invio di comunicazioni da parte dell Impresa al Beneficiario?
						"Esistenza Titolare effettivo",        //Esistenza Titolare effettiv
						"Cognome ",                            //Cognome
						"Nome",                                //Nome
						"Sesso ",                              //Sesso
						"Data di nascita ",                    //Data di nascita
						"Codice Fiscale ",                     //Codice Fiscale
						"Indirizzo di residenza ",             //Indirizzo di residenza
						"C.A.P. Residenza ",                   //C.A.P. Residenza
						"Comune Residenza ",                   //Comune Residenza
						"Provincia Residenza ",                //Provincia Residenza
						"Indirizzo di domicilio",              //Indirizzo di domicilio
						"C.A.P. Domicilio",                    //C.A.P. Domicilio
						"Comune Domicilio",                    //Comune Domicilio
						"Provincia Domicilio",                 //Provincia Domicilio
						"Stato occupazionale ",                //Stato occupazionale
						"Indirizzo e-mail ",                   //Indirizzo e-mail
						"Numero di cellulare ",                //Numero di cellulare
						"Luogo di nascita ",                   //Luogo di nascita
						"Provincia di nascita ",               //Provincia di nascita
						"Stato di residenza",                  //Stato di residenza
						"Tipo documento",                      //Tipo documento
						"Numero documento",                    //Numero documento
						"Data rilascio documento",             //Data rilascio documento
						"Ente rilascio documento",             //Ente rilascio documento
						"PEP - Persona Politicamente Esposta", //PEP - Persona Politicamente Esposta
						"Tipologia di PEP",                    //Tipologia di PEP
						"Esistenza Titolare effettivo",        //Esistenza Titolare effettivo
						"Cognome ",                            //Cognome
						"Nome",                                //Nome
						"Sesso ",                              //Sesso
						"Data di nascita ",                    //Data di nascita
						"Codice Fiscale ",                     //Codice Fiscale
						"Indirizzo di residenza ",             //Indirizzo di residenza
						"C.A.P. Residenza ",                   //C.A.P. Residenza
						"Comune Residenza ",                   //Comune Residenza
						"Provincia Residenza ",                //Provincia Residenza
						"Indirizzo di domicilio",              //Indirizzo di domicilio
						"C.A.P. Domicilio",                    //C.A.P. Domicilio
						"Comune Domicilio",                    //Comune Domicilio
						"Provincia Domicilio",                 //Provincia Domicilio
						"Stato occupazionale ",                //Stato occupazionale
						"Indirizzo e-mail ",                   //Indirizzo e-mail
						"Numero di cellulare ",                //Numero di cellulare
						"Luogo di nascita ",                   //Luogo di nascita
						"Provincia di nascita ",               //Provincia di nascita
						"Stato di residenza",                  //Stato di residenza
						"Tipo documento",                      //Tipo documento
						"Numero documento",                    //Numero documento
						"Data rilascio documento",             //Data rilascio documento
						"Ente rilascio documento",             //Ente rilascio documento
						"PEP - Persona Politicamente Esposta", //PEP - Persona Politicamente Esposta
						"Tipologia di PEP",                    //Tipologia di PEP
						"Esistenza Titolare effettivo",        //Esistenza Titolare effettivo
						"Cognome ",                            //Cognome
						"Nome",                                //Nome
						"Sesso ",                              //Sesso
						"Data di nascita ",                    //Data di nascita
						"Codice Fiscale ",                     //Codice Fiscale
						"Indirizzo di residenza ",             //Indirizzo di residenza
						"C.A.P. Residenza ",                   //C.A.P. Residenza
						"Comune Residenza ",                   //Comune Residenza
						"Provincia Residenza ",                //Provincia Residenza
						"Indirizzo di domicilio",              //Indirizzo di domicilio
						"C.A.P. Domicilio",                    //C.A.P. Domicilio
						"Comune Domicilio",                    //Comune Domicilio
						"Provincia Domicilio",                 //Provincia Domicilio
						"Stato occupazionale ",                //Stato occupazionale
						"Indirizzo e-mail ",                   //Indirizzo e-mail
						"Numero di cellulare ",                //Numero di cellulare
						"Luogo di nascita ",                   //Luogo di nascita
						"Provincia di nascita ",               //Provincia di nascita
						"Stato di residenza",                  //Stato di residenza
						"Tipo documento",                      //Tipo documento
						"Numero documento",                    //Numero documento
						"Data rilascio documento",             //Data rilascio documento
						"Ente rilascio documento",             //Ente rilascio documento
						"PEP - Persona Politicamente Esposta", //PEP - Persona Politicamente Esposta
						"Tipologia di PEP",                    //Tipologia di PEP
						"Esistenza Titolare effettivo",        //Esistenza Titolare effettivo
						"Cognome ",                            //Cognome
						"Nome",                                //Nome
						"Sesso ",                              //Sesso
						"Data di nascita ",                    //Data di nascita
						"Codice Fiscale ",                     //Codice Fiscale
						"Indirizzo di residenza ",             //Indirizzo di residenza
						"C.A.P. Residenza ",                   //C.A.P. Residenza
						"Comune Residenza ",                   //Comune Residenza
						"Provincia Residenza ",                //Provincia Residenza
						"Indirizzo di domicilio",              //Indirizzo di domicilio
						"C.A.P. Domicilio",                    //C.A.P. Domicilio
						"Comune Domicilio",                    //Comune Domicilio
						"Provincia Domicilio",                 //Provincia Domicilio
						"Stato occupazionale ",                //Stato occupazionale
						"Indirizzo e-mail ",                   //Indirizzo e-mail
						"Numero di cellulare ",                //umero di cellulare
						"Luogo di nascita ",                   //Luogo di nascita
						"Provincia di nascita ",               //Provincia di nascita
						"Stato di residenza",                  //Stato di residenza
						"Tipo documento",                      //Tipo documento"
						"Numero documento",                    //Numero documento
						"Data rilascio documento",             //Data rilascio documento
						"Ente rilascio documento",             //Ente rilascio documento
						"PEP - Persona Politicamente Esposta", //PEP - Persona Politicamente Esposta
						"Tipologia di PEP",                    //Tipologia di PEP
						"Cognome ",                            //Cognome
						"Nome",                                //Nome
						"Sesso",                               //Sesso
						"Data di nascita",                     //Data di nascita
						"Codice Fiscale ",                     //Codice Fiscale
						"Indirizzo di residenza",              //Indirizzo di residenza
						"C.A.P. Di residenza",                 //C.A.P. Di residenza
						"Comune di residenza",                 //Comune di residenza
						"Provincia di residenza",              //Provincia di residenza
						"Indirizzo di domicilio",              //Indirizzo di domicilio
						"C.A.P. Di domicilio",                 //C.A.P. Di domicilio
						"Comune di domicilio",                 //Comune di domicilio
						"Provincia di domicilio",              //Provincia di domicilio
						"Indirizzo e-mail",                    //Indirizzo e-mail
						"Numero di Cellulare",                 //Numero di Cellulare
						"Luogo di nascita dell’esecutore",     //Luogo di nascita dell’esecutore
						"Provincia di nascita dell’esecutore", //Provincia di nascita dell’esecutore
						"Stato di residenza dell’esecutore",   //Stato di residenza dell’esecutore
						"Tipo documento",                      //Tipo documento
						"Numero documento",                    //Numero documento
						"Data rilascio documento",             //Data rilascio documento"
						"Ente rilascio documento",             //
						"PEP - Persona Politicamente Esposta", //
						"Tipologia di PEP",                    //Tipologia di PEP
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

		}

		now := time.Now()
		layout2 := "2006-01-02"
		filepath := now.Format(layout2) + "-" + strconv.Itoa(time.Now().Nanosecond()) + ".xlsx"

		if os.Getenv("env") != "local" {
			//./serverless_function_source_code/

			//filepath = "../tmp/" + filepath
		}
		sourcest, e = CreateExcel(excel, "../tmp/"+filepath)
		lib.PutToStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "tway-fleet-axa/"+filepath, sourcest)
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
func mapBeneficiary(g models.Guarante, b int) (string, models.Beneficiary) {
	var (
		result      string
		resulStruct models.Beneficiary
	)
	if len(*g.Beneficiaries) > 0 && len(*g.Beneficiaries) >= b {
		result = ""
		if (*g.Beneficiaries)[b].IsLegitimateSuccessors || (*g.Beneficiaries)[b].IsFamilyMember {
			result = "GE"
		} else {
			result = (*g.Beneficiaries)[b].FiscalCode
			resulStruct = (*g.Beneficiaries)[b]
		}

	}

	return result, resulStruct
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
		"NUCLEO FAMILIARE",
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
		"NUCLEO FAMILIARE",
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
		"NUCLEO FAMILIARE",
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
	config := lib.SftpConfig{
		Username:     os.Getenv("AXA_LIFE_SFTP_USER"),
		Password:     "",                                                                                                          // required only if password authentication is to be used
		PrivateKey:   os.Getenv("AXA_LIFE_SFTP_PSW"),                                                                              //                           // required only if private key authentication is to be used
		Server:       os.Getenv("AXA_LIFE_SFTP_HOST"),                                                                             //
		KeyExchanges: []string{"diffie-hellman-group-exchange-sha1", "diffie-hellman-group1-sha1", "diffie-hellman-group14-sha1"}, // optional
		Timeout:      time.Second * 30,                                                                                            // 0 for not timeout
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
	//destination, e := client.Create(filePath)
	destination, e := client.Create("IN/" + filePath)
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

func CreateExcel(sheet [][]interface{}, filePath string) ([]byte, error) {
	log.Println("CreateExcel")
	f := excelize.NewFile()
	alfabet := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
	// Create a new sheet.
	index, err := f.NewSheet("Sheet1")
	for x, row := range sheet {
		for i, cel := range row {

			fmt.Println(cel)
			f.SetCellValue("Sheet1", alfabet[i]+""+strconv.Itoa(x+1), cel)
		}
	}
	//Set active sheet of the workbook.
	f.SetActiveSheet(index)

	//Save spreadsheet by the given path.
	err = f.SaveAs(filePath)

	resByte, err := f.WriteToBuffer()

	return resByte.Bytes(), err
}
