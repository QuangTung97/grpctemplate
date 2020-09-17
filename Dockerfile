FROM golang:1.15.2-buster as builder
WORKDIR /go/grpctemplate/

COPY install-protoc.sh ./
RUN ./install-protoc.sh

COPY go.mod go.sum Makefile ./
RUN go mod download && make install-tools
COPY . .
RUN make gen && make all

FROM debian:buster
WORKDIR /go/grpctemplate/
COPY --from=builder /go/grpctemplate/server .
