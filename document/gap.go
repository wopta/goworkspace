package document

import (
	"fmt"
	"github.com/go-pdf/fpdf"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"math"
	"sort"
	"strings"
	"time"
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

	gapHeader(pdf, policy)

	gapFooter(pdf)

	pdf.AddPage()

	vehicle := policy.Assets[0].Vehicle
	contractor := policy.Contractor
	vehicleOwner := policy.Assets[0].Person
	statements := *policy.Statements

	getParagraphTitle(pdf, "La tua assicurazione è operante sui dati sotto riportati, verifica la loro correttezza"+
		" e segnala eventuali inesattezze")
	pdf.Ln(3)

	gapVehicleDataTable(pdf, vehicle)

	gapPersonalInfoTable(pdf, contractor, *vehicleOwner)

	gapPolicyDataTable(pdf, policy)

	gapPriceTable(pdf, policy)

	pdf.Ln(3)

	gapStatements(pdf, statements[:len(statements)-1], policy.Company)

	companiesDescriptionSection(pdf, policy.Company)

	woptaGapHeader(pdf, *policy)

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

	location, err := time.LoadLocation("Europe/Rome")
	lib.CheckError(err)

	policyStartDate := policy.StartDate.In(location)
	policyEndDate := policy.EndDate.In(location)

	logoPath = lib.GetAssetPathByEnv(basePath) + "/logo_gap.png"
	productName = "Auto Valore Protetto"

	policyInfo := "Polizza Numero: " + policy.CodeCompany + "\n" +
		"Targa Veicolo: " + policy.Assets[0].Vehicle.Plate + "\n" +
		"Decorre dal: " + policyStartDate.Format(dateLayout) + " ore 24:00\n" +
		"Scade il: " + policyEndDate.Format(dateLayout) + " ore 24:00"

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

func gapFooter(pdf *fpdf.Fpdf) {
	footerText := "Wopta per te. Auto Valore Protetto è un prodotto assicurativo di Sogessur - Société Anonyme " +
		"– Capitale Sociale € 33 825 000 – Sede legale: Tour D2, 17bis Place des Reflets – 92919\n" +
		"Paris La Défense Cedex - 379 846 637 R.C.S. Nanterre - Francia - Sede secondaria: Via Tiziano 32, " +
		"20145 Milano - Italia - Registro delle Imprese di Milano, Lodi, Monza-Brianza\n" +
		"Codice Fiscale e P.IVA  07420570967  Iscritta nell’elenco I dell’Albo delle Imprese di Assicurazione tenuto " +
		"dall’IVASS al n. I00094"

	pdf.SetFooterFunc(func() {
		pdf.SetXY(10, -17)
		setPinkRegularFont(pdf, smallTextSize)
		pdf.MultiCell(0, 3, footerText, "", "", false)
		pdf.SetY(-8)
		setBlackRegularFont(pdf, smallTextSize)
		pdf.MultiCell(0, 3, "Auto Valore Protetto - VI - Settembre_2023", "", fpdf.AlignRight,
			false)
		pageNumber(pdf)
	})
}

func woptaGapHeader(pdf *fpdf.Fpdf, policy models.Policy) {
	location, err := time.LoadLocation("Europe/Rome")
	lib.CheckError(err)

	policyStartDate := policy.StartDate.In(location)
	policyEndDate := policy.EndDate.In(location)

	policyInfo := "Polizza Numero: " + policy.CodeCompany + "\n" +
		"Targa Veicolo: " + policy.Assets[0].Vehicle.Plate + "\n" +
		"Decorre dal: " + policyStartDate.Format(dateLayout) + " ore 24:00\n" +
		"Scade il: " + policyEndDate.Format(dateLayout) + " ore 24:00"

	pdf.SetHeaderFunc(func() {
		var opt fpdf.ImageOptions
		opt.ImageType = "png"
		pdf.ImageOptions(lib.GetAssetPathByEnv(basePath)+"/ARTW_LOGO_RGB_400px.png", 11, 6, 0, 10,
			false, opt, 0, "")

		setBlackBoldFont(pdf, standardTextSize)
		pdf.SetXY(11, 17)
		pdf.Cell(0, 3, "I dati della tua polizza")
		setBlackRegularFont(pdf, standardTextSize)
		pdf.SetXY(11, pdf.GetY()+3)
		pdf.MultiCell(0, 3.5, policyInfo, "", "", false)
		pdf.Ln(8)
	})

}

func gapVehicleDataTable(pdf *fpdf.Fpdf, vehicle *models.Vehicle) {
	tableRows := [][]string{
		{"Tipo Veicolo", vehicle.VehicleTypeDesc, "Data prima immatricolazione", vehicle.RegistrationDate.Format(dateLayout)},
		{"Marca", vehicle.Manufacturer, "Stato veicolo", vehicle.Condition},
		{"Modello", vehicle.Model, "Valore veicolo (*)", lib.HumanaizePriceEuro(vehicle.PriceValue)},
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

func gapPersonalInfoTable(pdf *fpdf.Fpdf, contractor, vehicleOwner models.User) {
	setWhiteBoldFont(pdf, standardTextSize)
	pdf.SetFillColor(229, 0, 117)
	pdf.SetDrawColor(229, 0, 117)
	pdf.CellFormat(30, 5, "Dati Personali", "1", 0, fpdf.AlignLeft, true, 0, "")
	pdf.CellFormat(65, 5, "Contraente", "1", 0, fpdf.AlignCenter, true, 0, "")
	pdf.CellFormat(95, 5, "Proprietario", "1", 1, fpdf.AlignCenter, true, 0, "")

	contractorBirthDate, err := time.Parse(time.RFC3339, contractor.BirthDate)
	lib.CheckError(err)
	vehicleOwnerBirthDate, err := time.Parse(time.RFC3339, vehicleOwner.BirthDate)
	lib.CheckError(err)

	tableRows := [][]string{
		{"Cognome e Nome", contractor.Surname + " " + contractor.Name, "Cognome e Nome",
			vehicleOwner.Surname + " " + vehicleOwner.Name},
		{"Residente in", contractor.Residence.StreetName + " " + contractor.Residence.StreetNumber + ", " +
			"" + contractor.Residence.PostalCode + ", " + contractor.Residence.City + "(" + contractor.Residence.
			CityCode + ")",
			"Residente in", vehicleOwner.Residence.StreetName + " " + vehicleOwner.Residence.StreetNumber + ", " +
				"" + vehicleOwner.Residence.PostalCode + ", " + vehicleOwner.Residence.City + "(" + vehicleOwner.Residence.
				CityCode + ")"},
		{"Mail", contractor.Mail, "Mail", "================"},
		{"Codice Fiscale", contractor.FiscalCode, "Codice Fiscale", vehicleOwner.FiscalCode},
		{"Data nascita", contractorBirthDate.Format(dateLayout), "Data nascita", vehicleOwnerBirthDate.Format(dateLayout)},
		{"Telefono", contractor.Phone, "Telefono", "================"},
	}

	lastRowBordersList := []string{"BL", "B", "B", "BR"}

	for x := 0; x < len(tableRows); x++ {
		bordersList := []string{"L", "B", "", "BR"}

		setBlackRegularFont(pdf, 8)
		numLines := math.Max(float64(len(pdf.SplitText(tableRows[x][1], 55))),
			float64(len(pdf.SplitText(tableRows[x][3], 55))))

		if x == len(tableRows)-1 {
			bordersList = lastRowBordersList
		}

		setPinkRegularFont(pdf, 8)
		pdf.CellFormat(40, 5*numLines, tableRows[x][0], bordersList[0], 0, fpdf.AlignLeft, false, 0, "")

		drawDynamicCell(pdf, 8, 5, 55, numLines, 95, tableRows[x][1], bordersList[1], fpdf.AlignLeft, false)

		setPinkRegularFont(pdf, 8)
		pdf.CellFormat(40, 5*numLines, tableRows[x][2], bordersList[2], 0, fpdf.AlignLeft, false, 0, "")

		drawDynamicCell(pdf, 8, 5, 55, numLines, 135, tableRows[x][3], bordersList[3], fpdf.AlignLeft, true)
	}
	pdf.Ln(5)
}

func gapPolicyDataTable(pdf *fpdf.Fpdf, policy *models.Policy) {
	offerMap := map[string]string{
		"base":     "Base",
		"complete": "Completa",
	}

	location, err := time.LoadLocation("Europe/Rome")
	lib.CheckError(err)

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
	pdf.CellFormat(55, 5, policy.StartDate.In(location).Format(dateLayout), "B", 0, fpdf.AlignLeft, false,
		0, "")
	setPinkRegularFont(pdf, 8)
	pdf.CellFormat(40, 5, "Ore", "B", 0, fpdf.AlignLeft, false, 0, "")
	setBlackRegularFont(pdf, 8)
	pdf.CellFormat(55, 5, "24:00", "BR", 1, fpdf.AlignLeft, false, 0, "")

	setPinkRegularFont(pdf, 8)
	pdf.CellFormat(40, 5, "Scadenza", "BL", 0, fpdf.AlignLeft, false, 0, "")
	setBlackRegularFont(pdf, 8)
	pdf.CellFormat(55, 5, policy.EndDate.In(location).Format(dateLayout), "B", 0, fpdf.AlignLeft, false,
		0, "")
	setPinkRegularFont(pdf, 8)
	pdf.CellFormat(40, 5, "Ore", "B", 0, fpdf.AlignLeft, false, 0, "")
	setBlackRegularFont(pdf, 8)
	pdf.CellFormat(55, 5, "24:00", "BR", 1, fpdf.AlignLeft, false, 0, "")

	duration := lib.MonthsDifference(policy.StartDate.In(location), policy.EndDate.In(location))
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
