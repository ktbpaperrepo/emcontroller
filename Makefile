GOOS=linux
GOARCH=amd64

GIT_COMMIT := $(shell git rev-list -1 HEAD)
BUILD_DATE := $(shell date)

.PHONY: emcontroller
emcontroller:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags "-X 'main.gitCommit=$(GIT_COMMIT)' -X 'main.buildDate=$(BUILD_DATE)'"

.PHONY: clean
clean:
	rm emcontroller