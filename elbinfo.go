package awsinfo

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/pkg/errors"
)

// ELBInfo contains a minimal amount of information about a classic ELB.
type ELBInfo struct {
	Name        string
	DNSName     string
	IPAddresses []string
}

// Matches determines if a value can be found in the data for the ELB.
func (i ELBInfo) Matches(value string) bool {
	if strings.Contains(i.Name, value) {
		return true
	}
	if strings.Contains(i.DNSName, value) {
		return true
	}
	for _, ip := range i.IPAddresses {
		if strings.Contains(ip, value) {
			return true
		}
	}
	return false
}

func (i ELBInfo) String() string {
	var msg string
	msg += fmt.Sprintf("Name:         %s\n", i.Name)
	msg += fmt.Sprintf("DNS Name:     %s\n", i.DNSName)
	msg += fmt.Sprintf("IP Addresses: %v\n", i.IPAddresses)
	return msg
}

// ELBInfoSlice is a slice of ELBInfo objects.
type ELBInfoSlice []ELBInfo

// Load gathers data from AWS about all the classic ELBs in the account.
func (e *ELBInfoSlice) Load(sess *session.Session) error {
	// Create a new EC2 service handle.
	svc := elb.New(sess)

	// Get information about all instances.
	v, err := svc.DescribeLoadBalancers(nil)
	if err != nil {
		return errors.Wrap(err, "failed to get load balancer info")
	}

	var r net.Resolver
	for _, lb := range v.LoadBalancerDescriptions {
		var dnsName, name string
		if lb.DNSName != nil {
			dnsName = *lb.DNSName
		}
		if lb.LoadBalancerName != nil {
			name = *lb.LoadBalancerName
		}
		addrs, err := r.LookupHost(context.Background(), dnsName)
		if err != nil {
			return errors.Wrapf(err, "failed to look up host: %v", dnsName)
		}

		*e = append(*e, ELBInfo{
			DNSName:     dnsName,
			Name:        name,
			IPAddresses: addrs,
		})
	}

	return nil
}
