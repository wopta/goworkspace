package document

import (
	"log"

	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
	lib "github.com/wopta/goworkspace/lib"
)

func (s Skin) lenToHeight(w string) float64 {

	if len(w) > s.DynamicHeightMin {
		log.Println((float64(len(w)) / 30.0))
		return (float64(len(w)) / s.DynamicHeightDiv)
	} else {
		return s.rowHeight
	}

}

func (s Skin) initDefault() pdf.Maroto {
	m := pdf.NewMaroto(consts.Portrait, consts.A4)
	log.Println("Document 2")
	m.SetPageMargins(10, 15, 10)
	m.SetBackgroundColor(color.NewWhite())

	m.SetFontLocation(lib.GetAssetPathByEnv("document"))
	m.AddUTF8Font("Montserrat", consts.Normal, "montserrat_regular.ttf")
	m.AddUTF8Font("Montserrat", consts.Bold, "montserrat_bold.ttf")
	m.AddUTF8Font("Montserrat", consts.Italic, "montserrat_italic.ttf")
	m.SetDefaultFontFamily("Montserrat")
	return m
}
func getVar() (Skin, props.Line, props.Text) {
	skin := Skin{
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
		DynamicHeightDiv:  25.0,
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
	return skin, linePropMagenta, textBold
}
