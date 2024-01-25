package command

import (
	"fmt"
	"log/slog"
	"os"
	"sort"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

func Main(name string, usage string, commands ...*cli.Command) {
	app := &cli.App{
		Name:     name,
		Usage:    usage,
		Commands: commands,
		Before: func(ctx *cli.Context) error {
			workdir := ctx.String("workdir")
			// Switch to new working directory if defined
			if workdir != "" {
				if err := os.Chdir(workdir); err != nil {
					return errors.Wrap(err, "could not change working directory")
				}
			}

			logLevel := ctx.String("log-level")

			var programLevel = new(slog.LevelVar)

			switch logLevel {
			case "debug":
				programLevel.Set(slog.LevelDebug)
			case "info":
				programLevel.Set(slog.LevelInfo)
			case "warn":
				programLevel.Set(slog.LevelWarn)
			case "error":
				programLevel.Set(slog.LevelError)
			}

			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "workdir",
				Value:   "",
				EnvVars: []string{"AMATL_WORKDIR"},
				Usage:   "The working directory",
			},
			&cli.BoolFlag{
				Name:    "debug",
				EnvVars: []string{"AMATL_DEBUG"},
				Usage:   "Enable debug mode",
			},
			&cli.StringFlag{
				Name:    "log-level",
				EnvVars: []string{"AMATL_LOG_LEVEL"},
				Usage:   "Set logging level",
				Value:   "info",
			},
		},
	}

	app.ExitErrHandler = func(ctx *cli.Context, err error) {
		if err == nil {
			return
		}

		debug := ctx.Bool("debug")

		if !debug {
			slog.ErrorContext(ctx.Context, err.Error())
		} else {
			slog.ErrorContext(ctx.Context, fmt.Sprintf("%+v", err))
		}
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}
