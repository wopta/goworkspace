package claim

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

type GetClaimDocumentReq struct {
	DocumentName string `json:"documentName"`
}

type GetClaimDocumentResp struct {
	Document string `json:"document"`
}

func GetClaimDocumentFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		request  GetClaimDocumentReq
		response GetClaimDocumentResp
	)
	log.Println("[GetClaimDocumentFx]")

	authToken, err := lib.VerifyUserIdToken(r.Header.Get("Authorization"))
	if err != nil {
		log.Printf("[GetClaimDocumentFx] invalid idToken, error %s", err.Error())
		return "", "", err
	}

	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()
	log.Println("[GetClaimDocumentFx] " + string(body))
	err = json.Unmarshal(body, &request)
	if err != nil {
		log.Printf("[GetClaimDocumentFx] error parsing body, error %s", err.Error())
		return "", "", err
	}

	res, err := getClaimDocument(r.Header.Get("Origin"), authToken.UID, r.Header.Get("claimUid"), request.DocumentName)
	if err != nil {
		log.Printf("[GetClaimDocumentFx] error getting document, error %s", err.Error())
		return "", "", err
	}

	response.Document = res

	jsonResponse, err := json.Marshal(response)

	return string(jsonResponse), response, err
}

func getClaimDocument(origin, userUid, claimUid, fileName string) (string, error) {
	var user models.User

	fireUser := lib.GetDatasetByEnv(origin, models.UserCollection)
	docsnap, err := lib.GetFirestoreErr(fireUser, userUid)
	if err != nil {
		log.Printf("[getClaimDocument] error retrieving user %s from database, error message %s", userUid, err.Error())
		return "", err
	}
	err = docsnap.DataTo(&user)
	if err != nil {
		log.Println("[getClaimDocument] error convert docsnap to user")
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

	return "", errors.New("[getClaimDocument] not found")
}
