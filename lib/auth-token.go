package lib

import "log"

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

	token, err := VerifyUserIdToken(idToken)
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

func (at *AuthToken) GetChannelByRoleV2() string {
	if at.IsNetworkNode {
		return NetworkChannel
	}

	if SliceContains([]string{UserRoleAdmin, UserRoleManager}, at.Role) {
		return MgaChannel
	}

	return ECommerceChannel
}
