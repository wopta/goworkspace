package document

import (
	"bytes"
	"encoding/base64"
	"log"
	"net/http"

	"github.com/go-pdf/fpdf"
)

func PaymentReceiptFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err error
	)

	log.SetPrefix("[PaymentReceiptFx] ")
	defer func() {
		if err != nil {
			log.Printf("error: %s", err.Error())
		}
		log.Println("Handler end -------------------------------------------------")
	}()

	doc, err := paymentReceipt()
	if err != nil {
		log.Printf("error: %s", err.Error())
		return "", nil, err
	}

	rawDoc := base64.StdEncoding.EncodeToString(doc)

	log.Println("Handler start -----------------------------------------------")

	return rawDoc, rawDoc, nil
}

func paymentReceipt() ([]byte, error) {
	var (
		err error
		buf bytes.Buffer
	)
	pdf := initFpdf()

	woptaHeader(pdf, false)

	woptaFooter(pdf)

	pdf.AddPage()

	pdf.SetX(130)

	text := "Egr./Gent.le/Spett.le\n" +
		"Mario Rossi\n" + // TODO: dynamic name and surname
		"Galleria del corso 1\n" + // TODO: dynamic address
		"20033 Milano (MI)\n" + // TODO: dynamic info
		"test@wopta - +393334455667" //TODO: dynamic mail and phone

	setBlackBoldFont(pdf, standardTextSize)
	pdf.MultiCell(0, 4, text, "", fpdf.AlignLeft, false)

	pdf.Ln(15)

	text = "Oggetto: Quietanza di pagamento polizza n. 100100" // TODO: dynamic policy companyCode

	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 4, text, "", fpdf.AlignLeft, false)

	pdf.Ln(5)

	text = "Gentile Cliente,"
	pdf.MultiCell(0, 4, text, "", fpdf.AlignLeft, false)

	pdf.Ln(5)

	text = "la presente ricevuta è da considerarsi valida come quietanza di pagamento, Le consigliamo pertanto di " +
		"conservarla con la documentazione del Suo contratto assicurativo."
	pdf.MultiCell(0, 4, text, "", fpdf.AlignLeft, false)

	pdf.Ln(5)

	text = "La ringraziamo e Le porgiamo i nostri più cordiali saluti."
	pdf.MultiCell(0, 4, text, "", fpdf.AlignLeft, false)

	pdf.Ln(17)

	text = "RICEVUTO DI PAGAMENTO DEL PREMIO"
	setBlackBoldFont(pdf, standardTextSize)
	pdf.MultiCell(0, 4, text, "", fpdf.AlignCenter, false)

	pdf.Ln(5)

	text = "Qui di seguito riportiamo i dati riepilogativi riferiti al pagamento in oggetto:"
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 4, text, "", fpdf.AlignLeft, false)

	pdf.Ln(3)

	text = "CONTRAENTE: Mario Rossi\n" + // TODO: dynamic contractor name
		"N. POLIZZA: 100100\n" + // TODO: dynamic policy codeCompany
		"EFFETTO COPERTURA: 03/10/2024\n" + // TODO: dynamic effective date
		"SCADENZA COPERTURA: 03/11/2024\n" + // TODO: dynamic transaction end date
		"PREMIO PAGATO: 100.00€\n" + // TODO: dynamic transaction priceGross
		"VALUTA INCASSO: 03/10/2024\n" + // TODO: dynamic info
		"PROSSIMO PAGAMENTO IL: 03/11/2024" // TODO: dynamic info

	setBlackBoldFont(pdf, standardTextSize)
	pdf.MultiCell(0, 12, text, "", fpdf.AlignLeft, false)

	pdf.Ln(5)

	text = "Milano, il 03/10/2024"
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 4, text, "", fpdf.AlignLeft, false)

	err = pdf.Output(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
