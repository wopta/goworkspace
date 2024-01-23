package lib

import (
	"context"
	"log"
	"net/http"
	"strings"
)

type Param struct {
	Key   string
	Value string
}
type Params []Param
type ParamsKey struct{}
type QueryParamsKey struct{}

func (router RouteData) RouterV3(w http.ResponseWriter, req *http.Request) {
	var (
		err        error
		ctx        context.Context
		response   string
		routeFound bool
	)

	for _, r := range router.Routes {
		if r.Method != req.Method {
			continue
		}

		route := r.Route
		splitQuery := strings.Split(req.RequestURI, "?")
		requestURIPaths := splitTrim(splitQuery[0], "/")
		queryPaths := splitTrim(splitQuery[1], "&")
		routePathsToMatch := requestURIPaths[1:]
		paths := splitTrim(route, "/")

		if len(paths) != len(routePathsToMatch) {
			continue
		}

		routeFound = foundRouteMatch(paths, routePathsToMatch)
		if !routeFound {
			continue
		}

		params := make(Params, 0)
		for index, path := range paths {
			if strings.Contains(path, ":") {
				params = append(params, Param{path[1:], routePathsToMatch[index]})
			}
		}

		if len(queryPaths) > 0 {
			queries := make(Params, 0)
			for _, query := range queryPaths {
				tmp := strings.Split(query, "=")
				queries = append(queries, Param{tmp[0], tmp[1]})
			}
			ctx = context.WithValue(req.Context(), QueryParamsKey{}, queries)
			req = req.WithContext(ctx)
		}

		if len(params) > 0 {
			ctx = context.WithValue(req.Context(), ParamsKey{}, params)
			req = req.WithContext(ctx)
		}

		response, _, err = r.Handler(w, req)
	}

	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), 500)
	}

	if !routeFound {
		http.NotFound(w, req)
	}

	_, err = w.Write([]byte(response))
	if err != nil {
		log.Printf("error writing response: %s", err.Error())
	}
}

func splitTrim(str string, separator string) []string {
	res := make([]string, 0)
	for _, s := range strings.Split(str, separator) {
		if s != "" {
			res = append(res, s)
		}
	}
	return res
}

func foundRouteMatch(baseRoutePaths, requestRoutePaths []string) bool {
	foundMatch := false
	for idx, p := range baseRoutePaths {
		if strings.Contains(p, ":") {
			continue
		}
		foundMatch = p == requestRoutePaths[idx]
	}
	return foundMatch
}
