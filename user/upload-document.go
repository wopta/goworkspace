package user

import (
	"encoding/base64"
	"encoding/json"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type UploadDocumentRequest struct {
	UserUID          string                   `json:"userUID"`
	IdentityDocument *models.IdentityDocument `json:"identityDocument"`
}

func UploadDocument(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		request UploadDocumentRequest
		user    models.User
	)

	log.Println("Update Document")

	body := lib.ErrorByte(io.ReadAll(r.Body))
	err := json.Unmarshal(body, &request)
	if err != nil {
		return "", nil, err
	}

	saveDocument(request.UserUID, request.IdentityDocument)

	updateUser(lib.GetDatasetByEnv(r.Header.Get("origin"), "users"), &request, &user)

	outJson, err := json.Marshal(request.IdentityDocument)

	return string(outJson), request.IdentityDocument, err
}

func updateUser(fireUsers string, request *UploadDocumentRequest, user *models.User) {
	docsnap := lib.GetFirestore(fireUsers, request.UserUID)
	docsnap.DataTo(&user)
	user.IdentityDocuments = append(user.IdentityDocuments, request.IdentityDocument)
	lib.SetFirestore(fireUsers, request.UserUID, user)
}

func saveDocument(userUID string, identityDocument *models.IdentityDocument) {
	bytes, err := base64.StdEncoding.DecodeString(identityDocument.FrontMedia.Base64Encoding)
	lib.CheckError(err)

	gsLink, err := lib.PutToGoogleStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "asset/"+userUID+"/"+
		identityDocument.Type+"_front."+identityDocument.FrontMedia.MimeType, bytes)
	lib.CheckError(err)
	identityDocument.FrontMedia.Link = gsLink

	gsLink, err = lib.PutToGoogleStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "asset/"+userUID+"/"+
		identityDocument.Type+"_back."+identityDocument.FrontMedia.MimeType, bytes)
	lib.CheckError(err)
	identityDocument.BackMedia.Link = gsLink

	identityDocument.LastUpdate = time.Now().UTC()
}
