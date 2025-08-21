package catnat

import (
	"gitlab.dev.wopta.it/goworkspace/document/internal/constants"
	"gitlab.dev.wopta.it/goworkspace/document/internal/domain"
	"gitlab.dev.wopta.it/goworkspace/document/internal/dto"
	"gitlab.dev.wopta.it/goworkspace/document/internal/engine"
)

func AddTableGuarantee(e *engine.Fpdf, guarantees dto.CatnatGuaranteeDTO) {
	const (
		widthColumn = 38
	)
	guranteeTable := func(lines [][]string) [][]domain.TableCell {
		result := make([][]domain.TableCell, 0, len(lines))
		var font domain.FontStyle
		for i, texts := range lines {
			if i == 0 {
				font = constants.BoldFontStyle
			} else {
				font = constants.RegularFontStyle
			}
			result = append(result, []domain.TableCell{
				{
					Text:      texts[0],
					Height:    constants.CellHeight,
					Width:     widthColumn,
					FontStyle: constants.BoldFontStyle,
					FontColor: constants.BlackColor,
					Fill:      false,
					FontSize:  constants.RegularFontSize,
					FillColor: domain.Color{},
					Align:     constants.LeftAlign,
					Border:    "1",
				},
				{
					Text:      texts[1],
					Height:    constants.CellHeight,
					FontStyle: font,
					Width:     widthColumn,
					FontColor: constants.BlackColor,
					FontSize:  constants.RegularFontSize,
					Fill:      false,
					FillColor: domain.Color{},
					Align:     constants.LeftAlign,
					Border:    "1",
				},
				{
					Text:      texts[2],
					Height:    constants.CellHeight,
					FontStyle: font,
					FontColor: constants.BlackColor,
					Width:     widthColumn,
					FontSize:  constants.RegularFontSize,
					Fill:      false,
					FillColor: domain.Color{},
					Align:     constants.LeftAlign,
					Border:    "1",
				},
				{
					Text:      texts[3],
					Height:    constants.CellHeight,
					FontStyle: font,
					FontColor: constants.BlackColor,
					FontSize:  constants.RegularFontSize,
					Fill:      false,
					Width:     widthColumn,
					FillColor: domain.Color{},
					Align:     constants.LeftAlign,
					Border:    "1",
				},
				{
					Text:      texts[4],
					Height:    constants.CellHeight,
					FontStyle: font,
					FontColor: constants.BlackColor,
					Width:     widthColumn,
					FontSize:  constants.RegularFontSize,
					Fill:      false,
					FillColor: domain.Color{},
					Align:     constants.LeftAlign,
					Border:    "1",
				},
			})
		}

		return result
	}
	guaranteeData := [][]string{
		{"Garanzie", "Somma Assicurata Fabricato €", "Somma Assicurata Contenuto €", "Somma Assicurata Merci €", "Importo Annuo €"},
		{"Terremoto", guarantees.EarthquakeGuarantee.Building, guarantees.EarthquakeGuarantee.Content, guarantees.EarthquakeGuarantee.Stock, guarantees.EarthquakeGuarantee.Total},
		{"Alluvione", guarantees.FloodGuarantee.Building, guarantees.FloodGuarantee.Content, guarantees.FloodGuarantee.Stock, guarantees.FloodGuarantee.Total},
		{"Frana", guarantees.LandslideGuarantee.Building, guarantees.LandslideGuarantee.Content, guarantees.LandslideGuarantee.Stock, guarantees.LandslideGuarantee.Total},
	}
	e.DrawTable(guranteeTable(guaranteeData))
}
