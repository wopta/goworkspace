package document

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
	lib "github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func (skin Skin) GlobalContract(m pdf.Maroto, data models.Policy) {
	layout2 := "2006-01-02"
	var (
		logo     string
		name     string
		nameSign string
	)
	if data.Name == "persona" {
		logo = "/persona.png"
		name = "Persona"
		m = skin.GetHeader(m, data, logo, name)
		m = skin.GetFooter(m, "/logo_global.png", "Wopta per te. Persona è un prodotto assicurativo di Global Assistance Compagnia di assicurazioni e riassicurazioni S.p.A, distribuito da Wopta Assicurazioni S.r.l")
		m = skin.Space(m, 5.0)
		m = skin.GetPersona(data, m)
		m = skin.CoveragesPersonTable(m, data)
		m.AddPage()
	}

	if data.Name == "pmi" {
		logo = "/pmi.png"
		name = "Artigiani & Imprese"
		m = skin.GetHeader(m, data, logo, name)
		m = skin.GetFooter(m, "/logo_global.png", "Wopta per te. Artigiani & Imprese è un prodotto assicurativo di Global Assistance Compagnia di assicurazioni e riassicurazioni S.p.A, distribuito da Wopta Assicurazioni S.r.l")
		m = skin.Space(m, 5.0)
		m = skin.GetPmi(data, m)
		m = skin.Space(m, 5.0)
		skin.GlobalEnterpriseTable(m, data)
		m.AddPage()
		skin.GlobalBuildingTable(m, data)

	}
	for _, as := range data.Assets {
		if as.Enterprise != nil {

			nameSign = as.Enterprise.Name
			break
		} else {
			name = data.Contractor.Name + " " + data.Contractor.Surname
		}
	}

	skin.checkPage(m)
	for i, A := range *data.Statements {

		skin.Stantement(m, A.Title, A)
		if i == 2 {
			skin.SignDouleLine(m, nameSign, "Global Assistance", strconv.Itoa(i+1), true)
		} else {
			skin.Sign(m, nameSign, "Assicurato ", strconv.Itoa(i+1), true)
		}
		skin.checkPage(m)
		m = skin.Space(m, 5.0)

	}

	skin.checkPage(m)
	skin.PriceTable(m, data)
	m = skin.Space(m, 5.0)
	m.Row(skin.RowHeight, func() {
		m.Col(12, func() {
			m.Text("In caso di sostituzione,il premio alla firma è al netto dell’ eventuale rimborso dei premi non goduti sulla polizza sostituita e tiene conto dell’ eventuale diversa durata rispetto alle rate successive. ", skin.NormaltextLeft)

		})

	})

	m = skin.Space(m, 5.0)
	title := "Come puoi pagare il premio "
	body := `I mezzi di pagamento consentiti nei confronti di Wopta sono esclusivamente 
 bonifico e strumenti di pagamento elettronico, quali ad esempio, carte di 
credito e/o carte di debito, incluse le carte prepagate.`
	skin.checkPage(m)
	m = skin.Space(m, 5.0)
	m = skin.Title(m, title, body, 10.0)
	title = "Emissione polizza e pagamento della prima rata "
	s := fmt.Sprintf("%.2f", data.PriceGross)
	body = `Polizza emessa a Milano il ` + time.Now().Format(layout2) + ` per un importo di euro ` + s + ` quale prima rata alla firma,
 il cui pagamento a saldo è da effettuarsi con i metodi di pagamento sopra indicati. 
Costituisce quietanza di pagamento la mail di conferma che Wopta invierà al Contraente. `
	skin.checkPage(m)
	m = skin.Title(m, title, body, 18.0)
	m = skin.RowCol1(m, "", consts.Normal)

	skin.checkPage(m)
	aboutUs := []Kv{{
		Key:   "Wopta Assicurazioni S.r.l.",
		Value: "Wopta Assicurazioni S.r.l. - intermediario assicurativo, soggetto al controllo dell’IVASS ed iscritto dal 14.02.2022 al Registro Unico degli Intermediari, in Sezione A nr. A000701923, avente sede legale in Galleria del Corso, 1 – 20122 Milano (MI). Capitale sociale Euro 120.000 - Codice Fiscale, Reg. Imprese e Partita IVA: 12072020964 - Iscritta al Registro delle imprese di Milano – REA MI 2638708",
	}, {Key: "Global Assistance Compagnia di assicurazioni e riassicurazioni S.p.A. a Socio Unico",
		Value: "Global Assistance Compagnia di assicurazioni e riassicurazioni S.p.A. a Socio Unico - Capitale Sociale: Euro 5.000.000 i.v. Codice Fiscale, Partita IVA e Registro Imprese di Milano n. 10086540159 R.E.A. n. 1345012 della C.C.I.A.A. di Milano. Sede e Direzione Generale Piazza Diaz 6 – 20123 Milano – Italia E-mail: global.assistance@globalassistance.it PEC: globalassistancespa@legalmail.it. Società soggetta all’attività di direzione e coordinamento di Ri-Fin S.r.l., iscritta all’Albo dei gruppi assicurativi presso l’Ivass al n. 014. La Società è autorizzata all’esercizio delle Assicurazioni e Riassicurazioni con D.M. del 2/8/93 n. 19619 (G.U. 7/8/93 n. 184) e successive autorizzazioni ed è iscritta all’Albo Imprese presso l’IVASS al n. 1.00111. La Società è soggetta alla vigilanza dell’IVASS; è possibile verificare la veridicità e la regolarità delle autorizzazioni mediante l'accesso al sito www.ivass.it"}}
	m = skin.AboutUs(m, "Chi siamo ", aboutUs)
	var consens string
	consens = "NON ACCONSENTO"
	if (*data.Contractor.Consens)[0].Answer {

		consens = "ACCONSENTO"
	}
	body = `Consenso per finalità commerciali. 
Il sottoscritto, letta e compresa l’informativa sul trattamento dei dati personali ` + consens + ` al trattamento dei propri dati personali da parte di Wopta Assicurazioni per l’invio di comunicazioni e proposte 
commerciali e di marketing, incluso l’invio di newsletter e ricerche di mercato, attraverso strumenti automatizzati 
(sms, mms, e-mail, ecc.) e non (posta cartacea e telefono con operatore)`
	m = skin.Space(m, 5.0)
	m = skin.Title(m, "Consenso per finalità commerciali. ", body, 10.0)
	skin.Sign(m, nameSign, "Assicurato ", "100", true)
	skin.Space(m, 5.0)
	m.RegisterHeader(func() {
		m.Row(15.0, func() {

			m.ColSpace(10)
			m.Col(2, func() {
				_ = m.FileImage(lib.GetAssetPathByEnv("document")+"/logo_global_02.png", props.Rect{
					Left:    1,
					Top:     1,
					Center:  false,
					Percent: 100,
				})
			})
		})

	})

	m.AddPage()
	m = skin.Space(m, 5.0)
	skin.TitleBlack(m, "DICHIARAZIONI E CONSENSI", "Io Sottoscritto, dichiaro di avere perso visione dell’Informativa Privacy ai sensi dell’art. 13 del GDPR (informativa resa all’interno del set documentale contenente anche la Documentazione Informativa Precontrattuale, il Glossario e le Condizioni di Assicurazione) e di averne compreso i contenuti:", 14.0)
	m = skin.Space(m, 5.0)
	m = skin.Sign(m, nameSign, "Assicurato ", "101", true)
	m = skin.Space(m, 5.0)
	m.Row(skin.RowHeight*2, func() {
		m.Col(12, func() {
			m.Text("Qui di seguito esprimo il mio consenso al trattamento dei dati personali particolari per le finalità sopra indicate, in conformità con quanto previsto all’interno dell’informativa:", skin.NormaltextLeft)

		})

	})
	m.Row(skin.RowHeight, func() {
		m.Col(12, func() {
			m.Text("1.	Consenso al trattamento dei miei dati al fine di perfezionamento dell’offerta assicurativa e ", skin.BoldtextLeft)

		})

	})
	m.Row(skin.RowHeight, func() {
		m.Col(12, func() {
			m.Text("	riassicurativa di cui alle lettere b) ed f) della presente informativa: ", skin.BoldtextLeft)

		})

	})
	m = skin.Space(m, 5.0)
	m = skin.Sign(m, nameSign, "Assicurato ", "102", true)
	m.RegisterFooter(func() {
		topv := 10.0
		t := props.Text{
			Top: topv,

			Size:  6,
			Style: consts.Normal,
			Align: consts.Left,
			Color: skin.TextColor}
		t1 := props.Text{
			Top: topv, Size: 6, Color: skin.SecondaryColor}

		m.Row(5, func() {
			m.Col(12, func() {
				m.Text("Global Assistance Compagnia di assicurazioni e riassicurazioni S.p.A. a Socio Unico", t)
				m.Text("Global Assistance Compagnia di assicurazioni e riassicurazioni S.p.A. a Socio Unico - Capitale Sociale: Euro 5.000.000 i.v. Codice Fiscale, Partita IVA e Registro Imprese di Milano n. 10086540159 R.E.A. n. 1345012 della C.C.I.A.A. di Milano. Sede e Direzione Generale Piazza Diaz 6 – 20123 Milano – ItaliaE-mail: global.assistance@globalassistance.it PEC: globalassistancespa@legalmail.it. Società soggetta all’attività di direzione e coordinamento di Ri-Fin S.r.l., iscritta all’Albo dei gruppi assicurativi presso l’IVASS al n. 014. La Società è autorizzata all’esercizio delle Assicurazioni e Riassicurazioni con D.M. del 2/8/93 n. 19619 (G.U. 7/8/93 n. 184) e successive autorizzazioni ed è iscritta all’Albo Imprese presso l’Ivass al n. 1.00111. La Società è soggetta alla vigilanza dell’IVASS; è possibile verificare la veridicità e la regolarità delle autorizzazioni mediante l'accesso al sito www.ivass.it", t1)
			})

		})

	})
}
func (s Skin) GlobalEnterpriseTable(m pdf.Maroto, data models.Policy) {
	textS := s.MagentaBoldtextLeft
	textS.Size = textS.Size - 3
	m.Row(s.RowTitleHeight, func() {
		m.Col(12, func() {
			m.Text("Le coperture assicurative che hai scelto ", s.MagentaBoldtextLeft)

		})
	})

	m.Row(s.RowTitleHeight-1, func() {
		m.Col(12, func() {
			m.Text("(operative se indicata Somma o Massimale e secondo le Opzioni/Estensioni attivate qualora indicato) ", textS)
		})
	})
	var (
		c                          [][]string
		SumInsuredLimitOfIndemnity string
		detail                     string
		title                      string
	)
	//-----------------------------------------------------------------------
	c = lib.Make2D[string](4, 8)

	c[0][4] = "Responsabilità civile terzi"
	c[1][4] = GetSumIndenity(data.Assets, "third-party-liability")

	c[2][0] = "Sono attive le seguenti opzioni / estensioni:"
	c[2][1] = "Danni ai veicoli in consegna e custodia: " + IfString(ExistAsset(data.Assets, "damage-to-goods-in-custody"), "SI", "NO")
	c[2][2] = "Responsabilità civile postuma officine: " + IfString(ExistAsset(data.Assets, "defect-liability-workmanships"), "SI", "NO")
	c[2][3] = "Responsabilità civile postuma 12 mesi: " + IfString(ExistAsset(data.Assets, "defect-liability-12-months"), "SI", "NO")
	c[2][4] = "Responsabilità civile postuma D.M.37/2008: " + IfString(ExistAsset(data.Assets, "defect-liability-dm-37-2008"), "SI", "NO")
	c[2][5] = "Danni da furto: " + IfString(ExistAsset(data.Assets, "property-damage-due-to-theft"), "SI", "NO")
	c[2][6] = "Danni alle cose sulle quali si eseguono i lavori: " + IfString(ExistAsset(data.Assets, "damage-to-goods-course-of-works"), "SI", "NO")
	c[2][7] = "RC impresa edile: " + IfString(ExistAsset(data.Assets, "third-party-liability-construction-company"), "SI", "NO")
	c[3][4] = GetPrice(data.Assets, "third-party-liability",
		"damage-to-goods-in-custody", "defect-liability-workmanships", "defect-liability-12-months", "defect-liability-dm-37-2008",
		"property-damage-due-to-theft", "damage-to-goods-course-of-works", "third-party-liability-construction-company")

	head1 := []string{"Garanzie ", "Massimale ", "Opzioni / Estensioni ", "Premio "}

	s.BackgroundColorRow(m, "Attività", s.SecondaryColor, s.WhiteTextCenter, s.RowTitleHeight)

	s.TableHeader(m, head1, true, 3, s.rowtableHeight+2, consts.Center, 0)
	s.MultiRow(m, c, true, []uint{4, 2, 4, 2}, 40)
	//-----------------------------------------------------------------------
	c = lib.Make2D[string](4, 1)
	c[0][0] = "Responsabilità civile addetti"
	c[1][0] = GetSumIndenity(data.Assets, "employers-liability")
	c[2][0] = "= ="
	c[3][0] = GetPrice(data.Assets, "employers-liability")
	s.MultiRow(m, c, true, []uint{4, 2, 4, 2}, 40)
	//-----------------------------------------------------------------------
	c = lib.Make2D[string](4, 1)
	c[0][0] = "Responsabilità civile prodotti"
	c[1][0] = GetSumIndenity(data.Assets, "product-liability")
	c[2][0] = "= ="
	c[3][0] = GetPrice(data.Assets, "product-liability")
	s.MultiRow(m, c, true, []uint{4, 2, 4, 2}, 40)
	//-----------------------------------------------------------------------
	c = lib.Make2D[string](4, 3)
	legalg := GetGuarante(data.Assets, "legal-defence")
	log.Println(legalg.LegalDefence)
	if legalg.LegalDefence == "basic" {
		SumInsuredLimitOfIndemnity = "€ " + humanize.FormatInteger("#.###,", 10000)
		detail = "Difesa Penale"
		title = "E’ attiva la garanzia:"
	} else if legalg.LegalDefence == "extended" {
		SumInsuredLimitOfIndemnity = "€ " + humanize.FormatInteger("#.###,", 25000)
		detail = "Difesa Penale Difesa Civile Circolazione"
		title = "E’ attiva la garanzia:"
	} else {
		SumInsuredLimitOfIndemnity = "= ="
		detail = "= ="
		title = ""
	}
	c[0][1] = "Tutela legale"
	c[1][1] = SumInsuredLimitOfIndemnity
	c[2][0] = title
	c[2][1] = detail
	c[3][1] = GetPrice(data.Assets, "legal-defence")
	s.MultiRow(m, c, true, []uint{4, 2, 4, 2}, 40)
	//-----------------------------------------------------------------------
	c = lib.Make2D[string](4, 1)

	c[0][0] = "Cyber risk"
	c[1][0] = GetSumIndenity(data.Assets, "cyber")
	c[2][0] = "= ="
	c[3][0] = GetPrice(data.Assets, "cyber")
	s.MultiRow(m, c, true, []uint{4, 2, 4, 2}, 40)
	//-----------------------------------------------------------------------

}
func (s Skin) GlobalBuildingTable(m pdf.Maroto, data models.Policy) {

	var (
		c [][]string
	)
	//-----------------------------------------------------------------------

	c = lib.Make2D[string](4, 12)

	c[0][5] = "Fabbricato"
	c[1][5] = GetSumIndenity(data.Assets, "building")
	c[0][7] = "Contenuto"
	c[1][7] = GetSumIndenity(data.Assets, "content")

	c[2][0] = "Forma di Assicurazione: VALORE INTERO"
	c[2][1] = "Formula di copertura: RISCHI NOMINATI"
	c[2][2] = "Sono attive le seguenti opzioni / estensioni:"
	c[2][3] = "Eventi Atmosferici: " + IfString(ExistAsset(data.Assets, "atmospheric-event"), "fino al 100% Somme Assicurate", "NO")
	c[2][4] = "Eventi Sociopolitici:  " + IfString(ExistAsset(data.Assets, "sociopolitical-event"), "fino al 80% Somme Assicurate", "NO")
	c[2][5] = "Atti di Terrorismo:  " + IfString(ExistAsset(data.Assets, "terrorism"), "fino al 50% Somme Assicurate", "NO")
	c[2][6] = "Terremoto:  " + IfString(ExistAsset(data.Assets, "earthquake"), "fino al 70% Somme Assicurate", "NO")
	c[2][7] = "Alluvione/Inondazioni: " + IfString(ExistAsset(data.Assets, "water-damage"), " fino al 70% Somme Assicurate", "NO")
	c[2][8] = "Danni da acqua:  " + IfString(ExistAsset(data.Assets, "burst-pipe"), "fino a "+GetSumIndenity(data.Assets, "burst-pipe"), "NO")
	c[2][9] = "Fenomeno Elettrico: " + IfString(ExistAsset(data.Assets, "power-surge"), "fino a "+GetSumIndenity(data.Assets, "power-surge"), "NO")
	c[2][10] = "Rotture Lastre: " + IfString(ExistAsset(data.Assets, "glass"), "fino a "+GetSumIndenity(data.Assets, "glass"), "NO")
	c[2][11] = "Guasto Macchine: " + IfString(ExistAsset(data.Assets, "machinery-breakdown"), "fino a "+GetSumIndenity(data.Assets, "machinery-breakdown"), "NO")

	c[3][4] = GetPrice(data.Assets, "building",
		"content", "atmospheric-event", "sociopolitical-event", "terrorism",
		"earthquake", "water-damage", "burst-pipe", "power-surge", "glass", "machinery-breakdown")

	head2 := []string{"Garanzie ", "Somma assicurata ", "Opzioni / Dettagli ", "Premio "}
	s.BackgroundColorRow(m, "Sede", s.SecondaryColor, s.WhiteTextCenter, s.RowTitleHeight)
	s.TableHeader(m, head2, true, 3, s.rowtableHeight+2, consts.Center, 0)

	s.MultiRow(m, c, true, []uint{4, 2, 4, 2}, 40)
	//-----------------------------------------------------------------------
	c = lib.Make2D[string](4, 1)
	c[0][0] = "Rischio Locativo"
	c[1][0] = GetSumIndenity(data.Assets, "lease-holders-interest")
	c[2][0] = "= ="
	c[3][0] = GetPrice(data.Assets, "lease-holders-interest")
	s.MultiRow(m, c, true, []uint{4, 2, 4, 2}, 40)
	//-----------------------------------------------------------------------
	c = lib.Make2D[string](4, 1)
	c[0][0] = "Ricorso Terzi Incendio"
	c[1][0] = GetSumIndenity(data.Assets, "third-party-recourse")
	c[2][0] = "= ="
	c[3][0] = GetPrice(data.Assets, "third-party-recourse")
	s.MultiRow(m, c, true, []uint{4, 2, 4, 2}, 40)
	//-----------------------------------------------------------------------
	c = lib.Make2D[string](4, 1)
	c[0][0] = "Responsabilità Civile Fabbricato"
	c[1][0] = GetSumIndenity(data.Assets, "property-owners-liability")
	c[2][0] = "= ="
	c[3][0] = GetPrice(data.Assets, "property-owners-liability")
	s.MultiRow(m, c, true, []uint{4, 2, 4, 2}, 40)
	//-----------------------------------------------------------------------
	c = lib.Make2D[string](4, 1)
	c[0][0] = "Responsabilità Civile Inquinamento"
	c[1][0] = GetSumIndenity(data.Assets, "environmental-liability")
	c[2][0] = "= ="
	c[3][0] = GetPrice(data.Assets, "environmental-liability")
	s.MultiRow(m, c, true, []uint{4, 2, 4, 2}, 40)
	//-----------------------------------------------------------------------
	c = lib.Make2D[string](4, 3)

	c[0][1] = "Furto, rapina ed estorsione"
	c[1][1] = GetSumIndenity(data.Assets, "theft")
	c[2][0] = "Sono attive le garanzie opzionali: "
	c[2][1] = "Furto valori e preziosi in cassaforte: " + IfString(ExistAsset(data.Assets, "valuables-in-safe-strongrooms"), "fino a "+GetSumIndenity(data.Assets, "valuables-in-safe-strongrooms"), "NO")
	c[2][2] = "Portavalori: " + IfString(ExistAsset(data.Assets, "valuables"), "fino a "+GetSumIndenity(data.Assets, "valuables"), "NO")
	c[3][1] = GetPrice(data.Assets, "theft", "valuables-in-safe-strongrooms", "valuables")
	s.MultiRow(m, c, true, []uint{4, 2, 4, 2}, 40)
	//-----------------------------------------------------------------------
	c = lib.Make2D[string](4, 4)

	c[0][1] = "Apparecchiature Elettroniche"
	c[1][1] = GetSumIndenity(data.Assets, "electronic-equipment")
	c[2][0] = "Sono attive le garanzie opzionali: "
	c[2][1] = "Maggiori costi  :" + IfString(ExistAsset(data.Assets, "increased-cost-of-working"), "fino a "+GetSumIndenity(data.Assets, "increased-cost-of-working"), "NO")
	c[2][2] = "Programmi in licenza d’uso: " + IfString(ExistAsset(data.Assets, "valuables"), "fino a "+GetSumIndenity(data.Assets, "software-under-license"), "NO")
	c[2][3] = "Supporto dati: " + IfString(ExistAsset(data.Assets, "restoration-of-data"), "fino a "+GetSumIndenity(data.Assets, "restoration-of-data"), "NO")
	c[3][1] = GetPrice(data.Assets, "electronic-equipment", "increased-cost-of-working", "restoration-of-data", "software-under-license")
	s.MultiRow(m, c, true, []uint{4, 2, 4, 2}, 40)
	//-----------------------------------------------------------------------
	c = lib.Make2D[string](4, 1)

	c[0][0] = "Business Interruption"
	c[1][0] = GetSumIndenity(data.Assets, "business-interruption")
	c[2][0] = "= ="
	c[3][0] = GetPrice(data.Assets, "business-interruption")
	s.MultiRow(m, c, true, []uint{4, 2, 4, 2}, 40)
	//-----------------------------------------------------------------------
	c = lib.Make2D[string](4, 1)
	c[0][0] = "Assistenza al Fabbricato"
	c[1][0] = IfString(ExistAsset(data.Assets, "assistance"), "inclusa", "= =")
	c[2][0] = "= ="
	c[3][0] = GetPrice(data.Assets, "assistance")
	s.MultiRow(m, c, true, []uint{4, 2, 4, 2}, 40)
	//-----------------------------------------------------------------------

}
