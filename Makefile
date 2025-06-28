Q ?= @
ifeq ($(VERBOSE),1)
	Q =
endif

build: clean
	$(Q)go build -o cmd/server/server ./cmd/server
	$(Q)go build -o cmd/client/client ./cmd/client

clean:
	$(Q)rm -f cmd/server/server cmd/client/client

test:
	$(Q)go test ./...

act-test:
	$(Q)act --container-architecture linux/amd664
