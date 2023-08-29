package user

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func OnboardUserFx(resp http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println("[OnboardUserFx] Handler start -------------------------------")

	var (
		onboardUserRequest OnboardUserDto
		user               *models.User
	)
	resp.Header().Set("Access-Control-Allow-Methods", "POST")

	reqBytes := lib.ErrorByte(io.ReadAll(r.Body))
	json.Unmarshal(reqBytes, &onboardUserRequest)
	log.Printf("[OnboardUserFx] Request email: '%s'", onboardUserRequest.Email)
	log.Printf("[OnboardUserFx] Request fiscalCode: %s", onboardUserRequest.FiscalCode)

	origin := r.Header.Get("Origin")
	fireUser := lib.GetDatasetByEnv(origin, models.UserCollection)

	canRegister, user, userId, email := CanUserRegisterUseCase(onboardUserRequest.FiscalCode)

	if !canRegister {
		log.Printf("User with fiscalCode %s cannot register", onboardUserRequest.FiscalCode)
		return `{"success": false}`, `{"success": false}`, nil
	}

	if email == nil {
		log.Printf("User with fiscalCode %s cannot register - email not found", onboardUserRequest.FiscalCode)
		return `{"success": false}`, `{"success": false}`, nil
	}

	dbEmailNormalized := strings.ToLower(strings.TrimSpace(*email))
	requestEmailNormalized := strings.ToLower(strings.TrimSpace(onboardUserRequest.Email))
	areEmailsEqual := strings.EqualFold(dbEmailNormalized, requestEmailNormalized)
	log.Printf("[OnboardUserFx] request email '%s' - db email '%s'", onboardUserRequest.Email, *email)
	log.Printf("[OnboardUserFx] normalized: request email '%s' - db email '%s' - equal %t", requestEmailNormalized, dbEmailNormalized, areEmailsEqual)

	if !areEmailsEqual {
		log.Printf("User with fiscalCode %s cannot register: emails doesn't match", onboardUserRequest.FiscalCode)
		return `{"success": false}`, `{"success": false}`, nil
	}

	dbUser, e := lib.CreateUserWithEmailAndPassword(requestEmailNormalized, onboardUserRequest.Password, userId)
	if e != nil {
		log.Printf("[OnboardUserFx] error creating auth user: %s", e.Error())
		return `{"success": false}`, `{"success": false}`, nil
	}

	if userId != nil {
		log.Printf("[OnboardUserFx] User with fiscalCode %s is being updated", onboardUserRequest.FiscalCode)
		e := lib.UpdateFirestoreErr(fireUser, dbUser.UID, map[string]interface{}{"authId": dbUser.UID,
			"role": models.UserRoleCustomer})
		if e != nil {
			log.Printf("[OnboardUserFx] error updating user: %s", e.Error())
		}
	} else {
		log.Printf("[OnboardUserFx] User with fiscalCode %s is being created", onboardUserRequest.FiscalCode)
		user.Uid = dbUser.UID
		user.AuthId = dbUser.UID
		user.Role = models.UserRoleCustomer
		lib.SetFirestore(fireUser, dbUser.UID, user)
	}

	err := user.BigquerySave(origin)
	if err != nil {
		log.Printf("[OnBoardUser] error save user %s bigquery: %s", user.Uid, err.Error())
	}

	log.Println("[OnboardUserFx] updating claims for user")
	lib.SetCustomClaimForUser(dbUser.UID, map[string]interface{}{
		"role": models.UserRoleCustomer,
	})

	return `{"success": true}`, `{"success": true}`, nil
}

type OnboardUserDto struct {
	FiscalCode string `json:"fiscalCode"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}
