[![Build Status](https://travis-ci.org/netlify/netlifyctl.svg?branch=master)](https://travis-ci.org/netlify/netlifyctl)

# netlifyctl

Legacy command line interface for managing and deploying sites on Netlify without leaving your terminal. Built in Go, on Netlify's [OpenAPI](https://github.com/netlify/open-api) definitions to interact with Netlify's API.

Netlify continues to support this version, but active development has moved to our new [Node-based CLI](https://github.com/netlify/cli). Find the full documentation at [netlify.com/docs/cli](https://www.netlify.com/docs/cli).

## Installation

### Homebrew (Mac)

If you're running the [Homebrew](https://brew.sh/) package manager for Mac, you can install netlifyctl with the following commands:

```bash
brew tap netlify/netlifyctl
brew install netlifyctl
```

When your installation completes, you can run `netlifyctl --help` for a list of available commands, or refer to this doc for more details.

### Scoop (Windows)

[Scoop](https://scoop.sh/) is a package manager for Windows. You can install it with a single PowerShell command, then use it to install other command-line tools. To install netlifyctl using Scoop, run the following commands:

```bash
scoop bucket add netlifyctl https://github.com/netlify/scoop-netlifyctl
scoop install netlifyctl
```

When your installation completes, you can run `netlifyctl --help` for a list of available commands, or refer to this doc for more details.


### Direct Binary Install (Linux, Mac, Windows)

Because netlifyctl is released as an executable binary file, you can download, extract, and run it from any directory you choose, with no dependencies.

You can download the latest version for your operating system directly from the following links:

  - Linux: https://cli.netlify.com/download/latest/linux
  - Mac: https://cli.netlify.com/download/latest/mac
  - Windows: https://cli.netlify.com/download/latest/windows

If you're working on a local machine with administrator permissions, you may want to add netlifyctl to your `PATH`, so that you can access it directly from any project folder. You could do this by extracting it to a directory that is already on your `PATH`. (For example, `/usr/local/bin` is a common choice for Mac and Linux.) Alternatively, you can extract the netlifyctl binary to a different folder, then follow your operating system instructions for adding a path to your `PATH` environment variable.

In cases where you can't or don't want to install globally, like in many Continuous Integration (CI) environments, you can run netlifyctl from the folder of your choice by calling the path to the binary file.

For example, you could run the following command to download and extract the binary file directly into the current directory in a Linux terminal:

```bash
wget -qO- 'https://cli.netlify.com/download/latest/linux' | tar xz
```

Then, to use netlifyctl in that directory, you would use the relative path to the binary: `./netlifyctl`. All `netlifyctl` commands in the rest of this document would follow the same pattern, for example:

```bash
./netlifyctl --help
```

### Installing from Source with `go get`

Use the following commands to install `netlifyctl` from source:

```sh
go get -d github.com/netlify/netlifyctl
cd $GOPATH/src/github.com/netlify/netlifyctl
make deps build
go install
```

## Authentication

Netlifyctl uses an access token to authenticate with Netlify. You can obtain this token via the command line or in the Netlify UI.

### Command-line Login

To authenticate and obtain an access token via the command line, enter the following command:

```bash
netlifyctl login
```

This will open a browser window, asking you to log in with Netlify and grant access to Netlify Cli.

![](site/docs-images/authorize-ui.png)

Once authorized, netlifyctl will store your access token in your home folder, under `.config/netlify`. Netlifyctl will use the token in this location automatically for all future commands.

If you'd like to store your token in a different location, you can remove it from the default location and add it manually to your commands by using the `-A` flag:

```bash
netlifyctl -A "YOUR_ACCESS_TOKEN" deploy
```

If you lose your token, you can repeat this process to generate a new one.

### Obtain a Token in the Netlify UI

You can generate an access token manually in your Netlify account settings under **OAuth applications**, at https://app.netlify.com/applications. 

1. Under **Personal access tokens**, select **New access token**.
2. Enter a description and select **Generate token**.
3. Copy the generated token to your clipboard. Once you navigate from the page, the token cannot be seen again.

You can add the access token to individual commands with the `-A` flag:

```bash
netlifyctl -A "YOUR_ACCESS_TOKEN" deploy
```

Alternatively, you can store the token locally, and netlifyctl will use it automatically. To do this, enter the following line in a file titled `netlify`:

```json
{"access_token": "YOUR_ACCESS_TOKEN"}
```

Store the file in a folder called `.config`, inside your home folder.

### Revoking Access

To revoke access to your account for netlifyctl, go to the **OAuth applications** section of your account settings, at https://app.netlify.com/applications. Find the appropriate token or application, and select **Revoke**.

## Continuous Deployment

With [continuous deployment](/docs/continuous-deployment), Netlify will automatically deploy new versions of your site when you push commits to your connected Git repository. This also enables features like Deploy Previews, branch deploys, and [split testing](/docs/split-testing). (Some of these features must be enabled in the Netlify UI.)

### Automated Setup

For repositories stored on GitHub or GitLab, you can use netlifyctl to connect your repository by running the following command from your local repository:

```bash
netlifyctl init
```

In order to connect your repository for continuous deployment, netlifyctl will need access to create a deploy key and a webhook on the repository. When you run the command above, you'll be prompted to log in to your GitHub account, which will create an account-level access token.

The access token will be stored in your home folder, under `.config/hub`. Your login password will never be stored. You can revoke the access token at any time from your GitHub account settings.

### Manual Setup

For repositories stored on other Git providers, or if you prefer to give more limited, repository-only access, you can connect your repository manually by adding the `--manual` flag. From your local repository, run the following command:

```bash
netlifyctl init --manual
```

The tool will prompt you for your deploy settings, then provide you with two items you will need to add to your repository settings with your Git provider:

- **Deploy/access key:** Netlify uses this key to fetch your repository via ssh for building and deploying.
    ![Sample terminal output reads: 'Give this Netlify SSH public key access to your repository,' and displays a key code.](site/docs-images/deploy-key-cli.png)
Copy the key printed in the command line, then add it as a deploy key in the repository settings on your Git Provider. The deploy key does not require write access. Note that if you have more than one site connected to a repo, you will need a unique key for each one.
- **Webhook:** Your Git provider will send a message to this webhook when you push changes to your repository, triggering a new deploy on Netlify.
    ![Sample terminal output reads: 'Configure the following webhook for your repository,' and displays a URL.](site/docs-images/webhook-cli.png)
Copy the webhook address printed in the command line, then add it as the Payload URL for a new webhook in the repository settings on your Git provider. If available, the **Content type** should be set to `application/json`. When selecting events to trigger the webhook, **Push** events will trigger production and branch deploys on watched branches, and **Pull/Merge request** events will trigger deploy previews.

## Manual Deploy

It's also possible to deploy a site manually, without continuous deployment. This method uploads files directly from your local project directory to your site on Netlify, without running a build step. It also works with directories that are not Git repositories. 

A common use case for this command is when you're using a separate Contiuous Integration (CI) tool, deploying prebuilt files to Netlify at the end of the CI tool tasks.

To deploy manually, run the following command from the base of your project directory:

```bash
netlifyctl deploy
```

Netlifyctl will deploy the site using the configuration settings in a [netlify.toml file](https://www.netlify.com/docs/netlify-toml-reference/) stored at the base of your project directory. If this file doesn't exist, netlifyctl will prompt you for your site settings, then create a new `netlify.toml` file to store them.

After the first deploy, you can run `netlifyctl deploy` again to update your site whenever you make changes. Only new and changed files will be uploaded.

### Draft Deploys

If you'd like to preview a manual deploy without changing it in production, you can use the `--draft` flag:

```bash
netlifyctl deploy --draft
```

This will run a deploy just like your production deploy, but at a unique address. The draft site URL will display in the command line when the deploy is done.

## Debugging

Netlifyctl generates debug logs with all the request and response interations when there is an error running any command. Those logs are stored in a file called `netlifyctl-debug.log` in the directory where you ran the command. These logs include your access token for the API! Please **make sure you don't share them with anyone without masking those first.**

You can force the CLI to generate these logs even when there are no errors with the `-D` flag: `netlifyctl -D deploy`.

## Additional Commands

For a full list of commands and global flags available with netlifyctl, run the following:

```bash
netlifyctl help
```

For more information about a specific command, run `help` with the name of the command.

```bash
netlifyctl help deploy
```

This also works for sub-commands.

```bash
netlifyctl help site update
```

## License

[MIT](LICENSE)
