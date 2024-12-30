package contract

import (
	"fmt"
	"strings"

	"github.com/wopta/goworkspace/document/internal/constants"
	"github.com/wopta/goworkspace/document/internal/domain"
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
	formulaGuaranteeSlug                string = "formula"
	diariaGuaranteeSlug                 string = "diaria"
	maggioriGuaranteeSlug               string = "maggiori"
	pigioniGuaranteeSlug                string = "pigioni"
	rctoGuaranteeSlug                   string = "rcto"
	rcpGuaranteeSlug                    string = "rcp"
	deoGuaranteeSlug                    string = "deo"
	cyberGuanrateeSlug                  string = "cyber"
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
		formulaGuaranteeSlug:                "Danni indiretti - Formula",
		diariaGuaranteeSlug:                 "Diaria Giornaliera",
		maggioriGuaranteeSlug:               "Maggiori costi",
		pigioniGuaranteeSlug:                "Perdita Pigioni",
		rctoGuaranteeSlug: "Responsabilità Civile verso Terzi (" +
			"RCT) e verso Prestatori di lavoro (RCO)",
		rcpGuaranteeSlug:   "Responsabilità Civile Prodotti (RCP) Ritiro prodotti",
		deoGuaranteeSlug:   "Responsabilitu Amministratori Sindaci Dirigenti(D&O)",
		cyberGuanrateeSlug: "Cyber Response e Data Security",
	}
)

type QBEGenerator struct {
	*baseGenerator
}

func NewQBEGenerator(engine *engine.Fpdf, isProposal bool) *QBEGenerator {
	return &QBEGenerator{
		&baseGenerator{engine: engine, isProposal: isProposal},
	}
}

