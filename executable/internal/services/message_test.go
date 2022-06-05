package services_test

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"os"
	"refurbedchallenge/executable/internal/services"
	"refurbedchallenge/executable/pkg"
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
		want string
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
			"",
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
			"Line '', error: error\nLine '', error: error\nLine '', error: error\nLine 'testing', error: error\nLine 'testing', error: error\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			notifierClientMock := mock_src.NewMockINotifierClient(ctrl)
			service := services.NewMessageService(notifierClientMock)

			for i, line := range tt.input.Lines {
				channel := make(chan constants.NotificationError)

				go func() {
					if i >= 0 && i < len(tt.input.ErrorLines) && line == tt.input.ErrorLines[i] {
						channel <- constants.NotificationError{Error: errors.New("error"), Message: line}
					} else {
						channel <- constants.NotificationError{}
					}
				}()

				notifierClientMock.EXPECT().Notify(line).Return(channel).Times(1)
			}

			//This function created for testing captures everything written to stdout
			startCapturingStdout := pkg.Capture()

			service.ProcessMessages(os.Stdout, tt.input.Lines, tt.input.Interval, false)

			//Match the output captured with the expected result
			output, err := startCapturingStdout()
			assert.Equal(t, nil, err)
			assert.Equal(t, tt.want, output)
		})
	}
}
