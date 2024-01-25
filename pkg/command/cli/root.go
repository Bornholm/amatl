package cli

import (
	"forge.cadoles.com/wpetit/amatl/pkg/command/cli/render"
	"github.com/urfave/cli/v2"
)

func Root() *cli.Command {
	return &cli.Command{
		Name: "cli",
		Subcommands: []*cli.Command{
			render.Root(),
		},
	}
}
