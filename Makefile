server:
	GOOS=darwin GOARCH=amd64 go build -o bin/server-darwin cmd/server/server.go
	GOOS=linux GOARCH=amd64 go build -o bin/server-linux cmd/server/server.go

client:
	GOOS=darwin GOARCH=amd64 go build -o bin/client-darwin cmd/client/client.go
	GOOS=linux GOARCH=amd64 go build -o bin/client-linux cmd/client/client.go

all: client server