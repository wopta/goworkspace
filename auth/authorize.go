package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	lib "github.com/wopta/goworkspace/lib"
)

func AuthorizeFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println("--------------------------AuthorizeFx-------------------------------------------")
	var (
		serviceAccountReq  ServiceAccount
		serviceAccountList []ServiceAccount
		tokenString        string
		e                  error
	)
	origin = r.Header.Get("Origin")

	rBody := lib.ErrorByte(io.ReadAll(r.Body))

	e = json.Unmarshal(rBody, &serviceAccountReq)
	credByte := lib.GetFilesByEnv("auth/clients-credential")
	e = json.Unmarshal(credByte, &serviceAccountList)
	for _, sa := range serviceAccountList {
		if sa.ClientId == serviceAccountReq.ClientId && sa.ClientSecret == serviceAccountReq.ClientSecret {
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"role": "3thParty",
				"nbf":  time.Now().Unix(),
			})
			// Sign and get the complete encoded token as a string using the secret
			tokenString, e = token.SignedString(os.Getenv("JWTSIGNKEY"))

		}
	}

	//log.Println("Proposal request proposal: ", string(j))
	defer r.Body.Close()
	return tokenString, nil, e
}
func TokenFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println("--------------------------TokenFx-------------------------------------------")
	var (
		e error
	)
	origin = r.Header.Get("Origin")
	tokenreq := r.Header.Get("Auth")

	token, v, e := Token(tokenreq)
	defer r.Body.Close()
	return token, v, e
}
func Token(tokenReq string) (string, interface{}, error) {
	log.Println("--------------------------Token-------------------------------------------")
	var (
		e error
	)

	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenReq, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return os.Getenv("JWTSIGNKEY"), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims["role"], claims["nbf"])
		res := strconv.FormatBool(token.Valid)
		return res, nil, e
	} else {
		fmt.Println(err)
	}

	return "", token.Valid, e
}

type ServiceAccount struct {
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	Role         string `json:"role"`
}
