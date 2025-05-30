package contract

import (
	"fmt"
	"math"
	"strings"
	"time"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"gitlab.dev.wopta.it/goworkspace/document/internal/constants"
	"gitlab.dev.wopta.it/goworkspace/document/internal/domain"
	"gitlab.dev.wopta.it/goworkspace/document/internal/engine"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/network"
)

const (
	tabDimension = 15
)

type baseGenerator struct {
	engine      *engine.Fpdf
	isProposal  bool
	now         time.Time
	signatureID uint32
	networkNode *models.NetworkNode
	policy      *models.Policy
}

func (bg *baseGenerator) emptyHeader() {
	bg.engine.SetHeader(func() {
		if bg.isProposal {
			bg.engine.DrawWatermark(constants.Proposal)
		}
	})
}

func (bg *baseGenerator) emptyFooter() {
	bg.engine.SetFooter(func() {})
}

func (bg *baseGenerator) woptaHeader() {
	bg.engine.SetHeader(func() {
		bg.engine.InsertImage(lib.GetAssetPathByEnvV2()+"logo_wopta.png", 10, 6, 0, 10)
		bg.engine.NewLine(10)

		if bg.isProposal {
			bg.engine.DrawWatermark(constants.Proposal)
		}
	})
}

func (bg *baseGenerator) woptaFooter() {
	const (
		rowHeight   = 3
		columnWidth = 50
	)

	bg.engine.SetFooter(func() {
		bg.engine.SetY(-30)

		currentY := bg.engine.GetY()

		bg.engine.DrawLine(11, currentY, 200, currentY, constants.RegularThickness, constants.PinkColor)
		bg.engine.NewLine(3)

		entries := [][]string{
			{"Wopta Assicurazioni s.r.l", " ", " ", "www.wopta.it"},
			{"Galleria del Corso, 1", "Numero REA: MI 2638708", "CF | P.IVA | n. iscr. Registro Imprese:",
				"info@wopta.it"},
			{"20122 - Milano (MI)", "Capitale Sociale: € 204.839,26 i.v.", "12072020964", "(+39) 02 91240346"},
		}

		table := make([][]domain.TableCell, 0, 3)

		for index, entry := range entries {
			textColor := constants.BlackColor
			if index == 0 {
				textColor = constants.PinkColor
			}
			row := make([]domain.TableCell, 0, 4)

			for _, cell := range entry {
				row = append(row, domain.TableCell{
					Text:      cell,
					Height:    rowHeight,
					Width:     columnWidth,
					FontSize:  constants.SmallFontSize,
					FontStyle: constants.RegularFontStyle,
					FontColor: textColor,
					Fill:      false,
					FillColor: domain.Color{},
					Align:     constants.LeftAlign,
					Border:    "",
				})
			}
			table = append(table, row)
		}

		bg.engine.DrawTable(table)

		bg.engine.NewLine(3)

		bg.engine.WriteText(domain.TableCell{
			Text: "Wopta Assicurazioni s.r.l. è un intermediario assicurativo soggetto alla vigilanza dell’IVASS" +
				" ed iscritto alla Sezione A del Registro Unico degli Intermediari Assicurativi con numero" +
				" A000701923. Consulta gli estremi dell’iscrizione al sito https://servizi.ivass.it/RuirPubblica/",
			Height:    rowHeight,
			Width:     constants.FullPageWidth,
			FontSize:  constants.SmallFontSize,
			FontStyle: constants.RegularFontStyle,
			FontColor: constants.BlackColor,
			Fill:      false,
			FillColor: domain.Color{},
			Align:     constants.LeftAlign,
			Border:    "",
		})

		bg.engine.SetY(-7)

		bg.engine.WriteText(domain.TableCell{
			Text:      fmt.Sprintf("%d", bg.engine.PageNumber()),
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

func (bg *baseGenerator) woptaPrivacySection() {
	const (
		rowHeight   = 3
		columnWidth = 190
	)

	type section struct {
		title       string
		subsections []string
	}

	bg.engine.WriteText(domain.TableCell{
		Text:      "COME RISPETTIAMO LA TUA PRIVACY",
		Height:    rowHeight,
		Width:     columnWidth,
		FontSize:  constants.LargeFontSize,
		FontStyle: constants.BoldFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.CenterAlign,
		Border:    "",
	})

	bg.engine.NewLine(rowHeight)

	bg.engine.WriteText(domain.TableCell{
		Text:      "Informativa sul trattamento dei dati personali",
		Height:    rowHeight,
		Width:     columnWidth,
		FontSize:  constants.RegularFontSize,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "",
	})

	bg.engine.NewLine(1)

	bg.engine.WriteText(domain.TableCell{
		Text: "Ai sensi del REGOLAMENTO (UE) 2016/679 (" +
			"relativo alla protezione delle persone fisiche con riguardo al trattamento dei dati personali, " +
			"nonché alla libera circolazione di tali dati) si informa l’“Interessato” (" +
			"contraente / aderente alla polizza collettiva o convenzione / assicurato / beneficiario / loro aventi" +
			" causa) di quanto segue.",
		Height:    rowHeight,
		Width:     columnWidth,
		FontSize:  constants.RegularFontSize,
		FontStyle: constants.BoldFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     "",
		Border:    "",
	})

	bg.engine.NewLine(rowHeight)

	bg.engine.WriteText(domain.TableCell{
		Text:      "1. TITOLARE DEL TRATTAMENTO",
		Height:    rowHeight,
		Width:     columnWidth,
		FontSize:  constants.LargeFontSize,
		FontStyle: constants.BoldFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "",
	})
	bg.engine.NewLine(1)
	bg.engine.WriteText(domain.TableCell{
		Text: "Titolare del trattamento è Wopta Assicurazioni, con sede legale in Milano, Galleria del Corso, " +
			"1 (di seguito “Titolare”), raggiungibile all’indirizzo e-mail: privacy@wopta.it",
		Height:    rowHeight,
		Width:     columnWidth,
		FontSize:  constants.RegularFontSize,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "",
	})

	bg.engine.NewLine(rowHeight)

	bg.engine.WriteText(domain.TableCell{
		Text:      "2. I DATI PERSONALI OGGETTO DI TRATTAMENTO, FINALITÀ E BASE GIURIDICA",
		Height:    rowHeight,
		Width:     columnWidth,
		FontSize:  constants.LargeFontSize,
		FontStyle: constants.BoldFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "",
	})
	bg.engine.NewLine(1)
	bg.engine.WriteText(domain.TableCell{
		Text:      "a) Finalità Contrattuali, normative, amministrative e giudiziali",
		Height:    rowHeight,
		Width:     columnWidth,
		FontSize:  constants.RegularFontSize,
		FontStyle: constants.BoldFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "",
	})
	bg.engine.NewLine(1)
	bg.engine.WriteText(domain.TableCell{
		Text: "Fermo restando quanto previsto dalla Privacy & Cookie Policy del Sito, ove " +
			"applicabile, i dati così conferiti potranno essere trattati, anche con strumenti elettronici, da parte del " +
			"Titolare per eseguire le prestazioni contrattuali, in qualità di intermediario, richieste dall’interessato, " +
			"o per adempiere ad obblighi normativi, contabili e fiscali, ovvero ancora per finalità di difesa in " +
			"giudizio, per il tempo strettamente necessario a tali attività.",
		Height:    rowHeight,
		Width:     columnWidth,
		FontSize:  constants.RegularFontSize,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "",
	})
	bg.engine.NewLine(1)
	bg.engine.WriteText(domain.TableCell{
		Text: "La base giuridica del trattamento di dati personali per le finalità di cui sopra " +
			"è l’art. 6.1 lett. b), c), f) del Regolamento in quanto i trattamenti sono necessari all'erogazione dei " +
			"servizi o per il riscontro di richieste dell’interessato, in conformità a quanto previsto dall’incarico " +
			"conferito all’intermediario, nonché ove il trattamento risulti necessario per l’adempimento di un preciso " +
			"obbligo di legge posto in capo al Titolare, o al fine di accertare, esercitare o difendere un diritto in " +
			"sede giudiziaria. Il conferimento dei dati personali per queste finalità è facoltativo, ma l'eventuale " +
			"mancato conferimento comporterebbe l'impossibilità per l’intermediario di eseguire le proprie obbligazioni " +
			"contrattuali.",
		Height:    rowHeight,
		Width:     columnWidth,
		FontSize:  constants.RegularFontSize,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "",
	})
	bg.engine.NewLine(1)
	bg.engine.WriteText(domain.TableCell{
		Text:      "b) Finalità commerciali",
		Height:    rowHeight,
		Width:     columnWidth,
		FontSize:  constants.RegularFontSize,
		FontStyle: constants.BoldFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "",
	})
	bg.engine.NewLine(1)
	bg.engine.WriteText(domain.TableCell{
		Text: "Inoltre, i Suoi dati personali potranno essere trattati al fine di inviarLe " +
			"comunicazioni e proposte commerciali, incluso l’invio di newsletter e ricerche di mercato, attraverso " +
			"strumenti automatizzati (sms, mms, email, messaggistica istantanea e chat) e non (posta cartacea, telefono); " +
			"si precisa che il Titolare raccoglie un unico consenso per le finalità di marketing qui descritte, ai sensi " +
			"del Provvedimento Generale del Garante per la Protezione dei Dati Personali \"Linee guida in materia di " +
			"attività promozionale e contrasto allo spam” del 4 luglio 2013; qualora, in ogni caso, Lei desiderasse " +
			"opporsi al trattamento dei Suoi dati per le finalità di marketing eseguito con i mezzi qui indicati, potrà " +
			"in qualunque momento farlo contattando il Titolare ai recapiti indicati nella sezione \"Contatti\" di " +
			"questa informativa, senza pregiudicare la liceità del trattamento effettuato prima dell’opposizione.",
		Height:    rowHeight,
		Width:     columnWidth,
		FontSize:  constants.RegularFontSize,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "",
	})
	bg.engine.NewLine(1)
	bg.engine.WriteText(domain.TableCell{
		Text: "I trattamenti eseguiti per la finalità di marketing, di cui al paragrafo che " +
			"precede, si basa sul rilascio del Suo consenso ai sensi dell’art. 6, par. 1, lett. a) ([…] l'interessato ha " +
			"espresso il consenso al trattamento dei propri dati personali per una o più specifiche finalità) del " +
			"Regolamento. Tale consenso è revocabile in qualsiasi momento senza pregiudizio alcuno della liceità del " +
			"trattamento effettuato anteriormente alla revoca in conformità a quanto previsto dall’art. 7 del " +
			"Regolamento. Il conferimento dei Suoi dati personali per queste finalità è quindi del tutto facoltativo e " +
			"non pregiudica la fruizione dei servizi. Qualora desiderasse opporsi al trattamento dei Suoi dati per le " +
			"finalità di marketing, potrà in qualunque momento farlo contattando il Titolare ai recapiti indicati nella " +
			"sezione \"Contatti\" di questa informativa.",
		Height:    rowHeight,
		Width:     columnWidth,
		FontSize:  constants.RegularFontSize,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "",
	})

	bg.engine.NewLine(rowHeight)

	sections := []section{
		{
			title: "3. DESTINATARI DEI DATI PERSONALI",
			subsections: []string{
				"I Suoi dati personali potranno essere condivisi, " +
					"per le finalità di cui alla sezione 2 della presente Policy, con:",
				"- Soggetti che agiscono tipicamente in qualità di Responsabili del trattamento " +
					"ex art. 28 del Regolamento per conto del Titolare, incaricati dell'erogazione dei Servizi (" +
					"a titolo esemplificativo: servizi tecnologici, " +
					"servizi di assistenza e consulenza in materia contabile, amministrativa, legale, " +
					"tributaria e finanziaria, manutenzione tecnica). Il Titolare conserva una lista aggiornata dei " +
					"responsabili del trattamento nominati e ne garantisce la presa visione all’interessato presso la" +
					" sede sopra indicata o previa richiesta indirizzata ai recapiti sopra indicati;",
				"- Persone autorizzate dal Titolare al trattamento dei dati personali ai sensi " +
					"degli artt. 29 e 2-quaterdecies del D.lgs. n. 196/2003 (“Codice “Privacy”) (ad es. " +
					"il personale dipendente addetto alla manutenzione del Sito, alla gestione del CRM, " +
					"alla gestione dei sistemi informativi ecc.);",
				"- Soggetti terzi, autonomi titolari del trattamento, a cui i dati potrebbero " +
					"essere trasmessi al fine di dare seguito a specifici servizi da Lei richiesti e/o  per dare" +
					" esecuzione alle attività di cui alla presente informativa, " +
					"e con i quali il Titolare abbia stipulato accordi commerciali; soggetti, " +
					"quali le imprese di assicurazione, che assumono il rischio di sottoscrizione della polizza, ai " +
					"quali sia obbligatorio comunicare i tuoi Dati personali in forza di obblighi contrattuali e di" +
					" disposizioni di legge e regolamentari sulla distribuzione di prodotti assicurativi;",
				"- Soggetti, enti od autorità a cui sia obbligatorio comunicare i Suoi dati personali in forza di" +
					" disposizioni di legge o di ordini delle autorità.",
				"Tali soggetti sono, di seguito, collettivamente definiti come “Destinatari”. " +
					"L'elenco completo dei responsabili del trattamento è disponibile inviando una richiesta scritta" +
					" al Titolare ai recapiti indicati nella sezione \"Contatti\" di questa informativa.",
			},
		},
		{
			title: "4. TRASFERIMENTI DEI DATI PERSONALI",
			subsections: []string{
				"Alcuni dei Suoi dati personali sono condivisi con Destinatari che si potrebbero " +
					"trovare al di fuori dello Spazio Economico Europeo. " +
					"Il Titolare assicura che il trattamento Suoi dati personali da parte di questi Destinatari" +
					" avviene nel rispetto degli artt. 44 - 49 del Regolamento. Invero, " +
					"per quanto concerne il trasferimento dei dati personali verso Paesi terzi, " +
					"il Titolare rende noto che il trattamento avverrà secondo una delle modalità consentite dalla" +
					" legge vigente, quali, ad esempio, il consenso dell’interessato, " +
					"l’adozione di Clausole Standard approvate dalla Commissione Europea, " +
					"la selezione di soggetti aderenti a programmi internazionali per la libera circolazione dei dati" +
					" o operanti in Paesi considerati sicuri dalla Commissione Europea sulla base di una decisione di" +
					" adeguatezza.",
				"Maggiori informazioni sono disponibili inviando una richiesta scritta al " +
					"Titolare ai recapiti indicati nella sezione \"Contatti\" di questa informativa.",
			},
		},
		{
			title: "5. CONSERVAZIONE DEI DATI PERSONALI",
			subsections: []string{
				"I Suoi dati personali saranno inseriti e conservati, in conformità ai principi " +
					"di minimizzazione e limitazione della conservazione di cui all’art. 5.1." +
					"c) ed e) del Regolamento, nei sistemi informativi del Titolare, " +
					"i cui server sono situati all’interno dello Spazio Economico Europeo.",
				"I dati personali trattati per le finalità di cui alle lettere a) e b) " +
					"saranno conservati per il tempo strettamente necessario a raggiungere quelle stesse finalit" +
					"à ovverossia per il tempo necessario all’esecuzione del contratto, " +
					"in conformità ai tempi di conservazione obbligatori per legge (vedi anche, in particolare, " +
					"art. 2946 c.c. e ss.).",
				"Per le finalità di cui alla lettera c), i suoi dati personali saranno invece " +
					"trattati fino alla revoca del suo consenso. Alla revoca del consenso, " +
					"i dati trattati per la finalità di cui sopra verranno cancellati o resi anonimi in modo" +
					" permanente.",
				"In generale, il Titolare si riserva in ogni caso di conservare i Suoi dati per " +
					"il tempo necessario ad adempiere ogni eventuale obbligo normativo cui lo stesso è soggetto o per" +
					" soddisfare eventuali esigenze difensive. " +
					"Resta infatti salva la possibilità per il Titolare di conservare i Suoi dati personali per il" +
					" periodo di tempo previsto e ammesso dalla legge Italiana a tutela dei propri interessi " +
					"(Art. 2947 c.c.).",
				"Maggiori informazioni in merito al periodo di conservazione dei dati e ai " +
					"criteri utilizzati per determinare tale periodo possono essere richieste inviando una richiesta" +
					" scritta al Titolare ai recapiti indicati nella sezione \"Contatti\" di questa informativa.",
			},
		},
		{
			title: "6. DIRITTI DELL’INTERESSATO",
			subsections: []string{
				"Lei ha il diritto di accedere in qualunque momento ai Dati Personali che La " +
					"riguardano, ai sensi degli artt. 15-22 del Regolamento. In particolare, " +
					"potrà chiedere la rettifica (ex art. 16), la cancellazione (ex art. 17), " +
					"la limitazione (ex art. 18) e la portabilità dei dati (ex art. 20), " +
					"di non essere sottoposto a una decisione basata unicamente sul trattamento automatizzato, " +
					"compresa la profilazione, " +
					"che produca effetti giuridici che La riguardano o che incida in modo analogo " +
					"significativamente sulla sua persona (ex art. 22), " +
					"nonché la revoca del consenso eventualmente prestato (ex art. 7, par. 3).",
				"Lei può formulare, inoltre, una richiesta di opposizione al trattamento dei " +
					"Suoi Dati Personali ex art. 21 del Regolamento nella quale dare evidenza delle ragioni che" +
					" giustifichino l’opposizione: il titolare si riserva di valutare la Sua istanza, " +
					"che non verrebbe accettata in caso di esistenza di motivi legittimi cogenti per procedere al" +
					" trattamento che prevalgano sui Suoi interessi, " +
					"diritti e libertà. Lei ha altresì il diritto di opporsi in ogni momento e senza" +
					" alcuna giustificazione all’invio di marketing diretto attraverso strumenti automatizzati (es. " +
					"sms, mms, e-mail, notifiche push, fax, " +
					"sistemi di chiamata automatizzati senza operatore) e non (posta cartacea, telefono con operatore). " +
					"Inoltre, con riguardo al marketing diretto, " +
					"resta salva la possibilità di esercitare tale diritto anche in parte, ossia, in tal caso, " +
					"opponendosi, ad esempio, " +
					"al solo invio di comunicazioni promozionali effettuato tramite strumenti automatizzati.",
				"Le richieste vanno rivolte per iscritto al" +
					" Titolare ai recapiti indicati nella sezione \"Contatti\" di questa informativa.",
				"Qualora Lei ritenga che il trattamento dei Suoi Dati personali effettuato dal " +
					"Titolare avvenga in violazione di quanto previsto dal GDPR, " +
					"ha il diritto di proporre reclamo al Garante Privacy, " +
					"come previsto dall'art. 77 del GDPR stesso, o di adire le opportune sedi giudiziarie " +
					"(art. 79 del GDPR).",
			},
		},
		{
			title: "7. CONTATTI",
			subsections: []string{
				"Per esercitare i diritti di cui sopra o per qualunque altra richiesta può " +
					"scrivere al Titolare del trattamento all’indirizzo: privacy@wopta.it.",
			},
		},
	}

	for index, s := range sections {
		if index == 2 {
			bg.engine.NewPage()
		}
		bg.engine.WriteText(domain.TableCell{
			Text:      s.title,
			Height:    rowHeight,
			Width:     columnWidth,
			FontSize:  constants.LargeFontSize,
			FontStyle: constants.BoldFontStyle,
			FontColor: constants.BlackColor,
			Fill:      false,
			FillColor: domain.Color{},
			Align:     constants.LeftAlign,
			Border:    "",
		})
		bg.engine.NewLine(1)
		for _, sub := range s.subsections {
			bg.engine.WriteText(domain.TableCell{
				Text:      sub,
				Height:    rowHeight,
				Width:     columnWidth,
				FontSize:  constants.RegularFontSize,
				FontStyle: constants.RegularFontStyle,
				FontColor: constants.BlackColor,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "",
			})
			bg.engine.NewLine(1)
		}
		bg.engine.NewLine(rowHeight)
	}
}

func (bg *baseGenerator) commercialConsentSection() {
	const (
		key = 2
	)

	var (
		consentText, notConsentText = " ", "X"
	)

	if bg.policy.Contractor.Consens != nil {
		consent, err := bg.policy.ExtractConsens(key)
		if err != nil {
			log.Println("Error extracting consens, ", key, err)
			panic(err)
		}

		if consent.Answer {
			consentText = "X"
			notConsentText = " "
		}
	}

	bg.engine.SetDrawColor(constants.BlackColor)
	bg.engine.WriteText(domain.TableCell{
		Text:      "Consenso per finalità commerciali.",
		Height:    3,
		Width:     constants.FullPageWidth,
		FontSize:  constants.LargeFontSize,
		FontStyle: constants.BoldFontStyle,
		FontColor: constants.PinkColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "",
	})
	bg.engine.NewLine(1)
	bg.engine.WriteText(domain.TableCell{
		Text:      "Il sottoscritto, letta e compresa l’informativa sul trattamento dei dati personali",
		Height:    3,
		Width:     constants.FullPageWidth,
		FontSize:  constants.RegularFontSize,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "",
	})
	bg.engine.NewLine(2)
	bg.engine.DrawTable([][]domain.TableCell{
		{
			{
				Text:      " ",
				Height:    4.5,
				Width:     5,
				FontSize:  constants.RegularFontSize,
				FontStyle: constants.RegularFontStyle,
				FontColor: constants.BlackColor,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "",
			},
			{
				Text:      consentText,
				Height:    4.5,
				Width:     4.5,
				FontSize:  constants.RegularFontSize,
				FontStyle: constants.RegularFontStyle,
				FontColor: constants.BlackColor,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.CenterAlign,
				Border:    "1",
			},
			{
				Text:      "ACCONSENTE",
				Height:    4.5,
				Width:     30,
				FontSize:  constants.RegularFontSize,
				FontStyle: constants.RegularFontStyle,
				FontColor: constants.BlackColor,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign + "M",
				Border:    "",
			},
			{
				Text:      " ",
				Height:    4.5,
				Width:     20,
				FontSize:  constants.RegularFontSize,
				FontStyle: constants.RegularFontStyle,
				FontColor: constants.BlackColor,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "",
			},
			{
				Text:      notConsentText,
				Height:    4.5,
				Width:     4.5,
				FontSize:  constants.RegularFontSize,
				FontStyle: constants.RegularFontStyle,
				FontColor: constants.BlackColor,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.CenterAlign,
				Border:    "1",
			},
			{
				Text:      "NON ACCONSENTE",
				Height:    4.5,
				Width:     125,
				FontSize:  constants.RegularFontSize,
				FontStyle: constants.RegularFontStyle,
				FontColor: constants.BlackColor,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign + "M",
				Border:    "",
			},
		},
	})
	bg.engine.NewLine(2)
	bg.engine.WriteText(domain.TableCell{
		Text: "al trattamento dei propri dati personali da parte di Wopta Assicurazioni per " +
			"l’invio di comunicazioni e proposte commerciali e di marketing, incluso l’invio di newsletter e ricerche di " +
			"mercato, attraverso strumenti automatizzati (sms, mms, e-mail, ecc.) e non (posta cartacea e telefono " +
			"con operatore).",
		Height:    3,
		Width:     constants.FullPageWidth,
		FontSize:  constants.RegularFontSize,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "",
	})
	bg.engine.NewLine(3)
	bg.engine.WriteText(domain.TableCell{
		Text:      bg.now.Format(constants.DayMonthYearFormat),
		Height:    3,
		Width:     30,
		FontSize:  constants.RegularFontSize,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "",
	})
	bg.signatureForm()
}

func (bg *baseGenerator) signatureForm() {
	if bg.isProposal {
		return
	}
	text := fmt.Sprintf("\"[[!sigField\"%d\":signer1:signature(sigType=\\\"Click2Sign\\\"):label"+
		"(\\\"firma qui\\\"):size(width=150,height=60)]]\"", bg.signatureID)

	bg.engine.SetX(-90)
	bg.engine.WriteText(domain.TableCell{
		Text:      "Firma del Contraente",
		Height:    3,
		Width:     50,
		FontSize:  constants.RegularFontSize,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.RightAlign,
		Border:    "",
	})
	bg.engine.NewLine(15)
	currentY := bg.engine.GetY()
	bg.engine.DrawLine(120, currentY, 190, currentY, constants.ThinThickness, constants.BlackColor)
	bg.engine.NewLine(2)
	bg.engine.SetX(-135)
	bg.engine.WriteText(domain.TableCell{
		Text:      text,
		Height:    3,
		Width:     130,
		FontSize:  constants.SmallFontSize,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.RightAlign,
		Border:    "",
	})
	bg.engine.NewLine(7.5)

	bg.signatureID++
}

func (bg *baseGenerator) AddMup() {
	bg.mup()
}

func (bg *baseGenerator) mup() {
	if bg.networkNode != nil && !bg.networkNode.HasAnnex && bg.networkNode.Type != models.PartnershipNetworkNodeType {
		return
	}

	if bg.networkNode == nil || bg.networkNode.Type == models.PartnershipNetworkNodeType || bg.networkNode.IsMgaProponent {
		bg.woptaHeader()

		bg.woptaFooter()
	} else {
		bg.emptyHeader()

		bg.emptyFooter()
	}

	producerInfo := bg.productInfo()
	proponentInfo := bg.proponentInfo()
	designationInfo := bg.designationInfo()
	mupSection2Info, mupSection5Info := bg.mupInfo()

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
	bg.engine.NewPage()
	bg.mupSectionVII()
}

func (bg *baseGenerator) mupTitle() {
	text := "ALLEGATO 3 - MODELLO UNICO PRECONTRATTUALE (MUP) PER I PRODOTTI ASSICURATIVI"
	bg.engine.WriteText(bg.engine.GetTableCell(text, constants.BoldFontStyle, constants.CenterAlign))

	bg.engine.NewLine(3)

	text = "Il distributore ha l’obbligo di consegnare/trasmettere al contraente il presente Modulo, " +
		"prima della sottoscrizione della proposta o del contratto di assicurazione. Il documento può " +
		"essere fornito con modalità non cartacea se appropriato rispetto alle modalità di distribuzione " +
		"del prodotto assicurativo e il contraente lo consente (art. 120-quater del " +
		"Codice delle Assicurazioni Private)"
	bg.engine.WriteText(bg.engine.GetTableCell(text, constants.RegularFontStyle))
}

func (bg *baseGenerator) mupSectionI(producerInfo, proponentInfo map[string]string, designation string) {
	text := "SEZIONE I - Informazioni generali sul distributore che entra in contatto con " +
		"il contraente"
	bg.engine.WriteText(bg.engine.GetTableCell(text, constants.BoldFontStyle))

	bg.woptaTable(producerInfo, proponentInfo, designation)

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

	text = "Nello svolgimento dell’attività di distribuzione, l’intermediario non presta " +
		"attività di consulenza prima della conclusione del contratto né fornisce al contraente una raccomandazione " +
		"personalizzata ai sensi dell’art. 119-ter, comma 3, del decreto legislativo n. 209/2005 " +
		"(Codice delle Assicurazioni Private)." + "\n" +
		"L'attività di distribuzione assicurativa è svolta in assenza di obblighi " +
		"contrattuali che impongano di offrire esclusivamente i contratti di una o più imprese di " +
		"assicurazioni."
	bg.engine.WriteText(bg.engine.GetTableCell(text, constants.RegularFontStyle))
}

func (bg *baseGenerator) mupSectionV(body string) {
	text := "SEZIONE V - Informazioni relative alle remunerazioni"
	bg.engine.WriteText(bg.engine.GetTableCell(text, constants.BoldFontStyle))

	text = body + "\n" +
		"L’informazione sopra resa riguarda i compensi complessivamente percepiti da tutti " +
		"gli intermediari coinvolti nella distribuzione del prodotto."
	bg.engine.WriteText(bg.engine.GetTableCell(text, constants.RegularFontStyle))
}

func (bg *baseGenerator) mupSectionVI() {
	text := "SEZIONE VI – Informazioni sul pagamento dei premi"
	bg.engine.WriteText(bg.engine.GetTableCell(text, constants.BoldFontStyle))

	text = "Relativamente a questo contratto i premi pagati dal Contraente " +
		"all’intermediario e le somme destinate ai risarcimenti o ai pagamenti dovuti dalle Imprese di Assicurazione, " +
		"se regolati per il tramite dell’intermediario costituiscono patrimonio autonomo e separato dal patrimonio " +
		"dello stesso."
	bg.engine.WriteText(bg.engine.GetTableCell(text, constants.RegularFontStyle))

	text = "Indicare le modalità di pagamento ammesse"
	bg.engine.WriteText(bg.engine.GetTableCell(text, constants.BoldFontStyle))

	text = "Sono consentiti, nei confronti dell'intermediario, esclusivamente bonifico e strumenti di " +
		"pagamento elettronico, quali ad esempio, carte di credito e/o carte di debito, incluse le carte " +
		"prepagate."
	bg.engine.WriteText(bg.engine.GetTableCell(text, constants.RegularFontStyle))
}

func (bg *baseGenerator) mupSectionVII() {
	text := "SEZIONE VII - Informazioni sugli strumenti di tutela del contraente"
	bg.engine.WriteText(bg.engine.GetTableCell(text, constants.BoldFontStyle))

	text = "L’attività di distribuzione è garantita da un contratto di assicurazione della " +
		"responsabilità civile che copre i danni arrecati ai contraenti da negligenze ed errori professionali " +
		"dell’intermediario o da negligenze, errori professionali ed infedeltà dei dipendenti, dei collaboratori o " +
		"delle persone del cui operato l’intermediario deve rispondere a norma di legge." + "\n" +
		"Il contraente ha la facoltà, ferma restando la possibilità di rivolgersi " +
		"all’Autorità Giudiziaria, di inoltrare reclamo per iscritto all’intermediario, via posta all’indirizzo di " +
		"sede legale o a mezzo mail alla PEC sopra indicati, oppure all’Impresa secondo le modalità e presso i " +
		"recapiti indicati nel DIP aggiuntivo nella relativa sezione, nonché la possibilità, qualora non dovesse " +
		"ritenersi soddisfatto dall’esito del reclamo o in caso di assenza di riscontro da parte dell’intermediario " +
		"o dell’impresa entro il termine di legge, di rivolgersi all’IVASS secondo quanto indicato nei DIP aggiuntivi." + "\n" +
		"Il contraente ha la facoltà di avvalersi di altri eventuali sistemi alternativi " +
		"di risoluzione delle controversie previsti dalla normativa vigente nonché quelli indicati nei" +
		" DIP aggiuntivi."
	bg.engine.WriteText(bg.engine.GetTableCell(text, constants.RegularFontStyle))
}

func (bg *baseGenerator) productInfo() map[string]string {
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

// TODO: private
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
			proponentNode = network.GetNetworkNodeByUid(proponentNode.WorksForUid)
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

func (bg *baseGenerator) designationInfo() string {
	var (
		designation                           string
		mgaProponentDirectDesignationFormat   = "%s %s"
		mgaRuiInfo                            = "Wopta Assicurazioni Srl, Società iscritta alla Sezione A del RUI con numero A000701923 in data 14/02/2022"
		designationDirectManager              = "Responsabile dell’attività di intermediazione assicurativa di"
		mgaProponentIndirectDesignationFormat = "%s di %s, iscritta in sezione E del RUI con numero %s in data %s, che opera per conto di %s"
		mgaEmitterDesignationFormat           = "%s dell’intermediario di %s iscritta alla sezione %s del RUI con numero %s in data %s"
	)

	if bg.networkNode == nil || bg.networkNode.Type == models.PartnershipNetworkNodeType {
		designation = fmt.Sprintf(mgaProponentDirectDesignationFormat, designationDirectManager, mgaRuiInfo)
	} else if bg.networkNode.IsMgaProponent {
		if bg.networkNode.WorksForUid == models.WorksForMgaUid {
			designation = fmt.Sprintf(mgaProponentDirectDesignationFormat, bg.networkNode.Designation, mgaRuiInfo)
		} else {
			worksForNode := bg.networkNode
			if bg.networkNode.WorksForUid != "" {
				worksForNode = network.GetNetworkNodeByUid(bg.networkNode.WorksForUid)
			}
			designation = fmt.Sprintf(
				mgaProponentIndirectDesignationFormat,
				bg.networkNode.Designation,
				worksForNode.Agency.Name,
				worksForNode.Agency.RuiCode,
				worksForNode.Agency.RuiRegistration.Format(constants.DayMonthYearFormat),
				mgaRuiInfo,
			)
		}
	} else {
		worksForNode := bg.networkNode
		if bg.networkNode.WorksForUid != "" {
			worksForNode = network.GetNetworkNodeByUid(bg.networkNode.WorksForUid)
		}
		designation = fmt.Sprintf(
			mgaEmitterDesignationFormat,
			bg.networkNode.Designation,
			worksForNode.Agency.Name,
			worksForNode.Agency.RuiSection,
			worksForNode.Agency.RuiCode,
			worksForNode.Agency.RuiRegistration.Format(constants.DayMonthYearFormat),
		)
	}

	log.Printf("designation info: %+v", designation)

	return designation
}

func (bg *baseGenerator) mupInfo() (section1Info, section3Info string) {
	const (
		mgaProponentFormat = "Secondo quanto indicato nel modulo di proposta/polizza e documentazione " +
			"precontrattuale ricevuta, la distribuzione  relativamente a questa proposta/contratto è svolta per " +
			"conto della seguente impresa di assicurazione: %s"
		mgaEmitterFormat = "Il contratto viene intermediato da %s, in qualità di soggetto proponente, che opera in " +
			"virtù della collaborazione con Wopta Assicurazioni Srl (intermediario emittente dell'Impresa di " +
			"Assicurazione %s, iscritto al RUI sezione A nr A000701923 dal 14.02.2022, ai sensi dell’articolo 22, " +
			"comma 10, del decreto legge 18 ottobre 2012, n. 179, convertito nella legge 17 dicembre 2012, n. 221"
		withoutConsultacy = "Per il prodotto intermediato, è corrisposto all’intermediario, da parte " +
			"dell’impresa di assicurazione, un compenso sotto forma di commissione inclusa nel premio " +
			"assicurativo."
		withConsultacyFormat = "Per il prodotto intermediato, è corrisposto un compenso all’intermediario, da parte " +
			"dell’impresa di assicurazione, sotto forma di commissione inclusa nel premio assicurativo " +
			"e un contributo per servizi di intermediazione, a carico del cliente, pari ad %s."
	)

	companyName := constants.CompanyMap[bg.policy.Company]

	if bg.policy.Channel != models.NetworkChannel || bg.networkNode == nil || bg.networkNode.IsMgaProponent {
		section1Info = fmt.Sprintf(
			mgaProponentFormat,
			companyName,
		)
	} else {
		worksForNode := bg.networkNode
		if bg.networkNode.WorksForUid != "" {
			worksForNode = network.GetNetworkNodeByUid(bg.networkNode.WorksForUid)
		}
		section1Info = fmt.Sprintf(
			mgaEmitterFormat,
			worksForNode.Agency.Name,
			companyName,
		)
	}

	section3Info = withoutConsultacy
	if bg.policy.ConsultancyValue.Price > 0 {
		section3Info = fmt.Sprintf(
			withConsultacyFormat,
			lib.HumanaizePriceEuro(bg.policy.ConsultancyValue.Price),
		)
	}

	return section1Info, section3Info
}

func (bg *baseGenerator) woptaTable(producerInfo, proponentInfo map[string]string, designation string) {
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
			title: "DATI DELLA PERSONA FISICA CHE ENTRA IN CONTATTO CON IL CONTRAENTE",
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
			body:  designation,
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
			body: "IVASS – Istituto per la Vigilanza sulle Assicurazioni - Via del Quirinale, " +
				"21 - 00187 Roma",
		},
	}
	table = append(table, parseEntries(entries, true)...)

	bg.engine.SetDrawColor(constants.PinkColor)
	bg.engine.DrawTable(table)
	bg.engine.NewLine(2)
}

func (bg *baseGenerator) annex3(producerInfo, proponentInfo map[string]string, designation string) {
	type section struct {
		title string
		body  []string
	}

	bg.engine.WriteText(domain.TableCell{
		Text:      "ALLEGATO 3 - INFORMATIVA SUL DISTRIBUTORE",
		Height:    3,
		Width:     constants.FullPageWidth,
		FontSize:  constants.LargeFontSize,
		FontStyle: constants.BoldFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.CenterAlign,
		Border:    "",
	})
	bg.engine.NewLine(3)

	bg.engine.WriteText(domain.TableCell{
		Text: "Il distributore ha l’obbligo di consegnare/trasmettere al contraente il presente" +
			" documento, prima della sottoscrizione della prima proposta o, qualora non prevista, del primo contratto di " +
			"assicurazione, di metterlo a disposizione del pubblico nei propri locali, anche mediante apparecchiature " +
			"tecnologiche, oppure di pubblicarlo sul proprio sito internet ove utilizzato per la promozione e collocamento " +
			"di prodotti assicurativi, dando avviso della pubblicazione nei propri locali. In occasione di rinnovo o " +
			"stipula di un nuovo contratto o di qualsiasi operazione avente ad oggetto un prodotto di investimento " +
			"assicurativo il distributore consegna o trasmette le informazioni di cui all’Allegato 3 solo in caso di " +
			"successive modifiche di rilievo delle stesse.",
		Height:    3,
		Width:     constants.FullPageWidth,
		FontSize:  constants.RegularFontSize,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "",
	})
	bg.engine.NewLine(3)

	bg.engine.WriteText(domain.TableCell{
		Text: "SEZIONE I - Informazioni generali sull’intermediario che entra in contatto con " +
			"il contraente",
		Height:    3,
		Width:     constants.FullPageWidth,
		FontSize:  constants.LargeFontSize,
		FontStyle: constants.BoldFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "",
	})
	bg.engine.NewLine(3)

	bg.woptaTable(producerInfo, proponentInfo, designation)

	bg.engine.WriteText(domain.TableCell{
		Text: "Gli estremi identificativi e di iscrizione dell’Intermediario e dei soggetti che " +
			"operano per lo stesso possono essere verificati consultando il Registro Unico degli Intermediari assicurativi " +
			"e riassicurativi sul sito internet dell’IVASS (www.ivass.it)",
		Height:    3,
		Width:     constants.FullPageWidth,
		FontSize:  constants.RegularFontSize,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "",
	})

	sections := []section{
		{
			title: "SEZIONE II - Informazioni sull’attività svolta dall’intermediario assicurativo",
			body: []string{
				proponentInfo["name"] + " comunica di aver messo a disposizione nei propri " +
					"locali l’elenco degli obblighi di comportamento cui adempie, come indicati nell’allegato 4-ter del Regolamento" +
					" IVASS n. 40/2018.",
				"Si comunica che nel caso di offerta fuori sede o nel caso in cui la fase " +
					"precontrattuale si svolga mediante tecniche di comunicazione a distanza il contraente riceve l’elenco " +
					"degli obblighi.",
			},
		},
		{
			title: "SEZIONE III - Informazioni relative a potenziali situazioni di conflitto " +
				"d’interessi",
			body: []string{
				proponentInfo["name"] + " ed i soggetti che operano per la stessa non sono " +
					"detentori di una partecipazione, diretta o indiretta, pari o superiore al 10% del capitale sociale o dei " +
					"diritti di voto di alcuna Impresa di assicurazione.",
				"Le Imprese di assicurazione o Imprese controllanti un’Impresa di assicurazione " +
					"non sono detentrici di una partecipazione, diretta o indiretta, pari o superiore al 10% del capitale sociale " +
					"o dei diritti di voto dell’Intermediario.",
			},
		},
		{
			title: "SEZIONE IV - Informazioni sugli strumenti di tutela del contraente",
			body: []string{
				"L’attività di distribuzione è garantita da un contratto di assicurazione della " +
					"responsabilità civile che copre i danni arrecati ai contraenti da negligenze ed errori professionali " +
					"dell’intermediario o da negligenze, errori professionali ed infedeltà dei dipendenti, dei collaboratori o " +
					"delle persone del cui operato l’intermediario deve rispondere a norma di legge.",
				"Il contraente ha la facoltà, ferma restando la possibilità di rivolgersi " +
					"all’Autorità Giudiziaria, di inoltrare reclamo per iscritto all’intermediario, via posta all’indirizzo di " +
					"sede legale o a mezzo mail alla PEC sopra indicati, oppure all’Impresa secondo le modalità e presso i " +
					"recapiti indicati nel DIP aggiuntivo nella relativa sezione, nonché la possibilità, qualora non dovesse " +
					"ritenersi soddisfatto dall’esito del reclamo o in caso di assenza di riscontro da parte dell’intermediario " +
					"o dell’impresa entro il termine di legge, di rivolgersi all’IVASS secondo quanto indicato nei DIP aggiuntivi.",
				"Il contraente ha la facoltà di avvalersi di altri eventuali sistemi alternativi " +
					"di risoluzione delle controversie previsti dalla normativa vigente nonché quelli indicati nei" +
					" DIP aggiuntivi.",
			},
		},
	}

	for _, s := range sections {
		bg.engine.NewLine(3)
		bg.engine.WriteText(domain.TableCell{
			Text:      s.title,
			Height:    3,
			Width:     constants.FullPageWidth,
			FontSize:  constants.LargeFontSize,
			FontStyle: constants.BoldFontStyle,
			FontColor: constants.BlackColor,
			Fill:      false,
			FillColor: domain.Color{},
			Align:     constants.LeftAlign,
			Border:    "",
		})
		bg.engine.NewLine(1)
		for _, b := range s.body {
			bg.engine.WriteText(domain.TableCell{
				Text:      b,
				Height:    3,
				Width:     constants.FullPageWidth,
				FontSize:  constants.RegularFontSize,
				FontStyle: constants.RegularFontStyle,
				FontColor: constants.BlackColor,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "",
			})
			bg.engine.NewLine(1)
		}
	}
}

func (bg *baseGenerator) annex4(producerInfo, proponentInfo map[string]string, designation, annex4Section1Info, annex4Section3Info string) {
	type section struct {
		title string
		body  []string
	}

	bg.engine.WriteText(domain.TableCell{
		Text:      "ALLEGATO 4 - INFORMAZIONI SULLA DISTRIBUZIONE\nDEL PRODOTTO ASSICURATIVO NON IBIP",
		Height:    3,
		Width:     constants.FullPageWidth,
		FontSize:  constants.LargeFontSize,
		FontStyle: constants.BoldFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.CenterAlign,
		Border:    "",
	})
	bg.engine.NewLine(3)

	bg.engine.WriteText(domain.TableCell{
		Text: "Il distributore ha l’obbligo di consegnare o trasmettere al contraente, prima " +
			"della sottoscrizione di ciascuna proposta o, qualora non prevista, di ciascun contratto assicurativo, il " +
			"presente documento, che contiene notizie sul modello e l’attività di distribuzione, sulla consulenza fornita " +
			"e sulle remunerazioni percepite.",
		Height:    3,
		Width:     constants.FullPageWidth,
		FontSize:  constants.RegularFontSize,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "",
	})
	bg.engine.NewLine(3)

	bg.engine.WriteText(domain.TableCell{
		Text: "SEZIONE I - Informazioni generali sull’intermediario che entra in contatto con " +
			"il contraente",
		Height:    3,
		Width:     constants.FullPageWidth,
		FontSize:  constants.LargeFontSize,
		FontStyle: constants.BoldFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "",
	})
	bg.engine.NewLine(3)

	bg.woptaTable(producerInfo, proponentInfo, designation)

	bg.engine.WriteText(domain.TableCell{
		Text: "Gli estremi identificativi e di iscrizione dell’Intermediario e dei soggetti che " +
			"operano per lo stesso possono essere verificati consultando il Registro Unico degli Intermediari assicurativi " +
			"e riassicurativi sul sito internet dell’IVASS (www.ivass.it)",
		Height:    3,
		Width:     constants.FullPageWidth,
		FontSize:  constants.RegularFontSize,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "",
	})

	sections := []section{
		{
			title: "SEZIONE I - Informazioni sul modello di distribuzione",
			body:  []string{annex4Section1Info},
		},
		{
			title: "SEZIONE II: Informazioni sull’attività di distribuzione e consulenza",
			body: []string{
				"Nello svolgimento dell’attività di distribuzione, l’intermediario non presta " +
					"attività di consulenza prima della conclusione del contratto né fornisce al contraente una raccomandazione " +
					"personalizzata ai sensi dell’art. 119-ter, comma 3, del decreto legislativo n. 209/2005 " +
					"(Codice delle Assicurazioni Private)",
				"L'attività di distribuzione assicurativa è svolta in assenza di obblighi " +
					"contrattuali che impongano di offrire esclusivamente i contratti di una o più imprese di " +
					"assicurazioni.",
			},
		},
		{
			title: "SEZIONE III - Informazioni relative alle remunerazioni",
			body: []string{
				annex4Section3Info,
				"L’informazione sopra resa riguarda i compensi complessivamente percepiti da tutti " +
					"gli intermediari coinvolti nella distribuzione del prodotto.",
			},
		},
		{
			title: "SEZIONE IV – Informazioni sul pagamento dei premi",
			body: []string{
				"Relativamente a questo contratto i premi pagati dal Contraente " +
					"all’intermediario e le somme destinate ai risarcimenti o ai pagamenti dovuti dalle Imprese di Assicurazione, " +
					"se regolati per il tramite dell’intermediario costituiscono patrimonio autonomo e separato dal patrimonio " +
					"dello stesso.",
				"Indicare le modalità di pagamento ammesse",
				"Sono consentiti, nei confronti dell'intermediario, esclusivamente bonifico e strumenti di " +
					"pagamento elettronico, quali ad esempio, carte di credito e/o carte di debito, incluse le carte " +
					"prepagate.",
			},
		},
	}

	for _, s := range sections {
		bg.engine.NewLine(3)
		bg.engine.WriteText(domain.TableCell{
			Text:      s.title,
			Height:    3,
			Width:     constants.FullPageWidth,
			FontSize:  constants.LargeFontSize,
			FontStyle: constants.BoldFontStyle,
			FontColor: constants.BlackColor,
			Fill:      false,
			FillColor: domain.Color{},
			Align:     constants.LeftAlign,
			Border:    "",
		})
		bg.engine.NewLine(1)
		for _, b := range s.body {
			bg.engine.WriteText(domain.TableCell{
				Text:      b,
				Height:    3,
				Width:     constants.FullPageWidth,
				FontSize:  constants.RegularFontSize,
				FontStyle: constants.RegularFontStyle,
				FontColor: constants.BlackColor,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "",
			})
			bg.engine.NewLine(1)
		}
	}
}

func (bg *baseGenerator) annex4Ter(producerInfo, proponentInfo map[string]string, designation string) {
	type section struct {
		title string
		body  []string
	}

	bg.engine.WriteText(domain.TableCell{
		Text:      "ALLEGATO 4 TER - ELENCO DELLE REGOLE DI COMPORTAMENTO DEL DISTRIBUTORE",
		Height:    3,
		Width:     constants.FullPageWidth,
		FontSize:  constants.LargeFontSize,
		FontStyle: constants.BoldFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.CenterAlign,
		Border:    "",
	})
	bg.engine.NewLine(3)

	bg.engine.WriteText(domain.TableCell{
		Text: "Il distributore ha l’obbligo di mettere a disposizione del pubblico il " +
			"presente documento nei propri locali, anche mediante apparecchiature tecnologiche, oppure pubblicarlo su " +
			"un sito internet ove utilizzato per la promozione e il collocamento di prodotti assicurativi, dando avviso " +
			"della pubblicazione nei propri locali. Nel caso di offerta fuori sede o nel caso in cui la fase " +
			"precontrattuale si svolga mediante tecniche di comunicazione a distanza, il distributore consegna o " +
			"trasmette al contraente il presente documento prima della sottoscrizione della proposta o, qualora non " +
			"prevista, del contratto di assicurazione.",
		Height:    3,
		Width:     constants.FullPageWidth,
		FontSize:  constants.RegularFontSize,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "",
	})
	bg.engine.NewLine(3)

	bg.woptaTable(producerInfo, proponentInfo, designation)

	bg.engine.WriteText(domain.TableCell{
		Text: "Gli estremi identificativi e di iscrizione dell’Intermediario e dei soggetti che " +
			"operano per lo stesso possono essere verificati consultando il Registro Unico degli Intermediari assicurativi " +
			"e riassicurativi sul sito internet dell’IVASS (www.ivass.it)",
		Height:    3,
		Width:     constants.FullPageWidth,
		FontSize:  constants.RegularFontSize,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "",
	})

	sections := []section{
		{
			title: "Sezione I - Regole generali per la distribuzione di prodotti assicurativi",
			body: []string{
				"a. obbligo di consegna al contraente dell’allegato 3 al Regolamento IVASS " +
					"n. 40 del 2 agosto 2018, prima della sottoscrizione della prima proposta o, qualora non prevista, del primo " +
					"contratto di assicurazione, di metterlo a disposizione del pubblico nei locali del distributore, anche " +
					"mediante apparecchiature tecnologiche, e di pubblicarlo sul sito internet, ove esistente",
				"b. obbligo di consegna dell’allegato 4 al Regolamento IVASS n. 40 del 2 agosto " +
					"2018, prima della sottoscrizione di ciascuna proposta di assicurazione o, qualora non prevista, del contratto " +
					"di assicurazione",
				"c. obbligo di consegnare copia della documentazione precontrattuale e " +
					"contrattuale prevista dalle vigenti disposizioni, copia della polizza e di ogni altro atto o documento " +
					"sottoscritto dal contraente",
				"d. obbligo di proporre o raccomandare contratti coerenti con le richieste e le " +
					"esigenze di copertura assicurativa e previdenziale del contraente o dell’assicurato, acquisendo a tal fine, " +
					"ogni utile informazione",
				"e. obbligo di valutare se il contraente rientra nel mercato di riferimento " +
					"identificato per il contratto di assicurazione proposto e non appartiene alle categorie di clienti per i quali " +
					"il prodotto non è compatibile, nonché l’obbligo di adottare opportune disposizioni per ottenere dai produttori" +
					" le informazioni di cui all’articolo 30-decies comma 5 del Codice e per comprendere le caratteristiche e il " +
					"mercato di riferimento individuato per ciascun prodotto",
				"f. obbligo di fornire in forma chiara e comprensibile le informazioni " +
					"oggettive sul prodotto, illustrandone le caratteristiche, la durata, i costi e i limiti della copertura ed " +
					"ogni altro elemento utile a consentire al contraente di prendere una decisione informata",
			},
		},
	}

	for _, s := range sections {
		bg.engine.NewLine(3)
		bg.engine.WriteText(domain.TableCell{
			Text:      s.title,
			Height:    3,
			Width:     constants.FullPageWidth,
			FontSize:  constants.LargeFontSize,
			FontStyle: constants.BoldFontStyle,
			FontColor: constants.BlackColor,
			Fill:      false,
			FillColor: domain.Color{},
			Align:     constants.LeftAlign,
			Border:    "",
		})
		bg.engine.NewLine(1)
		for _, b := range s.body {
			bg.engine.WriteText(domain.TableCell{
				Text:      b,
				Height:    3,
				Width:     constants.FullPageWidth,
				FontSize:  constants.RegularFontSize,
				FontStyle: constants.RegularFontStyle,
				FontColor: constants.BlackColor,
				Fill:      false,
				FillColor: domain.Color{},
				Align:     constants.LeftAlign,
				Border:    "",
			})
			bg.engine.NewLine(1)
		}
	}
}

func (bg *baseGenerator) howYouCanPaySection() {
	bg.engine.WriteText(domain.TableCell{
		Text:      "Come puoi pagare il premio",
		Height:    4.5,
		Width:     constants.FullPageWidth,
		FontSize:  constants.LargeFontSize,
		FontStyle: constants.BoldFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "",
	})
	bg.engine.NewLine(1)
	bg.engine.WriteText(domain.TableCell{
		Text: "I mezzi di pagamento consentiti, nei confronti di Wopta, " +
			"sono esclusivamente bonifico e strumenti di pagamento elettronico, quali ad esempio, " +
			"carte di credito e/o carte di debito, incluse le carte prepagate.",
		Height:    4.5,
		Width:     constants.FullPageWidth,
		FontSize:  constants.RegularFontSize,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "",
	})
}

func (bg *baseGenerator) emitResumeSection() {
	price := bg.policy.PriceGross
	if bg.policy.PaymentSplit == string(models.PaySplitMonthly) {
		price = bg.policy.PriceGrossMonthly
	}

	text := fmt.Sprintf("Polizza emessa a Milano il %s per un importo di euro %s quale prima rata alla firma, "+
		"il cui pagamento a saldo è da effettuarsi con i metodi di pagamento sopra indicati. "+
		"Costituisce quietanza di pagamento la mail di conferma che Wopta invierà al Contraente.",
		bg.policy.StartDate.Format(constants.DayMonthYearFormat), lib.HumanaizePriceEuro(price))

	bg.engine.WriteText(domain.TableCell{
		Text:      "Emissione Polizza e pagamento della prima rata",
		Height:    4.5,
		Width:     constants.FullPageWidth,
		FontSize:  constants.LargeFontSize,
		FontStyle: constants.BoldFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "",
	})
	bg.engine.NewLine(1)
	bg.engine.WriteText(domain.TableCell{
		Text:      text,
		Height:    4.5,
		Width:     constants.FullPageWidth,
		FontSize:  constants.RegularFontSize,
		FontStyle: constants.RegularFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "",
	})
}

func (bg *baseGenerator) companySignature() {
	type logoInfo struct {
		path                string
		x, y, width, height float64
	}

	type companyDetails struct {
		description string
		logo        logoInfo
	}

	companiesMap := map[string]companyDetails{
		models.AxaCompany: {
			description: "AXA France Vie\n(Rappresentanza Generale per l'Italia)",
			logo: logoInfo{
				path:   "signature_axa.png",
				x:      35,
				y:      9,
				width:  30,
				height: 8,
			},
		},
		models.GlobalCompany: {
			description: "Global Assistance",
			logo: logoInfo{
				path:   "signature_global.png",
				x:      25,
				y:      3,
				width:  40,
				height: 12,
			},
		},
		models.SogessurCompany: {
			description: "Sogessur SA\n(Rappresentanza Generale per l'Italia)",
			logo: logoInfo{
				path:   "signature_sogessur.png",
				x:      40,
				y:      9,
				width:  10,
				height: 10,
			},
		},
		models.QBECompany: {
			description: "QBE Europe SA/NV - Rappresentanza Generale per l’Italia",
			logo: logoInfo{
				path:   "signature_qbe.png",
				x:      35,
				y:      9,
				width:  30,
				height: 8,
			},
		},
	}

	logo := companiesMap[bg.policy.Company].logo

	bg.engine.WriteText(domain.TableCell{
		Text:      companiesMap[bg.policy.Company].description,
		Height:    3,
		Width:     70,
		FontSize:  constants.RegularFontSize,
		FontStyle: constants.BoldFontStyle,
		FontColor: constants.BlackColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.CenterAlign,
		Border:    "",
	})
	bg.engine.SetY(bg.engine.GetY() - 6)
	bg.engine.InsertImage(lib.GetAssetPathByEnvV2()+logo.path, logo.x, bg.engine.GetY()+logo.y, logo.width,
		logo.height)
}

func (bg *baseGenerator) checkSurveySpace(survey models.Survey) {
	var answer string
	leftMargin, _, rightMargin, _ := bg.engine.GetMargins()
	pageWidth, pageHeight := bg.engine.GetPageSize()
	availableWidth := pageWidth - leftMargin - rightMargin - 2
	requiredHeight := 5.0
	currentY := bg.engine.GetY()

	surveyTitle := survey.Title
	surveySubtitle := survey.Subtitle

	if surveyTitle != "" {
		bg.engine.SetFontStyle(constants.BoldFontStyle)
		bg.engine.SetFontSize(constants.LargeFontSize)
		lines := bg.engine.SplitText(surveyTitle, availableWidth)
		requiredHeight += 3.5 * float64(len(lines))
	}
	if surveySubtitle != "" {
		bg.engine.SetFontStyle(constants.BoldFontStyle)
		bg.engine.SetFontSize(constants.RegularFontSize)
		lines := bg.engine.SplitText(surveySubtitle, availableWidth)
		requiredHeight += 3.5 * float64(len(lines))
	}

	for _, question := range survey.Questions {
		availableWidth = pageWidth - leftMargin - rightMargin - 2

		questionText := question.Question

		if question.IsBold {
			bg.engine.SetFontStyle(constants.BoldFontStyle)
			bg.engine.SetFontSize(constants.LargeFontSize)
		} else {
			bg.engine.SetFontStyle(constants.RegularFontStyle)
			bg.engine.SetFontSize(constants.RegularFontSize)
		}
		if question.Indent {
			availableWidth -= tabDimension / 2
		}

		if question.HasAnswer {
			answer = "NO"
			if *question.Answer {
				answer = "SI"
			}
		}

		lines := bg.engine.SplitText(questionText+answer, availableWidth)
		requiredHeight += 3 * float64(len(lines))
	}

	if (!bg.isProposal && survey.ContractorSign) || survey.CompanySign {
		requiredHeight += 35
	}

	if (pageHeight-18)-currentY < requiredHeight {
		bg.engine.NewPage()
	}
}

func (bg *baseGenerator) printSurvey(survey models.Survey) error {
	var dotsString string
	leftMargin, _, rightMargin, _ := bg.engine.GetMargins()
	pageWidth, _ := bg.engine.GetPageSize()
	availableWidth := pageWidth - leftMargin - rightMargin - 2

	bg.checkSurveySpace(survey)

	surveyTitle := survey.Title
	surveySubtitle := survey.Subtitle

	bg.engine.SetFontStyle(constants.BoldFontStyle)
	bg.engine.SetFontSize(constants.RegularFontSize)
	if survey.HasAnswer {
		answer := "NO"
		if *survey.Answer {
			answer = "SI"
		}

		answerWidth := bg.engine.GetStringWidth(answer)
		dotWidth := bg.engine.GetStringWidth(".")

		var surveyWidth, paddingWidth float64
		var lines []string
		if surveyTitle != "" {
			lines = bg.engine.SplitText(surveyTitle+answer, availableWidth)
		} else if surveySubtitle != "" {
			lines = bg.engine.SplitText(surveySubtitle+answer, availableWidth)
		}

		surveyWidth = bg.engine.GetStringWidth(lines[len(lines)-1])
		paddingWidth = availableWidth - surveyWidth - answerWidth

		dotsString = strings.Repeat(".", int(math.Max((paddingWidth/dotWidth)-2, 0))) + answer
	}
	if surveyTitle != "" {
		bg.engine.WriteText(domain.TableCell{
			Text:      surveyTitle + dotsString,
			Height:    4,
			Width:     constants.FullPageWidth,
			FontSize:  constants.LargeFontSize,
			FontStyle: constants.BoldFontStyle,
			FontColor: constants.PinkColor,
			Fill:      false,
			FillColor: domain.Color{},
			Align:     constants.LeftAlign,
			Border:    "",
		})
	}
	if surveySubtitle != "" {
		bg.engine.WriteText(domain.TableCell{
			Text:      surveySubtitle + dotsString,
			Height:    3.5,
			Width:     availableWidth,
			FontSize:  constants.RegularFontSize,
			FontStyle: constants.BoldFontStyle,
			FontColor: constants.BlackColor,
			Fill:      false,
			FillColor: domain.Color{},
			Align:     constants.LeftAlign,
			Border:    "",
		})
	}

	for _, question := range survey.Questions {
		dotsString = ""
		availableWidth = pageWidth - leftMargin - rightMargin - 2
		fontStyle := constants.RegularFontStyle
		fontSize := constants.RegularFontSize

		if question.IsBold {
			fontStyle = constants.BoldFontStyle
		}

		if question.Indent {
			bg.engine.SetX(tabDimension)
			availableWidth -= tabDimension / 2
		}

		if question.HasAnswer {
			var questionWidth, paddingWidth float64
			answer := "NO"
			if *question.Answer {
				answer = "SI"
			}

			answerWidth := bg.engine.GetStringWidth(answer)
			dotWidth := bg.engine.GetStringWidth(".")

			lines := bg.engine.SplitText(question.Question+answer, availableWidth)

			questionWidth = bg.engine.GetStringWidth(lines[len(lines)-1])
			paddingWidth = availableWidth - questionWidth - answerWidth

			dotsString = strings.Repeat(".", int(math.Max((paddingWidth/dotWidth)-2, 0))) + answer
		}
		bg.engine.WriteText(domain.TableCell{
			Text:      question.Question + dotsString,
			Height:    3.5,
			Width:     availableWidth,
			FontSize:  fontSize,
			FontStyle: fontStyle,
			FontColor: constants.BlackColor,
			Fill:      false,
			FillColor: domain.Color{},
			Align:     constants.LeftAlign,
			Border:    "",
		})
	}
	bg.engine.NewLine(3)

	if survey.CompanySign {
		bg.companySignature()
		if bg.isProposal {
			bg.engine.NewLine(20)
		}
	}
	if !bg.isProposal && survey.ContractorSign {
		bg.signatureForm()
		bg.engine.NewLine(10)
	}
	return nil
}

func (bg *baseGenerator) checkStatementSpace(statement models.Statement) {
	leftMargin, _, rightMargin, _ := bg.engine.GetMargins()
	pageWidth, pageHeight := bg.engine.GetPageSize()
	availableWidth := pageWidth - leftMargin - rightMargin - 2
	requiredHeight := 5.0
	currentY := bg.engine.GetY()

	title := statement.Title
	subtitle := statement.Subtitle

	if title != "" {
		bg.engine.SetFontStyle(constants.BoldFontStyle)
		bg.engine.SetFontSize(constants.LargeFontSize)
		lines := bg.engine.SplitText(title, availableWidth)
		requiredHeight += 3.5 * float64(len(lines))
	}
	if subtitle != "" {
		bg.engine.SetFontStyle(constants.RegularFontStyle)
		bg.engine.SetFontSize(constants.RegularFontSize)
		lines := bg.engine.SplitText(subtitle, availableWidth)
		requiredHeight += 3.5 * float64(len(lines))
	}
	for _, question := range statement.Questions {
		availableWidth = pageWidth - leftMargin - rightMargin - 2

		text := question.Question

		if question.IsBold {
			bg.engine.SetFontSize(constants.RegularFontSize)
		} else {
			bg.engine.SetFontStyle(constants.RegularFontStyle)
		}
		bg.engine.SetFontSize(constants.RegularFontSize)

		if question.Indent {
			availableWidth -= tabDimension / 2
		}

		answer := ""
		if question.HasAnswer {
			answer = constants.No
			if *question.Answer {
				answer = constants.Yes
			}
		}

		lines := bg.engine.SplitText(text+answer, availableWidth)
		requiredHeight += 3 * float64(len(lines))
	}

	if (!bg.isProposal && statement.ContractorSign) || statement.CompanySign {
		requiredHeight += 35
	}

	if (pageHeight-18)-currentY < requiredHeight {
		bg.engine.NewPage()
	}
}

func (bg *baseGenerator) printStatement(statement models.Statement) {
	bg.checkStatementSpace(statement)

	title := statement.Title
	subtitle := statement.Subtitle

	if title != "" {
		bg.engine.WriteText(domain.TableCell{
			Text:      title,
			Height:    4,
			Width:     constants.FullPageWidth,
			FontSize:  constants.LargeFontSize,
			FontStyle: constants.BoldFontStyle,
			FontColor: constants.PinkColor,
			Fill:      false,
			FillColor: domain.Color{},
			Align:     constants.LeftAlign,
			Border:    "",
		})
	}
	if subtitle != "" {
		bg.engine.WriteText(domain.TableCell{
			Text:      subtitle,
			Height:    3.5,
			Width:     constants.FullPageWidth,
			FontSize:  constants.RegularFontSize,
			FontStyle: constants.RegularFontStyle,
			FontColor: constants.BlackColor,
			Fill:      false,
			FillColor: domain.Color{},
			Align:     constants.LeftAlign,
			Border:    "",
		})
	}
	for _, question := range statement.Questions {
		text := question.Question
		fontStyle := constants.RegularFontStyle
		fontSize := constants.RegularFontSize
		if question.IsBold {
			fontStyle = constants.BoldFontStyle
		}
		if question.Indent {
			bg.engine.SetX(tabDimension)
		}
		bg.engine.WriteText(domain.TableCell{
			Text:      text,
			Height:    3.5,
			Width:     constants.FullPageWidth,
			FontSize:  fontSize,
			FontStyle: fontStyle,
			FontColor: constants.BlackColor,
			Fill:      false,
			FillColor: domain.Color{},
			Align:     constants.LeftAlign,
			Border:    "",
		})
	}
	bg.engine.NewLine(3)

	if statement.CompanySign {
		bg.companySignature()
		if bg.isProposal {
			bg.engine.NewLine(20)
		}
	}
	if !bg.isProposal && statement.ContractorSign {
		bg.signatureForm()
		bg.engine.NewLine(10)
	}
}
