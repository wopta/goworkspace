package user

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

	reqBytes := lib.ErrorByte(ioutil.ReadAll(r.Body))
	json.Unmarshal(reqBytes, &onboardUserRequest)

	canRegister, user, userId, email := CanUserRegisterUseCase(onboardUserRequest.FiscalCode)

	if !canRegister {
		fmt.Printf("User with fiscalCode %s cannot register", onboardUserRequest.FiscalCode)
		return `{"success": false}`, `{"success": false}`, nil
	}

	if email != nil && *email != onboardUserRequest.Email {
		fmt.Printf("User with fiscalCode %s cannot register: emails doesn't match", onboardUserRequest.FiscalCode)
		return `{"success": false}`, `{"success": false}`, nil
	}

	fireUser, e := lib.CreateUserWithEmailAndPassword(onboardUserRequest.Email, onboardUserRequest.Password, userId)
	if e != nil {
		return `{"success": false}`, `{"success": false}`, nil
	}

	if userId != nil {
		fmt.Printf("User with fiscalCode %s is being updated", onboardUserRequest.FiscalCode)
		lib.UpdateFirestoreErr("users", fireUser.UID, map[string]interface{}{"authId": fireUser.UID})
	} else {
		fmt.Printf("User with fiscalCode %s is being created", onboardUserRequest.FiscalCode)
		user.Uid = fireUser.UID
		user.AuthId = fireUser.UID
		lib.SetFirestore("users", fireUser.UID, user)
	}

	return `{"success": true}`, `{"success": true}`, nil
}

type OnboardUserDto struct {
	FiscalCode string `json:"fiscalCode"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}