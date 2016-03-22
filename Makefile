build:
	go build -o "$@"

install:
	go install

$(GOPATH)/bin/godep:
	go get github.com/tools/godep

$(GOPATH)/bin/github-release:
	go get github.com/aktau/github-release

clean:
	rm -rf bin/

test:
	go test -v .

.PHONY: clean release install test update_internal_version
