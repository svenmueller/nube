package flavors

import (
	"github.com/svenmueller/nube/Godeps/_workspace/src/github.com/rackspace/gophercloud"
)

func getURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("flavors", id)
}

func listURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("flavors", "detail")
}
