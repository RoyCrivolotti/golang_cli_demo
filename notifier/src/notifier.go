package src

import "refurbedchallenge/notifier/constants"

//Library
//Create HTTP notifications client -> Configurable URL to which to propagate notifications
//Client gets message and sends them in the body of the POST request to the configured URL in a NON BLOCKER manner
//Should the Client get an error when sending the message, the caller should be able to handle the error

//Executable
//Consume library above
//Send new messages every interval (configurable)
//The exe reads standard input and every new line is a message to send to the library for propagation
//Consider shutdown by SIGINT

type INotifier interface {
	Notify(message string) chan constants.NotificationError
}

type notifier struct {
	url    string
	client IHttpClient
}

func NewNotifier(url string) INotifier {
	return &notifier{
		url:    url,
		client: NewHttpClient(),
	}
}

func (n *notifier) Notify(message string) chan constants.NotificationError {
	channel := make(chan constants.NotificationError)
	go func() {
		if _, err := n.client.Post(n.url, message); err != nil {
			channel <- constants.NotificationError{Error: err, Message: message}
		} else {
			channel <- constants.NotificationError{}
		}

		return
	}()

	return channel
}