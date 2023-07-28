package document

import (
	"fmt"
	"github.com/go-pdf/fpdf"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/question"
	"sort"
	"strings"
)

func GapContract(pdf *fpdf.Fpdf, origin string, policy *models.Policy) (string, []byte) {
	var (
		filename string
		out      []byte
	)

	filename, out = GapSogessur(pdf, origin, policy)

	return filename, out
}

func GapSogessur(pdf *fpdf.Fpdf, origin string, policy *models.Policy) (string, []byte) {
	signatureID = 0

	var statements []models.Statement

	if policy.Statements == nil {
		statements = question.GetStatements(*policy)
	}

	gapHeader(pdf, policy)

	mainFooter(pdf, policy.Name)

	pdf.AddPage()

	getParagraphTitle(pdf, "La tua assicurazione è operante sui dati sotto riportati, verifica la loro correttezza"+
		" e segnala eventuali inesattezze")
	pdf.Ln(3)

	vehicle := policy.Assets[0].Vehicle
	contractor := policy.Contractor
	vehicleOwner := policy.Assets[0].Person

	gapVehicleDataTable(pdf, vehicle)

	gapPersonalInfoTable(pdf, contractor, *vehicleOwner)

	gapPolicyDataTable(pdf, policy)

	gapPriceTable(pdf, policy)

	pdf.Ln(3)

	gapStatements(pdf, statements[:len(statements)-1], policy.Company)

	//contractWithdrawlSection(pdf)

	//gapClauseApprovalSection(pdf)

	//pdf.AddPage()

	companiesDescriptionSection(pdf, policy.Company)

	/*sogessurHeader(pdf)

	pdf.AddPage()

	sogessurFooter(pdf)

	gapConsentDeclaration(pdf)*/

	woptaHeader(pdf)

	pdf.AddPage()

	woptaFooter(pdf)

	printStatement(pdf, statements[len(statements)-1], policy.Company)

	woptaHeader(pdf)

	pdf.AddPage()

	woptaPrivacySection(pdf)

	personalDataHandlingSection(pdf, policy)

	filename, out := save(pdf, policy)
	return filename, out
}

func gapHeader(pdf *fpdf.Fpdf, policy *models.Policy) {
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
		pdf.ImageOptions(lib.GetAssetPathByEnv(basePath)+"/ARTW_LOGO_RGB_400px.png", 135, 8, 0, 6, false, opt, 0, "")
		pdf.SetX(pdf.GetX() + 126)
		pdf.SetDrawColor(229, 0, 117)
		pdf.SetLineWidth(0.5)
		pdf.Line(pdf.GetX(), 8, pdf.GetX(), 14)
		pdf.ImageOptions(lib.GetAssetPathByEnv(basePath)+"/logo_sogessur.png", 161, 8, 0, 6, false, opt, 0, "")

		setBlackBoldFont(pdf, standardTextSize)
		pdf.SetXY(11, 20)
		pdf.Cell(0, 3, "I dati della tua polizza")
		setBlackRegularFont(pdf, standardTextSize)
		pdf.SetXY(11, pdf.GetY()+3)
		pdf.MultiCell(0, 3.5, policyInfo, "", "", false)
		pdf.Ln(8)
	})
}

