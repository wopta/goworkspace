package lib

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net"
	"net/http"
	"time"
)

func RetryDo(req *http.Request, retry int) (*http.Response, error) {
	const (
		maxRetry = 5
	)
	var (
		resp *http.Response
		err  error
	)

	client := http.Client{
		Timeout: time.Second * 10,
	}

	for i := 0; i < retry && i < maxRetry; i++ {
		log.Printf("[RetryDo] sending request %d at time %s", i+1, time.Now().UTC())
		resp, err = client.Do(req)
		if err == nil || (i == maxRetry-1) {
			break
		}
		log.Printf("[RetryDo] error: %s", err.Error())
		time.Sleep(time.Duration(500*math.Pow(2, float64(i))) * time.Millisecond)
	}

	return resp, err
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
