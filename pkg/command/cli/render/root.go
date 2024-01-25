package render

import "github.com/urfave/cli/v2"

func Root() *cli.Command {
	return &cli.Command{
		Name: "render",
		Subcommands: []*cli.Command{
			HTML(),
			Markdown(),
			PDF(),
		},
	}
}
