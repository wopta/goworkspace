package win

import (
	"bytes"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func contractCallback(policy models.Policy) error {
	log.Println("win contract calback...")

	bucketFilepath := fmt.Sprintf("assets/users/%s/"+models.ContractDocumentFormat, policy.Contractor.Uid, policy.NameDesc, policy.CodeCompany)
	contractbyte, err := lib.GetFromGoogleStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), bucketFilepath)
	if err != nil {
		return err
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	defer writer.Close()

	fileField, _ := writer.CreateFormFile("file", filepath.Base(bucketFilepath))
	fileField.Write(contractbyte)
	idPraticaField, _ := writer.CreateFormField("idPratica")
	idPraticaField.Write([]byte(strconv.Itoa(policy.NumberCompany)))
	tipoAllegatoField, _ := writer.CreateFormField("tipoAllegato")
	tipoAllegatoField.Write([]byte("POLIZZA"))
	utenteField, _ := writer.CreateFormField("utente")
	utenteField.Write([]byte(policy.Contractor.Name + " " + policy.Contractor.Surname))

	client := &winClient{
		path: "/restba/extquote/emissione",
		headers: map[string]string{
			"Content-Type": writer.FormDataContentType(),
		},
	}
	res, err := client.Post(body)

	// TODO: should we do somethoing with the response?

	log.Println(res)

	return err
}
