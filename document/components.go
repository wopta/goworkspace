package document

import (
	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
	//model "github.com/wopta/goworkspace/models"
)

func Customer(m pdf.Maroto) {
	textBoldRight := props.Text{
		Top:   1.5,
		Size:  9,
		Style: consts.Bold,
		Align: consts.Center,
	}

	magenta := color.Color{
		Red:   229,
		Green: 0,
		Blue:  117,
	}
	textMagenta := props.Text{
		Color: magenta,
		Top:   0,
		Style: consts.Normal,
		Align: consts.Left,
		Size:  1,
	}
	m.Row(5, func() {
		m.Col(2, func() {
			m.Text("Assicurato:", textBoldRight)

		})

		m.Col(2, func() {
			m.Text("xxxx", props.Text{
				Top:   1.5,
				Size:  9,
				Style: consts.Normal,
				Align: consts.Center,
			})
		})
	})

	m.Row(1, func() {
		m.Col(2, func() {
			m.Text("_________________________________________________", textMagenta)

		})
	})
}
