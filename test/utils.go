package test

import (
	"bytes"
	"fmt"
	"github.com/go-pdf/fpdf"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"os"
	"strings"
)

const (
	thickLineWidth   = 0.4
	thinLineWidth    = 0.1
	smallTextSize    = 6
	standardTextSize = 9
	tabDimension     = 15
	layout           = "02/01/2006"
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

func save(pdf *fpdf.Fpdf) (string, error) {
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
	pdf.SetDrawColor(0, 0, 0)
	pdf.SetX(-80)
	pdf.Cell(0, 3, "Firma del Contraente/Assicurato")
	pdf.Ln(15)
	pdf.SetLineWidth(thickLineWidth)
	pdf.Line(130, pdf.GetY(), 190, pdf.GetY())
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

		var statementWidth, paddingWidth float64
		lines := pdf.SplitText(survey.Title, rowWidth)

		statementWidth = pdf.GetStringWidth(lines[len(lines)-1])
		paddingWidth = availableWidth - statementWidth - answerWidth

		survey.Title += strings.Repeat(".", int(paddingWidth/dotWidth)-1)
		survey.Title += answer
	}
	pdf.MultiCell(rowWidth, 3.5, survey.Title, "", fpdf.AlignLeft, false)

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

func indentedText(pdf *fpdf.Fpdf, content string) {
	pdf.SetX(tabDimension)
	pdf.MultiCell(0, 3, content, "", fpdf.AlignLeft, false)
}

func guaranteesToMap(guarantees []models.Guarante) map[string]models.Guarante {
	m := make(map[string]models.Guarante, 0)
	for _, guarantee := range guarantees {
		m[guarantee.Slug] = guarantee
	}
	return m
}

func extractGuarantee(guarantees []models.Guarante, guaranteeSlug string) (models.Guarante, error) {
	for _, guarantee := range guarantees {
		if guarantee.Slug == guaranteeSlug {
			return guarantee, nil
		}
	}
	return models.Guarante{}, fmt.Errorf("no %s guarantee found", guaranteeSlug)
}
