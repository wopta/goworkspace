package engine

import (
	"bytes"
	"os"
	"strings"
	"time"

	"github.com/go-pdf/fpdf"
	"github.com/mohae/deepcopy"
	"gitlab.dev.wopta.it/goworkspace/document/internal/constants"
	"gitlab.dev.wopta.it/goworkspace/document/internal/domain"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
)

var (
	alignmentMap = map[domain.TextAlign]string{
		constants.LeftAlign:   fpdf.AlignLeft,
		constants.CenterAlign: fpdf.AlignCenter,
		constants.RightAlign:  fpdf.AlignRight,
	}
)

type Fpdf struct {
	pdf         *fpdf.Fpdf
	signatureID int64
	font        domain.FontFamily
	style       domain.FontStyle
	size        domain.FontSize
	textColor   domain.Color
	drawColor   domain.Color
	fillColor   domain.Color
}

func NewFpdf() *Fpdf {
	pdf := fpdf.New(fpdf.OrientationPortrait, fpdf.UnitMillimeter, fpdf.PageSizeA4, "")

	pdf.SetMargins(10, 15, -1)
	pdf.SetAuthor("Wopta Assicurazioni s.r.l", true)
	pdf.SetCreationDate(time.Now().UTC())
	pdf.SetAutoPageBreak(true, 18)
	pdf.AddUTF8Font(string(constants.MontserratFont), string(constants.RegularFontStyle), lib.GetAssetPathByEnvV2()+"montserrat_light.ttf")
	pdf.AddUTF8Font(string(constants.MontserratFont), string(constants.BoldFontStyle), lib.GetAssetPathByEnvV2()+"montserrat_bold.ttf")
	pdf.AddUTF8Font(string(constants.MontserratFont), string(constants.ItalicFontStyle), lib.GetAssetPathByEnvV2()+"montserrat_italic.ttf")

	return &Fpdf{
		pdf:         pdf,
		signatureID: 0,
		font:        constants.MontserratFont,
		style:       constants.RegularFontStyle,
		size:        constants.RegularFontSize,
		textColor:   constants.BlackColor,
		drawColor:   constants.BlackColor,
		fillColor:   constants.WhiteColor,
	}
}

func (f *Fpdf) GetPdf() *fpdf.Fpdf {
	return f.pdf
}

func (f *Fpdf) NewPage() {
	f.pdf.AddPage()
	f.NewLine(5)
}

func (f *Fpdf) NewLine(space float64) {
	f.pdf.Ln(space)
}

func (f *Fpdf) ResetSignatureID() {
	f.signatureID = 0
}

func (f *Fpdf) IncrementSignatureID() {
	f.signatureID++
}

func (f *Fpdf) SetFontFamily(font domain.FontFamily) {
	f.font = font
	f.pdf.SetFont(string(f.font), string(f.style), float64(f.size))
}

func (f *Fpdf) SetFontStyle(style domain.FontStyle) {
	f.style = style
	f.pdf.SetFont(string(f.font), string(f.style), float64(f.size))
}

func (f *Fpdf) SetFontSize(size domain.FontSize) {
	f.size = size
	f.pdf.SetFont(string(f.font), string(f.style), float64(f.size))
}

func (f *Fpdf) SetFontColor(color domain.Color) {
	f.textColor = color
	f.pdf.SetTextColor(int(f.textColor.R), int(f.textColor.G), int(f.textColor.B))
}

func (f *Fpdf) SetDrawColor(color domain.Color) {
	f.drawColor = color
	f.pdf.SetDrawColor(int(f.drawColor.R), int(f.drawColor.G), int(f.drawColor.B))
}

func (f *Fpdf) SetFillColor(color domain.Color) {
	f.fillColor = color
	f.pdf.SetFillColor(int(f.fillColor.R), int(f.fillColor.G), int(f.fillColor.B))
}

func (f *Fpdf) SetX(x float64) {
	f.pdf.SetX(x)
}

func (f *Fpdf) GetX() float64 {
	return f.pdf.GetX()
}

