package document

import (
	"github.com/dustin/go-humanize"
	"github.com/go-pdf/fpdf"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"strconv"
	"strings"
	"time"
)

var (
	signatureID int
)

func Life(pdf *fpdf.Fpdf, policy models.Policy) (string, []byte) {
	signatureID = 0

	mainHeader(pdf, policy)
	mainFooter(pdf)

	pdf.AddPage()

	insuredInfoSection(pdf, policy)

	guaranteesTable(pdf, policy)

	avvertenzeBeneficiariSection(pdf)

	beneficiariesSection(pdf, policy)

	beneficiaryReferenceSection(pdf, policy)

	surveysSection(pdf, policy)

	statementsSection(pdf, policy)

	pdf.AddPage()

	offerResumeSection(pdf, policy)

	paymentResumeSection(pdf, policy)

	contractWithdrawlSection(pdf)

	paymentMethodSection(pdf)

	emitResumeSection(pdf, policy)

	woptaAxaCompanyDescriptionSection(pdf)

	axaHeader(pdf)

	pdf.AddPage()

	axaFooter(pdf)

	axaDeclarationsConsentSection(pdf, policy)

	pdf.AddPage()

	axaTableSection(pdf, policy)

	pdf.AddPage()

	GetAxaTablePart2Section(pdf, policy)

	pdf.Ln(15)

	GetAxaTablePart3Section(pdf)

	woptaHeader(pdf)

	pdf.AddPage()

	woptaFooter(pdf)

	GetAllegato3Section(pdf)

	pdf.AddPage()

	GetAllegato4Section(pdf)

	pdf.AddPage()

	GetAllegato4TerSection(pdf)

	pdf.AddPage()

	GetWoptaPrivacySection(pdf)

	GetPersonalDataHandlingSection(pdf, policy)

	filename, out := save(pdf, policy)
	return filename, out
}

func mainHeader(pdf *fpdf.Fpdf, policy models.Policy) {
	var (
		opt        fpdf.ImageOptions
		logoPath   string
		cfpi       string
		expiryInfo string
	)
	logoPath = lib.GetAssetPathByEnv(basePath) + "/logo_vita.png"

	contractor := policy.Contractor
	address := strings.ToUpper(contractor.Residence.StreetName+", "+contractor.Residence.StreetNumber+"\n"+
		contractor.Residence.PostalCode+" "+contractor.Residence.City+" ("+contractor.Residence.CityCode) + "\n"

	if contractor.VatCode == "" {
		cfpi = contractor.FiscalCode
	} else {
		cfpi = contractor.VatCode
	}

	if policy.PaymentSplit == "monthly" {
		expiryInfo = "Prima scandenza mensile il: " +
			policy.StartDate.AddDate(0, 1, 0).Format(dateLayout) + "\n"
	} else if policy.PaymentSplit == "yearly" {
		expiryInfo = "Prima scadenza annuale il: " +
			policy.StartDate.AddDate(1, 0, 0).Format(dateLayout) + "\n"
	}

	policyInfo := "Numero: " + policy.CodeCompany + "\n" +
		"Decorre dal: " + policy.StartDate.Format(dateLayout) + " ore 24:00\n" +
		"Scade il: " + policy.EndDate.Format(dateLayout) + " ore 24:00\n" +
		expiryInfo + "Non si rinnova a scadenza."

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
		pdf.Cell(10, 6, "Vita")
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
		pdf.Ln(8)
	})
}

func mainFooter(pdf *fpdf.Fpdf) {
	var opt fpdf.ImageOptions

	pdf.SetFooterFunc(func() {
		pdf.SetXY(10, -15)
		setPinkRegularFont(pdf, smallTextSize)
		pdf.MultiCell(0, 3, "Wopta per te. Vita è un prodotto assicurativo di AXA France Vie S.A. – Rappresentanza Generale per l’Italia\ndistribuito da Wopta Assicurazioni S.r.l.", "", "", false)
		opt.ImageType = "png"
		pdf.ImageOptions(lib.GetAssetPathByEnv(basePath)+"/axa/logo.png", 190, 281, 0, 8, false, opt, 0, "")
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
		pdf.SetY(-7)
		pageNumber(pdf)
	})
}

func insuredInfoSection(pdf *fpdf.Fpdf, policy models.Policy) {
	getParagraphTitle(pdf, "La tua assicurazione è operante per il seguente Assicurato e Garanzie")
	pdf.Ln(8)
	insuredInfoTable(pdf, policy.Assets[0].Person)
}

func insuredInfoTable(pdf *fpdf.Fpdf, insured *models.User) {
	residenceAddress := strings.ToUpper(insured.Residence.StreetName + ", " + insured.Residence.StreetNumber +
		" - " + insured.Residence.PostalCode + " " + insured.Residence.City + " (" + insured.Residence.CityCode + ")")

	drawPinkHorizontalLine(pdf, thickLineWidth)
	pdf.Ln(2)
	setBlackBoldFont(pdf, standardTextSize)
	pdf.Cell(20, 2, "Cognome e Nome")
	setBlackRegularFont(pdf, standardTextSize)
	pdf.SetX(pdf.GetX() + 24)
	pdf.Cell(20, 2, strings.ToUpper(insured.Surname+" "+insured.Name))
	pdf.SetX(pdf.GetX() + 60)
	setBlackBoldFont(pdf, standardTextSize)
	pdf.Cell(10, 2, "Codice fiscale:")
	pdf.SetX(pdf.GetX() + 20)
	setBlackRegularFont(pdf, standardTextSize)
	pdf.Cell(20, 2, strings.ToUpper(insured.FiscalCode))
	pdf.Ln(2.5)
	drawPinkHorizontalLine(pdf, thinLineWidth)
	pdf.Ln(2)
	setBlackBoldFont(pdf, standardTextSize)
	pdf.Cell(20, 2, "Residente in")
	setBlackRegularFont(pdf, standardTextSize)
	pdf.SetX(pdf.GetX() + 24)
	pdf.Cell(20, 2, residenceAddress)
	pdf.Ln(2.5)
	drawPinkHorizontalLine(pdf, thinLineWidth)
	pdf.Ln(2)
	setBlackBoldFont(pdf, standardTextSize)
	pdf.Cell(20, 2, "Mail")
	setBlackRegularFont(pdf, standardTextSize)
	pdf.SetX(pdf.GetX() + 24)
	pdf.Cell(20, 2, insured.Mail)
	pdf.SetX(pdf.GetX() + 60)
	setBlackBoldFont(pdf, standardTextSize)
	pdf.Cell(10, 2, "Telefono:")
	pdf.SetX(pdf.GetX() + 20)
	setBlackRegularFont(pdf, 9)
	pdf.Cell(20, 2, insured.Phone)
	pdf.Ln(3)
	drawPinkHorizontalLine(pdf, thickLineWidth)
	pdf.Ln(1)
}

