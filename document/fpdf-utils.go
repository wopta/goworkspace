package document

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"os"
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
	proposal         = "PROPOSTA"
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

func loadCustomFonts(pdf *fpdf.Fpdf) {
	pdf.AddUTF8Font("Montserrat", "", lib.GetAssetPathByEnvV2()+"montserrat_light.ttf")
	pdf.AddUTF8Font("Montserrat", "B", lib.GetAssetPathByEnvV2()+"montserrat_bold.ttf")
	pdf.AddUTF8Font("Montserrat", "I", lib.GetAssetPathByEnvV2()+"montserrat_italic.ttf")
}

func saveContract(pdf *fpdf.Fpdf, policy *models.Policy) (string, []byte) {
	var filename string
	if os.Getenv("env") == "local" {
		err := pdf.OutputFileAndClose("./document/contract.pdf")
		lib.CheckError(err)
	} else {
		var out bytes.Buffer
		err := pdf.Output(&out)
		lib.CheckError(err)
		filename = strings.ReplaceAll(fmt.Sprintf("%s/%s/"+models.ContractDocumentFormat, "temp", policy.Uid,
			policy.NameDesc, policy.CodeCompany), " ", "_")
		lib.PutToStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), filename, out.Bytes())
		lib.CheckError(err)
		return filename, out.Bytes()
	}
	return filename, nil
}

func saveProposal(pdf *fpdf.Fpdf, policy *models.Policy) (string, []byte) {
	var (
		filename string
		out      bytes.Buffer
	)

	err := pdf.Output(&out)
	lib.CheckError(err)
	filename = strings.ReplaceAll(fmt.Sprintf("%s/%s/"+models.ProposalDocumentFormat, "temp", policy.Uid,
		policy.NameDesc, policy.ProposalNumber), " ", "_")
	lib.PutToStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), filename, out.Bytes())
	lib.CheckError(err)
	return filename, out.Bytes()
}

func saveReservedDocument(pdf *fpdf.Fpdf, policy *models.Policy) (string, []byte) {
	var (
		filename string
		out      bytes.Buffer
	)

	if os.Getenv("env") == "local" {
		err := pdf.OutputFileAndClose("./document/reserved_document.pdf")
		lib.CheckError(err)
	} else {
		err := pdf.Output(&out)
		lib.CheckError(err)
		filename = strings.ReplaceAll(fmt.Sprintf("%s/%s/"+models.RvmInstructionsDocumentFormat, "temp",
			policy.Uid, policy.NameDesc, policy.ProposalNumber), " ", "_")
		lib.PutToStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), filename, out.Bytes())
		lib.CheckError(err)
		return filename, out.Bytes()
	}
	return filename, nil
}

func insertWatermark(pdf *fpdf.Fpdf, text string) {
	currentY := pdf.GetY()
	markFontHt := 115.0
	markLineHt := pdf.PointToUnitConvert(markFontHt)
	markY := (297.0 - markLineHt) / 2.0
	ctrX := 210.0 / 2.0
	ctrY := 297.0 / 2.0
	pdf.SetFont("Arial", "B", markFontHt)
	pdf.SetTextColor(206, 216, 232)
	pdf.SetXY(10, markY)
	pdf.TransformBegin()
	pdf.TransformRotate(-45, ctrX, ctrY)
	pdf.CellFormat(0, markLineHt, text, "", 0, "C", false, 0, "")
	pdf.TransformEnd()
	pdf.SetXY(10, currentY)
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

func checkSurveySpace(pdf *fpdf.Fpdf, survey models.Survey, isProposal bool) {
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

	if (!isProposal && survey.ContractorSign) || survey.CompanySign {
		requiredHeight += 35
	}

	if (pageHeight-18)-currentY < requiredHeight {
		pdf.AddPage()
	}
}

func printSurvey(pdf *fpdf.Fpdf, survey models.Survey, companyName string, isProposal bool) error {
	var dotsString string
	leftMargin, _, rightMargin, _ := pdf.GetMargins()
	pageWidth, _ := pdf.GetPageSize()
	availableWidth := pageWidth - leftMargin - rightMargin - 2

	checkSurveySpace(pdf, survey, isProposal)

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
		if isProposal {
			pdf.Ln(20)
		}
	}
	if !isProposal && survey.ContractorSign {
		drawSignatureForm(pdf)
		pdf.Ln(10)
	}
	return nil
}

