test:
	go test ./...

build:
	go env -w GOOS=linux
	go build -ldflags="-w -s" -trimpath ./cmd/emori
	go env -w GOOS=windows

run:
	go env -w CGO_ENABLED=0
	go run ./cmd/emori