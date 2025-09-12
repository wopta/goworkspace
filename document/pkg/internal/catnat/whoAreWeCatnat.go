package catnat

import (
	"gitlab.dev.wopta.it/goworkspace/document/internal/constants"
	"gitlab.dev.wopta.it/goworkspace/document/internal/engine"
)

func AddWhoAreWeCatnat(c *engine.Fpdf) {
	c.NewLine(5)
<<<<<<< HEAD
	c.WriteTexts(c.GetTableCell("Net Insurance S.p.a ", constants.BoldFontStyle), c.GetTableCell("compagnia assicurativa, Sede Legale e Direzione Generale via Giuseppe Antonio Guattani, 4 00161 Roma"))
=======
	c.WriteTexts(c.GetTableCell("Net Insurance S.p.a ", constants.BoldFontStyle), c.GetTableCell("impresa di assicurazione, Società per Azioni facente parte del Gruppo Assicurativo Poste Vita – Albo Gruppi Assicurativi IVASS n. 43 – Via Giuseppe Antonio Guattani n. 4, 00161 Roma, Tel. 06 89326.1 – Fax 06 89326.800; Sito internet: www.netinsurance.it; e-mail: info@netinsurance.it; PEC: netinsurance@pec.netinsurance.it. "))
>>>>>>> master
}
