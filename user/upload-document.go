package user

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func UploadDocument(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var identityDocument models.IdentityDocument

	log.Println("Upload User IdentityDocument")

	policyUID := r.Header.Get("policyUid")

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

	saveDocument(policyUID, &identityDocument)

	outJson, err := json.Marshal(identityDocument)

	return string(outJson), identityDocument, err
}

func saveDocument(policyUID string, identityDocument *models.IdentityDocument) {
	saveToStorage := func(policyUID string, documentSide, documentType string, media *models.Media) error {
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

		media.FileName = documentType + "_" + documentSide + "_" + timestamp + fileExtension
		gsLink, err := lib.PutToGoogleStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "temp/"+policyUID+"/"+
			media.FileName, bytes)
		media.Link = gsLink
		return err
	}

	documentType, err := getDocumentType(identityDocument)
	lib.CheckError(err)

	err = saveToStorage(policyUID, "front", documentType, identityDocument.FrontMedia)
	lib.CheckError(err)

	if identityDocument.BackMedia != nil {
		err = saveToStorage(policyUID, "back", documentType, identityDocument.BackMedia)
		lib.CheckError(err)
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
