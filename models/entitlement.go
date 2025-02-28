package models

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/wopta/goworkspace/lib"
	"google.golang.org/api/iterator"
)

// TODO: run add column script on bigquery table for networkNodes
type Entitlement struct {
	Slug string `json:"slug" firestore:"slug" bigquery:"slug"`
}

type EntitlementProfile struct {
	Slug         string        `json:"slug" firestore:"slug" bigquery:"slug"`
	Entitlements []Entitlement `json:"entilements" firestore:"entilements" bigquery:"entilements"`
	Data         string        `json:"-" firestore:"-" bigquery:"data"`
}

func (ep *EntitlementProfile) BigQueryParse() (err error) {
	bytes, err := json.Marshal(ep)
	if err != nil {
		return err
	}
	ep.Data = string(bytes)
	return nil
}

func (ep *EntitlementProfile) SaveFirestore() error {
	return lib.SetFirestoreErr(lib.EntitlementProfileCollection, ep.Slug, ep)
}

func (ep *EntitlementProfile) SaveBigQuery() error {
	if err := ep.BigQueryParse(); err != nil {
		return err
	}
	return lib.InsertRowsBigQuery(lib.WoptaDataset, lib.EntitlementProfileCollection, ep)
}

func (ep *EntitlementProfile) SaveAll() error {
	if err := ep.SaveFirestore(); err != nil {
		return err
	}
	if err := ep.SaveBigQuery(); err != nil {
		return err
	}
	return nil
}

func (ep *EntitlementProfile) GetFromFirestore(slug string) error {
	snap, err := lib.GetFirestoreErr(lib.EntitlementProfileCollection, slug)
	if err != nil {
		return err
	}
	if err := snap.DataTo(ep); err != nil {
		return err
	}
	return nil
}

type EntitlementProfileService struct{}

func NewEntitlementProfileService() *EntitlementProfileService {
	return &EntitlementProfileService{}
}

func (eps *EntitlementProfileService) GetAllFromFirestore(ctx context.Context) (map[string]EntitlementProfile, error) {
	client, err := firestore.NewClient(ctx, os.Getenv("GOOGLE_PROJECT_ID"))
	if err != nil {
		return nil, err
	}

	profileMap := make(map[string]EntitlementProfile)

	iter := client.Collection(lib.EntitlementProfileCollection).Documents(ctx)
	defer iter.Stop()
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var ep EntitlementProfile
		if err := doc.DataTo(&ep); err != nil {
			return nil, err
		}
		profileMap[ep.Slug] = ep
	}

	return profileMap, nil
}

func CreateProfiles() {
	public := []Entitlement{
		{"auth.authorize"},
		{"auth.token"},
		{"broker.get.policy.fiscalcode"},
		{"broker.get.policy.uid"},
		{"broker.lead"},
		{"broker.proposal"},
		{"broker.emit"},
		{"broker.update.policy"},
		{"broker.get.policy.attachment"},
		{"broker.init"},
		{"callback.sign"},
		{"callback.payment"},
		{"callback.payment.firstrate"},
		{"callback.payment.singlerate"},
		{"callback.email.verify"},
		{"companydata.global.transactions"},
		{"companydata.global.pmi.emit"},
		{"companydata.global.persona.emit"},
		{"companydata.axa.life.emit"},
		{"companydata.sogessur.gap.emit"},
		{"companydata.axa.life.delete"},
		{"companydata.emit"},
		{"companydata.axa.inclusive.bankaccount"},
		{"document.proposal"},
		{"document.contract"},
		{"document.reserved"},
		{"document.sign"},
		{"enrich.munichre.vat"},
		{"enrich.ateco"},
		{"enrich.cities"},
		{"enrich.works"},
		{"enrich.naics"},
		{"form.axa.fleet"},
		{"form.fleet.assistance"},
		{"form.fleet.assistance"},
		{"inclusive.hype.bankaccount"},
		{"inclusive.scalapay.bankaccount"},
		{"inclusive.hype.bankaccount.count"},
		{"inclusive.hybe.bankaccount.import"},
		{"mail.send"},
		{"mail.score"},
		{"mail.validate"},
		{"mga.get.products"},
		{"mga.get.product"},
		{"mga.consume.networknode.invite"},
	}
	internal := []Entitlement{
		{"auth.sso.jwt"},
	}
	node := []Entitlement{
		{"auth.sso.external.product"},
		{"broker.requestapproval"},
		{"broker.get.policy.transactions"},
		{"broker.get.transaction.receipt"},
		{"broker.get.portfolio"},
		{"mga.get.networknode"},
		{"mga.get.consens.undeclared"},
		{"mga.give.consent"},
		{"mga.get.warrants"},
		{"mga.get.quoter.life"},
	}
	admin := []Entitlement{
		{"broker.delete.renew"},
		{"broker.upload.policy.contract"},
		{"broker.duplicate.policy"},
		{"companydata.axa.life.import"},
		{"mga.modify.policy"},
	}
	admin_manager := []Entitlement{
		{"accounting.get.networktransactions"},
		{"accounting.put.networktransaction"},
		{"accounting.create.networktransaction"},
		{"broker.requestapproval"},
		{"broker.delete.policy"},
		{"broker.get.policy.transactions"},
		{"broker.acceptance"},
		{"broker.get.transaction.receipt"},
		{"broker.get.portfolio"},
		{"claim.create"},
		{"claim.get.attachment"},
		{"mga.get.networknode"},
		{"mga.create.networknode"},
		{"mga.update.networknode"},
		{"mga.get.networknodes"},
		{"mga.delete.networknode"},
		{"mga.create.networknode.invite"},
		{"mga.get.warrants"},
		{"mga.create.warrant"},
	}
	customer := []Entitlement{
		{"claim.create"},
		{"claim.get.attachment"},
	}
	manager := append(node, admin_manager...)

	customer = append(customer, public...)
	node = append(node, public...)
	internal = append(internal, public...)
	manager = append(manager, public...)
	admin_manager = append(admin_manager, public...)
	admin = append(admin, admin_manager...)


	profileMap := map[string]EntitlementProfile{
		lib.UserRoleAdmin: {
			Slug: lib.UserRoleAdmin,
			Entitlements: admin,
		},
	}

	for _, value := range profileMap {
		if err := value.SaveAll(); err != nil {
			log.Fatal(err)
		}
	}
}
