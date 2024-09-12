package document

import (
	"fmt"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/go-pdf/fpdf"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func mainHeader(pdf *fpdf.Fpdf, policy *models.Policy, isProposal bool) {
	var (
		opt                                                                   fpdf.ImageOptions
		logoPath, cfpi, policyInfoHeader, policyInfo, expiryInfo, productName string
	)

	location, err := time.LoadLocation("Europe/Rome")
	lib.CheckError(err)

	policyStartDate := policy.StartDate.In(location)
	policyEndDate := policy.EndDate.In(location)

	if policy.PaymentSplit == string(models.PaySplitMonthly) {
		expiryInfo = "Prima scandenza mensile il: " +
			policyStartDate.AddDate(0, 1, 0).Format(dateLayout) + "\n"
	} else if policy.PaymentSplit == string(models.PaySplitYear) || policy.PaymentSplit == string(models.PaySplitYearly) {
		expiryInfo = "Prima scadenza annuale il: " +
			policyStartDate.AddDate(1, 0, 0).Format(dateLayout) + "\n"
	}

	if isProposal {
		policyInfoHeader = "I dati della tua proposta"
		policyInfo = fmt.Sprintf("Numero: %d\n", policy.ProposalNumber)
	} else {
		policyInfoHeader = "I dati della tua polizza"
		policyInfo = fmt.Sprintf("Numero: %s\n", policy.CodeCompany)
	}

	policyInfo += "Decorre dal: " + policyStartDate.Format(dateLayout) + " ore 24:00\n" +
		"Scade il: " + policyEndDate.In(location).Format(dateLayout) + " ore 24:00\n"

	switch policy.Name {
	case models.LifeProduct:
		logoPath = lib.GetAssetPathByEnvV2() + "logo_vita.png"
		productName = "Vita"
		policyInfo += expiryInfo + "Non si rinnova a scadenza."
	case models.PmiProduct:
		logoPath = lib.GetAssetPathByEnvV2() + "logo_pmi.png"
		productName = "Artigiani & Imprese"
	case models.PersonaProduct:
		logoPath = lib.GetAssetPathByEnvV2() + "logo_persona.png"
		productName = "Persona"
		policyInfo += "Si rinnova a scadenza salvo disdetta da inviare 30 giorni prima\n" + "Prossimo pagamento "
		if policy.PaymentSplit == string(models.PaySplitMonthly) {
			policyInfo += policyStartDate.In(location).AddDate(0, 1, 0).Format(dateLayout) + "\n"
		} else if policy.PaymentSplit == string(models.PaySplitYear) {
			policyInfo += policyStartDate.In(location).AddDate(1, 0, 0).Format(dateLayout) + "\n"
		}
		policyInfo += "Sostituisce la polizza ========"
	}

	contractor := policy.Contractor
	address := strings.ToUpper(contractor.Residence.StreetName + ", " + contractor.Residence.StreetNumber + "\n" +
		contractor.Residence.PostalCode + " " + contractor.Residence.City + " (" + contractor.Residence.CityCode + ")\n")

	if contractor.VatCode == "" {
		cfpi = contractor.FiscalCode
	} else {
		cfpi = contractor.VatCode
	}

	contractorInfo := "Contraente: " + strings.ToUpper(contractor.Surname+" "+contractor.Name+"\n"+
		"C.F./P.IVA: "+cfpi) + "\n" +
		"Indirizzo: " + strings.ToUpper(address) + "Mail: " + contractor.Mail + "\n" +
		"Telefono: " + contractor.Phone

	pdf.SetHeaderFunc(func() {
		opt.ImageType = "png"
		pdf.ImageOptions(logoPath, 10, 6, 13, 13, false, opt, 0, "")
		pdf.SetXY(23, 7)
		setPinkBoldFont(pdf, 18)
		pdf.Cell(10, 6, "Wopta per te")
		setPinkItalicFont(pdf, 18)
		pdf.SetXY(23, 13)
		pdf.SetTextColor(92, 89, 92)
		pdf.Cell(10, 6, productName)
		pdf.ImageOptions(lib.GetAssetPathByEnvV2()+"logo_wopta.png", 170, 6, 0, 8, false, opt, 0, "")

		setBlackBoldFont(pdf, standardTextSize)
		pdf.SetXY(11, 20)
		pdf.Cell(0, 3, policyInfoHeader)
		setBlackRegularFont(pdf, standardTextSize)
		pdf.SetXY(11, pdf.GetY()+3)
		pdf.MultiCell(0, 3.5, policyInfo, "", "", false)

		setBlackBoldFont(pdf, standardTextSize)
		pdf.SetXY(-95, 20)
		pdf.Cell(0, 3, "I tuoi dati")
		setBlackRegularFont(pdf, standardTextSize)
		pdf.SetXY(-95, pdf.GetY()+3)
		pdf.MultiCell(0, 3.5, contractorInfo, "", "", false)
		pdf.Ln(5)
	})
}

func mainFooter(pdf *fpdf.Fpdf, productName string) {
	var (
		opt                  fpdf.ImageOptions
		footerText, logoPath string
		x, y, height         float64
	)

	switch productName {
	case models.LifeProduct:
		footerText = "Wopta per te. Vita è un prodotto assicurativo di AXA France Vie S.A. – Rappresentanza Generale per l’Italia\ndistribuito da Wopta Assicurazioni S.r.l."
		logoPath = lib.GetAssetPathByEnvV2() + "logo_axa.png"
		x = 190
		y = 282.5
		height = 8
	case models.PmiProduct:
		footerText = ""
		logoPath = ""
	case models.PersonaProduct:
		footerText = "Wopta per te. Persona è un prodotto assicurativo di Global Assistance Compagnia di assicurazioni" +
			" e riassicurazioni S.p.A,\ndistribuito da Wopta Assicurazioni S.r.l"
		logoPath = lib.GetAssetPathByEnvV2() + "logo_global.png"
		x = 180
		y = 280
		height = 10
	case models.GapProduct:
		footerText = "Wopta per te. Auto Valore Protetto è un prodotto assicurativo di Sogessur SA – Rappresentanza " +
			" Generale per l’Italia con sede in Via Tiziano, 32 – 20145 Milano – \nIscritta alla CCIAA di Milano P.I." +
			" 07420570967 – REA MI 1957443; distribuito da Wopta Assicurazioni Srl"
		logoPath = lib.GetAssetPathByEnvV2() + "logo_sogessur.png"
		x = 183
		y = 285
		height = 3
	}

	pdf.SetFooterFunc(func() {
		pdf.SetXY(10, -15)
		setPinkRegularFont(pdf, smallTextSize)
		pdf.MultiCell(0, 3, footerText, "", "", false)
		opt.ImageType = "png"
		pdf.ImageOptions(logoPath, x, y, 0, height, false, opt, 0, "")
		pdf.SetY(-7)
		pageNumber(pdf)
	})
}

func axaHeader(pdf *fpdf.Fpdf, isProposal bool) {
	pdf.SetHeaderFunc(func() {
		var opt fpdf.ImageOptions
		pdf.SetXY(-30, 7)
		opt.ImageType = "png"
		pdf.ImageOptions(lib.GetAssetPathByEnvV2()+"logo_axa.png", 190, 7, 0, 8, false, opt, 0, "")
		pdf.Ln(15)

		if isProposal {
			insertWatermark(pdf, proposal)
		}
	})
}

func axaFooter(pdf *fpdf.Fpdf) {
	pdf.SetFooterFunc(func() {
		pdf.SetXY(10, -25)
		setBlackRegularFont(pdf, smallTextSize)
		pdf.MultiCell(0, 3, "AXA France Vie (compagnia assicurativa del gruppo AXA). Indirizzo sede "+
			"legale in Francia: 313 Terrasses de l'Arche, 92727 NANTERRE CEDEX. Numero Iscrizione Registro delle "+
			"Imprese di Nanterre: 310499959. Autorizzata in Francia (Stato di origine) all'esercizio delle "+
			"assicurazioni, vigilata in Francia dalla Autorité de Contrôle Prudentiel et de Résolution (ACPR). "+
			"Numero Matricola Registre des organismes d'assurance: 5020051. // Indirizzo Rappresentanza Generale "+
			"per l'Italia: Corso Como n. 17, 20154 Milano - CF, P.IVA e N.Iscr. Reg. Imprese 08875230016 - "+
			"REA MI-2525395 - Telefono: 02-87103548 - Fax: 02-23331247 - PEC: axafrancevie@legalmail.it - sito "+
			"internet: www.clp.partners.axa/it. Ammessa ad operare in Italia in regime di stabilimento. Iscritta "+
			"all'Albo delle imprese di assicurazione tenuto dall'IVASS, in appendice Elenco I, nr. I.00149.",
			"", "", false)
		pdf.SetY(-7)
		pageNumber(pdf)
	})
}

