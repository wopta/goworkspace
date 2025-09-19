package user

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/models"
)

type updateUserReq struct {
	Role string `json:"role"`
}

func updateUserRoleFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err               error
		userUid, userRole string
		user              models.User
		request           updateUserReq
	)

	log.AddPrefix("UpdateUserRoleFx")
	defer log.PopPrefix()

	log.Println("Handler start -----------------------------------------------")

	userUid = chi.URLParam(r, "userUid")

	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

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
		err := fmt.Errorf("%s invalid user role", request.Role)
		log.Println(err.Error())
		return "", nil, err
	}

	log.Println("get user from firestore")
	fireUser := lib.UserCollection
	docsnap, err := lib.GetFirestoreErr(fireUser, userUid)
	lib.CheckError(err)
	err = docsnap.DataTo(&user)
	lib.CheckError(err)

	log.Println("set user role claim")
	lib.SetCustomClaimForUser(userUid, map[string]interface{}{
		"role": request.Role,
	})

	log.Println("updating user role in DB")
	user.Role = request.Role
	err = lib.SetFirestoreErr(fireUser, userUid, user)
	if err != nil {
		log.ErrorF("error save user %s firestore: %s", user.Role, err.Error())
		return "", nil, err
	}

	err = user.BigquerySave()

	log.Println("Handler end -------------------------------------------------")

	return "{}", nil, err
}
