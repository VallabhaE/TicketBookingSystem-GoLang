package api

import (
	"bytes"
	"net/http"
)

func __internalGet(url string) (*http.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func __internalPost(url string, body []byte) (*string, error) {
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return &response.Status, nil
}

func Get(url string) (*http.Response, error) {

	return __internalGet(url)
}

func Post(url, body string) (string, error) {

	data, err := __internalPost(url, []byte(body))
	
	return *data, err
}
