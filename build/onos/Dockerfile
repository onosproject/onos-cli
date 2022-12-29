# SPDX-FileCopyrightText: 2019-present Open Networking Foundation <info@opennetworking.org>
#
# SPDX-License-Identifier: Apache-2.0

FROM onosproject/golang-build:v1.3.0 as build

ENV GO111MODULE=on
WORKDIR /build

# copy only the files that are needed for the build, exclude everything else to make better use of the docker cache
COPY ./cmd /build/cmd
COPY ./pkg /build/pkg
COPY ./vendor /build/vendor
COPY ./go.mod /build
COPY ./go.sum /build

# build the executable(s)
RUN go build -mod=vendor -o build/_output/onos ./cmd/onos
RUN	go build -mod=vendor -o build/_output/gnmi_cli ./cmd/gnmi_cli
RUN go install github.com/atomix/cli/cmd/atomix@latest

FROM alpine:3.14
RUN apk add bash bash-completion curl libc6-compat gettext jq

RUN addgroup -S onos && adduser -S -G onos onos --home /home/onos
RUN /bin/bash -c "$(curl -sL https://get-gnmic.kmrd.dev)"
USER onos
WORKDIR /home/onos

COPY --from=build /build/build/_output/onos /usr/local/bin/onos
COPY --from=build /build/build/_output/gnmi_cli /usr/local/bin/gnmi_cli
COPY --from=build /build/pkg/certs/* /etc/ssl/certs/
COPY --from=build /go/bin/atomix /usr/local/bin/atomix

RUN mkdir /home/onos/.onos && \
    cp /etc/profile /home/onos/.bashrc && \
    onos completion bash > /home/onos/.onos/bash_completion.sh && \
    echo "source /home/onos/.onos/bash_completion.sh" >> /home/onos/.bashrc && \
    mkdir /home/onos/.atomix && \
    atomix completion bash > /home/onos/.atomix/bash_completion.sh && \
    echo "source /home/onos/.atomix/bash_completion.sh" >> /home/onos/.bashrc && \
    echo "eval \$(resize)" >> /home/onos/.bashrc

ENTRYPOINT ["/bin/bash"]
