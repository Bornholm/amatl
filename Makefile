GORELEASER_ARGS ?= --snapshot --clean

AMATL_LATEST_VERSION ?= $(shell git describe --tags --abbrev=0)

build:
	CGO_ENABLED=0 go build -o bin/amatl ./cmd/amatl

release:
	goreleaser $(GORELEASER_ARGS)

examples: example-document example-presentation

example-%: build
	bin/amatl --log-level debug render html --html-layout amatl://$*.html -o ./examples/$*/$*.html ./examples/$*/$*.md
	bin/amatl --log-level debug render pdf --html-layout amatl://$*.html -o ./examples/$*/$*.pdf ./examples/$*/$*.md

website: build
	mkdir -p dist/website
	echo '{"amatlVersion":"$(AMATL_LATEST_VERSION)"}' | bin/amatl \
		--log-level debug \
		render html \
		--link-replacements "file://$(PWD)::https://github.com/Bornholm/amatl/blob/$(AMATL_LATEST_VERSION)" \
		--vars stdin:// \
		--html-layout amatl://website.html \
		-o ./dist/website/index.html \
		./misc/website/index.md

include misc/make/*.mk