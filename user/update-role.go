package user

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
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

	log.SetPrefix("[UpdateUserRoleFx] ")
	defer log.SetPrefix("")

	log.Println("Handler start -----------------------------------------------")

	userUid = chi.URLParam(r, "userUid")
	origin := r.Header.Get("Origin")

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
	fireUser := lib.GetDatasetByEnv(origin, lib.UserCollection)
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
		log.Printf("error save user %s firestore: %s", user.Role, err.Error())
		return "", nil, err
	}

	err = user.BigquerySave(origin)

	log.Println("Handler end -------------------------------------------------")

	return "{}", nil, err
}
