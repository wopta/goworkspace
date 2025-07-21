package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	jwt "github.com/golang-jwt/jwt/v5"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
)

func JwtAuaFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		tokenString    string
		e              error
		node           []models.NetworkNode
		b              []byte
		responseSsoJwt ResponseSsoJwt
	)

	log.AddPrefix("JwtAuaFx")
	defer log.PopPrefix()

	log.Println("Handler start -----------------------------------------------")

	tokenReq := r.URL.Query().Get("jwt")

	log.Println("JwtFx request token:", tokenReq)

	claims, isvalid, e := verifyAuaJwt(tokenReq)

	if isvalid {
		q := lib.FireGenericQueries[models.NetworkNode]{
			Queries: []lib.Firequery{
				{
					Field:      "externalNetworkCode",
					Operator:   "==",
					QueryValue: claims.Id,
				},
			},
		}
		node, e = q.FireQuery(lib.NetworkNodesCollection)
		if len(node) > 0 {
			if node[0].AuthId == "" {
				userfire, _ := lib.CreateUserWithEmailAndPassword(node[0].Mail, os.Getenv("DEFAULT_PSW"), &node[0].Uid)
				node[0].AuthId = userfire.UID
				e = node[0].SaveFirestore()
				if e != nil {
					log.ErrorF("error updating node %s in Firestore: %s", node[0].Uid, e.Error())
					return "", nil, e
				}
				e = node[0].SaveBigQuery()
				if e != nil {
					log.ErrorF("error updating node %s in BigQuery: %s", node[0].Uid, e.Error())
					return "", nil, e
				}

			}
			tokenString, e = lib.CreateCustomJwt(node[0].Mail, node[0].Role, node[0].Type, node[0].AuthId)
			responseSsoJwt = ResponseSsoJwt{
				Token:    tokenString,
				Producer: node[0],
			}
			responseSsoJwt.Producer.JwtConfig = lib.JwtConfig{} // Do not expose inner configs to frontend
			responseSsoJwt.Producer.CallbackConfig = nil        // Do not expose inner configs to frontend
			b, e = json.Marshal(responseSsoJwt)
		}
	}

	log.Println("Handler end -------------------------------------------------")

	return string(b), responseSsoJwt, e
}

func verifyAuaJwt(tokenReq string) (*AuaClaims, bool, error) {
	log.Println("--------------------------Token-------------------------------------------")
	var (
		e error
	)

	token, e := jwt.ParseWithClaims(tokenReq, &AuaClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header)
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("AUAJWTSIGNKEY")), nil
	})
	if claims, ok := token.Claims.(*AuaClaims); ok && token.Valid {
		fmt.Println(claims)
		return claims, token.Valid, e
	} else {
		fmt.Println(e)
	}

	return nil, token.Valid, e
}

type ResponseSsoJwt struct {
	Token    string             `json:"token"`
	Producer models.NetworkNode `json:"producer"`
}
type AuaClaims struct {
	Id         string `json:"codSubAgent"`
	Name       string `json:"name"`
	Exp        int    `json:"exp"`
	Mail       string `json:"email"`
	AgencyCode string `json:"codiceagenzia"`
	jwt.RegisteredClaims
}
