package AppcheckProxy

import (
	"context"
	"log"
	"net/http"

	firebase "firebase.google.com/go/v4"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	lib "github.com/wopta/goworkspace/lib"
	//"google.golang.org/api/firebaseappcheck/v1"
)

func init() {
	log.Println("INIT AppcheckProxy")
	functions.HTTP("AppcheckProxy", AppcheckProxy)
}

func AppcheckProxy(w http.ResponseWriter, r *http.Request) {
	var (
		idToken string
	)
	ctx := context.Background()
	//firebaseappcheckService, err := firebaseappcheck.NewService(ctx)
	app, err := firebase.NewApp(context.Background(), nil)
	client, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}

	token, err := client.VerifyIDToken(ctx, idToken)
	if err != nil {
		log.Fatalf("error verifying ID token: %v\n", err)
	}

	log.Printf("Verified ID token: %v\n", token)

	lib.CheckError(err)

}
