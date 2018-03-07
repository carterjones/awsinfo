package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/service/ec2"
)

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

type instanceInfoSlice []instanceInfo

func (info *instanceInfoSlice) Import(o *ec2.DescribeInstancesOutput) {
	for _, reservation := range o.Reservations {
		for _, instance := range reservation.Instances {
			var name, publicIP, privateIP, instanceID, imageID, instanceType string
			var launchTime time.Time
			for _, tag := range instance.Tags {
				if *tag.Key == "Name" {
					name = *tag.Value
				}
			}
			if instance.PublicIpAddress != nil {
				publicIP = *instance.PublicIpAddress
			}
			if instance.PrivateIpAddress != nil {
				privateIP = *instance.PrivateIpAddress
			}
			if instance.InstanceId != nil {
				instanceID = *instance.InstanceId
			}
			if instance.ImageId != nil {
				imageID = *instance.ImageId
			}
			if instance.InstanceType != nil {
				instanceType = *instance.InstanceType
			}
			if instance.LaunchTime != nil {
				launchTime = *instance.LaunchTime
			}
			*info = append(*info, instanceInfo{
				Name:             name,
				PublicIPAddress:  publicIP,
				PrivateIPAddress: privateIP,
				InstanceID:       instanceID,
				ImageID:          imageID,
				InstanceType:     instanceType,
				LaunchTime:       launchTime,
			})
		}
	}
}
