package test

import (
	"github.com/dustin/go-humanize"
	"github.com/go-pdf/fpdf"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"strconv"
	"strings"
)

func GetMainHeader(pdf *fpdf.Fpdf, policy models.Policy) {
	var (
		opt      fpdf.ImageOptions
		logoPath string
		cfpi     string
		nextpay  string
	)
	pathPrefix := lib.GetAssetPathByEnv("test") + "/logo_"
	logoPath = pathPrefix + "vita.png"

	contractor := policy.Contractor

	if contractor.VatCode == "" {
		cfpi = contractor.FiscalCode
	} else {
		cfpi = contractor.VatCode
	}

	if policy.PaymentSplit == "monthly" {
		nextpay = policy.StartDate.AddDate(0, 1, 0).Format(layout)
	} else {
		nextpay = policy.StartDate.AddDate(1, 0, 0).Format(layout)
	}

	policyInfo := "Numero: " + policy.CodeCompany + "\n" +
		"Decorre dal: " + policy.StartDate.Format(layout) + " ore 24:00\n" +
		"Scade il: " + policy.EndDate.Format(layout) + " ore 24:00\n" +
		"Prima scadenza annuale il: " + nextpay + "\n" +
		"Non si rinnova a scadenza."

	contractorInfo := "Contraente: " + strings.ToUpper(contractor.Surname) + " " + strings.ToUpper(contractor.Name) + "\n" +
		"C.F./P.IVA: " + strings.ToUpper(cfpi) + "\n" +
		"Indirizzo: " + strings.ToUpper(contractor.Address) + ", " + strings.ToUpper(contractor.StreetNumber) + "\n" +
		strings.ToUpper(contractor.PostalCode) + " " + strings.ToUpper(contractor.City) + " (" +
		strings.ToUpper(contractor.Province) + ")\n" + "Mail: " + contractor.Mail + "\n" +
		"Telefono: " + contractor.Phone

	pdf.SetHeaderFunc(func() {
		opt.ImageType = "png"
		pdf.ImageOptions(logoPath, 10, 6, 13, 13, false, opt, 0, "")
		pdf.SetXY(23, 7)
		pdf.SetTextColor(229, 0, 117)
		pdf.SetFont("Montserrat", "B", 18)
		pdf.Cell(10, 6, "Wopta per te")
		pdf.SetFont("Montserrat", "I", 18)
		pdf.SetXY(23, 13)
		pdf.SetTextColor(92, 89, 92)
		pdf.Cell(10, 6, "Vita")
		pdf.ImageOptions(lib.GetAssetPathByEnv("test")+"/ARTW_LOGO_RGB_400px.png", 170, 6, 0, 8, false, opt, 0, "")

		pdf.SetTextColor(0, 0, 0)
		pdf.SetFont("Montserrat", "B", 8)
		pdf.SetXY(11, 20)
		pdf.Cell(0, 3, "I dati della tua polizza")
		pdf.SetFont("Montserrat", "", 8)
		pdf.SetXY(11, pdf.GetY()+3)
		pdf.MultiCell(0, 3, policyInfo, "", "", false)

		pdf.SetFont("Montserrat", "B", 8)
		pdf.SetXY(-95, 20)
		pdf.Cell(0, 3, "I tuoi dati")
		pdf.SetFont("Montserrat", "", 8)
		pdf.SetXY(-95, pdf.GetY()+3)
		pdf.MultiCell(0, 3, contractorInfo, "", "", false)
		pdf.Ln(10)
	})
}

func GetMainFooter(pdf *fpdf.Fpdf) {
	var opt fpdf.ImageOptions

	pdf.SetFooterFunc(func() {
		pdf.SetXY(10, -13)
		pdf.SetTextColor(229, 0, 117)
		pdf.SetFont("Montserrat", "B", 6)
		pdf.MultiCell(0, 3, "Wopta per te. Vita è un prodotto assicurativo di AXA France Vie S.A. – Rappresentanza Generale per l’Italia\ndistribuito da Wopta Assicurazioni S.r.l.", "", "", false)
		opt.ImageType = "png"
		pdf.ImageOptions(lib.GetAssetPathByEnv("test")+"/logo_axa.png", 190, 283, 8, 8, false, opt, 0, "")
	})
}

func GetAxaHeader(pdf *fpdf.Fpdf) {
	pdf.SetHeaderFunc(func() {
		var opt fpdf.ImageOptions
		pdf.SetXY(-30, 7)
		opt.ImageType = "png"
		pdf.ImageOptions(lib.GetAssetPathByEnv("test")+"/logo_axa.png", 190, 7, 8, 8, false, opt, 0, "")
		pdf.Ln(15)
	})
}

func GetAxaFooter(pdf *fpdf.Fpdf) {
	pdf.SetFooterFunc(func() {
		pdf.SetXY(10, -30)
		pdf.SetFont("Montserrat", "", 8)
		pdf.SetTextColor(0, 0, 0)
		pdf.MultiCell(0, 3, "AXA France Vie (compagnia assicurativa del gruppo AXA). Indirizzo sede "+
			"legale in Francia: 313 Terrasses de l'Arche, 92727 NANTERRE CEDEX. Numero Iscrizione Registro delle "+
			"Imprese di Nanterre: 310499959. Autorizzata in Francia (Stato di origine) all'esercizio delle "+
			"assicurazioni, vigilata in Francia dalla Autorité de Contrôle Prudentiel et de Résolution (ACPR). "+
			"Numero Matricola Registre des organismes d'assurance: 5020051. // Indirizzo Rappresentanza Generale "+
			"per l'Italia: Corso Como n. 17, 20154 Milano - CF, P.IVA e N.Iscr. Reg. Imprese 08875230016 - "+
			"REA MI-2525395 - Telefono: 02-87103548 - Fax: 02-23331247 - PEC: axafrancevie@legalmail.it - sito "+
			"internet: www.clp.partners.axa/it. Ammessa ad operare in Italia in regime di stabilimento. Iscritta "+
			"all'Albo delle imprese di assicurazione tenuto dall'IVASS, in appendice Elenco I, nr. I.00149.", "", "", false)
	})
}

func GetWoptaHeader(pdf *fpdf.Fpdf) {
	pdf.SetHeaderFunc(func() {
		var opt fpdf.ImageOptions
		opt.ImageType = "png"
		pdf.ImageOptions(lib.GetAssetPathByEnv("test")+"/ARTW_LOGO_RGB_400px.png", 10, 6, 0, 15, false, opt, 0, "")
		pdf.Ln(10)
	})
}

func GetWoptaFooter(pdf *fpdf.Fpdf) {
	pdf.SetFooterFunc(func() {
		pdf.SetY(-30)
		DrawPinkHorizontalLine(pdf, 0.4)
		pdf.Ln(5)
		pdf.SetFont("Montserrat", "B", 7)
		pdf.SetTextColor(229, 0, 117)
		pdf.Cell(pdf.GetStringWidth("Wopta Assicurazioni s.r.l"), 3, "Wopta Assicurazioni s.r.l")
		pdf.Cell(120, 3, "")
		pdf.Cell(pdf.GetStringWidth("www.wopta.it"), 3, "www.wopta.it")
		pdf.Ln(3)
		pdf.SetFont("Montserrat", "", 7)
		pdf.SetTextColor(0, 0, 0)
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
			"20143 - Milano (VI)", "", 0, "", false, 0, "")
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
	})
}

func GetContractorInfoSection(pdf *fpdf.Fpdf, contractor models.User) {
	GetParagraphTitle(pdf, "La tua assicurazione è operante per il seguente Assicurato e Garanzie")
	pdf.Ln(8)
	GetContractorInfoTable(pdf, contractor)
}

func GetContractorInfoTable(pdf *fpdf.Fpdf, contractor models.User) {
	DrawPinkHorizontalLine(pdf, thickLineWidth)
	pdf.Ln(2)
	SetBlackBoldFont(pdf, standardTextSize)
	pdf.Cell(20, 2, "Cognome e Nome")
	SetBlackRegularFont(pdf, standardTextSize)
	pdf.SetX(pdf.GetX() + 24)
	pdf.Cell(20, 2, strings.ToUpper(contractor.Surname+" "+contractor.Name))
	pdf.SetX(pdf.GetX() + 60)
	SetBlackBoldFont(pdf, standardTextSize)
	pdf.Cell(10, 2, "Codice fiscale:")
	pdf.SetX(pdf.GetX() + 20)
	SetBlackRegularFont(pdf, standardTextSize)
	pdf.Cell(20, 2, strings.ToUpper(contractor.FiscalCode))
	pdf.Ln(2.5)
	DrawPinkHorizontalLine(pdf, thinLineWidth)
	pdf.Ln(2)
	SetBlackBoldFont(pdf, standardTextSize)
	pdf.Cell(20, 2, "Residente in")
	SetBlackRegularFont(pdf, standardTextSize)
	pdf.SetX(pdf.GetX() + 24)
	pdf.Cell(20, 2, strings.ToUpper(contractor.Address+", "+contractor.StreetNumber+" - "+
		contractor.PostalCode+" "+contractor.City+" ("+contractor.Province+")"))
	pdf.Ln(2.5)
	DrawPinkHorizontalLine(pdf, thinLineWidth)
	pdf.Ln(2)
	SetBlackBoldFont(pdf, standardTextSize)
	pdf.Cell(20, 2, "Mail")
	SetBlackRegularFont(pdf, standardTextSize)
	pdf.SetX(pdf.GetX() + 24)
	pdf.Cell(20, 2, contractor.Mail)
	pdf.SetX(pdf.GetX() + 60)
	SetBlackBoldFont(pdf, standardTextSize)
	pdf.Cell(10, 2, "Telefono:")
	pdf.SetX(pdf.GetX() + 20)
	SetBlackRegularFont(pdf, 9)
	pdf.Cell(20, 2, contractor.Phone)
	pdf.Ln(3)
	DrawPinkHorizontalLine(pdf, thickLineWidth)
	pdf.Ln(1)
}

