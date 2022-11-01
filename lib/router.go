package lib

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type RouteData struct {
	Routes []Route
}
type Route struct {
	Route   string
	Hendler func(http.ResponseWriter, *http.Request) (string, interface{})
}

func (router RouteData) Router(w http.ResponseWriter, r *http.Request) {
	var result string
	var isFound bool
	log.Println(r.RequestURI)
	for _, v := range router.Routes {

		if strings.Contains(r.RequestURI, v.Route) {

			log.Println("found")
			result, _ = v.Hendler(w, r)
			isFound = true
		}

	}
	if !isFound {
		log.Println(" not found")
		fmt.Fprintf(w, `{"message":" select correct path "}`)
	}
	//io.Copy(w, result)
	fmt.Fprintf(w, result)
}
