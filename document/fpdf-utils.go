package document

import (
	"bytes"
	"fmt"
	"github.com/go-pdf/fpdf"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"os"
	"strings"
)

var (
	signatureID = 0
)

const (
	basePath         = "document"
	thickLineWidth   = 0.4
	thinLineWidth    = 0.1
	smallTextSize    = 6
	standardTextSize = 9
	titleTextSize    = 10
	tabDimension     = 15
	dateLayout       = "02/01/2006"
)

func initFpdf() *fpdf.Fpdf {
	pdf := fpdf.New(fpdf.OrientationPortrait, fpdf.UnitMillimeter, fpdf.PageSizeA4, "")
	pdf.SetMargins(10, 15, 10)
	loadCustomFonts(pdf)
	return pdf
}

func loadCustomFonts(pdf *fpdf.Fpdf) {
	pdf.AddUTF8Font("Montserrat", "", lib.GetAssetPathByEnv("test")+"/montserrat_light.ttf")
	pdf.AddUTF8Font("Montserrat", "B", lib.GetAssetPathByEnv("test")+"/montserrat_bold.ttf")
	pdf.AddUTF8Font("Montserrat", "I", lib.GetAssetPathByEnv("test")+"/montserrat_italic.ttf")
	pdf.AddUTF8Font("Noto", "", lib.GetAssetPathByEnv("test")+"/notosansmono.ttf")
}

func save(pdf *fpdf.Fpdf) (string, []byte) {
	filename := "document/contract.pdf"
	if os.Getenv("env") == "local" {
		err := pdf.OutputFileAndClose(filename)
		lib.CheckError(err)
	} else {
		var out bytes.Buffer
		err := pdf.Output(&out)
		lib.CheckError(err)
		filename := "temp/test_contract.pdf"
		lib.PutToStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), filename, out.Bytes())
		lib.CheckError(err)
		return filename, out.Bytes()
	}
	return filename, nil
}

func setBlackDrawColor(pdf *fpdf.Fpdf) {
	pdf.SetDrawColor(0, 0, 0)
}

func setBlackBoldFont(pdf *fpdf.Fpdf, fontSize float64) {
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Montserrat", "B", fontSize)
}

func setBlackRegularFont(pdf *fpdf.Fpdf, fontSize float64) {
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Montserrat", "", fontSize)
}

func setBlackMonospaceFont(pdf *fpdf.Fpdf, fontSize float64) {
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Noto", "", fontSize)
}

func setPinkBoldFont(pdf *fpdf.Fpdf, fontSize float64) {
	pdf.SetTextColor(229, 0, 117)
	pdf.SetFont("Montserrat", "B", fontSize)
}

func setPinkRegularFont(pdf *fpdf.Fpdf, fontSize float64) {
	pdf.SetTextColor(229, 0, 117)
	pdf.SetFont("Montserrat", "", fontSize)
}

func setPinkMonospaceFont(pdf *fpdf.Fpdf, fontSize float64) {
	pdf.SetTextColor(229, 0, 117)
	pdf.SetFont("Noto", "", fontSize)
}

func setWhiteBoldFont(pdf *fpdf.Fpdf, fontSize float64) {
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Montserrat", "B", fontSize)
}

func setWhiteRegularFont(pdf *fpdf.Fpdf, fontSize float64) {
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Montserrat", "", fontSize)
}

func drawBlackHorizontalLine(pdf *fpdf.Fpdf, width float64) {
	pdf.SetDrawColor(0, 0, 0)
	pdf.SetLineWidth(width)
	pdf.Line(11, pdf.GetY(), 200, pdf.GetY())
}

func drawPinkHorizontalLine(pdf *fpdf.Fpdf, lineWidth float64) {
	pdf.SetDrawColor(229, 0, 117)
	pdf.SetLineWidth(lineWidth)
	pdf.Line(11, pdf.GetY(), 200, pdf.GetY())
}

func drawSignatureForm(pdf *fpdf.Fpdf) {
	signatureID++
	text := fmt.Sprintf("\"[[!sigField\"%d\":signer1:signature(sigType=\\\"Click2Sign\\\"):label"+
		"(\\\"firma qui\\\"):size(width=150,height=60)]]\"", signatureID)
	pdf.SetDrawColor(0, 0, 0)
	pdf.SetX(-90)
	pdf.Cell(0, 3, "Firma del Contraente/Assicurato")
	pdf.Ln(15)
	pdf.SetLineWidth(thinLineWidth)
	pdf.Line(100, pdf.GetY(), 190, pdf.GetY())
	pdf.Ln(2)
	setBlackRegularFont(pdf, smallTextSize)
	pdf.CellFormat(0, 3, text, "", 0, fpdf.AlignRight, false, 0, "")
}

