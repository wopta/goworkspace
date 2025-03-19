package addedndum

import (
	"fmt"
	"strings"
	"time"

	"github.com/wopta/goworkspace/document/internal/constants"
	"github.com/wopta/goworkspace/document/internal/domain"
	"github.com/wopta/goworkspace/document/internal/engine"
	"github.com/wopta/goworkspace/models"
)

const (
	tabDimension = 15
)

type baseGenerator struct {
	engine      *engine.Fpdf
	now         time.Time
	networkNode *models.NetworkNode
	policy      *models.Policy
}

func (bg *baseGenerator) Save(rawDoc []byte) (string, error) {
	//filename := strings.ReplaceAll(fmt.Sprintf("%s/%s/"+models.AddendumDocumentFormat, "temp", bg.policy.Uid,
	//	bg.policy.NameDesc, bg.policy.ProposalNumber), " ", "_")
	filename := strings.ReplaceAll(fmt.Sprintf("%s/%s/"+"%s_TEST_ADDENDUM_%d__%s.pdf", "temp", bg.policy.Uid,
		bg.policy.NameDesc, bg.policy.ProposalNumber, time.Now().Format("2006-1-2_15:04:05")), " ", "_")
	return bg.engine.Save(rawDoc, filename)
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
