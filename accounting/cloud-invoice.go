package accounting

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	fattureincloud "github.com/fattureincloud/fattureincloud-go-sdk/v2/model"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/models"
)

const (
	baseurlInc = "https://api-v2.fattureincloud.it/"
)

func (invoiceData InvoiceInc) create(isPay bool, isProforma bool) (string, error) {
	companyId, err := getCompanyId()
	if err != nil {
		return "", err
	}
	token := "Bearer " + os.Getenv("FATTURE_INCLOUD_KEY")
	var (
		fcItems     []fattureincloud.IssuedDocumentItemsListItem
		status      fattureincloud.IssuedDocumentStatus
		Invoicetype fattureincloud.IssuedDocumentType
	)
	const (
		layout = "2006-01-02"
	)
	if isProforma {
		Invoicetype = fattureincloud.IssuedDocumentTypes.PROFORMA
	} else {
		Invoicetype = fattureincloud.IssuedDocumentTypes.INVOICE
	}
	if isPay {
		status = fattureincloud.IssuedDocumentStatuses.PAID
	} else {
		status = fattureincloud.IssuedDocumentStatuses.NOT_PAID
	}
	//set your company id
	if len(invoiceData.Items) == 0 {
		return "", errors.New("Error No invoice items founded")
	}
	for _, item := range invoiceData.Items {
		fcItems = append(fcItems, *fattureincloud.NewIssuedDocumentItemsListItem().
			//SetProductId(4).
			SetDescription(item.Desc).
			SetCode(item.Code).
			SetName(item.Name).
			SetGrossPrice(item.GrossPrice).
			SetCategory(item.Category).
			SetDiscount(0).
			SetQty(float32(item.Qty)).
			SetVat(*fattureincloud.NewVatType().SetId(8199995)))
	}
	entity := *fattureincloud.NewEntity().
		SetId(1).
		SetName(invoiceData.Name).
		SetVatNumber(invoiceData.VatNumber).
		SetTaxCode(invoiceData.TaxCode).
		SetAddressStreet(invoiceData.Address).
		SetAddressPostalCode(invoiceData.PostalCode).
		SetAddressCity(invoiceData.City).
		SetAddressProvince(invoiceData.CityCode).
		SetCountry(invoiceData.Country)

	invoice := *fattureincloud.NewIssuedDocument().
		SetUseGrossPrices(true).
		SetEntity(entity).
		SetType(Invoicetype).
		SetDate(invoiceData.Date.Format(layout)).
		//SetNumber(invoiceData.Qty).
		//SetNumeration("/fatt").
		//SetSubject("internal subject").
		//SetVisibleSubject("visible subject").
		SetCurrency(*fattureincloud.NewCurrency().SetId("EUR")).
		SetLanguage(*fattureincloud.NewLanguage().SetCode("it").SetName("italiano")).
		SetItemsList(fcItems).
		SetEInvoice(true).
		SetEiData(*fattureincloud.NewIssuedDocumentEiData().SetPaymentMethod(invoiceData.PaymentMethod)).
		SetPaymentsList([]fattureincloud.IssuedDocumentPaymentsListItem{
			*fattureincloud.NewIssuedDocumentPaymentsListItem().
				SetAmount(invoiceData.Amount).
				SetDueDate(invoiceData.Date.Format(layout)).
				SetPaidDate(invoiceData.PayDate.Format(layout)).
				SetStatus(status).
				SetPaymentAccount(*fattureincloud.NewPaymentAccount().SetId(1405562)),
		})
		// Here we add the payment method
		// List your payment methods: https://github.com/fattureincloud/fattureincloud-go-sdk/blob/master/docs/InfoApi.md#listpaymentmethods
		//SetPaymentMethod(*fattureincloud.NewPaymentMethod().SetId(386683))

	// Here we put our invoice in the request object
	createIssuedDocumentRequest := *fattureincloud.NewCreateIssuedDocumentRequest().SetData(invoice)

	uri := baseurlInc + "c/" + strconv.FormatInt(int64(companyId), 10) + "/issued_documents"
	bodyreq, err := createIssuedDocumentRequest.MarshalJSON()
	if err != nil {
		return "", err
	}
	req, _ := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(bodyreq))
	req.Header.Add("Authorization", token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Error on response %v ", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Error while reading the response bytes: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error response, expected:200, got: %v, message: %v", resp.Status, string(body))
	}
	var invresp InvoiceResponse
	err = json.Unmarshal([]byte(body), &invresp)
	if err != nil {
		return "", fmt.Errorf("Error unmarshalling: %v", err)
	}

	return invresp.Data.URL, nil
}

