package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/carterjones/awsinfo/info/instance"
)

func usage() {
	fmt.Println("usage: [-ips] instanceinfo <search-value>")
	os.Exit(1)
}

func main() {
	ipsPtr := flag.Bool("ips", false, "only print information about IPs")
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		usage()
	}

	searchValue := args[0]

	// Tell the SDK to load defaults from your ~/.aws/config file.
	os.Setenv("AWS_SDK_LOAD_CONFIG", "true")

	// Create a new session.
	sess, err := session.NewSession()
	panicIfErr(err)

	var infos instance.InfoSlice
	err = infos.Load(sess)
	panicIfErr(err)

	// Only print the info we care about.
	numMatches := 0
	for _, v := range infos {
		if v.Matches(searchValue) {
			if numMatches > 0 {
				fmt.Println()
			}

			if *ipsPtr {
				fmt.Print(v.IPs())
			} else {
				fmt.Print(v)
			}
			numMatches++
		}
	}
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
