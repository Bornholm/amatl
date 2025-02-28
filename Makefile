GORELEASER_ARGS ?= --snapshot --clean

build:
	CGO_ENABLED=0 go build -o bin/amatl ./cmd/amatl

release:
	goreleaser $(GORELEASER_ARGS)

examples: example-document example-presentation

example-%: build
	bin/amatl render html --html-layout amatl://$*.html -o ./examples/$*/$*.html ./examples/$*/$*.md
	bin/amatl render pdf --html-layout amatl://$*.html -o ./examples/$*/$*.pdf ./examples/$*/$*.md