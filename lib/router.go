package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime/debug"
	"slices"
	"strings"
	"time"

	"cloud.google.com/go/civil"
	firebase "firebase.google.com/go/v4"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
	env "gitlab.dev.wopta.it/goworkspace/lib/environment"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
)

type Route struct {
	Route       string
	Method      string
	Handler     http.HandlerFunc
	Middlewares []func(http.Handler) http.Handler
	Roles       []string
}

func setupEnv(w http.ResponseWriter, r *http.Request) {
	uuid := uuid.NewString()
	w.Header().Set("Execution-Id", uuid)
	os.Setenv("Execution-Id", uuid)
	idToken := r.Header.Get("Authorization")
	if idToken != "" {
		authToken, _ := GetAuthTokenFromIdToken(idToken)
		os.Setenv("User", fmt.Sprintf("email:%v,role:%v,id:%v", authToken.Email, authToken.Role, authToken.UserID))
	}
}

func ResponseLoggerWrapper(handler func(w http.ResponseWriter, r *http.Request) (string, any, error)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		setupEnv(w, r)
		logRequestMiddleware(w, r)
		str, _, err := handler(w, r)

		if err != nil {

			log.Error(err)
			resp := map[string]string{
				"errorMessage": err.Error(),
			}
			w.WriteHeader(http.StatusInternalServerError)
			if err = json.NewEncoder(w).Encode(resp); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		if w.Header().Get("Content-type") == "application/json" {
			log.Printf("Response: %s", str)
		}
		w.Write([]byte(str))
	}
}
func GetRouter(module string, routes []Route) *chi.Mux {
	var prefix string

	if env.IsLocal() {
		prefix = "/" + module
	}

	mux := chi.NewRouter()
	mux.Use(loggerConfig)
	mux.Use(recoverLogger)
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(corsMiddleware)
	mux.Use(middleware.SetHeader("Content-type", "application/json"))
	for _, route := range routes {
		mw := make([]func(http.Handler) http.Handler, 0)
		mw = append(mw,
			middleware.WithValue("roles", route.Roles),
			appCheckMiddleware,
			//		checkEntitlement,
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

func recoverLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				if rvr == http.ErrAbortHandler {
					// we don't recover http.ErrAbortHandler so the response
					// to the client is aborted, this should not be logged
					panic(rvr)
				}
				log.ErrorF("Panic %v", rvr)
				stackByte := debug.Stack()
				stack := strings.Split(string(stackByte), "\n")

				lines := []string{}

				for i := len(stack) - 1; i > 0; i-- {
					lines = append(lines, stack[i])
					if strings.HasPrefix(stack[i], "panic(") {
						lines = lines[0 : len(lines)-2] // remove boilerplate
						break
					}
				}
				slices.Reverse(lines)
				log.Log().CustomLog(strings.Join(lines, "\n"), log.PANIC)
				if r.Header.Get("Connection") != "Upgrade" {
					w.WriteHeader(http.StatusInternalServerError)
				}
			}

		}()
		next.ServeHTTP(w, r)
	})
}
func loggerConfig(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.ResetPrefix()
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
		log.ErrorF("error unmarshaling body fields: %s", err.Error())
		return []byte{}
	}

	for _, key := range forbiddenFields {
		if _, ok := temp[key]; ok {
			temp[key] = "**********"
		}
	}

	bb, err := json.Marshal(temp)
	if err != nil {
		log.ErrorF("error marshaling body fields: %s", err.Error())
		return []byte{}
	}

	return bb
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
		log.ErrorF("error creating audit log: %s", err.Error())
	}
	log.Printf("audit log: %v", audit)
	if err = audit.SaveToBigQuery(); err == nil {
		log.Printf("audit log saved!")
	}
}

func parseHttpRequest(r *http.Request, payload string) (AuditLog, error) {
	idToken := r.Header.Get("Authorization")
	authToken, err := GetAuthTokenFromIdToken(idToken)
	if err != nil {
		return AuditLog{}, fmt.Errorf("cannot retrieve the user's authorization token: %v", err)
	}

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

		if r.Method == http.MethodOptions || env.IsLocal() {
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
		roles := ctx.Value("roles").([]string)

		if len(roles) == 0 || env.IsLocal() || slices.Contains(roles, UserRoleInternal) {
			next.ServeHTTP(w, r)
			return
		}

		app, err := firebase.NewApp(ctx, &firebase.Config{ProjectID: os.Getenv("GOOGLE_PROJECT_ID")})
		if err != nil {
			log.ErrorF("error initializing app: %s", err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		appCheck, err := app.AppCheck(ctx)
		if err != nil {
			log.ErrorF("error initializing app: %s\n", err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		appCheckToken, ok := r.Header[http.CanonicalHeaderKey("X-Firebase-AppCheck")]
		if !ok {
			log.ErrorF("error missing token")
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		_, err = appCheck.VerifyToken(appCheckToken[0])
		if err != nil {
			log.ErrorF("error invalid token: %s", err.Error())
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func checkEntitlement(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		roles := ctx.Value("roles").([]string)

		if len(roles) == 0 || slices.Contains(roles, UserRoleInternal) || slices.Contains(roles, UserRoleAll) {
			next.ServeHTTP(w, r)
			return
		}

		idToken := strings.ReplaceAll(r.Header.Get("Authorization"), "Bearer ", "")
		if idToken == "" {
			log.ErrorF("empty token")
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		token, err := VerifyUserIdToken(idToken)
		if err != nil {
			log.ErrorF("verify id token error: %s", err.Error())
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		userRole := UserRoleCustomer
		if role, ok := token.Claims["role"].(string); ok {
			userRole = role
		}

		if !slices.Contains(roles, userRole) {
			log.WarningF("userRole '%s' not allowed", userRole)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func logRequestMiddleware(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	obfuscatedBody := obfuscateFields(body)
	if len(obfuscatedBody) > 0 {
		log.Printf("Request: %s", string(obfuscatedBody))
	}

	// rewrite body to request since it is a stream
	r.Body = io.NopCloser(bytes.NewReader(body))
}
