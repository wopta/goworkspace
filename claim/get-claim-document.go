package claim

import (
	"encoding/base64"
	"encoding/json"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"io"
	"log"
	"net/http"
	"os"
)

type GetClaimDocumentReq struct {
	DocumentName string `json:"documentName"`
}

type GetClaimDocumentResp struct {
	Document string `json:"document"`
}

func GetClaimDocumentFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		user     models.User
		request  GetClaimDocumentReq
		response GetClaimDocumentResp
	)
	log.Println("GetClaimDocument")

	idToken := r.Header.Get("Authorization")
	authToken, err := lib.VerifyUserIdToken(idToken)
	if err != nil {
		log.Printf("[GetClaimDocument] invalid idToken, error %s", err.Error())
		return "", "", err
	}

	claimUID := r.Header.Get("claimUid")

	fireUser := lib.GetDatasetByEnv(r.Header.Get("Origin"), models.UserCollection)
	docsnap, err := lib.GetFirestoreErr(fireUser, authToken.UID)
	if err != nil {
		log.Printf("[GetClaimDocument] error retrieving user %s from database, error message %s", authToken.UID, err.Error())
		return "", "", err
	}
	err = docsnap.DataTo(&user)
	if err != nil {
		log.Println("[GetClaimDocument] error convert docsnap to user")
		return "", "", err
	}

	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()
	log.Println("[GetClaimDocument] " + string(body))
	err = json.Unmarshal(body, &request)
	if err != nil {
		log.Printf("[GetClaimDocument] error parsing body, error %s", err.Error())
		return "", "", err
	}

	if user.Claims != nil {
		for _, userClaim := range *user.Claims {
			if userClaim.ClaimUid == claimUID {
				for _, document := range userClaim.Documents {
					if document.FileName == request.DocumentName {
						response.Document = base64.StdEncoding.EncodeToString(
							lib.GetFromStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "assets/users/"+authToken.UID+
								"/claims/"+claimUID+"/"+document.FileName, ""))
						break
					}
				}
			}
		}
	}

	jsonResponse, err := json.Marshal(response)

	return string(jsonResponse), response, err
}
