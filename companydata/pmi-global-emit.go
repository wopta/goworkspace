package companydata

import (
	"fmt"
	"net/http"
	"os"
	"time"

	lib "github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func PmiGlobalEmit(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	config := lib.SftpConfig{
		Username:     os.Getenv("GLOBAL_SFTP_USER"),
		Password:     os.Getenv("GLOBAL_SFTP_PSW"), // required only if password authentication is to be used
		PrivateKey:   "",                           // required only if private key authentication is to be used
		Server:       "ftps.globalassistance.it:222",
		KeyExchanges: []string{"diffie-hellman-group-exchange-sha1", "diffie-hellman-group1-sha1", "diffie-hellman-group14-sha1"}, // optional
		Timeout:      time.Second * 30,                                                                                            // 0 for not timeout
	}
	client, e := lib.NewSftpclient(config)
	reader, e := client.Download("wopta/")
	lib.ExcelRead(reader)
	q := lib.Firequeries{
		Queries: []lib.Firequery{{
			Field:      "companyEmit",
			Operator:   "==",
			QueryValue: true,
		},
			{
				Field:      "companyEmitted",
				Operator:   "==",
				QueryValue: false,
			},
		},
	}
	query := q.FirestoreWherefields("policy")
	policies := models.PolicyToListData(query)

	for _, policy := range policies {
		var (
			result         [][]string
			enterpriseName string
			employer       string
			revenue        string
		)

		startDate := lib.Dateformat(policy.StartDate)
		endDate := lib.Dateformat(policy.EndDate)
		crationDate := lib.Dateformat(policy.CreationDate)
		companyCode := policy.CodeCompany
		cityCode := policy.Contractor.CityCode
		city := policy.Contractor.City
		streetNumber := policy.Contractor.StreetNumber
		postalCode := policy.Contractor.PostalCode
		pi := policy.Contractor.VatCode
		fc := policy.Contractor.FiscalCode
		for _, asset := range policy.Assets {
			if asset.Building != nil {

			}
			if asset.Enterprise != nil {
				enterpriseName = asset.Enterprise.Name
				employer = fmt.Sprint(asset.Enterprise.Employer)
				revenue = fmt.Sprint(asset.Enterprise.Revenue)
			}
		}
		for _, asset := range policy.Assets {
			if asset.Building != nil {
				for _, g := range asset.Guarantees {
					//"TIPO OPERAZIONE",N. POLIZZA SOSTITUITA,	DENOMINAZIONE PRODOTTO,	NODO DI GESTIONE,	DATA EMISSIONE,	DATA EFFETTO,	PARTITA IVA CONTRAENTE,	CODICE FISCALE CONTRAENTE	NATURA GIURIDICA CONTRAENTE	RAGIONE SOCIALE CONTRAENTE	PROVINCIA CONTRAENTE	COMUNE CONTRAENTE	CAP CONTRAENTE	TOPONIMO CONTRAENTE	INDIRIZZO CONTRAENTE	NUMERO CIVICO CONTRAENTE	DATA SCADENZA	FRAZIONAMENTO	VINCOLO	NUMERO ADDETTI	COSA SI VUOLE ASSICURARE	DOMANDA 1	DOMANDA 2	DOMANDA 3	FATTURATO	FORMA DI COPERTURA	FORMULA INCENDIO	BENE	ANNO DI COSTRUZIONE FABBRICATO	MATERIALE COSTRUZIONE	NUMERO PIANI	PRESENZA ALLARME	PRESENZA POLIZZA CONDOMINIALE	TIPOLOGIA FABBRICATO	PROVINCIA UBICAZIONE	COMUNE UBICAZIONE	CAP UBICAZIONE	TOPONIMO UBICAZIONE	INDIRIZZO UBICAZIONE	NUMERO CIVICO UBICAZIONE	CODICE ATTIVITA' - BENI	CLASSE - SOLO BENI	SETTORE - BENI	TIPO - BENI	CLAUSOLA VINCOLO	TESTO CLAUSOLA VINCOLO	GARANZIE/PACCHETTI - BENI	FRANCHIGIA - BENI	SOMMA ASSICURATA - BENI	SCOPERTO - BENI	% SOMMA ASSICURATA INCENDIO FABBRICATO E CONTENUTO - BENI	MASSIMALE - BENI	DIARIA - BENI	CODICE ATTIVITA' - ATTIVITA'	CLASSE - ATTIVITA'	SETTORE - ATTIVITA'	TIPO - ATTIVITA'	GARANZIE/PACCHETTI - ATTIVITA'	FRANCHIGIA - ATTIVITA'	SCOPERTO - ATTIVITA'	MASSIMALE - ATTIVITA'	MASSIMALE PER EVENTO - ATTIVITA'	PREMIO ANNUO LORDO DI GARANZIA	SCONTO %	RATA ALLA FIRMA	RATA SUCCESSIVA	DATA SCADENZA I RATA	NUMERO POLIZZA

					//row := []string{"TIPO OPERAZIONE", "N. POLIZZA SOSTITUITA", "DENOMINAZIONE PRODOTTO", "NODO DI GESTIONE", "DATA EMISSIONE", "DATA EFFETTO", "PARTITA IVA CONTRAENTE", "CODICE FISCALE CONTRAENTE", "NATURA GIURIDICA CONTRAENTE", "RAGIONE SOCIALE CONTRAENTE", "PROVINCIA CONTRAENTE", "COMUNE CONTRAENTE", "CAP CONTRAENTE", "TOPONIMO CONTRAENTE", "INDIRIZZO CONTRAENTE", "NUMERO CIVICO CONTRAENTE", "DATA SCADENZA", "FRAZIONAMENTO", "VINCOLO", "NUMERO ADDETTI", "COSA SI VUOLE ASSICURARE", "DOMANDA 1", "DOMANDA 2", "DOMANDA 3", "FATTURATO", "FORMA DI COPERTURA", "FORMULA INCENDIO", "BENE", "ANNO DI COSTRUZIONE FABBRICATO", "MATERIALE COSTRUZIONE", "NUMERO PIANI", "PRESENZA ALLARME", "PRESENZA POLIZZA CONDOMINIALE", "TIPOLOGIA FABBRICATO", "PROVINCIA UBICAZIONE", "COMUNE UBICAZIONE", "CAP UBICAZIONE", "TOPONIMO UBICAZIONE", "INDIRIZZO UBICAZIONE", "NUMERO CIVICO UBICAZIONE", "CODICE ATTIVITA' - BENI", "CLASSE - SOLO BENI", "SETTORE - BENI", "TIPO - BENI", "CLAUSOLA VINCOLO", "TESTO CLAUSOLA VINCOLO", "GARANZIE/PACCHETTI - BENI", "FRANCHIGIA - BENI", "SOMMA ASSICURATA - BENI", "SCOPERTO - BENI", "% SOMMA ASSICURATA INCENDIO FABBRICATO E CONTENUTO - BENI", "MASSIMALE - BENI", "DIARIA - BENI", "CODICE ATTIVITA' - ATTIVITA'", "CLASSE - ATTIVITA'", "SETTORE - ATTIVITA'", "TIPO - ATTIVITA'", "GARANZIE/PACCHETTI - ATTIVITA'", "FRANCHIGIA - ATTIVITA'", "SCOPERTO - ATTIVITA'", "MASSIMALE - ATTIVITA'", "MASSIMALE PER EVENTO - ATTIVITA'", "PREMIO ANNUO LORDO DI GARANZIA", "SCONTO %", "RATA ALLA FIRMA", "RATA SUCCESSIVA", "DATA SCADENZA I RATA", "NUMERO POLIZZA"}
					row := []string{"Nuova emissione", "", "WOPTA PER TE. ARTIGIANI & IMPRESE", "0920", crationDate, startDate, pi, fc, "", "RAGIONE SOCIALE CONTRAENTE", "PROVINCIA CONTRAENTE", "COMUNE CONTRAENTE", "CAP CONTRAENTE", "TOPONIMO CONTRAENTE", "INDIRIZZO CONTRAENTE", "NUMERO CIVICO CONTRAENTE", "DATA SCADENZA", "FRAZIONAMENTO", "VINCOLO", "NUMERO ADDETTI", "COSA SI VUOLE ASSICURARE", "DOMANDA 1", "DOMANDA 2", "DOMANDA 3", "FATTURATO", "FORMA DI COPERTURA", "FORMULA INCENDIO", "BENE", "ANNO DI COSTRUZIONE FABBRICATO", "MATERIALE COSTRUZIONE", "NUMERO PIANI", "PRESENZA ALLARME", "PRESENZA POLIZZA CONDOMINIALE", "TIPOLOGIA FABBRICATO", "PROVINCIA UBICAZIONE", "COMUNE UBICAZIONE", "CAP UBICAZIONE", "TOPONIMO UBICAZIONE", "INDIRIZZO UBICAZIONE", "NUMERO CIVICO UBICAZIONE", "CODICE ATTIVITA' - BENI", "CLASSE - SOLO BENI", "SETTORE - BENI", "TIPO - BENI", "CLAUSOLA VINCOLO", "TESTO CLAUSOLA VINCOLO", "GARANZIE/PACCHETTI - BENI", "FRANCHIGIA - BENI", "SOMMA ASSICURATA - BENI", "SCOPERTO - BENI", "% SOMMA ASSICURATA INCENDIO FABBRICATO E CONTENUTO - BENI", "MASSIMALE - BENI", "DIARIA - BENI", "CODICE ATTIVITA' - ATTIVITA'", "CLASSE - ATTIVITA'", "SETTORE - ATTIVITA'", "TIPO - ATTIVITA'", "GARANZIE/PACCHETTI - ATTIVITA'", "FRANCHIGIA - ATTIVITA'", "SCOPERTO - ATTIVITA'", "MASSIMALE - ATTIVITA'", "MASSIMALE PER EVENTO - ATTIVITA'", "PREMIO ANNUO LORDO DI GARANZIA", "SCONTO %", "RATA ALLA FIRMA", "RATA SUCCESSIVA", "DATA SCADENZA I RATA", "NUMERO POLIZZA"}

					result = append(result, row)

				}

			}
			if asset.Enterprise != nil {
				for _, g := range asset.Guarantees {
					//TIPO OPERAZIONE	N. POLIZZA SOSTITUITA	DENOMINAZIONE PRODOTTO	NODO DI GESTIONE	DATA EMISSIONE	DATA EFFETTO	PARTITA IVA CONTRAENTE	CODICE FISCALE CONTRAENTE	NATURA GIURIDICA CONTRAENTE	RAGIONE SOCIALE CONTRAENTE	PROVINCIA CONTRAENTE	COMUNE CONTRAENTE	CAP CONTRAENTE	TOPONIMO CONTRAENTE	INDIRIZZO CONTRAENTE	NUMERO CIVICO CONTRAENTE	DATA SCADENZA	FRAZIONAMENTO	VINCOLO	NUMERO ADDETTI	COSA SI VUOLE ASSICURARE	DOMANDA 1	DOMANDA 2	DOMANDA 3	FATTURATO	FORMA DI COPERTURA	FORMULA INCENDIO	BENE	ANNO DI COSTRUZIONE FABBRICATO	MATERIALE COSTRUZIONE	NUMERO PIANI	PRESENZA ALLARME	PRESENZA POLIZZA CONDOMINIALE	TIPOLOGIA FABBRICATO	PROVINCIA UBICAZIONE	COMUNE UBICAZIONE	CAP UBICAZIONE	TOPONIMO UBICAZIONE	INDIRIZZO UBICAZIONE	NUMERO CIVICO UBICAZIONE	CODICE ATTIVITA' - BENI	CLASSE - SOLO BENI	SETTORE - BENI	TIPO - BENI	CLAUSOLA VINCOLO	TESTO CLAUSOLA VINCOLO	GARANZIE/PACCHETTI - BENI	FRANCHIGIA - BENI	SOMMA ASSICURATA - BENI	SCOPERTO - BENI	% SOMMA ASSICURATA INCENDIO FABBRICATO E CONTENUTO - BENI	MASSIMALE - BENI	DIARIA - BENI	CODICE ATTIVITA' - ATTIVITA'	CLASSE - ATTIVITA'	SETTORE - ATTIVITA'	TIPO - ATTIVITA'	GARANZIE/PACCHETTI - ATTIVITA'	FRANCHIGIA - ATTIVITA'	SCOPERTO - ATTIVITA'	MASSIMALE - ATTIVITA'	MASSIMALE PER EVENTO - ATTIVITA'	PREMIO ANNUO LORDO DI GARANZIA	SCONTO %	RATA ALLA FIRMA	RATA SUCCESSIVA	DATA SCADENZA I RATA	NUMERO POLIZZA
					//row := []string{"TIPO OPERAZIONE", "N. POLIZZA SOSTITUITA", "DENOMINAZIONE PRODOTTO", "NODO DI GESTIONE", "DATA EMISSIONE", "DATA EFFETTO", "PARTITA IVA CONTRAENTE", "CODICE FISCALE CONTRAENTE", "NATURA GIURIDICA CONTRAENTE", "RAGIONE SOCIALE CONTRAENTE", "PROVINCIA CONTRAENTE", "COMUNE CONTRAENTE", "CAP CONTRAENTE", "TOPONIMO CONTRAENTE", "INDIRIZZO CONTRAENTE", "NUMERO CIVICO CONTRAENTE", "DATA SCADENZA", "FRAZIONAMENTO", "VINCOLO", "NUMERO ADDETTI", "COSA SI VUOLE ASSICURARE", "DOMANDA 1", "DOMANDA 2", "DOMANDA 3", "FATTURATO", "FORMA DI COPERTURA", "FORMULA INCENDIO", "BENE", "ANNO DI COSTRUZIONE FABBRICATO", "MATERIALE COSTRUZIONE", "NUMERO PIANI", "PRESENZA ALLARME", "PRESENZA POLIZZA CONDOMINIALE", "TIPOLOGIA FABBRICATO", "PROVINCIA UBICAZIONE", "COMUNE UBICAZIONE", "CAP UBICAZIONE", "TOPONIMO UBICAZIONE", "INDIRIZZO UBICAZIONE", "NUMERO CIVICO UBICAZIONE", "CODICE ATTIVITA' - BENI", "CLASSE - SOLO BENI", "SETTORE - BENI", "TIPO - BENI", "CLAUSOLA VINCOLO", "TESTO CLAUSOLA VINCOLO", "GARANZIE/PACCHETTI - BENI", "FRANCHIGIA - BENI", "SOMMA ASSICURATA - BENI", "SCOPERTO - BENI", "% SOMMA ASSICURATA INCENDIO FABBRICATO E CONTENUTO - BENI", "MASSIMALE - BENI", "DIARIA - BENI", "CODICE ATTIVITA' - ATTIVITA'", "CLASSE - ATTIVITA'", "SETTORE - ATTIVITA'", "TIPO - ATTIVITA'", "GARANZIE/PACCHETTI - ATTIVITA'", "FRANCHIGIA - ATTIVITA'", "SCOPERTO - ATTIVITA'", "MASSIMALE - ATTIVITA'", "MASSIMALE PER EVENTO - ATTIVITA'", "PREMIO ANNUO LORDO DI GARANZIA", "SCONTO %", "RATA ALLA FIRMA", "RATA SUCCESSIVA", "DATA SCADENZA I RATA", "NUMERO POLIZZA"}
					row := []string{"Nuova emissione", "", "WOPTA PER TE. ARTIGIANI & IMPRESE", "0920", crationDate, startDate, pi, fc, "", "RAGIONE SOCIALE CONTRAENTE", "PROVINCIA CONTRAENTE", "COMUNE CONTRAENTE", "CAP CONTRAENTE", "TOPONIMO CONTRAENTE", "INDIRIZZO CONTRAENTE", "NUMERO CIVICO CONTRAENTE", "DATA SCADENZA", "FRAZIONAMENTO", "VINCOLO", "NUMERO ADDETTI", "COSA SI VUOLE ASSICURARE", "DOMANDA 1", "DOMANDA 2", "DOMANDA 3", "FATTURATO", "FORMA DI COPERTURA", "FORMULA INCENDIO", "BENE", "ANNO DI COSTRUZIONE FABBRICATO", "MATERIALE COSTRUZIONE", "NUMERO PIANI", "PRESENZA ALLARME", "PRESENZA POLIZZA CONDOMINIALE", "TIPOLOGIA FABBRICATO", "PROVINCIA UBICAZIONE", "COMUNE UBICAZIONE", "CAP UBICAZIONE", "TOPONIMO UBICAZIONE", "INDIRIZZO UBICAZIONE", "NUMERO CIVICO UBICAZIONE", "CODICE ATTIVITA' - BENI", "CLASSE - SOLO BENI", "SETTORE - BENI", "TIPO - BENI", "CLAUSOLA VINCOLO", "TESTO CLAUSOLA VINCOLO", "GARANZIE/PACCHETTI - BENI", "FRANCHIGIA - BENI", "SOMMA ASSICURATA - BENI", "SCOPERTO - BENI", "% SOMMA ASSICURATA INCENDIO FABBRICATO E CONTENUTO - BENI", "MASSIMALE - BENI", "DIARIA - BENI", "CODICE ATTIVITA' - ATTIVITA'", "CLASSE - ATTIVITA'", "SETTORE - ATTIVITA'", "TIPO - ATTIVITA'", "GARANZIE/PACCHETTI - ATTIVITA'", "FRANCHIGIA - ATTIVITA'", "SCOPERTO - ATTIVITA'", "MASSIMALE - ATTIVITA'", "MASSIMALE PER EVENTO - ATTIVITA'", "PREMIO ANNUO LORDO DI GARANZIA", "SCONTO %", "RATA ALLA FIRMA", "RATA SUCCESSIVA", "DATA SCADENZA I RATA", "NUMERO POLIZZA"}

					result = append(result, row)

				}

			}

		}

	}
	return "", nil, e
}
func getMapBuildingMaterial(key string) string {
	var res string
	mapGarante := map[string]string{
		"rsc": "",
		"r":   "",
		"gri": "",
		"adg": "",
	}

	if seconds, ok := mapGarante[key]; ok { // will be false if person is not in the map
		res = seconds
	}
	return res
}
func getMapRevenue(key int) int {
	var res int

	if key <= 200000 { // will be false if person is not in the map
		res = 1
	}
	if key > 200000 && key <= 500000 { // will be false if person is not in the map
		res = 2
	}
	if key > 500000 && key <= 1000000 { // will be false if person is not in the map
		res = 3
	}
	if key > 1000000 && key <= 1500000 { // will be false if person is not in the map
		res = 4
	}
	if key > 1500000 && key <= 5000000 { // will be false if person is not in the map
		res = 5
	}
	if key > 5000000 && key <= 7500000 { // will be false if person is not in the map
		res = 6
	}
	if key > 7500000 && key <= 10000000 { // will be false if person is not in the map
		res = 7
	}
	return res
}
