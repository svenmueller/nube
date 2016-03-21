# nube CLI

[![Build Status](https://travis-ci.org/svenmueller/nube.svg)](https://travis-ci.org/svenmueller/nube.svg?branch=master)

`nube` is a CLI for managing commercetools cloud resources on Rackspace/AWS. The word "nube" is taken from the spanish language and simply means "cloud".

## Installation

Clone and build yourself:

```
$ git clone
$ go get
```

Or using `go get`:

```
$ go get github.com/svenmueller/nube
```

## Usage

More details:

```
NAME:
   nube - commercetools command line interface for managing Rackspace/AWS resources.

USAGE:
   nube [global options] command [command options] [arguments...]

VERSION:
   0.0.1

COMMANDS:
   servers, s	Server commands.
   help, h	Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --api-key, -k 	Rackspace API key. [$RACKSPACE_API_KEY]
   --format, -f "yaml"	Format for output.
   --debug, -d		Turn on debug output.
   --help, -h		show help
   --version, -v	print the version
```

### Servers
```
NAME:
   nube servers - Server commands.

USAGE:
   nube servers [global options] command [command options] [arguments...]

VERSION:
   0.0.1

COMMANDS:
   list, l	List all available servers.
   create, c	Create a new server.
   destroy, d	[--id | <name>] Destroy a server.
   help, h	Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --rackspace-username 		The Rackspace API username. [$RACKSPACE_USERNAME]
   --rackspace-api-key 			The Rackspace API key. [$RACKSPACE_API_KEY]
   --rackspace-region-name "LON"	The Rackspace region name. [$RACKSPACE_REGION_NAME]
   --help, -h				show help
```
