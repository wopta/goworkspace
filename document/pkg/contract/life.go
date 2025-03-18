package contract

import (
	"fmt"
	"time"

	"github.com/wopta/goworkspace/document/internal/constants"
	"github.com/wopta/goworkspace/document/internal/domain"
	"github.com/wopta/goworkspace/document/internal/engine"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

type LifeGenerator struct {
	*baseGenerator
}

func NewLifeGenerator(engine *engine.Fpdf, policy *models.Policy, node *models.NetworkNode, product models.Product, isProposal bool) *LifeGenerator {
	return &LifeGenerator{
		baseGenerator: &baseGenerator{
			engine:      engine,
			isProposal:  isProposal,
			now:         time.Now(),
			signatureID: 0,
			networkNode: node,
			policy:      policy,
		},
	}
}

func (el *LifeGenerator) Generate() ([]byte, error) {
	el.addMainHeader()
	el.engine.NewPage()
	el.engine.WriteText(getTableCell("\n\nIl tuo Preventivo: cosa fare adesso?\n",constants.BoldFontStyle,constants.PinkColor,constants.LargeFontSize))
	el.addWelcomeSection()
	el.addEmailSection()
	el.addSignSection()
	el.addPolicyInformationSection()
	el.addSupportInformationSection()
	el.addGreatingsSection()
	return el.engine.RawDoc()
}

func (el *LifeGenerator) addMainHeader() {
	policy := el.policy
	const (
		firstColumnWidth  = 115
		secondColumnWidth = 75
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

	formatDate := func(t time.Time) string {
		location, _ := time.LoadLocation("Europe/Rome")
		time := t.In(location)
		return time.In(location).Format("02/01/2006")
	}

	addressFirstPart:=policy.Contractor.Residence.StreetName + ", " + policy.Contractor.Residence.StreetNumber
	addressSecondPart:=policy.Contractor.Residence.PostalCode + " " + policy.Contractor.Residence.City + " (" + policy.Contractor.Residence.CityCode + ")"
	rowsData := [][]string{
		{"I dati del tuo Preventivo", "I tuoi dati"},
		{fmt.Sprintf("Numero: %d", policy.Number), "Contraente: " + policy.Name + " " + policy.Contractor.Surname},
		{"Decore dal: " + formatDate(policy.EmitDate), "C.F./P.IVA: " + policy.Contractor.FiscalCode}, 
		{"Scade il: " + formatDate(policy.EndDate), "Indirizzo: " + addressFirstPart},
		{"Prima scadenza Annuale il: " + formatDate(policy.EndDate),addressSecondPart},
		{"Non si rinnova a scadenza.", "Mail: " + policy.Contractor.Mail},
		{fmt.Sprintf("Produttore: %s %s", policy.Contractor.Name, policy.Contractor.Surname), "Telefono: " + policy.Contractor.Phone},
	}

	el.engine.SetHeader(func() {
		el.engine.InsertImage(lib.GetAssetPathByEnvV2()+"logo_wopta.png", 10, 5, 35, 12)
		el.baseGenerator.engine.NewLine(3)
		el.engine.DrawTable(parseLogos([]string{"Wopta per te", "Vita"}))
		el.engine.InsertImage(lib.GetAssetPathByEnvV2()+"logo_vita.png", 180, 15, 13, 13)
		el.engine.DrawTable(parseData(rowsData))

		if el.isProposal {
			el.engine.DrawWatermark(constants.Proposal)
		}
	})
}

func (el *LifeGenerator) addWelcomeSection(){
	el.writeTexts(
		getTableCell(fmt.Sprintf("\n\nBuongiorno %v %v,\n", el.policy.Contractor.Name, el.policy.Contractor.Surname), constants.BlackColor),
		getTableCell("\nGrazie per aver fatto un preventivo per una polizza Vita, dimostrando volontà e interesse a tutelarti e/o proteggere le persone per te più importanti.\n", constants.BlackColor),
	)
}

func (el *LifeGenerator) addEmailSection(){
	el.engine.WriteText(getTableCell(
		"\nIn allegato trovi:\n"+
		"- modulo di Polizza\n"+
		"- informativa precontrattuale di Wopta, prevista per legge\n"+
		"- modulistica antiriciclaggio\n\n"+
		"- informativa e dichiarazioni privacy per l’Assicuratore\n"+
		"- informativa e dichiarazioni privacy per l’Intermediario\n", constants.BlackColor,
	))

	el.engine.WriteText(getTableCell("Verifica la correttezza di tutti i dati inseriti (anagrafici, indirizzi, codice fiscale, contatti) e delle prestazioni scelte (durata, importi, eventuali opzioni).\n", constants.BlackColor))

	if el.policy.Channel == models.ECommerceChannel {
		el.writeTexts(
			getTableCell("Riceverai anche due mail per procedere con la ", constants.BlackColor),
			getTableCell("firma ", constants.PinkColor,constants.BoldFontStyle),
			getTableCell("ed il ", constants.BlackColor),
			getTableCell("pagamento\n", constants.PinkColor,constants.BoldFontStyle),
		)
	}
}

func (el *LifeGenerator) addSignSection(){
	el.writeTexts(
		getTableCell("ATTENZIONE", constants.PinkColor,constants.BoldFontStyle),
		getTableCell(" :Solo una volta firmati i documenti ed effettuato il pagamento, la copertura assicurativa sarà attiva e così ti invieremo i documenti contrattuali da te firmati, che poi potrai visualizzare nell’area riservata ai clienti della nostra app e/o sito.", constants.BlackColor),
	)
}

func (el *LifeGenerator) addPolicyInformationSection(){
	getSplitLabel:=func (paymentSplit string)string{
		switch paymentSplit{
		case string(models.PaySplitMonthly):
			return "mensile"
		case string(models.PaySplitYearly):
			return "annuale"
		case string(models.PaySplitSingleInstallment):
			return "singolo"
		}
		return ""
	}

	if el.policy.ConsultancyValue.Price==0  {
		return
	}
	el.engine.NewLine(constants.CellHeight*2)
	text:=
		"Infine, ti ricordiamo la presente polizza prevede il pagamento dei seguenti costi:\n"+
		fmt.Sprintf("Premio di polizza: euro %v con frazionamento %v\n", lib.HumanaizePriceEuro(el.policy.PriceGross),getSplitLabel(el.policy.PaymentSplit))+
		fmt.Sprintf("- Contributo per servizi di intermediazione: euro %v vorrisposti con il pagamento della prima rata di polizza. Il documento contabile è scaricabile dall’app o nella tua area riservata\n", lib.HumanaizePriceEuro(el.policy.ConsultancyValue.Price))+
		fmt.Sprintf("- Per un totale annuo di euro %v", el.policy.PaymentComponents.PriceAnnuity.Total)

	el.engine.WriteText(getTableCell(text,constants.BlackColor))
}

func (el *LifeGenerator) addSupportInformationSection(){
	if el.policy.Channel == models.ECommerceChannel{
		text := "\nRestiamo a disposizione per ogni ulteriore informazione anche attraverso i canali di contatto che trovi a questo "
		el.writeTexts(
			getTableCell(text, constants.BlackColor),
			getTableCell("link", constants.PinkColor),
		)
		widthLink:=el.engine.GetStringWidth("link")
		el.engine.GetPdf().LinkString(1+el.engine.GetX()-widthLink, el.engine.GetY(), widthLink, constants.CellHeight, "https://www.wopta.it/it/vita/#contact-us")
		return
	}

	if el.policy.Channel == models.AgencyChannel {
		el.engine.WriteText(getTableCell("Se hai necessità di ulteriori informazioni e supporto, rivolgiti al tuo intermediario, che trovi in copia conoscenza alla mail accompagnatoria di questa comunicazione.", constants.BlackColor))
	}
}

func (el *LifeGenerator) addGreatingsSection(){
	el.engine.NewLine(constants.CellHeight)
	el.engine.WriteText(getTableCell("Cordiali saluti.", constants.BlackColor))
	el.engine.NewLine(constants.CellHeight*2)
	el.writeTexts(
		getTableCell("Anna di Wopta Assicurazioni\n", constants.BlackColor),
		getTableCell("Proteggiamo chi sei", constants.BlackColor),
	)
}

//get a tablecell personalized based on passed opts
func getTableCell(text string, opts ...any)domain.TableCell{
	tableCell:=domain.TableCell{}
	tableCell.Text=text
	tableCell.Height=constants.CellHeight
	tableCell.Align=constants.LeftAlign
	tableCell.FontStyle=constants.RegularFontStyle
	tableCell.FontSize=constants.RegularFontSize

	for _,opt := range opts{
		switch opt:=opt.(type){
		case domain.FontSize:
			tableCell.FontSize=opt
		case domain.FontStyle:
			tableCell.FontStyle=opt
		case domain.Color:
			tableCell.FontColor=opt
		}
	}
	return tableCell
}

func (el *LifeGenerator) writeTexts(tables ...domain.TableCell) {
	for _, text := range tables {
		el.engine.RawWriteText(text)
	}
}
