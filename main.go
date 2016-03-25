package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/svenmueller/nube/Godeps/_workspace/src/github.com/codegangsta/cli"
)

const AppVersion = "0.0.1"

var APIKey string
var OutputFormat string

func init() {
	log.SetFlags(0)
	log.SetPrefix("ctp> ")
}

func main() {
	app := buildApp()
	app.RunAndExitOnError()
}

func buildApp() *cli.App {
	app := cli.NewApp()
	app.Name = "nube"
	app.Version = AppVersion
	app.Usage = "commercetools command line interface for managing Rackspace/AWS resources."
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "api-key, k",
			Value:  "",
			Usage:  "Rackspace API key.",
			EnvVar: "RACKSPACE_API_KEY",
		},
		cli.StringFlag{Name: "format,f", Value: "yaml", Usage: "Format for output."},
		cli.BoolFlag{Name: "debug,d", Usage: "Turn on debug output."},
	}
	app.Before = func(ctx *cli.Context) error {

		if ctx.String("api-key") != "" {
			APIKey = ctx.String("api-key")
		}

		if APIKey == "" && !ctx.Bool("help") && !ctx.Bool("version") && !(ctx.NumFlags() == 0) {
			return errors.New("must provide API Key via RACKSPACE_API_KEY environment variable or via CLI argument.")
		}

		switch ctx.String("format") {
		case "json":
			OutputFormat = ctx.String("format")
		case "yaml":
			OutputFormat = ctx.String("format")
		default:
			return fmt.Errorf("invalid output format: %q, available output options: json, yaml.", ctx.String("format"))
		}

		return nil
	}
	app.Commands = []cli.Command{
		ServersCommand,
		DNSCommand,
	}

	return app
}
