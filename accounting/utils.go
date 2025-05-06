package accounting

import (
	"bytes"
	"net/http"
)

func HttpFileToByte( url string) (buffer bytes.Buffer, err error) {
	// Create the file
	var buf bytes.Buffer
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return buf, err
	}
	
	if err := resp.Write(&buf); err != nil {
		return buf, err
	}
	defer resp.Body.Close()
	// Writer the body to file
	return buf, err
}
