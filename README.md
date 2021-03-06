# UDocs

`udocs` is a Go CLI that lets developers easily build, deploy, and publish their app's documentation
guide. Documentation content is written in Markdown, and persisted in a `docs/` directory at the root of your project's Git repository.
From there, `udocs` can render the guide's content to HTML, and optionally serve it locally over HTTP for viewing, or send it to a remote UDocs instance for hosting.

> Note: UDocs currently has an Alpha release status, and is under active development.

---

## Features

- Easy-to-use CLI
- Fast Markdown to HTML rendering 
- Cross-document search (powered by [bleve](https://github.com/blevesearch/bleve))
- Live-reloading of UDocs server when making local document changes 
- MongoDB compatible 

--- 

## Installation

### Requirements

- Git 2.7+
- Go 1.6+
- bash
- Linux or Mac OS X

> Note: all `udocs`-related content (configs, libs, caches, binaries, etc.) will be installed under the directory `${HOME}/.udocs`. If you wish to uninstall `udocs` (and remove all related content from your system), simply `rm -rf ~/.udocs`. 

### Using the boostrap installer script (recommended)

```bash
$ curl https://raw.githubusercontent.com/UltimateSoftware/udocs/master/bin/bootstrap.sh | bash

$ udocs --help
```

### Clone and build from source

```bash
$ go get github.com/ultimatesoftware/udocs

$ cd $GOPATH/src/github.com/ultimatesoftware/udocs

# run the CLI install script that places UDocs in your $PATH, as well as all library/static files
# under the ~/.udocs directory in the local filesystem
$ ./bin/install.sh

# you can adjust UDocs configuration settings in two ways:
$ vim ~/.udocs/udocs.conf
# OR by setting the environment variables in:
$ vim ~/.udocs/.udocs_env && ~/.udocs/.udocs_env

# run UDocs!
$ udocs --help
```

---

## Usage

```
$ udocs --help                                                                                                                                                 [±master ●▴]

Description:
  UDocs is a CLI library for Go that easily renders Markdown documentation guides to HTML, and serves them over HTTP.

Usage:
  udocs [command]

Available Commands:
  build       Build a docs directory
  destroy     Destroy a docs directory from a remote UDocs server
  env         Show UDocs local environment information
  publish     Publish docs to a remote UDocs host
  pull        Pull docs from remote Git repository
  serve       Renders docs directories, and serves them locally over HTTP
  tar         Tar a docs directory
  validate    Validate a docs directory
  version     Show UDocs version

Use "udocs [command] --help" for more information about a command.
```

## Configuration 

`udocs` is configurable via the following environment varibles: 

- `UDOCS_ENTRY_POINT`
- `UDOCS_PORT`
- `UDOCS_BIND_ADDR`
- `UDOCS_EMAIL`
- `UDOCS_ROOT_ROUTE`
- `UDOCS_ROUTES`
- `UDOCS_ORGANIZATION`
- `UDOCS_SEARCH_PLACEHOLDER`
- `UDOCS_MONGO_URL`

Executing `udocs env` will output the state of your current, local environment.

## Vendored Dependencies

- https://github.com/blevesearch/bleve (Apache)
- https://github.com/dimfeld/httptreemux (MIT)
- https://github.com/mholt/archiver (MIT)
- https://github.com/shurcooL/github_flavored_markdown (MIT)
- https://github.com/spf13/cobra (Apache)
- https://github.com/fsnotify/fsnotify/tree/v1.4.1 (BSD-3-Clause)
- https://gopkg.in/mgo.v2 (BSD-2-Clause)
- http://fontawesome.io (http://fontawesome.io/license/)
- http://getbootstrap.com (MIT)
