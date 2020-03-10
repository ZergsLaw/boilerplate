lint:
	goimports -local -w .
	go fmt ./...
	golangci-lint -v run ./...

test:
	go test ./...

test-integration:
	docker-compose down --volumes
	docker-compose up --build -d postgres rabbit
	go test ./... -tags=integration
	docker-compose down --volumes