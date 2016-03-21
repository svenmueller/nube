package main

import (
  "fmt"
  "io/ioutil"
	"log"
  "errors"
  "strings"
  "os"
  "github.com/codegangsta/cli"
  "github.com/rackspace/gophercloud"
  "github.com/rackspace/gophercloud/rackspace"
  "github.com/rackspace/gophercloud/pagination"
  "github.com/rackspace/gophercloud/openstack/compute/v2/servers"
  
//  "github.com/mitchellh/goamz/aws"
//  "github.com/mitchellh/goamz/route53"
  
  "github.com/docker/docker/pkg/namesgenerator"
)

var RackspaceUsername string
var RackspaceAPIKey string
var RackspaceRegionName string
var RackspaceServiceClient *gophercloud.ServiceClient

var ServersCommand = cli.Command{
	Name:    "servers",
	Aliases: []string{"s"},
	Usage:   "Server commands.",
  Flags: []cli.Flag{
    cli.StringFlag{Name: "rackspace-username", Value: "", Usage: "The Rackspace API username.", EnvVar: "RACKSPACE_USERNAME", Destination: &RackspaceUsername},
    cli.StringFlag{Name: "rackspace-api-key", Value: "", Usage: "The Rackspace API key.", EnvVar: "RACKSPACE_API_KEY", Destination: &RackspaceAPIKey},
    cli.StringFlag{Name: "rackspace-region-name", Value: "LON", Usage: "The Rackspace region name.", EnvVar: "RACKSPACE_REGION_NAME", Destination: &RackspaceRegionName},
  },
	Action:  serversList,
  Before: func(ctx *cli.Context) error {

    if RackspaceUsername == "" {
      return errors.New("You must provide the Rackspace API username via RACKSPACE_USERNAME environment variable or via CLI argument.")
    }

    if RackspaceAPIKey == "" {
      return errors.New("You must provide the Rackspace API Key via RACKSPACE_API_KEY environment variable or via CLI argument.")
    }
    
    RackspaceServiceClient = newServiceClient(RackspaceUsername, RackspaceAPIKey, RackspaceRegionName)

    return nil
  },
	Subcommands: []cli.Command{
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "List all available servers.",
			Action:  serversList,
		},
    {
      Name:    "create",
      Aliases: []string{"c"},
      Usage:   "Create a new server.",
      Action:  serversCreate,
      Flags: []cli.Flag{
        cli.StringFlag{Name: "domain, d", Value: "ct-app.com", Usage: "Domain name to append to the hostname. (e.g. server01.example.com)"},
        cli.BoolFlag{Name: "add-random-name", Usage: "Append random name to server name. (e.g. server01-adjective-surname)"},
        cli.BoolFlag{Name: "add-region", Usage: "Append region to hostname. (e.g. server01.lon)"},
        cli.StringFlag{Name: "user-data, u", Value: "", Usage: "User data for creating server."},
        cli.StringFlag{Name: "user-data-file, uf", Value: "", Usage: "A path to a file for user data."},
        cli.StringFlag{Name: "flavor, f", Value: "1 GB Performance", Usage: "Flavor of server."},
        cli.StringFlag{Name: "image, i", Value: "Ubuntu 14.04 LTS (Trusty Tahr) (PVHVM)", Usage: "Image name of server."},
        cli.StringFlag{Name: "region, r", Value: "lon", Usage: "Region of server."},
//        cli.BoolFlag{Name: "backups, b", Usage: "Turn on backups."},
        cli.BoolFlag{Name: "wait-for-active", Usage: "Don't return until the create has succeeded or failed."},
      },
    },
    {
			Name:    "destroy",
			Aliases: []string{"d"},
			Usage:   "[--id | <name>] Destroy a server.",
			Action:  serversDestroy,
			Flags: []cli.Flag{
				cli.StringFlag{Name: "id", Usage: "ID for server. (e.g. 7a54885d-fd8a-4616-8e06-0bea3692582c)"},
			},
		},
	},
}

func newProviderClient(username string, apiKey string) *gophercloud.ProviderClient {
  authOpts := gophercloud.AuthOptions{
    Username: username,
    APIKey: apiKey,
  }

  provider, err := rackspace.AuthenticatedClient(authOpts)

  if err != nil {
    log.Fatal(err)
  }
  
  return provider
}

