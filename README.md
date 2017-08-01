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

### Overview

```sh
$ netlifyctl --help
```

### Quickstart

1. create an API token for your deploys here: https://app.netlify.com/applications (and save it)
2. find your site ID - most easily done at the bottom of your site's settings page
3. do one deploy manually, which will create a new config file (or add to an existing one): `netlify.toml`. You'll want to run this command on a system with both your code checked out AND a browser to interact with (for login purposes).

```sh
netlifyctl -A YOURAPITOKEN deploy
```

...and the interactive guides will take your site ID and deploy path and incorporate them into that config file.

Thereafter, you can run unattended and headless (though you should check the return code from running, in case there is some error):

```sh
netlifyctl -y -A YOURAPITOKEN deploy
```



## Contributions and Bug Reports

Contributions are welcome via Pull Request.

Bug Reports are welcome as Issues filed on this repository, but feel free to chat with support@netlify.com about issues too!


## License

[MIT](LICENSE)
