package catnat

import (
	"gitlab.dev.wopta.it/goworkspace/document/internal/constants"
	"gitlab.dev.wopta.it/goworkspace/document/internal/domain"
	"gitlab.dev.wopta.it/goworkspace/document/internal/dto"
	"gitlab.dev.wopta.it/goworkspace/document/internal/engine"
)

func AddBuildingInformation(e *engine.Fpdf, sede dto.BuildingCatnatDto, questions dto.QuestionsCatnatDto) {
	e.WriteText(e.GetTableCell("DATI SEDE DA ASSICURARE", constants.BoldFontStyle, constants.LargeFontSize))
	e.SetDrawColor(constants.PinkColor)
	e.NewLine(2)
	var (
		firstColumnWidth  float64 = 65
		secondColumnWidth float64 = 45
		thirdColumnWidth  float64 = 40
		fourthColumnWidth float64 = 40
	)
	parserTableInfoSede := func(lines [][]string) [][]domain.TableCell {
		result := make([][]domain.TableCell, 0, len(lines))
		for _, texts := range lines {
			result = append(result, []domain.TableCell{
				{
					Text:      texts[0],
					Height:    constants.CellHeight,
					Width:     firstColumnWidth,
					FontStyle: constants.BoldFontStyle,
					FontColor: constants.BlackColor,
					Fill:      false,
					FontSize:  constants.RegularFontSize,
					FillColor: domain.Color{},
					Align:     constants.LeftAlign,
					Border:    "T",
				},
				{
					Text:      texts[1],
					Height:    constants.CellHeight,
					FontStyle: constants.RegularFontStyle,
					Width:     secondColumnWidth,
					FontColor: constants.BlackColor,
					FontSize:  constants.RegularFontSize,
					Fill:      false,
					FillColor: domain.Color{},
					Align:     constants.LeftAlign,
					Border:    "T",
				},
				{
					Text:      texts[2],
					Height:    constants.CellHeight,
					FontStyle: constants.BoldFontStyle,
					FontColor: constants.BlackColor,
					Width:     thirdColumnWidth,
					FontSize:  constants.RegularFontSize,
					Fill:      false,
					FillColor: domain.Color{},
					Align:     constants.LeftAlign,
					Border:    "T",
				}, {
					Text:      texts[3],
					Height:    constants.CellHeight,
					FontStyle: constants.RegularFontStyle,
					FontColor: constants.BlackColor,
					Width:     fourthColumnWidth,
					FontSize:  constants.RegularFontSize,
					Fill:      false,
					FillColor: domain.Color{},
					Align:     constants.LeftAlign,
					Border:    "T",
				},
			})
		}

		return result
	}
	dataSede := [][]string{
		{"Indirizzo:", sede.Address, "Tipo Utilizzo Sede:", sede.Type},
		{"Anno di costruzione:", sede.BuildingYear, "Materiale di costruzione:", sede.BuildingMaterial},
		{"Numero di piani edificio oltre il piano terra:", sede.Floor, "", ""},
	}
	e.DrawTable(parserTableInfoSede(dataSede))
	e.SetX(e.GetX() - 3)
	e.SetY(e.GetY() - 3)
	firstColumnWidth = 170
	secondColumnWidth = 20

	parserTableQuestions := func(lines [][]string) [][]domain.TableCell {
		result := make([][]domain.TableCell, 0, len(lines))
		for _, texts := range lines {
			result = append(result, []domain.TableCell{
				{
					Text:      texts[0],
					Height:    constants.CellHeight,
					Width:     firstColumnWidth,
					FontColor: constants.BlackColor,
					Fill:      false,
					FontSize:  constants.RegularFontSize,
					FillColor: domain.Color{},
					Align:     constants.LeftAlign,
					Border:    "T",
				},
				{
					Text:      texts[1],
					Height:    constants.CellHeight,
					FontStyle: constants.RegularFontStyle,
					Width:     secondColumnWidth,
					FontColor: constants.BlackColor,
					FontSize:  constants.RegularFontSize,
					Fill:      false,
					FillColor: domain.Color{},
					Align:     constants.LeftAlign,
					Border:    "T",
				},
			})
		}

		return result
	}
	dataQuestion := [][]string{
		{"Il Fabbricato e il Terreno da assicurare sono già coperti per la Garanzia Terremoto?", questions.AlreadyEarthquake},
		{"Nel caso in cui il Fabbricato e il Terreno da assicurare posseggano già la Garanzia Terremoto, il Contraente intende acquistare la medesima copertura?", questions.WantEarthquake},
		{"Il Fabbricato e il Terreno da assicurare sono già coperti per la Garanzia Alluvione?", questions.AlreadyFlood},
		{"Nel caso in cui il Fabbricato e il Terreno da assicurare posseggano già la Garanzia Alluvione, il Contraente intende acquistare la medesima copertura?", questions.WantFlood},
	}
	e.DrawTable(parserTableQuestions(dataQuestion))
	e.DrawLine(10, e.GetY(), 200, e.GetY(), 0.25, constants.PinkColor)
}
