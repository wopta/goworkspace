package catnat

import (
	"gitlab.dev.wopta.it/goworkspace/document/internal/constants"
	"gitlab.dev.wopta.it/goworkspace/document/internal/engine"
)

func WhoAreWeCatnat(c *engine.Fpdf) {
	c.NewLine(5)
	c.WriteTexts(c.GetTableCell("Net Insurance S.p.a ", constants.BoldFontStyle), c.GetTableCell("compagnia assicurativa, Sede Legale e Direzione Generale via Giuseppe Antonio Guattani, 4 00161 Roma"))
}
