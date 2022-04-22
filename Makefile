install:
	go mod download
test:
	go clean --testcache
	go test ./...

server-start:
	go run cmd/tcpserver/main.go

client-start:
	go run cmd/tcpclient/main.go

compose-build:
	docker-compose build

compose-start:
	docker-compose up
