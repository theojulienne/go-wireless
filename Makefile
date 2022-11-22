
build:
	mkdir -p bin
	go build -o bin/api ./cmd/api
	go build -o bin/apscan ./cmd/apscan
	go build -o bin/connectap ./cmd/connectap
	go build -o bin/currentap ./cmd/currentap
	go build -o bin/wpalogs ./cmd/wpalogs
	go build -o bin/wpaspy ./cmd/wpaspy
	go build -o bin/wifistate ./cmd/wifistate
	go build -o bin/ifaces ./cmd/ifaces

test:
	go test ./... -race

wpapi_pkg:
	GOARCH=arm64 go build -o bin/wpapi.arm64 ./cmd/wpapi
	mkdir -p dpkg/wpapi/usr/bin/wpapi
	cp bin/wpapi.arm64 dpkg/wpapi/usr/bin/wpapi
	IAN_DIR=dpkg/wpapi ian set -a arm64
	IAN_DIR=dpkg/wpapi ian pkg

	GOARCH=amd64 go build -o bin/wpapi.amd64 ./cmd/wpapi
	mkdir -p dpkg/wpapi/usr/bin/wpapi
	cp bin/wpapi.amd64 dpkg/wpapi/usr/bin/wpapi
	IAN_DIR=dpkg/wpapi ian set -a amd64
	IAN_DIR=dpkg/wpapi ian pkg