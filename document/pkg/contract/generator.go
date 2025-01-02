package contract

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/wopta/goworkspace/document/internal/constants"
	"github.com/wopta/goworkspace/document/internal/domain"
	"github.com/wopta/goworkspace/document/internal/engine"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
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
			Width:     190,
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
		lib.CheckError(err)

		if consent.Answer {
			consentText = "X"
			notConsentText = " "
		}
	}

	bg.engine.WriteText(domain.TableCell{
		Text:      "Consenso per finalità commerciali.",
		Height:    3,
		Width:     190,
		FontSize:  constants.LargeFontSize,
		FontStyle: constants.BoldFontStyle,
		FontColor: constants.PinkColor,
		Fill:      false,
		FillColor: domain.Color{},
		Align:     constants.LeftAlign,
		Border:    "",
	})
	bg.engine.WriteText(domain.TableCell{
		Text:      "Il sottoscritto, letta e compresa l’informativa sul trattamento dei dati personali",
		Height:    3,
		Width:     190,
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
		Width:     190,
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
	if !bg.isProposal {
		bg.signatureForm()
	}
}

func (bg *baseGenerator) signatureForm() {
	text := fmt.Sprintf("\"[[!sigField\"%d\":signer1:signature(sigType=\\\"Click2Sign\\\"):label"+
		"(\\\"firma qui\\\"):size(width=150,height=60)]]\"", bg.signatureID)

	bg.engine.SetY(bg.engine.GetY() - 3)
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

func (bg *baseGenerator) annexSections() {
	if bg.networkNode != nil && !bg.networkNode.HasAnnex && bg.networkNode.Type != models.PartnershipNetworkNodeType {
		return
	}

	if bg.networkNode == nil || bg.networkNode.Type == models.PartnershipNetworkNodeType || bg.networkNode.
		IsMgaProponent {

		bg.woptaHeader()

		bg.woptaFooter()
	} else {
		bg.emptyHeader()

		bg.emptyFooter()
	}

	producerInfo := bg.productInfo()
	proponetInfo := bg.proponentInfo()
	designationInfo := bg.designationInfo()
	annex4Section1Info := bg.annex4Section1Info()

	log.Println(producerInfo, proponetInfo, designationInfo, annex4Section1Info)

	bg.engine.NewPage()

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

		proponentInfo["address"] = "====="
		proponentInfo["phone"] = "====="
		proponentInfo["email"] = "====="
		proponentInfo["pec"] = "====="
		proponentInfo["website"] = "====="

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

func (bg *baseGenerator) annex4Section1Info() string {
	var (
		section1Info       string
		mgaProponentFormat = "Secondo quanto indicato nel modulo di proposta/polizza e documentazione " +
			"precontrattuale ricevuta, la distribuzione  relativamente a questa proposta/contratto è svolta per " +
			"conto della seguente impresa di assicurazione: %s"
		mgaEmitterFormat = "Il contratto viene intermediato da %s, in qualità di soggetto proponente, che opera in " +
			"virtù della collaborazione con Wopta Assicurazioni Srl (intermediario emittente dell'Impresa di " +
			"Assicurazione %s, iscritto al RUI sezione A nr A000701923 dal 14.02.2022, ai sensi dell’articolo 22, " +
			"comma 10, del decreto legge 18 ottobre 2012, n. 179, convertito nella legge 17 dicembre 2012, n. 221"
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

	log.Printf("annex 4 section 1 info: %s", section1Info)

	return section1Info
}
