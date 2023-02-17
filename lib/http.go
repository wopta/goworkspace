package lib

import (
	"log"
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
