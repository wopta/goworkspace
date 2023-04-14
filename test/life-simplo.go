package test

import (
	"encoding/json"
	"github.com/go-pdf/fpdf"
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

func initFpdf() *fpdf.Fpdf {
	pdf := fpdf.New(fpdf.OrientationPortrait, fpdf.UnitMillimeter, fpdf.PageSizeA4, "")
	pdf.SetMargins(10, 15, 10)
	loadCustomFonts(pdf)
	return pdf
}

func loadCustomFonts(pdf *fpdf.Fpdf) {
	pdf.AddUTF8Font("Montserrat", "", lib.GetAssetPathByEnv("test")+"/montserrat_light.ttf")
	pdf.AddUTF8Font("Montserrat", "B", lib.GetAssetPathByEnv("test")+"/montserrat_bold.ttf")
	pdf.AddUTF8Font("Montserrat", "I", lib.GetAssetPathByEnv("test")+"/montserrat_italic.ttf")
	pdf.AddUTF8Font("Noto", "", lib.GetAssetPathByEnv("test")+"/notosansmono.ttf")
}

func Life(policy models.Policy) string {
	pdf := initFpdf()

	GetMainHeader(pdf, policy)
	GetMainFooter(pdf)

	pdf.AddPage()

	GetContractorInfoSection(pdf, policy.Contractor)

	GetGuaranteesTable(pdf)

	GetAvvertenzeBeneficiariSection(pdf)

	GetBeneficiariSection(pdf)

	GetReferenteTerzoSection(pdf)

	GetStatementsSection(pdf)

	DrawSignatureForm(pdf)

	pdf.AddPage()

	GetVisioneDocumentiSection(pdf)

	GetOfferResumeSection(pdf)

	GetPaymentResumeSection(pdf)

	GetContractWithdrawlSection(pdf)

	GetPaymentMethodSection(pdf)

	GetEmitResumeSection(pdf)

	GetPolicyDescriptionSection(pdf)

	GetWoptaAxaCompanyDescriptionSection(pdf)

	GetAxaHeader(pdf)

	pdf.AddPage()

	GetAxaFooter(pdf)

	GetAxaDeclarationsConsentSection(pdf)

	pdf.AddPage()

	GetAxaTableSection(pdf)

	pdf.AddPage()

	GetAxaTablePart2Section(pdf)

	pdf.Ln(15)

	pdf.SetFont("Montserrat", "B", 10)

	GetAxaTablePart3Section(pdf)

	GetWoptaHeader(pdf)

	pdf.AddPage()

	GetWoptaFooter(pdf)

	GetAllegato4Section(pdf)

	pdf.AddPage()

	GetAllegato4TerSection(pdf)

	pdf.AddPage()

	GetWoptaPrivacySection(pdf)

	pdf.AddPage()

	GetPersonalDataHandlingSection(pdf)

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

	filename, err := Save(pdf)
	lib.CheckError(err)
	log.Println(filename + " wrote successfully")
	/*err := pdf.OutputFileAndClose("test/test.pdf")
	log.Println(err)*/
	return filename
}
