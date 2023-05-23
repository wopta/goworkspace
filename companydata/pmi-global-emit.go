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
		result [][]string

		e error
	)

	layoutFilename := "20060102"
	location, e := time.LoadLocation("Europe/Rome")
	collection := "policy"
	fmt.Println(time.Now().In(location))
	now := time.Now().In(location).AddDate(0, 0, -1)
	filename := now.Format(layoutFilename) + "_EM_PMIW.XLSX"
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
				Field:      "company", //
				Operator:   "==",      //
				QueryValue: "global",
			},
			{
				Field:      "IsDeleted", //
				Operator:   "==",        //
				QueryValue: false,
			},
			{
				Field:      "name", //
				Operator:   "==",   //
				QueryValue: "pmi",
			},
		},
	}
	query, e := q.FirestoreWherefields(collection)
	policies := models.PolicyToListData(query)
	log.Println("len(policies):", len(policies))

	result = append(result, getPmiData(policies)...)
	log.Println("len(result):", len(result))
	filepath := "../tmp/" + filename
	excel, e := lib.CreateExcel(result, filepath, "Risultato")
	lib.PutToStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "track/global/pmi/emit/"+filepath, <-excel)
	//lib.PutGoogleStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "track/global/pmi/emit/"+filepath, source, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	if len(policies) > 0 {
		GlobalSftpDelete("/Wopta/" + filename)
		GlobalSftpUpload(filename, "/Wopta/")
	}
	return "", nil, e
}

func getPmiData(policies []models.Policy) [][]string {
	var (
		result                                    [][]string
		enterpriseName                            string
		employer, class, sector, atecoDesc, ateco string
		revenue                                   int
		sumlimitContentBuilding                   float64
		e                                         error
	)
	layout := "02/01/2006"
	for _, policy := range policies {

		for _, asset := range policy.Assets {
			if asset.Building != nil {
				for _, g := range asset.Guarantees {
					if g.Slug == "content" || g.Slug == "property-owners-liability" || g.Slug == "building" {
						sumlimitContentBuilding = sumlimitContentBuilding + g.SumInsuredLimitOfIndemnity
					}
				}
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

					sum, perc := getSumLimit(sumlimitContentBuilding, g)
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
						getDeductableMap(g.Deductible),                          //FRANCHIGIA – BENI
						sum,                                                     //SOMMA ASSICURATA – BENI

						getMapSelfInsurance(g.SelfInsuranceDesc), //SCOPERTO – BENI
						perc,                                     //% SOMMA ASSICURATA INCENDIO FABBRICATO E CONTENUTO – BENI
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
						getDeductableMap(g.Deductible),           //FRANCHIGIA - ATTIVITA'
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
			lib.SetFirestore(collection, policy.Uid, policy)
		}
	}
	return result
}
