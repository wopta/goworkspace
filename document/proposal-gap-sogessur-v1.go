package document

import (
	"github.com/go-pdf/fpdf"
	"github.com/wopta/goworkspace/models"
)

func gapSogessurProposalV1(pdf *fpdf.Fpdf, origin string, policy *models.Policy, networkNode *models.NetworkNode) (string, []byte) {
	gapHeaderV1(pdf, policy, networkNode, true)

	gapFooterV1(pdf, policy.NameDesc)

	pdf.AddPage()

	vehicle := policy.Assets[0].Vehicle
	contractor := policy.Contractor
	vehicleOwner := policy.Assets[0].Person
	statements := *policy.Statements

	pdf.SetTextColor(229, 0, 117)
	pdf.SetFont("Montserrat", "B", 9)
	pdf.MultiCell(0, 4, "La tua assicurazione è operante sui dati sotto riportati, verifica la loro correttezza"+
		" e segnala eventuali inesattezze", "", "", false)
	pdf.Ln(3)

	gapVehicleDataTableV1(pdf, vehicle)

	gapPersonalInfoTableV1(pdf, contractor, *vehicleOwner)

	gapPolicyDataTableV1(pdf, policy)

	gapPriceTableV1(pdf, policy)

	pdf.Ln(3)

	gapStatementsV1(pdf, statements[:len(statements)-1], policy.Company, true)

	companiesDescriptionSection(pdf, policy.Company)

	woptaGapHeader(pdf, *policy, true)

	pdf.AddPage()

	woptaFooter(pdf)

	printStatement(pdf, statements[len(statements)-1], policy.Company, true)

	woptaHeader(pdf, true)

	generatePolicyAnnex(pdf, origin, networkNode)

	pdf.AddPage()

	woptaPrivacySection(pdf)

	personalDataHandlingSection(pdf, policy, true)

	filename, out := saveProposal(pdf, policy)
	return filename, out
}
