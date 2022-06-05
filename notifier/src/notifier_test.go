package src_test

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"refurbedchallenge/notifier/constants"
	"refurbedchallenge/notifier/src"
	"testing"
)

func TestNotify(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	//Using httpmock to mock outbound HTTP requests
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	tests := []struct {
		description string
		input       struct {
			url      string
			method   string
			message  string
			response struct {
				status int
				body   string
			}
		}
		want constants.NotificationError
	}{
		{
			"TestCase: Happy path",
			struct {
				url      string
				method   string
				message  string
				response struct {
					status int
					body   string
				}
			}{
				"dummy.com",
				"POST",
				"this is a message",
				struct {
					status int
					body   string
				}{
					200,
					`{}`,
				},
			},
			constants.NotificationError{},
		},
		{
			"TestCase: HTTP POST request returns an error",
			struct {
				url      string
				method   string
				message  string
				response struct {
					status int
					body   string
				}
			}{
				"dummyerror.com",
				"POST",
				"another message",
				struct {
					status int
					body   string
				}{
					408,
					`{"error":"this is an error message"}`,
				},
			},
			constants.NotificationError{Error: errors.New(`{"error":"this is an error message"}`), Message: "another message"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			client := src.NewNotifierClient(tt.input.url)

			httpmock.RegisterResponder(
				tt.input.method,
				tt.input.url,
				httpmock.NewStringResponder(tt.input.response.status, tt.input.response.body),
			)

			c := client.Notify(tt.input.message)
			err := <-c

			assert.Equal(t, tt.want, err)
		})
	}
}
