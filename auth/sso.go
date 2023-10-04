package auth

import (
	"fmt"
	"log"
	"net/http"
	"os"

	jwt "github.com/golang-jwt/jwt/v5"
	lib "github.com/wopta/goworkspace/lib"
	models "github.com/wopta/goworkspace/models"
)

func JwtFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println("--------------------------JwtFx-------------------------------------------")
	var (
		tokenString string
		e           error
		node        []models.NetworkNode
	)
	origin = r.Header.Get("Origin")
	tokenReq := r.URL.Query().Get("jwt")
	log.Println("JwtFx request token:", tokenReq)
	log.Println("JwtFx AUAJWTSIGNKEY:", os.Getenv("AUAJWTSIGNKEY"))
	claims, isvalid, e := verifyAuaJwt(tokenReq)

	if isvalid {
		q := lib.FireGenericQueries[models.NetworkNode]{
			Queries: []lib.Firequery{
				{
					Field:      "networkCode",
					Operator:   "==",
					QueryValue: claims.Id,
				},
			},
		}
		node, e = q.FireQuery("networkNodes")
		if len(node) > 0 {

			tokenString, e = lib.CreateCustomJwt("", "", node[0].AuthId)
		}
	}
	//log.Println("Proposal request proposal: ", string(j))
	defer r.Body.Close()
	return tokenString, nil, e
}

func verifyAuaJwt(tokenReq string) (*AuaClaims, bool, error) {
	log.Println("--------------------------Token-------------------------------------------")
	var (
		e error
	)

	token, e := jwt.ParseWithClaims(tokenReq, &AuaClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return os.Getenv("AUAJWTSIGNKEY"), nil
	})
	if claims, ok := token.Claims.(*AuaClaims); ok && token.Valid {
		fmt.Println(claims)
		return claims, token.Valid, e
	} else {
		fmt.Println(e)
	}

	return nil, token.Valid, e
}

type AuaClaims struct {
	Id   string `json:"codSubAgent"`
	Name string `json:"name"`
	Exp  int    `json:"exp"`
	jwt.RegisteredClaims
}
