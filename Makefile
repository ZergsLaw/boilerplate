lint:
	golangci-lint -v run ./...

clean:
	goimports -local -w .
	go fmt ./...

test:
	time go test ./...

test-integration:
	docker-compose down --volumes
	docker-compose up --build -d postgres rabbit
	time go test ./... -tags=integration
	docker-compose down --volumes

start:
	rm -rf "bin"
	mkdir "bin"
	docker-compose down --volumes
	GOOS=linux go build -o "bin/" ./cmd/boilerplate
	docker-compose up --build
