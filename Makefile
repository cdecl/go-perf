
PROJECT=go-perf
BIN=$(CURDIR)/bin
EXEC=$(PROJECT)


all: build 

build:
	go build -o $(BIN)/$(EXEC).exe

test:
	go test -v 

dep:
	go mod tidy
	
cc:
	set GOOS=linux& set GOARCH=amd64& go build -o $(BIN)/$(EXEC) 
	set GOOS=windows& set GOARCH=amd64& go build -o $(BIN)/$(EXEC).exe