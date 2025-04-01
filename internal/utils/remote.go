package utils

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

var (
	errCouldntBuildHTTPRequest       = errors.New("couldn't build HTTP request")
	errCouldntSendHTTPRequest        = errors.New("couldn't send HTTP request")
	errRemoteServerSentNonOkResponse = errors.New("remote server responded with non success HTTP code")
	errCouldntReadHTTPResponse       = errors.New("couldn't read HTTP response")
)

func GetHTTPResponse(url string, headers [][2]string) ([]byte, error) {
	var zero []byte
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return zero, errCouldntBuildHTTPRequest
	}
	for _, header := range headers {
		req.Header.Set(header[0], header[1])
	}

	resp, err := client.Do(req)
	if err != nil {
		return zero, errCouldntSendHTTPRequest
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return zero, fmt.Errorf("%w: %d", errRemoteServerSentNonOkResponse, resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return zero, fmt.Errorf("%w: %s", errCouldntReadHTTPResponse, err.Error())
	}

	return bodyBytes, nil
}
