package main

import (
	"fmt"
  "github.com/rackspace/gophercloud"
  "github.com/rackspace/gophercloud/pagination"
  "github.com/rackspace/gophercloud/openstack/compute/v2/servers"
)

func FindServerByName(serviceClient *gophercloud.ServiceClient, name string) (*servers.Server, error) {
	opts := servers.ListOpts{}

  pager := servers.List(serviceClient, opts)
  
  var server *servers.Server

  // Define an anonymous function to be executed on each page's iteration
  err := pager.EachPage(func(page pagination.Page) (bool, error) {
    serverList, err := servers.ExtractServers(page)
    
    if err != nil {
      return false, err
    }
    
    // append the current page's servers to our list
    for _, s := range serverList {
      if s.Name == name {
        server = &s
        return false, err
      }
    }

    return true, err 
  })
  if server == nil {
    return nil, fmt.Errorf("Unable to find server: %s.", name)
  }
  
  return server, err
}
