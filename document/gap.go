package document

import (
	"github.com/dustin/go-humanize"
	"github.com/go-pdf/fpdf"
	"github.com/wopta/goworkspace/models"
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

	//pageWidth, _ := pdf.GetPageSize()

	mainMotorHeader(pdf, policy)

	mainFooter(pdf, policy.Name)

	pdf.AddPage()

	getParagraphTitle(pdf, "La tua assicurazione è operante sui dati sotto riportati, verifica la loro correttezza"+
		" e segnala eventuali inesattezze")
	pdf.Ln(5)

	vehicle := policy.Assets[0].Vehicle
	contractor := policy.Contractor
	insured := policy.Assets[0].Person

	vehicleDataTable(pdf, vehicle)
	pdf.Ln(5)

	gapPersonalInfoTable(pdf, contractor, *insured)

	woptaHeader(pdf)

	pdf.AddPage()

	woptaFooter(pdf)

	woptaPrivacySection(pdf)

	personalDataHandlingSection(pdf, policy)

	filename, out := save(pdf, policy)
	return filename, out
}

func vehicleDataTable(pdf *fpdf.Fpdf, vehicle *models.Vehicle) {
	tableRows := [][]string{
		{"Tipo Veicolo", vehicle.VehicleType, "Data prima immatricolazione", vehicle.RegistrationDate.Format(dateLayout)},
		{"Marca", vehicle.Manufacturer, "Stato veicolo", vehicle.Condition},
		{"Modello", vehicle.Model, "Valore veicolo", humanize.FormatFloat("#.###,", float64(vehicle.PriceValue))},
	}

	setWhiteBoldFont(pdf, standardTextSize)
	pdf.SetFillColor(229, 0, 117)
	pdf.SetDrawColor(229, 0, 117)
	pdf.CellFormat(95, 5, "Dati Veicolo", "TBL", 0, fpdf.AlignLeft, true, 0, "")
	pdf.CellFormat(95, 5, "Targa: "+vehicle.Plate, "TBR", 1, fpdf.AlignLeft, true, 0, "")

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
	pdf.MultiCell(0, 4, "*Veicolo Immatricolato in Italia ad uso privato, il peso a pieno carico non"+
		" eccede le 3,5 tonnellate ed è già coperto da una polizza furto e incendio.", "1",
		fpdf.AlignLeft, false)
}

func gapPersonalInfoTable(pdf *fpdf.Fpdf, contractor, insured models.User) {
	setWhiteBoldFont(pdf, standardTextSize)
	pdf.SetFillColor(229, 0, 117)
	pdf.SetDrawColor(229, 0, 117)
	pdf.CellFormat(30, 5, "Dati Personali", "TBL", 0, fpdf.AlignLeft, true, 0, "")
	pdf.CellFormat(65, 5, "Contraente", "TBL", 0, fpdf.AlignCenter, true, 0, "")
	pdf.CellFormat(95, 5, "Assicurato", "TBR", 1, fpdf.AlignCenter, true, 0, "")

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

}
