package document

import (
	"github.com/dustin/go-humanize"
	"github.com/go-pdf/fpdf"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"strings"
)

func mainHeader(pdf *fpdf.Fpdf, policy *models.Policy) {
	var (
		opt                                     fpdf.ImageOptions
		logoPath, cfpi, expiryInfo, productName string
	)

	policyInfo := "Numero: " + policy.CodeCompany + "\n" +
		"Decorre dal: " + policy.StartDate.Format(dateLayout) + " ore 24:00\n" +
		"Scade il: " + policy.EndDate.Format(dateLayout) + " ore 24:00\n"

	switch policy.Name {
	case "life":
		logoPath = lib.GetAssetPathByEnv(basePath) + "/logo_vita.png"
		productName = "Vita"
		policyInfo += expiryInfo + "Non si rinnova a scadenza."
	case "pmi":
		logoPath = lib.GetAssetPathByEnv(basePath) + "/pmi.png"
		productName = "Artigiani & Imprese"
	case "persona":
		logoPath = lib.GetAssetPathByEnv(basePath) + "/persona.png"
		productName = "Persona"
		policyInfo += "Si rinnova a scadenza salvo disdetta da inviare 30 giorni prima\n" + "Prossimo pagamento "
		if policy.PaymentSplit == string(models.PaySplitMonthly) {
			policyInfo += policy.StartDate.AddDate(0, 1, 0).Format(dateLayout) + "\n"
		} else if policy.PaymentSplit == string(models.PaySplitYear) {
			policyInfo += policy.StartDate.AddDate(1, 0, 0).Format(dateLayout) + "\n"
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

	if policy.PaymentSplit == string(models.PaySplitMonthly) {
		expiryInfo = "Prima scandenza mensile il: " +
			policy.StartDate.AddDate(0, 1, 0).Format(dateLayout) + "\n"
	} else if policy.PaymentSplit == string(models.PaySplitYear) {
		expiryInfo = "Prima scadenza annuale il: " +
			policy.StartDate.AddDate(1, 0, 0).Format(dateLayout) + "\n"
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
		pdf.ImageOptions(lib.GetAssetPathByEnv(basePath)+"/ARTW_LOGO_RGB_400px.png", 170, 6, 0, 8, false, opt, 0, "")

		setBlackBoldFont(pdf, standardTextSize)
		pdf.SetXY(11, 20)
		pdf.Cell(0, 3, "I dati della tua polizza")
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

func mainMotorHeader(pdf *fpdf.Fpdf, policy *models.Policy) {
	var (
		opt                   fpdf.ImageOptions
		logoPath, productName string
	)

	switch policy.Name {
	case "gap":
		logoPath = lib.GetAssetPathByEnv(basePath) + "/logo_gap.png"
		productName = "Auto Valore Protetto"
	}

	policyInfo := "Polizza Numero: " + policy.CodeCompany + "\n" +
		"Targa Veicolo: " + policy.Assets[0].Vehicle.Plate + "\n" +
		"Decorre dal: " + policy.StartDate.Format(dateLayout) + " ore 24:00\n" +
		"Scade il: " + policy.EndDate.Format(dateLayout) + " ore 24:00"

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
		pdf.ImageOptions(lib.GetAssetPathByEnv(basePath)+"/ARTW_LOGO_RGB_400px.png", 170, 6, 0, 8, false, opt, 0, "")

		setBlackBoldFont(pdf, standardTextSize)
		pdf.SetXY(11, 20)
		pdf.Cell(0, 3, "I dati della tua polizza")
		setBlackRegularFont(pdf, standardTextSize)
		pdf.SetXY(11, pdf.GetY()+3)
		pdf.MultiCell(0, 3.5, policyInfo, "", "", false)
		pdf.Ln(8)
	})
}

func mainFooter(pdf *fpdf.Fpdf, productName string) {
	var (
		opt                  fpdf.ImageOptions
		footerText, logoPath string
		x, y, height         float64
	)

	switch productName {
	case "life":
		footerText = "Wopta per te. Vita è un prodotto assicurativo di AXA France Vie S.A. – Rappresentanza Generale per l’Italia\ndistribuito da Wopta Assicurazioni S.r.l."
		logoPath = lib.GetAssetPathByEnv(basePath) + "/axa/logo.png"
		x = 190
		y = 282.5
		height = 8
	case "pmi":
		footerText = ""
		logoPath = ""
	case "persona":
		footerText = "Wopta per te. Persona è un prodotto assicurativo di Global Assistance Compagnia di assicurazioni" +
			" e riassicurazioni S.p.A,\ndistribuito da Wopta Assicurazioni S.r.l"
		logoPath = lib.GetAssetPathByEnv(basePath) + "/logo_global.png"
		x = 180
		y = 280
		height = 10
	case "gap":
		footerText = "Wopta per te. Auto Valore Protetto è un prodotto assicurativo di Sogessur SA – Rappresentanza" +
			" Generale per l’Italia con sede in\nVia Tiziano, 32 – 20145 Milano – Iscritta alla CCIAA di Milano P.I. " +
			"07420570967 – REA MI 1957443; distribuito da Wopta Assicurazioni Srl"
		logoPath = lib.GetAssetPathByEnv(basePath) + "/logo_sogessur.png"
		x = 151
		y = 281
		height = 7
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

func axaHeader(pdf *fpdf.Fpdf) {
	pdf.SetHeaderFunc(func() {
		var opt fpdf.ImageOptions
		pdf.SetXY(-30, 7)
		opt.ImageType = "png"
		pdf.ImageOptions(lib.GetAssetPathByEnv(basePath)+"/axa/logo.png", 190, 7, 0, 8, false, opt, 0, "")
		pdf.Ln(15)
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
		pdf.ImageOptions(lib.GetAssetPathByEnv(basePath)+"/logo_sogessur.png", 160, 7, 0, 6, false,
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

func woptaHeader(pdf *fpdf.Fpdf) {
	pdf.SetHeaderFunc(func() {
		var opt fpdf.ImageOptions
		opt.ImageType = "png"
		pdf.ImageOptions(lib.GetAssetPathByEnv(basePath)+"/ARTW_LOGO_RGB_400px.png", 10, 6, 0, 15, false, opt, 0, "")
		pdf.Ln(10)
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
			"Capitale Sociale: €120.000,00", "", 0, "", false, 0, "")
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
		pdf.ImageOptions(lib.GetAssetPathByEnv(basePath)+"/logo_global_02.png", 180, 7, 0, 15, false, opt, 0, "")
		pdf.Ln(15)
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

func paymentMethodSection(pdf *fpdf.Fpdf) {
	getParagraphTitle(pdf, "Come puoi pagare il premio")
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, "I mezzi di pagamento consentiti, nei confronti di Wopta, sono esclusivamente "+
		"bonifico e strumenti di pagamento elettronico, quali ad esempio, carte di credito e/o carte di debito, "+
		"incluse le carte prepagate. Oppure può essere pagato direttamente alla Compagnia alla "+
		"stipula del contratto, via bonifico o carta di credito.", "", "", false)
}

func emitResumeSection(pdf *fpdf.Fpdf, policy *models.Policy) {
	var offerPrice string
	emitDate := policy.EmitDate.Format(dateLayout)
	startDate := policy.StartDate.Format(dateLayout)
	if policy.PaymentSplit == "monthly" {
		offerPrice = humanize.FormatFloat("#.###,##", policy.PriceGrossMonthly)
	} else {
		offerPrice = humanize.FormatFloat("#.###,##", policy.PriceGross)
	}
	text := "Polizza emessa a Milano il " + emitDate + " per un importo di € " + offerPrice + " quale " +
		"prima rata alla firma, il cui pagamento a saldo è da effettuarsi con i metodi di pagamento sopra indicati."
	switch policy.Name {
	case "life":
		text += " Wopta conferma avvenuto incasso e copertura della polizza dal " + startDate + "."
	case "persona":
		text += "\nCostituisce quietanza di pagamento la mail di conferma che Wopta invierà al Contraente."

	}

	getParagraphTitle(pdf, "Emissione polizza e pagamento della prima rata")
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, text, "", "", false)
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
		pdf.MultiCell(0, 3, "Sogessur SA – Rappresentanza Generale per l’Italia con sede in Via Tiziano, "+
			"32 – 20145 Milano – Iscritta alla CCIAA di Milano P.I. 07420570967 – REA MI 1957443 Sogecap SA – "+
			"Rappresentanza Generale per l’Italia con sede in Via Tiziano, 32 – 20145 Milano – Iscritta alla CCIAA di "+
			"Milano P.I. 07160010968 ", "", "", false)
	}
}

func personalDataHandlingSection(pdf *fpdf.Fpdf, policy *models.Policy) {
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
	pdf.Cell(0, 3, policy.EmitDate.Format(dateLayout))
	drawSignatureForm(pdf)
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
	switch companyName {
	case "global":
		setBlackBoldFont(pdf, standardTextSize)
		pdf.CellFormat(70, 3, "Global Assistance", "", 0,
			fpdf.AlignCenter, false, 0, "")
		var opt fpdf.ImageOptions
		opt.ImageType = "png"
		pdf.ImageOptions(lib.GetAssetPathByEnv(basePath)+"/firma_global.png", 25, pdf.GetY()+3, 40, 12,
			false, opt, 0, "")
	case "axa":
		setBlackBoldFont(pdf, standardTextSize)
		pdf.MultiCell(70, 3, "AXA France Vie\n(Rappresentanza Generale per l'Italia)", "",
			fpdf.AlignCenter, false)
		pdf.SetY(pdf.GetY() - 6)
		var opt fpdf.ImageOptions
		opt.ImageType = "png"
		pdf.ImageOptions(lib.GetAssetPathByEnv(basePath)+"/firma_axa.png", 35, pdf.GetY()+9, 30, 8,
			false, opt, 0, "")
	}
}

func contractWithdrawlSection(pdf *fpdf.Fpdf) {
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
	pdf.Ln(3)
}
