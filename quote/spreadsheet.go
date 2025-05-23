package quote

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
	"google.golang.org/api/drive/v2"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type QuoteSpreadsheet struct {
	SheetName          string
	ExportedSheetName  string
	Id                 string
	DestinationSheetId string
	ExportFilePrefix   string
	InputCells         []Cell
	OutputCells        []Cell
	InitCells          []Cell
}

func (qs *QuoteSpreadsheet) Spreadsheets() ([]Cell, string) {
	var (
		path           []byte
		bucketSavePath = "test/download/"
	)

	path = lib.GetFilesByEnv("sa/positive-apex-350507-33284d6fdd55.json")

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

	err := clearUnwantedSheetsAndCopyToSpreadsheet(sheetClient, qs, ctx)
	if err != nil {
		log.Printf("unable to perform sheet operations: %v", err)
		return res, ""
	}

	doc, err := loadFromDrive(path, ctx, qs.DestinationSheetId)
	if err != nil {
		log.Printf("unable to load from GDrive: %v", err)
		return res, ""
	}
	gsLink, err := saveToBucket(fmt.Sprintf("%s%s_%s.xlsx", bucketSavePath, qs.ExportFilePrefix, time.Now().Format("2006-1-2_15:04:05")), doc)
	if err != nil {
		log.Printf("unable to save to bucket: %v", err)
		return res, ""
	}

	return res, gsLink
}

func clearUnwantedSheetsAndCopyToSpreadsheet(sheetClient *sheets.Service, qs *QuoteSpreadsheet, ctx context.Context) error {
	var originalSheetId, exportSheetId int64

	ssRes, _ := sheetClient.Spreadsheets.Get(qs.Id).Context(ctx).Do()
	for _, s := range ssRes.Sheets {
		if s.Properties.Title == qs.SheetName {
			originalSheetId = s.Properties.SheetId
		}
		if s.Properties.Title == qs.ExportedSheetName {
			exportSheetId = s.Properties.SheetId
		}
	}

	qs.export(sheetClient, ctx, originalSheetId, exportSheetId)

	clearUnwantedSheetsReq := make([]*sheets.Request, 0)
	clearLastSheetReq := make([]*sheets.Request, 0)
	ssRes, _ = sheetClient.Spreadsheets.Get(qs.DestinationSheetId).Context(ctx).Do()
	for i, s := range ssRes.Sheets {
		ds := sheets.DeleteSheetRequest{SheetId: s.Properties.SheetId}
		sr := sheets.Request{DeleteSheet: &ds}
		if i == 0 {
			clearLastSheetReq = append(clearLastSheetReq, &sr)
		} else {
			clearUnwantedSheetsReq = append(clearUnwantedSheetsReq, &sr)
		}
	}

	if len(clearUnwantedSheetsReq) != 0 {
		_, err := sheetClient.Spreadsheets.BatchUpdate(qs.DestinationSheetId, &sheets.BatchUpdateSpreadsheetRequest{Requests: clearUnwantedSheetsReq}).Context(ctx).Do()
		if err != nil {
			log.Printf("unable to delete sheets from spreadsheet: %v", err)
		}
	}

	_, err := sheetClient.Spreadsheets.Sheets.CopyTo(qs.Id, exportSheetId, &sheets.CopySheetToAnotherSpreadsheetRequest{
		DestinationSpreadsheetId: qs.DestinationSheetId,
	}).Context(ctx).Do()

	_, err = sheetClient.Spreadsheets.BatchUpdate(qs.DestinationSheetId, &sheets.BatchUpdateSpreadsheetRequest{Requests: clearLastSheetReq}).Context(ctx).Do()
	if err != nil {
		log.Printf("unable to delete sheets from spreadsheet: %v", err)
	}

	ssRes, _ = sheetClient.Spreadsheets.Get(qs.DestinationSheetId).Context(ctx).Do()
	sheetIdToRename := ssRes.Sheets[0].Properties.SheetId
	renameSheetReq := make([]*sheets.Request, 0)
	newSp := sheets.SheetProperties{SheetId: sheetIdToRename, Title: qs.ExportedSheetName}
	rs := sheets.UpdateSheetPropertiesRequest{Properties: &newSp, Fields: "Title"}
	dr := sheets.Request{UpdateSheetProperties: &rs}
	renameSheetReq = append(renameSheetReq, &dr)

	_, err = sheetClient.Spreadsheets.BatchUpdate(qs.DestinationSheetId, &sheets.BatchUpdateSpreadsheetRequest{Requests: renameSheetReq}).Context(ctx).Do()
	if err != nil {
		log.Printf("unable to rename sheet: %v", err)
	}

	return nil
}

func (qs *QuoteSpreadsheet) export(sheetClient *sheets.Service, ctx context.Context, originalSheetId, exportSheetId int64) {
	requests := make([]*sheets.Request, 0)
	requests = append(requests, &sheets.Request{
		CopyPaste: &sheets.CopyPasteRequest{
			Source: &sheets.GridRange{
				SheetId:          originalSheetId,
				StartColumnIndex: 0,
				EndColumnIndex:   100,
				StartRowIndex:    0,
				EndRowIndex:      300,
			},
			Destination: &sheets.GridRange{
				SheetId:          exportSheetId,
				StartColumnIndex: 0,
				EndColumnIndex:   100,
				StartRowIndex:    0,
				EndRowIndex:      300,
			},
			PasteType: "PASTE_VALUES",
		},
	})
	if _, err := sheetClient.Spreadsheets.BatchUpdate(qs.Id, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: requests,
	}).Context(ctx).Do(); err != nil {
		log.Error(err)
		panic(err)
	}

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
	col := map[string]int{"A": 0, "B": 1, "C": 2, "D": 3, "E": 4, "F": 5, "G": 6, "H": 7, "I": 8, "J": 9, "K": 10,
		"L": 11, "M": 12, "N": 13, "O": 14, "P": 15, "Q": 16, "R": 17, "S": 18}
	sheet, e := sheetClient.Spreadsheets.Values.Get(qs.Id, qs.SheetName+"!A:S").Do()
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

func saveToBucket(path string, file []byte) (string, error) {
	gsLink, err := lib.PutToGoogleStorageWithSpecificContentType(os.Getenv("GOOGLE_STORAGE_BUCKET"), path, file, "application/vnd.ms-excel")
	if err != nil {
		return "", fmt.Errorf("error uploading to bucket: %v", err)
	}

	return gsLink, nil
}
