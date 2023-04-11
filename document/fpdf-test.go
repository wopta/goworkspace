package document

import (
	"encoding/json"
	"github.com/go-pdf/fpdf"
	"github.com/wopta/goworkspace/models"
	"log"
	"net/http"
	"os"
)

type Statements struct {
	Statements []*models.Statement `json:"statements"`
	Text       string              `json:"text,omitempty"`
}

func FpdfHandler(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {

	pdf := initFpdf()

	GetHeader(pdf, "life")
	GetFooter(pdf)

	pdf.AddPage()

	GetParagraphTitle(pdf, "La tua assicurazione è operante per il seguente Assicurato e Garanzie")
	pdf.Ln(8)
	DrawPinkLine(pdf, 0.4)
	pdf.Ln(2)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.Cell(20, 2, "Cognome e Nome")
	pdf.SetFont("Montserrat", "", 9)
	pdf.SetX(pdf.GetX() + 24)
	pdf.Cell(20, 2, "HAMMAR YOUSEF")
	pdf.SetX(pdf.GetX() + 60)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.Cell(10, 2, "Codice fiscale:")
	pdf.SetX(pdf.GetX() + 20)
	pdf.SetFont("Montserrat", "", 9)
	pdf.Cell(20, 2, "HMMYSF94R07D912M")
	pdf.Ln(2.5)
	DrawPinkLine(pdf, 0.1)
	pdf.Ln(2)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.Cell(20, 2, "Residente in")
	pdf.SetFont("Montserrat", "", 9)
	pdf.SetX(pdf.GetX() + 24)
	pdf.Cell(20, 2, "Via Unicef 4 - 20033 Solaro (MI)")
	pdf.Ln(2.5)
	DrawPinkLine(pdf, 0.1)
	pdf.Ln(2)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.Cell(20, 2, "Mail")
	pdf.SetFont("Montserrat", "", 9)
	pdf.SetX(pdf.GetX() + 24)
	pdf.Cell(20, 2, "yousef.hammar@wopta.it")
	pdf.SetX(pdf.GetX() + 60)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.Cell(10, 2, "Telefono")
	pdf.SetX(pdf.GetX() + 20)
	pdf.SetFont("Montserrat", "", 9)
	pdf.Cell(20, 2, "+393451021004")
	pdf.Ln(2.5)
	DrawPinkLine(pdf, 0.4)
	pdf.Ln(1)

	pdf.SetFont("Montserrat", "B", 9)
	pdf.MultiCell(90, 3, "Garanzie", "", "CM", false)
	pdf.SetXY(pdf.GetX()+90, pdf.GetY()-3)
	pdf.MultiCell(30, 3, "Somma\nassicurata €", "", "CM", false)
	pdf.SetXY(pdf.GetX()+120, pdf.GetY()-6)
	pdf.MultiCell(20, 3, "Durata\nanni", "", "CM", false)
	pdf.SetXY(pdf.GetX()+143, pdf.GetY()-6)
	pdf.MultiCell(25, 3, "Scade il", "", "CM", false)
	pdf.SetXY(pdf.GetX()+169, pdf.GetY()-3)
	pdf.MultiCell(25, 3, "Premio annuale €", "", "CM", false)

	/*
		pdf.CellFormat(80, 6, "Garanzie", "", 0, "CM", false, 0, "")
		pdf.CellFormat(25, 6, "Somma\nassicurata €", "", 0, "CM", false, 0, "")
		pdf.CellFormat(30, 6, "Durata anni", "", 0, "CM", false, 0, "")
		pdf.CellFormat(25, 6, "Scade il", "", 0, "CM", false, 0, "")
		pdf.CellFormat(25, 6, "Premio annuale €", "", 0, "CM", false, 0, "")

	*/
	pdf.Ln(1)
	DrawPinkLine(pdf, 0.2)
	pdf.CellFormat(90, 6, "Decesso", "", 0, "", false, 0, "")
	pdf.SetFont("Montserrat", "", 9)
	pdf.CellFormat(25, 6, "100.000 €", "", 0, "RM", false, 0, "")
	pdf.CellFormat(30, 6, "20", "", 0, "CM", false, 0, "")
	pdf.CellFormat(25, 6, "04/03/2023", "", 0, "CM", false, 0, "")
	pdf.CellFormat(25, 6, "97,74 €     ", "", 0, "RM", false, 0, "")
	pdf.Ln(5)
	DrawPinkLine(pdf, 0.2)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.CellFormat(90, 6, "Invalidità Totale Permanente da Infortunio o Malattia", "", 0, "", false, 0, "")
	pdf.SetFont("Montserrat", "", 9)
	pdf.CellFormat(25, 6, "100.000 €", "", 0, "RM", false, 0, "")
	pdf.CellFormat(30, 6, "20", "", 0, "CM", false, 0, "")
	pdf.CellFormat(25, 6, "04/03/2023", "", 0, "CM", false, 0, "")
	pdf.CellFormat(25, 6, "12,02 € (*)", "", 0, "RM", false, 0, "")
	pdf.Ln(5)
	DrawPinkLine(pdf, 0.2)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.CellFormat(90, 6, "Inabilità Temporanea da Infortunio o Malattia", "", 0, "", false, 0, "")
	pdf.SetFont("Montserrat", "", 9)
	pdf.CellFormat(25, 6, "500 €", "", 0, "RM", false, 0, "")
	pdf.CellFormat(30, 6, "10", "", 0, "CM", false, 0, "")
	pdf.CellFormat(25, 6, "04/03/2023", "", 0, "CM", false, 0, "")
	pdf.CellFormat(25, 6, "10,00 € (*)", "", 0, "RM", false, 0, "")
	pdf.Ln(5)
	DrawPinkLine(pdf, 0.2)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.CellFormat(90, 6, "Malattie Gravi", "", 0, "", false, 0, "")
	pdf.SetFont("Montserrat", "", 9)
	pdf.CellFormat(25, 6, "10.000 €", "", 0, "RM", false, 0, "")
	pdf.CellFormat(30, 6, "10", "", 0, "CM", false, 0, "")
	pdf.CellFormat(25, 6, "04/03/2023", "", 0, "CM", false, 0, "")
	pdf.CellFormat(25, 6, "55,42 € (*)", "", 0, "RM", false, 0, "")
	pdf.Ln(5)
	DrawPinkLine(pdf, 0.2)
	pdf.Ln(0.5)
	pdf.SetFont("Montserrat", "", 6)
	pdf.Cell(80, 3, "(*) imposte assicurative di legge incluse nella misura del 2,50% del premio imponibile")
	pdf.Ln(3)

	GetParagraphTitle(pdf, "Nomina dei Beneficiari e Referente terzo, per il caso di garanzia Decesso (qualora sottoscritta)")
	pdf.Ln(8)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Montserrat", "", 9)
	pdf.MultiCell(0, 3, "AVVERTENZE: Può scegliere se designare nominativamente i beneficiari o se designare genericamente come beneficiari i suoi eredi legittimi e/o testamentari. In caso di mancata designazione nominativa, la Compagnia potrà incontrare, al decesso dell’Assicurato, maggiori difficoltà nell’identificazione e nella ricerca dei beneficiari. La modifica o revoca del/i beneficiario/i deve essere comunicata alla Compagnia in forma scritta.\nIn caso di specifiche esigenze di riservatezza, la Compagnia potrà rivolgersi ad un soggetto terzo (diverso dal Beneficiario)\nIn caso di Decesso al fine di contattare il Beneficiario designato.", "", "", false)
	pdf.Ln(0)

	GetParagraphTitle(pdf, "Beneficiario")
	pdf.Ln(8)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Montserrat", "", 9)
	pdf.CellFormat(0, 3, "Io sottoscritto Assicurato, con la sottoscrizione della presente polizza, in riferimento alla garanzia Decesso:", "", 0, "", false, 0, "")
	pdf.Ln(4)
	pdf.SetDrawColor(0, 0, 0)
	pdf.SetX(11)
	pdf.CellFormat(3, 3, "X", "1", 0, "CM", false, 0, "")
	pdf.CellFormat(0, 3, "Designo genericamente quali beneficiari della prestazione i miei eredi (legittimi e/o testamentari)", "", 0, "", false, 0, "")
	pdf.Ln(4)
	pdf.SetX(11)
	pdf.CellFormat(3, 3, "", "1", 0, "CM", false, 0, "")
	pdf.CellFormat(0, 3, "Designo nominativamente il/i seguente/i soggetto/i quale beneficiario/i della prestazione", "", 0, "", false, 0, "")
	pdf.Ln(5)
	DrawPinkLine(pdf, 0.4)
	pdf.Ln(1)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.Cell(50, 2, "Cognome e nome")
	pdf.SetFont("Montserrat", "", 9)
	pdf.Cell(20, 2, "=====")
	pdf.SetX(pdf.GetX() + 60)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.Cell(20, 2, "Cod. Fisc.: ")
	pdf.SetFont("Montserrat", "", 9)
	pdf.Cell(20, 2, "=====")
	pdf.Ln(3)
	DrawPinkLine(pdf, 0.1)
	pdf.Ln(2)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.Cell(50, 2, "Indirizzo")
	pdf.SetFont("Montserrat", "", 9)
	pdf.Cell(20, 2, "=====")
	pdf.SetX(pdf.GetX() + 60)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.Ln(3)
	DrawPinkLine(pdf, 0.1)
	pdf.Ln(2)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.Cell(50, 2, "Mail")
	pdf.SetFont("Montserrat", "", 9)
	pdf.Cell(20, 2, "=====")
	pdf.SetX(pdf.GetX() + 60)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.Cell(20, 2, "Telefono: ")
	pdf.SetFont("Montserrat", "", 9)
	pdf.Cell(20, 2, "=====")
	pdf.Ln(3)
	DrawPinkLine(pdf, 0.1)
	pdf.Ln(2)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.Cell(50, 2, "Relazione con Assicurato")
	pdf.SetFont("Montserrat", "", 9)
	pdf.Cell(20, 2, "=====")
	pdf.Ln(3)
	DrawPinkLine(pdf, 0.1)
	pdf.Ln(2)
	pdf.Cell(165, 2, "Consenso ad invio comunicazioni da parte della Compagnia al beneficiario, prima dell'evento Decesso:")
	pdf.Cell(20, 2, "=====")
	pdf.Ln(3)
	DrawPinkLine(pdf, 0.1)
	pdf.Ln(2)
	DrawPinkLine(pdf, 0.4)
	pdf.Ln(1)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.Cell(50, 2, "Cognome e nome")
	pdf.SetFont("Montserrat", "", 9)
	pdf.Cell(20, 2, "=====")
	pdf.SetX(pdf.GetX() + 60)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.Cell(20, 2, "Cod. Fisc.: ")
	pdf.SetFont("Montserrat", "", 9)
	pdf.Cell(20, 2, "=====")
	pdf.Ln(3)
	DrawPinkLine(pdf, 0.1)
	pdf.Ln(2)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.Cell(50, 2, "Indirizzo")
	pdf.SetFont("Montserrat", "", 9)
	pdf.Cell(20, 2, "=====")
	pdf.SetX(pdf.GetX() + 60)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.Ln(3)
	DrawPinkLine(pdf, 0.1)
	pdf.Ln(2)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.Cell(50, 2, "Mail")
	pdf.SetFont("Montserrat", "", 9)
	pdf.Cell(20, 2, "=====")
	pdf.SetX(pdf.GetX() + 60)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.Cell(20, 2, "Telefono: ")
	pdf.SetFont("Montserrat", "", 9)
	pdf.Cell(20, 2, "=====")
	pdf.Ln(3)
	DrawPinkLine(pdf, 0.1)
	pdf.Ln(2)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.Cell(50, 2, "Relazione con Assicurato")
	pdf.SetFont("Montserrat", "", 9)
	pdf.Cell(20, 2, "=====")
	pdf.Ln(3)
	DrawPinkLine(pdf, 0.1)
	pdf.Ln(2)
	pdf.Cell(165, 2, "Consenso ad invio comunicazioni da parte della Compagnia al beneficiario, prima dell'evento Decesso:")
	pdf.Cell(20, 2, "=====")
	pdf.Ln(3)
	DrawPinkLine(pdf, 0.1)
	pdf.Ln(2)

	GetParagraphTitle(pdf, "Referente terzo")
	pdf.Ln(8)
	pdf.SetDrawColor(229, 0, 117)
	DrawPinkLine(pdf, 0.4)
	pdf.Ln(1)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.Cell(50, 2, "Cognome e nome")
	pdf.SetFont("Montserrat", "", 9)
	pdf.Cell(20, 2, "=====")
	pdf.SetX(pdf.GetX() + 60)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.Cell(20, 2, "Cod. Fisc.: ")
	pdf.SetFont("Montserrat", "", 9)
	pdf.Cell(20, 2, "=====")
	pdf.Ln(3)
	DrawPinkLine(pdf, 0.1)
	pdf.Ln(2)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.Cell(50, 2, "Indirizzo")
	pdf.SetFont("Montserrat", "", 9)
	pdf.Cell(20, 2, "=====")
	pdf.SetX(pdf.GetX() + 60)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.Ln(3)
	DrawPinkLine(pdf, 0.1)
	pdf.Ln(2)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.Cell(50, 2, "Mail")
	pdf.SetFont("Montserrat", "", 9)
	pdf.Cell(20, 2, "=====")
	pdf.SetX(pdf.GetX() + 60)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.Cell(20, 2, "Telefono: ")
	pdf.SetFont("Montserrat", "", 9)
	pdf.Cell(20, 2, "=====")
	pdf.Ln(3)
	DrawPinkLine(pdf, 0.1)
	pdf.Ln(2)

	var statements Statements
	b, err := os.ReadFile("document/response.json")
	if err != nil {
		log.Println(err)
	}
	err = json.Unmarshal(b, &statements)
	if err != nil {
		log.Println(err)
	}

	GetParagraphTitle(pdf, "Dichiarazioni da leggere con attenzione prima di firmare")
	pdf.Ln(8)
	PrintStatement(pdf, statements.Statements[0])
	pdf.Ln(5)
	GetParagraphTitle(pdf, "Questionario Medico")
	pdf.Ln(8)
	for _, statement := range statements.Statements[1:] {
		PrintStatement(pdf, statement)
	}
	pdf.Ln(8)
	pdf.SetX(-80)
	pdf.Cell(0, 3, "Firma del Contraente/Assicurato")
	pdf.Ln(15)
	pdf.SetDrawColor(0, 0, 0)
	pdf.SetLineWidth(0.4)
	pdf.Line(130, pdf.GetY(), 190, pdf.GetY())

	pdf.AddPage()

	GetParagraphTitle(pdf, "Presa visione dei documenti precontrattuali e sottoscrizione Polizza")
	pdf.Ln(8)
	pdf.SetFont("Montserrat", "", 8)
	pdf.SetTextColor(0, 0, 0)
	email := "yousef.hammar@wopta.it"
	pdf.MultiCell(0, 3, "Ho scelto la ricezione della seguente documentazione via e-mail al seguente indirizzo "+
		"indirizzo"+email+", nonché all’utilizzo della stessa per l’invio delle comunicazioni in corso di contratto da "+
		"parte di Wopta e della Compagnia. Sono a conoscenza che, qualora volessi modificare questa mia scelta potrò "+
		"farlo scrivendo alla Compagnia, con le modalità previste "+
		"nelle Condizioni di Assicurazione.", "", "", false)
	pdf.SetFont("Montserrat", "B", 8)
	pdf.MultiCell(0, 3, "Confermo di aver ricevuto e preso visione, prima della conclusione del contratto:", "", "", false)
	pdf.SetX(15)
	pdf.MultiCell(0, 3, "1. degli Allegati 3, 4 e 4-ter, di cui al Regolamento IVASS n. 40/2018, relativi "+
		"agli obblighi informativi e di comportamento dell’Intermediario, inclusa l’informativa privacy "+
		"dell’intermediario (ai sensi dell’art. 13 del regolamento UE n. 2016/679);", "", "", false)
	pdf.SetX(15)
	pdf.MultiCell(0, 3, "2.	del Set informativo, identificato dal modello AX01.0323, contenente: "+
		"1) Documento informativo precontrattuale per i prodotti assicurativi vita diversi dai prodotti "+
		"d’investimento assicurativi (DIP Vita), Documento informativo per i prodotti assicurativi danni (DIP Danni), "+
		"Documento informativo precontrattuale aggiuntivo per i prodotti assicurativi multirischi "+
		"(DIP aggiuntivo Multirischi), di cui al Regolamento IVASS n. 41/2018; 2) Condizioni di Assicurazione "+
		"comprensive di Glossario, che dichiaro altresì di conoscere ed accettare.", "", "", false)
	pdf.SetX(40)
	pdf.Cell(30, 3, "AXA France Vie")
	pdf.SetX(-75)
	pdf.Cell(40, 3, "Firma del Contraente/Assicurato")
	pdf.Ln(2)
	pdf.SetX(25)
	pdf.Cell(20, 3, "(Rappresentanza Generale per l'Italia)")
	var opt fpdf.ImageOptions
	opt.ImageType = "png"
	pdf.ImageOptions("document/assets/firma_axa.png", 35, pdf.GetY()+5, 30, 8, false, opt, 0, "")
	pdf.Ln(10)
	pdf.SetDrawColor(0, 0, 0)
	pdf.SetLineWidth(0.3)
	pdf.Line(130, pdf.GetY(), 190, pdf.GetY())
	pdf.Ln(5)

	GetParagraphTitle(pdf, "Il premio per tutte le coperture assicurative attivate sulla polizza – Frazionamento: ANNUALE")
	pdf.Ln(8)
	pdf.SetFont("Montserrat", "", 7)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(40, 2, "Premio", "", 0, "", false, 0, "")
	pdf.SetX(pdf.GetX() + 20)
	pdf.CellFormat(40, 2, "Imponibile", "", 0, "", false, 0, "")
	pdf.SetX(pdf.GetX() + 15)
	pdf.CellFormat(40, 2, "Imposte Assicurative", "", 0, "", false, 0, "")
	pdf.SetX(pdf.GetX() + 15)
	pdf.CellFormat(40, 2, "Totale", "RM", 0, "", false, 0, "")
	pdf.Ln(3)
	DrawPinkLine(pdf, 0.1)
	pdf.Ln(1)
	pdf.CellFormat(40, 2, "Annuale firma del contratto", "", 0, "", false, 0, "")
	pdf.SetX(pdf.GetX() + 20)
	pdf.CellFormat(40, 2, "€ 173,30", "", 0, "", false, 0, "")
	pdf.SetX(pdf.GetX() + 15)
	pdf.CellFormat(40, 2, "€ 1,88", "", 0, "", false, 0, "")
	pdf.SetX(pdf.GetX() + 15)
	pdf.CellFormat(40, 2, "€ 175,18", "RM", 0, "", false, 0, "")
	pdf.Ln(3)
	DrawPinkLine(pdf, 0.1)
	pdf.Ln(5)

	GetParagraphTitle(pdf, "Pagamento dei premi successivi al primo")
	pdf.Ln(8)
	pdf.SetFont("Montserrat", "", 9)
	pdf.SetTextColor(0, 0, 0)
	pdf.MultiCell(0, 3, "Il Contraente è tenuto a pagare i Premi entro 30 giorni dalle relative scadenze. "+
		"In caso di mancato pagamento del premio entro 30 giorni dalla scadenza (c.d. termine di tolleranza) "+
		"l’assicurazione è sospesa. Il contratto è risolto automaticamente in caso di mancato pagamento "+
		"del Premio entro 90 giorni dalla scadenza.", "", "", false)
	pdf.Ln(3)
	DrawPinkLine(pdf, 0.4)
	pdf.Ln(1)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.CellFormat(pdf.GetStringWidth("Tipologia di premio"), 3, "Tipologia di premio:", "", 0, "", false, 0, "")
	pdf.SetX(pdf.GetX() + 5)
	pdf.SetFont("Montserrat", "", 9)
	pdf.SetDrawColor(0, 0, 0)
	pdf.CellFormat(3, 3, "", "1", 0, "", false, 0, "")
	pdf.SetX(pdf.GetX() + 1)
	pdf.CellFormat(pdf.GetStringWidth("natuale variabile annualmente"), 3, "naturale variabile annualmente", "", 0, "", false, 0, "")
	pdf.SetX(pdf.GetX() + 5)
	pdf.CellFormat(3, 3, "X", "1", 0, "CM", false, 0, "")
	pdf.SetX(pdf.GetX() + 1)
	pdf.CellFormat(pdf.GetStringWidth("fisso"), 3, "fisso", "", 0, "", false, 0, "")
	pdf.SetX(pdf.GetX() + 40)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.CellFormat(pdf.GetStringWidth("Frazionamento"), 3, "Frazionamento:", "", 0, "", false, 0, "")
	pdf.SetX(pdf.GetX() + 3)
	pdf.CellFormat(pdf.GetStringWidth("ANNUALE"), 3, "ANNUALE", "", 0, "", false, 0, "")
	pdf.Ln(4)
	DrawPinkLine(pdf, 0.1)
	pdf.Ln(1)
	pdf.SetFont("Montserrat", "", 9)
	pdf.Cell(0, 3, "Il Premio è dovuto alle diverse annualità di Polizza, alle date qui sotto indicate:")
	pdf.Ln(4)
	DrawPinkLine(pdf, 0.1)
	pdf.Ln(1)
	pdf.SetX(pdf.GetX() + 12)
	pdf.CellFormat(pdf.GetStringWidth("04/03/2030:"), 3, "Alla firma:", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 1)
	pdf.CellFormat(pdf.GetStringWidth("€ 175,18"), 3, "€ 175,18", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 12)
	pdf.CellFormat(pdf.GetStringWidth("04/03/2030:"), 3, "04/03/2028:", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 1)
	pdf.CellFormat(pdf.GetStringWidth("€ 175,18"), 3, "€ 175,18", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 12)
	pdf.CellFormat(pdf.GetStringWidth("04/03/2030:"), 3, "04/03/2033:", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 1)
	pdf.CellFormat(pdf.GetStringWidth("€ 175,18"), 3, "€ 175,18", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 12)
	pdf.CellFormat(pdf.GetStringWidth("04/03/2030:"), 3, "04/03/2038:", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 1)
	pdf.CellFormat(pdf.GetStringWidth("€ 175,18"), 3, "€ 175,18", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 12)
	pdf.Ln(4)
	DrawPinkLine(pdf, 0.1)
	pdf.Ln(1)
	pdf.SetX(pdf.GetX() + 12)
	pdf.CellFormat(pdf.GetStringWidth("04/03/2030:"), 3, "03/04/2024", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 1)
	pdf.CellFormat(pdf.GetStringWidth("€ 175,18"), 3, "€ 175,18", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 12)
	pdf.CellFormat(pdf.GetStringWidth("04/03/2030:"), 3, "04/03/2029:", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 1)
	pdf.CellFormat(pdf.GetStringWidth("€ 175,18"), 3, "€ 175,18", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 12)
	pdf.CellFormat(pdf.GetStringWidth("04/03/2030:"), 3, "04/03/2034:", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 1)
	pdf.CellFormat(pdf.GetStringWidth("€ 175,18"), 3, "€ 175,18", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 12)
	pdf.CellFormat(pdf.GetStringWidth("04/03/2030:"), 3, "04/03/2039:", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 1)
	pdf.CellFormat(pdf.GetStringWidth("€ 175,18"), 3, "€ 175,18", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 12)
	pdf.Ln(4)
	DrawPinkLine(pdf, 0.1)
	pdf.Ln(1)
	pdf.SetX(pdf.GetX() + 12)
	pdf.CellFormat(pdf.GetStringWidth("04/03/2030:"), 3, "03/04/2025", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 1)
	pdf.CellFormat(pdf.GetStringWidth("€ 175,18"), 3, "€ 175,18", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 12)
	pdf.CellFormat(pdf.GetStringWidth("04/03/2030:"), 3, "04/03/2030:", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 1)
	pdf.CellFormat(pdf.GetStringWidth("€ 175,18"), 3, "€ 175,18", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 12)
	pdf.CellFormat(pdf.GetStringWidth("04/03/2030:"), 3, "04/03/2035:", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 1)
	pdf.CellFormat(pdf.GetStringWidth("€ 175,18"), 3, "€ 175,18", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 12)
	pdf.CellFormat(pdf.GetStringWidth("04/03/2030:"), 3, "04/03/2040:", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 1)
	pdf.CellFormat(pdf.GetStringWidth("€ 175,18"), 3, "€ 175,18", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 12)
	pdf.Ln(4)
	DrawPinkLine(pdf, 0.1)
	pdf.Ln(1)
	pdf.SetX(pdf.GetX() + 12)
	pdf.CellFormat(pdf.GetStringWidth("04/03/2030:"), 3, "03/04/2026", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 1)
	pdf.CellFormat(pdf.GetStringWidth("€ 175,18"), 3, "€ 175,18", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 12)
	pdf.CellFormat(pdf.GetStringWidth("04/03/2030:"), 3, "04/03/2031:", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 1)
	pdf.CellFormat(pdf.GetStringWidth("€ 175,18"), 3, "€ 175,18", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 12)
	pdf.CellFormat(pdf.GetStringWidth("04/03/2030:"), 3, "04/03/2036:", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 1)
	pdf.CellFormat(pdf.GetStringWidth("€ 175,18"), 3, "€ 175,18", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 12)
	pdf.CellFormat(pdf.GetStringWidth("04/03/2030:"), 3, "04/03/2041:", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 1)
	pdf.CellFormat(pdf.GetStringWidth("€ 175,18"), 3, "€ 175,18", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 12)
	pdf.Ln(4)
	DrawPinkLine(pdf, 0.1)
	pdf.Ln(1)
	pdf.SetX(pdf.GetX() + 12)
	pdf.CellFormat(pdf.GetStringWidth("04/03/2030:"), 3, "03/04/2027", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 1)
	pdf.CellFormat(pdf.GetStringWidth("€ 175,18"), 3, "€ 175,18", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 12)
	pdf.CellFormat(pdf.GetStringWidth("04/03/2030:"), 3, "04/03/2032:", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 1)
	pdf.CellFormat(pdf.GetStringWidth("€ 175,18"), 3, "€ 175,18", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 12)
	pdf.CellFormat(pdf.GetStringWidth("04/03/2030:"), 3, "04/03/2037:", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 1)
	pdf.CellFormat(pdf.GetStringWidth("€ 175,18"), 3, "€ 175,18", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 12)
	pdf.CellFormat(pdf.GetStringWidth("04/03/2030:"), 3, "04/03/2042:", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 1)
	pdf.CellFormat(pdf.GetStringWidth("€ 175,18"), 3, "€ 175,18", "", 0, "RM", false, 0, "")
	pdf.SetX(pdf.GetX() + 12)
	pdf.Ln(4)
	DrawPinkLine(pdf, 0.1)
	pdf.Ln(1)
	pdf.MultiCell(0, 3, "In caso di frazionamento mensile i Premi sopra riportati sono dovuti, alle date "+
		"indicate e con successiva frequenza mensile, in misura di 1/12 per ogni mensilità. Non sono previsti oneri "+
		"o interessi di frazionamento.", "", "", false)
	pdf.Ln(5)

	GetParagraphTitle(pdf, "Informativa sul diritto di recesso")
	pdf.Ln(8)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.MultiCell(0, 3, "Diritto di recesso entro i primi 30 giorni dalla stipula ("+
		"diritto di ripensamento)", "", "", false)
	pdf.SetFont("Montserrat", "", 9)
	pdf.MultiCell(0, 3, "Il Contraente può recedere dal contratto entro il termine di 30 giorni dalla "+
		"decorrenza dell’assicurazione (diritto di ripensamento). In tal caso, l’assicurazione si intende come mai "+
		"entrata in vigore e la Compagnia, per il tramite dell’intermediario, provvederà a rimborsare al Contraente "+
		"l’importo di Premio già versato (al netto delle imposte).", "", "", false)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.MultiCell(0, 3, "Diritto di recesso annuale (disdetta alla annualità)", "", "", false)
	pdf.SetFont("Montserrat", "", 9)
	pdf.MultiCell(0, 3, "Il Contraente può recedere dal contratto annualmente, entro il termine di 30 "+
		"giorni dalla scadenza annuale della polizza (disdetta alla annualità). In tal caso, l’assicurazione cessa alle "+
		"ore 24:00 dell’ultimo giorno della annualità in corso. È possibile disdettare singolarmente una o più delle "+
		"coperture attivate in fase di sottoscrizione.", "", "", false)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.MultiCell(0, 3, "Modalità per l’esercizio del diritto di recesso", "", "", false)
	pdf.SetFont("Montserrat", "", 9)
	pdf.MultiCell(0, 3, "Il Contraente è tenuto ad esercitare il diritto di recesso mediante invio di una "+
		"lettera raccomandata a.r. al seguente indirizzo: Wopta Assicurazioni srl – Gestione Portafoglio – Galleria del "+
		"Corso, 1 – 201212 Milano (MI) oppure via posta elettronica certificata (PEC) all’indirizzo "+
		"email: woptaassicurazioni@legalmail.it", "", "", false)
	pdf.Ln(5)

	GetParagraphTitle(pdf, "Come puoi pagare il premio")
	pdf.Ln(8)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Montserrat", "", 9)
	pdf.MultiCell(0, 3, "I mezzi di pagamento consentiti, nei confronti di Wopta, sono esclusivamente "+
		"bonifico e strumenti di pagamento elettronico, quali ad esempio, carte di credito e/o carte di debito, "+
		"incluse le carte prepagate. Oppure può essere pagato direttamente alla Compagnia alla "+
		"stipula del contratto, via bonifico o carta di credito.", "", "", false)
	pdf.Ln(5)

	GetParagraphTitle(pdf, "Emissione polizza e pagamento della prima rata")
	pdf.Ln(8)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Montserrat", "", 9)
	pdf.MultiCell(0, 3, "Polizza emessa a Milano il 03/04/2023 per un importo di €  175,18 quale prima "+
		"rata alla firma, il cui pagamento a saldo è da effettuarsi con i metodi di pagamento sopra indicati. "+
		"Wopta conferma avvenuto incasso e copertura della polizza dal 03/04/2023.", "", "", false)
	pdf.Ln(5)

	GetParagraphTitle(pdf, "Per noi questa polizza fa al caso tuo")
	pdf.Ln(8)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.MultiCell(0, 3, "Richieste ed esigenze di copertura assicurativa del contraente", "", "", false)
	pdf.SetFont("Montserrat", "", 9)
	pdf.MultiCell(0, 3, "In Polizza sono riportate le tue dichiarazioni relative al rischio. Sulla base di "+
		"tali dichiarazioni, esigenze e richieste, le soluzioni assicurative individuate in coerenza con esse, sono:", "", "", false)
	pdf.SetX(15)
	pdf.MultiCell(0, 3, "- tutelare dei soggetti cari, in caso di decesso dell’Assicurato nel corso della "+
		"durata della copertura, attraverso un sostegno economico indennizzato ai Beneficiari, che "+
		"sono stati indicati dal Contraente;", "", "", false)
	pdf.SetX(15)
	pdf.MultiCell(0, 3, "-	difendere il proprio reddito, potendo beneficiare di un capitale in caso di "+
		"perdita, da parte dell’Assicurato, definitiva ed irrimediabile, da Infortunio o Malattia, della capacità di "+
		"attendere a un qualsiasi lavoro proficuo, in misura totale (almeno del 60%);", "", "", false)
	pdf.SetX(15)
	pdf.MultiCell(0, 3, "-	difendere il proprio reddito attraverso una indennità mensile, in caso di "+
		"perdita totale, ma in via temporanea, delle capacità dell’Assicurato di attendere alla propria professione o "+
		"attività lavorativa a seguito di Infortunio o Malattia", "", "", false)
	pdf.SetX(15)
	pdf.MultiCell(0, 3, "-	tutelare la propria salute attraverso una indennità in caso di diagnosi, in "+
		"capo all’Assicurato, di una delle seguenti malattie: cancro, attacco cardiaco (infarto del miocardio), "+
		"chirurgia aorto-coronarica (bypass), ictus, insufficienza renale (fase finale di malattia renale), trapianto "+
		"di organi principali (cuore, polmone, fegato, pancreas, rene o midollo osseo);", "", "", false)
	pdf.SetX(15)
	pdf.MultiCell(0, 3, "- le coperture operano con un orizzonte l’orizzonte temporale di protezione "+
		"indicato a pagina 1 per ogni garanzia sottoscritta e tale durata viene valutata congrua "+
		"con le esigenze di protezione;", "", "", false)
	pdf.MultiCell(0, 3, "non rilevando interesse per altre eventuali coperture previste dal prodotto, ma "+
		"non incluse in questa Polizza. Il Contraente è stato informato che la Polizza può prevedere, in relazione "+
		"alle garanzie che precedono, l’applicazione di Scoperti, Franchigie, Limiti di indennizzo ed esclusioni, "+
		"meglio riportate nelle Condizioni Generali di Assicurazione, e che sono stati da te valutati in linea con la "+
		"capacità finanziaria di sostenere in proprio tale livello di danno e rischio.", "", "", false)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.MultiCell(0, 3, "Con la seguente sottoscrizione dichiari che quanto sopra corrisponde a quanto "+
		"illustrato dall’Intermediario, il quale ha fornito ogni altro elemento utile a consentirti di prendere una "+
		"decisione informata e coerente con le esigenze espresse.", "", "", false)
	pdf.Ln(5)
	pdf.SetX(-80)
	pdf.Cell(0, 3, "Firma del Contraente/Assicurato")
	pdf.Ln(15)
	pdf.SetDrawColor(0, 0, 0)
	pdf.SetLineWidth(0.4)
	pdf.Line(130, pdf.GetY(), 190, pdf.GetY())
	pdf.Ln(5)

	GetParagraphTitle(pdf, "Chi siamo")
	pdf.Ln(8)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Montserrat", "", 8)
	pdf.MultiCell(0, 3, "Wopta Assicurazioni S.r.l. - intermediario assicurativo, soggetto al controllo "+
		"dell’IVASS ed iscritto dal 14.02.2022 al Registro Unico degli Intermediari, in Sezione A nr. A000701923, "+
		"avente sede legale in Galleria del Corso, 1 – 20122 Milano (MI). Capitale sociale Euro 120.000 - "+
		"Codice Fiscale, Reg. Imprese e Partita IVA: 12072020964 - Iscritta al Registro delle imprese di Milano – "+
		"REA MI 2638708", "", "", false)
	pdf.Ln(5)
	pdf.MultiCell(0, 3, "AXA France Vie (compagnia assicurativa del gruppo AXA). Indirizzo sede legale in "+
		"Francia: 313 Terrasses de l'Arche, 92727 NANTERRE CEDEX. Numero Iscrizione Registro delle Imprese di "+
		"Nanterre: 310499959. Autorizzata in Francia (Stato di origine) all’esercizio delle assicurazioni, vigilata "+
		"in Francia dalla Autorité de Contrôle Prudentiel et de Résolution (ACPR). Numero Matricola Registre des "+
		"organismes d’assurance: 5020051. // Indirizzo Rappresentanza Generale per l’Italia: Corso Como n. 17, 20154 "+
		"Milano - CF, P.IVA e N.Iscr. Reg. Imprese 08875230016 - REA MI-2525395 - Telefono: 02-87103548 - "+
		"Fax: 02-23331247 - PEC: axafrancevie@legalmail.it – sito internet: www.clp.partners.axa/it. Ammessa ad "+
		"operare in Italia in regime di stabilimento. Iscritta all’Albo delle imprese di assicurazione tenuto "+
		"dall’IVASS, in appendice Elenco I, nr. I.00149.", "", "", false)

	pdf.SetHeaderFunc(func() {
		pdf.SetXY(-30, 7)
		opt.ImageType = "png"
		pdf.ImageOptions("document/assets/logo_axa.png", 190, 7, 8, 8, false, opt, 0, "")
		pdf.Ln(25)
	})

	pdf.AddPage()

	pdf.SetFooterFunc(func() {
		pdf.SetXY(10, -30)
		pdf.SetFont("Montserrat", "", 8)
		pdf.SetTextColor(0, 0, 0)
		pdf.MultiCell(0, 3, "AXA France Vie (compagnia assicurativa del gruppo AXA). Indirizzo sede "+
			"legale in Francia: 313 Terrasses de l'Arche, 92727 NANTERRE CEDEX. Numero Iscrizione Registro delle "+
			"Imprese di Nanterre: 310499959. Autorizzata in Francia (Stato di origine) all'esercizio delle "+
			"assicurazioni, vigilata in Francia dalla Autorité de Contrôle Prudentiel et de Résolution (ACPR). "+
			"Numero Matricola Registre des organismes d'assurance: 5020051. // Indirizzo Rappresentanza Generale "+
			"per l'Italia: Corso Como n. 17, 20154 Milano - CF, P.IVA e N.Iscr. Reg. Imprese 08875230016 - "+
			"REA MI-2525395 - Telefono: 02-87103548 - Fax: 02-23331247 - PEC: axafrancevie@legalmail.it - sito "+
			"internet: www.clp.partners.axa/it. Ammessa ad operare in Italia in regime di stabilimento. Iscritta "+
			"all'Albo delle imprese di assicurazione tenuto dall'IVASS, in appendice Elenco I, nr. I.00149.", "", "", false)
	})

	pdf.SetFont("Montserrat", "B", 9)
	pdf.SetTextColor(0, 0, 0)
	pdf.Cell(0, 3, "DICHIARAZIONI E CONSENSI")
	pdf.Ln(3)
	pdf.SetFont("Montserrat", "", 9)
	pdf.MultiCell(0, 3, "Io sottoscritto, dopo aver letto l’Informativa Privacy della compagnia titolare "+
		"del trattamento redatta ai sensi del Regolamento (UE) 2016/679 (relativo alla protezione delle persone "+
		"fisiche con riguardo al trattamento dei dati personali), della quale confermo ricezione, PRESTO IL CONSENSO "+
		"al trattamento dei miei dati personali, ivi inclusi quelli eventualmente da me conferiti in riferimento al "+
		"mio stato di salute, per le finalità indicate nell’informativa, nonché alla loro comunicazione, per "+
		"successivo trattamento, da parte dei soggetti indicati nella informativa predetta.", "", "", false)
	pdf.Ln(3)
	pdf.SetFont("Montserrat", "B", 8)
	pdf.Cell(0, 3, "Resta inteso che in caso di negazione del consenso non sarà possibile "+
		"finalizzare il rapporto contrattuale assicurativo.")
	pdf.Ln(3)
	pdf.SetDrawColor(0, 0, 0)
	pdf.SetLineWidth(0.1)
	pdf.Line(11, pdf.GetY(), 180, pdf.GetY())
	pdf.Ln(5)
	pdf.Cell(0, 3, "03/04/2023")
	pdf.Ln(8)
	pdf.SetX(-80)
	pdf.Cell(0, 3, "Firma del Contraente/Assicurato")
	pdf.Ln(15)
	pdf.SetDrawColor(0, 0, 0)
	pdf.SetLineWidth(0.4)
	pdf.Line(130, pdf.GetY(), 190, pdf.GetY())

	pdf.AddPage()

	pdf.SetFont("Montserrat", "B", 12)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFillColor(229, 0, 117)
	pdf.MultiCell(0, 6, "MODULO PER L’IDENTIFICAZIONE E L’ADEGUATA VERIFICA DELLA CLIENTELA", "LTR", "CM", true)
	pdf.SetFont("Montserrat", "B", 8)
	pdf.MultiCell(0, 4, "POLIZZA DI RAMO VITA I  - Polizza “Wopta per te. Vita”", "LR", "CM", true)
	pdf.SetFont("Montserrat", "I", 6)
	pdf.MultiCell(0, 3, "(da compilarsi in caso di scelta da parte del Contraente/Assicurato della garanzia Decesso)", "LBR", "CM", true)

	err = pdf.OutputFileAndClose("document/test.pdf")
	log.Println(err)
	return "", nil, err
}

func initFpdf() *fpdf.Fpdf {
	pdf := fpdf.New(fpdf.OrientationPortrait, fpdf.UnitMillimeter, fpdf.PageSizeA4, "")
	pdf.SetMargins(10, 15, 10)
	loadCustomFonts(pdf)
	return pdf
}

func loadCustomFonts(pdf *fpdf.Fpdf) {
	pdf.AddUTF8Font("Montserrat", "", "document/assets/montserrat_regular.ttf")
	pdf.AddUTF8Font("Montserrat", "B", "document/assets/montserrat_bold.ttf")
	pdf.AddUTF8Font("Montserrat", "I", "document/assets/montserrat_italic.ttf")
}

func DrawPinkLine(pdf *fpdf.Fpdf, lineWidth float64) {
	pdf.SetDrawColor(229, 0, 117)
	pdf.SetLineWidth(lineWidth)
	pdf.Line(11, pdf.GetY(), 200, pdf.GetY())
}

func GetHeader(pdf *fpdf.Fpdf, name string) {
	var opt fpdf.ImageOptions
	var product, logoPath string
	pathPrefix := "document/assets/logo_"

	switch name {
	case "life":
		product = "Vita"
		logoPath = pathPrefix + "vita.png"
	}

	pdf.SetHeaderFunc(func() {
		opt.ImageType = "png"
		pdf.ImageOptions(logoPath, 10, 6, 13, 13, false, opt, 0, "")
		pdf.SetXY(23, 7)
		pdf.SetTextColor(229, 0, 117)
		pdf.SetFont("Montserrat", "B", 18)
		pdf.Cell(10, 6, "Wopta per te")
		pdf.SetFont("Montserrat", "I", 18)
		pdf.SetXY(23, 13)
		pdf.SetTextColor(92, 89, 92)
		pdf.Cell(10, 6, product)
		pdf.ImageOptions("document/assets/ARTW_LOGO_RGB_400px.png", 158, 6, 0, 10, false, opt, 0, "")

		pdf.SetTextColor(0, 0, 0)
		pdf.SetFont("Montserrat", "B", 8)
		pdf.SetXY(11, 20)
		pdf.Cell(0, 3, "I dati della tua polizza")
		pdf.SetFont("Montserrat", "", 8)
		pdf.SetXY(11, pdf.GetY()+3)
		pdf.MultiCell(0, 3, "Numero: 12345\nDecorre dal: 03/04/2023 ore 24:00\nScade il: 03/04/2043 ore 24:00\nPrima scadenza annuale il: 04/04/2024\nNon si rinnova a scadenza.", "", "", false)

		pdf.SetFont("Montserrat", "B", 8)
		pdf.SetXY(-90, 20)
		pdf.Cell(0, 3, "I tuoi dati")
		pdf.SetFont("Montserrat", "", 8)
		pdf.SetXY(-90, pdf.GetY()+3)
		pdf.MultiCell(0, 3, "Contraente: HAMMAR YOUSEF\nC.F./P.IVA: HMMYSF94R07D912M\nIndirizzo: Via Unicef, 4\n20033 SOLARO (MI)\nMail: yousef.hammar@wopta.it\nTelefono: +393451031004", "", "", false)
		pdf.Ln(10)
	})
}

func GetFooter(pdf *fpdf.Fpdf) {
	var opt fpdf.ImageOptions

	pdf.SetFooterFunc(func() {
		pdf.SetXY(10, -13)
		pdf.SetTextColor(229, 0, 117)
		pdf.SetFont("Montserrat", "B", 6)
		pdf.MultiCell(0, 3, "Wopta per te. Vita è un prodotto assicurativo di AXA France Vie S.A. – Rappresentanza Generale per l’Italia\ndistribuito da Wopta Assicurazioni S.r.l.", "", "", false)
		opt.ImageType = "png"
		pdf.ImageOptions("document/assets/logo_axa.png", 190, 283, 8, 8, false, opt, 0, "")
	})
}

func GetParagraphTitle(pdf *fpdf.Fpdf, title string) {
	pdf.SetTextColor(229, 0, 117)
	pdf.SetFont("Montserrat", "B", 10)
	pdf.Cell(0, 10, title)
}

func PrintStatement(pdf *fpdf.Fpdf, statement *models.Statement) {
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.MultiCell(0, 4, statement.Title, "", "", false)
	for _, question := range statement.Questions {
		if question.IsBold {
			pdf.SetFont("Montserrat", "B", 6)
		}
		if question.Indent {
			pdf.SetX(15)
		}
		pdf.MultiCell(0, 4, question.Question, "", "", false)
	}
}
