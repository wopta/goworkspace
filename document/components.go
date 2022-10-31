package document

import (
	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
	//model "github.com/wopta/goworkspace/models"
)

type Skin struct {
	LineColor         color.Color
	TextColor         color.Color
	TitleColor        color.Color
	rowHeight         float64
	rowtableHeight    float64
	LineHeight        float64
	Size              float64
	SizeTitle         float64
	TableHeight       float64
	rowtableHeightMin float64
}

func (s Skin) CustomerRow(m pdf.Maroto, k string, v string) pdf.Maroto {
	textBoldRight := props.Text{
		Top:   1.5,
		Size:  s.Size,
		Style: consts.Bold,
		Align: consts.Right,
		Color: s.TextColor,
	}
	textRight := props.Text{
		Top:   1.5,
		Size:  s.Size,
		Style: consts.Bold,
		Align: consts.Left,
		Color: s.TextColor,
	}
	textMagenta := props.Text{
		Color: s.LineColor,
		Top:   0,
		Style: consts.Normal,
		Align: consts.Left,
		Size:  5,
	}
	m.Row(3, func() {
		m.Col(4, func() {
			m.Text(k, textBoldRight)

		})

		m.Col(4, func() {
			m.Text(v, textRight)
		})
	})

	m.Row(1, func() {
		m.Col(2, func() {
			m.Text("_________________________________________________", textMagenta)

		})
	})
	return m
}
func (s Skin) TableRow(m pdf.Maroto, colText []string, isLine bool, colspace uint, h float64) pdf.Maroto {

	textBold := props.Text{
		Top:   1.5,
		Size:  s.Size,
		Style: consts.Bold,
		Align: consts.Left,
		Color: s.TextColor,
	}
	text := props.Text{
		Top:   1.5,
		Size:  s.Size,
		Style: consts.Normal,
		Align: consts.Left,
		Color: s.TextColor,
	}
	textLast := props.Text{
		Top:   1.5,
		Size:  s.Size,
		Style: consts.Normal,
		Align: consts.Right,
		Color: s.TextColor,
	}
	linePropMagenta := props.Line{
		Color: s.LineColor,
		Style: consts.Solid,
		Width: 0.2,
	}
	m.Row(h, func() {

		for i, k := range colText {
			var prop props.Text
			prop = text
			//log.Println(len(colText))
			if len(colText) > 2 {
				if i == 0 {
					prop = textBold
				} else if i == len(colText)-1 {
					prop = textLast
				} else {
					prop = text
				}
			}
			m.Col(colspace, func() {

				m.Text(k, prop)

			})
		}

	})
	if isLine {
		m.Line(s.LineHeight, linePropMagenta)
	}

	return m
}
func (s Skin) TableHeader(m pdf.Maroto, colText []string, isLine bool, colspace uint, h float64, c consts.Align) pdf.Maroto {

	textBold := props.Text{
		Top:   1.5,
		Size:  s.SizeTitle,
		Style: consts.Bold,
		Align: c,
		Color: s.TextColor,
	}

	linePropMagenta := props.Line{
		Color: s.LineColor,
		Style: consts.Solid,
		Width: 0.2,
	}
	m.Row(h, func() {

		for _, k := range colText {

			m.Col(colspace, func() {

				m.Text(k, textBold)

			})
		}

	})
	if isLine {
		m.Line(s.LineHeight, linePropMagenta)
	}

	return m
}
func (s Skin) RowBullet(m pdf.Maroto, k string, v string, style consts.Style) pdf.Maroto {
	rowh := lenToHeight(v)
	prop := props.Text{
		Top:   1.5,
		Size:  s.Size,
		Style: style,
		Align: consts.Left,
		Color: s.TextColor,
	}
	propR := props.Text{
		Top:   1.5,
		Size:  s.Size,
		Style: style,
		Align: consts.Right,
		Color: s.TextColor,
	}

	m.Row(rowh, func() {
		m.Col(1, func() {
			m.Text(k, propR)

		})

		m.Col(11, func() {
			m.Text(v, prop)
		})
	})

	return m
}
func (s Skin) RowCol1(m pdf.Maroto, v string, style consts.Style) pdf.Maroto {
	rowh := lenToHeight(v)
	prop := props.Text{
		Top:             0.5,
		Size:            s.Size,
		Style:           style,
		Align:           consts.Left,
		Color:           s.TextColor,
		VerticalPadding: 0.0,
	}

	m.Row(rowh, func() {

		m.Col(12, func() {
			m.Text(v, prop)
		})
	})

	return m
}
func (s Skin) Space(m pdf.Maroto, h float64) pdf.Maroto {

	m.Row(h, func() {
		m.ColSpace(12)
	})

	return m
}