func GetGuaranteesTable(pdf *fpdf.Fpdf, policy models.Policy) {
	const (
		death               = "death"
		permanentDisability = "permanent-disability"
		temporaryDisability = "temporary-disability"
		seriousIll          = "serious-ill"
	)

	slugs := []string{death, permanentDisability, temporaryDisability, seriousIll}

	guarantees := map[string]map[string]string{
		death: {
			"name":                       "Decesso",
			"sumInsuredLimitOfIndemnity": "======",
			"duration":                   "==",
			"endDate":                    "====",
			"price":                      "====",
		},
		permanentDisability: {
			"name":                       "Invalidità Totale Permanente da Infortunio o Malattia",
			"sumInsuredLimitOfIndemnity": "======",
			"duration":                   "==",
			"endDate":                    "====",
			"price":                      "====",
		},
		temporaryDisability: {
			"name":                       "Inabilità Temporanea da Infortunio o Malattia",
			"sumInsuredLimitOfIndemnity": "======",
			"duration":                   "==",
			"endDate":                    "====",
			"price":                      "====",
		},
		seriousIll: {
			"name":                       "Malattie Gravi",
			"sumInsuredLimitOfIndemnity": "======",
			"duration":                   "==",
			"endDate":                    "====",
			"price":                      "====",
		},
	}

	for _, guarantee := range GuaranteesToMap(policy.Assets[0].Guarantees) {
		guarantees[guarantee.Slug]["sumInsuredLimitOfIndemnity"] = humanize.FormatFloat("#.###,", guarantee.Value.SumInsuredLimitOfIndemnity) + " €"
		guarantees[guarantee.Slug]["duration"] = strconv.Itoa(guarantee.Value.Duration.Year)
		guarantees[guarantee.Slug]["endDate"] = policy.StartDate.AddDate(guarantee.Value.Duration.Year, 0, 0).Format(layout)
		guarantees[guarantee.Slug]["price"] = humanize.FormatFloat("#.###,##", guarantee.Value.PremiumGrossYearly) + " €"
		if guarantee.Slug != death {
			guarantees[guarantee.Slug]["price"] += " (*)"
		}
	}

	SetBlackBoldFont(pdf, standardTextSize)
	pdf.MultiCell(90, 3, "Garanzie", "", "CM", false)
	pdf.SetXY(pdf.GetX()+90, pdf.GetY()-3)
	pdf.MultiCell(30, 3, "Somma\nassicurata €", "", "CM", false)
	pdf.SetXY(pdf.GetX()+117, pdf.GetY()-6)
	pdf.MultiCell(20, 3, "Durata\nanni", "", "CM", false)
	pdf.SetXY(pdf.GetX()+140, pdf.GetY()-6)
	pdf.MultiCell(25, 3, "Scade il", "", "CM", false)
	pdf.SetXY(pdf.GetX()+169, pdf.GetY()-3)
	pdf.MultiCell(0, 3, "Premio annuale €", "", "CM", false)
	pdf.Ln(1)
	DrawPinkHorizontalLine(pdf, thinLineWidth)

	for _, slug := range slugs {
		SetBlackBoldFont(pdf, standardTextSize)
		pdf.CellFormat(90, 6, guarantees[slug]["name"], "", 0, "", false, 0, "")
		SetBlackRegularFont(pdf, standardTextSize)
		pdf.CellFormat(25, 6, guarantees[slug]["sumInsuredLimitOfIndemnity"],
			"", 0, "RM", false, 0, "")
		pdf.CellFormat(25, 6, guarantees[slug]["duration"], "", 0, "CM",
			false, 0, "")
		pdf.CellFormat(25, 6, guarantees[slug]["endDate"], "", 0, "CM", false, 0, "")
		pdf.CellFormat(0, 6, guarantees[slug]["price"], "",
			0, "RM", false, 0, "")
		pdf.Ln(5)
		DrawPinkHorizontalLine(pdf, thinLineWidth)
	}
	pdf.Ln(0.5)
	SetBlackRegularFont(pdf, smallTextSize)
	pdf.Cell(80, 3, "(*) imposte assicurative di legge incluse nella misura del 2,50% del premio imponibile")
	pdf.Ln(3)
}

func GetAvvertenzeBeneficiariSection(pdf *fpdf.Fpdf) {
	GetParagraphTitle(pdf, "Nomina dei Beneficiari e Referente terzo, per il caso di garanzia Decesso (qualora sottoscritta)")
	pdf.Ln(8)
	SetBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, "AVVERTENZE: Può scegliere se designare nominativamente i beneficiari o se "+
		"designare genericamente come beneficiari i suoi eredi legittimi e/o testamentari. In caso di mancata "+
		"designazione nominativa, la Compagnia potrà incontrare, al decesso dell’Assicurato, maggiori difficoltà "+
		"nell’identificazione e nella ricerca dei beneficiari. La modifica o revoca del/i beneficiario/i deve essere "+
		"comunicata alla Compagnia in forma scritta.\nIn caso di specifiche esigenze di riservatezza, la Compagnia "+
		"potrà rivolgersi ad un soggetto terzo (diverso dal Beneficiario) in caso di Decesso al fine di contattare "+
		"il Beneficiario designato.", "", "", false)
}

func GetBeneficiariSection(pdf *fpdf.Fpdf) {
	GetParagraphTitle(pdf, "Beneficiario")
	pdf.Ln(8)
	SetBlackRegularFont(pdf, standardTextSize)
	pdf.CellFormat(0, 3, "Io sottoscritto Assicurato, con la sottoscrizione della presente polizza, in "+
		"riferimento alla garanzia Decesso:", "", 0, "", false, 0, "")
	pdf.Ln(4)
	SetBlackDrawColor(pdf)
	pdf.SetX(11)
	pdf.CellFormat(3, 3, "X", "1", 0, "CM", false, 0, "")
	pdf.CellFormat(0, 3, "Designo genericamente quali beneficiari della prestazione i miei eredi "+
		"(legittimi e/o testamentari)", "", 0, "", false, 0, "")
	pdf.Ln(4)
	pdf.SetX(11)
	pdf.CellFormat(3, 3, "", "1", 0, "CM", false, 0, "")
	pdf.CellFormat(0, 3, "Designo nominativamente il/i seguente/i soggetto/i quale beneficiario/i della "+
		"prestazione", "", 0, "", false, 0, "")
	GetBeneficiariTable(pdf)
}

func GetBeneficiariTable(pdf *fpdf.Fpdf) {
	pdf.Ln(5)
	DrawPinkHorizontalLine(pdf, thickLineWidth)
	pdf.Ln(1)
	SetBlackBoldFont(pdf, standardTextSize)
	pdf.Cell(50, 2, "Cognome e nome")
	SetBlackRegularFont(pdf, standardTextSize)
	pdf.Cell(20, 2, "=====")
	pdf.SetX(pdf.GetX() + 60)
	SetBlackBoldFont(pdf, standardTextSize)
	pdf.Cell(20, 2, "Cod. Fisc.: ")
	SetBlackRegularFont(pdf, standardTextSize)
	pdf.Cell(20, 2, "=====")
	pdf.Ln(3)
	DrawPinkHorizontalLine(pdf, thinLineWidth)
	pdf.Ln(2)
	SetBlackBoldFont(pdf, standardTextSize)
	pdf.Cell(50, 2, "Indirizzo")
	SetBlackRegularFont(pdf, standardTextSize)
	pdf.Cell(20, 2, "=====")
	pdf.SetX(pdf.GetX() + 60)
	SetBlackBoldFont(pdf, standardTextSize)
	pdf.Ln(3)
	DrawPinkHorizontalLine(pdf, thinLineWidth)
	pdf.Ln(2)
	SetBlackBoldFont(pdf, standardTextSize)
	pdf.Cell(50, 2, "Mail")
	SetBlackRegularFont(pdf, standardTextSize)
	pdf.Cell(20, 2, "=====")
	pdf.SetX(pdf.GetX() + 60)
	SetBlackBoldFont(pdf, standardTextSize)
	pdf.Cell(20, 2, "Telefono: ")
	SetBlackRegularFont(pdf, standardTextSize)
	pdf.Cell(20, 2, "=====")
	pdf.Ln(3)
	DrawPinkHorizontalLine(pdf, thinLineWidth)
	pdf.Ln(2)
	SetBlackBoldFont(pdf, standardTextSize)
	pdf.Cell(50, 2, "Relazione con Assicurato")
	SetBlackRegularFont(pdf, standardTextSize)
	pdf.Cell(20, 2, "=====")
	pdf.Ln(3)
	DrawPinkHorizontalLine(pdf, thinLineWidth)
	pdf.Ln(2)
	pdf.Cell(165, 2, "Consenso ad invio comunicazioni da parte della Compagnia al beneficiario, prima "+
		"dell'evento Decesso:")
	pdf.Cell(20, 2, "=====")
	pdf.Ln(3)
	DrawPinkHorizontalLine(pdf, thinLineWidth)
	pdf.Ln(2)
	DrawPinkHorizontalLine(pdf, thickLineWidth)
	pdf.Ln(1)
	SetBlackBoldFont(pdf, standardTextSize)
	pdf.Cell(50, 2, "Cognome e nome")
	SetBlackRegularFont(pdf, standardTextSize)
	pdf.Cell(20, 2, "=====")
	pdf.SetX(pdf.GetX() + 60)
	SetBlackBoldFont(pdf, standardTextSize)
	pdf.Cell(20, 2, "Cod. Fisc.: ")
	SetBlackRegularFont(pdf, standardTextSize)
	pdf.Cell(20, 2, "=====")
	pdf.Ln(3)
	DrawPinkHorizontalLine(pdf, thinLineWidth)
	pdf.Ln(2)
	SetBlackBoldFont(pdf, standardTextSize)
	pdf.Cell(50, 2, "Indirizzo")
	SetBlackRegularFont(pdf, standardTextSize)
	pdf.Cell(20, 2, "=====")
	pdf.SetX(pdf.GetX() + 60)
	SetBlackBoldFont(pdf, standardTextSize)
	pdf.Ln(3)
	DrawPinkHorizontalLine(pdf, thinLineWidth)
	pdf.Ln(2)
	SetBlackBoldFont(pdf, standardTextSize)
	pdf.Cell(50, 2, "Mail")
	SetBlackRegularFont(pdf, standardTextSize)
	pdf.Cell(20, 2, "=====")
	pdf.SetX(pdf.GetX() + 60)
	SetBlackBoldFont(pdf, standardTextSize)
	pdf.Cell(20, 2, "Telefono: ")
	SetBlackRegularFont(pdf, standardTextSize)
	pdf.Cell(20, 2, "=====")
	pdf.Ln(3)
	DrawPinkHorizontalLine(pdf, thinLineWidth)
	pdf.Ln(2)
	SetBlackBoldFont(pdf, standardTextSize)
	pdf.Cell(50, 2, "Relazione con Assicurato")
	SetBlackRegularFont(pdf, standardTextSize)
	pdf.Cell(20, 2, "=====")
	pdf.Ln(3)
	DrawPinkHorizontalLine(pdf, thinLineWidth)
	pdf.Ln(2)
	pdf.Cell(165, 2, "Consenso ad invio comunicazioni da parte della Compagnia al beneficiario, prima "+
		"dell'evento Decesso:")
	pdf.Cell(20, 2, "=====")
	pdf.Ln(3)
	DrawPinkHorizontalLine(pdf, thinLineWidth)
	pdf.Ln(1)
}

func GetReferenteTerzoSection(pdf *fpdf.Fpdf) {
	GetParagraphTitle(pdf, "Referente terzo")
	pdf.Ln(8)
	GetReferenteTerzoTable(pdf)
	pdf.Ln(2)
}

func GetReferenteTerzoTable(pdf *fpdf.Fpdf) {
	DrawPinkHorizontalLine(pdf, thickLineWidth)
	pdf.Ln(1)
	SetBlackBoldFont(pdf, standardTextSize)
	pdf.Cell(50, 2, "Cognome e nome")
	SetBlackRegularFont(pdf, standardTextSize)
	pdf.Cell(20, 2, "=====")
	pdf.SetX(pdf.GetX() + 60)
	SetBlackBoldFont(pdf, standardTextSize)
	pdf.Cell(20, 2, "Cod. Fisc.: ")
	SetBlackRegularFont(pdf, standardTextSize)
	pdf.Cell(20, 2, "=====")
	pdf.Ln(3)
	DrawPinkHorizontalLine(pdf, thinLineWidth)
	pdf.Ln(2)
	SetBlackBoldFont(pdf, standardTextSize)
	pdf.Cell(50, 2, "Indirizzo")
	SetBlackRegularFont(pdf, standardTextSize)
	pdf.Cell(20, 2, "=====")
	pdf.SetX(pdf.GetX() + 60)
	SetBlackBoldFont(pdf, standardTextSize)
	pdf.Ln(3)
	DrawPinkHorizontalLine(pdf, thinLineWidth)
	pdf.Ln(2)
	SetBlackBoldFont(pdf, standardTextSize)
	pdf.Cell(50, 2, "Mail")
	SetBlackRegularFont(pdf, standardTextSize)
	pdf.Cell(20, 2, "=====")
	pdf.SetX(pdf.GetX() + 60)
	SetBlackBoldFont(pdf, standardTextSize)
	pdf.Cell(20, 2, "Telefono: ")
	SetBlackRegularFont(pdf, standardTextSize)
	pdf.Cell(20, 2, "=====")
	pdf.Ln(3)
	DrawPinkHorizontalLine(pdf, thinLineWidth)
}

func GetStatementsSection(pdf *fpdf.Fpdf, policy models.Policy) {
	statements := *policy.Statements

	GetParagraphTitle(pdf, "Dichiarazioni da leggere con attenzione prima di firmare")
	pdf.Ln(8)
	PrintStatement(pdf, statements[0])
	pdf.Ln(5)
	GetParagraphTitle(pdf, "Questionario Medico")
	pdf.Ln(8)
	for _, statement := range statements[1:] {
		PrintStatement(pdf, statement)
	}
	pdf.Ln(8)
}

