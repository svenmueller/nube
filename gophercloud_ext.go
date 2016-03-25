package main

import (
	"fmt"

	"github.com/svenmueller/nube/Godeps/_workspace/src/github.com/rackspace/gophercloud"
	"github.com/svenmueller/nube/Godeps/_workspace/src/github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/svenmueller/nube/Godeps/_workspace/src/github.com/rackspace/gophercloud/pagination"
)

func FindServerByName(serviceClient *gophercloud.ServiceClient, name string) (*servers.Server, error) {
	opts := servers.ListOpts{}

	pager := servers.List(serviceClient, opts)

	var server *servers.Server

	// Define an anonymous function to be executed on each page's iteration
	err := pager.EachPage(func(page pagination.Page) (bool, error) {
		serverList, ExtractServersError := servers.ExtractServers(page)

		if ExtractServersError != nil {
			return false, ExtractServersError
		}

		// append the current page's servers to our list
		for _, s := range serverList {
			if s.Name == name {
				server = &s
				return false, nil
			}
		}

		return true, nil
	})
	if server == nil {
		return nil, fmt.Errorf("Unable to find server: %s.", name)
	}

	return server, err
}
