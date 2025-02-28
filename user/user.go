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

var userRoutes []lib.Route = []lib.Route{
	{
		Route:       "/fiscalCode/v1/{fiscalcode}",
		Fn:          GetUserByFiscalCodeFx,
		Method:      http.MethodGet,
		Roles:       []string{lib.UserRoleAll},
		Entitlement: "user.get.user.fiscalcode",
	},
	{
		Route:       "/mail/v1/{mail}",
		Fn:          GetUserByMailFx,
		Method:      http.MethodGet,
		Roles:       []string{lib.UserRoleAll},
		Entitlement: "user.get.user.mail",
	},
	{
		Route:       "/authId/v1/{authId}",
		Fn:          GetUserByAuthIdFx,
		Method:      http.MethodGet,
		Roles:       []string{lib.UserRoleAll},
		Entitlement: "user.get.user.authid",
	},
	{
		Route:       "/onboarding/v1",
		Fn:          OnboardUserFx,
		Method:      http.MethodPost,
		Roles:       []string{lib.UserRoleAll},
		Entitlement: "user.onboard",
	},
	{
		Route:       "/document/v1/{policyUid}",
		Fn:          UploadDocumentFx,
		Method:      http.MethodPost,
		Roles:       []string{lib.UserRoleAll},
		Entitlement: "user.upload.user.document",
	},
	{
		Route:       "/fiscalcode/v1/it/{operation}",
		Fn:          FiscalCodeFx,
		Method:      http.MethodPost,
		Roles:       []string{lib.UserRoleAll},
		Entitlement: "user.calculate.user.fiscalcode",
	},
	{
		Route:       "/invite/v1/create",
		Fn:          CreateInviteFx,
		Method:      http.MethodPost,
		Roles:       []string{lib.UserRoleAdmin},
		Entitlement: "user.invite.create",
	},
	{
		Route:       "/invite/v1/consume",
		Fn:          ConsumeInviteFx,
		Method:      http.MethodPost,
		Roles:       []string{lib.UserRoleAll},
		Entitlement: "user.invite.consume",
	},
	{
		Route:       "/role/v1/{userUid}",
		Fn:          UpdateUserRoleFx,
		Method:      http.MethodPatch,
		Roles:       []string{lib.UserRoleAdmin},
		Entitlement: "user.update.user.role",
	},
	{
		Route:       "/v1",
		Fn:          GetUsersFx,
		Method:      http.MethodPost,
		Roles:       []string{lib.UserRoleAdmin},
		Entitlement: "user.get.users",
	},
	{
		Route:       "/v2",
		Fn:          GetUsersV2Fx,
		Method:      http.MethodGet,
		Roles:       []string{lib.UserRoleAdmin},
		Entitlement: "user.get.users",
	},
}

func init() {
	log.Println("INIT User")
	functions.HTTP("User", User)
}

func User(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix)

	router := models.GetExtendedRouter("user", userRoutes)
	router.ServeHTTP(w, r)
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
		policy.Contractor.UpdatedDate = policy.Contractor.CreationDate
		fireUsers := lib.GetDatasetByEnv(origin, lib.UserCollection)
		lib.SetFirestore(fireUsers, userUID, policy.Contractor)
		err = policy.Contractor.BigquerySave(origin)
		if err != nil {
			log.Printf("[SetUserIntoPolicyContractor] error save user %s bigquery\n", policy.Contractor.Uid)
		}
		return
	}

	user := policy.Contractor.ToUser()
	if user == nil {
		return
	}

	_, err = models.UpdateUserByFiscalCode(origin, *user)
	lib.CheckError(err)
}
