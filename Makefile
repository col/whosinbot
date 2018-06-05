build:
	dep ensure
	GOOS=linux go build -o bin/telegram_lambda cmd/telegram/lambda/main.go
	GOOS=linux go build -o bin/telegram_http cmd/telegram/http/main.go