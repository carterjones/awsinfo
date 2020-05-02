package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/carterjones/awsinfo/info/r53"
)

func main() {
	onlyIPInfo := flag.Bool("ips", false, "if true, only show IP address information")
	flag.Parse()
	tail := flag.Args()
	if len(tail) != 1 {
		fmt.Println("usage: r53info [-ips] <search-value>")
		os.Exit(1)
	}

	searchValue := tail[0]

	// Tell the SDK to load defaults from your ~/.aws/config file.
	os.Setenv("AWS_SDK_LOAD_CONFIG", "true")

	// Create a new session.
	sess, err := session.NewSession()
	panicIfErr(err)

	// Load the route53 info.
	var infos r53.InfoSlice
	err = infos.Load(sess)
	panicIfErr(err)

	// Print the matches.
	justPrintedSomething := false
	for _, v := range infos {
		if v.Matches(searchValue) {
			if justPrintedSomething {
				fmt.Println()
			}

			if *onlyIPInfo {
				msg := v.IPInfo()
				if msg != "" {
					fmt.Print(msg)
					justPrintedSomething = true
				} else {
					justPrintedSomething = false
				}
			} else {
				fmt.Print(v)
				justPrintedSomething = true
			}
		}
	}
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