func sogessurHeader(pdf *fpdf.Fpdf) {
	pdf.SetHeaderFunc(func() {
		var opt fpdf.ImageOptions
		pdf.SetXY(-30, 7)
		opt.ImageType = "png"
		pdf.ImageOptions(lib.GetAssetPathByEnvV2()+"logo_sogessur.png", 160, 7, 0, 6, false,
			opt, 0, "")
		pdf.Ln(15)
	})
}

func sogessurFooter(pdf *fpdf.Fpdf) {
	pdf.SetFooterFunc(func() {
		pdf.SetXY(10, -12)
		setBlackRegularFont(pdf, smallTextSize)
		pdf.MultiCell(0, 3, "Sogecap SA – Rappresentanza Generale per l’Italia con sede in Via Tiziano, "+
			"32 – 20145 Milano – Iscritta alla CCIAA di Milano P.I. 07160010968 – REA MI 1939709", "",
			"", false)
		pdf.SetY(-7)
		pageNumber(pdf)
	})
}

func woptaHeader(pdf *fpdf.Fpdf, isProposal bool) {
	pdf.SetHeaderFunc(func() {
		var opt fpdf.ImageOptions
		opt.ImageType = "png"
		pdf.ImageOptions(lib.GetAssetPathByEnvV2()+"logo_wopta.png", 10, 6, 0, 10,
			false, opt, 0, "")
		pdf.Ln(10)

		if isProposal {
			insertWatermark(pdf, proposal)
		}
	})
}

func woptaFooter(pdf *fpdf.Fpdf) {
	pdf.SetFooterFunc(func() {
		pdf.SetY(-30)
		drawPinkHorizontalLine(pdf, 0.4)
		pdf.Ln(5)
		setPinkRegularFont(pdf, smallTextSize)
		pdf.Cell(pdf.GetStringWidth("Wopta Assicurazioni s.r.l"), 3, "Wopta Assicurazioni s.r.l")
		pdf.Cell(112, 3, "")
		pdf.Cell(pdf.GetStringWidth("www.wopta.it"), 3, "www.wopta.it")
		pdf.Ln(3)
		setBlackRegularFont(pdf, smallTextSize)
		pdf.CellFormat(pdf.GetStringWidth("Galleria del Corso, 1"), 3,
			"Galleria del Corso, 1", "", 0, "", false, 0, "")
		pdf.CellFormat(20, 3, "", "", 0, "", false, 0, "")
		pdf.CellFormat(pdf.GetStringWidth("Numero REA: MI 2638708"), 3,
			"Numero REA: MI 2638708", "", 0, "", false, 0, "")
		pdf.CellFormat(20, 3, "", "", 0, "", false, 0, "")
		pdf.CellFormat(pdf.GetStringWidth("CF | P.IVA | n. iscr. Registro Imprese:"), 3,
			"CF | P.IVA | n. iscr. Registro Imprese:", "", 0, "", false, 0, "")
		pdf.CellFormat(13, 3, "", "", 0, "", false, 0, "")
		pdf.CellFormat(30, 3, "info@wopta.it", "", 1, "", false, 0, "")
		pdf.CellFormat(pdf.GetStringWidth("Galleria del Corso, 1"), 3,
			"20122 - Milano (MI)", "", 0, "", false, 0, "")
		pdf.CellFormat(20, 3, "", "", 0, "", false, 0, "")
		pdf.CellFormat(pdf.GetStringWidth("Numero REA: MI 2638708"), 3,
			"Capitale Sociale: € 204.839,26 i.v.", "", 0, "", false, 0, "")
		pdf.CellFormat(20, 3, "", "", 0, "", false, 0, "")
		pdf.CellFormat(pdf.GetStringWidth("CF | P.IVA | n. iscr. Registro Imprese:"), 3,
			"12072020964", "", 0, "", false, 0, "")
		pdf.CellFormat(13, 3, "", "", 0, "", false, 0, "")
		pdf.CellFormat(30, 3, "(+39) 02 91240346", "", 1, "", false, 0, "")
		pdf.Ln(3)
		pdf.MultiCell(0, 3, "Wopta Assicurazioni s.r.l. è un intermediario assicurativo soggetto alla "+
			"vigilanza dell’IVASS ed iscritto alla Sezione A del Registro Unico degli Intermediari Assicurativi "+
			"con numero A000701923. Consulta gli estremi dell’iscrizione al sito "+
			"https://servizi.ivass.it/RuirPubblica/", "", "", false)
		pdf.SetY(-7)
		pageNumber(pdf)
	})
}

func globalHeader(pdf *fpdf.Fpdf) {
	pdf.SetHeaderFunc(func() {
		var opt fpdf.ImageOptions
		pdf.SetXY(-30, 7)
		opt.ImageType = "png"
		pdf.ImageOptions(lib.GetAssetPathByEnvV2()+"logo_global_02.png", 180, 7, 0, 15, false, opt, 0, "")
		pdf.Ln(17)
	})
}

func globalFooter(pdf *fpdf.Fpdf) {
	pdf.SetFooterFunc(func() {
		pdf.SetXY(10, -25)
		setBlackRegularFont(pdf, smallTextSize)
		pdf.MultiCell(0, 3, "Global Assistance Compagnia di assicurazioni e riassicurazioni S.p.A. a "+
			"Socio Unico - Capitale Sociale: Euro 5.000.000 i.v. Codice Fiscale, Partita IVA e Registro Imprese di "+
			"Milano n. 10086540159 R.E.A. n. 1345012 della C.C.I.A.A. di Milano. Sede e Direzione Generale Piazza "+
			"Diaz 6 – 20123 Milano – ItaliaE-mail: global.assistance@globalassistance.it PEC: "+
			"globalassistancespa@legalmail.it. Società soggetta all’attività di direzione e coordinamento di Ri-Fin "+
			"S.r.l., iscritta all’Albo dei gruppi assicurativi presso l’IVASS al n. 014. La Società è autorizzata "+
			"all’esercizio delle Assicurazioni e Riassicurazioni con D.M. del 2/8/93 n. 19619 (G.U. 7/8/93 n. 184) e"+
			" successive autorizzazioni ed è iscritta all’Albo Imprese presso l’Ivass al n. 1.00111. La Società è"+
			" soggetta alla vigilanza dell’IVASS; è possibile verificare la veridicità e la regolarità delle"+
			" autorizzazioni mediante l'accesso al sito www.ivass.it", "", "", false)
		pdf.SetY(-7)
		pageNumber(pdf)
	})
}

