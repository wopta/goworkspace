package document

import (
	"bytes"
	"fmt"

	"github.com/dustin/go-humanize"
	"github.com/go-pdf/fpdf"
)

type PolicyInfo struct {
	Company            string
	ProductDescription string
	Code               string
}

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
	EffectiveDate  string
	ExpirationDate string
	PriceGross     float64
}

type ReceiptInfo struct {
	PolicyInfo   PolicyInfo
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

	pdf.SetX(115)

	text := "Egr./Gent.le/Spett.le\n" +
		info.CustomerInfo.Fullname + "\n" +
		info.CustomerInfo.Address + "\n" +
		info.CustomerInfo.PostalCode + " " + info.CustomerInfo.City + " (" + info.CustomerInfo.Province + ")\n" +
		info.CustomerInfo.Email + " - " + info.CustomerInfo.Phone

	setBlackBoldFont(pdf, standardTextSize)
	pdf.MultiCell(0, 4, text, "", fpdf.AlignLeft, false)

	pdf.Ln(20)

	text = "Oggetto: Quietanza di pagamento polizza n. " + info.PolicyInfo.Code

	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 4, text, "", fpdf.AlignLeft, false)

	pdf.Ln(15)

	text = "Gentile Cliente,"
	pdf.MultiCell(0, 4, text, "", fpdf.AlignLeft, false)

	pdf.Ln(5)

	text = "la presente ricevuta è da considerarsi valida come quietanza di pagamento, Le consigliamo pertanto di " +
		"conservarla con la documentazione del Suo contratto assicurativo."
	pdf.MultiCell(0, 4, text, "", fpdf.AlignLeft, false)

	pdf.Ln(5)

	text = "La ringraziamo e Le porgiamo i nostri più cordiali saluti."
	pdf.MultiCell(0, 4, text, "", fpdf.AlignLeft, false)

	pdf.Ln(15)

	text = "RICEVUTA DI PAGAMENTO DEL PREMIO"
	setBlackBoldFont(pdf, standardTextSize)
	pdf.MultiCell(0, 4, text, "", fpdf.AlignCenter, false)

	pdf.Ln(5)

	text = "Qui di seguito riportiamo i dati riepilogativi riferiti al pagamento in oggetto:"
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 4, text, "", fpdf.AlignLeft, false)

	pdf.Ln(7)

	setBlackBoldFont(pdf, standardTextSize)

	formattedPriceGross := "€ " + humanize.FormatFloat("#.###,##", info.Transaction.PriceGross)
	table := [][]tableCell{
		{
			{
				text:      "CONTRAENTE:",
				height:    5,
				width:     pdf.GetStringWidth("CONTRAENTE:"),
				textBold:  true,
				fill:      false,
				fillColor: rgbColor{},
				align:     fpdf.AlignLeft,
				border:    "",
			},
			{
				text:      info.CustomerInfo.Fullname,
				height:    5,
				width:     100 - pdf.GetStringWidth("CONTRAENTE:"),
				textBold:  false,
				fill:      false,
				fillColor: rgbColor{},
				align:     fpdf.AlignLeft,
				border:    "",
			},
		},
		{
			{
				text:      "COMPAGNIA:",
				height:    5,
				width:     pdf.GetStringWidth("COMPAGNIA:"),
				textBold:  true,
				fill:      false,
				fillColor: rgbColor{},
				align:     fpdf.AlignLeft,
				border:    "",
			},
			{
				text:      info.PolicyInfo.Company,
				height:    5,
				width:     100 - pdf.GetStringWidth("COMPAGNIA:"),
				textBold:  false,
				fill:      false,
				fillColor: rgbColor{},
				align:     fpdf.AlignLeft,
				border:    "",
			},
		},
		{
			{
				text:      "N. POLIZZA:",
				height:    5,
				width:     pdf.GetStringWidth("N. POLIZZA:"),
				textBold:  true,
				fill:      false,
				fillColor: rgbColor{},
				align:     fpdf.AlignLeft,
				border:    "",
			},
			{
				text:      info.PolicyInfo.Code,
				height:    5,
				width:     90 - pdf.GetStringWidth("N. POLIZZA:"),
				textBold:  false,
				fill:      false,
				fillColor: rgbColor{},
				align:     fpdf.AlignLeft,
				border:    "",
			},
		},
		{
			{
				text:      "DESCRIZIONE:",
				height:    5,
				width:     pdf.GetStringWidth("DESCRIZIONE:"),
				textBold:  true,
				fill:      false,
				fillColor: rgbColor{},
				align:     fpdf.AlignLeft,
				border:    "",
			},
			{
				text:      info.PolicyInfo.ProductDescription,
				height:    5,
				width:     90 - pdf.GetStringWidth("DESCRIZIONE:"),
				textBold:  false,
				fill:      false,
				fillColor: rgbColor{},
				align:     fpdf.AlignLeft,
				border:    "",
			},
		},
		{
			{
				text:      "DECORRENZA:",
				height:    5,
				width:     pdf.GetStringWidth("DECORRENZA:"),
				textBold:  true,
				fill:      false,
				fillColor: rgbColor{},
				align:     fpdf.AlignLeft,
				border:    "",
			},
			{
				text:      info.Transaction.EffectiveDate,
				height:    5,
				width:     100 - pdf.GetStringWidth("DECORRENZA:"),
				textBold:  false,
				fill:      false,
				fillColor: rgbColor{},
				align:     fpdf.AlignLeft,
				border:    "",
			},
		},
		{
			{
				text:      "VALIDITA’ COPERTURA FINO AL:",
				height:    5,
				width:     pdf.GetStringWidth("VALIDITA’ COPERTURA FINO AL: "),
				textBold:  true,
				fill:      false,
				fillColor: rgbColor{},
				align:     fpdf.AlignLeft,
				border:    "",
			},
			{
				text:      info.Transaction.ExpirationDate,
				height:    5,
				width:     90 - pdf.GetStringWidth("VALIDITA’ COPERTURA FINO AL: "),
				textBold:  false,
				fill:      false,
				fillColor: rgbColor{},
				align:     fpdf.AlignLeft,
				border:    "",
			},
		},
		{
			{
				text:      "PREMIO PAGATO:",
				height:    5,
				width:     pdf.GetStringWidth("PREMIO PAGATO:"),
				textBold:  true,
				fill:      false,
				fillColor: rgbColor{},
				align:     fpdf.AlignLeft,
				border:    "",
			},
			{
				text:      formattedPriceGross,
				height:    5,
				width:     190 - pdf.GetStringWidth("PREMIO PAGATO:"),
				textBold:  false,
				fill:      false,
				fillColor: rgbColor{},
				align:     fpdf.AlignLeft,
				border:    "",
			},
		},
	}

	tableDrawer(pdf, table)

	pdf.Ln(7)

	text = fmt.Sprintf("Il premio relativo alla presente quietanza, pari a %s è stato incassato il "+
		"_____._________._____ in ___________________", formattedPriceGross)
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 4, text, "", fpdf.AlignLeft, false)

	pdf.Ln(15)

	text = "L'intermediario ______________________________"
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 4, text, "", fpdf.AlignLeft, false)

	err = pdf.Output(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
