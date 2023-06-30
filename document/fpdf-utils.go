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
	pdf.MultiCell(0, titleTextSize, title, "", "", false)
}

func printSurvey(pdf *fpdf.Fpdf, survey models.Survey, companyName string) error {
	var dotsString string
	leftMargin, _, rightMargin, _ := pdf.GetMargins()
	pageWidth, _ := pdf.GetPageSize()
	availableWidth := pageWidth - leftMargin - rightMargin - 2

	checkSurveySpace(pdf, survey)

	surveyTitle := survey.Title
	surveySubtitle := survey.Subtitle

	if survey.SimploTitle != "" {
		surveyTitle = survey.SimploTitle
	}
	if survey.SimploSubtitle != "" {
		surveySubtitle = survey.SimploSubtitle
	}

	setBlackBoldFont(pdf, standardTextSize)
	if survey.HasAnswer {
		answer := "NO"
		if *survey.Answer {
			answer = "SI"
		}

		answerWidth := pdf.GetStringWidth(answer)
		dotWidth := pdf.GetStringWidth(".")

		var surveyWidth, paddingWidth float64
		var lines []string
		if surveyTitle != "" {
			lines = pdf.SplitText(surveyTitle+answer, availableWidth)
		} else if surveySubtitle != "" {
			lines = pdf.SplitText(surveySubtitle+answer, availableWidth)
		}

		surveyWidth = pdf.GetStringWidth(lines[len(lines)-1])
		paddingWidth = availableWidth - surveyWidth - answerWidth

		dotsString = strings.Repeat(".", int(paddingWidth/dotWidth)-2) + answer
	}
	if surveyTitle != "" {
		getParagraphTitle(pdf, surveyTitle+dotsString)
	}
	if surveySubtitle != "" {
		setBlackBoldFont(pdf, standardTextSize)
		pdf.MultiCell(availableWidth, 3.5, surveySubtitle+dotsString, "", fpdf.AlignLeft, false)
	}

	for _, question := range survey.Questions {
		dotsString = ""
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

		questionText := question.Question
		if question.SimploQuestion != "" {
			questionText = question.SimploQuestion
		}

		if question.HasAnswer {
			var questionWidth, paddingWidth float64
			if question.SimploQuestion != "" {

			}

			answer := "NO"
			if *question.Answer {
				answer = "SI"
			}

			answerWidth := pdf.GetStringWidth(answer)
			dotWidth := pdf.GetStringWidth(".")

			lines := pdf.SplitText(questionText+answer, availableWidth)

			questionWidth = pdf.GetStringWidth(lines[len(lines)-1])
			paddingWidth = availableWidth - questionWidth - answerWidth

			dotsString = strings.Repeat(".", int(paddingWidth/dotWidth)-2) + answer
		}
		pdf.MultiCell(availableWidth, 3.5, questionText+dotsString, "", fpdf.AlignLeft, false)
	}

	if survey.CompanySign || survey.ContractorSign {
		pdf.Ln(5)
	}

	if survey.CompanySign {
		companySignature(pdf, companyName)
	}
	if survey.ContractorSign {
		drawSignatureForm(pdf)
		pdf.Ln(10)
	}
	return nil
}

func checkSurveySpace(pdf *fpdf.Fpdf, survey models.Survey) {
	var answer string
	leftMargin, _, rightMargin, bottomMargin := pdf.GetMargins()
	pageWidth, pageHeight := pdf.GetPageSize()
	availableWidth := pageWidth - leftMargin - rightMargin - 2
	requiredHeight := 0.0

	surveyTitle := survey.Title
	surveySubtitle := survey.Subtitle

	if survey.SimploTitle != "" {
		surveyTitle = survey.SimploTitle
	}
	if survey.SimploSubtitle != "" {
		surveySubtitle = survey.SimploSubtitle
	}

	if surveyTitle != "" {
		setPinkBoldFont(pdf, titleTextSize)
		lines := pdf.SplitText(surveyTitle, availableWidth)
		requiredHeight += float64(titleTextSize * len(lines))
	}
	if surveySubtitle != "" {
		setBlackBoldFont(pdf, standardTextSize)
		lines := pdf.SplitText(surveySubtitle, availableWidth)
		requiredHeight += float64(standardTextSize * len(lines))
	}

	for _, question := range survey.Questions {
		availableWidth = pageWidth - leftMargin - rightMargin - 2

		questionText := question.Question
		if question.SimploQuestion != "" {
			questionText = question.SimploQuestion
		}

		if question.IsBold {
			setBlackBoldFont(pdf, standardTextSize)
		} else {
			setBlackRegularFont(pdf, standardTextSize)
		}
		if question.Indent {
			availableWidth -= tabDimension / 2
		}

		if question.HasAnswer {
			answer = "NO"
			if *question.Answer {
				answer = "SI"
			}
		}

		lines := pdf.SplitText(questionText+answer, availableWidth)
		requiredHeight += float64(standardTextSize * len(lines))
	}

	if survey.ContractorSign || survey.CompanySign {
		requiredHeight += 10
	}

	if (pageHeight-bottomMargin)-pdf.GetY() < requiredHeight {
		pdf.AddPage()
	}
}

