package broker

import (
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/user"
)

func HasUserAccessToPolicy(idToken string, policyID string, origin string) (bool, error) {

	authID, err := lib.GetUserIdFromIdToken(idToken)
	lib.CheckError(err)

	usr, err := user.GetUserByAuthId(authID)
	lib.CheckError(err)

	policies := GetPoliciesFromFirebase(usr.FiscalCode, origin)

	for _, policy := range policies {
		if policy.Uid == policyID {
			return true, err
		}
	}

	return false, err
}
