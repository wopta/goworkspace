package models

import (
	"log"

	"github.com/wopta/goworkspace/lib"
)

type AuthToken struct {
	Role          string `json:"role"`
	Type          string `json:"type"`
	UserID        string `json:"userId"`
	Email         string `json:"email"`
	IsNetworkNode bool   `json:"isNetworkNode"`
}

func GetAuthTokenFromIdToken(idToken string) (AuthToken, error) {
	if idToken == "" {
		return AuthToken{
			Role:   UserRoleAll,
			Type:   "",
			UserID: "",
			Email:  "",
		}, nil
	}

	token, err := lib.VerifyUserIdToken(idToken)
	if err != nil {
		log.Printf("[GetAuthTokenFromIdToken] idToken: %s , err: %v", idToken, err)
		return AuthToken{}, err
	}

	nodeType := ""
	if token.Claims["type"] != nil {
		nodeType = token.Claims["type"].(string)
	}

	isNetworkNode := false
	if token.Claims["isNetworkNode"] != nil {
		isNetworkNode = token.Claims["isNetworkNode"].(bool)
	}

	return AuthToken{
		Role:          token.Claims["role"].(string),
		Type:          nodeType,
		UserID:        token.Claims["user_id"].(string),
		Email:         token.Claims["email"].(string),
		IsNetworkNode: isNetworkNode,
	}, nil
}

// DEPRECATED: remove once product versioning completed
func (at *AuthToken) GetChannelByRole() string {
	channel := ECommerceChannel

	switch at.Role {
	case UserRoleAdmin, UserRoleManager:
		channel = MgaChannel
	case UserRoleAgency:
		channel = AgencyChannel
	case UserRoleAgent:
		channel = AgentChannel
	}

	return channel
}

func (at *AuthToken) GetChannelByRoleV2() string {
	if at.IsNetworkNode {
		return NetworkChannel
	}

	if lib.SliceContains([]string{UserRoleAdmin, UserRoleManager}, at.Role) {
		return MgaChannel
	}

	return ECommerceChannel
}
