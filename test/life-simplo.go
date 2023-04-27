package test

import (
	"encoding/json"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"io"
	"log"
	"net/http"
)

type Statements struct {
	Statements []*models.Statement `json:"statements"`
	Text       string              `json:"text,omitempty"`
}

func LifeSimploHandler(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		policy models.Policy
	)

	b := lib.ErrorByte(io.ReadAll(r.Body))
	err := json.Unmarshal(b, &policy)
	lib.CheckError(err)

	filename := Life(policy)
	log.Println(filename)

	return "", nil, nil
}

func Life(policy models.Policy) string {
	pdf := initFpdf()

	GetMainHeader(pdf, policy)
	GetMainFooter(pdf)

	pdf.AddPage()

	GetContractorInfoSection(pdf, policy)

	GetGuaranteesTable(pdf, policy)

	GetAvvertenzeBeneficiariSection(pdf)

	GetBeneficiariSection(pdf, policy)

	GetBeneficiaryReferenceSection(pdf, policy)

	GetSurveysSection(pdf, policy)

	drawSignatureForm(pdf)
	pdf.Ln(5)

	GetStatementsSection(pdf, policy)

	pdf.AddPage()

	//GetVisioneDocumentiSection(pdf, policy)

	GetOfferResumeSection(pdf, policy)

	GetPaymentResumeSection(pdf, policy)

	GetContractWithdrawlSection(pdf)

	GetPaymentMethodSection(pdf)

	GetEmitResumeSection(pdf, policy)

	//GetPolicyDescriptionSection(pdf)

	GetWoptaAxaCompanyDescriptionSection(pdf)

	GetAxaHeader(pdf)

	pdf.AddPage()

	GetAxaFooter(pdf)

	GetAxaDeclarationsConsentSection(pdf, policy)

	pdf.AddPage()

	GetAxaTableSection(pdf, policy)

	pdf.AddPage()

	GetAxaTablePart2Section(pdf, policy)

	pdf.Ln(15)

	GetAxaTablePart3Section(pdf)

	GetWoptaHeader(pdf)

	pdf.AddPage()

	GetWoptaFooter(pdf)

	GetAllegato3Section(pdf)

	pdf.AddPage()

	GetAllegato4Section(pdf)

	pdf.AddPage()

	GetAllegato4TerSection(pdf)

	pdf.AddPage()

	GetWoptaPrivacySection(pdf)

	GetPersonalDataHandlingSection(pdf, policy)

	/*tpl := new(fpdf.FpdfTpl)

	fb, err := os.ReadFile("document/assets/template.pdf")
	if err != nil {
		return "", nil, err
	}

	err = tpl.GobDecode(fb)
	if err != nil {
		return "", nil, err
	}

	template, _ := tpl.FromPage(1)

	pdf.UseTemplate(template)*/

	filename, err := save(pdf)
	lib.CheckError(err)
	log.Println(filename + " wrote successfully")
	/*err := pdf.OutputFileAndClose("test/test.pdf")
	log.Println(err)*/
	return filename
}
