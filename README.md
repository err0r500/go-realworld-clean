# go-realworld-clean

[![Build Status](https://travis-ci.org/err0r500/go-realworld-clean.svg?branch=master)](https://travis-ci.org/err0r500/go-realworld-clean)
[![BCH compliance](https://bettercodehub.com/edge/badge/err0r500/go-realworld-clean?branch=master)](https://bettercodehub.com/)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/107e5849415b40f4ae9c235afecebf56)](https://www.codacy.com/app/Err0r500/go-realworld-clean?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=err0r500/go-realworld-clean&amp;utm_campaign=Badge_Grade)
[![codecov](https://codecov.io/gh/err0r500/go-realworld-clean/branch/master/graph/badge.svg)](https://codecov.io/gh/err0r500/go-realworld-clean)

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/err0r500/go-realworld-clean/blob/master/LICENSE)

## Clean Architecture :
Layers ( from the most abstract to the most concrete ) :
- domain : abstract data structures
- uc : "use cases", the pure business logic
- implem : implementations of the interfaces used in the business logic (uc layer)
- infra : setup/configuration of the implementation

Golden rule : a layer never imports something from a layer below it

## Install and Hooks

Make sure you installed [dep](https://github.com/golang/dep) and 
[golangci-lint](https://github.com/golangci/golangci-lint). Then install the
dependencies using the `dep` command:

```
$ go get -u github.com/golang/dep/cmd/dep
$ go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
$ dep ensure
```

Then generate the mock to pass tests:

```
$ make mock
```

To enable the pre-commit hook:

```
$ git config core.hooksPath .githooks
```

## Configuration

## Make Targets

The version is either `0.1.0` if no tag has ever been defined or the latest
tag defined. The build number is the SHA1 of the latest commit.

- **make**: Builds and injects version/build in binary
- **make init**: Sets the pre-commit hook in the repository
- **make docker**: Build docker image and tag it with both `latest` and version
- **make latest**: Build docker image and tag it only with `latest`
- **make test**: Executes the test suite
- **make mock**: Generate the necessary mocks
- **make clean**: Removes the built binary if present