GOC=go build
GOFLAGS=-a -ldflags '-s'
CGOR=CGO_ENABLED=0

all: build

build:
	$(GOC) leaf.go

run:
	go run leaf.go

stat:
	$(CGOR) $(GOC) $(GOFLAGS) leaf.go

fmt:
	gofmt -w .

clean:
	rm leaf