func GetVisioneDocumentiSection(pdf *fpdf.Fpdf) {
	GetParagraphTitle(pdf, "Presa visione dei documenti precontrattuali e sottoscrizione Polizza")
	pdf.Ln(8)
	SetBlackRegularFont(pdf, standardTextSize)
	email := "biciodallavalle86@gmail.com"
	pdf.MultiCell(0, 3, "Ho scelto la ricezione della seguente documentazione via e-mail al seguente indirizzo "+
		"indirizzo"+email+", nonché all’utilizzo della stessa per l’invio delle comunicazioni in corso di contratto da "+
		"parte di Wopta e della Compagnia. Sono a conoscenza che, qualora volessi modificare questa mia scelta potrò "+
		"farlo scrivendo alla Compagnia, con le modalità previste "+
		"nelle Condizioni di Assicurazione.", "", "", false)
	SetBlackBoldFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, "Confermo di aver ricevuto e preso visione, prima della conclusione del contratto:", "", "", false)
	IndentedText(pdf, "1. degli Allegati 3, 4 e 4-ter, di cui al Regolamento IVASS n. 40/2018, relativi "+
		"agli obblighi informativi e di comportamento dell’Intermediario, inclusa l’informativa privacy "+
		"dell’intermediario (ai sensi dell’art. 13 del regolamento UE n. 2016/679);")
	IndentedText(pdf, "2. del Set informativo, identificato dal modello AX01.0323, contenente: "+
		"1) Documento informativo precontrattuale per i prodotti assicurativi vita diversi dai prodotti "+
		"d’investimento assicurativi (DIP Vita), Documento informativo per i prodotti assicurativi danni (DIP Danni), "+
		"Documento informativo precontrattuale aggiuntivo per i prodotti assicurativi multirischi "+
		"(DIP aggiuntivo Multirischi), di cui al Regolamento IVASS n. 41/2018; 2) Condizioni di Assicurazione "+
		"comprensive di Glossario, che dichiaro altresì di conoscere ed accettare.")
	pdf.Ln(2)
	pdf.SetX(40)
	pdf.Cell(30, 3, "AXA France Vie")
	pdf.SetX(-75)
	pdf.Cell(40, 3, "Firma del Contraente/Assicurato")
	pdf.Ln(2)
	pdf.SetX(25)
	pdf.Cell(20, 3, "(Rappresentanza Generale per l'Italia)")
	var opt fpdf.ImageOptions
	opt.ImageType = "png"
	pdf.ImageOptions(lib.GetAssetPathByEnv("test")+"/firma_axa.png", 35, pdf.GetY()+5, 30, 8, false, opt, 0, "")
	pdf.Ln(10)
	DrawBlackLine(pdf, thickLineWidth)
	pdf.Ln(5)
}

func GetOfferResumeSection(pdf *fpdf.Fpdf) {
	GetParagraphTitle(pdf, "Il premio per tutte le coperture assicurative attivate sulla polizza – Frazionamento: ANNUALE")
	pdf.Ln(8)
	SetBlackRegularFont(pdf, standardTextSize)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(40, 2, "Premio", "", 0, "", false, 0, "")
	pdf.SetX(pdf.GetX() + 20)
	pdf.CellFormat(40, 2, "Imponibile", "", 0, "", false, 0, "")
	pdf.SetX(pdf.GetX() + 15)
	pdf.CellFormat(40, 2, "Imposte Assicurative", "", 0, "", false, 0, "")
	pdf.SetX(pdf.GetX() + 15)
	pdf.CellFormat(40, 2, "Totale", "RM", 0, "", false, 0, "")
	pdf.Ln(3)
	DrawPinkHorizontalLine(pdf, 0.1)
	pdf.Ln(1)
	pdf.CellFormat(40, 2, "Annuale firma del contratto", "", 0, "", false, 0, "")
	pdf.SetX(pdf.GetX() + 20)
	pdf.CellFormat(40, 2, "€ 173,30", "", 0, "", false, 0, "")
	pdf.SetX(pdf.GetX() + 15)
	pdf.CellFormat(40, 2, "€ 1,88", "", 0, "", false, 0, "")
	pdf.SetX(pdf.GetX() + 15)
	pdf.CellFormat(40, 2, "€ 175,18", "RM", 0, "", false, 0, "")
	pdf.Ln(3)
	DrawPinkHorizontalLine(pdf, 0.1)
	pdf.Ln(3)
}

func GetPaymentResumeSection(pdf *fpdf.Fpdf) {
	GetParagraphTitle(pdf, "Pagamento dei premi successivi al primo")
	pdf.Ln(8)
	SetBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, "Il Contraente è tenuto a pagare i Premi entro 30 giorni dalle relative scadenze. "+
		"In caso di mancato pagamento del premio entro 30 giorni dalla scadenza (c.d. termine di tolleranza) "+
		"l’assicurazione è sospesa. Il contratto è risolto automaticamente in caso di mancato pagamento "+
		"del Premio entro 90 giorni dalla scadenza.", "", "", false)
	pdf.Ln(3)
	DrawPinkHorizontalLine(pdf, 0.4)
	pdf.Ln(1)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.CellFormat(pdf.GetStringWidth("Tipologia di premio"), 3, "Tipologia di premio:", "", 0, "", false, 0, "")
	pdf.SetX(pdf.GetX() + 5)
	pdf.SetFont("Montserrat", "", 9)
	pdf.SetDrawColor(0, 0, 0)
	pdf.CellFormat(3, 3, "", "1", 0, "", false, 0, "")
	pdf.SetX(pdf.GetX() + 1)
	pdf.CellFormat(pdf.GetStringWidth("naturale variabile annualmente"), 3, "naturale variabile annualmente", "", 0, "", false, 0, "")
	pdf.SetX(pdf.GetX() + 5)
	pdf.CellFormat(3, 3, "X", "1", 0, "CM", false, 0, "")
	pdf.SetX(pdf.GetX() + 1)
	pdf.CellFormat(pdf.GetStringWidth("fisso"), 3, "fisso", "", 0, "", false, 0, "")
	pdf.SetX(pdf.GetX() + 40)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.CellFormat(pdf.GetStringWidth("Frazionamento"), 3, "Frazionamento:", "", 0, "", false, 0, "")
	pdf.SetX(pdf.GetX() + 3)
	pdf.CellFormat(pdf.GetStringWidth("ANNUALE"), 3, "ANNUALE", "", 0, "", false, 0, "")
	pdf.Ln(4)
	DrawPinkHorizontalLine(pdf, 0.1)
	pdf.Ln(1)
	pdf.SetFont("Montserrat", "", 9)
	pdf.Cell(0, 3, "Il Premio è dovuto alle diverse annualità di Polizza, alle date qui sotto indicate:")
	pdf.Ln(4)
	DrawPinkHorizontalLine(pdf, 0.1)
	pdf.Ln(1)
	pdf.SetX(pdf.GetX() + 12)
	pdf.CellFormat(pdf.GetStringWidth("04/03/2030:"), 3, "Alla firma:", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 1)
	pdf.CellFormat(pdf.GetStringWidth("€ 175,18"), 3, "€ 175,18", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 12)
	pdf.CellFormat(pdf.GetStringWidth("04/03/2030:"), 3, "04/03/2028:", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 1)
	pdf.CellFormat(pdf.GetStringWidth("€ 175,18"), 3, "€ 175,18", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 12)
	pdf.CellFormat(pdf.GetStringWidth("04/03/2030:"), 3, "04/03/2033:", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 1)
	pdf.CellFormat(pdf.GetStringWidth("€ 175,18"), 3, "€ 175,18", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 12)
	pdf.CellFormat(pdf.GetStringWidth("04/03/2030:"), 3, "04/03/2038:", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 1)
	pdf.CellFormat(pdf.GetStringWidth("€ 175,18"), 3, "€ 175,18", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 12)
	pdf.Ln(4)
	DrawPinkHorizontalLine(pdf, 0.1)
	pdf.Ln(1)
	pdf.SetX(pdf.GetX() + 12)
	pdf.CellFormat(pdf.GetStringWidth("04/03/2030:"), 3, "03/04/2024", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 1)
	pdf.CellFormat(pdf.GetStringWidth("€ 175,18"), 3, "€ 175,18", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 12)
	pdf.CellFormat(pdf.GetStringWidth("04/03/2030:"), 3, "04/03/2029:", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 1)
	pdf.CellFormat(pdf.GetStringWidth("€ 175,18"), 3, "€ 175,18", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 12)
	pdf.CellFormat(pdf.GetStringWidth("04/03/2030:"), 3, "04/03/2034:", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 1)
	pdf.CellFormat(pdf.GetStringWidth("€ 175,18"), 3, "€ 175,18", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 12)
	pdf.CellFormat(pdf.GetStringWidth("04/03/2030:"), 3, "04/03/2039:", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 1)
	pdf.CellFormat(pdf.GetStringWidth("€ 175,18"), 3, "€ 175,18", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 12)
	pdf.Ln(4)
	DrawPinkHorizontalLine(pdf, 0.1)
	pdf.Ln(1)
	pdf.SetX(pdf.GetX() + 12)
	pdf.CellFormat(pdf.GetStringWidth("04/03/2030:"), 3, "03/04/2025", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 1)
	pdf.CellFormat(pdf.GetStringWidth("€ 175,18"), 3, "€ 175,18", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 12)
	pdf.CellFormat(pdf.GetStringWidth("04/03/2030:"), 3, "04/03/2030:", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 1)
	pdf.CellFormat(pdf.GetStringWidth("€ 175,18"), 3, "€ 175,18", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 12)
	pdf.CellFormat(pdf.GetStringWidth("04/03/2030:"), 3, "04/03/2035:", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 1)
	pdf.CellFormat(pdf.GetStringWidth("€ 175,18"), 3, "€ 175,18", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 12)
	pdf.CellFormat(pdf.GetStringWidth("04/03/2030:"), 3, "04/03/2040:", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 1)
	pdf.CellFormat(pdf.GetStringWidth("€ 175,18"), 3, "€ 175,18", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 12)
	pdf.Ln(4)
	DrawPinkHorizontalLine(pdf, 0.1)
	pdf.Ln(1)
	pdf.SetX(pdf.GetX() + 12)
	pdf.CellFormat(pdf.GetStringWidth("04/03/2030:"), 3, "03/04/2026", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 1)
	pdf.CellFormat(pdf.GetStringWidth("€ 175,18"), 3, "€ 175,18", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 12)
	pdf.CellFormat(pdf.GetStringWidth("04/03/2030:"), 3, "04/03/2031:", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 1)
	pdf.CellFormat(pdf.GetStringWidth("€ 175,18"), 3, "€ 175,18", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 12)
	pdf.CellFormat(pdf.GetStringWidth("04/03/2030:"), 3, "04/03/2036:", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 1)
	pdf.CellFormat(pdf.GetStringWidth("€ 175,18"), 3, "€ 175,18", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 12)
	pdf.CellFormat(pdf.GetStringWidth("04/03/2030:"), 3, "04/03/2041:", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 1)
	pdf.CellFormat(pdf.GetStringWidth("€ 175,18"), 3, "€ 175,18", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 12)
	pdf.Ln(4)
	DrawPinkHorizontalLine(pdf, 0.1)
	pdf.Ln(1)
	pdf.SetX(pdf.GetX() + 12)
	pdf.CellFormat(pdf.GetStringWidth("04/03/2030:"), 3, "03/04/2027", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 1)
	pdf.CellFormat(pdf.GetStringWidth("€ 175,18"), 3, "€ 175,18", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 12)
	pdf.CellFormat(pdf.GetStringWidth("04/03/2030:"), 3, "04/03/2032:", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 1)
	pdf.CellFormat(pdf.GetStringWidth("€ 175,18"), 3, "€ 175,18", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 12)
	pdf.CellFormat(pdf.GetStringWidth("04/03/2030:"), 3, "04/03/2037:", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 1)
	pdf.CellFormat(pdf.GetStringWidth("€ 175,18"), 3, "€ 175,18", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 12)
	pdf.CellFormat(pdf.GetStringWidth("04/03/2030:"), 3, "04/03/2042:", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 1)
	pdf.CellFormat(pdf.GetStringWidth("€ 175,18"), 3, "€ 175,18", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 12)
	pdf.Ln(4)
	DrawPinkHorizontalLine(pdf, 0.1)
	pdf.Ln(1)
	pdf.MultiCell(0, 3, "In caso di frazionamento mensile i Premi sopra riportati sono dovuti, alle date "+
		"indicate e con successiva frequenza mensile, in misura di 1/12 per ogni mensilità. Non sono previsti oneri "+
		"o interessi di frazionamento.", "", "", false)
	pdf.Ln(5)
}