func (qb *QBEGenerator) mainHeader(policy *models.Policy) {
	type policyInfo struct {
		code         string
		startDate    string
		endDate      string
		paymentSplit string
		nextPayment  string
		hasBond      string
	}

	type contractorInfo struct {
		name       string
		fiscalCode string
		vatCode    string
		address    string
		mail       string
		phone      string
	}

	plcInfo := policyInfo{
		code:         "=======",
		startDate:    "=======",
		endDate:      "=======",
		paymentSplit: "=======",
		nextPayment:  "=======",
		hasBond:      "NO",
	}

	ctrInfo := contractorInfo{
		name:       "=======",
		fiscalCode: "=======",
		vatCode:    "=======",
		address:    "=======",
		mail:       "=======",
		phone:      "=======",
	}

	if policy.CodeCompany != "" {
		plcInfo.code = policy.CodeCompany
	}

	if !policy.StartDate.IsZero() {
		plcInfo.startDate = policy.StartDate.Format(constants.DayMonthYearFormat)
	}

	if !policy.EndDate.IsZero() {
		plcInfo.endDate = policy.EndDate.Format(constants.DayMonthYearFormat)
	}

	if _, ok := constants.PaymentSplitMap[policy.PaymentSplit]; ok {
		plcInfo.paymentSplit = constants.PaymentSplitMap[policy.PaymentSplit]
	}

	nextPayDate := lib.AddMonths(policy.StartDate.AddDate(policy.Annuity, 0, 0), 12)
	if !nextPayDate.After(policy.EndDate) {
		plcInfo.nextPayment = nextPayDate.Format(constants.DayMonthYearFormat)
	} else {
		plcInfo.nextPayment = plcInfo.endDate
	}

	if policy.HasBond {
		plcInfo.hasBond = "SI"
	}

	if len(policy.Contractor.Name) != 0 {
		ctrInfo.name = policy.Contractor.Name
	}

	if len(policy.Contractor.VatCode) != 0 {
		ctrInfo.vatCode = policy.Contractor.VatCode
	}

	if len(policy.Contractor.FiscalCode) != 0 {
		ctrInfo.fiscalCode = policy.Contractor.FiscalCode
	}

	if policy.Contractor.CompanyAddress != nil {
		ctrInfo.address = fmt.Sprintf("%s %s\n%s %s (%s)", policy.Contractor.CompanyAddress.StreetName,
			policy.Contractor.CompanyAddress.StreetNumber, policy.Contractor.CompanyAddress.PostalCode,
			policy.Contractor.CompanyAddress.City, policy.Contractor.CompanyAddress.CityCode)
	}

	if len(policy.Contractor.Mail) != 0 {
		ctrInfo.mail = policy.Contractor.Mail
	}

	if len(policy.Contractor.Phone) != 0 {
		ctrInfo.phone = policy.Contractor.Phone
	}

	table := [][]domain.TableCell{
		{
			{
				Text:      "I dati della tua Polizza nr. " + plcInfo.code,
				Height:    constants.CellHeight,
				Width:     115,
				FontStyle: constants.BoldFontStyle,
				FontColor: constants.BlackColor,
				FontSize:  constants.RegularFontSize,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "",
			},
			{
				Text:      "I tuoi dati",
				Height:    constants.CellHeight,
				Width:     75,
				FontStyle: constants.BoldFontStyle,
				FontColor: constants.BlackColor,
				FontSize:  constants.RegularFontSize,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "",
			},
		},
		{
			{
				Text:      "Decorre dal: " + plcInfo.startDate + " ore 24:00",
				Height:    constants.CellHeight,
				Width:     115,
				FontStyle: constants.RegularFontStyle,
				FontColor: constants.BlackColor,
				FontSize:  constants.RegularFontSize,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "",
			},
			{
				Text:      "Contraente: " + ctrInfo.name,
				Height:    constants.CellHeight,
				Width:     75,
				FontStyle: constants.RegularFontStyle,
				FontColor: constants.BlackColor,
				FontSize:  constants.RegularFontSize,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "",
			},
		},
		{
			{
				Text:      "Scade il: " + plcInfo.endDate + " ore 24:00",
				Height:    constants.CellHeight,
				Width:     115,
				FontStyle: constants.RegularFontStyle,
				FontColor: constants.BlackColor,
				FontSize:  constants.RegularFontSize,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "",
			},
			{
				Text:      "P.IVA: " + ctrInfo.vatCode,
				Height:    constants.CellHeight,
				Width:     75,
				FontStyle: constants.RegularFontStyle,
				FontColor: constants.BlackColor,
				FontSize:  constants.RegularFontSize,
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
				FontStyle: constants.RegularFontStyle,
				FontColor: constants.BlackColor,
				FontSize:  constants.RegularFontSize,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "",
			},
			{
				Text:      "Codice Fiscale: " + ctrInfo.fiscalCode,
				Height:    constants.CellHeight,
				Width:     75,
				FontStyle: constants.RegularFontStyle,
				FontColor: constants.BlackColor,
				FontSize:  constants.RegularFontSize,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "",
			},
		},
		{
			{
				Text:      "Frazionamento: " + plcInfo.paymentSplit,
				Height:    constants.CellHeight,
				Width:     115,
				FontStyle: constants.RegularFontStyle,
				FontColor: constants.BlackColor,
				FontSize:  constants.RegularFontSize,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "",
			},
			{
				Text:      strings.Split("Indirizzo: "+ctrInfo.address, "\n")[0],
				Height:    constants.CellHeight,
				Width:     75,
				FontStyle: constants.RegularFontStyle,
				FontColor: constants.BlackColor,
				FontSize:  constants.RegularFontSize,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "",
			},
		},
		{
			{
				Text:      "Prossimo pagamento il: " + plcInfo.nextPayment,
				Height:    constants.CellHeight,
				Width:     115,
				FontStyle: constants.RegularFontStyle,
				FontColor: constants.BlackColor,
				FontSize:  constants.RegularFontSize,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "",
			},
			{
				Text:      strings.Split(ctrInfo.address, "\n")[1],
				Height:    constants.CellHeight,
				Width:     75,
				FontStyle: constants.RegularFontStyle,
				FontColor: constants.BlackColor,
				FontSize:  constants.RegularFontSize,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "",
			},
		},
		{
			{
				Text:      "Sostituisce la Polizza: ======",
				Height:    constants.CellHeight,
				Width:     115,
				FontStyle: constants.RegularFontStyle,
				FontColor: constants.BlackColor,
				FontSize:  constants.RegularFontSize,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "",
			},
			{
				Text:      "Mail: " + ctrInfo.mail,
				Height:    constants.CellHeight,
				Width:     75,
				FontStyle: constants.RegularFontStyle,
				FontColor: constants.BlackColor,
				FontSize:  constants.RegularFontSize,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "",
			},
		},
		{
			{
				Text:      "Presenza Vincolo: " + plcInfo.hasBond + " Convenzione: NO",
				Height:    constants.CellHeight,
				Width:     115,
				FontStyle: constants.RegularFontStyle,
				FontColor: constants.BlackColor,
				FontSize:  constants.RegularFontSize,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "",
			},
			{
				Text:      "Telefono: " + ctrInfo.phone,
				Height:    constants.CellHeight,
				Width:     75,
				FontStyle: constants.RegularFontStyle,
				FontColor: constants.BlackColor,
				FontSize:  constants.RegularFontSize,
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

		if qb.isProposal {
			qb.engine.DrawWatermark("PROPOSTA")
		}
	})
}

func (qb *QBEGenerator) mainFooter() {
	text := "QBE Europe SA/NV, Rappresentanza Generale per l’Italia, Via Melchiorre Gioia 8 – 20124 Milano. R.E.A. MI-2538674. Codice fiscale/P.IVA 10532190963 Autorizzazione IVASS n. I.00147\n" +
		"QBE Europe SA/NV è autorizzata dalla Banca Nazionale del Belgio con licenza numero 3093. Sede legale Place du Champ de Mars 5, BE 1050, Bruxelles, Belgio.   N. di registrazione 0690537456."

	qb.engine.SetFooter(func() {
		qb.engine.SetX(10)
		qb.engine.SetY(-17)
		qb.engine.WriteText(domain.TableCell{
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
		qb.engine.WriteText(domain.TableCell{
			Text:      fmt.Sprintf("%d", qb.engine.PageNumber()),
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

func (qb *QBEGenerator) introTable() {
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
	qb.engine.DrawTable(introTable)
}

func (qb *QBEGenerator) whoWeAreTable() {
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
	qb.engine.DrawTable(whoWeAreTable)
}

func (qb *QBEGenerator) insuredDetailsSection(policy *models.Policy) {
	type buildingInfo struct {
		address          string
		buildingMaterial string
		hasSandwichPanel string
		hasAlarm         string
		hasSprinkler     string
		naics            string
		naicsDetail      string
	}

	type enterpriseInfo struct {
		revenue                  string
		northAmericanMarket      string
		employer                 string
		workEmployerRemunaration string
		totalBilled              string
		ownerTotalBilled         string
	}

	newBuildingInfo := func() *buildingInfo {
		return &buildingInfo{
			address:          "======",
			buildingMaterial: "======",
			hasSandwichPanel: "======",
			hasAlarm:         "======",
			hasSprinkler:     "======",
			naics:            "======",
			naicsDetail:      "======",
		}
	}
	newEnterpriseInfo := func() *enterpriseInfo {
		return &enterpriseInfo{
			revenue:                  "======",
			northAmericanMarket:      "======",
			employer:                 "======",
			workEmployerRemunaration: "======",
			totalBilled:              "======",
			ownerTotalBilled:         "======",
		}
	}

	buildings := make([]*buildingInfo, 5)
	for i := 0; i < 5; i++ {
		buildings[i] = newBuildingInfo()
	}

	enterprise := newEnterpriseInfo()

	index := 0
	for _, asset := range policy.Assets {
		if asset.Building == nil {
			continue
		}

		buildings[index].address = fmt.Sprintf("%s, %s - %s %s (%s)", asset.Building.Address,
			asset.Building.StreetNumber, asset.Building.PostalCode, asset.Building.City, asset.Building.CityCode)
		buildings[index].buildingMaterial = asset.Building.BuildingMaterial
		buildings[index].hasSandwichPanel = "NO"
		if asset.Building.HasSandwichPanel {
			buildings[index].hasSandwichPanel = "SI"
		}

		buildings[index].hasAlarm = "NO"
		if asset.Building.HasAlarm {
			buildings[index].hasAlarm = "SI"
		}

		buildings[index].hasSprinkler = "NO"
		if asset.Building.HasSprinkler {
			buildings[index].hasSprinkler = "SI"
		}

		buildings[index].naics = asset.Building.Naics
		buildings[index].naicsDetail = asset.Building.NaicsDetail

		index++
	}

	for _, asset := range policy.Assets {
		if asset.Enterprise == nil {
			continue
		}

		if asset.Enterprise.Revenue != 0.0 {
			enterprise.revenue = lib.HumanaizePriceEuro(asset.Enterprise.Revenue)
		}

		if asset.Enterprise.NorthAmericanMarket != 0.0 {
			enterprise.northAmericanMarket = lib.HumanaizePriceEuro(asset.Enterprise.NorthAmericanMarket)
		}

		if asset.Enterprise.Employer != 0 {
			enterprise.employer = fmt.Sprintf("%d", asset.Enterprise.Employer)
		}

		if asset.Enterprise.WorkEmployersRemuneration != 0.0 {
			enterprise.workEmployerRemunaration = lib.HumanaizePriceEuro(asset.Enterprise.WorkEmployersRemuneration)
		}

		if asset.Enterprise.TotalBilled != 0.0 {
			enterprise.totalBilled = lib.HumanaizePriceEuro(asset.Enterprise.TotalBilled)
		}

		// TODO: add check on OwnerTotalBilled field

	}

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
					building.address, building.buildingMaterial, building.hasSandwichPanel, building.hasAlarm,
					building.hasSprinkler, building.naics, building.naicsDetail),
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

	activityRow := []domain.TableCell{
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
			Text: fmt.Sprintf("Fatturato: %s di cui verso USA e Canada: %s\nPrestatori di lavoro nr: %s"+
				" - Retribuzioni: %s\nTotal Asset: %s di cui capitale proprio: %s",
				enterprise.revenue, enterprise.northAmericanMarket, enterprise.employer,
				enterprise.workEmployerRemunaration, enterprise.totalBilled, enterprise.ownerTotalBilled),
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
	table = append(table, activityRow)

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

	qb.engine.DrawTable(table)

}

// TODO: parse policy info
func (qb *QBEGenerator) guaranteesDetailsSection(policy *models.Policy) {
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
		guastiGuaranteeSlug, elettronicaGuaranteeSlug, furtoGuaranteeSlug, formulaGuaranteeSlug,
		diariaGuaranteeSlug, maggioriGuaranteeSlug,
		pigioniGuaranteeSlug, rctoGuaranteeSlug, rcpGuaranteeSlug, deoGuaranteeSlug, cyberGuanrateeSlug}

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

	for _, slug := range enterpriseSlugs[:6] {
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

	for _, slug := range enterpriseSlugs[6:11] {
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

	for i := 11; i < len(enterpriseSlugs); i += 2 {
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
				Border:    "T",
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
				Border:    "T",
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
				Border:    "T",
			},
			{
				Text:      enterpriseData[enterpriseSlugs[i+1]].sumInsuredLimitOfIndemnity,
				Height:    4.5,
				Width:     25,
				FontSize:  constants.RegularFontSize,
				FontStyle: constants.RegularFontStyle,
				FontColor: constants.BlackColor,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.RightAlign,
				Border:    "T",
			},
		}

		table = append(table, row)
	}

	qb.engine.DrawTable(table)

	qb.engine.NewLine(5)

	sumInsuredLimitOfIndemnity := "5.000.000"
	if enterpriseData[rctoGuaranteeSlug].sumInsuredLimitOfIndemnity != "3000000" {
		sumInsuredLimitOfIndemnity = "7.500.000"
	}

	qb.engine.WriteText(domain.TableCell{
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
	qb.engine.WriteText(domain.TableCell{
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

func (qb *QBEGenerator) deductibleSection() {
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

	qb.engine.WriteText(domain.TableCell{
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

	qb.engine.NewLine(3)

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
						{"Art. A/30", "Merci in refrigerazione", "Somma assicurata", "t10% min € 1.500"},
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
		qb.engine.DrawTable(parsedTable)
		if t.newPage {
			qb.engine.NewPage()
			continue
		} else if index < len(rawTables)-1 {
			qb.engine.NewLine(10)
		}
	}
	qb.engine.WriteText(domain.TableCell{
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

func (qb *QBEGenerator) Contract(policy *models.Policy) ([]byte, error) {
	qb.mainHeader(policy)

	qb.engine.NewPage()

	qb.mainFooter()

	qb.engine.NewLine(10)

	qb.introTable()

	qb.engine.NewLine(10)

	qb.whoWeAreTable()

	qb.engine.NewLine(10)

	qb.insuredDetailsSection(policy)

	qb.engine.NewPage()

	qb.guaranteesDetailsSection(policy)

	qb.engine.NewPage()

	qb.deductibleSection()

	return qb.engine.RawDoc()
}
