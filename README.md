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

Or using Docker:

```bash
docker run --rm -v ~/.nube.yaml:/root/.nube.yaml svenmueller/nube
```

Set an alias for repeated calls
```
alias nube='docker run --rm -v ~/.nube.yaml:/root/.nube.yaml svenmueller/nube'
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
nube servers instance list
```

```bash
# Create 3 server instances (using defaults for flavor/size etc.)
nube servers instance create server1 server2 server3
```

```bash
# Create 1 server instance using custom flavor/image and user-data file
nube servers instance create server1 \
--flavor "2 GB Performance" \
--image "Ubuntu 14.04 LTS (Trusty Tahr) (PVHVM)" \
--user-data-file cloud-config.yaml
```

```bash
# Destroy 3 server instances using ID or name
nube servers instance destroy cdf0bb56 server2.example.com server3.example.com
```

```bash
# List all hosted zones
nube dns zones list
```

```bash
# List all resource records of for given hosted zone ID
$ nube dns records list /hostedzone/XXXXXX
```

## Configuration

There are multiple ways to set values for `nube` CLI. All values are looked up in the following order:

- Configuration file (default path is `~/.nube.yaml`)
- Environment variable (upper case, prefix `NUBE_`, `-` replaced by `_`, e.g. `NUBE_RACKSPACE_USERNAME`)
- Flags (e.g. `--rackspace-username`)

Example configuration file (`~/.nube.yaml`)
```yaml
default:
  rackspace-username: bart.simpson
  rackspace-api-key: 12121212121
  rackspace-region: LON

  aws-access-key-id: 113131313131313
  aws-secret-access-key: 00000000000

  # to avoid having to pass --hosted-zone-id option to create resource record set
  # hosted-zone-id: /hostedzone/XXXXXXX
```

### Named Profiles

The nube CLI supports named profiles stored in the config files. You can configure additional profiles by adding entries to the config files.

The following example shows a config file with two profiles:
```yaml
default:
  ...
  ...
production:
  ...
  ...
```

To use a named profile, add the `--profile` option to your command. The following example lists server instances using the `production` profile from example above.

```bash
# List all hosted zones
nube dns zones list
```
By default, the named profile `default` is used if no `--profile` option is found.
