build:
	dep ensure
	GOOS=linux go build -o bin/telegram telegram/main.go
