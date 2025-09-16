APP_NAME=weather-app

build:
	go build -o $(APP_NAME)

deps:
	go get -v ./...

run:
	make build
	./$(APP_NAME)

test:
	go test -v -cover ./internal/...

test-integration:
	go test -v -cover ./tests/integration/...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

mock:
	mockery --case snake --all

clean:
	go clean
	rm -f $(APP_NAME)
	rm -f coverage.out

docker-start:
	docker compose up --build -d

docker-stop:
	docker compose down
