ARG ONOS_BUILD_VERSION=stable

FROM onosproject/golang-build:$ONOS_BUILD_VERSION AS build
ENV GO111MODULE=on
RUN go get github.com/openconfig/gnmi/cmd/gnmi_cli && \
    go get github.com/atomix/atomix-cli/cmd/atomix
COPY . /go/src/github.com/onosproject/onos-cli
RUN cd /go/src/github.com/onosproject/onos-cli && GOFLAGS=-mod=vendor make build

FROM python:alpine3.10 as pybuilder

RUN apk update
RUN apk add git python3-dev build-base
RUN pip3 install --install-option="--prefix=/install" grpcio protobuf
RUN pip3 install grpcio-tools

ENV PYTHONPATH="/install/lib/python3.7/site-packages"

# Change back to official one after they fix this
# RUN git clone https://github.com/openconfig/gnmi
RUN git clone https://github.com/bocon13/gnmi
RUN git clone https://github.com/bocon13/Yi-s-gNMI-tool.git
RUN git clone https://github.com/google/protobuf

ENV proto_imports=.:/protobuf/src

RUN cd gnmi/proto && \
    python -m grpc_tools.protoc -I=$proto_imports --python_out=.                     gnmi_ext/gnmi_ext.proto && \
    python -m grpc_tools.protoc -I=$proto_imports --python_out=. --grpc_python_out=. gnmi/gnmi.proto && \
    python -m grpc_tools.protoc -I=$proto_imports --python_out=.                     target/target.proto

RUN mv /gnmi/proto/gnmi $PYTHONPATH/ && \
    touch $PYTHONPATH/gnmi/__init__.py
RUN mv /gnmi/proto/gnmi_ext $PYTHONPATH/ && \
    touch $PYTHONPATH/gnmi_ext/__init_.py
RUN mv /gnmi/proto/target $PYTHONPATH/ && \
    touch $PYTHONPATH/target

# Python is only added here so that we can use the Stratum gnmi-cli.py
# When Stratum implements TLS we will be able to use the go gnmi_cli
# and this should be removed
FROM python:3.7-alpine3.10
RUN apk add --no-cache bash bash-completion curl libc6-compat libstdc++

RUN addgroup -S onos && adduser -S -G onos onos
USER onos
WORKDIR /home/onos

COPY --from=build /go/src/github.com/onosproject/onos-cli/build/_output/onos /usr/local/bin/onos
COPY --from=build /go/bin/atomix /usr/local/bin/atomix
COPY --from=build /go/bin/gnmi_cli /usr/local/bin/gnmi_cli
COPY --from=build /go/src/github.com/onosproject/onos-cli/pkg/certs/* /etc/ssl/certs/

COPY --from=pybuilder /install /usr/local
COPY --from=pybuilder Yi-s-gNMI-tool/gnmi-cli.py /usr/local/bin/gnmi-cli-stratum

RUN mkdir /home/onos/.onos && \
    cp /etc/profile /home/onos/.bashrc && \
    onos completion bash > /home/onos/.onos/bash_completion.sh && \
    echo "source /home/onos/.onos/bash_completion.sh" >> /home/onos/.bashrc && \
    mkdir /home/onos/.atomix && \
    atomix completion bash > /home/onos/.atomix/bash_completion.sh && \
    echo "source /home/onos/.atomix/bash_completion.sh" >> /home/onos/.bashrc && \
    echo "eval \$(resize)" >> /home/onos/.bashrc

ENTRYPOINT ["/bin/bash"]