func GetContractWithdrawlSection(pdf *fpdf.Fpdf) {
	GetParagraphTitle(pdf, "Informativa sul diritto di recesso")
	pdf.Ln(8)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.MultiCell(0, 3, "Diritto di recesso entro i primi 30 giorni dalla stipula ("+
		"diritto di ripensamento)", "", "", false)
	pdf.SetFont("Montserrat", "", 9)
	pdf.MultiCell(0, 3, "Il Contraente può recedere dal contratto entro il termine di 30 giorni dalla "+
		"decorrenza dell’assicurazione (diritto di ripensamento). In tal caso, l’assicurazione si intende come mai "+
		"entrata in vigore e la Compagnia, per il tramite dell’intermediario, provvederà a rimborsare al Contraente "+
		"l’importo di Premio già versato (al netto delle imposte).", "", "", false)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.MultiCell(0, 3, "Diritto di recesso annuale (disdetta alla annualità)", "", "", false)
	pdf.SetFont("Montserrat", "", 9)
	pdf.MultiCell(0, 3, "Il Contraente può recedere dal contratto annualmente, entro il termine di 30 "+
		"giorni dalla scadenza annuale della polizza (disdetta alla annualità). In tal caso, l’assicurazione cessa alle "+
		"ore 24:00 dell’ultimo giorno della annualità in corso. È possibile disdettare singolarmente una o più delle "+
		"coperture attivate in fase di sottoscrizione.", "", "", false)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.MultiCell(0, 3, "Modalità per l’esercizio del diritto di recesso", "", "", false)
	pdf.SetFont("Montserrat", "", 9)
	pdf.MultiCell(0, 3, "Il Contraente è tenuto ad esercitare il diritto di recesso mediante invio di una "+
		"lettera raccomandata a.r. al seguente indirizzo: Wopta Assicurazioni srl – Gestione Portafoglio – Galleria del "+
		"Corso, 1 – 201212 Milano (MI) oppure via posta elettronica certificata (PEC) all’indirizzo "+
		"email: woptaassicurazioni@legalmail.it", "", "", false)
	pdf.Ln(5)
}

func GetPaymentMethodSection(pdf *fpdf.Fpdf) {
	GetParagraphTitle(pdf, "Come puoi pagare il premio")
	pdf.Ln(8)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Montserrat", "", 9)
	pdf.MultiCell(0, 3, "I mezzi di pagamento consentiti, nei confronti di Wopta, sono esclusivamente "+
		"bonifico e strumenti di pagamento elettronico, quali ad esempio, carte di credito e/o carte di debito, "+
		"incluse le carte prepagate. Oppure può essere pagato direttamente alla Compagnia alla "+
		"stipula del contratto, via bonifico o carta di credito.", "", "", false)
	pdf.Ln(5)
}

func GetEmitResumeSection(pdf *fpdf.Fpdf) {
	GetParagraphTitle(pdf, "Emissione polizza e pagamento della prima rata")
	pdf.Ln(8)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Montserrat", "", 9)
	pdf.MultiCell(0, 3, "Polizza emessa a Milano il 03/04/2023 per un importo di €  175,18 quale prima "+
		"rata alla firma, il cui pagamento a saldo è da effettuarsi con i metodi di pagamento sopra indicati. "+
		"Wopta conferma avvenuto incasso e copertura della polizza dal 03/04/2023.", "", "", false)
	pdf.Ln(5)
}

func GetPolicyDescriptionSection(pdf *fpdf.Fpdf) {
	GetParagraphTitle(pdf, "Per noi questa polizza fa al caso tuo")
	pdf.Ln(8)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.MultiCell(0, 3, "Richieste ed esigenze di copertura assicurativa del contraente", "", "", false)
	pdf.SetFont("Montserrat", "", 9)
	pdf.MultiCell(0, 3, "In Polizza sono riportate le tue dichiarazioni relative al rischio. Sulla base di "+
		"tali dichiarazioni, esigenze e richieste, le soluzioni assicurative individuate in coerenza con esse, sono:", "", "", false)
	pdf.SetX(15)
	pdf.MultiCell(0, 3, "- tutelare dei soggetti cari, in caso di decesso dell’Assicurato nel corso della "+
		"durata della copertura, attraverso un sostegno economico indennizzato ai Beneficiari, che "+
		"sono stati indicati dal Contraente;", "", "", false)
	pdf.SetX(15)
	pdf.MultiCell(0, 3, "- difendere il proprio reddito, potendo beneficiare di un capitale in caso di "+
		"perdita, da parte dell’Assicurato, definitiva ed irrimediabile, da Infortunio o Malattia, della capacità di "+
		"attendere a un qualsiasi lavoro proficuo, in misura totale (almeno del 60%);", "", "", false)
	pdf.SetX(15)
	pdf.MultiCell(0, 3, "- difendere il proprio reddito attraverso una indennità mensile, in caso di "+
		"perdita totale, ma in via temporanea, delle capacità dell’Assicurato di attendere alla propria professione o "+
		"attività lavorativa a seguito di Infortunio o Malattia", "", "", false)
	pdf.SetX(15)
	pdf.MultiCell(0, 3, "- tutelare la propria salute attraverso una indennità in caso di diagnosi, in "+
		"capo all’Assicurato, di una delle seguenti malattie: cancro, attacco cardiaco (infarto del miocardio), "+
		"chirurgia aorto-coronarica (bypass), ictus, insufficienza renale (fase finale di malattia renale), trapianto "+
		"di organi principali (cuore, polmone, fegato, pancreas, rene o midollo osseo);", "", "", false)
	pdf.SetX(15)
	pdf.MultiCell(0, 3, "- le coperture operano con un orizzonte l’orizzonte temporale di protezione "+
		"indicato a pagina 1 per ogni garanzia sottoscritta e tale durata viene valutata congrua "+
		"con le esigenze di protezione;", "", "", false)
	pdf.MultiCell(0, 3, "non rilevando interesse per altre eventuali coperture previste dal prodotto, ma "+
		"non incluse in questa Polizza. Il Contraente è stato informato che la Polizza può prevedere, in relazione "+
		"alle garanzie che precedono, l’applicazione di Scoperti, Franchigie, Limiti di indennizzo ed esclusioni, "+
		"meglio riportate nelle Condizioni Generali di Assicurazione, e che sono stati da te valutati in linea con la "+
		"capacità finanziaria di sostenere in proprio tale livello di danno e rischio.", "", "", false)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.MultiCell(0, 3, "Con la seguente sottoscrizione dichiari che quanto sopra corrisponde a quanto "+
		"illustrato dall’Intermediario, il quale ha fornito ogni altro elemento utile a consentirti di prendere una "+
		"decisione informata e coerente con le esigenze espresse.", "", "", false)
	pdf.Ln(5)
	pdf.SetX(-80)
	pdf.Cell(0, 3, "Firma del Contraente/Assicurato")
	pdf.Ln(15)
	pdf.SetDrawColor(0, 0, 0)
	pdf.SetLineWidth(0.4)
	pdf.Line(130, pdf.GetY(), 190, pdf.GetY())
	pdf.Ln(5)
}

func GetWoptaAxaCompanyDescriptionSection(pdf *fpdf.Fpdf) {
	GetParagraphTitle(pdf, "Chi siamo")
	pdf.Ln(8)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Montserrat", "", 8)
	pdf.MultiCell(0, 3, "Wopta Assicurazioni S.r.l. - intermediario assicurativo, soggetto al controllo "+
		"dell’IVASS ed iscritto dal 14.02.2022 al Registro Unico degli Intermediari, in Sezione A nr. A000701923, "+
		"avente sede legale in Galleria del Corso, 1 – 20122 Milano (MI). Capitale sociale Euro 120.000 - "+
		"Codice Fiscale, Reg. Imprese e Partita IVA: 12072020964 - Iscritta al Registro delle imprese di Milano – "+
		"REA MI 2638708", "", "", false)
	pdf.Ln(5)
	pdf.MultiCell(0, 3, "AXA France Vie (compagnia assicurativa del gruppo AXA). Indirizzo sede legale in "+
		"Francia: 313 Terrasses de l'Arche, 92727 NANTERRE CEDEX. Numero Iscrizione Registro delle Imprese di "+
		"Nanterre: 310499959. Autorizzata in Francia (Stato di origine) all’esercizio delle assicurazioni, vigilata "+
		"in Francia dalla Autorité de Contrôle Prudentiel et de Résolution (ACPR). Numero Matricola Registre des "+
		"organismes d’assurance: 5020051. // Indirizzo Rappresentanza Generale per l’Italia: Corso Como n. 17, 20154 "+
		"Milano - CF, P.IVA e N.Iscr. Reg. Imprese 08875230016 - REA MI-2525395 - Telefono: 02-87103548 - "+
		"Fax: 02-23331247 - PEC: axafrancevie@legalmail.it – sito internet: www.clp.partners.axa/it. Ammessa ad "+
		"operare in Italia in regime di stabilimento. Iscritta all’Albo delle imprese di assicurazione tenuto "+
		"dall’IVASS, in appendice Elenco I, nr. I.00149.", "", "", false)
}

func GetAxaDeclarationsConsentSection(pdf *fpdf.Fpdf) {
	pdf.SetFont("Montserrat", "B", 9)
	pdf.SetTextColor(0, 0, 0)
	pdf.Cell(0, 3, "DICHIARAZIONI E CONSENSI")
	pdf.Ln(3)
	pdf.SetFont("Montserrat", "", 9)
	pdf.MultiCell(0, 3, "Io sottoscritto, dopo aver letto l’Informativa Privacy della compagnia titolare "+
		"del trattamento redatta ai sensi del Regolamento (UE) 2016/679 (relativo alla protezione delle persone "+
		"fisiche con riguardo al trattamento dei dati personali), della quale confermo ricezione, PRESTO IL CONSENSO "+
		"al trattamento dei miei dati personali, ivi inclusi quelli eventualmente da me conferiti in riferimento al "+
		"mio stato di salute, per le finalità indicate nell’informativa, nonché alla loro comunicazione, per "+
		"successivo trattamento, da parte dei soggetti indicati nella informativa predetta.", "", "", false)
	pdf.Ln(3)
	pdf.SetFont("Montserrat", "B", 8)
	pdf.Cell(0, 3, "Resta inteso che in caso di negazione del consenso non sarà possibile "+
		"finalizzare il rapporto contrattuale assicurativo.")
	pdf.Ln(3)
	pdf.SetDrawColor(0, 0, 0)
	pdf.SetLineWidth(0.1)
	pdf.Line(11, pdf.GetY(), 180, pdf.GetY())
	pdf.Ln(5)
	pdf.Cell(0, 3, "03/04/2023")
	pdf.Ln(8)
	pdf.SetX(-80)
	pdf.Cell(0, 3, "Firma del Contraente/Assicurato")
	pdf.Ln(15)
	pdf.SetDrawColor(0, 0, 0)
	pdf.SetLineWidth(0.4)
	pdf.Line(130, pdf.GetY(), 190, pdf.GetY())
}

