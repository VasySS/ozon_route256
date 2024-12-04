ifeq ($(OS),Windows_NT)
  EXT := .exe
else
  EXT :=
endif

APP_NAME := pvz
BUILD_DIR := ./build
SRC_DIR := ./cmd/interactive
DOCKER_YML := ./docker-compose.yml

LOCAL_BIN:=$(CURDIR)/bin
OUT_PATH:=$(CURDIR)/pkg
PROTOS_PATH=./api/pvz/v1/*.proto

PG_URL := postgres://postgres:postgres@localhost:5432/route256?sslmode=disable

all: deps build run

build:
	@echo "Building the application..."
	go build -o $(BUILD_DIR)/$(APP_NAME) $(SRC_DIR)/main.go

deps: bin-deps
	go mod download

deps-update:
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy

run:
	@echo "Running the application..."
	$(BUILD_DIR)/$(APP_NAME)

run-prometheus:
	prometheus --config.file config/prometheus.yaml

lint:
	@echo "Running golangci-lint..."
	golangci-lint run

clean:
	@echo "Cleaning build directory..."
	rm -rf $(BUILD_DIR)

depgraph-install:
	go install github.com/kisielk/godepgraph@latest

depgraph-build:
	godepgraph ./cmd/cobra | dot -Tpng -o godepgraph.png

depgraph:
	make depgraph-install
	make depgraph-build

test:
	@echo "Running tests..."
	go test ./...

test-coverage:
	@echo "Generating coverage report..."
	go test -coverprofile=$(BUILD_DIR)/coverage.out ./...
	go tool cover -html=$(BUILD_DIR)/coverage.out -o $(BUILD_DIR)/coverage.html
	@echo "Coverage report generated at $(BUILD_DIR)/coverage.html"

# ---------------------------
# Запуск базы данных в Docker
# ---------------------------

compose-up:
	docker-compose -f $(DOCKER_YML) up -d

compose-down:
	docker-compose -f $(DOCKER_YML) down

compose-stop:
	docker-compose -f $(DOCKER_YML) stop

compose-start:
	docker-compose -f $(DOCKER_YML) start

compose-ps:
	docker-compose -f $(DOCKER_YML) ps

# ---------------------------
# Запуск миграций через Goose
# ---------------------------

goose-install:
	go install github.com/pressly/goose/v3/cmd/goose@latest

goose-add:
	goose -dir ./migrations postgres "${PG_URL}" create hw-migration sql

goose-up:
	goose -dir ./migrations postgres "${PG_URL}" up

goose-down:
	goose -dir ./migrations postgres "${PG_URL}" down-to 0

goose-status:
	goose -dir ./migrations postgres "${PG_URL}" status

squawk:
	squawk ./migrations/*

# ---------------------------------
# Запуск кодогенерации через protoc
# ---------------------------------

bin-deps: .vendor-proto
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	GOBIN=$(LOCAL_BIN) go install github.com/envoyproxy/protoc-gen-validate@latest
	GOBIN=$(LOCAL_BIN) go install github.com/rakyll/statik@latest

generate:
	mkdir -p $(OUT_PATH)
	protoc --proto_path api --proto_path vendor.protogen \
		--plugin=protoc-gen-go=$(LOCAL_BIN)/protoc-gen-go$(EXT) --go_out=${OUT_PATH} --go_opt=paths=source_relative \
		--plugin=protoc-gen-go-grpc=$(LOCAL_BIN)/protoc-gen-go-grpc$(EXT) --go-grpc_out=${OUT_PATH} --go-grpc_opt=paths=source_relative \
		--plugin=protoc-gen-grpc-gateway=$(LOCAL_BIN)/protoc-gen-grpc-gateway$(EXT) --grpc-gateway_out=${OUT_PATH} --grpc-gateway_opt paths=source_relative \
		--plugin=protoc-gen-openapiv2=$(LOCAL_BIN)/protoc-gen-openapiv2$(EXT) --openapiv2_out=${OUT_PATH} \
		--plugin=protoc-gen-validate=$(LOCAL_BIN)/protoc-gen-validate$(EXT) --validate_out="lang=go,paths=source_relative:${OUT_PATH}" \
		$(PROTOS_PATH)

.vendor-proto: .vendor-proto/google/protobuf .vendor-proto/google/api .vendor-proto/protoc-gen-openapiv2/options .vendor-proto/validate

.vendor-proto/protoc-gen-openapiv2/options:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/grpc-ecosystem/grpc-gateway vendor.protogen/grpc-ecosystem && \
 		cd vendor.protogen/grpc-ecosystem && \
		git sparse-checkout set --no-cone protoc-gen-openapiv2/options && \
		git checkout
		mkdir -p vendor.protogen/protoc-gen-openapiv2
		mv vendor.protogen/grpc-ecosystem/protoc-gen-openapiv2/options vendor.protogen/protoc-gen-openapiv2
		rm -rf vendor.protogen/grpc-ecosystem

.vendor-proto/google/protobuf:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
		https://github.com/protocolbuffers/protobuf vendor.protogen/protobuf &&\
		cd vendor.protogen/protobuf &&\
		git sparse-checkout set --no-cone src/google/protobuf &&\
		git checkout
		mkdir -p vendor.protogen/google
		mv vendor.protogen/protobuf/src/google/protobuf vendor.protogen/google
		rm -rf vendor.protogen/protobuf

.vendor-proto/google/api:
	git clone -b master --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/googleapis/googleapis vendor.protogen/googleapis && \
 		cd vendor.protogen/googleapis && \
		git sparse-checkout set --no-cone google/api && \
		git checkout
		mkdir -p  vendor.protogen/google
		mv vendor.protogen/googleapis/google/api vendor.protogen/google
		rm -rf vendor.protogen/googleapis

.vendor-proto/validate:
	git clone -b main --single-branch --depth=2 --filter=tree:0 \
		https://github.com/bufbuild/protoc-gen-validate vendor.protogen/tmp && \
		cd vendor.protogen/tmp && \
		git sparse-checkout set --no-cone validate &&\
		git checkout
		mkdir -p vendor.protogen/validate
		mv vendor.protogen/tmp/validate vendor.protogen/
		rm -rf vendor.protogen/tmp


.PHONY: all build deps deps-update run lint clean \
depgraph depgraph-install depgraph-build test test-coverage\
 compose-up compose-down compose-stop compose-start compose-ps\
  goose-install goose-add goose-up goose-down goose-status squawk\
  bin-deps generate .vendor-proto .vendor-proto/protoc-gen-openapiv2/options\
  .vendor-proto/google/protobuf .vendor-proto/google/api .vendor-proto/validate
