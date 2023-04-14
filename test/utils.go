package test

import (
	"bytes"
	"github.com/go-pdf/fpdf"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"os"
	"strings"
)

func Save(pdf *fpdf.Fpdf) (string, error) {
	filename := "test/contract.pdf"
	if os.Getenv("env") == "local" {
		err := pdf.OutputFileAndClose(filename)
		lib.CheckError(err)
	} else {
		var buf bytes.Buffer
		err := pdf.Output(&buf)
		lib.CheckError(err)
		filename := "temp/test_contract.pdf"
		lib.PutToStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), filename, buf.Bytes())
		lib.CheckError(err)
	}
	return filename, nil
}

func SetBlackDrawColor(pdf *fpdf.Fpdf) {
	pdf.SetDrawColor(0, 0, 0)
}

func SetBlackBoldFont(pdf *fpdf.Fpdf, fontSize float64) {
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Montserrat", "B", fontSize)
}

func SetBlackRegularFont(pdf *fpdf.Fpdf, fontSize float64) {
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Montserrat", "", fontSize)
}

func SetBlackMonospaceFont(pdf *fpdf.Fpdf, fontSize float64) {
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Noto", "", fontSize)
}

func DrawBlackLine(pdf *fpdf.Fpdf, width float64) {
	pdf.SetDrawColor(0, 0, 0)
	pdf.SetLineWidth(width)
	pdf.Line(130, pdf.GetY(), 190, pdf.GetY())
}

func DrawPinkHorizontalLine(pdf *fpdf.Fpdf, lineWidth float64) {
	pdf.SetDrawColor(229, 0, 117)
	pdf.SetLineWidth(lineWidth)
	pdf.Line(11, pdf.GetY(), 200, pdf.GetY())
}

func DrawSignatureForm(pdf *fpdf.Fpdf) {
	pdf.SetX(-80)
	pdf.Cell(0, 3, "Firma del Contraente/Assicurato")
	pdf.Ln(15)
	DrawBlackLine(pdf, thickLineWidth)
}

func GetParagraphTitle(pdf *fpdf.Fpdf, title string) {
	pdf.SetTextColor(229, 0, 117)
	pdf.SetFont("Montserrat", "B", 10)
	pdf.Cell(0, 10, title)
}

func PrintStatement(pdf *fpdf.Fpdf, statement *models.Statement) {

	//yesWidth := pdf.GetStringWidth("YES")
	leftMargin, _, rightMargin, _ := pdf.GetMargins()
	pageWidth, _ := pdf.GetPageSize()
	availableWidth := pageWidth - leftMargin - rightMargin - 2
	rowWidth := pageWidth - leftMargin - rightMargin - 1

	SetBlackBoldFont(pdf, standardTextSize)
	if statement.HasAnswer {
		answer := "NO"
		if *statement.Answer {
			answer = "SI"
		}

		answerWidth := pdf.GetStringWidth(answer)
		dotWidth := pdf.GetStringWidth(".")

		var statementWidth, paddingWidth float64
		lines := pdf.SplitText(statement.Title, rowWidth)

		statementWidth = pdf.GetStringWidth(lines[len(lines)-1])
		paddingWidth = availableWidth - statementWidth - answerWidth

		statement.Title += strings.Repeat(".", int(paddingWidth/dotWidth)-1)
		statement.Title += answer
	}
	pdf.MultiCell(rowWidth, 3.5, statement.Title, "", fpdf.AlignLeft, false)

	for _, question := range statement.Questions {
		availableWidth = pageWidth - leftMargin - rightMargin - 2

		if question.IsBold {
			SetBlackBoldFont(pdf, standardTextSize)
		} else {
			SetBlackRegularFont(pdf, standardTextSize)
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

func IndentedText(pdf *fpdf.Fpdf, content string) {
	pdf.SetX(tabDimension)
	pdf.MultiCell(0, 3, content, "", fpdf.AlignLeft, false)
}
