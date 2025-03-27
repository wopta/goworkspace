package accounting

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	fattureincloudapi "github.com/fattureincloud/fattureincloud-go-sdk/v2/api"
	fattureincloud "github.com/fattureincloud/fattureincloud-go-sdk/v2/model"
	oauth "github.com/fattureincloud/fattureincloud-go-sdk/v2/oauth2"
)

func (invoiceData *InvoiceInc) CreateInvoice(isPay bool, isProforma bool) {
	log.SetPrefix("CreateInvoice")
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

	apiClient, auth, id := getClient()
	//set your company id
	companyId := id
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
	apiClient.UserAPI.ListUserCompanies(auth)
	// Now we are all set for the final call
	// Create the invoice: https://github.com/fattureincloud/fattureincloud-go-sdk/blob/master/docs/IssuedDocumentsApi.md#createIssuedDocument
	resp, r, err := apiClient.IssuedDocumentsAPI.CreateIssuedDocument(auth, companyId).CreateIssuedDocumentRequest(createIssuedDocumentRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `IssuedDocumentsAPI.CreateIssuedDocument``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	json.NewEncoder(os.Stdout).Encode(resp)
}

func getClient() (*fattureincloudapi.APIClient, context.Context, int32) {
	log.SetPrefix("CreateInvoice getClient")
	redirectUri := "http://localhost:3000/oauth"
	auth := oauth.NewOAuth2AuthorizationCodeManager("EZVpwY4saebHSo293egZqSi3I5nyy1fK", os.Getenv("FATTURE_INCLOUD_KEY"), redirectUri)
	oauth.NewOAuth2AuthorizationCodeParams()

	scopes := []oauth.Scope{oauth.Scopes.SETTINGS_ALL, oauth.Scopes.ISSUED_DOCUMENTS_INVOICES_ALL}
	url := auth.GetAuthorizationUrl(scopes, "state")
	log.Println("GetAuthorizationUrl: ", url)
	params, _ := auth.GetParamsFromUrl(url)
	code := params.AuthorizationCode
	//state := params.State
	auth1 := context.WithValue(context.Background(), fattureincloudapi.ContextAccessToken, code)
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
