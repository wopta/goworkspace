package contract

import (
	"fmt"
	"strings"
	"time"

	"github.com/wopta/goworkspace/document/internal/constants"
	"github.com/wopta/goworkspace/document/internal/domain"
	"github.com/wopta/goworkspace/document/internal/dto"
	"github.com/wopta/goworkspace/document/internal/engine"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

const (
	fabbricatoGuaranteeSlug             string = "fabbricato"
	rischioLocativoGuaranteeSlug        string = "rischio locativo"
	macchinariGuaranteeSlug             string = "macchinari"
	merciGuaranteeSlug                  string = "merci"
	merciAumentoGuaranteeSlug           string = "merci in aumento"
	merciAumentoEnterpriseGuaranteeSlug string = "merci in aumento-enterprise"
	ricorsoGuaranteeSlug                string = "ricorso"
	fenomenoGuaranteeSlug               string = "fenomeno"
	merciRefrigerazioneGuaranteeSlug    string = "merci-refri"
	guastiGuaranteeSlug                 string = "guasti"
	elettronicaGuaranteeSlug            string = "elettronica"
	furtoGuaranteeSlug                  string = "furto"
	danniGuaranteeSlug                  string = "danni"
	diariaGuaranteeSlug                 string = "diaria"
	maggioriGuaranteeSlug               string = "maggiori"
	pigioniGuaranteeSlug                string = "pigioni"
	rctoGuaranteeSlug                   string = "rcto"
	rcpGuaranteeSlug                    string = "rcp"
	managementOrganizationGuaranteeSlug string = "management-organization"
	cyberGuanrateeSlug                  string = "cyber"
	rctGuaranteeSlug                    string = "rct"
	ritiroGuaranteeSlug                 string = "ritiro"
)

var (
	guaranteeNamesMap = map[string]string{
		fabbricatoGuaranteeSlug:             "Fabbricato",
		rischioLocativoGuaranteeSlug:        "Rischio Locativo (in aumento A/24)",
		macchinariGuaranteeSlug:             "Macchinari",
		merciGuaranteeSlug:                  "Merci (importi fissi)",
		merciAumentoGuaranteeSlug:           "Merci (Aumento temporaneo A/29)",
		merciAumentoEnterpriseGuaranteeSlug: "Merci (Aumento temporaneo A/29) - giorni",
		ricorsoGuaranteeSlug:                "Ricordo Terzi (in aumento A/25)",
		fenomenoGuaranteeSlug:               "Fenomeno Elettrico (in aumento A/23)",
		merciRefrigerazioneGuaranteeSlug:    "Merci in refrigerazione",
		guastiGuaranteeSlug:                 "Guasti alle macchine (in aumento A/27)",
		elettronicaGuaranteeSlug:            "Apparecch.re Elettroniche (in aumento A/26)",
		furtoGuaranteeSlug:                  "Furto, rapina, estorsione (in aumento C/1)",
		danniGuaranteeSlug:                  "Danni indiretti - Formula",
		diariaGuaranteeSlug:                 "Diaria Giornaliera",
		maggioriGuaranteeSlug:               "Maggiori costi",
		pigioniGuaranteeSlug:                "Perdita Pigioni",
		rctoGuaranteeSlug: "Responsabilità Civile verso Terzi (" +
			"RCT) e verso Prestatori di lavoro (RCO)",
		rcpGuaranteeSlug:                    "Responsabilità Civile Prodotti (RCP) Ritiro prodotti",
		managementOrganizationGuaranteeSlug: "Responsabilità Amministratori Sindaci Dirigenti (D&O)",
		cyberGuanrateeSlug:                  "Cyber Response e Data Security",
	}
)

type CommercialCombinedGenerator struct {
	*baseGenerator
	dto *dto.CommercialCombinedDTO
}

func NewCommercialCombinedGenerator(engine *engine.Fpdf, policy *models.Policy, node *models.NetworkNode,
	isProposal bool) *CommercialCombinedGenerator {
	commercialCombinedDTO := dto.NewCommercialCombinedDto()
	commercialCombinedDTO.FromPolicy(*policy, isProposal)
	return &CommercialCombinedGenerator{
		baseGenerator: &baseGenerator{
			engine:      engine,
			isProposal:  isProposal,
			now:         time.Now(),
			signatureID: 0,
			networkNode: node,
			policy:      policy,
		},
		dto: commercialCombinedDTO,
	}
}

func (ccg *CommercialCombinedGenerator) Contract() ([]byte, error) {
	ccg.mainHeader()

	ccg.engine.NewPage()

	ccg.mainFooter()

	ccg.engine.NewLine(10)

	ccg.introSection()

	ccg.engine.NewLine(10)

	ccg.whoWeAreSection()

	ccg.engine.NewLine(10)

	ccg.insuredDetailsSection()

	ccg.engine.NewPage()

	ccg.guaranteesDetailsSection()

	ccg.engine.NewPage()

	ccg.deductibleSection()

	ccg.engine.NewLine(10)

	ccg.dynamicDeductibleSection()

	ccg.engine.NewLine(5)

	ccg.detailsSection()

	ccg.engine.NewLine(5)

	ccg.specialConditionsSection()

	ccg.engine.NewLine(5)

	ccg.bondSection()

	ccg.engine.NewPage()

	ccg.resumeSection()

	ccg.engine.NewLine(5)

	ccg.howYouCanPaySection()

	ccg.engine.NewLine(5)

	ccg.emitResumeSection()

	ccg.engine.NewLine(5)

	ccg.statementsFirstPart()

	ccg.engine.NewPage()

	ccg.claimsStatement()

	ccg.statementsSecondPart()

	ccg.engine.NewLine(3)

	ccg.qbePrivacySection()

	ccg.engine.NewPage()

	ccg.qbePersonalDataSection()

	ccg.engine.NewLine(5)

	ccg.commercialConsentSection()

	ccg.annexSections()

	ccg.woptaHeader()

	ccg.woptaFooter()

	ccg.woptaPrivacySection()

	return ccg.engine.RawDoc()
}

func (ccg *CommercialCombinedGenerator) mainHeader() {
	const (
		firstColumnWidth  = 115
		secondColumnWidth = 75
	)

	contractDTO := ccg.dto.ContractDTO
	contractorDTO := ccg.dto.ContractorDTO

	parser := func(rows [][]string) [][]domain.TableCell {
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
					FontSize:  constants.RegularFontSize,
					Fill:      false,
					FillColor: domain.Color{},
					Align:     constants.LeftAlign,
					Border:    "",
				},
				{
					Text:      row[1],
					Height:    constants.CellHeight,
					Width:     secondColumnWidth,
					FontStyle: fontStyle,
					FontColor: constants.BlackColor,
					FontSize:  constants.RegularFontSize,
					Fill:      false,
					FillColor: domain.Color{},
					Align:     constants.LeftAlign,
					Border:    "",
				},
			})
		}
		return result
	}

	address := strings.Split(fmt.Sprintf("%s %s\n%s %s (%s)", contractorDTO.StreetName,
		contractorDTO.StreetNumber, contractorDTO.PostalCode,
		contractorDTO.City, contractorDTO.CityCode), "\n")

	rows := [][]string{
		{contractDTO.CodeHeading + " " + contractDTO.Code, "I tuoi dati"},
		{"Decorre dal: " + contractDTO.StartDate + " ore 24:00",
			"Contraente: " + contractorDTO.Name + " " + contractorDTO.Surname},
		{
			"Scade il: " + contractDTO.EndDate + " ore 24:00", "P.IVA: " + contractorDTO.VatCode},
		{"Si rinnova a scadenza, salvo disdetta da inviare 30 giorni prima", "Codice Fiscale: " + contractorDTO.FiscalCode},
		{"Frazionamento: " + contractDTO.PaymentSplit, address[0]},
		{"Prossimo pagamento il: " + contractDTO.NextPay, address[1]},
		{"Sostituisce la Polizza: ======", "Mail: " + contractorDTO.Mail},
		{"Presenza Vincolo: " + contractDTO.HasBond + " Convenzione: NO", "Telefono: " + contractorDTO.Phone},
	}

	table := parser(rows)

	ccg.engine.SetHeader(func() {
		ccg.engine.InsertImage(lib.GetAssetPathByEnvV2()+"logo_qbe.png", 75, 6.5, 22, 8)
		ccg.engine.DrawLine(102, 6.25, 102, 15, 0.25, constants.BlackColor)
		ccg.engine.InsertImage(lib.GetAssetPathByEnvV2()+"logo_wopta.png", 107.5, 5, 35, 12)
		ccg.engine.NewLine(7)
		ccg.engine.DrawTable(table)

		if ccg.isProposal {
			ccg.engine.DrawWatermark(constants.Proposal)
		}
	})
}

func (ccg *CommercialCombinedGenerator) mainFooter() {
	text := "QBE Europe SA/NV, Rappresentanza Generale per l’Italia, Via Melchiorre Gioia 8 – 20124 Milano. R.E.A. MI-2538674. Codice fiscale/P.IVA 10532190963 Autorizzazione IVASS n. I.00147\n" +
		"QBE Europe SA/NV è autorizzata dalla Banca Nazionale del Belgio con licenza numero 3093. Sede legale Place du Champ de Mars 5, BE 1050, Bruxelles, Belgio.   N. di registrazione 0690537456."

	ccg.engine.SetFooter(func() {
		ccg.engine.SetX(10)
		ccg.engine.SetY(-17)
		ccg.engine.WriteText(domain.TableCell{
			Text:      text,
			Height:    3,
			Width:     190,
			FontStyle: constants.RegularFontStyle,
			FontColor: constants.BlackColor,
			FontSize:  constants.SmallFontSize,
			Fill:      false,
			FillColor: domain.Color{},
			Align:     constants.LeftAlign,
			Border:    "",
		})
		ccg.engine.WriteText(domain.TableCell{
			Text:      fmt.Sprintf("%d", ccg.engine.PageNumber()),
			Height:    3,
			Width:     0,
			FontStyle: constants.RegularFontStyle,
			Fill:      false,
			FillColor: domain.Color{},
			Align:     constants.RightAlign,
			Border:    "",
		})
	})
}

func (ccg *CommercialCombinedGenerator) introSection() {
	introTable := [][]domain.TableCell{
		{
			{
				Text:      " \nScheda di Polizza\n ", // TODO: find better solution
				Height:    3.5,
				Width:     95,
				FontStyle: constants.BoldFontStyle,
				FontColor: constants.BlackColor,
				FontSize:  constants.RegularFontSize,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.CenterAlign,
				Border:    "TB",
			},
			{
				Text:      "COMMERCIAL COMBINED\nAssicurazione Multigaranzia per le imprese\nSet informativo - Edizione 10/2022",
				Height:    3.5,
				Width:     95,
				FontStyle: constants.BoldFontStyle,
				FontColor: constants.BlackColor,
				FontSize:  constants.RegularFontSize,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.CenterAlign,
				Border:    "TB",
			},
		},
	}
	ccg.engine.DrawTable(introTable)
}