func (f *Fpdf) SetY(y float64) {
	f.pdf.SetY(y)
}

func (f *Fpdf) GetY() float64 {
	return f.pdf.GetY()
}

func (f *Fpdf) DrawWatermark(text string) {
	currentY := f.pdf.GetY()
	markFontHt := 115.0
	markLineHt := f.pdf.PointToUnitConvert(markFontHt)
	markY := (297.0 - markLineHt) / 2.0
	ctrX := 210.0 / 2.0
	ctrY := 297.0 / 2.0
	f.pdf.SetFont(string(constants.ArialFont), string(constants.BoldFontStyle), markFontHt)
	f.pdf.SetTextColor(int(constants.WatermarkColor.R), int(constants.WatermarkColor.G), int(constants.WatermarkColor.B))
	f.pdf.SetXY(10, markY)
	f.pdf.TransformBegin()
	f.pdf.TransformRotate(-45, ctrX, ctrY)
	f.pdf.CellFormat(0, markLineHt, text, "", 0, "C", false, 0, "")
	f.pdf.TransformEnd()
	f.pdf.SetXY(10, currentY)
}

func (f *Fpdf) PageNumber() int {
	return f.pdf.PageNo()
}

func (f *Fpdf) WriteText(cell domain.TableCell) {
	oldFontStyle := f.style
	oldFillColor := f.fillColor

	if cell.Fill {
		f.SetFillColor(cell.FillColor)
	}

	f.SetFontStyle(cell.FontStyle)
	f.SetFontColor(cell.FontColor)
	f.SetFontSize(cell.FontSize)
	f.SetFontFamily(constants.MontserratFont)

	f.pdf.MultiCell(cell.Width, cell.Height, cell.Text, cell.Border, alignmentMap[cell.Align], cell.Fill)

	f.SetFillColor(oldFillColor)
	f.SetFontStyle(oldFontStyle)
}

func (f *Fpdf) RawWriteText(cell domain.TableCell) {
	oldFontStyle := f.style
	oldFillColor := f.fillColor

	if cell.Fill {
		f.SetFillColor(cell.FillColor)
	}

	f.SetFontStyle(cell.FontStyle)
	f.SetFontColor(cell.FontColor)
	f.SetFontSize(cell.FontSize)
	f.SetFontFamily(constants.MontserratFont)

	f.pdf.Write(cell.Height, cell.Text)

	f.SetFillColor(oldFillColor)
	f.SetFontStyle(oldFontStyle)
}

func (f *Fpdf) WriteLink(url string, cell domain.TableCell) {
	oldFontStyle := f.style
	oldFillColor := f.fillColor

	if cell.Fill {
		f.SetFillColor(cell.FillColor)
	}

	f.SetFontStyle(cell.FontStyle)
	f.SetFontColor(cell.FontColor)
	f.SetFontSize(cell.FontSize)
	f.SetFontFamily(constants.MontserratFont)

	f.pdf.WriteLinkString(cell.Height, cell.Text, url)

	f.SetFillColor(oldFillColor)
	f.SetFontStyle(oldFontStyle)
}

func (f *Fpdf) DrawTable(table [][]domain.TableCell) {
	for _, row := range table {
		maxNumLines := 1
		for _, cell := range row {
			f.SetFontStyle(cell.FontStyle)
			f.SetFontSize(cell.FontSize)
			f.SetFontFamily(constants.MontserratFont)

			numLines := len(f.pdf.SplitText(cell.Text, cell.Width))
			if numLines > maxNumLines {
				maxNumLines = numLines
			}
		}

		for index, cell := range row {
			numLines := len(f.pdf.SplitText(cell.Text, cell.Width))
			currentX, currentY := f.pdf.GetXY()

			hasBottomBorder := strings.Contains(cell.Border, "B") || cell.Border == "1"

			if numLines < maxNumLines {
				switch cell.Border {
				case "1":
					cell.Border = "TLR"
				default:
					cell.Border = strings.ReplaceAll(cell.Border, "B", "")
				}
			}
			f.WriteText(cell)

			emptyCell := deepcopy.Copy(cell).(domain.TableCell)
			emptyCell.Text = ""
			emptyCell.Border = strings.ReplaceAll(emptyCell.Border, "T", "")

			for i := numLines; i < maxNumLines; i++ {
				if i == maxNumLines-1 && hasBottomBorder {
					emptyCell.Border += "B"
				}
				f.pdf.SetX(currentX)
				f.WriteText(emptyCell)
			}

			if index < len(row)-1 {
				f.pdf.SetXY(currentX+cell.Width, currentY)
			}
		}
	}
}

