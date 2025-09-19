package claim

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"

	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/models"
)

type GetClaimDocumentReq struct {
	DocumentName string `json:"documentName"`
}

type GetClaimDocumentResp struct {
	Document string `json:"document"`
}

func getClaimDocumentFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		request  GetClaimDocumentReq
		response GetClaimDocumentResp
	)
	log.AddPrefix("GetClaimDocumentFx")
	defer log.PopPrefix()

	log.Println("Handler start -----------------------------------------------")

	authToken, err := lib.VerifyUserIdToken(r.Header.Get("Authorization"))
	if err != nil {
		log.ErrorF("invalid idToken, error %s", err.Error())
		return "", "", err
	}

	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	err = json.Unmarshal(body, &request)
	if err != nil {
		log.ErrorF("error parsing body, error %s", err.Error())
		return "", "", err
	}

	res, err := getClaimDocument(authToken.UID, chi.URLParam(r, "claimUid"), request.DocumentName)
	if err != nil {
		log.ErrorF("error getting document, error %s", err.Error())
		return "", "", err
	}

	response.Document = res

	jsonResponse, err := json.Marshal(response)

	log.Println("Handler end -------------------------------------------------")

	return string(jsonResponse), response, err
}

func getClaimDocument(userUid, claimUid, fileName string) (string, error) {
	var user models.User
	log.AddPrefix("getClaimDocument")
	defer log.PopPrefix()
	fireUser := lib.UserCollection
	docsnap, err := lib.GetFirestoreErr(fireUser, userUid)
	if err != nil {
		log.Printf("error retrieving user %s from database, error message %s", userUid, err.Error())
		return "", err
	}
	err = docsnap.DataTo(&user)
	if err != nil {
		log.Println("error convert docsnap to user")
		return "", err
	}

	if user.Claims != nil {
		for _, userClaim := range *user.Claims {
			if userClaim.ClaimUid == claimUid {
				for _, document := range userClaim.Documents {
					if document.FileName == fileName {
						return base64.StdEncoding.EncodeToString(
							lib.GetFromStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "assets/users/"+userUid+
								"/claims/"+claimUid+"/"+document.FileName, "")), nil
					}
				}
			}
		}
	}

	return "", errors.New("not found")
}
