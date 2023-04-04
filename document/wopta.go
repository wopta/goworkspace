package document

import (
	"github.com/dustin/go-humanize"
	"github.com/johnfercher/maroto/pkg/pdf"
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
		tablePremium = append(tablePremium, []string{"Rata Mensile", "€ " + humanize.FormatFloat("#.###,##", (data.PriceNett/12)), "€ " + humanize.FormatFloat("#.###,##", ((data.PriceGross-data.PriceNett)/12)), "€ " + humanize.FormatFloat("#.###,##", (data.PriceGross/12))})
		tablePremium = append(tablePremium, []string{"Rata alla firma della polizza", "€ " + humanize.FormatFloat("#.###,##", (data.PriceNett/12)), "€ " + humanize.FormatFloat("#.###,##", ((data.PriceGross-data.PriceNett)/12)), "€ " + humanize.FormatFloat("#.###,##", (data.PriceGross/12))})

	}
	if data.PaymentSplit == "year" {
		tablePremium = append(tablePremium, []string{"Annuale", "€ " + humanize.FormatFloat("#.###,##", data.PriceNett), "€ " + humanize.FormatFloat("#.###,##", data.PriceGross-data.PriceNett), "€ " + humanize.FormatFloat("#.###,##", data.PriceGross)})
		tablePremium = append(tablePremium, []string{"Rata alla firma della polizza", "€ " + humanize.FormatFloat("#.###,##", data.PriceNett), "€ " + humanize.FormatFloat("#.###,##", data.PriceGross-data.PriceNett), "€ " + humanize.FormatFloat("#.###,##", data.PriceGross)})

	}
	skin.Space(m, 10.0)
	skin.TableLine(m, h, tablePremium)

}
