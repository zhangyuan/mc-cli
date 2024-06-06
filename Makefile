build:
	go build

clean:
	rm -rf mc-cli
	rm -rf bin/mc-cli-*

build-macos:
	GOOS=darwin GOARCH=amd64 go build -o bin/mc-cli_darwin-amd64
	GOOS=darwin GOARCH=arm64 go build -o bin/mc-cli_darwin-arm64

build-linux:
	GOOS=linux GOARCH=amd64 go build -o bin/mc-cli_linux-amd64

build-windows:
	GOOS=windows GOARCH=amd64 go build -o bin/mc-cli_windows-amd64

build-all: clean build-macos build-linux build-windows

compress-linux:
	upx ./bin/mc-cli_linux*
