package main

import (
	"github.com/Bornholm/amatl/pkg/command"
	"github.com/Bornholm/amatl/pkg/command/cli"

	// Import resolvers
	_ "github.com/Bornholm/amatl/pkg/resolver/file"
	_ "github.com/Bornholm/amatl/pkg/resolver/http"
	_ "github.com/Bornholm/amatl/pkg/resolver/stdin"
)

func main() {
	command.Main(
		"amatl",
		"a markdown to markdown/html/pdf compiler",
		cli.Root().Subcommands...,
	)
}
