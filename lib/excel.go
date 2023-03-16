package lib

import (
	"fmt"
	"io"
	"log"
	"strconv"

	"github.com/xuri/excelize/v2"
)

func CreateExcel(sheet [][]interface{}, filePath string) ([]byte, error) {
	log.Println("CreateExcel")
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// Create a new sheet.
	index, err := f.NewSheet("Sheet1")
	for x, row := range sheet {
		for i, cel := range row {
			alfabet := rune('A' - 1 + i)
			fmt.Println(string(alfabet) + "" + strconv.Itoa(x))
			fmt.Println(cel)
			f.SetCellValue("Sheet1", string(alfabet)+""+strconv.Itoa(x), cel)
		}
	}
	//Set active sheet of the workbook.
	f.SetActiveSheet(index)
	//Save spreadsheet by the given path.
	err = f.SaveAs(filePath)
	resByte, err := f.WriteToBuffer()
	return resByte.Bytes(), err
}

func ExcelRead(r io.Reader) {
	// f, err := excelize.OpenFile("Book1.xlsx")
	f, err := excelize.OpenReader(r, excelize.Options{})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// Get value from cell by given worksheet name and cell reference.
	cell, err := f.GetCellValue("Sheet1", "B2")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(cell)
	// Get all the rows in the Sheet1.
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, row := range rows {
		for _, colCell := range row {
			fmt.Print(colCell, "\t")
		}
		fmt.Println()
	}
}
