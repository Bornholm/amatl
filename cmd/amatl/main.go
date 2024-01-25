package main

import (
	"forge.cadoles.com/wpetit/amatl/pkg/command"
	"forge.cadoles.com/wpetit/amatl/pkg/command/cli"
)

func main() {
	command.Main(
		"amatl",
		"",
		cli.Root().Subcommands...,
	)
}
