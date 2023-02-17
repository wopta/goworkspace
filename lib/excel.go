package lib

import (
	"fmt"
	"log"
	"strconv"

	"github.com/xuri/excelize/v2"
)

func CreateExcel(sheet [][]interface{}, filePath string) {
	log.Println("CreateExcel")
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// Create a new sheet.
	index, err := f.NewSheet("Sheet1")
	if err != nil {
		fmt.Println(err)
		return
	}
	// Set value of a cell.

	for x, row := range sheet {
		for i, cel := range row {
			alfabet := rune('A' - 1 + i)
			f.SetCellValue("Sheet1", string(alfabet)+""+strconv.Itoa(x), cel)
		}
	}

	// Set active sheet of the workbook.
	f.SetActiveSheet(index)
	// Save spreadsheet by the given path.

	if err := f.SaveAs(filePath); err != nil {
		fmt.Println(err)
	}
	Files("")
}
