package document

import (
	"fmt"
	"github.com/go-pdf/fpdf"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"log"
	"strings"
	"time"
)

func LifeReserved(pdf *fpdf.Fpdf, origin string, policy models.Policy) (string, []byte) {
	var (
	//cfpi, expiryInfo string
	)

	log.Println("[LifeReserved]")

	lifeReservedFooter(pdf)

	pdf.AddPage()

	lifeReservedHeader(pdf, policy)

	insuredInfoSection(pdf, &policy)

	guaranteesMap, slugs := loadLifeGuarantees(&policy)

	lifeGuaranteesTable(pdf, guaranteesMap, slugs)

	insuranceLimitSection(pdf)

	instructionsSection(pdf, policy)

	filename, out := saveReservedDocument(pdf, &policy)
	return filename, out
}

func lifeReservedHeader(pdf *fpdf.Fpdf, policy models.Policy) {
	var (
		opt                        fpdf.ImageOptions
		logoPath, cfpi, expiryInfo string
	)

	logoPath = lib.GetAssetPathByEnv(basePath) + "/logo_vita.png"

	location, err := time.LoadLocation("Europe/Rome")
	lib.CheckError(err)

	policyStartDate := policy.StartDate.In(location)
	policyEndDate := policy.EndDate.In(location)

	if policy.PaymentSplit == string(models.PaySplitMonthly) {
		expiryInfo = "Prima scandenza mensile il: " +
			policyStartDate.AddDate(0, 1, 0).Format(dateLayout)
	} else if policy.PaymentSplit == string(models.PaySplitYear) {
		expiryInfo = "Prima scadenza annuale il: " +
			policyStartDate.AddDate(1, 0, 0).Format(dateLayout)
	}

	proposalInfo := fmt.Sprintf("Numero: %d\nDecorrere dal: %s\nScade il: %s\n%s\nNon si rinnova a scadenza",
		policy.ProposalNumber, policyStartDate.Format(dateLayout), policyEndDate.Format(dateLayout), expiryInfo)

	contractor := policy.Contractor
	address := strings.ToUpper(contractor.Residence.StreetName + ", " + contractor.Residence.StreetNumber + "\n" +
		contractor.Residence.PostalCode + " " + contractor.Residence.City + " (" + contractor.Residence.CityCode + ")")

	if contractor.VatCode == "" {
		cfpi = contractor.FiscalCode
	} else {
		cfpi = contractor.VatCode
	}

	if policy.PaymentSplit == string(models.PaySplitMonthly) {
		expiryInfo = "Prima scandenza mensile il: " +
			policyStartDate.AddDate(0, 1, 0).Format(dateLayout) + "\n"
	} else if policy.PaymentSplit == string(models.PaySplitYear) {
		expiryInfo = "Prima scadenza annuale il: " +
			policyStartDate.AddDate(1, 0, 0).Format(dateLayout) + "\n"
	}

	contractorInfo := fmt.Sprintf("Contraente: %s\nC.F./P.IVA%s\nIndirizzo: %s\nMail: %s\nTelefono: %s",
		strings.ToUpper(contractor.Surname+" "+contractor.
			Name), cfpi, strings.ToUpper(address), contractor.Mail, contractor.Phone)

	opt.ImageType = "png"
	pdf.ImageOptions(lib.GetAssetPathByEnv(basePath)+"/axa/logo.png", 180, 15, 0, 20,
		false, opt, 0, "")
	pdf.SetY(pdf.GetY() + 22)

	setBlackBoldFont(pdf, 20)
	pdf.MultiCell(0, 3, "RAPPORTO DI VISITA MEDICA", "", fpdf.AlignCenter, false)
	pdf.Ln(5)
	setPinkBoldFont(pdf, 20)
	pdf.MultiCell(0, 3, policy.NameDesc, "", fpdf.AlignCenter, false)
	pdf.ImageOptions(logoPath, 100, pdf.GetY()+5, 15, 15, false, opt, 0, "")

	y := pdf.GetY() + 30
	setBlackBoldFont(pdf, standardTextSize)
	pdf.SetXY(11, y)
	pdf.Cell(0, 3, "I dati della tua proposta")
	setBlackRegularFont(pdf, standardTextSize)
	pdf.SetXY(11, pdf.GetY()+3)
	pdf.MultiCell(0, 3.5, proposalInfo, "", "", false)

	setBlackBoldFont(pdf, standardTextSize)
	pdf.SetXY(-90, y)
	pdf.Cell(0, 3, "I tuoi dati")
	setBlackRegularFont(pdf, standardTextSize)
	pdf.SetXY(-90, pdf.GetY()+3)
	pdf.MultiCell(0, 3.5, contractorInfo, "", "", false)
	pdf.Ln(10)
}

