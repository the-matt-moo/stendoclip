APP := stendoclip
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -X main.version=$(VERSION)

.PHONY: build test run release clean

build:
	go build -ldflags "$(LDFLAGS)" -o $(APP).exe ./cmd/stendoclip

test:
	go test ./...

run:
	go run -ldflags "$(LDFLAGS)" ./cmd/stendoclip

release:
	go build -trimpath -ldflags "-s -w -H=windowsgui $(LDFLAGS)" -o $(APP).exe ./cmd/stendoclip

clean:
	go clean
	$(RM) $(APP).exe
