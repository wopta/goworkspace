package addendum

import (
	"time"

	"gitlab.dev.wopta.it/goworkspace/document/internal/constants"
	"gitlab.dev.wopta.it/goworkspace/document/internal/domain"
	"gitlab.dev.wopta.it/goworkspace/document/internal/dto"
	"gitlab.dev.wopta.it/goworkspace/document/internal/engine"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
)

type CatnatAddendumGenerator struct {
	dto *dto.AddendumCatnatDTO
	*baseGenerator
}

func NewCatnatAddendumGenerator(engine *engine.Fpdf, policy *models.Policy) *CatnatAddendumGenerator {
	now := time.Now()
	dto := dto.NewCatnatAddendumDto()
	dto.FromPolicy(policy, now)
	return &CatnatAddendumGenerator{
		baseGenerator: &baseGenerator{
			engine: engine,
			now:    now,
			policy: policy,
		},
		dto: dto,
	}
}

func (lag *CatnatAddendumGenerator) Generate() {
	lag.mainHeader()

	lag.engine.NewPage()

	lag.contract()

	lag.engine.NewLine(6)

	lag.declarations()

	lag.engine.NewLine(6)

	lag.contractor()

	lag.engine.NewLine(6)

	lag.signer()

	lag.engine.NewLine(6)

	lag.contractorSignature()

	lag.woptaFooter()

}

func (lag *CatnatAddendumGenerator) mainHeader() {
	lag.engine.SetHeader(func() {
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
			Text:      "Catastrofali Azienda",
			Height:    7,
			Width:     60,
			FontStyle: constants.RegularFontStyle,
			FontColor: constants.BlackColor,
			FontSize:  15,
			Fill:      false,
			FillColor: domain.Color{},
			Align:     constants.RightAlign,
			Border:    "",
		}

		lag.engine.InsertImage(lib.GetAssetPathByEnvV2()+"logo_catnat.png", 12, 6.5, 12, 12)
		lag.engine.InsertImage(lib.GetAssetPathByEnvV2()+"logo_wopta.png", 160, 6.5, 35, 10)
		lag.engine.NewLine(7)
		origY := lag.engine.GetY()
		lag.engine.SetY(origY - 15)
		lag.engine.WriteText(first)
		lag.engine.SetY(lag.engine.GetY() - 1)
		lag.engine.SetX(15)
		lag.engine.WriteText(second)
		lag.engine.SetY(origY)
	})
}

func (lag *CatnatAddendumGenerator) contract() {
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

func (lag *CatnatAddendumGenerator) declarations() {
	lag.engine.WriteText(domain.TableCell{
		Text:      "Dichiarazione di Variazione dati Anagrafici Contraente-Firmatario",
		Height:    constants.CellHeight,
		Width:     constants.FullPageWidth,
		FontStyle: constants.BoldFontStyle,
		FontColor: constants.PinkColor,
	})
	lag.engine.NewLine(1)
	lag.engine.WriteText(domain.TableCell{
		Text:      "Le modifiche non sono attive in assenza di firma da parte del Contraente, con allegata copia documento di riconoscimento (carta identità, passaporto, patente, in corso di validità alla data di firma).",
		Height:    constants.CellHeight,
		Width:     constants.FullPageWidth,
		FontStyle: constants.BoldFontStyle,
	})
}

func (lag *CatnatAddendumGenerator) contractor() {
	cDTO := lag.dto.Contractor
	checked := " "
	var rows1 [][]string
	var rows2 [][]string
	if cDTO.Name != "" || cDTO.CompanyName != "" {
		checked = "X"
	}
	rows1 = [][]string{
		{"Tipo Soggetto:", cDTO.Type + " ", " ", " "},
		{"Cognome e Nome:", cDTO.Surname + " " + cDTO.Name, "Cod. Fisc: ", cDTO.FiscalCode + " "},
		{"Denominazione sociale:", cDTO.CompanyName, "Partita Iva:", cDTO.VatCode},
		{"Indirizzo Sede Legale:", cDTO.Address, "Codice Ateco:", cDTO.Ateco},
		{"Mail:", cDTO.Mail, "Telefono: ", cDTO.Phone},
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

	lag.engine.DrawTable(title)
	lag.engine.NewLine(2)
	lag.engine.DrawLine(10, lag.engine.GetY(), 200, lag.engine.GetY(), 0.25, constants.BlackColor)
	lag.engine.NewLine(2)
	lag.engine.DrawTable(table1)
	lag.engine.DrawTable(table2)
}

func (lag *CatnatAddendumGenerator) signer() {
	iDTO := lag.dto.Signer

	checked := " "
	var rows1 [][]string
	var rows2 [][]string
	if iDTO.FiscalCode != "" {
		checked = "X"
	}
	rows1 = [][]string{
		{"Cognome e Nome:", iDTO.Surname + " " + iDTO.Name, "Cod. Fisc:", iDTO.FiscalCode + " "},
		{"Luogo di Nascita:", iDTO.GetBirthAddress(), "Data nascita:", iDTO.BirthDate},
		{"Indirizzo Residenza:", iDTO.GetResidenceAddress(), "Sesso:", iDTO.Gender},
		{"Mail:", iDTO.Mail, "Telefono:", iDTO.Phone},
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
			Text:      "  Dati Firmatario",
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

	lag.engine.DrawTable(title)
	lag.engine.NewLine(2)
	lag.engine.DrawLine(10, lag.engine.GetY(), 200, lag.engine.GetY(), 0.25, constants.BlackColor)
	lag.engine.NewLine(2)
	lag.engine.DrawTable(table1)
	lag.engine.DrawTable(table2)
}

func (lag *CatnatAddendumGenerator) contractorSignature() {
	var (
		cellHeight   = 5
		colWidth     = float64(50)
		spacingWidth = constants.FullPageWidth - (2 * colWidth)
	)
	row := []domain.TableCell{
		{
			Text:   lag.dto.Contract.IssueDate,
			Height: float64(cellHeight),
			Width:  colWidth,
		},
		{
			Text:   " ",
			Height: float64(cellHeight),
			Width:  spacingWidth,
		},
		{
			Text:   "Firma Contraente",
			Height: float64(cellHeight),
			Width:  colWidth,
			Align:  constants.CenterAlign,
			Border: constants.BorderTop,
		},
	}
	table := make([][]domain.TableCell, 0, 1)
	table = append(table, row)

	lag.engine.SetY(-40)

	lag.engine.DrawTable(table)
}
