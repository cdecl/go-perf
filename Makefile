
GOPATH=$(CURDIR)
GOBIN=$(GOPATH)/bin
GOFILES=$(wildcard src/*.go)
EXEC=perf

build:
	@env GOPATH=$(GOPATH) GOBIN=$(GOBIN) go build -o $(GOBIN)/$(EXEC).exe $(GOFILES)

run:
	@env GOPATH=$(GOPATH) GOBIN=$(GOBIN) go run $(GOFILES)

get:
	-env GOPATH=$(GOPATH) GOBIN=$(GOBIN) go get github.com/shirou/gopsutil
	# -env GOPATH=$(GOPATH) GOBIN=$(GOBIN) go get github.com/StackExchange/wmi
	# -env GOPATH=$(GOPATH) GOBIN=$(GOBIN) go get golang.org/x/sys

cc:
	@env GOPATH=$(GOPATH) GOBIN=$(GOBIN)  GOOS=linux GOARCH=amd64 go build -o $(GOBIN)/$(EXEC) $(GOFILES)
	@env GOPATH=$(GOPATH) GOBIN=$(GOBIN)  GOOS=windows GOARCH=amd64 go build -o $(GOBIN)/$(EXEC).exe $(GOFILES)