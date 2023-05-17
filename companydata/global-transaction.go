package companydata

import (
	"log"
	"net/http"
	"os"
	"time"

	lib "github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func GlobalTransaction(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		result [][]string

		e error
	)
	//layout := "02/01/2006"
	layoutFilename := "20060102"
	location, e := time.LoadLocation("Europe/Rome")
	collection := "transactions"

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
		},
	}
	query, e := q.FirestoreWherefields(collection)
	transactions := models.TransactionToListData(query)
	log.Println("len(policies):", len(transactions))
	for _, transaction := range transactions {

		//row := []string{"TIPO OPERAZIONE", "N. POLIZZA SOSTITUITA", "DENOMINAZIONE PRODOTTO", "NODO DI GESTIONE", "DATA EMISSIONE", "DATA EFFETTO", "PARTITA IVA CONTRAENTE", "CODICE FISCALE CONTRAENTE", "NATURA GIURIDICA CONTRAENTE", "RAGIONE SOCIALE CONTRAENTE", "PROVINCIA CONTRAENTE", "COMUNE CONTRAENTE", "CAP CONTRAENTE", "TOPONIMO CONTRAENTE", "INDIRIZZO CONTRAENTE", "NUMERO CIVICO CONTRAENTE", "DATA SCADENZA", "FRAZIONAMENTO", "VINCOLO", "NUMERO ADDETTI", "COSA SI VUOLE ASSICURARE", "DOMANDA 1", "DOMANDA 2", "DOMANDA 3", "FATTURATO", "FORMA DI COPERTURA", "FORMULA INCENDIO", "BENE", "ANNO DI COSTRUZIONE FABBRICATO", "MATERIALE COSTRUZIONE", "NUMERO PIANI", "PRESENZA ALLARME", "PRESENZA POLIZZA CONDOMINIALE", "TIPOLOGIA FABBRICATO", "PROVINCIA UBICAZIONE", "COMUNE UBICAZIONE", "CAP UBICAZIONE", "TOPONIMO UBICAZIONE", "INDIRIZZO UBICAZIONE", "NUMERO CIVICO UBICAZIONE", "CODICE ATTIVITA' - BENI", "CLASSE - SOLO BENI", "SETTORE - BENI", "TIPO - BENI", "CLAUSOLA VINCOLO", "TESTO CLAUSOLA VINCOLO", "GARANZIE/PACCHETTI - BENI", "FRANCHIGIA - BENI", "SOMMA ASSICURATA - BENI", "SCOPERTO - BENI", "% SOMMA ASSICURATA INCENDIO FABBRICATO E CONTENUTO - BENI", "MASSIMALE - BENI", "DIARIA - BENI", "CODICE ATTIVITA' - ATTIVITA'", "CLASSE - ATTIVITA'", "SETTORE - ATTIVITA'", "TIPO - ATTIVITA'", "GARANZIE/PACCHETTI - ATTIVITA'", "FRANCHIGIA - ATTIVITA'", "SCOPERTO - ATTIVITA'", "MASSIMALE - ATTIVITA'", "MASSIMALE PER EVENTO - ATTIVITA'", "PREMIO ANNUO LORDO DI GARANZIA", "SCONTO %", "RATA ALLA FIRMA", "RATA SUCCESSIVA", "DATA SCADENZA I RATA", "NUMERO POLIZZA"}
		row := []string{
			"", //DATA INCASSO
			"", //IMPORTO PREMIO LORDO
			"", //IMPORTO PROVVIGIONI
			"", //CONTRAENTE
			"", //NUMERO POLIZZA
			"", //DATA EFFETTO

		}
		result = append(result, row)

		if e == nil {
			transaction.IsEmit = true
			//lib.SetFirestore(collection, policy.Agent.Uid, policy)
		}
	}
	log.Println("len(result):", len(result))
	filepath := filename

	excel, e := lib.CreateExcel(result, "../tmp/"+filepath, "Risultato")
	//source, _ := ioutil.ReadFile("../tmp/" + filepath)

	lib.PutToStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "track/global/pmi/emit/"+filepath, <-excel)
	//lib.PutGoogleStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "track/global/pmi/emit/"+filepath, source, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	return "", nil, e
}
