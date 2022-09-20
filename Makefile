PRODUCT_BINARY_NAME=product.out

all: build test

build:
	go build -o ${PRODUCT_BINARY_NAME} cmd/product/main.go

run:
	go mod tidy && go mod download && \
	CGO_ENABLED=0 go run -tags migrate cmd/product/main.go
.PHONY: run

test:
	go test -v main.go

compose-up: ### Run docker-compose
	docker-compose up --build -d postgresql && docker-compose logs -f
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