package pkg

import (
	"bufio"
	"os"
	"strings"
)

//ReadLinesFromStdin returns every line from stdin as a new entry of a string slice, and an error should one have occurred at any point
func ReadLinesFromStdin() ([]string, error) {
	r := bufio.NewReader(os.Stdin)

	var bytes []byte
	var lines []string
	var err error = nil

	for {
		if stats, _err := os.Stdin.Stat(); _err != nil {
			err = _err
			break
		} else if stats.Size() <= 0 {
			break
		}

		line, isPrefix, _err := r.ReadLine()
		if _err != nil {
			err = _err
			break
		}

		bytes = append(bytes, line...)
		if !isPrefix {
			str := strings.TrimSpace(string(bytes))
			if len(str) > 0 {
				lines = append(lines, str)
				bytes = []byte{}
			}
		}
	}

	if len(bytes) > 0 {
		lines = append(lines, string(bytes))
	}

	return lines, err
}
