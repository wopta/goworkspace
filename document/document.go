package document

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"

	lib "github.com/wopta/goworkspace/lib"
)

func init() {
	log.Println("INIT Document")

	functions.HTTP("Document", Document)
}

func Document(w http.ResponseWriter, r *http.Request) {
	var (
		templ *template.Template
		err   error
	)
	data := PdfData{name: ""}
	// use Go's default HTML template generation tools to generate your HTML
	log.Println("use Go's default HTML template generation tools to generate your HTML")
	templ, err = template.ParseFiles("document/template.html")
	lib.CheckError(err)
	// apply the parsed HTML template data and keep the result in a Buffer
	log.Println(" apply the parsed HTML template data and keep the result in a Buffer")
	var body bytes.Buffer
	err = templ.Execute(&body, data)
	lib.CheckError(err)
	// Create object from reader.
	inFile, err := os.Open("sample2.html")
	lib.CheckError(err)
	defer inFile.Close()
	//begin := time.Now()
	magenta := color.Color{
		// Red is the amount of red
		Red: 0,
		// Green is the amount of red
		Green: 0,
		// Blue is the amount of red
		Blue: 0,
	}
	darkGrayColor := color.NewBlack()
	grayColor := color.NewBlack()
	whiteColor := color.NewWhite()
	header := []string{""}
	contents := [][]string{{""}}

	m := pdf.NewMaroto(consts.Portrait, consts.A4)
	m.SetPageMargins(10, 15, 10)
	//m.SetBorder(true)

	m.RegisterHeader(func() {
	})

	m.RegisterFooter(func() {
	})

	m.Row(10, func() {
		m.Col(12, func() {
			m.Text("Invoice ABC123456789", props.Text{
				Top:   3,
				Style: consts.Bold,
				Align: consts.Center,
			})
		})
		m.SetBackgroundColor(magenta)
	})

	m.SetBackgroundColor(darkGrayColor)

	m.Row(7, func() {
		m.Col(3, func() {
			m.Text("Transactions", props.Text{
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

	err = m.OutputFileAndClose("internal/examples/pdfs/billing.pdf")
	if err != nil {

		os.Exit(1)
	}

	//end := time.Now()
	//fmt.Println(end.Sub(begin))

}

type PdfData struct {
	name string
}
