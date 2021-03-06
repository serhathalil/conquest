package main

import (
	"flag"
	"fmt"
	"os"
	"time"
	
	"github.com/brsyuksel/conquest/conquest"
)

var (
	users                       uint64
	timeout, configfile, output string
	sequential, verbose         bool
)

func init() {
	flag.Uint64Var(&users, "u", 10, "concurrent users.")
	flag.StringVar(&timeout, "t", "30s",
		"duration for performing transactions. Use s, m, h modifiers")
	flag.StringVar(&output, "o", "", "output file for summary")
	flag.StringVar(&configfile, "c", "conquest.js", "conquest js file path")
	flag.BoolVar(&sequential, "s", false, "do transactions in sequential mode")
	flag.BoolVar(&verbose, "v", false, "print failed requests")
}

func main() {
	var err error
	flag.Parse()

	fmt.Println("conquest", "v" + conquest.VERSION, "\n")

	if _, err := os.Stat(configfile); os.IsNotExist(err) {
		fmt.Println(configfile, "file not found")
		os.Exit(1)
	}

	conq, err := conquest.RunScript(configfile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "u":
			conq.TotalUsers = users
		case "t":
			var duration time.Duration
			duration, err = time.ParseDuration(timeout)
			if err == nil {
				conq.Duration = duration
			}
		case "s":
			conq.Sequential = sequential
		}
	})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("performing transactions...\n")

	var fo *os.File
	if output == "" {
		fo = os.Stdout
	} else {
		fo, err = os.Create(output)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	reporter := conquest.NewReporter(fo, verbose)

	err = conquest.Perform(conq, reporter)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	<-reporter.C.Done
}
