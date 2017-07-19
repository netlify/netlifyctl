[![Build Status](https://travis-ci.org/netlify/netlifyctl.svg?branch=master)](https://travis-ci.org/netlify/netlifyctl)
# Introduction

netlifyctl is a proof of concept to rewrite Netlify's CLI in Go.
It uses the [OpenAPI](https://github.com/netlify/open-api) definitions
to interact with Netlify's API.


## Installation

### Source

netlifyctl can be installed from source:

```sh
$ go get github.com/netlify/netlifyctl
```

### Binary Distribution

#### GitHub Releases

Prebuilt [binaries are available for osx and linux](https://github.com/netlify/netlifyctl/releases). Other architectures available upon request to support@netlify.com.

#### Homebrew

```sh
brew tap netlify/netlifyctl
brew install netlifyctl
```

## Usage

```sh
$ netlifyctl --help
```


## Contributions and Bug Reports

Contributions are welcome via Pull Request.

Bug Reports are welcome as Issues filed on this repository, but feel free to chat with support@netlify.com about issues too!


## License

[MIT](LICENSE)
