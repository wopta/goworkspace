package document

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/johnfercher/maroto/pkg/consts"
	lib "github.com/wopta/goworkspace/lib"
	model "github.com/wopta/goworkspace/models"
)

func Contract(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println("Contract")
	//lib.Files("./serverless_function_source_code")
	req := lib.ErrorByte(ioutil.ReadAll(r.Body))
	var data model.Policy
	defer r.Body.Close()
	err := json.Unmarshal([]byte(req), &data)
	lib.CheckError(err)
	respObj := <-ContractObj(data)
	resp, err := json.Marshal(respObj)

	lib.CheckError(err)
	return string(resp), respObj, nil
}

func ContractObj(data model.Policy) <-chan DodumentResponse {
	r := make(chan DodumentResponse)

	//now := time.Now()
	//next := now.AddDate(0, 0, 4)
	//layout := "2006-01-02T15:04:05.000Z"
	layout2 := "2006-01-02"
	go func() {
		skin, _, _, _, _ := getVar()
		m := skin.initDefault()

		var (
			logo, name string
			//coverages  pdf.Maroto
			//assets     pdf.Maroto
		)
		if data.Name == "persona" {
			logo = "/persona.png"
			name = "Persona"
			m = skin.GetHeader(m, data, logo, name)
			m = skin.GetFooter(m, data, logo, name)
			m = skin.Space(m, 5.0)
			m = skin.GetPersona(data, m)
			m = skin.CoveragesPersonTable(m, data)
			m.AddPage()
		}

		if data.Name == "pmi" {
			logo = "/pmi.png"
			name = "Artigiani & Imprese"
			m = skin.GetHeader(m, data, logo, name)
			m = skin.GetFooter(m, data, logo, name)
			m = skin.Space(m, 5.0)
			m = skin.GetPmi(data, m)
			m = skin.Space(m, 5.0)
			m = skin.CoveragesPmiTable(m, data)
			skin.checkPage(m)

		}

		//var stantments []Kv
		var questions []Kv

		var alfabet []rune
		for r := 'a'; r < 'z'; r++ {
			//R := unicode.ToUpper(r)
			alfabet = append(alfabet, r)
			//log.Println(R)
		}
		for _, A := range data.Statements {

			skin.checkPage(m)
			m = skin.Title(m, A.Title, A.Question, float64(getRowHeight(A.Question, 120, 6)))
		}
		skin.checkPage(m)
		for x, A := range data.Survay {
			alfa := strconv.QuoteRune(alfabet[x])
			t := strings.Replace(alfa, "'", "", -1)
			question := Kv{
				Key:   t + ") ",
				Value: A.Question,
			}
			questions = append(questions, question)
			//m = skin.Title(m, A.Title, A.Question, 18.0)
		}

		skin.checkPage(m)
		m = skin.Stantements(m, "Ai fini dell’efficacia di tutte delle le garanzie, ", questions)
		m = skin.Space(m, 5.0)
		skin.SignDouleLine(m, data.Contractor.Name+" "+data.Contractor.Surname, "Global Assistance", "1", true)
		m = skin.Space(m, 5.0)
		skin.checkPage(m)

		title := "Le clausole della Polizza da approvare in modo specifico "
		body := `Ai sensi degli artt. 1341 e 1342 Codice Civile, dichiaro di
	 approvare in modo specifico, le disposizioni indicate nelle condizioni di
	  assicurazione con particolare riguardo agli articoli dei seguenti capitoli: 
	Art. 5 Foro competente; Art. 30 Denuncia e obblighi in caso di Sinistro Infortuni; 
	Art. 32 Controversie: Arbitrato irrituale; Art. 35.1 Invalidità Permanente da Infortunio; 
	Art. 36.1 Gestione del caso assicurativo; Art. 38 Denuncia e obblighi in caso di sinistro 
	Invalidità Permanente da Malattia; Art. 38.3 Criteri di liquidazione dell’Invalidità Permanente da Malattia; 
	Art. 38.4 Valutazione del danno – ricorso all’Arbitrato`
		skin.checkPage(m)
		m = skin.Title(m, title, body, 25.0)

		h := []string{"Premio ", "Imponibile  ", "Imposte Assicurative ", "Totale"}
		var tablePremium [][]string

		tablePremium = append(tablePremium, []string{"Annuale", "€ " + humanize.FormatFloat("#.###,##", data.PriceNett), "€ " + humanize.FormatFloat("#.###,##", data.PriceGross-data.PriceNett), "€ " + humanize.FormatFloat("#.###,##", data.PriceGross)})
		if data.PaymentSplit == "monthly" {
			tablePremium = append(tablePremium, []string{"Mensile", "€ " + humanize.FormatFloat("#.###,##", (data.PriceNett/12)), "€ " + humanize.FormatFloat("#.###,##", ((data.PriceGross-data.PriceNett)/12)), "€ " + humanize.FormatFloat("#.###,##", (data.PriceGross/12))})

		}
		m = skin.Space(m, 10.0)
		m = skin.TableLine(m, h, tablePremium)

		title = "Come puoi pagare il premio "
		body = `I mezzi di pagamento consentiti nei confronti di Wopta sono esclusivamente 
	bonifico e strumenti di pagamento elettronico, quali ad esempio, carte di 
	credito e/o carte di debito, incluse le carte prepagate.`
		skin.checkPage(m)
		m = skin.Space(m, 5.0)
		m = skin.Title(m, title, body, 25.0)
		title = "Emissione polizza e pagamento della prima rata "
		s := fmt.Sprintf("%.2f", data.PriceGross)
		body = `Polizza emessa a Milano il ` + data.StartDate.Format(layout2) + ` 00/00/0000 per un importo di euro ` + s + ` quale prima rata alla firma,
	 il cui pagamento a saldo è da effettuarsi con i metodi di pagamento sopra indicati. 
	Costituisce quietanza di pagamento la mail di conferma che Wopta invierà al Contraente. `
		skin.checkPage(m)
		m = skin.Title(m, title, body, 25.0)
		m = skin.Sign(m, "Wopta Assicurazioni", "Wopta Assicurazioni", "2", false)
		skin.checkPage(m)
		m = skin.RowCol1(m, "", consts.Normal)
		m = skin.Sign(m, data.Contractor.Name+" "+data.Contractor.Surname, "Assicurato ", "2", true)
		skin.checkPage(m)
		aboutUs := []Kv{{
			Key:   "Wopta Assicurazioni S.r.l.",
			Value: " intermediario assicurativo, soggetto al controllo dell’IVASS ed iscritto dal 14.02.2022 al Registro Unico degli Intermediari, in Sezione A nr. A000701923, avente sede legale in Galleria del Corso, 1 – 20122 Milano (MI). Capitale sociale Euro 120.000 - Codice Fiscale, Reg. Imprese e Partita IVA: 12072020964 - Iscritta al Registro delle imprese di Milano – REA MI 2638708 ",
		}, {Key: "Global Assistance Compagnia di assicurazioni e riassicurazioni S.p.A. a Socio Unico",
			Value: "5.000.000 i.v. Codice Fiscale, Partita IVA e Registro Imprese di Milano n. 10086540159 R.E.A. n. 1345012 della C.C.I.A.A. di Milano. Sede e Direzione Generale Piazza Diaz 6 – 20123 Milano – Italia E-mail: global.assistance@globalassistance.it PEC: globalassistancespa@legalmail.it. Società soggetta all’attività di direzione e coordinamento di Ri-Fin S.r.l., iscritta all’Albo dei gruppi assicurativi presso l’Ivass al n. 014. La Società è autorizzata all’esercizio delle Assicurazioni e Riassicurazioni con D.M. del 2/8/93 n. 19619 (G.U. 7/8/93 n. 184) e successive autorizzazioni ed è iscritta all’Albo Imprese presso l’IVASS al n. 1.00111. La Società è soggetta alla vigilanza dell’IVASS; è possibile verificare la veridicità e la regolarità delle autorizzazioni mediante l'accesso al sito www.ivass.it"}}
		m = skin.AboutUs(m, "Chi siamo ", aboutUs)

		//-----------Save file
		if os.Getenv("env") == "local" {
			err := m.OutputFileAndClose("document/contract.pdf")
			lib.CheckError(err)
		} else {
			out, err := m.Output()
			lib.CheckError(err)
			now := time.Now()
			timestamp := strconv.FormatInt(now.Unix(), 10)
			filename := "temp/" + data.Contractor.Name + "_" + data.Contractor.Surname + "_" + timestamp + "_contract.pdf"
			lib.PutToStorage("function-data", filename, out.Bytes())
			lib.CheckError(err)

			data.DocumentName = filename

			r <- DodumentResponse{
				LinkGcs: filename,
				Bytes:   base64.StdEncoding.EncodeToString(out.Bytes()),
			}

		}
		log.Println(data.Uid + " ContractObj end")
	}()
	return r
}
