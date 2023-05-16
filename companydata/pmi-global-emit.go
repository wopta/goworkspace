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
)

func PmiGlobalEmit(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		result                                    [][]string
		enterpriseName                            string
		employer, class, sector, atecoDesc, ateco string
		revenue                                   int
		e                                         error
	)

	layout := "02/01/2006"
	layoutFilename := "20060102"
	//client, e := lib.NewSftpclient(config)
	location, e := time.LoadLocation("Europe/Rome")
	fmt.Println(time.Now().In(location))
	now := time.Now().In(location).AddDate(0, 0, -1)
	filename := now.Format(layoutFilename) + "_EM_PMIW.XLSX"
	//println(config)
	println("filename: ", filename)
	_, reader, e := GlobalSftpDownload("./"+filename, "track/in/global/emit/", "")
	excelsource, e := lib.ExcelRead(reader)
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
	query, e := q.FirestoreWherefields("policy")
	policies := models.PolicyToListData(query)
	log.Println("len(policies):", len(policies))
	for _, policy := range policies {

		for _, asset := range policy.Assets {
			if asset.Building != nil {

			}
			if asset.Enterprise != nil {
				enterpriseName = asset.Enterprise.Name
				employer = fmt.Sprint(asset.Enterprise.Employer)
				revenue, _ = strconv.Atoi(asset.Enterprise.Revenue)
				sector = asset.Enterprise.AtecoMacro
				class = asset.Enterprise.AtecoSub
				atecoDesc = asset.Enterprise.AtecoDesc
				ateco = asset.Enterprise.Ateco
			}
		}
		for _, asset := range policy.Assets {
			if asset.Building != nil {
				for _, g := range asset.Guarantees {
					//"TIPO OPERAZIONE",N. POLIZZA SOSTITUITA,	DENOMINAZIONE PRODOTTO,	NODO DI GESTIONE,	DATA EMISSIONE,	DATA EFFETTO,	PARTITA IVA CONTRAENTE,	CODICE FISCALE CONTRAENTE	NATURA GIURIDICA CONTRAENTE	RAGIONE SOCIALE CONTRAENTE	PROVINCIA CONTRAENTE	COMUNE CONTRAENTE	CAP CONTRAENTE	TOPONIMO CONTRAENTE	INDIRIZZO CONTRAENTE	NUMERO CIVICO CONTRAENTE	DATA SCADENZA	FRAZIONAMENTO	VINCOLO	NUMERO ADDETTI	COSA SI VUOLE ASSICURARE	DOMANDA 1	DOMANDA 2	DOMANDA 3	FATTURATO	FORMA DI COPERTURA	FORMULA INCENDIO	BENE	ANNO DI COSTRUZIONE FABBRICATO	MATERIALE COSTRUZIONE	NUMERO PIANI	PRESENZA ALLARME	PRESENZA POLIZZA CONDOMINIALE	TIPOLOGIA FABBRICATO	PROVINCIA UBICAZIONE	COMUNE UBICAZIONE	CAP UBICAZIONE	TOPONIMO UBICAZIONE	INDIRIZZO UBICAZIONE	NUMERO CIVICO UBICAZIONE	CODICE ATTIVITA' - BENI	CLASSE - SOLO BENI	SETTORE - BENI	TIPO - BENI	CLAUSOLA VINCOLO	TESTO CLAUSOLA VINCOLO	GARANZIE/PACCHETTI - BENI	FRANCHIGIA - BENI	SOMMA ASSICURATA - BENI	SCOPERTO - BENI	% SOMMA ASSICURATA INCENDIO FABBRICATO E CONTENUTO - BENI	MASSIMALE - BENI	DIARIA - BENI	CODICE ATTIVITA' - ATTIVITA'	CLASSE - ATTIVITA'	SETTORE - ATTIVITA'	TIPO - ATTIVITA'	GARANZIE/PACCHETTI - ATTIVITA'	FRANCHIGIA - ATTIVITA'	SCOPERTO - ATTIVITA'	MASSIMALE - ATTIVITA'	MASSIMALE PER EVENTO - ATTIVITA'	PREMIO ANNUO LORDO DI GARANZIA	SCONTO %	RATA ALLA FIRMA	RATA SUCCESSIVA	DATA SCADENZA I RATA	NUMERO POLIZZA
					fmt.Println(g)
					//row := []string{"TIPO OPERAZIONE", "N. POLIZZA SOSTITUITA", "DENOMINAZIONE PRODOTTO", "NODO DI GESTIONE", "DATA EMISSIONE", "DATA EFFETTO", "PARTITA IVA CONTRAENTE", "CODICE FISCALE CONTRAENTE", "NATURA GIURIDICA CONTRAENTE", "RAGIONE SOCIALE CONTRAENTE", "PROVINCIA CONTRAENTE", "COMUNE CONTRAENTE", "CAP CONTRAENTE", "TOPONIMO CONTRAENTE", "INDIRIZZO CONTRAENTE", "NUMERO CIVICO CONTRAENTE", "DATA SCADENZA", "FRAZIONAMENTO", "VINCOLO", "NUMERO ADDETTI", "COSA SI VUOLE ASSICURARE", "DOMANDA 1", "DOMANDA 2", "DOMANDA 3", "FATTURATO", "FORMA DI COPERTURA", "FORMULA INCENDIO", "BENE", "ANNO DI COSTRUZIONE FABBRICATO", "MATERIALE COSTRUZIONE", "NUMERO PIANI", "PRESENZA ALLARME", "PRESENZA POLIZZA CONDOMINIALE", "TIPOLOGIA FABBRICATO", "PROVINCIA UBICAZIONE", "COMUNE UBICAZIONE", "CAP UBICAZIONE", "TOPONIMO UBICAZIONE", "INDIRIZZO UBICAZIONE", "NUMERO CIVICO UBICAZIONE", "CODICE ATTIVITA' - BENI", "CLASSE - SOLO BENI", "SETTORE - BENI", "TIPO - BENI", "CLAUSOLA VINCOLO", "TESTO CLAUSOLA VINCOLO", "GARANZIE/PACCHETTI - BENI", "FRANCHIGIA - BENI", "SOMMA ASSICURATA - BENI", "SCOPERTO - BENI", "% SOMMA ASSICURATA INCENDIO FABBRICATO E CONTENUTO - BENI", "MASSIMALE - BENI", "DIARIA - BENI", "CODICE ATTIVITA' - ATTIVITA'", "CLASSE - ATTIVITA'", "SETTORE - ATTIVITA'", "TIPO - ATTIVITA'", "GARANZIE/PACCHETTI - ATTIVITA'", "FRANCHIGIA - ATTIVITA'", "SCOPERTO - ATTIVITA'", "MASSIMALE - ATTIVITA'", "MASSIMALE PER EVENTO - ATTIVITA'", "PREMIO ANNUO LORDO DI GARANZIA", "SCONTO %", "RATA ALLA FIRMA", "RATA SUCCESSIVA", "DATA SCADENZA I RATA", "NUMERO POLIZZA"}
					row := []string{
						"Nuova emissione",                   //TIPO OPERAZIONE
						"",                                  //N. POLIZZA SOSTITUITA
						"WOPTA PER TE. ARTIGIANI & IMPRESE", //DENOMINAZIONE PRODOTTO
						"0920",                              //NODO DI GESTIONE
						policy.CreationDate.Format(layout),  //DATA EMISSIONE
						policy.StartDate.Format(layout),     //DATA EFFETTO
						policy.Contractor.VatCode,           //PARTITA IVA CONTRAENTE
						policy.Contractor.FiscalCode,        //CODICE FISCALE CONTRAENTE
						"",                                  //NATURA GIURIDICA CONTRAENTE
						enterpriseName,                      //RAGIONE SOCIALE CONTRAENTE
						policy.Contractor.Surname,           //COGNOME CONTRANTE
						policy.Contractor.Name,              //NOME CONTRANTE
						policy.Contractor.CityCode,          //PROVINCIA CONTRAENTE
						policy.Contractor.Locality,          //COMUNE CONTRAENTE
						policy.Contractor.PostalCode,        //CAP CONTRAENTE
						"",                                  //TOPONIMO CONTRAENTE
						policy.Contractor.Address,           //INDIRIZZO CONTRAENTE
						policy.Contractor.StreetNumber,      //NUMERO CIVICO CONTRAENTE
						policy.EndDate.Format(layout),       //DATA SCADENZA
						getMapSplit(policy.PaymentSplit),    //FRAZIONAMENTO
						"NO",                                //VINCOLO
						"",                                  //CONVENZIONE
						"Diretto online",                    //CANALE
						"",                                  //DEROGA
						employer,                            //NUMERO ADDETTI
						"3",                                 //COSA SI VUOLE ASSICURARE
						"1",                                 //DOMANDA 1
						"1",                                 //DOMANDA 2
						"1",                                 //DOMANDA 3
						getMapRevenue(revenue),              //FATTURATO
						"1",                                 //FORMA DI COPERTURA ------------------------------------------bENI
						"2",                                 //FORMULA INCENDIO
						"1",                                 //BENE -----------------------------------------------------BENI 1 FABBRICATO
						getMapBuildingYear(asset.Building.BuildingYear),         //ANNO DI COSTRUZIONE FABBRICATO
						getMapBuildingMaterial(asset.Building.BuildingMaterial), //MATERIALE COSTRUZIONE
						getMapBuildingFloor(asset.Building.Floor),               //NUMERO PIANI
						getOneTwo(asset.Building.IsAllarm),                      //PRESENZA ALLARME
						"",                                                      //PRESENZA POLIZZA CONDOMINIALE
						getOneTwo(asset.Building.IsHolder),                      //TIPOLOGIA FABBRICATO
						asset.Building.CityCode,                                 //PROVINCIA UBICAZIONE
						asset.Building.Locality,                                 //COMUNE UBICAZIONE
						asset.Building.PostalCode,                               //CAP UBICAZIONE
						"",                                                      //TOPONIMO UBICAZIONE
						asset.Building.Address,                                  //INDIRIZZO UBICAZIONE
						asset.Building.StreetNumber,                             //NUMERO CIVICO UBICAZIONE
						ateco,                                                   //CODICE ATTIVITA' – BENI
						class,                                                   //CLASSE - SOLO BENI
						sector,                                                  //SETTORE – BENI
						atecoDesc,                                               //TIPO – BENI
						"",                                                      //CLAUSOLA VINCOLO
						"",                                                      //TESTO CLAUSOLA VINCOLO
						g.CompanyCodec,                                          //GARANZIE/PACCHETTI – BENI
						"",                                                      //ESTENSIONE RC DM 37/2008
						"",                                                      //CLAUSOLA BENI - BENE
						"",                                                      //CLAUSOLA BENI - GARANZIE
						g.Deductible,                                            //FRANCHIGIA – BENI
						strconv.Itoa(int(g.SumInsuredLimitOfIndemnity)), //SOMMA ASSICURATA – BENI

						getMapSelfInsurance(g.SelfInsuranceDesc), //SCOPERTO – BENI
						"",                                       //% SOMMA ASSICURATA INCENDIO FABBRICATO E CONTENUTO – BENI
						strconv.Itoa(int(g.SumInsuredLimitOfIndemnity)), //MASSIMALE - BENI
						"",                                //DIARIA – BENI
						"",                                //CODICE ATTIVITA' - ATTIVITA' -------------------------------------------------ATTIVITA 2 ATTIVITA
						"",                                //CLASSE - ATTIVITA'
						"",                                //SETTORE - ATTIVITA'
						"",                                //TIPO - ATTIVITA'
						"",                                //GARANZIE/PACCHETTI - ATTIVITA'
						"",                                //CLAUSOLA ATTIVITA' - BENE
						"",                                //CLAUSOLA ATTIVITA' - GARANZIE
						"",                                //FRANCHIGIA - ATTIVITA'
						"",                                //SCOPERTO - ATTIVITA'
						"",                                //MASSIMALE - ATTIVITA'
						"",                                //MASSIMALE PER EVENTO - ATTIVITA'
						fmt.Sprintf("%.2f", g.PriceGross), //PREMIO ANNUO LORDO DI GARANZIA
						"0",                               //SCONTO %
						fmt.Sprintf("%.2f", getInstallament(policy.PaymentSplit, g.PriceGross)), //RATA ALLA FIRMA
						fmt.Sprintf("%.2f", getInstallament(policy.PaymentSplit, g.PriceGross)), //RATA SUCCESSIVA
						getInstallamentDate(policy, layout),                                     //DATA SCADENZA I RATA
						policy.CodeCompany,
					}
					result = append(result, row)

				}

			}
			if asset.Enterprise != nil {
				for _, g := range asset.Guarantees {
					fmt.Println(g)
					//TIPO OPERAZIONE	N. POLIZZA SOSTITUITA	DENOMINAZIONE PRODOTTO	NODO DI GESTIONE	DATA EMISSIONE	DATA EFFETTO	PARTITA IVA CONTRAENTE	CODICE FISCALE CONTRAENTE	NATURA GIURIDICA CONTRAENTE	RAGIONE SOCIALE CONTRAENTE	PROVINCIA CONTRAENTE	COMUNE CONTRAENTE	CAP CONTRAENTE	TOPONIMO CONTRAENTE	INDIRIZZO CONTRAENTE	NUMERO CIVICO CONTRAENTE	DATA SCADENZA	FRAZIONAMENTO	VINCOLO	NUMERO ADDETTI	COSA SI VUOLE ASSICURARE	DOMANDA 1	DOMANDA 2	DOMANDA 3	FATTURATO	FORMA DI COPERTURA	FORMULA INCENDIO	BENE	ANNO DI COSTRUZIONE FABBRICATO	MATERIALE COSTRUZIONE	NUMERO PIANI	PRESENZA ALLARME	PRESENZA POLIZZA CONDOMINIALE	TIPOLOGIA FABBRICATO	PROVINCIA UBICAZIONE	COMUNE UBICAZIONE	CAP UBICAZIONE	TOPONIMO UBICAZIONE	INDIRIZZO UBICAZIONE	NUMERO CIVICO UBICAZIONE	CODICE ATTIVITA' - BENI	CLASSE - SOLO BENI	SETTORE - BENI	TIPO - BENI	CLAUSOLA VINCOLO	TESTO CLAUSOLA VINCOLO	GARANZIE/PACCHETTI - BENI	FRANCHIGIA - BENI	SOMMA ASSICURATA - BENI	SCOPERTO - BENI	% SOMMA ASSICURATA INCENDIO FABBRICATO E CONTENUTO - BENI	MASSIMALE - BENI	DIARIA - BENI	CODICE ATTIVITA' - ATTIVITA'	CLASSE - ATTIVITA'	SETTORE - ATTIVITA'	TIPO - ATTIVITA'	GARANZIE/PACCHETTI - ATTIVITA'	FRANCHIGIA - ATTIVITA'	SCOPERTO - ATTIVITA'	MASSIMALE - ATTIVITA'	MASSIMALE PER EVENTO - ATTIVITA'	PREMIO ANNUO LORDO DI GARANZIA	SCONTO %	RATA ALLA FIRMA	RATA SUCCESSIVA	DATA SCADENZA I RATA	NUMERO POLIZZA
					//row := []string{"TIPO OPERAZIONE", "N. POLIZZA SOSTITUITA", "DENOMINAZIONE PRODOTTO", "NODO DI GESTIONE", "DATA EMISSIONE", "DATA EFFETTO", "PARTITA IVA CONTRAENTE", "CODICE FISCALE CONTRAENTE", "NATURA GIURIDICA CONTRAENTE", "RAGIONE SOCIALE CONTRAENTE", "PROVINCIA CONTRAENTE", "COMUNE CONTRAENTE", "CAP CONTRAENTE", "TOPONIMO CONTRAENTE", "INDIRIZZO CONTRAENTE", "NUMERO CIVICO CONTRAENTE", "DATA SCADENZA", "FRAZIONAMENTO", "VINCOLO", "NUMERO ADDETTI", "COSA SI VUOLE ASSICURARE", "DOMANDA 1", "DOMANDA 2", "DOMANDA 3", "FATTURATO", "FORMA DI COPERTURA", "FORMULA INCENDIO", "BENE", "ANNO DI COSTRUZIONE FABBRICATO", "MATERIALE COSTRUZIONE", "NUMERO PIANI", "PRESENZA ALLARME", "PRESENZA POLIZZA CONDOMINIALE", "TIPOLOGIA FABBRICATO", "PROVINCIA UBICAZIONE", "COMUNE UBICAZIONE", "CAP UBICAZIONE", "TOPONIMO UBICAZIONE", "INDIRIZZO UBICAZIONE", "NUMERO CIVICO UBICAZIONE", "CODICE ATTIVITA' - BENI", "CLASSE - SOLO BENI", "SETTORE - BENI", "TIPO - BENI", "CLAUSOLA VINCOLO", "TESTO CLAUSOLA VINCOLO", "GARANZIE/PACCHETTI - BENI", "FRANCHIGIA - BENI", "SOMMA ASSICURATA - BENI", "SCOPERTO - BENI", "% SOMMA ASSICURATA INCENDIO FABBRICATO E CONTENUTO - BENI", "MASSIMALE - BENI", "DIARIA - BENI", "CODICE ATTIVITA' - ATTIVITA'", "CLASSE - ATTIVITA'", "SETTORE - ATTIVITA'", "TIPO - ATTIVITA'", "GARANZIE/PACCHETTI - ATTIVITA'", "FRANCHIGIA - ATTIVITA'", "SCOPERTO - ATTIVITA'", "MASSIMALE - ATTIVITA'", "MASSIMALE PER EVENTO - ATTIVITA'", "PREMIO ANNUO LORDO DI GARANZIA", "SCONTO %", "RATA ALLA FIRMA", "RATA SUCCESSIVA", "DATA SCADENZA I RATA", "NUMERO POLIZZA"}
					row := []string{
						"Nuova emissione",                        //TIPO OPERAZIONE
						"",                                       //N. POLIZZA SOSTITUITA
						"WOPTA PER TE. ARTIGIANI & IMPRESE",      //DENOMINAZIONE PRODOTTO
						"0920",                                   //NODO DI GESTIONE
						policy.CreationDate.Format(layout),       //DATA EMISSIONE
						policy.StartDate.Format(layout),          //DATA EFFETTO
						policy.Contractor.VatCode,                //PARTITA IVA CONTRAENTE
						policy.Contractor.FiscalCode,             //CODICE FISCALE CONTRAENTE
						"",                                       //NATURA GIURIDICA CONTRAENTE
						enterpriseName,                           //RAGIONE SOCIALE CONTRAENTE
						policy.Contractor.Name,                   //COGNOME CONTRANTE
						policy.Contractor.Surname,                //NOME CONTRANTE
						policy.Contractor.City,                   //PROVINCIA CONTRAENTE
						policy.Contractor.Locality,               //COMUNE CONTRAENTE
						policy.Contractor.CityCode,               //CAP CONTRAENTE
						"",                                       //TOPONIMO CONTRAENTE
						policy.Contractor.Address,                //INDIRIZZO CONTRAENTE
						policy.Contractor.StreetNumber,           //NUMERO CIVICO CONTRAENTE
						policy.EndDate.Format(layout),            //DATA SCADENZA
						getMapSplit(policy.PaymentSplit),         //FRAZIONAMENTO
						"NO",                                     //VINCOLO
						"",                                       //CONVENZIONE
						"Diretto online",                         //CANALE
						"",                                       //DEROGA
						employer,                                 //NUMERO ADDETTI
						"3",                                      //COSA SI VUOLE ASSICURARE
						"1",                                      //DOMANDA 1
						"1",                                      //DOMANDA 2
						"1",                                      //DOMANDA 3
						getMapRevenue(revenue),                   //FATTURATO
						"1",                                      //FORMA DI COPERTURA ------------------------------------------bENI
						"2",                                      //FORMULA INCENDIO
						"2",                                      //BENE -----------------------------------------------------BENI 1 FABBRICATO
						"",                                       //ANNO DI COSTRUZIONE FABBRICATO
						"",                                       //MATERIALE COSTRUZIONE
						"",                                       //NUMERO PIANI
						"",                                       //PRESENZA ALLARME
						"",                                       //PRESENZA POLIZZA CONDOMINIALE
						"",                                       //TIPOLOGIA FABBRICATO
						"",                                       //PROVINCIA UBICAZIONE
						"",                                       //COMUNE UBICAZIONE
						"",                                       //CAP UBICAZIONE
						"",                                       //TOPONIMO UBICAZIONE
						"",                                       //INDIRIZZO UBICAZIONE
						"",                                       //NUMERO CIVICO UBICAZIONE
						"",                                       //CODICE ATTIVITA' – BENI
						"",                                       //CLASSE - SOLO BENI
						"",                                       //SETTORE – BENI
						"",                                       //TIPO – BENI
						"",                                       //CLAUSOLA VINCOLO
						"",                                       //TESTO CLAUSOLA VINCOLO
						"",                                       //GARANZIE/PACCHETTI – BENI
						"",                                       //ESTENSIONE RC DM 37/2008
						"",                                       //CLAUSOLA BENI - BENE
						"",                                       //CLAUSOLA BENI - GARANZIE
						"",                                       //FRANCHIGIA – BENI
						"",                                       //SOMMA ASSICURATA – BENI
						"",                                       //SCOPERTO – BENI
						"",                                       //% SOMMA ASSICURATA INCENDIO FABBRICATO E CONTENUTO – BENI
						"",                                       //MASSIMALE - BENI
						"",                                       //DIARIA – BENI
						ateco,                                    //CODICE ATTIVITA' - ATTIVITA' -------------------------------------------------ATTIVITA 2 ATTIVITA
						class,                                    //CLASSE - SOLO BENI
						sector,                                   //SETTORE – BENI
						atecoDesc,                                //TIPO - ATTIVITA'
						g.CompanyCodec,                           //GARANZIE/PACCHETTI - ATTIVITA'
						"",                                       //CLAUSOLA ATTIVITA' - BENE
						"",                                       //CLAUSOLA ATTIVITA' - GARANZIE
						g.Deductible,                             //FRANCHIGIA - ATTIVITA'
						getMapSelfInsurance(g.SelfInsuranceDesc), //SCOPERTO - ATTIVITA'
						strconv.Itoa(int(g.SumInsuredLimitOfIndemnity)), //MASSIMALE - ATTIVITA'
						"",                                //MASSIMALE PER EVENTO - ATTIVITA'
						fmt.Sprintf("%.2f", g.PriceGross), //PREMIO ANNUO LORDO DI GARANZIA
						"0",                               //SCONTO %
						fmt.Sprintf("%.2f", getInstallament(policy.PaymentSplit, g.PriceGross)), //RATA ALLA FIRMA
						fmt.Sprintf("%.2f", getInstallament(policy.PaymentSplit, g.PriceGross)), //RATA SUCCESSIVA
						getInstallamentDate(policy, layout),                                     //DATA SCADENZA I RATA
						policy.CodeCompany,                                                      //NUMERO POLIZZA
					}

					result = append(result, row)

				}

			}

		}

		if e == nil {
			policy.CompanyEmitted = true
			//lib.SetFirestore("policy", policy.Agent.Uid, policy)
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

func getInstallamentDate(p models.Policy, layout string) string {
	var res string
	res = p.EndDate.Format(layout)
	if p.PaymentSplit == "monthly" {
		res = p.StartDate.AddDate(0, 1, 0).Format(layout)
	}

	return res
}
func getInstallament(key string, price float64) float64 {
	var res float64
	res = price
	if key == "monthly" {
		res = price / 12
	}
	return res
}
func getYesNo(key bool) string {
	var res string
	mapGarante := map[bool]string{true: "SI", false: "NO"}

	if seconds, ok := mapGarante[key]; ok { // will be false if person is not in the map
		res = seconds
	}
	return res
}
func getOneTwo(key bool) string {
	var res string
	mapGarante := map[bool]string{true: "1", false: "2"}

	if seconds, ok := mapGarante[key]; ok { // will be false if person is not in the map
		res = seconds
	}
	return res
}
func getMapBuildingYear(key string) string {
	var res string
	mapGarante := map[string]string{"before1972": "1", "1972between2009": "2", "after2009": "3"}

	if seconds, ok := mapGarante[key]; ok { // will be false if person is not in the map
		res = seconds
	}
	return res
}
func getMapBuildingFloor(key string) string {
	var res string
	mapGarante := map[string]string{"ground_floor": "1", "first": "2", "second": "3", "greater_than_second": "4"}

	if seconds, ok := mapGarante[key]; ok { // will be false if person is not in the map
		res = seconds
	}
	return res
}
func getMapBuildingMaterial(key string) string {
	var res string
	mapGarante := map[string]string{"masonry": "1", "reinforcedConcrete": "2", "antiSeismicLaminatedTimber": "3", "steel": "4"}

	if seconds, ok := mapGarante[key]; ok { // will be false if person is not in the map
		res = seconds
	}
	return res
}
func getMapSplit(key string) string {
	var res string
	res = "1"
	if key == "monthly" {
		res = "12"
	}
	return res
}
func getBuildingType(key string) string {
	var res string
	res = "1"
	if key == "montly" {
		res = "12"
	}
	return res
}
func getMapRevenue(key int) string {
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
	return strconv.Itoa(res)
}
func getMapSelfInsurance(key string) string {
	var res int

	if key == "5% - minimo € 500" { // will be false if person is not in the map
		res = 1
	}
	if key == "5% - minimo € 1.000" { // will be false if person is not in the map
		res = 2
	}
	if key == "5% - minimo € 1.500" { // will be false if person is not in the map
		res = 3
	}
	if key == "10% - minimo € 500" { // will be false if person is not in the map
		res = 4
	}
	if key == "10% - minimo € 1.000" { // will be false if person is not in the map
		res = 5
	}
	if key == "10% - minimo € 1.500" { // will be false if person is not in the map
		res = 6
	}
	if key == "10% - minimo € 2.000" { // will be false if person is not in the map
		res = 7
	}
	if key == "10% - minimo € 3.000" { // will be false if person is not in the map
		res = 8
	}
	if key == "10% - minimo € 5.000" { // will be false if person is not in the map
		res = 9
	}
	if key == "15% - minimo € 5.000" { // will be false if person is not in the map
		res = 10
	}
	if key == "10% - minimo € 10.000" { // will be false if person is not in the map
		res = 11
	}
	if key == "10% - minimo € 20.000" { // will be false if person is not in the map
		res = 12
	}
	if key == "10% - minimo € 25.000" { // will be false if person is not in the map
		res = 13
	}
	if key == "10% - minimo € 30.000" { // will be false if person is not in the map
		res = 14
	}

	return strconv.Itoa(res)
}
