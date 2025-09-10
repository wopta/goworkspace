package quote

import (
	"fmt"
	"time"

	"gitlab.dev.wopta.it/goworkspace/document/internal/constants"
	"gitlab.dev.wopta.it/goworkspace/document/internal/domain"
	"gitlab.dev.wopta.it/goworkspace/document/internal/dto"
	"gitlab.dev.wopta.it/goworkspace/document/internal/engine"
	"gitlab.dev.wopta.it/goworkspace/document/pkg/internal/utils"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
)

const (
	deathGuaranteeSlug               string = "death"
	permanentDisabilityGuaranteeSlug string = "permanent-disability"
	temporaryDisabilityGuaranteeSlug string = "temporary-disability"
	seariousIllnessGuaranteeSlug     string = "serious-ill"
)

type LifeGenerator struct {
	*baseGenerator
	dto *dto.QuoteLifeDTO
}

func NewLifeGenerator(engine *engine.Fpdf, policy *models.Policy,
	product *models.Product) *LifeGenerator {
	quoteDTO := dto.NewQuoteLifeDTO()
	quoteDTO.FromData(policy, product)
	return &LifeGenerator{
		baseGenerator: &baseGenerator{
			engine:  engine,
			now:     time.Now(),
			policy:  policy,
			product: product,
			dto:     quoteDTO.QuoteBaseDTO,
		},
		dto: quoteDTO,
	}
}

func (lg *LifeGenerator) Exec() ([]byte, error) {
	const (
		normalSpacing = 5
		largeSpacing  = 10
	)

	lg.engine.SetMargins(10, 10, -1)

	lg.baseGenerator.mainHeader()

	lg.mainFooter()

	lg.engine.NewPage()

	lg.baseGenerator.heading()

	lg.engine.NewLine(largeSpacing)

	lg.insuredDataTable()

	lg.guaranteeTable()

	lg.baseGenerator.priceSummary()

	lg.engine.NewLine(largeSpacing)

	lg.disclaimer()

	lg.engine.NewLine(normalSpacing)

	lg.whoWeAre()

	lg.engine.NewLine(largeSpacing)

	lg.engine.CrossRemainingSpace()

	return lg.engine.RawDoc()
}

func (lg *LifeGenerator) mainFooter() {
	var (
		text = "Wopta per te. Vita è un prodotto assicurativo di AXA France Vie S.A. - " +
			"Rappresentanza Generale per l’Italia\ndistribuito da Wopta Assicurazioni S.r.l."
		logoPath = lib.GetAssetPathByEnvV2() + constants.CompanyLogoMap[lg.policy.Company]
	)

	lg.engine.SetFooter(func() {
		lg.engine.SetX(10)
		lg.engine.SetY(-17)
		lg.engine.WriteText(domain.TableCell{
			Text:      text,
			Height:    constants.CellHeight,
			Width:     constants.FullPageWidth,
			FontColor: constants.PinkColor,
			FontSize:  constants.SmallFontSize,
		})
		lg.engine.WriteText(domain.TableCell{
			Text:   fmt.Sprintf("%d", lg.engine.PageNumber()),
			Height: constants.CellHeight,
			Width:  0,
			Align:  constants.RightAlign,
		})
		lg.engine.InsertImage(logoPath, constants.FullPageWidth, 279, 13, 9)
	})
}

