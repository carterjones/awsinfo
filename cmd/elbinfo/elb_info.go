package main

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/aws/aws-sdk-go/service/elb"
)

type elbInfo struct {
	Name        string
	DNSName     string
	IPAddresses []string
}

func (i elbInfo) Matches(value string) bool {
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

func (i elbInfo) String() string {
	var msg string
	msg += fmt.Sprintf("Name:         %s\n", i.Name)
	msg += fmt.Sprintf("DNS Name:     %s\n", i.DNSName)
	msg += fmt.Sprintf("IP Addresses: %v", i.IPAddresses)
	return msg
}

type elbInfoSlice []elbInfo

func (e *elbInfoSlice) Import(o *elb.DescribeLoadBalancersOutput) {
	var r net.Resolver
	for _, lb := range o.LoadBalancerDescriptions {
		var dnsName, name string
		if lb.DNSName != nil {
			dnsName = *lb.DNSName
		}
		if lb.LoadBalancerName != nil {
			name = *lb.LoadBalancerName
		}
		addrs, err := r.LookupHost(context.Background(), dnsName)
		panicIfErr(err)

		*e = append(*e, elbInfo{
			DNSName:     dnsName,
			Name:        name,
			IPAddresses: addrs,
		})
	}
}