func gapVehicleDataTable(pdf *fpdf.Fpdf, vehicle *models.Vehicle) {
	tableRows := [][]string{
		{"Tipo Veicolo", vehicle.VehicleType, "Data prima immatricolazione", vehicle.RegistrationDate.Format(dateLayout)},
		{"Marca", vehicle.Manufacturer, "Stato veicolo", vehicle.Condition},
		{"Modello", vehicle.Model, "Valore veicolo (*)", lib.HumanaizePriceEuro(float64(vehicle.PriceValue))},
	}

	setWhiteBoldFont(pdf, standardTextSize)
	pdf.SetFillColor(229, 0, 117)
	pdf.SetDrawColor(229, 0, 117)
	pdf.CellFormat(95, 5, "Dati Veicolo", "1", 0, fpdf.AlignLeft, true, 0, "")
	pdf.CellFormat(95, 5, "Targa: "+vehicle.Plate, "1", 1, fpdf.AlignLeft, true, 0, "")

	for x := 0; x < len(tableRows); x++ {
		setPinkRegularFont(pdf, 8)
		pdf.CellFormat(40, 5, tableRows[x][0], "L", 0, fpdf.AlignLeft, false, 0, "")
		setBlackRegularFont(pdf, 8)
		pdf.CellFormat(55, 5, tableRows[x][1], "B", 0, fpdf.AlignLeft, false, 0, "")
		setPinkRegularFont(pdf, 8)
		pdf.CellFormat(45, 5, tableRows[x][2], "", 0, fpdf.AlignLeft, false, 0, "")
		setBlackRegularFont(pdf, 8)
		pdf.CellFormat(50, 5, tableRows[x][3], "BR", 1, fpdf.AlignLeft, false, 0, "")
	}
	setBlackRegularFont(pdf, 7)
	pdf.MultiCell(0, 4, "Il veicolo deve essere Immatricolato in Italia ad uso privato, con peso a pieno "+
		"carico non eccedente le 3,5 tonnellate ed essere già coperto da una polizza furto e incendio. L’elenco di "+
		"tutte le condizioni di assicurabilità è presente nel Set Informativo", "LR",
		fpdf.AlignLeft, false)
	setPinkRegularFont(pdf, 7)
	pdf.CellFormat(pdf.GetStringWidth("(*) "), 4, "(*) ", "L", 0, fpdf.AlignLeft,
		false, 0, "")
	setBlackRegularFont(pdf, 7)
	pdf.MultiCell(0, 4, "Valore Veicolo si intende:", "R", fpdf.AlignLeft, false)
	pdf.MultiCell(0, 4, "- il valore di fattura se l’acquisto della polizza è contestuale all’acquisto "+
		"del veicolo;", "LR", fpdf.AlignLeft, false)
	pdf.MultiCell(0, 4, "- il valore commerciale al momento della sottoscrizione se l’acquisto della "+
		"polizza è differito dall’acquisto del veicolo.", "BLR", fpdf.AlignLeft, false)
	pdf.Ln(5)
}

func gapPersonalInfoTable(pdf *fpdf.Fpdf, contractor, insured models.User) {
	setWhiteBoldFont(pdf, standardTextSize)
	pdf.SetFillColor(229, 0, 117)
	pdf.SetDrawColor(229, 0, 117)
	pdf.CellFormat(30, 5, "Dati Personali", "1", 0, fpdf.AlignLeft, true, 0, "")
	pdf.CellFormat(65, 5, "Contraente", "1", 0, fpdf.AlignCenter, true, 0, "")
	pdf.CellFormat(95, 5, "Proprietario", "1", 1, fpdf.AlignCenter, true, 0, "")

	tableRows := [][]string{
		{"Cognome e Nome", contractor.Surname + " " + contractor.Name, "Cognome e Nome",
			insured.Surname + " " + insured.Name},
		{"Residente in", contractor.Residence.StreetName + " " + contractor.Residence.StreetNumber + ", " +
			"" + contractor.Residence.PostalCode + ", " + contractor.Residence.City + "(" + contractor.Residence.
			CityCode + ")",
			"Residente in", insured.Residence.StreetName + " " + insured.Residence.StreetNumber + ", " +
				"" + insured.Residence.PostalCode + ", " + insured.Residence.City + "(" + insured.Residence.
				CityCode + ")"},
		{"Mail", contractor.Mail, "Mail", insured.Mail},
		{"Codice Fiscale", contractor.FiscalCode, "Codice Fiscale", insured.FiscalCode},
		{"Data nascita", contractor.BirthDate, "Data nascita", insured.BirthDate},
		{"Telefono", contractor.Phone, "Telefono", insured.Phone},
	}

	for x := 0; x < len(tableRows); x++ {
		if x != len(tableRows)-1 {
			setPinkRegularFont(pdf, 8)
			pdf.CellFormat(40, 5, tableRows[x][0], "L", 0, fpdf.AlignLeft, false, 0, "")
			setBlackRegularFont(pdf, 8)
			pdf.CellFormat(55, 5, tableRows[x][1], "B", 0, fpdf.AlignLeft, false, 0, "")
			setPinkRegularFont(pdf, 8)
			pdf.CellFormat(40, 5, tableRows[x][2], "", 0, fpdf.AlignLeft, false, 0, "")
			setBlackRegularFont(pdf, 8)
			pdf.CellFormat(55, 5, tableRows[x][3], "BR", 1, fpdf.AlignLeft, false, 0, "")
		} else {
			setPinkRegularFont(pdf, 8)
			pdf.CellFormat(40, 5, tableRows[x][0], "BL", 0, fpdf.AlignLeft, false, 0, "")
			setBlackRegularFont(pdf, 8)
			pdf.CellFormat(55, 5, tableRows[x][1], "B", 0, fpdf.AlignLeft, false, 0, "")
			setPinkRegularFont(pdf, 8)
			pdf.CellFormat(40, 5, tableRows[x][2], "B", 0, fpdf.AlignLeft, false, 0, "")
			setBlackRegularFont(pdf, 8)
			pdf.CellFormat(55, 5, tableRows[x][3], "BR", 1, fpdf.AlignLeft, false, 0, "")
		}
	}
	pdf.Ln(5)
}

