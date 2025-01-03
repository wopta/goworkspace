package inclusive

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/civil"
	"github.com/google/uuid"
	lib "github.com/wopta/goworkspace/lib"
)

// TO DO security,payload,error,fasature
func BankAccountScalapayFx(resp http.ResponseWriter, r *http.Request) (string, interface{}, error) {

	var (
		e   error
		obj BankAccountMovement
	)
	e = CheckScalapayApikey(r)
	if e != nil {
		return "", nil, e
	}
	obj, e = CheckScalapayData(r)
	if e != nil {
		return "", nil, e
	}
	obj = SetScalapayData(obj)
	e = lib.InsertRowsBigQuery("wopta", dataMovement, obj)
	if obj.MovementType == "insert" {

		e = lib.InsertRowsBigQuery("wopta", dataBanckAccount, obj)
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
func CheckScalapayApikey(r *http.Request) error {
	apikey := os.Getenv("HYPE_APIKEY")
	apikeyReq := r.Header.Get("api_key")
	if apikey != apikeyReq {
		return GetErrorJson(401, "Unauthorized", "")
	}
	return nil
}
func CheckScalapayData(r *http.Request) (BankAccountMovement, error) {
	req := lib.ErrorByte(io.ReadAll(r.Body))
	log.Println(r.Header)
	log.Println(string(req))
	var obj BankAccountMovement
	defer r.Body.Close()
	json.Unmarshal([]byte(req), &obj)
	if obj.Name == "" {
		return obj, GetErrorJson(400, "Bad request", "field name miss")
	}
	if obj.Tenant == "" {
		return obj, GetErrorJson(400, "Bad request", "field tenant miss")
	}
	if obj.Surname == "" {
		return obj, GetErrorJson(400, "Bad request", "field Surname miss")
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
		res, _ := QueryRowsBigQuery[BankAccountMovement]("wopta",
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

	if obj.Id == "" {
		return obj, GetErrorJson(400, "Bad request", "field Id miss")
	}
	if obj.GuaranteesCode != "base" && obj.GuaranteesCode != "premium" {
		return obj, GetErrorJson(400, "Bad request", "field GuaranteesCode out of enum")
	}

	return obj, nil
}
func SetScalapayData(obj BankAccountMovement) BankAccountMovement {

	obj.BigStartDate = civil.DateTimeOf(obj.StartDate)
	obj.BigEndDate = civil.DateTimeOf(obj.EndDate)
	if obj.GuaranteesCode == "base" {
		obj.PolicyNumber = "180623"
		obj.Uid = uuid.New().String()
		obj.Customer = "Scalapay"
		obj.Company = "axa"
		obj.PolicyType = ""
		obj.PolicyUid = ""
		obj.AssetType = ""
		obj.PolicyName = "Scalapay base"
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
