package contract

import (
	"fmt"
	"time"

	"gitlab.dev.wopta.it/goworkspace/document/internal/constants"
	"gitlab.dev.wopta.it/goworkspace/document/internal/domain"
	"gitlab.dev.wopta.it/goworkspace/document/internal/dto"
	"gitlab.dev.wopta.it/goworkspace/document/internal/engine"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/network"
)

type LifeGenerator struct {
	*baseGenerator
	dtoLife dto.LifeDTO
}

func NewLifeGenerator(engine *engine.Fpdf, policy *models.Policy, node *models.NetworkNode, product models.Product, isProposal bool) *LifeGenerator {
	dto := dto.NewLifeDto()
	dto.FromPolicy(policy, node)

	var worksForNode *models.NetworkNode
	if node != nil && node.WorksForUid != "" {
		worksForNode = network.GetNetworkNodeByUid(node.WorksForUid)
	}

	return &LifeGenerator{
		baseGenerator: &baseGenerator{
			engine:       engine,
			isProposal:   isProposal,
			now:          time.Now(),
			signatureID:  0,
			networkNode:  node,
			policy:       policy,
			worksForNode: worksForNode,
		},
		dtoLife: dto,
	}
}

func (el *LifeGenerator) Generate() {
	el.woptaFooter()
	el.addMainHeader()

	el.engine.NewPage()
	el.engine.NewLine(constants.CellHeight)
	el.addHeading()

	el.engine.NewLine(constants.CellHeight * 2)
	el.addWelcomeSection()
	el.engine.NewLine(constants.CellHeight)

	el.addEmailSection()
	el.engine.NewLine(constants.CellHeight)

	el.addSignSection()

	el.addPolicyInformationSection()

	el.engine.NewLine(constants.CellHeight)
	el.addSupportInformationSection()

	el.engine.NewLine(constants.CellHeight)
	el.addGreatingsSection()
}

func (el *LifeGenerator) addMainHeader() {
	const (
		firstColumnWidth  = 115
		secondColumnWidth = constants.FullPageWidth - firstColumnWidth
	)
	parseData := func(rows [][]string) [][]domain.TableCell {
		result := make([][]domain.TableCell, 0, len(rows))

		for index, row := range rows {
			fontStyle := constants.RegularFontStyle
			if index == 0 {
				fontStyle = constants.BoldFontStyle
			}
			result = append(result, []domain.TableCell{
				{
					Text:      row[0],
					Height:    constants.CellHeight,
					Width:     firstColumnWidth,
					FontStyle: fontStyle,
					FontColor: constants.BlackColor,
					FontSize:  constants.MediumFontSize,
					Fill:      false,
					FillColor: domain.Color{},
					Align:     constants.LeftAlign,
					Border:    "",
				},
				{
					Text:      row[1],
					Height:    constants.CellHeight,
					FontStyle: fontStyle,
					Width:     secondColumnWidth,
					FontColor: constants.BlackColor,
					FontSize:  constants.MediumFontSize,
					Fill:      false,
					FillColor: domain.Color{},
					Align:     constants.LeftAlign,
					Border:    "",
				},
			})
		}
		return result
	}

	parseLogos := func(texts []string) [][]domain.TableCell {
		result := make([][]domain.TableCell, 0, len(texts))

		result = append(result, []domain.TableCell{
			{
				Text:      "",
				Height:    constants.CellHeight,
				Width:     firstColumnWidth,
				FontStyle: constants.RegularFontStyle,
				FontColor: constants.BlackColor,
				FontSize:  constants.RegularFontSize,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "",
			},
			{
				Text:      texts[0],
				Height:    constants.CellHeight,
				FontStyle: constants.BoldFontStyle,
				FontColor: constants.PinkColor,
				FontSize:  18,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "",
			},
		})

		result = append(result, []domain.TableCell{
			{
				Text:      "",
				Height:    constants.CellHeight,
				Width:     firstColumnWidth,
				FontStyle: constants.RegularFontStyle,
				FontColor: constants.BlackColor,
				FontSize:  constants.RegularFontSize,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "",
			},
			{
				Text:      texts[1],
				Height:    constants.CellHeight + 5,
				FontStyle: constants.RegularFontStyle,
				FontColor: constants.BlackColor,
				FontSize:  16,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "",
			},
		})
		return result
	}

	rowsData := [][]string{
		{"I dati del tuo Preventivo", "I tuoi dati"},
		{"Numero: " + el.dtoLife.ProposalNumber, "Contraente: " + el.dtoLife.Contractor.GetFullNameContractor()},
		{"Decorre dal: " + el.dtoLife.ValidityDate.StartDate, "C.F./P.IVA: " + el.dtoLife.Contractor.FiscalCode},
		{"Scade il: " + el.dtoLife.ValidityDate.EndDate, "Indirizzo: " + el.dtoLife.GetAddressFirstPart()},
		{"Prima scadenza Annuale il: " + el.dtoLife.ValidityDate.FirstAnnuityExpiry, el.dtoLife.GetAddressSecondPart()},
		{"Non si rinnova a scadenza.", "Mail: " + el.dtoLife.Contractor.Mail},
		{"Produttore: " + el.dtoLife.ProductorName, "Telefono: " + el.dtoLife.Contractor.Phone},
	}

	el.engine.SetHeader(func() {
		el.engine.InsertImage(lib.GetAssetPathByEnvV2()+"logo_wopta.png", 10, 5, 35, 10)
		el.baseGenerator.engine.NewLine(3)
		el.engine.DrawTable(parseLogos([]string{"Wopta per te", "Vita"}))
		el.engine.InsertImage(lib.GetAssetPathByEnvV2()+"logo_vita.png", 180, 15, 13, 13)

		el.engine.DrawTable(parseData(rowsData))
	})
}

