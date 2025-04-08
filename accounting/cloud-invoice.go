package accounting

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	fattureincloudapi "github.com/fattureincloud/fattureincloud-go-sdk/v2/api"
	fattureincloud "github.com/fattureincloud/fattureincloud-go-sdk/v2/model"
	oauth "github.com/fattureincloud/fattureincloud-go-sdk/v2/oauth2"
)

const (
	baseurlInc = "http://api-v2.fattureincloud.it/"
)

var (
	token = "Bearer " + os.Getenv("FATTURE_INCLOUD_KEY")
)

func (invoiceData *InvoiceInc) CreateInvoice(isPay bool, isProforma bool) {
	log.SetPrefix("CreateInvoice")
	companyId := getCompanyId()
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

	for _, item := range invoiceData.Items {
		fcItems = append(fcItems, *fattureincloud.NewIssuedDocumentItemsListItem().
			SetProductId(4).
			SetCode(item.Code).
			SetName(item.Name).
			SetNetPrice(item.NetPrice).
			SetCategory(item.Category).
			SetDiscount(0).
			SetQty(1).
			SetVat(*fattureincloud.NewVatType().SetId(0)))
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
		SetEntity(entity).
		SetType(Invoicetype).
		SetDate(invoiceData.Date.Format(layout)).
		SetNumber(1).
		SetNumeration("/fatt").
		SetSubject("internal subject").
		SetVisibleSubject("visible subject").
		SetCurrency(*fattureincloud.NewCurrency().SetId("EUR")).
		SetLanguage(*fattureincloud.NewLanguage().SetCode("it").SetName("italiano")).
		SetItemsList(fcItems).
		SetPaymentsList([]fattureincloud.IssuedDocumentPaymentsListItem{
			*fattureincloud.NewIssuedDocumentPaymentsListItem().
				SetAmount(invoiceData.Amount).
				SetDueDate(invoiceData.Date.Format(layout)).
				SetPaidDate(invoiceData.PayDate.Format(layout)).
				SetStatus(status).
				SetPaymentAccount(*fattureincloud.NewPaymentAccount().SetId(110)),
		}).
		// Here we add the payment method
		// List your payment methods: https://github.com/fattureincloud/fattureincloud-go-sdk/blob/master/docs/InfoApi.md#listpaymentmethods
		SetPaymentMethod(*fattureincloud.NewPaymentMethod().SetId(386683))

	// Here we put our invoice in the request object
	createIssuedDocumentRequest := *fattureincloud.NewCreateIssuedDocumentRequest().SetData(invoice)

	uri := baseurlInc + "/c/" + strconv.FormatInt(int64(companyId), 10) + "/issued_documents"
	bodyreq, e := createIssuedDocumentRequest.MarshalJSON()
	log.Println(e)
	req, _ := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(bodyreq))
	req.Header.Add("Authorization", token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error while reading the response bytes:", err)
	}
	log.Println(string([]byte(body)))
}

func getClient() (*fattureincloudapi.APIClient, context.Context, int32) {
	log.SetPrefix("CreateInvoice getClient")
	redirectUri := "http://localhost:3000/oauth"
	auth := oauth.NewOAuth2AuthorizationCodeManager("EZVpwY4saebHSo293egZqSi3I5nyy1fK", os.Getenv("FATTURE_INCLOUD_SECRET"), redirectUri)

	scopes := []oauth.Scope{oauth.Scopes.SETTINGS_ALL, oauth.Scopes.ISSUED_DOCUMENTS_INVOICES_ALL}
	url := auth.GetAuthorizationUrl(scopes, "state")
	log.Println("GetAuthorizationUrl: ", url)
	//params, _ := auth.GetParamsFromUrl(url)
	//code := params.AuthorizationCode
	//state := params.State
	log.Println("os.Getenv(FATTURE_INCLOUD_KEY): ", os.Getenv("FATTURE_INCLOUD_KEY"))
	auth1 := context.WithValue(context.Background(), fattureincloudapi.ContextAccessToken, os.Getenv("FATTURE_INCLOUD_KEY"))
	configuration := fattureincloudapi.NewConfiguration()
	apiClient := fattureincloudapi.NewAPIClient(configuration)
	// Retrieve the first company id
	userCompaniesResponse, _, err := apiClient.UserAPI.ListUserCompanies(auth1).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `UserAPI.ListUserCompanies``: %v\n", err)
	}

	firstCompanyId := userCompaniesResponse.GetData().Companies[0].GetId()
	log.Println("firstCompanyId: ", firstCompanyId)
	return fattureincloudapi.NewAPIClient(configuration), auth1, firstCompanyId
}
func getCompanyId() int32 {
	var (
		listCompany *fattureincloud.ListUserCompaniesResponse
	)
	// for this example we define the token as string, but you should have obtained it in the previous steps

	uri := baseurlInc + "user/companies"
	req, _ := http.NewRequest("GET", uri, nil)
	req.Header.Add("Authorization", token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error while reading the response bytes:", err)
	}
	log.Println(string([]byte(body)))

	e := json.Unmarshal(body, listCompany)
	log.Println(e)

	companyId := listCompany.GetData().Companies[0].Id
	log.Println("companyId:", companyId)
	return int32(*companyId.Get())
}
func putInvoive() int32 {
	var (
		listCompany *fattureincloud.ListUserCompaniesResponse
	)
	// for this example we define the token as string, but you should have obtained it in the previous steps

	uri := baseurlInc + "user/companies"
	req, _ := http.NewRequest("GET", uri, nil)
	req.Header.Add("Authorization", token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error while reading the response bytes:", err)
	}
	log.Println(string([]byte(body)))

	e := json.Unmarshal(body, listCompany)
	log.Println(e)

	companyId := listCompany.GetData().Companies[0].Id
	log.Println("companyId:", companyId)
	return int32(*companyId.Get())
}
