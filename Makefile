PRODUCT_BINARY_NAME=product.out
PROXY_BINARY_NAME=proxy.out

all: build test

build-product:
	go build -tags migrate -o ./cmd/product/${PRODUCT_BINARY_NAME} github.com/thangchung/go-coffeeshop/cmd/product

build-proxy:
	go build -tags migrate -o ./cmd/proxy/${PROXY_BINARY_NAME} github.com/thangchung/go-coffeeshop/cmd/proxy

run-product:
	cd cmd/product && go mod tidy && go mod download && \
	CGO_ENABLED=0 go run -tags migrate github.com/thangchung/go-coffeeshop/cmd/product
.PHONY: run-product

run-counter:
	cd cmd/counter && go mod tidy && go mod download && \
	CGO_ENABLED=0 go run -tags migrate github.com/thangchung/go-coffeeshop/cmd/counter
.PHONY: run-counter

run-barista:
	cd cmd/counter && go mod tidy && go mod download && \
	CGO_ENABLED=0 go run -tags migrate github.com/thangchung/go-coffeeshop/cmd/barista
.PHONY: run-barista

run-proxy:
	cd cmd/proxy && go mod tidy && go mod download && \
	CGO_ENABLED=0 go run -tags migrate github.com/thangchung/go-coffeeshop/cmd/proxy
.PHONY: run-proxy

test:
	go test -v main.go

package:
	docker-compose down --remove-orphans -v
	docker-compose build
.PHONY: package

compose-up: ### Run docker-compose
	docker-compose up --build -d postgres && docker-compose logs -f
.PHONY: compose-up

compose-down: ### Down docker-compose
	docker-compose down --remove-orphans
.PHONY: compose-down

docker-rm-volume: ### remove docker volume
	docker volume rm go-clean-template_pg-data
.PHONY: docker-rm-volume

linter-golangci: ### check by golangci linter
	golangci-lint run
.PHONY: linter-golangci

clean:
	go clean
	rm ${PRODUCT_BINARY_NAME}