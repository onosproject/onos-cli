ARG ONOS_BUILD_VERSION=stable

FROM onosproject/golang-build:$ONOS_BUILD_VERSION AS build
ENV GO111MODULE=on
COPY . /go/src/github.com/onosproject/onos-cli
RUN cd /go/src/github.com/onosproject/onos-cli && GOFLAGS=-mod=vendor make build

FROM alpine:3.9
RUN apk add libc6-compat

USER nobody

COPY --from=build /go/src/github.com/onosproject/onos-cli/build/_output/onos /usr/local/bin/onos

ENTRYPOINT ["onos"]
