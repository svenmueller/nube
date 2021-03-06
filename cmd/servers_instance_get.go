// Copyright © 2016 Sven Mueller <skanky@gmx.de>
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

	"github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/spf13/cobra"
	"github.com/svenmueller/nube/common"
	"github.com/svenmueller/nube/util"
)

var servers_instance_getCmd = &cobra.Command{
	Use:   "get ID [ID|Name ...]",
	Short: "Get server instance",
	Run: func(cmd *cobra.Command, args []string) {
		err := serversInstanceGet(cmd, args)
		common.HandleError(err, cmd)
	},
}

func init() {
	servers_instanceCmd.AddCommand(servers_instance_getCmd)
}

func serversInstanceGet(cmd *cobra.Command, args []string) error {

	if len(args) < 1 {
		return common.NewMissingArgumentsError(cmd)
	}

	rackspaceServiceClient, err := util.NewRackspaceService(Cfg.GetString("rackspace-username"), Cfg.GetString("rackspace-api-key"), Cfg.GetString("rackspace-region"))

	if err != nil {
		return fmt.Errorf("Unable to establish connection: %v", err)
	}

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

		util.WriteOutput(matchedServer, Cfg.GetString("output"))
	}

	return nil
}
