// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/spf13/cobra"
	"github.com/svenmueller/nube/common"
	"github.com/svenmueller/nube/util"
)

var dns_zones_listCmd = &cobra.Command{
	Use:   "list",
	Short: "List hosted zones",
	Run: func(cmd *cobra.Command, args []string) {
		err := dnsHostedZonesList(cmd, args)
		common.HandleError(err, cmd)
	},
}

func init() {
	dns_zonesCmd.AddCommand(dns_zones_listCmd)
}

func dnsHostedZonesList(cmd *cobra.Command, args []string) error {

	serviceClient := util.NewRoute53Service()
	params := &route53.ListHostedZonesInput{}
	resp, err := serviceClient.ListHostedZones(params)

	if err != nil {
		return fmt.Errorf("Unable to list hosted zones: %v", err)
	}

	cliOut := util.NewCLIOutput()
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

	return nil
}
