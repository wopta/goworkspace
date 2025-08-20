package contract

import (
	"fmt"
	"strings"
	"time"

	"gitlab.dev.wopta.it/goworkspace/document/internal/constants"
	"gitlab.dev.wopta.it/goworkspace/document/internal/domain"
	"gitlab.dev.wopta.it/goworkspace/document/internal/dto"
	"gitlab.dev.wopta.it/goworkspace/document/internal/engine"
	"gitlab.dev.wopta.it/goworkspace/document/pkg/internal/catnat"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/network"
)

type CatnatGenerator struct {
	*baseGenerator
	dtoCatnat dto.CatnatDTO
}

func NewCatnatGenerator(engine *engine.Fpdf, policy *models.Policy, node *models.NetworkNode, product models.Product, isProposal bool) *CatnatGenerator {
	var worksForNode *models.NetworkNode
	if node != nil && node.WorksForUid != "" {
		worksForNode = network.GetNetworkNodeByUid(node.WorksForUid)
	}
	dto := dto.CatnatDTO{}
	dto.FromPolicy(policy)
	return &CatnatGenerator{
		baseGenerator: &baseGenerator{
			engine:       engine,
			isProposal:   isProposal,
			now:          time.Now(),
			signatureID:  0,
			networkNode:  node,
			worksForNode: worksForNode,
			policy:       policy,
		},
		dtoCatnat: dto,
	}
}
func (el *CatnatGenerator) Generate() {
	el.woptaFooter()
	el.addMainHeader()
	el.engine.NewPage()
	el.engine.NewLine(constants.CellHeight)
	if el.isProposal {
		el.generateProposal()
	} else {
		el.generateContract()
	}
}
func (el *CatnatGenerator) generateProposal() {
	el.engine.NewLine(3)
	el.addInitialGreatings()
	el.engine.NewLine(3)
	catnat.AddBuildingInformation(el.engine, el.dtoCatnat.SedeDaAssicurare, el.dtoCatnat.Questions)
	el.engine.NewLine(3)
	catnat.AddTableGuarantee(el.engine, el.dtoCatnat.Guarantee)
	el.engine.NewLine(2)
	el.addPolicyInformationSection()
	el.engine.NewLine(2)
	el.addLinks()
	el.engine.NewLine(2)
	el.addFinalGreatings()
	el.engine.NewPage()
	el.addStatements()
	el.whoWeAre()
	el.AddMup()
	el.woptaPrivacySection()
}

func (el *CatnatGenerator) addFinalGreatings() {
	el.engine.WriteText(el.engine.GetTableCell("Cordiali saluti."))
	el.engine.NewLine(2)
	el.engine.WriteText(el.engine.GetTableCell("Anna di Wopta Assicurazioni\nProteggiamo chi sei"))
}

// TODO: insert rght link for catnat form
func (el *CatnatGenerator) addLinks() {
	if el.policy.Channel == models.ECommerceChannel || true {
		el.engine.WriteTexts(el.engine.GetTableCell("Qualora avessi necessità di ulteriori informazioni e supporto, la risposta che stai cercando potrebbe essere tra le nostre FAQ che trovi a questo "), domain.TableCell{Text: "link", Link: models.FAQ})
		el.engine.WriteTexts(el.engine.GetTableCell("Se non la trovi, puoi contattarci attraverso questo "), domain.TableCell{Text: "Form", Link: "forma.com"}, el.engine.GetTableCell(" oppure scriverci a questa "), domain.TableCell{Text: "e-mail", Link: "mailto:info@wopta.it"}, el.engine.GetTableCell(".In entrambi i casi un nostro esperto si prenderà cura della tua richiesta. "))
	} else if el.policy.Channel == models.NetworkChannel {
		el.engine.WriteText(el.engine.GetTableCell("Se hai necessità di ulteriori informazioni e supporto, contatta il tuo intermediario con il quale hai realizzato la polizza. "))
	}
}

