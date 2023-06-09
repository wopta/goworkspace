package lib

import (
	"encoding/json"
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
func getIP(req *http.Request) net.IP {

	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		log.Println(err)

	}

	userIP := net.ParseIP(ip)
	return userIP
}

func CheckPayload[T any](body []byte, payload *T, fields []string) error {
	if len(body) == 0 {
		return fmt.Errorf("Missing payload")
	}

	var bodyJson map[string]interface{}

	err := json.Unmarshal(body, &bodyJson)
	CheckError(err)

	for _, param := range fields {
		if _, ok := bodyJson[param]; !ok {
			return fmt.Errorf("Missing paramenter %s", param)
		}
	}

	return json.Unmarshal(body, payload)
}
