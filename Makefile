SERVICE_NAME=user-api

deps:
	go mod download

build:
	make clean & go build -o bin/${SERVICE_NAME}

build-linux:
	GOOS=linux GOARCH=amd64 go build -o bin/${SERVICE_NAME}

run:
	./bin/${SERVICE_NAME}

start:
	make build & make run

clean:
	rm -rf bin
