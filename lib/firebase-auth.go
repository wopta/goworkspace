package lib

import (
	"context"
	"log"
	"os"

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

func GetUserIdFromIdToken(idToken string) (string, error) {
	client, ctx := getClient()

	token, err := client.VerifyIDToken(ctx, idToken)
	if err != nil {
		return "", err
	}

	return token.Claims["user_id"].(string), err
}

func SetCustomClaimForUser(uid string, claims map[string]interface{}) {
	// Get an auth client from the firebase.App
	client, ctx := getClient()

	err := client.SetCustomUserClaims(ctx, uid, claims)
	if err != nil {
		log.Fatalf("error setting custom claims %v\n", err)
	}
}
