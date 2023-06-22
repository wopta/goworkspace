package user

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

const (
	invitesCollection = "invites"
	usersCollection   = "users"
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
				Roles:   []string{models.UserRoleAll},
			},
			{
				Route:   "/v1/mail/:mail",
				Handler: GetUserByMailFx,
				Method:  "GET",
				Roles:   []string{models.UserRoleAll},
			},
			{
				Route:   "/v1/authId/:authId",
				Handler: GetUserByAuthIdFx,
				Method:  "GET",
				Roles:   []string{models.UserRoleAll},
			},
			{
				Route:   "/v1/onboarding",
				Handler: OnboardUserFx,
				Method:  "POST",
				Roles:   []string{models.UserRoleAll},
			},
			{
				Route:   "/document/v1/:policyUid",
				Handler: UploadDocument,
				Method:  http.MethodPost,
				Roles:   []string{models.UserRoleAll},
			},
			{
				Route:   "/fiscalcode/v1/it/:operation",
				Handler: FiscalCode,
				Method:  http.MethodPost,
				Roles:   []string{models.UserRoleAll},
			},
			{
				Route:   "/invite/v1/create",
				Handler: CreateInviteFx,
				Method:  http.MethodPost,
				Roles:   []string{models.UserRoleAdmin},
			},
			{
				Route:   "/invite/v1/consume",
				Handler: ConsumeInviteFx,
				Method:  http.MethodPost,
				Roles:   []string{models.UserRoleAll},
			},
			{
				Route:   "/role/v1/:userUid",
				Handler: UpdateUserRoleFx,
				Method:  http.MethodPatch,
				Roles:   []string{models.UserRoleAdmin},
			},
			{
				Route:   "/user/v1",
				Handler: GetUsersFx,
				Method:  http.MethodPost,
				Roles:   []string{models.UserRoleAll},
			},
		},
	}
	route.Router(w, r)

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
	userFirebase := lib.WhereLimitFirestore(usersCollection, "fiscalCode", "==", fiscalCode, 1)
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
	userFirebase := lib.WhereLimitFirestore(usersCollection, "mail", "==", mail, 1)
	var user models.User
	user, err := models.FirestoreDocumentToUser(userFirebase)
	return user, err
}

// Consider moving into policy, as User is a dependency of Policy
// and User does not need to know what a Policy is.
func SetUserIntoPolicyContractor(policy *models.Policy, origin string) {
	userUID, newUser, err := models.GetUserUIDByFiscalCode(origin, policy.Contractor.FiscalCode)
	lib.CheckError(err)

	policy.Contractor.Uid = userUID
	log.Println("SetUserIntoPolicyContractor::Contractor UID: ", userUID)
	log.Println("SetUserIntoPolicyContractor::Policy Contractor UID: ", policy.Contractor.Uid)

	// Move user identity documents to user folder on Google Storage
	for _, identityDocument := range policy.Contractor.IdentityDocuments {
		frontMediaBytes, e := lib.GetFromGoogleStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"),
			"temp/"+policy.Uid+"/"+identityDocument.FrontMedia.FileName)
		lib.CheckError(e)
		frontGsLink, e := lib.PutToGoogleStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "assets/users/"+
			userUID+"/"+identityDocument.FrontMedia.FileName, frontMediaBytes)
		log.Println("SetUserIntoPolicyContractor::frontGsLink: ", frontGsLink)
		identityDocument.FrontMedia.Link = frontGsLink

		if identityDocument.BackMedia != nil {
			backMediaBytes, e := lib.GetFromGoogleStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"),
				"temp/"+policy.Uid+"/"+identityDocument.BackMedia.FileName)
			lib.CheckError(e)
			backGsLink, e := lib.PutToGoogleStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "assets/users/"+
				userUID+"/"+identityDocument.FrontMedia.FileName, backMediaBytes)
			log.Println("SetUserIntoPolicyContractor::backGsLink: ", backGsLink)
			identityDocument.BackMedia.Link = backGsLink
		}
	}

	if newUser {
		policy.Contractor.CreationDate = time.Now().UTC()
		fireUsers := lib.GetDatasetByEnv(origin, "users")
		lib.SetFirestore(fireUsers, userUID, policy.Contractor)
	} else {
		_, err = models.UpdateUserByFiscalCode(origin, policy.Contractor)
		lib.CheckError(err)
	}
}

func GetAuthUserByMail(mail string) (models.User, error) {
	var user models.User

	authId, err := lib.GetAuthUserIdByEmail(mail)
	if err != nil {
		return user, err
	}

	return GetUserByAuthId(authId)
}
