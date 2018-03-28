DIR := bin
VERSION := v1.0.0
PLATFORMS := darwin linux freebsd windows
os = $(word 1, $@)

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

.PHONY: generate
generate:
	go generate ./...

.PHONY: test
test:
	go test -cover ./...

.PHONY: bench
bench:
	go test -bench=. -benchmem ./...

.PHONY: $(PLATFORMS)
$(PLATFORMS): generate fmt
	GOOS=$(os) GOARCH=$(ARCH) go build -ldflags "-s -w" -o $(DIR)/$(ARCH)/short-$(VERSION)-$(os)-$(ARCH)

.PHONY: clean
clean:
	rm -rf bin/
	rm -f models/*_ffjson.go