func checkStatementSpace(pdf *fpdf.Fpdf, statement models.Statement, isProposal bool) {
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

	if (!isProposal && statement.ContractorSign) || statement.CompanySign {
		requiredHeight += 35
	}

	if (pageHeight-18)-currentY < requiredHeight {
		pdf.AddPage()
	}
}

func printStatement(pdf *fpdf.Fpdf, statement models.Statement, companyName string, isProposal bool) {
	checkStatementSpace(pdf, statement, isProposal)

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
		if isProposal {
			pdf.Ln(20)
		}
	}
	if !isProposal && statement.ContractorSign {
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

func woptaInfoTable(pdf *fpdf.Fpdf, producerInfo map[string]string) {
	drawPinkHorizontalLine(pdf, 0.1)
	setBlackRegularFont(pdf, smallTextSize)
	pdf.MultiCell(0, 5, "DATI DELLA PERSONA FISICA CHE ENTRA IN CONTATTO CON IL "+
		"CONTRAENTE", "", "", false)
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 5, producerInfo["name"]+" iscritto alla Sezione "+
		producerInfo["ruiSection"]+" del RUI con numero "+producerInfo["ruiCode"]+" in data "+
		producerInfo["ruiRegistration"], "", "", false)
	drawPinkHorizontalLine(pdf, 0.1)
	setBlackRegularFont(pdf, smallTextSize)
	pdf.MultiCell(0, 5, "QUALIFICA", "", "", false)
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3.5, "Responsabile dell’attività di intermediazione assicurativa di Wopta "+
		"Assicurazioni Srl, Società iscritta alla Sezione A del RUI con numero A000701923 in data "+
		"14.02.2022", "", "", false)
	drawPinkHorizontalLine(pdf, 0.1)
	setBlackRegularFont(pdf, smallTextSize)
	pdf.MultiCell(0, 5, "SEDE LEGALE", "", "", false)
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 5, "Galleria del Corso, 1 – 20122 MILANO (MI)", "", "", false)
	drawPinkHorizontalLine(pdf, 0.1)
	setBlackRegularFont(pdf, smallTextSize)
	pdf.CellFormat(90, 5, "RECAPITI TELEFONICI", "", 0, fpdf.AlignLeft, false, 0, "")
	pdf.CellFormat(90, 5, "E-MAIL", "", 1, fpdf.AlignLeft, false, 0, "")
	setBlackRegularFont(pdf, standardTextSize)
	pdf.CellFormat(90, 5, "02.91.24.03.46", "", 0, fpdf.AlignLeft, false, 0, "")
	pdf.CellFormat(90, 5, "info@wopta.it", "", 1, fpdf.AlignLeft, false, 0, "")
	drawPinkHorizontalLine(pdf, 0.1)
	setBlackRegularFont(pdf, smallTextSize)
	pdf.CellFormat(90, 5, "PEC", "", 0, fpdf.AlignLeft, false, 0, "")
	pdf.CellFormat(90, 5, "SITO INTERNET", "", 1, fpdf.AlignLeft, false, 0, "")
	setBlackRegularFont(pdf, standardTextSize)
	pdf.CellFormat(90, 5, "woptaassicurazioni@legalmail.it", "", 0, fpdf.AlignLeft, false, 0, "")
	pdf.CellFormat(90, 5, "wopta.it", "", 1, fpdf.AlignLeft, false, 0, "")
	drawPinkHorizontalLine(pdf, 0.1)
	setBlackRegularFont(pdf, smallTextSize)
	pdf.MultiCell(0, 5, "AUTORITÀ COMPETENTE ALLA VIGILANZA DELL’ATTIVITÀ SVOLTA",
		"", "", false)
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 5, "IVASS – Istituto per la Vigilanza sulle Assicurazioni - Via del Quirinale, "+
		"21 - 00187 Roma", "", "", false)
	drawPinkHorizontalLine(pdf, 0.1)
}

func getProducerName(networkNode *models.NetworkNode) string {
	var producerName string

	if networkNode == nil {
		return producerName
	}

	switch networkNode.Type {
	case models.AgentNetworkNodeType:
		producerName = strings.ToUpper(fmt.Sprintf("%s %s", networkNode.Agent.Surname, networkNode.Agent.Name))
	case models.AgencyNetworkNodeType:
		producerName = strings.ToUpper(fmt.Sprintf("%s", networkNode.Agency.Name))
	case models.BrokerNetworkNodeType:
		producerName = strings.ToUpper(fmt.Sprintf("%s", networkNode.Broker.Name))
	}
	return producerName
}
