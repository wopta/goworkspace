package proforma

import (
	"fmt"
	"strings"
	"time"

	"github.com/wopta/goworkspace/document/internal/constants"
	"github.com/wopta/goworkspace/document/internal/domain"
	"github.com/wopta/goworkspace/document/internal/dto"
	"github.com/wopta/goworkspace/document/internal/engine"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

type baseGenerator struct {
	engine      *engine.Fpdf
	now         time.Time
	networkNode *models.NetworkNode
	policy      *models.Policy
}

func (bg *baseGenerator) Save(rawDoc []byte) (string, error) {
	userPathFormat := "assets/users/%s/"
	path := fmt.Sprintf(userPathFormat, bg.policy.Contractor.Uid)
	filename := strings.ReplaceAll(fmt.Sprintf(path+models.ProformaDocumentFormat,
		bg.policy.NameDesc, bg.policy.CodeCompany, bg.policy.Annuity, time.Now().Format(constants.DayMonthYearFormat)), " ", "_")
	return bg.engine.Save(rawDoc, filename)
}

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

func (pg *ProformaGenerator) Generate() ([]byte, error) {
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
		"Contraente: " + contr.NameAndSurname,
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
		oneLineWidth      = constants.FullPageWidth - emptyColumnWidth
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

func (bg *baseGenerator) woptaFooter() {
	const (
		rowHeight   = 3
		columnWidth = 50
	)

	bg.engine.SetFooter(func() {
		bg.engine.SetY(-30)
		bg.engine.NewLine(3)

		entries := [][]string{
			{"Wopta Assicurazioni s.r.l", " ", " ", "www.wopta.it"},
			{"Galleria del Corso, 1", "Numero REA: MI 2638708", "CF | P.IVA | n. iscr. Registro Imprese:",
				"info@wopta.it"},
			{"20122 - Milano (MI)", "Capitale Sociale: € 204.839,26 i.v.", "12072020964", "(+39) 02 91240346"},
		}

		table := make([][]domain.TableCell, 0, 3)

		for index, entry := range entries {
			textColor := constants.BlackColor
			fontStyle := constants.RegularFontStyle
			if index == 0 {
				textColor = constants.PinkColor
				fontStyle = constants.BoldFontStyle
			}
			row := make([]domain.TableCell, 0, 4)

			for _, cell := range entry {
				row = append(row, domain.TableCell{
					Text:      cell,
					Height:    rowHeight,
					Width:     columnWidth,
					FontSize:  constants.SmallFontSize,
					FontStyle: fontStyle,
					FontColor: textColor,
					Fill:      false,
					FillColor: domain.Color{},
					Align:     constants.LeftAlign,
					Border:    "",
				})
			}
			table = append(table, row)
		}

		bg.engine.DrawTable(table)

		bg.engine.NewLine(3)

		bg.engine.WriteText(domain.TableCell{
			Text: "Wopta Assicurazioni s.r.l. è un intermediario assicurativo soggetto alla vigilanza dell’IVASS" +
				" ed iscritto alla Sezione A del Registro Unico degli Intermediari Assicurativi con numero" +
				" A000701923. Consulta gli estremi dell’iscrizione al sito https://servizi.ivass.it/RuirPubblica/",
			Height:    rowHeight,
			Width:     constants.FullPageWidth,
			FontSize:  constants.SmallFontSize,
			FontStyle: constants.RegularFontStyle,
			FontColor: constants.BlackColor,
			Fill:      false,
			FillColor: domain.Color{},
			Align:     constants.LeftAlign,
			Border:    "",
		})

		bg.engine.SetY(-7)

		bg.engine.WriteText(domain.TableCell{
			Text:      fmt.Sprintf("%d", bg.engine.PageNumber()),
			Height:    3,
			Width:     0,
			FontStyle: constants.RegularFontStyle,
			Fill:      false,
			FillColor: domain.Color{},
			Align:     constants.RightAlign,
			Border:    "",
		})
	})
}
