package user

import (
	"encoding/json"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"io"
	"log"
	"net/http"
)

type UpdateUserReq struct {
	Role string `json:"role"`
}

func UpdateUserRoleFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		userUid string
		user    models.User
		request UpdateUserReq
		err     error
	)

	log.Println("Update User Role")

	defer r.Body.Close()

	userUid = r.Header.Get("userUid")

	body := lib.ErrorByte(io.ReadAll(r.Body))
	err = json.Unmarshal(body, &request)
	lib.CheckError(err)

	log.Println("UpdateUserRole: get user from firestore")
	fireUser := lib.GetDatasetByEnv(r.Header.Get("origin"), "users")
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
	lib.SetFirestore(fireUser, userUid, user)

	return `{"success":true}`, `{"success":true}`, err
}
