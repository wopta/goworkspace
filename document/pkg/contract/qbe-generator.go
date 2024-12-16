package contract

import (
	"github.com/wopta/goworkspace/document/internal/constants"
	"github.com/wopta/goworkspace/document/internal/domain"
	"github.com/wopta/goworkspace/document/internal/engine"
	"github.com/wopta/goworkspace/lib"
)

type QBEGenerator struct {
	*baseGenerator
}

func NewQBEGenerator(engine *engine.Fpdf) *QBEGenerator {
	return &QBEGenerator{
		&baseGenerator{engine: engine},
	}
}

func (qb *QBEGenerator) mainHeader() {
	table := [][]domain.TableCell{
		{
			{
				Text:      "I dati della tua Polizza nr. 100100", // TODO: add dynamic code company
				Height:    constants.CellHeight,
				Width:     115,
				TextBold:  true,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "",
			},
			{
				Text:      "I tuoi dati",
				Height:    constants.CellHeight,
				Width:     75,
				TextBold:  true,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "",
			},
		},
		{
			{
				Text:      "Decorre dal 13/12/2024 ore 24:00", // TODO: add dynamic startDate
				Height:    constants.CellHeight,
				Width:     115,
				TextBold:  false,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "",
			},
			{
				Text:      "Contraente Wopta Assicurazioni S.R.L", // TODO: add dynamic contractor name
				Height:    constants.CellHeight,
				Width:     75,
				TextBold:  false,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "",
			},
		},
		{
			{
				Text:      "Si rinnova a scadenza, salvo disdetta da inviare 30 giorni prima",
				Height:    constants.CellHeight,
				Width:     115,
				TextBold:  false,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "",
			},
			{
				Text:      "P.IVA: 012345678910", // TODO: add dynamic vatCode
				Height:    constants.CellHeight,
				Width:     75,
				TextBold:  false,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "",
			},
		},
		{
			{
				Text:      "Frazionamento: MENSILE", // TODO: add dynamic payment split
				Height:    constants.CellHeight,
				Width:     115,
				TextBold:  false,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "",
			},
			{
				Text:      "Codice Fiscale: HMMYSF94R07D912M", // TODO: add dynamic fiscalCode
				Height:    constants.CellHeight,
				Width:     75,
				TextBold:  false,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "",
			},
		},
		{
			{
				Text:      "Prossimo pagamento il: 16/12/2025", // TODO: add dynamic nextPayment date
				Height:    constants.CellHeight,
				Width:     115,
				TextBold:  false,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "",
			},
			{
				Text:      "Indirizzo: Galleria del corso 1, Milano (MI)", // TODO: add dynamic address
				Height:    constants.CellHeight,
				Width:     75,
				TextBold:  false,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "",
			},
		},
		{
			{
				Text:      "Presenza Vincolo: NO", // TODO: add dynamic hasBond
				Height:    constants.CellHeight,
				Width:     115,
				TextBold:  false,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "",
			},
			{
				Text:      "Telefono: +393334455667", // TODO: add dynamic phone
				Height:    constants.CellHeight,
				Width:     75,
				TextBold:  false,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "",
			},
		},
		{
			{
				Text:      "Sostituisce la Polizza: 1234567", // TODO: add dynamic old codeCompany
				Height:    constants.CellHeight,
				Width:     115,
				TextBold:  false,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "",
			},
			{
				Text:      "Mail: wopta@wopta.it", // TODO: add dynamic mail
				Height:    constants.CellHeight,
				Width:     75,
				TextBold:  false,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "",
			},
		},
		{
			{
				Text:      "Convenzione: NO",
				Height:    constants.CellHeight,
				Width:     190,
				TextBold:  false,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "",
			},
		},
	}

	qb.engine.SetHeader(func() {
		qb.engine.InsertImage(lib.GetAssetPathByEnvV2()+"logo_qbe.png", 75, 6.5, 22, 8)
		qb.engine.DrawLine(102, 6.25, 102, 15, 0.25, constants.BlackColor)
		qb.engine.InsertImage(lib.GetAssetPathByEnvV2()+"logo_wopta.png", 107.5, 5, 35, 12)
		qb.engine.NewLine(7)
		qb.engine.DrawTable(table)
	})
}

func (qb *QBEGenerator) mainFooter() {

}

func (qb *QBEGenerator) Contract() ([]byte, error) {
	qb.mainHeader()

	return qb.engine.RawDoc()
}