func guaranteesTable(pdf *fpdf.Fpdf, policy models.Policy) {
	const (
		death               = "death"
		permanentDisability = "permanent-disability"
		temporaryDisability = "temporary-disability"
		seriousIll          = "serious-ill"
	)
	var (
		price float64
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
			"price":                      "==== (*)",
		},
		temporaryDisability: {
			"name":                       "Inabilità Temporanea da Infortunio o Malattia",
			"sumInsuredLimitOfIndemnity": "======",
			"duration":                   "==",
			"endDate":                    "====",
			"price":                      "==== (*)",
		},
		seriousIll: {
			"name":                       "Malattie Gravi",
			"sumInsuredLimitOfIndemnity": "======",
			"duration":                   "==",
			"endDate":                    "====",
			"price":                      "==== (*)",
		},
	}

	for _, guarantee := range policy.GuaranteesToMap() {
		guarantees[guarantee.Slug]["sumInsuredLimitOfIndemnity"] = humanize.FormatFloat("#.###,", guarantee.Value.SumInsuredLimitOfIndemnity) + " €"
		guarantees[guarantee.Slug]["duration"] = strconv.Itoa(guarantee.Value.Duration.Year)
		guarantees[guarantee.Slug]["endDate"] = policy.StartDate.AddDate(guarantee.Value.Duration.Year, 0, 0).Format(dateLayout)
		if policy.PaymentSplit == "monthly" {
			price = guarantee.Value.PremiumGrossMonthly * 12
		} else {
			price = guarantee.Value.PremiumGrossYearly
		}
		guarantees[guarantee.Slug]["price"] = humanize.FormatFloat("#.###,##", price) + " €"
		if guarantee.Slug != death {
			guarantees[guarantee.Slug]["price"] += " (*)"
		}
	}

	setBlackBoldFont(pdf, standardTextSize)
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
	drawPinkHorizontalLine(pdf, thinLineWidth)

	for _, slug := range slugs {
		setBlackBoldFont(pdf, standardTextSize)
		pdf.CellFormat(90, 6, guarantees[slug]["name"], "", 0, "", false, 0, "")
		setBlackRegularFont(pdf, standardTextSize)
		pdf.CellFormat(25, 6, guarantees[slug]["sumInsuredLimitOfIndemnity"],
			"", 0, "RM", false, 0, "")
		pdf.CellFormat(25, 6, guarantees[slug]["duration"], "", 0, "CM",
			false, 0, "")
		pdf.CellFormat(25, 6, guarantees[slug]["endDate"], "", 0, "CM", false, 0, "")
		pdf.CellFormat(0, 6, guarantees[slug]["price"], "",
			0, "RM", false, 0, "")
		pdf.Ln(5)
		drawPinkHorizontalLine(pdf, thinLineWidth)
	}
	pdf.Ln(0.5)
	setBlackRegularFont(pdf, smallTextSize)
	pdf.Cell(80, 3, "(*) imposte assicurative di legge incluse nella misura del 2,50% del premio imponibile")
	pdf.Ln(3)
}

func avvertenzeBeneficiariSection(pdf *fpdf.Fpdf) {
	getParagraphTitle(pdf, "Nomina dei Beneficiari e Referente terzo, per il caso di garanzia Decesso "+
		"(qualora sottoscritta)")
	pdf.Ln(8)
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, "AVVERTENZE: Può scegliere se designare nominativamente i beneficiari o se "+
		"designare genericamente come beneficiari i suoi eredi legittimi e/o testamentari. In caso di mancata "+
		"designazione nominativa, la Compagnia potrà incontrare, al decesso dell’Assicurato, maggiori difficoltà "+
		"nell’identificazione e nella ricerca dei beneficiari. La modifica o revoca del/i beneficiario/i deve essere "+
		"comunicata alla Compagnia in forma scritta.\nIn caso di specifiche esigenze di riservatezza, la Compagnia "+
		"potrà rivolgersi ad un soggetto terzo (diverso dal Beneficiario) in caso di Decesso al fine di contattare "+
		"il Beneficiario designato.", "", "", false)
}

func beneficiariesSection(pdf *fpdf.Fpdf, policy models.Policy) {
	legitimateSuccessorsChoice := "X"
	designatedSuccessorsChoice := ""
	beneficiaries := [2]map[string]string{
		{
			"name":     "=====",
			"fiscCode": "=====",
			"address":  "=====",
			"mail":     "=====",
			"phone":    "=====",
			"relation": "=====",
			"consent":  "=====",
		},
		{
			"name":           "=====",
			"fiscCode":       "=====",
			"address":        "=====",
			"mail":           "=====",
			"phone":          "=====",
			"relation":       "=====",
			"contactConsent": "=====",
		},
	}

	deathGuarantee, err := policy.ExtractGuarantee("death")
	lib.CheckError(err)

	if deathGuarantee.Beneficiaries != nil && !(*deathGuarantee.Beneficiaries)[0].IsLegitimateSuccessors {
		legitimateSuccessorsChoice = ""
		designatedSuccessorsChoice = "X"

		for index, beneficiary := range *deathGuarantee.Beneficiaries {
			address := strings.ToUpper(beneficiary.Residence.StreetName + ", " + beneficiary.Residence.StreetNumber +
				" - " + beneficiary.Residence.PostalCode + " " + beneficiary.Residence.City +
				" (" + beneficiary.Residence.CityCode + ")")
			beneficiaries[index]["name"] = strings.ToUpper(beneficiary.Surname + " " + beneficiary.Name)
			beneficiaries[index]["fiscCode"] = strings.ToUpper(beneficiary.FiscalCode)
			beneficiaries[index]["address"] = address
			beneficiaries[index]["mail"] = beneficiary.Mail
			beneficiaries[index]["phone"] = beneficiary.Phone
			if beneficiary.IsFamilyMember {
				beneficiaries[index]["relation"] = "Nucleo familiare (rapporto di parentela, coniuge, unione civile, " +
					"convivenza more uxorio)"
			} else {
				beneficiaries[index]["relation"] = "Altro (no rapporto parentela)"
			}
			if beneficiary.IsContactable {
				beneficiaries[index]["contactConsent"] = "SI"
			} else {
				beneficiaries[index]["contactConsent"] = "NO"
			}
		}
	}

	getParagraphTitle(pdf, "Beneficiario")
	pdf.Ln(8)
	setBlackRegularFont(pdf, standardTextSize)
	pdf.CellFormat(0, 3, "Io sottoscritto Assicurato, con la sottoscrizione della presente polizza, in "+
		"riferimento alla garanzia Decesso:", "", 0, "", false, 0, "")
	pdf.Ln(4)
	setBlackDrawColor(pdf)
	pdf.SetX(11)
	pdf.CellFormat(3, 3, legitimateSuccessorsChoice, "1", 0, "CM", false, 0, "")
	pdf.CellFormat(0, 3, "Designo genericamente quali beneficiari della prestazione i miei eredi "+
		"(legittimi e/o testamentari)", "", 0, "", false, 0, "")
	pdf.Ln(4)
	pdf.SetX(11)
	pdf.CellFormat(3, 3, designatedSuccessorsChoice, "1", 0, "CM", false, 0, "")
	pdf.CellFormat(0, 3, "Designo nominativamente il/i seguente/i soggetto/i quale beneficiario/i della "+
		"prestazione", "", 0, "", false, 0, "")
	pdf.Ln(5)
	beneficiariesTable(pdf, beneficiaries)
}

