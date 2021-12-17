FROM golang:1.14 as builder
## docker build -t xxx .
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

RUN curl -s "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh"  | bash \
   && chmod +x kustomize \
   && mv kustomize /usr/local/bin/


RUN curl -L https://github.com/a8m/envsubst/releases/download/v1.2.0/envsubst-`uname -s`-`uname -m` -o envsubst \
    && chmod +x envsubst \
    && sudo mv envsubst /usr/local/bin