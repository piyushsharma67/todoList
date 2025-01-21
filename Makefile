.PHONY: todolist
todolist:
	go build -o build/todolist ./cmd/todolist/

.PHONY: test
test:
	go test ./...

ifndef $(GOPATH)
    GOPATH=$(shell go env GOPATH)
    export GOPATH
endif

.PHONY: staticcheck
staticcheck:
	go install honnef.co/go/tools/cmd/staticcheck@latest
	$(GOPATH)/bin/staticcheck ./...