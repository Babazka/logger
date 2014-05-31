BIN=logger
.PHONY: $(BIN) test doc

export GOPATH:=$(shell pwd)
export CWD:=$(shell pwd)

$(BIN):
	go install $(BIN)

run: $(BIN)
	bin/$(BIN)

test:
	go test $(BIN)

doc:
	godoc -http=:6060 -goroot=`pwd`


