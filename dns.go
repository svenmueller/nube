package main

import (
	"fmt"
	"log"

	"github.com/svenmueller/nube/Godeps/_workspace/src/github.com/aws/aws-sdk-go/aws"
	"github.com/svenmueller/nube/Godeps/_workspace/src/github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/svenmueller/nube/Godeps/_workspace/src/github.com/aws/aws-sdk-go/aws/session"
	"github.com/svenmueller/nube/Godeps/_workspace/src/github.com/aws/aws-sdk-go/service/route53"
	"github.com/svenmueller/nube/Godeps/_workspace/src/github.com/codegangsta/cli"
)

var AWSAccessKey string
var AWSSecretKey string

var DNSCommand = cli.Command{
	Name:    "dns",
	Aliases: []string{"d"},
	Usage:   "DNS commands.",
	Flags: []cli.Flag{
		cli.StringFlag{Name: "aws-access-key", Value: "", Usage: "The AWS access key.", EnvVar: "AWS_ACCESS_KEY", Destination: &AWSAccessKey},
		cli.StringFlag{Name: "aws-secret-key", Value: "", Usage: "The AWS secret key.", EnvVar: "AWS_SECRET_KEY", Destination: &AWSSecretKey},
	},
	Subcommands: []cli.Command{
		{
			Name:    "zones",
			Aliases: []string{"z"},
			Usage:   "Hosted Zones commands.",
			Subcommands: []cli.Command{
				{
					Name:    "list",
					Aliases: []string{"l"},
					Usage:   "List all available servers.",
					Action:  dnsHostedZonesList,
				},
			},
		},
	},
}

func validatDNSRequiredArgs() {
	if AWSAccessKey == "" {
		log.Fatal("You must provide the AWS API Access Key via AWS_ACCESS_KEY environment variable or via CLI argument.")
	}

	if AWSSecretKey == "" {
		log.Fatal("You must provide the AWS secret key via AWS_SECRET_KEY environment variable or via CLI argument.")
	}
}

func dnsHostedZonesList(ctx *cli.Context) {

	validatDNSRequiredArgs()

	provider := credentials.StaticProvider{
		Value: credentials.Value{
			AccessKeyID:     AWSAccessKey,
			SecretAccessKey: AWSSecretKey,
			SessionToken:    "",
		},
	}

	serviceClient := route53.New(session.New(aws.NewConfig().WithCredentials(credentials.NewCredentials(&provider))))
	params := &route53.ListHostedZonesInput{}
	resp, err := serviceClient.ListHostedZones(params)

	if err != nil {
		log.Fatalf("Error: Unable to list hosted zones: %s.", err)
	}

	cliOut := NewCLIOutput()
	defer cliOut.Flush()
	cliOut.Header("Caller Reference", "Id", "Name", "ResourceRecordSetCount", "PrivateZone", "Comment")
	for _, hostedZone := range resp.HostedZones {

		comment := ""
		if hostedZone.Config.Comment != nil {
			comment = *hostedZone.Config.Comment
		}

		cliOut.Writeln("%s\t%s\t%s\t%d\t%t\t%s\n",
			*hostedZone.CallerReference, *hostedZone.Id, *hostedZone.Name, *hostedZone.ResourceRecordSetCount, *hostedZone.Config.PrivateZone, comment)
	}
}

func describe(i interface{}) {
	fmt.Printf("(%v, %T)\n", i, i)
}
