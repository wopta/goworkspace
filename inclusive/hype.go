package inclusive

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"net/http"
	"os"
	"time"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/civil"
	"github.com/google/uuid"
	lib "gitlab.dev.wopta.it/goworkspace/lib"
)

const (
	dataMovement     = "inclusive_bank_account_movement"
	dataBanckAccount = "inclusive_bank_account"
	suspended        = "suspended"
	insert           = "insert"
	delete           = "delete"
	active           = "active"
)

// TO DO security,payload,error,fasature
func BankAccountHypeFx(resp http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.AddPrefix("BankAccountHypeFx ")
	defer log.PopPrefix()
	var (
		e   error
		obj BankAccountMovement
	)
	e = CheckHypeApikey(r)
	if e != nil {
		return "", nil, e
	}
	obj, e = CheckHypeData(r)
	if e != nil {
		return "", nil, e
	}
	obj = SetHypeData(obj)
	e = lib.InsertRowsBigQuery("wopta", dataMovement, obj)
	if obj.MovementType == "insert" {

		//e = lib.InsertRowsBigQuery("wopta", dataBanckAccount, obj)
	}
	/*
		layout := "2006-01-02"
		if obj.MovementType == "delete" || obj.MovementType == "suspended" {
			e = lib.UpdateRowBigQuery("wopta", dataBanckAccount, map[string]string{
				"status":  obj.Status,
				"endDate": obj.EndDate.Format(layout) + " 00:00:00",
			}, "fiscalCode='"+obj.FiscalCode+"' and guaranteesCode='"+obj.GuaranteesCode+"'")

		}
	*/
	return `{"woptaUid":"` + obj.Uid + `"}`, nil, e
}
func CountHypeFx(resp http.ResponseWriter, r *http.Request) (string, interface{}, error) {

	HypeCount("", "", "")
	return ``, nil, nil
}
func HypeImportMovementbankAccountFx(resp http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println("---------------HypeImportMovementbankAccountFx -------------------------------")
	HypeImportMovementbankAccount()
	return ``, nil, nil
}

func CheckHypeApikey(r *http.Request) error {
	apikey := os.Getenv("HYPE_APIKEY")
	apikeyReq := r.Header.Get("api_key")
	if apikey != apikeyReq {
		return GetErrorJson(401, "Unauthorized", "")
	}
	return nil
}
func CheckHypeData(r *http.Request) (BankAccountMovement, error) {
	req := lib.ErrorByte(io.ReadAll(r.Body))
	log.Println(r.Header)
	log.Println(string(req))
	var obj BankAccountMovement
	defer r.Body.Close()
	json.Unmarshal([]byte(req), &obj)
	if obj.Name == "" {
		return obj, GetErrorJson(400, "Bad request", "field name miss")
	}
	if obj.FiscalCode == "" {
		return obj, GetErrorJson(400, "Bad request", "field name miss")
	}
	if obj.Surname == "" {
		return obj, GetErrorJson(400, "Bad request", "field name miss")
	}
	if obj.MovementType != "insert" && obj.MovementType != "delete" && obj.MovementType != "suspended" {
		return obj, GetErrorJson(400, "Bad request", "field MovementType out of enum")
	}
	if obj.MovementType == "insert" {
		if obj.StartDate.IsZero() {
			return obj, GetErrorJson(400, "Bad request", "field StartDate miss")
		}
	}
	if obj.MovementType == "delete" || obj.MovementType == "suspended" {
		res := QueryRowsBigQuery[BankAccountMovement]("wopta",
			"inclusive_axa_bank_account",
			"select * from `wopta."+dataMovement+"` where fiscalCode='"+obj.FiscalCode+"' and guaranteesCode ='"+obj.GuaranteesCode+"'")
		log.Println(len(res))
		if len(res) == 0 {
			return obj, GetErrorJson(400, "Bad request", "insert movement miss")
		}
		if obj.StartDate.IsZero() {

			return obj, GetErrorJson(400, "Bad request", "field StartDate miss")
		}
		if obj.EndDate.IsZero() {
			return obj, GetErrorJson(400, "Bad request", "field EndDate miss")
		}
	}

	if obj.HypeId == "" {
		return obj, GetErrorJson(400, "Bad request", "field HypeId miss")
	}
	if obj.GuaranteesCode != "next" && obj.GuaranteesCode != "premium" && obj.GuaranteesCode != "QUICK2CASH" {
		return obj, GetErrorJson(400, "Bad request", "field GuaranteesCode out of enum")
	}

	return obj, nil
}
func SetHypeData(obj BankAccountMovement) BankAccountMovement {

	obj.BigStartDate = civil.DateTimeOf(obj.StartDate)
	obj.BigEndDate = civil.DateTimeOf(obj.EndDate)
	obj.Customer = "hype"
	obj.Company = "axa"
	obj.PolicyType = ""
	obj.PolicyUid = ""
	obj.AssetType = ""
	obj.PostalCode = ""
	obj.CityCode = ""
	obj.Address = ""
	obj.Tenant = "hype"
	obj.City = ""
	if obj.GuaranteesCode == "next" {
		obj.PolicyNumber = "180623"
		obj.Uid = uuid.New().String()
		obj.Customer = "hype"
		obj.PolicyName = "Hype Next"
	}
	if obj.GuaranteesCode == "premium" {
		obj.PolicyNumber = "191123"
		obj.PolicyName = "Hype Premium"

	}
	//TODO:CHANGE NUMBER
	if obj.GuaranteesCode == "lottomatica" {
		obj.PolicyNumber = "270125"
		obj.PolicyName = "Hype QUICK2CASH"
	}
	obj.CustomerId = obj.HypeId
	if obj.MovementType == "insert" {

		obj.Status = "active"
	}
	if obj.MovementType == "delete" {
		obj.Status = "delete"

	}
	if obj.MovementType == "suspended" {
		obj.Status = "suspended"

	}

	return obj
}

