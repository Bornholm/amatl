package render

import (
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

func PDF() *cli.Command {
	return &cli.Command{
		Name:  "pdf",
		Flags: withCommonFlags(),
		Action: func(ctx *cli.Context) error {
			return errors.New("unimplemented")
		},
	}
}
