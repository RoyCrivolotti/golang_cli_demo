package main

import (
	"clidemo/executable/internal/services"
	"clidemo/executable/pkg"
	"clidemo/notifier/src"
	"flag"
	"fmt"
	url2 "net/url"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	//Exiting gracefully in case of SIGINT
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Print("\r") //Using carriage return character to avoid the console printing '^C'
		fmt.Println("SIGINT captured: Exiting program...")
		os.Exit(1)
	}()

	fmt.Println("Executing program")

	//Getting the custom configuration from flags
	url := flag.String("url", "http://localhost:8080/notify", "Configure endpoint where events will be notified")
	debug := flag.Bool("debug", false, "More printing statements are used when in debug mode")
	interval := flag.Int64("i", 100, "Configure interval of time for messages in stdin to be processed (milliseconds)")

	//Printing arguments passed for debugging
	if *debug == true {
		for i, arg := range os.Args {
			fmt.Printf("Argument at %d: %s\n", i, arg)
		}
	}

	//Parsing flags
	flag.Parse()

	//Validating the URL value
	_, urlParsingError := url2.ParseRequestURI(*url)
	if *url == "" || urlParsingError != nil {
		fmt.Println("The URL flag has to be a valid, non-empty URL")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *debug == true {
		//Printing values being used after parsing
		fmt.Printf("URL is %s\n", *url)
		fmt.Printf("Interval is %d\n", *interval)
	}

	//Reading input stream from stdin
	lines, err := pkg.ReadLinesFromStdin()

	if len(lines) == 0 || err != nil {
		fmt.Println("Failed to read input; please input a valid stream of text lines")

		if *debug == true {
			fmt.Printf("Error: %s\n", err.Error())
		}

		os.Exit(1)
	}

	//Instantiate notifier service
	messageService := services.NewMessageService(src.NewNotifierClient(*url))
	messageService.ProcessMessagesChannel(os.Stdout, lines, *interval)
}