func (f *Fpdf) RawDoc() ([]byte, error) {
	var out bytes.Buffer
	err := f.pdf.Output(&out)
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func (f *Fpdf) SetHeader(header func()) {
	f.pdf.SetHeaderFunc(header)
}

func (f *Fpdf) SetFooter(footer func()) {
	f.pdf.SetFooterFunc(footer)
}

func (f *Fpdf) InsertImage(imagePath string, x, y, width, height float64) {
	var opt fpdf.ImageOptions

	opt.ImageType = "png"
	f.pdf.ImageOptions(imagePath, x, y, width, height, false, opt, 0, "")
}

func (f *Fpdf) DrawLine(startX, startY, endX, endY, lineWidth float64, color domain.Color) {
	f.SetDrawColor(color)
	f.pdf.SetLineWidth(lineWidth)
	f.pdf.Line(startX, startY, endX, endY)
}

func (f *Fpdf) GetMargins() (float64, float64, float64, float64) {
	return f.pdf.GetMargins()
}

func (f *Fpdf) GetPageSize() (float64, float64) {
	return f.pdf.GetPageSize()
}

func (f *Fpdf) SplitText(text string, width float64) []string {
	return f.pdf.SplitText(text, width)
}

func (f *Fpdf) GetStringWidth(text string) float64 {
	return f.pdf.GetStringWidth(text)
}

// TODO: this shouldnt be a method from the generator but an action invoked by the callee
func (f *Fpdf) Save(rawDoc []byte, filename string) (string, error) {
	return lib.PutToStorageErr(os.Getenv("GOOGLE_STORAGE_BUCKET"), filename, rawDoc)
}

// get a tablecell personalized based on passed opts
func (e *Fpdf) GetTableCell(text string, opts ...any) domain.TableCell {
	tableCell := domain.TableCell{}
	tableCell.Text = text
	tableCell.Height = constants.CellHeight
	tableCell.Align = constants.LeftAlign
	tableCell.FontStyle = constants.RegularFontStyle
	tableCell.FontSize = constants.RegularFontSize

	for _, opt := range opts {
		switch opt := opt.(type) {
		case domain.FontSize:
			tableCell.FontSize = opt
		case domain.FontStyle:
			tableCell.FontStyle = opt
		case domain.Color:
			tableCell.FontColor = opt
		case domain.TextAlign:
			tableCell.Align = opt
		}
	}
	return tableCell
}

func (e *Fpdf) WriteTexts(tables ...domain.TableCell) {
	for _, text := range tables {
		log.ErrorF("prova", text.Link)
		if text.Link == "" {
			e.RawWriteText(text)
		} else {
			e.WriteLink(text.Link, e.GetTableCell(text.Text, constants.PinkColor))
		}
	}
	e.NewLine(constants.CellHeight)
}

func (f *Fpdf) CrossRemainingSpace() {
	var (
		minimumAreaHeight     = float64(5)
		currentX, currentY    = f.pdf.GetXY()
		_, _, _, bottomMargin = f.pdf.GetMargins()
		_, pageHeight         = f.GetPageSize()
		maximumY              = pageHeight - bottomMargin - minimumAreaHeight
	)

	if currentY > maximumY {
		return
	}

	f.DrawLine(currentX, currentY, constants.FullPageWidth, maximumY, 0.25, constants.BlackColor)
}

func (f *Fpdf) SetMargins(left, top, right float64) {
	f.pdf.SetMargins(left, top, right)
}
