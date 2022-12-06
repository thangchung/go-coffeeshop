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

compose-up:
	docker-compose up --build
.PHONY: compose-up

compose-down:
	docker-compose down --remove-orphans -v
.PHONY: compose-down

compose-build:
	docker-compose down --remove-orphans -v
	docker-compose build
.PHONY: package

test:
	go test -v main.go

linter-golangci: ### check by golangci linter
	golangci-lint run
.PHONY: linter-golangci

clean:
	go clean