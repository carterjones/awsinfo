package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: instanceinfo <search-value>")
		os.Exit(1)
	}

	searchValue := os.Args[1]

	// Tell the SDK to load defaults from your ~/.aws/config file.
	os.Setenv("AWS_SDK_LOAD_CONFIG", "true")

	// Create a new session.
	sess, err := session.NewSession()
	panicIfErr(err)

	// Create a new EC2 service handle.
	svc := ec2.New(sess)

	// Get information about all instances.
	v, err := svc.DescribeInstances(nil)
	panicIfErr(err)

	// Extract the info we care about.
	var info instanceInfoSlice
	info.Import(v)

	// Find matches.
	var matches []instanceInfo
	for _, instance := range info {
		if instance.Matches(searchValue) {
			matches = append(matches, instance)
		}
	}

	// Only print the info we care about.
	for i, instance := range matches {
		if i > 0 {
			fmt.Println()
		}

		fmt.Print(instance)
	}
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
