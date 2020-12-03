package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/aserto-demo/aserto-seed-auth0/cmd/aserto-seed-auth0/cmd"
	"github.com/aserto-demo/aserto-seed-auth0/pkg/version"
	_ "github.com/joho/godotenv/autoload"
	"github.com/urfave/cli/v2"
)

const (
	appName  = "aserto-seed-auth0"
	appUsage = "seed Auth0 user data"
)

func main() {
	log.SetOutput(ioutil.Discard)

	app := cli.NewApp()
	app.Name = appName
	app.Usage = appUsage
	app.HideVersion = true
	app.HideHelpCommand = true
	app.Version = version.GetInfo().String()
	app.Flags = []cli.Flag{}
	app.Commands = []*cli.Command{
		cmd.SeedCommand(),
		cmd.ResetCommand(),
		cmd.VersionCommand(),
	}

	ctx := context.Background()

	if err := app.RunContext(ctx, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error %v\n", err)
		os.Exit(1)
	}
}
