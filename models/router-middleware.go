package models

import (
	"context"
	"log"
	"net/http"
	"slices"

	"github.com/wopta/goworkspace/lib"
)

func MiddlewareCheckEntitlement(next http.Handler) http.Handler {
	// Simulating Firestore. This info will be DB based
	// var entitlementMapping = map[string][]Entitlement{
	// 	lib.UserRoleAdmin:       {{lib.EntitlementBrokerAcceptance}, {lib.EntitlementBrokerDuplicate}},
	// 	lib.UserRoleAgent:       {{lib.EntitlementBrokerLead}, {lib.EntitlementBrokerEmit}},
	// 	lib.UserRoleAgency:      {{lib.EntitlementBrokerLead}, {lib.EntitlementBrokerEmit}},
	// 	lib.UserRoleAreaManager: {},
	// 	lib.UserRoleManager:     {},
	// 	lib.UserRoleAll:         {{lib.EntitlementBrokerLead}, {lib.EntitlementBrokerEmit}},
	// 	lib.UserRoleCustomer:    {},
	// 	lib.UserRoleInternal:    {},
	// }

	epg := NewEntitlementProfileGenerator()
	entitlementMapping, err := epg.GetAllFromFirestore(context.Background())
	if err != nil {
		log.Printf("error retrieving entitlement profiles: %s", err.Error())
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		authToken := ctx.Value(lib.CtxAuthToken).(AuthToken)
		entitlement := ctx.Value(lib.CtxEntitlement).(string)

		userEntitlements := entitlementMapping[authToken.Role].Entitlements
		
		if authToken.IsNetworkNode {
			snap, err := lib.GetFirestoreErr(lib.NetworkNodesCollection, authToken.UserID)
			if err != nil {
				log.Printf("error extracting network node: %s", err.Error())
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}
			var nn NetworkNode
			snap.DataTo(&nn)

			r = r.WithContext(context.WithValue(r.Context(), lib.CtxNetworkNode, nn))

			if nn.EntitlementProfile.Slug != "" {
				userEntitlements = entitlementMapping[nn.EntitlementProfile.Slug].Entitlements
			}

			if len(nn.EntitlementProfile.Entitlements) > 0 {
				userEntitlements = nn.EntitlementProfile.Entitlements
			}
		}

		if !slices.ContainsFunc(userEntitlements, func(e Entitlement) bool {
			return e.Slug == entitlement
		}) {
			log.Printf("user does not have '%s' entitlement", entitlement)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
