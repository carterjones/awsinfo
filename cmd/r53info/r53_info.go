package main

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/pkg/errors"
)

type r53info struct {
	Zone        string
	Name        string
	Values      []string
	AliasTarget string
}

func (i r53info) String() string {
	var msg string
	msg += fmt.Sprintf("Name:         %s\n", i.Name)
	msg += fmt.Sprintf("Zone:         %s\n", i.Zone)
	msg += fmt.Sprintf("Value:        %s\n", strings.Join(i.Values, ", "))
	if i.AliasTarget != "" {
		msg += fmt.Sprintf("Alias Target: %s", i.AliasTarget)
	}
	return msg
}

func (i r53info) Matches(value string) bool {
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

type r53infoSlice []r53info

func (r53infos *r53infoSlice) Load(sess *session.Session) error {
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
				info := r53info{
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

				*r53infos = append(*r53infos, info)
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
