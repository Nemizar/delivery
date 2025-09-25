ARG BUILD_ARG_GO_VERSION=1.24
ARG BUILD_ARG_ALPINE_VERSION=3.20

FROM golang:${BUILD_ARG_GO_VERSION}-alpine${BUILD_ARG_ALPINE_VERSION} AS builder

ENV OAPI_CODEGEN_VERSION=2.5.0
ENV PROTOC_GO_VERSION=1.31.0
ENV PROTOC_GRPC_VERSION=1.3.0
ENV PROTOC_VERSION=31.1

RUN apk add --update --no-cache curl unzip

RUN go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v${OAPI_CODEGEN_VERSION} && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@v${PROTOC_GO_VERSION} && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v${PROTOC_GRPC_VERSION}

RUN curl -sSL "https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-linux-x86_64.zip" \
    -o /tmp/protoc.zip && \
    mkdir -p /usr/local/include && \
    unzip /tmp/protoc.zip -d /usr/local && \
    rm /tmp/protoc.zip

ENV PATH="/usr/local/bin:/usr/local/include:$PATH"

WORKDIR /src
