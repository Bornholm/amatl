package main

import (
	"github.com/Bornholm/amatl/pkg/command"
	"github.com/Bornholm/amatl/pkg/command/cli"
)

func main() {
	command.Main(
		"amatl",
		"a markdown to markdown/html/pdf compiler",
		cli.Root().Subcommands...,
	)
}
