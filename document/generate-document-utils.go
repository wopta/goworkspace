package document

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-pdf/fpdf"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/models"
)

type DocumentGenerated struct {
	ParentPath string
	//Default name of the file
	FileName string
	Bytes    []byte
}

func NewDocumentGenerated(parentPath string, filename string, out []byte, err error) (doc DocumentGenerated, errOut error) {
	if err != nil {
		log.ErrorF("error generating contract: %v", err)
		return doc, err
	}
	doc = DocumentGenerated{
		ParentPath: parentPath,
		FileName:   filename,
		Bytes:      out,
	}
	return doc, err
}

// save document in fullpath and return a documentResponse
func (d DocumentGenerated) Save() (result DocumentResponse, err error) {
	log.Printf("Saving document '%v/%v'", d.ParentPath, d.FileName)
	return d.SaveWithName(d.FileName)
}

func (d DocumentGenerated) SaveWithName(name string) (result DocumentResponse, err error) {
	log.Printf("Saving document '%v' with name %v", d.ParentPath, name)
	linkGcs, err := lib.PutToGoogleStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), d.ParentPath+"/"+name, d.Bytes)
	if err != nil {
		return result, err
	}
	return DocumentResponse{
		LinkGcs:  linkGcs,
		FullPath: d.ParentPath + "/" + name,
		Bytes:    base64.StdEncoding.EncodeToString(d.Bytes),
	}, nil
}

func generateContractDocument(pdf *fpdf.Fpdf, policy *models.Policy) (DocumentGenerated, error) {
	var (
		res DocumentGenerated
		out bytes.Buffer
	)
	err := pdf.Output(&out)
	if err != nil {
		return res, err
	}
	res.FileName = fmt.Sprint(policy.NameDesc, "  Contratto ", policy.CodeCompany)
	res.ParentPath = fmt.Sprintf("temp/%s", policy.Uid)
	res.Bytes = out.Bytes()
	return res, nil
}

func generateProposalDocument(pdf *fpdf.Fpdf, policy *models.Policy) (DocumentGenerated, error) {
	var (
		res DocumentGenerated
		out bytes.Buffer
	)
	err := pdf.Output(&out)
	if err != nil {
		return res, err
	}
	res.FileName = fmt.Sprintf(models.ProposalDocumentFormat, policy.NameDesc, policy.ProposalNumber)
	res.ParentPath = fmt.Sprintf("temp/%s", policy.Uid)
	res.Bytes = out.Bytes()
	return res, nil
}

func generateReservedDocument(pdf *fpdf.Fpdf, policy *models.Policy) (DocumentGenerated, error) {
	var (
		res DocumentGenerated
		out bytes.Buffer
	)
	err := pdf.Output(&out)
	if err != nil {
		return res, err
	}
	res.FileName = fmt.Sprintf(models.RvmInstructionsDocumentFormat, policy.ProposalNumber)
	res.ParentPath = fmt.Sprintf("temp/%s", policy.Uid)
	res.Bytes = out.Bytes()
	return res, nil
}
func generateAddendumDocument(pdf *fpdf.Fpdf, policy *models.Policy) (DocumentGenerated, error) {
	var (
		res DocumentGenerated
		out bytes.Buffer
	)
	err := pdf.Output(&out)
	if err != nil {
		return res, err
	}
	res.FileName = strings.ReplaceAll(fmt.Sprintf(addendumDocumentFormat, policy.NameDesc,
		policy.CodeCompany, time.Now().Format("2006-01-02_15:04:05")), " ", "_")
	res.ParentPath = fmt.Sprintf("assets/users/%s/%s", policy.Contractor.Uid, res.FileName)

	res.Bytes = out.Bytes()
	return res, nil
}
