package server

import (
	"context"
    "fmt"
	"net"
	"os"
	"testing"

	"github.com/hosackm/proglog/api"
	"github.com/hosackm/proglog/internal/log"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func TestServer(t *testing.T) {
	scenarios := map[string]func(t *testing.T, client api.LogClient){
		"produce/consume a message to/from the log succeeds": testProduceConsume,
		"produce/consume stream suceeds":                     testProduceConsumeStream,
		"consume past log boundary fails": testConsumePastBoundary,
	}
	for scenario, fn := range scenarios {
		t.Run(scenario, func(t *testing.T) {
			client, _, teardown := setupTestServer(t, nil)
			defer teardown()
			fn(t, client)
		})
	}
}

func setupTestServer(t *testing.T, fn func(*Config)) (
	client api.LogClient,
	cfg *Config,
	teardown func(),
) {
	t.Helper()
	l, err := net.Listen("tcp", ":0")
	require.NoError(t, err)

	clientOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	cc, err := grpc.NewClient(l.Addr().String(), clientOptions...)
	require.NoError(t, err)

	dir, err := os.MkdirTemp("", "server-test")
	require.NoError(t, err)

	clog, err := log.NewLog(dir, log.Config{})
	require.NoError(t, err)

	cfg = &Config{CommitLog: clog}
	if fn != nil {
		fn(cfg)
	}

	server, err := NewGrpcServer(cfg)
	require.NoError(t, err)

	go func() {
		err = server.Serve(l)
        if err != nil {
            fmt.Println(err)
        }
	}()

	client = api.NewLogClient(cc)
	return client, cfg, (func() {
		server.Stop()
		cc.Close()
		l.Close()
		err = clog.Remove()
        if err != nil {
            fmt.Println(err)
        }
	})
}

func testProduceConsume(t *testing.T, client api.LogClient) {
	ctx := context.Background()
	want := &api.Record{
		Value: []byte("hello world"),
	}
	produce, err := client.Produce(ctx, &api.ProduceRequest{Record: want})
	require.NoError(t, err)

	consume, err := client.Consume(ctx, &api.ConsumeRequest{Offset: produce.Offset})
	require.NoError(t, err)
	require.Equal(t, want.Value, consume.Record.Value)
	require.Equal(t, want.Offset, consume.Record.Offset)
}

func testConsumePastBoundary(t *testing.T, client api.LogClient) {
	ctx := context.Background()
	produce, err := client.Produce(
		ctx,
		&api.ProduceRequest{
			Record: &api.Record{
				Value: []byte("hello world"),
			},
		},
	)
	require.NoError(t, err)

	consume, err := client.Consume(ctx, &api.ConsumeRequest{
		Offset: produce.Offset + 1,
	})
	if consume != nil {
		t.Fatalf("consume past boundary returned non-nil")
	}

	got := status.Code(err)
	want := status.Code(api.ErrOffsetOutOfRange{}.GRPCStatus().Err())
	require.Equal(t, got, want)
}

func testProduceConsumeStream(t *testing.T, client api.LogClient) {
    ctx := context.Background()
    records := []*api.Record{
        {Value: []byte("first message"), Offset: 0},
        {Value: []byte("second message"), Offset: 1},
    }

    {
        stream, err := client.ProduceStream(ctx)
        require.NoError(t, err)

        for offset, record := range records {
            err = stream.Send(&api.ProduceRequest{Record: record})
            require.NoError(t ,err)
            res, err := stream.Recv()
            require.NoError(t, err)
            require.Equal(t, res.Offset, uint64(offset))
        }
    }

    {
        stream, err := client.ConsumeStream(ctx, &api.ConsumeRequest{Offset: 0})
        require.NoError(t, err)

        for offset, record := range records {
            res, err := stream.Recv()
            require.NoError(t, err)
            require.Equal(t, res.Record, &api.Record{
                Value: record.Value,
                Offset: uint64(offset),
            })
       }
    }
}

