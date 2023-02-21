package document

import (
	"log"
	"strconv"

	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
	lib "github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func (s Skin) lenToHeight(w string) float64 {

	if len(w) > s.DynamicHeightMin {
		//log.Println((float64(len(w)) / 30.0))
		return (float64(len(w)) / s.DynamicHeightDiv)
	} else {
		return s.rowHeight
	}

}

func (s Skin) initDefault() pdf.Maroto {
	m := pdf.NewMaroto(consts.Portrait, consts.A4)
	log.Println("initDefault()")
	m.SetPageMargins(10, 15, 10)
	m.SetBackgroundColor(color.NewWhite())

	m.SetFontLocation(lib.GetAssetPathByEnv("document"))
	m.AddUTF8Font("Montserrat", consts.Normal, "montserrat_regular.ttf")
	m.AddUTF8Font("Montserrat", consts.Bold, "montserrat_bold.ttf")
	m.AddUTF8Font("Montserrat", consts.Italic, "montserrat_italic.ttf")
	m.SetDefaultFontFamily("Montserrat")

	return m
}
func getVar() (Skin, props.Line, props.Text, props.Text, props.Text) {

	skin := Skin{
		PrimaryColor: color.Color{
			Red:   229,
			Green: 0,
			Blue:  117,
		},
		SecondaryColor: color.Color{
			Red:   92,
			Green: 89,
			Blue:  92,
		},
		LineColor: color.Color{
			Red:   229,
			Green: 0,
			Blue:  117,
		},
		TextColor: color.Color{
			Red:   88,
			Green: 90,
			Blue:  93,
		},
		Size:              9,
		SizeTitle:         12,
		rowHeight:         7.0,
		rowtableHeight:    5.0,
		rowtableHeightMin: 2.0,
		LineHeight:        1.0,
		DynamicHeightMin:  90,
		DynamicHeightDiv:  20.0,
	}
	skin.Line = props.Line{
		Color: skin.LineColor,
		Style: consts.Solid,
		Width: 0.2,
	}
	skin.MagentaBoldtextLeft = props.Text{
		Top:   1,
		Size:  skin.SizeTitle,
		Style: consts.Bold,
		Align: consts.Left,
		Color: skin.LineColor,
	}
	skin.WhiteTextCenter = props.Text{
		Top:   1,
		Size:  skin.SizeTitle,
		Style: consts.Bold,
		Align: consts.Center,
		Color: color.NewWhite(),
	}
	skin.MagentaBoldtextRight = props.Text{
		Top:   1,
		Size:  skin.SizeTitle,
		Style: consts.Bold,
		Align: consts.Right,
		Color: skin.LineColor,
	}
	skin.MagentaTextLeft = props.Text{
		Top:   1,
		Size:  skin.SizeTitle,
		Style: consts.Normal,
		Align: consts.Left,
		Color: skin.LineColor,
	}
	skin.TitletextLeft = props.Text{
		Top:   1,
		Size:  skin.SizeTitle,
		Style: consts.Normal,
		Align: consts.Left,
		Color: skin.TextColor,
	}
	skin.NormaltextLeft = props.Text{
		Top:   1,
		Size:  skin.Size,
		Style: consts.Normal,
		Align: consts.Left,
		Color: skin.TextColor,
	}
	skin.NormaltextLeftExt = props.Text{
		Top:         1,
		Size:        skin.Size,
		Style:       consts.Normal,
		Align:       consts.Left,
		Color:       skin.TextColor,
		Extrapolate: true,
	}
	magenta := props.Text{
		Top:   1,
		Size:  skin.SizeTitle,
		Style: consts.Bold,
		Align: consts.Left,
		Color: skin.LineColor,
	}
	linePropMagenta := props.Line{
		Color: skin.LineColor,
		Style: consts.Solid,
		Width: 0.2,
	}

	textBold := props.Text{
		Top:   1,
		Style: consts.Bold,
		Align: consts.Center,
	}
	normal := props.Text{
		Top:   1,
		Size:  skin.Size,
		Style: consts.Bold,
		Align: consts.Left,
		Color: skin.TextColor,
	}
	return skin, linePropMagenta, textBold, normal, magenta
}
func getRowHeight(data string, base int, lineh int) int {
	charsNum := len(data)
	multi := charsNum / base
	res := multi * lineh
	return res
}
func sumStringFloat(data []string, price float64) float64 {

	var sum float64
	for _, v := range data {
		s, _ := strconv.ParseFloat(v, 64)

		sum = sum + s
	}
	res := sum + price

	return res
}
func (s Skin) checkPage(m pdf.Maroto) {
	current := m.GetCurrentOffset()
	_, sizeh := m.GetPageSize()

	if current > (sizeh * 0.61) {
		log.Println("Contrat add page")
		m.AddPage()
		s.Space(m, 10.0)

	}
}
func ExistGuarance(list []models.Guarantee, find string) bool {
	var res bool
	for _, g := range list {
		if g.Slug == find {
			return true
		}
	}

	return res
}
