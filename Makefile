all: test vet fmt lint build_all

test:
	go test ./...

vet:
	go vet ./...

fmt:
	go list -f '{{.Dir}}' ./... | grep -v /vendor/ | xargs -L1 gofmt -l
	test -z $$(go list -f '{{.Dir}}' ./... | grep -v /vendor/ | xargs -L1 gofmt -l)

lint:
	go list ./... | grep -v /vendor/ | xargs -L1 revive -set_exit_status
	# go list ./... | grep -v /vendor/ | xargs -L1 golint -set_exit_status

install-lint-revive:
	go install github.com/mgechev/revive@latest

build:
	go build -o bin/chip8 ./cmd

build_all: build