func beneficiariesTable(pdf *fpdf.Fpdf, beneficiaries [2]map[string]string) {
	for _, beneficiary := range beneficiaries {
		drawPinkHorizontalLine(pdf, thickLineWidth)
		pdf.Ln(1.5)
		setBlackBoldFont(pdf, standardTextSize)
		pdf.Cell(50, 2, "Cognome e nome")
		setBlackRegularFont(pdf, standardTextSize)
		pdf.Cell(20, 2, beneficiary["name"])
		pdf.SetX(pdf.GetX() + 60)
		setBlackBoldFont(pdf, standardTextSize)
		pdf.Cell(20, 2, "Cod. Fisc.: ")
		setBlackRegularFont(pdf, standardTextSize)
		pdf.Cell(20, 2, beneficiary["fiscCode"])
		pdf.Ln(3)
		drawPinkHorizontalLine(pdf, thinLineWidth)
		pdf.Ln(2)
		setBlackBoldFont(pdf, standardTextSize)
		pdf.Cell(50, 2, "Indirizzo")
		setBlackRegularFont(pdf, standardTextSize)
		pdf.Cell(20, 2, beneficiary["address"])
		pdf.SetX(pdf.GetX() + 60)
		setBlackBoldFont(pdf, standardTextSize)
		pdf.Ln(3)
		drawPinkHorizontalLine(pdf, thinLineWidth)
		pdf.Ln(2)
		setBlackBoldFont(pdf, standardTextSize)
		pdf.Cell(50, 2, "Mail")
		setBlackRegularFont(pdf, standardTextSize)
		pdf.Cell(20, 2, beneficiary["mail"])
		pdf.SetX(pdf.GetX() + 60)
		setBlackBoldFont(pdf, standardTextSize)
		pdf.Cell(20, 2, "Telefono: ")
		setBlackRegularFont(pdf, standardTextSize)
		pdf.Cell(20, 2, beneficiary["phone"])
		pdf.Ln(3)
		drawPinkHorizontalLine(pdf, thinLineWidth)
		pdf.Ln(2)
		setBlackBoldFont(pdf, standardTextSize)
		pdf.Cell(50, 2, "Relazione con Assicurato")
		setBlackRegularFont(pdf, standardTextSize)
		pdf.Cell(20, 2, beneficiary["relation"])
		pdf.Ln(3)
		drawPinkHorizontalLine(pdf, thinLineWidth)
		pdf.Ln(2)
		pdf.Cell(165, 2, "Consenso ad invio comunicazioni da parte della Compagnia al beneficiario, prima "+
			"dell'evento Decesso:")
		pdf.Cell(20, 2, beneficiary["contactConsent"])
		pdf.Ln(3)
		drawPinkHorizontalLine(pdf, thinLineWidth)
		pdf.Ln(2)
	}
}

func beneficiaryReferenceSection(pdf *fpdf.Fpdf, policy models.Policy) {
	beneficiaryReference := map[string]string{
		"name":     "=====",
		"fiscCode": "=====",
		"address":  "=====",
		"mail":     "=====",
		"phone":    "=====",
	}

	deathGuarantee, err := policy.ExtractGuarantee("death")
	lib.CheckError(err)

	if deathGuarantee.BeneficiaryReferance != nil {
		beneficiary := deathGuarantee.BeneficiaryReferance
		address := strings.ToUpper(beneficiary.Residence.StreetName + ", " + beneficiary.Residence.StreetNumber +
			" - " + beneficiary.Residence.PostalCode + " " + beneficiary.Residence.City +
			" (" + beneficiary.Residence.CityCode + ")")
		beneficiaryReference["name"] = strings.ToUpper(beneficiary.Surname + " " + beneficiary.Name)
		beneficiaryReference["fiscCode"] = strings.ToUpper(beneficiary.FiscalCode)
		beneficiaryReference["address"] = address
		beneficiaryReference["mail"] = beneficiary.Mail
		beneficiaryReference["phone"] = beneficiary.Phone
	}

	getParagraphTitle(pdf, "Referente terzo")
	pdf.Ln(8)
	beneficiaryReferenceTable(pdf, beneficiaryReference)
	pdf.Ln(2)
}

func beneficiaryReferenceTable(pdf *fpdf.Fpdf, beneficiaryReference map[string]string) {
	drawPinkHorizontalLine(pdf, thickLineWidth)
	pdf.Ln(1.5)
	setBlackBoldFont(pdf, standardTextSize)
	pdf.Cell(50, 2, "Cognome e nome")
	setBlackRegularFont(pdf, standardTextSize)
	pdf.Cell(20, 2, beneficiaryReference["name"])
	pdf.SetX(pdf.GetX() + 60)
	setBlackBoldFont(pdf, standardTextSize)
	pdf.Cell(20, 2, "Cod. Fisc.: ")
	setBlackRegularFont(pdf, standardTextSize)
	pdf.Cell(20, 2, beneficiaryReference["fiscCode"])
	pdf.Ln(3)
	drawPinkHorizontalLine(pdf, thinLineWidth)
	pdf.Ln(2)
	setBlackBoldFont(pdf, standardTextSize)
	pdf.Cell(50, 2, "Indirizzo")
	setBlackRegularFont(pdf, standardTextSize)
	pdf.Cell(20, 2, beneficiaryReference["address"])
	pdf.SetX(pdf.GetX() + 60)
	setBlackBoldFont(pdf, standardTextSize)
	pdf.Ln(3)
	drawPinkHorizontalLine(pdf, thinLineWidth)
	pdf.Ln(2)
	setBlackBoldFont(pdf, standardTextSize)
	pdf.Cell(50, 2, "Mail")
	setBlackRegularFont(pdf, standardTextSize)
	pdf.Cell(20, 2, beneficiaryReference["mail"])
	pdf.SetX(pdf.GetX() + 60)
	setBlackBoldFont(pdf, standardTextSize)
	pdf.Cell(20, 2, "Telefono: ")
	setBlackRegularFont(pdf, standardTextSize)
	pdf.Cell(20, 2, beneficiaryReference["phone"])
	pdf.Ln(3)
	drawPinkHorizontalLine(pdf, thinLineWidth)
}

func surveysSection(pdf *fpdf.Fpdf, policy models.Policy) {
	surveys := *policy.Surveys

	getParagraphTitle(pdf, "Dichiarazioni da leggere con attenzione prima di firmare")
	pdf.Ln(8)
	printSurvey(pdf, surveys[0])
	pdf.Ln(5)
	getParagraphTitle(pdf, "Questionario Medico")
	pdf.Ln(8)
	for _, survey := range surveys[1:] {
		err := printSurvey(pdf, survey)
		lib.CheckError(err)
	}
	pdf.Ln(8)
	drawSignatureForm(pdf)
	pdf.Ln(5)
}

func statementsSection(pdf *fpdf.Fpdf, policy models.Policy) {
	statements := *policy.Statements
	pdf.Ln(8)
	for _, statement := range statements {
		printStatement(pdf, statement)
	}
	pdf.SetY(pdf.GetY() - 28)
	setBlackBoldFont(pdf, standardTextSize)
	pdf.MultiCell(70, 3, "AXA France Vie\n(Rappresentanza Generale per l'Italia)", "",
		fpdf.AlignCenter, false)
	var opt fpdf.ImageOptions
	opt.ImageType = "png"
	pdf.ImageOptions(lib.GetAssetPathByEnv(basePath)+"/firma_axa.png", 35, pdf.GetY()+3, 30, 8,
		false, opt, 0, "")
	pdf.Ln(15)
}

