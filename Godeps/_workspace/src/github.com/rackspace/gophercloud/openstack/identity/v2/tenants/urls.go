package tenants

import "github.com/svenmueller/nube/Godeps/_workspace/src/github.com/rackspace/gophercloud"

func listURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("tenants")
}
