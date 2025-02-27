package models

import (
	"context"
	"encoding/json"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/wopta/goworkspace/lib"
	"google.golang.org/api/iterator"
)

// TODO: move me to lib
const EntitlementProfileCollection string = "entitlementProfile"

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
	return lib.SetFirestoreErr(EntitlementProfileCollection, ep.Slug, ep)
}

func (ep *EntitlementProfile) SaveBigQuery() error {
	if err := ep.BigQueryParse(); err != nil {
		return err
	}
	return lib.InsertRowsBigQuery(lib.WoptaDataset, EntitlementProfileCollection, ep)
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

func (ep* EntitlementProfile) GetFromFirestore(slug string) error {
	snap, err := lib.GetFirestoreErr(EntitlementProfileCollection, slug)
	if err != nil {
		return err
	}
	if err := snap.DataTo(ep); err != nil {
		return err
	}
	return nil
}

type EntitlementProfileGenerator struct{}

func NewEntitlementProfileGenerator() *EntitlementProfileGenerator {
	return &EntitlementProfileGenerator{}
}

func (epg *EntitlementProfileGenerator) GetAllFromFirestore(ctx context.Context) (map[string]EntitlementProfile, error) {
	client, err := firestore.NewClient(ctx, os.Getenv("GOOGLE_PROJECT_ID"))
	if err != nil {
		return nil, err
	}

	eps := make(map[string]EntitlementProfile)

	iter := client.Collection(EntitlementProfileCollection).Documents(ctx)
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
		eps[ep.Slug] = ep
	}

	return eps, nil
}
