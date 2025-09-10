package utils

import (
	"gitlab.dev.wopta.it/goworkspace/document/internal/constants"
	"gitlab.dev.wopta.it/goworkspace/document/internal/domain"
	"gitlab.dev.wopta.it/goworkspace/document/internal/engine"
)

func AddWhoWeAre(e *engine.Fpdf) {
	e.WriteText(domain.TableCell{
		Text:      "Chi siamo",
		Height:    constants.CellHeight,
		FontStyle: constants.BoldFontStyle,
		FontColor: constants.PinkColor,
		FontSize:  constants.LargeFontSize,
	})

	e.NewLine(2)

	e.RawWriteText(domain.TableCell{
		Text:      "Wopta Assicurazioni S.r.l.",
		Height:    constants.CellHeight,
		FontColor: constants.BlackColor,
		FontStyle: constants.BoldFontStyle,
		FontSize:  constants.RegularFontSize,
	})
	e.RawWriteText(domain.TableCell{
		Text:      " - intermediario assicurativo, soggetto al controllo dell’IVASS ed iscritto dal 14.02.2022 al Registro Unico degli Intermediari, in Sezione A nr. A000701923, avente sede legale in Galleria del Corso, 1 – 20122 Milano (MI). Capitale sociale Euro 120.000 - Codice Fiscale, Reg. Imprese e Partita IVA: 12072020964 - Iscritta al Registro delle imprese di Milano – REA MI 2638708",
		Height:    constants.CellHeight,
		FontColor: constants.BlackColor,
		FontSize:  constants.RegularFontSize,
	})
}
