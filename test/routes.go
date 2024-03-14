package test

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
)

var mux = newMux()

func init() {
	log.Println("INIT Test")
	funcframework.RegisterHTTPFunction("Test", mux.ServeHTTP)
}

func newMux() *http.ServeMux {
	log.Println("Creating Test mux...")
	mux := http.NewServeMux()
	mux.HandleFunc("/test/test1", test1)
	mux.HandleFunc("/test/test2/param", test2)
	return mux
}

func test1(w http.ResponseWriter, r *http.Request) {
	log.Println("test1 handler!")
	log.Printf("Request: %s", r.RequestURI)
}

func test2(w http.ResponseWriter, r *http.Request) {
	log.Println("test2 handler!")
	log.Printf("Request: %s", r.RequestURI)
}
