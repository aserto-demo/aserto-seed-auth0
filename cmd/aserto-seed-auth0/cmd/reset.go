package cmd

import (
	"github.com/aserto-demo/aserto-seed-auth0/pkg/auth0"
	"github.com/urfave/cli/v2"
)

func resetHandler(c *cli.Context) (err error) {
	helper := auth0.NewHelper(
		c.Path(flagInputFile),
	)

	if err := helper.Init(); err != nil {
		return err
	}

	helper.Dryrun(c.Bool(flagDryRun))
	helper.Spew(c.Bool(flagSpew))

	if err := helper.Reset(); err != nil {
		return err
	}

	return nil
}