func globalPrivacySection(pdf *fpdf.Fpdf, survey models.Survey) {
	type row struct {
		text   string
		isBold bool
	}

	pages := [][][]row{
		{
			{
				{
					text: "Informativa resa all’interessato per il trattamento assicurativo di dati personali comuni, particolari " +
						"e dei dati relativi a condanne penali e reati",
					isBold: true,
				},
				{
					text: "Ai sensi dell’art. 13 del Regolamento Europeo n. 2016/679 " +
						"(General Data Protection Regulation – GDPR) ed in relazione ai dati personali che si " +
						"intendono trattare, La informiamo di quanto segue:",
					isBold: false,
				},
			},
			{
				{
					text:   "1. CATEGORIE DI DATI PERSONALI TRATTATI",
					isBold: true,
				},
				{
					text: "Il \"dato personale\" è \"qualsiasi informazione riguardante una persona fisica identificata o " +
						"identificabile (\"interessato\")”. Ai fini della presente Informativa il Titolare tratta i seguenti " +
						"dati personali: nome, cognome, indirizzo, e-mail, numero telefonico, codice fiscale o P. IVA " +
						"dell’interessato e dei soggetti da lui indicati per la copertura assicurativa. Oltre alle categorie " +
						"di dati indicati potranno anche essere trattati, previo consenso espresso dell’interessato, anche " +
						"per conto degli altri soggetti inclusi nella copertura assicurativa, dati particolari di cui all’art. " +
						"9 del GDPR (dati sanitari) e dati relativi a condanne penali e reati di cui all’art. 10 del GDPR.",
					isBold: false,
				},
			},
			{
				{
					text:   "2. FINALITÀ DEL TRATTAMENTO DEI DATI",
					isBold: true,
				},
				{
					text: "Il trattamento è diretto all’espletamento da parte del Titolare delle seguenti finalità:\n" +
						"- Procedere all’elaborazione di preventivi Assicurativi, sulla base delle informazioni ricevute;\n" +
						"- Procedere alla valutazione dei requisiti per l’assicurabilità dei soggetti interessati alla " +
						"stipula del contratto;\n" +
						"- Procedere alla conclusione, gestione ed esecuzione di contratti assicurativi e gestione e " +
						"liquidazione dei sinistri relativi ai medesimi contratti;\n" +
						"- Adempiere ad eventuali obblighi previsti dalla legge, da regolamenti, dalla normativa " +
						"comunitaria o da un ordine dell’Autorità;\n" +
						"- Esercitare i diritti del Titolare, ad esempio il diritto di difesa in giudizio;\n" +
						"- Perfezionare le offerte contrattuali sulla base dell’analisi della domanda di mercato e " +
						"delle caratteristiche degli assicurati e dei soggetti interessati " +
						"alla stipula di prodotti assicurativi, elaborando tali informazioni anche in combinazione con " +
						"informazioni provenienti da banche dati pubbliche.\n" +
						"Il trattamento avviene nell’ambito di attività assicurativa e riassicurativa, a cui il Titolare è " +
						"autorizzato ai sensi delle vigenti disposizioni di legge.",
					isBold: false,
				},
			},
			{
				{
					text:   "3. MODALITÀ DEL TRATTAMENTO DEI DATI",
					isBold: true,
				},
				{
					text: "Il trattamento dei Vostri dati personali, inclusi i dati particolari ai sensi degli artt. " +
						"9 e 10 GDPR, è realizzato per mezzo delle operazioni indicate all’art. 4 comma 1 n. 2) del GDPR " +
						"e precisamente: raccolta, registrazione, organizzazione, conservazione, consultazione, " +
						"elaborazione, modificazione, selezione, estrazione, raffronto, utilizzo, interconnessione, " +
						"blocco, comunicazione, cancellazione e distruzione dei dati. I Vostri dati personali sono " +
						"sottoposti a trattamento in formato sia cartaceo che elettronico.",
					isBold: false,
				},
			},
			{
				{
					text:   "4. NATURA DEL CONFERIMENTO DEI DATI E CONSEGUENZE DEL RIFIUTO",
					isBold: true,
				},
				{
					text: "Ferma l’autonomia personale dell’interessato, il conferimento dei dati può essere:\n" +
						"a) Obbligatorio in base ad una legge, regolamento o normativa comunitaria (ad esempio " +
						"Antiriciclaggio, Casellario Centrale Infortuni, Motorizzazione Civile);\n" +
						"b) Strettamente necessario alla redazione di preventivi assicurativi;\n" +
						"c) Strettamente necessario alla conclusione, gestione, ed esecuzione di contratti assicurativi " +
						"e gestione e liquidazione dei sinistri relativi ai medesimi contratti\n" +
						"L’eventuale rifiuto dell’interessato di conferire i dati personali in relazione alle finalità " +
						"di trattamento a), b), c), d) ed e) di cui al punto 2 della presente informativa comporta " +
						"l’impossibilità di procedere alla conclusione, gestione, ed esecuzione di contratti assicurativi " +
						"e gestione e\nliquidazione dei sinistri relativi ai medesimi contratti.\n" +
						"Il mancato consenso al trattamento dei dati sanitari comporterà l’impossibilità di includere la " +
						"copertura del rischio infortuni all’interno del contratto e il mancato consenso al trattamento " +
						"dei dati relativi a condanne penali o reati comporterà l’impossibilità di includere la " +
						"copertura della Tutela Legale all’interno del contratto.",
					isBold: false,
				},
			},
			{
				{
					text:   "5. CONSERVAZIONE",
					isBold: true,
				},
				{
					text: "I dati personali conferiti per le finalità sopra esposte saranno conservati per il periodo di " +
						"validità contrattuale assicurativa e successivamente per un periodo di 10 anni. Decorso tale " +
						"termine i dati personali saranno cancellati.",
					isBold: false,
				},
			},
			{
				{
					text:   "6. ACCESSO AI DATI",
					isBold: true,
				},
				{
					text: "I Vostri dati personali potranno essere resi accessibili per le finalità di cui sopra:\n" +
						"a) A dipendenti e collaboratori del Titolare, nella loro qualità di soggetti designati;\n" +
						"b) A intermediari assicurativi per finalità di conclusione gestione, ed esecuzione di " +
						"contratti assicurativi e gestione dei sinistri relativi ai medesimi contratti;\n" +
						"c) A soggetti esterni che forniscono servizi in outsourcing al Titolare;\n" +
						"d) A riassicuratori con i quali il Titolare sottoscriva specifici trattati per la copertura " +
						"dei rischi riferiti al contratto assicurativo, tra cui espressamente la società Munich Re.",
					isBold: false,
				},
			},
			{
				{
					text:   "7. COMUNICAZIONE DEI DATI",
					isBold: true,
				},
				{
					text: "Il Titolare potrà comunicare i Vostri dati, per le finalità di cui al punto 2 precedente e " +
						"per essere sottoposti a trattamenti aventi le medesime finalità o obbligatori per legge, a " +
						"terzi soggetti operanti nel settore assicurativo, società di servizi informatici o società a " +
						"cui il Titolare ha affidato attività in outsourcing o altri soggetti nei confronti dei quali la " +
						"comunicazione è obbligatoria.",
					isBold: false,
				},
			},
		},
		{
			{
				{
					text:   "8. DIFFUSIONE",
					isBold: true,
				},
				{
					text:   "I dati personali di cui alla presente informativa non sono soggetti a diffusione.",
					isBold: false,
				},
			},
			{
				{
					text:   "9. TRASFERIMENTO DATI ALL’ESTERO",
					isBold: true,
				},
				{
					text: "La gestione e la conservazione dei dati personali avverranno su server ubicati all’interno " +
						"del territorio italiano o comunque dell’Unione Europea. I dati non saranno oggetto di " +
						"trasferimento all’esterno dell’Unione Europea.",
					isBold: false,
				},
			},
			{
				{
					text:   "10. DIRITTI DELL’INTERESSATO",
					isBold: true,
				},
				{
					text: "In qualità di interessati, avete i diritti riconosciuti dall’art. 15 del GDPR, in " +
						"particolare di:\n" +
						"a) Ottenere la conferma dell’esistenza o meno dei dati personali che vi riguardano;\n" +
						"b) Ottenere l’indicazione: a) dell’origine dei dati personali; b) delle finalità e modalità " +
						"del trattamento; c) della logica applicata in caso di trattamento effettuato con l’ausilio " +
						"di strumenti elettronici; d) degli estremi identificativi del Titolare, degli eventuali " +
						"responsabili e dell’eventuale\nrappresentante designati ai sensi dell’art. 3 comma 1 del " +
						"GDPR; e) dei soggetti e delle categorie di soggetti ai quali i dati personali possono essere " +
						"comunicati o che possono venirne a conoscenza in qualità di responsabili o incaricati; " +
						"c) Ottenere: a) l’aggiornamento, la rettifica ovvero, quanto avete interesse, l’integrazione " +
						"dei dati; b) la cancellazione, la trasformazione in forma anonima o il blocco dei dati " +
						"trattati in violazione di legge, compresi quelli di cui non è necessaria la conservazione " +
						"in relazione agli scopi per i quali i dati sono stati raccolti o successivamente trattati; " +
						"c) l’attestazione che le operazioni di cui alle lettere a) e b) sono state portate a " +
						"conoscenza, anche per quanto riguarda il loro contenuto, di coloro ai quali i dati son o " +
						"stati comunicati o diffusi, eccettuato il caso in cui tale adempimento si " +
						"riveli impossibile o comporti un impiego di mezzi manifestamente sproporzionato rispetto al " +
						"diritto tutelato; " +
						"d) Opporsi, in tutto o in parte: a) per motivi legittimi al trattamento dei dati personali " +
						"che vi riguardano, ancorché pertinenti allo scopo della raccolta; b) al trattamento di " +
						"dati personali che vi riguardano a fini di invio di materiale pubblicitario o di vendita " +
						"diretta o per il compimento di ricerche di mercato o di comunicazione commerciale. Ove " +
						"applicabili, avete altresì i diritti di cui agli articoli 16 – 21 del GDPR (Diritto di " +
						"rettifica, diritto all’oblio, diritto di limitazione di trattamento, diritto alla " +
						"portabilità dei dati contrattuali e grezzi di navigazione, diritto di opposizione), " +
						"nonché il diritto di reclamo all’Autorità Garante.",
					isBold: false,
				},
			},
			{
				{
					text:   "11. TITOLARE DEL TRATTAMENTO",
					isBold: true,
				},
				{
					text: "Il titolare dei trattamenti per le finalità indicate al punto 2 della presente " +
						"informativa è:\n" +
						"Global Assistance Compagnia di Assicurazioni e Riassicurazioni S.p.A. (Global)\n" +
						"Piazza Armando Diaz n. 6\n" +
						"20123 – Milano\n" +
						"E-mail: global.assistance@globalassistance.it\n" +
						"PEC: globalassistancespa@legalmail.it\n" +
						"Fax: 02/43335020\n" +
						"Limitatamente alle finalità di cui alle lettere b) e f) del punto 2 della presente " +
						"informativa è titolare del trattamento anche:\n" +
						"Münchener Rückversicherungs-Gesellschaft (Munich Re)\n" +
						"Rappresentanza Generale per l'Italia\n" +
						"Via Pola, 9\n" +
						"20124 - Milano\n" +
						"E-mail: mritalia@munichre.com\n" +
						"PEC: munchenerruck@legalmail.it",
					isBold: false,
				},
			},
			{
				{
					text:   "12. MODALITA’ DI ESERCIZIO DEI DIRITTI",
					isBold: true,
				},
				{
					text: "Potrete in qualsiasi momento esercitare i Vostri diritti inviando una e-mail, una PEC, " +
						"un fax o una raccomandata A.R. all’indirizzo del Titolare.\n" +
						"È possibile contattare direttamente il Responsabile della Protezione dei Dati– RPD o Data " +
						"Protection Officer – DPO di Global al seguente indirizzo e-mail: info@lext.it.\n" +
						"È possibile contattare direttamente il Responsabile della Protezione dei Dati– RPD o Data " +
						"Protection Officer – DPO di Munich Re al seguente indirizzo e-mail: datenschutz@munichre.com.",
					isBold: false,
				},
			},
			{
				{
					text: "Informativa resa all’interessato per il trattamento assicurativo di dati personali comuni, " +
						"particolari e dei dati relativi a condanne penali e reati.",
					isBold: true,
				},
			},
			{
				{
					text:   "DICHIARAZIONI E CONSENSI",
					isBold: true,
				},
				{
					text: "Io Sottoscritto, dichiaro di avere perso visione dell’Informativa Privacy ai sensi " +
						"dell’art. 13 del GDPR (informativa resa all’interno del set documentale contenente anche la " +
						"Documentazione Informativa Precontrattuale, il Glossario e le Condizioni di Assicurazione) e " +
						"di averne compreso i contenuti.",
					isBold: false,
				},
			},
		},
	}

	for pageIndex, page := range pages {
		for _, paragraph := range page {
			for _, r := range paragraph {
				setBlackRegularFont(pdf, standardTextSize)
				if r.isBold {
					setBlackBoldFont(pdf, standardTextSize)
				}
				pdf.MultiCell(0, 3.75, r.text, "", fpdf.AlignLeft, false)
			}
			pdf.Ln(2.5)
		}
		if pageIndex < len(pages)-1 {
			pdf.AddPage()
		}
	}

	drawSignatureForm(pdf)

	pdf.AddPage()

	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3.75, "Qui di seguito esprimo il mio consenso al trattamento dei dati personali "+
		"particolari per le finalità sopra indicate, in conformità con quanto previsto all’interno dell’informativa: ",
		"", fpdf.AlignLeft, false)

	pdf.Ln(3)

	privacyConsent := [][]string{
		{"", "X"},
		{"", "X"},
	}

	questions := survey.Questions[len(survey.Questions)-2:]
	for questionIndex, question := range questions {
		if question.Answer != nil && *question.Answer {
			privacyConsent[questionIndex][0] = "X"
			privacyConsent[questionIndex][1] = ""
		}
	}

	table := [][]tableCell{
		{
			{
				text:      "Consenso al trattamento dei miei dati particolari (sanitari) di cui all’art. 9 del GDPR:",
				height:    3,
				width:     160,
				textBold:  true,
				fill:      false,
				fillColor: rgbColor{},
				align:     fpdf.AlignLeft,
				border:    "",
			},
			{
				text:      privacyConsent[0][0],
				height:    5,
				width:     5,
				textBold:  true,
				fill:      false,
				fillColor: rgbColor{},
				align:     fpdf.AlignCenter,
				border:    "1",
			},
			{
				text:      "SI",
				height:    5,
				width:     10,
				textBold:  false,
				fill:      false,
				fillColor: rgbColor{},
				align:     fpdf.AlignCenter,
				border:    "",
			},
			{
				text:      privacyConsent[0][1],
				height:    5,
				width:     5,
				textBold:  true,
				fill:      false,
				fillColor: rgbColor{},
				align:     fpdf.AlignCenter,
				border:    "1",
			},
			{
				text:      "NO",
				height:    5,
				width:     10,
				textBold:  false,
				fill:      false,
				fillColor: rgbColor{},
				align:     fpdf.AlignCenter,
				border:    "",
			},
		},
	}

	tableDrawer(pdf, table)

	pdf.Ln(3)

	table = [][]tableCell{
		{
			{
				text: "Consenso al trattamento dei miei dati al fine di perfezionamento dell’offerta assicurativa " +
					"e riassicurativa di cui alle lettere b) ed f) della presente informativa:",
				height:    3,
				width:     160,
				textBold:  true,
				fill:      false,
				fillColor: rgbColor{},
				align:     fpdf.AlignLeft,
				border:    "",
			},
			{
				text:      privacyConsent[1][0],
				height:    5,
				width:     5,
				textBold:  true,
				fill:      false,
				fillColor: rgbColor{},
				align:     fpdf.AlignCenter,
				border:    "1",
			},
			{
				text:      "SI",
				height:    5,
				width:     10,
				textBold:  false,
				fill:      false,
				fillColor: rgbColor{},
				align:     fpdf.AlignCenter,
				border:    "",
			},
			{
				text:      privacyConsent[1][1],
				height:    5,
				width:     5,
				textBold:  true,
				fill:      false,
				fillColor: rgbColor{},
				align:     fpdf.AlignCenter,
				border:    "1",
			},
			{
				text:      "NO",
				height:    5,
				width:     10,
				textBold:  false,
				fill:      false,
				fillColor: rgbColor{},
				align:     fpdf.AlignCenter,
				border:    "",
			},
		},
	}

	tableDrawer(pdf, table)

	pdf.Ln(7)

	drawSignatureForm(pdf)

}

