package utility

import (
	"errors"
	"fmt"
	"strings"

	"gitlab.dev.wopta.it/goworkspace/document"
	"gitlab.dev.wopta.it/goworkspace/document/namirial"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/payment"
	"gitlab.dev.wopta.it/goworkspace/payment/consultancy"
	"gitlab.dev.wopta.it/goworkspace/transaction"
)

func SignFiles(policy *models.Policy, product *models.Product, networkNode *models.NetworkNode, sendEmail bool, origin string) error {
	log.AddPrefix("emitSign")
	defer log.PopPrefix()
	log.Printf("Policy Uid %s", policy.Uid)

	policy.IsSign = false
	policy.Status = models.PolicyStatusToSign
	policy.StatusHistory = append(policy.StatusHistory, models.PolicyStatusContact, models.PolicyStatusToSign)

	namirialInput := namirial.NamirialInput{
		Policy:            *policy,
		DocumentsFullPath: make([]string, 0),
		SendEmail:         sendEmail,
		Origin:            origin,
	}
	//TODO: remove this 'if' after catnat document is done
	if policy.Name != models.CatNatProduct {
		p := <-document.ContractObj(origin, *policy, networkNode, product)
		document, err := p.SaveWithName(models.ContractAttachmentName)
		if err != nil {
			return err
		}
		policy.DocumentName = document.LinkGcs
		namirialInput.DocumentsFullPath = append(namirialInput.DocumentsFullPath, document.LinkGcs)
	}
	//Preparing dto for namirial
	basePathForDocument := strings.ReplaceAll(fmt.Sprintf("temp/%s/namirial/", policy.Uid), " ", "_")
	fullPathDocumentToSign, err := lib.ListGoogleStorageFolderContent(basePathForDocument)
	if err != nil {
		return err
	}
	for _, path := range fullPathDocumentToSign {
		namirialInput.DocumentsFullPath = append(namirialInput.DocumentsFullPath, path)
	}

	if len(namirialInput.DocumentsFullPath) == 0 {
		log.ErrorF("nothing to sign")
		return nil
	}

	envelope, err := namirial.Sign(namirialInput)
	if err != nil {
		return err
	}
	policy.IdSign = envelope.IdEnvelope
	policy.SignUrl = envelope.Url
	policy.ContractFileId = envelope.FileIds[0] //this field is deprecated
	policy.DocumentName = basePathForDocument   //this field is deprecated
	return nil
}

func EmitPay(policy *models.Policy, origin string, productP, mgaProductP *models.Product, networkNode *models.NetworkNode) {
	log.AddPrefix("emitPay")
	defer log.PopPrefix()
	log.Printf("Policy Uid %s", policy.Uid)

	policy.IsPay = false
	payUrl, err := CreatePolicyTransactions(policy, productP, mgaProductP, networkNode)
	if err != nil {
		return
	}
	policy.PayUrl = payUrl
}

func CreatePolicyTransactions(policy *models.Policy, product *models.Product, mgaProduct *models.Product, networkNode *models.NetworkNode) (string, error) {
	transactions := transaction.CreateTransactions(*policy, *mgaProduct, func() string { return lib.NewDoc(models.TransactionsCollection) })
	if len(transactions) == 0 {
		log.Println("no transactions created")
		return "", errors.New("no transactions created")
	}

	client := payment.NewClient(policy.Payment, *policy, *product, transactions, false, "")
	payUrl, updatedTransactions, err := client.NewBusiness()
	if err != nil {
		log.ErrorF("error emitPay policy %s: %s", policy.Uid, err.Error())
		return "", err
	}

	for index, tr := range updatedTransactions {
		err = lib.SetFirestoreErr(models.TransactionsCollection, tr.Uid, tr)
		if err != nil {
			log.ErrorF("error saving transaction %s to firestore: %s", tr.Uid, err.Error())
			return "", err
		}
		tr.BigQuerySave("")

		if tr.IsPay {
			err = transaction.CreateNetworkTransactions(policy, &updatedTransactions[index], networkNode, mgaProduct)
			if err != nil {
				log.ErrorF("error creating network transactions: %s", err.Error())
				return "", err
			}
			if err := consultancy.GenerateInvoice(*policy, tr); err != nil {
				log.Printf("error handling consultancy: %s", err.Error())
			}
		}
	}
	return payUrl, err
}

func SetAdvance(policy *models.Policy, origin string, product *models.Product, mgaProduct *models.Product, networkNode *models.NetworkNode, paymentSplit string, paymentMode string) {
	policy.Payment = models.ManualPaymentProvider
	policy.IsPay = true
	policy.Status = models.PolicyStatusPay
	policy.StatusHistory = append(policy.StatusHistory, models.PolicyStatusToPay, models.PolicyStatusPay)

	//TODO: fix me someday in the future
	if paymentSplit != "" && policy.PaymentSplit == "" {
		policy.PaymentSplit = paymentSplit
	}
	if paymentMode != "" && policy.PaymentMode == "" {
		policy.PaymentMode = paymentMode
	}
	CreatePolicyTransactions(policy, product, mgaProduct, networkNode)
}