func GetAxaTableSection(pdf *fpdf.Fpdf) {
	pdf.SetFont("Montserrat", "B", 12)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFillColor(229, 0, 117)
	pdf.MultiCell(0, 6, "MODULO PER L’IDENTIFICAZIONE E L’ADEGUATA VERIFICA DELLA CLIENTELA", "LTR", "CM", true)
	pdf.SetFont("Montserrat", "B", 8)
	pdf.MultiCell(0, 4, "POLIZZA DI RAMO VITA I  - Polizza “Wopta per te. Vita”", "LR", "CM", true)
	pdf.SetFont("Montserrat", "I", 6)
	pdf.MultiCell(0, 3, "(da compilarsi in caso di scelta da parte del Contraente/Assicurato della garanzia Decesso)", "LBR", "CM", true)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Montserrat", "", 8)
	pdf.MultiCell(0, 2, "", "LR", "", false)
	pdf.MultiCell(0, 3, "AVVERTENZA PRELIMINARE - Al fine di adempiere agli obblighi previsti dal "+
		"Decreto Legislativo 21 novembre 2007 n. 231 (di seguito il “Decreto”), in materia di prevenzione "+
		"del fenomeno del riciclaggio e del finanziamento del terrorismo, il Cliente (il soggetto Contraente/Assicurato "+
		"alla polizza “Wopta per te. Vita”) è tenuto a compilare e sottoscrivere il presente Modulo. Le "+
		"disposizioni del Decreto richiedono infatti, per una completa identificazione ed una adeguata conoscenza del "+
		"cliente e dell’eventuale titolare effettivo, la raccolta di informazioni ulteriori rispetto a quelle "+
		"anagrafiche già raccolte. La menzionata normativa impone al cliente di fornire, sotto la propria "+
		"responsabilità, tutte le informazioni necessarie ed aggiornate per consentire all’Intermediario di adempiere "+
		"agli obblighi di adeguata verifica e prevede specifiche sanzioni nel caso in cui le informazioni non "+
		"vengano fornite o risultino false.", "LR", "", false)
	pdf.MultiCell(0, 3, "", "LR", "", false)
	pdf.MultiCell(0, 3, "Il conferimento dei dati e delle informazioni personali per l’identificazione "+
		"del Cliente e per la compilazione della presente sezione è obbligatorio per legge e, in caso di loro mancato "+
		"rilascio, la Compagnia Assicurativa non potrà procedere ad instaurare il rapporto (c.d. obbligo di "+
		"astensione), e dovrà valutare se effettuare una segnalazione alle autorità competenti (Unità di "+
		"Informazione Finanziaria presso Banca d’Italia e Guardia di Finanza). I dati saranno trattati per le "+
		"finalità di assolvimento degli obblighi previsti dalla normativa antiriciclaggio e, pertanto, tale "+
		"trattamento non richiede il consenso dell’interessato.", "LR", "", false)
	pdf.MultiCell(0, 3, "", "LR", "", false)
	pdf.MultiCell(0, 3, "Io sottoscritto DALLA VALLE FABRIZIO (Contraente/Assicurato), letta l’Avvertenza "+
		"Preliminare di cui sopra e l’Informativa sui Riferimenti Normativi Antiriciclaggio (in calce al presente "+
		"modulo), al fine di permettere all’Intermediario di assolvere agli obblighi di adeguata verifica di cui al "+
		"D.Lgs. n. 231/2007 in materia di prevenzione dei fenomeni di riciclaggio e di finanziamento del terrorismo, "+
		"in relazione all’instaurazione del rapporto assicurativo di cui al contratto di assicurazione “Wopta per te. "+
		"Vita” - che prevede una garanzia di ramo vita emessa dall’impresa AXA France VIE S.A. (Rappresentanza "+
		"Generale per l’Italia):", "LR", "", false)
	pdf.MultiCell(0, 4, "", "LR", "", false)
	pdf.MultiCell(0, 4, "A. dichiaro che i seguenti dati riportati relativi alla mia persona "+
		"corrispondono al vero ", "LR", "", false)
	pdf.CellFormat(5, 4, "", "LR", 0, "", false, 0, "")
	pdf.SetFont("Montserrat", "B", 9)
	pdf.SetTextColor(255, 255, 255)
	pdf.CellFormat(180, 4, "DATI IDENTIFICATIVI DEL CLIENTE (CONTRAENTE/ASSICURATO)", "TLR",
		0, "CM", true, 0, "")
	pdf.CellFormat(5, 4, "", "LR", 1, "", false, 0, "")
	pdf.CellFormat(5, 4, "", "LR", 0, "", false, 0, "")
	pdf.SetFont("Montserrat", "B", 8)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(90, 4, "Nome: FABRIZIO", "TLR", 0, "", false, 0, "")
	pdf.CellFormat(90, 4, "Cognome:  DALLA VALLE", "TLR", 0, "", false, 0, "")
	pdf.CellFormat(5, 4, "", "LR", 1, "", false, 0, "")
	pdf.CellFormat(5, 4, "", "LR", 0, "", false, 0, "")
	pdf.CellFormat(90, 4, "Data di nascita: 07/10/1994", "TLR", 0, "", false, 0, "")
	pdf.CellFormat(90, 4, "Codice Fiscale: DLLFRZ86T09E970X", "TLR", 0, "", false, 0, "")
	pdf.CellFormat(5, 4, "", "LR", 1, "", false, 0, "")
	pdf.CellFormat(5, 4, "", "LR", 0, "", false, 0, "")
	pdf.CellFormat(90, 4, "Comune di nascita: PIANEZZE", "TLR", 0, "", false, 0, "")
	pdf.CellFormat(45, 4, "CAP: 36060", "TLR", 0, "", false, 0, "")
	pdf.CellFormat(45, 4, "Prov.: MI", "TLR", 0, "", false, 0, "")
	pdf.CellFormat(5, 4, "", "LR", 1, "", false, 0, "")
	pdf.CellFormat(5, 4, "", "LR", 0, "", false, 0, "")
	pdf.CellFormat(90, 4, "Comune di residenza: PIANEZZE", "TLR", 0, "", false, 0, "")
	pdf.CellFormat(45, 4, "CAP: 36060", "TLR", 0, "", false, 0, "")
	pdf.CellFormat(45, 4, "Prov.: VI", "TLR", 0, "", false, 0, "")
	pdf.CellFormat(5, 4, "", "LR", 1, "", false, 0, "")
	pdf.CellFormat(5, 4, "", "LR", 0, "", false, 0, "")
	pdf.CellFormat(180, 4, "Indirizzo di residenza: VIA TEZZE, 45/A", "TLR", 0, "", false, 0, "")
	pdf.CellFormat(5, 4, "", "LR", 1, "", false, 0, "")
	pdf.CellFormat(5, 4, "", "LR", 0, "", false, 0, "")
	pdf.CellFormat(180, 4, "Comune di domicilio (se diverso dalla residenza:", "TLR", 0, "", false, 0, "")
	pdf.CellFormat(5, 4, "", "LR", 1, "", false, 0, "")
	pdf.CellFormat(5, 4, "", "LR", 0, "", false, 0, "")
	pdf.CellFormat(180, 4, "Indirizzo di domicilio (se diverso dalla residenza:", "LR", 0, "", false, 0, "")
	pdf.CellFormat(5, 4, "", "LR", 1, "", false, 0, "")
	pdf.CellFormat(5, 4, "", "LR", 0, "", false, 0, "")
	pdf.CellFormat(180, 4, "Status occupazinale: Lavoratore/dipendente", "TLR", 0, "", false, 0, "")
	pdf.CellFormat(5, 4, "", "LR", 1, "", false, 0, "")
	pdf.CellFormat(5, 4, "", "LR", 0, "", false, 0, "")
	pdf.CellFormat(180, 4, "Se Altro (specificare):", "BLR", 0, "", false, 0, "")
	pdf.CellFormat(5, 4, "", "LR", 1, "", false, 0, "")
	pdf.SetFont("Montserrat", "", 8)
	pdf.MultiCell(0, 1, "", "LR", "", false)
	pdf.MultiCell(0, 4, "B. allego una fotocopia fronte/retro del mio documento di identità non scaduto "+
		"avente i seguenti estremi, confermando la veridicità dei dati sotto riportati: ", "LR", "", false)
	pdf.CellFormat(5, 4, "", "LR", 0, "", false, 0, "")
	pdf.CellFormat(90, 4, "Tipo documento: 01 = Carta di identità", "TLR", 0, "", false, 0, "")
	pdf.CellFormat(90, 4, "Nr. Documento: AT9045321", "TLR", 0, "", false, 0, "")
	pdf.CellFormat(5, 4, "", "LR", 1, "", false, 0, "")
	pdf.CellFormat(5, 4, "", "LR", 0, "", false, 0, "")
	pdf.CellFormat(90, 4, "Ente di rilascio: Comune", "TLR", 0, "", false, 0, "")
	pdf.CellFormat(90, 4, "Data di rilascio: 20/10/2013", "TLR", 0, "", false, 0, "")
	pdf.CellFormat(5, 4, "", "LR", 1, "", false, 0, "")
	pdf.CellFormat(5, 4, "", "LR", 0, "", false, 0, "")
	pdf.CellFormat(90, 4, "Località di rilascio: PIANEZZE", "1", 0, "", false, 0, "")
	pdf.CellFormat(90, 4, "Data di scadenza: 07/10/2023", "1", 0, "", false, 0, "")
	pdf.CellFormat(5, 4, "", "LR", 1, "", false, 0, "")
	pdf.MultiCell(0, 1, "", "LR", "", false)
	pdf.MultiCell(0, 4, "C. dichiaro di NON essere una Persona Politicamente Esposta", "LR", "", false)
	pdf.CellFormat(4, 4, "", "L", 0, "", false, 0, "")
	pdf.CellFormat(0, 4, "In caso di risposta affermativa indicare la tipologia:", "R", 1, "", false, 0, "")
	pdf.MultiCell(0, 4, "D. dichiaro di NON essere destinatario di misure di congelamento dei fondi e risorse economiche", "LR", "", false)
	pdf.CellFormat(4, 4, "", "L", 0, "", false, 0, "")
	pdf.CellFormat(0, 4, "In caso di risposta affermativa indicare il motivo:", "R", 1, "", false, 0, "")
	pdf.MultiCell(0, 4, "E. dichiaro di NON essere sottoposto a procedimenti o di NON aver subito condanne "+
		"per reati in materia economica/ finanziaria/tributaria/societaria", "LR", "", false)
	pdf.CellFormat(4, 4, "", "L", 0, "", false, 0, "")
	pdf.CellFormat(0, 4, "In caso di risposta affermativa indicare il motivo:", "R", 1, "", false, 0, "")
	pdf.MultiCell(0, 4, "F. dichiaro ai fini dell'identificazione del Titolare Effettivo, di essere una "+
		"persona fisica che agisce in nome e per conto proprio, di essere il soggetto Contraente/Assicurato, e "+
		"quindi che non esiste il titolare effettivo", "LR", "", false)
	pdf.MultiCell(0, 4, "G. fornisco, con riferimento allo scopo e alla natura prevista del rapporto "+
		"continuativo, le seguenti informazioni", "LR", "", false)
	pdf.CellFormat(4, 8, "", "L", 0, "", false, 0, "")
	pdf.MultiCell(0, 4, "i. Tipologia di rapporto continuativo (informazione immediatamente desunta dal "+
		"rapporto): Stipula di un contratto di assicurazione di puro rischio che prevede garanzia di ramo vita "+
		"(caso morte Assicurato)", "R", "", false)
	pdf.CellFormat(4, 8, "", "L", 0, "", false, 0, "")
	pdf.MultiCell(0, 4, "ii. Scopo prevalente del rapporto continuativo in riferimento alle garanzie vita"+
		" (informazione immediatamente desunta dal rapporto):Protezione assicurativa al fine di garantire ai "+
		"beneficiari un capitale qualora si verifichi l’evento oggetto di copertura", "R", "", false)
	pdf.CellFormat(4, 4, "", "L", 0, "", false, 0, "")
	pdf.MultiCell(0, 4, "iii.  Origine dei fondi utilizzati per il pagamento dei premi assicurativi: "+
		"Proprie risorse economiche", "R", "", false)
	pdf.CellFormat(0, 2, "", "BLR", 1, "", false, 0, "")
}

