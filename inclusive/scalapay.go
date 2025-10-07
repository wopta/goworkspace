package inclusive

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"cloud.google.com/go/civil"
	"github.com/google/uuid"
	lib "gitlab.dev.wopta.it/goworkspace/lib"
)

const (
	dataset       = "wopta_inclusive"
	movementTable = "bank_account_movement_scalapay"
	usersTable    = "bank_account_users_scalapay"
	layout        = "2006-01-02"
	layoutQuery   = "2006-01-02"
)

// TO DO security,payload,error,fasature
func BankAccountScalapayFx(resp http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.AddPrefix("BankAccountScalapayFx ")
	defer log.PopPrefix()
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

	if obj.MovementType == "insert" {
		e = lib.InsertRowsBigQuery(dataset, movementTable, obj)
		log.Println(e)
		res := QueryRowsBigQuery[BankAccountMovement](dataset,
			usersTable,
			"select * from `"+dataset+"."+usersTable+"` where guaranteesCode ='"+obj.GuaranteesCode+"' and id ='"+obj.Id+"'")
		log.Println(len(res))
		if len(res) == 0 {
			e = lib.InsertRowsBigQuery(dataset, usersTable, obj)
		} else {
			e = lib.UpdateRowBigQuery(dataset, usersTable, map[string]string{
				"status":    obj.Status,
				"startDate": obj.StartDate.Format(layout) + " 00:00:00",
			}, " id='"+obj.Id+"' and guaranteesCode='"+obj.GuaranteesCode+"'")
		}

	}

	if obj.MovementType == "delete" || obj.MovementType == "suspended" {

		refDay := time.Now()
		res := QueryRowsBigQuery[BankAccountMovement](dataset,
			usersTable,
			"select * from `"+dataset+"."+movementTable+"` where guaranteesCode ='"+obj.GuaranteesCode+"' and id ='"+obj.Id+"' and StartDate ='"+refDay.Format(layoutQuery)+"' and movementType ='insert'")
		log.Println(len(res))

		if len(res) == 0 {
			e = lib.InsertRowsBigQuery(dataset, movementTable, obj)
			log.Println(e)

			e = lib.UpdateRowBigQuery(dataset, usersTable, map[string]string{
				"status":  obj.Status,
				"endDate": obj.EndDate.Format(layout) + " 00:00:00",
			}, "id ='"+obj.Id+"' and guaranteesCode='"+obj.GuaranteesCode+"'")
		} else {
			log.Println("400, Bad request field Movement insert same day for")
			return "", nil, GetErrorJson(400, "Bad request", "field Movement insert same day for "+obj.Id)
		}

	}

	if e != nil {
		log.Error(e)
		return "", nil, GetErrorJson(500, "internal server error", "")
	} else {
		log.Println(`200 {"woptaUid":"` + obj.Uid + `"}`)
	}
	return `{"woptaUid":"` + obj.Uid + `"}`, nil, e
}
func CheckScalapayApikey(r *http.Request) error {
	apikey := os.Getenv("HYPE_APIKEY")
	apikeyReq := r.Header.Get("api_key")
	log.Println("apikeyReq: ", apikeyReq, " apikey: ", apikey)
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
		res := QueryRowsBigQuery[BankAccountMovement](dataset,
			usersTable,
			"select * from `"+dataset+"."+movementTable+"` where guaranteesCode ='"+obj.GuaranteesCode+"' and id ='"+obj.Id+"' and movementType ='insert'")
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
		obj.PolicyNumber = "051124"
		obj.Uid = uuid.New().String()
		obj.Customer = "Scalapay"
		obj.Company = "axa"
		obj.PolicyType = ""
		obj.PolicyUid = ""
		obj.AssetType = ""
		obj.PolicyName = "Scalapay base"
	}

	obj.CustomerId = obj.Id
	if obj.MovementType == "insert" {
		obj.Status = "active"
		obj.Daystart = strconv.Itoa(obj.StartDate.Day())
		if obj.StartDate.Day() == 29 && int(obj.StartDate.Month()) == 2 {
			obj.Daystart = "28"
		}
	}
	if obj.MovementType == "delete" {
		obj.Status = "delete"

	}
	if obj.MovementType == "suspended" {
		obj.Status = "suspended"

	}

	return obj
}
