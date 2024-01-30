package main

import (
	"github.com/Bornholm/amatl/pkg/command"
	"github.com/Bornholm/amatl/pkg/command/cli"
)

func main() {
	command.Main(
		"amatl",
		"",
		cli.Root().Subcommands...,
	)
}