func HypeCount(date string, fiscalCode string, guaranteesCode string) {
	var (
		countResponseModel CountResponseModel
	)

	refday := time.Now().AddDate(0, 0, -1)
	refdayString := refday.Format("2006-01-02")
	stringquery := "with Mov AS(SELECT distinct fiscalCode,* from `wopta." + dataMovement + "` where fiscalCode='" + fiscalCode + "' and guaranteesCode ='" + guaranteesCode + "and _PARTITIONTIME ='" + refdayString + "' SELECT Mov.movementType ,count(*)as count FROM Mov group by Mov.movementType"
	log.Println(len(stringquery))
	queryWopta := QueryRowsBigQuery[bigquery.Value]("wopta", "inclusive_axa_bank_account", stringquery)
	log.Println(len(queryWopta))
	for _, mov := range queryWopta {
		log.Println(mov)
	}
	requestUrl := os.Getenv("HYPE_PLATHFORM_PATH") + "/external/wopta/v1/next/amount/" + refdayString + "/" + refdayString

	req, err := http.NewRequest(http.MethodGet, requestUrl, nil)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("ext-wopta-service-key", os.Getenv("HYPE_APIKEY_OUT"))
	res := lib.Httpclient(req)
	reqAll := lib.ErrorByte(io.ReadAll(res.Body))
	json.Unmarshal(reqAll, &countResponseModel)
	log.Println(res)

}

func HypeReconciliation(date string, fiscalCode string, guaranteesCode string) {

}

/*
		name,surname,fiscalCode,hypeId,guaranteesCode,startDate
	    Luca,Barbieri,BRBLCU81H803F205Q,123789,next,2023-07-15
*/
func HypeImportMovementbankAccount() {
	log.Println("---------------HypeImportMovementbankAccount -------------------------------")

	data := lib.GetFromStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "track/in/inclusive/bank-account/hype/20240227_esportazione_wopta_premium.csv", "")
	df := lib.CsvToDataframe(data)
	log.Println("HypeImportMovementbankAccount  row", df.Nrow())
	log.Println("HypeImportMovementbankAccount  col", df.Ncol())
	var result [][]string
	var movList []BankAccountMovement
	count := 0
	for i, d := range df.Records() {

		log.Println("HypeImportMovementbankAccount  num ", i)
		if i > 0 {
			uid := uuid.New().String()
			start := time.Now()

			mov := BankAccountMovement{
				Uid:            uid,
				Status:         active,
				Name:           d[0],
				Surname:        d[1],
				FiscalCode:     d[2],
				GuaranteesCode: "premium",
				HypeId:         d[3],
				BigStartDate:   civil.DateTimeOf(start),
				BigEndDate:     civil.DateTimeOf(start),
				PolicyNumber:   "191123",
				Customer:       "hype",
				Company:        "axa",
				PolicyName:     "Hype Premium",
			}
			result = append(result, []string{d[0], d[1], d[2], d[3], d[4], uid})
			movList = append(movList, mov)
			count++
			if count == 500 {
				count = 0
				e := lib.InsertRowsBigQuery("wopta", dataMovement, movList)
				e = lib.InsertRowsBigQuery("wopta", dataBanckAccount, movList)
				log.Println("HypeImportMovementbankAccount error InsertRowsBigQuery: ", e)
				movList = []BankAccountMovement{}
			}
		}
	}
	e := lib.InsertRowsBigQuery("wopta", dataMovement, movList)
	e = lib.InsertRowsBigQuery("wopta", dataBanckAccount, movList)
	log.Println("HypeImportMovementbankAccount error InsertRowsBigQuery: ", e)
	filepath := "result_02_premium.csv"
	lib.WriteCsv("../tmp/"+filepath, result, ',')
	source, _ := ioutil.ReadFile("../tmp/" + filepath)
	lib.PutToStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "track/in/inclusive/bank-account/hype/"+filepath, source)

}