func (ccg *CommercialCombinedGenerator) whoWeAreSection() {
	whoWeAreTable := [][]domain.TableCell{
		{
			{
				Text:      "Chi siamo",
				Height:    5,
				Width:     190,
				FontSize:  constants.LargeFontSize,
				FontStyle: constants.BoldFontStyle,
				FontColor: constants.BlackColor,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "",
			},
		},
		{
			{
				Text:      "QBE Europe SA/NV Rappresentanza generale per l’Italia,",
				Height:    3.5,
				Width:     93.5,
				FontSize:  constants.RegularFontSize,
				FontStyle: constants.BoldFontStyle,
				FontColor: constants.BlackColor,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "",
			},
			{
				Text:      "impresa di assicurazione operante in Italia in regime di",
				Height:    3.5,
				Width:     96.5,
				FontSize:  constants.RegularFontSize,
				FontStyle: constants.RegularFontStyle,
				FontColor: constants.BlackColor,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "",
			},
		},
		{
			{
				Text:      "libertà di stabilimento, autorizzata dalla Banca Nazionale del Belgio con licenza numero 3093, con sede legale in Place du Champ de Mars 5, BE 1050, Bruxelles, Belgio e sede secondaria in Italia, Via Melchiorre Gioia, 8, 20124, Milano (MI), R.E.A. MI-2538674, codice fiscale e p. iva 10532190963, Autorizzazione IVASS n. I.00147.",
				Height:    3.5,
				Width:     190,
				FontSize:  constants.RegularFontSize,
				FontStyle: constants.RegularFontStyle,
				FontColor: constants.BlackColor,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "",
			},
		},
		{
			{
				Text:      " ",
				Height:    5,
				Width:     190,
				FontSize:  constants.RegularFontSize,
				FontStyle: constants.RegularFontStyle,
				FontColor: constants.BlackColor,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     "",
				Border:    "",
			},
		},
		{
			{
				Text:      "Wopta Assicurazioni S.r.l.",
				Height:    3.5,
				Width:     45,
				FontSize:  constants.RegularFontSize,
				FontStyle: constants.BoldFontStyle,
				FontColor: constants.PinkColor,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "",
			},
			{
				Text:      "(nel testo anche “Wopta”) - intermediario assicurativo, soggetto al controllo dell’IVASS ed",
				Height:    3.5,
				Width:     145,
				FontSize:  constants.RegularFontSize,
				FontStyle: constants.RegularFontStyle,
				FontColor: constants.BlackColor,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "",
			},
		},
		{
			{
				Text:      "iscritto dal 14.02.2022 al Registro Unico degli Intermediari, in Sezione A nr. A000701923, avente sede legale in Galleria del Corso, 1 – 20122 Milano (MI). Capitale sociale euro 120.000 - Codice Fiscale, Reg. Imprese e Partita IVA: 12072020964 - Iscritta al Registro delle imprese di Milano – REA MI 2638708",
				Height:    3.5,
				Width:     190,
				FontSize:  constants.RegularFontSize,
				FontStyle: constants.RegularFontStyle,
				FontColor: constants.BlackColor,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "",
			},
		},
		{
			{
				Text:      " ",
				Height:    5,
				Width:     190,
				FontSize:  constants.RegularFontSize,
				FontStyle: constants.RegularFontStyle,
				FontColor: constants.BlackColor,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     "",
				Border:    "",
			},
		},
		{
			{
				Text:      "Commercial Combined ",
				Height:    3.5,
				Width:     40,
				FontSize:  constants.RegularFontSize,
				FontStyle: constants.BoldFontStyle,
				FontColor: constants.BlackColor,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "",
			},
			{
				Text:      "Assicurazione multigaranzia per le imprese è un prodotto assicurativo di QBE Europe SA/NV",
				Height:    3.5,
				Width:     150,
				FontSize:  constants.RegularFontSize,
				FontStyle: constants.RegularFontStyle,
				FontColor: constants.BlackColor,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "",
			},
		},
		{
			{
				Text:      "Rappresentanza Generale per l’Italia distribuito da Wopta Assicurazioni S.r.l.",
				Height:    3.5,
				Width:     190,
				FontSize:  constants.RegularFontSize,
				FontStyle: constants.RegularFontStyle,
				FontColor: constants.BlackColor,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "",
			},
		},
	}
	ccg.engine.DrawTable(whoWeAreTable)
}

func (ccg *CommercialCombinedGenerator) insuredDetailsSection() {
	buildings := ccg.dto.BuildingsDTO
	enterprise := ccg.dto.EnterpriseDTO

	table := make([][]domain.TableCell, 0)

	titleRow := []domain.TableCell{
		{
			Text:      "L'assicurazione è prestata per",
			Height:    3.5,
			Width:     190,
			FontSize:  constants.LargeFontSize,
			FontStyle: constants.BoldFontStyle,
			FontColor: constants.BlackColor,
			Fill:      false,
			FillColor: domain.Color{},
			Align:     constants.LeftAlign,
			Border:    "",
		},
	}
	table = append(table, titleRow)

	for i := 0; i < 5; i++ {
		building := buildings[i]
		border := "B"
		if i == 0 {
			border = "TB"
		}

		address := fmt.Sprintf("%s, %s - %s %s (%s)", building.StreetName,
			building.StreetNumber, building.PostalCode, building.City, building.CityCode)

		row := []domain.TableCell{
			{
				Text:      fmt.Sprintf(" \nSede %d\n ", i+1),
				Height:    4.5,
				Width:     40,
				FontSize:  constants.RegularFontSize,
				FontStyle: constants.BoldFontStyle,
				FontColor: constants.BlackColor,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.CenterAlign,
				Border:    border,
			},
			{
				Text: fmt.Sprintf("Indirizzo: %s\nFabbricato in %s, "+
					"pannelli sandwich: %s; Allarme antifurto: %s, Sprinkler: %s,"+
					"\nAttività NAICS codice: %s Descrizione: %s",
					address, building.BuildingMaterial, building.HasSandwichPanel, building.HasAlarm,
					building.HasSprinkler, building.Naics, building.NaicsDetail),
				Height:    4.5,
				Width:     150,
				FontSize:  constants.RegularFontSize,
				FontStyle: constants.RegularFontStyle,
				FontColor: constants.BlackColor,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    border,
			},
		}
		table = append(table, row)
	}

	limitOfIndemnity := "======"
	sumInsured := "======"

	if val, ok := enterprise.Guarantees[managementOrganizationGuaranteeSlug]; ok {
		if val.LimitOfIndemnity != 0 {
			limitOfIndemnity = lib.HumanaizePriceEuro(val.LimitOfIndemnity)
		}
		if val.SumInsured != 0 {
			sumInsured = lib.HumanaizePriceEuro(val.SumInsured)
		}
	}

	enterpriseRow := []domain.TableCell{
		{
			Text:      " \nAttività\n ",
			Height:    4.5,
			Width:     40,
			FontSize:  constants.RegularFontSize,
			FontStyle: constants.BoldFontStyle,
			FontColor: constants.BlackColor,
			Fill:      false,
			FillColor: domain.Color{},
			Align:     constants.CenterAlign,
			Border:    "TB",
		},
		{
			Text: fmt.Sprintf("Fatturato: %s di cui verso USA e Canada: %s\nPrestatori di lavoro nr: %d"+
				" - Retribuzioni: %s\nTotal Asset: %s di cui capitale proprio: %s",
				lib.HumanaizePriceEuro(enterprise.Revenue), lib.HumanaizePriceEuro(enterprise.NorthAmericanMarket),
				enterprise.Employer, lib.HumanaizePriceEuro(enterprise.WorkEmployersRemuneration),
				limitOfIndemnity, sumInsured),
			Height:    4.5,
			Width:     150,
			FontSize:  constants.RegularFontSize,
			FontStyle: constants.RegularFontStyle,
			FontColor: constants.BlackColor,
			Fill:      false,
			FillColor: domain.Color{},
			Align:     constants.LeftAlign,
			Border:    "TB",
		},
	}
	table = append(table, enterpriseRow)

	riskDescriptionRow := []domain.TableCell{
		{
			Text:      " \n \nDescrizione del rischio\n \n ",
			Height:    4.5,
			Width:     40,
			FontSize:  constants.RegularFontSize,
			FontStyle: constants.BoldFontStyle,
			FontColor: constants.BlackColor,
			Fill:      false,
			FillColor: domain.Color{},
			Align:     constants.CenterAlign,
			Border:    "TB",
		},
		{
			Text: "Stabilimento costituito da uno o più corpi di Fabbricati, " +
				"prevalentemente costruiti in materiali incombustibili, " +
				"nel quale i processi di lavorazione sono quelli che la tecnica inerente all’attività svolta insegna" +
				" e consiglia di usare o che l'Assicurato intende adottare. " +
				"S'intendono altresì compresi i depositi e tutte le dipendenze necessarie per la conduzione dell" +
				"'attività, incluse le abitazioni e le attività di carattere assistenziale e/o commerciale.",
			Height:    4.5,
			Width:     150,
			FontSize:  constants.RegularFontSize,
			FontStyle: constants.RegularFontStyle,
			FontColor: constants.BlackColor,
			Fill:      false,
			FillColor: domain.Color{},
			Align:     constants.LeftAlign,
			Border:    "TB",
		},
	}
	table = append(table, riskDescriptionRow)

	ccg.engine.DrawTable(table)
}

