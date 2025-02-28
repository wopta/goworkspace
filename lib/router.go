package lib

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"

	"cloud.google.com/go/civil"
	firebase "firebase.google.com/go/v4"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type Route struct {
	Route       string
	Method      string
	Handler     http.HandlerFunc
	Fn          func(http.ResponseWriter, *http.Request) (string, any, error)
	Middlewares []RouteMiddleware
	Roles       []string
	Entitlement string
}

type RouteMiddleware = func(http.Handler) http.Handler

type Ctxkey string

const (
	roles          = Ctxkey("roles")
	CtxEntitlement = Ctxkey("entitlement")
	CtxAuthToken   = Ctxkey("authToken")
)

func responseLoggerWrapper(handler func(w http.ResponseWriter, r *http.Request) (string, any, error)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		str, _, err := handler(w, r)
		if err != nil {
			log.Printf("Error: %s", err.Error())
			resp := map[string]string{
				"errorMessage": err.Error(),
			}
			w.WriteHeader(http.StatusInternalServerError)
			if err = json.NewEncoder(w).Encode(resp); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		log.Printf("Response: %s", str)
		w.Write([]byte(str))
	}
}

func GetRouter(module string, routes []Route) *chi.Mux {
	var prefix string

	if IsLocal() {
		prefix = "/" + module
	}

	mux := chi.NewRouter()
	mux.Use(loggerConfig)
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.SetHeader("Content-type", "application/json"))
	mux.Use(corsMiddleware)
	mux.Use(logRequestMiddleware)

	for idx, route := range routes {
		routes[idx].Handler = responseLoggerWrapper(routes[idx].Fn)
		mw := make([]func(http.Handler) http.Handler, 0)
		mw = append(mw,
			withAuthToken,
			middleware.WithValue(roles, route.Roles),
			middleware.WithValue(CtxEntitlement, route.Entitlement),
			appCheckMiddleware,
			checkRoles,
		)

		if slices.Contains(route.Roles, UserRoleAdmin) {
			mw = append(mw, auditLogMiddleware)
		}

		mw = append(mw, route.Middlewares...)
		mux.With(mw...).Method(route.Method, prefix+route.Route, route.Handler)
	}

	return mux
}

// MIDDLEWARES

func loggerConfig(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix)

		defer func() {
			log.SetPrefix("")
		}()

		next.ServeHTTP(w, r)
	})
}

func auditLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		obfuscatedBody := obfuscateFields(body)

		defer func() {
			createAuditLog(r, string(obfuscatedBody))
		}()

		// rewrite body to request since it is a stream
		r.Body = io.NopCloser(bytes.NewReader(body))
		next.ServeHTTP(w, r)
	})
}

func obfuscateFields(body []byte) []byte {
	var (
		forbiddenFields []string = []string{"password"}
		temp            map[string]any
	)

	if len(body) == 0 {
		return []byte{}
	}

	err := json.Unmarshal(body, &temp)
	if err != nil {
		log.Printf("error unmarshaling body fields: %s", err.Error())
		return []byte{}
	}

	for _, key := range forbiddenFields {
		if _, ok := temp[key]; ok {
			temp[key] = "**********"
		}
	}

	bb, err := json.Marshal(temp)
	if err != nil {
		log.Printf("error marshaling body fields: %s", err.Error())
		return []byte{}
	}

	return bb
}

func withAuthToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idToken := r.Header.Get("Authorization")
		authToken, err := GetAuthTokenFromIdToken(idToken)
		if err != nil {
			log.Printf("error extracting authToken: %s", err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), CtxAuthToken, authToken))

		next.ServeHTTP(w, r)
	})
}

type AuditLog struct {
	Payload  string         `bigquery:"payload"`
	Date     civil.DateTime `bigquery:"date"`
	UserUid  string         `bigquery:"userUid"`
	Method   string         `bigquery:"method"`
	Endpoint string         `bigquery:"endpoint"`
	Role     string         `bigquery:"role"`
}

