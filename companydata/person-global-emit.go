package companydata

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	lib "github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func PersonGlobalEmit(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		result [][]string

		e error
	)

	//layout := "02/01/2006"
	layoutFilename := "20060102"
	//client, e := lib.NewSftpclient(config)
	location, e := time.LoadLocation("Europe/Rome")
	collection := "policy"
	fmt.Println(time.Now().In(location))

	now := time.Now().In(location).AddDate(0, 0, -1)

	filename := now.Format(layoutFilename) + "_EM_PMIW.XLSX"
	//println(config)
	println("filename: ", filename)
	GlobalSftpDownload(""+filename, "track/in/global/emit/", "/Wopta/")
	excelsource, _ := lib.ExcelReadFile("../tmp/" + filename)
	for k, v := range excelsource {
		println("key shhet name: ", k)
		result = v
	}
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
			{
				Field:      "company", //
				Operator:   "==",      //
				QueryValue: "global",
			},
			{
				Field:      "name", //
				Operator:   "==",   //
				QueryValue: "pmi",
			},
			{
				Field:      "IsDeleted", //
				Operator:   "==",        //
				QueryValue: false,
			},
		},
	}
	query, e := q.FirestoreWherefields("policy")
	policies := models.PolicyToListData(query)
	log.Println("len(policies):", len(policies))
	for _, policy := range policies {

		for _, asset := range policy.Assets {
			if asset.Building != nil {
				for _, g := range asset.Guarantees {
					//"TIPO OPERAZIONE",N. POLIZZA SOSTITUITA,	DENOMINAZIONE PRODOTTO,	NODO DI GESTIONE,	DATA EMISSIONE,	DATA EFFETTO,	PARTITA IVA CONTRAENTE,	CODICE FISCALE CONTRAENTE	NATURA GIURIDICA CONTRAENTE	RAGIONE SOCIALE CONTRAENTE	PROVINCIA CONTRAENTE	COMUNE CONTRAENTE	CAP CONTRAENTE	TOPONIMO CONTRAENTE	INDIRIZZO CONTRAENTE	NUMERO CIVICO CONTRAENTE	DATA SCADENZA	FRAZIONAMENTO	VINCOLO	NUMERO ADDETTI	COSA SI VUOLE ASSICURARE	DOMANDA 1	DOMANDA 2	DOMANDA 3	FATTURATO	FORMA DI COPERTURA	FORMULA INCENDIO	BENE	ANNO DI COSTRUZIONE FABBRICATO	MATERIALE COSTRUZIONE	NUMERO PIANI	PRESENZA ALLARME	PRESENZA POLIZZA CONDOMINIALE	TIPOLOGIA FABBRICATO	PROVINCIA UBICAZIONE	COMUNE UBICAZIONE	CAP UBICAZIONE	TOPONIMO UBICAZIONE	INDIRIZZO UBICAZIONE	NUMERO CIVICO UBICAZIONE	CODICE ATTIVITA' - BENI	CLASSE - SOLO BENI	SETTORE - BENI	TIPO - BENI	CLAUSOLA VINCOLO	TESTO CLAUSOLA VINCOLO	GARANZIE/PACCHETTI - BENI	FRANCHIGIA - BENI	SOMMA ASSICURATA - BENI	SCOPERTO - BENI	% SOMMA ASSICURATA INCENDIO FABBRICATO E CONTENUTO - BENI	MASSIMALE - BENI	DIARIA - BENI	CODICE ATTIVITA' - ATTIVITA'	CLASSE - ATTIVITA'	SETTORE - ATTIVITA'	TIPO - ATTIVITA'	GARANZIE/PACCHETTI - ATTIVITA'	FRANCHIGIA - ATTIVITA'	SCOPERTO - ATTIVITA'	MASSIMALE - ATTIVITA'	MASSIMALE PER EVENTO - ATTIVITA'	PREMIO ANNUO LORDO DI GARANZIA	SCONTO %	RATA ALLA FIRMA	RATA SUCCESSIVA	DATA SCADENZA I RATA	NUMERO POLIZZA
					fmt.Println(g)
					//row := []string{"TIPO OPERAZIONE", "N. POLIZZA SOSTITUITA", "DENOMINAZIONE PRODOTTO", "NODO DI GESTIONE", "DATA EMISSIONE", "DATA EFFETTO", "PARTITA IVA CONTRAENTE", "CODICE FISCALE CONTRAENTE", "NATURA GIURIDICA CONTRAENTE", "RAGIONE SOCIALE CONTRAENTE", "PROVINCIA CONTRAENTE", "COMUNE CONTRAENTE", "CAP CONTRAENTE", "TOPONIMO CONTRAENTE", "INDIRIZZO CONTRAENTE", "NUMERO CIVICO CONTRAENTE", "DATA SCADENZA", "FRAZIONAMENTO", "VINCOLO", "NUMERO ADDETTI", "COSA SI VUOLE ASSICURARE", "DOMANDA 1", "DOMANDA 2", "DOMANDA 3", "FATTURATO", "FORMA DI COPERTURA", "FORMULA INCENDIO", "BENE", "ANNO DI COSTRUZIONE FABBRICATO", "MATERIALE COSTRUZIONE", "NUMERO PIANI", "PRESENZA ALLARME", "PRESENZA POLIZZA CONDOMINIALE", "TIPOLOGIA FABBRICATO", "PROVINCIA UBICAZIONE", "COMUNE UBICAZIONE", "CAP UBICAZIONE", "TOPONIMO UBICAZIONE", "INDIRIZZO UBICAZIONE", "NUMERO CIVICO UBICAZIONE", "CODICE ATTIVITA' - BENI", "CLASSE - SOLO BENI", "SETTORE - BENI", "TIPO - BENI", "CLAUSOLA VINCOLO", "TESTO CLAUSOLA VINCOLO", "GARANZIE/PACCHETTI - BENI", "FRANCHIGIA - BENI", "SOMMA ASSICURATA - BENI", "SCOPERTO - BENI", "% SOMMA ASSICURATA INCENDIO FABBRICATO E CONTENUTO - BENI", "MASSIMALE - BENI", "DIARIA - BENI", "CODICE ATTIVITA' - ATTIVITA'", "CLASSE - ATTIVITA'", "SETTORE - ATTIVITA'", "TIPO - ATTIVITA'", "GARANZIE/PACCHETTI - ATTIVITA'", "FRANCHIGIA - ATTIVITA'", "SCOPERTO - ATTIVITA'", "MASSIMALE - ATTIVITA'", "MASSIMALE PER EVENTO - ATTIVITA'", "PREMIO ANNUO LORDO DI GARANZIA", "SCONTO %", "RATA ALLA FIRMA", "RATA SUCCESSIVA", "DATA SCADENZA I RATA", "NUMERO POLIZZA"}
					row := []string{
						"", //TIPO OPERAZIONE
						"", //N. POLIZZA SOSTITUITA
						"", //DENOMINAZIONE PRODOTTO
						"", //NODO DI GESTIONE
						"", //DATA EMISSIONE
						"", //DATA EFFETTO
						"", //CODICE FISCALE CONTRAENTE
						"", //PIVA CONTRAENTE
						"", //COGNOME/RAGIONE SOCIALE CONTRAENTE
						"", //NOME CONTRENTE
						"", //SESSO CONTRAENTE
						"", //DATA DI NASCITA CONTRAENTE
						"", //PROVINCIA DI NASCITA CONTRAENTE
						"", //COMUNE DI NASCITA CONTRAENTE
						"", //PROVINCIA CONTRAENTE
						"", //COMUNE CONTRAENTE
						"", //CAP CONTRAENTE
						"", //TOPONIMO CONTRAENTE
						"", //INDIRIZZO CONTRAENTE
						"", //NUMERO CIVICO CONTRAENTE
						"", //DATA SCADENZA
						"", //FRAZIONAMENTO
						"", //CANALE
						"", //SCELTA OPZIONI
						"", //CODICE FISCALE ASSICURATO
						"", //PIVA ASSICURATO
						"", //COGNOME/RAGIONE SOCIALE ASSICURATO
						"", //NOME ASSICURATO
						"", //SESSO ASSICURATO
						"", //DATA DI NASCITA ASSICURATO
						"", //PROVINCIA DI NASCITA ASSICURATO
						"", //COMUNE DI NASCITA ASSICURATO
						"", //PROVINCIA ASSICURATO
						"", //COMUNE ASSICURATO
						"", //CAP ASSICURATO
						"", //TOPONIMO ASSICURATO
						"", //INDIRIZZO ASSICURATO
						"", //NUMERO CIVICO ASSICURATO
						"", //CODICE ATTIVITA' ASSICURATO
						"", //SETTORE ASSICURATO
						"", //TIPO ASSICURATO
						"", //GARANZIE/PACCHETTI
						"", //KEY MAN
						"", //ESTENSIONE SUPERVALUTAZIONE ARTI SUPERIORI
						"", //FRANCHIGIA
						"", //MASSIMALE
						"", //TIPO COPERTURA INFORTUNI
						"", //PREMIO ANNUO LORDO DI GARANZIA
						"", //SCONTO %
						"", //RATA ALLA FIRMA
						"", //RATA SUCCESSIVA
						"", //DATA SCADENZA I RATA
						"", //NUMERO POLIZZA

					}
					result = append(result, row)

				}

			}

		}

		if e == nil {
			policy.CompanyEmitted = true
			lib.SetFirestore(collection, policy.Uid, policy)
		}
	}
	log.Println("len(result):", len(result))
	filepath := "../tmp/" + filename
	excel, e := lib.CreateExcel(result, filepath, "Risultato")
	lib.PutToStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "track/global/person/emit/"+filepath, <-excel)
	//lib.PutGoogleStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "track/global/pmi/emit/"+filepath, source, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	if len(policies) > 0 {
		GlobalSftpDelete("/Wopta/" + filename)
		GlobalSftpUpload(filename, "/Wopta/")
	}
	return "", nil, e
}