func paymentMethodSection(pdf *fpdf.Fpdf) {
	getParagraphTitle(pdf, "Come puoi pagare il premio")
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, "I mezzi di pagamento consentiti, nei confronti di Wopta, sono esclusivamente "+
		"bonifico e strumenti di pagamento elettronico, quali ad esempio, carte di credito e/o carte di debito, "+
		"incluse le carte prepagate. Oppure può essere pagato direttamente alla Compagnia alla "+
		"stipula del contratto, via bonifico o carta di credito.", "", "", false)
	pdf.Ln(3)
}

func emitResumeSection(pdf *fpdf.Fpdf, policy *models.Policy) {
	var offerPrice string
	emitDate := time.Now().UTC().Format(dateLayout)
	if policy.PaymentSplit == "monthly" {
		offerPrice = humanize.FormatFloat("#.###,##", policy.PriceGrossMonthly)
	} else {
		offerPrice = humanize.FormatFloat("#.###,##", policy.PriceGross)
	}
	text := "Polizza emessa a Milano il " + emitDate + " per un importo di € " + offerPrice + " quale " +
		"prima rata alla firma, il cui pagamento a saldo è da effettuarsi con i metodi di pagamento sopra indicati."
	switch policy.Name {
	case models.PersonaProduct:
		text += "\nCostituisce quietanza di pagamento la mail di conferma che Wopta invierà al Contraente."

	}

	getParagraphTitle(pdf, "Emissione polizza e pagamento della prima rata")
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, text, "", "", false)
	pdf.Ln(3)
}

func companiesDescriptionSection(pdf *fpdf.Fpdf, companyName string) {
	getParagraphTitle(pdf, "Chi siamo")
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, "Wopta Assicurazioni S.r.l. - intermediario assicurativo, soggetto al controllo "+
		"dell’IVASS ed iscritto dal 14.02.2022 al Registro Unico degli Intermediari, in Sezione A nr. A000701923, "+
		"avente sede legale in Galleria del Corso, 1 – 20122 Milano (MI). Capitale sociale Euro 120.000 - "+
		"Codice Fiscale, Reg. Imprese e Partita IVA: 12072020964 - Iscritta al Registro delle imprese di Milano – "+
		"REA MI 2638708", "", "", false)
	pdf.Ln(5)

	switch companyName {
	case "axa":
		pdf.MultiCell(0, 3, "AXA France Vie (compagnia assicurativa del gruppo AXA). Indirizzo sede legale in "+
			"Francia: 313 Terrasses de l'Arche, 92727 NANTERRE CEDEX. Numero Iscrizione Registro delle Imprese di "+
			"Nanterre: 310499959. Autorizzata in Francia (Stato di origine) all’esercizio delle assicurazioni, vigilata "+
			"in Francia dalla Autorité de Contrôle Prudentiel et de Résolution (ACPR). Numero Matricola Registre des "+
			"organismes d’assurance: 5020051. // Indirizzo Rappresentanza Generale per l’Italia: Corso Como n. 17, 20154 "+
			"Milano - CF, P.IVA e N.Iscr. Reg. Imprese 08875230016 - REA MI-2525395 - Telefono: 02-87103548 - "+
			"Fax: 02-23331247 - PEC: axafrancevie@legalmail.it – sito internet: www.clp.partners.axa/it. Ammessa ad "+
			"operare in Italia in regime di stabilimento. Iscritta all’Albo delle imprese di assicurazione tenuto "+
			"dall’IVASS, in appendice Elenco I, nr. I.00149.", "", "", false)
	case "global":
		pdf.MultiCell(0, 3, "Global Assistance Compagnia di assicurazioni e riassicurazioni S.p.A. a "+
			"Socio Unico - Capitale Sociale: Euro 5.000.000 i.v. Codice Fiscale, Partita IVA e Registro Imprese di "+
			"Milano n. 10086540159 R.E.A. n. 1345012 della C.C.I.A.A. di Milano. Sede e Direzione Generale Piazza "+
			"Diaz 6 – 20123 Milano – Italia E-mail: global.assistance@globalassistance.it PEC: "+
			"globalassistancespa@legalmail.it. Società soggetta all’attività di direzione e coordinamento di Ri-Fin "+
			"S.r.l., iscritta all’Albo dei gruppi assicurativi presso l’Ivass al n. 014. La Società è autorizzata "+
			"all’esercizio delle Assicurazioni e Riassicurazioni con D.M. del 2/8/93 n. 19619 (G.U. 7/8/93 n. 184) e "+
			"successive autorizzazioni ed è iscritta all’Albo Imprese presso l’IVASS al n. 1.00111. La Società è "+
			"soggetta alla vigilanza dell’IVASS; è possibile verificare la veridicità e la regolarità delle "+
			"autorizzazioni mediante l'accesso al sito www.ivass.it", "", "", false)
	case "sogessur":
		pdf.MultiCell(0, 3, "Sogessur SA – Sogessur - Société Anonyme – Capitale Sociale € 33 825 000 "+
			"– Sede legale: Tour D2, 17bis Place des Reflets – 92919 Paris La Défense Cedex - 379 846 637 R.C.S. "+
			"Nanterre - Francia - Sede secondaria: Via Tiziano 32, 20145 Milano - Italia - Registro delle Imprese di "+
			"Milano, Lodi,Monza-Brianza Codice Fiscale e P.IVA  07420570967  Iscritta nell’elenco I dell’Albo delle "+
			"Imprese di Assicurazione tenuto dall’IVASS al n. I00094", "", "", false)
	}
	pdf.Ln(3)
}

