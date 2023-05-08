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
)

func UploadDocument(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var identityDocument models.IdentityDocument

	log.Println("Update Document")

	b := lib.ErrorByte(io.ReadAll(r.Body))
	err := json.Unmarshal(b, identityDocument)
	if err != nil {
		return `{"link":"", "success":"false"}`, `{"link":"", "success":"false"}`, err
	}

	link := saveDocument("1234", &identityDocument)

	return `{"link":"` + link + `", "success":"true"}`, `{"link":"` + link + `", "success":"true"}`, err
}

func saveDocument(userUID string, identityDocument *models.IdentityDocument) string {
	var filename, link string

	bytes, err := base64.StdEncoding.DecodeString(identityDocument.Base64Encoding)
	lib.CheckError(err)

	filename = "asset/" + userUID + "_" + identityDocument.Type + "." + identityDocument.MimeType
	link = lib.PutToStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), filename, bytes)
	return link
}
