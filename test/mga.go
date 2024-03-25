package test

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	// "reflect"
	"slices"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/wopta/goworkspace/mga"
	"github.com/wopta/goworkspace/models"
)

var mgaRoutes []Route = []Route{
	{
		Route:   "/products/v1",
		Handler: handlerWrapper(mga.GetProductsListByChannelFx),
		Method:  http.MethodGet,
		Roles:   []string{models.UserRoleAdmin},
	},
	{
		Route:   "/products/v1",
		Handler: handlerWrapper(mga.GetProductByChannelFx),
		Method:  http.MethodPost,
		Roles:   []string{models.UserRoleAll},
		// Middlewares: nil,
		// RequestType: &ProductRequest{},
	},
	{
		Route:   "/network/node/v1/{uid}",
		Handler: handlerWrapper(mga.GetNetworkNodeByUidFx),
		Method:  http.MethodGet,
		Roles:   []string{models.UserRoleAll},
	},
	{
		Route:   "/network/node/v1",
		Handler: handlerWrapper(mga.CreateNetworkNodeFx),
		Method:  http.MethodPost,
		Roles:   []string{models.UserRoleAdmin, models.UserRoleManager},
		// RequestType: &NodeRequest{},
	},
	{
		Route:   "/network/node/v1",
		Handler: handlerWrapper(mga.UpdateNetworkNodeFx),
		Method:  http.MethodPut,
		Roles:   []string{models.UserRoleAdmin, models.UserRoleManager},
	},
	{
		Route:   "/network/nodes/v1",
		Handler: handlerWrapper(mga.GetAllNetworkNodesFx),
		Method:  http.MethodGet,
		Roles:   []string{models.UserRoleAdmin, models.UserRoleManager},
	},
	{
		Route:   "/network/node/v1/:uid",
		Handler: handlerWrapper(mga.DeleteNetworkNodeFx),
		Method:  http.MethodDelete,
		Roles:   []string{models.UserRoleAdmin, models.UserRoleManager},
	},
	{
		Route:   "/network/invite/v1/create",
		Handler: handlerWrapper(mga.CreateNetworkNodeInviteFx),
		Method:  http.MethodPost,
		Roles:   []string{models.UserRoleAdmin, models.UserRoleManager},
	},
	{
		Route:   "/network/invite/v1/consume",
		Handler: handlerWrapper(mga.ConsumeNetworkNodeInviteFx),
		Method:  http.MethodPost,
		Roles:   []string{models.UserRoleAll},
	},
	{
		Route:   "/warrants/v1",
		Handler: handlerWrapper(mga.GetWarrantsFx),
		Method:  http.MethodGet,
		Roles:   []string{models.UserRoleAdmin, models.UserRoleManager},
	},
	{
		Route:   "/warrant/v1",
		Handler: handlerWrapper(mga.CreateWarrantFx),
		Method:  http.MethodPut,
		Roles:   []string{models.UserRoleAdmin, models.UserRoleManager},
	},
	{
		Route:   "/policy/v1",
		Handler: handlerWrapper(mga.ModifyPolicyFx),
		Method:  http.MethodPatch,
		Roles:   []string{models.UserRoleAdmin},
	},
}

func Mga(w http.ResponseWriter, r *http.Request) {
	router := getModuleRouter("test", mgaRoutes)
	router.ServeHTTP(w, r)
}

// lib
type Route struct {
	Route       string
	Method      string
	Handler     http.HandlerFunc
	Middlewares []func(http.Handler) http.Handler
	Roles       []string
	// RequestType RouteRequest
}

func handlerWrapper(handler func(w http.ResponseWriter, r *http.Request) (string, any, error)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		str, _, err := handler(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte(str))
	}
}

func getModuleRouter(module string, routes []Route) *chi.Mux {
	prefix := "/"

	if os.Getenv("env") == "local" {
		prefix += module
	}

	mux := chi.NewRouter()
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.SetHeader("Content-type", "application/json"))
	mux.Use(CorsMiddleware)
	mux.Use(LogRequestObMiddleware)

	for _, route := range routes {
		mw := make([]func(http.Handler) http.Handler, 0)
		// if route.RequestType != nil {
		// 	mw = append(mw,
		// 		middleware.WithValue("requestType", route.RequestType),
		// 		LogRequestMiddleware,
		// 	)
		// }
		mw = append(mw,
			middleware.WithValue("roles", route.Roles),
			AppCheckMiddleware,
			CheckEntitlement,
		)

		if slices.Contains(route.Roles, models.UserRoleAdmin) {
			mw = append(mw, AuditLogMiddleware)
		}

		mw = append(mw, route.Middlewares...)
		mux.With(mw...).Method(route.Method, prefix+route.Route, route.Handler)
	}

	return mux
}

func AuditLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		var t map[string]any

		json.Unmarshal(body, &t)
		for _, key := range []string{"password", "pwd"} {
			if _, ok := t[key]; ok {
				t[key] = "**********"
			}
		}
		bb, _ := json.Marshal(t)

		defer func() {
			models.CreateAuditLog(r, string(bb))
		}()

		// rewrite body to request since it is a stream
		r.Body = io.NopCloser(bytes.NewReader(body))
		next.ServeHTTP(w, r)
	})
}

// func AuditLogMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		t2 := r.Context().Value("requestType")

// 		if t2 != nil {
// 			t := t2.(RouteRequest)
// 			req, temp := t.Parse(r)

// 			request := RemoveObfuscated(temp)
// 			defer func() {
// 				models.CreateAuditLog(r, string(request))
// 			}()

// 			// rewrite body to request since it is a stream
// 			r.Body = io.NopCloser(bytes.NewReader(req))
// 		} else {
// 			defer func() {
// 				models.CreateAuditLog(r, "")
// 			}()
// 		}

// 		next.ServeHTTP(w, r)
// 	})
// }

// type RouteRequest interface {
// 	Parse(*http.Request) ([]byte, any)
// }

// type ProductRequest struct {
// 	ProductName string `json:"name" obfuscate:"true"`
// 	CompanyName string `json:"company"` // DEPRECATED
// 	Version     string `json:"version"` // DEPRECATED
// }

// func (r *ProductRequest) Parse(req *http.Request) ([]byte, any) {
// 	bytes, _ := io.ReadAll(req.Body)
// 	defer req.Body.Close()

// 	var value ProductRequest

// 	json.Unmarshal(bytes, &value)

// 	return bytes, value
// }

// type NodeRequest struct {
// 	Code string `json:"code"`
// 	Mail string `json:"mail" obfuscate:"true"`
// }

// func (r *NodeRequest) Parse(req *http.Request) ([]byte, any) {
// 	bytes, _ := io.ReadAll(req.Body)
// 	defer req.Body.Close()

// 	var value NodeRequest

// 	json.Unmarshal(bytes, &value)

// 	return bytes, value
// }

// func LogRequestMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		t := r.Context().Value("requestType").(RouteRequest)

// 		req, temp := t.Parse(r)

// 		request := RemoveObfuscated(temp)
// 		log.Printf("Request: %s", string(request))

// 		// rewrite body to request since it is a stream
// 		r.Body = io.NopCloser(bytes.NewReader(req))
// 		next.ServeHTTP(w, r)
// 	})
// }

// func RemoveObfuscated[T any](data T) []byte {
// 	val := reflect.ValueOf(&data).Elem()
// 	overridableCopy := reflect.New(val.Elem().Type()).Elem()

// 	overridableCopy.Set(val.Elem())
// 	for i := 0; i < overridableCopy.NumField(); i++ {
// 		field := overridableCopy.Field(i)
// 		fieldType := overridableCopy.Type().Field(i)

// 		tag := fieldType.Tag.Get("obfuscate")

// 		if tag != "" {
// 			field.SetString("*******")
// 		}
// 	}
// 	val.Set(overridableCopy)

// 	bytes, _ := json.Marshal(data)
// 	return bytes
// }

func LogRequestObMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		var t map[string]any

		json.Unmarshal(body, &t)
		for _, key := range []string{"password", "pwd"} {
			if _, ok := t[key]; ok {
				t[key] = "**********"
			}
		}
		bb, _ := json.Marshal(t)
		log.Printf("Request: %v", string(bb))

		// rewrite body to request since it is a stream
		r.Body = io.NopCloser(bytes.NewReader(body))
		next.ServeHTTP(w, r)
	})
}
