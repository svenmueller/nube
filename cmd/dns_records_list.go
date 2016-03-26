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

var dns_records_listCmd = &cobra.Command{
	Use:   "list NAME",
	Short: "List resource record sets",
	Run: func(cmd *cobra.Command, args []string) {
		err := dnsResourceRecordSetsList(cmd, args)
		common.HandleError(err, cmd)
	},
}

func init() {
	dns_recordsCmd.AddCommand(dns_records_listCmd)
}

func dnsResourceRecordSetsList(cmd *cobra.Command, args []string) error {

	if len(args) < 1 {
		return common.NewMissingArgumentsError(cmd)
	}

	hostedZoneId := args[0]

	serviceClient := util.NewRoute53Service()

	params := &route53.ListResourceRecordSetsInput{
		HostedZoneId: &hostedZoneId,
	}

	resp, err := serviceClient.ListResourceRecordSets(params)

	if err != nil {
		return fmt.Errorf("Unable to list resource record sets for hosted zone with ID %q: %v", hostedZoneId, err)
	}

	cliOut := util.NewCLIOutput()
	defer cliOut.Flush()
	cliOut.Header("Name", "Type", "TTL", "Resource Records")
	for _, resourceRecordSet := range resp.ResourceRecordSets {
		cliOut.Writeln("%s\t%s\t%d\t%v\n",
			*resourceRecordSet.Name, *resourceRecordSet.Type, *resourceRecordSet.TTL, *resourceRecordSet.ResourceRecords[0].Value)
	}

	return nil
}
