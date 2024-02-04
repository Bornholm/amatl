GORELEASER_ARGS ?= --snapshot --clean

build:
	CGO_ENABLED=0 go build -o bin/amatl ./cmd/amatl

release:
	goreleaser $(GORELEASER_ARGS)