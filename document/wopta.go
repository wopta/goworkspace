package document

import (
	"github.com/dustin/go-humanize"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func (skin Skin) PriceTable(m pdf.Maroto, data models.Policy) {
	m.Row(skin.RowHeight, func() {
		m.Col(12, func() {
			m.Text("Il premio per tutte le coperture assicurative attivate sulla polizza", skin.MagentaBoldtextLeft)

		})

	})

	skin.checkPage(m)
	h := []string{"Premio ", "Imponibile  ", "Imposte Assicurative ", "Totale"}
	var tablePremium [][]string

	if data.PaymentSplit == "monthly" {
		tablePremium = append(tablePremium, []string{"Rata Mensile", lib.HumanaizePriceEuro(data.PriceNett),
			lib.HumanaizePriceEuro(data.PriceGross - data.PriceNett), lib.HumanaizePriceEuro(data.PriceGross)})
	}
	if data.PaymentSplit == "year" {
		tablePremium = append(tablePremium, []string{"Annuale", lib.HumanaizePriceEuro(data.PriceNett), "€ " + humanize.FormatFloat("#.###,##", data.PriceGross-data.PriceNett), lib.HumanaizePriceEuro(data.PriceGross)})
	}
	tablePremium = append(tablePremium, []string{"Rata alla firma della polizza", lib.HumanaizePriceEuro(data.PriceNett), lib.HumanaizePriceEuro(data.PriceGross - data.PriceNett), lib.HumanaizePriceEuro(data.PriceGross)})

	skin.Space(m, 10.0)
	skin.TableLine(m, h, tablePremium)

}
