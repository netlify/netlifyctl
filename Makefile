.PONY: all build deps image test

help: ## Show this help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

all: deps test build ## Run the tests and build the binary.

os = darwin
arch = amd64

build: test
	@echo "Making netlifyctl for $(os)/$(arch)"
	GOOS=$(os) GOARCH=$(arch) go build -ldflags "-X github.com/netlify/netlifyctl/commands.Version=`git rev-parse HEAD`"

build_linux: override os=linux
build_linux: build

package: build
	tar -czf netlifyctl-$(os)-$(arch).tar.gz netlifyctl

package_linux: override os=linux
package_linux: package

deps: ## Install dependencies.
	go get -u github.com/Masterminds/glide && glide install

test: ## Run tests.
	go test -v `go list ./... | grep -v /vendor/`
