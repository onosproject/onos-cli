FROM onosproject/golang-build:v1.0 AS build
ENV GO111MODULE=on
RUN go get github.com/atomix/cli/cmd/atomix
COPY . /go/src/github.com/onosproject/onos-cli
RUN cd /go/src/github.com/onosproject/onos-cli && GOFLAGS=-mod=vendor make build

FROM alpine:3.14
RUN apk add bash bash-completion curl libc6-compat gettext jq

RUN addgroup -S onos && adduser -S -G onos onos --home /home/onos
USER onos
WORKDIR /home/onos

COPY --from=build /go/src/github.com/onosproject/onos-cli/build/_output/onos /usr/local/bin/onos
COPY --from=build /go/src/github.com/onosproject/onos-cli/build/_output/gnmi_cli /usr/local/bin/gnmi_cli
COPY --from=build /go/bin/atomix /usr/local/bin/atomix
COPY --from=build /go/src/github.com/onosproject/onos-cli/pkg/certs/* /etc/ssl/certs/

RUN mkdir /home/onos/.onos && \
    cp /etc/profile /home/onos/.bashrc && \
    onos completion bash > /home/onos/.onos/bash_completion.sh && \
    echo "source /home/onos/.onos/bash_completion.sh" >> /home/onos/.bashrc && \
    mkdir /home/onos/.atomix && \
    atomix completion bash > /home/onos/.atomix/bash_completion.sh && \
    echo "source /home/onos/.atomix/bash_completion.sh" >> /home/onos/.bashrc && \
    echo "eval \$(resize)" >> /home/onos/.bashrc

ENTRYPOINT ["/bin/bash"]
