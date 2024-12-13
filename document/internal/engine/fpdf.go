package engine

import (
	"bytes"
	"os"
	"time"

	"github.com/go-pdf/fpdf"
	"github.com/mohae/deepcopy"
	"github.com/wopta/goworkspace/document/internal/constants"
	"github.com/wopta/goworkspace/document/internal/domain"
	"github.com/wopta/goworkspace/lib"
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

	pdf.SetMargins(10, 15, 10)
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
		size:        constants.RegularFontsize,
		textColor:   constants.BlackColor,
		drawColor:   constants.BlackColor,
		fillColor:   constants.WhiteColor,
	}
}

func (f *Fpdf) NewPage() {
	f.pdf.AddPage()
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

func (f *Fpdf) SetY(y float64) {
	f.pdf.SetY(y)
}

func (f *Fpdf) DrawWatermark(text string) {
	currentY := f.pdf.GetY()
	markFontHt := 115.0
	markLineHt := f.pdf.PointToUnitConvert(markFontHt)
	markY := (297.0 - markLineHt) / 2.0
	ctrX := 210.0 / 2.0
	ctrY := 297.0 / 2.0
	f.pdf.SetFont("Arial", "B", markFontHt)
	f.pdf.SetTextColor(206, 216, 232)
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
		f.fillColor = cell.FillColor
		f.SetFillColor(cell.FillColor)
	}
	if cell.TextBold {
		f.style = constants.BoldFontStyle
		f.SetFontStyle(constants.BoldFontStyle)
	}

	f.pdf.MultiCell(cell.Width, cell.Height, cell.Text, cell.Border, cell.Align, cell.Fill)

	f.SetFillColor(oldFillColor)
	f.SetFontStyle(oldFontStyle)
}

func (f *Fpdf) DrawTable(table [][]domain.TableCell) {
	for _, row := range table {
		maxNumLines := 1
		for _, cell := range row {
			f.SetFontStyle(constants.RegularFontStyle)
			if cell.TextBold {
				f.SetFontStyle(constants.BoldFontStyle)
			}
			numLines := len(f.pdf.SplitText(cell.Text, cell.Width))
			if numLines > maxNumLines {
				maxNumLines = numLines
			}
		}

		for index, cell := range row {
			numLines := len(f.pdf.SplitText(cell.Text, cell.Width))
			currentX, currentY := f.pdf.GetXY()

			f.SetFontFamily(constants.MontserratFont)

			f.WriteText(cell)

			emptyCell := deepcopy.Copy(cell).(domain.TableCell)
			emptyCell.Text = ""
			emptyCell.Border = ""

			for i := numLines; i < maxNumLines; i++ {
				f.pdf.SetX(currentX)
				f.WriteText(emptyCell)
			}

			if index < len(row)-1 {
				f.pdf.SetXY(currentX+cell.Width, currentY)
			}
		}
	}
}

func (f *Fpdf) Save(filePath string) error {
	var out bytes.Buffer
	err := f.pdf.Output(&out)
	if err != nil {
		return err
	}
	lib.PutToStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), filePath, out.Bytes())
	return nil
}