func (lg *LifeGenerator) insuredDataTable() {
	lg.engine.SetDrawColor(constants.PinkColor)

	lg.engine.WriteText(domain.TableCell{
		Text:      "DATI DEL PREVENTIVO",
		Height:    5,
		Width:     constants.FullPageWidth,
		FontStyle: constants.BoldFontStyle,
		FontColor: constants.PinkColor,
	})

	lg.engine.NewLine(2)

	parser := func(rows [][]string) [][]domain.TableCell {
		result := make([][]domain.TableCell, 0, len(rows))

		const (
			cellHeight     = 5
			firstColWidth  = 32
			thirdColWidth  = 24
			fourthColWidth = 25
			secondColWidth = constants.FullPageWidth - firstColWidth - thirdColWidth - fourthColWidth
		)

		for idx, row := range rows {
			fourthColBorder := constants.BorderBottom
			if idx == 2 {
				fourthColBorder = constants.NoBorder
			}

			result = append(result, []domain.TableCell{
				{
					Text:      row[0],
					Height:    cellHeight,
					Width:     firstColWidth,
					FontStyle: constants.BoldFontStyle,
					FontColor: constants.BlackColor,
					FontSize:  constants.RegularFontSize,
				},
				{
					Text:      row[1],
					Height:    cellHeight,
					Width:     secondColWidth,
					FontColor: constants.BlackColor,
					FontSize:  constants.RegularFontSize,
					Border:    constants.BorderBottom,
				},
				{
					Text:      row[2],
					Height:    cellHeight,
					Width:     thirdColWidth,
					FontStyle: constants.BoldFontStyle,
					FontColor: constants.BlackColor,
					FontSize:  constants.RegularFontSize,
				},
				{
					Text:      row[3],
					Height:    cellHeight,
					Width:     fourthColWidth,
					FontColor: constants.BlackColor,
					FontSize:  constants.RegularFontSize,
					Border:    fourthColBorder,
				},
			})
		}

		return result
	}

	rows := [][]string{
		{"Cognome e Nome", constants.EmptyField, " Cod.Fisc.", constants.EmptyField},
		{"Residente in", constants.EmptyField, " Data nascita", lg.dto.Contractor.BirthDate},
		{"Domicilio", constants.EmptyField, " ", " "},
		{"Mail", constants.EmptyField, " Telefono", constants.EmptyField},
	}

	table := parser(rows)

	lg.engine.DrawTable(table)
}

func (lg *LifeGenerator) guaranteeTable() {
	const (
		firstColWidth       = float64(90)
		otherColWidth       = float64(25)
		headingRowHeight    = 5
		guaranteesRowHeight = 8
	)

	lg.engine.SetDrawColor(constants.PinkColor)

	table := make([][]domain.TableCell, 0, 5)

	heading := []string{
		"Garanzie", "Somma assicurata €", "Durata anni", "Scade il",
		"Premio Annuale €"}
	header := make([]domain.TableCell, 0, 5)
	for idx, h := range heading {
		width := otherColWidth
		if idx == 0 {
			width = firstColWidth
		}
		header = append(header, domain.TableCell{
			Text:      h,
			Height:    headingRowHeight,
			Width:     width,
			Align:     constants.CenterAlign,
			FontStyle: constants.BoldFontStyle,
			Border:    constants.BorderTopBottom,
		})
	}

	table = append(table, header)

	guaranteeList := []string{
		deathGuaranteeSlug, permanentDisabilityGuaranteeSlug,
		temporaryDisabilityGuaranteeSlug, seariousIllnessGuaranteeSlug}

	for _, slug := range guaranteeList {
		row := make([]domain.TableCell, 0, 5)
		row = append(row,
			domain.TableCell{
				Text:      lg.dto.Guarantees[slug].Description,
				Height:    guaranteesRowHeight,
				Width:     firstColWidth,
				FontStyle: constants.BoldFontStyle,
				Border:    constants.BorderTopBottom,
			},
			domain.TableCell{
				Text:   lg.dto.Guarantees[slug].SumInsuredLimitOfIndemnity.Text,
				Height: guaranteesRowHeight,
				Width:  otherColWidth,
				Border: constants.BorderTopBottom,
			},
			domain.TableCell{
				Text:   lg.dto.Guarantees[slug].Duration.Text,
				Height: guaranteesRowHeight,
				Width:  otherColWidth,
				Align:  constants.CenterAlign,
				Border: constants.BorderTopBottom,
			},
			domain.TableCell{
				Text:   lg.dto.Guarantees[slug].ExpiryDate,
				Height: guaranteesRowHeight,
				Width:  otherColWidth,
				Align:  constants.CenterAlign,
				Border: constants.BorderTopBottom,
			},
			domain.TableCell{
				Text:   lg.dto.Guarantees[slug].PremiumGrossYearly.Text,
				Height: guaranteesRowHeight,
				Width:  otherColWidth,
				Border: constants.BorderTopBottom,
			},
		)
		table = append(table, row)
	}

	lg.engine.DrawTable(table)
}

