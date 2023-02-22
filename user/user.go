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
	log.Println("INIT AppcheckProxy")
	functions.HTTP("User", User)
}

func User(w http.ResponseWriter, r *http.Request) {

	log.Println("Product")
	lib.EnableCors(&w, r)
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	route := lib.RouteData{
		Routes: []lib.Route{
			{
				Route:   "/v1/fiscalcode/:fiscalcode",
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
				Method:  "POST",
			},

			{
				Route:   "/v1/onboarding",
				Handler: OnboardUserFx,
				Method:  "POST",
			},

			{
				Route:   "/v1/login",
				Handler: PutFx,
				Method:  "POST",
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
	reqBytes := lib.ErrorByte(ioutil.ReadAll(r.Body))
	json.Unmarshal(reqBytes, &onboardUserRequest)

	canRegister, userId := CanUserRegisterUseCase(onboardUserRequest.FiscalCode)

	if canRegister {
		_, e := lib.CreateUserWithEmailAndPassword(onboardUserRequest.Email, onboardUserRequest.Password, userId)
		if e != nil {
			result = `{"success": false}`
		} else {
			result = `{"success": true}`
		}
	} else {
		result = `{"success": false}`
	}

	return result, result, nil
}

func GetName(name string) (models.Product, error) {

	productFire := lib.WhereFirestore("pruducts", "name", "==", name)

	products := models.ProductToListData(productFire)

	return products[0], nil

}
func PutFx(resp http.ResponseWriter, r *http.Request) (string, interface{}, error) {

	request := lib.ErrorByte(ioutil.ReadAll(r.Body))
	pr, e := models.UnmarshalProduct([]byte(request))
	p, e := Put(pr)
	return "{}", p, e
}
func Put(p models.Product) (models.Product, error) {

	r, _, e := lib.PutFirestoreErr("pruducts", p)
	log.Println(r.ID)

	return p, e
}


func GetUserByAuthIdFx(resp http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println(r.Header.Get("uid"))
	p, e := GetUserByAuthId(r.Header.Get("uid"))
	jsonString, e := p.Marshal()
	return string(jsonString), p, e
}

func GetUserByAuthId(uid string) (models.User, error) {
	log.Println(uid)
	userFirebase := lib.WhereLimitFirestore("users", "authId", "==", uid, 1)
	var user models.User
	user, err := models.FirestoreDocumentToUser(userFirebase)
	return user, err
}

func GetUserByFiscalCodeFx(resp http.ResponseWriter, r *http.Request) (string, interface{}, error) {
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
