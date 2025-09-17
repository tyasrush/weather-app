APP_NAME=weather-app
WORKER_APP_NAME=weather-app-worker

build-api:
	go build -o $(APP_NAME) ./cmd/restapi

build-worker:
	go build -o $(WORKER_APP_NAME) ./cmd/worker

deps:
	go get -v ./...

run-order-query:
	go run ./cmd/orders/main.go

run-api:
	make build-api
	./$(APP_NAME)

run-worker:
	make build-worker
	./$(WORKER_APP_NAME)

test:
	go test -v -cover ./internal/...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

mock:
	mockery --case snake --all

clean:
	go clean
	rm -f $(APP_NAME)
	rm -f $(WORKER_APP_NAME)
	rm -f coverage.out

docker-compose-up:
	docker compose up -d

docker-compose-down:
	docker compose down
