package quote

import (
	"time"

	"gitlab.dev.wopta.it/goworkspace/document/internal/constants"
	"gitlab.dev.wopta.it/goworkspace/document/internal/domain"
	"gitlab.dev.wopta.it/goworkspace/document/internal/dto"
	"gitlab.dev.wopta.it/goworkspace/document/internal/engine"
	"gitlab.dev.wopta.it/goworkspace/document/pkg/internal/catnat"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
)

type CatnatGenerator struct {
	*baseGenerator
	dto dto.QuoteCatnatDTO
}

func NewCatnatGenerator(engine *engine.Fpdf, policy *models.Policy, product models.Product) *CatnatGenerator {
	dto := dto.NewCatnatDto()
	dto.FromPolicy(policy)

	return &CatnatGenerator{
		baseGenerator: &baseGenerator{
			engine: engine,
			now:    time.Now(),
			policy: policy,
		},
		dto: dto,
	}
}

func (c *CatnatGenerator) Exec() ([]byte, error) {
	c.addCatnatHeader()
	c.engine.NewPage()
	c.engine.NewLine(5)
	catnat.AddBuildingInformation(c.engine, c.dto.Sede, c.dto.Questions)
	c.engine.NewLine(10)
	catnat.AddTableGuarantee(c.engine, c.dto.Guarantees)
	c.engine.NewLine(3)
	c.addFrazionamento()
	c.engine.NewLine(4)
	c.addSetInformativoInfo()
	c.engine.NewLine(4)
	c.addWhoAreWeCatnat()

	return c.engine.RawDoc()
}

func (c *CatnatGenerator) addCatnatHeader() {
	const (
		firstColumnWidth  = 15
		secondColumnWidth = constants.FullPageWidth - firstColumnWidth
	)
	parseLogos := func(texts []string) [][]domain.TableCell {
		result := make([][]domain.TableCell, 0, len(texts))

		result = append(result, []domain.TableCell{
			{
				Text:      "",
				Height:    constants.CellHeight,
				Width:     firstColumnWidth,
				FontStyle: constants.RegularFontStyle,
				FontColor: constants.BlackColor,
				Fill:      false,
				FontSize:  constants.RegularFontSize,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "",
			},
			{
				Text:      texts[0],
				Height:    constants.CellHeight,
				FontStyle: constants.BoldFontStyle,
				FontColor: constants.PinkColor,
				FontSize:  18,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "",
			},
		})

		result = append(result, []domain.TableCell{
			{
				Text:      "",
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
				Text:      texts[1],
				Height:    constants.CellHeight + 5,
				FontStyle: constants.BoldFontStyle,
				FontColor: constants.BlackColor,
				FontSize:  16,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "",
			},
		})
		return result
	}
	c.engine.SetHeader(func() {
		c.engine.DrawTable(parseLogos([]string{"				Wopta per te", "Catastrofali Azienda"}))
		c.engine.InsertImage(lib.GetAssetPathByEnvV2()+"logo_vita.png", 10, 15, 13, 13)
		c.engine.NewLine(4)
		c.engine.WriteText(c.engine.GetTableCell("PREVENTIVO\nANONIMO", domain.FontSize(10), constants.BoldFontStyle))
		c.engine.InsertImage(lib.GetAssetPathByEnvV2()+"logo_wopta.png", 165, 15, 35, 10)
	})
}

func (c *CatnatGenerator) addFrazionamento() {
	c.engine.WriteText(c.engine.GetTableCell("Frazionamento: "+c.dto.PaymentSplit, constants.BoldFontStyle))
	c.engine.NewLine(8)
	c.engine.WriteTexts(c.engine.GetTableCell("Premio di polizza ", constants.BoldFontStyle), c.engine.GetTableCell(c.dto.Prize.Gross.Text))
	c.engine.WriteTexts(c.engine.GetTableCell("Contributo servizi di intermediazione annuale ", constants.BoldFontStyle), c.engine.GetTableCell(c.dto.Prize.Consultancy.Text))
	c.engine.WriteTexts(c.engine.GetTableCell("Totale da pagare ", constants.BoldFontStyle), c.engine.GetTableCell(c.dto.Prize.Total.Text))
	c.engine.NewLine(4)
	c.engine.WriteText(c.engine.GetTableCell("Milano, "+c.now.Format(constants.DayMonthYearFormat), constants.BoldFontStyle))
}

func (c *CatnatGenerator) addSetInformativoInfo() {
	c.engine.WriteText(c.engine.GetTableCell("Il presente preventivo non ha validità di proposta assicurativa.Ha valore esclusivamente nel giorno di emissione e non impegna la compagnia alla sottoscrizione del rischio.", constants.BoldFontStyle))
	c.engine.NewLine(2)
	c.engine.WriteText(c.engine.GetTableCell("Prima della sottoscrizione leggere il set informativo.", constants.BoldFontStyle))
}
func (c *CatnatGenerator) addWhoAreWeCatnat() {
	c.WhoWeAre()
	c.engine.NewLine(5)
	c.engine.WriteTexts(c.engine.GetTableCell("Net Insurance S.p.a ", constants.BoldFontStyle), c.engine.GetTableCell("compagnia assicurativa, Sede Legale e Direzione Generale via Giuseppe Antonio Guattani, 4 00161 Roma"))
}
