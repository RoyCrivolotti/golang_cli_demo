package main

import (
	"flag"
	"fmt"
	url2 "net/url"
	"os"
	"os/signal"
	"refurbedexe/pkg"
	"syscall"
	"time"
)

//TODO: Build async notifier that allows the end client to handle errors and document assumption taken (ie. return all errors at the end for them to manually handle it?)
//TODO: Consume notifier and do integration test, then write unit tests..
//TODO: Abstract code to corresponding pkgs, use Docker to automatize getting the dependencies, building the executable and running the app, maybe add a makefile to automatize running every unit test

func main() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Print("\r") //Using carriage return character to avoid the console printing '^C'
		fmt.Println("SIGINT captured: Exiting program...")
		os.Exit(1)
	}()

	fmt.Println("Executing program")

	//Printing arguments passed for debugging
	for i, arg := range os.Args {
		fmt.Printf("Argument at %d: %s\n", i, arg)
	}

	//Getting the custom configuration from flags
	url := flag.String("url", "http://localhost:8080/notify", "Configure endpoint where events will be notified")
	interval := flag.Int64("i", 1000, "Configure interval of time for messages in stdin to be processed (milliseconds)")

	//Parsing flags
	flag.Parse()

	//Validating the URL value
	_, urlParsingError := url2.ParseRequestURI(*url)
	if *url == "" || urlParsingError != nil {
		fmt.Println("The URL flag has to be a valid, non-empty URL")
		flag.PrintDefaults()
		os.Exit(1)
	}

	fmt.Printf("URL is %s\n", *url)
	fmt.Printf("Interval is %d\n", *interval)

	//Reading input stream from stdin
	lines, _ := pkg.ReadLinesFromStdin() //TODO handle err and document the assumption taken
	if len(lines) == 0 {
		fmt.Println("Failed to read input; please input a valid stream of text lines")
		os.Exit(1)
	}

	ticker := time.NewTicker(time.Duration(*interval) * time.Millisecond)
	quit := make(chan struct{})

	for {
		select {
		case <-ticker.C:
			//Getting first element from lines slice
			line := lines[0]
			//Popping the element
			lines = lines[1:]
			//Printing the popped element
			fmt.Printf("Current line: %s\n", line)
			fmt.Printf("Waiting %d milliseconds\n", *interval)
		case <-quit:
			ticker.Stop()
			return
		}
	}
}
