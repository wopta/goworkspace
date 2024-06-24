package quote

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"google.golang.org/api/drive/v2"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type QuoteSpreadsheet struct {
	SheetName   string
	Id          string
	InputCells  []Cell
	OutputCells []Cell
	InitCells   []Cell
}

func SpreadsheetsFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.SetPrefix("[SpreadsheetsFx] ")
	defer log.SetPrefix("")

	log.Println("Handler start -----------------------------------------------")

	qs := QuoteSpreadsheet{Id: "tn0Jqce-r_JKdecExFOFVEJdGUaPYdGo31A9FOgvt-Y"}
	res := qs.Spreadsheets()
	log.Println(res)
	log.Println("Handler end -------------------------------------------------")

	return "", nil, nil
}

func (qs *QuoteSpreadsheet) Spreadsheets() []Cell {
	var (
		path []byte
	)

	switch os.Getenv("env") {
	case "local":
		path = lib.ErrorByte(os.ReadFile("function-data/sa/positive-apex-350507-33284d6fdd55.json"))
	case "dev":
		path = lib.GetFromStorage("function-data", "sa/positive-apex-350507-33284d6fdd55.json", "")
	case "prod":
		path = lib.GetFromStorage("core-350507-function-data", "sa/positive-apex-350507-33284d6fdd55.json", "")

	default:

	}
	ctx := context.Background()
	spreadsheet := &GoogleSpreadsheet{
		CredentialsByte: path,
		Ctx:             ctx,
	}

	sheetClient, e := GoogleClient[*sheets.Service](spreadsheet)
	lib.CheckError(e)
	fmt.Printf("sheetClient: %v\n", sheetClient)
	qs.setInitCells(sheetClient, ctx)
	qs.setInputCells(sheetClient, ctx)
	return qs.getOutput(sheetClient)
}
func (qs *QuoteSpreadsheet) setInitCells(sheetClient *sheets.Service, ctx context.Context) {

	rb := &sheets.BatchUpdateValuesRequest{
		ValueInputOption: "USER_ENTERED",
	}

	for k, cell := range qs.InitCells {
		fmt.Printf("%s -> %s\n", k, cell)
		/*
			cel := &sheets.ValueRange{
				Values: [][]interface{}{{cell.Value}},
			}
			_, e = sheetClient.Spreadsheets.Values.Update(qs.Id, cell.Cell+":"+cell.Cell, cel).ValueInputOption("USER_ENTERED").Context(ctx).Do()
		*/
		rangeData := qs.SheetName + "!" + cell.Cell + ":" + cell.Cell
		rb.Data = append(rb.Data, &sheets.ValueRange{
			Range:  rangeData,
			Values: [][]interface{}{{cell.Value}},
		})

	}
	_, e := sheetClient.Spreadsheets.Values.BatchUpdate(qs.Id, rb).Context(ctx).Do()
	lib.CheckError(e)
}
func (qs *QuoteSpreadsheet) setInputCells(sheetClient *sheets.Service, ctx context.Context) {

	rb := &sheets.BatchUpdateValuesRequest{
		ValueInputOption: "USER_ENTERED",
	}

	for k, cell := range qs.InputCells {
		fmt.Printf("%s -> %s\n", k, cell)
		/*cel := &sheets.ValueRange{
			Values: [][]interface{}{{cell.Value}},
		}
		_, e = sheetClient.Spreadsheets.Values.Update(qs.Id, cell.Cell+":"+cell.Cell, cel).ValueInputOption("USER_ENTERED").Context(ctx).Do()
		*/
		rangeData := qs.SheetName + "!" + cell.Cell + ":" + cell.Cell
		rb.Data = append(rb.Data, &sheets.ValueRange{
			Range:  rangeData,
			Values: [][]interface{}{{cell.Value}},
		})
	}
	_, e := sheetClient.Spreadsheets.Values.BatchUpdate(qs.Id, rb).Context(ctx).Do()
	lib.CheckError(e)
}
func (qs *QuoteSpreadsheet) getOutput(sheetClient *sheets.Service) []Cell {
	var (
		res []Cell
	)
	col := map[string]int{"A": 0, "B": 1, "C": 2, "E": 3, "F": 4, "G": 5}
	sheet, e := sheetClient.Spreadsheets.Values.Get(qs.Id, qs.SheetName+"!A:G").Do()
	lib.CheckError(e)
	for k, cell := range qs.OutputCells {
		fmt.Printf("%s -> %s\n", k, cell)
		row, e := strconv.Atoi(string(string(cell.Cell[1:])))
		colum := cell.Cell[0:1]
		lib.CheckError(e)
		fmt.Printf("len(sheet.Values): %v\n", len(sheet.Values))
		fmt.Printf("len(sheet.Values): %v\n", sheet.Values)
		fmt.Printf("row: %v\n", row)
		fmt.Printf("row: %v\n", colum)
		fmt.Printf("value: %v\n", sheet.Values[row-1][col[colum]])
		rescell := Cell{
			Cell:  cell.Cell,
			Value: sheet.Values[row][col[colum]],
		}
		res = append(res, rescell)

	}

	return res
}

type DriveService struct {
	Svc *drive.Service
}
type GoogleSpreadsheet struct {
	CredentialsByte []byte
	Svc             *sheets.Service
	Ctx             context.Context
}
type GoogleDrive struct {
	CredentialsByte []byte
	Svc             *drive.Service
	Ctx             context.Context
}

func (s *GoogleSpreadsheet) NewClient() (*sheets.Service, error) {
	var svc *sheets.Service
	var err error

	svc, err = sheets.NewService(s.Ctx, option.WithCredentialsJSON(s.CredentialsByte), option.WithScopes(sheets.SpreadsheetsScope))

	if err != nil {
		return nil, err
	}

	return svc, nil
}

func (s *GoogleDrive) NewClient() (*DriveService, error) {
	var svc *drive.Service
	var err error
	svc, err = drive.NewService(s.Ctx, option.WithCredentialsJSON(s.CredentialsByte), option.WithScopes(drive.DriveScope))
	if err != nil {
		return nil, err
	}
	res := &DriveService{
		Svc: svc,
	}
	return res, nil
}

type GoogleService[T any] interface {
	NewClient() (T, error)
}

func GoogleClient[T any](g GoogleService[T]) (T, error) {
	return g.NewClient()
}

func getProductMock() (models.Product, error) {

	prod := models.Product{}
	return prod, nil
}
func CopySpreadsheet(path []byte, ctx context.Context, id string) (string, error) {
	googleDrive := &GoogleDrive{
		CredentialsByte: path,
		Ctx:             ctx,
	}

	driveClient, e := GoogleClient[*DriveService](googleDrive)
	lib.CheckError(e)
	fmt.Printf("driveClient: %v\n", driveClient)
	permission := &drive.Permission{
		Type:         "user",
		Role:         "writer",
		EmailAddress: "woptaassicurazioni@gmail.com",
	}
	f, e := driveClient.Svc.Files.Copy(id, &drive.File{
		Permissions: []*drive.Permission{permission},
	}).Do()
	fmt.Printf("f.Id: %v\n", e)
	fmt.Printf("f.Id: %v\n", f.Id)
	return f.Id, nil
}
