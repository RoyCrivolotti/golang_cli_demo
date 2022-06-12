# golang_cli_demo

This repository contains a personal demo.

## Ecosystem in which the application was run and tested

- MacBook Pro 13-inch 2020, `macOS Monterrey Version 12.3.1`
- Go version [`go1.18.3.darwin-amd64.pkg`](https://go.dev/dl/)

## Contents

- Create a library that implements an HTTP notification client. The client can be configured with a URL to which the notifications are sent through `HTTP POST` requests, with the message content in the request body. This operation should be non-blocking for the caller.
The program should implement graceful shutdown on SIGINT.
- Create a small command-line program that reads stdin and sends each line of the input every n milliseconds (configurable) using the afformentioned library.

**Example usage**

`> ./bin/exe -i=1000 -url=http://localhost:8080/notify < test_files/test.txt`

## Assumptions and choices

- While generally, in a production-ready scenario, one would upload code that is meant to be used as a library to a different repository, allowing multiple projects to import it and make use of it, in this specific case I chose to use a monorepo with multiple Go `modules`. This allows me to have a single deliverable and makes it easier to test the integration of both projects locally while debugging.
- I assumed that the executable did not need to be interactive, at least not in a first version (see [Next steps](#next-steps)).
- I assumed that the URL where the messages are being sent handles receiving multiple requests at once: if the service throws a throttling error, there is no retry logic on the notifier's side at this time.
- I implemented two versions of the notification/messaging service:
    - One using `channels`: The **benefit** of this implementation is that the library itself, which makes the `POST` requests, is async. As for the **negative** aspects, though:
        - They are less readable and intuitive than `WaitGroups`
        - They require both parts to handle the async implementation, the caller and the callee, which means that both have to handle `channels`
        - The difference in efficiency is extremely marginal, unless working with very costly algorithms such as complex graph traversals
        - They make testing incredibly more complicated
    - Another one using `WaitGroups`: This implementation is more readable and easier to understand, the **negative** aspect of it being that the library itself, which makes the `POST` requests, is not in itself `async` in this version; this means that the consumer has to ensure that each call is asynchronous 100% on their side.

## How to test app locally

From the project's root folder: `bash setup.sh`, and then run the following command from the project's root folder:

`> ./executable/bin/exe -i=1000 -url=http://url.com < ./executable/test_files/test.txt`

If you want to run the program with your own test file/input (after downloading the dependencies, that is):

- First, create the test `txt` file that is meant to be streamed to the CLI
- From the executable module's root folder (that is: `./executable`):
    - Build the executable in the following directory: `./cmd/cli` -> `go build -o bin/exe ./cmd/cli/main.go`
    - Run the executable with the desired parameters : `./bin/exe -i=1000 -url=http://url.com < path_to_test_file.txt`
    - You can ask what parameters the CLI accepts by running: `./bin/exe -help`

## How to run unit tests

From the project's root folder: `bash test.sh`, which does the following:

- Using a terminal from the project's root folder (note: this is not a go module), run the following
  command: `find . -name go.mod -execdir go test ./... \;`
- Or, from within each module, respectively: `go test ./...`

## Next steps

- Instead of adding `replace` directives to work across multiple modules, use Go
  1.18's `workspaces`: https://go.dev/doc/tutorial/workspaces
- Instead of printing the messages that encountered an error to the user's screen, maintain an interactive state in the
  terminal and give the user options to retry the messages, and consider saving each message that failed in a log file (
  in a format that makes it easy to process it programmatically)
  to make it easier to collect the data after the execution
- Improve the styles of the messages being printed for them to be more readable
- As hinted at in the third point of the [assumptions](#assumptions) section, consider handling a throttling error
  within the notifier itself
- Create a spy for the `notifier`, which would finally allow me to test the `channels` implementation of the message
  processing (`ProcessMessagesChannel` method).

