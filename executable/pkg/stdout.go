package pkg

import (
	"io"
	"os"
	"strings"
)

//Capture replaces os.Stdout with a writer that buffers any data written to os.Stdout.
//Call the returned function to clean up and get the data as a string.
//Written mostly to be used to test the executable's services
func Capture() func() (string, error) {
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	done := make(chan error, 1)

	save := os.Stdout
	os.Stdout = w

	var buf strings.Builder

	go func() {
		_, err := io.Copy(&buf, r)
		_ = r.Close()
		done <- err
	}()

	return func() (string, error) {
		os.Stdout = save
		_ = w.Close()
		err := <-done
		return buf.String(), err
	}
}
