package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var ALL_FUNCTIONS []string = []string{
	"Auth",
	"Broker",
	"Callback",
	"Claim",
	"Companydata",
	"Enrich",
	"Form",
	"Mail",
	"Mga",
	"Network",
	"Partnership",
	"Payment",
	"Policy",
	"Question",
	"Quote",
	"Renew",
	"Reserved",
	"Rules",
	"Sellable",
	"Transaction",
	"User",
	"Inclusive",
	"Document",
}

func GetLocalIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddress := conn.LocalAddr().(*net.UDPAddr)

	return localAddress.IP
}
func main() {
	var bind map[string]string = make(map[string]string)
	var err error
	var start int
	if len(os.Args) < 2 {
		start = 8080
	} else {
		start, err = strconv.Atoi(os.Args[1])
		if err != nil {
			panic(err)
		}
	}
	var serversErrors []string
	for i, function := range ALL_FUNCTIONS {
		cmd := exec.Command(os.Getenv("GOWORKSPACE") + "/../bin/api")
		cmd.Env = append(cmd.Env, "FUNCTION_TARGET="+function, " ")
		cmd.Env = append(cmd.Env, "PORT="+fmt.Sprint(start+i+1))
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		bind[function] = fmt.Sprint(start + i + 1)
		if e := cmd.Start(); e != nil {
			serversErrors = append(serversErrors, e.Error())
		}
		defer cmd.Process.Kill()
	}
	time.Sleep(time.Second)
	fmt.Printf("\nServer started %s:%v\n", GetLocalIP(), start)
	if len(serversErrors) > 0 {
		fmt.Printf("Errors: %v\n", serversErrors)
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization,Content-Type,X-Firebase-Appcheck")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Max-Age", "600")
			w.WriteHeader(http.StatusNoContent)
			return
		}
		url := strings.Split(r.URL.String(), "/")[1]
		funcName := strings.ToUpper(string(url[0])) + url[1:]
		url = "http://localhost:" + bind[funcName] + r.URL.String()
		nR, _ := http.NewRequest(r.Method, url, r.Body)
		nR.Header = r.Header
		nR.Body = r.Body
		nR.Host = r.Host
		client := &http.Client{
			Timeout: time.Second * 10, // Timeout each requests
		}
		fmt.Printf("Calling function %s...\n\n", funcName)
		resp, e := client.Do(nR)
		if e != nil {
			w.WriteHeader(resp.StatusCode)
			w.Write([]byte(e.Error()))
			return
		}
		defer resp.Body.Close()
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Content-Type", r.Header.Get("Content-Type"))
		w.WriteHeader(resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		w.Write(body)
		fmt.Printf("\nServer is running at %s:%v\n", GetLocalIP(), start)

	})

	http.ListenAndServe(":"+fmt.Sprint(start), nil)

}
