package models

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"slices"

	"github.com/go-chi/chi/v5"
	"github.com/wopta/goworkspace/lib"
)

const (
	CtxRequesterNetworkNode = lib.Ctxkey("requesternetworknode")
)

func GetExtendedRouter(module string, routes []lib.Route) *chi.Mux {
	for idx := range routes {
		routes[idx].Middlewares = append(routes[idx].Middlewares, withRequesterNetworkNode, withCheckEntitlement)
	}

	return lib.GetRouter(module, routes)
}

func withRequesterNetworkNode(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		authToken := ctx.Value(lib.CtxAuthToken).(AuthToken)

		if authToken.IsNetworkNode {
			snap, err := lib.GetFirestoreErr(lib.NetworkNodesCollection, authToken.UserID)
			if err != nil {
				log.Printf("error extracting network node: %s", err.Error())
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}
			var nn NetworkNode
			if err := snap.DataTo(&nn); err != nil {
				log.Printf("error parsing network node: %s", err.Error())
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}

			r = r.WithContext(context.WithValue(r.Context(), CtxRequesterNetworkNode, nn))
		}
		
		next.ServeHTTP(w, r)
	})
}

func withCheckEntitlement(next http.Handler) http.Handler {
	// eps := NewEntitlementProfileService()
	// entitlementMapping, err := eps.GetAllFromFirestore(context.Background())
	// if err != nil {
	// 	log.Printf("error retrieving entitlement profiles: %s", err.Error())
	// }
	entitlementMapping := CreateProfiles()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		authToken := ctx.Value(lib.CtxAuthToken).(AuthToken)
		entitlement := ctx.Value(lib.CtxEntitlement).(string)

		var userEntitlements []Entitlement
		if _, ok := entitlementMapping[authToken.Role]; ok {
			userEntitlements = entitlementMapping[authToken.Role].Entitlements
		}

		if authToken.IsNetworkNode {
			nn := ctx.Value(CtxRequesterNetworkNode).(NetworkNode)

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

		resp := struct{ Response string `json:"response"`}{
			Response: "passed entitlement check!!",
		}
			
		log.Println(resp.Response)
		w.WriteHeader(200)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "something went wrong", http.StatusServiceUnavailable)
			return
		}
		return

		next.ServeHTTP(w, r)
	})
}
