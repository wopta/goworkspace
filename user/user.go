package user

import (
	"net/http"
	"os"
	"time"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
)

var userRoutes []lib.Route = []lib.Route{
	{
		Route:   "/fiscalCode/v1/{fiscalcode}",
		Handler: lib.ResponseLoggerWrapper(GetUserByFiscalCodeFx),
		Method:  http.MethodGet,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/mail/v1/{mail}",
		Handler: lib.ResponseLoggerWrapper(GetUserByMailFx),
		Method:  http.MethodGet,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/authId/v1/{authId}",
		Handler: lib.ResponseLoggerWrapper(GetUserByAuthIdFx),
		Method:  http.MethodGet,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/onboarding/v1",
		Handler: lib.ResponseLoggerWrapper(OnboardUserFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/document/v1/{policyUid}",
		Handler: lib.ResponseLoggerWrapper(UploadDocumentFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/fiscalcode/v1/it/{operation}",
		Handler: lib.ResponseLoggerWrapper(FiscalCodeFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/fiscalcode/v2/check/{fiscalCode}",
		Handler: lib.ResponseLoggerWrapper(FiscalCodeCheckFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/invite/v1/create",
		Handler: lib.ResponseLoggerWrapper(CreateInviteFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAdmin},
	},
	{
		Route:   "/invite/v1/consume",
		Handler: lib.ResponseLoggerWrapper(ConsumeInviteFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/role/v1/{userUid}",
		Handler: lib.ResponseLoggerWrapper(UpdateUserRoleFx),
		Method:  http.MethodPatch,
		Roles:   []string{lib.UserRoleAdmin},
	},
	{
		Route:   "/v1",
		Handler: lib.ResponseLoggerWrapper(GetUsersFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAdmin},
	},
	{
		Route:   "/v2",
		Handler: lib.ResponseLoggerWrapper(GetUsersV2Fx),
		Method:  http.MethodGet,
		Roles:   []string{lib.UserRoleAdmin},
	},
}

func init() {
	log.Println("INIT User")
	functions.HTTP("User", User)
}

func User(w http.ResponseWriter, r *http.Request) {

	router := lib.GetRouter("user", userRoutes)
	router.ServeHTTP(w, r)
}

// Consider moving into policy, as User is a dependency of Policy
// and User does not need to know what a Policy is.
func SetUserIntoPolicyContractor(policy *models.Policy) {
	log.AddPrefix("SetUserIntoPolicyContractor")
	defer log.PopPrefix()
	userUID, newUser, err := models.GetUserUIDByFiscalCode(policy.Contractor.FiscalCode)
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
		policy.Contractor.UpdatedDate = policy.Contractor.CreationDate
		fireUsers := lib.UserCollection
		lib.SetFirestore(fireUsers, userUID, policy.Contractor)
		err = policy.Contractor.BigquerySave()
		if err != nil {
			log.ErrorF("error save user %s bigquery\n", policy.Contractor.Uid)
		}
		return
	}

	user := policy.Contractor.ToUser()
	if user == nil {
		return
	}

	_, err = models.UpdateUserByFiscalCode(*user)
	lib.CheckError(err)
}
