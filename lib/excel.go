package lib

import (
	"fmt"
	"io"
	"log"
	"strconv"

	"github.com/xuri/excelize/v2"
)

func CreateExcel(sheet [][]interface{}, filePath string, sheetName string) ([]byte, error) {
	log.Println("CreateExcel")
	f := excelize.NewFile()
	alfabet := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
	// Create a new sheet.
	index, err := f.NewSheet(sheetName)
	for x, row := range sheet {
		for i, cel := range row {

			fmt.Println(cel)
			err = f.SetCellValue(sheetName, alfabet[i]+""+strconv.Itoa(x+1), cel)
		}
	}
	//Set active sheet of the workbook.
	f.SetActiveSheet(index)
	//Save spreadsheet by the given path.
	err = f.SaveAs(filePath)

	resByte, err := f.WriteToBuffer()

	return resByte.Bytes(), err
}

func ExcelRead(r io.Reader) (map[string][][]string, error) {
	// f, err := excelize.OpenFile("Book1.xlsx")
	var res map[string][][]string
	var rows [][]string
	var err error
	f, err := excelize.OpenReader(r, excelize.Options{})

	// Get value from cell by given worksheet name and cell reference.
	cell, err := f.GetCellValue("Sheet1", "B2")

	fmt.Println(cell)
	// Get all the rows in the Sheet1.

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
