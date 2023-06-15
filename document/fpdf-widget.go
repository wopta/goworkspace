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
		y = 281
		height = 8
	case "pmi":
		footerText = ""
		logoPath = ""
	case "persona":
		footerText = "Wopta per te. Persona è un prodotto assicurativo di Global Assistance Compagnia di assicurazioni e riassicurazioni S.p.A, distribuito da Wopta Assicurazioni S.r.l"
		logoPath = lib.GetAssetPathByEnv(basePath) + "/logo_global.png"
		x = 180
		y = 280
		height = 10
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
			"all'Albo delle imprese di assicurazione tenuto dall'IVASS, in appendice Elenco I, nr. I.00149.", "", "", false)
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
		pdf.Cell(120, 3, "")
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
		offerPrice = humanize.FormatFloat("#.###,##", policy.PriceGross*12)
	} else {
		offerPrice = humanize.FormatFloat("#.###,##", policy.PriceGross)
	}
	text := "Polizza emessa a Milano il " + emitDate + " per un importo di € " + offerPrice + " quale " +
		"prima rata alla firma, il cui pagamento a saldo è da effettuarsi con i metodi di pagamento sopra indicati. " +
		"Wopta conferma avvenuto incasso e copertura della polizza dal " + startDate + "."
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
		pdf.MultiCell(0, 3, "Global Assistance Compagnia di assicurazioni e riassicurazioni S.p.A. a"+
			" Socio Unico - Capitale Sociale: Euro 5.000.000 i.v. Codice Fiscale, Partita IVA e Registro Imprese di"+
			" Milano n. 10086540159 R.E.A. n. 1345012 della C.C.I.A.A. di Milano. Sede e Direzione Generale Piazza"+
			" Diaz 6 – 20123 Milano – Italia E-mail: global.assistance@globalassistance.it PEC: "+
			"globalassistancespa@legalmail.it. Società soggetta all’attività di direzione e coordinamento di Ri-Fin "+
			"S.r.l., iscritta all’Albo dei gruppi assicurativi presso l’Ivass al n. 014. La Società è autorizzata "+
			"all’esercizio delle Assicurazioni e Riassicurazioni con D.M. del 2/8/93 n. 19619 (G.U. 7/8/93 n. 184) e"+
			" successive autorizzazioni ed è iscritta all’Albo Imprese presso l’IVASS al n. 1.00111. La Società è"+
			" soggetta alla vigilanza dell’IVASS; è possibile verificare la veridicità e la regolarità delle "+
			"autorizzazioni mediante l'accesso al sito www.ivass.it", "", "", false)
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
