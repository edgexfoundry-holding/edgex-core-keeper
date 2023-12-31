#  ----------------------------------------------------------------------------------
#  Copyright 2021 Intel Corp.
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
#
#  SPDX-License-Identifier: Apache-2.0
#  ----------------------------------------------------------------------------------

ARG BUILDER_BASE=golang:1.17-alpine3.15
FROM ${BUILDER_BASE} AS builder

WORKDIR /edgex-go

RUN sed -e 's/dl-cdn[.]alpinelinux.org/nl.alpinelinux.org/g' -i~ /etc/apk/repositories

RUN apk add --update --no-cache make git

COPY go.mod vendor* ./
RUN [ ! -d "vendor" ] && go mod download all || echo "skipping..."

COPY . .
RUN make cmd/security-bootstrapper/security-bootstrapper

FROM alpine:3.14

LABEL license='SPDX-License-Identifier: Apache-2.0' \
      copyright='Copyright (c) 2021 Intel Corporation'

RUN apk add --update --no-cache dumb-init su-exec

ENV SECURITY_INIT_DIR /edgex-init
ARG BOOTSTRAP_REDIS_DIR=${SECURITY_INIT_DIR}/bootstrap-redis

RUN mkdir -p ${SECURITY_INIT_DIR} \
    && mkdir -p ${BOOTSTRAP_REDIS_DIR}

WORKDIR ${SECURITY_INIT_DIR}

# copy all entrypoint scripts into shared folder
COPY --from=builder /edgex-go/cmd/security-bootstrapper/entrypoint-scripts/ ${SECURITY_INIT_DIR}/
RUN chmod +x ${SECURITY_INIT_DIR}/*.sh

COPY --from=builder /edgex-go/Attribution.txt /
COPY --from=builder /edgex-go/cmd/security-bootstrapper/security-bootstrapper .
COPY --from=builder /edgex-go/cmd/security-bootstrapper/res/configuration.toml ./res/

# needed for bootstrapping Redis db
COPY --from=builder /edgex-go/cmd/security-bootstrapper/res-bootstrap-redis/configuration.toml ${BOOTSTRAP_REDIS_DIR}/res/

# copy Consul ACL related configs
COPY --from=builder /edgex-go/cmd/security-bootstrapper/consul-acl/ ${SECURITY_INIT_DIR}/consul-bootstrapper/

# Expose the file directory as a volume since there's long-running state
VOLUME ${SECURITY_INIT_DIR}

# setup entry point script
COPY --from=builder /edgex-go/cmd/security-bootstrapper/entrypoint.sh /
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]

# gate is one subcommand for security-bootstrapper to do security bootstrapping
CMD ["gate"]
