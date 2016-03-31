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

	"github.com/svenmueller/nube/Godeps/_workspace/src/github.com/aws/aws-sdk-go/aws"
	"github.com/svenmueller/nube/Godeps/_workspace/src/github.com/aws/aws-sdk-go/service/route53"
	"github.com/svenmueller/nube/Godeps/_workspace/src/github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/svenmueller/nube/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/svenmueller/nube/Godeps/_workspace/src/github.com/spf13/viper"
	"github.com/svenmueller/nube/common"
	"github.com/svenmueller/nube/util"
)

var servers_instance_destroyCmd = &cobra.Command{
	Use:   "destroy ID [ID|Name ...]",
	Short: "Destroy server instance by id or name",
	Run: func(cmd *cobra.Command, args []string) {
		err := serversInstanceDestroy(cmd, args)
		common.HandleError(err, cmd)
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("hosted-zone-id", cmd.Flags().Lookup("hosted-zone-id"))
	},
}

func init() {
	servers_instanceCmd.AddCommand(servers_instance_destroyCmd)
	servers_instance_destroyCmd.Flags().StringP("hosted-zone-id", "a", "", "Delete DNS resource record with same name in hosted zone.")
}

func serversInstanceDestroy(cmd *cobra.Command, args []string) error {

	if len(args) < 1 {
		return common.NewMissingArgumentsError(cmd)
	}

	rackspaceServiceClient, err := util.NewRackspaceService()

	if err != nil {
		return fmt.Errorf("Unable to establish connection: %v", err)
	}

	awsServiceClient := util.NewRoute53Service()

	var list []servers.Server
	listInitialized := false

	for _, idOrName := range args {
		if !listInitialized {
			list, err = util.ListAllServers(rackspaceServiceClient)

			if err != nil {
				return fmt.Errorf("Unable to build server list: %v", err)
			}

			listInitialized = true
		}

		var matchedServer *servers.Server
		for _, server := range list {
			if server.Name == idOrName || server.ID == idOrName {
				matchedServer = &server
				break
			}
		}

		if matchedServer == nil {
			return fmt.Errorf("Unable to find server %q", idOrName)
		}
		result := servers.Delete(rackspaceServiceClient, matchedServer.ID)

		if result.ErrResult.Result.Err != nil {
			return fmt.Errorf("Unable to delete server %q: %v", matchedServer.ID, result.ErrResult.Result.Err)
		}

		fmt.Printf("Destroyed server %q (ID: %q)\n", matchedServer.Name, matchedServer.ID)

		if viper.GetString("hosted-zone-id") != "" {
			params := &route53.ChangeResourceRecordSetsInput{
				HostedZoneId: aws.String(viper.GetString("hosted-zone-id")),
				ChangeBatch: &route53.ChangeBatch{
					Changes: []*route53.Change{
						{
							Action: aws.String("DELETE"),
							ResourceRecordSet: &route53.ResourceRecordSet{
								Name: aws.String(fmt.Sprintf("%s.", matchedServer.Name)),
								Type: aws.String("A"),
								TTL:  aws.Int64(3600),
								ResourceRecords: []*route53.ResourceRecord{
									{
										Value: aws.String(matchedServer.AccessIPv4),
									},
								},
							},
						},
					},
				},
			}

			resp, err := awsServiceClient.ChangeResourceRecordSets(params)

			if err != nil {
				return fmt.Errorf("Unable to delete resource record set %s: %v", fmt.Sprintf("%s.", matchedServer.Name), err.Error())
			}

			fmt.Printf("Deleted resource record set %q from hosted zone %q \n", fmt.Sprintf("%s.", matchedServer.Name), viper.GetString("hosted-zone-id"))
		}
	}

	return nil
}
