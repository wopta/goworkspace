package document

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-pdf/fpdf"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/models"
)

func lifeReserved(policy *models.Policy, product *models.Product) (DocumentGenerated, error) {
	log.AddPrefix("lifeReserved")
	defer log.PopPrefix()
	log.Println("function start ------------------------------")

	pdf := initFpdf()

	lifeReservedFooter(pdf)

	pdf.AddPage()

	lifeReservedHeader(pdf, policy)

	lifeReservedInsuredInfoSection(pdf, policy)

	guaranteesMap, slugs := loadLifeGuarantees(policy, product)

	lifeGuaranteesTableV2(pdf, guaranteesMap, slugs)

	lifeReservedInsuranceLimitSection(pdf)

	lifeReservedInstructionsSection(pdf, policy)

	log.Println("function end --------------------------------")

	return generateReservedDocument(pdf, policy)
}

func lifeReservedHeader(pdf *fpdf.Fpdf, policy *models.Policy) {
	var (
		opt                  fpdf.ImageOptions
		logoPath, expiryInfo string
	)

	logoPath = lib.GetAssetPathByEnvV2() + "logo_vita.png"

	location, err := time.LoadLocation("Europe/Rome")
	lib.CheckError(err)

	policyStartDate := policy.StartDate.In(location)
	policyEndDate := policy.EndDate.In(location)

	if policy.PaymentSplit == string(models.PaySplitMonthly) {
		expiryInfo = "Prima scandenza mensile il: " +
			policyStartDate.AddDate(0, 1, 0).Format(dateLayout)
	} else if policy.PaymentSplit == string(models.PaySplitYear) || policy.PaymentSplit == string(models.PaySplitYearly) {
		expiryInfo = "Prima scadenza annuale il: " +
			policyStartDate.AddDate(1, 0, 0).Format(dateLayout)
	}

	proposalInfo := fmt.Sprintf("Numero: %d\nDecorrere dal: %s\nScade il: %s\n%s\nNon si rinnova a scadenza",
		policy.ProposalNumber, policyStartDate.Format(dateLayout), policyEndDate.Format(dateLayout), expiryInfo)

	contractor := policy.Contractor
	address := strings.ToUpper(contractor.Residence.StreetName + ", " + contractor.Residence.StreetNumber + "\n" +
		contractor.Residence.PostalCode + " " + contractor.Residence.City + " (" + contractor.Residence.CityCode + ")")

	if policy.PaymentSplit == string(models.PaySplitMonthly) {
		expiryInfo = "Prima scandenza mensile il: " +
			policyStartDate.AddDate(0, 1, 0).Format(dateLayout) + "\n"
	} else if policy.PaymentSplit == string(models.PaySplitYear) {
		expiryInfo = "Prima scadenza annuale il: " +
			policyStartDate.AddDate(1, 0, 0).Format(dateLayout) + "\n"
	}

	contractorInfo := fmt.Sprintf("Contraente: %s\nC.F./P.IVA: %s\nIndirizzo: %s\nMail: %s\nTelefono: %s",
		strings.ToUpper(contractor.Surname+" "+contractor.
			Name), contractor.FiscalCode, strings.ToUpper(address), contractor.Mail, contractor.Phone)

	opt.ImageType = "png"
	pdf.ImageOptions(lib.GetAssetPathByEnvV2()+"logo_axa.png", 180, 10, 0, 15,
		false, opt, 0, "")
	pdf.SetY(pdf.GetY() + 18)

	setBlackBoldFont(pdf, 20)
	pdf.MultiCell(0, 3, "SCHEDA RAPPORTO VISITA MEDICA", "", fpdf.AlignCenter, false)
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
	pdf.Ln(6)
}

func lifeReservedFooter(pdf *fpdf.Fpdf) {
	pdf.SetFooterFunc(func() {
		pdf.SetY(-8)
		pageNumber(pdf)
	})
}

func lifeReservedInsuredInfoSection(pdf *fpdf.Fpdf, policy *models.Policy) {
	lifeInsuredInfoSectionV2(pdf, policy, true)
}

func lifeReservedInsuranceLimitSection(pdf *fpdf.Fpdf) {
	text := "Limiti assuntivi:\n\n" +
		"Decesso: 75 anni a scadenza - max 500.000 euro\n" +
		"Invalidità Totale Permanente: 75 anni a scadenza - max 500.000 euro\n" +
		"Inabilità Totale Temporanea: 75 anni a scadenza  - max 3.000\n" +
		"Malattie Gravi: 65 anni a scadenza - max 100.000 euro"

	setBlackBoldFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3.5, text, "", fpdf.AlignLeft, false)
	pdf.Ln(5)
}

func lifeReservedInstructionsSection(pdf *fpdf.Fpdf, policy *models.Policy) {
	setBlackDrawColor(pdf)
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3.5, "", "LTR", fpdf.AlignCenter, false)
	pdf.MultiCell(0, 3.5, "Da restituire alla compagnia assicurativa, unitamente al Rapporto Visita "+
		"Medica compilato e firmato da te e da un medico", "LR", fpdf.AlignCenter, false)

	for _, contact := range policy.ReservedInfo.Contacts {
		setBlackBoldFont(pdf, standardTextSize)
		pdf.MultiCell(0, 3.5, "", "LR", fpdf.AlignCenter, false)
		pdf.MultiCell(0, 3.5, "", "LR", fpdf.AlignCenter, false)
		pdf.MultiCell(0, 3.5, contact.Title, "LR", fpdf.AlignCenter, false)
		setBlackRegularFont(pdf, standardTextSize)
		pdf.MultiCell(0, 3.5, contact.Address, "LR", fpdf.AlignCenter, false)
		pdf.MultiCell(0, 3.5, contact.Subject,
			"LR", fpdf.AlignCenter,
			false)
	}

	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3.5, "", "LR", fpdf.AlignCenter, false)
	pdf.MultiCell(0, 3.5, "", "LR", fpdf.AlignCenter, false)
	pdf.MultiCell(0, 3.5, fmt.Sprintf("1 - %s", policy.ReservedInfo.RequiredExams[0]), "LR", fpdf.AlignLeft, false)
	pdf.MultiCell(0, 3.5, "", "LR", fpdf.AlignCenter, false)

	if len(policy.ReservedInfo.RequiredExams) > 1 {
		setBlackBoldFont(pdf, standardTextSize)
		pdf.MultiCell(0, 3.5, "In caso di capitali assicurati tra i €400.000,00 ed €500.000,00 allegare "+
			"altresì i seguenti esami medici:", "LR", fpdf.AlignLeft, false)

		for index, medicalDocument := range policy.ReservedInfo.RequiredExams[1:] {
			setBlackRegularFont(pdf, standardTextSize)
			pdf.MultiCell(0, 3.5, fmt.Sprintf("%d - %s", index+2, medicalDocument), "LR", fpdf.AlignLeft, false)
		}
	}

	pdf.MultiCell(0, 3.5, "", "LBR", fpdf.AlignLeft, false)
}
