package render

import "github.com/urfave/cli/v2"

func Root() *cli.Command {
	return &cli.Command{
		Name: "render",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "configuration file to use",
			},
		},
		Subcommands: []*cli.Command{
			HTML(),
			Markdown(),
			PDF(),
		},
	}
}
