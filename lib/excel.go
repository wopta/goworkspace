package lib

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strconv"

	"github.com/xuri/excelize/v2"
)

func ExcelRead(r io.Reader) (map[string][][]string, error) {
	var res map[string][][]string
	var rows [][]string
	var err error
	f, err := excelize.OpenReader(r, excelize.Options{})
	for _, sheet := range f.GetSheetList() {
		rows, err = f.GetRows(sheet)
		res[sheet] = rows
		for _, colCell := range rows {
			fmt.Print(colCell, "\t")
		}
		fmt.Println()
	}
	return res, err
}
func ExcelReadFile(filePath string) (map[string][][]string, error) {
	// f, err := excelize.OpenFile("Book1.xlsx")
	var res map[string][][]string
	var rows [][]string
	var err error
	f, err := excelize.OpenFile(filePath)
	for _, sheet := range f.GetSheetList() {
		rows, err = f.GetRows(sheet)
		res[sheet] = rows
		for _, colCell := range rows {
			fmt.Print(colCell, "\t")
		}
		fmt.Println()
	}
	return res, err
}
func CreateExcel(sheet [][]string, filePath string, sheetName string) (<-chan []byte, error) {
	ch := make(chan []byte)
	var err error
	var index int
	var resByte *bytes.Buffer
	go func() {
		log.Println("CreateExcel")
		f := excelize.NewFile()
		alfabet := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
			"AA", "AB", "AC", "AD", "AE", "AF", "AG", "AH", "AI", "AJ", "AK", "AL", "AM", "AN", "AO", "AP", "AQ", "AR", "AS", "AT", "AU", "AV", "AW", "AX", "AY", "AZ",
			"BA", "BB", "BC", "BD", "BE", "BF", "BG", "BH", "BI", "BJ", "BK", "BL", "BM", "BN", "BO", "BP", "BQ", "BR", "BS", "BT", "BU", "BV", "BW", "BX", "BY", "BZ"}
		// Create a new sheet.
		f.SetSheetName("Sheet1", sheetName)
		index, err = f.NewSheet(sheetName)
		for x, row := range sheet {

			for i, cel := range row {

				f.SetCellValue(sheetName, alfabet[i]+""+strconv.Itoa(x+1), cel)
			}
		}
		//Set active sheet of the workbook.
		f.SetActiveSheet(index)

		//Save spreadsheet by the given path.
		err = f.SaveAs(filePath)

		resByte, err = f.WriteToBuffer()
		ch <- resByte.Bytes()
	}()
	return ch, err
}