func personalDataHandlingSection(pdf *fpdf.Fpdf, policy *models.Policy, isProposal bool) {
	consentText := ""
	notConsentText := "X"

	if policy.Contractor.Consens != nil {
		consent, err := policy.ExtractConsens(2)
		lib.CheckError(err)

		if consent.Answer {
			consentText = "X"
			notConsentText = ""
		}
	}

	getParagraphTitle(pdf, "Consenso per finalità commerciali.")
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, "Il sottoscritto, letta e compresa l’informativa sul trattamento dei dati personali",
		"", "", false)
	pdf.Ln(1)
	setBlackDrawColor(pdf)
	pdf.Cell(5, 3, "")
	pdf.CellFormat(3, 3, consentText, "1", 0, "CM", false, 0, "")
	pdf.CellFormat(20, 3, "ACCONSENTE", "", 0, "", false, 0, "")
	pdf.Cell(20, 3, "")
	pdf.CellFormat(3, 3, notConsentText, "1", 0, "CM", false, 0, "")
	pdf.CellFormat(20, 3, "NON ACCONSENTE", "", 1, "", false, 0, "")
	pdf.Ln(1)
	pdf.MultiCell(0, 3, "al trattamento dei propri dati personali da parte di Wopta Assicurazioni per "+
		"l’invio di comunicazioni e proposte commerciali e di marketing, incluso l’invio di newsletter e ricerche di "+
		"mercato, attraverso strumenti automatizzati (sms, mms, e-mail, ecc.) e non (posta cartacea e telefono "+
		"con operatore).", "", "", false)
	pdf.Ln(3)
	pdf.Cell(0, 3, time.Now().UTC().Format(dateLayout))
	if !isProposal {
		drawSignatureForm(pdf)
	}
}

