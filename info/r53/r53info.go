package r53

import (
	"fmt"
	"net"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/pkg/errors"
)

// Info contains minimal data about a Route53 entry.
type Info struct {
	Zone        string
	Name        string
	Values      []string
	AliasTarget string
}

// IPInfo returns only information related to IP addresses.
func (i Info) IPInfo() string {
	var msg string
	msg += fmt.Sprintf("Name:         %s\n", i.Name)
	valuesHaveIPs := false
	for _, v := range i.Values {
		if net.ParseIP(v) != nil {
			valuesHaveIPs = true
			break
		}
	}
	if valuesHaveIPs {
		msg += fmt.Sprintf("Value:        %s\n", strings.Join(i.Values, ", "))
		return msg
	}

	return ""
}

func (i Info) String() string {
	var msg string
	msg += fmt.Sprintf("Name:         %s\n", i.Name)
	msg += fmt.Sprintf("Zone:         %s\n", i.Zone)
	if i.AliasTarget != "" {
		msg += fmt.Sprintf("Alias Target: %s", i.AliasTarget)
	}
	msg += fmt.Sprintf("Value:        %s\n", strings.Join(i.Values, ", "))
	return msg
}

// Matches determines if a value can be found in the data for the Route53 entry.
func (i Info) Matches(value string) bool {
	if strings.Contains(i.Name, value) {
		return true
	}
	if strings.Contains(i.Zone, value) {
		return true
	}
	if strings.Contains(i.AliasTarget, value) {
		return true
	}
	for _, v := range i.Values {
		if strings.Contains(v, value) {
			return true
		}
	}
	return false
}

// InfoSlice is a slice of Info objects.
type InfoSlice []Info

// Load gathers data from AWS about all the Route53 entries in the account.
func (Infos *InfoSlice) Load(sess *session.Session) error {
	// Create a new route53 service handle.
	svc := route53.New(sess)

	// Get all the hosted zones.
	var zones []*route53.HostedZone
	handleZones := func(out *route53.ListHostedZonesOutput, ok bool) bool {
		zones = append(zones, out.HostedZones...)
		return *out.IsTruncated
	}
	err := svc.ListHostedZonesPages(nil, handleZones)
	if err != nil {
		return errors.Wrap(err, "failed to load hosted zones")
	}

	// Get information about all route53 entries.
	for _, zone := range zones {
		in := &route53.ListResourceRecordSetsInput{
			HostedZoneId: zone.Id,
		}
		handleRecords := func(out *route53.ListResourceRecordSetsOutput, ok bool) bool {
			for _, rs := range out.ResourceRecordSets {
				zoneID := strings.Replace(*zone.Id, "/hostedzone/", "", -1)
				info := Info{
					Zone: *zone.Name + " (" + zoneID + ")",
					Name: *rs.Name,
				}

				target := rs.AliasTarget
				if target != nil {
					info.AliasTarget = *target.DNSName
				}

				for _, v := range rs.ResourceRecords {
					info.Values = append(info.Values, *v.Value)
				}

				*Infos = append(*Infos, info)
			}
			return *out.IsTruncated
		}

		err = svc.ListResourceRecordSetsPages(in, handleRecords)
		if err != nil {
			return errors.Wrap(err, "failed to load resource records")
		}
	}

	return nil
}
