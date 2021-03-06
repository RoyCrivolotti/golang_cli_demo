package src

import (
	"clidemo/notifier/clients"
	"clidemo/notifier/constants"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

//go:generate mockgen -source=./notifier.go -destination=./mock/notifier_mock.go

type INotifierClient interface {
	NotifySync(message string) constants.NotificationError
	NotifyChannel(message string, c chan constants.NotificationError)
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

//NotifySync sends the message to the configured URL via POST message
func (n *notifier) NotifySync(message string) constants.NotificationError {
	if res, err := n.client.Post(n.url, message); err != nil {
		return constants.NotificationError{Error: err, Message: message}
	} else if res.StatusCode >= 400 && res.StatusCode <= 599 {
		errMsg := handleHttpError(res)
		return constants.NotificationError{Error: errMsg, Message: message}
	} else {
		return constants.NotificationError{}
	}
}

//NotifyChannel sends the message to the configured URL via POST message and the channel receives the error, should one be encountered
func (n *notifier) NotifyChannel(message string, c chan constants.NotificationError) {
	go func() {
		if res, err := n.client.Post(n.url, message); err != nil {
			c <- constants.NotificationError{Error: err, Message: message}
		} else if res.StatusCode >= 400 && res.StatusCode <= 599 {
			errMsg := handleHttpError(res)
			c <- constants.NotificationError{Error: errMsg, Message: message}
		} else {
			c <- constants.NotificationError{}
		}
	}()
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
