package tokens

import "github.com/svenmueller/nube/Godeps/_workspace/src/github.com/rackspace/gophercloud"

func tokenURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("auth", "tokens")
}
