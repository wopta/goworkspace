package document

import (
	"bytes"
	"strconv"
	"time"

	"gitlab.dev.wopta.it/goworkspace/document/internal/engine"
	"gitlab.dev.wopta.it/goworkspace/document/pkg/contract"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/models"
)

func ContractObj(origin string, data models.Policy, networkNode *models.NetworkNode, product *models.Product) <-chan DocumentGenerated {
	r := make(chan DocumentGenerated)
	log.AddPrefix("ContractObj")
	defer log.PopPrefix()

	log.Println("function start -------------------------------")

	rawPolicy, _ := data.Marshal()
	log.Printf("policy: %s", string(rawPolicy))

	go func() {
		var (
			err      error
			document DocumentGenerated
		)

		switch data.Name {
		case models.PmiProduct:
			var buffer bytes.Buffer
			skin := getVar()
			m := skin.initDefault()
			skin.GlobalContract(m, data)
			//-----------Save file
			//TODO: why is this different?
			//filename, out = Save(m, data)

			now := time.Now()
			timestamp := strconv.FormatInt(now.Unix(), 10)

			buffer, err = m.Output()
			out := buffer.Bytes()
			document = DocumentGenerated{
				ParentPath: "temp/" + data.Uid,
				FileName:   data.Contractor.Name + "_" + data.Contractor.Surname + "_" + timestamp + "_contract.pdf",
				Bytes:      out,
			}
		case models.LifeProduct:
			pdf := engine.NewFpdf()
			document, err = lifeContract(pdf, origin, &data, networkNode, product)
		case models.CatNatProduct:
			//TODO: to change
			//filename, out = "prova catnat contratto", []byte{}
		case models.PersonaProduct:
			pdf := engine.NewFpdf()
			generator := contract.NewPersonaGenerator(pdf, &data, networkNode, *product, false)
			personaGlobalContractV1(pdf.GetPdf(), &data, networkNode, product)
			generator.AddMup()
			document, err = generateContractDocument(pdf.GetPdf(), &data)
		case models.GapProduct:
			pdf := initFpdf()
			document, err = gapSogessurContractV1(pdf, origin, &data, networkNode)
		}
		if err != nil {
			log.ErrorF("error generating contract: %v", err)
			return
		}

		log.Println(data.Uid + " ContractObj end")
		r <- document
	}()

	log.Println("function end -------------------------------..")

	return r
}

func lifeContract(enginePdf *engine.Fpdf, origin string, policy *models.Policy, networkNode *models.NetworkNode, product *models.Product) (DocumentGenerated, error) {
	var (
		document DocumentGenerated
		err      error
	)

	log.AddPrefix("LifeContract")
	defer log.PopPrefix()
	log.Println("function start ------------------------------")

	switch policy.ProductVersion {
	case models.ProductV1:
		log.Println("life v1")
		pdf := enginePdf.GetPdf()
		document, err = lifeAxaContractV1(pdf, origin, policy, networkNode, product)
	case models.ProductV2:
		log.Println("life v2")
		pdf := enginePdf.GetPdf()
		generator := contract.NewLifeGenerator(enginePdf, policy, networkNode, *product, false)
		lifeAxaContractV2(pdf, origin, policy, networkNode, product)
		generator.AddMup()
		document, err = generateContractDocument(pdf, policy)
	}

	log.Println("function end --------------------------------")

	return document, err
}