// TODO: parse policy info
func (ccg *CommercialCombinedGenerator) guaranteesDetailsSection() {
	const emptyInfo string = "======"

	type guaranteeInfo struct {
		sumInsuredLimitOfIndemnity string
		startDate                  string
		duration                   string
	}

	buildingsData := make(map[string][]string, 5)
	enterpriseData := make(map[string]guaranteeInfo)

	buildingsSlugs := []string{fabbricatoGuaranteeSlug, rischioLocativoGuaranteeSlug, macchinariGuaranteeSlug,
		merciGuaranteeSlug, merciAumentoGuaranteeSlug}

	enterpriseSlugs := []string{merciAumentoEnterpriseGuaranteeSlug, ricorsoGuaranteeSlug, fenomenoGuaranteeSlug,
		merciRefrigerazioneGuaranteeSlug,
		guastiGuaranteeSlug, elettronicaGuaranteeSlug, furtoGuaranteeSlug, danniGuaranteeSlug,
		diariaGuaranteeSlug, maggioriGuaranteeSlug,
		pigioniGuaranteeSlug, rctoGuaranteeSlug, rcpGuaranteeSlug, managementOrganizationGuaranteeSlug, cyberGuanrateeSlug}

	for _, slug := range buildingsSlugs {
		i := 0
		emptyData := make([]string, 5)
		for i < 5 {
			emptyData[i] = emptyInfo
			i++
		}

		buildingsData[slug] = emptyData
	}

	for _, slug := range enterpriseSlugs {
		enterpriseData[slug] = guaranteeInfo{
			sumInsuredLimitOfIndemnity: emptyInfo,
			startDate:                  emptyInfo,
			duration:                   emptyInfo,
		}
	}

	// TODO: fetch building data from policy

	// TODO: fetch enterprise data from policy

	table := make([][]domain.TableCell, 0)

	tableHeader := [][]domain.TableCell{
		{
			{
				Text:      "Le coperture assicurative che hai scelto (operative se indicata la Somma o il Massimale)",
				Height:    4.5,
				Width:     190,
				FontSize:  constants.LargeFontSize,
				FontStyle: constants.BoldFontStyle,
				FontColor: constants.BlackColor,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "",
			},
		},
		{
			{
				Text:      "PARTITE, MASSIMALI, SOMME ASSICURATE – DANNI DIRETTI – sezione A - C",
				Height:    4.5,
				Width:     190,
				FontSize:  constants.LargeFontSize,
				FontStyle: constants.BoldFontStyle,
				FontColor: constants.BlackColor,
				Fill:      true,
				FillColor: constants.GreyColor,
				Align:     constants.CenterAlign,
				Border:    "",
			},
		},
		{
			{
				Text:      "Valori in euro",
				Height:    4.5,
				Width:     65,
				FontSize:  constants.SmallFontSize,
				FontStyle: constants.RegularFontStyle,
				FontColor: constants.BlackColor,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "T",
			},
			{
				Text:      "Sede 1",
				Height:    4.5,
				Width:     25,
				FontSize:  constants.RegularFontSize,
				FontStyle: constants.BoldFontStyle,
				FontColor: constants.BlackColor,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.RightAlign,
				Border:    "T",
			},
			{
				Text:      "Sede 2",
				Height:    4.5,
				Width:     25,
				FontSize:  constants.RegularFontSize,
				FontStyle: constants.BoldFontStyle,
				FontColor: constants.BlackColor,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.RightAlign,
				Border:    "T",
			},
			{
				Text:      "Sede 3",
				Height:    4.5,
				Width:     25,
				FontSize:  constants.RegularFontSize,
				FontStyle: constants.BoldFontStyle,
				FontColor: constants.BlackColor,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.RightAlign,
				Border:    "T",
			},
			{
				Text:      "Sede 4",
				Height:    4.5,
				Width:     25,
				FontSize:  constants.RegularFontSize,
				FontStyle: constants.BoldFontStyle,
				FontColor: constants.BlackColor,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.RightAlign,
				Border:    "T",
			},
			{
				Text:      "Sede 5",
				Height:    4.5,
				Width:     25,
				FontSize:  constants.RegularFontSize,
				FontStyle: constants.BoldFontStyle,
				FontColor: constants.BlackColor,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.RightAlign,
				Border:    "T",
			},
		},
	}
	table = append(table, tableHeader...)

	for slug, data := range buildingsData {
		row := make([]domain.TableCell, 6)
		row[0] = domain.TableCell{
			Text:      guaranteeNamesMap[slug],
			Height:    4.5,
			Width:     65,
			FontSize:  constants.RegularFontSize,
			FontStyle: constants.BoldFontStyle,
			FontColor: constants.BlackColor,
			Fill:      false,
			FillColor: domain.Color{},
			Align:     constants.RightAlign,
			Border:    "T",
		}
		for _, building := range data {
			row = append(row, domain.TableCell{
				Text:      building,
				Height:    4.5,
				Width:     25,
				FontSize:  constants.RegularFontSize,
				FontStyle: constants.RegularFontStyle,
				FontColor: constants.BlackColor,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.RightAlign,
				Border:    "T",
			})
		}
		table = append(table, row)
	}

	for _, slug := range enterpriseSlugs[:7] {
		row := []domain.TableCell{
			{
				Text:      guaranteeNamesMap[slug],
				Height:    4.5,
				Width:     65,
				FontSize:  constants.RegularFontSize,
				FontStyle: constants.BoldFontStyle,
				FontColor: constants.BlackColor,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.RightAlign,
				Border:    "T",
			},
		}

		info := enterpriseData[slug].sumInsuredLimitOfIndemnity
		if slug == merciAumentoEnterpriseGuaranteeSlug {
			info += " a partire dal " + enterpriseData[slug].startDate + " di ogni anno"
		}

		row = append(row, domain.TableCell{
			Text:      info,
			Height:    4.5,
			Width:     125,
			FontSize:  constants.RegularFontSize,
			FontStyle: constants.RegularFontStyle,
			FontColor: constants.BlackColor,
			Fill:      false,
			FillColor: domain.Color{},
			Align:     constants.LeftAlign,
			Border:    "T",
		})
		table = append(table, row)
	}

	table = append(table, []domain.TableCell{
		{
			Text:      "GARANZIE E SOMME ASSICURATE – DANNI INDIRETTI – sezione B",
			Height:    4.5,
			Width:     190,
			FontSize:  constants.LargeFontSize,
			FontStyle: constants.BoldFontStyle,
			FontColor: constants.BlackColor,
			Fill:      true,
			FillColor: constants.GreyColor,
			Align:     constants.CenterAlign,
			Border:    "T",
		},
	})

	for _, slug := range enterpriseSlugs[7:11] {
		row := []domain.TableCell{
			{
				Text:      guaranteeNamesMap[slug],
				Height:    4.5,
				Width:     65,
				FontSize:  constants.RegularFontSize,
				FontStyle: constants.BoldFontStyle,
				FontColor: constants.BlackColor,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.RightAlign,
				Border:    "T",
			},
		}

		info := enterpriseData[slug].sumInsuredLimitOfIndemnity
		if slug == diariaGuaranteeSlug {
			info += " Periodo di indennizzo " + enterpriseData[slug].duration + " giorni"
		}

		row = append(row, domain.TableCell{
			Text:      info,
			Height:    4.5,
			Width:     125,
			FontSize:  constants.RegularFontSize,
			FontStyle: constants.RegularFontStyle,
			FontColor: constants.BlackColor,
			Fill:      false,
			FillColor: domain.Color{},
			Align:     constants.LeftAlign,
			Border:    "T",
		})
		table = append(table, row)
	}

	table = append(table, []domain.TableCell{
		{
			Text:      "GARANZIE E MASSIMALI RESPONSABILITA’ CIVILE E CYBER – sezioni D E F G H I",
			Height:    4.5,
			Width:     190,
			FontSize:  constants.LargeFontSize,
			FontStyle: constants.BoldFontStyle,
			FontColor: constants.BlackColor,
			Fill:      true,
			FillColor: constants.GreyColor,
			Align:     constants.CenterAlign,
			Border:    "T",
		},
	})

	border := "T"

	for i := 11; i < len(enterpriseSlugs); i += 2 {
		if i == len(enterpriseSlugs)-2 {
			border = "TB"
		}
		row := []domain.TableCell{
			{
				Text:      guaranteeNamesMap[enterpriseSlugs[i]],
				Height:    4.5,
				Width:     65,
				FontSize:  constants.RegularFontSize,
				FontStyle: constants.BoldFontStyle,
				FontColor: constants.BlackColor,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    border,
			},
			{
				Text:      enterpriseData[enterpriseSlugs[i]].sumInsuredLimitOfIndemnity,
				Height:    4.5,
				Width:     25,
				FontSize:  constants.RegularFontSize,
				FontStyle: constants.RegularFontStyle,
				FontColor: constants.BlackColor,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.RightAlign,
				Border:    border,
			},
			{
				Text:      " ",
				Height:    4.5,
				Width:     5,
				FontSize:  constants.RegularFontSize,
				FontStyle: constants.RegularFontStyle,
				FontColor: constants.BlackColor,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.RightAlign,
				Border:    border,
			},
			{
				Text:      guaranteeNamesMap[enterpriseSlugs[i+1]],
				Height:    4.5,
				Width:     65,
				FontSize:  constants.RegularFontSize,
				FontStyle: constants.BoldFontStyle,
				FontColor: constants.BlackColor,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    border,
			},
			{
				Text:      enterpriseData[enterpriseSlugs[i+1]].sumInsuredLimitOfIndemnity,
				Height:    4.5,
				Width:     30,
				FontSize:  constants.RegularFontSize,
				FontStyle: constants.RegularFontStyle,
				FontColor: constants.BlackColor,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.RightAlign,
				Border:    border,
			},
		}

		table = append(table, row)
	}

	ccg.engine.DrawTable(table)

	ccg.engine.NewLine(5)

	sumInsuredLimitOfIndemnity := "5.000.000"
	if enterpriseData[rctoGuaranteeSlug].sumInsuredLimitOfIndemnity != "3000000" {
		sumInsuredLimitOfIndemnity = "7.500.000"
	}

	ccg.engine.WriteText(domain.TableCell{
		Text:      "Il limite di esposizione massima annua della Compagnia, per la sezione Responsabilità Civile, è pari ad € " + sumInsuredLimitOfIndemnity,
		Height:    3.5,
		Width:     190,
		FontSize:  constants.RegularFontSize,
		FontStyle: constants.BoldFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "",
	})
	ccg.engine.WriteText(domain.TableCell{
		Text:      "a valere per tutte le garanzie e complessivamente per tutti gli Assicurati, per ogni sinistro, qualunque sia il numero delle persone decedute o lese o che abbiano subito danni a cose di loro proprietà.",
		Height:    3.5,
		Width:     190,
		FontSize:  constants.RegularFontSize,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "",
	})
}

