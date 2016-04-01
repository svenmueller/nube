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
	"io/ioutil"
	"strings"
	"sync"

	"github.com/svenmueller/nube/Godeps/_workspace/src/github.com/aws/aws-sdk-go/aws"
	"github.com/svenmueller/nube/Godeps/_workspace/src/github.com/aws/aws-sdk-go/service/route53"
	"github.com/svenmueller/nube/Godeps/_workspace/src/github.com/docker/docker/pkg/namesgenerator"
	"github.com/svenmueller/nube/Godeps/_workspace/src/github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/svenmueller/nube/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/svenmueller/nube/common"
	"github.com/svenmueller/nube/util"
)

var servers_instance_createCmd = &cobra.Command{
	Use:   "create NAME [NAME ...]",
	Short: "Create server instance",
	Run: func(cmd *cobra.Command, args []string) {
		err := serversInstanceCreate(cmd, args)
		common.HandleError(err, cmd)
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		Cfg.BindPFlag("hosted-zone-id", cmd.Flags().Lookup("hosted-zone-id"))
		Cfg.BindPFlag("domain", cmd.Flags().Lookup("domain"))
		Cfg.BindPFlag("add-random-name", cmd.Flags().Lookup("add-random-name"))
		Cfg.BindPFlag("add-region", cmd.Flags().Lookup("add-region"))
		Cfg.BindPFlag("user-data", cmd.Flags().Lookup("user-data"))
		Cfg.BindPFlag("user-data-file", cmd.Flags().Lookup("user-data-file"))
		Cfg.BindPFlag("flavor", cmd.Flags().Lookup("flavor"))
		Cfg.BindPFlag("image", cmd.Flags().Lookup("image"))
		Cfg.BindPFlag("wait-for-active", cmd.Flags().Lookup("wait-for-active"))
	},
}

func init() {
	servers_instanceCmd.AddCommand(servers_instance_createCmd)

	servers_instance_createCmd.Flags().StringP("domain", "d", "ct-app.com", "Domain name to append to the hostname (e.g. server01.example.com)")
	servers_instance_createCmd.Flags().BoolP("add-random-name", "n", true, "Append random name to server name (e.g. server01-adjective-surname)")
	servers_instance_createCmd.Flags().BoolP("add-region", "r", true, "Append region to hostname (e.g. server01.lon)")
	servers_instance_createCmd.Flags().StringP("user-data", "u", "", "User data for creating server")
	servers_instance_createCmd.Flags().StringP("user-data-file", "p", "", "A path to a file for user data")
	servers_instance_createCmd.Flags().StringP("flavor", "f", "1 GB Performance", "Flavor of server")
	servers_instance_createCmd.Flags().StringP("image", "i", "Ubuntu 14.04 LTS (Trusty Tahr) (PVHVM)", "Image name of server")
	servers_instance_createCmd.Flags().BoolP("wait-for-active", "w", true, "Don't return until the create has succeeded or failed")
	servers_instance_createCmd.Flags().StringP("hosted-zone-id", "z", "", "Add DNS resource record (type A, public IPv4 address) to hosted zone.")

}

func serversInstanceCreate(cmd *cobra.Command, args []string) error {

	if len(args) < 1 {
		return common.NewMissingArgumentsError(cmd)
	}

	rackspaceServiceClient, err := util.NewRackspaceService(Cfg.GetString("rackspace-username"), Cfg.GetString("rackspace-api-key"), Cfg.GetString("rackspace-region"))

	if err != nil {
		return fmt.Errorf("Unable to establish connection: %v", err)
	}

	awsServiceClient := util.NewRoute53Service(Cfg.GetString("aws-access-key-id"), Cfg.GetString("aws-secret-access-key"))

	userData := ""
	filename := Cfg.GetString("user-data-file")
	if filename != "" {
		userDataFile, err := ioutil.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("Error reading user data file %q: %v.", filename, err)
		}
		userData = string(userDataFile)
	} else {
		userData = Cfg.GetString("user-data")
	}

	var waitGroup sync.WaitGroup
	errs := make(chan error, len(args))

	for _, name := range args {
		// Add domain to end if available.
		if Cfg.GetBool("add-random-name") {
			randomServerName := strings.Replace(namesgenerator.GetRandomName(0), "_", "-", -1)
			name = fmt.Sprintf("%s-%s", name, randomServerName)
		}
		if Cfg.GetBool("add-region") {
			name = fmt.Sprintf("%s.%s", name, Cfg.GetString("rackspace-region"))
		}
		if Cfg.GetString("domain") != "" {
			name = fmt.Sprintf("%s.%s", name, Cfg.GetString("domain"))
		}

		name = strings.ToLower(name)

		configDrive := false
		if userData != "" {
			configDrive = true
		}

		opts := &servers.CreateOpts{
			Name:        name,
			ImageName:   Cfg.GetString("image"),
			FlavorName:  Cfg.GetString("flavor"),
			UserData:    []byte(userData),
			ConfigDrive: configDrive,
		}

		waitGroup.Add(1)

		go func() {
			defer waitGroup.Done()

			server, err := servers.Create(rackspaceServiceClient, opts).Extract()

			if err != nil {
				errs <- err
				return
			}

			fmt.Printf("Creating server %q\n", opts.Name)

			if Cfg.GetBool("wait-for-active") || Cfg.GetString("hosted-zone-id") != "" {
				fmt.Printf("Waiting for server %q\n", opts.Name)
				err = servers.WaitForStatus(rackspaceServiceClient, server.ID, "ACTIVE", 600)

				if err != nil {
					errs <- err
					return
				}

				server, err = servers.Get(rackspaceServiceClient, server.ID).Extract()

				if err != nil {
					errs <- err
					return
				}
			}

			fmt.Printf("Created server %q\n\n", server.Name)
			util.WriteOutput(server, Cfg.GetString("output"))

			if Cfg.GetString("hosted-zone-id") != "" {

				params := &route53.ChangeResourceRecordSetsInput{
					HostedZoneId: aws.String(Cfg.GetString("hosted-zone-id")),
					ChangeBatch: &route53.ChangeBatch{
						Changes: []*route53.Change{
							{
								Action: aws.String("CREATE"),
								ResourceRecordSet: &route53.ResourceRecordSet{
									Name: aws.String(fmt.Sprintf("%s.", server.Name)),
									Type: aws.String("A"),
									TTL:  aws.Int64(3600),
									ResourceRecords: []*route53.ResourceRecord{
										{
											Value: aws.String(server.AccessIPv4),
										},
									},
								},
							},
						},
					},
				}

				fmt.Printf("Creating DNS resource record set %q\n", server.Name)
				_, err := awsServiceClient.ChangeResourceRecordSets(params)

				if err != nil {
					errs <- err
					return
				}

				fmt.Printf("Created DNS resource record set %q (Type: %s, TTL: %d, Value: %s)\n", server.Name, "A", 3600, server.AccessIPv4)
			}
		}()
	}

	waitGroup.Wait()
	close(errs)

	for err := range errs {
		if err != nil {
			return err
		}
	}

	return nil
}
