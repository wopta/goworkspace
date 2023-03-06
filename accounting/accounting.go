package accounting

import (
	"bufio"
	"encoding/csv"
	"log"
	"os"
	"strings"

	lib "github.com/wopta/goworkspace/lib"
)

func main() {
	GetFx()
}
func csvL() {

	//dfres := dataframe.NewDataFrame()
	//dfres := lib.NewDf()
	if _, err := os.Stat("output.csv"); err == nil {
		log.Printf("File exists\n")
		e := os.Remove("output.csv")
		lib.CheckError(e)
	} else {
		log.Printf("File does not exist\n")
	}

	filesPath := lib.Files("")
	f, err := os.Create("output.csv")
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()
	w.Write([]string{"FC", "Compagnia", "Ramo", "Num. polizza", "Tipo", "titolo", "Subagente", "Data Titolo", "Data Pagamento", "Mezzo Pagamento", "Premio tot.", "Prv. totali", "Prv. acq.", "Contr./Descr.", "Data ins.", "Data rif. cont.", "Abbin.", "Valuta", "Modello", "RischioCodice", "Codice Azienda", "Partita IVA", "Premio Imponibile", "Provvigioni %", "check provv"})
	for _, f := range filesPath {
		if strings.Contains(strings.ToUpper(f), "GLOBAL") {
			global(f, w)

		}
		if strings.Contains(strings.ToUpper(f), "ELBA") {
			elba(f, w)

		}

	}
	w.Flush()
	f.Close()
	lib.CheckError(err)

}
func global(f string, w *csv.Writer) {
	readFile, err := os.Open(f)
	defer readFile.Close()
	lib.CheckError(err)
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	for fileScanner.Scan() {
		cels := strings.Split(fileScanner.Text(), "      ")
		if !strings.Contains(fileScanner.Text(), "WOPTA ASSICURAZIONI SRL - 0920") && len(cels) > 5 {

			log.Println("len cles:", len(cels))
			log.Println(cels)
			var celsClean []string
			log.Println("-----------------------CLEAN DATA----------------------------------")
			for i, c := range cels {
				if !(len(strings.TrimSpace(c)) == 0) {
					log.Println("append: ", strings.TrimSpace(c))

					celsClean = append(celsClean, strings.TrimSpace(c))
				}

				log.Println(i)
				log.Println(c)
			}
			log.Println("-----------------------MAP TO RESULT----------------------------------")
			if strings.Contains(celsClean[1], "G") {
				s := celsClean[2]

				w.Write([]string{"", "Global", "", celsClean[1], "", "", "", "", s[24:34], "", strings.Replace(s[:len(s)-50], ".", ",", 1), "", "", celsClean[4], "", "", "", "EUR", "", "", "", "", "", "", ""})
			}
			for i, c := range celsClean {

				log.Println(i)
				log.Println(c)
			}

		}
	}
	readFile.Close()
	log.Println("---------------END FILE GLOBAL--------------------------------")
}
func elba(f string, w *csv.Writer) {
	readFile, err := os.Open(f)
	defer readFile.Close()
	lib.CheckError(err)
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	count := 0
	for fileScanner.Scan() {
		text := strings.Replace(fileScanner.Text(), "\"", "", -1)
		log.Println("clen text: ", text)
		cels := strings.Split(text, ";")
		if count > 6 && len(cels) > 10 {
			if !strings.Contains(fileScanner.Text(), "Totale premi") || !strings.Contains(fileScanner.Text(), "Totale provvigioni") || !strings.Contains(fileScanner.Text(), "Saldo") {

				log.Println(cels)
				w.Write([]string{"", "Elba", cels[1], cels[2], cels[3], "", "", "", cels[6], cels[7], cels[8], cels[9], "", cels[11], cels[12], cels[13], cels[14], cels[15], "", "", "", "", "", "", ""})
				log.Println("-----------------------MAP TO RESULT--- ELBA----------------------------------")
				for i, c := range cels {
					log.Println(i)
					log.Println(c)
				}

			}
		}
		count++
	}
	readFile.Close()
	log.Println("---------------END FILE ELBA--------------------------------")
}
