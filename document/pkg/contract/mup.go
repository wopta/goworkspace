package contract

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"time"

	"gitlab.dev.wopta.it/goworkspace/document/internal/constants"
	"gitlab.dev.wopta.it/goworkspace/document/internal/domain"
	"gitlab.dev.wopta.it/goworkspace/document/internal/engine"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/network"
)

func GenerateMup(companyName string, consultancyPrice float64, channel string, nodeUid string) (out bytes.Buffer, err error) {
	node := network.GetNetworkNodeByUid(nodeUid)
	if node == nil {
		return bytes.Buffer{}, errors.New("node not found")
	}
	workFor := network.GetNetworkNodeByUid(node.WorksForUid)
	if workFor == nil {
		log.WarningF("workfor %s node not found", node.WorksForUid)
	}
	mockPolicy := models.Policy{
		Channel: channel,
	}
	generator := &baseGenerator{
		engine:       engine.NewFpdf(),
		isProposal:   false,
		now:          time.Now(),
		signatureID:  0,
		networkNode:  node,
		policy:       &mockPolicy,
		worksForNode: workFor,
	}
	generator.mup(true, companyName, consultancyPrice, channel)

	err = generator.engine.GetPdf().Output(&out)
	return out, err
}

func (bg *baseGenerator) mup(isManualGenerated bool, companyName string, consultancyPrice float64, channel string) {
	if !isManualGenerated && (bg.networkNode != nil && !bg.networkNode.HasAnnex && bg.networkNode.Type != models.PartnershipNetworkNodeType) {
		return
	}

	if bg.networkNode == nil || bg.networkNode.Type == models.PartnershipNetworkNodeType || bg.networkNode.IsMgaProponent {
		bg.woptaHeader()

		bg.woptaFooter()
	} else {
		bg.emptyHeader()

		bg.emptyFooter()
	}

	producerInfo := bg.producerInfo()
	proponentInfo := bg.proponentInfo()
	designationInfo := bg.designationInfo()
	mupSection2Info, mupSection5Info := bg.mupInfo(companyName, consultancyPrice, channel)

	bg.engine.NewPage()

	bg.mupTitle()
	bg.engine.NewLine(3)
	bg.mupSectionI(producerInfo, proponentInfo, designationInfo)
	bg.engine.NewLine(3)
	bg.mupSectionII(mupSection2Info)
	bg.engine.NewLine(3)
	bg.mupSectionIII(proponentInfo["name"])
	bg.engine.NewLine(3)
	bg.mupSectionIV()
	bg.engine.NewLine(3)
	bg.mupSectionV(mupSection5Info)
	bg.engine.NewLine(3)
	bg.mupSectionVI()
	if !isManualGenerated {
		bg.engine.NewPage()
	}
	bg.mupSectionVII()
}

func (bg *baseGenerator) mupTitle() {
	text := "ALLEGATO 3 - MODELLO UNICO PRECONTRATTUALE (MUP) PER I PRODOTTI ASSICURATIVI"
	bg.engine.WriteText(bg.engine.GetTableCell(text, constants.BoldFontStyle, constants.CenterAlign))

	bg.engine.NewLine(3)

	text = "L’Intermediario ha l’obbligo di consegnare/trasmettere al contraente il presente Modulo, " +
		"prima della sottoscrizione della proposta o del contratto di assicurazione. Il documento può " +
		"essere fornito con modalità non cartacea se appropriato rispetto alle modalità di distribuzione " +
		"del prodotto assicurativo e il contraente lo consente (art. 120-quater del " +
		"Codice delle Assicurazioni Private)"
	bg.engine.WriteText(bg.engine.GetTableCell(text, constants.RegularFontStyle))
}

func (bg *baseGenerator) mupSectionI(producerInfo, proponentInfo map[string]string, designation string) {
	text := "SEZIONE I - Informazioni generali sull’Intermediario che entra in contatto con " +
		"il contraente"
	bg.engine.WriteText(bg.engine.GetTableCell(text, constants.BoldFontStyle))

	bg.mupProducerTable(producerInfo, proponentInfo, designation)

	text = "Gli estremi identificativi e di iscrizione dell’Intermediario e dei soggetti che " +
		"operano per lo stesso possono essere verificati consultando il Registro Unico degli Intermediari assicurativi " +
		"e riassicurativi sul sito internet dell’IVASS (www.ivass.it)"
	bg.engine.WriteText(bg.engine.GetTableCell(text, constants.RegularFontStyle))
}