func (a *AuditLog) SaveToBigQuery() error {
	if err := InsertRowsBigQuery(WoptaDataset, AuditsCollection, a); err != nil {
		return fmt.Errorf("cannot save the audit log: %v", err)
	}
	return nil
}

func createAuditLog(r *http.Request, payload string) {
	log.Println("saving audit trail...")
	audit, err := parseHttpRequest(r, payload)
	if err != nil {
		log.Printf("error creating audit log: %s", err.Error())
	}
	log.Printf("audit log: %v", audit)
	if err = audit.SaveToBigQuery(); err == nil {
		log.Printf("audit log saved!")
	}
}

func parseHttpRequest(r *http.Request, payload string) (AuditLog, error) {
	authToken := r.Context().Value(CtxAuthToken).(AuthToken)

	return AuditLog{
		Payload:  payload,
		Date:     civil.DateTimeOf(time.Now().UTC()),
		UserUid:  authToken.UserID,
		Method:   r.Method,
		Endpoint: r.RequestURI,
		Role:     authToken.Role,
	}, nil
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		options := cors.Options{
			AllowedHeaders: []string{"*"},
			AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodHead},
		}

		if r.Method == http.MethodOptions || IsLocal() {
			options.AllowedOrigins = []string{"*"}
		}

		if r.Method == http.MethodOptions {
			options.AllowedMethods = []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete}
			options.AllowCredentials = true
			options.MaxAge = 3600
		} else {
			options.AllowCredentials = false
		}

		c := cors.New(options)
		handler := c.Handler(next)
		handler.ServeHTTP(w, r)
	})
}

func appCheckMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		roles := ctx.Value(roles).([]string)

		if len(roles) == 0 || IsLocal() || slices.Contains(roles, UserRoleInternal) {
			next.ServeHTTP(w, r)
			return
		}

		app, err := firebase.NewApp(ctx, &firebase.Config{ProjectID: os.Getenv("GOOGLE_PROJECT_ID")})
		if err != nil {
			log.Printf("error initializing app: %s", err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		appCheck, err := app.AppCheck(ctx)
		if err != nil {
			log.Printf("error initializing app: %s\n", err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		appCheckToken, ok := r.Header[http.CanonicalHeaderKey("X-Firebase-AppCheck")]
		if !ok {
			log.Printf("error missing token")
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		_, err = appCheck.VerifyToken(appCheckToken[0])
		if err != nil {
			log.Printf("error invalid token: %s", err.Error())
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func checkRoles(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		roles := ctx.Value(roles).([]string)

		// TODO review me - no useful role to check entitlement
		if len(roles) == 0 || slices.Contains(roles, UserRoleInternal) || slices.Contains(roles, UserRoleAll) {
			next.ServeHTTP(w, r)
			return
		}

		idToken := strings.ReplaceAll(r.Header.Get("Authorization"), "Bearer ", "")
		if idToken == "" {
			log.Println("empty token")
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		token := ctx.Value(CtxAuthToken).(AuthToken)

		// token, err := VerifyUserIdToken(idToken)
		// if err != nil {
		// 	log.Printf("verify id token error: %s", err.Error())
		// 	http.Error(w, "unauthorized", http.StatusUnauthorized)
		// 	return
		// }

		// userRole := UserRoleCustomer
		// if role, ok := token.Claims["role"].(string); ok {
		// 	userRole = role
		// }

		// r = r.WithContext(context.WithValue(r.Context(), role, userRole))

		if !slices.Contains(roles, token.Role) {
			log.Printf("userRole '%s' not allowed", token.Role)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func logRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		obfuscatedBody := obfuscateFields(body)
		if len(obfuscatedBody) > 0 {
			log.Printf("Request: %s", string(obfuscatedBody))
		}

		// rewrite body to request since it is a stream
		r.Body = io.NopCloser(bytes.NewReader(body))
		next.ServeHTTP(w, r)
	})
}
