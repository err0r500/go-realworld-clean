.PHONY: all clean test docker latest init

export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64

BINARY=go-realworld-clean
VERSION=$(shell git describe --abbrev=0 --tags 2> /dev/null || echo "0.1.0")
BUILD=$(shell git rev-parse HEAD 2> /dev/null || echo "undefined")
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.Build=$(BUILD)"

all:
	GO111MODULE=on go build -o $(BINARY) $(LDFLAGS)

init:
	git config core.hooksPath .githooks

docker:
	docker build \
		-t $(BINARY):latest \
		-t $(BINARY):$(VERSION) \
		--build-arg build=$(BUILD) --build-arg version=$(VERSION) \
		-f Dockerfile --no-cache .

latest:
	docker build \
		-t $(BINARY):latest \
		--build-arg build=$(BUILD) --build-arg version=$(VERSION) \
		-f Dockerfile --no-cache .

test:
	GO111MODULE=on go test ./...

clean:
	if [ -f $(BINARY) ] ; then rm $(BINARY) ; fi

mock:
	$(GOPATH)/bin/mockgen -source=./uc/INTERACTOR.go -destination=./implem/uc.mock/interactor.go -package=mock && \
    $(GOPATH)/bin/mockgen -source=./uc/HANDLER.go -destination=./implem/uc.mock/handler.go -package=mock