func (bg *baseGenerator) mupSectionII(body string) {
	text := "SEZIONE II - Informazioni sul modello di distribuzione"
	bg.engine.WriteText(bg.engine.GetTableCell(text, constants.BoldFontStyle))

	bg.engine.WriteText(bg.engine.GetTableCell(body, constants.RegularFontStyle))
}

func (bg *baseGenerator) mupSectionIII(proponent string) {
	text := "SEZIONE III - Informazioni relative a potenziali situazioni di conflitto " +
		"d’interessi"
	bg.engine.WriteText(bg.engine.GetTableCell(text, constants.BoldFontStyle))

	text = proponent + " ed i soggetti che operano per la stessa non sono " +
		"detentori di una partecipazione, diretta o indiretta, pari o superiore al 10% del capitale sociale o dei " +
		"diritti di voto di alcuna Impresa di assicurazione." + "\n" +
		"Le Imprese di assicurazione o Imprese controllanti un’Impresa di assicurazione " +
		"non sono detentrici di una partecipazione, diretta o indiretta, pari o superiore al 10% del capitale sociale " +
		"o dei diritti di voto dell’Intermediario."
	bg.engine.WriteText(bg.engine.GetTableCell(text, constants.RegularFontStyle))
}

func (bg *baseGenerator) mupSectionIV() {
	text := "SEZIONE IV - Informazioni sull’attività di distribuzione e consulenza"
	bg.engine.WriteText(bg.engine.GetTableCell(text, constants.BoldFontStyle))

	text = "Nello svolgimento dell’attività di distribuzione, l’Intermediario non presta " +
		"attività di consulenza prima della conclusione del contratto né fornisce al contraente una raccomandazione " +
		"personalizzata ai sensi dell’art. 119-ter, comma 3 e 4, del decreto legislativo n. 209/2005 " +
		"(Codice delle Assicurazioni Private)." + "\n" +
		"L'attività di distribuzione assicurativa è svolta in assenza di obblighi " +
		"contrattuali che impongano di offrire esclusivamente i contratti di una o più imprese di " +
		"assicurazioni."
	bg.engine.WriteText(bg.engine.GetTableCell(text, constants.RegularFontStyle))

	companyListText := "L’elenco delle Imprese con cui l’Intermediario ha rapporti d’affari diretti è"
	contractorsFacultyDisclaimer := "È facoltà del contraente chiedere la consegna o la trasmissione di tale elenco."

	if bg.networkNode != nil && !bg.networkNode.IsMgaProponent {
		text = companyListText + " affisso nei propri locali. " + contractorsFacultyDisclaimer
		bg.engine.WriteText(bg.engine.GetTableCell(text, constants.RegularFontStyle))
		return
	}

	website := "https://www.wopta.it/it/information-sets/"
	text = companyListText + " pubblicato sul proprio sito internet "
	bg.engine.RawWriteText(bg.engine.GetTableCell(text, constants.BlackColor))
	bg.engine.WriteLink(website, bg.engine.GetTableCell(website, constants.PinkColor))
	text = ". " + contractorsFacultyDisclaimer
	bg.engine.RawWriteText(bg.engine.GetTableCell(text, constants.BlackColor))
	bg.engine.NewLine(constants.CellHeight)
}

func (bg *baseGenerator) mupSectionV(body string) {
	text := "SEZIONE V - Informazioni relative alle remunerazioni"
	bg.engine.WriteText(bg.engine.GetTableCell(text, constants.BoldFontStyle))

	text = body + "\n" +
		"Le informazioni sopra rese riguardano i compensi complessivamente percepiti da tutti " +
		"gli intermediari coinvolti nella distribuzione del prodotto."
	bg.engine.WriteText(bg.engine.GetTableCell(text, constants.RegularFontStyle))
}

