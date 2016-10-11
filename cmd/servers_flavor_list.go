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

	"github.com/rackspace/gophercloud/openstack/compute/v2/flavors"
	"github.com/rackspace/gophercloud/pagination"
	"github.com/spf13/cobra"
	"github.com/svenmueller/nube/common"
	"github.com/svenmueller/nube/util"
)

var servers_flavor_listCmd = &cobra.Command{
	Use:   "list",
	Short: "List server flavors",
	Run: func(cmd *cobra.Command, args []string) {
		err := serversFlavorList(cmd, args)
		common.HandleError(err, cmd)
	},
}

func init() {
	servers_flavorCmd.AddCommand(servers_flavor_listCmd)
}

func serversFlavorList(cmd *cobra.Command, args []string) error {

	// We have the option of filtering the server list. If we want the full
	// collection, leave it as an empty struct
	opts := flavors.ListOpts{}

	rackspaceServiceClient, err := util.NewRackspaceService(Cfg.GetString("rackspace-username"), Cfg.GetString("rackspace-api-key"), Cfg.GetString("rackspace-region"))

	if err != nil {
		return fmt.Errorf("Unable to establish connection: %v", err)
	}

	// Retrieve a pager (i.e. a paginated collection)
	pager := flavors.ListDetail(rackspaceServiceClient, opts)

	// Define an anonymous function to be executed on each page's iteration
	err = pager.EachPage(func(page pagination.Page) (bool, error) {
		flavorList, ExtractFlavorError := flavors.ExtractFlavors(page)

		if ExtractFlavorError != nil {
			return false, ExtractFlavorError
		}

		cliOut := util.NewCLIOutput()
		defer cliOut.Flush()
		cliOut.Header("ID", "Name", "VCPUs", "RAM", "Disk", "Swap", "RxTxFactor")
		for _, flavor := range flavorList {
			cliOut.Writeln("%s\t%s\t%d\t%d\t%d\t%d\t%f\n",
				flavor.ID, flavor.Name, flavor.VCPUs, flavor.RAM, flavor.Disk, flavor.Swap, flavor.RxTxFactor)
		}
		return true, nil
	})

	if err != nil {
		return fmt.Errorf("Unable to list server flavors: %v", err)
	}

	return nil
}
