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
	"os"
	"strings"
	"sync"

	"github.com/svenmueller/nube/Godeps/_workspace/src/github.com/docker/docker/pkg/namesgenerator"
	"github.com/svenmueller/nube/Godeps/_workspace/src/github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/svenmueller/nube/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/svenmueller/nube/Godeps/_workspace/src/github.com/spf13/viper"
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

	viper.BindPFlag("domain", servers_instance_createCmd.Flags().Lookup("domain"))
	viper.BindPFlag("add-random-name", servers_instance_createCmd.Flags().Lookup("add-random-name"))
	viper.BindPFlag("add-region", servers_instance_createCmd.Flags().Lookup("add-region"))
	viper.BindPFlag("user-data", servers_instance_createCmd.Flags().Lookup("user-data"))
	viper.BindPFlag("user-data-file", servers_instance_createCmd.Flags().Lookup("user-data-file"))
	viper.BindPFlag("flavor", servers_instance_createCmd.Flags().Lookup("flavor"))
	viper.BindPFlag("image", servers_instance_createCmd.Flags().Lookup("image"))
	viper.BindPFlag("wait-for-active", servers_instance_createCmd.Flags().Lookup("wait-for-active"))

}

func serversInstanceCreate(cmd *cobra.Command, args []string) error {

	if len(args) < 1 {
		return common.NewMissingArgumentsError(cmd)
	}

	rackspaceServiceClient, err := util.NewRackspaceService()

	if err != nil {
		return fmt.Errorf("Unable to establish connection: %v", err)
	}

	userData := ""
	userDataPath := viper.GetString("user-data-file")
	if userDataPath != "" {
		file, errOpen := os.Open(userDataPath)
		if errOpen != nil {
			return fmt.Errorf("Error opening user data file %q: %v.", userDataPath, errOpen)
		}

		userDataFile, errRead := ioutil.ReadAll(file)
		if errRead != nil {
			return fmt.Errorf("Error reading user data file %q: %v.", userDataPath, errRead)
		}
		userData = string(userDataFile)
	} else {
		userData = viper.GetString("user-data")
	}

	var waitGroup sync.WaitGroup
	errs := make(chan error, len(args))

	for _, name := range args {
		// Add domain to end if available.
		if viper.GetBool("add-random-name") {
			randomServerName := strings.Replace(namesgenerator.GetRandomName(0), "_", "-", -1)
			name = fmt.Sprintf("%s-%s", name, randomServerName)
		}
		if viper.GetBool("add-region") {
			name = fmt.Sprintf("%s.%s", name, viper.GetString("rackspace-region"))
		}
		if viper.GetString("domain") != "" {
			name = fmt.Sprintf("%s.%s", name, viper.GetString("domain"))
		}

		name = strings.ToLower(name)

		opts := &servers.CreateOpts{
			Name:       name,
			ImageName:  viper.GetString("image"),
			FlavorName: viper.GetString("flavor"),
			UserData:   []byte(userData),
		}

		waitGroup.Add(1)

		go func() {
			defer waitGroup.Done()

			server, err := servers.Create(rackspaceServiceClient, opts).Extract()
			fmt.Printf("Creating server %q\n", opts.Name)

			if err != nil {
				errs <- err
				return
			}

			if viper.GetBool("wait-for-active") {
				fmt.Printf("Waiting for state %q of server %q\n", "ACTIVE", opts.Name)
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

			util.WriteOutput(server, viper.GetString("output"))
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
