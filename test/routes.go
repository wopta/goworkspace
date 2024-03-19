package test

import (
	"context"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// for local testing only
func init() {
	log.Println("INIT Test")
	functions.HTTP("Test", Test)
}

func Test(w http.ResponseWriter, r *http.Request) {
	prefix := "/"

	if os.Getenv("env") == "local" {
		prefix = "/test/"
	}

	mux := chi.NewRouter()
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(CorsMiddleware)
	w.Header().Add("Content-type", "application/json")

	mux.Route(prefix, func(r chi.Router) {
		r.Mount("/public", publicRouter())
		r.Mount("/network", networkRouter())
		r.Mount("/internal", internalRouter())
		r.Mount("/admin", adminRouter())
	})

	mux.ServeHTTP(w, r)
}

func publicRouter() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.WithValue("roles", []string{"all"}))
	r.Use(AppCheckMiddleware)
	r.Use(CheckEntitlement)

	r.Route("/test2", func(r chi.Router) {
		r.Get("/", test2)
		r.Get("/{param}", test2)
	})

	return r
}

func networkRouter() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.WithValue("roles", []string{"admin", "agent", "agency", "manager"}))
	r.Use(AppCheckMiddleware)
	r.Use(CheckEntitlement)

	r.Route("/v1/policy", func(r chi.Router) {
		r.Get("/", test1)
		r.Route("/{policyUid}", func(r chi.Router) {
			r.Get("/", test2)
			r.Get("/transactions", test2)
		})
	})

	return r
}

func internalRouter() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.WithValue("roles", []string{"internal"}))
	r.Use(AppCheckMiddleware)
	r.Use(CheckEntitlement)

	r.Get("/", test1)

	return r
}

func adminRouter() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.WithValue("roles", []string{"admin"}))
	r.Use(AppCheckMiddleware)
	r.Use(CheckEntitlement)

	r.Route("/v1/policy", func(r chi.Router) {
		r.Get("/", test1)
		r.Route("/{policyUid}/transactions", func(r chi.Router) {
			r.Get("/", test2)
			r.Get("/{transactionUid}", test2)
		})
	})

	return r
}

func test1(w http.ResponseWriter, r *http.Request) {
	log.Println("test1 handler!")
	log.Printf("Request: %s", r.RequestURI)
	w.Write([]byte(`{}`))
}

func test2(w http.ResponseWriter, r *http.Request) {
	log.Println("test2 handler!")
	log.Printf("Request: %s", r.RequestURI)

	policyUid := chi.URLParam(r, "policyUid")
	log.Printf("PolicyUid: %s", policyUid)
	transactionUid := chi.URLParam(r, "transactionUid")
	log.Printf("TransactionUid: %s", transactionUid)

	filterQuery := r.URL.Query().Get("filter")
	log.Printf("FilterQuery: %s", filterQuery)

	w.Write([]byte(`{"success":true}`))
}

// lib
func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		options := cors.Options{
			AllowedHeaders: []string{"*"},
		}

		if r.Method == http.MethodOptions || os.Getenv("env") == "local" {
			options.AllowedOrigins = []string{"*"}
		}

		if r.Method == http.MethodOptions {
			options.AllowedMethods = []string{"*"}
			options.AllowCredentials = true
			options.MaxAge = 3600
			w.WriteHeader(http.StatusNoContent)
		} else {
			options.AllowCredentials = false
		}

		cors.Handler(options)(next).ServeHTTP(w, r)
	})
}

func AppCheckMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		roles := ctx.Value("roles").([]string)
		log.Printf("AppCheckMiddleware roles: %v", roles)

		if os.Getenv("env") == "local" || slices.Contains(roles, "internal") {
			next.ServeHTTP(w, r)
			return
		}

		app, err := firebase.NewApp(ctx, &firebase.Config{ProjectID: os.Getenv("GOOGLE_PROJECT_ID")})
		if err != nil {
			log.Printf("error initializing app: %s", err.Error())
			http.Error(w, "get out of here", http.StatusForbidden)
			return
		}

		appCheck, err := app.AppCheck(ctx)
		if err != nil {
			log.Printf("error initializing app: %s\n", err.Error())
			http.Error(w, "get out of here", http.StatusForbidden)
			return
		}

		appCheckToken, ok := r.Header[http.CanonicalHeaderKey("X-Firebase-AppCheck")]
		if !ok {
			log.Printf("error missing token")
			http.Error(w, "get out of here", http.StatusForbidden)
			return
		}

		_, err = appCheck.VerifyToken(appCheckToken[0])
		if err != nil {
			log.Printf("error invalid token: %s", err.Error())
			http.Error(w, "get out of here", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func VerifyUserIdToken(idToken string) (*auth.Token, error) {
	ctx := context.Background()
	app, err := firebase.NewApp(ctx, &firebase.Config{ProjectID: os.Getenv("GOOGLE_PROJECT_ID")})
	if err != nil {
		log.Fatalf("error creating app: %v\n", err)
	}
	client, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}
	token, err := client.VerifyIDTokenAndCheckRevoked(ctx, idToken)
	return token, err
}

func CheckEntitlement(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		roles := ctx.Value("roles").([]string)
		log.Printf("CheckEntitlement roles: %v", roles)

		if /*len(roles) == 0 || os.Getenv("env") == "local" ||*/ slices.Contains(roles, "internal") || slices.Contains(roles, "all") {
			next.ServeHTTP(w, r)
			return
		}

		idToken := strings.ReplaceAll(r.Header.Get("Authorization"), "Bearer ", "")
		if idToken == "" {
			log.Println("VerifyAuthorization: empty token")
			http.Error(w, "who do you think you are", http.StatusUnauthorized)
			return
		}

		token, err := VerifyUserIdToken(idToken)
		if err != nil {
			log.Printf("VerifyAuthorization: verify id token error: %s", err.Error())
			http.Error(w, "who do you think you are", http.StatusUnauthorized)
			return
		}

		userRole := "customer"
		if role, ok := token.Claims["role"].(string); ok {
			userRole = role
		}

		if !slices.Contains(roles, userRole) {
			log.Printf("VerifyAuthorization: userRole '%s' not allowed", userRole)
			http.Error(w, "who do you think you are", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
