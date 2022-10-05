PRODUCT_BINARY_NAME=product.out

all: build test

build:
	go build -tags migrate -o ./cmd/product/${PRODUCT_BINARY_NAME} github.com/thangchung/go-coffeeshop/cmd/product

run:
	cd cmd/product && go mod tidy && go mod download && \
	CGO_ENABLED=0 go run -tags migrate github.com/thangchung/go-coffeeshop/cmd/product
.PHONY: run

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