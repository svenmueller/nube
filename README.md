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
NAME:
   nube - commercetools command line interface for managing Rackspace/AWS resources.

USAGE:
   nube [global options] command [command options] [arguments...]

VERSION:
   0.0.1

COMMANDS:
   servers, s	Server commands.
   dns, d	DNS commands.
   help, h	Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --api-key, -k 	Rackspace API key. [$RACKSPACE_API_KEY]
   --format, -f "yaml"	Format for output.
   --debug, -d		Turn on debug output.
   --help, -h		show help
   --version, -v	print the version
```
