package services_test

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"os"
	"refurbedchallenge/executable/internal/services"
	"refurbedchallenge/notifier/constants"
	mock_src "refurbedchallenge/notifier/src/mock"
	"testing"
)

func TestProcessMessagesWaitGroup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		description string
		input       struct {
			Interval   int64
			Lines      []string
			ErrorLines []string
		}
		want []constants.NotificationError
	}{
		{
			"TestCase: Happy path",
			struct {
				Interval   int64
				Lines      []string
				ErrorLines []string
			}{
				100,
				[]string{"this", "is", "", "", "testing"},
				[]string{},
			},
			[]constants.NotificationError(nil),
		},
		{
			"TestCase: Errors due to, for instance, invalid URL for the HTTP POST request",
			struct {
				Interval   int64
				Lines      []string
				ErrorLines []string
			}{
				100,
				[]string{"this", "is", "", "", "testing"},
				[]string{"this", "is", "", "", "testing"},
			},
			[]constants.NotificationError{{Error: errors.New("error"), Message: "this"}, {Error: errors.New("error"), Message: "is"}, {Error: errors.New("error"), Message: ""}, {Error: errors.New("error"), Message: ""}, {Error: errors.New("error"), Message: "testing"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			notifierClientMock := mock_src.NewMockINotifierClient(ctrl)
			service := services.NewMessageService(notifierClientMock)

			for i, line := range tt.input.Lines {
				var err constants.NotificationError

				if i >= 0 && i < len(tt.input.ErrorLines) && line == tt.input.ErrorLines[i] {
					err = constants.NotificationError{Error: errors.New("error"), Message: line}
				} else {
					err = constants.NotificationError{}
				}

				notifierClientMock.EXPECT().NotifySync(line).Return(err).Times(1)
			}

			errs := service.ProcessMessagesWaitGroup(os.Stdout, tt.input.Lines, tt.input.Interval)

			assert.Equal(t, tt.want, errs)
		})
	}
}
