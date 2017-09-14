jobserveclient
===
jobserveclient accepts TaskRequests from the Scheduler on 0.0.0.0:3000 and responds with TaskResults


Getting Started
---

### Prerequisities

Building and running jobserveclient requires that you have Go, git, and make.

#### Go

jobserveclient expects Go v1.8 or later to be installed

#### Source

Create the following directory and cd into it:

`$GOPATH/src/`

then, clone/copy the jobserveclient repo.

After cloning the repo, run `make setup` to install needed development tools.

### Validate the install

Build the binary and start the jobserveclient with:

`make run`

From there you can run the cmd/scheduler binaries to check for functional correctness. The binary at `bin/jsclient` can be run on darwin systems and is listening on port `3000`.

Usage
---

Common commands are defined as make targets

* `make` - builds all executables
* `make setup` - installs utilities necessary for development.
* `make run` - builds and runs the `jsclient` binary locally
* `make clean` - remove build artifacts
* `make unittest` - run unit tests
* `make functionaltest` - run functional tests
* `make lint` - checks code for formatting
