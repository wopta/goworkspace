package user

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func UploadDocument(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var identityDocument models.IdentityDocument

	log.Println("Update Document")

	userUID := r.Header.Get("userUid")

	body := lib.ErrorByte(io.ReadAll(r.Body))
	err := json.Unmarshal(body, &identityDocument)
	if err != nil {
		return "", nil, err
	}

	saveDocument(userUID, &identityDocument)

	updateUser(lib.GetDatasetByEnv(r.Header.Get("origin"), "users"), userUID, &identityDocument)

	identityDocument.FrontMedia.Base64Bytes = ""
	if identityDocument.BackMedia != nil {
		identityDocument.BackMedia.Base64Bytes = ""
	}

	outJson, err := json.Marshal(identityDocument)

	return string(outJson), identityDocument, err
}

func updateUser(fireUsers string, userUID string, identityDocument *models.IdentityDocument) {
	var user models.User
	docsnap := lib.GetFirestore(fireUsers, userUID)
	docsnap.DataTo(&user)
	user.IdentityDocuments = append(user.IdentityDocuments, identityDocument)
	lib.SetFirestore(fireUsers, userUID, user)
}

func saveDocument(userUID string, identityDocument *models.IdentityDocument) {
	now := time.Now()
	timestamp := strconv.FormatInt(now.Unix(), 10)

	documentType, err := getDocumentType(identityDocument)
	lib.CheckError(err)

	bytes, err := base64.StdEncoding.DecodeString(identityDocument.FrontMedia.Base64Bytes)
	lib.CheckError(err)

	fileExtension, err := getFileExtension(identityDocument.FrontMedia.MimeType)
	lib.CheckError(err)

	identityDocument.FrontMedia.Filename = documentType + "_front_" + timestamp + fileExtension
	gsLink, err := lib.PutToGoogleStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "assets/users/"+userUID+"/"+
		identityDocument.FrontMedia.Filename, bytes)
	lib.CheckError(err)
	identityDocument.FrontMedia.Link = gsLink

	if identityDocument.BackMedia != nil {
		bytes, err = base64.StdEncoding.DecodeString(identityDocument.BackMedia.Base64Bytes)
		lib.CheckError(err)

		fileExtension, err = getFileExtension(identityDocument.BackMedia.MimeType)
		lib.CheckError(err)

		identityDocument.BackMedia.Filename = documentType + "_back" + timestamp + fileExtension
		gsLink, err = lib.PutToGoogleStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "assets/users/"+userUID+"/"+
			identityDocument.BackMedia.Filename, bytes)
		lib.CheckError(err)
		identityDocument.BackMedia.Link = gsLink
	}

	identityDocument.LastUpdate = time.Now().UTC()
}

func getDocumentType(identityDocument *models.IdentityDocument) (string, error) {
	switch identityDocument.Code {
	case "01":
		return "identity_document", nil
	case "02":
		return "license", nil
	case "03":
		return "passport", nil
	}
	return "", fmt.Errorf("invalid identity document code")
}

func getFileExtension(mimeType string) (string, error) {
	extensions := map[string]string{
		"application/pdf": ".pdf",
		"image/jpeg":      ".jpeg",
		"image/jpg":       ".jpg",
		"image/png":       ".png",
		"image/webp":      ".webp",
	}

	for mime, extension := range extensions {
		if mime == mimeType {
			return extension, nil
		}
	}
	return "", fmt.Errorf("invalid mime type")
}
