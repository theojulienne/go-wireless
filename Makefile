
build:
	mkdir -p bin
	go build -o bin/api ./cmd/api
	go build -o bin/apscan ./cmd/apscan
	go build -o bin/connectap ./cmd/connectap
	go build -o bin/wpalogs ./cmd/wpalogs
	go build -o bin/wpaspy ./cmd/wpaspy

test:
	go test ./... -race