func offerResumeSection(pdf *fpdf.Fpdf, policy models.Policy) {
	var (
		paymentSplit string
		tableInfo    [][]string
	)

	switch policy.PaymentSplit {
	case "monthly":
		paymentSplit = "MENSILE"
		tableInfo = [][]string{
			{
				"Mensile firma del contratto",
				lib.HumanaizePriceEuro(policy.OffersPrices["default"]["monthly"].Net),
				lib.HumanaizePriceEuro(policy.OffersPrices["default"]["monthly"].Tax),
				lib.HumanaizePriceEuro(policy.OffersPrices["default"]["monthly"].Gross),
			},
			{
				"Pari ad un premio Annuale",
				lib.HumanaizePriceEuro(policy.OffersPrices["default"]["monthly"].Net * 12),
				lib.HumanaizePriceEuro(policy.OffersPrices["default"]["monthly"].Tax * 12),
				lib.HumanaizePriceEuro(policy.OffersPrices["default"]["monthly"].Gross * 12),
			},
		}
	case "yearly":
		paymentSplit = "ANNUALE"
		tableInfo = [][]string{
			{
				"Annuale firma del contratto",
				lib.HumanaizePriceEuro(policy.OffersPrices["default"]["yearly"].Net),
				lib.HumanaizePriceEuro(policy.OffersPrices["default"]["yearly"].Tax),
				lib.HumanaizePriceEuro(policy.OffersPrices["default"]["yearly"].Gross),
			},
		}
	}

	getParagraphTitle(pdf, "Il premio per tutte le coperture assicurative attivate sulla polizza – Frazionamento: "+paymentSplit)
	pdf.Ln(8)
	setBlackRegularFont(pdf, standardTextSize)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(40, 2, "Premio", "", 0, "", false, 0, "")
	pdf.SetX(pdf.GetX() + 20)
	pdf.CellFormat(40, 2, "Imponibile", "", 0, "", false, 0, "")
	pdf.SetX(pdf.GetX() + 15)
	pdf.CellFormat(40, 2, "Imposte Assicurative", "", 0, "", false, 0, "")
	pdf.SetX(pdf.GetX() + 15)
	pdf.CellFormat(40, 2, "Totale", "", 0, "", false, 0, "")
	pdf.Ln(3)
	drawPinkHorizontalLine(pdf, thinLineWidth)
	pdf.Ln(1)
	for _, info := range tableInfo {
		pdf.CellFormat(40, 2, info[0], "", 0, "", false, 0,
			"")
		pdf.SetX(pdf.GetX() + 8)
		pdf.CellFormat(40, 2, info[1], "", 0,
			"CM", false, 0, "")
		pdf.SetX(pdf.GetX() + 20)
		pdf.CellFormat(40, 2, info[2], "", 0,
			"CM", false, 0, "")
		pdf.SetX(pdf.GetX() + 18)
		pdf.CellFormat(20, 2, info[3], "",
			0, "CM", false, 0, "")
		pdf.Ln(3)
		drawPinkHorizontalLine(pdf, thinLineWidth)
		pdf.Ln(1)
	}

}

func paymentResumeSection(pdf *fpdf.Fpdf, policy models.Policy) {
	payments := make([]float64, 20)
	var paymentSplit string
	policyStartDate := policy.StartDate

	cellWidth := pdf.GetStringWidth("00/00/0000:") + pdf.GetStringWidth("€ ###.###,##")

	if policy.PaymentSplit == "yearly" {
		paymentSplit = "ANNUALE"
		for _, guarantee := range policy.Assets[0].Guarantees {
			for i := 0; i < guarantee.Value.Duration.Year; i++ {
				payments[i] += guarantee.Value.PremiumGrossYearly
			}
		}
	} else if policy.PaymentSplit == "monthly" {
		paymentSplit = "MENSILE"
		for _, guarantee := range policy.Assets[0].Guarantees {
			for i := 0; i < guarantee.Value.Duration.Year; i++ {
				for y := 0; y < 12; y++ {
					payments[i] += guarantee.Value.PremiumGrossMonthly
				}
			}
		}
	}

	getParagraphTitle(pdf, "Pagamento dei premi successivi al primo")
	pdf.Ln(8)
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, "Il Contraente è tenuto a pagare i Premi entro 30 giorni dalle relative scadenze. "+
		"In caso di mancato pagamento del premio entro 30 giorni dalla scadenza (c.d. termine di tolleranza) "+
		"l’assicurazione è sospesa. Il contratto è risolto automaticamente in caso di mancato pagamento "+
		"del Premio entro 90 giorni dalla scadenza.", "", "", false)
	pdf.Ln(3)
	drawPinkHorizontalLine(pdf, 0.4)
	pdf.Ln(1)
	setBlackBoldFont(pdf, standardTextSize)
	pdf.CellFormat(pdf.GetStringWidth("Tipologia di premio"), 3, "Tipologia di premio:", "",
		0, "", false, 0, "")
	pdf.SetX(pdf.GetX() + 5)
	setBlackRegularFont(pdf, standardTextSize)
	setBlackDrawColor(pdf)
	pdf.CellFormat(3, 3, "", "1", 0, "", false, 0, "")
	pdf.SetX(pdf.GetX() + 1)
	pdf.CellFormat(pdf.GetStringWidth("naturale variabile annualmente"), 3, "naturale variabile annualmente",
		"", 0, "", false, 0, "")
	pdf.SetX(pdf.GetX() + 5)
	pdf.CellFormat(3, 3, "X", "1", 0, "CM", false, 0, "")
	pdf.SetX(pdf.GetX() + 1)
	pdf.CellFormat(pdf.GetStringWidth("fisso"), 3, "fisso", "", 0, "", false, 0,
		"")
	pdf.SetX(pdf.GetX() + 40)
	setBlackBoldFont(pdf, standardTextSize)
	pdf.CellFormat(pdf.GetStringWidth("Frazionamento"), 3, "Frazionamento:", "", 0, "",
		false, 0, "")
	pdf.SetX(pdf.GetX() + 3)
	pdf.CellFormat(pdf.GetStringWidth(paymentSplit), 3, paymentSplit, "", 0, "", false,
		0, "")
	pdf.Ln(4)
	drawPinkHorizontalLine(pdf, 0.1)
	pdf.Ln(1)
	setBlackRegularFont(pdf, standardTextSize)
	pdf.Cell(0, 3, "Il Premio è dovuto alle diverse annualità di Polizza, alle date qui sotto indicate:")
	pdf.Ln(4)
	drawPinkHorizontalLine(pdf, 0.1)

	for x := 0; x < len(payments)/4; x++ {
		pdf.Ln(1)
		for y := 0; y < 4; y++ {
			pdf.SetX(pdf.GetX() + 4)

			if payments[x+(5*y)] != 0 {
				var date string
				if x == 0 && y == 0 {
					date = "Alla firma:"
				} else {
					date = policyStartDate.AddDate(x+(5*y), 0, 0).Format(dateLayout) + ":"
				}
				price := lib.HumanaizePriceEuro(payments[x+(5*y)])

				pdf.CellFormat(cellWidth, 3, date+" "+price, "", 0, fpdf.AlignRight, false,
					0, "")

			} else {
				pdf.CellFormat(cellWidth, 3, "===========", "", 0, fpdf.AlignRight, false,
					0, "")
			}

		}
		pdf.Ln(4)
		drawPinkHorizontalLine(pdf, 0.1)
	}

	pdf.Ln(1)
	setBlackRegularFont(pdf, smallTextSize)
	pdf.MultiCell(0, 3, "In caso di frazionamento mensile i Premi sopra riportati sono dovuti, alle date "+
		"indicate e con successiva frequenza mensile, in misura di 1/12 per ogni mensilità. Non sono previsti oneri "+
		"o interessi di frazionamento.", "", "", false)
}

func contractWithdrawlSection(pdf *fpdf.Fpdf) {
	getParagraphTitle(pdf, "Informativa sul diritto di recesso")
	pdf.Ln(8)
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
}

func paymentMethodSection(pdf *fpdf.Fpdf) {
	getParagraphTitle(pdf, "Come puoi pagare il premio")
	pdf.Ln(8)
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, "I mezzi di pagamento consentiti, nei confronti di Wopta, sono esclusivamente "+
		"bonifico e strumenti di pagamento elettronico, quali ad esempio, carte di credito e/o carte di debito, "+
		"incluse le carte prepagate. Oppure può essere pagato direttamente alla Compagnia alla "+
		"stipula del contratto, via bonifico o carta di credito.", "", "", false)
}

