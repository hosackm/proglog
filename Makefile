BIN=./bin
EXE=$(BIN)/server

all: $(EXE)

$(EXE): cmd/server/main.go
	go build -o $@ $^

compile:
	protoc api/*.proto \
		--go_out=. \
		--python_out=. \
		--pyi_out=. \
		--go-grpc_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		--proto_path=.

