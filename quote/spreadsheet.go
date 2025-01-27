package quote

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

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

var (
	originalSheetId, exportSheetId int64
)

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
		path               []byte
		destinationSheetId = "1tMi7NYFZu7AnV4WkVrD0yzy1Dt3d-wVs0iZwlOcxLrg"
		bucketSavePath     = "test/download/"
	)

	switch os.Getenv("env") {
	case "local":
		path = lib.ErrorByte(os.ReadFile("../function-data/dev/sa/positive-apex-350507-33284d6fdd55.json"))
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
	res := qs.getOutput(sheetClient)

	ssRes, _ := sheetClient.Spreadsheets.Get(qs.Id).Context(ctx).Do()
	for _, s := range ssRes.Sheets {
		if s.Properties.Title == qs.SheetName {
			originalSheetId = s.Properties.SheetId
		}
		if s.Properties.Title == "Export" {
			exportSheetId = s.Properties.SheetId
		}
	}

	sheetsToClear := make([]*sheets.Request, 0)
	ssRes, _ = sheetClient.Spreadsheets.Get(destinationSheetId).Context(ctx).Do()
	for _, s := range ssRes.Sheets {
		if s.Properties.Title != "Sheet1" {
			ds := sheets.DeleteSheetRequest{SheetId: s.Properties.SheetId}
			sr := sheets.Request{DeleteSheet: &ds}
			sheetsToClear = append(sheetsToClear, &sr)
		}
	}
	_, err := sheetClient.Spreadsheets.BatchUpdate(destinationSheetId, &sheets.BatchUpdateSpreadsheetRequest{Requests: sheetsToClear}).Context(ctx).Do()
	if err != nil {
		log.Printf("unable to delete sheets from spreadsheet: %v", err)
	}

	_, err = sheetClient.Spreadsheets.Sheets.CopyTo(qs.Id, exportSheetId, &sheets.CopySheetToAnotherSpreadsheetRequest{
		DestinationSpreadsheetId: destinationSheetId,
	}).Context(ctx).Do()

	// load from drive and save to bucket
	doc, err := loadFromDrive(path, ctx, destinationSheetId)
	if err != nil {
		log.Printf("unable to load from GDrive: %v", err)
		return res
	}
	err = saveToBucket(bucketSavePath+"proposta_"+time.Now().Format("2006-1-2_15:4:5")+".xls", doc)
	if err != nil {
		log.Printf("unable to save to bucket: %v", err)
		return res
	}

	return res
}

func (qs *QuoteSpreadsheet) setInitCells(sheetClient *sheets.Service, ctx context.Context) {

	rb := &sheets.BatchUpdateValuesRequest{
		ValueInputOption: "USER_ENTERED",
	}

	for _, cell := range qs.InitCells {
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

	for _, cell := range qs.InputCells {
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

	for _, cell := range qs.OutputCells {
		row, e := strconv.Atoi(string(string(cell.Cell[1:])))
		colum := cell.Cell[0:1]
		lib.CheckError(e)
		fmt.Printf("value: %v\n", sheet.Values[row-1][col[colum]])
		rescell := Cell{
			Cell:  cell.Cell,
			Value: sheet.Values[row-1][col[colum]],
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

func loadFromDrive(path []byte, ctx context.Context, fileId string) ([]byte, error) {
	googleDrive := &GoogleDrive{
		CredentialsByte: path,
		Ctx:             ctx,
	}

	driveClient, err := GoogleClient[*DriveService](googleDrive)
	if err != nil {
		return nil, fmt.Errorf("unable to create GDrive client: %v", err)
	}

	resp, err := driveClient.Svc.Files.Export(fileId, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet").Download()
	if err != nil {
		return nil, fmt.Errorf("error exporting file from GDrive: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting response: http status is %d (%s)", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading resp body: %v", err)
	}

	return body, nil
}

func saveToBucket(path string, file []byte) error {
	_, err := lib.PutToGoogleStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), path, file)
	if err != nil {
		return fmt.Errorf("error uploading to bucket: %v", err)
	}

	// Set Content-Type, otherwise it will be saved as application/zip
	err = lib.SetGoogleStorageObjectContentType(path, "application/vnd.ms-excel")
	if err != nil {
		return fmt.Errorf("error setting content type for obj %s: %v", path, err)
	}
	return nil
}