func woptaPrivacySection(pdf *fpdf.Fpdf) {
	setBlackBoldFont(pdf, titleTextSize)
	pdf.MultiCell(0, 3, "COME RISPETTIAMO LA TUA PRIVACY", "", "CM", false)
	pdf.Ln(3)
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, "Informativa sul trattamento dei dati personali", "", "", false)
	pdf.Ln(1)
	setBlackBoldFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, "Ai sensi del REGOLAMENTO (UE) 2016/679 "+
		"(relativo alla protezione delle persone fisiche con riguardo al trattamento dei dati personali, nonché alla "+
		"libera circolazione di tali dati) si informa l’ “Interessato” (contraente / aderente alla polizza collettiva o "+
		"convenzione / assicurato / beneficiario / loro aventi causa) di quanto segue.", "", "", false)
	pdf.Ln(3)
	setBlackBoldFont(pdf, titleTextSize)
	pdf.MultiCell(0, 3, "1. TITOLARE DEL TRATTAMENTO", "", "", false)
	pdf.Ln(1)
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, "Titolare del trattamento è Wopta Assicurazioni, con sede legale in Milano, "+
		"Galleria del Corso, 1 (di seguito “Titolare”), raggiungibile all’indirizzo e-mail: "+
		"privacy@wopta.it", "", "", false)
	pdf.Ln(3)
	setBlackBoldFont(pdf, titleTextSize)
	pdf.MultiCell(0, 3, "2. I DATI PERSONALI OGGETTO DI TRATTAMENTO, FINALITÀ E BASE "+
		"GIURIDICA", "", "", false)
	pdf.Ln(1)
	setBlackBoldFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, "a) Finalità Contrattuali, normative, amministrative e giudiziali", "", "", false)
	pdf.Ln(1)
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, "Fermo restando quanto previsto dalla Privacy & Cookie Policy del Sito, ove "+
		"applicabile, i dati così conferiti potranno essere trattati, anche con strumenti elettronici, da parte del "+
		"Titolare per eseguire le prestazioni contrattuali, in qualità di intermediario, richieste dall’interessato, "+
		"o per adempiere ad obblighi normativi, contabili e fiscali, ovvero ancora per finalità di difesa in "+
		"giudizio, per il tempo strettamente necessario a tali attività.", "", "", false)
	pdf.Ln(1)
	pdf.MultiCell(0, 3, "La base giuridica del trattamento di dati personali per le finalità di cui sopra "+
		"è l’art. 6.1 lett. b), c), f) del Regolamento in quanto i trattamenti sono necessari all'erogazione dei "+
		"servizi o per il riscontro di richieste dell’interessato, in conformità a quanto previsto dall’incarico "+
		"conferito all’intermediario, nonché ove il trattamento risulti necessario per l’adempimento di un preciso "+
		"obbligo di legge posto in capo al Titolare, o al fine di accertare, esercitare o difendere un diritto in "+
		"sede giudiziaria. Il conferimento dei dati personali per queste finalità è facoltativo, ma l'eventuale "+
		"mancato conferimento comporterebbe l'impossibilità per l’intermediario di eseguire le proprie obbligazioni "+
		"contrattuali.", "", "", false)
	pdf.Ln(1)
	setBlackBoldFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, "b) Finalità commerciali", "", "", false)
	pdf.Ln(1)
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, "Inoltre, i Suoi dati personali potranno essere trattati al fine di inviarLe "+
		"comunicazioni e proposte commerciali, incluso l’invio di newsletter e ricerche di mercato, attraverso "+
		"strumenti automatizzati (sms, mms, email, messaggistica istantanea e chat) e non (posta cartacea, telefono); "+
		"si precisa che il Titolare raccoglie un unico consenso per le finalità di marketing qui descritte, ai sensi "+
		"del Provvedimento Generale del Garante per la Protezione dei Dati Personali \"Linee guida in materia di "+
		"attività promozionale e contrasto allo spam” del 4 luglio 2013; qualora, in ogni caso, Lei desiderasse "+
		"opporsi al trattamento dei Suoi dati per le finalità di marketing eseguito con i mezzi qui indicati, potrà "+
		"in qualunque momento farlo contattando il Titolare ai recapiti indicati nella sezione \"Contatti\" di "+
		"questa informativa, senza pregiudicare la liceità del trattamento effettuato prima dell’opposizione.",
		"", "", false)
	pdf.Ln(1)
	pdf.MultiCell(0, 3, "I trattamenti eseguiti per la finalità di marketing, di cui al paragrafo che "+
		"precede, si basa sul rilascio del Suo consenso ai sensi dell’art. 6, par. 1, lett. a) ([…] l'interessato ha "+
		"espresso il consenso al trattamento dei propri dati personali per una o più specifiche finalità) del "+
		"Regolamento. Tale consenso è revocabile in qualsiasi momento senza pregiudizio alcuno della liceità del "+
		"trattamento effettuato anteriormente alla revoca in conformità a quanto previsto dall’art. 7 del "+
		"Regolamento. Il conferimento dei Suoi dati personali per queste finalità è quindi del tutto facoltativo e "+
		"non pregiudica la fruizione dei servizi. Qualora desiderasse opporsi al trattamento dei Suoi dati per le "+
		"finalità di marketing, potrà in qualunque momento farlo contattando il Titolare ai recapiti indicati nella "+
		"sezione \"Contatti\" di questa informativa.", "", "", false)
	pdf.Ln(3)
	setBlackBoldFont(pdf, titleTextSize)
	pdf.MultiCell(0, 3, "3. DESTINATARI DEI DATI PERSONALI", "", "", false)
	pdf.Ln(1)
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, "I Suoi dati personali potranno essere condivisi, per le finalità di cui alla "+
		"sezione 2 della presente Policy, con:", "", "", false)
	pdf.Ln(1)
	pdf.MultiCell(0, 3, "- Soggetti che agiscono tipicamente in qualità di Responsabili del trattamento "+
		"ex art. 28 del Regolamento per conto del Titolare, incaricati dell'erogazione dei Servizi (a titolo "+
		"esemplificativo: servizi tecnologici, servizi di assistenza e consulenza in materia contabile, amministrativa, "+
		"legale, tributaria e finanziaria, manutenzione tecnica). Il Titolare conserva una lista aggiornata dei "+
		"responsabili del trattamento nominati e ne garantisce la presa visione all’interessato presso la sede sopra "+
		"indicata o previa richiesta indirizzata ai recapiti sopra indicati;", "", "", false)
	pdf.Ln(1)
	pdf.MultiCell(0, 3, "- Persone autorizzate dal Titolare al trattamento dei dati personali ai sensi "+
		"degli artt. 29 e 2-quaterdecies del D.lgs. n. 196/2003 (“Codice “Privacy”) (ad es. il personale dipendente "+
		"addetto alla manutenzione del Sito, alla gestione del CRM, alla gestione dei sistemi informativi ecc.); ",
		"", "", false)
	pdf.Ln(1)
	pdf.MultiCell(0, 3, "- Soggetti terzi, autonomi titolari del trattamento, a cui i dati potrebbero "+
		"essere trasmessi al fine di dare seguito a specifici servizi da Lei richiesti e/o  per dare esecuzione alle "+
		"attività di cui alla presente informativa, e con i quali il Titolare abbia stipulato accordi commerciali; "+
		"soggetti, quali le imprese di assicurazione, che assumono il rischio di sottoscrizione della polizza, ai "+
		"quali sia obbligatorio comunicare i tuoi Dati personali in forza di obblighi contrattuali e di disposizioni "+
		"di legge e regolamentari sulla distribuzione di prodotti assicurativi;", "", "", false)
	pdf.Ln(1)
	pdf.MultiCell(0, 3, "- Soggetti, enti od autorità a cui sia obbligatorio comunicare i Suoi dati "+
		"personali in forza di disposizioni di legge o di ordini delle autorità.", "", "", false)
	pdf.Ln(1)
	pdf.MultiCell(0, 3, "Tali soggetti sono, di seguito, collettivamente definiti come “Destinatari”. "+
		"L'elenco completo dei responsabili del trattamento è disponibile inviando una richiesta scritta al Titolare "+
		"ai recapiti indicati nella sezione \"Contatti\" di questa informativa.", "", "", false)
	pdf.Ln(3)
	setBlackBoldFont(pdf, titleTextSize)
	pdf.MultiCell(0, 3, "4. TRASFERIMENTI DEI DATI PERSONALI", "", "", false)
	pdf.Ln(1)
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, "Alcuni dei Suoi dati personali sono condivisi con Destinatari che si potrebbero "+
		"trovare al di fuori dello Spazio Economico Europeo. Il Titolare assicura che il trattamento Suoi dati "+
		"personali da parte di questi Destinatari avviene nel rispetto degli artt. 44 - 49 del Regolamento. Invero, "+
		"per quanto concerne il trasferimento dei dati personali verso Paesi terzi, il Titolare rende noto che il "+
		"trattamento avverrà secondo una delle modalità consentite dalla legge vigente, quali, ad esempio, il "+
		"consenso dell’interessato, l’adozione di Clausole Standard approvate dalla Commissione Europea, la selezione"+
		" di soggetti aderenti a programmi internazionali per la libera circolazione dei dati o operanti in Paesi "+
		"considerati sicuri dalla Commissione Europea sulla base di una decisione di adeguatezza.",
		"", "", false)
	pdf.Ln(1)
	pdf.MultiCell(0, 3, "Maggiori informazioni sono disponibili inviando una richiesta scritta al "+
		"Titolare ai recapiti indicati nella sezione \"Contatti\" di questa informativa.",
		"", "", false)
	setBlackBoldFont(pdf, titleTextSize)
	pdf.AddPage()
	pdf.MultiCell(0, 3, "5. CONSERVAZIONE DEI DATI PERSONALI", "", "", false)
	pdf.Ln(1)
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, "I Suoi dati personali saranno inseriti e conservati, in conformità ai principi "+
		"di minimizzazione e limitazione della conservazione di cui all’art. 5.1.c) ed e) del Regolamento, nei "+
		"sistemi informativi del Titolare, i cui server sono situati all’interno dello Spazio Economico Europeo.",
		"", "", false)
	pdf.Ln(1)
	pdf.MultiCell(0, 3, "I dati personali trattati per le finalità di cui alle lettere a) e b) "+
		"saranno conservati per il tempo strettamente necessario a raggiungere quelle stesse finalità ovverossia per "+
		"il tempo necessario all’esecuzione del contratto, in conformità ai tempi di conservazione obbligatori per "+
		"legge (vedi anche, in particolare, art. 2946 c.c. e ss.).",
		"", "", false)
	pdf.Ln(1)
	pdf.MultiCell(0, 3, "Per le finalità di cui alla lettera c), i suoi dati personali saranno invece "+
		"trattati fino alla revoca del suo consenso. Alla revoca del consenso, i dati trattati per la finalità di cui"+
		" sopra verranno cancellati o resi anonimi in modo permanente.",
		"", "", false)
	pdf.Ln(1)
	pdf.MultiCell(0, 3, "In generale, il Titolare si riserva in ogni caso di conservare i Suoi dati per "+
		"il tempo necessario ad adempiere ogni eventuale obbligo normativo cui lo stesso è soggetto o per soddisfare "+
		"eventuali esigenze difensive. Resta infatti salva la possibilità per il Titolare di conservare i Suoi dati "+
		"personali per il periodo di tempo previsto e ammesso dalla legge Italiana a tutela dei propri interessi "+
		"(Art. 2947 c.c.).", "", "", false)
	pdf.Ln(1)
	pdf.MultiCell(0, 3, "Maggiori informazioni in merito al periodo di conservazione dei dati e ai "+
		"criteri utilizzati per determinare tale periodo possono essere richieste inviando una richiesta scritta al "+
		"Titolare ai recapiti indicati nella sezione \"Contatti\" di questa informativa. ",
		"", "", false)
	pdf.Ln(3)
	setBlackBoldFont(pdf, titleTextSize)
	pdf.MultiCell(0, 3, "6. DIRITTI DELL’INTERESSATO", "", "", false)
	pdf.Ln(1)
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, "Lei ha il diritto di accedere in qualunque momento ai Dati Personali che La "+
		"riguardano, ai sensi degli artt. 15-22 del Regolamento. In particolare, potrà chiedere la rettifica "+
		"(ex art. 16), la cancellazione (ex art. 17), la limitazione (ex art. 18) e la portabilità dei dati "+
		"(ex art. 20), di non essere sottoposto a una decisione basata unicamente sul trattamento automatizzato, "+
		"compresa la profilazione, che produca effetti giuridici che La riguardano o che incida in modo analogo "+
		"significativamente sulla sua persona (ex art. 22), nonché la revoca del consenso eventualmente prestato "+
		"(ex art. 7, par. 3).", "", "", false)
	pdf.Ln(1)
	pdf.MultiCell(0, 3, "Lei può formulare, inoltre, una richiesta di opposizione al trattamento dei "+
		"Suoi Dati Personali ex art. 21 del Regolamento nella quale dare evidenza delle ragioni che giustifichino "+
		"l’opposizione: il titolare si riserva di valutare la Sua istanza, che non verrebbe accettata in caso di "+
		"esistenza di motivi legittimi cogenti per procedere al trattamento che prevalgano sui Suoi interessi, "+
		"diritti e libertà. Lei ha altresì il diritto di opporsi in ogni momento e senza alcuna giustificazione "+
		"all’invio di marketing diretto attraverso strumenti automatizzati (es. sms, mms, e-mail, notifiche push, "+
		"fax, sistemi di chiamata automatizzati senza operatore) e non (posta cartacea, telefono con operatore). "+
		"Inoltre, con riguardo al marketing diretto, resta salva la possibilità di esercitare tale diritto anche "+
		"in parte, ossia, in tal caso, opponendosi, ad esempio, al solo invio di comunicazioni promozionali "+
		"effettuato tramite strumenti automatizzati.", "", "", false)
	pdf.Ln(1)
	pdf.MultiCell(0, 3, "Le richieste vanno rivolte per iscritto al"+
		" Titolare ai recapiti indicati nella sezione \"Contatti\" di questa informativa.", "", "",
		false)
	pdf.Ln(1)
	pdf.MultiCell(0, 3, "Qualora Lei ritenga che il trattamento dei Suoi Dati personali effettuato dal "+
		"Titolare avvenga in violazione di quanto previsto dal GDPR, ha il diritto di proporre reclamo al Garante "+
		"Privacy, come previsto dall'art. 77 del GDPR stesso, o di adire le opportune sedi giudiziarie "+
		"(art. 79 del GDPR).", "", "", false)
	pdf.Ln(3)
	setBlackBoldFont(pdf, titleTextSize)
	pdf.MultiCell(0, 3, "7. CONTATTI", "", "", false)
	pdf.Ln(1)
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, "Per esercitare i diritti di cui sopra o per qualunque altra richiesta può "+
		"scrivere al Titolare del trattamento all’indirizzo: privacy@wopta.it.", "", "", false)
	pdf.Ln(3)
}

