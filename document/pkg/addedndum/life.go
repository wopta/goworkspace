package addedndum

import (
	"time"

	"github.com/wopta/goworkspace/document/internal/constants"
	"github.com/wopta/goworkspace/document/internal/domain"
	"github.com/wopta/goworkspace/document/internal/dto"
	"github.com/wopta/goworkspace/document/internal/engine"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

type LifeAddendumGenerator struct {
	*baseGenerator
	dto *dto.BeneficiariesDTO
}

func NewLifeAddendumGenerator(engine *engine.Fpdf, policy *models.Policy, node *models.NetworkNode,
	product models.Product) *LifeAddendumGenerator {
	LifeAddendumDTO := dto.NewBeneficiariesDto()
	LifeAddendumDTO.FromPolicy(*policy, product)
	return &LifeAddendumGenerator{
		baseGenerator: &baseGenerator{
			engine:      engine,
			now:         time.Now(),
			signatureID: 0,
			networkNode: node,
			policy:      policy,
		},
		dto: LifeAddendumDTO,
	}
}

func (lag *LifeAddendumGenerator) Contract() ([]byte, error) {
	lag.mainHeader()

	lag.engine.NewPage()

	lag.engine.NewLine(6)

	lag.declarations()

	lag.engine.NewLine(6)

	lag.contractor()

	lag.engine.NewLine(6)

	lag.insured()

	lag.engine.NewLine(6)

	lag.beneficiaries()

	lag.engine.NewLine(6)

	lag.beneficiaryReference()

	lag.engine.NewLine(6)

	lag.woptaFooter()

	return lag.engine.RawDoc()
}

func (lag *LifeAddendumGenerator) mainHeader() {
	const (
		firstColumnWidth  = 140
		secondColumnWidth = 50
	)

	contractDTO := lag.dto.Contract

	parser := func(rows [][]string) [][]domain.TableCell {
		result := make([][]domain.TableCell, 0, len(rows))

		for index, row := range rows {
			fontStyle := constants.RegularFontStyle
			if index == 0 {
				fontStyle = constants.BoldFontStyle
			}
			result = append(result, []domain.TableCell{
				{
					Text:      row[0],
					Height:    constants.CellHeight,
					Width:     firstColumnWidth,
					FontStyle: fontStyle,
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

	rows := [][]string{
		{contractDTO.CodeHeading, ""},
		{"Numero: " + contractDTO.Code, ""},
		{"Decorre dal: " + contractDTO.StartDate + " ore 24:00", ""},
		{"Scade il: " + contractDTO.EndDate + " ore 24:00", ""},
		{"Non si rinnova a scadenza", ""},
		{"Produttore: " + contractDTO.Producer, ""},
	}

	table := parser(rows)

	lag.engine.SetHeader(func() {
		lag.engine.InsertImage(lib.GetAssetPathByEnvV2()+"logo_vita.png", 12, 6.5, 12, 12)
		//lag.engine.DrawLine(102, 6.25, 102, 15, 0.25, constants.BlackColor)
		lag.engine.InsertImage(lib.GetAssetPathByEnvV2()+"logo_wopta.png", 160, 6.5, 35, 10)
		lag.engine.NewLine(7)
		lag.engine.DrawTable(table)

	})
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

	title := domain.TableCell{
		Text:      "Dati Contraente",
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

	const (
		firstColumnWidth  = 33
		secondColumnWidth = 70
		thirdColumnWidth  = 25
		fourthColumnWidth = 60
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

	rows := [][]string{
		{"Cognome e Nome ", cDTO.Surname + " " + cDTO.Name, "Cod. Fisc: ", cDTO.FiscalCode},
		{"Residente in ", cDTO.StreetName + " " + cDTO.StreetNumber + " " + cDTO.City, "Data nascita: ", cDTO.BirthDate},
		{"Mail ", cDTO.Mail, "Telefono: ", cDTO.Phone},
	}
	table := parser(rows)

	lag.engine.WriteText(title)
	lag.engine.NewLine(2)
	lag.engine.DrawLine(10, lag.engine.GetY(), 200, lag.engine.GetY(), 0.25, constants.BlackColor)
	lag.engine.NewLine(2)
	lag.engine.DrawTable(table)
}

func (lag *LifeAddendumGenerator) insured() {
	iDTO := lag.dto.Insured
	title := domain.TableCell{
		Text:      "Dati Assicurato",
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

	const (
		firstColumnWidth  = 33
		secondColumnWidth = 70
		thirdColumnWidth  = 25
		fourthColumnWidth = 60
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

	rows := [][]string{
		{"Cognome e Nome ", iDTO.Surname + " " + iDTO.Name, "Cod. Fisc: ", iDTO.FiscalCode},
		{"Residente in ", iDTO.StreetName + " " + iDTO.StreetNumber + " " + iDTO.City, "Data nascita: ", iDTO.BirthDate},
		{"Mail ", iDTO.Mail, "Telefono: ", iDTO.Phone},
	}
	table := parser(rows)

	lag.engine.WriteText(title)
	lag.engine.NewLine(2)
	lag.engine.DrawLine(10, lag.engine.GetY(), 200, lag.engine.GetY(), 0.25, constants.BlackColor)
	lag.engine.NewLine(2)
	lag.engine.DrawTable(table)
}

func (lag *LifeAddendumGenerator) beneficiaries() {
	bDTO := lag.dto.Beneficiaries
	title := domain.TableCell{
		Text:      "Dati Beneficiari",
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

	const (
		firstColumnWidth  = 33
		secondColumnWidth = 70
		thirdColumnWidth  = 25
		fourthColumnWidth = 60
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

	lag.engine.WriteText(title)
	lag.engine.NewLine(2)
	lag.engine.DrawLine(10, lag.engine.GetY(), 200, lag.engine.GetY(), 0.25, constants.BlackColor)
	lag.engine.NewLine(2)
	for i := 0; i < 2; i++ {
		rows := [][]string{
			{"Cognome e Nome ", (*bDTO)[i].Surname + " " + (*bDTO)[i].Name, "Cod. Fisc: ", (*bDTO)[i].FiscalCode},
			{"Residente in ", (*bDTO)[i].StreetName + " " + (*bDTO)[i].StreetNumber + " " + (*bDTO)[i].City, "Data nascita: ", (*bDTO)[i].BirthDate},
			{"Mail ", (*bDTO)[i].Mail, "Telefono ", (*bDTO)[i].Phone},
		}
		table := parser(rows)
		lag.engine.DrawTable(table)
		lag.engine.NewLine(2)
		conf := "No"
		if (*bDTO)[i].Contactable {
			conf = "Sì"
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
	title := domain.TableCell{
		Text:      "Referente terzo",
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

	const (
		firstColumnWidth  = 33
		secondColumnWidth = 70
		thirdColumnWidth  = 25
		fourthColumnWidth = 60
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

	rows := [][]string{
		{"Cognome e Nome ", brDTO.Surname + " " + brDTO.Name, "Cod. Fisc: ", brDTO.FiscalCode},
		{"Residente in ", brDTO.StreetName + " " + brDTO.StreetNumber + " " + brDTO.City, "Data nascita: ", brDTO.BirthDate},
		{"Mail ", brDTO.Mail, "Telefono: ", brDTO.Phone},
	}
	table := parser(rows)

	lag.engine.WriteText(title)
	lag.engine.NewLine(2)
	lag.engine.DrawLine(10, lag.engine.GetY(), 200, lag.engine.GetY(), 0.25, constants.BlackColor)
	lag.engine.NewLine(2)
	lag.engine.DrawTable(table)
}
