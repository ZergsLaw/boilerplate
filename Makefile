lint:
	golangci-lint -v run ./...

clean:
	goimports -local -w .
	go fmt ./...

test:
	time go test ./...

restart-dependencies:
	docker-compose down --volumes
	docker-compose up --build -d postgres rabbit

build:
	rm -rf "bin"
	mkdir "bin"
	GOOS=linux go build -o "bin/" ./cmd/boilerplate

test-integration: restart-dependencies
	time go test ./... -tags=integration
	docker-compose down --volumes

start: build
	docker-compose down --volumes
	docker-compose up --build
