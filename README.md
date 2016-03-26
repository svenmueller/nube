# nube CLI

[![Build Status](https://travis-ci.org/svenmueller/nube.svg?branch=master)](https://travis-ci.org/svenmueller/nube)

`nube` is a CLI for managing commercetools cloud resources on Rackspace/AWS. The word "nube" is taken from the spanish language and simply means "cloud".

## Installation

Using `go get`:

```
$ go get github.com/svenmueller/nube
```

Or clone and build yourself:

```
$ git clone
$ go get
```

## Usage

```
Manage Rackspace and AWS resources

Usage:
  nube [command]

Available Commands:
  dns         Manage AWS Route53 DNS resources
  servers     Manage Rackspace Cloud Server resources

Flags:
      --config string   config file (default is $HOME/.nube.yaml)
  -h, --help            help for nube
  -o, --output string   output format [yaml|json] (default "yaml")
  -t, --toggle          Help message for toggle

Use "nube [command] --help" for more information about a command.
```

### Examples
```bash
# List all server instances
$ nube servers instance list
```

```bash
# Create 3 new server instances (using defaults for flavor/size etc.)
$ nube instance create server1 server2 server3
```

```bash
# Destroy 3 server instances using ID or name
$ nube instance destroy cdf0bb56 server2.example.com server3.example.com
```
