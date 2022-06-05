package src

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"refurbedchallenge/notifier/clients"
	"refurbedchallenge/notifier/constants"
)

//go:generate mockgen -source=./notifier.go -destination=./mock/notifier_mock.go

type INotifierClient interface {
	Notify(message string) chan constants.NotificationError
}

type notifier struct {
	url    string
	client clients.IHttpClient
}

func NewNotifierClient(url string) INotifierClient {
	return &notifier{
		url:    url,
		client: clients.NewHttpClient(),
	}
}

//Notify sends the message to the configured URL via POST message, using channels to avoid blocking operations
func (n *notifier) Notify(message string) chan constants.NotificationError {
	channel := make(chan constants.NotificationError)
	go func() {
		if res, err := n.client.Post(n.url, message); err != nil {
			channel <- constants.NotificationError{Error: err, Message: message}
		} else if res.StatusCode >= 400 && res.StatusCode <= 599 {
			errMsg := handleHttpError(res)
			channel <- constants.NotificationError{Error: errMsg, Message: message}
		} else {
			channel <- constants.NotificationError{}
		}

		return
	}()

	return channel
}

func handleHttpError(res *http.Response) error {
	resBody := new(map[string]interface{})

	//Return the response's error's body if possible, otherwise return the low level error
	err := json.NewDecoder(res.Body).Decode(resBody)
	if err != nil {
		return errors.New(fmt.Sprintf("Unexpected error: Status code %d", res.StatusCode))
	}

	b, err := json.Marshal(resBody)
	if err != nil {
		return errors.New(fmt.Sprintf("Unexpected error: Status code: %d - %s", res.StatusCode, err.Error()))
	} else {
		return errors.New(string(b))
	}
}
