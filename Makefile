test: test-main

build:
	@if [ ! -f $(shell which box) ]; \
	then \
		echo "Need to install box to build the docker images we use. Requires root access."; \
		curl -sSL box-builder.sh | sudo bash; \
	fi
	box --no-tty -t migrator box.rb

test-main: build
	docker run --rm -it -v "${GOPATH}/src:/go/src" migrator bash -c "go get -t github.com/erikh/migrator/... && go test -race -v github.com/erikh/migrator -check.v"

shell: build
	docker run --rm -it -w /go/src/github.com/erikh/migrator -v "${GOPATH}/src:/go/src" migrator bash
