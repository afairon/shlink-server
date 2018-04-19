DIR := bin
BIN := shlink
VERSION := v1.0.0
PLATFORMS := darwin linux freebsd windows
os = $(word 1, $@)

GOVERSION := `go version | cut -d ' ' -f 3`
GOPLATFORM := `go version | cut -d ' ' -f 4`

ARCH ?= amd64

.PHONY: default
default: $(PLATFORMS)

vendor: Gopkg.toml
	go get -u -v github.com/golang/dep/cmd/dep
	go get -u -v github.com/pquerna/ffjson
	dep ensure

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: generate
generate:
	go generate ./...

.PHONY: test
test: fmt
	go test -cover ./...

.PHONY: $(PLATFORMS)
$(PLATFORMS): generate fmt vet
	GOOS=$(os) GOARCH=$(ARCH) go build -ldflags "-s -w -X main.version=$(VERSION) -X main.goVersion=$(GOVERSION) -X main.goPlatform=$(GOPLATFORM)" -o $(DIR)/$(ARCH)/$(BIN)-$(VERSION)-$(os)-$(ARCH)

.PHONY: clean
clean:
	rm -rf bin/
	rm -rf models/ffjson*
	rm -rf models/*_ffjson.go