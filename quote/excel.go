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
	fmt.Println("-------Excel---------")
	filePath := "quote/excel/testFx.xlsx"
	filePathOut := "../tmp/temp.xlsx"
	excelBytes := lib.GetFilesByEnv(filePath)
	xlsx, err := excelize.OpenReader(bytes.NewReader(excelBytes))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		// Close the spreadsheet.
		if err := xlsx.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	xlsx.SetCellValue("Tabelle1", "A1", 100)
	// Get value from cell by given worksheet name and cell reference.
	cell, err := xlsx.GetCellValue("Tabelle1", "E1")
	fmt.Println("excel get value. ", cell)
	err = xlsx.UpdateLinkedValue()

	err = xlsx.SaveAs(filePathOut)
	fmt.Println(err)
	fmt.Println("excel get value: ", cell)
	xlsxOut, err := excelize.OpenFile(filePathOut)

	cell, err = xlsxOut.GetCellValue("Tabelle1", "E1")
	fmt.Println(err)
	fmt.Println("excel get value out: ", cell)

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(cell)
	// Get all the rows in the Sheet1.

}
