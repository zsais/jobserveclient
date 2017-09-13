PACKAGE = server
PACKAGE_PATH = github.com/zsais/jobserveclient/cmd

all: $(PACKAGE)

$(PACKAGE):
	@mkdir -p bin
	CGO_ENABLED=0 GOOS=darwin go build -a -installsuffix cgo -o bin/jsclient $(PACKAGE_PATH)/$(PACKAGE)

setup:
	go get -u github.com/Masterminds/glide
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install

clean:
	rm -rf bin/

run: $(PACKAGE)
	bin/jsclient

test: unittest functionaltest

unittest:
	go test ./jsclient/ -cover

functionaltest:
	./functionaltests/scheduler-darwin-amd64

lint:
	gometalinter ./jsclient ./cmd/...

.PHONY: $(PACKAGE) run unittest lint
