# refurbed_challenge

This repository contains my challenge for `refurbed`.

### Ecosystem in which the application it was run and tested

- MacBook Pro 13-inch 2020, `macOS Monterrey Version 12.3.1`
- Go version [`go1.18.3.darwin-amd64.pkg`](https://go.dev/dl/)

### Explaining the use of a monorepo

While generally, in a production-ready scenario, one would upload code that is meant to be used as a library in a
different repository from the executable, allowing multiple projects to import it and make use of it, in this specific
case I chose to use a monorepo with multiple Go `modules` because they are part of the same challenge. This allows me to
have a single deliverable and makes it easier to test locally the integration of both projects locally while debugging.

### Assumptions and choices

- As mentioned above ([Explaining the use of a monorepo](#explaining-the-use-of-a-monorepo)), I chose to do the
  challenge within a single repository, using multiple Go `modules`. This also allows me to track my own changes and
  send a deliverable as a repository instead of a zip file.
- I assumed that the executable did not need to be interactive, at least not in this first version (
  see [Next steps](#next-steps))
- I assumed that the URL where the messages are being sent handles receiving multiple requests at once: if the service
  throws a throttling error, there is no retry logic on the notifier's side at this time

### How to test app locally

- First, create the test `txt` file that is meant to be streamed to the CLI, and save it in the following
  directory: `./executable/test_files/FILE_NAME.txt`
- From the executable module's root folder (that is: `./executable`):
    - Build the executable in the following directory: `./executable/cmd/cli` -> `go build -o bin/exe ./cmd/cli/main.go`
    - Run the executable with the desired
      parameters : `./bin/exe -i=1000 -url=http://url.com < test_files/FILE_NAME.txt`
    - You can ask what parameters the CLI accepts by running: `./bin/exe -help`

### How to run unit tests

- Using a terminal from the project's root folder (note: this is not a go module), run the following
  command: `find . -name go.mod -execdir go test ./... \;`
- Or, from within each module, respectively: `go test ./...`

### Next steps

- Instead of adding `replace` directives to work across multiple modules, use Go
  1.18's `workspaces`: https://go.dev/doc/tutorial/workspaces
- Instead of printing the messages that encountered an error to the user's screen, maintain an interactive state in the
  terminal and give the user options to retry the messages, and consider saving each message that failed in a log file (
  in a format that makes it easy to process it programmatically)
  to make it easier to collect the data after the execution
- Improve the styles of the messages being printed for them to be more readable
- As hinted at in the third point of the [assumptions](#assumptions) section, consider handling a throttling error
  within the notifier itself