func (ccg *CommercialCombinedGenerator) deductibleSection() {
	const (
		descriptionColumnWidth = 90
		otherColumnWidth       = 50
	)

	type tableSection struct {
		title   string
		entries [][]string
	}

	type table struct {
		header     []string
		subHeaders []string
		emptyLine  bool
		sections   []tableSection
		newPage    bool
	}

	parseSections := func(sections []tableSection) [][]domain.TableCell {
		result := make([][]domain.TableCell, 0)

		for sectionIndex, section := range sections {
			result = append(result, []domain.TableCell{{
				Text:      section.title,
				Height:    4.5,
				Width:     190,
				FontSize:  constants.RegularFontSize,
				FontStyle: constants.BoldFontStyle,
				FontColor: constants.BlackColor,
				Fill:      true,
				FillColor: constants.GreyColor,
				Align:     constants.LeftAlign,
				Border:    "TLR",
			}})

			borders := []string{"TL", "TLR"}

			for entriesIndex, entries := range section.entries {
				if sectionIndex == len(sections)-1 && entriesIndex == len(entries)-1 {
					borders = []string{"TLB", "1"}
				}

				row := []domain.TableCell{
					{
						Text:      entries[0],
						Height:    4.5,
						Width:     25,
						FontSize:  constants.MediumFontSize,
						FontStyle: constants.RegularFontStyle,
						FontColor: constants.BlackColor,
						Fill:      false,
						FillColor: domain.Color{},
						Align:     constants.RightAlign,
						Border:    borders[0],
					},
					{
						Text:      entries[1],
						Height:    4.5,
						Width:     65,
						FontSize:  constants.MediumFontSize,
						FontStyle: constants.RegularFontStyle,
						FontColor: constants.BlackColor,
						Fill:      false,
						FillColor: domain.Color{},
						Align:     constants.LeftAlign,
						Border:    borders[0],
					},
					{
						Text:      entries[2],
						Height:    4.5,
						Width:     otherColumnWidth,
						FontSize:  constants.MediumFontSize,
						FontStyle: constants.RegularFontStyle,
						FontColor: constants.BlackColor,
						Fill:      false,
						FillColor: domain.Color{},
						Align:     constants.RightAlign,
						Border:    borders[0],
					},
					{
						Text:      entries[3],
						Height:    4.5,
						Width:     otherColumnWidth,
						FontSize:  constants.MediumFontSize,
						FontStyle: constants.RegularFontStyle,
						FontColor: constants.BlackColor,
						Fill:      false,
						FillColor: domain.Color{},
						Align:     constants.RightAlign,
						Border:    borders[1],
					},
				}
				result = append(result, row)
			}
		}

		return result
	}

	parseTable := func(t table) [][]domain.TableCell {
		result := make([][]domain.TableCell, 0)

		if len(t.header) > 0 {
			headerRow := []domain.TableCell{
				{
					Text:      t.header[0],
					Height:    4.5,
					Width:     descriptionColumnWidth,
					FontSize:  constants.RegularFontSize,
					FontStyle: constants.RegularFontStyle,
					FontColor: constants.BlackColor,
					Fill:      false,
					FillColor: domain.Color{},
					Align:     constants.LeftAlign,
					Border:    "",
				},
				{
					Text:      t.header[1],
					Height:    4.5,
					Width:     otherColumnWidth,
					FontSize:  constants.RegularFontSize,
					FontStyle: constants.BoldFontStyle,
					FontColor: constants.BlackColor,
					Fill:      true,
					FillColor: constants.GreyColor,
					Align:     constants.CenterAlign,
					Border:    "TL",
				},
				{
					Text:      t.header[2],
					Height:    4.5,
					Width:     otherColumnWidth,
					FontSize:  constants.RegularFontSize,
					FontStyle: constants.BoldFontStyle,
					FontColor: constants.BlackColor,
					Fill:      true,
					FillColor: constants.GreyColor,
					Align:     constants.CenterAlign,
					Border:    "TLR",
				},
			}
			result = append(result, headerRow)
		}

		if len(t.subHeaders) > 0 {
			subHeaderRow := []domain.TableCell{
				{
					Text:      t.subHeaders[0],
					Height:    4.5,
					Width:     descriptionColumnWidth,
					FontSize:  constants.RegularFontSize,
					FontStyle: constants.BoldFontStyle,
					FontColor: constants.BlackColor,
					Fill:      true,
					FillColor: constants.GreyColor,
					Align:     constants.LeftAlign,
					Border:    "TL",
				},
				{
					Text:      t.subHeaders[1],
					Height:    4.5,
					Width:     otherColumnWidth,
					FontSize:  constants.MediumFontSize,
					FontStyle: constants.RegularFontStyle,
					FontColor: constants.BlackColor,
					Fill:      false,
					FillColor: domain.Color{},
					Align:     constants.RightAlign,
					Border:    "TL",
				},
				{
					Text:      t.subHeaders[2],
					Height:    4.5,
					Width:     otherColumnWidth,
					FontSize:  constants.MediumFontSize,
					FontStyle: constants.RegularFontStyle,
					FontColor: constants.BlackColor,
					Fill:      false,
					FillColor: domain.Color{},
					Align:     constants.RightAlign,
					Border:    "TLR",
				},
			}
			result = append(result, subHeaderRow)
		}

		if t.emptyLine {
			result = append(result, []domain.TableCell{
				{
					Text:      " ",
					Height:    4.5,
					Width:     190,
					FontSize:  constants.RegularFontSize,
					FontStyle: constants.RegularFontStyle,
					FontColor: constants.BlackColor,
					Fill:      false,
					FillColor: domain.Color{},
					Align:     constants.LeftAlign,
					Border:    "TLR",
				},
			})
		}

		sections := parseSections(t.sections)
		result = append(result, sections...)

		return result
	}

	ccg.engine.WriteText(domain.TableCell{
		Text:      "Franchigie, scoperti, limiti di indennizzo",
		Height:    4.5,
		Width:     190,
		FontSize:  constants.LargeFontSize,
		FontStyle: constants.BoldFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "",
	})

	ccg.engine.NewLine(3)

	rawTables := []table{
		{
			header:     []string{" ", "Limiti di indennizzo", "Scoperto/Franchigia"},
			subHeaders: []string{"Limite di Polizza e franchigia frontale", "Somma assicurata", "€ 1.500"},
			emptyLine:  true,
			sections: []tableSection{
				{
					title: "SEZIONE A - INCENDIO E \"TUTTI I RISCHI\"",
					entries: [][]string{
						{"Art. A/01 - a", "Costi per demolire, sgomberare, trattare e trasportare.", "10% indennizzo max € 500." +
							"000", "nessuna"},
						{" ", "Sottolimite per Tossici e Nocivi", "€ 20.000", "nessuna"},
						{"Art. A/01 - b", "Costi/Oneri di urbanizzazione", " 10% indennizzo max € 50.000", "nessuna"},
						{"Art. A/01 - c", "Costi per rimuovere, trasportare, ricollocare i beni", "10% SA Contenuto max € 50." +
							"000", "nessuna"},
						{"Art. A/01 - d", "Onorari dei periti\nOnorari progettisti/consulenti/professionisti",
							"10% indennizzo max €25.000\n10% indennizzo max € 25.000", "nessuna"},
						{"Art. A/03.1", "Cose speciali come previste in Polizza (Disegni,modelli)",
							"10% SA macchinari max €50.000", "Franchigia frontale"},
						{"Art. A/03.2", "Eventi atmosferici", "70% SA max € 10.0000.000", "10% min € 5.0000"},
						{"Art. A/03.3", "Grandine su fragili: lastre, fabbricati aperti da più lati", "€ 200.000",
							"Franchigia frontale"},
						{"Art. A/03.4", "Eventi sociopolitici (escluso terrorismo)", "80% SA max € 10.000.000",
							"10% min. € 2.5000"},
						{"Art. A/03.6", "Sovraccarico neve", "50% SA max € 2.000.000", "10% min. € 5.000"},
						{"Art. A/03.7", "Valori", "€ 5.000", "€ 1.000"},
						{"Art. A/03.8", "Acqua Condotta\nRicerca del guasto", "€ 500.000\n€ 10.000", "10% min € 1.000"},
						{"Art. A/03.9", "Gelo", "€ 50.000", "Franchigia frontale"},
					},
				},
				{
					title: "SEZIONE A - CONDIZIONI PARTICOLARI SEMPRE OPERANTI",
					entries: [][]string{
						{"Art. A/11", "Acqua Piovana", "€ 20.000", "Franchigia frontale"},
						{"Art. A/12", "Dispersione liquidi", "€ 50.000", "€ 5.000"},
						{"Art. A/13", "Rigurgiti e Traboccamenti di Fogna", "€ 20.000", "Franchigia frontale"},
						{"Art. A/14", "Rottura lastre", "€ 5.000", "€ 500"},
						{"Art. A/15", "Decentramento merci e macchinari", "€ 500.000", "Franchigia frontale"},
						{"Art. A/16", "Miscelazione accidentale delle merci", "€ 50.000", "€ 5.000"},
						{"Art. A/17", "Colaggio da Impianti Automatici di Estinzione", "€ 100.000", "10% min € 2.500"},
						{"Art. A/18", "Inondazione, Alluvione", "50% SA max € 7.000.000", "10% min € 10,000"},
						{"Art. A/19", "Allagamento", "€ 500.000", "10% min € 2.500"},
						{"Art. A/20", "Terremoto", " 50% SA max € 7.000.000", "10% min € 10.000"},
						{"Art. A/21", "Terrorismo e Sabotaggio", "50% SA max € 5.000.000", "10% min € 5.000"},
						{"Art. A/22", "Movimentazione Interna / Urto Veicoli", "€ 5.000", "€ 1.000"},
						{"Art. A/23", "Fenomeno elettrico", "€ 10.000", "€ 1.500"},
						{"Art. A/24", "Rischio Locativo", "€ 100.000", "nessuna"},
						{"Art. A/25", "Ricorso dei Terzi", "€ 250.000", "nessuna"},
						{"Art. A/26", "Apparecchiature Elettroniche", "50.000", "10% min € 1.500"},
						{"Art. A/27", "Guasti alle Macchine", "50.000", "10% min € 5.000"},
						{"Art. A/28", "Fuoriuscita materiale fuso", "10% merci max € 100.000", "10% min € 1.500"},
					},
				},
			},
			newPage: true,
		},
		{
			header:     []string{" ", "Limiti di indennizzo", "Scoperto/Franchigia"},
			subHeaders: []string{},
			emptyLine:  false,
			sections: []tableSection{
				{
					title: "SEZIONE A - GARANZIE AGGIUNTIVE",
					entries: [][]string{
						{"Art. A/29", "Aumento temporaneo merci", "Somma assicurata", "nessuna"},
						{"Art. A/30", "Merci in refrigerazione", "Somma assicurata", "10% min € 1.500"},
					},
				},
				{
					title: "SEZIONE B - DANNI INDIRETTI",
					entries: [][]string{
						{"Art. B/1", "Indennità aggiuntiva a percentuale", "Somma assicurata", "nessuna"},
						{"Art. B/2", "Diaria Giornaliera", " ", "3 giorni"},
						{"Art. B/3", "Maggiori costi", "Somma assicurata", "nessuna"},
						{"Art. B/4", "Perdita pigioni", "Somma assicurata", "nessuna"},
					},
				},
				{
					title: "SEZIONE C - FURTO",
					entries: [][]string{
						{"Art. C/1 - 1/2", "Furto - Rapina - Estorsione", "€ 20.000", "€ 1.000"},
						{" ", "Con i seguenti sottolimiti:", " ", " "},
						{"Art. C/1 - 3", "Danni ai beni assicurati", "10% SA max € 5.000", "€ 500"},
						{"Art. C/1 - 4", "Atti vandalici", "10% SA max € 5.000", "€ 500"},
						{"Art. C/1 - 5", "Guasti cagionati dai ladri", "10% SA max € 10.000", "scoperto 10%"},
						{"Art. C/1 - 6\na.\nb.", "Valori:\nOvunque riposti\nIn cassaforte o armadio corazzato",
							" \n5% SA max € 5.000\n10% SA max € 10.000", " \n€ 250\n€1.000"},
					},
				},
			},
			newPage: false,
		},
		{
			header:     []string{},
			subHeaders: []string{},
			emptyLine:  false,
			sections: []tableSection{
				{
					title: "SEZIONE C - FURTO",
					entries: [][]string{
						{"Art. C/1 - 7", "Furto commesso da dipendenti", "10% SA max € 10.000", "10% min € 1.000 "},
						{"Art. C/1 - 8", "Quadri, tappeti, oggetti d'arte", "10% SA max € 5.000 per oggetto",
							"10% min € 500"},
						{"Art. C/1 - 8", "Merci ed attrezzature presso terzi", "10% SA max € 20.000", "10% min € 1.000"},
						{"Art. C/3", "Strumenti di chiusura dei locali", " ", "20% min € 1.000"},
						{"Art. C/7", "Beni posti all'aperto", "5% SA max € 5.000", "€ 500"},
						{"Art. C/8", "Beni presso mostre e fiere", "20% SA max € 20.000", "10% min € 1.000"},
						{"Art. C/9", "Portavalori", "20% SA max € 10.000", "10%"},
					},
				},
			},
			newPage: false,
		},
	}

	for index, t := range rawTables {
		parsedTable := parseTable(t)
		ccg.engine.DrawTable(parsedTable)
		if t.newPage {
			ccg.engine.NewPage()
			continue
		} else if index < len(rawTables)-1 {
			ccg.engine.NewLine(10)
		}
	}
	ccg.engine.WriteText(domain.TableCell{
		Text:      "SA = Somma Assicurata",
		Height:    4.5,
		Width:     190,
		FontSize:  constants.MediumFontSize,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "",
	})
}