func GetAxaTablePart2Section(pdf *fpdf.Fpdf) {
	pdf.CellFormat(0, 2, "", "TLR", 1, "", false, 0, "")
	pdf.MultiCell(0, 4, "Il sottoscritto, ai sensi degli artt. 22 e 55 comma 3 del d.lgs. 231/2007, "+
		"consapevole della responsabilità penale derivante da omesse e/o mendaci affermazioni, dichiara che tutte le "+
		"informazioni fornite (anche in riferimento al titolare effettivo), le dichiarazioni rilasciate il documento "+
		"di identità che allego, ed i dati riprodotti negli appositi campi del Modulo di Polizza corrispondono al "+
		"vero. Il sottoscritto si assume tutte le responsabilità di natura civile, amministrativa e penale per "+
		"dichiarazioni non veritiere. Il sottoscritto si impegna a comunicare senza ritardo a AXA France VIE S.A. "+
		"(Rappresentanza Generale per l’Italia) ogni eventuale integrazione o variazione che si dovesse verificare "+
		"in relazione ai dati ed alle informazioni forniti con il presente modulo.", "LR", "", false)
	pdf.SetFont("Montserrat", "B", 8)
	pdf.CellFormat(0, 4, "", "LR", 1, "", false, 0, "")
	pdf.CellFormat(30, 4, "Data 03/04/2023", "L", 0, "CM", false, 0, "")
	pdf.CellFormat(100, 4, "", "", 0, "", false, 0, "")
	pdf.CellFormat(60, 4, "Firma del contraente/Assicurato", "R", 1, "CM", false, 0, "")
	pdf.CellFormat(0, 30, "", "BLR", 1, "", false, 0, "")
	pdf.SetDrawColor(0, 0, 0)
	pdf.SetLineWidth(0.4)
	pdf.Line(147, pdf.GetY()-20, 195, pdf.GetY()-20)
}

func GetAxaTablePart3Section(pdf *fpdf.Fpdf) {
	SetBlackBoldFont(pdf, 10)
	pdf.MultiCell(0, 2, "", "TLR", "", false)
	pdf.MultiCell(0, 4, "Informativa antiriciclaggio (articoli di riferimento) - "+
		"(Decreto legislativo n. 231/2007)", "LR", "CM", false)
	pdf.MultiCell(0, 3, "", "LR", "", false)
	pdf.SetFont("Montserrat", "B", 6)
	pdf.MultiCell(0, 2.5, "Obbligo di astensione – art. 42", "LR", "", false)
	pdf.SetFont("Montserrat", "", 6)
	pdf.MultiCell(0, 2.5, "1. I soggetti obbligati che si trovano nell’impossibilità oggettiva di "+
		"effettuare l'adeguata verifica della clientela, ai sensi delle disposizioni di cui all'articolo 19, "+
		"comma 1, lettere a), b) e c), si astengono dall'instaurare, eseguire ovvero proseguire il rapporto, la "+
		"prestazione professionale e le operazioni e valutano se effettuare una segnalazione di operazione sospetta "+
		"alla UIF a norma dell'articolo 35.", "LR", "", false)
	pdf.SetFont("Montserrat", "", 6)
	pdf.MultiCell(0, 2.5, "2. I soggetti obbligati si astengono dall'instaurare il rapporto continuativo, "+
		"eseguire operazioni o prestazioni professionali e pongono fine al rapporto continuativo o alla prestazione "+
		"professionale già in essere di cui siano, direttamente o indirettamente, parte società fiduciarie, trust, "+
		"società anonime o controllate attraverso azioni al portatore aventi sede in Paesi terzi ad alto rischio. "+
		"Tali misure si applicano anche nei confronti delle ulteriori entità giuridiche, altrimenti denominate, "+
		"aventi sede nei suddetti Paesi, di cui non è possibile identificare il titolare effettivo ne' verificarne "+
		"l’identità.", "LR", "", false)
	pdf.MultiCell(0, 2.5, "3. (…).", "LR", "", false)
	pdf.MultiCell(0, 2.5, "4. È fatta in ogni caso salva l'applicazione dell'articolo 35, comma 2, nei "+
		"casi in cui l'operazione debba essere eseguita in quanto sussiste un obbligo di legge di ricevere "+
		"l'atto.", "LR", "", false)
	pdf.SetFont("Montserrat", "B", 6)
	pdf.MultiCell(0, 2.5, "Obblighi del cliente / sanzioni", "LR", "", false)
	pdf.SetFont("Montserrat", "", 6)
	pdf.MultiCell(0, 2.5, "Art. 22, comma 1 - I clienti forniscono per iscritto, sotto la propria "+
		"responsabilità, tutte le informazioni necessarie e aggiornate per consentire ai soggetti obbligati di "+
		"adempiere agli obblighi di adeguata verifica.", "LR", "", false)
	pdf.MultiCell(0, 2.5, "Art. 55, comma 3 - Salvo che il fatto costituisca più grave reato, chiunque "+
		"essendo obbligato, ai sensi del presente decreto, a fornire i dati e le informazioni necessarie ai fini "+
		"dell'adeguata verifica della clientela, fornisce dati falsi o informazioni non veritiere, e' punito con la "+
		"reclusione da sei mesi a tre anni e con la multa da 10.000 euro a 30.000 "+
		"euro", "LR", "", false)
	pdf.SetFont("Montserrat", "B", 6)
	pdf.MultiCell(0, 2.5, "Nozione di titolare effettivo", "LR", "", false)
	pdf.MultiCell(0, 2.5, "Art.1, comma 2, lett. pp) del D. Lgs. n.231/2007 ", "LR", "", false)
	pdf.SetFont("Montserrat", "", 6)
	pdf.MultiCell(0, 2.5, "la persona fisica o le persone fisiche, diverse dal cliente, nell'interesse "+
		"della quale  o  delle  quali,  in ultima istanza, il rapporto continuativo è istaurato, la prestazione "+
		"professionale è resa o l'operazione è eseguita.", "LR", "", false)
	pdf.SetFont("Montserrat", "B", 6)
	pdf.MultiCell(0, 2.5, "Nozione di persona politicamente esposta", "LR", "", false)
	pdf.MultiCell(0, 2.5, "Art. 1, comma 1, lettera dd) D. Lgs. 231/2007 così come modificato dal D. Lgs. 125/2019", "LR", "", false)
	pdf.SetFont("Montserrat", "", 6)
	pdf.MultiCell(0, 2.5, "Persone politicamente esposte: le persone fisiche che occupano o hanno "+
		"cessato di occupare da meno di un anno importanti cariche pubbliche, nonché i loro familiari e coloro "+
		"che con i predetti soggetti intrattengono notoriamente stretti legami, come di "+
		"seguito elencate:", "LR", "", false)
	pdf.MultiCell(0, 2.5, "1) sono persone fisiche che occupano o hanno occupato importanti cariche "+
		"pubbliche coloro che ricoprono o hanno ricoperto la carica di:", "LR", "", false)
	pdf.CellFormat(5, 5, "", "L", 0, "", false, 0, "")
	pdf.MultiCell(0, 2.5, "1.1 Presidente della Repubblica, Presidente del Consiglio, Ministro, "+
		"Vice-Ministro e Sottosegretario, Presidente di Regione, assessore regionale, Sindaco di capoluogo di "+
		"provincia o città metropolitana, Sindaco di comune con popolazione non inferiore a 15.000 abitanti "+
		"nonché cariche analoghe in Stati esteri;", "R", "", false)
	pdf.CellFormat(5, 2.5, "", "L", 0, "", false, 0, "")
	pdf.MultiCell(0, 2.5, "1.2 deputato, senatore, parlamentare europeo, consigliere regionale "+
		"nonché cariche analoghe in Stati esteri;", "R", "", false)
	pdf.CellFormat(5, 2.5, "", "L", 0, "", false, 0, "")
	pdf.MultiCell(0, 2.5, "1.3 membro degli organi direttivi centrali di partiti politici;", "R", "", false)
	pdf.CellFormat(5, 5, "", "L", 0, "", false, 0, "")
	pdf.MultiCell(0, 2.5, "1.4 giudice della Corte Costituzionale, magistrato della Corte di Cassazione "+
		"o della Corte dei conti, consigliere di Stato e altri componenti del Consiglio di Giustizia Amministrativa "+
		"per la Regione siciliana nonché cariche analoghe in Stati esteri;", "R", "", false)
	pdf.CellFormat(5, 2.5, "", "L", 0, "", false, 0, "")
	pdf.MultiCell(0, 2.5, "1.5 membro degli organi direttivi delle banche centrali e delle autorità "+
		"indipendenti;", "R", "", false)
	pdf.CellFormat(5, 2.5, "", "L", 0, "", false, 0, "")
	pdf.MultiCell(0, 2.5, "1.6 ambasciatore, incaricato d’affari ovvero cariche equivalenti in Stati "+
		"esteri, ufficiale di grado apicale delle forze armate ovvero cariche analoghe in "+
		"Stati esteri;", "R", "", false)
	pdf.CellFormat(5, 7.5, "", "L", 0, "", false, 0, "")
	pdf.MultiCell(0, 2.5, "1.7 componente degli organi di amministrazione, direzione o controllo delle "+
		"imprese controllate, anche indirettamente, dallo Stato italiano o da uno Stato estero ovvero partecipate, "+
		"in misura prevalente o totalitaria, dalle Regioni, da comuni capoluoghi di provincia e città metropolitane "+
		"e da comuni con popolazione complessivamente non inferiore a 15.000 "+
		"abitanti;", "R", "", false)
	pdf.CellFormat(5, 2.5, "", "L", 0, "", false, 0, "")
	pdf.MultiCell(0, 2.5, "1.8 direttore generale di ASL e di azienda ospedaliera, di azienda ospedaliera "+
		"universitaria e degli altri enti del servizio sanitario nazionale.", "R", "", false)
	pdf.CellFormat(5, 2.5, "", "L", 0, "", false, 0, "")
	pdf.MultiCell(0, 2.5, "1.9 direttore, vicedirettore e membro dell’organo di gestione o soggetto "+
		"svolgenti funzioni equivalenti in organizzazioni internazionali;", "R", "", false)
	pdf.MultiCell(0, 2.5, "2) sono familiari di persone politicamente esposte: i genitori, il coniuge o "+
		"la persona legata in unione civile o convivenza di fatto o istituti assimilabili alla persona politicamente "+
		"esposta, i figli e i loro coniugi nonché le persone legate ai figli in unione civile o convivenza di fatto "+
		"o istituti assimilabili;", "LR", "", false)
	pdf.MultiCell(0, 2.5, "3) sono soggetti con i quali le persone politicamente esposte intrattengono "+
		"notoriamente stretti legami:", "LR", "", false)
	pdf.CellFormat(5, 5, "", "L", 0, "", false, 0, "")
	pdf.MultiCell(0, 2.5, "3.1 le persone fisiche che ai sensi del presente decreto detengono, "+
		"congiuntamente alla persona politicamente esposta, la titolarità effettiva di enti giuridici, trust  e "+
		"istituti giuridici affini ovvero che intrattengono con la persona politicamente esposta stretti rapporti "+
		"di affari;", "R", "", false)
	pdf.CellFormat(5, 5, "", "L", 0, "", false, 0, "")
	pdf.MultiCell(0, 2.5, "3.2 le persone fisiche che detengono solo formalmente il controllo totalitario "+
		"di un’entità notoriamente costituita, di fatto, nell’interesse e a beneficio di una persona politicamente "+
		"esposta.", "R", "", false)
	pdf.MultiCell(0, 2, "", "BLR", "", false)
}

