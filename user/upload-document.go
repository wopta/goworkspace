package user

import (
	"encoding/base64"
	"encoding/json"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"strconv"
	"time"
)

type UploadDocumentRequest struct {
	UserUID          string                   `json:"userUID"`
	IdentityDocument *models.IdentityDocument `json:"identityDocument"`
}

func UploadDocument(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var request *UploadDocumentRequest

	log.Println("Update Document")

	body := lib.ErrorByte(io.ReadAll(r.Body))
	err := json.Unmarshal(body, request)
	if err != nil {
		return "", nil, err
	}

	saveDocument(request.UserUID, request.IdentityDocument)

	updateUser(lib.GetDatasetByEnv(r.Header.Get("origin"), "users"), request)

	request.IdentityDocument.FrontMedia.Base64Encoding = ""
	if request.IdentityDocument.BackMedia != nil {
		request.IdentityDocument.BackMedia.Base64Encoding = ""
	}

	outJson, err := json.Marshal(request.IdentityDocument)

	return string(outJson), request.IdentityDocument, err
}

func updateUser(fireUsers string, request *UploadDocumentRequest) {
	var user models.User

	docsnap := lib.GetFirestore(fireUsers, request.UserUID)
	docsnap.DataTo(&user)
	user.IdentityDocuments = append(user.IdentityDocuments, request.IdentityDocument)
	lib.SetFirestore(fireUsers, request.UserUID, user)
}

func saveDocument(userUID string, identityDocument *models.IdentityDocument) {
	now := time.Now()
	timestamp := strconv.FormatInt(now.Unix(), 10)

	bytes, err := base64.StdEncoding.DecodeString(identityDocument.FrontMedia.Base64Encoding)
	lib.CheckError(err)

	fileExtensions, err := mime.ExtensionsByType(identityDocument.FrontMedia.MimeType)
	lib.CheckError(err)
	gsLink, err := lib.PutToGoogleStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "asset/"+userUID+"/"+
		identityDocument.Type+timestamp+"_front."+fileExtensions[0], bytes)
	lib.CheckError(err)
	identityDocument.FrontMedia.Link = "gs://" + os.Getenv("GOOGLE_STORAGE_BUCKET") + gsLink

	if identityDocument.BackMedia != nil {
		bytes, err = base64.StdEncoding.DecodeString(identityDocument.BackMedia.Base64Encoding)
		lib.CheckError(err)

		fileExtensions, err = mime.ExtensionsByType(identityDocument.BackMedia.MimeType)
		lib.CheckError(err)
		gsLink, err = lib.PutToGoogleStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "asset/"+userUID+"/"+
			identityDocument.Type+timestamp+"_back."+fileExtensions[0], bytes)
		lib.CheckError(err)
		identityDocument.BackMedia.Link = "gs://" + os.Getenv("GOOGLE_STORAGE_BUCKET") + gsLink
	}

	identityDocument.LastUpdate = time.Now().UTC()
}
