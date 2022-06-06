package services

import (
	"fmt"
	"io"
	"refurbedchallenge/notifier/constants"
	"refurbedchallenge/notifier/src"
	"time"
)

type IMessageService interface {
	ProcessMessages(w io.Writer, lines []string, interval int64, debug bool) []constants.NotificationError
}

type service struct {
	client src.INotifierClient
}

func NewMessageService(notifier src.INotifierClient) IMessageService {
	return &service{
		client: notifier,
	}
}

//ProcessMessages receives a slice of messages to process and the interval of time to pause between each one.
//For every message that fails to be processed, an error message is printed
func (s *service) ProcessMessages(w io.Writer, lines []string, interval int64, debug bool) []constants.NotificationError {
	ticker := time.NewTicker(time.Duration(interval) * time.Millisecond)
	quit := make(chan struct{})

	var errors []constants.NotificationError

Loop:
	for {
		select {
		case <-ticker.C:
			if len(lines) == 0 {
				break Loop
			}

			//Getting first element from lines slice
			line := lines[0]
			//Popping the element
			lines = lines[1:]
			//Letting the user know what element is being processed currently
			_, _ = fmt.Fprintf(w, "Current line: '%s' - Waiting %d milliseconds\n", line, interval)

			//Messages are sent to be processed every n milliseconds, regardless of how much an HTTP request takes
			c := make(chan constants.NotificationError)
			go s.client.Notify(line, c)

			err := <-c

			//Printing lines that encountered an error for the user to manually handle them
			if err != (constants.NotificationError{}) {
				errors = append(errors, err)
				_, _ = fmt.Fprintf(w, "%+v\n", err)
			}

			//channels = append(channels, c)
		case <-quit:
			ticker.Stop()
			break Loop
		}
	}

	return errors
}
