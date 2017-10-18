TARGETS_NOVENDOR=$(shell glide novendor)

all: msl

msl: cmd/msl/*.go server/*.go jsonrpc/*.go config/*.go wlog/*.go *.go
	go build cmd/msl/msl.go

bundle:
	glide install

check:
	go test $(TARGETS_NOVENDOR)

fmt:
	@echo $(TARGETS_NOVENDOR) | xargs go fmt

clean:
	rm -rf msl
