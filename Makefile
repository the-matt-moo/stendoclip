APP := stendoclip
VERSION ?= $(shell cat VERSION 2>/dev/null || echo dev)
LDFLAGS := -X main.version=$(VERSION)
RSRC := github.com/akavel/rsrc@v0.10.2

.PHONY: build test run resources release clean

build:
	go build -ldflags "$(LDFLAGS)" -o $(APP).exe ./cmd/stendoclip

test:
	go test ./...

run:
	go run -ldflags "$(LDFLAGS)" ./cmd/stendoclip

resources:
	go run $(RSRC) -arch amd64 -ico assets/watergun_icon.ico,assets/cow.ico -manifest assets/stendoclip.manifest -o cmd/stendoclip/rsrc_windows_amd64.syso
	go run $(RSRC) -arch 386 -ico assets/watergun_icon.ico,assets/cow.ico -manifest assets/stendoclip.manifest -o cmd/stendoclip/rsrc_windows_386.syso

release:
	go build -trimpath -ldflags "-s -w -H=windowsgui $(LDFLAGS)" -o $(APP).exe ./cmd/stendoclip

clean:
	go clean
	$(RM) $(APP).exe
