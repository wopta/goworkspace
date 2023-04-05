package document

import (
	"strconv"

	"github.com/dustin/go-humanize"
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

func (s Skin) Question(m pdf.Maroto, data models.Question) {
	var prop props.Text
	var rh float64
	if data.IsBold {
		prop = s.BoldtextLeft
		rh = s.RowHeight + 0.0
	} else {
		prop = s.NormaltextLeft
		rh = s.RowHeight
	}
	if data.Indent {
		m.Row(s.getRowHeight(data.Question, s.CharForRow, rh), func() {
			m.ColSpace(1)
			m.Col(11, func() {
				m.Text(data.Question, prop)

			})

		})
	} else {
		m.Row(s.getRowHeight(data.Question, s.CharForRow, rh), func() {
			m.Col(12, func() {
				m.Text(data.Question, prop)

			})
		})

		m = s.Space(m, 0.3)

	}

}

func (s Skin) Stantement(m pdf.Maroto, title string, data models.Statement) {
	//d := data.Questions
	m.Row(s.RowTitleHeight, func() {
		m.Col(12, func() {
			m.Text(title, s.MagentaBoldtextLeft)

		})

		//m.SetBackgroundColor(magenta)
	})
	for _, v := range data.Questions {
		//qlen := len(*data.Questions)
		//nextI := 1
		//s.checkPageNext(m, (*d)[nextI].Question)
		s.Question(m, *v)

	}

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
	m.Row(s.RowHeight, func() {
		m.Col(12, func() {
			m.Text("Dichiarazioni da leggere con attenzione prima di firmare  ", textBoldMagenta)
		})
		//m.SetBackgroundColor(magenta)
	})
	m.Row(s.RowHeight, func() {
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

func (s Skin) TitleList(m pdf.Maroto, title string, data []string) pdf.Maroto {
	textBold := props.Text{
		Top:   1.5,
		Size:  s.SizeTitle,
		Style: consts.Bold,
		Align: consts.Left,
		Color: s.TextColor,
	}

	m.Row(s.RowHeight, func() {
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

	m.Row(s.RowHeight, func() {
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

	m.Row(s.RowHeight, func() {
		m.Col(12, func() {
			m.Text(title, s.MagentaBoldtextLeft)

		})

	})
	for _, k := range sub {

		m.Row(10, func() {
			m.Col(12, func() {
				m.Text(k.Key, s.NormaltextLeft)
				m.Text(k.Value, s.NormaltextLeftBlack)
			})

		})
	}

	return m
}
func (s Skin) OverBold(m pdf.Maroto, sub []Kv, prop props.Text) pdf.Maroto {

	for _, k := range sub {

		m.Row(10, func() {
			m.Col(12, func() {
				m.Text(k.Key, s.NormaltextLeft)
				m.Text(k.Value, s.NormaltextLeftBlack)
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
	var (
		deductibleValue      string
		constructionYear     string
		constructionMaterial string
		floor                string
	)
	m.Row(s.RowTitleHeight, func() {
		m.Col(8, func() {
			m.Text("La tua assicurazione è prestata per", s.MagentaBoldtextLeft)
		})
	})
	m.Line(1.0, s.Line)
	//buildingcount := 0
	enterprise := GetEnterprise(data.Assets)
	ateco := enterprise.Ateco
	atecodesc := enterprise.AtecoDesc

	for _, A := range data.Assets {
		var (
			c [][]string
			d []string
		)
		if A.Building != nil {

			alarm := "senza"
			holder := "in affitto"
			build := A.Building
			if build.IsHolder {
				holder = "di proprietà "
			}
			if build.IsAllarm {
				alarm = "con"
			} else {
				alarm = "senza"
			}

			switch os := build.BuildingYear; os {
			case "before1972":
				constructionYear = "prima del 1972"
			case "1972between2009":
				constructionYear = "tra 1972 e 2009"
			case "after2009":
				constructionYear = "dopo il 2009"

			}
			switch os := build.Floor; os {
			case "ground_floor":
				floor = "solo piano terra"
			case "first":
				floor = "un piano"
			case "second":
				floor = "di 2 piani"
			case "greater_than_second":
				floor = "di oltre 2 piani"

			}

			switch os := build.BuildingMaterial; os {
			case "masonry":
				constructionMaterial = "muratura o CA in appoggio"
			case "reinforcedConcrete":
				constructionMaterial = "Cemento armato legato"
			case "antiSeismicLaminatedTimber":
				constructionMaterial = "Legno lamellare"
			case "steel":
				constructionMaterial = "Acciaio"

			}
			c = append(c, []string{"", "", "Sede "})
			d = append(d, build.Address+" "+build.StreetNumber+" - "+build.PostalCode+" "+build.City+" ("+build.CityCode+")")
			d = append(d, "Fabbricato "+constructionMaterial+" construito "+constructionYear+", "+floor+", "+alarm+" antifurto, "+holder)
			d = append(d, "Attività ATECO codice: "+ateco)

			if len(atecodesc) > 100 {
				d = append(d, "Descrizione: "+atecodesc[:100])
				d = append(d, atecodesc[100:])
			} else if len(atecodesc) > 200 {
				d = append(d, "Descrizione: "+atecodesc[:100])
				d = append(d, atecodesc[100:200])
				d = append(d, atecodesc[200:])
			} else {
				d = append(d, "Descrizione: "+atecodesc)
			}
			c = append(c, d)

		}
		if A.Enterprise != nil {

			e := A.Enterprise
			rev, err := strconv.Atoi(e.Revenue)
			lib.CheckError(err)
			c = append(c, []string{"", "", "Attività"})
			d = append(d, "Fatturato: € "+humanize.FormatInteger("#.###,", int(rev)))
			d = append(d, "Addetti nr: "+strconv.Itoa(int(e.Employer)))
			d = append(d, "Attività ATECO codice: "+ateco)
			if len(atecodesc) > 100 {
				d = append(d, "Descrizione: "+atecodesc[:100])
				d = append(d, atecodesc[100:])
			} else if len(atecodesc) > 200 {
				d = append(d, "Descrizione: "+atecodesc[:100])
				d = append(d, atecodesc[100:200])
				d = append(d, atecodesc[200:])

			} else {
				d = append(d, "Descrizione: "+atecodesc)
			}
			d = append(d, "Ubicazione Attività: Ogni sede sede assicurata indicata in polizza")
			c = append(c, d)
		}

		for _, k := range A.Guarantees {
			if k.Slug == "third-party-liability" {
				switch os := k.Deductible; os {
				case "0":
					deductibleValue = "MINIMO"
				case "500":
					deductibleValue = "BASSO"
				case "1000":
					deductibleValue = "MEDIO-BASSO"
				case "1500":
					deductibleValue = "MEDIO-BASSO"
				case "2000":
					deductibleValue = "MEDIO-ALTO"
				case "3000":
					deductibleValue = "ALTO"
				case "5000":
					deductibleValue = "MASSIMO"
				}
			}

		}
		m = s.MultiRow(m, c, true, []uint{2, 10}, 30)
	}

	m.Row(20, func() {
		m.Col(2, func() {
			m.Text("Franchigia e Scoperto  ", s.BoldtextLeft)
		})
		m.Col(8, func() {
			m.Text("Il livello scelto è: "+deductibleValue+". Per ogni garanzia, nella Tabella “Scoperti e Franchigie” alla voce "+deductibleValue+" troverai il dettaglio di tutti gli Scoperti e Franchigie in caso di Sinistro, di cui l’importo qui indicato costituisce, in ogni caso, il minimo applicato se non diversamente specificato", s.NormaltextLeft)
			m.Text("                               "+deductibleValue+"                                                                                                                                                                                                                                                                                                                   ", s.NormaltextLeftBlack)
		})
	})
	return m

}

// ----------------------------------------------------------------------------------

func (s Skin) CoveragesPersonTable(m pdf.Maroto, data models.Policy) pdf.Maroto {
	h := []string{"Garanzie ", "Somma assicurata ", "Opzioni / Dettagli ", "Premio "}
	var table [][]string
	for _, A := range data.Assets {
		for _, k := range A.Guarantees {
			r := []string{k.Name, strconv.Itoa(int(k.SumInsuredLimitOfIndemnity)), k.SelfInsurance, strconv.Itoa(int(k.PriceGross))}
			table = append(table, r)
		}
	}

	m = s.CoveragesTable(m, h, table)
	return m
}
func (skin Skin) GetHeader(m pdf.Maroto, data models.Policy, logo string, nameProd string) pdf.Maroto {

	h := []string{"I dati della tua Polizza ", "I tuoi dati"}
	layout := "02/01/2006"
	m.RegisterHeader(func() {
		var (
			tablePremium [][]string
			nextpay      string
			cfpi         string
			tie          func()
			name         string
		)
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

				m.Text(nameProd, props.Text{
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
		skin.Space(m, 5.0)

		if data.PaymentSplit == "montly" {
			nextpay = data.StartDate.AddDate(0, 1, 0).Format(layout)
		} else {
			nextpay = data.EndDate.Format(layout)
		}

		if data.Contractor.VatCode == "" {
			cfpi = data.Contractor.FiscalCode
		} else {
			cfpi = data.Contractor.VatCode
		}
		if data.Name == "pmi" {
			tie = func() {

				tablePremium = append(tablePremium, []string{"Presenza Vincolo: NO  - Convenzione: NO", ""})
			}

			for _, as := range data.Assets {
				if as.Enterprise != nil {

					name = as.Enterprise.Name
					break
				} else {
					name = data.Contractor.Name + " " + data.Contractor.Surname
				}
			}
		}
		tablePremium = append(tablePremium, []string{"Numero: " + data.CodeCompany, "Contraente: " + name})
		tablePremium = append(tablePremium, []string{"Decorre dal: " + data.StartDate.Format(layout) + " ore 24:00", "C.F. / P.IVA: " + cfpi})
		tablePremium = append(tablePremium, []string{"Scade il: " + data.EndDate.Format(layout) + " ore 24:00", "Indirizzo: " + data.Contractor.Address + " " + data.Contractor.StreetNumber})
		tablePremium = append(tablePremium, []string{"Si rinnova a scadenza, salvo disdetta da inviare 30 giorni prima", data.Contractor.PostalCode + " " + data.Contractor.City + " (" + data.Contractor.CityCode + ")"})
		tablePremium = append(tablePremium, []string{"Prossimo pagamento il: " + nextpay, "Mail:  " + data.Contractor.Mail})
		tablePremium = append(tablePremium, []string{"Sostituisce la polizza: = = = = = = = =", "Telefono: " + data.Contractor.Phone})
		tie()
		//tablePremium = append(tablePremium, []string{"Presenza Vincolo: NO  - Convenzione: NO", ""})
		skin.Table(m, h, tablePremium, 6, skin.RowHeight-0.5)
		skin.Space(m, 5.0)
		//m.Line(skin.LineHeight, skin.Lin
	})

	return m
}
func (skin Skin) GetFooter(m pdf.Maroto, logo string, name string) pdf.Maroto {
	m.RegisterFooter(func() {
		m.Row(15.0, func() {
			m.Col(8, func() {
				m.Text(name, props.Text{
					Top:         25,
					Style:       consts.Bold,
					Align:       consts.Left,
					Color:       skin.LineColor,
					Size:        skin.Size - 1,
					Extrapolate: false,
				})
			})
			m.Col(2, func() {
				_ = m.FileImage(lib.GetAssetPathByEnv("document")+logo, props.Rect{
					Left:    5,
					Top:     10,
					Center:  false,
					Percent: 100,
				})
			})
		})
	})
	return m
}
