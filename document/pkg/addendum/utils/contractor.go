package utils

import (
	"strings"

	"gitlab.dev.wopta.it/goworkspace/document/internal/constants"
	"gitlab.dev.wopta.it/goworkspace/document/internal/domain"
	"gitlab.dev.wopta.it/goworkspace/document/internal/dto"
	"gitlab.dev.wopta.it/goworkspace/document/internal/engine"
)

func AddContractor(contractor dto.ContractorDTO, engine *engine.Fpdf) {
	checked := " "
	var rows1 [][]string
	var rows2 [][]string
	if contractor.Name != "" || contractor.CompanyName != "" {
		checked = "X"
	}
	rows1 = [][]string{
		{"Tipo Soggetto:", contractor.Type + " ", " ", " "},
		{"Cognome e Nome:", contractor.Surname + " " + contractor.Name, "Cod. Fisc: ", contractor.FiscalCode + " "},
		{"Denominazione sociale:", contractor.CompanyName, "Partita Iva:", contractor.VatCode},
		{"Indirizzo Sede Legale:", strings.ReplaceAll(contractor.Address, "\n", ""), "Codice Ateco:", contractor.Ateco},
		{"Mail:", contractor.Mail, "Telefono:", contractor.Phone},
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

	engine.DrawTable(title)
	engine.NewLine(2)
	engine.DrawLine(10, engine.GetY(), 200, engine.GetY(), 0.25, constants.BlackColor)
	engine.NewLine(2)
	engine.DrawTable(table1)
	engine.DrawTable(table2)
}
