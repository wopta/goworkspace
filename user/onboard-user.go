package user

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	lib "github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func OnboardUserFx(resp http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		onboardUserRequest OnboardUserDto
		user               *models.User
	)
	resp.Header().Set("Access-Control-Allow-Methods", "POST")

	reqBytes := lib.ErrorByte(io.ReadAll(r.Body))
	json.Unmarshal(reqBytes, &onboardUserRequest)

	origin := r.Header.Get("Origin")
	fireUser := lib.GetDatasetByEnv(origin, models.UserCollection)

	canRegister, user, userId, email := CanUserRegisterUseCase(onboardUserRequest.FiscalCode)

	if !canRegister {
		fmt.Printf("User with fiscalCode %s cannot register", onboardUserRequest.FiscalCode)
		return `{"success": false}`, `{"success": false}`, nil
	}

	if email != nil && *email != onboardUserRequest.Email {
		fmt.Printf("User with fiscalCode %s cannot register: emails doesn't match", onboardUserRequest.FiscalCode)
		return `{"success": false}`, `{"success": false}`, nil
	}

	dbUser, e := lib.CreateUserWithEmailAndPassword(onboardUserRequest.Email, onboardUserRequest.Password, userId)
	if e != nil {
		return `{"success": false}`, `{"success": false}`, nil
	}

	if userId != nil {
		fmt.Printf("User with fiscalCode %s is being updated", onboardUserRequest.FiscalCode)
		lib.UpdateFirestoreErr(fireUser, dbUser.UID, map[string]interface{}{"authId": dbUser.UID,
			"role": models.UserRoleCustomer})
	} else {
		fmt.Printf("User with fiscalCode %s is being created", onboardUserRequest.FiscalCode)
		user.Uid = dbUser.UID
		user.AuthId = dbUser.UID
		user.Role = models.UserRoleCustomer
		lib.SetFirestore(fireUser, dbUser.UID, user)
	}

	err := user.BigquerySave(origin)
	if err != nil {
		log.Printf("[OnBoardUser] error save user %s bigquery: %s", user.Uid, err.Error())
	}

	// update the user custom claim
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
