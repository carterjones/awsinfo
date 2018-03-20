package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/carterjones/awsinfo"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: elbinfo <search-value>")
		os.Exit(1)
	}

	searchValue := os.Args[1]

	// Tell the SDK to load defaults from your ~/.aws/config file.
	os.Setenv("AWS_SDK_LOAD_CONFIG", "true")

	// Create a new session.
	sess, err := session.NewSession()
	panicIfErr(err)

	var infos awsinfo.ELBInfoSlice
	infos.Load(sess)

	// Print the matches.
	numMatches := 0
	for _, lb := range infos {
		if lb.Matches(searchValue) {
			if numMatches > 0 {
				fmt.Println()
			}

			fmt.Println(lb)
			numMatches++
		}
	}
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
