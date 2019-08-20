ARG ONOS_BUILD_VERSION=stable

FROM onosproject/golang-build:$ONOS_BUILD_VERSION AS build
ENV GO111MODULE=on
RUN go get -u github.com/openconfig/gnmi/cmd/gnmi_cli
COPY . /go/src/github.com/onosproject/onos-cli
RUN cd /go/src/github.com/onosproject/onos-cli && GOFLAGS=-mod=vendor make build

FROM alpine:3.9
RUN apk add bash bash-completion libc6-compat

RUN addgroup -S onos && adduser -S -G onos onos
USER onos
WORKDIR /home/onos

COPY --from=build /go/src/github.com/onosproject/onos-cli/build/_output/onos /usr/local/bin/onos
COPY --from=build /go/bin/gnmi_cli /usr/local/bin/gnmi_cli
COPY --from=build /go/src/github.com/onosproject/onos-cli/pkg/certs/* /etc/ssl/certs/

RUN mkdir /home/onos/.onos && \
    cp /etc/profile /home/onos/.bashrc && \
    onos completion bash > /home/onos/.onos/bash_completion.sh && \
    echo "source /home/onos/.onos/bash_completion.sh" >> /home/onos/.bashrc

ENTRYPOINT ["/bin/bash"]
