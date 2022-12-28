include .env
export

all: build test

run: run-product run-counter run-barista run-kitchen run-proxy run-web

run-product:
	cd cmd/product && go mod tidy && go mod download && \
	CGO_ENABLED=0 go run github.com/thangchung/go-coffeeshop/cmd/product
.PHONY: run-product

run-counter:
	cd cmd/counter && go mod tidy && go mod download && \
	CGO_ENABLED=0 go run -tags migrate github.com/thangchung/go-coffeeshop/cmd/counter
.PHONY: run-counter

run-barista:
	cd cmd/barista && go mod tidy && go mod download && \
	CGO_ENABLED=0 go run -tags migrate github.com/thangchung/go-coffeeshop/cmd/barista
.PHONY: run-barista

run-kitchen:
	cd cmd/kitchen && go mod tidy && go mod download && \
	CGO_ENABLED=0 go run -tags migrate github.com/thangchung/go-coffeeshop/cmd/kitchen
.PHONY: run-kitchen

run-proxy:
	cd cmd/proxy && go mod tidy && go mod download && \
	CGO_ENABLED=0 go run -tags migrate github.com/thangchung/go-coffeeshop/cmd/proxy
.PHONY: run-proxy

run-web:
	cd cmd/web && go mod tidy && go mod download && \
	CGO_ENABLED=0 go run github.com/thangchung/go-coffeeshop/cmd/web
.PHONY: run-web

docker-compose: docker-compose-stop docker-compose-start
.PHONY: docker-compose

docker-compose-start:
	docker-compose up --build
.PHONY: docker-compose-start

docker-compose-stop:
	docker-compose down --remove-orphans -v
.PHONY: docker-compose-stop

docker-compose-core: docker-compose-core-stop docker-compose-core-start

docker-compose-core-start:
	docker-compose -f docker-compose-core.yaml up --build
.PHONY: docker-compose-core-start

docker-compose-core-stop:
	docker-compose -f docker-compose-core.yaml down --remove-orphans -v
.PHONY: docker-compose-core-stop

docker-compose-build:
	docker-compose down --remove-orphans -v
	docker-compose build
.PHONY: docker-compose-build

wire:
	cd internal/barista/app && wire && cd - && \
	cd internal/counter/app && wire && cd - && \
	cd internal/kitchen/app && wire && cd - && \
	cd internal/product/app && wire && cd -
.PHONY: wire

sqlc:
	sqlc generate
.PHONY: sqlc

test:
	go test -v main.go

linter-golangci: ### check by golangci linter
	golangci-lint run
.PHONY: linter-golangci

clean:
	go clean