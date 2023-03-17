package document

import (
	"fmt"
	"io/ioutil"
	"log"
	"sort"

	"os"
	"strconv"

	"github.com/dustin/go-humanize"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
	lib "github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	p "github.com/wopta/goworkspace/product"
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

func (s Skin) Stantements(m pdf.Maroto, title string, data []Kv) pdf.Maroto {
	prop := props.Text{
		Top:   1.5,
		Size:  s.Size,
		Style: consts.Bold,
		Align: consts.Left,
		Color: s.TextColor,
	}

	m.Row(s.RowTitleHeight, func() {
		m.Col(6, func() {
			m.Text(title, prop)

		})
		m.Col(4, func() {

		})

		//m.SetBackgroundColor(magenta)
	})
	for _, v := range data {
		s.RowBullet(m, v.Key, v.Value, consts.Normal)

	}
	return m
}
func (s Skin) Stantement(m pdf.Maroto, title string, data models.Statement) pdf.Maroto {
	d := data.Questions
	m.Row(s.RowTitleHeight, func() {
		m.Col(12, func() {
			m.Text(title, s.MagentaBoldtextLeft)

		})

		//m.SetBackgroundColor(magenta)
	})
	for i, v := range *data.Questions {
		qlen := len(*data.Questions)
		nextI := 1
		s.checkPageNext(m, (*d)[nextI].Question)
		if qlen == i-1 {
			nextI = i - 1
		} else {
			nextI = i
		}
		var prop props.Text
		var rh float64
		if v.Isbold {
			prop = s.BoldtextLeft
			rh = s.RowHeight + 0.0
		} else {
			prop = s.NormaltextLeft
			rh = s.RowHeight
		}
		m.Row(s.getRowHeight(v.Question, s.CharForRow, rh), func() {
			m.Col(12, func() {
				m.Text(v.Question, prop)

			})

			//m.SetBackgroundColor(magenta)
		})
		m = s.Space(m, 0.3)

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
	var deductibleValue string
	var constructionYear string
	var constructionMaterial string
	m.Row(s.RowTitleHeight, func() {
		m.Col(8, func() {
			m.Text("L’assicurazione è prestata per", s.MagentaBoldtextLeft)
		})
	})
	m.Line(1.0, s.Line)
	//buildingcount := 0
	for _, A := range data.Assets {
		c := [][]string{{""}, {""}, {""}, {""}}
		if A.Building != nil {
			alarm := "senza"
			holder := "in affitto"
			build := A.Building
			if build.IsHolder {
				holder = "di proprietà "
			}
			if build.IsAllarm {
				alarm = "con"
			}

			switch os := build.BuildingYear; os {
			case "before1972":
				constructionYear = "prima dle 1972"
			case "1972between2009":
				constructionYear = "tra 1972 e 2009"
			case "after2009":
				constructionYear = "dopo il 2009"

			}
			switch os := build.BuildingMaterial; os {
			case "masonry":
				constructionMaterial = "muratura"
			case "reinforcedConcrete":
				constructionMaterial = " CA in appoggio"
			case "antiSeismicLaminatedTimber":
				constructionMaterial = "Legno lamellare"
			case "steel":
				constructionMaterial = "Acciaio"

			}

			c = [][]string{{"", "Sede "},
				{build.Address,
					"Fabbricato " + constructionMaterial + "construito " + constructionYear + ", " + alarm + " antifurto, " + holder,
					"Attività ATECO codice: " + build.Ateco,
					"Descrizione: " + build.AtecoDesc}}

		}

		if A.Enterprise != nil {
			e := A.Enterprise
			c = [][]string{{"", "", "Attivita"},
				{"Fatturato: € " + e.Revenue,
					"Adetti nr:" + strconv.Itoa(int(e.Employer)),
					"Attività ATECO codice: " + e.Ateco,
					"Descrizione: " + e.AtecoDesc,
					"Ubicazione Attività: : " + e.Address}}
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
func (s Skin) CoveragesPmiTable(m pdf.Maroto, data models.Policy) pdf.Maroto {
	existGuarance := func(str bool) string {
		res := "NO"
		if str {
			res = "SI"
		}
		return res
	}
	textS := s.MagentaBoldtextLeft
	textS.Size = textS.Size - 3
	m.Row(s.RowTitleHeight, func() {
		m.Col(12, func() {
			m.Text("Le coperture assicurative che hai scelto ", s.MagentaBoldtextLeft)

		})
	})
	m.Row(s.RowTitleHeight-1, func() {
		m.Col(12, func() {
			m.Text("(operative se indicata Somma o Massimale e secondo le Opzioni/Estensioni attivate qualora indicato) ", textS)
		})
	})

	head1 := []string{"Garanzie ", "Somma assicurata ", "Opzioni / Estensioni ", "Premio "}
	head2 := []string{"Garanzie ", "Somma assicurata ", "Opzioni / Dettagli ", "Premio "}
	//var table [][][]string
	var product models.Product
	var e error
	if os.Getenv("env") == "local" {
		productFile := lib.ErrorByte(ioutil.ReadFile("function-data/products/" + data.Name + ".json"))
		product, e = models.UnmarshalProduct(productFile)
		lib.CheckError(e)

	} else {
		product, e = p.GetName(data.Name, "v"+fmt.Sprint(data.ProductVersion))
		lib.CheckError(e)
	}

	for _, A := range data.Assets {
		//m = s.Space(m, 10.0)
		s.checkPage(m)
		var isOptionShow bool
		mapg := make(map[string][][]string)
		mapprice := make(map[string]float64)
		var firecount int
		sort.Slice(A.Guarantees, func(i, j int) bool {

			return product.Companies[0].GuaranteesMap[A.Guarantees[i].Slug].OrderAsset < product.Companies[0].GuaranteesMap[A.Guarantees[j].Slug].OrderAsset
		})
		log.Println("-----------------------------------------------------------------------------------------")

		for _, k := range A.Guarantees {

			guarance := product.Companies[0].GuaranteesMap[k.Slug]
			group := guarance.Group

			if len(mapg[group]) == 0 {

				mapg[group] = [][]string{{}, {}, {}, {}}
			}

			if group == "RCT" {
				mapg[group][2] = make([]string, 8)
				if !guarance.IsExtension {
					mapg[group][0] = make([]string, 8)
					mapg[group][1] = make([]string, 8)
					mapg[group][0][4] = guarance.CompanyName
					mapg[group][1][4] = "€ " + humanize.FormatInteger("#.###,", int(k.SumInsuredLimitOfIndemnity))
				}
				mapg[group][2][0] = "Sono attive le seguenti opzioni / estensioni:"
				mapg[group][2][1] = "Danni ai veicoli in consegna e custodia: " + existGuarance(ExistGuarance(A.Guarantees, "damage-to-goods-in-custody"))
				mapg[group][2][2] = "Responsabilità civile postuma officine: " + existGuarance(ExistGuarance(A.Guarantees, "defect-liability-workmanships"))
				mapg[group][2][3] = "Responsabilità civile postuma 12 mesi: " + existGuarance(ExistGuarance(A.Guarantees, "defect-liability-12-months"))
				mapg[group][2][4] = "Responsabilità civile postuma D.M.37/2008: " + existGuarance(ExistGuarance(A.Guarantees, "defect-liability-dm-37-2008"))
				mapg[group][2][5] = "Danni da furto: " + existGuarance(ExistGuarance(A.Guarantees, "property-damage-due-to-theft"))
				mapg[group][2][6] = "Danni alle cose sulle quali si eseguono i lavori: " + existGuarance(ExistGuarance(A.Guarantees, "damage-to-goods-course-of-works"))
				mapg[group][2][7] = "RC impresa edile: " + existGuarance(ExistGuarance(A.Guarantees, "third-party-liability-construction-company"))
				//mapg[group][2] = append(mapg[group][2], k.CompanyName+": "+existGuarance(ExistGuarance(A.Guarantees, k.Slug)))
			} else if group == "LEGAL" {
				var SumInsuredLimitOfIndemnity float64
				var detail string
				if k.LegalDefence == "basic" {
					SumInsuredLimitOfIndemnity = 10000
					detail = "Difesa Penale"
				} else {
					SumInsuredLimitOfIndemnity = 25000
					detail = "Difesa Penale Difesa Civile Circolazione"
				}

				mapg[group][0] = append(mapg[group][0], guarance.CompanyName)
				mapg[group][1] = append(mapg[group][1], "€ "+humanize.FormatInteger("#.###,", int(SumInsuredLimitOfIndemnity)))

				mapg[group][2] = append(mapg[group][2], "E’ attiva la garanzia:")
				mapg[group][2] = append(mapg[group][2], detail)

				if k.LegalDefence == "extended" {
					mapg[group][2] = append(mapg[group][2], "Difesa Penale, Civile e Circolazione ")
				}

			} else if group == "FIRE" && !guarance.IsExtension {
				if firecount == 0 {
					mapg[group][0] = append(mapg[group][0], "")
					mapg[group][1] = append(mapg[group][1], "")
					mapg[group][0] = append(mapg[group][0], "")
					mapg[group][1] = append(mapg[group][1], "")
					mapg[group][0] = append(mapg[group][0], "")
					mapg[group][1] = append(mapg[group][1], "")
					mapg[group][0] = append(mapg[group][0], "")
					mapg[group][1] = append(mapg[group][1], "")
					mapg[group][0] = append(mapg[group][0], "")
					mapg[group][1] = append(mapg[group][1], "")
					firecount++
				}
				mapg[group][0] = append(mapg[group][0], guarance.CompanyName)
				mapg[group][1] = append(mapg[group][1], "€ "+humanize.FormatInteger("#.###,", int(k.SumInsuredLimitOfIndemnity)))

				if !isOptionShow {
					mapg[group][2] = append(mapg[group][2], "Forma di Assicurazione: VALORE INTERO ")
					mapg[group][2] = append(mapg[group][2], "Formula di copertura: RISCHI NOMINATI ")
					mapg[group][2] = append(mapg[group][2], "Sono attive le garanzie opzionali:")
					isOptionShow = true
				}
			} else if group == "FIRE" && guarance.IsExtension {

				var text string
				if !isOptionShow {
					mapg[group][2] = append(mapg[group][2], "Forma di Assicurazione: VALORE INTERO ")
					mapg[group][2] = append(mapg[group][2], "Formula di copertura: RISCHI NOMINATI ")
					mapg[group][2] = append(mapg[group][2], "Sono attive le garanzie opzionali:")
					isOptionShow = true
				}
				if k.SumInsuredLimitOfIndemnity <= 1 {
					text = k.CompanyName + ": fino al  " + strconv.FormatFloat(k.SumInsuredLimitOfIndemnity*100.00, 'f', 0, 64) + "% "
				} else {
					text = k.CompanyName + ": fino a  " + humanize.FormatFloat("#.###,##", k.SumInsuredLimitOfIndemnity) + "€ "
				}
				mapg[group][2] = append(mapg[group][2], text)
			} else if group == "THEFT" {
				mapg[group][0] = append(mapg[group][0], guarance.CompanyName)
				mapg[group][1] = append(mapg[group][1], "€ "+humanize.FormatInteger("#.###,", int(k.SumInsuredLimitOfIndemnity)))
				if len(mapg[group][2]) == 1 {
					mapg[group][2] = append(mapg[group][2], "Sono attive le garanzie opzionali:")
				}
				mapg[group][2] = append(mapg[group][2], k.CompanyName+": fino a  "+humanize.FormatInteger("#.###,", int(k.SumInsuredLimitOfIndemnity))+"€ ")

			} else if group == "ELETRONIC" {
				mapg[group][0] = append(mapg[group][0], guarance.CompanyName)
				mapg[group][1] = append(mapg[group][1], "€ "+humanize.FormatInteger("#.###,", int(k.SumInsuredLimitOfIndemnity)))
				if len(mapg[group][2]) == 1 {
					mapg[group][2] = append(mapg[group][2], "Sono attive le garanzie opzionali:")
				}

				mapg[group][2] = append(mapg[group][2], k.CompanyName+": fino a  "+humanize.FormatInteger("#.###,", int(k.SumInsuredLimitOfIndemnity))+"€ ")

			} else if group == "BUSINESS INTERRUPTTION" {
				mapg[group][0] = append(mapg[group][0], guarance.CompanyName)
				mapg[group][1] = append(mapg[group][1], "€ "+humanize.FormatInteger("#.###,", int(k.SumInsuredLimitOfIndemnity)))
				mapg[group][2] = append(mapg[group][2], "La garanzia opera con una franchigia di 10 giorni ")
				mapg[group][2] = append(mapg[group][2], "ed un massimo indennizzo di 1.000 € al giorno ")

			} else if group == "ASSISTANCE" {
				mapg[group][0] = append(mapg[group][0], guarance.CompanyName)
				mapg[group][1] = append(mapg[group][1], "Inclusa")
				mapg[group][2] = append(mapg[group][2], "= = = = = = = = = = = = = = = =")

			} else {
				if !guarance.IsExtension {
					mapg[group][0] = append(mapg[group][0], guarance.CompanyName)
					mapg[group][1] = append(mapg[group][1], "€ "+humanize.FormatInteger("#.###,", int(k.SumInsuredLimitOfIndemnity)))
					mapg[group][2] = append(mapg[group][2], "= = = = = = = = = = = = = = = =")
				}

			}

			mapprice[group] = mapprice[group] + k.PriceGross
			l := len(mapg[group][2]) / 2.0
			p := lib.RoundFloat(float64(l), 0)

			mapg[group][3] = make([]string, int(p)+1)
			mapg[group][3][int(p)] = "€ " + humanize.FormatFloat("#.###,##", mapprice[group])
		}
		var listOrder []string
		var mapgOrder [][][]string
		if A.Enterprise != nil {
			listOrder = []string{"RCT", "CO", "RCP", "LEGAL", "CYBER"}
			for _, c := range listOrder {

				if v, found := mapg[c]; found {

					mapgOrder = append(mapgOrder, v)
				}

			}

			m = s.BackgroundColorRow(m, "Attività", s.SecondaryColor, s.WhiteTextCenter, s.RowTitleHeight)
			s.TableHeader(m, head1, true, 3, s.rowtableHeight+2, consts.Center, 0)
		}
		if A.Building != nil {
			log.Println("Building")
			listOrder = []string{"FIRE", "RL", "RT", "RCF", "RI", "THEFT", "ELETRONIC", "ASSISTANCE"}
			for _, c := range listOrder {

				if v, found := mapg[c]; found {
					mapgOrder = append(mapgOrder, v)
				}
			}
			m = s.BackgroundColorRow(m, "Sede", s.SecondaryColor, s.WhiteTextCenter, s.RowTitleHeight)
			s.TableHeader(m, head2, true, 3, s.rowtableHeight+2, consts.Center, 0)
		}
		log.Println("------------------------------------------------------------------------------------")

		for _, c := range mapgOrder {

			m = s.MultiRow(m, c, true, []uint{4, 2, 4, 2}, 40)
		}
	}

	return m
}
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
		skin.Space(m, 5.0)
		h := []string{"I dati della tua Polizza ", "I tuoi dati"}
		var tablePremium [][]string
		layout := "02/01/2006"
		if data.PaymentSplit == "montly" {

		}
		tablePremium = append(tablePremium, []string{"Numero: " + data.CodeCompany, "Contraente: " + data.Contractor.Name + " " + data.Contractor.Surname})
		tablePremium = append(tablePremium, []string{"Decorre dal: " + data.StartDate.Format(layout) + " ore 24:00", "C.F. / P.IVA: " + data.Contractor.VatCode})
		tablePremium = append(tablePremium, []string{"Scade il: " + data.EndDate.Format(layout) + " ore 24:00", "Indirizzo: " + data.Contractor.Address})
		tablePremium = append(tablePremium, []string{"Si rinnova a scadenza, salvo disdetta da inviare 30 giorni prima", data.Contractor.PostalCode + " " + data.Contractor.City + " (" + data.Contractor.CityCode + ")"})
		tablePremium = append(tablePremium, []string{"Prossimo pagamento il: " + data.EndDate.Format(layout), "Mail:  " + data.Contractor.Mail})
		tablePremium = append(tablePremium, []string{"Sostituisce la polizza: = = = = = = = =", "Telefono: " + data.Contractor.Phone})
		m = skin.Table(m, h, tablePremium, 6, skin.RowHeight-0.5)
		skin.Space(m, 5.0)
		//m.Line(skin.LineHeight, skin.Lin
	})

	return m
}
func (skin Skin) GetFooter(m pdf.Maroto, data models.Policy, logo string, name string) pdf.Maroto {
	m.RegisterFooter(func() {
		m.Row(15.0, func() {
			m.Col(8, func() {
				m.Text("Wopta per te. Persona è un prodotto assicurativo di Global Assistance Compagnia di assicurazioni e riassicurazioni S.p.A, distribuito da Wopta Assicurazioni S.r.l", props.Text{
					Top:         25,
					Style:       consts.Bold,
					Align:       consts.Left,
					Color:       skin.LineColor,
					Size:        skin.Size - 1,
					Extrapolate: false,
				})
			})
			m.Col(2, func() {
				_ = m.FileImage(lib.GetAssetPathByEnv("document")+"/logo_global.png", props.Rect{
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
