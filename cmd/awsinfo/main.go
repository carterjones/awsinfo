package main

import (
	"flag"
	"fmt"
	"os"
	"reflect"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/carterjones/awsinfo/info/elb"
	"github.com/carterjones/awsinfo/info/instance"
	"github.com/carterjones/awsinfo/info/r53"
)

type matchIPInfoer interface {
	Matches(string) bool
	IPInfo() string
}

type loader interface {
	Load(*session.Session) error
}

func main() {
	onlyIPInfo := flag.Bool("ips", false, "if true, only show IP address information")
	flag.Parse()
	tail := flag.Args()

	if len(tail) != 1 {
		fmt.Println("usage: awsinfo [-ips] <search-value>")
		os.Exit(1)
	}

	searchValue := tail[0]

	// Tell the SDK to load defaults from your ~/.aws/config file.
	os.Setenv("AWS_SDK_LOAD_CONFIG", "true")

	// Create a new session.
	sess, err := session.NewSession()
	panicIfErr(err)

	var loaders = []loader{
		&instance.InfoSlice{},
		&elb.InfoSlice{},
		&r53.InfoSlice{},
	}

	for _, l := range loaders {
		err = l.Load(sess)
		panicIfErr(err)
	}

	// Print the matches.
	for i, l := range loaders {
		if i > 0 {
			fmt.Println()
		}
		var matchIPInfoers []matchIPInfoer
		switch l.(type) {
		case *instance.InfoSlice:
			infos := l.(*instance.InfoSlice)
			if len(*infos) == 0 {
				continue
			}
			for _, i := range *infos {
				matchIPInfoers = append(matchIPInfoers, i)
			}
			fmt.Println("Instances:")
		case *elb.InfoSlice:
			infos := l.(*elb.InfoSlice)
			if len(*infos) == 0 {
				continue
			}
			for _, i := range *infos {
				matchIPInfoers = append(matchIPInfoers, i)
			}
			fmt.Println("ELBs:")
		case *r53.InfoSlice:
			infos := l.(*r53.InfoSlice)
			if len(*infos) == 0 {
				continue
			}
			for _, i := range *infos {
				matchIPInfoers = append(matchIPInfoers, i)
			}
			fmt.Println("Route53 entries:")
		default:
			panic("invalid type detected: " + reflect.TypeOf(l).String())
		}

		justPrintedSomething := false
		for _, v := range matchIPInfoers {
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
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