func emitResumeSection(pdf *fpdf.Fpdf, policy models.Policy) {
	var offerPrice string
	emitDate := policy.EmitDate.Format(dateLayout)
	startDate := policy.StartDate.Format(dateLayout)
	if policy.OffersPrices["default"]["monthly"] != nil {
		offerPrice = humanize.FormatFloat("#.###,##", policy.OffersPrices["default"]["monthly"].Gross*12)
	} else {
		offerPrice = humanize.FormatFloat("#.###,##", policy.OffersPrices["default"]["yearly"].Gross)
	}
	text := "Polizza emessa a Milano il " + emitDate + " per un importo di € " + offerPrice + " quale " +
		"prima rata alla firma, il cui pagamento a saldo è da effettuarsi con i metodi di pagamento sopra indicati. " +
		"Wopta conferma avvenuto incasso e copertura della polizza dal " + startDate + "."
	getParagraphTitle(pdf, "Emissione polizza e pagamento della prima rata")
	pdf.Ln(8)
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, text, "", "", false)
}

func woptaAxaCompanyDescriptionSection(pdf *fpdf.Fpdf) {
	getParagraphTitle(pdf, "Chi siamo")
	pdf.Ln(8)
	setBlackRegularFont(pdf, standardTextSize)
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

func axaDeclarationsConsentSection(pdf *fpdf.Fpdf, policy models.Policy) {
	setBlackBoldFont(pdf, standardTextSize)
	pdf.Cell(0, 3, "DICHIARAZIONI E CONSENSI")
	pdf.Ln(3)
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, "Io sottoscritto, dopo aver letto l’Informativa Privacy della compagnia titolare "+
		"del trattamento redatta ai sensi del Regolamento (UE) 2016/679 (relativo alla protezione delle persone "+
		"fisiche con riguardo al trattamento dei dati personali), della quale confermo ricezione, PRESTO IL CONSENSO "+
		"al trattamento dei miei dati personali, ivi inclusi quelli eventualmente da me conferiti in riferimento al "+
		"mio stato di salute, per le finalità indicate nell’informativa, nonché alla loro comunicazione, per "+
		"successivo trattamento, da parte dei soggetti indicati nella informativa predetta.", "", "", false)
	pdf.Ln(3)
	setBlackBoldFont(pdf, standardTextSize)
	pdf.Cell(0, 3, "Resta inteso che in caso di negazione del consenso non sarà possibile "+
		"finalizzare il rapporto contrattuale assicurativo.")
	pdf.Ln(3)
	setBlackDrawColor(pdf)
	drawBlackHorizontalLine(pdf, thinLineWidth)
	pdf.Ln(5)
	pdf.Cell(0, 3, policy.EmitDate.Format(dateLayout))
	drawSignatureForm(pdf)
}

