package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

type instanceInfo struct {
	ImageID          string
	InstanceID       string
	InstanceType     string
	LaunchTime       time.Time
	PrivateIPAddress string
	PublicIPAddress  string
	Name             string
}

func (i instanceInfo) String() string {
	var msg string
	if i.Name != "" {
		msg += fmt.Sprintf("Name:        %s\n", i.Name)
	}
	if i.PublicIPAddress != "" {
		msg += fmt.Sprintf("Public IP:   %s\n", i.PublicIPAddress)
	}
	msg += fmt.Sprintf("Private IP:  %s\n", i.PrivateIPAddress)
	msg += fmt.Sprintf("ID:          %s\n", i.InstanceID)
	msg += fmt.Sprintf("AMI:         %s\n", i.ImageID)
	msg += fmt.Sprintf("Type:        %s\n", i.InstanceType)
	msg += fmt.Sprintf("Launch Time: %v\n", i.LaunchTime)
	return msg
}

func (i instanceInfo) Matches(value string) bool {
	if strings.Contains(i.Name, value) {
		return true
	}
	if strings.Contains(i.ImageID, value) {
		return true
	}
	if strings.Contains(i.InstanceID, value) {
		return true
	}
	if strings.Contains(i.InstanceType, value) {
		return true
	}
	if strings.Contains(i.LaunchTime.String(), value) {
		return true
	}
	if strings.Contains(i.PrivateIPAddress, value) {
		return true
	}
	if strings.Contains(i.PublicIPAddress, value) {
		return true
	}
	return false
}

func extractInstanceInfo(o *ec2.DescribeInstancesOutput) []instanceInfo {
	var info []instanceInfo
	for _, reservation := range o.Reservations {
		for _, instance := range reservation.Instances {
			var name string
			for _, tag := range instance.Tags {
				if *tag.Key == "Name" {
					name = *tag.Value
				}
			}
			var publicIp string
			if instance.PublicIpAddress != nil {
				publicIp = *instance.PublicIpAddress
			}
			info = append(info, instanceInfo{
				Name:             name,
				ImageID:          *instance.ImageId,
				InstanceID:       *instance.InstanceId,
				InstanceType:     *instance.InstanceType,
				LaunchTime:       *instance.LaunchTime,
				PrivateIPAddress: *instance.PrivateIpAddress,
				PublicIPAddress:  publicIp,
			})
		}
	}
	return info
}

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
	info := extractInstanceInfo(v)

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