func GetWoptaInfoTable(pdf *fpdf.Fpdf) {
	DrawPinkHorizontalLine(pdf, 0.1)
	pdf.Ln(0.5)
	pdf.SetFont("Montserrat", "", 6)
	pdf.MultiCell(0, 3, "DATI DELLA PERSONA FISICA CHE ENTRA IN CONTATTO CON IL "+
		"CONTRAENTE", "", "", false)
	pdf.Ln(1)
	pdf.SetFont("Montserrat", "", 9)
	pdf.MultiCell(0, 3, "LOMAZZI MICHELE iscritto alla Sezione A del RUI con numero "+
		"A000703480 in data 02.03.2022", "", "", false)
	pdf.Ln(0.5)
	DrawPinkHorizontalLine(pdf, 0.1)
	pdf.Ln(0.5)
	pdf.SetFont("Montserrat", "", 6)
	pdf.MultiCell(0, 3, "QUALIFICA", "", "", false)
	pdf.Ln(1)
	pdf.SetFont("Montserrat", "", 9)
	pdf.MultiCell(0, 3, "Responsabile dell’attività di intermediazione assicurativa di Wopta "+
		"Assicurazioni Srl, Società iscritta alla Sezione A del RUI con numero A000701923 in data "+
		"14.02.2022", "", "", false)
	pdf.Ln(0.5)
	DrawPinkHorizontalLine(pdf, 0.1)
	pdf.Ln(0.5)
	pdf.SetFont("Montserrat", "", 6)
	pdf.MultiCell(0, 3, "SEDE LEGALE", "", "", false)
	pdf.Ln(1)
	pdf.SetFont("Montserrat", "", 9)
	pdf.MultiCell(0, 3, "Galleria del Corso, 1 – 20122 MILANO (VI)", "", "", false)
	pdf.Ln(0.5)
	DrawPinkHorizontalLine(pdf, 0.1)
	pdf.Ln(0.5)
	pdf.SetFont("Montserrat", "", 6)
	pdf.Cell(50, 3, "RECAPITI TELEFONICI")
	pdf.Cell(40, 3, "")
	pdf.MultiCell(50, 3, "E-MAIL", "", "1", false)
	pdf.Ln(1)
	pdf.SetFont("Montserrat", "", 9)
	pdf.Cell(50, 3, "02.91.24.03.46")
	pdf.Cell(40, 3, "")
	pdf.MultiCell(50, 3, "info@wopta.it", "", "1", false)
	pdf.Ln(0.5)
	DrawPinkHorizontalLine(pdf, 0.1)
	pdf.Ln(0.5)
	pdf.SetFont("Montserrat", "", 6)
	pdf.Cell(50, 3, "PEC ")
	pdf.Cell(40, 3, "")
	pdf.MultiCell(50, 3, "SITO INTERNET", "", "1", false)
	pdf.Ln(1)
	pdf.SetFont("Montserrat", "", 9)
	pdf.Cell(50, 3, "woptaassicurazioni@legalmail.it")
	pdf.Cell(40, 3, "")
	pdf.MultiCell(50, 3, "wopta.it", "", "1", false)
	pdf.Ln(0.5)
	DrawPinkHorizontalLine(pdf, 0.1)
	pdf.Ln(0.5)
	pdf.SetFont("Montserrat", "", 6)
	pdf.MultiCell(0, 3, "AUTORITÀ COMPETENTE ALLA VIGILANZA DELL’ATTIVITÀ SVOLTA",
		"", "", false)
	pdf.Ln(1)
	pdf.SetFont("Montserrat", "", 9)
	pdf.MultiCell(0, 3, "IVASS – Istituto per la Vigilanza sulle Assicurazioni - Via del Quirinale, "+
		"21 - 00187 Roma", "", "", false)
	pdf.Ln(0.5)
	DrawPinkHorizontalLine(pdf, 0.1)
}

func GetAllegato3Section(pdf *fpdf.Fpdf) {
	pdf.SetFont("Montserrat", "B", 11)
	pdf.MultiCell(0, 3, "ALLEGATO 3 - INFORMATIVA SUL DISTRIBUTORE", "", "CM", false)
	pdf.Ln(3)
	pdf.SetFont("Montserrat", "", 9)
	pdf.MultiCell(0, 3, "Il distributore ha l’obbligo di consegnare/trasmettere al contraente il presente"+
		" documento, prima della sottoscrizione della prima proposta o, qualora non prevista, del primo contratto di "+
		"assicurazione, di metterlo a disposizione del pubblico nei propri locali, anche mediante apparecchiature "+
		"tecnologiche, oppure di pubblicarlo sul proprio sito internet ove utilizzato per la promozione e collocamento "+
		"di prodotti assicurativi, dando avviso della pubblicazione nei propri locali. In occasione di rinnovo o "+
		"stipula di un nuovo contratto o di qualsiasi operazione avente ad oggetto un prodotto di investimento "+
		"assicurativo il distributore consegna o trasmette le informazioni di cui all’Allegato 3 solo in caso di "+
		"successive modifiche di rilievo delle stesse.", "", "", false)
	pdf.Ln(3)

	pdf.SetFont("Montserrat", "B", 10)
	pdf.MultiCell(0, 3, "SEZIONE I - Informazioni generali sull’intermediario che entra in contatto con "+
		"il contraente", "", "", false)
	pdf.Ln(1)

	GetWoptaInfoTable(pdf)
	pdf.Ln(1)

	SetBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, "Gli estremi identificativi e di iscrizione dell’Intermediario e dei soggetti che "+
		"operano per lo stesso possono essere verificati consultando il Registro Unico degli Intermediari assicurativi "+
		"e riassicurativi sul sito internet dell’IVASS (www.ivass.it)", "", fpdf.AlignLeft, false)
	pdf.Ln(3)
	SetBlackBoldFont(pdf, 10)
	pdf.MultiCell(0, 3, "SEZIONE II - Informazioni sull’attività svolta dall’intermediario assicurativo ",
		"", fpdf.AlignLeft, false)
	pdf.Ln(1)

	SetBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, "La Wopta Assicurazioni Srl comunica di aver messo a disposizione nei propri "+
		"locali l’elenco degli obblighi di comportamento cui adempie, come indicati nell’allegato 4-ter del Regolamento"+
		" IVASS n. 40/2018.", "", "", false)
	pdf.Ln(1)
	pdf.MultiCell(0, 3, "Si comunica che nel caso di offerta fuori sede o nel caso in cui la fase "+
		"precontrattuale si svolga mediante tecniche di comunicazione a distanza il contraente riceve l’elenco "+
		"degli obblighi.", "", "", false)
	pdf.Ln(3)
	SetBlackBoldFont(pdf, 10)
	pdf.MultiCell(0, 3, "SEZIONE III - Informazioni relative a potenziali situazioni di conflitto "+
		"d’interessi", "", "", false)
	pdf.Ln(1)
	SetBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, "Wopta Assicurazioni Srl ed i soggetti che operano per la stessa non sono "+
		"detentori di una partecipazione, diretta o indiretta, pari o superiore al 10% del capitale sociale o dei "+
		"diritti di voto di alcuna Impresa di assicurazione.", "", "", false)
	pdf.Ln(1)
	pdf.MultiCell(0, 3, "Le Imprese di assicurazione o Imprese controllanti un’Impresa di assicurazione "+
		"non sono detentrici di una partecipazione, diretta o indiretta, pari o superiore al 10% del capitale sociale "+
		"o dei diritti di voto dell’Intermediario.", "", "", false)
	pdf.Ln(3)
	SetBlackBoldFont(pdf, 10)
	pdf.MultiCell(0, 3, "SEZIONE IV - Informazioni sugli strumenti di tutela del contraente",
		"", "", false)
	pdf.Ln(1)
	SetBlackRegularFont(pdf, standardTextSize)
	pdf.SetFont("Montserrat", "", 9)
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

func GetAllegato4Section(pdf *fpdf.Fpdf) {
	pdf.SetFont("Montserrat", "B", 11)
	pdf.MultiCell(0, 3, "ALLEGATO 4 - INFORMAZIONI SULLA DISTRIBUZIONE\nDEL PRODOTTO ASSICURATIVO NON IBIP",
		"", "CM", false)
	pdf.Ln(3)
	pdf.SetFont("Montserrat", "", 9)
	pdf.MultiCell(0, 3, "Il distributore ha l’obbligo di consegnare o trasmettere al contraente, prima "+
		"della sottoscrizione di ciascuna proposta o, qualora non prevista, di ciascun contratto assicurativo, il "+
		"presente documento, che contiene notizie sul modello e l’attività di distribuzione, sulla consulenza fornita "+
		"e sulle remunerazioni percepite.", "", "", false)
	pdf.Ln(1)

	GetWoptaInfoTable(pdf)
	pdf.Ln(3)

	pdf.SetFont("Montserrat", "B", 10)
	pdf.MultiCell(0, 3, "SEZIONE I - Informazioni sul modello di distribuzione", "", "", false)
	pdf.Ln(1)
	pdf.SetFont("Montserrat", "", 9)
	pdf.MultiCell(0, 3, "Secondo quanto indicato nel modulo di proposta/polizza e documentazione "+
		"precontrattuale ricevuta, la distribuzione relativamente a questa proposta/contratto è svolta per conto "+
		"della seguente impresa di assicurazione: AXA FRANCE VIE S.A.", "", "", false)
	pdf.Ln(3)
	pdf.SetFont("Montserrat", "B", 10)
	pdf.MultiCell(0, 3, "SEZIONE II: Informazioni sull’attività di distribuzione e consulenza", "", "", false)
	pdf.Ln(1)
	pdf.SetFont("Montserrat", "", 9)
	pdf.MultiCell(0, 3, "Nello svolgimento dell’attività di distribuzione, l’intermediario non presta "+
		"attività di consulenza prima della conclusione del contratto né fornisce al contraente una raccomandazione "+
		"personalizzata ai sensi dell’art. 119-ter, comma 3, del decreto legislativo n. 209/2005 "+
		"(Codice delle Assicurazioni Private)", "", "", false)
	pdf.Ln(1)
	pdf.MultiCell(0, 3, "L'attività di distribuzione assicurativa è svolta in assenza di obblighi "+
		"contrattuali che impongano di offrire esclusivamente i contratti di una o più imprese di "+
		"assicurazioni.", "", "", false)
	pdf.Ln(3)
	pdf.SetFont("Montserrat", "B", 10)
	pdf.MultiCell(0, 3, "SEZIONE III - Informazioni relative alle remunerazioni", "", "", false)
	pdf.Ln(1)
	pdf.SetFont("Montserrat", "", 9)
	pdf.MultiCell(0, 3, "Per il prodotto intermediato, è corrisposto all’intermediario, da parte "+
		"dell’impresa di assicurazione, un compenso sotto forma di commissione inclusa nel premio "+
		"assicurativo.", "", "", false)
	pdf.Ln(1)
	pdf.MultiCell(0, 3, "L’informazione sopra resa riguarda i compensi complessivamente percepiti da tutti "+
		"gli intermediari coinvolti nella distribuzione del prodotto.", "", "", false)
	pdf.Ln(3)
	pdf.SetFont("Montserrat", "B", 10)
	pdf.MultiCell(0, 3, "SEZIONE IV – Informazioni sul pagamento dei premi", "", "", false)
	pdf.Ln(1)
	pdf.SetFont("Montserrat", "", 9)
	pdf.MultiCell(0, 3, "Relativamente a questo contratto i premi pagati dal Contraente "+
		"all’intermediario e le somme destinate ai risarcimenti o ai pagamenti dovuti dalle Imprese di Assicurazione, "+
		"se regolati per il tramite dell’intermediario costituiscono patrimonio autonomo e separato dal patrimonio "+
		"dello stesso.", "", "", false)
	pdf.Ln(3)
	pdf.SetFont("Montserrat", "B", 10)
	pdf.MultiCell(0, 3, "Indicare le modalità di pagamento ammesse ", "", "", false)
	pdf.Ln(1)
	pdf.SetFont("Montserrat", "", 9)
	pdf.MultiCell(0, 3, "Sono consentiti, nei confronti di Wopta, esclusivamente bonifico e strumenti di "+
		"pagamento elettronico, quali ad esempio, carte di credito e/o carte di debito, incluse le carte "+
		"prepagate.", "", "", false)
	pdf.Ln(3)
}

