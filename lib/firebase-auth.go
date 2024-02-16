package lib

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
)

func getClient() (*auth.Client, context.Context) {
	ctx := context.Background()
	app, err := firebase.NewApp(ctx, &firebase.Config{ProjectID: os.Getenv("GOOGLE_PROJECT_ID")})
	client, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}
	return client, ctx
}

func CreateCustomJwt(email, role, nodeType, id string) (string, error) {
	client, ctx := getClient()
	claims := map[string]interface{}{
		"role":          role,
		"type":          nodeType,
		"email":         email,
		"isNetworkNode": true,
	}
	token, err := client.CustomTokenWithClaims(ctx, id, claims)
	if err != nil {
		log.Fatalf("error minting custom token: %v\n", err)
	}
	log.Printf("Got custom token: %v\n", token)
	return token, err
}

func CreateUserWithEmailAndPassword(email string, password string, id *string) (*auth.UserRecord, error) {
	client, ctx := getClient()
	params := (&auth.UserToCreate{}).
		Email(email).
		Password(password)
	if id != nil && len(*id) > 0 {
		params.UID(*id)
	}
	u, err := client.CreateUser(ctx, params)
	if err == nil {
		log.Printf("Successfully created user: %v\n", u)
	} else {
		log.Printf("Error creating user: %v\n", err)
	}
	return u, err
}

func UpdateUserEmail(uid, email string) (*auth.UserRecord, error) {
	client, ctx := getClient()
	params := (&auth.UserToUpdate{}).
		Email(email).
		EmailVerified(true)
	userRecord, err := client.UpdateUser(ctx, uid, params)
	if err != nil {
		log.Printf("[UpdateUserEmail] error updating user: %v\n", err)
		return nil, err
	}
	log.Printf("[UpdateUserEmail] successfully updated user: %v\n", userRecord)
	return userRecord, err
}

func HandleUserAuthenticationStatus(uid string, disabled bool) error {
	client, ctx := getClient()
	params := (&auth.UserToUpdate{}).
		Disabled(disabled)
	userRecord, err := client.UpdateUser(ctx, uid, params)
	if err != nil {
		log.Printf("[UpdateUserEmail] error updating user: %v\n", err)
		return err
	}
	log.Printf("[UpdateUserEmail] successfully updated user: %v\n", userRecord)
	return err
}

func VerifyUserIdToken(idToken string) (*auth.Token, error) {
	client, ctx := getClient()
	token, err := client.VerifyIDTokenAndCheckRevoked(ctx, strings.ReplaceAll(idToken, "Bearer ", ""))
	return token, err
}

func GetUserIdFromIdToken(idToken string) (string, error) {
	token, err := VerifyUserIdToken(idToken)
	if err != nil {
		return "", err
	}
	return token.Claims["user_id"].(string), err
}

func GetUserRoleFromIdToken(idToken string) (string, error) {
	token, err := VerifyUserIdToken(idToken)
	if err != nil {
		return "", err
	}
	return token.Claims["role"].(string), err
}

func VerifyAuthorization(handler func(w http.ResponseWriter, r *http.Request) (string, interface{}, error), roles ...string) func(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	wrappedHandler := func(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
		errorHandler := func(w http.ResponseWriter) (string, interface{}, error) {
			return "", nil, fmt.Errorf("not found")
		}

		if len(roles) == 0 || os.Getenv("env") == "local" || SliceContains(roles, "internal") {
			return handler(w, r)
		}

		if !SliceContains(roles, "all") {
			idToken := strings.ReplaceAll(r.Header.Get("Authorization"), "Bearer ", "")
			if idToken == "" {
				return errorHandler(w)
			}

			token, err := VerifyUserIdToken(idToken)
			if err != nil {
				log.Println("VerifyAuthorization: verify id token error: ", err)
				return errorHandler(w)
			}

			userRole := "customer"
			if role, ok := token.Claims["role"].(string); ok {
				userRole = role
			}

			if !SliceContains(roles, userRole) {
				return errorHandler(w)
			}
		}

		return VerifyAppcheck(handler)(w, r)
	}

	return wrappedHandler
}

func SetCustomClaimForUser(uid string, claims map[string]interface{}) {
	client, ctx := getClient()

	err := client.SetCustomUserClaims(ctx, uid, claims)
	if err != nil {
		log.Fatalf("error setting custom claims %v\n", err)
	}
}

func GetAuthUserIdByEmail(mail string) (string, error) {
	client, ctx := getClient()

	user, err := client.GetUserByEmail(ctx, mail)

	if err != nil {
		return "", err
	}
	return user.UID, nil
}

func VerifyAppcheck(handler func(w http.ResponseWriter, r *http.Request) (string, interface{}, error)) func(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {

	wrappedHandler := func(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
		errorHandler := func(w http.ResponseWriter) (string, interface{}, error) {
			log.Println("[VerifyAppcheck]: unauthenticated.")
			return "", nil, fmt.Errorf("Unavailable")
		}

		ctx := context.Background()
		app, err := firebase.NewApp(ctx, &firebase.Config{ProjectID: os.Getenv("GOOGLE_PROJECT_ID")})
		if err != nil {
			log.Fatalf("error initializing app: %v\n", err)
			return errorHandler(w)
		}

		appCheck, err := app.AppCheck(context.Background())
		if err != nil {
			log.Fatalf("error initializing app: %v\n", err)
			return errorHandler(w)
		}

		appCheckToken, ok := r.Header[http.CanonicalHeaderKey("X-Firebase-AppCheck")]
		if !ok {
			return errorHandler(w)
		}

		_, err = appCheck.VerifyToken(appCheckToken[0])
		if err != nil {
			return errorHandler(w)
		}

		return handler(w, r)
	}
	return wrappedHandler
}
