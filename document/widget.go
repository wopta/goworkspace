package document

import (
	"strconv"

	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
	lib "github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	//model "github.com/wopta/goworkspace/models"
)

func (s Skin) Customer(m pdf.Maroto, customer []Kv) pdf.Maroto {
	for _, v := range customer {
		m = s.CustomerRow(m, v.Key, v.Value)
	}
	return m
}
func (s Skin) CoveragesTable(m pdf.Maroto, head []string, content [][]string) pdf.Maroto {
	s.TableHeader(m, head, true, 3, s.rowtableHeight+2, consts.Center, 0)
	for _, v := range content {
		s.TableRow(m, v, true, 3, s.rowtableHeight, 0, consts.Center)

	}
	return m
}
func (s Skin) Table(m pdf.Maroto, head []string, content [][]string, col uint, h float64) pdf.Maroto {
	s.TableHeader(m, head, false, col, h+2, consts.Left, 1)
	for _, v := range content {
		s.TableRow(m, v, false, col, h, 1, consts.Left)

	}
	return m
}
func (s Skin) TableLine(m pdf.Maroto, head []string, content [][]string) pdf.Maroto {
	s.TableHeader(m, head, true, 4, s.rowtableHeight+2, consts.Center, 0)
	for _, v := range content {
		s.TableRow(m, v, true, 4, s.rowtableHeight, 0, consts.Center)

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
		Style: consts.Bold,
		Align: consts.Left,
		Color: s.TextColor,
	}
	m.Row(s.rowHeight+2.0, func() {
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

	})
	for _, k := range sub {
		m.Row(15, func() {
			m.Col(12, func() {
				m.Text(k.Key, bold)

			})

			//m.SetBackgroundColor(magenta)
		})
		m.Row(25, func() {
			m.Col(12, func() {
				m.Text(k.Value, normal)

			})

		})
	}

	return m
}
func (s Skin) GetPersona(data models.Policy, m pdf.Maroto) pdf.Maroto {
	linePropMagenta := props.Line{
		Color: s.LineColor,
		Style: consts.Solid,
		Width: 0.2,
	}
	m.Row(10, func() {
		m.Col(12, func() {
			m.Text("La tua assicurazione è operante per il seguente Assicurato e Garanzie ", props.Text{
				Color: s.LineColor,
				Top:   5,
				Style: consts.Bold,
				Align: consts.Left,
				Size:  s.SizeTitle,
			})
		})

	})
	m.Line(1.0, linePropMagenta)
	customer := []Kv{
		{
			Key:   "Assicurato: ",
			Value: "1"},
		{
			Key:   "Cognome e Nome: ",
			Value: data.Contractor.Name + " " + data.Contractor.Surname},
		{
			Key:   "Codice Fiscale: ",
			Value: data.Contractor.FiscalCode},
		{
			Key:   "Professione: ",
			Value: data.Contractor.Work},
		{
			Key:   "Tipo professione: ",
			Value: data.Contractor.WorkType},
		{
			Key:   "Classe rischio: ",
			Value: data.Contractor.RiskClass},
		{
			Key:   "Forma di copertura: ",
			Value: data.CoverageType},
	}
	m = s.Customer(m, customer)
	return m

}
func (s Skin) GetPmi(data models.Policy, m pdf.Maroto) pdf.Maroto {

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

	linePropMagenta := props.Line{
		Color: s.LineColor,
		Style: consts.Solid,
		Width: 0.2,
	}
	m.Row(20, func() {
		m.Col(8, func() {
			m.Text("Le coperture assicurative che hai scelto ", prop)
		})
	})

	m.Line(1.0, linePropMagenta)
	m.Row(20, func() {
		m.Col(2, func() {
			m.Text("Attività", bold)

		})
		m.Col(8, func() {
			m.Text(`title title`, normal)
			m.Text(`title title`, props.Text{
				Top:         6,
				Style:       consts.Italic,
				Align:       consts.Left,
				Color:       s.TextColor,
				Size:        s.Size,
				Extrapolate: true,
			})

		})

	})
	m.Line(1.0, linePropMagenta)
	m.Row(20, func() {
		m.Col(2, func() {
			m.Text("Sede 1", bold)

		})
		m.Col(8, func() {
			m.Text("title", normal)

		})

	})

	m.Line(1.0, linePropMagenta)
	m.Row(20, func() {
		m.Col(2, func() {
			m.Text("Franchigia e Scoperto  ", bold)

		})
		m.Col(8, func() {
			m.Text("Il livello scelto è: MINIMO. Per ogni garanzia, nella Tabella “Scoperti e Franchigie” alla voce MINIMO troverai il dettaglio di tutti gli Scoperti e Franchigie in caso di Sinistro, di cui l’importo qui indicato costituisce, in ogni caso, il minimo applicato se non diversamente specificato", normal)

		})

	})
	return m

}
func (s Skin) CoveragesPmiTable(m pdf.Maroto, data models.Policy) pdf.Maroto {

	m.Row(25, func() {
		m.Col(12, func() {
			m.Text("Le coperture assicurative che hai scelto ", s.MagentaBoldtextLeft)

		})

	})
	m.Row(25, func() {
		m.Col(12, func() {
			m.Text("(operative se indicata Somma o Massimale e secondo le Opzioni/Estensioni attivate qualora indicato) ", s.MagentaBoldtextLeft)

		})

	})

	head1 := []string{"Garanzie ", "Somma assicurata ", "Opzioni / Estensioni ", "Premio €"}
	var table [][]string
	for _, A := range data.Assets {
		for _, k := range A.Guarantees {
			r := []string{k.Name, strconv.Itoa(int(k.SumInsuredLimitOfIndemnity)), k.SelfInsurance, strconv.Itoa(int(k.Price))}
			table = append(table, r)
		}
	}
	m = s.CoveragesTable(m, head1, table)

	var table2 [][]string
	head2 := []string{"Garanzie ", "Somma assicurata ", "Opzioni / Dettagli ", "Premio €"}
	for _, A := range data.Assets {
		for _, k := range A.Guarantees {
			r := []string{k.Name, strconv.Itoa(int(k.SumInsuredLimitOfIndemnity)), k.SelfInsurance, strconv.Itoa(int(k.Price))}
			table = append(table, r)
		}
	}
	m = s.CoveragesTable(m, head2, table2)
	return m
}
func (s Skin) CoveragesPersonTable(m pdf.Maroto, data models.Policy) pdf.Maroto {
	h := []string{"Garanzie ", "Somma assicurata ", "Opzioni / Dettagli ", "Premio "}
	var table [][]string
	for _, A := range data.Assets {
		for _, k := range A.Guarantees {
			r := []string{k.Name, strconv.Itoa(int(k.SumInsuredLimitOfIndemnity)), k.SelfInsurance, strconv.Itoa(int(k.Price))}
			table = append(table, r)
		}
	}

	m = s.CoveragesTable(m, h, table)
	return m
}
func (skin Skin) GetHeader(m pdf.Maroto, data models.Policy, logo string, name string) pdf.Maroto {
	m.RegisterHeader(func() {
		m.Row(15.0, func() {
			m.Col(2, func() {

				_ = m.FileImage(lib.GetAssetPathByEnv("document")+logo, props.Rect{
					Left:    1,
					Top:     1,
					Center:  false,
					Percent: 100,
				})
			})
			m.Col(1, func() {
				m.Text("Wopta per te", props.Text{
					Color:       skin.LineColor,
					Top:         1,
					Style:       consts.Bold,
					Align:       consts.Left,
					Size:        skin.SizeTitle + 3,
					Extrapolate: true,
				})

				m.Text(name, props.Text{
					Top:         6,
					Style:       consts.Italic,
					Align:       consts.Left,
					Color:       skin.TextColor,
					Size:        skin.SizeTitle + 3,
					Extrapolate: true,
				})
			})
			m.ColSpace(6)
			m.Col(2, func() {
				_ = m.FileImage(lib.GetAssetPathByEnv("document")+"/ARTW_LOGO_RGB_400px.png", props.Rect{
					Left:    1,
					Top:     1,
					Center:  false,
					Percent: 100,
				})
			})
		})
		h := []string{"I dati della tua Polizza ", "I tuoi dati"}
		var tablePremium [][]string
		tablePremium = append(tablePremium, []string{"Numero: " + data.ID, "Contraente: " + data.Contractor.Name + " " + data.Contractor.Surname})
		tablePremium = append(tablePremium, []string{"Decorre dal: " + data.StartDate.String() + " ore 24:00", "C.F. / P.IVA: " + data.Contractor.Surname})
		tablePremium = append(tablePremium, []string{"Scade il: " + data.EndDate.String() + " ore 24:00", "Indirizzo: " + data.Contractor.Address})
		tablePremium = append(tablePremium, []string{"Si rinnova a scadenza, salvo disdetta da inviare 30 giorni prima", "XXXXX  XXXXXXXXXXXXXXXXXXX (XX)"})
		tablePremium = append(tablePremium, []string{"Prossimo pagamento il: " + data.EndDate.String(), "Mail:  " + data.Contractor.Mail})
		tablePremium = append(tablePremium, []string{"Sostituisce la polizza: = = = = = = = =", "Telefono: " + data.Contractor.Phone})
		m = skin.Table(m, h, tablePremium, 6, 3.0)
	})

	return m
}
func (skin Skin) GetFooter(m pdf.Maroto, data models.Policy, logo string, name string) pdf.Maroto {
	m.RegisterFooter(func() {
		m.Row(15.0, func() {
			m.Col(8, func() {
				m.Text("Wopta per te. Persona è un prodotto assicurativo di Global Assistance Compagnia di assicurazioni e riassicurazioni S.p.A, distribuito da Wopta Assicurazioni S.r.l", props.Text{
					Top:         1,
					Style:       consts.Bold,
					Align:       consts.Left,
					Color:       skin.LineColor,
					Size:        skin.Size - 1,
					Extrapolate: false,
				})
			})
			m.Col(2, func() {
				_ = m.FileImage(lib.GetAssetPathByEnv("document")+"/logo_global.png", props.Rect{
					Left:    1,
					Top:     1,
					Center:  false,
					Percent: 100,
				})
			})
		})
	})
	return m
}
