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

	"github.com/svenmueller/nube/Godeps/_workspace/src/github.com/rackspace/gophercloud/openstack/compute/v2/images"
	"github.com/svenmueller/nube/Godeps/_workspace/src/github.com/rackspace/gophercloud/pagination"
	"github.com/svenmueller/nube/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/svenmueller/nube/common"
	"github.com/svenmueller/nube/util"
)

var servers_image_listCmd = &cobra.Command{
	Use:   "list",
	Short: "List server images",
	Run: func(cmd *cobra.Command, args []string) {
		err := serversImageList(cmd, args)
		common.HandleError(err, cmd)
	},
}

func init() {
	servers_imageCmd.AddCommand(servers_image_listCmd)
}

func serversImageList(cmd *cobra.Command, args []string) error {

	// We have the option of filtering the server list. If we want the full
	// collection, leave it as an empty struct
	opts := images.ListOpts{}

	rackspaceServiceClient, err := util.NewRackspaceService(Cfg.GetString("rackspace-username"), Cfg.GetString("rackspace-api-key"), Cfg.GetString("rackspace-region"))

	if err != nil {
		return fmt.Errorf("Unable to establish connection: %v", err)
	}

	// Retrieve a pager (i.e. a paginated collection)
	pager := images.ListDetail(rackspaceServiceClient, opts)

	// Define an anonymous function to be executed on each page's iteration
	err = pager.EachPage(func(page pagination.Page) (bool, error) {
		imageList, ExtractImagesError := images.ExtractImages(page)

		if ExtractImagesError != nil {
			return false, ExtractImagesError
		}

		cliOut := util.NewCLIOutput()
		defer cliOut.Flush()
		cliOut.Header("ID", "Name", "MinRAM", "MinRAM", "Status", "Created", "Updated", "Progress")
		for _, image := range imageList {
			cliOut.Writeln("%s\t%s\t%d\t%d\t%s\t%s\t%s\t%d\n",
				image.ID, image.Name, image.MinRAM, image.MinRAM, image.Status, image.Created, image.Updated, image.Progress)
		}
		return true, nil
	})

	if err != nil {
		return fmt.Errorf("Unable to list server images: %v", err)
	}

	return nil
}
