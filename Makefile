#
# Copyright (C) 2022 IOTech Ltd
#
# SPDX-License-Identifier: Apache-2.0
#


.PHONY: build clean docker run

# change the following boolean flag to include or exclude the delayed start libs for builds for most of core services except support services
INCLUDE_DELAYED_START_BUILD_CORE:="false"
# change the following boolean flag to include or exclude the delayed start libs for builds for support services exculsively
INCLUDE_DELAYED_START_BUILD_SUPPORT:="true"

GO=CGO_ENABLED=0 GO111MODULE=on go

# see https://shibumi.dev/posts/hardening-executables
CGO_CPPFLAGS="-D_FORTIFY_SOURCE=2"
CGO_CFLAGS="-O2 -pipe -fno-plt"
CGO_CXXFLAGS="-O2 -pipe -fno-plt"
CGO_LDFLAGS="-Wl,-O1,–sort-common,–as-needed,-z,relro,-z,now"
GOCGO=CGO_ENABLED=1 GO111MODULE=on go

CGOFLAGS=-ldflags "-linkmode=external -X github.com/edgexfoundry/edgex-go.Version=$(VERSION)" -trimpath -mod=readonly -buildmode=pie

VERSION=$(shell cat ./VERSION 2>/dev/null || echo 0.0.0)
DOCKER_TAG=$(VERSION)-dev

GIT_SHA=$(shell git rev-parse HEAD)

ARCH=$(shell uname -m)

# DO NOT change the following flag, as it is automatically set based on the boolean switch INCLUDE_DELAYED_START_BUILD_CORE
NON_DELAYED_START_GO_BUILD_TAG_FOR_CORE:=non_delayedstart
ifeq ($(INCLUDE_DELAYED_START_BUILD_CORE),"true")
	NON_DELAYED_START_GO_BUILD_TAG_FOR_CORE:=
endif
NON_DELAYED_START_GO_BUILD_TAG_FOR_SUPPORT:=
ifeq ($(INCLUDE_DELAYED_START_BUILD_SUPPORT),"false")
	NON_DELAYED_START_GO_BUILD_TAG_FOR_SUPPORT:=non_delayedstart
endif

NO_MESSAGEBUS_GO_BUILD_TAG:=no_messagebus
NO_DTO_VALIDATOR_GO_BUILD_TAG:=no_dto_validator
NO_ZMQ:=no_zmq

tidy:
	go mod tidy -compat=1.17

export EDGEX_SECURITY_SECRET_STORE=false

cmd/core-data/core-data:
	$(GOCGO) build -tags "$(NON_DELAYED_START_GO_BUILD_TAG_FOR_CORE)" $(CGOFLAGS) -o $@ ./cmd/core-data

cmd/core-keeper/core-keeper:
	$(GO) build -tags "$(NO_DTO_VALIDATOR_GO_BUILD_TAG) $(NON_DELAYED_START_GO_BUILD_TAG_FOR_CORE) $(NO_ZMQ)" $(CGOFLAGS) -o $@ ./cmd/core-keeper

build: cmd/core-keeper/core-keeper

run: cmd/core-keeper/core-keeper
	cd ./cmd/core-keeper; \
	./core-keeper

run_core_data: cmd/core-data/core-data
	cd ./cmd/core-data; \
	./core-data --registry -cp=keeper.http://localhost:59883

clean:
	rm -f cmd/core-keeper/core-keeper

docker: docker_core_keeper

docker_core_keeper:
	docker build \
	    --build-arg http_proxy \
	    --build-arg https_proxy \
		-f cmd/core-keeper/Dockerfile \
		--label "git_sha=$(GIT_SHA)" \
		-t edgexfoundry/core-keeper:$(GIT_SHA) \
		-t edgexfoundry/core-keeper:$(DOCKER_TAG) \
		.

vendor:
	$(GO) mod vendor