func (el *LifeGenerator) addHeading() {
	el.engine.WriteText(el.engine.GetTableCell("Il tuo Preventivo: cosa fare adesso?", constants.BoldFontStyle, constants.PinkColor, constants.LargeFontSize))
}
func (el *LifeGenerator) addWelcomeSection() {
	el.engine.WriteTexts(
		el.engine.GetTableCell(fmt.Sprintf("Buongiorno %v %v,\n", el.dtoLife.Contractor.Name, el.dtoLife.Contractor.Surname), constants.BlackColor),
		el.engine.GetTableCell("Grazie per aver fatto un preventivo per una polizza Vita, dimostrando volontà e interesse a tutelarti e/o proteggere le persone per te più importanti.", constants.BlackColor),
	)
}

func (el *LifeGenerator) addEmailSection() {
	el.engine.WriteText(el.engine.GetTableCell(
		"In allegato trovi:\n"+
			"- modulo di Polizza\n"+
			"- informativa precontrattuale di Wopta, prevista per legge\n"+
			"- modulistica antiriciclaggio\n\n"+
			"- informativa e dichiarazioni privacy per l’Assicuratore\n"+
			"- informativa e dichiarazioni privacy per l’Intermediario", constants.BlackColor,
	))

	el.engine.WriteText(el.engine.GetTableCell("Verifica la correttezza di tutti i dati inseriti (anagrafici, indirizzi, codice fiscale, contatti) e delle prestazioni scelte (durata, importi, eventuali opzioni).", constants.BlackColor))

	el.engine.NewLine(constants.CellHeight)
	if el.dtoLife.Channel != models.NetworkChannel {
		el.engine.WriteTexts(
			el.engine.GetTableCell("Riceverai anche due mail per procedere con la ", constants.BlackColor),
			el.engine.GetTableCell("firma ", constants.PinkColor, constants.BoldFontStyle),
			el.engine.GetTableCell("ed il ", constants.BlackColor),
			el.engine.GetTableCell("pagamento", constants.PinkColor, constants.BoldFontStyle),
		)
	}
}

func (el *LifeGenerator) addSignSection() {
	el.engine.WriteTexts(
		el.engine.GetTableCell("ATTENZIONE", constants.PinkColor, constants.BoldFontStyle),
		el.engine.GetTableCell(": Solo una volta firmati i documenti ed effettuato il pagamento, la copertura assicurativa sarà attiva e così ti invieremo i documenti contrattuali da te firmati, che poi potrai visualizzare nell’area riservata ai clienti della nostra app e/o sito.", constants.BlackColor),
	)
}

func (el *LifeGenerator) addPolicyInformationSection() {
	if el.dtoLife.ConsultancyValue.Price.ValueFloat == 0 {
		return
	}
	el.engine.NewLine(constants.CellHeight)
	text :=
		"Infine, ti ricordiamo la presente polizza prevede il pagamento dei seguenti costi:\n" +
			fmt.Sprintf("- Premio di polizza: euro %v con frazionamento %v\n", el.dtoLife.Prizes.Gross.Text, el.dtoLife.Prizes.Split) +
			fmt.Sprintf("- Contributo servizi di intermediazione annuale: euro %v corrisposti con il pagamento della prima rata di polizza\n", el.dtoLife.ConsultancyValue.Price.Text) +
			fmt.Sprintf("- Per un totale annuo di euro %v", el.dtoLife.PriceAnnuity)

	el.engine.WriteText(el.engine.GetTableCell(text, constants.BlackColor))
}

func (el *LifeGenerator) addSupportInformationSection() {
	if el.dtoLife.Channel != models.NetworkChannel {
		el.engine.RawWriteText(
			el.engine.GetTableCell("Restiamo a disposizione per ogni ulteriore informazione anche attraverso i canali di contatto che trovi a questo ", constants.BlackColor),
		)
		el.engine.WriteLink("https://www.wopta.it/it/vita/#contact-us", el.engine.GetTableCell("link", constants.PinkColor))
		el.engine.RawWriteText(el.engine.GetTableCell(".", constants.BlackColor))
	} else {
		el.engine.RawWriteText(
			el.engine.GetTableCell("Se hai necessità di ulteriori informazioni e supporto, rivolgiti al tuo intermediario, che trovi in copia conoscenza alla mail accompagnatoria di questa comunicazione.", constants.BlackColor))
	}
	el.engine.NewLine(constants.CellHeight)
}

func (el *LifeGenerator) addGreatingsSection() {
	el.engine.WriteText(el.engine.GetTableCell("Cordiali saluti.", constants.BlackColor))
	el.engine.NewLine(constants.CellHeight)
	el.engine.WriteTexts(
		el.engine.GetTableCell("Anna di Wopta Assicurazioni\n", constants.BlackColor),
		el.engine.GetTableCell("Proteggiamo chi sei", constants.BlackColor),
	)
}
