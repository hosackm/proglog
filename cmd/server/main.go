package main

import (
	"log"
    "net"
    "os"

	plog "github.com/hosackm/proglog/internal/log"
    "github.com/hosackm/proglog/internal/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// To test, use grpcurl:
// ./bin/server
//
// check the port the is used, it is randomly allocated by OS.
// PORT=55467
//
// echo -n "value 2" | base64
//   dmFsdWUgMg==
//
// Produce a log record
// grpcurl -plaintext -d '{"record": {"value": "dmFsdWUgMg==", "offset": 0}}' \
//    -import-path api -proto api/log.proto localhost:$PORT log.v1.Log/Produce
//
// Consume the log record
// grpcurl -plaintext -d '{"offset": 0}' -import-path api -proto api/log.proto \
// localhost:$PORT log.v1.Log/Consume | jq -r .record.value | base64 --decode
//
// For streaming service you can start the consume stream with:
// grpcurl -plaintext -import-path api -proto api/log.proto localhost:$PORT log.v1.Log/ConsumeStream
// then you can grpcurl to the Produce service
//
// grpcurl -plaintext -d '{"record": {"value": "dmFsdWUgMx=="}}' -import-path api -proto api/log.proto localhost:$PORT log.v1.Log/Produce

func main() {
	l, err := net.Listen("tcp", ":0")
    if err != nil {
        log.Fatal(err)
    }

	clientOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	cc, err := grpc.NewClient(l.Addr().String(), clientOptions...)
    if err != nil {
        log.Fatal(err)
    }

	dir, err := os.MkdirTemp("", "server-test")
    if err != nil {
        log.Fatal(err)
    }

	clog, err := plog.NewLog(dir, plog.Config{})
    if err != nil {
        log.Fatal(err)
    }

    cfg := &server.Config{CommitLog: clog}
	server, err := server.NewGrpcServer(cfg)
    if err != nil {
        log.Fatal(err)
    }

    defer func() {
        server.Stop()
        cc.Close()
        l.Close()
        err = clog.Remove()
        if err != nil {
            log.Fatal(err)
        }
    }()

    log.Println("Serving on:", l.Addr().String())
    err = server.Serve(l)
    if err != nil {
        log.Fatal(err)
    }
}

