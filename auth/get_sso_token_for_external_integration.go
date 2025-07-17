package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/network"
)

func GetTokenForExternalIntegrationFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		e        error
		b        []byte
		response GetTokenForExternalIntegrationResponse
	)

	log.AddPrefix("GetTokenForExternalIntegrationFx")
	defer log.PopPrefix()

	log.Println("Handler start -----------------------------------------------")

	productName := chi.URLParam(r, "productName")
	token := r.Header.Get("Authorization")
	authToken, e := lib.GetAuthTokenFromIdToken(token)

	response.Token, e = getTokenForExternalIntegration(productName, authToken.UserID)

	b, e = json.Marshal(response)

	log.Println("Handler end -------------------------------------------------")

	return string(b), response, e
}

func getTokenForExternalIntegration(productName string, userUid string) (string, error) {
	var code string

	log.Println("--------------------------getTokenForExternalIntegration-------------------------------------------")

	node := network.GetNetworkNodeByUid(userUid)
	code = node.ExternalNetworkCode
	if code == "" {
		code = node.Code
	}

	// verify if user has access to the product
	warrant := node.GetWarrant()
	if !warrant.HasProductByName(productName) {
		return "", fmt.Errorf("node does not have access to the product")
	}

	// Define the signing key
	signingKey := []byte(os.Getenv(fmt.Sprintf("%s_EXTERNAL_SIGNING_KEY", lib.ToUpper(strings.ReplaceAll(productName, "-", "_")))))
	if len(signingKey) == 0 {
		return "", fmt.Errorf("signing key not found")
	}

	// Set the expiration time
	expirationTime := time.Now().Add(30 * time.Minute)

	// Set the not-before time
	notBeforeTime := time.Now()

	// Create the JWT token with claims
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"codSubAgent": code,
		"name":        node.GetName(),
		"email":       node.Mail,
		"exp":         expirationTime.Unix(),
		"nbf":         notBeforeTime.Unix(),
	})

	tokenString, err := jwtToken.SignedString(signingKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

type GetTokenForExternalIntegrationResponse struct {
	Token string `json:"token"`
}
