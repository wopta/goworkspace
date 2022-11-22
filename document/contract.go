package document

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/props"
	lib "github.com/wopta/goworkspace/lib"
	model "github.com/wopta/goworkspace/models"
)

func Contract(w http.ResponseWriter, r *http.Request) (string, interface{}) {
	log.Println("Contract")
	//lib.Files("./serverless_function_source_code")
	req := lib.ErrorByte(ioutil.ReadAll(r.Body))
	var data model.Policy
	defer r.Body.Close()
	err := json.Unmarshal([]byte(req), &data)
	lib.CheckError(err)
	resp, respObj := ContractObj(data)
	return resp, respObj
}

func ContractObj(data model.Policy) (string, DodumentResponse) {
	skin, _, textBold, _, _ := getVar()
	log.Println(textBold)
	log.Println("Document 1")
	m := skin.initDefault()
	m.RegisterHeader(func() {
		m.Row(15.0, func() {
			m.Col(2, func() {

				_ = m.FileImage(lib.GetAssetPathByEnv("document")+"/logo_persona.png", props.Rect{
					Left:    1,
					Top:     1,
					Center:  false,
					Percent: 100,
				})
			})
			m.Col(1, func() {
				m.Text("Wopta per te", props.Text{
					Color:       skin.LineColor,
					Top:         1,
					Style:       consts.Bold,
					Align:       consts.Left,
					Size:        skin.SizeTitle + 3,
					Extrapolate: true,
				})

				m.Text("Persona", props.Text{
					Top:         6,
					Style:       consts.Italic,
					Align:       consts.Left,
					Color:       skin.TextColor,
					Size:        skin.SizeTitle + 3,
					Extrapolate: true,
				})
			})
			m.ColSpace(6)
			m.Col(2, func() {
				_ = m.FileImage(lib.GetAssetPathByEnv("document")+"/ARTW_LOGO_RGB_400px.png", props.Rect{
					Left:    1,
					Top:     1,
					Center:  false,
					Percent: 100,
				})
			})
		})
		h := []string{"I dati della tua Polizza ", "I tuoi dati"}
		var tablePremium [][]string
		tablePremium = append(tablePremium, []string{"Numero: " + data.ID, "Contraente: " + data.Contractor.Name + " " + data.Contractor.Surname})
		tablePremium = append(tablePremium, []string{"Decorre dal: " + data.StartDate.String() + " ore 24:00", "C.F. / P.IVA: " + data.Contractor.Surname})
		tablePremium = append(tablePremium, []string{"Scade il: " + data.EndDate.String() + " ore 24:00", "Indirizzo: " + data.Contractor.Address})
		tablePremium = append(tablePremium, []string{"Si rinnova a scadenza, salvo disdetta da inviare 30 giorni prima", "XXXXX  XXXXXXXXXXXXXXXXXXX (XX)"})
		tablePremium = append(tablePremium, []string{"Prossimo pagamento il: " + data.EndDate.String(), "Mail:  " + data.Contractor.Mail})
		tablePremium = append(tablePremium, []string{"Sostituisce la polizza: = = = = = = = =", "Telefono: " + data.Contractor.Phone})
		m = skin.Table(m, h, tablePremium, 6, 3.0)
	})

	m.RegisterFooter(func() {
		m.Row(15.0, func() {
			m.Col(8, func() {
				m.Text("Wopta per te. Persona è un prodotto assicurativo di Global Assistance Compagnia di assicurazioni e riassicurazioni S.p.A, distribuito da Wopta Assicurazioni S.r.l", props.Text{
					Top:         1,
					Style:       consts.Bold,
					Align:       consts.Left,
					Color:       skin.LineColor,
					Size:        skin.Size - 1,
					Extrapolate: false,
				})
			})
			m.Col(2, func() {
				_ = m.FileImage(lib.GetAssetPathByEnv("document")+"/logo_global.png", props.Rect{
					Left:    1,
					Top:     1,
					Center:  false,
					Percent: 100,
				})
			})
		})
	})
	log.Println("Document 3")
	m = skin.Space(m, 10.0)

	log.Println("Document 4")
	if data.Name == "persona" {
		m = skin.GetPersona(data, m)
	}

	if data.Name == "pmi" {
		m = skin.GetPmi(data, m)
	}
	m = skin.Space(m, 10.0)
	var table [][]string
	h := []string{"Garanzie ", "Somma assicurata ", "Opzioni / Dettagli ", "Premio "}
	for _, A := range data.Assets {
		for _, k := range A.Guarantees {
			r := []string{k.Name, strconv.Itoa(int(k.SumInsuredLimitOfIndemnity)), k.SelfInsurance, strconv.Itoa(int(k.Price))}
			table = append(table, r)
		}
	}

	m = skin.CoveragesTable(m, h, table)

	m.AddPage()
	articles := []Kv{
		{
			Key:   "1. ",
			Value: "le dichiarazioni non veritiere, inesatte o reticenti, da me rese, possono compromettere il diritto alla prestazione (come da art. 1892, 1893, 1894 c.c.)"},
		{
			Key:   "2.",
			Value: "nel caso di coperture che richiedono di acquisire informazioni sullo stato di salute dell’assicurato, come nel presente contratto: a) prima della sottoscrizione, ho verificato l’esattezza e rispondenza a verità delle mie dichiarazioni qui riportate; b) sono a conoscenza di poter chiedere di essere sottoposto a visita medica per certificare l’effettivo mio stato di salute, con costi a mio carico; "},
	}
	stantments := []Kv{
		{
			Key:   "a) ",
			Value: "di NON essere affetto da infermità gravi quali: alcoolismo, tossicodipendenza, sindrome da immunodeficienza acquisita (AIDS), ovvero infermità dovute a malattie del sistema nervoso o della psiche (schizofrenia, psicosi, depressione, nevrosi, insufficienza mentale, demenza, Alzheimer, Parkinson, SLA, sclerosi multipla, cerebropatie, paresi, paralisi, epilessia); "},
		{
			Key:   "b) ",
			Value: "di NON essere affetto da Difetti Fisici gravi ed invalidanti, da infermità e/o Invalidità Permanenti con postumi valutati in misura superiore al 50%;  "},
		{
			Key:   "c) ",
			Value: "di NON aver avuto precedenti Polizze infortuni annullate, per iniziativa di compagnie, prima della loro naturale scadenza;  "},
		{
			Key:   "d) ",
			Value: "rispetto al Contraente di essere: socio, membro del consiglio di amministrazione, dipendente, collaboratore (anche esterno);   "},
	}

	stantments2 := []Kv{
		{
			Key:   "a. ",
			Value: "sono alto XXX cm e peso YY Kg  "},
		{
			Key:   "b. ",
			Value: "NON assumo o NON ho assunto negli ultimi 15 anni sostanze stupefacenti "},
		{
			Key:   "c. ",
			Value: "NON consumo abitualmente alcolici in misura pari o superiore a ad 1 litro di vino e/o di birra e/o un quarto di litro di superalcolici (bevande oltre 21 gradi alcolici) al giorno  "},
		{
			Key:   "d. ",
			Value: "NON fumo più di 10 sigarette al giorno "},
		{
			Key:   "e. ",
			Value: "NON ho subito infortuni o malattie, negli ultimi cinque anni, che mi hanno impedito di svolgere la mia professione per più di due settimane. L’inabilità o invalidità è insorta nel e durata per: xxxxxxxxxxxxxxxxxxxxxxxx  "},
		{
			Key:   "f. ",
			Value: "NESSUNA/una malattia e/o infortunio (o loro postumi), attualmente mi impedisce di svolgere la tua mia professione. Nel dettaglio la malattia /infortunio all’origine dell’inabilità o invalidità è: xxxxxxxxxxxxxxxxxxxx "},
		{
			Key:   "g. ",
			Value: "NON soffro di malattia acuta o cronica del sistema cardiocircolatorio, dell’apparato respiratorio, del sistema nervoso, dell’apparato digerente, del sangue, delle vie urinarie e genitourinarie, del sistema endocrino metabolico, dell’apparato muscolo-scheletrico, di tumori maligni. Nel dettaglio la malattia di cui soffro è: xxxxxxxxxxxxxxxxxx  "},
	}
	m = skin.Space(m, 20.0)
	m = skin.Articles(m, articles)

	m = skin.Stantements(m, "Ai fini dell’efficacia di tutte delle le garanzie, ", stantments)

	m = skin.Stantements(m, "Con specifico riferimento alla garanzia Invalidità Permanente da Malattia, ", stantments2)
	m = skin.Space(m, 5.0)
	m.Row(10, func() {
		m.Col(12, func() {
			m.Text("OPZIONE KEY MAN: DICHIARO di fornire il mio consenso ai sensi dell’art. 1919 c.c. a designare il Contraente come Beneficiario, per le garanzie Invalidità Permanente e morte, nella percentuale del XXX% indicata nella scheda delle garanzie, alle condizioni previste dagli Artt. 21.3 e 22.3 della delle Condizioni di Assicurazione della presente Polizza.  ", textBold)
		})
		//m.SetBackgroundColor(magenta)
	})
	m = skin.Space(m, 1.0)
	m = skin.Sign(m, data.Contractor.Name+" "+data.Contractor.Surname, "Assicurato ")
	m.AddPage()
	m = skin.Space(m, 5.0)
	title := "Condizioni Speciali in deroga alle Condizioni Generali di Assicurazione "
	sub := "In deroga a quanto riportato nelle Condizioni Generali di Assicurazione, si concorda tra le Parti che: "
	body := "TXT libero Fermo tutto il resto non derogato da quanto precede.  "
	m = skin.TitleSub(m, title, sub, body)
	title = "Presa visione dei documenti precontrattuali e sottoscrizione Polizza "
	body = "Ho scelto la ricezione della seguente documentazione su supporto cartaceo / via e-mail al seguente indirizzo: XXXXXXXXXX. Sono a conoscenza che, anche le future comunicazioni avverranno con questo mezzo e che qualora volessi modificare questa mia scelta potrò farlo scrivendo a Global Assistance, con le modalità previste nelle Condizioni Generali di Assicurazione.  "
	m = skin.Title(m, title, body, 18.0)
	confirmationRecepit := []string{
		"1. degli Allegati 3, 4 e 4-ter, di cui al Regolamento IVASS n. 40/2018, relativi agli obblighi informativi e di comportamento dell’Intermediario, inclusa l’informativa privacy dell’intermediario (ai sensi dell’art. 13 del regolamento UE n. 2016/679); ",
		"2. del Set informativo, identificato dal modello XXXXXXXX ed. 2022, contenente: 1) documento informativo per i prodotti assicurativi danni (DIP Danni) e documento informativo precontrattuale aggiuntivo per i prodotti assicurativi danni (DIP Aggiuntivo danni) cui al Regolamento IVASS n. 41/2018; 2) Condizioni di Assicurazione comprensive di Glossario, che dichiaro altresì di conoscere ed accettare. ",
	}

	m = skin.TitleList(m, "", confirmationRecepit)
	m = skin.Space(m, 3.0)
	m = skin.Sign(m, "data.Name"+" "+"data.Surname", "Global Assistance")
	m = skin.Space(m, 3.0)
	title = "Le clausole della Polizza da approvare in modo specifico "
	body = `Ai sensi degli artt. 1341 e 1342 Codice Civile, dichiaro di
	 approvare in modo specifico, le disposizioni indicate nelle condizioni di
	  assicurazione con particolare riguardo agli articoli dei seguenti capitoli: 
	Art. 5 Foro competente; Art. 30 Denuncia e obblighi in caso di Sinistro Infortuni; 
	Art. 32 Controversie: Arbitrato irrituale; Art. 35.1 Invalidità Permanente da Infortunio; 
	Art. 36.1 Gestione del caso assicurativo; Art. 38 Denuncia e obblighi in caso di sinistro 
	Invalidità Permanente da Malattia; Art. 38.3 Criteri di liquidazione dell’Invalidità Permanente da Malattia; 
	Art. 38.4 Valutazione del danno – ricorso all’Arbitrato`

	m = skin.Title(m, title, body, 25.0)

	h = []string{"Premio ", "Imponibile  ", "Imposte Assicurative ", "Totale"}
	var tablePremium [][]string
	tablePremium = append(tablePremium, []string{"Annuale", strconv.Itoa(int(data.PriceNett)), strconv.Itoa(int(data.TaxAmount)), strconv.Itoa(int(data.PriceGross))})
	tablePremium = append(tablePremium, []string{"Mensile", strconv.Itoa(int(data.PriceNett)), strconv.Itoa(int(data.TaxAmount)), strconv.Itoa(int(data.PriceGross))})
	m = skin.Space(m, 10.0)
	m = skin.TableLine(m, h, tablePremium)
	title = "Come puoi pagare il premio "
	body = `I mezzi di pagamento consentiti nei confronti di Wopta sono esclusivamente 
	bonifico e strumenti di pagamento elettronico, quali ad esempio, carte di 
	credito e/o carte di debito, incluse le carte prepagate.`
	m = skin.Space(m, 5.0)
	m = skin.Title(m, title, body, 25.0)
	title = "Emissione polizza e pagamento della prima rata "
	body = `Polizza emessa a Milano il 00/00/0000 per un importo di euro XXX,XX quale prima rata alla firma,
	 il cui pagamento a saldo è da effettuarsi con i metodi di pagamento sopra indicati. 
	Costituisce quietanza di pagamento la mail di conferma che Wopta invierà al Contraente. `

	m = skin.Title(m, title, body, 25.0)
	m = skin.Sign(m, "data.Name"+" "+"data.Surname", "Wopta Assicurazioni")

	m.AddPage()
	m = skin.Space(m, 10.0)
	title = "Per noi questa polizza fa al caso tuo "
	sub = "Richieste ed esigenze di copertura assicurativa del contraente "
	body = "In Polizza sono riportate le tue dichiarazioni relative al rischio. Sulla base di tali dichiarazioni, esigenze e richieste, le soluzioni assicurative individuate in coerenza con esse, per ogni Assicurato, sono: .  "
	m = skin.TitleSub(m, title, sub, body)
	customerList := []Kv{{Key: "1. ", Value: "Cognome e nome: copertura per gli infortuni occorsi AA in qualità di BB, che offre CC1 FR CC2 DD EE FF GG HH II JJ KK LL "},
		{Key: "", Value: ""}}
	m = skin.BulletList(m, customerList)
	m = skin.RowCol1(m, "", consts.Normal)
	m = skin.Sign(m, data.Contractor.Name+" "+data.Contractor.Surname, "Assicurato ")
	aboutUs := []Kv{{
		Key:   "Wopta Assicurazioni S.r.l.",
		Value: " intermediario assicurativo, soggetto al controllo dell’IVASS ed iscritto dal 14.02.2022 al Registro Unico degli Intermediari, in Sezione A nr. A000701923, avente sede legale in Galleria del Corso, 1 – 20122 Milano (MI). Capitale sociale Euro 120.000 - Codice Fiscale, Reg. Imprese e Partita IVA: 12072020964 - Iscritta al Registro delle imprese di Milano – REA MI 2638708 ",
	}, {Key: "Global Assistance Compagnia di assicurazioni e riassicurazioni S.p.A. a Socio Unico",
		Value: "5.000.000 i.v. Codice Fiscale, Partita IVA e Registro Imprese di Milano n. 10086540159 R.E.A. n. 1345012 della C.C.I.A.A. di Milano. Sede e Direzione Generale Piazza Diaz 6 – 20123 Milano – Italia E-mail: global.assistance@globalassistance.it PEC: globalassistancespa@legalmail.it. Società soggetta all’attività di direzione e coordinamento di Ri-Fin S.r.l., iscritta all’Albo dei gruppi assicurativi presso l’Ivass al n. 014. La Società è autorizzata all’esercizio delle Assicurazioni e Riassicurazioni con D.M. del 2/8/93 n. 19619 (G.U. 7/8/93 n. 184) e successive autorizzazioni ed è iscritta all’Albo Imprese presso l’IVASS al n. 1.00111. La Società è soggetta alla vigilanza dell’IVASS; è possibile verificare la veridicità e la regolarità delle autorizzazioni mediante l'accesso al sito www.ivass.it"}}
	m = skin.AboutUs(m, "Chi siamo ", aboutUs)
	log.Println("Document 8")
	//m.Output()
	var resp []byte
	var respObj DodumentResponse
	if os.Getenv("env") == "local" {
		err := m.OutputFileAndClose("document/contract.pdf")
		lib.CheckError(err)
	} else {
		out, err := m.Output()
		lib.CheckError(err)
		now := time.Now() // current local time
		layout := "2006-01-02"
		t, _ := time.Parse(layout, now.String())

		filename := "temp/" + data.Contractor.Name + "_" + data.Contractor.Surname + "_" + t.Format("02-Jan-2006") + "_contract.pdf"
		result := lib.PutToStorage("function-data", filename, out.Bytes())
		lib.CheckError(err)
		log.Println(result)
		data.DocumentName = filename
		log.Println(filename)
		respObj = DodumentResponse{
			LinkGcs: filename,
			Bytes:   base64.StdEncoding.EncodeToString(out.Bytes()),
		}
		resp, err = json.Marshal(respObj)
		lib.CheckError(err)
	}
	return string(resp), respObj
}
