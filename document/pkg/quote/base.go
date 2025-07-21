package quote

import (
	"time"

	"gitlab.dev.wopta.it/goworkspace/document/internal/constants"
	"gitlab.dev.wopta.it/goworkspace/document/internal/domain"
	"gitlab.dev.wopta.it/goworkspace/document/internal/dto"
	"gitlab.dev.wopta.it/goworkspace/document/internal/engine"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
)

type baseGenerator struct {
	engine  *engine.Fpdf
	now     time.Time
	policy  *models.Policy
	product *models.Product
	dto     *dto.QuoteBaseDTO
}

func (bg *baseGenerator) mainHeader() {
	var (
		productLogo      string
		productName      string
		woptaLogo        string  = lib.GetAssetPathByEnvV2() + constants.WoptaLogo
		productLogoWidth float64 = 15
		cellHeight       float64 = productLogoWidth / 2
		woptaLogoWidth   float64 = 41
		logoExtraSpacing float64 = 1
	)

	if path, ok := constants.ProductLogoMap[bg.policy.Name]; ok {
		productLogo = lib.GetAssetPathByEnvV2() + path
	}

	if name, ok := constants.ProductNameMap[bg.policy.Name]; ok {
		productName = name
	}

	bg.engine.SetHeader(func() {
		if productLogo != "" {
			bg.engine.InsertImage(productLogo, 10, 10, productLogoWidth, productLogoWidth)
		}
		if productName != "" {
			var (
				firstColWidth  = productLogoWidth + logoExtraSpacing
				secondColWidth = constants.FullPageWidth - productLogoWidth - logoExtraSpacing - woptaLogoWidth
			)

			bg.engine.DrawTable([][]domain.TableCell{
				{
					{
						Text:      "",
						Height:    cellHeight,
						Width:     firstColWidth,
						FontSize:  constants.ExtraLargeFontSize,
						FontStyle: constants.BoldFontStyle,
						FontColor: constants.PinkColor,
					},
					{
						Text:      "Wopta per te",
						Height:    cellHeight,
						Width:     secondColWidth,
						FontSize:  constants.ExtraLargeFontSize,
						FontStyle: constants.BoldFontStyle,
						FontColor: constants.PinkColor,
					},
				},
				{
					{
						Text:      "",
						Height:    cellHeight,
						Width:     firstColWidth,
						FontSize:  constants.ExtraLargeFontSize,
						FontStyle: constants.BoldFontStyle,
						FontColor: constants.PinkColor,
					},
					{
						Text:      productName,
						Height:    cellHeight,
						Width:     secondColWidth,
						FontSize:  constants.ExtraLargeFontSize,
						FontStyle: constants.ItalicFontStyle,
						FontColor: constants.BlackColor,
					},
				},
			})
		}
		bg.engine.InsertImage(woptaLogo, 159, 10, woptaLogoWidth, 12)
		bg.engine.NewLine(5)
	})
}

func (bg *baseGenerator) heading() {
	bg.engine.WriteText(domain.TableCell{
		Text:      "PREVENTIVO\nANONIMO",
		Height:    constants.CellHeight,
		Width:     constants.FullPageWidth,
		FontSize:  constants.LargeFontSize,
		FontStyle: constants.BoldFontStyle,
		FontColor: constants.BlackColor,
	})
}

func (bg *baseGenerator) priceSummary() {
	if bg.dto.Price.Consultancy.ValueFloat == 0 {
		return
	}

	bg.engine.NewLine(5)
	bg.engine.RawWriteText(domain.TableCell{
		Text:      "Premio di Polizza ",
		Height:    constants.CellHeight,
		FontStyle: constants.BoldFontStyle,
		FontColor: constants.BlackColor,
		FontSize:  constants.LargeFontSize,
	})
	bg.engine.RawWriteText(domain.TableCell{
		Text:      bg.dto.Price.Gross.Text,
		Height:    constants.CellHeight,
		FontColor: constants.BlackColor,
		FontSize:  constants.LargeFontSize,
	})
	bg.engine.NewLine(5)
	bg.engine.RawWriteText(domain.TableCell{
		Text:      "Contributo servizi di intermediazione annuale ",
		Height:    constants.CellHeight,
		FontStyle: constants.BoldFontStyle,
		FontColor: constants.BlackColor,
		FontSize:  constants.LargeFontSize,
	})
	bg.engine.RawWriteText(domain.TableCell{
		Text:      bg.dto.Price.Consultancy.Text,
		Height:    constants.CellHeight,
		FontColor: constants.BlackColor,
		FontSize:  constants.LargeFontSize,
	})
	bg.engine.NewLine(5)
	bg.engine.RawWriteText(domain.TableCell{
		Text:      "Totale da pagare ",
		Height:    constants.CellHeight,
		FontStyle: constants.BoldFontStyle,
		FontColor: constants.BlackColor,
		FontSize:  constants.LargeFontSize,
	})
	bg.engine.RawWriteText(domain.TableCell{
		Text:      bg.dto.Price.Total.Text,
		Height:    constants.CellHeight,
		FontColor: constants.BlackColor,
		FontSize:  constants.LargeFontSize,
	})
}

func (bg *baseGenerator) WhoWeAre() {
	bg.engine.WriteText(domain.TableCell{
		Text:      "Chi siamo",
		Height:    constants.CellHeight,
		FontStyle: constants.BoldFontStyle,
		FontColor: constants.PinkColor,
		FontSize:  constants.LargeFontSize,
	})

	bg.engine.NewLine(2)

	bg.engine.RawWriteText(domain.TableCell{
		Text:      "Wopta Assicurazioni S.r.l.",
		Height:    constants.CellHeight,
		FontStyle: constants.BoldFontStyle,
		FontColor: constants.BlackColor,
		FontSize:  constants.RegularFontSize,
	})
	bg.engine.RawWriteText(domain.TableCell{
		Text:      " - intermediario assicurativo, soggetto al controllo dell’IVASS ed iscritto dal 14.02.2022 al Registro Unico degli Intermediari, in Sezione A nr. A000701923, avente sede legale in Galleria del Corso, 1 – 20122 Milano (MI). Capitale sociale Euro 120.000 - Codice Fiscale, Reg. Imprese e Partita IVA: 12072020964 - Iscritta al Registro delle imprese di Milano – REA MI 2638708",
		Height:    constants.CellHeight,
		FontColor: constants.BlackColor,
		FontSize:  constants.RegularFontSize,
	})
}
