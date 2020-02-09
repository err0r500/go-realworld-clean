# ![RealWorld Example App](https://github.com/gothinkster/realworld-starter-kit/raw/master/logo.png)

> ### Golang clean-architecture codebase containing real world examples (CRUD, auth, advanced patterns, etc) that adheres to the [RealWorld](https://github.com/gothinkster/realworld) spec and API.


### [Demo](https://github.com/gothinkster/realworld)&nbsp;&nbsp;&nbsp;&nbsp;[RealWorld](https://github.com/gothinkster/realworld)


This codebase was created to demonstrate a fully fledged fullstack application built with  go including CRUD operations, authentication, routing, pagination, and more.

We've gone to great lengths to adhere to the go community styleguides & best practices.

For more information on how to this works with other frontends/backends, head over to the [RealWorld](https://github.com/gothinkster/realworld) repo.


[![Build Status](https://travis-ci.org/err0r500/go-realworld-clean.svg?branch=master)](https://travis-ci.org/err0r500/go-realworld-clean)
[![BCH compliance](https://bettercodehub.com/edge/badge/err0r500/go-realworld-clean?branch=master)](https://bettercodehub.com/)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/107e5849415b40f4ae9c235afecebf56)](https://www.codacy.com/app/Err0r500/go-realworld-clean?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=err0r500/go-realworld-clean&amp;utm_campaign=Badge_Grade)
[![codecov](https://codecov.io/gh/err0r500/go-realworld-clean/branch/master/graph/badge.svg)](https://codecov.io/gh/err0r500/go-realworld-clean)


# How it works

## Clean Architecture :
Layers ( from the most abstract to the most concrete ) :
- domain : abstract data structures
- uc : "use cases", the pure business logic
- implem : implementations of the interfaces used in the business logic (uc layer)
- infra : setup/configuration of the implementation

#### Golden rules : 
- a layer never imports something from a layer below it
- 3rd-party libraries are forbidden in the 2 topmost layers

#### Benefits :
- flexibility
- testability

# Getting started
### Build the app 
```
make
```
### Run the app
```bash
./go-realworld-clean
```

### Run the integration tests
Start the server with an existing user
```
./go-realworld-clean --populate=true
```

In another terminal, run the tests against the API
```
newman run api/Conduit.postman_collection.json \
  -e api/Conduit.postman_integration_test_environment.json \
  --global-var "EMAIL=joe@what.com" \
  --global-var "PASSWORD=password"
```
# Additional
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