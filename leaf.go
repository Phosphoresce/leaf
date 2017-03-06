package main

import (
	"flag"
	"fmt"
	"net"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

func main() {
	var extAddress string
	regionPtr := flag.String("r", "us-east-1", "AWS Region: default works for global.")
	domainPtr := flag.String("d", "", "DNS Domain name")
	flag.Parse()

	// Get external ip via first interface without a private prefix
	ifaces, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for i := range ifaces {
		if !strings.HasPrefix(ifaces[i].String(), "127") &&
			!strings.HasPrefix(ifaces[i].String(), "10") &&
			!strings.HasPrefix(ifaces[i].String(), "192") &&
			!strings.HasPrefix(ifaces[i].String(), "172") &&
			!strings.Contains(ifaces[i].String(), ":") {
			extAddress = strings.Split(ifaces[i].String(), "/")[0]
		}
	}
	if extAddress == "" {
		fmt.Println("No external IPs found.")
		return
	}

	// Route 53 client
	client := route53.New(session.New(&aws.Config{Region: aws.String(*regionPtr)}))

	// find the hosted zone by name
	params := &route53.ListHostedZonesByNameInput{
		DNSName:  aws.String(*domainPtr),
		MaxItems: aws.String("1"),
	}
	zone, err := client.ListHostedZonesByName(params)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// update records
	params2 := &route53.ChangeResourceRecordSetsInput{
		HostedZoneId: zone.HostedZones[0].Id,
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action: aws.String("UPSERT"),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: zone.DNSName,
						Type: aws.String("A"),
						TTL:  aws.Int64(600),
						ResourceRecords: []*route53.ResourceRecord{
							{
								Value: aws.String(extAddress),
							},
						},
					},
				},
			},
		},
	}
	resp, err := client.ChangeResourceRecordSets(params2)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(resp.ChangeInfo)
}
