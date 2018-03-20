package awsinfo

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
)

// InstanceInfo contains a minimal set of information about an EC2 instance.
type InstanceInfo struct {
	ImageID          string
	InstanceID       string
	InstanceType     string
	LaunchTime       time.Time
	PrivateIPAddress string
	PublicIPAddress  string
	Name             string
}

func (i InstanceInfo) String() string {
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

// Matches determines if a value can be found in the data for the EC2 instance.
func (i InstanceInfo) Matches(value string) bool {
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

// InstanceInfoSlice is a slice of InstanceInfo objects.
type InstanceInfoSlice []InstanceInfo

// Load gathers data from AWS about all the EC2 instances in the account.
func (info *InstanceInfoSlice) Load(sess *session.Session) error {
	// Create a new EC2 service handle.
	svc := ec2.New(sess)

	// Get information about all instances.
	v, err := svc.DescribeInstances(nil)
	if err != nil {
		return errors.Wrap(err, "could not get instance info")
	}

	for _, reservation := range v.Reservations {
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
			*info = append(*info, InstanceInfo{
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

	return nil
}
