package document

import (
	"bytes"
	"fmt"
	"github.com/go-pdf/fpdf"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"os"
	"strconv"
	"strings"
	"time"
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
	pdf.SetAuthor("Wopta Assicurazioni s.r.l", true)
	pdf.SetCreationDate(time.Now().UTC())
	loadCustomFonts(pdf)
	return pdf
}

func loadCustomFonts(pdf *fpdf.Fpdf) {
	pdf.AddUTF8Font("Montserrat", "", lib.GetAssetPathByEnv(basePath)+"/montserrat_light.ttf")
	pdf.AddUTF8Font("Montserrat", "B", lib.GetAssetPathByEnv(basePath)+"/montserrat_bold.ttf")
	pdf.AddUTF8Font("Montserrat", "I", lib.GetAssetPathByEnv(basePath)+"/montserrat_italic.ttf")
}

func save(pdf *fpdf.Fpdf, policy *models.Policy) (string, []byte) {
	var filename string
	if os.Getenv("env") == "local" {
		err := pdf.OutputFileAndClose(basePath + "/contract.pdf")
		lib.CheckError(err)
	} else {
		var out bytes.Buffer
		err := pdf.Output(&out)
		lib.CheckError(err)
		now := time.Now()
		timestamp := strconv.FormatInt(now.Unix(), 10)
		filename = "temp/" + policy.Uid + "/" + policy.Contractor.Name + "_" + policy.Contractor.Surname + "_" + timestamp + "_contract.pdf"
		lib.PutToStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), filename, out.Bytes())
		lib.CheckError(err)
		return filename, out.Bytes()
	}
	return filename, nil
}

func pageNumber(pdf *fpdf.Fpdf) {
	setBlackRegularFont(pdf, smallTextSize)
	pdf.MultiCell(0, 3, fmt.Sprintf("pagina %d", pdf.PageNo()), "", fpdf.AlignRight, false)
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

func setBlackItalicFont(pdf *fpdf.Fpdf, fontSize float64) {
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Montserrat", "I", fontSize)
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

func setPinkItalicFont(pdf *fpdf.Fpdf, fontSize float64) {
	pdf.SetTextColor(229, 0, 117)
	pdf.SetFont("Montserrat", "I", fontSize)
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

func setWhiteItalicFont(pdf *fpdf.Fpdf, fontSize float64) {
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Montserrat", "I", fontSize)
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
	pdf.SetFont(consts.Arial, "", 6)
	pdf.CellFormat(0, 3, text, "", 0, fpdf.AlignRight, false, 0, "")
}

func getParagraphTitle(pdf *fpdf.Fpdf, title string) {
	pdf.SetTextColor(229, 0, 117)
	pdf.SetFont("Montserrat", "B", titleTextSize)
	pdf.MultiCell(0, 4, title, "", "", false)
}

func printSurvey(pdf *fpdf.Fpdf, survey models.Survey) error {
	var text string
	leftMargin, _, rightMargin, _ := pdf.GetMargins()
	pageWidth, _ := pdf.GetPageSize()
	availableWidth := pageWidth - leftMargin - rightMargin - 2

	setBlackBoldFont(pdf, standardTextSize)
	if survey.HasAnswer {
		if *survey.Answer != *survey.ExpectedAnswer {
			return fmt.Errorf("%s: answer not equal expected answer", survey.Title)
		}
		answer := "NO"
		if *survey.Answer {
			answer = "SI"
		}

		answerWidth := pdf.GetStringWidth(answer)
		dotWidth := pdf.GetStringWidth(".")

		var surveyWidth, paddingWidth float64
		lines := pdf.SplitText(survey.Title+answer, availableWidth)

		surveyWidth = pdf.GetStringWidth(lines[len(lines)-1])
		paddingWidth = availableWidth - surveyWidth - answerWidth

		text = strings.Repeat(".", int(paddingWidth/dotWidth)-2) + answer
	}
	if survey.Title != "" {
		pdf.MultiCell(availableWidth, 3.5, survey.Title+text, "", fpdf.AlignLeft, false)
	}
	if survey.Subtitle != "" {
		pdf.MultiCell(availableWidth, 3.5, survey.Subtitle+text, "", fpdf.AlignLeft, false)
	}

	for _, question := range survey.Questions {
		text = ""
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
			if *question.Answer != *question.ExpectedAnswer {
				return fmt.Errorf("%s: answer not equal expected answer", question.Question)
			}

			answer := "NO"
			if *question.Answer {
				answer = "SI"
			}

			answerWidth := pdf.GetStringWidth(answer)
			dotWidth := pdf.GetStringWidth(".")

			lines := pdf.SplitText(question.Question+answer, availableWidth)

			questionWidth = pdf.GetStringWidth(lines[len(lines)-1])
			paddingWidth = availableWidth - questionWidth - answerWidth

			text = strings.Repeat(".", int(paddingWidth/dotWidth)-2) + answer
		}
		pdf.MultiCell(availableWidth, 3.5, question.Question+text, "", fpdf.AlignLeft, false)
	}
	return nil
}

func printStatement(pdf *fpdf.Fpdf, statement models.Statement) {
	setPinkBoldFont(pdf, titleTextSize)
	if statement.Title != "" {
		pdf.MultiCell(0, 3.5, statement.Title, "", fpdf.AlignLeft, false)
		pdf.Ln(3)
	}
	setBlackBoldFont(pdf, standardTextSize)
	if statement.Subtitle != "" {
		pdf.MultiCell(0, 3.5, statement.Subtitle, "", fpdf.AlignLeft, false)
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
		pdf.MultiCell(0, 3.5, question.Question, "", fpdf.AlignLeft, false)
	}
	pdf.Ln(5)
	if statement.HasAnswer {
		drawSignatureForm(pdf)
		pdf.Ln(10)
	}
	checkPage(pdf)
}

func indentedText(pdf *fpdf.Fpdf, content string) {
	pdf.SetX(tabDimension)
	pdf.MultiCell(0, 3, content, "", fpdf.AlignLeft, false)
}

func checkPage(pdf *fpdf.Fpdf) {
	_, pageHeight := pdf.GetPageSize()
	if pageHeight-pdf.GetY() < 30 {
		pdf.AddPage()
	}
}
