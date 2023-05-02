package companydata

import (
	"fmt"
	"net/http"
	"os"
	"time"

	lib "github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func LifeAxalEmit(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	config := lib.SftpConfig{
		Username:     os.Getenv("AXA_LIFE_SFTP_USER"),
		Password:     "",                                                                                                          // required only if password authentication is to be used
		PrivateKey:   os.Getenv("AXA_LIFE_SFTP_PSW"),                                                                              //                           // required only if private key authentication is to be used
		Server:       os.Getenv("AXA_LIFE_SFTP_HOST"),                                                                             //
		KeyExchanges: []string{"diffie-hellman-group-exchange-sha1", "diffie-hellman-group1-sha1", "diffie-hellman-group14-sha1"}, // optional
		Timeout:      time.Second * 30,                                                                                            // 0 for not timeout
	}
	client, e := lib.NewSftpclient(config)
	reader, e := client.Download("wopta/")
	lib.ExcelRead(reader)
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

		//startDate := lib.Dateformat(policy.StartDate)
		//endDate := lib.Dateformat(policy.EndDate)
		//crationDate := lib.Dateformat(policy.CreationDate)
		//companyCode := policy.CodeCompany
		//cityCode := policy.Contractor.CityCode
		//city := policy.Contractor.City
		//streetNumber := policy.Contractor.StreetNumber
		//postalCode := policy.Contractor.PostalCode
		//pi := policy.Contractor.VatCode
		//fc := policy.Contractor.FiscalCode
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
					fmt.Println(g)
					row := []string{
						"Codice schema",                                //Codice schema
						"N° adesione individuale univoco",              //N° adesione individuale univoco
						"Tipo di Transazione",                          //Tipo di Transazione
						"Data di decorrenza",                           //Data di decorrenza
						"Data di rinnovo",                              //"Data di rinnovo"
						"Durata copertura assicurativa",                //"Durata copertura assicurativa"
						"Durata complessiva",                           //"Durata complessiva"
						"Premio assicurativo lordo",                    //
						"Importo Assicurato",                           //"Importo Assicurato"
						"indennizzo mensile",                           //indennizzo mensile
						"campo disponibile",                            //campo disponibile
						"% di sovrappremio da applicare alla garanzia", //% di sovrappremio da applicare alla garanzia
						"Codice Concessionario /dipendenti (iscr.E)",   //Codice Concessionario /dipendenti (iscr.E)
						"Codice Banca",                                 //Codice Banca
						"Codice Campagna",                              //Codice Campagna
						"Copertura Assicurativa: Totale o Pro quota",   //Copertura Assicurativa: Totale o Pro quota
						"% assicurata dell'assicurato ",                //% assicurata dell'assicurato
						"campo disponibile",                            //campo disponibile
						"Maxi rata finale/Valore riscatto",             //Maxi rata finale/Valore riscatto
						"Stato occupazionale dell'Assicurato",          //Stato occupazionale dell'Assicurato
						"Tasso di Interesse",                           //Tasso di Interesse
						"Canale di vendita ",                           //Canale di vendita
						"Tipo contraente / Contraente",                 //Tipo contraente / Contraente
						"Denominazione Sociale o Cognome",              //Denominazione Sociale o Cognome
						"campo vuoto o nome",                           //campo vuoto o nome
						"Sesso",                                        //Sesso
						"Data di nascita",                              //Data di nascita
						"Codice Fiscale ",                              //Codice Fiscale
						"Indirizzo di residenza",                       //Indirizzo di residenza
						"C.A.P. Di residenza",                          //C.A.P. Di residenza
						"Comune di residenza",                          //Comune di residenza
						"Provincia di residenza",                       //Provincia di residenza
						"Indirizzo e-mail",                             //Indirizzo e-mail
						"Numero di Cellulare",                          //Numero di Cellulare
						"Cognome Assicurato ",                          //Cognome Assicurato
						"Nome",                                         //Nome
						"Sesso ",                                       //Sesso
						"Data di nascita ",                             //Data di nascita
						"Codice Fiscale ",                              //Codice Fiscale
						"Codice Fiscale Beneficiario 1",                //Codice Fiscale Beneficiario
						"Codice Fiscale Beneficiario 2",                //Codice Fiscale Beneficiario 2
						"Codice Fiscale Beneficiario 3",                //Codice Fiscale Beneficiario 3
						"AML DATA",                                     //AML DATA
						"Natura del rapporto ",                         //Natura del rapporto
						"Scopo del rapporto ",                          //Scopo del rapporto
						"Modalità di pagamento del premio assicurativo (all'intermediario) ", //Modalità di pagamento del premio assicurativo (all'intermediario)
						"contraente = Assicurato?",                                        //contraente = Assicurato?
						"Indirizzo di domicilio",                                          //Indirizzo di domicilio
						"C.A.P. Di domicilio",                                             //C.A.P. Di domicilio
						"Comune di domicilio",                                             //Comune di domicilio
						"Provincia di domicilio",                                          //Provincia di domicilio
						"Luogo di nascita dell’contraente persona fisica",                 //Luogo di nascita dell’contraente persona fisica
						"Provincia di nascita dell’contraente persona fisica",             //Provincia di nascita dell’contraente persona fisica
						"Stato di residenza dell’contraente ",                             //Stato di residenza dell’contraente
						"Cab della città di residenza dell’contraente",                    //Cab della città di residenza dell’contraente
						"Sottogruppo attività economica",                                  //Sottogruppo attività economica
						"Ramo gruppo attività economica",                                  //Ramo gruppo attività economica
						"Tipo documento dell'contraente persona fisica ",                  //Tipo documento dell'contraente persona fisica
						"Numero documento dell'contraente persona fisica ",                //Numero documento dell'contraente persona fisica
						"Data rilascio documento dell'contraente persona fisica ",         //Data rilascio documento dell'contraente persona fisica
						"Ente rilascio documento dell'contraente persona fisica ",         //Ente rilascio documento dell'contraente persona fisica
						"PEP - Persona Politicamente Esposta",                             //PEP - Persona Politicamente Esposta
						"Tipologia di PEP",                                                //Tipologia di PEP
						"Modalità di comunicazione prescelta tra Compagnia ed contraente", //Modalità di comunicazione prescelta tra Compagnia ed contraente
						"Indirizzo di residenza ",                                         //Indirizzo di residenza
						"C.A.P. Residenza ",                                               //C.A.P. Residenza
						"Comune Residenza ",                                               //Comune Residenza
						"Provincia Residenza ",                                            //Provincia Residenza
						"Indirizzo di domicilio",                                          //Indirizzo di domicilio
						"C.A.P. Domicilio",                                                //C.A.P. Domicilio
						"Comune Domicilio",                                                //Comune Domicilio
						"Provincia Domicilio",                                             //Provincia Domicilio
						"Indirizzo e-mail ",                                               //Indirizzo e-mail
						"Numero di cellulare ",                                            //Numero di cellulare
						"Luogo di nascita ",                                               //Luogo di nascita
						"Provincia di nascita ",                                           //Provincia di nascita
						"Stato di residenza",                                              //Stato di residenza
						"Tipo documento",                                                  //Tipo documento
						"Numero documento",                                                //Numero documento
						"Data rilascio documento",                                         //Data rilascio documento
						"Ente rilascio documento",                                         //Ente rilascio documento
						"PEP - Persona Politicamente Esposta",                             //PEP - Persona Politicamente Esposta
						"Tipologia di PEP",                                                //Tipologia di PEP
						"Eredi designati nominativamente o genericamente?",                //Eredi designati nominativamente o genericamente?
						"Cognome Beneficiario 1",                                          //Cognome Beneficiario 1
						"Nome",                                                            //Nome
						"Codice Fiscale ",                                                 //Codice Fiscale
						"Numero di Telefono del Beneficiario",                             //Numero di Telefono del Beneficiario
						"Indirizzo di residenza ",                                         //Indirizzo di residenza
						"Città /Comune di Residenza",                                      //Città /Comune di Residenza
						"CAP",                                                             //CAP
						"Provincia",                                                       //CAP
						"Email ",                                                          //Email
						"Legame del Cliente col Beneficiario",                             //Legame del Cliente col Beneficiario
						"NUCLEO FAMILIARE",                                                //NUCLEO FAMILIARE
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
						"Carica ricoperta dell'esecutore ",    //Carica ricoperta dell'esecutore
						"Cognome ",                            //Cognome
						"Nome",                                //Nome
						"Indirizzo di residenza ",             //Indirizzo di residenza
						"Città /Comune di Residenza",          //Città /Comune di Residenza
						"CAP",                                 //CAP
						"Codice Fiscale",                      //Codice Fiscale
						"Numero di Telefono",                  //Numero di Telefono
						"Email",                               //Email
					}

					result = append(result, row)

				}

			}

		}

	}
	return "", nil, e
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
		"AML DATA",
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