func (ccg *CommercialCombinedGenerator) dynamicDeductibleSection() {
	const (
		descriptionColumnWidth = 100
		otherColumnWidth       = 45
		target                 = 300000
	)

	var (
		rctSumInsuredLimitOfIndemnity, rcpSumInsuredLimitOfIndemnity float64
	)

	type section struct {
		headers []string
		entries [][]string
	}

	parseSections := func(sections []section) [][]domain.TableCell {
		result := make([][]domain.TableCell, 0)

		for sectionIndex, s := range sections {
			if len(s.headers) > 0 {
				row := []domain.TableCell{
					{
						Text:      s.headers[0],
						Height:    4.5,
						Width:     descriptionColumnWidth,
						FontSize:  constants.RegularFontSize,
						FontStyle: constants.BoldFontStyle,
						FontColor: constants.BlackColor,
						Fill:      true,
						FillColor: constants.GreyColor,
						Align:     constants.LeftAlign,
						Border:    "TL",
					},
					{
						Text:      s.headers[1],
						Height:    4.5,
						Width:     otherColumnWidth,
						FontSize:  constants.RegularFontSize,
						FontStyle: constants.BoldFontStyle,
						FontColor: constants.BlackColor,
						Fill:      true,
						FillColor: constants.GreyColor,
						Align:     constants.CenterAlign,
						Border:    "TL",
					},
					{
						Text:      s.headers[2],
						Height:    4.5,
						Width:     otherColumnWidth,
						FontSize:  constants.RegularFontSize,
						FontStyle: constants.BoldFontStyle,
						FontColor: constants.BlackColor,
						Fill:      true,
						FillColor: constants.GreyColor,
						Align:     constants.CenterAlign,
						Border:    "TLR",
					},
				}
				result = append(result, row)
			}

			borders := []string{"TL", "TLR"}

			if len(s.entries) > 0 {
				for entryIndex, entry := range s.entries {
					if sectionIndex == len(sections)-1 && entryIndex == len(s.entries)-1 {
						borders = []string{"TLB", "1"}
					}
					row := []domain.TableCell{
						{
							Text:      entry[0],
							Height:    4.5,
							Width:     descriptionColumnWidth,
							FontSize:  constants.MediumFontSize,
							FontStyle: constants.RegularFontStyle,
							FontColor: constants.BlackColor,
							Fill:      false,
							FillColor: domain.Color{},
							Align:     constants.LeftAlign,
							Border:    borders[0],
						},
						{
							Text:      entry[1],
							Height:    4.5,
							Width:     otherColumnWidth,
							FontSize:  constants.MediumFontSize,
							FontStyle: constants.RegularFontStyle,
							FontColor: constants.BlackColor,
							Fill:      false,
							FillColor: domain.Color{},
							Align:     constants.RightAlign,
							Border:    borders[0],
						},
						{
							Text:      entry[2],
							Height:    4.5,
							Width:     otherColumnWidth,
							FontSize:  constants.MediumFontSize,
							FontStyle: constants.RegularFontStyle,
							FontColor: constants.BlackColor,
							Fill:      false,
							FillColor: domain.Color{},
							Align:     constants.RightAlign,
							Border:    borders[1],
						},
					}
					result = append(result, row)
				}
			}
		}

		return result
	}

	for _, asset := range ccg.policy.Assets {
		if asset.Enterprise != nil {
			for _, guarantee := range asset.Guarantees {
				if guarantee.Slug == rctGuaranteeSlug {
					rctSumInsuredLimitOfIndemnity = guarantee.Value.SumInsuredLimitOfIndemnity
					//rctStartDate = guarantee.Value.StartDate.Format(constants.DayMonthYearFormat)
				} else if guarantee.Slug == rcpGuaranteeSlug {
					rcpSumInsuredLimitOfIndemnity = guarantee.Value.SumInsuredLimitOfIndemnity
					//rcpStartDate = guarantee.Value.StartDate.Format(constants.DayMonthYearFormat)
					//rcpStartDateUSA = guarantee.Value.StartDate.Format(constants.
					//	DayMonthYearFormat) // TODO: get startDate USA
				}
			}
		}
	}

	rctSumInsuredLimitOfIndemnityString := lib.HumanaizePriceEuro(rctSumInsuredLimitOfIndemnity)
	rcpSumInsuredLimitOfIndemnityString := lib.HumanaizePriceEuro(rcpSumInsuredLimitOfIndemnity)
	// TODO: find better names
	sumInsuredLimitA := rctSumInsuredLimitOfIndemnityString
	sumInsuredLimitB := "€ 600.000"
	sumInsuredLimitC := "€ 1.500.00"
	sumInsuredLimitD := "€ 150.000"
	sumInsuredLimitE := "€ 100.000"
	sumInsuredLimitF := "€ 750.000"
	sumInsuredLimitG := "€ 500.000"
	if rctSumInsuredLimitOfIndemnity != target {
		sumInsuredLimitA = "€ 500.000"
		sumInsuredLimitB = "€ 1.000.000"
		sumInsuredLimitC = "€ 2.500.000"
		sumInsuredLimitD = "€ 250.000"
		sumInsuredLimitE = "€ 150.000"
		sumInsuredLimitF = "€ 1.500.000"
		sumInsuredLimitG = "€ 1.000.000"
	}

	sections := []section{
		{
			headers: []string{"SEZIONE D - RESPONSABILITÀ CIVILE VERSO TERZI (RCT)", "Massimale per sinistro",
				"Scoperto/Franchigia"},
			entries: [][]string{
				{"Art. D/1 Responsabilità Civile Terzi", rctSumInsuredLimitOfIndemnityString, "€ 1.000"},
			},
		},
		{
			headers: []string{"CONDIZIONI AGGIUNTIVE – SEMPRE OPERANTI", "Sottolimiti sinistro/anno",
				"Scoperto / Franchigia"},
			entries: [][]string{
				{"Art. D/5.1  Committenza auto ed altri veicoli", rctSumInsuredLimitOfIndemnityString, "€ 1.000"},
				{"Art. D/5.2  Aree adibite a parcheggi", rctSumInsuredLimitOfIndemnityString, "€ 1.000"},
				{"Art. D/5.3 Installazione e/o Manutenzione", rctSumInsuredLimitOfIndemnityString, "€ 1.000"},
				{"Art. D/5.4 Committenza lavori straordinari e committenza lavori", rctSumInsuredLimitOfIndemnityString,
					"€ 1.000"},
				{"Art. D/5.5 Cessione di lavori in appalto/subappalto – Resp.tà da Committenza",
					rctSumInsuredLimitOfIndemnityString, "€ 1.000"},
				{"Art. D/5.6 Cessione di lavori in appalto/subappalto – infortuni subiti da subappaltatori e loro" +
					" dipendenti", rctSumInsuredLimitOfIndemnityString, "€ 1.000"},
				{"Art. D/5.7  Danni da interruzioni o sospensioni di attività", sumInsuredLimitA, "€ 1.000"},
				{"Art. D/5.8  Danni da furto", sumInsuredLimitA, "€ 1.000"},
				{"Art. D/5.9  Danni alle cose di Terzi in ambito lavori", sumInsuredLimitB, "€ 1.000"},
				{"Art. D/5.10 Responsabilità in materia di salute e sicurezza su lavoro",
					rctSumInsuredLimitOfIndemnityString, "€ 1.000"},
				{"Art. D/5.11 Responsabilità in materia di protezione dei dati personali",
					rctSumInsuredLimitOfIndemnityString, "€ 1.000"},
				{"Art. D/5.12 Danni da Incendio", sumInsuredLimitB, "€ 1.000"},
				{"Art. D/5.13 Danni da circolazione all’interno del perimetro aziendale", sumInsuredLimitC, "€ 1.000"},
			},
		},
	}

	parsedTable := parseSections(sections)
	ccg.engine.DrawTable(parsedTable)

	ccg.engine.NewPage()

	sections = []section{
		{
			headers: []string{"CONDIZIONI AGGIUNTIVE – SEMPRE OPERANTI", "Sottolimiti sinistro/anno",
				"Scoperto/Franchigia"},
			entries: [][]string{
				{"Art. D/5.14 Danni a mezzi sotto carico e scarico", sumInsuredLimitB, "€ 1.000"},
				{"Art. D/5.15 Responsabilità civile personale prestatori di lavoro",
					rctSumInsuredLimitOfIndemnityString, "€ 1.000"},
				{"Art. D/5.16 Danni a veicoli", sumInsuredLimitB, "€ 1.000"},
				{"Art. D/5.17 Proprietà e/o conduzione di fabbricati", rctSumInsuredLimitOfIndemnityString, "€ 1.000"},
				{"Art. D/5.18 Beni in leasing", rctSumInsuredLimitOfIndemnityString, "€ 1.000"},
				{"Art. D/5.19 Danni a cose di prestatori di lavoro", sumInsuredLimitD, "€ 1.000"},
				{"Art. D/5.20 Prestatori di lavoro terzi per crollo totale e/o parziale dei fabbricati",
					rctSumInsuredLimitOfIndemnityString, "€ 1.000"},
				{"Art. D/5.21 Mancato o insufficiente servizio di vigilanza", rctSumInsuredLimitOfIndemnityString,
					"€ 1.000"},
				{"Art. D/5.22 Macchinari e impianti azionati da persone non abilitate", rctSumInsuredLimitOfIndemnityString,
					"€ 1.000"},
				{"Art. D/5.23 Nuovi insediamenti, fusioni, incorporazioni, acquisti", rctSumInsuredLimitOfIndemnityString,
					"€ 1.000"},
				{"Art. D/5.24 Sorveglianza pulizia manutenzione riparazione e collaudo", rctSumInsuredLimitOfIndemnityString,
					"€ 1.000"},
				{"Art. D/5.25 Amministratori Terzi", rctSumInsuredLimitOfIndemnityString, "€ 1.000"},
				{"Art. D/5.26 Qualifica di terzi a fornitori e clienti", rctSumInsuredLimitOfIndemnityString, "€ 1.000"},
			},
		},
		{
			headers: []string{"CONDIZIONI PARTICOLARI – SEMPRE OPERANTI", "Sottolimiti sinistro/anno",
				"Scoperto/Franchigia"},
			entries: [][]string{
				{"Art. D/6.1 Danni a condutture tubature e impianti sotterranei", sumInsuredLimitE, "€ 1.000"},
				{"Art. D/6.2 Danni da vibrazione, cedimento o franamento del terreno", sumInsuredLimitE, "€ 1.000"},
				{"Art. D/6.3 Cose in consegna e custodia", sumInsuredLimitE, "€ 1.000"},
				{"Art. D/6.4 Inquinamento improvviso ed accidentale", sumInsuredLimitE, "€ 1.000"},
				{"Art. D/6.5 Cose di terzi sollevate, caricate, scaricate", sumInsuredLimitE, "€ 1.000"},
				{"Art. D/6.6 Lavori di scavo e reinterro", sumInsuredLimitE, "€ 1.000"},
				{"Art. D/6.7 Danni alle cose sulle quali si eseguono i lavori", sumInsuredLimitE, "€ 1.000"},
				{"Art. D/6.8 Danni alle cose di Terzi in ambito lavori che per volume o peso possono essere rimosse",
					sumInsuredLimitE, "€ 1.000"},
				{"Art. D/6.9  Responsabilità civile postuma", sumInsuredLimitA, "€ 1.000"},
				{"Art. D/6.10 Qualifica di Assicurato", rctSumInsuredLimitOfIndemnityString, "€ 1.000"},
				{"Art. D/6.11 Smercio prodotti", rctSumInsuredLimitOfIndemnityString, "€ 1.000"},
				{"Art. D/6.12 Integrativa auto", sumInsuredLimitC, "€ 1.000"},
				{"Art. D/6.13 Fonti radioattive", sumInsuredLimitE, "€ 1.000"},
				{"Art. D/6.14 Responsabilità civile incrociata", sumInsuredLimitE, "€ 1.000"},
				{"Art. D/6.15 Cessione di lavori in appalto/subappalto – Responsabilità dell’Assicurato e dei" +
					" subappaltatori", sumInsuredLimitC, "€ 1.000"},
			},
		},
		{
			headers: []string{"SEZIONE E – RESPONS. CIVILE VERSO PRESTATORI DI LAVORO (RCO)",
				"Massimale per sinistro", "Franch. per Prest.re Lavoro"},
			entries: [][]string{
				{"Art E/1 Responsabilità Civile Verso prestatori di Lavoro\nper sinistro" +
					"\nSottolimite per Prestatore di Lavoro",
					"\n" + rctSumInsuredLimitOfIndemnityString + "\n" + sumInsuredLimitC,
					"€ 2.500"},
			},
		},
		{
			headers: []string{"GARANZIA", "Sottolimiti sinistro/anno", "Franchigia"},
			entries: [][]string{
				{"Malattie Professionali\nSottolimite per Prestatore di Lavoro",
					sumInsuredLimitF + "\n" + sumInsuredLimitG, "€ 2.500"},
			},
		},
		{
			headers: []string{"SEZIONE F - RESPONSABILITÀ CIVILE DA PRODOTTI DIFETTOSI (RCP)",
				"Massimale sinistro/anno", "Franchigia"},
			entries: [][]string{
				{"Per sinistri verificatisi in qualsiasi paese, esclusi USA e Canada",
					rcpSumInsuredLimitOfIndemnityString, "€ 5.000"},
				{"Per sinistri verificatisi in USA e Canada e territori loro sotto la loro giurisdizione (" +
					"esportazione occulta e/o indiretta)", "Nell'ambito massimale RC Prodotti Difettosi",
					"10% min € 25.000"},
			},
		},
	}

	parsedTable = parseSections(sections)
	ccg.engine.DrawTable(parsedTable)

	ccg.engine.NewPage()

	sections = []section{
		{
			headers: []string{"RISCHI INCLUSI (se attiva la relativa sezione)", "Sottolimiti sinistro/anno",
				"Scoperto/Franchigia"},
			entries: [][]string{
				{"Art. F/4 -1 Inquinamento accidentale da prodotto, esclusi USA e Canada", "€ 500.000",
					"come da franchigie sopra indicate per territorialità"},
				{"Art. F/4 - 2 Danni da incendio", "€ 1.000.000", "come da franchigie sopra indicate per territorialità"},
				{"Art. F/4 - 3 Danni al prodotto finito", "€ 1.500.000",
					"come da franchigie sopra indicate per  territorialità"},
				{"Art. F/4 -4 Danni al contenuto", "€ 1.500.000", "come da franchigie sopra indicate per  territorialità"},
				{"Art. F4 - 5 Prodotti promozionali", "€ 3.000.000",
					"come da franchigie sopra indicate per  territorialità"},
			},
		},
		{
			headers: []string{"CONDIZIONI PARTICOLARI (se attiva la relativa sezione)", "Sottolimiti sinistor/anno",
				"Scoperto/Franchigia"},
			entries: [][]string{
				{"Art. F/6 Responsabilità Civile Postuma da installazione", sumInsuredLimitA, "10% min € 5.000"},
				{"Art. F/7 Estensione validità territoriale Usa e Canada", "€ 1.000.000", "10% min € 25.000"},
				{"Art. F/8 Danni patrimoniali puri", "€ 100.000", "10% min € 15.000"},
			},
		},
		{
			headers: []string{"SEZIONE G - RITIRO PRODOTTI (se attiva la relativa sezione)",
				"Massimale sinistro/anno", "Scoperto/Franchigia"},
			entries: [][]string{
				{"Art . G/1 per sinistro, sinistro in serie, anno", "€ 250.000", "€ 10.000"},
			},
		},
	}

	parsedTable = parseSections(sections)
	ccg.engine.DrawTable(parsedTable)

	sections = []section{
		{
			headers: []string{"SEZIONE H – RESPONSABILITÀ AMMINISTRATORI SINDACI DIRIGENTI (D&0)",
				"Massimale sinistro/anno", "Scoperto/Franchigia"},
			entries: [][]string{
				{"Spese legali", "25% del Massimale D&O", "nessuna"},
				{"Spese di Difesa per Inquinamento", "€ 50.000", "nessuna"},
				{"Spese per Pubbliche relazioni", "€ 75.000", "nessuna"},
				{"Costi di difesa in relazione a procedimenti di estradizione", "€ 75.000", "nessuna"},
				{"Spese in sede cautelare o d'urgenza", "€ 75.000", "nessuna"},
				{"Costi per la garanzia finanziaria sostitutiva della cauzione", "€ 50.000", "nessuna"},
				{"Spese di emergenza", "€ 50.000", "nessuna"},
				{"Presenza ad indagini ed esami", "€ 50.000", "nessuna"},
			},
		},
		{
			headers: []string{"SEZIONE I – CYBER RESPONSE E DATA SECURITY", "Massimale sinistro/anno",
				"Scoperto/Franchigia"},
			entries: [][]string{
				{"Sezione I/A e I/B", "Massimale di Polizza", "€ 2.500"},
			},
		},
	}

	parsedTable = parseSections(sections)
	ccg.engine.DrawTable(parsedTable)

	ccg.engine.NewLine(5)
}

