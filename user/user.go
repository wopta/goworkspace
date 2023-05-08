package user

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	lib "github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func init() {
	log.Println("INIT User")
	functions.HTTP("User", User)
}

func User(w http.ResponseWriter, r *http.Request) {

	log.Println("User")
	lib.EnableCors(&w, r)
	route := lib.RouteData{
		Routes: []lib.Route{
			{
				Route:   "/v1/fiscalCode/:fiscalcode",
				Handler: GetUserByFiscalCodeFx,
				Method:  "GET",
			},
			{
				Route:   "/v1/mail/:mail",
				Handler: GetUserByMailFx,
				Method:  "GET",
			},
			{
				Route:   "/v1/authId/:authId",
				Handler: GetUserByAuthIdFx,
				Method:  "GET",
			},
			{
				Route:   "/v1/onboarding",
				Handler: OnboardUserFx,
				Method:  "POST",
			},
			{
				Route:   "/v1/document",
				Handler: UploadDocument,
				Method:  http.MethodPost,
			},
		},
	}
	route.Router(w, r)

}

func OnboardUserFx(resp http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		onboardUserRequest OnboardUserDto
		result             string
	)
	resp.Header().Set("Access-Control-Allow-Methods", "POST")

	reqBytes := lib.ErrorByte(ioutil.ReadAll(r.Body))
	json.Unmarshal(reqBytes, &onboardUserRequest)

	canRegister, userId, email := CanUserRegisterUseCase(onboardUserRequest.FiscalCode)
	log.Println(email)

	if canRegister {
		fireUser, e := lib.CreateUserWithEmailAndPassword(onboardUserRequest.Email, onboardUserRequest.Password, userId)
		if e != nil {
			result = `{"success": false}`
		} else {
			lib.UpdateFirestoreErr("users", *userId, map[string]interface{}{"authId": fireUser.UID})
			result = `{"success": true}`
		}
	} else {
		result = `{"success": false}`
	}

	return result, result, nil
}

func GetUserByAuthIdFx(resp http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	resp.Header().Set("Access-Control-Allow-Methods", "GET")
	log.Println(r.Header.Get("authId"))
	user, e := GetUserByAuthId(r.Header.Get("authId"))
	jsonString, e := user.Marshal()
	return string(jsonString), user, e
}

func GetUserByAuthId(authId string) (models.User, error) {
	log.Println(authId)
	userFirebase := lib.WhereLimitFirestore("users", "authId", "==", authId, 1)
	var user models.User
	user, err := models.FirestoreDocumentToUser(userFirebase)
	return user, err
}

func GetUserByFiscalCodeFx(resp http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	resp.Header().Set("Access-Control-Allow-Methods", "GET")
	log.Println(r.Header.Get("fiscalCode"))
	p, e := GetUserByFiscalCode(r.Header.Get("fiscalCode"))
	jsonString, e := p.Marshal()
	return string(jsonString), p, e
}

func GetUserByFiscalCode(fiscalCode string) (models.User, error) {
	log.Println(fiscalCode)
	userFirebase := lib.WhereLimitFirestore("users", "fiscalCode", "==", fiscalCode, 1)
	var user models.User
	user, err := models.FirestoreDocumentToUser(userFirebase)
	return user, err
}

func GetUserByMailFx(resp http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	resp.Header().Set("Access-Control-Allow-Methods", "GET")
	log.Println(r.Header.Get("mail"))
	p, e := GetUserByMail(r.Header.Get("mail"))
	jsonString, e := p.Marshal()
	return string(jsonString), p, e
}

func GetUserByMail(mail string) (models.User, error) {
	log.Println(mail)
	userFirebase := lib.WhereLimitFirestore("users", "mail", "==", mail, 1)
	var user models.User
	user, err := models.FirestoreDocumentToUser(userFirebase)
	return user, err
}

type OnboardUserDto struct {
	FiscalCode string `json:"fiscalCode"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}
