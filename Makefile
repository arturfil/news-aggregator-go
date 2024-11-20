include .env
 
build:
	go build -o ${BINARY} ./cmd/script/main.go

run:
	./${BINARY}

all: build run
