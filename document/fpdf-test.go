package document

import (
	"encoding/json"
	"github.com/go-pdf/fpdf"
	"github.com/wopta/goworkspace/models"
	"log"
	"net/http"
	"os"
)

type Statements struct {
	Statements []*models.Statement `json:"statements"`
	Text       string              `json:"text,omitempty"`
}

func FpdfHandler(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {

	pdf := initFpdf()

	GetHeader(pdf, "life")
	GetFooter(pdf)

	pdf.AddPage()

	GetParagraphTitle(pdf, "La tua assicurazione è operante per il seguente Assicurato e Garanzie")
	pdf.Ln(8)
	DrawPinkLine(pdf, 0.4)
	pdf.Ln(2)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.Cell(20, 2, "Cognome e Nome")
	pdf.SetFont("Montserrat", "", 9)
	pdf.SetX(pdf.GetX() + 24)
	pdf.Cell(20, 2, "HAMMAR YOUSEF")
	pdf.SetX(pdf.GetX() + 60)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.Cell(10, 2, "Codice fiscale:")
	pdf.SetX(pdf.GetX() + 20)
	pdf.SetFont("Montserrat", "", 9)
	pdf.Cell(20, 2, "HMMYSF94R07D912M")
	pdf.Ln(2.5)
	DrawPinkLine(pdf, 0.1)
	pdf.Ln(2)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.Cell(20, 2, "Residente in")
	pdf.SetFont("Montserrat", "", 9)
	pdf.SetX(pdf.GetX() + 24)
	pdf.Cell(20, 2, "Via Unicef 4 - 20033 Solaro (MI)")
	pdf.Ln(2.5)
	DrawPinkLine(pdf, 0.1)
	pdf.Ln(2)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.Cell(20, 2, "Mail")
	pdf.SetFont("Montserrat", "", 9)
	pdf.SetX(pdf.GetX() + 24)
	pdf.Cell(20, 2, "yousef.hammar@wopta.it")
	pdf.SetX(pdf.GetX() + 60)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.Cell(10, 2, "Telefono")
	pdf.SetX(pdf.GetX() + 20)
	pdf.SetFont("Montserrat", "", 9)
	pdf.Cell(20, 2, "+393451021004")
	pdf.Ln(2.5)
	DrawPinkLine(pdf, 0.4)
	pdf.Ln(1)

	pdf.SetFont("Montserrat", "B", 9)
	pdf.CellFormat(80, 6, "Garanzie", "", 0, "CM", false, 0, "")
	pdf.CellFormat(25, 6, "Somma\nassicurata €", "", 0, "CM", false, 0, "")
	pdf.CellFormat(30, 6, "Durata anni", "", 0, "CM", false, 0, "")
	pdf.CellFormat(25, 6, "Scade il", "", 0, "CM", false, 0, "")
	pdf.CellFormat(25, 6, "Premio annuale €", "", 0, "CM", false, 0, "")
	pdf.Ln(5)
	DrawPinkLine(pdf, 0.2)
	pdf.CellFormat(80, 6, "Decesso", "", 0, "", false, 0, "")
	pdf.SetFont("Montserrat", "", 9)
	pdf.CellFormat(25, 6, "100.000 €", "", 0, "RM", false, 0, "")
	pdf.CellFormat(30, 6, "20", "", 0, "CM", false, 0, "")
	pdf.CellFormat(25, 6, "04/03/2023", "", 0, "CM", false, 0, "")
	pdf.CellFormat(25, 6, "97,74 €     ", "", 0, "RM", false, 0, "")
	pdf.Ln(5)
	DrawPinkLine(pdf, 0.2)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.CellFormat(80, 6, "Invalidità Totale Permanente da Infortunio o Malattia", "", 0, "", false, 0, "")
	pdf.SetFont("Montserrat", "", 9)
	pdf.CellFormat(25, 6, "100.000 €", "", 0, "RM", false, 0, "")
	pdf.CellFormat(30, 6, "20", "", 0, "CM", false, 0, "")
	pdf.CellFormat(25, 6, "04/03/2023", "", 0, "CM", false, 0, "")
	pdf.CellFormat(25, 6, "12,02 € (*)", "", 0, "RM", false, 0, "")
	pdf.Ln(5)
	DrawPinkLine(pdf, 0.2)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.CellFormat(80, 6, "Inabilità Temporanea da Infortunio o Malattia", "", 0, "", false, 0, "")
	pdf.SetFont("Montserrat", "", 9)
	pdf.CellFormat(25, 6, "500 €", "", 0, "RM", false, 0, "")
	pdf.CellFormat(30, 6, "10", "", 0, "CM", false, 0, "")
	pdf.CellFormat(25, 6, "04/03/2023", "", 0, "CM", false, 0, "")
	pdf.CellFormat(25, 6, "10,00 € (*)", "", 0, "RM", false, 0, "")
	pdf.Ln(5)
	DrawPinkLine(pdf, 0.2)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.CellFormat(80, 6, "Malattie Gravi", "", 0, "", false, 0, "")
	pdf.SetFont("Montserrat", "", 9)
	pdf.CellFormat(25, 6, "10.000 €", "", 0, "RM", false, 0, "")
	pdf.CellFormat(30, 6, "10", "", 0, "CM", false, 0, "")
	pdf.CellFormat(25, 6, "04/03/2023", "", 0, "CM", false, 0, "")
	pdf.CellFormat(25, 6, "55,42 € (*)", "", 0, "RM", false, 0, "")
	pdf.Ln(5)
	DrawPinkLine(pdf, 0.2)
	pdf.Ln(0.5)
	pdf.SetFont("Montserrat", "", 6)
	pdf.Cell(80, 3, "(*) imposte assicurative di legge incluse nella misura del 2,50% del premio imponibile")
	pdf.Ln(3)

	GetParagraphTitle(pdf, "Nomina dei Beneficiari e Referente terzo, per il caso di garanzia Decesso (qualora sottoscritta)")
	pdf.Ln(8)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Montserrat", "", 9)
	pdf.MultiCell(0, 3, "AVVERTENZE: Può scegliere se designare nominativamente i beneficiari o se designare genericamente come beneficiari i suoi eredi legittimi e/o testamentari. In caso di mancata designazione nominativa, la Compagnia potrà incontrare, al decesso dell’Assicurato, maggiori difficoltà nell’identificazione e nella ricerca dei beneficiari. La modifica o revoca del/i beneficiario/i deve essere comunicata alla Compagnia in forma scritta.\nIn caso di specifiche esigenze di riservatezza, la Compagnia potrà rivolgersi ad un soggetto terzo (diverso dal Beneficiario)\nIn caso di Decesso al fine di contattare il Beneficiario designato.", "", "", false)
	pdf.Ln(0)

	GetParagraphTitle(pdf, "Beneficiario")
	pdf.Ln(8)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Montserrat", "", 9)
	pdf.CellFormat(0, 3, "Io sottoscritto Assicurato, con la sottoscrizione della presente polizza, in riferimento alla garanzia Decesso:", "", 0, "", false, 0, "")
	pdf.Ln(4)
	pdf.SetDrawColor(0, 0, 0)
	pdf.SetX(11)
	pdf.CellFormat(3, 3, "X", "1", 0, "CM", false, 0, "")
	pdf.CellFormat(0, 3, "Designo genericamente quali beneficiari della prestazione i miei eredi (legittimi e/o testamentari)", "", 0, "", false, 0, "")
	pdf.Ln(4)
	pdf.SetX(11)
	pdf.CellFormat(3, 3, "", "1", 0, "CM", false, 0, "")
	pdf.CellFormat(0, 3, "Designo nominativamente il/i seguente/i soggetto/i quale beneficiario/i della prestazione", "", 0, "", false, 0, "")
	pdf.Ln(5)
	DrawPinkLine(pdf, 0.4)
	pdf.Ln(1)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.Cell(50, 2, "Cognome e nome")
	pdf.SetFont("Montserrat", "", 9)
	pdf.Cell(20, 2, "=====")
	pdf.SetX(pdf.GetX() + 60)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.Cell(20, 2, "Cod. Fisc.: ")
	pdf.SetFont("Montserrat", "", 9)
	pdf.Cell(20, 2, "=====")
	pdf.Ln(3)
	DrawPinkLine(pdf, 0.1)
	pdf.Ln(2)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.Cell(50, 2, "Indirizzo")
	pdf.SetFont("Montserrat", "", 9)
	pdf.Cell(20, 2, "=====")
	pdf.SetX(pdf.GetX() + 60)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.Ln(3)
	DrawPinkLine(pdf, 0.1)
	pdf.Ln(2)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.Cell(50, 2, "Mail")
	pdf.SetFont("Montserrat", "", 9)
	pdf.Cell(20, 2, "=====")
	pdf.SetX(pdf.GetX() + 60)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.Cell(20, 2, "Telefono: ")
	pdf.SetFont("Montserrat", "", 9)
	pdf.Cell(20, 2, "=====")
	pdf.Ln(3)
	DrawPinkLine(pdf, 0.1)
	pdf.Ln(2)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.Cell(50, 2, "Relazione con Assicurato")
	pdf.SetFont("Montserrat", "", 9)
	pdf.Cell(20, 2, "=====")
	pdf.Ln(3)
	DrawPinkLine(pdf, 0.1)
	pdf.Ln(2)
	pdf.Cell(165, 2, "Consenso ad invio comunicazioni da parte della Compagnia al beneficiario, prima dell'evento Decesso:")
	pdf.Cell(20, 2, "=====")
	pdf.Ln(3)
	DrawPinkLine(pdf, 0.1)
	pdf.Ln(2)
	DrawPinkLine(pdf, 0.4)
	pdf.Ln(1)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.Cell(50, 2, "Cognome e nome")
	pdf.SetFont("Montserrat", "", 9)
	pdf.Cell(20, 2, "=====")
	pdf.SetX(pdf.GetX() + 60)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.Cell(20, 2, "Cod. Fisc.: ")
	pdf.SetFont("Montserrat", "", 9)
	pdf.Cell(20, 2, "=====")
	pdf.Ln(3)
	DrawPinkLine(pdf, 0.1)
	pdf.Ln(2)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.Cell(50, 2, "Indirizzo")
	pdf.SetFont("Montserrat", "", 9)
	pdf.Cell(20, 2, "=====")
	pdf.SetX(pdf.GetX() + 60)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.Ln(3)
	DrawPinkLine(pdf, 0.1)
	pdf.Ln(2)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.Cell(50, 2, "Mail")
	pdf.SetFont("Montserrat", "", 9)
	pdf.Cell(20, 2, "=====")
	pdf.SetX(pdf.GetX() + 60)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.Cell(20, 2, "Telefono: ")
	pdf.SetFont("Montserrat", "", 9)
	pdf.Cell(20, 2, "=====")
	pdf.Ln(3)
	DrawPinkLine(pdf, 0.1)
	pdf.Ln(2)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.Cell(50, 2, "Relazione con Assicurato")
	pdf.SetFont("Montserrat", "", 9)
	pdf.Cell(20, 2, "=====")
	pdf.Ln(3)
	DrawPinkLine(pdf, 0.1)
	pdf.Ln(2)
	pdf.Cell(165, 2, "Consenso ad invio comunicazioni da parte della Compagnia al beneficiario, prima dell'evento Decesso:")
	pdf.Cell(20, 2, "=====")
	pdf.Ln(3)
	DrawPinkLine(pdf, 0.1)
	pdf.Ln(2)

	GetParagraphTitle(pdf, "Referente terzo")
	pdf.Ln(8)
	pdf.SetDrawColor(229, 0, 117)
	DrawPinkLine(pdf, 0.4)
	pdf.Ln(1)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.Cell(50, 2, "Cognome e nome")
	pdf.SetFont("Montserrat", "", 9)
	pdf.Cell(20, 2, "=====")
	pdf.SetX(pdf.GetX() + 60)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.Cell(20, 2, "Cod. Fisc.: ")
	pdf.SetFont("Montserrat", "", 9)
	pdf.Cell(20, 2, "=====")
	pdf.Ln(3)
	DrawPinkLine(pdf, 0.1)
	pdf.Ln(2)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.Cell(50, 2, "Indirizzo")
	pdf.SetFont("Montserrat", "", 9)
	pdf.Cell(20, 2, "=====")
	pdf.SetX(pdf.GetX() + 60)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.Ln(3)
	DrawPinkLine(pdf, 0.1)
	pdf.Ln(2)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.Cell(50, 2, "Mail")
	pdf.SetFont("Montserrat", "", 9)
	pdf.Cell(20, 2, "=====")
	pdf.SetX(pdf.GetX() + 60)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.Cell(20, 2, "Telefono: ")
	pdf.SetFont("Montserrat", "", 9)
	pdf.Cell(20, 2, "=====")
	pdf.Ln(3)
	DrawPinkLine(pdf, 0.1)
	pdf.Ln(2)

	var statements Statements
	b, err := os.ReadFile("document/response.json")
	if err != nil {
		log.Println(err)
	}
	err = json.Unmarshal(b, &statements)
	if err != nil {
		log.Println(err)
	}

	GetParagraphTitle(pdf, "Dichiarazioni da leggere con attenzione prima di firmare")
	pdf.Ln(8)
	PrintStatement(pdf, statements.Statements[0])
	pdf.Ln(5)
	GetParagraphTitle(pdf, "Questionario Medico")
	pdf.Ln(8)
	for _, statement := range statements.Statements[1:] {
		PrintStatement(pdf, statement)
	}
	pdf.Ln(8)
	pdf.SetX(-80)
	pdf.Cell(0, 3, "Firma del Contraente/Assicurato")
	pdf.Ln(15)
	pdf.SetDrawColor(0, 0, 0)
	pdf.SetLineWidth(0.4)
	pdf.Line(130, pdf.GetY(), 190, pdf.GetY())
	//GetFooter(pdf)

	//pdf.AddPage()

	err = pdf.OutputFileAndClose("document/test.pdf")
	log.Println(err)
	return "", nil, err
}

func initFpdf() *fpdf.Fpdf {
	pdf := fpdf.New(fpdf.OrientationPortrait, fpdf.UnitMillimeter, fpdf.PageSizeA4, "")
	pdf.SetMargins(10, 15, 10)
	loadCustomFonts(pdf)
	return pdf
}

func loadCustomFonts(pdf *fpdf.Fpdf) {
	pdf.AddUTF8Font("Montserrat", "", "document/assets/montserrat_regular.ttf")
	pdf.AddUTF8Font("Montserrat", "B", "document/assets/montserrat_bold.ttf")
	pdf.AddUTF8Font("Montserrat", "I", "document/assets/montserrat_italic.ttf")
}

func DrawPinkLine(pdf *fpdf.Fpdf, lineWidth float64) {
	pdf.SetDrawColor(229, 0, 117)
	pdf.SetLineWidth(lineWidth)
	pdf.Line(11, pdf.GetY(), 200, pdf.GetY())
}

func GetHeader(pdf *fpdf.Fpdf, name string) {
	var opt fpdf.ImageOptions
	var product, logoPath string
	pathPrefix := "document/assets/logo_"

	switch name {
	case "life":
		product = "Vita"
		logoPath = pathPrefix + "vita.png"
	}

	pdf.SetHeaderFunc(func() {
		opt.ImageType = "png"
		pdf.ImageOptions(logoPath, 10, 6, 13, 13, false, opt, 0, "")
		pdf.SetXY(23, 7)
		pdf.SetTextColor(229, 0, 117)
		pdf.SetFont("Montserrat", "B", 18)
		pdf.Cell(10, 6, "Wopta per te")
		pdf.SetFont("Montserrat", "I", 18)
		pdf.SetXY(23, 13)
		pdf.SetTextColor(92, 89, 92)
		pdf.Cell(10, 6, product)
		pdf.ImageOptions("document/assets/ARTW_LOGO_RGB_400px.png", 158, 6, 0, 10, false, opt, 0, "")

		pdf.SetTextColor(0, 0, 0)
		pdf.SetFont("Montserrat", "B", 8)
		pdf.SetXY(11, 20)
		pdf.Cell(0, 3, "I dati della tua polizza")
		pdf.SetFont("Montserrat", "", 8)
		pdf.SetXY(11, pdf.GetY()+3)
		pdf.MultiCell(0, 3, "Numero: 12345\nDecorre dal: 03/04/2023 ore 24:00\nScade il: 03/04/2043 ore 24:00\nPrima scadenza annuale il: 04/04/2024\nNon si rinnova a scadenza.", "", "", false)

		pdf.SetFont("Montserrat", "B", 8)
		pdf.SetXY(-90, 20)
		pdf.Cell(0, 3, "I tuoi dati")
		pdf.SetFont("Montserrat", "", 8)
		pdf.SetXY(-90, pdf.GetY()+3)
		pdf.MultiCell(0, 3, "Contraente: HAMMAR YOUSEF\nC.F./P.IVA: HMMYSF94R07D912M\nIndirizzo: Via Unicef, 4\n20033 SOLARO (MI)\nMail: yousef.hammar@wopta.it\nTelefono: +393451031004", "", "", false)
		pdf.Ln(10)
	})
}

func GetFooter(pdf *fpdf.Fpdf) {
	var opt fpdf.ImageOptions

	pdf.SetFooterFunc(func() {
		pdf.SetXY(10, -13)
		pdf.SetTextColor(229, 0, 117)
		pdf.SetFont("Montserrat", "B", 6)
		pdf.MultiCell(0, 3, "Wopta per te. Vita è un prodotto assicurativo di AXA France Vie S.A. – Rappresentanza Generale per l’Italia\ndistribuito da Wopta Assicurazioni S.r.l.", "", "", false)
		opt.ImageType = "png"
		pdf.ImageOptions("document/assets/logo_axa.png", 190, 283, 8, 8, false, opt, 0, "")
	})
}

func GetParagraphTitle(pdf *fpdf.Fpdf, title string) {
	pdf.SetTextColor(229, 0, 117)
	pdf.SetFont("Montserrat", "B", 10)
	pdf.Cell(0, 10, title)
}

func PrintStatement(pdf *fpdf.Fpdf, statement *models.Statement) {
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.MultiCell(0, 3, statement.Title, "", "", false)
	for _, question := range statement.Questions {
		if question.IsBold {
			pdf.SetFont("Montserrat", "B", 9)
		}
		if question.Indent {
			pdf.SetX(15)
		}
		pdf.MultiCell(0, 3, question.Question, "", "", false)
	}
}
