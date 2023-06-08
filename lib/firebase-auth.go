package lib

import (
	"context"
	"fmt"
	"github.com/wopta/goworkspace/models"
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

func VerifyUserIdToken(idToken string) (*auth.Token, error) {
	client, ctx := getClient()

	token, err := client.VerifyIDToken(ctx, idToken)

	return token, err
}

func GetUserIdFromIdToken(idToken string) (string, error) {
	token, err := VerifyUserIdToken(idToken)
	if err != nil {
		return "", err
	}

	return token.Claims["user_id"].(string), err
}

func VerifyAuthorization(handler func(w http.ResponseWriter, r *http.Request) (string, interface{}, error), roles ...string) func(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	wrappedHandler := func(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
		errorHandler := func(w http.ResponseWriter, err error) (string, interface{}, error) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Not Found"))
			return "", nil, err
		}

		idToken := strings.ReplaceAll(r.Header.Get("Authorization"), "Bearer ", "")
		token, err := VerifyUserIdToken(idToken)
		if err != nil {
			return errorHandler(w, err)
		}

		userRole := token.Claims["role"].(string)
		if len(roles) == 1 && roles[0] == models.UserRoleAll || SliceContains(roles, userRole) {
			return handler(w, r)
		}

		return errorHandler(w, fmt.Errorf("not found"))
	}

	return wrappedHandler
}
