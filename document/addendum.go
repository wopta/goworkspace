package document

import (
	"log"

	"github.com/wopta/goworkspace/document/internal/engine"
	"github.com/wopta/goworkspace/document/pkg/addedndum"
	"github.com/wopta/goworkspace/models"
)

type AddendumResponse struct {
	LinkGcs  string `json:"linkGcs"`
	Filename string `json:"fileName"`
}

func AddendumObj(origin string, data models.Policy, networkNode *models.NetworkNode, product *models.Product) <-chan AddendumResponse {
	r := make(chan AddendumResponse)

	log.Println("[AddendumObj] function start -------------------------------")

	rawPolicy, _ := data.Marshal()
	log.Printf("[AddendumObj] policy: %s", string(rawPolicy))

	go func() {
		var (
			err      error
			filename string
			out      []byte
		)

		switch data.Name {
		case models.PmiProduct:

		case models.LifeProduct:
			prod := models.Product{}
			generator := addedndum.NewLifeAddendumGenerator(engine.NewFpdf(), &data, networkNode, prod)
			out, err = generator.Contract()
			if err != nil {
				log.Printf("error generating addendum: %v", err)
				return
			}
			filename, err = generator.Save(out)
			if err != nil {
				log.Printf("error saving addendum: %v", err)
				return
			}
		case models.PersonaProduct:

		case models.GapProduct:

		case models.CommercialCombinedProduct:

		}

		data.DocumentName = filename
		log.Println(data.Uid + " AddendumObj end")
		r <- AddendumResponse{
			LinkGcs:  filename,
			Filename: filename,
		}
	}()

	log.Println("[AddendumObj] function end -------------------------------..")

	return r
}

/*func GetGenerator(origin string, data models.Policy, networkNode *models.NetworkNode, product *models.Product) {
	log.Println("[GetGenerator] function start -------------------------------")

	rawPolicy, _ := data.Marshal()
	log.Printf("[GetGenerator] policy: %s", string(rawPolicy))

	switch data.Name {
	case models.PmiProduct:

	case models.LifeProduct:
		generator := addedndum.NewLifeAddendumGenerator(engine.NewFpdf(), &data, networkNode, *product)
		out, err := generator.Contract()
		if err != nil {
			log.Printf("error generating addendum: %v", err)
			return
		}
		filename, err := generator.Save(out)
		if err != nil {
			log.Printf("error saving addendum: %v", err)
			return
		}
	case models.PersonaProduct:

	case models.GapProduct:

	case models.CommercialCombinedProduct:

	}
}*/
