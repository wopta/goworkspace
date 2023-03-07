package document

import (
	"log"
	"math"
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
		return s.RowHeight
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
func getVar() Skin {

	skin := Skin{
		CharForRow: 138,
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
		Size:              8,
		SizeTitle:         10,
		RowHeight:         3.4,
		RowTitleHeight:    5.0,
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
	skin.BoldtextLeft = props.Text{
		Top:   0,
		Size:  skin.Size,
		Style: consts.Bold,
		Align: consts.Left,
		Color: skin.TextColor,
	}
	skin.MagentaBoldtextLeft = props.Text{
		Top:   0,
		Size:  skin.SizeTitle,
		Style: consts.Bold,
		Align: consts.Left,
		Color: skin.LineColor,
	}
	skin.WhiteTextCenter = props.Text{
		Top:   0,
		Size:  skin.SizeTitle,
		Style: consts.Bold,
		Align: consts.Center,
		Color: color.NewWhite(),
	}
	skin.MagentaBoldtextRight = props.Text{
		Top:   0,
		Size:  skin.SizeTitle,
		Style: consts.Bold,
		Align: consts.Right,
		Color: skin.LineColor,
	}
	skin.MagentaTextLeft = props.Text{
		Top:   0,
		Size:  skin.SizeTitle,
		Style: consts.Normal,
		Align: consts.Left,
		Color: skin.LineColor,
	}
	skin.TitletextLeft = props.Text{
		Top:   0,
		Size:  skin.SizeTitle,
		Style: consts.Normal,
		Align: consts.Left,
		Color: skin.TextColor,
	}
	skin.NormaltextLeft = props.Text{
		Top:   0,
		Size:  skin.Size,
		Style: consts.Normal,
		Align: consts.Left,
		Color: skin.TextColor,
	}
	skin.NormaltextLeftBlack = props.Text{
		Top:   0,
		Size:  skin.Size,
		Style: consts.Normal,
		Align: consts.Left,
		Color: color.NewBlack(),
	}
	skin.NormaltextLeftExt = props.Text{
		Top:         0,
		Size:        skin.Size,
		Style:       consts.Normal,
		Align:       consts.Left,
		Color:       skin.TextColor,
		Extrapolate: true,
	}

	return skin
}
func (s Skin) getRowHeight(data string, base int, lineh float64) float64 {
	charsNum := len(data)
	var res float64
	if charsNum > base {
		multi := float64(charsNum) / float64(base)
		res = math.Ceil(multi) * lineh
	} else {
		res = s.RowHeight
	}
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
func (s Skin) checkPagePerc(m pdf.Maroto, perc float64) {
	current := m.GetCurrentOffset()
	_, sizeh := m.GetPageSize()

	if current > (sizeh * perc) {
		log.Println("Contrat add page")
		m.AddPage()
		s.Space(m, 10.0)

	}
}
func (s Skin) checkPageNext(m pdf.Maroto, next string) {
	len := len(next)
	log.Println(len)
	var perc float64
	perc = 0.90
	if len > 300 {
		perc = 0.90
	}
	if len > 400 {
		perc = 0.80
	}
	if len > 500 {
		perc = 0.70
	}
	if len > 600 {
		perc = 0.60
	}
	current := m.GetCurrentOffset()
	_, sizeh := m.GetPageSize()

	if current > (sizeh * perc) {
		log.Println("Contrat add page")
		m.AddPage()
		s.Space(m, 10.0)

	}
}
func (s Skin) checkIfAddPage(m pdf.Maroto, perc float64) {
	current := m.GetCurrentOffset()
	_, sizeh := m.GetPageSize()

	if current > (sizeh * perc) {
		log.Println("Contrat add page")
		m.AddPage()
		s.Space(m, 10.0)

	}
}
func ExistGuarance(list []models.Guarante, find string) bool {
	var res bool
	for _, g := range list {
		if g.Slug == find {
			return true
		}
	}

	return res
}
