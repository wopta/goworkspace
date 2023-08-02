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

const collection = "transactions"

func GlobalTransaction(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println("----------------GlobalTransaction-----------------")
	var (
		result [][]string

		e error
	)
	//layout := "02/01/2006"
	layoutFilename := "20060102"
	location, e := time.LoadLocation("Europe/Rome")

	now := time.Now().In(location).AddDate(0, 0, -1)
	filename := now.Format(layoutFilename) + "_EM_PMIW.XLSX"
	//println(config)
	println("filename: ", filename)
	GlobalSftpDownload(""+filename, "track/in/global/transactions/", "/Wopta/")
	excelsource, _ := lib.ExcelReadFile("../tmp/" + filename)
	for k, v := range excelsource {
		println("key shhet name: ", k)
		result = v
	}
	q := lib.Firequeries{
		Queries: []lib.Firequery{
			{
				Field:      "isEmit",
				Operator:   "==",
				QueryValue: false,
			},
			{
				Field:      "isPay",
				Operator:   "==",
				QueryValue: true,
			},
			{
				Field:      "company", //
				Operator:   "==",      //
				QueryValue: "global",
			},
		},
	}
	query, e := q.FirestoreWherefields(collection)
	transactions := models.TransactionToListData(query)
	log.Println("len(policies):", len(transactions))
	result = append(result, getTransData(transactions)...)
	log.Println("len(result):", len(result))
	filepath := "../tmp/" + filename
	excel, e := lib.CreateExcel(result, filepath, "Risultato")
	lib.PutToStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "track/global/transactions/"+filepath, <-excel)
	//lib.PutGoogleStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "track/global/pmi/emit/"+filepath, source, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	if len(transactions) > 0 {
		GlobalSftpDelete("/Wopta/" + filename)
		GlobalSftpUpload(filename, "/Wopta/")
	}
	return "", nil, e
}
func getTransData(transactions []models.Transaction) [][]string {
	var (
		result [][]string
	)
	for _, transaction := range transactions {

		//row := []string{"TIPO OPERAZIONE", "N. POLIZZA SOSTITUITA", "DENOMINAZIONE PRODOTTO", "NODO DI GESTIONE", "DATA EMISSIONE", "DATA EFFETTO", "PARTITA IVA CONTRAENTE", "CODICE FISCALE CONTRAENTE", "NATURA GIURIDICA CONTRAENTE", "RAGIONE SOCIALE CONTRAENTE", "PROVINCIA CONTRAENTE", "COMUNE CONTRAENTE", "CAP CONTRAENTE", "TOPONIMO CONTRAENTE", "INDIRIZZO CONTRAENTE", "NUMERO CIVICO CONTRAENTE", "DATA SCADENZA", "FRAZIONAMENTO", "VINCOLO", "NUMERO ADDETTI", "COSA SI VUOLE ASSICURARE", "DOMANDA 1", "DOMANDA 2", "DOMANDA 3", "FATTURATO", "FORMA DI COPERTURA", "FORMULA INCENDIO", "BENE", "ANNO DI COSTRUZIONE FABBRICATO", "MATERIALE COSTRUZIONE", "NUMERO PIANI", "PRESENZA ALLARME", "PRESENZA POLIZZA CONDOMINIALE", "TIPOLOGIA FABBRICATO", "PROVINCIA UBICAZIONE", "COMUNE UBICAZIONE", "CAP UBICAZIONE", "TOPONIMO UBICAZIONE", "INDIRIZZO UBICAZIONE", "NUMERO CIVICO UBICAZIONE", "CODICE ATTIVITA' - BENI", "CLASSE - SOLO BENI", "SETTORE - BENI", "TIPO - BENI", "CLAUSOLA VINCOLO", "TESTO CLAUSOLA VINCOLO", "GARANZIE/PACCHETTI - BENI", "FRANCHIGIA - BENI", "SOMMA ASSICURATA - BENI", "SCOPERTO - BENI", "% SOMMA ASSICURATA INCENDIO FABBRICATO E CONTENUTO - BENI", "MASSIMALE - BENI", "DIARIA - BENI", "CODICE ATTIVITA' - ATTIVITA'", "CLASSE - ATTIVITA'", "SETTORE - ATTIVITA'", "TIPO - ATTIVITA'", "GARANZIE/PACCHETTI - ATTIVITA'", "FRANCHIGIA - ATTIVITA'", "SCOPERTO - ATTIVITA'", "MASSIMALE - ATTIVITA'", "MASSIMALE PER EVENTO - ATTIVITA'", "PREMIO ANNUO LORDO DI GARANZIA", "SCONTO %", "RATA ALLA FIRMA", "RATA SUCCESSIVA", "DATA SCADENZA I RATA", "NUMERO POLIZZA"}
		row := []string{
			transaction.ScheduleDate,                     //DATA INCASSO
			fmt.Sprintf("%.2f", transaction.Amount),      //IMPORTO PREMIO LORDO
			fmt.Sprintf("%.2f", transaction.Commissions), //IMPORTO PROVVIGIONI
			transaction.Name,                             //CONTRAENTE
			transaction.NumberCompany,                    //NUMERO POLIZZA
			transaction.ScheduleDate,                     //DATA EFFETTO

		}
		result = append(result, row)

		transaction.IsEmit = true
		lib.SetFirestore(collection, transaction.Uid, transaction)

	}
	return result
}
