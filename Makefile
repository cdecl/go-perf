
PROJECT=go-perf
BIN=$(CURDIR)/bin
EXEC=$(PROJECT)


all: build 

build:
	go build -o $(BIN)/$(EXEC)

test:
	go test -v 

dep:
	go mod tidy
	
cc:
	CGO_ENABLED=0 GOOS=linux    GOARCH=amd64  go build -o $(BIN)/$(EXEC) 
	CGO_ENABLED=0 GOOS=windows  GOARCH=amd64  go build -o $(BIN)/$(EXEC).exe

