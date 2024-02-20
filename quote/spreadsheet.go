package quote

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/wopta/goworkspace/lib"
	"google.golang.org/api/drive/v2"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type QuoteSpreadsheet struct {
	SheetName, filename string
	Id                  string
	InputCells          []InputCell
}

func SpreadsheetsFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	qs := QuoteSpreadsheet{}
	qs.Spreadsheets()
	return "", nil, nil
}

func (qs *QuoteSpreadsheet) Spreadsheets() {
	var (
		path []byte
		file *drive.File
	)
	switch os.Getenv("env") {
	case "local":
		path = lib.ErrorByte(ioutil.ReadFile("function-data/sa/positive-apex-350507-33284d6fdd55.json"))
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
	googleDrive := &GoogleDrive{
		CredentialsByte: path,
		Ctx:             ctx,
	}

	driveClient, e := GoogleClient[*DriveService](googleDrive)
	lib.CheckError(e)
	fmt.Printf("driveClient: %v\n", driveClient)
	sheetClient, e := GoogleClient[*sheets.Service](spreadsheet)
	lib.CheckError(e)
	fmt.Printf("sheetClient: %v\n", sheetClient)
	f, e := driveClient.Svc.Files.Copy(qs.Id, file).Do()
	sheet, e := sheetClient.Spreadsheets.Values.Get(f.Id, "A:J").Do()
	fmt.Printf("file: %v\n", sheet.Values[99][3])
	sheet.Values[40][2] = "10000000"
	fmt.Printf("file: %v\n", sheet.Values[99][3])
	

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
	svc, err = drive.NewService(s.Ctx, option.WithCredentialsJSON(s.CredentialsByte), option.WithScopes(drive.DriveScriptsScope))
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
