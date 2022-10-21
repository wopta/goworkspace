package document

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
	lib "github.com/wopta/goworkspace/lib"
	//model "github.com/wopta/goworkspace/models"
)

func init() {
	log.Println("INIT Document")
	functions.HTTP("Document", Document)
}

func Document(w http.ResponseWriter, r *http.Request) {
	log.Println("Document")
	lib.EnableCors(&w, r)
	req := lib.ErrorByte(ioutil.ReadAll(r.Body))
	log.Println(string(req))
	var data PdfData
	// Unmarshal or Decode the JSON to the interface.
	//json.NewDecoder(req).Decode(&send)
	defer r.Body.Close()

	json.Unmarshal([]byte(req), &data)
	//data := PdfData{name: ""}

	//begin := time.Now()

	magenta := color.Color{
		Red:   229,
		Green: 0,
		Blue:  117,
	}
	gray := color.Color{
		Red:   88,
		Green: 90,
		Blue:  93,
	}
	log.Println(gray)
	textBoldRight := props.Text{
		Top:   1.5,
		Size:  9,
		Style: consts.Bold,
		Align: consts.Center,
	}
	linePropMagenta := props.Line{
		Color: magenta,
		Style: consts.Solid,
		Width: 0.2,
	}

	//darkGrayColor := color.NewBlack()
	rowHeight := 5.0
	rowHeightSlim := 1.0
	//blackColor := color.NewBlack()
	whiteColor := color.NewWhite()
	textBold := props.Text{
		Top:   3,
		Style: consts.Bold,
		Align: consts.Center,
	}

	textBoldMagenta := props.Text{
		Color: magenta,
		Top:   3,
		Style: consts.Bold,
		Align: consts.Center,
	}
	textMagenta := props.Text{
		Color: magenta,
		Top:   0,
		Style: consts.Normal,
		Align: consts.Left,
		Size:  1,
	}

	log.Println(textBold)
	log.Println("Document 1")
	m := pdf.NewMaroto(consts.Portrait, consts.A4)
	log.Println("Document 2")
	m.SetPageMargins(10, 15, 10)
	m.SetBackgroundColor(whiteColor)
	m.SetFontLocation("document")

	// Define font to all styles.
	m.AddUTF8Font("Montserrat", consts.Normal, "Montserrat-Regular.ttf")
	m.AddUTF8Font("Montserrat", consts.Bold, "Montserrat-Bold.ttf")
	m.SetDefaultFontFamily("Montserrat")
	//m.SetBorder(true)

	m.RegisterHeader(func() {
		m.Row(rowHeight, func() {
			m.Col(12, func() {
				_ = m.FileImage("document/ARTW_LOGO_RGB_400px.png", props.Rect{
					Left:    5,
					Top:     5,
					Center:  true,
					Percent: 85,
				})
			})
		})
	})

	m.RegisterFooter(func() {
	})
	log.Println("Document 3")

	m.Row(10, func() {
		m.Col(12, func() {
			m.Text("La tua assicurazione è operante per il seguente Assicurato e Garanzie ", textBoldMagenta)
		})
		//m.SetBackgroundColor(magenta)
	})
	m.Line(1.0, linePropMagenta)
	//m.SetBackgroundColor(darkGrayColor)
	log.Println("Document 4")
	m.Row(rowHeight, func() {
		m.Col(2, func() {
			m.Text("Assicurato:", textBoldRight)

		})

		m.Col(2, func() {
			m.Text("xxxx", props.Text{
				Top:   1.5,
				Size:  9,
				Style: consts.Normal,
				Align: consts.Center,
			})
		})
	})

	m.Row(rowHeightSlim, func() {
		m.Col(2, func() {
			m.Text("_________________________________________________", textMagenta)

		})
	})

	log.Println("Document 8")
	//m.Output()
	err := m.OutputFileAndClose("document/billing.pdf")
	lib.CheckError(err)

	//end := time.Now()
	//fmt.Println(end.Sub(begin))

}

type PdfData struct {
	Name         string  `json:"name"`
	Surname      string  `json:"surname"`
	FiscalCode   string  `json:"fiscalCode"`
	Work         string  `json:"work"`
	WorkType     string  `json:"workType"`
	Class        string  `json:"class"`
	CoverageType string  `json:"coverageType"`
	Price        float64 `json:"price"`
	PriceNett    float64 `json:"priceNett"`
	Coverages    []struct {
		Deductible                 string  `json:"deductible"`
		SelfInsurance              string  `json:"selfInsurance"`
		SumInsuredLimitOfIndemnity float64 `json:"sumInsuredLimitOfIndemnity"`
		Price                      float64 `json:"price"`
		PriceNett                  float64 `json:"priceNett"`
	} `json:"coverages"`
	Statements []struct {
		Text string `json:"text"`
	} `json:"statements"`
	SpecialConditions []struct {
		Text string `json:"text"`
	} `json:"specialConditions"`
}
