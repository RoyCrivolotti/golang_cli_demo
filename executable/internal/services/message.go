package services

import (
	"clidemo/notifier/constants"
	"clidemo/notifier/src"
	"fmt"
	"io"
	"sync"
	"time"
)

type IMessageService interface {
	ProcessMessagesChannel(w io.Writer, lines []string, interval int64) []constants.NotificationError
	ProcessMessagesWaitGroup(w io.Writer, lines []string, interval int64) []constants.NotificationError
}

type service struct {
	client src.INotifierClient
}

func NewMessageService(notifier src.INotifierClient) IMessageService {
	return &service{
		client: notifier,
	}
}

//ProcessMessagesWaitGroup receives a slice of messages to process and the interval of time to pause between each one.
//For every message that fails to be processed, an error message is printed. WaitGroups are used for a simpler implementation.
func (s *service) ProcessMessagesWaitGroup(w io.Writer, lines []string, interval int64) []constants.NotificationError {
	var errors []constants.NotificationError

	var wg sync.WaitGroup

	// Increment the wait group counter
	wg.Add(len(lines))

	for i, line := range lines {
		go func(i int, line string) {
			// Decrement the counter when the go routine completes
			defer wg.Done()

			//Letting the user know what element is being processed currently
			_, _ = fmt.Fprintf(w, "Current line: '%s' - Waiting %d milliseconds\n", line, interval)

			//Messages are sent to be processed every n milliseconds, regardless of how much an HTTP request takes
			err := s.client.NotifySync(line)

			//Printing lines that encountered an error for the user to manually handle them
			if err != (constants.NotificationError{}) {
				errors = append(errors, err)
				_, _ = fmt.Fprintf(w, "%+v\n", err)
			}
		}(i, line)

		time.Sleep(time.Millisecond * time.Duration(interval))
	}

	// Wait for all the requests to notifications to finish
	wg.Wait()

	return errors
}

//ProcessMessagesChannel receives a slice of messages to process and the interval of time to pause between each one.
//For every message that fails to be processed, an error message is printed. Channels are used to be compatible with an async notifier library.
func (s *service) ProcessMessagesChannel(w io.Writer, lines []string, interval int64) []constants.NotificationError {

	ticker := time.NewTicker(time.Duration(interval) * time.Millisecond)
	quit := make(chan struct{})

	var errors []constants.NotificationError

Loop:
	for len(lines) > 0 {
		select {
		case <-ticker.C:
			//Getting first element from lines slice
			line := lines[0]
			//Popping the element
			lines = lines[1:]
			//Letting the user know what element is being processed currently
			_, _ = fmt.Fprintf(w, "Current line: '%s' - Waiting %d milliseconds\n", line, interval)

			//Messages are sent to be processed every n milliseconds, regardless of how much an HTTP request takes
			c := make(chan constants.NotificationError)
			go func() {
				s.client.NotifyChannel(line, c)
				err := <-c

				//Printing lines that encountered an error for the user to manually handle them
				if err != (constants.NotificationError{}) {
					errors = append(errors, err)
					_, _ = fmt.Fprintf(w, "%+v\n", err)
				}
			}()
		case <-quit:
			ticker.Stop()
			break Loop
		}
	}

	return errors
}
