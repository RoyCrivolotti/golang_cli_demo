package services

import (
	"fmt"
	"refurbedchallenge/notifier/constants"
	"refurbedchallenge/notifier/src"
	"time"
)

type IMessageService interface {
	ProcessMessages(lines []string, interval int64)
}

type service struct {
	client src.INotifierClient
}

func NewMessageService(url string) IMessageService {
	return &service{
		client: src.NewNotifierClient(url),
	}
}

//ProcessMessages receives a slice of messages to process and the interval of time to pause between each one.
//For every message that fails to be processed, an error message is printed
func (s *service) ProcessMessages(lines []string, interval int64) {
	ticker := time.NewTicker(time.Duration(interval) * time.Millisecond)
	quit := make(chan struct{})

	var channels []chan constants.NotificationError

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
			//Printing the popped element
			fmt.Printf("Current line: %s - Waiting %d milliseconds\n", line, interval)

			c := s.client.Notify(line)
			channels = append(channels, c)
		case <-quit:
			ticker.Stop()
			break Loop
		}
	}

	//Printing lines that encountered an error for the user to manually handle them
	for _, c := range channels {
		err := <-c
		fmt.Printf("Error: Line '%s', error: %v\n", err.Message, err.Error)
	}
}
