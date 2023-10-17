package document

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-pdf/fpdf"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/ttacon/libphonenumber"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
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
	pdf.SetAutoPageBreak(true, 18)
	loadCustomFonts(pdf)
	return pdf
}

func downloadAssets() error {
	const (
		folderPath = "../tmp/assets/"
	)

	env := os.Getenv("env")
	if env == "local" {
		return nil
	}

	bucket := os.Getenv("GOOGLE_STORAGE_BUCKET")

	filesList, err := lib.ListGoogleStorageFolderContent("assets/documents")
	if err != nil {
		return err
	}
	if len(filesList) == 0 {
		return fmt.Errorf("no files found")
	}

	err = os.Mkdir(folderPath, 0750)
	if err != nil {
		return err
	}

	for _, file := range filesList {
		rawFile := lib.GetFromStorage(bucket, file, "")
		filePath := fmt.Sprintf("%s%s", folderPath, strings.SplitN(file, "/", 3)[2])
		log.Printf("[downloadAssets] write file to: %s", filePath)
		err = os.WriteFile(filePath, rawFile, 0666)
		if err != nil {
			return err
		}
	}

	files, err := os.ReadDir(folderPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		log.Println(file.Name())
	}

	return nil
}

func loadCustomFonts(pdf *fpdf.Fpdf) {
	pdf.AddUTF8Font("Montserrat", "", lib.GetAssetPathByEnvV2()+"montserrat_light.ttf")
	pdf.AddUTF8Font("Montserrat", "B", lib.GetAssetPathByEnvV2()+"montserrat_bold.ttf")
	pdf.AddUTF8Font("Montserrat", "I", lib.GetAssetPathByEnvV2()+"montserrat_italic.ttf")
}

func saveContract(pdf *fpdf.Fpdf, policy *models.Policy) (string, []byte) {
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

func saveReservedDocument(pdf *fpdf.Fpdf, policy *models.Policy) (string, []byte) {
	var (
		gsLink string
		out    bytes.Buffer
	)

	err := pdf.Output(&out)
	lib.CheckError(err)

	now := time.Now()
	timestamp := strconv.FormatInt(now.Unix(), 10)
	filename := "temp/" + policy.Uid + "/" + policy.Contractor.Name + "_" + policy.Contractor.Surname + "_" +
		timestamp + "_reserved_document.pdf"
	gsLink = lib.PutToStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), filename, out.Bytes())
	lib.CheckError(err)

	return gsLink, out.Bytes()

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
	pdf.SetTextColor(229, 9, 117)
	pdf.SetFont("Montserrat", "B", fontSize)
}

func setPinkRegularFont(pdf *fpdf.Fpdf, fontSize float64) {
	pdf.SetTextColor(229, 9, 117)
	pdf.SetFont("Montserrat", "", fontSize)
}

func setPinkItalicFont(pdf *fpdf.Fpdf, fontSize float64) {
	pdf.SetTextColor(229, 9, 117)
	pdf.SetFont("Montserrat", "I", fontSize)
}

func setPinkMonospaceFont(pdf *fpdf.Fpdf, fontSize float64) {
	pdf.SetTextColor(229, 9, 117)
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
	pdf.SetDrawColor(229, 9, 117)
	pdf.SetLineWidth(lineWidth)
	pdf.Line(11, pdf.GetY(), 200, pdf.GetY())
}

func drawSignatureForm(pdf *fpdf.Fpdf) {
	signatureID++
	text := fmt.Sprintf("\"[[!sigField\"%d\":signer1:signature(sigType=\\\"Click2Sign\\\"):label"+
		"(\\\"firma qui\\\"):size(width=150,height=60)]]\"", signatureID)
	pdf.SetDrawColor(0, 0, 0)
	pdf.SetX(-90)
	pdf.Cell(0, 3, "Firma del Contraente")
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

func checkSurveySpace(pdf *fpdf.Fpdf, survey models.Survey) {
	var answer string
	leftMargin, _, rightMargin, _ := pdf.GetMargins()
	pageWidth, pageHeight := pdf.GetPageSize()
	availableWidth := pageWidth - leftMargin - rightMargin - 2
	requiredHeight := 5.0
	currentY := pdf.GetY()

	surveyTitle := survey.Title
	surveySubtitle := survey.Subtitle

	if surveyTitle != "" {
		setPinkBoldFont(pdf, titleTextSize)
		lines := pdf.SplitText(surveyTitle, availableWidth)
		requiredHeight += 3.5 * float64(len(lines))
	}
	if surveySubtitle != "" {
		setBlackBoldFont(pdf, standardTextSize)
		lines := pdf.SplitText(surveySubtitle, availableWidth)
		requiredHeight += 3.5 * float64(len(lines))
	}

	for _, question := range survey.Questions {
		availableWidth = pageWidth - leftMargin - rightMargin - 2

		questionText := question.Question

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
		requiredHeight += 3 * float64(len(lines))
	}

	if survey.ContractorSign || survey.CompanySign {
		requiredHeight += 35
	}

	if (pageHeight-18)-currentY < requiredHeight {
		pdf.AddPage()
	}
}

func printSurvey(pdf *fpdf.Fpdf, survey models.Survey, companyName string) error {
	var dotsString string
	leftMargin, _, rightMargin, _ := pdf.GetMargins()
	pageWidth, _ := pdf.GetPageSize()
	availableWidth := pageWidth - leftMargin - rightMargin - 2

	checkSurveySpace(pdf, survey)

	surveyTitle := survey.Title
	surveySubtitle := survey.Subtitle

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

		dotsString = strings.Repeat(".", int(math.Max((paddingWidth/dotWidth)-2, 0))) + answer
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

		if question.HasAnswer {
			var questionWidth, paddingWidth float64
			answer := "NO"
			if *question.Answer {
				answer = "SI"
			}

			answerWidth := pdf.GetStringWidth(answer)
			dotWidth := pdf.GetStringWidth(".")

			lines := pdf.SplitText(question.Question+answer, availableWidth)

			questionWidth = pdf.GetStringWidth(lines[len(lines)-1])
			paddingWidth = availableWidth - questionWidth - answerWidth

			dotsString = strings.Repeat(".", int(math.Max((paddingWidth/dotWidth)-2, 0))) + answer
		}
		pdf.MultiCell(availableWidth, 3.5, question.Question+dotsString, "", fpdf.AlignLeft, false)
	}
	pdf.Ln(3)

	if survey.CompanySign {
		companySignature(pdf, companyName)
	}
	if survey.ContractorSign {
		drawSignatureForm(pdf)
		pdf.Ln(10)
	}
	return nil
}