func getCompanyId() (int32, error) {
	var (
		listCompany CompanyResponseData
	)
	// for this example we define the token as string, but you should have obtained it in the previous steps

	token := "Bearer " + os.Getenv("FATTURE_INCLOUD_KEY")
	uri := baseurlInc + "user/companies"
	req, _ := http.NewRequest("GET", uri, nil)
	req.Header.Add("Authorization", token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("Error on response %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("Error while reading the response bytes %v", err)
	}
	log.Println(string([]byte(body)))

	e := json.Unmarshal(body, &listCompany)
	if e != nil {
		return 0, e
	}
	if len(listCompany.Data.Companies) == 0 {
		return 0, errors.New("Error no companies founded")
	}
	companyId := listCompany.Data.Companies[0].ID
	return int32(companyId), nil
}
func (invoiceData InvoiceInc) save(url string, path string) (string, error) {
	out, err := HttpFileToByte(url)
	if err != nil {
		return "", err
	}
	return lib.PutToStorageErr(os.Getenv("GOOGLE_STORAGE_BUCKET"), path, out.Bytes())
}

type CompanyResponseData struct {
	Data Data `json:"data,omitempty"`
}
type Permissions struct {
	FicSituation         string `json:"fic_situation,omitempty"`
	FicClients           string `json:"fic_clients,omitempty"`
	FicSuppliers         string `json:"fic_suppliers,omitempty"`
	FicProducts          string `json:"fic_products,omitempty"`
	FicIssuedDocuments   string `json:"fic_issued_documents,omitempty"`
	FicReceivedDocuments string `json:"fic_received_documents,omitempty"`
	FicReceipts          string `json:"fic_receipts,omitempty"`
	FicCalendar          string `json:"fic_calendar,omitempty"`
	FicArchive           string `json:"fic_archive,omitempty"`
	FicTaxes             string `json:"fic_taxes,omitempty"`
	FicStock             string `json:"fic_stock,omitempty"`
	FicCashbook          string `json:"fic_cashbook,omitempty"`
	FicSettings          string `json:"fic_settings,omitempty"`
	FicEmails            string `json:"fic_emails,omitempty"`
	DicEmployees         string `json:"dic_employees,omitempty"`
	DicTimesheet         string `json:"dic_timesheet,omitempty"`
	DicSettings          string `json:"dic_settings,omitempty"`
}
type Companies struct {
	ID                  int         `json:"id,omitempty"`
	Name                string      `json:"name,omitempty"`
	Email               string      `json:"email,omitempty"`
	Alias               any         `json:"alias,omitempty"`
	VatNumber           string      `json:"vat_number,omitempty"`
	TaxCode             string      `json:"tax_code,omitempty"`
	Type                string      `json:"type,omitempty"`
	ConnectionID        int         `json:"connection_id,omitempty"`
	ConnectionRole      string      `json:"connection_role,omitempty"`
	ControlledCompanies []any       `json:"controlled_companies,omitempty"`
	FicPlan             string      `json:"fic_plan,omitempty"`
	DicPlan             int         `json:"dic_plan,omitempty"`
	Fic                 bool        `json:"fic,omitempty"`
	Dic                 bool        `json:"dic,omitempty"`
	FicLicenseExpire    string      `json:"fic_license_expire,omitempty"`
	Permissions         Permissions `json:"permissions,omitempty"`
}
type InvoiceResponse struct {
	Data struct {
		ID                               int     `json:"id,omitempty"`
		Type                             string  `json:"type,omitempty"`
		Year                             int     `json:"year,omitempty"`
		Numeration                       string  `json:"numeration,omitempty"`
		Subject                          string  `json:"subject,omitempty"`
		VisibleSubject                   string  `json:"visible_subject,omitempty"`
		RcCenter                         string  `json:"rc_center,omitempty"`
		AmountRivalsa                    float64 `json:"amount_rivalsa,omitempty"`
		AmountRivalsaTaxable             float64 `json:"amount_rivalsa_taxable,omitempty"`
		AmountGlobalCassaTaxable         float64 `json:"amount_global_cassa_taxable,omitempty"`
		AmountCassa                      float64 `json:"amount_cassa,omitempty"`
		AmountCassaTaxable               float64 `json:"amount_cassa_taxable,omitempty"`
		AmountCassa2                     float64 `json:"amount_cassa2,omitempty"`
		AmountCassa2Taxable              float64 `json:"amount_cassa2_taxable,omitempty"`
		AmountWithholdingTax             float64 `json:"amount_withholding_tax,omitempty"`
		AmountWithholdingTaxTaxable      float64 `json:"amount_withholding_tax_taxable,omitempty"`
		AmountOtherWithholdingTax        float64 `json:"amount_other_withholding_tax,omitempty"`
		AmountEnasarcoTaxable            float64 `json:"amount_enasarco_taxable,omitempty"`
		AmountOtherWithholdingTaxTaxable float64 `json:"amount_other_withholding_tax_taxable,omitempty"`
		EiCassaType                      any     `json:"ei_cassa_type,omitempty"`
		EiCassa2Type                     any     `json:"ei_cassa2_type,omitempty"`
		EiWithholdingTaxCausal           any     `json:"ei_withholding_tax_causal,omitempty"`
		EiOtherWithholdingTaxType        any     `json:"ei_other_withholding_tax_type,omitempty"`
		EiOtherWithholdingTaxCausal      any     `json:"ei_other_withholding_tax_causal,omitempty"`
		StampDuty                        float64 `json:"stamp_duty,omitempty"`
		UseGrossPrices                   bool    `json:"use_gross_prices,omitempty"`
		EInvoice                         bool    `json:"e_invoice,omitempty"`
		AgyoCompanyID                    any     `json:"agyo_company_id,omitempty"`
		AgyoID                           any     `json:"agyo_id,omitempty"`
		AgyoSentAt                       any     `json:"agyo_sent_at,omitempty"`
		DeliveryNote                     bool    `json:"delivery_note,omitempty"`
		AccompanyingInvoice              bool    `json:"accompanying_invoice,omitempty"`
		AmountNet                        float64 `json:"amount_net,omitempty"`
		AmountVat                        float64 `json:"amount_vat,omitempty"`
		AmountGross                      float64 `json:"amount_gross,omitempty"`
		AmountDueDiscount                float64 `json:"amount_due_discount,omitempty"`
		PermanentToken                   string  `json:"permanent_token,omitempty"`
		HMargins                         float64 `json:"h_margins,omitempty"`
		VMargins                         float64 `json:"v_margins,omitempty"`
		ShowPaymentMethod                bool    `json:"show_payment_method,omitempty"`
		ShowPayments                     bool    `json:"show_payments,omitempty"`
		ShowTotals                       string  `json:"show_totals,omitempty"`
		ShowNotificationButton           bool    `json:"show_notification_button,omitempty"`
		IsMarked                         bool    `json:"is_marked,omitempty"`
		CreatedAt                        string  `json:"created_at,omitempty"`
		UpdatedAt                        string  `json:"updated_at,omitempty"`
		AttachPdfToXML                   bool    `json:"attach_pdf_to_xml,omitempty"`
		PriceListID                      any     `json:"price_list_id,omitempty"`
		Entity                           struct {
			Name              string `json:"name,omitempty"`
			VatNumber         string `json:"vat_number,omitempty"`
			TaxCode           string `json:"tax_code,omitempty"`
			AddressStreet     string `json:"address_street,omitempty"`
			AddressPostalCode string `json:"address_postal_code,omitempty"`
			AddressCity       string `json:"address_city,omitempty"`
			AddressProvince   string `json:"address_province,omitempty"`
			AddressExtra      string `json:"address_extra,omitempty"`
			Country           string `json:"country,omitempty"`
			CertifiedEmail    string `json:"certified_email,omitempty"`
			EiCode            string `json:"ei_code,omitempty"`
			EntityType        string `json:"entity_type,omitempty"`
			Type              any    `json:"type,omitempty"`
		} `json:"entity,omitempty"`
		Date     string `json:"date,omitempty"`
		Number   int    `json:"number,omitempty"`
		Currency struct {
			ID           string `json:"id,omitempty"`
			ExchangeRate string `json:"exchange_rate,omitempty"`
			Symbol       string `json:"symbol,omitempty"`
		} `json:"currency,omitempty"`
		Language struct {
			Code string `json:"code,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"language,omitempty"`
		Notes                      string  `json:"notes,omitempty"`
		Rivalsa                    float64 `json:"rivalsa,omitempty"`
		RivalsaTaxable             float64 `json:"rivalsa_taxable,omitempty"`
		GlobalCassaTaxable         float64 `json:"global_cassa_taxable,omitempty"`
		Cassa                      float64 `json:"cassa,omitempty"`
		CassaTaxable               float64 `json:"cassa_taxable,omitempty"`
		Cassa2                     float64 `json:"cassa2,omitempty"`
		Cassa2Taxable              float64 `json:"cassa2_taxable,omitempty"`
		WithholdingTax             float64 `json:"withholding_tax,omitempty"`
		WithholdingTaxTaxable      float64 `json:"withholding_tax_taxable,omitempty"`
		OtherWithholdingTax        float64 `json:"other_withholding_tax,omitempty"`
		OtherWithholdingTaxTaxable float64 `json:"other_withholding_tax_taxable,omitempty"`
		PaymentMethod              struct {
			ID      any    `json:"id,omitempty"`
			Name    string `json:"name,omitempty"`
			Details []struct {
				Title       string `json:"title,omitempty"`
				Description string `json:"description,omitempty"`
			} `json:"details,omitempty"`
		} `json:"payment_method,omitempty"`
		UseSplitPayment  bool `json:"use_split_payment,omitempty"`
		MergedIn         any  `json:"merged_in,omitempty"`
		OriginalDocument any  `json:"original_document,omitempty"`
		ItemsList        []struct {
			ProductID             any     `json:"product_id,omitempty"`
			Code                  string  `json:"code,omitempty"`
			Name                  string  `json:"name,omitempty"`
			Measure               string  `json:"measure,omitempty"`
			Category              string  `json:"category,omitempty"`
			ID                    int     `json:"id,omitempty"`
			ApplyWithholdingTaxes bool    `json:"apply_withholding_taxes,omitempty"`
			Discount              float64 `json:"discount,omitempty"`
			DiscountHighlight     bool    `json:"discount_highlight,omitempty"`
			InDn                  bool    `json:"in_dn,omitempty"`
			Qty                   int     `json:"qty,omitempty"`
			NetPrice              float64 `json:"net_price,omitempty"`
			Vat                   struct {
				ID          int     `json:"id,omitempty"`
				Value       float64 `json:"value,omitempty"`
				Description string  `json:"description,omitempty"`
			} `json:"vat,omitempty"`
			Stock       bool    `json:"stock,omitempty"`
			Description string  `json:"description,omitempty"`
			GrossPrice  float64 `json:"gross_price,omitempty"`
			NotTaxable  bool    `json:"not_taxable,omitempty"`
		} `json:"items_list,omitempty"`
		PaymentsList  []any `json:"payments_list,omitempty"`
		AttachmentURL any   `json:"attachment_url,omitempty"`
		SeenDate      any   `json:"seen_date,omitempty"`
		NextDueDate   any   `json:"next_due_date,omitempty"`
		Template      struct {
			ID   int    `json:"id,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"template,omitempty"`
		ExtraData              any    `json:"extra_data,omitempty"`
		URL                    string `json:"url,omitempty"`
		Locked                 bool   `json:"locked,omitempty"`
		EiLocked               bool   `json:"ei_locked,omitempty"`
		HasTsPayPendingPayment bool   `json:"has_ts_pay_pending_payment,omitempty"`
		ShowTspayButton        bool   `json:"show_tspay_button,omitempty"`
		PayWithTspayURL        any    `json:"pay_with_tspay_url,omitempty"`
		HasAttachment          bool   `json:"has_attachment,omitempty"`
	} `json:"data,omitempty"`
}
type Data struct {
	Companies []Companies `json:"companies,omitempty"`
}

var paymentMethodMap map[string]string = map[string]string{
	// codes from https://www.fatturapa.gov.it/export/documenti/fatturapa/v1.3/Rappresentazione-tabellare-fattura-ordinaria.pdf
	"creditcard": "MP08",
	"transfer":   "MP05",
	"sdd":        "MP20",
	"remittance": "MP05",
}

func MapPolicyInvoiceInc(policy models.Policy, tr models.Transaction, desc string) InvoiceInc {
	var amount float32
	for _, item := range tr.Items {
		if item.Type == "consultancy" {
			amount = float32(item.AmountGross)
		}
	}
	inv := InvoiceInc{
		Name:          policy.Contractor.Name + " " + policy.Contractor.Surname,
		VatNumber:     policy.Contractor.VatCode,
		TaxCode:       policy.Contractor.FiscalCode,
		Address:       policy.Contractor.Residence.StreetName + " " + policy.Contractor.Residence.StreetNumber,
		PostalCode:    policy.Contractor.Residence.PostalCode,
		City:          policy.Contractor.Residence.City,
		CityCode:      policy.Contractor.Residence.CityCode,
		Country:       "Italia",
		Mail:          policy.Contractor.Mail,
		Amount:        amount,
		Date:          time.Now().UTC(),
		PayDate:       tr.PayDate,
		PaymentMethod: paymentMethodMap[tr.PaymentMethod],
		Items: []Items{{
			Desc:       desc,
			Name:       policy.Name,
			Code:       policy.Name,
			Qty:        1,
			ProductId:  0,
			GrossPrice: amount,
			Category:   policy.Name,
			Date:       tr.CreationDate}}}
	return inv

}
