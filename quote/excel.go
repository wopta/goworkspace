package quote

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/wopta/goworkspace/lib"
	"github.com/xuri/excelize/v2"
)

func ExcelFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.SetPrefix("[ExcelFx] ")
	defer log.SetPrefix("")
	log.Println("Handler start -----------------------------------------------")

	Excel()

	log.Println("Handler end -------------------------------------------------")

	return "", nil, nil
}

type Cell struct {
	Cell  string
	Value interface{}
}

type QuoteExcel struct {
	SheetName, filename string
	InputCells          []Cell
}

func Excel() {
	fmt.Println("-------Excel---------")
	filePath := "quote/excel/qbeRatingModel.xlsx"
	//filePathOut := "../tmp/temp.xlsx"
	sheet := "Input dati Polizza"

	excelBytes := lib.GetFilesByEnv(filePath)
	xlsx, err := excelize.OpenReader(bytes.NewReader(excelBytes))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer xlsx.Close()
	xlsx.SetCellValue(sheet, "C24", 20)

	// Get value from cell by given worksheet name and cell reference.

	cell1, err := xlsx.GetCellValue(sheet, "C91")
	// <-SaveExcel(xlsx, filePathOut)
	fmt.Println(err)
	fmt.Println("excel get value E1: ", cell1)
	err = xlsx.UpdateLinkedValue()
	fmt.Println(err)
	cell1, err = xlsx.GetCellValue(sheet, "C91")
	// <-SaveExcel(xlsx, filePathOut)
	fmt.Println(err)
	fmt.Println("excel get value UpdateLinkedValue E1: ", cell1)
	calc, err := xlsx.CalcCellValue(sheet, "C91")
	//<-SaveExcel(xlsx, filePathOut)
	fmt.Println(err)

	fmt.Println("excel get calc value E1: ", calc)

	//

	//xlsxOut, err := excelize.OpenFile(filePathOut)
	fmt.Println(err)

	fmt.Println(err)
	//fmt.Println("excel get value out: ", cell)

	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Println(cell)
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
		lib.PutToStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "quote/excel/testout.xlsx", resByte.Bytes())
		ch <- resByte.Bytes()
	}()
	return ch
}
