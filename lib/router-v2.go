package lib

import (
	"io"
	"log"
	"net/http"
	"strings"
)

type RouteDataV2 struct {
	Routes []Route
}
type RouteV2[Req any, Resp any] struct {
	Route    string
	Method   string
	Handler  func(http.ResponseWriter, *http.Request) (string, interface{}, error)
	Roles    []string
	ReqModel *Req
	ResModel *Resp
}

func (router RouteDataV2) Router(w http.ResponseWriter, r *http.Request) {
	log.Println("Router")
	var result string
	var isFound bool
	var route string
	var e error

	//reqUri := r.RequestURI
	log.Printf("[RouterV2] request URI: %s", r.RequestURI)

	for _, v := range router.Routes {
		route = v.Route
		if strings.Contains(route, "?") {
			//i := strings.Index(route, ":")
			routeSplit := strings.Split(route, ":")
			base := routeSplit[0]
			key := routeSplit[1]
			reqUris := strings.Split(r.RequestURI, "/")
			value := reqUris[len(reqUris)-1]

			r.Header.Add(key, value)
			route = base
		}
		if strings.Contains(route, ":") {
			//i := strings.Index(route, ":")
			routeSplit := strings.Split(route, ":")
			base := routeSplit[0]
			key := routeSplit[1]
			reqUris := strings.Split(r.RequestURI, "/")
			value := reqUris[len(reqUris)-1]

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