func checkStatementSpace(pdf *fpdf.Fpdf, statement models.Statement) {
	leftMargin, _, rightMargin, _ := pdf.GetMargins()
	pageWidth, pageHeight := pdf.GetPageSize()
	availableWidth := pageWidth - leftMargin - rightMargin - 2
	requiredHeight := 5.0
	currentY := pdf.GetY()

	title := statement.Title
	subtitle := statement.Subtitle

	if title != "" {
		setPinkBoldFont(pdf, titleTextSize)
		lines := pdf.SplitText(title, availableWidth)
		requiredHeight += 3.5 * float64(len(lines))
	}
	if subtitle != "" {
		setBlackBoldFont(pdf, standardTextSize)
		lines := pdf.SplitText(subtitle, availableWidth)
		requiredHeight += 3.5 * float64(len(lines))
	}
	for _, question := range statement.Questions {
		availableWidth = pageWidth - leftMargin - rightMargin - 2

		text := question.Question

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
		requiredHeight += 3 * float64(len(lines))
	}

	if statement.ContractorSign || statement.CompanySign {
		requiredHeight += 35
	}

	if (pageHeight-18)-currentY < requiredHeight {
		pdf.AddPage()
	}
}

func printStatement(pdf *fpdf.Fpdf, statement models.Statement, companyName string) {
	checkStatementSpace(pdf, statement)

	title := statement.Title
	subtitle := statement.Subtitle

	if title != "" {
		getParagraphTitle(pdf, title)
	}
	if subtitle != "" {
		setBlackBoldFont(pdf, standardTextSize)
		pdf.MultiCell(0, 3.5, subtitle, "", fpdf.AlignLeft, false)
	}
	for _, question := range statement.Questions {
		text := question.Question
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
	pdf.Ln(3)

	if statement.CompanySign {
		companySignature(pdf, companyName)
	}
	if statement.ContractorSign {
		drawSignatureForm(pdf)
		pdf.Ln(10)
	}
}

func indentedText(pdf *fpdf.Fpdf, content string) {
	pdf.SetX(tabDimension)
	pdf.MultiCell(0, 3, content, "", fpdf.AlignLeft, false)
}

func drawDynamicCell(pdf *fpdf.Fpdf, fontSize, cellHeight, cellWidth, rowLines, nextX float64, cellText,
	innerCellBorder, outerCellBorder,
	align string, rightMost bool) {
	cellSplittedText := pdf.SplitText(cellText, cellWidth)
	cellNumLines := float64(len(cellSplittedText))

	setXY := func() {}
	ln := 1

	if !rightMost {
		setXY = func() {
			pdf.SetXY(pdf.GetX()+nextX, pdf.GetY()-(cellHeight*rowLines))
		}
		ln = 0
	}

	setBlackRegularFont(pdf, fontSize)
	if cellNumLines > 1 {
		if cellNumLines < rowLines {
			cellSplittedText = append(cellSplittedText, strings.Repeat("", int(rowLines-cellNumLines)))
		}

		for index, text := range cellSplittedText {
			if index < int(rowLines-1) {
				pdf.CellFormat(cellWidth, cellHeight, text, innerCellBorder, 2, fpdf.AlignLeft, false, 0, "")
			} else {
				pdf.CellFormat(cellWidth, cellHeight, text, outerCellBorder, 1, fpdf.AlignLeft, false, 0, "")
			}
		}
		setXY()
	} else {
		pdf.CellFormat(cellWidth, cellHeight*rowLines, cellText, outerCellBorder, ln, align, false, 0, "")
	}
}

func formatPhoneNumber(phone string) string {
	num, err := libphonenumber.Parse(phone, "IT")
	if err != nil {
		log.Printf("[DisplayPhoneNumber] error parsing phone %s", phone)
		return "================"
	}
	return libphonenumber.Format(num, libphonenumber.INTERNATIONAL)
}
