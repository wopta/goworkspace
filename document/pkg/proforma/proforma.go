package proforma

import (
	"time"

	"github.com/wopta/goworkspace/document/internal/constants"
	"github.com/wopta/goworkspace/document/internal/domain"
	"github.com/wopta/goworkspace/document/internal/dto"
	"github.com/wopta/goworkspace/document/internal/engine"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

type ProformaGenerator struct {
	*baseGenerator
	dto *dto.ProformaDTO
}

func NewProformaGenerator(engine *engine.Fpdf, policy *models.Policy, node *models.NetworkNode,
	product models.Product) *ProformaGenerator {
	ProformaDTO := dto.NewProformaDTO()
	ProformaDTO.FromPolicy(*policy, product)
	return &ProformaGenerator{
		baseGenerator: &baseGenerator{
			engine:      engine,
			now:         time.Now(),
			networkNode: node,
			policy:      policy,
		},
		dto: ProformaDTO,
	}
}

func (pg *ProformaGenerator) Contract() ([]byte, error) {
	pg.mainHeader()

	pg.engine.NewPage()

	pg.engine.NewLine(20)

	pg.contractor()

	pg.engine.NewLine(6)

	pg.body()

	pg.woptaFooter()

	return pg.engine.RawDoc()
}

func (pg *ProformaGenerator) mainHeader() {
	pg.engine.SetHeader(func() {
		pg.engine.InsertImage(lib.GetAssetPathByEnvV2()+"logo_wopta.png", 20, 15, 40, 13)
		pg.engine.NewLine(7)
	})
}

func (pg *ProformaGenerator) contractor() {
	contr := pg.dto.Contractor
	const (
		firstColumnWidth  = 110
		secondColumnWidth = 80
	)
	parser := func(rows []string) [][]domain.TableCell {
		result := make([][]domain.TableCell, 0, len(rows))

		for i, row := range rows {
			fontStyle := constants.RegularFontStyle
			if i == 0 {
				fontStyle = constants.BoldFontStyle
			}
			result = append(result, []domain.TableCell{
				{
					Text:      "",
					Height:    constants.CellHeight,
					Width:     firstColumnWidth,
					FontStyle: constants.BoldFontStyle,
					FontColor: constants.BlackColor,
					FontSize:  constants.RegularFontSize,
					Fill:      false,
					FillColor: domain.Color{},
					Align:     constants.LeftAlign,
					Border:    "",
				},
				{
					Text:      row,
					Height:    constants.CellHeight,
					Width:     secondColumnWidth,
					FontStyle: fontStyle,
					FontColor: constants.BlackColor,
					FontSize:  constants.RegularFontSize,
					Fill:      false,
					FillColor: domain.Color{},
					Align:     constants.LeftAlign,
					Border:    "",
				},
			})
		}
		return result
	}

	data := []string{
		"I tuoi dati",
		"Contraente: " + contr.Name + " " + contr.Surname,
		"CF/P.IVA: " + contr.FiscalOrVatCode,
		"Indirizzo: " + contr.StreetNameAndNumber,
		contr.PostalCodeAndCity,
		"Mail: " + contr.Mail,
		"Telefono: " + contr.Phone,
	}

	text := parser(data)
	pg.engine.DrawTable(text)
}

func (pg *ProformaGenerator) body() {
	body := pg.dto.Body

	const (
		emptyColumnWidth  = 10
		oneLineWidth      = 180
		firstColumnWidth  = 100
		secondColumnWidth = 40
	)

	twoColParser := func(rows []string, boldFirstLine bool) [][]domain.TableCell {
		result := make([][]domain.TableCell, 0, len(rows))

		for i, row := range rows {
			fontStyle := constants.RegularFontStyle
			if i == 0 && boldFirstLine {
				fontStyle = constants.BoldFontStyle
			}
			result = append(result, []domain.TableCell{
				{
					Text:      " ",
					Height:    constants.CellHeight,
					Width:     emptyColumnWidth,
					FontStyle: constants.RegularFontStyle,
					FontColor: constants.BlackColor,
					FontSize:  constants.RegularFontSize,
					Fill:      false,
					FillColor: domain.Color{},
					Align:     constants.LeftAlign,
					Border:    "",
				},
				{
					Text:      row,
					Height:    constants.CellHeight,
					Width:     oneLineWidth,
					FontStyle: fontStyle,
					FontColor: constants.BlackColor,
					FontSize:  constants.RegularFontSize,
					Fill:      false,
					FillColor: domain.Color{},
					Align:     constants.LeftAlign,
					Border:    "",
				},
			})
		}
		return result
	}

	threeColParser := func(rows [][]string) [][]domain.TableCell {
		result := make([][]domain.TableCell, 0, len(rows))

		for _, row := range rows {
			result = append(result, []domain.TableCell{
				{
					Text:      " ",
					Height:    constants.CellHeight,
					Width:     emptyColumnWidth,
					FontStyle: constants.RegularFontStyle,
					FontColor: constants.BlackColor,
					FontSize:  constants.RegularFontSize,
					Fill:      false,
					FillColor: domain.Color{},
					Align:     constants.LeftAlign,
					Border:    "",
				},
				{
					Text:      row[0],
					Height:    constants.CellHeight,
					Width:     firstColumnWidth,
					FontStyle: constants.RegularFontStyle,
					FontColor: constants.BlackColor,
					FontSize:  constants.RegularFontSize,
					Fill:      false,
					FillColor: domain.Color{},
					Align:     constants.LeftAlign,
					Border:    "",
				},
				{
					Text:      row[1],
					Height:    constants.CellHeight,
					Width:     secondColumnWidth,
					FontStyle: constants.RegularFontStyle,
					FontColor: constants.BlackColor,
					FontSize:  constants.RegularFontSize,
					Fill:      false,
					FillColor: domain.Color{},
					Align:     constants.RightAlign,
					Border:    "",
				},
			})
		}
		return result
	}

	title := []string{
		"Proforma del " + body.Date,
	}

	data := [][]string{
		{"Contributo per servizi di intermediazione:", body.Gross + " Euro"},
		{"Imponibile:", body.Net + " Euro"},
		{"IVA (esente ex art. 10 c9 D.P.R. n. 633 del 1972:", body.Vat + " Euro"},
		{"Totale lordo da pagare:", body.Gross + " Euro"},
	}

	middle := []string{
		"L’originale della fattura elettronica, emessa in conformità a quanto disposto dalla Legge n. 205 del 27 dicembre 2017 (Legge di Bilancio), è disponibile nella Sua area riservata del sito web dell’Agenzia delle Entrate nella sezione “Fatture e Corrispettivi”)",
	}

	end := []string{
		"Scadenza pagamento: " + body.PayDate,
		"Condizioni di Pagamento: PAGATA",
	}

	titleCell := twoColParser(title, true)
	dataCell := threeColParser(data)
	middleCell := twoColParser(middle, false)
	endCell := twoColParser(end, false)

	pg.engine.DrawTable(titleCell)
	pg.engine.NewLine(10)
	pg.engine.DrawTable(dataCell)
	pg.engine.NewLine(6)
	pg.engine.DrawTable(middleCell)
	pg.engine.NewLine(6)
	pg.engine.DrawTable(endCell)

}

/*
func (lag *LifeAddendumGenerator) WoptaPerTe() {
	first := domain.TableCell{
		Text:      "Wopta per te",
		Height:    7,
		Width:     57,
		FontStyle: constants.BoldFontStyle,
		FontColor: constants.PinkColor,
		FontSize:  17,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.RightAlign,
		Border:    "",
	}
	second := domain.TableCell{
		Text:      "Vita",
		Height:    7,
		Width:     28,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		FontSize:  15,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.RightAlign,
		Border:    "",
	}

	origY := lag.engine.GetY()
	lag.engine.SetY(origY - 20)
	lag.engine.WriteText(first)
	lag.engine.SetY(lag.engine.GetY() - 1)
	lag.engine.WriteText(second)
	lag.engine.SetY(origY)
}

func (lag *LifeAddendumGenerator) contract() {
	const (
		firstColumnWidth  = 140
		secondColumnWidth = 50
	)

	contractDTO := lag.dto.Contract

	dataParser := func(rows []string) []domain.TableCell {
		result := make([]domain.TableCell, 0, len(rows))

		for index, row := range rows {
			fontStyle := constants.RegularFontStyle
			if index == 0 {
				fontStyle = constants.BoldFontStyle
			}
			result = append(result, domain.TableCell{

				Text:      row,
				Height:    constants.CellHeight,
				Width:     firstColumnWidth,
				FontStyle: fontStyle,
				FontColor: constants.BlackColor,
				FontSize:  constants.MediumFontSize,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "",
			})
		}
		return result
	}

	data := []string{
		contractDTO.CodeHeading,
		"Numero: " + contractDTO.Code,
		"Decorre dal: " + contractDTO.StartDate + " ore 24:00",
		"Scade il: " + contractDTO.EndDate + " ore 24:00",
		"Non si rinnova a scadenza",
		"Produttore: " + contractDTO.Producer,
	}

	text := dataParser(data)
	for _, v := range text {
		lag.engine.WriteText(v)
		//lag.engine.NewLine(2)
	}
}

func (lag *LifeAddendumGenerator) declarations() {
	first := domain.TableCell{
		Text:      "Dichiarazione di Variazione dati anagrafici Contraente-Assicurato-Beneficiario",
		Height:    3.5,
		Width:     190,
		FontStyle: constants.BoldFontStyle,
		FontColor: constants.PinkColor,
		FontSize:  constants.RegularFontSize,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "",
	}
	second := domain.TableCell{
		Text:      "Come da richiesta sono state trasmesse all’assicuratore AXA France Vie S.A. – Rappresentanza Generale per l’Italia le seguenti variazioni Anagrafiche di Polizza:",
		Height:    3.5,
		Width:     190,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		FontSize:  constants.RegularFontSize,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "",
	}
	lag.engine.WriteText(first)
	lag.engine.NewLine(3)
	lag.engine.WriteText(second)
}

func (lag *LifeAddendumGenerator) contractor() {
	cDTO := lag.dto.Contractor
	checked := " "
	var rows1 [][]string
	var rows2 [][]string
	var domTxt [][]string
	if cDTO.FiscalOrVatCode != constants.EmptyField {
		checked = "X"
		rows1 = [][]string{
			{"Cognome e Nome ", cDTO.Surname + " " + cDTO.Name, "Cod. Fisc: ", cDTO.FiscalOrVatCode},
			{"Residente in ", cDTO.StreetName + " " + cDTO.StreetNumber + " " + cDTO.City + " (" + cDTO.Province + ")", "Data nascita: ", cDTO.BirthDate},
		}
		rows2 = [][]string{
			{"Mail ", cDTO.Mail, "Telefono: ", cDTO.Phone},
		}
		domTxt = [][]string{
			{"Domicilio ", cDTO.DomStreetName + " " + cDTO.DomStreetNumber + " " + cDTO.DomCity + " (" + cDTO.DomProvince + ")"},
		}
	} else {
		rows1 = [][]string{
			{"Cognome e Nome ", " ", "Cod. Fisc: ", " "},
			{"Residente in ", " ", "Data nascita: ", " "},
		}
		rows2 = [][]string{
			{"Mail ", " ", "Telefono: ", " "},
		}
		domTxt = [][]string{
			{"Domicilio ", " "},
		}
	}

	titleT := []domain.TableCell{
		{
			Text:      checked,
			Height:    4.5,
			Width:     4.5,
			FontSize:  constants.LargeFontSize,
			FontStyle: constants.BoldFontStyle,
			FontColor: constants.BlackColor,
			Fill:      false,
			FillColor: domain.Color{},
			Align:     constants.CenterAlign,
			Border:    "1",
		},
		{
			Text:      "  Dati Contraente",
			Height:    4.5,
			Width:     190,
			FontStyle: constants.BoldFontStyle,
			FontColor: constants.PinkColor,
			FontSize:  constants.RegularFontSize,
			Fill:      false,
			FillColor: domain.Color{},
			Align:     constants.LeftAlign,
			Border:    "",
		},
	}
	title := make([][]domain.TableCell, 0)
	title = append(title, titleT)

	const (
		firstColumnWidth  = 35
		secondColumnWidth = 95
		thirdColumnWidth  = 25
		fourthColumnWidth = 35
	)
	parser := func(rows [][]string) [][]domain.TableCell {
		result := make([][]domain.TableCell, 0, len(rows))

		for _, row := range rows {

			result = append(result, []domain.TableCell{
				{
					Text:      row[0],
					Height:    constants.CellHeight,
					Width:     firstColumnWidth,
					FontStyle: constants.BoldFontStyle,
					FontColor: constants.BlackColor,
					FontSize:  constants.RegularFontSize,
					Fill:      false,
					FillColor: domain.Color{},
					Align:     constants.LeftAlign,
					Border:    "",
				},
				{
					Text:      row[1],
					Height:    constants.CellHeight,
					Width:     secondColumnWidth,
					FontStyle: constants.RegularFontStyle,
					FontColor: constants.BlackColor,
					FontSize:  constants.RegularFontSize,
					Fill:      false,
					FillColor: domain.Color{},
					Align:     constants.LeftAlign,
					Border:    "B",
				},
				{
					Text:      row[2],
					Height:    constants.CellHeight,
					Width:     thirdColumnWidth,
					FontStyle: constants.BoldFontStyle,
					FontColor: constants.BlackColor,
					FontSize:  constants.RegularFontSize,
					Fill:      false,
					FillColor: domain.Color{},
					Align:     constants.LeftAlign,
					Border:    "",
				},
				{
					Text:      row[3],
					Height:    constants.CellHeight,
					Width:     fourthColumnWidth,
					FontStyle: constants.RegularFontStyle,
					FontColor: constants.BlackColor,
					FontSize:  constants.RegularFontSize,
					Fill:      false,
					FillColor: domain.Color{},
					Align:     constants.LeftAlign,
					Border:    "B",
				},
			})
		}
		return result
	}

	table1 := parser(rows1)
	table2 := parser(rows2)

	const (
		domFirstColumnWidth  = 35
		domSecondColumnWidth = 155
	)
	domParser := func(rows [][]string) [][]domain.TableCell {
		result := make([][]domain.TableCell, 0, len(rows))

		for _, row := range rows {

			result = append(result, []domain.TableCell{
				{
					Text:      row[0],
					Height:    constants.CellHeight,
					Width:     domFirstColumnWidth,
					FontStyle: constants.BoldFontStyle,
					FontColor: constants.BlackColor,
					FontSize:  constants.RegularFontSize,
					Fill:      false,
					FillColor: domain.Color{},
					Align:     constants.LeftAlign,
					Border:    "",
				},
				{
					Text:      row[1],
					Height:    constants.CellHeight,
					Width:     domSecondColumnWidth,
					FontStyle: constants.RegularFontStyle,
					FontColor: constants.BlackColor,
					FontSize:  constants.RegularFontSize,
					Fill:      false,
					FillColor: domain.Color{},
					Align:     constants.LeftAlign,
					Border:    "B",
				},
			})
		}
		return result
	}

	dom := domParser(domTxt)

	lag.engine.DrawTable(title)
	lag.engine.NewLine(2)
	lag.engine.DrawLine(10, lag.engine.GetY(), 200, lag.engine.GetY(), 0.25, constants.BlackColor)
	lag.engine.NewLine(2)
	lag.engine.DrawTable(table1)
	lag.engine.DrawTable(dom)
	lag.engine.DrawTable(table2)
}

func (lag *LifeAddendumGenerator) insured() {
	iDTO := lag.dto.Insured

	checked := " "
	var rows1 [][]string
	var rows2 [][]string
	var domTxt [][]string
	if iDTO.FiscalOrVatCode != constants.EmptyField {
		checked = "X"
		rows1 = [][]string{
			{"Cognome e Nome ", iDTO.Surname + " " + iDTO.Name, "Cod. Fisc: ", iDTO.FiscalOrVatCode},
			{"Residente in ", iDTO.StreetName + " " + iDTO.StreetNumber + " " + iDTO.City + " (" + iDTO.Province + ")", "Data nascita: ", iDTO.BirthDate},
		}
		rows2 = [][]string{
			{"Mail ", iDTO.Mail, "Telefono: ", iDTO.Phone},
		}
		domTxt = [][]string{
			{"Domicilio ", iDTO.DomStreetName + " " + iDTO.DomStreetNumber + " " + iDTO.DomCity + " (" + iDTO.DomProvince + ")"},
		}
	} else {
		rows1 = [][]string{
			{"Cognome e Nome ", " ", "Cod. Fisc: ", " "},
			{"Residente in ", " ", "Data nascita: ", " "},
		}
		rows2 = [][]string{
			{"Mail ", " ", "Telefono: ", " "},
		}
		domTxt = [][]string{
			{"Domicilio ", " "},
		}
	}

	titleT := []domain.TableCell{
		{
			Text:      checked,
			Height:    4.5,
			Width:     4.5,
			FontSize:  constants.LargeFontSize,
			FontStyle: constants.BoldFontStyle,
			FontColor: constants.BlackColor,
			Fill:      false,
			FillColor: domain.Color{},
			Align:     constants.CenterAlign,
			Border:    "1",
		},
		{
			Text:      "  Dati Assicurato",
			Height:    4.5,
			Width:     190,
			FontStyle: constants.BoldFontStyle,
			FontColor: constants.PinkColor,
			FontSize:  constants.RegularFontSize,
			Fill:      false,
			FillColor: domain.Color{},
			Align:     constants.LeftAlign,
			Border:    "",
		},
	}
	title := make([][]domain.TableCell, 0)
	title = append(title, titleT)

	const (
		firstColumnWidth  = 35
		secondColumnWidth = 95
		thirdColumnWidth  = 25
		fourthColumnWidth = 35
	)
	parser := func(rows [][]string) [][]domain.TableCell {
		result := make([][]domain.TableCell, 0, len(rows))

		for _, row := range rows {

			result = append(result, []domain.TableCell{
				{
					Text:      row[0],
					Height:    constants.CellHeight,
					Width:     firstColumnWidth,
					FontStyle: constants.BoldFontStyle,
					FontColor: constants.BlackColor,
					FontSize:  constants.RegularFontSize,
					Fill:      false,
					FillColor: domain.Color{},
					Align:     constants.LeftAlign,
					Border:    "",
				},
				{
					Text:      row[1],
					Height:    constants.CellHeight,
					Width:     secondColumnWidth,
					FontStyle: constants.RegularFontStyle,
					FontColor: constants.BlackColor,
					FontSize:  constants.RegularFontSize,
					Fill:      false,
					FillColor: domain.Color{},
					Align:     constants.LeftAlign,
					Border:    "B",
				},
				{
					Text:      row[2],
					Height:    constants.CellHeight,
					Width:     thirdColumnWidth,
					FontStyle: constants.BoldFontStyle,
					FontColor: constants.BlackColor,
					FontSize:  constants.RegularFontSize,
					Fill:      false,
					FillColor: domain.Color{},
					Align:     constants.LeftAlign,
					Border:    "",
				},
				{
					Text:      row[3],
					Height:    constants.CellHeight,
					Width:     fourthColumnWidth,
					FontStyle: constants.RegularFontStyle,
					FontColor: constants.BlackColor,
					FontSize:  constants.RegularFontSize,
					Fill:      false,
					FillColor: domain.Color{},
					Align:     constants.LeftAlign,
					Border:    "B",
				},
			})
		}
		return result
	}

	table1 := parser(rows1)

	table2 := parser(rows2)

	const (
		domFirstColumnWidth  = 35
		domSecondColumnWidth = 155
	)
	domParser := func(rows [][]string) [][]domain.TableCell {
		result := make([][]domain.TableCell, 0, len(rows))

		for _, row := range rows {

			result = append(result, []domain.TableCell{
				{
					Text:      row[0],
					Height:    constants.CellHeight,
					Width:     domFirstColumnWidth,
					FontStyle: constants.BoldFontStyle,
					FontColor: constants.BlackColor,
					FontSize:  constants.RegularFontSize,
					Fill:      false,
					FillColor: domain.Color{},
					Align:     constants.LeftAlign,
					Border:    "",
				},
				{
					Text:      row[1],
					Height:    constants.CellHeight,
					Width:     domSecondColumnWidth,
					FontStyle: constants.RegularFontStyle,
					FontColor: constants.BlackColor,
					FontSize:  constants.RegularFontSize,
					Fill:      false,
					FillColor: domain.Color{},
					Align:     constants.LeftAlign,
					Border:    "B",
				},
			})
		}
		return result
	}

	dom := domParser(domTxt)

	lag.engine.DrawTable(title)
	lag.engine.NewLine(2)
	lag.engine.DrawLine(10, lag.engine.GetY(), 200, lag.engine.GetY(), 0.25, constants.BlackColor)
	lag.engine.NewLine(2)
	lag.engine.DrawTable(table1)
	lag.engine.DrawTable(dom)
	lag.engine.DrawTable(table2)
}

func (lag *LifeAddendumGenerator) beneficiaries() {
	bDTO := lag.dto.Beneficiaries

	checked := " "
	var rows [][]string
	var relTxt [][]string
	if bDTO != nil && len(*bDTO) != 0 {
		for _, v := range *bDTO {
			if v.FiscalOrVatCode != constants.EmptyField {
				checked = "X"
			}
		}

	}

	titleT := []domain.TableCell{
		{
			Text:      checked,
			Height:    4.5,
			Width:     4.5,
			FontSize:  constants.LargeFontSize,
			FontStyle: constants.BoldFontStyle,
			FontColor: constants.BlackColor,
			Fill:      false,
			FillColor: domain.Color{},
			Align:     constants.CenterAlign,
			Border:    "1",
		},
		{
			Text:      "  Dati Beneficiari",
			Height:    4.5,
			Width:     190,
			FontStyle: constants.BoldFontStyle,
			FontColor: constants.PinkColor,
			FontSize:  constants.RegularFontSize,
			Fill:      false,
			FillColor: domain.Color{},
			Align:     constants.LeftAlign,
			Border:    "",
		},
	}
	title := make([][]domain.TableCell, 0)
	title = append(title, titleT)

	const (
		firstColumnWidth  = 35
		secondColumnWidth = 95
		thirdColumnWidth  = 25
		fourthColumnWidth = 35
	)
	parser := func(rows [][]string) [][]domain.TableCell {
		result := make([][]domain.TableCell, 0, len(rows))

		for _, row := range rows {

			result = append(result, []domain.TableCell{
				{
					Text:      row[0],
					Height:    constants.CellHeight,
					Width:     firstColumnWidth,
					FontStyle: constants.BoldFontStyle,
					FontColor: constants.BlackColor,
					FontSize:  constants.RegularFontSize,
					Fill:      false,
					FillColor: domain.Color{},
					Align:     constants.LeftAlign,
					Border:    "",
				},
				{
					Text:      row[1],
					Height:    constants.CellHeight,
					Width:     secondColumnWidth,
					FontStyle: constants.RegularFontStyle,
					FontColor: constants.BlackColor,
					FontSize:  constants.RegularFontSize,
					Fill:      false,
					FillColor: domain.Color{},
					Align:     constants.LeftAlign,
					Border:    "B",
				},
				{
					Text:      row[2],
					Height:    constants.CellHeight,
					Width:     thirdColumnWidth,
					FontStyle: constants.BoldFontStyle,
					FontColor: constants.BlackColor,
					FontSize:  constants.RegularFontSize,
					Fill:      false,
					FillColor: domain.Color{},
					Align:     constants.LeftAlign,
					Border:    "",
				},
				{
					Text:      row[3],
					Height:    constants.CellHeight,
					Width:     fourthColumnWidth,
					FontStyle: constants.RegularFontStyle,
					FontColor: constants.BlackColor,
					FontSize:  constants.RegularFontSize,
					Fill:      false,
					FillColor: domain.Color{},
					Align:     constants.LeftAlign,
					Border:    "B",
				},
			})
		}
		return result
	}

	const (
		domFirstColumnWidth  = 35
		domSecondColumnWidth = 155
	)
	relParser := func(rows [][]string) [][]domain.TableCell {
		result := make([][]domain.TableCell, 0, len(rows))

		for _, row := range rows {

			result = append(result, []domain.TableCell{
				{
					Text:      row[0],
					Height:    constants.CellHeight,
					Width:     domFirstColumnWidth,
					FontStyle: constants.BoldFontStyle,
					FontColor: constants.BlackColor,
					FontSize:  constants.RegularFontSize,
					Fill:      false,
					FillColor: domain.Color{},
					Align:     constants.LeftAlign,
					Border:    "",
				},
				{
					Text:      row[1],
					Height:    constants.CellHeight,
					Width:     domSecondColumnWidth,
					FontStyle: constants.RegularFontStyle,
					FontColor: constants.BlackColor,
					FontSize:  constants.RegularFontSize,
					Fill:      false,
					FillColor: domain.Color{},
					Align:     constants.LeftAlign,
					Border:    "B",
				},
			})
		}
		return result
	}

	lag.engine.DrawTable(title)
	lag.engine.NewLine(2)
	lag.engine.DrawLine(10, lag.engine.GetY(), 200, lag.engine.GetY(), 0.25, constants.BlackColor)
	lag.engine.NewLine(2)
	for i := 0; i < 2; i++ {
		if checked == "X" {
			rows = [][]string{
				{"Cognome e Nome ", (*bDTO)[i].Surname + " " + (*bDTO)[i].Name, "Cod. Fisc: ", (*bDTO)[i].FiscalOrVatCode},
				{"Residente in ", (*bDTO)[i].StreetName + " " + (*bDTO)[i].StreetNumber + " " + (*bDTO)[i].City + " (" + (*bDTO)[i].Province + ")", "Data nascita: ", (*bDTO)[i].BirthDate},
				{"Mail ", (*bDTO)[i].Mail, "Telefono ", (*bDTO)[i].Phone},
			}
			relTxt = [][]string{
				{"Relazione con Assicurato ", (*bDTO)[i].Relation},
			}
		} else {
			rows = [][]string{
				{"Cognome e Nome ", " ", "Cod. Fisc: ", " "},
				{"Residente in ", " ", "Data nascita: ", " "},
				{"Mail ", " ", "Telefono ", " "},
			}
			relTxt = [][]string{
				{"Relazione con Assicurato ", " "},
			}
		}

		table := parser(rows)
		rel := relParser(relTxt)
		lag.engine.DrawTable(table)
		lag.engine.DrawTable(rel)
		lag.engine.NewLine(2)
		conf := " "
		if checked == "X" {
			conf = "No"
			if (*bDTO)[i].Contactable {
				conf = "Sì"
			}
		}

		cons := domain.TableCell{
			Text:      "Consenso ad invio comunicazioni da parte della compagnia ai beneficiari, prima dell'evento decesso: " + conf,
			Height:    3.5,
			Width:     190,
			FontStyle: constants.RegularFontStyle,
			FontColor: constants.BlackColor,
			FontSize:  constants.MediumFontSize,
			Fill:      false,
			FillColor: domain.Color{},
			Align:     constants.LeftAlign,
			Border:    "",
		}
		lag.engine.WriteText(cons)
		lag.engine.DrawLine(10, lag.engine.GetY(), 200, lag.engine.GetY(), 0.25, constants.BlackColor)
		if i == 0 {
			lag.engine.NewLine(4)
		}
	}
}

func (lag *LifeAddendumGenerator) beneficiaryReference() {
	brDTO := lag.dto.BeneficiaryReference
	checked := " "
	var rows [][]string
	if brDTO.FiscalOrVatCode != constants.EmptyField {
		checked = "X"
		rows = [][]string{
			{"Cognome e Nome ", brDTO.Surname + " " + brDTO.Name, "Cod. Fisc: ", brDTO.FiscalOrVatCode},
			{"Residente in ", brDTO.StreetName + " " + brDTO.StreetNumber + " " + brDTO.City + " (" + brDTO.Province + ")", "Data nascita: ", brDTO.BirthDate},
			{"Mail ", brDTO.Mail, "Telefono: ", brDTO.Phone},
		}
	} else {
		rows = [][]string{
			{"Cognome e Nome ", " ", "Cod. Fisc: ", " "},
			{"Residente in ", " ", "Data nascita: ", " "},
			{"Mail ", " ", "Telefono: ", " "},
		}
	}
	titleT := []domain.TableCell{
		{
			Text:      checked,
			Height:    4.5,
			Width:     4.5,
			FontSize:  constants.LargeFontSize,
			FontStyle: constants.BoldFontStyle,
			FontColor: constants.BlackColor,
			Fill:      false,
			FillColor: domain.Color{},
			Align:     constants.CenterAlign,
			Border:    "1",
		},
		{
			Text:      "  Referente terzo",
			Height:    4.5,
			Width:     190,
			FontStyle: constants.BoldFontStyle,
			FontColor: constants.PinkColor,
			FontSize:  constants.RegularFontSize,
			Fill:      false,
			FillColor: domain.Color{},
			Align:     constants.LeftAlign,
			Border:    "",
		},
	}
	title := make([][]domain.TableCell, 0)
	title = append(title, titleT)

	const (
		firstColumnWidth  = 35
		secondColumnWidth = 95
		thirdColumnWidth  = 25
		fourthColumnWidth = 35
	)
	parser := func(rows [][]string) [][]domain.TableCell {
		result := make([][]domain.TableCell, 0, len(rows))

		for _, row := range rows {

			result = append(result, []domain.TableCell{
				{
					Text:      row[0],
					Height:    constants.CellHeight,
					Width:     firstColumnWidth,
					FontStyle: constants.BoldFontStyle,
					FontColor: constants.BlackColor,
					FontSize:  constants.RegularFontSize,
					Fill:      false,
					FillColor: domain.Color{},
					Align:     constants.LeftAlign,
					Border:    "",
				},
				{
					Text:      row[1],
					Height:    constants.CellHeight,
					Width:     secondColumnWidth,
					FontStyle: constants.RegularFontStyle,
					FontColor: constants.BlackColor,
					FontSize:  constants.RegularFontSize,
					Fill:      false,
					FillColor: domain.Color{},
					Align:     constants.LeftAlign,
					Border:    "B",
				},
				{
					Text:      row[2],
					Height:    constants.CellHeight,
					Width:     thirdColumnWidth,
					FontStyle: constants.BoldFontStyle,
					FontColor: constants.BlackColor,
					FontSize:  constants.RegularFontSize,
					Fill:      false,
					FillColor: domain.Color{},
					Align:     constants.LeftAlign,
					Border:    "",
				},
				{
					Text:      row[3],
					Height:    constants.CellHeight,
					Width:     fourthColumnWidth,
					FontStyle: constants.RegularFontStyle,
					FontColor: constants.BlackColor,
					FontSize:  constants.RegularFontSize,
					Fill:      false,
					FillColor: domain.Color{},
					Align:     constants.LeftAlign,
					Border:    "B",
				},
			})
		}
		return result
	}

	table := parser(rows)

	lag.engine.DrawTable(title)
	lag.engine.NewLine(2)
	lag.engine.DrawLine(10, lag.engine.GetY(), 200, lag.engine.GetY(), 0.25, constants.BlackColor)
	lag.engine.NewLine(2)
	lag.engine.DrawTable(table)
}
*/
