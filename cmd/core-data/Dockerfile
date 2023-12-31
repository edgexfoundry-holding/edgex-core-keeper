#  ----------------------------------------------------------------------------------
#  Copyright 2018 Dell Technologies, Inc.
#  Copyright 2018 Cavium
#
#  Licensed under the Apache License, Version 2.0 (the "License");
#  you may not use this file except in compliance with the License.
#  You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
#   Unless required by applicable law or agreed to in writing, software
#   distributed under the License is distributed on an "AS IS" BASIS,
#   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#   See the License for the specific language governing permissions and
#   limitations under the License.
# 
#  ----------------------------------------------------------------------------------

# Docker image for Golang Core Data micro service
ARG BUILDER_BASE=golang:1.17-alpine3.15
FROM ${BUILDER_BASE} AS builder

WORKDIR /edgex-go

# The main mirrors are giving us timeout issues on builds periodically.
# So we can try these.

RUN sed -e 's/dl-cdn[.]alpinelinux.org/nl.alpinelinux.org/g' -i~ /etc/apk/repositories

RUN apk add --update --no-cache zeromq-dev libsodium-dev pkgconfig build-base git

COPY go.mod vendor* ./
RUN [ ! -d "vendor" ] && go mod download all || echo "skipping..."

COPY . .
RUN make cmd/core-data/core-data

#Next image - Copy built Go binary into new workspace
FROM alpine:3.14

LABEL license='SPDX-License-Identifier: Apache-2.0' \
      copyright='Copyright (c) 2018: Dell, Cavium'

ENV APP_PORT=59880
#expose core data port
EXPOSE $APP_PORT

# The main mirrors are giving us timeout issues on builds periodically.
# So we can try these.
RUN sed -e 's/dl-cdn[.]alpinelinux.org/nl.alpinelinux.org/g' -i~ /etc/apk/repositories

RUN apk add --update --no-cache zeromq dumb-init
COPY --from=builder /edgex-go/Attribution.txt /
COPY --from=builder /edgex-go/cmd/core-data/core-data /
COPY --from=builder /edgex-go/cmd/core-data/res/configuration.toml /res/configuration.toml

ENTRYPOINT ["/core-data"]
CMD ["-cp=consul.http://edgex-core-consul:8500", "--registry", "--confdir=/res"]