func (bg *baseGenerator) mupSectionVI() {
	text := "SEZIONE VI - Informazioni sul pagamento dei premi"
	bg.engine.WriteText(bg.engine.GetTableCell(text, constants.BoldFontStyle))

	text = "a) Relativamente al prodotto intermediato, i premi pagati dal Contraente " +
		"all’Intermediario e le somme destinate ai risarcimenti o ai pagamenti dovuti dalle Imprese di Assicurazione, " +
		"se regolati per il tramite dell’Intermediario, costituiscono patrimonio autonomo e separato dal patrimonio " +
		"dello stesso."
	bg.engine.WriteText(bg.engine.GetTableCell(text, constants.RegularFontStyle))
	bg.engine.NewLine(constants.CellHeight / 2)
	text = "b) Le modalità di pagamento consentite, nei confronti dell’Intermediario, sono esclusivamente mediante " +
		"bonifico e strumenti di pagamento elettronico, quali ad esempio, carte di credito e/o carte di debito, " +
		"incluse le carte prepagate."
	bg.engine.WriteText(bg.engine.GetTableCell(text, constants.RegularFontStyle))
}

func (bg *baseGenerator) mupSectionVII() {
	text := "SEZIONE VII - Informazioni sugli strumenti di tutela del contraente"
	bg.engine.WriteText(bg.engine.GetTableCell(text, constants.BoldFontStyle))

	text = "L’attività di distribuzione è garantita da un contratto di assicurazione della " +
		"responsabilità civile che copre i danni arrecati ai contraenti da negligenze ed errori professionali " +
		"dell’Intermediario o da negligenze, errori professionali ed infedeltà dei dipendenti, dei collaboratori o " +
		"delle persone del cui operato l’Intermediario deve rispondere a norma di legge." + "\n" +
		"Il contraente ha la facoltà, ferma restando la possibilità di rivolgersi " +
		"all’Autorità Giudiziaria, di inoltrare reclamo per iscritto all’Intermediario, via posta all’indirizzo di " +
		"sede legale o a mezzo mail alla PEC sopra indicati, oppure all’Impresa secondo le modalità e presso i " +
		"recapiti indicati nel DIP aggiuntivo nella relativa sezione, nonché la possibilità, qualora non dovesse " +
		"ritenersi soddisfatto dall’esito del reclamo o in caso di assenza di riscontro da parte dell’Intermediario " +
		"o dell’Impresa entro il termine di legge, di rivolgersi all’IVASS secondo quanto indicato nei DIP aggiuntivi." + "\n" +
		"Il contraente ha la facoltà di avvalersi di altri eventuali sistemi alternativi " +
		"di risoluzione delle controversie previsti dalla normativa vigente nonché quelli indicati nei" +
		" DIP aggiuntivi."
	bg.engine.WriteText(bg.engine.GetTableCell(text, constants.RegularFontStyle))
}

func (bg *baseGenerator) producerInfo() map[string]string {
	producer := map[string]string{
		"name":            "LOMAZZI MICHELE",
		"ruiSection":      "A",
		"ruiCode":         "A000703480",
		"ruiRegistration": "02.03.2022",
	}

	defer log.Printf("producer info: %+v", producer)

	if bg.networkNode == nil || strings.EqualFold(bg.networkNode.Type, models.PartnershipNetworkNodeType) {
		return producer
	}

	switch bg.networkNode.Type {
	case models.AgentNetworkNodeType:
		producer["name"] = strings.ToUpper(bg.networkNode.Agent.Surname) + " " + strings.ToUpper(bg.networkNode.Agent.
			Name)
		producer["ruiSection"] = bg.networkNode.Agent.RuiSection
		producer["ruiCode"] = bg.networkNode.Agent.RuiCode
		producer["ruiRegistration"] = bg.networkNode.Agent.RuiRegistration.Format("02.01.2006")
	case models.AgencyNetworkNodeType:
		producer["name"] = strings.ToUpper(bg.networkNode.Agency.Manager.Name) + " " + strings.ToUpper(bg.networkNode.
			Agency.Manager.Surname)
		producer["ruiSection"] = bg.networkNode.Agency.Manager.RuiSection
		producer["ruiCode"] = bg.networkNode.Agency.Manager.RuiCode
		producer["ruiRegistration"] = bg.networkNode.Agency.Manager.RuiRegistration.Format("02.01.2006")
	}
	return producer
}