func gapPolicyDataTable(pdf *fpdf.Fpdf, policy *models.Policy) {
	offerMap := map[string]string{
		"base":     "Base",
		"complete": "Completa",
	}

	sort.Slice(policy.Assets[0].Guarantees, func(i, j int) bool {
		return policy.Assets[0].Guarantees[i].Order < policy.Assets[0].Guarantees[j].Order
	})

	var guaranteesNames []string

	for _, guarantee := range policy.Assets[0].Guarantees {
		guaranteesNames = append(guaranteesNames, guarantee.CompanyName)
	}

	setWhiteBoldFont(pdf, standardTextSize)
	pdf.SetFillColor(229, 0, 117)
	pdf.SetDrawColor(229, 0, 117)
	pdf.CellFormat(190, 5, "Dati di polizza", "1", 1, fpdf.AlignLeft, true, 0,
		"")

	setPinkRegularFont(pdf, 8)
	pdf.CellFormat(40, 5, "Decorrenza", "BL", 0, fpdf.AlignLeft, false, 0, "")
	setBlackRegularFont(pdf, 8)
	pdf.CellFormat(55, 5, policy.StartDate.Format(dateLayout), "B", 0, fpdf.AlignLeft, false,
		0, "")
	setPinkRegularFont(pdf, 8)
	pdf.CellFormat(40, 5, "Ore", "B", 0, fpdf.AlignLeft, false, 0, "")
	setBlackRegularFont(pdf, 8)
	pdf.CellFormat(55, 5, "24:00", "BR", 1, fpdf.AlignLeft, false, 0, "")

	setPinkRegularFont(pdf, 8)
	pdf.CellFormat(40, 5, "Scadenza", "BL", 0, fpdf.AlignLeft, false, 0, "")
	setBlackRegularFont(pdf, 8)
	pdf.CellFormat(55, 5, policy.EndDate.Format(dateLayout), "B", 0, fpdf.AlignLeft, false,
		0, "")
	setPinkRegularFont(pdf, 8)
	pdf.CellFormat(40, 5, "Ore", "B", 0, fpdf.AlignLeft, false, 0, "")
	setBlackRegularFont(pdf, 8)
	pdf.CellFormat(55, 5, "24:00", "BR", 1, fpdf.AlignLeft, false, 0, "")

	duration := lib.MonthsDifference(policy.StartDate, policy.EndDate)
	setPinkRegularFont(pdf, 8)
	pdf.CellFormat(40, 5, "Durata", "BL", 0, fpdf.AlignLeft, false, 0, "")
	setBlackRegularFont(pdf, 8)
	pdf.CellFormat(0, 5, fmt.Sprintf("%d mesi", duration), "BR", 1, fpdf.AlignLeft, false,
		0, "")
	setPinkBoldFont(pdf, 8)
	pdf.CellFormat(0, 5, "Opzione  di prodotto selezionata: "+offerMap[policy.OfferlName], "BLR", 1, fpdf.AlignLeft,
		false,
		0, "")
	setBlackRegularFont(pdf, 8)
	pdf.MultiCell(0, 4, "Include le seguenti garanzie: "+strings.Join(
		guaranteesNames, ", ")+"\nPer il dettaglio delle garanzie vedere il Set Informativo.",
		"BLR", fpdf.AlignLeft, false)

	pdf.Ln(5)
}

