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

type LifeAddendumGenerator struct {
	*baseGenerator
	dto *dto.AddendumBeneficiariesDTO
}

func NewLifeAddendumGenerator(engine *engine.Fpdf, policy *models.Policy) *LifeAddendumGenerator {
	now := time.Now()
	LifeAddendumDTO := dto.NewBeneficiariesDto()
	LifeAddendumDTO.FromPolicy(policy, now)
	return &LifeAddendumGenerator{
		baseGenerator: &baseGenerator{
			engine: engine,
			now:    now,
			policy: policy,
		},
		dto: LifeAddendumDTO,
	}
}

func (lag *LifeAddendumGenerator) Generate() {
	lag.mainHeader()

	lag.engine.NewPage()

	lag.contract()

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

	lag.contractorSignature()

	lag.woptaFooter()

}

func (lag *LifeAddendumGenerator) mainHeader() {
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

		lag.engine.InsertImage(lib.GetAssetPathByEnvV2()+"logo_vita.png", 12, 6.5, 12, 12)
		lag.engine.InsertImage(lib.GetAssetPathByEnvV2()+"logo_wopta.png", 160, 6.5, 35, 10)
		lag.engine.NewLine(7)
		origY := lag.engine.GetY()
		lag.engine.SetY(origY - 15)
		lag.engine.WriteText(first)
		lag.engine.SetY(lag.engine.GetY() - 1)
		lag.engine.WriteText(second)
		lag.engine.SetY(origY)
	})
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
	lag.engine.WriteText(domain.TableCell{
		Text: "Dichiarazione di Variazione dati Anagrafici Contraente-" +
			"Assicurato-Beneficiario-Referente Terzo",
		Height:    constants.CellHeight,
		Width:     constants.FullPageWidth,
		FontStyle: constants.BoldFontStyle,
		FontColor: constants.PinkColor,
	})
	lag.engine.NewLine(3)
	lag.engine.WriteText(domain.TableCell{
		Text: "Come da richiesta sono state trasmesse all’assicuratore " +
			"AXA France Vie S.A. - Rappresentanza Generale per l’Italia le " +
			"seguenti variazioni alla polizza sopra meglio evidenziata",
		Height: constants.CellHeight,
		Width:  constants.FullPageWidth,
	})
	lag.engine.NewLine(1)
	lag.engine.WriteText(domain.TableCell{
		Text: "Le modifiche non sono attive in assenza di firma da parte del " +
			"Contraente, con allegata copia documento di riconoscimento " +
			"(carta identità, passaporto, patente, in corso di validità alla " +
			"data della firma).",
		Height:    constants.CellHeight,
		Width:     constants.FullPageWidth,
		FontStyle: constants.BoldFontStyle,
	})
}

func (lag *LifeAddendumGenerator) contractor() {
	cDTO := lag.dto.Contractor
	checked := " "
	var rows1 [][]string
	var rows2 [][]string
	var domTxt [][]string
	if cDTO.FiscalCode != constants.EmptyField {
		checked = "X"
		rows1 = [][]string{
			{"Cognome e Nome ", cDTO.Surname + " " + cDTO.Name, "Cod. Fisc: ", cDTO.FiscalCode},
			{"Residente in ", cDTO.StreetName + " " + cDTO.StreetNumber + ", " + cDTO.PostalCode + " " + cDTO.City + " (" + cDTO.Province + ")", "Data nascita: ", cDTO.BirthDate},
		}
		rows2 = [][]string{
			{"Mail ", cDTO.Mail, "Telefono: ", cDTO.Phone},
		}
		domTxt = [][]string{
			{"Domicilio ", cDTO.DomStreetName + " " + cDTO.DomStreetNumber + ", " + cDTO.DomPostalCode + " " + cDTO.DomCity + " (" + cDTO.DomProvince + ")"},
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
	if iDTO.FiscalCode != constants.EmptyField {
		checked = "X"
		rows1 = [][]string{
			{"Cognome e Nome ", iDTO.Surname + " " + iDTO.Name, "Cod. Fisc: ", iDTO.FiscalCode},
			{"Residente in ", iDTO.StreetName + " " + iDTO.StreetNumber + ", " + iDTO.PostalCode + " " + iDTO.City + " (" + iDTO.Province + ")", "Data nascita: ", iDTO.BirthDate},
		}
		rows2 = [][]string{
			{"Mail ", iDTO.Mail, "Telefono: ", iDTO.Phone},
		}
		domTxt = [][]string{
			{"Domicilio ", iDTO.DomStreetName + " " + iDTO.DomStreetNumber + ", " + iDTO.DomPostalCode + " " + iDTO.DomCity + " (" + iDTO.DomProvince + ")"},
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
			if v.FiscalCode != constants.EmptyField {
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
				{"Cognome e Nome ", (*bDTO)[i].Surname + " " + (*bDTO)[i].Name, "Cod. Fisc: ", (*bDTO)[i].FiscalCode},
				{"Residente in ", (*bDTO)[i].StreetName + " " + (*bDTO)[i].StreetNumber + ", " + (*bDTO)[i].PostalCode + " " + (*bDTO)[i].City + " (" + (*bDTO)[i].Province + ")", "Data nascita: ", (*bDTO)[i].BirthDate},
				{"Mail ", (*bDTO)[i].Mail, "Telefono ", (*bDTO)[i].Phone},
				{"Relazione con Assicurato ", (*bDTO)[i].Relation, "Quota indennizzo", " "},
			}
			relTxt = [][]string{}
		} else {
			rows = [][]string{
				{"Cognome e Nome ", " ", "Cod. Fisc: ", " "},
				{"Residente in ", " ", "Data nascita: ", " "},
				{"Mail ", " ", "Telefono ", " "},
				{"Relazione con Assicurato ", " ", "Quota indennizzo", " "},
			}
			relTxt = [][]string{}
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
	if brDTO.FiscalCode != constants.EmptyField {
		checked = "X"
		rows = [][]string{
			{"Cognome e Nome ", brDTO.Surname + " " + brDTO.Name, "Cod. Fisc: ", brDTO.FiscalCode},
			{"Residente in ", brDTO.StreetName + " " + brDTO.StreetNumber + ", " + brDTO.PostalCode + " " + brDTO.City + " (" + brDTO.Province + ")", "Data nascita: ", brDTO.BirthDate},
			{"Mail ", brDTO.Mail, "Telefono: ", brDTO.Phone},
		}
	} else {
		checked = "X"
		rows = [][]string{
			{"Cognome e Nome ", constants.EmptyField, "Cod. Fisc: ", constants.EmptyField},
			{"Residente in ", constants.EmptyField, "Data nascita: ", constants.EmptyField},
			{"Mail ", constants.EmptyField, "Telefono: ", constants.EmptyField},
		}
	}
	if brDTO.Name == "" {
		checked = " "
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

func (lag *LifeAddendumGenerator) contractorSignature() {
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
