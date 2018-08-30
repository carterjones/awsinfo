package main

import (
	"fmt"
	"os"
	"reflect"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/carterjones/awsinfo/info/elb"
	"github.com/carterjones/awsinfo/info/instance"
	"github.com/carterjones/awsinfo/info/r53"
)

type matcher interface {
	Matches(string) bool
}

type loader interface {
	Load(*session.Session) error
}

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
		var matchers []matcher
		switch l.(type) {
		case *instance.InfoSlice:
			infos := l.(*instance.InfoSlice)
			for _, i := range *infos {
				matchers = append(matchers, i)
			}
			fmt.Println("Instances:")
		case *elb.InfoSlice:
			infos := l.(*elb.InfoSlice)
			for _, i := range *infos {
				matchers = append(matchers, i)
			}
			fmt.Println("ELBs:")
		case *r53.InfoSlice:
			infos := l.(*r53.InfoSlice)
			for _, i := range *infos {
				matchers = append(matchers, i)
			}
			fmt.Println("Route53 entries:")
		default:
			panic("invalid type detected: " + reflect.TypeOf(l).String())
		}

		numMatches := 0
		for _, m := range matchers {
			if m.Matches(searchValue) {
				if numMatches > 0 {
					fmt.Println()
				}

				fmt.Print(m)
				numMatches++
			}
		}
	}
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
