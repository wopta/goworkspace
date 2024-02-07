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
	defer xlsx.Close()
	xlsx.SetCellValue("Tabelle1", "A1", 100)
	// Get value from cell by given worksheet name and cell reference.
	cell, err := xlsx.GetCellValue("Tabelle1", "E1")
	cell1, err := xlsx.GetCellValue("Tabelle1", "A1")
	fmt.Println("excel get value E1: ", cell)
	fmt.Println("excel get value A1: ", cell1)
	err = xlsx.UpdateLinkedValue()

	<-SaveExcel(xlsx, filePathOut)
	fmt.Println(err)
	fmt.Println("excel get value: ", cell)
	xlsxOut, err := excelize.OpenFile(filePathOut)
	fmt.Println(err)
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
func SaveExcel(xlsx *excelize.File, filePath string) <-chan []byte {
	ch := make(chan []byte)
	var err error

	var resByte *bytes.Buffer
	go func() {

		//Save spreadsheet by the given path.
		err = xlsx.SaveAs(filePath)
		fmt.Println(err)
		resByte, err = xlsx.WriteToBuffer()
		fmt.Println("excel Saved Excel ")
		ch <- resByte.Bytes()
	}()
	return ch
}
