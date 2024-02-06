package quote

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/wopta/goworkspace/lib"
	"github.com/xuri/excelize/v2"
)

func ExcelFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	Excel()
	return "", nil, nil
}

type InputCell struct {
	cell  string
	value interface{}
}

type QuoteExcel struct {
	SheetName, filename string
	InputCells          []InputCell
}

func Excel() {
	filePath := "quote/excel/testFx.xlsx"
	excelBytes := lib.GetFilesByEnv(filePath)
	f, err := excelize.OpenReader(bytes.NewReader(excelBytes))
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
	f.SetCellValue("Tabelle1", "A1", 100)
	// Get value from cell by given worksheet name and cell reference.
	cell, err := f.GetCellValue("Tabelle1", "E1")
	fmt.Println(cell)
	err = f.UpdateLinkedValue()
	fmt.Println(err)
	cell, err = f.GetCellValue("Tabelle1", "E1")
	fmt.Println(cell)
	f.Save()
	fmt.Println(cell)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(cell)
	// Get all the rows in the Sheet1.

}
