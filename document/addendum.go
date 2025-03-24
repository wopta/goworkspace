package document

import (
	"fmt"
	"log"

	"github.com/wopta/goworkspace/document/internal/engine"
	"github.com/wopta/goworkspace/document/pkg/addedndum"
	"github.com/wopta/goworkspace/models"
)

type AddendumResponse struct {
	LinkGcs  string `json:"linkGcs"`
	Filename string `json:"fileName"`
}

func Addendum(origin string, data models.Policy, networkNode *models.NetworkNode, product *models.Product) (AddendumResponse, error) {
	var (
		err      error
		filename string
		out      []byte
	)

	log.Println("[AddendumObj] function start -------------------------------")

	rawPolicy, _ := data.Marshal()
	log.Printf("[AddendumObj] policy: %s", string(rawPolicy))

	switch data.Name {
	case models.LifeProduct:
		prod := models.Product{}
		generator := addedndum.NewLifeAddendumGenerator(engine.NewFpdf(), &data, networkNode, prod)
		out, err = generator.Generate()
		if err != nil {
			log.Printf("error generating addendum: %v", err)
			return AddendumResponse{}, err
		}
		filename, err = generator.Save(out)
		if err != nil {
			log.Printf("error saving addendum: %v", err)
			return AddendumResponse{}, err
		}
	default:
		return AddendumResponse{}, fmt.Errorf("addendum not implemented for product %s", data.Name)
	}

	data.DocumentName = filename
	log.Println(data.Uid + " AddendumObj end")
	res := AddendumResponse{
		LinkGcs:  filename,
		Filename: filename,
	}

	log.Println("[AddendumObj] function end -------------------------------..")

	return res, nil
}
