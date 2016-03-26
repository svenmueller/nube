package util

import (
	"fmt"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/rackspace/gophercloud/pagination"
	"github.com/rackspace/gophercloud/rackspace"
	"github.com/spf13/viper"
)

func newRackspaceProviderClient() (*gophercloud.ProviderClient, error) {
	authOpts := gophercloud.AuthOptions{
		Username: viper.GetString("rackspace-username"),
		APIKey:   viper.GetString("rackspace-api-key"),
	}

	provider, err := rackspace.AuthenticatedClient(authOpts)

	return provider, err
}

func NewRackspaceService() (*gophercloud.ServiceClient, error) {

	provider, err := newRackspaceProviderClient()

	if err != nil {
		return nil, err
	}

	client, err := rackspace.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: viper.GetString("rackspace-region"),
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
		return nil, fmt.Errorf("Unable to build list of servers: %v", err)
	}

	return list, err
}
