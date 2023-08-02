package user

import (
	"encoding/json"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"io"
	"log"
	"net/http"
	"strings"
)

type updateUserReq struct {
	Role string `json:"role"`
}

func UpdateUserRoleFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err               error
		userUid, userRole string
		user              models.User
		request           updateUserReq
	)

	log.Println("Update User Role")

	defer r.Body.Close()

	userUid = r.Header.Get("userUid")
	origin := r.Header.Get("origin")

	body := lib.ErrorByte(io.ReadAll(r.Body))
	err = json.Unmarshal(body, &request)
	lib.CheckError(err)

	roles := models.GetAllRoles()
	for _, role := range roles {
		if strings.EqualFold(request.Role, role) {
			userRole = role
			break
		}
	}
	if userRole == "" {
		log.Printf("UpdateUserRole: %s invalid user role", request.Role)
		return `{"success":false}`, `{"success":false}`, nil
	}

	log.Println("UpdateUserRole: get user from firestore")
	fireUser := lib.GetDatasetByEnv(origin, usersCollection)
	docsnap, err := lib.GetFirestoreErr(fireUser, userUid)
	lib.CheckError(err)
	err = docsnap.DataTo(&user)
	lib.CheckError(err)

	log.Println("UpdateUserRole: set user role claim")
	lib.SetCustomClaimForUser(userUid, map[string]interface{}{
		"role": request.Role,
	})

	log.Println("UpdateUserRole: updating user role in DB")
	user.Role = request.Role
	err = lib.SetFirestoreErr(fireUser, userUid, user)
	if err != nil {
		log.Printf("UpdateUserRole: error save user %s firestore: %s", user.Role, err.Error())
		return `{"success":false}`, `{"success":false}`, nil
	}

	err = user.BigquerySave(origin)

	return `{"success":true}`, `{"success":true}`, err
}
