package document

import (
	"log"
	"net/http"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
	lib "github.com/wopta/goworkspace/lib"
	model "github.com/wopta/goworkspace/models"
)

func init() {
	log.Println("INIT Document")
	functions.HTTP("Document", Document)
}

func Document(w http.ResponseWriter, r *http.Request) {
	log.Println("Document")
	lib.EnableCors(&w, r)

	//data := PdfData{name: ""}

	//begin := time.Now()
	magenta := color.Color{
		// Red is the amount of red
		Red: 0,
		// Green is the amount of red
		Green: 0,
		// Blue is the amount of red
		Blue: 0,
	}
	textBoldRight := props.Text{
		Top:   1.5,
		Size:  9,
		Style: consts.Bold,
		Align: consts.Center,
	}
	lineProp := props.Line{
		Color: magenta,
		Style: consts.Solid,
		Width: 1.0,
	}

	//darkGrayColor := color.NewBlack()
	rowHeight := 5.0
	grayColor := color.NewBlack()
	whiteColor := color.NewWhite()
	textBold := props.Text{
		Top:   3,
		Style: consts.Bold,
		Align: consts.Center,
	}
	header := []string{""}
	contents := [][]string{{""}}
	log.Println("Document 1")
	m := pdf.NewMaroto(consts.Portrait, consts.A4)
	log.Println("Document 2")
	m.SetPageMargins(10, 15, 10)
	//m.SetBorder(true)

	m.RegisterHeader(func() {
	})

	m.RegisterFooter(func() {
	})
	log.Println("Document 3")
	m.Line(1.0, lineProp)
	m.Row(10, func() {
		m.Col(12, func() {
			m.Text("La tua assicurazione è operante per il seguente Assicurato e Garanzie ", textBold)
		})
		//m.SetBackgroundColor(magenta)
	})
	m.Line(1.0, props.Line{
		Color: color.Color{
			Red:   255,
			Green: 100,
			Blue:  50,
		},
		Style: consts.Dotted,
		Width: 1.0,
	})
	//m.SetBackgroundColor(darkGrayColor)
	log.Println("Document 4")
	m.Row(rowHeight, func() {
		m.Col(2, func() {
			m.Text("Assicurato", textBoldRight)
			m.Line(1.0, lineProp)
		})
		m.Col(2, func() {
			m.Text("xxxx", props.Text{
				Top:   1.5,
				Size:  9,
				Style: consts.Bold,
				Align: consts.Center,
			})
		})
		m.ColSpace(9)
	})

	m.SetBackgroundColor(whiteColor)

	m.TableList(header, contents, props.TableList{
		HeaderProp: props.TableListContent{
			Size:      9,
			GridSizes: []uint{3, 4, 2, 3},
		},
		ContentProp: props.TableListContent{
			Size:      8,
			GridSizes: []uint{3, 4, 2, 3},
		},
		Align:                consts.Center,
		AlternatedBackground: &grayColor,
		HeaderContentSpace:   1,
		Line:                 false,
	})

	m.Row(20, func() {
		m.ColSpace(7)
		m.Col(2, func() {
			m.Text("Total:", props.Text{
				Top:   5,
				Style: consts.Bold,
				Size:  8,
				Align: consts.Right,
			})
		})
		m.Col(3, func() {
			m.Text("R$ 2.567,00", props.Text{
				Top:   5,
				Style: consts.Bold,
				Size:  8,
				Align: consts.Center,
			})
		})
	})

	m.Row(15, func() {
		m.Col(6, func() {
			_ = m.Barcode("5123.151231.512314.1251251.123215", props.Barcode{
				Percent: 0,
				Proportion: props.Proportion{
					Width:  20,
					Height: 2,
				},
			})
			m.Text("5123.151231.512314.1251251.123215", props.Text{
				Top:    12,
				Family: "",
				Style:  consts.Bold,
				Size:   9,
				Align:  consts.Center,
			})
		})
		m.ColSpace(6)
	})
	log.Println("Document 8")
	//m.Output()
	err := m.OutputFileAndClose("document/billing.pdf")
	if err != nil {

		os.Exit(1)
	}

	//end := time.Now()
	//fmt.Println(end.Sub(begin))

}

type PdfData struct {
	Tite  string
	Users model.User
	name  string
}
