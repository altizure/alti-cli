alti-cli: *.go */*.go
	go build -o alti-cli

alti-cli-linux: *.go */*.go
	env GOOS=linux GOARCH=amd64 go build -o alti-cli-linux

dev: *.go */*.go
	go build -race -o alti-cli

lint:
	golint ./...
	golangci-lint run

all: alti-cli alti-cli-linux

clean:
	rm alti-cli alti-cli-linux
