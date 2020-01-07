
Kaspanet Faucet
====
Warning: This is pre-alpha software. There's no guarantee anything works.
====

[![ISC License](http://img.shields.io/badge/license-ISC-blue.svg)](http://copyfree.org)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/kaspanet/faucet)

Kaspanet Faucet is a faucet implementation for Kaspa written in Go (golang).

This project is currently under active development and is in a pre-Alpha state. 
Some things still don't work and APIs are far from finalized. The code is provided for reference only.

## Requirements

Latest version of [Go](http://golang.org) (currently 1.13).

## Installation

#### Build from Source

- Install Go according to the installation instructions here:
  http://golang.org/doc/install

- Ensure Go was installed properly and is a supported version:

```bash
$ go version
$ go env GOROOT GOPATH
```

NOTE: The `GOROOT` and `GOPATH` above must not be the same path. It is
recommended that `GOPATH` is set to a directory in your home directory such as
`~/dev/go` to avoid write permission issues. It is also recommended to add
`$GOPATH/bin` to your `PATH` at this point.

- Run the following commands to obtain and install faucet including all dependencies:

```bash
$ git clone https://github.com/kaspanet/faucet $GOPATH/src/github.com/kaspanet/faucet
$ cd $GOPATH/src/github.com/kaspanet/faucet
$ go install ./...
```

- faucet should now be installed in `$GOPATH/bin`. If you did
  not already add the bin directory to your system path during Go installation,
  you are encouraged to do so now.


## Getting Started

Faucet expects to have access to the following systems:
- A [kasparovd](https://github.com/kaspanet/kasparov) instance
- A MySQL database

### Linux/BSD/POSIX/Source

```bash
$ ./faucet --dbuser=user --dbpass=pass --dbaddress=localhost:3306 --dbname=faucet --migrate --testnet
$ ./faucet --dbuser=user --dbpass=pass --dbaddress=localhost:3306 --dbname=faucet --fee-rate=5 --private-key=00000000000000000000000000000000000000000000 --kasparovd-url=http://localhost:8080 --testnet
```

## Discord
Join our discord server using the following link: https://discord.gg/WmGhhzk

## Issue Tracker

The [integrated github issue tracker](https://github.com/kaspanet/faucet/issues)
is used for this project.

## License

Faucet is licensed under the [copyfree](http://copyfree.org) ISC License.

