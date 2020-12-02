package main

import (
	"fmt"

	"github.com/aserto-demo/aserto-seed-auth0/pkg/auth0"
	ver "github.com/aserto-demo/aserto-seed-auth0/pkg/version"

	_ "github.com/joho/godotenv/autoload"
	flag "github.com/spf13/pflag"
)

func main() {
	var (
		seed    *bool   = flag.Bool("seed", false, "seed")
		reset   *bool   = flag.Bool("reset", false, "reset")
		spew    *bool   = flag.Bool("spew", false, "spew")
		dryrun  *bool   = flag.Bool("dryrun", false, "dryrun")
		version *bool   = flag.Bool("version", false, "version")
		input   *string = flag.String("input", "", "inputfile")
	)

	flag.Parse()

	if *version {
		fmt.Printf("%s\n", ver.GetInfo().String())

		return
	}

	if input == nil || *input == "" {
		fmt.Printf("--input not set\n")

		return
	}

	helper := auth0.NewHelper(*input)
	if err := helper.Init(); err != nil {
		fmt.Printf("error %+v\n", err)

		return
	}

	helper.Spew(*spew)
	helper.Dryrun(*dryrun)

	switch {
	case *seed:
		if err := helper.Seed(); err != nil {
			fmt.Printf("error %+v\n", err)
		}
		return

	case *reset:
		if err := helper.Reset(); err != nil {
			fmt.Printf("error %+v\n", err)
		}
		return

	default:
		fmt.Printf("no valid arguments passed in")
		return
	}
}