func (lg *LifeGenerator) disclaimer() {
	lg.engine.WriteText(domain.TableCell{
		Text:      fmt.Sprintf("Milano, il %s", lg.now.Format(constants.DayMonthYearFormat)),
		Height:    constants.CellHeight,
		FontStyle: constants.BoldFontStyle,
		FontColor: constants.BlackColor,
		FontSize:  constants.LargeFontSize,
	})

	lg.engine.NewLine(5)

	lg.engine.WriteText(domain.TableCell{
		Text: "Il presente preventivo non ha validità di proposta assicurativa. Ha valore esclusivamente nel giorno " +
			"di emissione e non impegna la compagnia alla sottoscrizione del rischio.",
		Height:    constants.CellHeight,
		FontStyle: constants.BoldFontStyle,
		FontColor: constants.BlackColor,
		FontSize:  constants.LargeFontSize,
	})

	lg.engine.NewLine(5)

	lg.engine.WriteText(domain.TableCell{
		Text: "La sottoscrizione del rischio richiede la preventiva valutazione dello stato di salute dell’Assicurato " +
			"attraverso la compilazione di un Questionario medico e/o Rapporto di Visita medica con eventuali esami " +
			"o visite richieste dalla Compagnia in relazione ad età, somme assicurate e stato di salute rilevabile " +
			"dal suddetto Questionario e/o Rapporto di Visita medica.",
		Height:    constants.CellHeight,
		FontStyle: constants.BoldFontStyle,
		FontColor: constants.BlackColor,
		FontSize:  constants.LargeFontSize,
	})

	lg.engine.NewLine(5)

	lg.engine.WriteText(domain.TableCell{
		Text:      "Prima della sottoscrizione leggere il set informativo.",
		Height:    constants.CellHeight,
		FontStyle: constants.BoldFontStyle,
		FontColor: constants.BlackColor,
		FontSize:  constants.LargeFontSize,
	})
}

func (lg *LifeGenerator) whoWeAre() {
	utils.WhoWeAre(lg.engine)
	lg.engine.NewLine(5)

	lg.engine.RawWriteText(domain.TableCell{
		Text:      "AXA France Vie",
		Height:    constants.CellHeight,
		FontStyle: constants.BoldFontStyle,
		FontColor: constants.BlackColor,
		FontSize:  constants.RegularFontSize,
	})
	lg.engine.RawWriteText(domain.TableCell{
		Text:      " (compagnia assicurativa del gruppo AXA). Indirizzo sede legale in Francia: 313 Terrasses de l'Arche, 92727 NANTERRE CEDEX. Numero Iscrizione Registro delle Imprese di Nanterre: 310499959. Autorizzata in Francia (Stato di origine) all’esercizio delle assicurazioni, vigilata in Francia dalla Autorité de Contrôle Prudentiel et de Résolution (ACPR). Numero Matricola Registre des organismes d’assurance: 5020051. // Indirizzo Rappresentanza Generale per l’Italia: Corso Como n. 17, 20154 Milano - CF, P.IVA e N.Iscr. Reg. Imprese 08875230016 - REA MI-2525395 - Telefono: 02-87103548 - Fax: 02-23331247 - PEC: axafrancevie@legalmail.it – sito internet: www.clp.partners.axa. Ammessa ad operare in Italia in regime di stabilimento. Iscritta all’Albo delle imprese di assicurazione tenuto dall’IVASS, in appendice Elenco I, nr. I.00149.",
		Height:    constants.CellHeight,
		FontColor: constants.BlackColor,
		FontSize:  constants.RegularFontSize,
	})
}
