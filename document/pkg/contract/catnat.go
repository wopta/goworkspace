package contract

import (
	"strings"
	"time"

	"gitlab.dev.wopta.it/goworkspace/document/internal/constants"
	"gitlab.dev.wopta.it/goworkspace/document/internal/domain"
	"gitlab.dev.wopta.it/goworkspace/document/internal/dto"
	"gitlab.dev.wopta.it/goworkspace/document/internal/engine"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/network"
)

type CatnatGenerator struct {
	*baseGenerator
	dto dto.CatnatDTO
}

func NewCatnatGenerator(engine *engine.Fpdf, policy *models.Policy, node *models.NetworkNode, product models.Product, isProposal bool) *CatnatGenerator {
	dto := dto.NewCatnatDto()
	dto.FromPolicy(policy, node)

	var worksForNode *models.NetworkNode
	if node != nil && node.WorksForUid != "" {
		worksForNode = network.GetNetworkNodeByUid(node.WorksForUid)
	}

	return &CatnatGenerator{
		baseGenerator: &baseGenerator{
			engine:       engine,
			isProposal:   isProposal,
			now:          time.Now(),
			signatureID:  0,
			networkNode:  node,
			policy:       policy,
			worksForNode: worksForNode,
		},
		dto: dto,
	}
}

func (c *CatnatGenerator) Generate() {
	c.addCatnatHeader()
	c.engine.NewPage()
	c.engine.NewLine(5)
	c.addCatnatHeading()
	c.engine.NewLine(10)
	c.addTableGuarantee()
	c.engine.NewLine(3)
	c.addFrazionamento()
	c.engine.NewLine(4)
	c.addSetInformativoInfo()
	c.engine.NewLine(4)
	c.addWhoAreWeCatnat()

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
		c.engine.WriteText(c.engine.GetTableCell(""))
		c.engine.DrawTable(parseLogos([]string{"				Wopta per te", "Catastrofali Azienda"}))
		c.engine.InsertImage(lib.GetAssetPathByEnvV2()+"logo_vita.png", 10, 15, 13, 13)
		c.engine.NewLine(4)
		c.engine.WriteText(c.engine.GetTableCell("PREVENTIVO\nANONIMO", domain.FontSize(10), constants.BoldFontStyle))
		c.engine.InsertImage(lib.GetAssetPathByEnvV2()+"logo_wopta.png", 165, 15, 35, 10)
	})
}

func (c *CatnatGenerator) addCatnatHeading() {
	c.engine.WriteText(c.engine.GetTableCell("DATI DEL PREVENIVO", constants.BoldFontStyle, constants.PinkColor, domain.FontSize(10)))
	c.engine.NewLine(5)
	c.engine.WriteText(c.engine.GetTableCell("DATI SEDE DA ASSICURARE", constants.BoldFontStyle, domain.FontSize(12)))
	c.engine.DrawLine(10, c.engine.GetY(), 200, c.engine.GetY(), 0.25, constants.PinkColor)
	c.engine.NewLine(2)
	const (
		firstColumnWidth  = 35
		secondColumnWidth = 95
		thirdColumnWidth  = 25
		fourthColumnWidth = 35
	)
	dataSede := [][]string{
		{"Indirizzo", "bla bla indirizzo", "Tipo Utilizzo Sede", "conduttore"},
		{"Anno di costruzione", "bla bla costruzione", "Materiale di costruzione", "materiale"},
		{"Numero di piani edificio oltre il piano terra", "bla bla numero", "", ""},
	}
	spaces := []int{
		100, 72, 0,
	}
	for i, line := range dataSede {
		space := strings.Repeat(" ", int(spaces[i]))
		c.baseGenerator.engine.WriteTexts(c.engine.GetTableCell(line[0]+"  ", constants.BoldFontStyle), c.engine.GetTableCell(line[1]), c.engine.GetTableCell(space+line[2]+"  ", constants.BoldFontStyle), c.engine.GetTableCell(line[3]))
		c.engine.DrawLine(10, c.engine.GetY(), 200, c.engine.GetY(), 0.25, constants.PinkColor)
	}
}

func (c *CatnatGenerator) addTableGuarantee() {
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
		{"Terremoto", "100.000 €", "50.000 €", "====", "150,94 €"},
		{"Alluvione", "100.000 €", "50.000 €", "====", "150,94 €"},
		{"Frane", "100.000 €", "50.000 €", "====", "150,94 €"},
	}
	c.engine.DrawTable(guranteeTable(guaranteeData))
}

func (c *CatnatGenerator) addFrazionamento() {
	c.engine.WriteText(c.engine.GetTableCell("Frazionamento: "+"pippo", constants.BoldFontStyle))
	c.engine.NewLine(8)
	c.engine.WriteTexts(c.engine.GetTableCell("Premio di polizza ", constants.BoldFontStyle), c.engine.GetTableCell("wwee"))
	c.engine.WriteTexts(c.engine.GetTableCell("Contributo servizi di intermediazione annuale ", constants.BoldFontStyle), c.engine.GetTableCell("fjdsklfds"))
	c.engine.WriteTexts(c.engine.GetTableCell("Totale da pagare ", constants.BoldFontStyle), c.engine.GetTableCell("fjdsklfds"))
	c.engine.NewLine(4)
	c.engine.WriteText(c.engine.GetTableCell("Milano, il pippo/pippo/pippo", constants.BoldFontStyle))
}

func (c *CatnatGenerator) addSetInformativoInfo() {
	c.engine.WriteText(c.engine.GetTableCell("Il presente preventivo non ha validità di proposta assicurativa. Ha valore esclusivamente nel giorno di emissione e non impegna la compagnia alla sottoscrizione del rischio.", constants.BoldFontStyle))
	c.engine.NewLine(2)
	c.engine.WriteText(c.engine.GetTableCell("Prima della sottoscrizione leggere il set informativo.", constants.BoldFontStyle))
}
func (c *CatnatGenerator) addWhoAreWeCatnat() {
	c.whoWeAre()
	c.engine.NewLine(5)
	c.engine.WriteTexts(c.engine.GetTableCell("Net Insurance S.p.a ", constants.BoldFontStyle), c.engine.GetTableCell("compagnia assicurativa, Sede Legale e Direzione Generale via Giuseppe Antonio Guattani, 4 00161 Roma"))
}