func (bg *baseGenerator) proponentInfo() map[string]string {
	proponentInfo := make(map[string]string)

	defer log.Printf("proponent info: %+v", proponentInfo)

	proponentInfo["name"] = "Wopta Assicurazioni Srl"

	if bg.networkNode == nil || bg.networkNode.Type == models.PartnershipNetworkNodeType || bg.networkNode.
		IsMgaProponent {
		proponentInfo["address"] = "Galleria del Corso, 1 - 20122 MILANO (MI)"
		proponentInfo["phone"] = "02.91.24.03.46"
		proponentInfo["email"] = "info@wopta.it"
		proponentInfo["pec"] = "woptaassicurazioni@legalmail.it"
		proponentInfo["website"] = "wopta.it"
	} else {
		proponentNode := bg.networkNode
		if proponentNode.WorksForUid != "" {
			proponentNode = bg.worksForNode
			if proponentNode == nil {
				panic("could not find node for proponent with uid " + bg.networkNode.WorksForUid)
			}
		}

		proponentInfo["address"] = constants.EmptyField
		proponentInfo["phone"] = constants.EmptyField
		proponentInfo["email"] = constants.EmptyField
		proponentInfo["pec"] = constants.EmptyField
		proponentInfo["website"] = constants.EmptyField

		if name := proponentNode.Agency.Name; name != "" {
			proponentInfo["name"] = name
		}
		if address := proponentNode.GetAddress(); address != "" {
			proponentInfo["address"] = address
		}
		if phone := proponentNode.Agency.Phone; phone != "" {
			proponentInfo["phone"] = phone
		}
		if email := proponentNode.Mail; email != "" {
			proponentInfo["email"] = email
		}
		if pec := proponentNode.Agency.Pec; pec != "" {
			proponentInfo["pec"] = pec
		}
		if website := proponentNode.Agency.Website; website != "" {
			proponentInfo["website"] = website
		}
	}

	return proponentInfo
}

func (bg *baseGenerator) mupInfo(companyName string, consultancyPrice float64, channel string) (section2Info, section5Info string) {
	const (
		mgaProponentFormat = "La distribuzione relativamente a questa proposta/contratto è svolta per conto della " +
			"seguente Impresa di assicurazione: %s"
		mgaEmitterFormat = "Il contratto viene intermediato da %s, in qualità di soggetto proponente, che opera in " +
			"virtù della collaborazione con Wopta Assicurazioni Srl (intermediario emittente dell'Impresa di " +
			"Assicurazione %s, iscritto al RUI sezione A nr A000701923 dal 14.02.2022), ai sensi dell’articolo 22, " +
			"comma 10, del decreto legge 18 ottobre 2012, n. 179, convertito nella legge 17 dicembre 2012, n. 221"
		withoutConsultacy = "Per il prodotto intermediato, è corrisposto all’Intermediario un compenso da parte " +
			"dell’Impresa di assicurazione, sotto forma di commissione inclusa nel premio assicurativo."
		withConsultacyFormat = "Per il prodotto intermediato, è corrisposto all’Intermediario un compenso da parte " +
			"dell’Impresa di assicurazione, sotto forma di commissione inclusa nel premio assicurativo, " +
			"e un compenso direttamente dal Contraente, pari ad %s."
	)

	if tCompanyName := constants.CompanyMap[companyName]; tCompanyName != "" {
		companyName = tCompanyName
	}

	if channel != models.NetworkChannel || bg.networkNode == nil || bg.networkNode.IsMgaProponent {
		section2Info = fmt.Sprintf(
			mgaProponentFormat,
			companyName,
		)
	} else {
		worksForNode := bg.networkNode
		if bg.networkNode.WorksForUid != "" {
			worksForNode = bg.worksForNode
		}
		section2Info = fmt.Sprintf(
			mgaEmitterFormat,
			worksForNode.Agency.Name,
			companyName,
		)
	}

	section5Info = withoutConsultacy
	if consultancyPrice > 0 {
		section5Info = fmt.Sprintf(
			withConsultacyFormat,
			lib.HumanaizePriceEuro(consultancyPrice),
		)
	}

	return section2Info, section5Info
}