func (el *CatnatGenerator) addPolicyInformationSection() {
	if el.dtoCatnat.Price.Consultancy.ValueFloat == 0 {
		return
	}
	el.engine.NewLine(constants.CellHeight)
	text :=
		"Infine, ti ricordiamo la presente polizza prevede il pagamento dei seguenti costi:\n" +
			fmt.Sprintf("- Premio di polizza: euro %v con frazionamento %v\n", el.dtoCatnat.Price.Gross.Text, el.dtoCatnat.Price.Split) +
			fmt.Sprintf("- Contributo servizi di intermediazione annuale: euro %v corrisposti con il pagamento della prima rata di polizza\n", el.dtoCatnat.Price.Consultancy.Text) +
			fmt.Sprintf("- Per un totale annuo di euro %v", el.dtoCatnat.Price.Total.Text)

	el.engine.WriteText(el.engine.GetTableCell(text, constants.BlackColor))
}
func (el *CatnatGenerator) addInitialGreatings() {
	el.engine.WriteText(el.engine.GetTableCell("Documento informativo di Wopta Assicurazioni relativo alla polizza Wopta per te Catastrofali Azienda", constants.PinkColor, constants.BoldFontStyle))
	el.engine.NewLine(3)
	el.engine.WriteText(el.engine.GetTableCell("Buongiorno " + el.dtoCatnat.Contractor.Name + " " + el.dtoCatnat.Contractor.Surname + ","))
	el.engine.WriteText(el.engine.GetTableCell("Confermiamo emissione della polizza in oggetto con i seguenti dati."))
	el.engine.WriteText(el.engine.GetTableCell("Verifica i dati e le informazioni qui riportate."))
	el.engine.WriteText(el.engine.GetTableCell("Riceverai mail per la firma della polizza e, una volta effettuato anche il pagamento, questa risulterà valida ed efficace."))
	el.engine.NewLine(2)
	el.engine.WriteText(el.engine.GetTableCell("Nel seguito sono riportate le informazioni, le dichiarazioni ed i consensi che ti saranno richiesti e che saranno poi da firmare con i documenti che provvederemo ad inviarti."))
	el.engine.NewLine(2)
	el.engine.WriteText(el.engine.GetTableCell("Il presente documento vale a titolo informativo precontrattuale e non ha validità ai fini contrattuali.", constants.BoldFontStyle))
	el.engine.NewLine(2)
}
func (el *CatnatGenerator) generateContract() {
	el.addContractorInformation()
	el.engine.NewLine(constants.CellHeight)
	el.addStatements()
	el.addAttachmentsInformation()
	el.AddMup()
	el.woptaPrivacySection()
	el.addElectronicSignPolicy()
	el.addOtpSignPolicy()
}
func (el *CatnatGenerator) addStatements() {
	if el.policy.Statements == nil {
		return
	}
	for _, statement := range *el.policy.Statements {
		el.printStatement(statement)
	}
}
func (el *CatnatGenerator) addMainHeader() {
	var (
		firstColumnWidth  float64 = 115
		secondColumnWidth         = constants.FullPageWidth - firstColumnWidth
	)
	parseLogos := func(texts []string) [][]domain.TableCell {
		result := make([][]domain.TableCell, 0, len(texts))

		result = append(result, []domain.TableCell{
			{
				Text:      "",
				Height:    constants.CellHeight,
				Width:     firstColumnWidth,
				FontStyle: constants.RegularFontStyle,
				FontColor: constants.BlackColor,
				Fill:      false,
				FontSize:  constants.RegularFontSize,
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
				FontStyle: constants.BoldFontStyle,
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

	rowsData := [][]string{
		{"I dati del tuo Polizza", "I tuoi dati"},
		{"Numero: " + fmt.Sprint(el.policy.NumberCompany), "Contraente: " + el.dtoCatnat.Contractor.Name + " " + el.dtoCatnat.Contractor.Surname},
		{"Decorre dal: " + el.dtoCatnat.ValidityDate.StartDate, "C.F./P.IVA: " + el.dtoCatnat.Contractor.FiscalCode_VatCode},
		{"Scade il: " + el.dtoCatnat.ValidityDate.EndDate, "Sede Legale: " + strings.ReplaceAll(el.dtoCatnat.Contractor.Address, "\n", "")},
		{"Si rinnova a scadenza, salvo disdetta da inviare 30 giorni prima", "Sede Assicurata: " + strings.ReplaceAll(el.dtoCatnat.SedeDaAssicurare.Address, "\n", "")},
		{"Produttore: Michele Lomazzi", " "},
	}
	el.engine.SetHeader(func() {
		firstColumnWidth = 15
		secondColumnWidth = constants.FullPageWidth - firstColumnWidth
		el.engine.DrawTable(parseLogos([]string{"				Wopta per te", "Catastrofali Azienda"}))
		el.engine.InsertImage(lib.GetAssetPathByEnvV2()+"logo_catnat.png", 10, 15, 13, 13)
		el.engine.NewLine(4)
		el.engine.InsertImage(lib.GetAssetPathByEnvV2()+"logo_wopta.png", 165, 15, 35, 10)
		firstColumnWidth = 115
		secondColumnWidth = constants.FullPageWidth - firstColumnWidth
		el.engine.DrawTable(parseData(rowsData))
	})
}

func (el *CatnatGenerator) addContractorInformation() {
	el.engine.WriteText(el.engine.GetTableCell("In relazione alla polizza sopra meglio identificata, il contraente, nella figura del suo rappresentante legale:"))
	const (
		firstColumnWidth  float64 = 30
		secondColumnWidth         = constants.FullPageWidth - firstColumnWidth
	)
	parseData := func(rows [][]string) [][]domain.TableCell {
		result := make([][]domain.TableCell, 0, len(rows))

		for _, row := range rows {
			fontStyle := constants.RegularFontStyle
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
	rowsData := [][]string{
		{"cognome e nome: ", el.dtoCatnat.Contractor.Name + " " + el.dtoCatnat.Contractor.Surname},
		{"codice fiscale: ", el.dtoCatnat.Contractor.FiscalCode},
		{"ruolo", "xxxxx"},
	}
	el.engine.DrawTable(parseData(rowsData))
	el.engine.WriteText(domain.TableCell{
		Text:  "dichiara quanto segue",
		Align: constants.CenterAlign,
		Width: constants.FullPageWidth,
	})
	el.engine.NewLine(2)
}

func (el *CatnatGenerator) addAttachmentsInformation() {
	el.engine.WriteText(el.engine.GetTableCell("Il presente documento include:"))
	el.engine.WriteText(el.engine.GetTableCell("- ALLEGATO 3 Modulo Unico Precontrattuale per i prodotti assicurativi, così come previsto dal regolamento 40/2018 IVASS"))
	el.engine.WriteText(el.engine.GetTableCell("- Informativa privacy dell’intermediario Wopta Assicurazioni Srl"))
	el.engine.WriteText(el.engine.GetTableCell("- \"Condizioni Generali di Servizio per l'utilizzazione della Firma Elettronica Avanzata\" e l'annessa \"Scheda Tecnica Illustrativa\""))
}
