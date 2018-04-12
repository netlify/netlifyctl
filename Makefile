.PHONY: all build deps image release test

help: ## Show this help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

all: deps test build ## Run tests and build the binary.

os = darwin
arch = amd64

build: test
	@echo "Making netlifyctl for $(os)/$(arch)"
	GOOS=$(os) GOARCH=$(arch) CGO_ENABLED=0 go build -ldflags "-X github.com/netlify/netlifyctl/commands.Version=${TAG} -X github.com/netlify/netlifyctl/commands.SHA=`git rev-parse HEAD`"

build_linux: override os=linux ## Build the binary for Linux hosts.
build_linux: build

build_windows: override os=windows ## Build the binary for Windows hosts.
build_windows: build

deps: ## Install dependencies.
	go get -u github.com/Masterminds/glide && glide install

package: build ## Package binary for the default OS.
	tar -czf netlifyctl-$(os)-$(arch).tar.gz netlifyctl

package_linux: override os=linux ## Package Linux binary.
package_linux: package

package_windows: override os=windows ## Package Windows binary.
package_windows: package

release: ## Upload release to GitHub releases.
	mkdir -p builds/darwin-${TAG}
	GOOS=darwin GOARCH=$(arch) CGO_ENABLED=0 go build -ldflags "-X github.com/netlify/netlifyctl/commands.Version=${TAG} -X github.com/netlify/netlifyctl/commands.SHA=`git rev-parse HEAD`" -o builds/darwin-${TAG}/netlifyctl
	mkdir -p builds/linux-${TAG}
	GOOS=linux GOARCH=$(arch) CGO_ENABLED=0 go build -ldflags "-X github.com/netlify/netlifyctl/commands.Version=${TAG} -X github.com/netlify/netlifyctl/commands.SHA=`git rev-parse HEAD`" -o builds/linux-${TAG}/netlifyctl
	mkdir -p builds/windows-${TAG}
	GOOS=windows GOARCH=$(arch) CGO_ENABLED=0 go build -ldflags "-X github.com/netlify/netlifyctl/commands.Version=${TAG} -X github.com/netlify/netlifyctl/commands.SHA=`git rev-parse HEAD`" -o builds/windows-${TAG}/netlifyctl.exe
	@rm -rf releases/${TAG}
	mkdir -p releases/${TAG}
	tar -czf releases/${TAG}/netlifyctl-darwin-$(arch)-${TAG}.tar.gz -C builds/darwin-${TAG} netlifyctl
	tar -czf releases/${TAG}/netlifyctl-linux-$(arch)-${TAG}.tar.gz -C builds/linux-${TAG} netlifyctl
	zip -j releases/${TAG}/netlifyctl-windows-$(arch)-${TAG}.zip builds/windows-${TAG}/netlifyctl.exe
	@hub release create -a releases/${TAG}/netlifyctl-darwin-$(arch)-${TAG}.tar.gz -a releases/${TAG}/netlifyctl-linux-$(arch)-${TAG}.tar.gz -a releases/${TAG}/netlifyctl-windows-$(arch)-${TAG}.zip v${TAG}

test: ## Run tests.
	go test -v `go list ./... | grep -v /vendor/`
