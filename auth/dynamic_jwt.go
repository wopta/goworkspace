package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
)

func DynamicJwtFx(w http.ResponseWriter, r *http.Request) (string, any, error) {
	var (
		tokenString    string
		err            error
		node           []models.NetworkNode
		bytes          []byte
		responseSsoJwt ResponseSsoJwt
	)

	log.AddPrefix("DynamicJwtFx")
	log.Println("Handler start -----------------------------------------------")

	defer func() {
		log.Println("Handler end ---------------------------------------------")
		log.PopPrefix()
	}()

	jwt := r.URL.Query().Get("jwt")
	provider := strings.ReplaceAll(chi.URLParam(r, "provider"), "-", "_")

	// TODO: remove me when all AUA nodes have their own jwtConfig
	if strings.EqualFold(provider, "aua") {
		return JwtAuaFx(w, r)
	}

	key := os.Getenv(fmt.Sprintf("%s_SIGNING_KEY", lib.ToUpper(provider)))

	if key == "" {
		log.Printf("unhandled provider: %s", provider)
		return "", nil, fmt.Errorf("unhandled provider")
	}

	claims, isValid, err := verifyJwt(jwt, key)

	if isValid {
		q := lib.FireGenericQueries[models.NetworkNode]{
			Queries: []lib.Firequery{
				{
					Field:      "externalNetworkCode",
					Operator:   "==",
					QueryValue: claims.Id,
				},
			},
		}
		node, err = q.FireQuery(lib.NetworkNodesCollection)
		if len(node) > 0 {
			if node[0].AuthId == "" {
				userfire, _ := lib.CreateUserWithEmailAndPassword(node[0].Mail, os.Getenv("DEFAULT_PSW"), &node[0].Uid)
				node[0].AuthId = userfire.UID
				err = node[0].SaveFirestore()
				if err != nil {
					log.ErrorF("error updating node %s in Firestore: %s", node[0].Uid, err.Error())
					return "", nil, err
				}
				err = node[0].SaveBigQuery()
				if err != nil {
					log.ErrorF("error updating node %s in BigQuery: %s", node[0].Uid, err.Error())
					return "", nil, err
				}

			}
			tokenString, err = lib.CreateCustomJwt(node[0].Mail, node[0].Role, node[0].Type, node[0].AuthId)
			if err != nil {
				log.ErrorF("error creating token: %s", err.Error())
				return "", nil, err
			}
			responseSsoJwt = ResponseSsoJwt{
				Token:    tokenString,
				Producer: node[0],
			}
			responseSsoJwt.Producer.JwtConfig = lib.JwtConfig{} // Do not expose inner configs to frontend
			responseSsoJwt.Producer.CallbackConfig = nil        // Do not expose inner configs to frontend
			bytes, err = json.Marshal(responseSsoJwt)
		}
	}

	return string(bytes), responseSsoJwt, err
}

type Claims = AuaClaims

func verifyJwt(jwtData, key string) (claims *Claims, isValid bool, err error) {
	token, e := jwt.ParseWithClaims(jwtData, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header)
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(key), nil
	})
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		log.Printf("parsed claims: %+v", claims)
		return claims, token.Valid, e
	} else {
		log.Printf("claims error: %s", e.Error())
	}

	return nil, token.Valid, e
}
