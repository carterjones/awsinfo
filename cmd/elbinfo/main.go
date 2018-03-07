package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elb"
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

	// Create a new EC2 service handle.
	svc := elb.New(sess)

	// Get information about all instances.
	v, err := svc.DescribeLoadBalancers(nil)
	panicIfErr(err)

	var info elbInfoSlice
	info.Import(v)

	for _, lb := range info {
		if lb.Matches(searchValue) {
			fmt.Println(lb)
		}
	}
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