func axaTableSection(pdf *fpdf.Fpdf, policy models.Policy) {
	contractor := policy.Contractor

	identityDocumentInfo := map[string]string{
		"type":             "=====",
		"number":           "=====",
		"issuingAuthority": "=====",
		"dateOfIssue":      "=====",
		"placeOfIssue":     "=====",
		"expiryDate":       "=====",
	}
	identityDocument := contractor.GetIdentityDocument()
	if identityDocument != nil {
		identityDocumentInfo["type"] = identityDocument.Type
		identityDocumentInfo["number"] = identityDocument.Number
		identityDocumentInfo["issuingAuthority"] = identityDocument.IssuingAuthority
		identityDocumentInfo["dateOfIssue"] = identityDocument.DateOfIssue.Format(dateLayout)
		identityDocumentInfo["placeOfIssue"] = identityDocument.PlaceOfIssue
		identityDocumentInfo["expiryDate"] = identityDocument.ExpiryDate.Format(dateLayout)
	}

	insured := policy.Assets[0].Person
	domicileCity := ""
	domicileAddress := ""
	if insured.Domicile != nil {
		domicileCity = strings.ToUpper(insured.Domicile.City + " (" + insured.Domicile.CityCode + ")")
		domicileAddress = strings.ToUpper(insured.Domicile.StreetName + " " + insured.Domicile.StreetNumber)
	}

	birthDate, err := time.Parse(time.RFC3339, insured.BirthDate)
	lib.CheckError(err)

	setWhiteBoldFont(pdf, 12)
	pdf.SetFillColor(229, 0, 117)
	pdf.MultiCell(0, 6, "MODULO PER L’IDENTIFICAZIONE E L’ADEGUATA VERIFICA DELLA CLIENTELA", "LTR", "CM", true)
	setWhiteBoldFont(pdf, 8)
	pdf.MultiCell(0, 4, "POLIZZA DI RAMO VITA I  - Polizza “Wopta per te. Vita”", "LR", "CM", true)
	setWhiteItalicFont(pdf, smallTextSize)
	pdf.MultiCell(0, 3, "(da compilarsi in caso di scelta da parte del Contraente/Assicurato della garanzia Decesso)", "LBR", "CM", true)
	pdf.Ln(2)
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, "AVVERTENZA PRELIMINARE - Al fine di adempiere agli obblighi previsti dal "+
		"Decreto Legislativo 21 novembre 2007 n. 231 (di seguito il “Decreto”), in materia di prevenzione "+
		"del fenomeno del riciclaggio e del finanziamento del terrorismo, il Cliente (il soggetto Contraente/Assicurato "+
		"alla polizza “Wopta per te. Vita”) è tenuto a compilare e sottoscrivere il presente Modulo. Le "+
		"disposizioni del Decreto richiedono infatti, per una completa identificazione ed una adeguata conoscenza del "+
		"cliente e dell’eventuale titolare effettivo, la raccolta di informazioni ulteriori rispetto a quelle "+
		"anagrafiche già raccolte. La menzionata normativa impone al cliente di fornire, sotto la propria "+
		"responsabilità, tutte le informazioni necessarie ed aggiornate per consentire all’Intermediario di adempiere "+
		"agli obblighi di adeguata verifica e prevede specifiche sanzioni nel caso in cui le informazioni non "+
		"vengano fornite o risultino false.", "", "", false)
	pdf.Ln(3)

	pdf.MultiCell(0, 3, "Il conferimento dei dati e delle informazioni personali per l’identificazione "+
		"del Cliente e per la compilazione della presente sezione è obbligatorio per legge e, in caso di loro mancato "+
		"rilascio, la Compagnia Assicurativa non potrà procedere ad instaurare il rapporto (c.d. obbligo di "+
		"astensione), e dovrà valutare se effettuare una segnalazione alle autorità competenti (Unità di "+
		"Informazione Finanziaria presso Banca d’Italia e Guardia di Finanza). I dati saranno trattati per le "+
		"finalità di assolvimento degli obblighi previsti dalla normativa antiriciclaggio e, pertanto, tale "+
		"trattamento non richiede il consenso dell’interessato.", "", "", false)
	pdf.Ln(3)

	pdf.MultiCell(0, 3, "Io sottoscritto "+strings.ToUpper(insured.Surname+" "+insured.Name)+
		" (Contraente/Assicurato), letta l’Avvertenza Preliminare di cui sopra e l’Informativa sui Riferimenti Normativi"+
		" Antiriciclaggio (in calce al presente "+
		"modulo), al fine di permettere all’Intermediario di assolvere agli obblighi di adeguata verifica di cui al "+
		"D.Lgs. n. 231/2007 in materia di prevenzione dei fenomeni di riciclaggio e di finanziamento del terrorismo, "+
		"in relazione all’instaurazione del rapporto assicurativo di cui al contratto di assicurazione “Wopta per te. "+
		"Vita” - che prevede una garanzia di ramo vita emessa dall’impresa AXA France VIE S.A. (Rappresentanza "+
		"Generale per l’Italia):", "", "", false)
	pdf.Ln(3)

	pdf.MultiCell(0, 3, "A. dichiaro che i seguenti dati riportati relativi alla mia persona "+
		"corrispondono al vero ", "", "", false)
	pdf.CellFormat(5, 3, "", "", 0, "", false, 0, "")
	setWhiteBoldFont(pdf, standardTextSize)
	pdf.CellFormat(180, 4, "DATI IDENTIFICATIVI DEL CLIENTE (CONTRAENTE/ASSICURATO)", "TLR",
		0, "CM", true, 0, "")
	pdf.CellFormat(5, 4, "", "", 1, "", false, 0, "")
	pdf.CellFormat(5, 4, "", "", 0, "", false, 0, "")
	setBlackBoldFont(pdf, standardTextSize)
	pdf.CellFormat(90, 4, "Nome: "+strings.ToUpper(insured.Name), "TLR", 0, "",
		false, 0, "")
	pdf.CellFormat(90, 4, "Cognome:  "+strings.ToUpper(insured.Surname), "TLR", 0, "",
		false, 0, "")
	pdf.CellFormat(5, 4, "", "", 1, "", false, 0, "")
	pdf.CellFormat(5, 4, "", "", 0, "", false, 0, "")
	pdf.CellFormat(90, 4, "Data di nascita: "+birthDate.Format(dateLayout), "TLR", 0, "",
		false, 0, "")
	pdf.CellFormat(90, 4, "Codice Fiscale: "+strings.ToUpper(insured.FiscalCode), "TLR", 0,
		"", false, 0, "")
	pdf.CellFormat(5, 4, "", "", 1, "", false, 0, "")
	pdf.CellFormat(5, 4, "", "", 0, "", false, 0, "")
	pdf.CellFormat(90, 4, "Comune di nascita: "+strings.ToUpper(insured.BirthCity), "TLR", 0,
		"", false, 0, "")
	pdf.CellFormat(45, 4, "CAP: "+insured.PostalCode, "TLR", 0, "", false,
		0, "")
	pdf.CellFormat(45, 4, "Prov.: "+insured.BirthProvince, "TLR", 0, "", false,
		0, "")
	pdf.CellFormat(5, 4, "", "", 1, "", false, 0, "")
	pdf.CellFormat(5, 4, "", "", 0, "", false, 0, "")
	pdf.CellFormat(90, 4, "Comune di residenza: "+strings.ToUpper(insured.Residence.City), "TLR",
		0, "", false, 0, "")
	pdf.CellFormat(45, 4, "CAP: "+insured.Residence.PostalCode, "TLR", 0,
		"", false, 0, "")
	pdf.CellFormat(45, 4, "Prov.: "+strings.ToUpper(insured.Residence.CityCode),
		"TLR", 0, "", false, 0, "")
	pdf.CellFormat(5, 4, "", "", 1, "", false, 0, "")
	pdf.CellFormat(5, 4, "", "", 0, "", false, 0, "")
	pdf.CellFormat(180, 4, "Indirizzo di residenza: "+strings.ToUpper(insured.Address+" "+
		insured.StreetNumber), "TLR", 0, "", false, 0, "")
	pdf.CellFormat(5, 4, "", "", 1, "", false, 0, "")
	pdf.CellFormat(5, 4, "", "", 0, "", false, 0, "")
	pdf.CellFormat(180, 4, "Comune di domicilio (se diverso dalla residenza): "+domicileCity,
		"TLR", 0, "", false, 0, "")
	pdf.CellFormat(5, 4, "", "", 1, "", false, 0, "")
	pdf.CellFormat(5, 4, "", "", 0, "", false, 0, "")
	pdf.CellFormat(180, 4, "Indirizzo di domicilio (se diverso dalla residenza): "+domicileAddress,
		"LR", 0, "", false, 0, "")
	pdf.CellFormat(5, 4, "", "", 1, "", false, 0, "")
	pdf.CellFormat(5, 4, "", "", 0, "", false, 0, "")
	pdf.CellFormat(180, 4, "Status occupazionale: "+insured.WorkType, "TLR", 0, "",
		false, 0, "")
	pdf.CellFormat(5, 4, "", "", 1, "", false, 0, "")
	pdf.CellFormat(5, 4, "", "", 0, "", false, 0, "")
	pdf.CellFormat(180, 4, "Se Altro (specificare):", "BLR", 0, "", false, 0, "")
	pdf.CellFormat(5, 4, "", "L", 1, "", false, 0, "")
	pdf.Ln(1)

	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, "B. allego una fotocopia fronte/retro del mio documento di identità non scaduto "+
		"avente i seguenti estremi, confermando la veridicità dei dati sotto riportati: ", "", "", false)
	pdf.CellFormat(5, 4, "", "", 0, "", false, 0, "")
	pdf.CellFormat(90, 4, "Tipo documento: "+identityDocumentInfo["type"], "TLR", 0, "", false, 0, "")
	pdf.CellFormat(90, 4, "Nr. Documento: "+identityDocumentInfo["number"], "TLR", 0, "", false, 0, "")
	pdf.CellFormat(5, 4, "", "", 1, "", false, 0, "")
	pdf.CellFormat(5, 4, "", "", 0, "", false, 0, "")
	pdf.CellFormat(90, 4, "Ente di rilascio: "+identityDocumentInfo["issuingAuthority"], "TLR", 0, "", false, 0, "")
	pdf.CellFormat(90, 4, "Data di rilascio: "+identityDocumentInfo["dateOfIssue"], "TLR", 0, "", false, 0, "")
	pdf.CellFormat(5, 4, "", "", 1, "", false, 0, "")
	pdf.CellFormat(5, 4, "", "", 0, "", false, 0, "")
	pdf.CellFormat(90, 4, "Località di rilascio: "+identityDocumentInfo["placeOfIssue"], "1", 0, "", false, 0, "")
	pdf.CellFormat(90, 4, "Data di scadenza: "+identityDocumentInfo["expiryDate"], "1", 1, "", false, 0, "")
	pdf.Ln(1)

	pdf.MultiCell(0, 3, "C. dichiaro di NON essere una Persona Politicamente Esposta", "",
		"", false)
	pdf.CellFormat(4, 3, "", "", 0, "", false, 0, "")
	pdf.CellFormat(0, 3, "In caso di risposta affermativa indicare la tipologia:", "", 1,
		"", false, 0, "")

	pdf.MultiCell(0, 3, "D. dichiaro di NON essere destinatario di misure di congelamento dei fondi e "+
		"risorse economiche", "", "", false)
	pdf.CellFormat(4, 3, "", "", 0, "", false, 0, "")
	pdf.CellFormat(0, 3, "In caso di risposta affermativa indicare il motivo:", "", 1,
		"", false, 0, "")

	pdf.MultiCell(0, 3, "E. dichiaro di NON essere sottoposto a procedimenti o di NON aver subito condanne "+
		"per reati in materia economica/ finanziaria/tributaria/societaria", "", "", false)
	pdf.CellFormat(4, 3, "", "", 0, "", false, 0, "")
	pdf.CellFormat(0, 3, "In caso di risposta affermativa indicare il motivo:", "", 1,
		"", false, 0, "")

	pdf.MultiCell(0, 3, "F. dichiaro ai fini dell'identificazione del Titolare Effettivo, di essere una "+
		"persona fisica che agisce in nome e per conto proprio, di essere il soggetto Contraente/Assicurato, e "+
		"quindi che non esiste il titolare effettivo", "", "", false)

	pdf.MultiCell(0, 3, "G. fornisco, con riferimento allo scopo e alla natura prevista del rapporto "+
		"continuativo, le seguenti informazioni", "", "", false)
	pdf.CellFormat(4, 8, "", "", 0, "", false, 0, "")
	pdf.MultiCell(0, 3, "i. Tipologia di rapporto continuativo (informazione immediatamente desunta dal "+
		"rapporto): Stipula di un contratto di assicurazione di puro rischio che prevede garanzia di ramo vita "+
		"(caso morte Assicurato)", "", "", false)
	pdf.CellFormat(4, 12, "", "", 0, "", false, 0, "")
	pdf.MultiCell(0, 3, "ii. Scopo prevalente del rapporto continuativo in riferimento alle garanzie vita"+
		" (informazione immediatamente desunta dal rapporto):Protezione assicurativa al fine di garantire ai "+
		"beneficiari un capitale qualora si verifichi l’evento oggetto di copertura", "", "", false)
	pdf.CellFormat(4, 3, "", "", 0, "", false, 0, "")
	pdf.MultiCell(0, 3, "iii.  Origine dei fondi utilizzati per il pagamento dei premi assicurativi: "+
		"Proprie risorse economiche", "", "", false)
	pdf.CellFormat(0, 2, "", "", 1, "", false, 0, "")
}

