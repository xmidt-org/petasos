# petasos
(pronounced "pet-uh-sos")

[![Build Status](https://travis-ci.com/xmidt-org/petasos.svg?branch=master)](https://travis-ci.com/xmidt-org/petasos)
[![codecov.io](http://codecov.io/github/xmidt-org/petasos/coverage.svg?branch=master)](http://codecov.io/github/xmidt-org/petasos?branch=master)
[![Code Climate](https://codeclimate.com/github/xmidt-org/petasos/badges/gpa.svg)](https://codeclimate.com/github/xmidt-org/petasos)
[![Issue Count](https://codeclimate.com/github/xmidt-org/petasos/badges/issue_count.svg)](https://codeclimate.com/github/xmidt-org/petasos)
[![Go Report Card](https://goreportcard.com/badge/github.com/xmidt-org/petasos)](https://goreportcard.com/report/github.com/xmidt-org/petasos)
[![Apache V2 License](http://img.shields.io/badge/license-Apache%20V2-blue.svg)](https://github.com/xmidt-org/petasos/blob/master/LICENSE)
[![GitHub release](https://img.shields.io/github/release/xmidt-org/petasos.svg)](CHANGELOG.md)

## Summary
Petasos is the HTTP redirector component. Petasos will redirect http requests, to a
talaria depending on the the device id and talaria service discovery configuration.

## Details
There is only one endpoint with petasos `/api/v2/device`. Petasos will return a
http 307 redirect to the corresponding talaria.


## Build

### Source

In order to build from the source, you need a working Go environment with
version 1.11 or greater. Find more information on the [Go website](https://golang.org/doc/install).

You can directly use `go get` to put the petasos binary into your `GOPATH`:
```bash
GO111MODULE=on go get github.com/xmidt-org/petasos
```

You can also clone the repository yourself and build using make:

```bash
mkdir -p $GOPATH/src/github.com/xmidt-org
cd $GOPATH/src/github.com/xmidt-org
git clone git@github.com:xmidt-org/petasos.git
cd petasos
make build
```

### Makefile

The Makefile has the following options you may find helpful:
* `make build`: builds the petasos binary
* `make rpm`: builds an rpm containing petasos
* `make docker`: builds a docker image for petasos, making sure to get all
   dependencies
* `make local-docker`: builds a docker image for petasos with the assumption
   that the dependencies can be found already
* `make test`: runs unit tests with coverage for petasos
* `make clean`: deletes previously-built binaries and object files

### Docker

The docker image can be built either with the Makefile or by running a docker
command.  Either option requires first getting the source code.

See [Makefile](#Makefile) on specifics of how to build the image that way.

For running a command, either you can run `docker build` after getting all
dependencies, or make the command fetch the dependencies.  If you don't want to
get the dependencies, run the following command:
```bash
docker build -t petasos:local -f deploy/Dockerfile .
```
If you want to get the dependencies then build, run the following commands:
```bash
GO111MODULE=on go mod vendor
docker build -t petasos:local -f deploy/Dockerfile.local .
```

For either command, if you want the tag to be a version instead of `local`,
then replace `local` in the `docker build` command.

### Kubernetes

WIP. TODO: add info

## Deploy

For deploying a XMiDT cluster refer to [getting started](https://xmidt.io/docs/operating/getting_started/).

For running locally, ensure you have the binary [built](#Source).  If it's in
your `GOPATH`, run:
```
petasos
```
If the binary is in your current folder, run:
```
./petasos
```

## Contributing

Refer to [CONTRIBUTING.md](CONTRIBUTING.md).
