BIN_DIR=./build

default: build

build-dir:
	mkdir -p ${BIN_DIR}

build: build-dir
	go build -a -installsuffix cgo -o ${BIN_DIR}/chitchat ./main.go

release: clean
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 make

clean:
	go clean
	rm -rf ${BIN_DIR}

test:
	go test ./...

lint:
	golangci-lint run

