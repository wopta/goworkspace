package lib

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

func RetryDo(req *http.Request, retry int) (*http.Response, error) {
	var (
		resp *http.Response
		e    error
	)

	client := http.Client{
		Timeout: time.Millisecond * 10,
	}

	for i := 1; i <= retry; i++ {
		resp, e = client.Do(req)
		if e != nil {
			log.Printf("error sending the first time: %v\n", e)
			time.Sleep(5000)
		} else {
			e = nil

		}

	}
	return resp, e
}
func getIP(w http.ResponseWriter, req *http.Request) net.IP {

	ip, port, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		//return nil, fmt.Errorf("userip: %q is not IP:port", req.RemoteAddr)

		fmt.Fprintf(w, "userip: %q is not IP:port", req.RemoteAddr)
	}

	userIP := net.ParseIP(ip)
	if userIP == nil {
		//return nil, fmt.Errorf("userip: %q is not IP:port", req.RemoteAddr)
		fmt.Fprintf(w, "userip: %q is not IP:port", req.RemoteAddr)
		return userIP
	}

	// This will only be defined when site is accessed via non-anonymous proxy
	// and takes precedence over RemoteAddr
	// Header.Get is case-insensitive
	forward := req.Header.Get("X-Forwarded-For")

	fmt.Fprintf(w, "<p>IP: %s</p>", ip)
	fmt.Fprintf(w, "<p>Port: %s</p>", port)
	fmt.Fprintf(w, "<p>Forwarded for: %s</p>", forward)
	return userIP
}
