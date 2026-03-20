package mcp

import (
	"os"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

// Root returns the MCP CLI command.
func Root() *cli.Command {
	return &cli.Command{
		Name:  "mcp",
		Usage: "Model Context Protocol server commands",
		Subcommands: []*cli.Command{
			{
				Name:  "serve",
				Usage: "Start the MCP stdio server (JSON-RPC 2.0 over stdin/stdout)",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "workspace",
						Aliases: []string{"w"},
						Usage:   "Root directory accessible to the MCP server (default: current directory)",
						EnvVars: []string{"AMATL_MCP_WORKSPACE"},
					},
				},
				Action: func(ctx *cli.Context) error {
					workspace := ctx.String("workspace")
					if workspace == "" {
						var err error
						workspace, err = os.Getwd()
						if err != nil {
							return errors.Wrap(err, "could not determine current directory")
						}
					}
					if err := serve(workspace); err != nil {
						return errors.Wrap(err, "MCP server error")
					}
					return nil
				},
			},
		},
	}
}