func (bg *baseGenerator) mupProducerTable(producerInfo, proponentInfo map[string]string, designation string) {
	type entry struct {
		title string
		body  string
	}

	table := make([][]domain.TableCell, 0)

	parseEntries := func(entries []entry, last bool) [][]domain.TableCell {
		result := make([][]domain.TableCell, 0)
		borders := []string{"T", ""}
		for index, e := range entries {
			if last && index == len(entries)-1 {
				borders = []string{"T", "B"}
			}
			row := [][]domain.TableCell{
				{
					{
						Text:      e.title,
						Height:    5,
						Width:     constants.FullPageWidth,
						FontSize:  constants.SmallFontSize,
						FontStyle: constants.RegularFontStyle,
						FontColor: constants.BlackColor,
						Fill:      false,
						FillColor: domain.Color{},
						Align:     constants.LeftAlign,
						Border:    borders[0],
					},
				},
				{
					{
						Text:      e.body,
						Height:    5,
						Width:     constants.FullPageWidth,
						FontSize:  constants.RegularFontSize,
						FontStyle: constants.RegularFontStyle,
						FontColor: constants.BlackColor,
						Fill:      false,
						FillColor: domain.Color{},
						Align:     constants.LeftAlign,
						Border:    borders[1],
					},
				},
			}
			result = append(result, row...)
		}
		return result
	}

	entries := []entry{
		{
			title: "DATI SULL’IDENTIFICAZIONE DELL’INTERMEDIARIO",
			body: producerInfo["name"] + " iscritto alla Sezione " +
				producerInfo["ruiSection"] + " del RUI con numero " + producerInfo["ruiCode"] + " in data " +
				producerInfo["ruiRegistration"],
		},
		{
			title: "QUALIFICA",
			body:  designation,
		},
		{
			title: "SEDE LEGALE",
			body:  proponentInfo["address"],
		},
	}

	table = append(table, parseEntries(entries, false)...)

	table = append(table, [][]domain.TableCell{
		{
			{
				Text:      "RECAPITI TELEFONICI",
				Height:    5,
				Width:     95,
				FontSize:  constants.SmallFontSize,
				FontStyle: constants.RegularFontStyle,
				FontColor: constants.BlackColor,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "T",
			},
			{
				Text:      "E-MAIL",
				Height:    5,
				Width:     95,
				FontSize:  constants.SmallFontSize,
				FontStyle: constants.RegularFontStyle,
				FontColor: constants.BlackColor,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "T",
			},
		},
		{
			{
				Text:      proponentInfo["phone"],
				Height:    5,
				Width:     95,
				FontSize:  constants.RegularFontSize,
				FontStyle: constants.RegularFontStyle,
				FontColor: constants.BlackColor,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "",
			},
			{
				Text:      proponentInfo["email"],
				Height:    5,
				Width:     95,
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
				Text:      "PEC",
				Height:    5,
				Width:     95,
				FontSize:  constants.SmallFontSize,
				FontStyle: constants.RegularFontStyle,
				FontColor: constants.BlackColor,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "T",
			},
			{
				Text:      "SITO INTERNET",
				Height:    5,
				Width:     95,
				FontSize:  constants.SmallFontSize,
				FontStyle: constants.RegularFontStyle,
				FontColor: constants.BlackColor,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "T",
			},
		},
		{
			{
				Text:      proponentInfo["pec"],
				Height:    5,
				Width:     95,
				FontSize:  constants.RegularFontSize,
				FontStyle: constants.RegularFontStyle,
				FontColor: constants.BlackColor,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "",
			},
			{
				Text:      proponentInfo["website"],
				Height:    5,
				Width:     95,
				FontSize:  constants.RegularFontSize,
				FontStyle: constants.RegularFontStyle,
				FontColor: constants.BlackColor,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "",
			},
		},
	}...)

	entries = []entry{
		{
			title: "AUTORITÀ COMPETENTE ALLA VIGILANZA DELL’ATTIVITÀ SVOLTA",
			body: "IVASS - Istituto per la Vigilanza sulle Assicurazioni - Via del Quirinale, " +
				"21 - 00187 Roma",
		},
	}
	table = append(table, parseEntries(entries, true)...)

	bg.engine.SetDrawColor(constants.PinkColor)
	bg.engine.DrawTable(table)
	bg.engine.NewLine(2)
}
