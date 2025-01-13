BIN=./bin
EXE=$(BIN)/server
CONFIG_PATH=${HOME}/.proglog

.PHONY: init
init:
	mkdir -p ${CONFIG_PATH}

.PHONY: gencert
gencert:
	cfssl gencert -initca certs/ca-csr.json | cfssljson -bare ca
	cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=certs/ca-config.json -profile=server certs/server-csr.json | cfssljson -bare server
	cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=certs/ca-config.json -profile=client certs/client-csr.json | cfssljson -bare client
	cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=certs/ca-config.json -profile=client -cn="root" certs/client-csr.json | cfssljson -bare root-client
	cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=certs/ca-config.json -profile=client -cn="nobody" certs/client-csr.json | cfssljson -bare nobody-client
	mv *.pem *.csr ${CONFIG_PATH}

$(CONFIG_PATH)/model.conf: init
	cp ./certs/model.conf $(CONFIG_PATH)

$(CONFIG_PATH)/policy.csv: init
	cp ./certs/policy.csv $(CONFIG_PATH)

.PHONY: test
test: $(CONFIG_PATH)/policy.csv $(CONFIG_PATH)/model.conf
	go test -race ./...

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

