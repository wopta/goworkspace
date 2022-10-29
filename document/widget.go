package document

import (
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
	//model "github.com/wopta/goworkspace/models"
)

func (s Skin) Customer(m pdf.Maroto, customer []Kv) pdf.Maroto {
	for _, v := range customer {
		m = s.CustomerRow(m, v.Key, v.Value)
	}
	return m
}
func (s Skin) CoveragesTable(m pdf.Maroto, head []string, content [][]string) pdf.Maroto {
	s.TableHeader(m, head, true, 3, 5.0, consts.Center)
	for _, v := range content {
		s.TableRow(m, v, true, 3, s.rowtableHeight)

	}
	return m
}
func (s Skin) Table(m pdf.Maroto, head []string, content [][]string, col uint, h float64) pdf.Maroto {
	s.TableHeader(m, head, false, col, h, consts.Left)
	for _, v := range content {
		s.TableRow(m, v, false, col, h)

	}
	return m
}
func (s Skin) TableLine(m pdf.Maroto, head []string, content [][]string) pdf.Maroto {
	s.TableHeader(m, head, true, 4, s.rowtableHeight, consts.Center)
	for _, v := range content {
		s.TableRow(m, v, true, 4, s.rowtableHeight)

	}
	return m
}
func (s Skin) Stantements(m pdf.Maroto, title string, data []Kv) pdf.Maroto {
	prop := props.Text{
		Top:   1.5,
		Size:  s.Size,
		Style: consts.Normal,
		Align: consts.Left,
		Color: s.TextColor,
	}
	bold := props.Text{
		Top:   1.5,
		Size:  s.SizeTitle,
		Style: consts.Normal,
		Align: consts.Left,
		Color: s.TextColor,
	}

	m.Row(10, func() {
		m.Col(6, func() {
			m.Text(title, prop)

		})
		m.Col(4, func() {

			m.Text(" DICHIARO:  ", bold)
		})

		//m.SetBackgroundColor(magenta)
	})
	for _, v := range data {
		s.RowBullet(m, v.Key, v.Value, consts.Normal)

	}
	return m
}
func (s Skin) Articles(m pdf.Maroto, data []Kv) pdf.Maroto {
	textBold := props.Text{
		Top:   1.5,
		Size:  s.SizeTitle,
		Style: consts.Bold,
		Align: consts.Left,
		Color: s.TextColor,
	}
	textBoldMagenta := props.Text{
		Top:   1.5,
		Size:  s.SizeTitle,
		Style: consts.Bold,
		Align: consts.Left,
		Color: s.LineColor,
	}
	m.Row(s.rowHeight, func() {
		m.Col(12, func() {
			m.Text("Dichiarazioni da leggere con attenzione prima di firmare  ", textBoldMagenta)
		})
		//m.SetBackgroundColor(magenta)
	})
	m.Row(s.rowHeight, func() {
		m.Col(12, func() {
			m.Text("Premesso di essere a conoscenza che:  ", textBold)
		})
		//m.SetBackgroundColor(magenta)
	})
	for _, v := range data {
		s.RowBullet(m, v.Key, v.Value, consts.Bold)

	}
	return m
}
func (s Skin) Sign(m pdf.Maroto, name string, label string) pdf.Maroto {
	prop := props.Text{
		Top:   15,
		Size:  s.Size,
		Style: consts.Bold,
		Align: consts.Center,
		Color: s.TextColor,
	}
	m.Row(10, func() {
		m.Col(8, func() {
			m.Text("", prop)

		})
		m.Col(4, func() {
			m.Text(label+" "+name, prop)

		})

		//m.SetBackgroundColor(magenta)
	})

	m.Row(10, func() {
		m.Col(8, func() {
			m.Text("", prop)

		})

		m.Col(4, func() {

			m.Text("---------------", prop)
		})
		//m.SetBackgroundColor(magenta)
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
		Size:  s.SizeTitle,
		Style: consts.Bold,
		Align: consts.Left,
		Color: s.TextColor,
	}
	normal := props.Text{
		Top:   1,
		Size:  s.Size,
		Style: consts.Bold,
		Align: consts.Left,
		Color: s.TextColor,
	}
	m.Row(s.rowHeight, func() {
		m.Col(12, func() {
			m.Text(title, prop)

		})

		//m.SetBackgroundColor(magenta)
	})
	m.Row(s.rowHeight, func() {
		m.Col(12, func() {
			m.Text(subtitle, bold)

		})

		//m.SetBackgroundColor(magenta)
	})
	m.Row(s.rowHeight, func() {
		m.Col(12, func() {
			m.Text(body, normal)

		})

		//m.SetBackgroundColor(magenta)
	})
	return m
}
func (s Skin) Title(m pdf.Maroto, title string, body string) pdf.Maroto {
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
		Style: consts.Bold,
		Align: consts.Left,
		Color: s.TextColor,
	}
	m.Row(s.rowHeight+2.0, func() {
		m.Col(12, func() {
			m.Text(title, prop)

		})

	})
	m.Row(s.rowHeight, func() {
		m.Col(10, func() {
			m.Text(body, normal)

		})

	})
	return m
}
func (s Skin) TitleList(m pdf.Maroto, title string, data []string) pdf.Maroto {
	textBold := props.Text{
		Top:   1.5,
		Size:  s.SizeTitle,
		Style: consts.Bold,
		Align: consts.Left,
		Color: s.TextColor,
	}

	m.Row(s.rowHeight, func() {
		m.Col(12, func() {
			m.Text(title, textBold)
		})
		//m.SetBackgroundColor(magenta)
	})
	for _, v := range data {
		s.RowCol1(m, v, consts.Bold)

	}
	return m
}
func (s Skin) TitleBulletList(m pdf.Maroto, title string, data []Kv) pdf.Maroto {
	textBold := props.Text{
		Top:   1.5,
		Size:  s.SizeTitle,
		Style: consts.Bold,
		Align: consts.Left,
		Color: s.TextColor,
	}

	m.Row(s.rowHeight, func() {
		m.Col(12, func() {
			m.Text(title, textBold)
		})
		//m.SetBackgroundColor(magenta)
	})
	for _, v := range data {
		s.RowBullet(m, v.Key, v.Value, consts.Bold)

	}
	return m
}
func (s Skin) BulletList(m pdf.Maroto, content []Kv) pdf.Maroto {

	for _, v := range content {
		s.RowBullet(m, v.Key, v.Value, consts.Normal)

	}
	return m
}
func (s Skin) AboutUs(m pdf.Maroto, title string, sub []Kv) pdf.Maroto {
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
		Style: consts.Bold,
		Align: consts.Left,
		Color: s.TextColor,
	}
	m.Row(10, func() {
		m.Col(12, func() {
			m.Text(title, prop)

		})

		//m.SetBackgroundColor(magenta)
	})
	for _, k := range sub {
		m.Row(s.rowHeight, func() {
			m.Col(12, func() {
				m.Text(k.Key, bold)

			})

			//m.SetBackgroundColor(magenta)
		})
		m.Row(s.rowHeight+2, func() {
			m.Col(12, func() {
				m.Text(k.Value, normal)

			})

			//m.SetBackgroundColor(magenta)
		})
	}

	return m
}