func gapPriceTable(pdf *fpdf.Fpdf, policy *models.Policy) {
	setWhiteBoldFont(pdf, standardTextSize)
	pdf.SetFillColor(229, 0, 117)
	pdf.SetDrawColor(229, 0, 117)
	pdf.CellFormat(30, 5, "Premio", "1", 0, fpdf.AlignLeft, true, 0, "")
	pdf.CellFormat(0, 5, "Unico Anticipato", "1", 1, fpdf.AlignCenter, true, 0,
		"")

	setPinkRegularFont(pdf, 8)
	pdf.CellFormat(20, 5, "", "BL", 0, fpdf.AlignLeft, false, 0, "")
	pdf.CellFormat(56, 5, "Imponibile", "B", 0, fpdf.AlignCenter, false, 0, "")
	pdf.CellFormat(56, 5, "Imposte Assicurative", "B", 0, fpdf.AlignCenter, false,
		0, "")
	pdf.CellFormat(0, 5, "Totale", "BR", 1, fpdf.AlignCenter, false, 0, "")

	pdf.CellFormat(20, 5, "Alla firma", "BL", 0, fpdf.AlignLeft, false, 0, "")
	setBlackRegularFont(pdf, 8)
	pdf.CellFormat(56, 5, lib.HumanaizePriceEuro(policy.PriceNett), "B", 0, fpdf.AlignCenter,
		false, 0, "")
	pdf.CellFormat(56, 5, lib.HumanaizePriceEuro(policy.PriceGross-policy.PriceNett), "B",
		0, fpdf.AlignCenter, false,
		0, "")
	pdf.CellFormat(0, 5, lib.HumanaizePriceEuro(policy.PriceGross), "BR", 1, fpdf.AlignCenter,
		false, 0, "")

	pdf.Ln(5)
}

func gapStatements(pdf *fpdf.Fpdf, statements []models.Statement, companyName string) {
	for _, statement := range statements {
		printStatement(pdf, statement, companyName)
	}
}

func gapClauseApprovalSection(pdf *fpdf.Fpdf) {
	getParagraphTitle(pdf, "Le clausole della Polizza da approvare in modo specifico")
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, "Il sottoscritto dichiara di approvare, ai sensi e per gli effetti degli artt. "+
		"1341 e 1342 del codice civile, i seguenti articoli delle Condizioni di Assicurazione: INSERIRE VESSATORIE",
		"", fpdf.AlignLeft, false)
	pdf.Ln(5)
	drawSignatureForm(pdf)
	pdf.Ln(5)
}

func gapConsentDeclaration(pdf *fpdf.Fpdf) {
	setBlackBoldFont(pdf, standardTextSize)
	pdf.SetDrawColor(0, 0, 0)
	pdf.MultiCell(0, 3, "Consenso al trattemento dei dati personali", "", fpdf.AlignLeft, false)
	pdf.Ln(1)
	pdf.SetLineWidth(thinLineWidth)
	pdf.Line(11, pdf.GetY(), 80, pdf.GetY())
	pdf.Ln(1)
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, "Il sottoscritto, dopo aver ricevuto copia e preso visione dell’Informativa "+
		"della Compagnia sul trattamento dei dati personali, ai sensi della normativa sulla privacy "+
		"(Reg. UE 2016/679) acconsente al trattamento dei propri dati personali, anche sensibili (particolari?), da "+
		"parte di Sogessur S.A. - Rappresentanza Generale per l’Italia, per le finalità, secondo le modalità e "+
		"mediante i soggetti indicati nella predetta informativa.", "", fpdf.AlignLeft, false)
	setBlackBoldFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, "Sono consapevole che il mancato consenso al trattamento dei dati "+
		"personali, necessari alla Compagnia per le finalità ivi illustrate, comporta l’impossibilità di "+
		"dare esecuzione al rapporto contrattuale.\"", "", fpdf.AlignLeft, false)
	pdf.Ln(3)
	drawSignatureForm(pdf)
}
