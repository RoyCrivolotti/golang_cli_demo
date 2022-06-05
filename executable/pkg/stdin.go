package pkg

import (
	"bufio"
	"os"
)

//ReadLinesFromStdin returns every line from stdin as a new entry of a string slice, and an error should one have occurred at any point
func ReadLinesFromStdin() ([]string, error) {
	var lines []string
	var err error = nil

	if stats, _err := os.Stdin.Stat(); _err != nil || stats.Size() <= 0 {
		err = _err
		return lines, err
	}

	reader := bufio.NewReader(os.Stdin)
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		currLine := scanner.Text()
		lines = append(lines, currLine)
	}

	return lines, err
}