func companySignature(pdf *fpdf.Fpdf, companyName string) {
	var opt fpdf.ImageOptions
	opt.ImageType = "png"

	switch companyName {
	case "global":
		setBlackBoldFont(pdf, standardTextSize)
		pdf.CellFormat(70, 3, "Global Assistance", "", 0,
			fpdf.AlignCenter, false, 0, "")
		pdf.ImageOptions(lib.GetAssetPathByEnvV2()+"signature_global.png", 25, pdf.GetY()+3, 40, 12,
			false, opt, 0, "")
	case "axa":
		setBlackBoldFont(pdf, standardTextSize)
		pdf.MultiCell(70, 3, "AXA France Vie\n(Rappresentanza Generale per l'Italia)", "",
			fpdf.AlignCenter, false)
		pdf.SetY(pdf.GetY() - 6)
		pdf.ImageOptions(lib.GetAssetPathByEnvV2()+"signature_axa.png", 35, pdf.GetY()+9, 30, 8,
			false, opt, 0, "")
	case "sogessur":
		setBlackBoldFont(pdf, standardTextSize)
		pdf.MultiCell(70, 3, "Sogessur SA\n(Rappresentanza Generale per l'Italia)", "",
			fpdf.AlignCenter, false)
		pdf.SetY(pdf.GetY() - 6)
		pdf.ImageOptions(lib.GetAssetPathByEnvV2()+"signature_sogessur.png", 40, pdf.GetY()+9, 10, 10,
			false, opt, 0, "")
	}
}

func contractWithdrawlSection(pdf *fpdf.Fpdf, isProposal bool) {
	getParagraphTitle(pdf, "Informativa sul diritto di recesso")
	setBlackBoldFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, "Diritto di recesso entro i primi 30 giorni dalla stipula ("+
		"diritto di ripensamento)", "", "", false)
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, "Il Contraente può recedere dal contratto entro il termine di 30 giorni dalla "+
		"decorrenza dell’assicurazione (diritto di ripensamento). In tal caso, l’assicurazione si intende come mai "+
		"entrata in vigore e la Compagnia, per il tramite dell’intermediario, provvederà a rimborsare al Contraente "+
		"l’importo di Premio già versato (al netto delle imposte).", "", "", false)
	setBlackBoldFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, "Diritto di recesso annuale (disdetta alla annualità)", "", "",
		false)
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, "Il Contraente può recedere dal contratto annualmente, entro il termine di 30 "+
		"giorni dalla scadenza annuale della polizza (disdetta alla annualità). In tal caso, l’assicurazione cessa alle "+
		"ore 24:00 dell’ultimo giorno della annualità in corso. È possibile disdettare singolarmente una o più delle "+
		"coperture attivate in fase di sottoscrizione.", "", "", false)
	setBlackBoldFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, "Modalità per l’esercizio del diritto di recesso", "", "", false)
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, "Il Contraente è tenuto ad esercitare il diritto di recesso mediante invio di una "+
		"lettera raccomandata a.r. al seguente indirizzo: Wopta Assicurazioni srl – Gestione Portafoglio – Galleria del "+
		"Corso, 1 – 201212 Milano (MI) oppure via posta elettronica certificata (PEC) all’indirizzo "+
		"email: woptaassicurazioni@legalmail.it", "", "", false)
	if !isProposal {
		pdf.Ln(5)
		drawSignatureForm(pdf)
		pdf.Ln(5)
	}
}

func allegato3Section(pdf *fpdf.Fpdf, producerInfo, proponentInfo map[string]string, designation string) {
	setBlackBoldFont(pdf, titleTextSize)
	pdf.MultiCell(0, 3, "ALLEGATO 3 - INFORMATIVA SUL DISTRIBUTORE", "", "CM", false)
	pdf.Ln(3)
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, "Il distributore ha l’obbligo di consegnare/trasmettere al contraente il presente"+
		" documento, prima della sottoscrizione della prima proposta o, qualora non prevista, del primo contratto di "+
		"assicurazione, di metterlo a disposizione del pubblico nei propri locali, anche mediante apparecchiature "+
		"tecnologiche, oppure di pubblicarlo sul proprio sito internet ove utilizzato per la promozione e collocamento "+
		"di prodotti assicurativi, dando avviso della pubblicazione nei propri locali. In occasione di rinnovo o "+
		"stipula di un nuovo contratto o di qualsiasi operazione avente ad oggetto un prodotto di investimento "+
		"assicurativo il distributore consegna o trasmette le informazioni di cui all’Allegato 3 solo in caso di "+
		"successive modifiche di rilievo delle stesse.", "", "", false)
	pdf.Ln(3)

	setBlackBoldFont(pdf, titleTextSize)
	pdf.MultiCell(0, 3, "SEZIONE I - Informazioni generali sull’intermediario che entra in contatto con "+
		"il contraente", "", "", false)
	pdf.Ln(1)

	woptaInfoTable(pdf, producerInfo, proponentInfo, designation)
	pdf.Ln(1)

	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, "Gli estremi identificativi e di iscrizione dell’Intermediario e dei soggetti che "+
		"operano per lo stesso possono essere verificati consultando il Registro Unico degli Intermediari assicurativi "+
		"e riassicurativi sul sito internet dell’IVASS (www.ivass.it)", "", fpdf.AlignLeft, false)
	pdf.Ln(3)
	setBlackBoldFont(pdf, titleTextSize)
	pdf.MultiCell(0, 3, "SEZIONE II - Informazioni sull’attività svolta dall’intermediario assicurativo ",
		"", fpdf.AlignLeft, false)
	pdf.Ln(1)

	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, "La Wopta Assicurazioni Srl comunica di aver messo a disposizione nei propri "+
		"locali l’elenco degli obblighi di comportamento cui adempie, come indicati nell’allegato 4-ter del Regolamento"+
		" IVASS n. 40/2018.", "", "", false)
	pdf.Ln(1)
	pdf.MultiCell(0, 3, "Si comunica che nel caso di offerta fuori sede o nel caso in cui la fase "+
		"precontrattuale si svolga mediante tecniche di comunicazione a distanza il contraente riceve l’elenco "+
		"degli obblighi.", "", "", false)
	pdf.Ln(3)
	setBlackBoldFont(pdf, titleTextSize)
	pdf.MultiCell(0, 3, "SEZIONE III - Informazioni relative a potenziali situazioni di conflitto "+
		"d’interessi", "", "", false)
	pdf.Ln(1)
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, "Wopta Assicurazioni Srl ed i soggetti che operano per la stessa non sono "+
		"detentori di una partecipazione, diretta o indiretta, pari o superiore al 10% del capitale sociale o dei "+
		"diritti di voto di alcuna Impresa di assicurazione.", "", "", false)
	pdf.Ln(1)
	pdf.MultiCell(0, 3, "Le Imprese di assicurazione o Imprese controllanti un’Impresa di assicurazione "+
		"non sono detentrici di una partecipazione, diretta o indiretta, pari o superiore al 10% del capitale sociale "+
		"o dei diritti di voto dell’Intermediario.", "", "", false)
	pdf.Ln(3)
	setBlackBoldFont(pdf, titleTextSize)
	pdf.MultiCell(0, 3, "SEZIONE IV - Informazioni sugli strumenti di tutela del contraente",
		"", "", false)
	pdf.Ln(1)
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, "L’attività di distribuzione è garantita da un contratto di assicurazione della "+
		"responsabilità civile che copre i danni arrecati ai contraenti da negligenze ed errori professionali "+
		"dell’intermediario o da negligenze, errori professionali ed infedeltà dei dipendenti, dei collaboratori o "+
		"delle persone del cui operato l’intermediario deve rispondere a norma di legge.",
		"", "", false)
	pdf.Ln(1)
	pdf.MultiCell(0, 3, "Il contraente ha la facoltà, ferma restando la possibilità di rivolgersi "+
		"all’Autorità Giudiziaria, di inoltrare reclamo per iscritto all’intermediario, via posta all’indirizzo di "+
		"sede legale o a mezzo mail alla PEC sopra indicati, oppure all’Impresa secondo le modalità e presso i "+
		"recapiti indicati nel DIP aggiuntivo nella relativa sezione, nonché la possibilità, qualora non dovesse "+
		"ritenersi soddisfatto dall’esito del reclamo o in caso di assenza di riscontro da parte dell’intermediario "+
		"o dell’impresa entro il termine di legge, di rivolgersi all’IVASS secondo quanto indicato nei DIP aggiuntivi.",
		"", "", false)
	pdf.Ln(1)
	pdf.MultiCell(0, 3, "Il contraente ha la facoltà di avvalersi di altri eventuali sistemi alternativi "+
		"di risoluzione delle controversie previsti dalla normativa vigente nonché quelli indicati nei DIP aggiuntivi.",
		"", "", false)
}

