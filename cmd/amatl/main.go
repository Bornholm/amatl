package main

import (
	"github.com/Bornholm/amatl/pkg/command"
	"github.com/Bornholm/amatl/pkg/command/cli"
	"github.com/Bornholm/amatl/pkg/resolver"

	// Import resolvers
	_ "github.com/Bornholm/amatl/pkg/resolver/all"
)

func main() {
	resolver.SetDefault("file")

	command.Main(
		"amatl",
		"a markdown to markdown/html/pdf compiler",
		cli.Root().Subcommands...,
	)
}
