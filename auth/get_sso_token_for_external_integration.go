package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	models "github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
)

func GetTokenForExternalIntegrationFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		e              error
		b              []byte
		response GetTokenForExternalIntegrationResponse
	)

	log.Println("--------------------------GetTokenForExternalIntegrationFx-------------------------------------------")

	origin = r.Header.Get("Origin")
	productName := r.Header.Get("productName")
	token := r.Header.Get("Authorization")
	authToken, e := models.GetAuthTokenFromIdToken(token)

	response.Token, e = getTokenForExternalIntegration(productName, authToken.UserID)

	b, e = json.Marshal(response)
	return string(b), response, e
}

func getTokenForExternalIntegration(productName string, userUid string) (string, error) {
	log.Println("--------------------------getTokenForExternalIntegration-------------------------------------------")

	node := network.GetNetworkNodeByUid(userUid)

	// verify if user has access to the product
	warrant := node.GetWarrant()
	if !warrant.HasProductByName(productName) {
		return "", fmt.Errorf("node does not have access to the product")
	}

	// Define the signing key
	signingKey := []byte(os.Getenv("AUAJWTSIGNKEY"))

	// Set the expiration time
	expirationTime := time.Now().Add(30 * time.Minute)

	// Set the not-before time
	notBeforeTime := time.Now()

	// Create the JWT token with claims
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"codSubAgent": node.Code,
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
