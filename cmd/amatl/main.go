package main

import (
	"github.com/Bornholm/amatl/pkg/command"
	"github.com/Bornholm/amatl/pkg/command/cli"
	"github.com/Bornholm/amatl/pkg/resolver"

	// Import resolvers
	_ "github.com/Bornholm/amatl/pkg/resolver/all"
)

var (
	version = "unknown"
)

func main() {
	resolver.SetDefault("file")

	command.Main(
		"amatl",
		version,
		"a markdown to markdown/html/pdf compiler",
		cli.Root().Subcommands...,
	)
}
