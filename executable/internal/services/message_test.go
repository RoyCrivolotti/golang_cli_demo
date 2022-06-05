package services_test

import (
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"os"
	"refurbedchallenge/executable/internal/services"
	"refurbedchallenge/notifier/constants"
	mock_src "refurbedchallenge/notifier/src/mock"
	"testing"
)

func TestProcessMessages(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		description string
		input       struct {
			Interval   int64
			Lines      []string
			ErrorLines []string
		}
		want []string
	}{
		{
			"TestCase: Happy path",
			struct {
				Interval   int64
				Lines      []string
				ErrorLines []string
			}{
				1000,
				[]string{"this", "is", "", "", "testing"},
				[]string{},
			},
			[]string(nil),
		},
		{
			"TestCase: Errors due to, for instance, invalid URL for the HTTP POST request",
			struct {
				Interval   int64
				Lines      []string
				ErrorLines []string
			}{
				1000,
				[]string{"this", "is", "", "", "testing"},
				[]string{"this", "is", "", "", "testing"},
			},
			[]string([]string{"Line 'this', error: error\n", "Line 'is', error: error\n", "Line '', error: error\n", "Line '', error: error\n", "Line 'testing', error: error\n"}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			notifierClientMock := mock_src.NewMockINotifierClient(ctrl)
			service := services.NewMessageService(notifierClientMock)

			for i, line := range tt.input.Lines {
				channel := make(chan constants.NotificationError)

				go func(i int, line string) {
					if i >= 0 && i < len(tt.input.ErrorLines) && line == tt.input.ErrorLines[i] {
						fmt.Println(line)
						channel <- constants.NotificationError{Error: errors.New("error"), Message: line}
					} else {
						channel <- constants.NotificationError{}
					}
				}(i, line)

				notifierClientMock.EXPECT().Notify(line).Return(channel).Times(1)
			}

			errors := service.ProcessMessages(os.Stdout, tt.input.Lines, tt.input.Interval, false)

			assert.Equal(t, tt.want, errors)
		})
	}
}
