BINARY_NAME=bca

build:
	go build -o ${BINARY_NAME} ./cmd/bca/main.go

run: build
	./${BINARY_NAME}

clean:
	go clean
	rm ${BINARY_NAME}