func getParagraphTitle(pdf *fpdf.Fpdf, title string) {
	pdf.SetTextColor(229, 0, 117)
	pdf.SetFont("Montserrat", "B", 10)
	pdf.Cell(0, 10, title)
}

func printSurvey(pdf *fpdf.Fpdf, survey models.Survey) {
	leftMargin, _, rightMargin, _ := pdf.GetMargins()
	pageWidth, _ := pdf.GetPageSize()
	availableWidth := pageWidth - leftMargin - rightMargin - 2
	rowWidth := pageWidth - leftMargin - rightMargin - 1

	setBlackBoldFont(pdf, standardTextSize)
	if survey.HasAnswer {
		answer := "NO"
		if *survey.Answer {
			answer = "SI"
		}

		answerWidth := pdf.GetStringWidth(answer)
		dotWidth := pdf.GetStringWidth(".")

		var surveyWidth, paddingWidth float64
		lines := pdf.SplitText(survey.Title, rowWidth)

		surveyWidth = pdf.GetStringWidth(lines[len(lines)-1])
		paddingWidth = availableWidth - surveyWidth - answerWidth

		survey.Title += strings.Repeat(".", int(paddingWidth/dotWidth)-1)
		survey.Title += answer
	}
	if survey.Title != "" {
		pdf.MultiCell(rowWidth, 3.5, survey.Title, "", fpdf.AlignLeft, false)
	}
	if survey.Subtitle != "" {
		pdf.MultiCell(rowWidth, 3.5, survey.Subtitle, "", fpdf.AlignLeft, false)
	}

	for _, question := range survey.Questions {
		availableWidth = pageWidth - leftMargin - rightMargin - 2

		if question.IsBold {
			setBlackBoldFont(pdf, standardTextSize)
		} else {
			setBlackRegularFont(pdf, standardTextSize)
		}
		if question.Indent {
			pdf.SetX(tabDimension)
			availableWidth -= tabDimension / 2
		}

		if question.HasAnswer {
			var questionWidth, paddingWidth float64

			answer := "NO"
			if *question.Answer {
				answer = "SI"
			}

			answerWidth := pdf.GetStringWidth(answer)
			dotWidth := pdf.GetStringWidth(".")

			lines := pdf.SplitText(question.Question, rowWidth)

			questionWidth = pdf.GetStringWidth(lines[len(lines)-1])
			paddingWidth = availableWidth - questionWidth - answerWidth

			question.Question += strings.Repeat(".", int(paddingWidth/dotWidth)+1)
			question.Question += answer
		}
		pdf.MultiCell(rowWidth, 3.5, question.Question, "", fpdf.AlignLeft, false)
	}
}

func printStatement(pdf *fpdf.Fpdf, statement models.Statement) {
	leftMargin, _, rightMargin, _ := pdf.GetMargins()
	pageWidth, _ := pdf.GetPageSize()
	rowWidth := pageWidth - leftMargin - rightMargin
	hasAnswer := false

	setPinkBoldFont(pdf, titleTextSize)
	if statement.Title != "" {
		pdf.MultiCell(rowWidth, 3.5, statement.Title, "", fpdf.AlignLeft, false)
		pdf.Ln(3)
	}
	setBlackBoldFont(pdf, standardTextSize)
	if statement.Subtitle != "" {
		pdf.MultiCell(rowWidth, 3.5, statement.Subtitle, "", fpdf.AlignLeft, false)
	}
	setBlackRegularFont(pdf, standardTextSize)
	for _, question := range statement.Questions {
		if question.IsBold {
			setBlackBoldFont(pdf, standardTextSize)
		} else {
			setBlackRegularFont(pdf, standardTextSize)
		}
		if question.Indent {
			pdf.SetX(tabDimension)
		}
		if question.HasAnswer {
			hasAnswer = true
		}
		pdf.MultiCell(rowWidth, 3.5, question.Question, "", fpdf.AlignLeft, false)
	}
	pdf.Ln(5)
	if hasAnswer {
		drawSignatureForm(pdf)
		pdf.Ln(10)
	}
}

func indentedText(pdf *fpdf.Fpdf, content string) {
	pdf.SetX(tabDimension)
	pdf.MultiCell(0, 3, content, "", fpdf.AlignLeft, false)
}
