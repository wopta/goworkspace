package AppcheckProxy

import (
	"context"
	"log"
	"net/http"

	firebaseAdmin "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/appcheck"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	lib "github.com/wopta/goworkspace/lib"
)

func init() {
	log.Println("INIT AppcheckProxy")
	functions.HTTP("AppcheckProxy", AppcheckProxy)
}

func AppcheckProxy(w http.ResponseWriter, r *http.Request) {
	var (
		//idToken  string
		appCheck *appcheck.Client
	)
	app, err := firebaseAdmin.NewApp(context.Background(), nil)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	appCheck, err = app.AppCheck(context.Background())
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
	appCheckToken, ok := r.Header[http.CanonicalHeaderKey("X-Firebase-AppCheck")]
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized."))
		return
	}

	_, err = appCheck.VerifyToken(appCheckToken[0])
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized."))
		return
	}

	// If VerifyToken() succeeds, continue with the provided handler.

	lib.CheckError(err)

}
