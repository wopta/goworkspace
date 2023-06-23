package models

import (
	"log"

	"github.com/wopta/goworkspace/lib"
)

type AuthToken struct {
	Role   string `json:"role"`
	UserID string `json:"userId"`
	Email  string `json:"email"`
}

func GetAuthTokenFromIdToken(idToken string) (AuthToken, error) {
	if idToken == "" {
		return AuthToken{
			Role:   UserRoleAll,
			UserID: "",
			Email:  "",
		}, nil
	}

	token, err := lib.VerifyUserIdToken(idToken)
	if err != nil {
		log.Printf("[GetAuthTokenFromIdToken] idToken: %s , err: %v", idToken, err)
		return AuthToken{}, err
	}

	return AuthToken{
		Role:   token.Claims["role"].(string),
		UserID: token.Claims["user_id"].(string),
		Email:  token.Claims["email"].(string),
	}, nil
}