func allegato4Section(pdf *fpdf.Fpdf, producerInfo, proponentInfo map[string]string, designation, section1Info string) {
	setBlackBoldFont(pdf, titleTextSize)
	pdf.MultiCell(0, 3, "ALLEGATO 4 - INFORMAZIONI SULLA DISTRIBUZIONE\nDEL PRODOTTO ASSICURATIVO NON IBIP",
		"", "CM", false)
	pdf.Ln(3)
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, "Il distributore ha l’obbligo di consegnare o trasmettere al contraente, prima "+
		"della sottoscrizione di ciascuna proposta o, qualora non prevista, di ciascun contratto assicurativo, il "+
		"presente documento, che contiene notizie sul modello e l’attività di distribuzione, sulla consulenza fornita "+
		"e sulle remunerazioni percepite.", "", "", false)
	pdf.Ln(1)

	woptaInfoTable(pdf, producerInfo, proponentInfo, designation)
	pdf.Ln(3)

	setBlackBoldFont(pdf, titleTextSize)
	pdf.MultiCell(0, 3, "SEZIONE I - Informazioni sul modello di distribuzione", "",
		"", false)
	pdf.Ln(1)
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, section1Info, "", "", false)
	pdf.Ln(3)
	setBlackBoldFont(pdf, titleTextSize)
	pdf.MultiCell(0, 3, "SEZIONE II: Informazioni sull’attività di distribuzione e consulenza",
		"", "", false)
	pdf.Ln(1)
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, "Nello svolgimento dell’attività di distribuzione, l’intermediario non presta "+
		"attività di consulenza prima della conclusione del contratto né fornisce al contraente una raccomandazione "+
		"personalizzata ai sensi dell’art. 119-ter, comma 3, del decreto legislativo n. 209/2005 "+
		"(Codice delle Assicurazioni Private)", "", "", false)
	pdf.Ln(1)
	pdf.MultiCell(0, 3, "L'attività di distribuzione assicurativa è svolta in assenza di obblighi "+
		"contrattuali che impongano di offrire esclusivamente i contratti di una o più imprese di "+
		"assicurazioni.", "", "", false)
	pdf.Ln(3)
	setBlackBoldFont(pdf, titleTextSize)
	pdf.MultiCell(0, 3, "SEZIONE III - Informazioni relative alle remunerazioni", "", "", false)
	pdf.Ln(1)
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, "Per il prodotto intermediato, è corrisposto all’intermediario, da parte "+
		"dell’impresa di assicurazione, un compenso sotto forma di commissione inclusa nel premio "+
		"assicurativo.", "", "", false)
	pdf.Ln(1)
	pdf.MultiCell(0, 3, "L’informazione sopra resa riguarda i compensi complessivamente percepiti da tutti "+
		"gli intermediari coinvolti nella distribuzione del prodotto.", "", "", false)
	pdf.Ln(3)
	setBlackBoldFont(pdf, titleTextSize)
	pdf.MultiCell(0, 3, "SEZIONE IV – Informazioni sul pagamento dei premi", "", "", false)
	pdf.Ln(1)
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, "Relativamente a questo contratto i premi pagati dal Contraente "+
		"all’intermediario e le somme destinate ai risarcimenti o ai pagamenti dovuti dalle Imprese di Assicurazione, "+
		"se regolati per il tramite dell’intermediario costituiscono patrimonio autonomo e separato dal patrimonio "+
		"dello stesso.", "", "", false)
	pdf.Ln(3)
	setBlackBoldFont(pdf, titleTextSize)
	pdf.MultiCell(0, 3, "Indicare le modalità di pagamento ammesse ", "", "", false)
	pdf.Ln(1)
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, "Sono consentiti, nei confronti di Wopta, esclusivamente bonifico e strumenti di "+
		"pagamento elettronico, quali ad esempio, carte di credito e/o carte di debito, incluse le carte "+
		"prepagate.", "", "", false)
	pdf.Ln(3)
}

func allegato4TerSection(pdf *fpdf.Fpdf, producerInfo, proponentInfo map[string]string, designation string) {
	setBlackBoldFont(pdf, titleTextSize)
	pdf.MultiCell(0, 3, "ALLEGATO 4 TER - ELENCO DELLE REGOLE DI COMPORTAMENTO DEL DISTRIBUTORE",
		"", fpdf.AlignCenter, false)
	pdf.Ln(3)
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, "Il distributore ha l’obbligo di mettere a disposizione del pubblico il "+
		"presente documento nei propri locali, anche mediante apparecchiature tecnologiche, oppure pubblicarlo su "+
		"un sito internet ove utilizzato per la promozione e il collocamento di prodotti assicurativi, dando avviso "+
		"della pubblicazione nei propri locali. Nel caso di offerta fuori sede o nel caso in cui la fase "+
		"precontrattuale si svolga mediante tecniche di comunicazione a distanza, il distributore consegna o "+
		"trasmette al contraente il presente documento prima della sottoscrizione della proposta o, qualora non "+
		"prevista, del contratto di assicurazione.", "", "", false)
	pdf.Ln(1)

	woptaInfoTable(pdf, producerInfo, proponentInfo, designation)
	pdf.Ln(3)

	setBlackBoldFont(pdf, titleTextSize)
	pdf.MultiCell(0, 3, "Sezione I - Regole generali per la distribuzione di prodotti assicurativi",
		"", "", false)
	pdf.Ln(1)
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, "a. obbligo di consegna al contraente dell’allegato 3 al Regolamento IVASS "+
		"n. 40 del 2 agosto 2018, prima della sottoscrizione della prima proposta o, qualora non prevista, del primo "+
		"contratto di assicurazione, di metterlo a disposizione del pubblico nei locali del distributore, anche "+
		"mediante apparecchiature tecnologiche, e di pubblicarlo sul sito internet, ove esistente",
		"", "", false)
	pdf.MultiCell(0, 3, "b. obbligo di consegna dell’allegato 4 al Regolamento IVASS n. 40 del 2 agosto "+
		"2018, prima della sottoscrizione di ciascuna proposta di assicurazione o, qualora non prevista, del contratto "+
		"di assicurazione", "", "", false)
	pdf.MultiCell(0, 3, "c. obbligo di consegnare copia della documentazione precontrattuale e "+
		"contrattuale prevista dalle vigenti disposizioni, copia della polizza e di ogni altro atto o documento "+
		"sottoscritto dal contraente", "", "", false)
	pdf.MultiCell(0, 3, "d. obbligo di proporre o raccomandare contratti coerenti con le richieste e le "+
		"esigenze di copertura assicurativa e previdenziale del contraente o dell’assicurato, acquisendo a tal fine, "+
		"ogni utile informazione", "", "", false)
	pdf.MultiCell(0, 3, "e. obbligo di valutare se il contraente rientra nel mercato di riferimento "+
		"identificato per il contratto di assicurazione proposto e non appartiene alle categorie di clienti per i quali "+
		"il prodotto non è compatibile, nonché l’obbligo di adottare opportune disposizioni per ottenere dai produttori"+
		" le informazioni di cui all’articolo 30-decies comma 5 del Codice e per comprendere le caratteristiche e il "+
		"mercato di riferimento individuato per ciascun prodotto", "", "", false)
	pdf.MultiCell(0, 3, "f. obbligo di fornire in forma chiara e comprensibile le informazioni "+
		"oggettive sul prodotto, illustrandone le caratteristiche, la durata, i costi e i limiti della copertura ed "+
		"ogni altro elemento utile a consentire al contraente di prendere una decisione informata",
		"", "", false)
}

func generatePolicyAnnex(pdf *fpdf.Fpdf, origin string, networkNode *models.NetworkNode, policy *models.Policy) {
	if networkNode == nil || networkNode.HasAnnex || networkNode.Type == models.PartnershipNetworkNodeType {
		producerInfo := loadProducerInfo(origin, networkNode)
		proponentInfo := loadProponentInfo(networkNode)
		designation := loadDesignation(networkNode)
		annex4Section1Info := loadAnnex4Section1Info(policy, networkNode)

		pdf.AddPage()

		woptaFooter(pdf)

		allegato3Section(pdf, producerInfo, proponentInfo, designation)

		pdf.AddPage()

		allegato4Section(pdf, producerInfo, proponentInfo, designation, annex4Section1Info)

		pdf.AddPage()

		allegato4TerSection(pdf, producerInfo, proponentInfo, designation)
	}
}
