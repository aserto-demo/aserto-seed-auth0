package cmd

import (
	"github.com/aserto-demo/aserto-seed-auth0/pkg/auth0"
	"github.com/aserto-demo/aserto-seed-auth0/pkg/config"
	"github.com/urfave/cli/v2"
)

func seedHandler(c *cli.Context) (err error) {
	// get config from context
	cfg := config.FromContext(c.Context)

	mgr := auth0.NewManager(
		cfg,
		c.Path(flagInputFile),
	)

	if err := mgr.Init(); err != nil {
		return err
	}

	mgr.Dryrun(c.Bool(flagDryRun))
	mgr.Spew(c.Bool(flagSpew))

	if err := mgr.Seed(); err != nil {
		return err
	}

	return nil
}
