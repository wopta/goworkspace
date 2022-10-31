package lib

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func Router(w http.ResponseWriter, r *http.Request, route map[string]func(http.ResponseWriter, *http.Request) string) {
	var result string
	log.Println(r.RequestURI)
	for k, v := range route {
		var isFound bool
		if strings.Contains(r.RequestURI, k) {

			log.Println("found")
			result = v(w, r)
			isFound = true
		}
		if !isFound {
			fmt.Fprintf(w, `{"message":" select correct path "}`)
		}

	}
	fmt.Fprintf(w, result)
}