func lifeReservedFooter(pdf *fpdf.Fpdf) {
	pdf.SetFooterFunc(func() {
		pdf.SetY(-7)
		pageNumber(pdf)
	})
}

func insuranceLimitSection(pdf *fpdf.Fpdf) {
	text := "Limiti assuntivi:\n\n" +
		"Decesso: 75 anni a scadenza - max 500.000 euro\n" +
		"Invalidità Totale Permanente: 75 anni a scadenza - max 500.000 euro\n" +
		"Inabilità Totale Temporanea: 75 anni a scadenza  - max 3.000\n" +
		"Malattie Gravi: 65 anni a scadenza - max 100.000 euro"

	setBlackBoldFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, text, "", fpdf.AlignLeft, false)
	pdf.Ln(5)
}

func instructionsSection(pdf *fpdf.Fpdf, policy models.Policy) {
	setBlackDrawColor(pdf)
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3.5, "", "LTR", fpdf.AlignCenter, false)
	pdf.MultiCell(0, 3.5, "Da restituire in busta chiusa alla compagnia assicurativa, unitamente alle "+
		"schede “dati Polizza”,\n“Questionario Medico” e “Antiriciclaggio” compilate e sottoscritte in ogni sua parte, "+
		"alternativamente a:", "LR", fpdf.AlignCenter, false)

	setBlackBoldFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3.5, "", "LR", fpdf.AlignCenter, false)
	pdf.MultiCell(0, 3.5, "", "LR", fpdf.AlignCenter, false)
	pdf.MultiCell(0, 3.5, "Tramite Posta a:", "LR", fpdf.AlignCenter, false)
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3.5, "AXA PARTNERS", "LR", fpdf.AlignCenter, false)
	pdf.MultiCell(0, 3.5, "Ufficio Underwriting Medico – Corso Como n. 17 – 20154 MILANO",
		"LR", fpdf.AlignCenter,
		false)
	setBlackBoldFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3.5, "", "LR", fpdf.AlignCenter, false)
	pdf.MultiCell(0, 3.5, "", "LR", fpdf.AlignCenter, false)
	pdf.MultiCell(0, 3.5, "Tramite e-mail:", "LR", fpdf.AlignCenter,
		false)
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3.5, fmt.Sprintf("clp.it.sinistri@partners."+
		"axa – oggetto dell’email: %s %d –\nUNDERWRITING MEDICO – %s %s", policy.NameDesc,
		policy.ProposalNumber, strings.ToUpper(policy.Contractor.Surname), strings.ToUpper(policy.Contractor.Name)), "LR",
		fpdf.AlignCenter, false)

	setBlackBoldFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3.5, "", "LR", fpdf.AlignCenter, false)
	pdf.MultiCell(0, 3.5, "", "LR", fpdf.AlignCenter, false)
	pdf.MultiCell(0, 3.5, "In caso di capitali assicurati tra i €400.000,00 ed €500.000,00 allegare "+
		"altresì i seguenti esami medici:", "LR", fpdf.AlignLeft, false)
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3.5, "- Analisi del sangue: esame emocromocitometrico - piastrine - Velocità di eritro"+
		" sedimentazione – Glicemia – creatinina – uricemia - colesterolo totale - HDL ("+
		"Lipoproteine ad alta densità) - LDL (Lipoproteine a bassa densità) – trigliceridi - transaminasi GOT/GPT"+
		" - Gammaglutammiltransferasi - anticorpi anti HIV 1 e 2 - sierologia epatite virale B ("+
		"antigeni HB – anti-HBs – anti HBc) - sierologia epatite virale C (anti VHC), Hba1c;", "LR",
		fpdf.AlignLeft, false)
	pdf.MultiCell(0, 3.5, "- Esame cardiovascolare con resoconto medico; ", "LR",
		fpdf.AlignLeft, false)
	pdf.MultiCell(0, 3.5, "- Elettrocardiogramma;", "LR", fpdf.AlignLeft, false)
	pdf.MultiCell(0, 3.5, "- Analisi del PSA (semenogelasi/antigene prostatico specifico) "+
		"esclusivamente per gli uomini la cui età all’Adesione supera i 50 anni.", "LR",
		fpdf.AlignLeft, false)
	pdf.MultiCell(0, 3.5, "", "LBR",
		fpdf.AlignLeft, false)
}
