IMAGE ?= core-svc
TAG ?= local

.PHONY: all local-build local-run db-build db-up db-down start docs

all: db-build db-up local-build local-run

start: build run

build: gen-proto
	@echo "Building core svc..."
	@mkdir -p .bin
	@cd $(CURDIR) && go build -o .bin/main cmd/main.go

run:
	@echo "Starting core svc..."
	@cd $(CURDIR) && go run cmd/main.go configs/config.yaml

local: local-build local-run

db-build:
	@echo "DB build"
	@cd $(CURDIR) && docker compose build --no-cache

db-up:
	@echo "DB up"
	@cd $(CURDIR) && docker compose up db -d

local-build:
	@echo "Building core svc..."
	@cd $(CURDIR) && go build -o .bin/main cmd/main.go

local-run:
	@echo "Starting core svc..."
	@cd $(CURDIR) && go run cmd/main.go configs/local.yaml

		
local-stop:
	@echo "Stopping core svc..."
	@-lsof -ti:8082 | xargs kill -9 2>/dev/null
	@echo "Port 8082 is free"

db-down:
	@echo "DB down"
	@cd $(CURDIR) && docker compose down -v

gen-proto:
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
	  proto/auth.proto

docs:
	@echo "Generating swagger docs..."
	@swag init -g cmd/main.go -o docs --parseInternal --parseDependency


.PHONY: docker-build
docker-build:
	@echo "Building docker image..."
	@docker build -t $(IMAGE):$(TAG) .