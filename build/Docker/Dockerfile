FROM golang:1.14 as builder

WORKDIR /go/grpc-gateway

RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn \
    && echo export PATH="$PATH:$(go env GOPATH)/bin" >> ~/.bashrc

RUN apt-get update\
    && apt install -y protobuf-compiler \
    && go mod init gateway\
    && go get -v -u github.com/grpc-ecosystem/grpc-gateway/v2@v2.6.0 \
    && go install \
           github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
           github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
           google.golang.org/protobuf/cmd/protoc-gen-go \
           google.golang.org/grpc/cmd/protoc-gen-go-grpc



(x-12/2)*1+(x-10)