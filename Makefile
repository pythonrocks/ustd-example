default: omni

PWD := $(shell pwd)
USER := $(shell id -u $$USER)
PWD := $(shell pwd)
USER := $(shell id -u $$USER)
PROTOC := docker run --rm -v $(PWD):$(PWD) -u $(USER) -w $(PWD) znly/protoc \
	--proto_path=:.

.phony: interface
interface:
	$(PROTOC) \
		--proto_path=:. \
		--go_out=,plugins=grpc:internal/api \
		service.proto

.phony: descriptor
descriptor:
	$(PROTOC) \
	  	--include_imports \
        --include_source_info \
        --descriptor_set_out=./api_descriptor.pb \
        service.proto

.phony: run
run:
	go run ./cmd/example-api/main.go
