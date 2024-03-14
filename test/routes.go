package test

import (
	"log"
	"net/http"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

// for local testing only
func init() {
	log.Println("INIT Test")
	functions.HTTP("Test", Test)
}

func newMux() *http.ServeMux {
	prefix := ""

	if os.Getenv("env") == "local" {
		prefix = "/test"
	}

	mux := http.NewServeMux()
	mux.HandleFunc(prefix+"/test1", test1)
	mux.HandleFunc(prefix+"/test2/param", test2)
	return mux
}

func Test(w http.ResponseWriter, r *http.Request) {
	mux := newMux()
	mux.ServeHTTP(w, r)
}

func test1(w http.ResponseWriter, r *http.Request) {
	log.Println("test1 handler!")
	log.Printf("Request: %s", r.RequestURI)
	w.Write([]byte(`{}`))
}

func test2(w http.ResponseWriter, r *http.Request) {
	log.Println("test2 handler!")
	log.Printf("Request: %s", r.RequestURI)
	w.Header().Add("Content-type", "application/json")
	w.Write([]byte(`{"success":true}`))
}
