package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

// SeedCommand -- seed, add entries from the seed file.
func SeedCommand() *cli.Command {
	return &cli.Command{
		Name:  "seed",
		Usage: "seed",
		Flags: []cli.Flag{
			InputFileFlag(),
			SpewFlag(),
			DryrunFlag(),
		},
		Before: configRetriever,
		Action: seedHandler,
	}
}

// ResetCommand -- reset state, removes all entries added from the seed file.
func ResetCommand() *cli.Command {
	return &cli.Command{
		Name:  "reset",
		Usage: "reset",
		Flags: []cli.Flag{
			InputFileFlag(),
			SpewFlag(),
			DryrunFlag(),
		},
		Before: configRetriever,
		Action: resetHandler,
	}
}

// VersionCommand -- version command definition.
func VersionCommand() *cli.Command {
	return &cli.Command{
		Name:  "version",
		Usage: "display verion information",
		Action: func(c *cli.Context) (err error) {
			fmt.Fprintf(c.App.Writer, "%s - %s\n",
				c.App.Name,
				c.App.Version,
			)
			return nil
		},
	}
}
