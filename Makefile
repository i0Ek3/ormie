.PHONY: build test clean

GO=go

build:
	@$(GO) build -o ormie-darwin

test:
	@$(GO) test -v .

clean:
	@rm ormie-darwin
