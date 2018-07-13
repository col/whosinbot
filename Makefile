
UNIT_TEST_PACKAGES=$(shell  go list ./... | grep -v "vendor")

build:
	dep ensure
	GOOS=linux go build -o bin/telegram_lambda cmd/telegram/lambda/main.go
	GOOS=linux go build -o bin/telegram_http cmd/telegram/http/main.go

test:
	ENVIRONMENT=test go test $(UNIT_TEST_PACKAGES) -race

deploy: build
	serverless deploy