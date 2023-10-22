package document

import (
	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
	lib "github.com/wopta/goworkspace/lib"
	//model "github.com/wopta/goworkspace/models"
)

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
func (s Skin) TableRow(m pdf.Maroto, colText []string, isLine bool, colspace uint, h float64, minus float64, c consts.Align) pdf.Maroto {

	textBold := props.Text{
		Top:   1.5,
		Size:  s.Size - minus,
		Style: consts.Bold,
		Align: c,
		Color: s.TextColor,
	}
	text := props.Text{
		Top:   1.5,
		Size:  s.Size - minus,
		Style: consts.Normal,
		Align: c,
		Color: s.TextColor,
	}
	textLast := props.Text{
		Top:   1.5,
		Size:  s.Size - minus,
		Style: consts.Normal,
		Align: c,
		Color: s.TextColor,
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
		m.Line(s.LineHeight, s.Line)
	}

	return m
}

func (s Skin) MultiRow(m pdf.Maroto, colTexts [][]string, isLine bool, colspace []uint, h float64) pdf.Maroto {

	textS := s.NormaltextLeftExt
	textS.Size = textS.Size - 2
	var maxRow int
	for _, k := range colTexts {
		if maxRow < len(k) {
			maxRow = len(k)
		}
	}
	m.Row(float64(maxRow*4), func() {
		for i, k := range colTexts {
			m.Col(colspace[i], func() {
				for a, t := range k {
					top := a * 4
					textS.Top = float64(top)
					if i == 0 {
						textS.Style = consts.Bold
					} else {
						textS.Style = consts.Normal
					}
					m.Text(t, textS)
				}

			})
		}

	})
	if isLine {
		m.Line(s.LineHeight, s.Line)
	}

	return m
}
func (s Skin) TableHeader(m pdf.Maroto, colText []string, isLine bool, colspace uint, h float64, c consts.Align, minus float64) pdf.Maroto {

	textBold := props.Text{
		Top:   1.5,
		Size:  s.SizeTitle - minus,
		Style: consts.Bold,
		Align: c,
		Color: s.TextColor,
	}

	m.Row(h, func() {

		for _, k := range colText {

			m.Col(colspace, func() {

				m.Text(k, textBold)

			})
		}

	})
	if isLine {
		m.Line(s.LineHeight, s.Line)
	}

	return m
}
func (s Skin) RowBullet(m pdf.Maroto, k string, v string, style consts.Style) pdf.Maroto {
	rowh := s.lenToHeight(v)
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
	rowh := s.lenToHeight(v)
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
func (s Skin) BackgroundColorRow(m pdf.Maroto, text string, colorBg color.Color, p props.Text, height float64) pdf.Maroto {

	m.SetBackgroundColor(colorBg)
	m.Row(height, func() {
		m.Col(12, func() {
			m.Text(text, p)

		})

	})
	m.SetBackgroundColor(color.NewWhite())

	return m
}
func (s Skin) Table(m pdf.Maroto, head []string, content [][]string, col uint, h float64) pdf.Maroto {
	s.TableHeader(m, head, false, col, h, consts.Left, 1)
	for _, v := range content {
		s.TableRow(m, v, false, col, h, 1, consts.Left)

	}
	return m
}
func (s Skin) TableLine(m pdf.Maroto, head []string, content [][]string) pdf.Maroto {
	s.TableHeader(m, head, true, 3, s.rowtableHeight+2, consts.Center, 0)
	for _, v := range content {
		s.TableRow(m, v, true, 3, s.rowtableHeight, 0, consts.Center)

	}
	return m
}
func (s Skin) Sign(m pdf.Maroto, name string, label string, id string, isTag bool) pdf.Maroto {

	prop := props.Text{

		Top:    6,
		Size:   6,
		Style:  consts.Normal,
		Align:  consts.Center,
		Color:  color.NewBlack(),
		Family: consts.Arial,
	}

	m.Row(12, func() {
		m.Col(4, func() {

		})
		m.Col(2, func() {

		})
		m.Col(6, func() {

			m.Signature(name, props.Font{
				Size:   12.0,
				Style:  consts.BoldItalic,
				Family: consts.Courier,
				Color:  s.TextColor,
			})
		})

	})
	m.Row(4, func() {
		m.Col(6, func() {

		})
		m.Col(6, func() {
			if isTag {
				m.Text("[[!sigField"+id+":signer1:signature(sigType=\"Click2Sign\"):label(\"firma qui\"):size(width=150,height=60)]]", prop)
			}

		})

	})

	return m
}
func (s Skin) SignDouleLine(m pdf.Maroto, name string, name2 string, id string, isTag bool) pdf.Maroto {

	prop := props.Text{

		Top:    2,
		Size:   6,
		Style:  consts.Normal,
		Align:  consts.Center,
		Color:  color.NewBlack(),
		Family: consts.Arial,
	}
	signProps := props.Font{
		Size:   11.0,
		Style:  consts.Normal,
		Family: consts.Courier,
		Color:  s.TextColor,
	}

	m.Row(15, func() {
		m.Col(6, func() {
			m.Signature(name, signProps)

		})
		m.Col(6, func() {

			_ = m.FileImage(lib.GetAssetPathByEnvV2()+"signature_global.png", props.Rect{
				Left:    20,
				Top:     -20,
				Center:  false,
				Percent: 100,
			})

			m.Signature(name2, signProps)
		})

	})
	m.Row(10, func() {
		m.Col(8, func() {
			if isTag {
				m.Text("[[!sigField"+id+":signer1:signature(sigType=\"Click2Sign\"):label(\"Firma qui\"):size(width=150,height=60)]]", prop)
			}

		})
		m.Col(4, func() {

		})

	})

	return m
}
func (s Skin) TitleSub(m pdf.Maroto, title string, subtitle string, body string) pdf.Maroto {
	prop := props.Text{
		Top:   1,
		Size:  s.SizeTitle,
		Style: consts.Bold,
		Align: consts.Left,
		Color: s.LineColor,
	}
	bold := props.Text{
		Top:   1,
		Size:  s.Size,
		Style: consts.Bold,
		Align: consts.Left,
		Color: s.TextColor,
	}
	normal := props.Text{
		Top:   1,
		Size:  s.Size,
		Style: consts.Normal,
		Align: consts.Left,
		Color: s.TextColor,
	}
	m.Row(s.RowHeight, func() {
		m.Col(12, func() {
			m.Text(title, prop)

		})

		//m.SetBackgroundColor(magenta)
	})
	m.Row(s.RowHeight, func() {
		m.Col(12, func() {
			m.Text(subtitle, bold)

		})

		//m.SetBackgroundColor(magenta)
	})
	m.Row(s.RowHeight, func() {
		m.Col(12, func() {
			m.Text(body, normal)

		})

		//m.SetBackgroundColor(magenta)
	})
	return m
}
func (s Skin) Title(m pdf.Maroto, title string, body string, bodyHeight float64) pdf.Maroto {
	prop := props.Text{
		Top:   1.5,
		Size:  s.SizeTitle,
		Style: consts.Bold,
		Align: consts.Left,
		Color: s.LineColor,
	}

	normal := props.Text{
		Top:   1.5,
		Size:  s.Size,
		Style: consts.Normal,
		Align: consts.Left,
		Color: s.TextColor,
	}
	m.Row(s.RowHeight, func() {
		m.Col(12, func() {
			m.Text(title, prop)

		})

	})
	m.Row(bodyHeight, func() {
		m.Col(10, func() {
			m.Text(body, normal)

		})

	})
	return m
}
func (s Skin) TitleBlack(m pdf.Maroto, title string, body string, bodyHeight float64) pdf.Maroto {
	prop := props.Text{
		Top:   1.5,
		Size:  s.SizeTitle,
		Style: consts.Bold,
		Align: consts.Left,
		Color: s.TextColor,
	}

	normal := props.Text{
		Top:   1.5,
		Size:  s.Size,
		Style: consts.Normal,
		Align: consts.Left,
		Color: s.TextColor,
	}
	m.Row(s.RowHeight, func() {
		m.Col(12, func() {
			m.Text(title, prop)

		})

	})
	m.Row(bodyHeight, func() {
		m.Col(10, func() {
			m.Text(body, normal)

		})

	})
	return m
}