func printStatement(pdf *fpdf.Fpdf, statement models.Statement, companyName string) {
	checkStatementSpace(pdf, statement)

	title := statement.Title
	if statement.SimploTitle != "" {
		title = statement.SimploTitle
	}
	subtitle := statement.Subtitle
	if statement.SimploSubtitle != "" {
		subtitle = statement.SimploSubtitle
	}

	if title != "" {
		getParagraphTitle(pdf, title)
	}
	if subtitle != "" {
		setBlackBoldFont(pdf, standardTextSize)
		pdf.MultiCell(0, 3.5, subtitle, "", fpdf.AlignLeft, false)
	}
	for _, question := range statement.Questions {
		text := question.Question
		if question.SimploQuestion != "" {
			text = question.SimploQuestion
		}
		if question.IsBold {
			setBlackBoldFont(pdf, standardTextSize)
		} else {
			setBlackRegularFont(pdf, standardTextSize)
		}
		if question.Indent {
			pdf.SetX(tabDimension)
		}
		pdf.MultiCell(0, 3.5, text, "", fpdf.AlignLeft, false)
	}
	pdf.Ln(5)
	if statement.CompanySign {
		companySignature(pdf, companyName)
	}
	if statement.ContractorSign {
		drawSignatureForm(pdf)
		pdf.Ln(10)
	}
}

func checkStatementSpace(pdf *fpdf.Fpdf, statement models.Statement) {
	leftMargin, _, rightMargin, bottomMargin := pdf.GetMargins()
	pageWidth, pageHeight := pdf.GetPageSize()
	availableWidth := pageWidth - leftMargin - rightMargin - 2
	requiredHeight := 0.0

	title := statement.Title
	if statement.SimploTitle != "" {
		title = statement.SimploTitle
	}
	subtitle := statement.Subtitle
	if statement.SimploSubtitle != "" {
		subtitle = statement.SimploSubtitle
	}

	if title != "" {
		setPinkBoldFont(pdf, titleTextSize)
		lines := pdf.SplitText(title, availableWidth)
		requiredHeight += float64(titleTextSize * len(lines))
	}
	if subtitle != "" {
		setBlackBoldFont(pdf, standardTextSize)
		lines := pdf.SplitText(subtitle, availableWidth)
		requiredHeight += float64(standardTextSize * len(lines))
	}
	for _, question := range statement.Questions {
		availableWidth = pageWidth - leftMargin - rightMargin - 2

		text := question.Question
		if question.SimploQuestion != "" {
			text = question.SimploQuestion
		}

		if question.IsBold {
			setBlackBoldFont(pdf, standardTextSize)
		} else {
			setBlackRegularFont(pdf, standardTextSize)
		}
		if question.Indent {
			availableWidth -= tabDimension / 2
		}

		answer := ""
		if question.HasAnswer {
			answer = "NO"
			if *question.Answer {
				answer = "SI"
			}
		}

		lines := pdf.SplitText(text+answer, availableWidth)
		requiredHeight += float64(standardTextSize * len(lines))
	}

	if statement.ContractorSign || statement.CompanySign {
		requiredHeight += 15
	}

	if (pageHeight-bottomMargin)-pdf.GetY() < requiredHeight {
		pdf.AddPage()
	}
}

func indentedText(pdf *fpdf.Fpdf, content string) {
	pdf.SetX(tabDimension)
	pdf.MultiCell(0, 3, content, "", fpdf.AlignLeft, false)
}
