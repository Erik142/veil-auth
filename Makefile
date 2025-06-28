Q ?= @
ifeq ($(VERBOSE),1)
	Q =
endif

build: clean proto
	$(Q)go build -o cmd/server/server ./cmd/server
	$(Q)go build -o cmd/client/client ./cmd/client

proto:
	$(Q)protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative pkg/grpc/auth/auth.proto

clean:
	$(Q)rm -f cmd/server/server cmd/client/client
	$(Q)rm -f pkg/grpc/auth/*.pb.go

test:
	$(Q)go test ./...

act-test:
	$(Q)act --container-architecture linux/amd664
