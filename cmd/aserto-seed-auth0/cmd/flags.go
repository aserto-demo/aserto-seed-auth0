package cmd

import "github.com/urfave/cli/v2"

// InputFileFlag -- input file flag.
func InputFileFlag() cli.Flag {
	return &cli.PathFlag{
		Name:     flagInputFile,
		Aliases:  []string{},
		Usage:    usageInputFile,
		Required: true,
	}
}

// DryrunFlag -- dryrun flag.
func DryrunFlag() cli.Flag {
	return &cli.BoolFlag{
		Name:    flagDryRun,
		Aliases: []string{},
		Usage:   usageDryRun,
	}
}

// SpewFlag -- spew flag.
func SpewFlag() cli.Flag {
	return &cli.BoolFlag{
		Name:    flagSpew,
		Aliases: []string{},
		Usage:   usageSpew,
	}
}
