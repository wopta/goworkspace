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
				FontSize:  constants.RegularFontsize,
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
				FontSize:  constants.RegularFontsize,
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
				FontSize:  constants.RegularFontsize,
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
				FontSize:  constants.RegularFontsize,
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
				FontSize:  constants.RegularFontsize,
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
				FontSize:  constants.RegularFontsize,
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
				FontSize:  constants.RegularFontsize,
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
				FontSize:  constants.RegularFontsize,
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
				FontSize:  constants.RegularFontsize,
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
				FontSize:  constants.RegularFontsize,
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
				FontSize:  constants.RegularFontsize,
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
				FontSize:  constants.RegularFontsize,
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
				FontSize:  constants.RegularFontsize,
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
				FontSize:  constants.RegularFontsize,
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
				FontSize:  constants.RegularFontsize,
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
				FontSize:  constants.RegularFontsize,
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
				FontSize:  constants.RegularFontsize,
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
				FontSize:  constants.RegularFontsize,
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
				FontSize:  constants.RegularFontsize,
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
				FontSize:  constants.RegularFontsize,
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
				FontSize:  constants.RegularFontsize,
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
				FontSize:  constants.RegularFontsize,
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
				FontSize:  constants.RegularFontsize,
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
				FontSize:  constants.RegularFontsize,
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
				FontSize:  constants.RegularFontsize,
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
				FontSize:  constants.RegularFontsize,
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
				FontSize:  constants.RegularFontsize,
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
				FontSize:  constants.RegularFontsize,
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
				FontSize:  constants.RegularFontsize,
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
				FontSize:  constants.RegularFontsize,
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
				FontSize:  constants.RegularFontsize,
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
			FontSize:  constants.RegularFontsize,
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
			FontSize:  constants.RegularFontsize,
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
			FontSize:  constants.RegularFontsize,
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
			FontSize:  constants.RegularFontsize,
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

	return qb.engine.RawDoc()
}
