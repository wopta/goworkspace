package document

import (
	"bytes"
	"time"

	"github.com/go-pdf/fpdf"
)

type CustomerInfo struct {
	Fullname   string
	Address    string
	PostalCode string
	City       string
	Province   string
	Email      string
	Phone      string
}

type TransactionInfo struct {
	PolicyCode     string
	EffectiveDate  string
	ExpirationDate string
	PriceGross     string
	NextPayment    string
}

type ReceiptInfo struct {
	CustomerInfo CustomerInfo
	Transaction  TransactionInfo
}

func PaymentReceipt(info ReceiptInfo) ([]byte, error) {
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
		info.CustomerInfo.Fullname + "\n" +
		info.CustomerInfo.Address + "\n" +
		info.CustomerInfo.PostalCode + " " + info.CustomerInfo.City + " (" + info.CustomerInfo.Province + ")\n" +
		info.CustomerInfo.Email + " - " + info.CustomerInfo.Phone

	setBlackBoldFont(pdf, standardTextSize)
	pdf.MultiCell(0, 4, text, "", fpdf.AlignLeft, false)

	pdf.Ln(15)

	text = "Oggetto: Quietanza di pagamento polizza n. " + info.Transaction.PolicyCode

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

	text = "RICEVUTA DI PAGAMENTO DEL PREMIO"
	setBlackBoldFont(pdf, standardTextSize)
	pdf.MultiCell(0, 4, text, "", fpdf.AlignCenter, false)

	pdf.Ln(5)

	text = "Qui di seguito riportiamo i dati riepilogativi riferiti al pagamento in oggetto:"
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 4, text, "", fpdf.AlignLeft, false)

	pdf.Ln(3)

	text = "CONTRAENTE: " + info.CustomerInfo.Fullname + "\n" +
		"N. POLIZZA: " + info.Transaction.PolicyCode + "\n" +
		"EFFETTO COPERTURA: " + info.Transaction.EffectiveDate + "\n" +
		"SCADENZA COPERTURA: " + info.Transaction.ExpirationDate + "\n" +
		"PREMIO PAGATO: " + info.Transaction.PriceGross + "\n" +
		"PROSSIMO PAGAMENTO IL: " + info.Transaction.NextPayment

	setBlackBoldFont(pdf, standardTextSize)
	pdf.MultiCell(0, 12, text, "", fpdf.AlignLeft, false)

	pdf.Ln(5)

	text = "Milano, il " + time.Now().UTC().Format("02/01/2006")
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 4, text, "", fpdf.AlignLeft, false)

	err = pdf.Output(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