func GetAxaTablePart2Section(pdf *fpdf.Fpdf, policy models.Policy) {
	pdf.MultiCell(0, 3, "Il sottoscritto, ai sensi degli artt. 22 e 55 comma 3 del d.lgs. 231/2007, "+
		"consapevole della responsabilità penale derivante da omesse e/o mendaci affermazioni, dichiara che tutte le "+
		"informazioni fornite (anche in riferimento al titolare effettivo), le dichiarazioni rilasciate il documento "+
		"di identità che allego, ed i dati riprodotti negli appositi campi del Modulo di Polizza corrispondono al "+
		"vero. Il sottoscritto si assume tutte le responsabilità di natura civile, amministrativa e penale per "+
		"dichiarazioni non veritiere. Il sottoscritto si impegna a comunicare senza ritardo a AXA France VIE S.A. "+
		"(Rappresentanza Generale per l’Italia) ogni eventuale integrazione o variazione che si dovesse verificare "+
		"in relazione ai dati ed alle informazioni forniti con il presente modulo.", "", "", false)
	pdf.Ln(4)

	setBlackBoldFont(pdf, standardTextSize)
	pdf.CellFormat(30, 3, "Data "+policy.EmitDate.Format(dateLayout), "", 0, "CM",
		false, 0, "")
	drawSignatureForm(pdf)
}

func GetAxaTablePart3Section(pdf *fpdf.Fpdf) {
	setBlackBoldFont(pdf, titleTextSize)
	pdf.MultiCell(0, 4, "Informativa antiriciclaggio (articoli di riferimento) - "+
		"(Decreto legislativo n. 231/2007)", "", "CM", false)
	pdf.Ln(4)

	setBlackBoldFont(pdf, smallTextSize)
	pdf.MultiCell(0, 3, "Obbligo di astensione – art. 42", "", "", false)
	setBlackRegularFont(pdf, smallTextSize)
	pdf.MultiCell(0, 3, "1. I soggetti obbligati che si trovano nell’impossibilità oggettiva di "+
		"effettuare l'adeguata verifica della clientela, ai sensi delle disposizioni di cui all'articolo 19, "+
		"comma 1, lettere a), b) e c), si astengono dall'instaurare, eseguire ovvero proseguire il rapporto, la "+
		"prestazione professionale e le operazioni e valutano se effettuare una segnalazione di operazione sospetta "+
		"alla UIF a norma dell'articolo 35.", "", "", false)
	pdf.MultiCell(0, 3, "2. I soggetti obbligati si astengono dall'instaurare il rapporto continuativo, "+
		"eseguire operazioni o prestazioni professionali e pongono fine al rapporto continuativo o alla prestazione "+
		"professionale già in essere di cui siano, direttamente o indirettamente, parte società fiduciarie, trust, "+
		"società anonime o controllate attraverso azioni al portatore aventi sede in Paesi terzi ad alto rischio. "+
		"Tali misure si applicano anche nei confronti delle ulteriori entità giuridiche, altrimenti denominate, "+
		"aventi sede nei suddetti Paesi, di cui non è possibile identificare il titolare effettivo ne' verificarne "+
		"l’identità.", "", "", false)
	pdf.MultiCell(0, 3, "3. (…).", "", "", false)
	pdf.MultiCell(0, 3, "4. È fatta in ogni caso salva l'applicazione dell'articolo 35, comma 2, nei "+
		"casi in cui l'operazione debba essere eseguita in quanto sussiste un obbligo di legge di ricevere "+
		"l'atto.", "", "", false)
	setBlackBoldFont(pdf, smallTextSize)
	pdf.MultiCell(0, 3, "Obblighi del cliente / sanzioni", "", "", false)
	setBlackRegularFont(pdf, smallTextSize)
	pdf.MultiCell(0, 3, "Art. 22, comma 1 - I clienti forniscono per iscritto, sotto la propria "+
		"responsabilità, tutte le informazioni necessarie e aggiornate per consentire ai soggetti obbligati di "+
		"adempiere agli obblighi di adeguata verifica.", "", "", false)
	pdf.MultiCell(0, 3, "Art. 55, comma 3 - Salvo che il fatto costituisca più grave reato, chiunque "+
		"essendo obbligato, ai sensi del presente decreto, a fornire i dati e le informazioni necessarie ai fini "+
		"dell'adeguata verifica della clientela, fornisce dati falsi o informazioni non veritiere, e' punito con la "+
		"reclusione da sei mesi a tre anni e con la multa da 10.000 euro a 30.000 "+
		"euro", "", "", false)
	setBlackBoldFont(pdf, smallTextSize)
	pdf.MultiCell(0, 3, "Nozione di titolare effettivo", "", "", false)
	pdf.MultiCell(0, 3, "Art.1, comma 2, lett. pp) del D. Lgs. n.231/2007 ", "", "",
		false)
	setBlackRegularFont(pdf, smallTextSize)
	pdf.MultiCell(0, 3, "la persona fisica o le persone fisiche, diverse dal cliente, nell'interesse "+
		"della quale  o  delle  quali,  in ultima istanza, il rapporto continuativo è istaurato, la prestazione "+
		"professionale è resa o l'operazione è eseguita.", "", "", false)
	setBlackBoldFont(pdf, smallTextSize)
	pdf.MultiCell(0, 3, "Nozione di persona politicamente esposta", "", "", false)
	pdf.MultiCell(0, 3, "Art. 1, comma 1, lettera dd) D. Lgs. 231/2007 così come modificato dal D. Lgs."+
		" 125/2019", "", "", false)
	setBlackRegularFont(pdf, smallTextSize)
	pdf.MultiCell(0, 3, "Persone politicamente esposte: le persone fisiche che occupano o hanno "+
		"cessato di occupare da meno di un anno importanti cariche pubbliche, nonché i loro familiari e coloro "+
		"che con i predetti soggetti intrattengono notoriamente stretti legami, come di "+
		"seguito elencate:", "", "", false)

	pdf.MultiCell(0, 3, "1) sono persone fisiche che occupano o hanno occupato importanti cariche "+
		"pubbliche coloro che ricoprono o hanno ricoperto la carica di:", "", "", false)
	indentedText(pdf, "1.1 Presidente della Repubblica, Presidente del Consiglio, Ministro, "+
		"Vice-Ministro e Sottosegretario, Presidente di Regione, assessore regionale, Sindaco di capoluogo di "+
		"provincia o città metropolitana, Sindaco di comune con popolazione non inferiore a 15.000 abitanti "+
		"nonché cariche analoghe in Stati esteri;")
	indentedText(pdf, "1.2 deputato, senatore, parlamentare europeo, consigliere regionale "+
		"nonché cariche analoghe in Stati esteri;")
	indentedText(pdf, "1.3 membro degli organi direttivi centrali di partiti politici;")
	indentedText(pdf, "1.4 giudice della Corte Costituzionale, magistrato della Corte di Cassazione "+
		"o della Corte dei conti, consigliere di Stato e altri componenti del Consiglio di Giustizia Amministrativa "+
		"per la Regione siciliana nonché cariche analoghe in Stati esteri;")
	indentedText(pdf, "1.5 membro degli organi direttivi delle banche centrali e delle autorità indipendenti;")
	indentedText(pdf, "1.6 ambasciatore, incaricato d’affari ovvero cariche equivalenti in Stati "+
		"esteri, ufficiale di grado apicale delle forze armate ovvero cariche analoghe in "+
		"Stati esteri;")
	indentedText(pdf, "1.7 componente degli organi di amministrazione, direzione o controllo delle "+
		"imprese controllate, anche indirettamente, dallo Stato italiano o da uno Stato estero ovvero partecipate, "+
		"in misura prevalente o totalitaria, dalle Regioni, da comuni capoluoghi di provincia e città metropolitane "+
		"e da comuni con popolazione complessivamente non inferiore a 15.000 "+
		"abitanti;")
	indentedText(pdf, "1.8 direttore generale di ASL e di azienda ospedaliera, di azienda ospedaliera "+
		"universitaria e degli altri enti del servizio sanitario nazionale.")
	indentedText(pdf, "1.9 direttore, vicedirettore e membro dell’organo di gestione o soggetto "+
		"svolgenti funzioni equivalenti in organizzazioni internazionali;")

	pdf.MultiCell(0, 3, "2) sono familiari di persone politicamente esposte: i genitori, il coniuge o "+
		"la persona legata in unione civile o convivenza di fatto o istituti assimilabili alla persona politicamente "+
		"esposta, i figli e i loro coniugi nonché le persone legate ai figli in unione civile o convivenza di fatto "+
		"o istituti assimilabili;", "", "", false)

	pdf.MultiCell(0, 3, "3) sono soggetti con i quali le persone politicamente esposte intrattengono "+
		"notoriamente stretti legami:", "", "", false)
	indentedText(pdf, "3.1 le persone fisiche che ai sensi del presente decreto detengono, "+
		"congiuntamente alla persona politicamente esposta, la titolarità effettiva di enti giuridici, trust  e "+
		"istituti giuridici affini ovvero che intrattengono con la persona politicamente esposta stretti rapporti "+
		"di affari;")
	indentedText(pdf, "3.2 le persone fisiche che detengono solo formalmente il controllo totalitario "+
		"di un’entità notoriamente costituita, di fatto, nell’interesse e a beneficio di una persona politicamente "+
		"esposta.")
}