func newServiceClient(username string, apiKey string, region string) *gophercloud.ServiceClient {

  provider := newProviderClient(username, apiKey)

  client, err := rackspace.NewComputeV2(provider, gophercloud.EndpointOpts{
  	Region: region,
  })

  if err != nil {
    log.Fatal(err)
  }
  
  return client
}

func serversList(ctx *cli.Context) {

  // We have the option of filtering the server list. If we want the full
  // collection, leave it as an empty struct
  opts := servers.ListOpts{}

  // Retrieve a pager (i.e. a paginated collection)
  pager := servers.List(RackspaceServiceClient, opts)

  // Define an anonymous function to be executed on each page's iteration
  err := pager.EachPage(func(page pagination.Page) (bool, error) {
  	serverList, err := servers.ExtractServers(page)

    cliOut := NewCLIOutput()
    defer cliOut.Flush()
    cliOut.Header("ID", "Name", "Status")
    for _, server := range serverList {
      cliOut.Writeln("%s\t%s\t%s\n",
        server.ID, server.Name, server.Status)
    }
    return true, err 
  })

	if err != nil {
		log.Fatal(err)
	}

}

func serversCreate(ctx *cli.Context) {
  
  if len(ctx.Args()) != 1 {
    log.Fatal("Error: Must provide name for server.")
  }

  // Add domain to end if available.
  serverName := ctx.Args().First()
  if ctx.Bool("add-random-name") {
    randomServerName := strings.Replace(namesgenerator.GetRandomName(0), "_", "-", -1)
    serverName = fmt.Sprintf("%s-%s", serverName, randomServerName)
  }
  if ctx.Bool("add-region") {
    serverName = fmt.Sprintf("%s.%s", serverName, ctx.String("region"))
  }
  if ctx.String("domain") != "" {
    serverName = fmt.Sprintf("%s.%s", serverName, ctx.String("domain"))
  }
  
  userData := ""
  userDataPath := ctx.String("user-data-file")
  if userDataPath != "" {
    file, err := os.Open(userDataPath)
    if err != nil {
      log.Fatalf("Error opening user data file: %s.", err)
    }

    userDataFile, err := ioutil.ReadAll(file)
    if err != nil {
      log.Fatalf("Error reading user data file: %s.", err)
    }
    userData = string(userDataFile)
  } else {
    userData = ctx.String("user-data")
  }

  server, err := servers.Create(RackspaceServiceClient, servers.CreateOpts{
    Name:      serverName,
    ImageName:  ctx.String("image"),
    FlavorName: ctx.String("flavor"),
    UserData: []byte(userData),
  }).Extract()
  
  log.Println("Creating server. This process will take some seconds.")
  
  if err != nil {
    log.Fatalf("Unable to create server: %s.", err)
  }

  if ctx.Bool("wait-for-active") {
    log.Println("Waiting for server to be in state 'ACTIVE'.")
    err = servers.WaitForStatus(RackspaceServiceClient, server.ID, "ACTIVE", 600)

    if err != nil {
      log.Fatal(err)
    }
    
    server, err = servers.Get(RackspaceServiceClient, server.ID).Extract()
  }

  WriteOutput(server)
}

func serversDestroy(ctx *cli.Context) {
	if ctx.Int("id") == 0 && len(ctx.Args()) != 1 {
		log.Fatal("Error: Must provide ID or name for server to destroy.")
	}

	id := ctx.String("id")
	if id == "" {
		server, err := FindServerByName(RackspaceServiceClient, ctx.Args()[0])
		if err != nil {
			log.Fatal(err)
		} else {
			id = server.ID
		}
	}

	server, err := servers.Get(RackspaceServiceClient, id).Extract()
	if err != nil {
		log.Fatalf("Unable to find server: %s.", err)
	}

	result := servers.Delete(RackspaceServiceClient, server.ID)
  
  if result.ErrResult.Result.Err != nil {
    log.Fatalf("Unable to delete server: %s.", result.ErrResult.Result.Err)
  }

	log.Fatalf("Server %s destroyed.", server.Name)
}
