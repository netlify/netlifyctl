[![Build Status](https://travis-ci.org/netlify/netlifyctl.svg?branch=master)](https://travis-ci.org/netlify/netlifyctl)

# Introduction

Netlify's CLI to manage and deploy sites on Netlify without leaving your terminal.

It uses the [OpenAPI](https://github.com/netlify/open-api) definitions to interact with Netlify's API.

## Installation

### Binary Distribution

#### GitHub Releases

Prebuilt [binaries are available for Windows, OS X and Linux](https://github.com/netlify/netlifyctl/releases).

#### Homebrew

```sh
brew tap netlify/netlifyctl
brew install netlifyctl
```

### Source

Use the following commands to install `netlifyctl` from source:

```sh
go get -d github.com/netlify/netlifyctl
cd $GOPATH/src/github.com/netlify/netlifyctl
make deps build
go install
```

## Usage

### Overview

```sh
netlifyctl --help
```

or to get details on a subcommand:

```sh
netlifyctl site update --help
```

## Quickstart

1.  Use `netlifyctl login` to create an API token for your personal use. This command requires you to have access to a browser. Your access token will be stored in %HOME%/.config/netlify when you run the command directly.

2.  Use `netlifyctl sites` to display the list of sites you have access to.

3.  Use `netlifyctl deploy` to deploy changes on a site. This command must run from the root directory where you have your site's source code. The interactive guides will take your site ID and deploy path and incorporate them into that config file.

Thereafter, you can run unattended and headless using the flag `-y` to auto confirm the current settings: `netlifyctl -y deploy`.

## Debugging

Netlifyctl generates debug logs with all the request and response interations when there is an error running any command. Those logs are stored in a file called `netlifyctl-debug.log` in the directory where you ran the command. These logs include your access token for the API! Please **make sure you don't share them with anyone without masking those first.**

You can force the CLI to generate these logs even when there are no errors with the `-D` flag: `netlifyctl -D deploy`.

### Logging in & creating authentication tokens via a browser instead of the command line

You can get an access token from https://app.netlify.com/applications. Once you've created it, you can store it in your computer. The default location where netlifyctl tries to find this token is within your home directory, inside a file called `netlify` within a directory called `.config`. The path in a Unix system looks like `~/.config/netlify`. This file uses JSON formatting, you can see an example below:

```json
{
  "access_token": "my secret access token"
}
```

You can also set this token with the flag `-A` in each command call if you don't want to store it in a file: `netlifyctl -A "my secret access token" deploy`.

### Redirects/headers in the `netlify.toml` file does not work after deploy

If the redirects/headers in the `netlify.toml` file don't work, this is because the `netlifyctl deploy` command doesn't automatically include the netlify.toml file in the deploy. You will need to manually copy the file to your publish directory before running the netlifyctl command.

## Contributions and Bug Reports

Contributions are welcome via Pull Request.

Bug Reports are welcome as Issues filed on this repository, but feel free to chat with [support@netlify.com](mailto:support@netlify.com) about issues as well - often we'll have faster advice to help you succeed.

## License

[MIT](LICENSE)