func GetWoptaInfoTable(pdf *fpdf.Fpdf) {
	drawPinkHorizontalLine(pdf, 0.1)
	pdf.Ln(0.5)
	setBlackRegularFont(pdf, smallTextSize)
	pdf.MultiCell(0, 3, "DATI DELLA PERSONA FISICA CHE ENTRA IN CONTATTO CON IL "+
		"CONTRAENTE", "", "", false)
	pdf.Ln(1)
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, "LOMAZZI MICHELE iscritto alla Sezione A del RUI con numero "+
		"A000703480 in data 02.03.2022", "", "", false)
	pdf.Ln(0.5)
	drawPinkHorizontalLine(pdf, 0.1)
	pdf.Ln(0.5)
	setBlackRegularFont(pdf, smallTextSize)
	pdf.MultiCell(0, 3, "QUALIFICA", "", "", false)
	pdf.Ln(1)
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, "Responsabile dell’attività di intermediazione assicurativa di Wopta "+
		"Assicurazioni Srl, Società iscritta alla Sezione A del RUI con numero A000701923 in data "+
		"14.02.2022", "", "", false)
	pdf.Ln(0.5)
	drawPinkHorizontalLine(pdf, 0.1)
	pdf.Ln(0.5)
	setBlackRegularFont(pdf, smallTextSize)
	pdf.MultiCell(0, 3, "SEDE LEGALE", "", "", false)
	pdf.Ln(1)
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, "Galleria del Corso, 1 – 20122 MILANO (VI)", "", "", false)
	pdf.Ln(0.5)
	drawPinkHorizontalLine(pdf, 0.1)
	pdf.Ln(0.5)
	setBlackRegularFont(pdf, smallTextSize)
	pdf.Cell(50, 3, "RECAPITI TELEFONICI")
	pdf.Cell(40, 3, "")
	pdf.MultiCell(50, 3, "E-MAIL", "", "1", false)
	pdf.Ln(1)
	setBlackRegularFont(pdf, standardTextSize)
	pdf.Cell(50, 3, "02.91.24.03.46")
	pdf.Cell(40, 3, "")
	pdf.MultiCell(50, 3, "info@wopta.it", "", "1", false)
	pdf.Ln(0.5)
	drawPinkHorizontalLine(pdf, 0.1)
	pdf.Ln(0.5)
	setBlackRegularFont(pdf, smallTextSize)
	pdf.Cell(50, 3, "PEC ")
	pdf.Cell(40, 3, "")
	pdf.MultiCell(50, 3, "SITO INTERNET", "", "1", false)
	pdf.Ln(1)
	setBlackRegularFont(pdf, standardTextSize)
	pdf.Cell(50, 3, "woptaassicurazioni@legalmail.it")
	pdf.Cell(40, 3, "")
	pdf.MultiCell(50, 3, "wopta.it", "", "1", false)
	pdf.Ln(0.5)
	drawPinkHorizontalLine(pdf, 0.1)
	pdf.Ln(0.5)
	setBlackRegularFont(pdf, smallTextSize)
	pdf.MultiCell(0, 3, "AUTORITÀ COMPETENTE ALLA VIGILANZA DELL’ATTIVITÀ SVOLTA",
		"", "", false)
	pdf.Ln(1)
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, "IVASS – Istituto per la Vigilanza sulle Assicurazioni - Via del Quirinale, "+
		"21 - 00187 Roma", "", "", false)
	pdf.Ln(0.5)
	drawPinkHorizontalLine(pdf, 0.1)
}

func GetAllegato3Section(pdf *fpdf.Fpdf) {
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

	GetWoptaInfoTable(pdf)
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

func GetAllegato4Section(pdf *fpdf.Fpdf) {
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

	GetWoptaInfoTable(pdf)
	pdf.Ln(3)

	setBlackBoldFont(pdf, titleTextSize)
	pdf.MultiCell(0, 3, "SEZIONE I - Informazioni sul modello di distribuzione", "",
		"", false)
	pdf.Ln(1)
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, "Secondo quanto indicato nel modulo di proposta/polizza e documentazione "+
		"precontrattuale ricevuta, la distribuzione relativamente a questa proposta/contratto è svolta per conto "+
		"della seguente impresa di assicurazione: AXA FRANCE VIE S.A.", "", "", false)
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

func GetAllegato4TerSection(pdf *fpdf.Fpdf) {
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

	GetWoptaInfoTable(pdf)
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

func GetWoptaPrivacySection(pdf *fpdf.Fpdf) {
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

func GetPersonalDataHandlingSection(pdf *fpdf.Fpdf, policy models.Policy) {
	consentText := "X"
	notConsentText := ""

	if policy.Contractor.Consens != nil {
		consent, err := policy.ExtractConsens(1)
		lib.CheckError(err)

		if !consent.Answer {
			consentText = ""
			notConsentText = "X"
		}
	}

	setPinkBoldFont(pdf, titleTextSize)
	pdf.MultiCell(0, 3, "Consenso per finalità commerciali.", "", "", false)
	pdf.Ln(1)
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
