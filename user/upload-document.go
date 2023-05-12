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

	if identityDocument.DateOfIssue.After(identityDocument.ExpiryDate) {
		return "", nil, fmt.Errorf("date of issue cannot be after expiration date")
	}
	if identityDocument.IsExpired() {
		return "", nil, fmt.Errorf("identity document expired")
	}

	saveDocument(userUID, &identityDocument)

	updateUser(lib.GetDatasetByEnv(r.Header.Get("origin"), "users"), userUID, &identityDocument)

	outJson, err := json.Marshal(identityDocument)

	return string(outJson), identityDocument, err
}

func saveDocument(userUID string, identityDocument *models.IdentityDocument) {
	saveToStorage := func(userUID string, documentSide, documentType string, media *models.Media) error {
		now := time.Now()
		timestamp := strconv.FormatInt(now.Unix(), 10)

		bytes, err := base64.StdEncoding.DecodeString(media.Base64Bytes)
		if err != nil {
			return err
		}
		media.Base64Bytes = ""

		fileExtension, err := getFileExtension(media.MimeType)
		if err != nil {
			return err
		}

		media.Filename = documentType + "_" + documentSide + "_" + timestamp + fileExtension
		gsLink, err := lib.PutToGoogleStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "assets/users/"+userUID+"/"+
			media.Filename, bytes)
		media.Link = gsLink
		return err
	}

	documentType, err := getDocumentType(identityDocument)
	lib.CheckError(err)

	err = saveToStorage(userUID, "front", documentType, identityDocument.FrontMedia)
	lib.CheckError(err)

	if identityDocument.BackMedia != nil {
		err = saveToStorage(userUID, "back", documentType, identityDocument.BackMedia)
		lib.CheckError(err)
	}

	identityDocument.LastUpdate = time.Now().UTC()
}

func updateUser(fireUsers string, userUID string, identityDocument *models.IdentityDocument) {
	var user models.User
	docsnap := lib.GetFirestore(fireUsers, userUID)
	docsnap.DataTo(&user)
	user.IdentityDocuments = append(user.IdentityDocuments, identityDocument)
	user.UpdatedDate = time.Now().UTC()
	lib.SetFirestore(fireUsers, userUID, user)
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
