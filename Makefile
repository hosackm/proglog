BIN=./bin
EXE=$(BIN)/server

all: $(EXE)

$(EXE):
	go build -o $(EXE) ./cmd/server/main.go

compile:
	protoc api/*.proto \
		--go_out=. \
		--python_out=. \
		--pyi_out=. \
		--go-grpc_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		--proto_path=.

