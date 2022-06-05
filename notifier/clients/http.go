package clients

import (
	"bytes"
	"encoding/json"
	"net/http"
)

//go:generate mockgen -source=./http.go -destination=./mock/http_mock.go

type IHttpClient interface {
	Post(url string, data interface{}) (*http.Response, error)
}

type httpClient struct {
}

func NewHttpClient() IHttpClient {
	return &httpClient{}
}

//Post sends the data to the url and returns the response's body as a byte array and an error, should one occur
func (h *httpClient) Post(url string, data interface{}) (*http.Response, error) {
	requestBody := map[string]interface{}{
		"message": data,
	}

	buffer := new(bytes.Buffer)
	if err := json.NewEncoder(buffer).Encode(requestBody); err != nil {
		return nil, err
	}

	response, err := http.Post(url, "application/json", buffer)

	if err != nil {
		return nil, err
	}

	return response, nil
}