func (ccg *CommercialCombinedGenerator) detailsSection() {
	const (
		emptyField        = "====="
		firstColumnWidth  = 65
		secondColumnWidth = 45
		thirdColumnWidth  = 80
		secondColumnX     = 75.0
		thirdColumnX      = 120.0
	)

	type guaranteeStartDateInfo struct {
		rct    string
		rcp    string
		rcpUsa string
		ritiro string
		deo    string
	}

	startDateInfo := guaranteeStartDateInfo{
		rct:    emptyField,
		rcp:    emptyField,
		rcpUsa: emptyField,
		ritiro: emptyField,
		deo:    emptyField,
	}

	for _, asset := range ccg.policy.Assets {
		if asset.Enterprise == nil {
			continue
		}

		for _, guarantee := range asset.Guarantees {
			if guarantee.Value.StartDate == nil {
				continue
			}
			switch guarantee.Slug {
			case rctGuaranteeSlug:
				startDateInfo.rct = guarantee.Value.StartDate.Format(constants.DayMonthYearFormat)
			case rcpGuaranteeSlug:
				startDateInfo.rcp = guarantee.Value.StartDate.Format(constants.DayMonthYearFormat)
				startDateInfo.rcpUsa = guarantee.Value.StartDate.Format(constants.DayMonthYearFormat)
			case ritiroGuaranteeSlug:
				startDateInfo.ritiro = guarantee.Value.StartDate.Format(constants.DayMonthYearFormat)
			case managementOrganizationGuaranteeSlug:
				startDateInfo.deo = guarantee.Value.StartDate.Format(constants.DayMonthYearFormat)
			}
		}
	}

	ccg.engine.WriteText(domain.TableCell{
		Text:      "Dettagli di alcune sezioni",
		Height:    4.5,
		Width:     190,
		FontSize:  constants.LargeFontSize,
		FontStyle: constants.BoldFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "",
	})

	// First three rows
	ccg.engine.DrawTable([][]domain.TableCell{
		{
			{
				Text:      "SEZIONE",
				Height:    4.5,
				Width:     firstColumnWidth,
				FontSize:  constants.LargeFontSize,
				FontStyle: constants.BoldFontStyle,
				FontColor: constants.BlackColor,
				Fill:      true,
				FillColor: constants.LightGreyColor,
				Align:     constants.LeftAlign,
				Border:    "TL",
			},
			{
				Text:      "REQUISITO",
				Height:    4.5,
				Width:     secondColumnWidth,
				FontSize:  constants.LargeFontSize,
				FontStyle: constants.BoldFontStyle,
				FontColor: constants.BlackColor,
				Fill:      true,
				FillColor: constants.LightGreyColor,
				Align:     constants.CenterAlign,
				Border:    "TL",
			},
			{
				Text:      "CONDIZIONE",
				Height:    4.5,
				Width:     thirdColumnWidth,
				FontSize:  constants.LargeFontSize,
				FontStyle: constants.BoldFontStyle,
				FontColor: constants.BlackColor,
				Fill:      true,
				FillColor: constants.LightGreyColor,
				Align:     constants.CenterAlign,
				Border:    "TLR",
			},
		},
		{
			{
				Text:      "RESPONSABILITA’ CIVILE E VERSO TERZI E PRESTATORI DI LAVORO",
				Height:    4.5,
				Width:     firstColumnWidth,
				FontSize:  constants.MediumFontSize,
				FontStyle: constants.RegularFontStyle,
				FontColor: constants.BlackColor,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "TL",
			},
			{
				Text:      "Regime copertura",
				Height:    4.5,
				Width:     secondColumnWidth,
				FontSize:  constants.MediumFontSize,
				FontStyle: constants.RegularFontStyle,
				FontColor: constants.BlackColor,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "TL",
			},
			{
				Text:      "LOSS OCCURRENCE",
				Height:    4.5,
				Width:     thirdColumnWidth,
				FontSize:  constants.MediumFontSize,
				FontStyle: constants.RegularFontStyle,
				FontColor: constants.BlackColor,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "TLR",
			},
		},
		{
			{
				Text:      "MALATTIE PROFESSIONALI",
				Height:    4.5,
				Width:     firstColumnWidth,
				FontSize:  constants.MediumFontSize,
				FontStyle: constants.RegularFontStyle,
				FontColor: constants.BlackColor,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "TL",
			},
			{
				Text:      "Retroattività",
				Height:    4.5,
				Width:     secondColumnWidth,
				FontSize:  constants.MediumFontSize,
				FontStyle: constants.RegularFontStyle,
				FontColor: constants.BlackColor,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "TL",
			},
			{
				Text: "L’Assicurazione vale per le conseguenze di fatti colposi commessi dopo la data del" +
					" " + startDateInfo.rct,
				Height:    4.5,
				Width:     thirdColumnWidth,
				FontSize:  constants.MediumFontSize,
				FontStyle: constants.RegularFontStyle,
				FontColor: constants.BlackColor,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "TLR",
			},
		},
	})

	// Fourth row
	fourthRowY := ccg.engine.GetY()
	// First column
	ccg.engine.WriteText(domain.TableCell{
		Text:      " \n \nRESPONSABILITA' CIVILE\nDA PRODOTTI DIFETTOSI\n \n ",
		Height:    4.5,
		Width:     firstColumnWidth,
		FontSize:  constants.MediumFontSize,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "TL",
	})
	// Second column
	ccg.engine.SetY(fourthRowY)
	ccg.engine.SetX(secondColumnX)
	ccg.engine.WriteText(domain.TableCell{
		Text:      "Regime copertura",
		Height:    4.5,
		Width:     secondColumnWidth,
		FontSize:  constants.MediumFontSize,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "TL",
	})
	ccg.engine.SetX(secondColumnX)
	ccg.engine.WriteText(domain.TableCell{
		Text:      "Retroattività – Mondo escluso USA Canada",
		Height:    4.5,
		Width:     secondColumnWidth,
		FontSize:  constants.MediumFontSize,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "TL",
	})
	ccg.engine.SetX(secondColumnX)
	ccg.engine.WriteText(domain.TableCell{
		Text:      " \nRetroattività – USA Canada\n ",
		Height:    4.5,
		Width:     secondColumnWidth,
		FontSize:  constants.MediumFontSize,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "TL",
	})
	// Third column
	ccg.engine.SetY(fourthRowY)
	ccg.engine.SetX(thirdColumnX)
	ccg.engine.WriteText(domain.TableCell{
		Text:      "CLAIMS MADE",
		Height:    4.5,
		Width:     thirdColumnWidth,
		FontSize:  constants.MediumFontSize,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "TLR",
	})
	ccg.engine.SetX(thirdColumnX)
	ccg.engine.WriteText(domain.TableCell{
		Text:      "L'Assicurazione vale per i danni verificatisi dopo la data del " + startDateInfo.rcp,
		Height:    4.5,
		Width:     thirdColumnWidth,
		FontSize:  constants.MediumFontSize,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "TLR",
	})
	ccg.engine.SetX(thirdColumnX)
	ccg.engine.WriteText(domain.TableCell{
		Text: "L'Assicurazione vale per i danni verificatisi dopo la data del " + startDateInfo.rcpUsa + " purch" +
			"é  relativi a prodotti descritti in Polizza consegnati a terzi dopo la stessa data.",
		Height:    4.5,
		Width:     thirdColumnWidth,
		FontSize:  constants.MediumFontSize,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "TLR",
	})

	// Fifth row
	ccg.engine.DrawTable([][]domain.TableCell{
		{
			{
				Text: "RESPONSABILITA’ CIVILE\nVERSO TERZI E PRESTATORI DI LAVORO\nRESPONSABILITA' CIVILE\nDA" +
					" PRODOTTI DIFETTOSI\n",
				Height:    4.5,
				Width:     firstColumnWidth,
				FontSize:  constants.MediumFontSize,
				FontStyle: constants.RegularFontStyle,
				FontColor: constants.BlackColor,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "TLB",
			},
			{
				Text:      " \nPremio minimo\n ",
				Height:    4.5,
				Width:     secondColumnWidth,
				FontSize:  constants.MediumFontSize,
				FontStyle: constants.RegularFontStyle,
				FontColor: constants.BlackColor,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "TLB",
			},
			{
				Text: " \nPremio minimo indicato in Polizza calcolato sui parametri di fatturato e Prestatori di" +
					" lavoro dichiarati\n ",
				Height:    4.5,
				Width:     thirdColumnWidth,
				FontSize:  constants.MediumFontSize,
				FontStyle: constants.RegularFontStyle,
				FontColor: constants.BlackColor,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "1",
			},
		},
	})

	ccg.engine.NewPage()

	ccg.engine.DrawTable([][]domain.TableCell{
		{
			{
				Text:      "SEZIONE",
				Height:    4.5,
				Width:     firstColumnWidth,
				FontSize:  constants.LargeFontSize,
				FontStyle: constants.BoldFontStyle,
				FontColor: constants.BlackColor,
				Fill:      true,
				FillColor: constants.LightGreyColor,
				Align:     constants.LeftAlign,
				Border:    "TL",
			},
			{
				Text:      "REQUISITO",
				Height:    4.5,
				Width:     secondColumnWidth,
				FontSize:  constants.LargeFontSize,
				FontStyle: constants.BoldFontStyle,
				FontColor: constants.BlackColor,
				Fill:      true,
				FillColor: constants.LightGreyColor,
				Align:     constants.CenterAlign,
				Border:    "TL",
			},
			{
				Text:      "CONDIZIONE",
				Height:    4.5,
				Width:     thirdColumnWidth,
				FontSize:  constants.LargeFontSize,
				FontStyle: constants.BoldFontStyle,
				FontColor: constants.BlackColor,
				Fill:      true,
				FillColor: constants.LightGreyColor,
				Align:     constants.CenterAlign,
				Border:    "TLR",
			},
		},
	})

	// Sixth row
	sixthRowY := ccg.engine.GetY()
	// First Column
	ccg.engine.WriteText(domain.TableCell{
		Text:      " \n RITIRO PRODOTTI\n \n ",
		Height:    4.5,
		Width:     firstColumnWidth,
		FontSize:  constants.MediumFontSize,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "TL",
	})
	// Second Column
	ccg.engine.SetY(sixthRowY)
	ccg.engine.SetX(secondColumnX)
	ccg.engine.WriteText(domain.TableCell{
		Text:      "Regime di copertura",
		Height:    4.5,
		Width:     secondColumnWidth,
		FontSize:  constants.MediumFontSize,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "TL",
	})
	ccg.engine.SetX(secondColumnX)
	ccg.engine.WriteText(domain.TableCell{
		Text:      " \nRetroattività\n ",
		Height:    4.5,
		Width:     secondColumnWidth,
		FontSize:  constants.MediumFontSize,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "TL",
	})
	// Third Column
	ccg.engine.SetY(sixthRowY)
	ccg.engine.SetX(thirdColumnX)
	ccg.engine.WriteText(domain.TableCell{
		Text:      "CLAIMS MADE",
		Height:    4.5,
		Width:     thirdColumnWidth,
		FontSize:  constants.MediumFontSize,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "TLR",
	})
	ccg.engine.SetX(thirdColumnX)
	ccg.engine.WriteText(domain.TableCell{
		Text: "L'Assicurazione vale per i danni verificatisi dopo la data del " + startDateInfo.
			ritiro + " purché relativi a prodotti descritti in Polizza consegnati a terzi dopo la stessa data.",
		Height:    4.5,
		Width:     thirdColumnWidth,
		FontSize:  constants.MediumFontSize,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "TLR",
	})

	// Seventh row
	seventhRowY := ccg.engine.GetY()
	// First Column
	ccg.engine.WriteText(domain.TableCell{
		Text:      " \n \n \n \nRESPONSABILITÀ AMMINISTRATORI\nSINDACI DIRIGENTI (D&0)\n \n \n \n ",
		Height:    4.5,
		Width:     firstColumnWidth,
		FontSize:  constants.MediumFontSize,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "TL",
	})
	// Second Column
	ccg.engine.SetY(seventhRowY)
	ccg.engine.SetX(secondColumnX)
	ccg.engine.WriteText(domain.TableCell{
		Text:      "Territorialità",
		Height:    4.5,
		Width:     secondColumnWidth,
		FontSize:  constants.MediumFontSize,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "TL",
	})
	ccg.engine.SetX(secondColumnX)
	ccg.engine.WriteText(domain.TableCell{
		Text:      "Retroattività",
		Height:    4.5,
		Width:     secondColumnWidth,
		FontSize:  constants.MediumFontSize,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "TL",
	})
	ccg.engine.SetX(secondColumnX)
	ccg.engine.WriteText(domain.TableCell{
		Text:      "Data di continuità",
		Height:    4.5,
		Width:     secondColumnWidth,
		FontSize:  constants.MediumFontSize,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "TL",
	})
	ccg.engine.SetX(secondColumnX)
	ccg.engine.WriteText(domain.TableCell{
		Text:      " \n \nPremio addizionale per il maggior termine di notifica\n ",
		Height:    4.5,
		Width:     secondColumnWidth,
		FontSize:  constants.MediumFontSize,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "TL",
	})
	ccg.engine.SetX(secondColumnX)
	ccg.engine.WriteText(domain.TableCell{
		Text:      "Maggior termine di notifica per amministratori cessati",
		Height:    4.5,
		Width:     secondColumnWidth,
		FontSize:  constants.MediumFontSize,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "TL",
	})
	// Third Column
	ccg.engine.SetY(seventhRowY)
	ccg.engine.SetX(thirdColumnX)
	ccg.engine.WriteText(domain.TableCell{
		Text:      "Unione Economica Europea",
		Height:    4.5,
		Width:     thirdColumnWidth,
		FontSize:  constants.MediumFontSize,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "TLR",
	})
	ccg.engine.SetX(thirdColumnX)
	ccg.engine.WriteText(domain.TableCell{
		Text:      "Illimitata",
		Height:    4.5,
		Width:     thirdColumnWidth,
		FontSize:  constants.MediumFontSize,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "TLR",
	})
	ccg.engine.SetX(thirdColumnX)
	ccg.engine.WriteText(domain.TableCell{
		Text:      startDateInfo.deo,
		Height:    4.5,
		Width:     thirdColumnWidth,
		FontSize:  constants.MediumFontSize,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "TLR",
	})
	ccg.engine.SetX(thirdColumnX)
	ccg.engine.WriteText(domain.TableCell{
		Text: "12 mesi al 30% dell'ultimo premio pagato\n24 mesi al 60% dell'ultimo premio pagato\n36 mesi al 90% dell" +
			"'ultimo premio pagato\n48 mesi al 120% dell'ultimo premio pagato\n60 mesi al 150% dell'ultimo premio" +
			" pagato",
		Height:    4.5,
		Width:     thirdColumnWidth,
		FontSize:  constants.MediumFontSize,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "TLR",
	})
	ccg.engine.SetX(thirdColumnX)
	ccg.engine.WriteText(domain.TableCell{
		Text:      "60 mesi\n ",
		Height:    4.5,
		Width:     thirdColumnWidth,
		FontSize:  constants.MediumFontSize,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "TLR",
	})

	// Eighth row
	eighthRowY := ccg.engine.GetY()
	// First Column
	ccg.engine.WriteText(domain.TableCell{
		Text:      " \n \n \n \n \n \nCYBER RESPONSE E DATA SECURITY\n \n \n \n \n \n",
		Height:    4.5,
		Width:     firstColumnWidth,
		FontSize:  constants.MediumFontSize,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "TLB",
	})
	// Second Column
	ccg.engine.SetY(eighthRowY)
	ccg.engine.SetX(secondColumnX)
	ccg.engine.WriteText(domain.TableCell{
		Text:      "Periodo di carenza Art. I/3\nCyber Business interruption",
		Height:    4.5,
		Width:     secondColumnWidth,
		FontSize:  constants.MediumFontSize,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "TL",
	})
	ccg.engine.SetX(secondColumnX)
	ccg.engine.WriteText(domain.TableCell{
		Text:      "Territorialità",
		Height:    4.5,
		Width:     secondColumnWidth,
		FontSize:  constants.MediumFontSize,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "TL",
	})
	ccg.engine.SetX(secondColumnX)
	ccg.engine.WriteText(domain.TableCell{
		Text:      "Retroattività",
		Height:    4.5,
		Width:     secondColumnWidth,
		FontSize:  constants.MediumFontSize,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "TL",
	})
	ccg.engine.SetX(secondColumnX)
	ccg.engine.WriteText(domain.TableCell{
		Text:      " \n \n \nIncident Response (*)\n \n \n \n ",
		Height:    4.5,
		Width:     secondColumnWidth,
		FontSize:  constants.MediumFontSize,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "TLB",
	})
	// Third Column
	ccg.engine.SetY(eighthRowY)
	ccg.engine.SetX(thirdColumnX)
	ccg.engine.WriteText(domain.TableCell{
		Text:      "12 ore per ciascun sinistro\n ",
		Height:    4.5,
		Width:     thirdColumnWidth,
		FontSize:  constants.MediumFontSize,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "TLR",
	})
	ccg.engine.SetX(thirdColumnX)
	ccg.engine.WriteText(domain.TableCell{
		Text:      "Unione Economica Europea",
		Height:    4.5,
		Width:     thirdColumnWidth,
		FontSize:  constants.MediumFontSize,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "TLR",
	})
	ccg.engine.SetX(thirdColumnX)
	ccg.engine.WriteText(domain.TableCell{
		Text:      "Illimitata",
		Height:    4.5,
		Width:     thirdColumnWidth,
		FontSize:  constants.MediumFontSize,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "TLR",
	})
	ccg.engine.SetX(thirdColumnX)
	ccg.engine.WriteText(domain.TableCell{
		Text: "One Network Firm: Advant Nctm\nNumber+39 (02) 38.592.788\nEmail: OneCyberResponseLine." +
			"Italy@clydeco.com\n \n (*) Nel caso in cui venisse scoperto un presunto evento\ninformatico, " +
			"il Contraente potrà contattare il\ncentralino, 24H su 24H, " +
			"al numero o all'indirizzo mail\nsopraindicato. Il servizio è offerto da ADVANT Nctm",
		Height:    4.5,
		Width:     thirdColumnWidth,
		FontSize:  constants.MediumFontSize,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "1",
	})
}

