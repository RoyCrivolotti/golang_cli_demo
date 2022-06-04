package main

//Library
//Create HTTP notifications client -> Configurable URL to which to propagate notifications
//Client gets message and sends them in the body of the POST request to the configured URL in a NON BLOCKER manner
//Should the Client get an error when sending the message, the caller should be able to handle the error

//Executable
//Consume library above
//Send new messages every interval (configurable)
//The exe reads standard input and every new line is a message to send to the library for propagation
//Consider shutdown by SIGINT
func main() {

}
