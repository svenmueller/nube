package util

import (
	"github.com/svenmueller/nube/Godeps/_workspace/src/github.com/rackspace/gophercloud"
	"github.com/svenmueller/nube/Godeps/_workspace/src/github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/svenmueller/nube/Godeps/_workspace/src/github.com/rackspace/gophercloud/pagination"
	"github.com/svenmueller/nube/Godeps/_workspace/src/github.com/rackspace/gophercloud/rackspace"
)

func newRackspaceProviderClient(username string, apiKey string) (*gophercloud.ProviderClient, error) {
	authOpts := gophercloud.AuthOptions{
		Username: username,
		APIKey:   apiKey,
	}

	provider, err := rackspace.AuthenticatedClient(authOpts)

	return provider, err
}

func NewRackspaceService(username string, apiKey string, region string) (*gophercloud.ServiceClient, error) {

	provider, err := newRackspaceProviderClient(username, apiKey)

	if err != nil {
		return nil, err
	}

	client, err := rackspace.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: region,
	})

	return client, err
}

func ListAllServers(serviceClient *gophercloud.ServiceClient) ([]servers.Server, error) {
	opts := servers.ListOpts{}
	pager := servers.List(serviceClient, opts)

	var list []servers.Server

	// Define an anonymous function to be executed on each page's iteration
	err := pager.EachPage(func(page pagination.Page) (bool, error) {
		serverList, ExtractServersError := servers.ExtractServers(page)

		if ExtractServersError != nil {
			return false, ExtractServersError
		}
		list = append(list, serverList...)

		return true, nil
	})

	if err != nil {
		return nil, err
	}

	return list, err
}