// TODO: parse clause
func (ccg *CommercialCombinedGenerator) specialConditionsSection() {
	clause := "======"

	ccg.engine.WriteText(domain.TableCell{
		Text:      "Condizioni Speciali in deroga alle Condizioni di Assicurazione",
		Height:    4.5,
		Width:     190,
		FontSize:  constants.LargeFontSize,
		FontStyle: constants.BoldFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "",
	})
	ccg.engine.NewLine(1)

	ccg.engine.WriteText(domain.TableCell{
		Text:      "In deroga a quanto riportato nelle Condizioni di Assicurazione, si concorda tra le Parti che:",
		Height:    4.5,
		Width:     190,
		FontSize:  constants.RegularFontSize,
		FontStyle: constants.BoldFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "",
	})
	ccg.engine.NewLine(1)

	ccg.engine.WriteText(domain.TableCell{
		Text:      clause,
		Height:    4.5,
		Width:     190,
		FontSize:  constants.RegularFontSize,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "",
	})
	ccg.engine.NewLine(1)

	ccg.engine.WriteText(domain.TableCell{
		Text:      "Fermo tutto il resto non derogato da quanto precede.",
		Height:    4.5,
		Width:     190,
		FontSize:  constants.RegularFontSize,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "",
	})
}

// TODO: parse bond
func (ccg *CommercialCombinedGenerator) bondSection() {
	hasBondText := "NO"
	bondText := "======"

	ccg.engine.WriteText(domain.TableCell{
		Text:      "Clausola di Vincolo assicurativo " + hasBondText,
		Height:    4.5,
		Width:     190,
		FontSize:  constants.LargeFontSize,
		FontStyle: constants.BoldFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "",
	})
	ccg.engine.NewLine(1)
	ccg.engine.WriteText(domain.TableCell{
		Text: "La presente Polizza, a decorrere dal suo effetto, " +
			"si intende vincolata a favore dell’Istituto Vincolatario " + bondText + " , " +
			"alle condizioni tutte di cui all’Art. 28.",
		Height:    4.5,
		Width:     190,
		FontSize:  constants.RegularFontSize,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "",
	})
}

// TODO: parse data from policy
func (ccg *CommercialCombinedGenerator) resumeSection() {
	const (
		emptyField       = "======"
		no               = "NO"
		yes              = "SI"
		descriptionWidth = 90
		cellWidth        = 25
	)

	type section struct {
		description string
		active      string
		priceNet    string
		taxAmount   string
		priceGross  string
	}

	sections := []section{
		{
			description: "A - INCENDIO E \" TUTTI I RISCHI\"",
			active:      no,
			priceNet:    emptyField,
			taxAmount:   emptyField,
			priceGross:  emptyField,
		},
		{
			description: "B - DANNI INDIRETTI",
			active:      no,
			priceNet:    emptyField,
			taxAmount:   emptyField,
			priceGross:  emptyField,
		},
		{
			description: "C - FURTO",
			active:      no,
			priceNet:    emptyField,
			taxAmount:   emptyField,
			priceGross:  emptyField,
		},
		{
			description: "D - RESPONSABILITA' CIVILE VERSO TERZI (RCT)",
			active:      no,
			priceNet:    emptyField,
			taxAmount:   emptyField,
			priceGross:  emptyField,
		},
		{
			description: "E - RESP. CIVILE VERSO PRESTATORI DI LAVORO (RCO)",
			active:      no,
			priceNet:    emptyField,
			taxAmount:   emptyField,
			priceGross:  emptyField,
		},
		{
			description: "F - RESP. CIVILE DA PRODOTTI DIFETTOSI (RCP)",
			active:      no,
			priceNet:    emptyField,
			taxAmount:   emptyField,
			priceGross:  emptyField,
		},
		{
			description: "G - RITIRO PRODOTTI",
			active:      no,
			priceNet:    emptyField,
			taxAmount:   emptyField,
			priceGross:  emptyField,
		},
		{
			description: "H - RESP. AMMINISTRATORI SINDACI DIRIGENTI (D&O)",
			active:      no,
			priceNet:    emptyField,
			taxAmount:   emptyField,
			priceGross:  emptyField,
		},
		{
			description: "I - CYBER RESPONSE E DATA SECURITY",
			active:      no,
			priceNet:    emptyField,
			taxAmount:   emptyField,
			priceGross:  emptyField,
		},
		{
			description: "TOTALE PREMIO ANNUALE",
			active:      " ",
			priceNet:    emptyField,
			taxAmount:   emptyField,
			priceGross:  emptyField,
		},
		{
			description: "Rata alla firma della polizza",
			active:      " ",
			priceNet:    emptyField,
			taxAmount:   emptyField,
			priceGross:  emptyField,
		},
		{
			description: "Rate successive alla prima",
			active:      " ",
			priceNet:    emptyField,
			taxAmount:   emptyField,
			priceGross:  emptyField,
		},
	}

	table := make([][]domain.TableCell, 0)

	table = append(table, []domain.TableCell{
		{
			Text:      "SEZIONI",
			Height:    4.5,
			Width:     descriptionWidth,
			FontSize:  constants.RegularFontSize,
			FontStyle: constants.BoldFontStyle,
			FontColor: constants.BlackColor,
			Fill:      false,
			FillColor: domain.Color{},
			Align:     constants.LeftAlign,
			Border:    "",
		},
		{
			Text:      "ATTIVATA",
			Height:    4.5,
			Width:     cellWidth,
			FontSize:  constants.RegularFontSize,
			FontStyle: constants.BoldFontStyle,
			FontColor: constants.BlackColor,
			Fill:      false,
			FillColor: domain.Color{},
			Align:     constants.CenterAlign,
			Border:    "",
		},
		{
			Text:      "Imponibile Euro",
			Height:    4.5,
			Width:     cellWidth,
			FontSize:  constants.MediumFontSize,
			FontStyle: constants.BoldFontStyle,
			FontColor: constants.BlackColor,
			Fill:      false,
			FillColor: domain.Color{},
			Align:     constants.LeftAlign,
			Border:    "",
		},
		{
			Text:      "Imposte Euro",
			Height:    4.5,
			Width:     cellWidth,
			FontSize:  constants.MediumFontSize,
			FontStyle: constants.BoldFontStyle,
			FontColor: constants.BlackColor,
			Fill:      false,
			FillColor: domain.Color{},
			Align:     constants.RightAlign,
			Border:    "",
		},
		{
			Text:      "Totale Euro",
			Height:    4.5,
			Width:     cellWidth,
			FontSize:  constants.MediumFontSize,
			FontStyle: constants.BoldFontStyle,
			FontColor: constants.BlackColor,
			Fill:      false,
			FillColor: domain.Color{},
			Align:     constants.RightAlign,
			Border:    "",
		},
	})

	border := "T"
	fontStyle := constants.RegularFontStyle
	for index, s := range sections {
		if index == len(sections)-1 {
			border = "TB"
		}
		if index > len(sections)-3 {
			fontStyle = constants.BoldFontStyle
		}
		table = append(table, []domain.TableCell{
			{
				Text:      s.description,
				Height:    4.5,
				Width:     descriptionWidth,
				FontSize:  constants.RegularFontSize,
				FontStyle: fontStyle,
				FontColor: constants.BlackColor,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    border,
			},
			{
				Text:      s.active,
				Height:    4.5,
				Width:     cellWidth,
				FontSize:  constants.RegularFontSize,
				FontStyle: fontStyle,
				FontColor: constants.BlackColor,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.CenterAlign,
				Border:    border,
			},
			{
				Text:      s.priceNet,
				Height:    4.5,
				Width:     cellWidth,
				FontSize:  constants.RegularFontSize,
				FontStyle: fontStyle,
				FontColor: constants.BlackColor,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.RightAlign,
				Border:    border,
			},
			{
				Text:      s.taxAmount,
				Height:    4.5,
				Width:     cellWidth,
				FontSize:  constants.RegularFontSize,
				FontStyle: fontStyle,
				FontColor: constants.BlackColor,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.RightAlign,
				Border:    border,
			},
			{
				Text:      s.priceGross,
				Height:    4.5,
				Width:     cellWidth,
				FontSize:  constants.RegularFontSize,
				FontStyle: fontStyle,
				FontColor: constants.BlackColor,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.RightAlign,
				Border:    border,
			},
		})
	}

	ccg.engine.WriteText(domain.TableCell{
		Text:      "Il Premio per tutte le coperture assicurative attivate sulla Polizza",
		Height:    4.5,
		Width:     190,
		FontSize:  constants.LargeFontSize,
		FontStyle: constants.BoldFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "",
	})
	ccg.engine.NewLine(1)
	ccg.engine.DrawTable(table)
	ccg.engine.NewLine(0.75)
	ccg.engine.WriteText(domain.TableCell{
		Text: "In caso di sostituzione, " +
			"il premio alla firma è al netto dell’eventuale rimborso dei premi non goduti sulla polizza sostituita." +
			"\nIn ogni caso il premio alla firma può tener conto dell’eventuale diversa durata rispetto alle rate" +
			" successive.",
		Height:    2.5,
		Width:     190,
		FontSize:  constants.MediumFontSize,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "",
	})
}

