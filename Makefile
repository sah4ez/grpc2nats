GIT_SHA?=$(shell git log -n 1 | head -n 1 | awk '{print $$2}')
BUILD_STAMP=`date +%FT%T%z`
VERSION=v0.0.1
LDFLAGS=-ldflags "-extldflags "-static" -X main.GitSHA=$(GIT_SHA) -X main.BuildStamp=$(BUILD_STAMP) -X main.Version=$(VERSION)"
BIN_OUT=./bin


build:
	@echo "$(APP_BINARY_NAME) ..."
	CGO_ENABLED=0 go build $(LDFLAGS) -v -o $(BIN_OUT) -a -tags="$(BUILD_TAGS)" ./cmd/sender/.
	CGO_ENABLED=0 go build $(LDFLAGS) -v -o $(BIN_OUT) -a -tags="$(BUILD_TAGS)" ./cmd/receiver/.
	@chown $(UID):$(GID) -R ./bin go.*

gen:
	go generate ./...
