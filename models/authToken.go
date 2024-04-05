package models

import "github.com/wopta/goworkspace/lib"

// TODO: delete me when all modules read directly form lib

type AuthToken = lib.AuthToken

var GetAuthTokenFromIdToken func(idToken string) (lib.AuthToken, error) = lib.GetAuthTokenFromIdToken
