package lib

import (
	"io"
	"log"
	"net/http"
	"strings"
)

type RouteData struct {
	Routes []Route
}
type Route struct {
	Route    string
	Method   string
	Handler  func(http.ResponseWriter, *http.Request) (string, interface{}, error)
	Roles    []string
	ReqModel *interface{}
	ResModel *interface{}
}

func (router RouteData) Router(w http.ResponseWriter, r *http.Request) {
	var (
		result  string
		isFound bool
		route   string
		e       error
	)

	log.Println("Router")
	log.Printf("[Router] request URI: %s", r.RequestURI)

	for _, v := range router.Routes {
		route = v.Route
		if strings.Contains(route, "?") {
			log.Println("[Router] splitting route for query params")
			routeSplit := strings.Split(route, ":")
			base := routeSplit[0]
			key := routeSplit[1]
			reqUris := strings.Split(r.RequestURI, "/")
			value := reqUris[len(reqUris)-1]
			r.Header.Add(key, value)
			route = base
		}
		if strings.Contains(route, ":") {
			log.Println("[Router] splitting route for dynamic params")
			routeSplit := strings.Split(route, ":")
			base := routeSplit[0]
			key := routeSplit[1]
			reqUris := strings.Split(r.RequestURI, "/")
			value := reqUris[len(reqUris)-1]
			value = strings.Split(value, "?")[0]
			r.Header.Add(key, value)
			route = base
		}
		if strings.Contains(r.RequestURI, route) && v.Method == r.Method {
			log.Println("found")
			log.Println(r.RequestURI)
			result, _, e = VerifyAuthorization(v.Handler, v.Roles...)(w, r)
			isFound = true
			break
		}

	}
	if e != nil {
		log.Println("Router error")
		log.Println(e.Error())
		http.Error(w, e.Error(), 500)
	}
	if !isFound {
		log.Println("Router not found")
		http.NotFound(w, r)
	}

	reader := strings.NewReader(result)
	io.Copy(w, reader)
}
