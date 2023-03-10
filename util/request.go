package util

import (
	"io"
	"log"
	"net/http"
	"strings"
)

func MyRequest(method string, url string, headers []string, data string) ([]byte, error) {
	req, err := http.NewRequest(method, url, strings.NewReader(data))
	if checkRequestError(err) == true {
		return nil, err
	}

	for _, header := range headers {
		parts := strings.SplitN(header, ":", 2)
		if len(parts) == 2 {
			req.Header.Set(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if checkRequestError(err) == true {
		return nil, err
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Printf("MyGET Error closing response body: %s", err.Error())
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if checkRequestError(err) == true {
		return nil, err
	}

	return body, nil
}

func checkRequestError(err error) bool {
	if err != nil {
		//log.Printf("Request请求出错: %v", err)
		return true
	}
	return false
}
