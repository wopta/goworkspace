package document

import (
	"log"

	"github.com/wopta/goworkspace/document/internal/engine"
	"github.com/wopta/goworkspace/document/pkg/proforma"
	"github.com/wopta/goworkspace/models"
)

type ProformaResponse struct {
	FileName string
	LinkGcs  string `json:"linkGcs"`
}

func Proforma(policy models.Policy, networkNode *models.NetworkNode, product *models.Product) (ProformaResponse, error) {
	var (
		err      error
		filename string
		out      []byte
	)

	log.Println("[ProformaObj] function start -------------------------------")

	generator := proforma.NewProformaGenerator(engine.NewFpdf(), &policy, networkNode, *product)
	out, err = generator.Generate()
	if err != nil {
		log.Printf("error generating proforma: %v", err)
		return ProformaResponse{}, err
	}
	filename, err = generator.Save(out)
	if err != nil {
		log.Printf("error saving proforma: %v", err)
		return ProformaResponse{}, err
	}

	log.Println(policy.Uid + " ProformaObj end")
	res := ProformaResponse{
		LinkGcs:  filename,
		FileName: filename,
	}

	log.Println("[ProformaObj] function end -------------------------------..")

	return res, nil
}
