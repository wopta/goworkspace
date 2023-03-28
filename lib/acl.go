package lib

import (
	"github.com/wopta/goworkspace/broker"
	"github.com/wopta/goworkspace/user"
)

func HasUserAccessToPolicy(jwt string, policyID string) (bool, error) {
	client, ctx := getClient()

	token, err := client.VerifyIDToken(ctx, jwt)
	if err != nil {
		return false, err
	}

	authID := token.Claims["user_id"]

	usr, err := user.GetUserByAuthId(authID.(string))
	CheckError(err)

	policies := broker.GetPoliciesFromFirebase(usr.FiscalCode)

	for _, policy := range policies {
		if policy.ID == policyID {
			return true, err
		}
	}

	return false, err
}
