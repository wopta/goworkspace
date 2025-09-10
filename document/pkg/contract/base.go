package contract

import (
	"fmt"
	"time"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"gitlab.dev.wopta.it/goworkspace/document/internal/constants"
	"gitlab.dev.wopta.it/goworkspace/document/internal/domain"
	"gitlab.dev.wopta.it/goworkspace/document/internal/engine"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
)

const (
	tabDimension = 15
)

type baseGenerator struct {
	engine       *engine.Fpdf
	isProposal   bool
	now          time.Time
	signatureID  uint32
	networkNode  *models.NetworkNode
	worksForNode *models.NetworkNode
	policy       *models.Policy
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
	bg.emptyHeader()
	bg.engine.NewPage()
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
			"conferito all’Intermediario, nonché ove il trattamento risulti necessario per l’adempimento di un preciso " +
			"obbligo di legge posto in capo al Titolare, o al fine di accertare, esercitare o difendere un diritto in " +
			"sede giudiziaria. Il conferimento dei dati personali per queste finalità è facoltativo, ma l'eventuale " +
			"mancato conferimento comporterebbe l'impossibilità per l’Intermediario di eseguire le proprie obbligazioni " +
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
	bg.emptyHeader()
	bg.mup(false, bg.policy.Company, bg.policy.ConsultancyValue.Price, bg.policy.Channel)
}