func GetAllegato4TerSection(pdf *fpdf.Fpdf) {
	pdf.SetFont("Montserrat", "B", 11)
	pdf.MultiCell(0, 3, "ALLEGATO 4 TER - ELENCO DELLE REGOLE DI COMPORTAMENTO DEL DISTRIBUTORE",
		"", "", false)
	pdf.Ln(3)
	pdf.SetFont("Montserrat", "", 9)
	pdf.MultiCell(0, 3, "Il distributore ha l’obbligo di mettere a disposizione del pubblico il "+
		"presente documento nei propri locali, anche mediante apparecchiature tecnologiche, oppure pubblicarlo su "+
		"un sito internet ove utilizzato per la promozione e il collocamento di prodotti assicurativi, dando avviso "+
		"della pubblicazione nei propri locali. Nel caso di offerta fuori sede o nel caso in cui la fase "+
		"precontrattuale si svolga mediante tecniche di comunicazione a distanza, il distributore consegna o "+
		"trasmette al contraente il presente documento prima della sottoscrizione della proposta o, qualora non "+
		"prevista, del contratto di assicurazione.", "", "", false)
	pdf.Ln(1)
	DrawPinkHorizontalLine(pdf, 0.1)
	pdf.Ln(0.5)
	pdf.SetFont("Montserrat", "", 6)
	pdf.MultiCell(0, 3, "DATI DELLA PERSONA FISICA CHE ENTRA IN CONTATTO CON IL "+
		"CONTRAENTE", "", "", false)
	pdf.Ln(1)
	pdf.SetFont("Montserrat", "", 9)
	pdf.MultiCell(0, 3, "LOMAZZI MICHELE iscritto alla Sezione A del RUI con numero "+
		"A000703480 in data 02.03.2022", "", "", false)
	pdf.Ln(0.5)
	DrawPinkHorizontalLine(pdf, 0.1)
	pdf.Ln(0.5)
	pdf.SetFont("Montserrat", "", 6)
	pdf.MultiCell(0, 3, "QUALIFICA", "", "", false)
	pdf.Ln(1)
	pdf.SetFont("Montserrat", "", 9)
	pdf.MultiCell(0, 3, "Responsabile dell’attività di intermediazione assicurativa di Wopta "+
		"Assicurazioni Srl, Società iscritta alla Sezione A del RUI con numero A000701923 in data "+
		"14.02.2022", "", "", false)
	pdf.Ln(0.5)
	DrawPinkHorizontalLine(pdf, 0.1)
	pdf.Ln(0.5)
	pdf.SetFont("Montserrat", "", 6)
	pdf.MultiCell(0, 3, "SEDE LEGALE", "", "", false)
	pdf.Ln(1)
	pdf.SetFont("Montserrat", "", 9)
	pdf.MultiCell(0, 3, "Galleria del Corso, 1 – 20122 MILANO (VI)", "", "", false)
	pdf.Ln(0.5)
	DrawPinkHorizontalLine(pdf, 0.1)
	pdf.Ln(0.5)
	pdf.SetFont("Montserrat", "", 6)
	pdf.Cell(50, 3, "RECAPITI TELEFONICI")
	pdf.Cell(40, 3, "")
	pdf.MultiCell(50, 3, "E-MAIL", "", "1", false)
	pdf.Ln(1)
	pdf.SetFont("Montserrat", "", 9)
	pdf.Cell(50, 3, "02.91.24.03.46")
	pdf.Cell(40, 3, "")
	pdf.MultiCell(50, 3, "info@wopta.it", "", "1", false)
	pdf.Ln(0.5)
	DrawPinkHorizontalLine(pdf, 0.1)
	pdf.Ln(0.5)
	pdf.SetFont("Montserrat", "", 6)
	pdf.Cell(50, 3, "PEC ")
	pdf.Cell(40, 3, "")
	pdf.MultiCell(50, 3, "SITO INTERNET", "", "1", false)
	pdf.Ln(1)
	pdf.SetFont("Montserrat", "", 9)
	pdf.Cell(50, 3, "woptaassicurazioni@legalmail.it")
	pdf.Cell(40, 3, "")
	pdf.MultiCell(50, 3, "wopta.it", "", "1", false)
	pdf.Ln(0.5)
	DrawPinkHorizontalLine(pdf, 0.1)
	pdf.Ln(0.5)
	pdf.SetFont("Montserrat", "", 6)
	pdf.MultiCell(0, 3, "AUTORITÀ COMPETENTE ALLA VIGILANZA DELL’ATTIVITÀ SVOLTA",
		"", "", false)
	pdf.Ln(1)
	pdf.SetFont("Montserrat", "", 9)
	pdf.MultiCell(0, 3, "IVASS – Istituto per la Vigilanza sulle Assicurazioni - Via del Quirinale, "+
		"21 - 00187 Roma", "", "", false)
	pdf.Ln(0.5)
	DrawPinkHorizontalLine(pdf, 0.1)
	pdf.Ln(3)
	pdf.SetFont("Montserrat", "B", 10)
	pdf.MultiCell(0, 3, "Sezione I - Regole generali per la distribuzione di prodotti assicurativi",
		"", "", false)
	pdf.Ln(1)
	pdf.SetFont("Montserrat", "", 9)
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

func GetWoptaPrivacySection(pdf *fpdf.Fpdf) {
	pdf.SetFont("Montserrat", "B", 11)
	pdf.MultiCell(0, 3, "COME RISPETTIAMO LA TUA PRIVACY", "", "CM", false)
	pdf.Ln(3)
	pdf.SetFont("Montserrat", "", 9)
	pdf.MultiCell(0, 3, "Informativa sul trattamento dei dati personali", "", "", false)
	pdf.Ln(1)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.MultiCell(0, 3, "Ai sensi del REGOLAMENTO (UE) 2016/679 "+
		"(relativo alla protezione delle persone fisiche con riguardo al trattamento dei dati personali, nonché alla "+
		"libera circolazione di tali dati) si informa l’ “Interessato” (contraente / aderente alla polizza collettiva o "+
		"convenzione / assicurato / beneficiario / loro aventi causa) di quanto segue.", "", "", false)
	pdf.Ln(3)
	pdf.SetFont("Montserrat", "B", 10)
	pdf.MultiCell(0, 3, "1. TITOLARE DEL TRATTAMENTO", "", "", false)
	pdf.Ln(1)
	pdf.SetFont("Montserrat", "", 9)
	pdf.MultiCell(0, 3, "Titolare del trattamento è Wopta Assicurazioni, con sede legale in Milano, "+
		"Galleria del Corso, 1 (di seguito “Titolare”), raggiungibile all’indirizzo e-mail: "+
		"privacy@wopta.it", "", "", false)
	pdf.Ln(3)
	pdf.SetFont("Montserrat", "B", 10)
	pdf.MultiCell(0, 3, "2. I DATI PERSONALI OGGETTO DI TRATTAMENTO, FINALITÀ E BASE "+
		"GIURIDICA", "", "", false)
	pdf.Ln(1)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.MultiCell(0, 3, "a) Finalità Contrattuali, normative, amministrative e giudiziali", "", "", false)
	pdf.Ln(1)
	pdf.SetFont("Montserrat", "", 9)
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
	pdf.SetFont("Montserrat", "B", 9)
	pdf.MultiCell(0, 3, "b) Finalità commerciali", "", "", false)
	pdf.Ln(1)
	pdf.SetFont("Montserrat", "", 9)
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
	pdf.SetFont("Montserrat", "B", 11)
	pdf.MultiCell(0, 3, "3. DESTINATARI DEI DATI PERSONALI", "", "", false)
	pdf.Ln(1)
	pdf.SetFont("Montserrat", "", 9)
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
	pdf.SetFont("Montserrat", "B", 11)
	pdf.MultiCell(0, 3, "4. TRASFERIMENTI DEI DATI PERSONALI", "", "", false)
	pdf.Ln(1)
	pdf.SetFont("Montserrat", "", 9)
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
	pdf.SetFont("Montserrat", "B", 11)
	pdf.MultiCell(0, 3, "5. CONSERVAZIONE DEI DATI PERSONALI", "", "", false)
	pdf.Ln(1)
	pdf.SetFont("Montserrat", "", 9)
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
	pdf.SetFont("Montserrat", "B", 11)
	pdf.MultiCell(0, 3, "6. DIRITTI DELL’INTERESSATO", "", "", false)
	pdf.Ln(1)
	pdf.SetFont("Montserrat", "", 9)
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
	pdf.MultiCell(0, 3, "Qualora Lei ritenga che il trattamento dei Suoi Dati personali effettuato dal "+
		"Titolare avvenga in violazione di quanto previsto dal GDPR, ha il diritto di proporre reclamo al Garante "+
		"Privacy, come previsto dall'art. 77 del GDPR stesso, o di adire le opportune sedi giudiziarie "+
		"(art. 79 del GDPR).", "", "", false)
	pdf.Ln(3)
	pdf.SetFont("Montserrat", "B", 11)
	pdf.MultiCell(0, 3, "7. CONTATTI", "", "", false)
	pdf.Ln(1)
	pdf.SetFont("Montserrat", "", 9)
	pdf.MultiCell(0, 3, "Per esercitare i diritti di cui sopra o per qualunque altra richiesta può "+
		"scrivere al Titolare del trattamento all’indirizzo: privacy@wopta.it.", "", "", false)
	pdf.Ln(3)
}

func GetPersonalDataHandlingSection(pdf *fpdf.Fpdf) {
	pdf.SetTextColor(229, 0, 117)
	pdf.SetFont("Montserrat", "B", 10)
	pdf.MultiCell(0, 3, "Consenso per finalità commerciali.", "", "", false)
	pdf.Ln(1)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Montserrat", "", 9)
	pdf.MultiCell(0, 3, "Il sottoscritto, letta e compresa l’informativa sul trattamento dei dati personali",
		"", "", false)
	pdf.Ln(1)
	pdf.SetDrawColor(0, 0, 0)
	pdf.Cell(5, 3, "")
	pdf.CellFormat(3, 3, "X", "1", 0, "CM", false, 0, "")
	pdf.CellFormat(20, 3, "ACCONSENTE", "", 0, "", false, 0, "")
	pdf.Cell(20, 3, "")
	pdf.CellFormat(3, 3, "", "1", 0, "CM", false, 0, "")
	pdf.CellFormat(20, 3, "NON ACCONSENTE", "", 1, "", false, 0, "")
	pdf.Ln(1)
	pdf.MultiCell(0, 3, "al trattamento dei propri dati personali da parte di Wopta Assicurazioni per "+
		"l’invio di comunicazioni e proposte commerciali e di marketing, incluso l’invio di newsletter e ricerche di "+
		"mercato, attraverso strumenti automatizzati (sms, mms, e-mail, ecc.) e non (posta cartacea e telefono "+
		"con operatore).", "", "", false)
	pdf.Ln(3)
	pdf.Cell(0, 3, "03/04/2023")
	pdf.Ln(3)
	pdf.SetX(-80)
	pdf.Cell(0, 3, "Firma del Contraente/Assicurato")
	pdf.Ln(15)
	pdf.SetDrawColor(0, 0, 0)
	pdf.SetLineWidth(0.4)
	pdf.Line(130, pdf.GetY(), 190, pdf.GetY())
}