func (ccg *CommercialCombinedGenerator) claimsStatement() {
	const (
		zeroClaims             = "0"
		zeroClaimValue         = "€ 0,00"
		descriptionColumnWidth = 50
		quantityColumnWidth    = 30
		valueColumnWidth       = 50
		tabWidth               = 40
	)

	type claim struct {
		description string
		quantity    string
		value       string
	}

	claimsMap := map[string]claim{
		danniGuaranteeSlug: {
			description: "Danni ai beni (escluso Furto)",
			quantity:    zeroClaims,
			value:       zeroClaimValue,
		},
		rcpGuaranteeSlug: {
			description: "Responsabilità Civile",
			quantity:    zeroClaims,
			value:       zeroClaimValue,
		},
		furtoGuaranteeSlug: {
			description: "Furto",
			quantity:    zeroClaims,
			value:       zeroClaimValue,
		},
		managementOrganizationGuaranteeSlug: {
			description: "D&O",
			quantity:    zeroClaims,
			value:       zeroClaimValue,
		},
		cyberGuanrateeSlug: {
			description: "Cyber",
			quantity:    zeroClaims,
			value:       zeroClaimValue,
		},
	}

	for _, declaredClaim := range ccg.policy.DeclaredClaims {
		if _, ok := claimsMap[declaredClaim.GuaranteeSlug]; !ok {
			// TODO: improve this case handling
			continue
		}
		quantity, value := 0, 0.0
		for _, history := range declaredClaim.History {
			quantity += history.Quantity
			value += history.Value
		}

		claimsMap[declaredClaim.GuaranteeSlug] = claim{
			quantity: fmt.Sprintf("%d", quantity),
			value:    lib.HumanaizePriceEuro(value),
		}
	}

	ccg.engine.WriteText(domain.TableCell{
		Text:      "Dichiarazioni da leggere con attenzione prima di firmare",
		Height:    4.5,
		Width:     190,
		FontSize:  constants.LargeFontSize,
		FontStyle: constants.BoldFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "",
	})
	ccg.engine.NewLine(1)
	ccg.engine.WriteText(domain.TableCell{
		Text: "Premesso di essere a conoscenza che le dichiarazioni non veritiere, inesatte o reticenti, " +
			"da me rese, possono compromettere il diritto alla prestazione (come da art. 1892, 1893, 1894 c.c.), " +
			"ai fini dell’efficacia delle garanzie",
		Height:    4.5,
		Width:     190,
		FontSize:  constants.RegularFontSize,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "",
	})
	ccg.engine.NewLine(1)
	ccg.engine.SetX(95)
	ccg.engine.WriteText(domain.TableCell{
		Text:      "DICHIARO",
		Height:    4.5,
		Width:     18.75,
		FontSize:  constants.RegularFontSize,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "B",
	})
	ccg.engine.NewLine(1)
	ccg.engine.WriteText(domain.TableCell{
		Text: "1. di agire in qualità di proprietario dei Beni indicati nella presente Scheda di Polizza o per conto" +
			" altrui o di chi spetta;\n2. che i Beni descritti nella presente Scheda di Polizza non sono assicurati" +
			" presso altre compagnie di assicurazioni;\n3. che l’attività assicurata, " +
			"lo stato e le dichiarazioni relative al rischio, " +
			"rispondono a quanto riportato nella presente Scheda di Polizza;\n4. che, " +
			"sui medesimi rischi assicurati con la presente Polizza, nel quinquennio precedente:",
		Height:    4.5,
		Width:     190,
		FontSize:  constants.RegularFontSize,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "",
	})
	ccg.engine.NewLine(1)
	ccg.engine.SetX(14)
	ccg.engine.WriteText(domain.TableCell{
		Text: "a) non vi sono state coperture assicurative annullate dall’assicuratore;\nb) si sono verificati" +
			" eventi dannosi di importo liquidato come indicato nella tabella che segue:\n",
		Height:    4.5,
		Width:     180,
		FontSize:  constants.RegularFontSize,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "",
	})
	ccg.engine.NewLine(1)

	ccg.engine.SetX(tabWidth)
	ccg.engine.DrawTable([][]domain.TableCell{
		{
			{
				Text:      " ",
				Height:    4.5,
				Width:     descriptionColumnWidth,
				FontSize:  constants.MediumFontSize,
				FontStyle: constants.BoldFontStyle,
				FontColor: constants.BlackColor,
				Fill:      true,
				FillColor: constants.LightGreyColor,
				Align:     constants.LeftAlign,
				Border:    "TL",
			},
			{
				Text:      "Numero",
				Height:    4.5,
				Width:     quantityColumnWidth,
				FontSize:  constants.MediumFontSize,
				FontStyle: constants.BoldFontStyle,
				FontColor: constants.BlackColor,
				Fill:      true,
				FillColor: constants.LightGreyColor,
				Align:     constants.CenterAlign,
				Border:    "TL",
			},
			{
				Text:      "Importo complessivo (euro)",
				Height:    4.5,
				Width:     valueColumnWidth,
				FontSize:  constants.MediumFontSize,
				FontStyle: constants.BoldFontStyle,
				FontColor: constants.BlackColor,
				Fill:      true,
				FillColor: constants.LightGreyColor,
				Align:     constants.CenterAlign,
				Border:    "TLR",
			},
		},
	})

	guaranteeSlugs := []string{danniGuaranteeSlug, rcpGuaranteeSlug, furtoGuaranteeSlug, managementOrganizationGuaranteeSlug,
		cyberGuanrateeSlug}

	borders := []string{"TL", "TL", "TLR"}
	for index, slug := range guaranteeSlugs {
		if index == len(guaranteeSlugs)-1 {
			borders = []string{"TLB", "TLB", "1"}
		}
		ccg.engine.SetX(tabWidth)
		ccg.engine.DrawTable([][]domain.TableCell{
			{
				{
					Text:      claimsMap[slug].description,
					Height:    4.5,
					Width:     descriptionColumnWidth,
					FontSize:  constants.MediumFontSize,
					FontStyle: constants.RegularFontStyle,
					FontColor: constants.BlackColor,
					Fill:      false,
					FillColor: domain.Color{},
					Align:     constants.LeftAlign,
					Border:    borders[0],
				},
				{
					Text:      claimsMap[slug].quantity,
					Height:    4.5,
					Width:     quantityColumnWidth,
					FontSize:  constants.MediumFontSize,
					FontStyle: constants.RegularFontStyle,
					FontColor: constants.BlackColor,
					Fill:      false,
					FillColor: domain.Color{},
					Align:     constants.RightAlign,
					Border:    borders[1],
				},
				{
					Text:      claimsMap[slug].value,
					Height:    4.5,
					Width:     valueColumnWidth,
					FontSize:  constants.MediumFontSize,
					FontStyle: constants.RegularFontStyle,
					FontColor: constants.BlackColor,
					Fill:      false,
					FillColor: domain.Color{},
					Align:     constants.RightAlign,
					Border:    borders[2],
				},
			},
		})
	}

	ccg.engine.NewLine(1)
	ccg.engine.WriteText(domain.TableCell{
		Text: "5. al momento della stipula di questa Polizza non ha ricevuto comunicazioni, " +
			"richieste e notifiche che possano configurare un sinistro relativo alle garanzie assicurate e di non" +
			" essere a conoscenza di eventi o circostanze che possano dare origine ad una richiesta di risarcimento.",
		Height:    4.5,
		Width:     190,
		FontSize:  constants.RegularFontSize,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "",
	})
	ccg.engine.NewLine(5)

	ccg.signatureForm()
}

func (ccg *CommercialCombinedGenerator) statementsFirstPart() {
	const id = 1
	statements := lib.SliceFilter(*ccg.policy.Statements, func(statement models.Statement) bool {
		return statement.Id == id
	})

	ccg.printStatement(statements[0])
}

func (ccg *CommercialCombinedGenerator) statementsSecondPart() {
	const id = 1
	statements := lib.SliceFilter(*ccg.policy.Statements, func(statement models.Statement) bool {
		return statement.Id != id
	})

	for _, statement := range statements {
		ccg.printStatement(statement)
	}
}

func (ccg *CommercialCombinedGenerator) qbePrivacySection() {
	ccg.engine.WriteText(domain.TableCell{
		Text:      "DICHIARAZIONI E CONSENSI PRIVACY - assicuratore",
		Height:    4.5,
		Width:     190,
		FontSize:  constants.LargeFontSize,
		FontStyle: constants.BoldFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "",
	})
	ccg.engine.NewLine(1)
	ccg.engine.WriteText(domain.TableCell{
		Text:      "Io sottoscritto dichiaro di avere perso visione dell’Informativa sul trattamento dei dati personali di QBE Europe SA/NV\nRappresentanza generale per l’Italia ai sensi dell’art. 13 del Regolamento UE n. 2016/679 (informativa resa all’interno del\nSet informativo contenente anche la Documentazione Informativa Precontrattuale, il Glossario e le Condizioni di\nAssicurazione) e di averne compreso i contenuti.",
		Height:    4.5,
		Width:     190,
		FontSize:  constants.RegularFontSize,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "",
	})
	ccg.engine.NewLine(3)
	ccg.signatureForm()
}

func (ccg *CommercialCombinedGenerator) qbePersonalDataSection() {
	ccg.engine.WriteText(domain.TableCell{
		Text:      "CONSENSO AL TRATTAMENTO DEI DATI PERSONALI E PARTICOLARI - assicuratore",
		Height:    4.5,
		Width:     190,
		FontSize:  constants.LargeFontSize,
		FontStyle: constants.BoldFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "",
	})
	ccg.engine.NewLine(1)
	ccg.engine.WriteText(domain.TableCell{
		Text: "Presa visione dell’Informativa sul trattamento dei dati personali di QBE Europe SA/NV" +
			" Rappresentanza generale per l’Italia, " +
			"dichiaro di essere consapevole che il trattamento dei dati personali - anche relativi alla mia salute" +
			" - eventualmente forniti da parte di QBE Europe SA/NV Rappresentanza generale per l’Italia in qualità di" +
			" Titolare del trattamento è necessario per l'adempimento delle Finalità Assicurative di cui all" +
			"’Informativa sul trattamento dei dati personali e, pertanto, presto il consenso a tale trattamento. " +
			"QBE Europe SA/NV Rappresentanza generale per l’Italia informa il Contraente della possibilità di" +
			" revocare il suo consenso in qualsiasi momento. Tuttavia, in caso di revoca del consenso, " +
			"il contratto assicurativo non potrà essere eseguito e/o concluso.",
		Height:    4.5,
		Width:     190,
		FontSize:  constants.RegularFontSize,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "",
	})
	ccg.engine.NewLine(3)
	ccg.signatureForm()
}
