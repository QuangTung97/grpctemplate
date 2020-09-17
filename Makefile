.PHONY: all

PROTO_DIR := proto
RPC_DIR := rpc

CURRENT_DIR := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

GOOGLE_API_PATH := ${GOPATH}/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.14.8/third_party/googleapis/

define generate
	mkdir -p ${RPC_DIR}/$(1) && \
		cd ${PROTO_DIR}/$(1) && \
		protoc -I.:${CURRENT_DIR}/${PROTO_DIR}:${GOOGLE_API_PATH} \
			--go_out=paths=source_relative:${CURRENT_DIR}/${RPC_DIR}/$(1) \
			--go-grpc_out=paths=source_relative:${CURRENT_DIR}/${RPC_DIR}/$(1) \
			--grpc-gateway_out=logtostderr=true,paths=source_relative:${CURRENT_DIR}/${RPC_DIR}/$(1) \
			$(2)
endef

all:
	go build -o server cmd/server/main.go

gen:
	rm -rf rpc
	$(call generate,lib/v1/,enum_campaign_type.proto)
	$(call generate,backend/v1/,backend.proto)

install-tools:
	go install github.com/fzipp/gocyclo
	go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
	go install github.com/kisielk/errcheck
	go install golang.org/x/lint/golint
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go install google.golang.org/protobuf/cmd/protoc-gen-go

lint:
	go fmt
	golint ./...
	go vet ./...
	errcheck ./...
	gocyclo -over 12 .