func (bg *baseGenerator) designationInfo() string {
	var (
		designation                           string
		agentContact                          string
		mgaProponentDirectDesignationFormat   = "%s %s"
		mgaRuiInfo                            = "Wopta Assicurazioni Srl, Società iscritta alla Sezione A del RUI con numero A000701923 in data 14/02/2022"
		designationDirectManager              = "Responsabile dell’attività di intermediazione assicurativa di"
		mgaProponentIndirectDesignationFormat = "%s di %s con sede legale in %s, iscritta in sezione E del RUI con numero %s in data %s, che opera per conto di %s"
		mgaEmitterDesignationFormat           = "%s dell’Intermediario di %s iscritta alla sezione %s del RUI con numero %s in data %s"
	)

	if bg.networkNode == nil || bg.networkNode.Type == models.PartnershipNetworkNodeType {
		designation = fmt.Sprintf(mgaProponentDirectDesignationFormat, designationDirectManager, mgaRuiInfo)
	} else if bg.networkNode.IsMgaProponent {
		if bg.networkNode.WorksForUid == models.WorksForMgaUid {
			designation = fmt.Sprintf(mgaProponentDirectDesignationFormat, bg.networkNode.Designation, mgaRuiInfo)
		} else {
			worksForNode := bg.networkNode
			if bg.networkNode.WorksForUid != "" {
				worksForNode = bg.worksForNode
			}
			designation = fmt.Sprintf(
				mgaProponentIndirectDesignationFormat,
				bg.networkNode.Designation,
				worksForNode.GetName(),
				worksForNode.GetAddress(),
				worksForNode.GetRuiCode(),
				worksForNode.GetRuiRegistration().Format(constants.DayMonthYearFormat),
				mgaRuiInfo,
			)
		}
	} else {
		worksForNode := bg.networkNode
		if bg.networkNode.WorksForUid != "" {
			worksForNode = bg.worksForNode
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

	if bg.policy.Channel == lib.NetworkChannel {
		agentContact = fmt.Sprintf(". Contatti Intermediario - mail: %s", bg.networkNode.Mail)
	}

	designation += agentContact

	log.Printf("designation info: %+v", designation)

	return designation
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
		models.NetInsuranceCompany: {
			description: "Net company",
			logo: logoInfo{
				path:   "signature_axa.png",
				x:      35,
				y:      9,
				width:  30,
				height: 8,
			},
		},
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
		if question.Answer != nil && question.HasAnswer {
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
func (bg *baseGenerator) whoWeAre() {
	bg.engine.WriteText(domain.TableCell{
		Text:      "Chi siamo",
		Height:    constants.CellHeight,
		FontStyle: constants.BoldFontStyle,
		FontColor: constants.PinkColor,
		FontSize:  constants.LargeFontSize,
	})

	bg.engine.NewLine(2)

	bg.engine.RawWriteText(domain.TableCell{
		Text:      "Wopta Assicurazioni S.r.l.",
		Height:    constants.CellHeight,
		FontStyle: constants.BoldFontStyle,
		FontColor: constants.BlackColor,
		FontSize:  constants.RegularFontSize,
	})
	bg.engine.RawWriteText(domain.TableCell{
		Text:      " - intermediario assicurativo, soggetto al controllo dell’IVASS ed iscritto dal 14.02.2022 al Registro Unico degli Intermediari, in Sezione A nr. A000701923, avente sede legale in Galleria del Corso, 1 – 20122 Milano (MI). Capitale sociale Euro 120.000 - Codice Fiscale, Reg. Imprese e Partita IVA: 12072020964 - Iscritta al Registro delle imprese di Milano – REA MI 2638708",
		Height:    constants.CellHeight,
		FontColor: constants.BlackColor,
		FontSize:  constants.RegularFontSize,
	})
}

func (bg *baseGenerator) addElectronicSignPolicy() {
	bg.emptyHeader()
	bg.engine.NewPage()
	bg.engine.WriteText(bg.engine.GetTableCell("Servizio di Firma Elettronica Avanzata (“FEA”)", constants.BoldFontStyle, constants.CenterAlign))
	bg.engine.WriteText(bg.engine.GetTableCell("in modalità One Time Password (“OTP”) erogato da Wopta Assicurazioni srl", constants.BoldFontStyle, constants.CenterAlign))
	bg.engine.NewLine(2)

	bg.engine.WriteText(bg.engine.GetTableCell("Wopta Assicurazioni S.r.l. (sede legale in Galleria del Corso, 1 – 20122 Milano (MI), partita iva 12072020964) (di seguito Wopta), ha scelto di avvalersi di processi digitali basati su documenti informatici per la sottoscrizione dei prodotti assicurativi e la documentazione operativa, dando un sostanzioso contributo alla tutela ambientale, evitando la stampa di una notevole mole di documenti."))
	bg.engine.WriteText(bg.engine.GetTableCell("L’utilizzo di documenti informatici è consentito e incentivato dalla normativa (Art. 62 Reg. 40 - Utilizzo della firma elettronica avanzata, della firma elettronica qualificata e della firma digitale- 1. I distributori favoriscono l’utilizzo da parte dei contraenti della tecnologia di firma elettronica avanzata, di firma elettronica qualificata e di firma digitale per la sottoscrizione della documentazione relativa al contratto di assicurazione. 2. La polizza può essere formata come documento informatico sottoscritto con firma elettronica avanzata, con firma elettronica qualificata o con firma digitale, nel rispetto delle disposizioni normative vigenti in materia)."))
	bg.engine.NewLine(2)
	bg.engine.WriteText(bg.engine.GetTableCell("Condizioni Generali Servizio", constants.BoldFontStyle, constants.CenterAlign))
	bg.engine.NewLine(2)
	bg.engine.WriteTexts(bg.engine.GetTableCell("Il Servizio consente l’utilizzo della Firma Elettronica Autorizzata (da ora in poi FEA) solo dopo l’identificazione del Firmatario e l’esplicita autorizzazione all’utilizzo dello stesso. Il Firmatario in fase di identificazione dovrà fornire necessariamente un numero di telefono mobile e indirizzo mail validi, questi dovranno essere ad uso esclusivo e sotto il pieno controllo del Firmatario. Le comunicazioni del processo di firma avvengono tramite invio di Email all’indirizzo comunicato dal Firmatario. Il documento da firmare verrà inviato al Firmatario che potrà visionarlo. La FEA verrà apposta, sul documento ricevuto dal Firmatario accettando i punti firma indicati e confermando di aver letto il documento, con l’inserimento della One Time Password (OTP) quando richiesto. Si rimanda alla Sezione "), bg.engine.GetTableCell("\"Scheda Tecnica servizio\"", constants.BoldFontStyle), bg.engine.GetTableCell("per ulteriori informazioni. Il Firmatario accettando la FEA, dichiara che sia il numero telefonico mobile, l’indirizzo email comunicati e gli strumenti utilizzati per il processo di FEA sono sotto il proprio esclusivo controllo."))
	bg.engine.NewLine(2)
	bg.engine.WriteText(bg.engine.GetTableCell("L’utilizzo e la validità della FEA è sempre subordinata alla corretta conclusione del processo di identificazione del Firmatario."))
	bg.engine.NewLine(2)
	bg.engine.WriteText(bg.engine.GetTableCell("-	Attivazione del servizio", constants.BoldFontStyle))
	bg.engine.NewLine(2)
	bg.engine.WriteText(bg.engine.GetTableCell("Il nuovo Firmatario attiva e aderisce al servizio di FEA nel processo di identificazione, con contestuale accettazione delle presenti Condizioni Generali come indicato nel percorso di firma. Nel caso il Firmatario avesse già utilizzato il servizio di FEA in precedenza e quindi identificato, per ogni firma successiva alla prima verrà invitato a verificare l’esattezza dei dati personali archiviati e gli estremi dell’accettazione della FEA. In caso di variazioni dei dati, il Firmatario dovrà aggiornarli e accettare nuovamente la FEA. Il Nuovo Firmatario attiva la FEA mediante l’inserimento del codice OTP ricevuto via SMS/mail e attiva il processo di firma in conformità a quanto previsto dalla normativa vigente."))
	bg.engine.NewLine(2)
	bg.engine.WriteText(bg.engine.GetTableCell("-	Revoca del servizio", constants.BoldFontStyle))
	bg.engine.NewLine(2)
	bg.engine.WriteText(bg.engine.GetTableCell("Viene richiesto di rinnovare la scelta ad ogni processo di firma."))
	bg.engine.NewLine(2)
	bg.engine.WriteText(bg.engine.GetTableCell("-	Apposizione della firma elettronica avanzata", constants.BoldFontStyle))
	bg.engine.NewLine(2)
	bg.engine.WriteText(bg.engine.GetTableCell("L’ operazione di accettazione del presente documento “Servizio di Firma Elettronica Avanzata (“FEA”) in modalità One Time Password (“OTP”) erogato da Wopta Assicurazioni srl” viene effettuato ad ogni firma. L’apposizione della FEA sul documento informatico avviene con la validazione di tutti i punti firma richiesti e la dichiarazione di aver letto tutto il documento da firmare. A seguire, quando richiesto, dovrà essere inserito il codice OTP inviato al Firmatario attraverso il canale (SMS/Email) scelto dal Firmatario. Questa operazione conclude il proVcesso di firma crittografica del documento, necessario a rendere il documento elettronico integro, non modificabile, leggibile e autentico. Il processo di FEA può ritenersi concluso solo quando il Firmatario riceve l’email con il messaggio di completamento dell’operazione e contenente il documento firmato."))
	bg.engine.WriteText(bg.engine.GetTableCell("Successivamente, il Firmatario potrà verificare in qualsiasi momento i dettagli del documento firmato su una pagina specifica nella homepage “Area Riservata”"))
	bg.engine.NewLine(2)
	bg.engine.WriteText(bg.engine.GetTableCell("-	Limitazioni d’uso", constants.BoldFontStyle))
	bg.engine.NewLine(2)
	bg.engine.WriteText(bg.engine.GetTableCell("Il Firmatario potrà utilizzare la firma elettronica avanzata esclusivamente per sottoscrivere i documenti per i quali ha accettato di aderire."))
	bg.engine.NewLine(2)
	bg.engine.WriteText(bg.engine.GetTableCell("-	Recupero e verifica del documento informatico sottoscritto", constants.BoldFontStyle))
	bg.engine.NewLine(2)
	bg.engine.WriteText(bg.engine.GetTableCell("Il cliente può recuperare il documento informatico sottoscritto attraverso l’apposita funzionalità descritta al punto precedente “Apposizione della firma elettronica avanzata”. Il cliente può verificare l’integrità, la non modifica e l’autenticità del documento e delle firme elettroniche attraverso i passi: Salvataggio sul proprio computer del documento informatico sottoscritto (formato pdf). Inoltre ad ogni documento sottoscritto è associato un Report dettagliato delle operazioni che hanno portato alla firma del documento stesso."))
	bg.engine.NewPage()

	bg.engine.WriteText(bg.engine.GetTableCell("-	Conservazione del documento informatico sottoscritto", constants.BoldFontStyle))
	bg.engine.NewLine(2)
	bg.engine.WriteText(bg.engine.GetTableCell("Wopta provvede all’archiviazione digitale del processo e dei dati del Firmatario con firma elettronica avanzata in ottemperanza alle indicazioni normative. Wopta si occupa della conservazione del documento e dell’archiviazione di tutte le evidenze informatiche necessarie a comprovarne l’integrità, la leggibilità, l’assenza di modifiche dopo l’apposizione delle firme e l’autenticità delle firme apposte, il tutto in conformità alla normativa tempo per tempo vigente."))
}

func (bg *baseGenerator) addOtpSignPolicy() {
	bg.emptyHeader()
	bg.engine.NewPage()
	bg.engine.WriteText(bg.engine.GetTableCell("Scheda Tecnica Servizio di FEA in modalità OTP erogato da Wopta", constants.BoldFontStyle, constants.CenterAlign))
	bg.engine.WriteText(bg.engine.GetTableCell("- DPCM 23 Febbraio 2013 art. 57 comma 1 lett. E) ed f) delle Regole Tecniche -", constants.CenterAlign))
	bg.engine.NewLine(2)

	bg.engine.WriteText(bg.engine.GetTableCell("Per la consueta operatività Wopta si avvale, nel rispetto della normativa, di “documenti elettronici” in sostituzione di quelli cartacei."))
	bg.engine.WriteText(bg.engine.GetTableCell("L'utilizzo dei documenti informatici è possibile grazie ad una tecnologia che Le permette di indicare, laddove richiesto, la Sua volontà di apporre la firma utilizzando una procedura gratuita e guidata autorizzata tramite un codice usa e getta."))
	bg.engine.WriteText(bg.engine.GetTableCell("Il servizio di Firma Elettronica Avanzata (di seguito anche “FEA”) è attuato dalla Società nel rispetto delle vigenti disposizioni in materia."))

	bg.engine.NewLine(2)
	bg.engine.WriteText(bg.engine.GetTableCell("-	Il controllo esclusivo del firmatario del sistema di generazione della firma", constants.BoldFontStyle))
	bg.engine.NewLine(2)

	bg.engine.WriteText(bg.engine.GetTableCell("Il Cliente ha il controllo esclusivo del sistema del processo di firma, avendo sempre la possibilità, per ogni singola firma apposta di:"))
	bg.engine.WriteText(bg.engine.GetTableCell("     -     visualizzare il contenuto dei documenti al momento della firma in modo da aver evidenza di quanto sta per sottoscrivere"))
	bg.engine.WriteText(bg.engine.GetTableCell("     -     utilizzare gli appositi controlli per ingrandire/rimpicciolire il testo del documento"))
	bg.engine.WriteText(bg.engine.GetTableCell("     -     identificare in modo intuitivo tutte le parti del documento dove è prevista una firma"))
	bg.engine.WriteText(bg.engine.GetTableCell("     -     confermare la firma apposta"))
	bg.engine.WriteText(bg.engine.GetTableCell("     -     annullare l’operazione di firma"))
	bg.engine.WriteText(bg.engine.GetTableCell("Per garantire il massimo livello di sicurezza possibile relativamente a questa scelta tecnologica, Wopta Assicurazioni S.r.l. ha adottato le migliori soluzioni certificate sul mercato dotate di sofisticate tecnologie atte ad impedire manomissioni informatiche."))
	bg.engine.WriteText(bg.engine.GetTableCell("La piattaforma utilizzata da Wopta e fornita da Namirial S.p.A. è integrata per la digitalizzazione delle transazioni di business e dei workflow documentali eSignAnyWhere (eSAW). Il sistema eSAW implementa funzionalità di apposizione delle firme digitali. La modalità standard è quella interattiva in cui i firmatari accedono ad un Web Signature Tray (Viewer) all’interno del quale è possibile visualizzare e firmare i documenti proposti seguendo un wizard particolarmente intuitivo. La modalità interattiva prevede un flusso secondo il quale un’applicazione oppure un sito internet richiedono la firma digitale di un documento attraverso la creazione di una pratica (o transazione) di firma."))
	bg.engine.WriteText(bg.engine.GetTableCell("I documenti da firmare sono protetti da un meccanismo di autenticazione per assicurarsi che solamente il firmatario designato possa accedervi. Il meccanismo di autenticazione utilizza un OTP (One Time Password) inviato sul numero di cellulare del firmatario (dispositivo personale). Il sistema eSAW crea la pratica di firma e genera una URL univoca che la identifica e può essere direttamente utilizzata nel browser del firmatario designato. L’URL viene inviato all’email indicata dal firmatario in sede di censimento dati anagrafici. Accedendo con un browser alla URL della transazione di firma, viene mostrato il Viewer, ovvero un’interfaccia che presenta:"))
	bg.engine.WriteText(bg.engine.GetTableCell("     -     Accettazione di utilizzo del servizio di FEA;"))
	bg.engine.WriteText(bg.engine.GetTableCell("     -     Un’anteprima del documento da firmare;"))
	bg.engine.WriteText(bg.engine.GetTableCell("     -     I controlli di firma;"))
	bg.engine.WriteText(bg.engine.GetTableCell("     -     Il wizard di firma;"))
	bg.engine.WriteText(bg.engine.GetTableCell("     -     Menu con Guida e azioni accessorie per l’utente;"))

	bg.engine.NewLine(2)
	bg.engine.WriteText(bg.engine.GetTableCell("-	La connessione univoca della firma al firmatario", constants.BoldFontStyle))
	bg.engine.NewLine(2)
	bg.engine.WriteText(bg.engine.GetTableCell("I dati di ogni firma raccolta vengono poi cifrati, utilizzando uno specifico certificato digitale, ed inclusi all'interno del documento sottoscritto. Il Cliente ha sempre la possibilità di verificare che il documento informatico sottoscritto non abbia subito modifiche dopo l’apposizione della firma; presso AgID all’ URL http://www.agid.gov.it/identita-digitali/firme-elettroniche/software-verifica.S"))
}
