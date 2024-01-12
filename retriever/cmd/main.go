package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/urfave/cli"
	pb "github.com/zero-gravity-labs/zerog-data-avail/api/grpc/retriever"
	"github.com/zero-gravity-labs/zerog-data-avail/clients"
	"github.com/zero-gravity-labs/zerog-data-avail/common/healthcheck"
	"github.com/zero-gravity-labs/zerog-data-avail/common/logging"
	"github.com/zero-gravity-labs/zerog-data-avail/core/encoding"
	"github.com/zero-gravity-labs/zerog-data-avail/retriever"
	"github.com/zero-gravity-labs/zerog-data-avail/retriever/flags"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	Version   = ""
	GitCommit = ""
	GitDate   = ""
)

func main() {
	app := cli.NewApp()
	app.Version = fmt.Sprintf("%s-%s-%s", Version, GitCommit, GitDate)
	app.Name = "retriever"
	app.Usage = "ZGDA Retriever"
	app.Description = "Service for collecting coded chunks and decode the original data"
	app.Flags = flags.Flags
	app.Action = RetrieverMain
	if err := app.Run(os.Args); err != nil {
		log.Fatalf("application failed: %v", err)
	}

	select {}
}

func RetrieverMain(ctx *cli.Context) error {
	log.Println("Initializing Retriever")
	hostname := ctx.String(flags.HostnameFlag.Name)
	port := ctx.String(flags.GrpcPortFlag.Name)
	addr := fmt.Sprintf("%s:%s", hostname, port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln("could not start tcp listener", err)
	}

	opt := grpc.MaxRecvMsgSize(1024 * 1024 * 300)
	gs := grpc.NewServer(
		opt,
		grpc.ChainUnaryInterceptor(
		// TODO(ian-shim): Add interceptors
		// correlation.UnaryServerInterceptor(),
		// logger.UnaryServerInterceptor(*s.logger.Logger),
		),
	)

	config := retriever.NewConfig(ctx)
	logger, err := logging.GetLogger(config.LoggerConfig)
	if err != nil {
		return err
	}

	encoder, err := encoding.NewEncoder(config.EncoderConfig)
	if err != nil {
		log.Fatalln("could not start tcp listener", err)
	}

	retrievalClient, err := clients.NewRetrievalClient(logger, encoder, config.NumConnections, config.StorageNodeConfig)
	if err != nil {
		log.Fatalln("could not start tcp listener", err)
	}

	retrieverServiceServer := retriever.NewServer(config, logger, retrievalClient, encoder)
	if err = retrieverServiceServer.Start(context.Background()); err != nil {
		log.Fatalln("failed to start retriever service server", err)
	}

	// Register reflection service on gRPC server
	// This makes "grpcurl -plaintext localhost:9000 list" command work
	reflection.Register(gs)

	pb.RegisterRetrieverServer(gs, retrieverServiceServer)

	// Register Server for Health Checks
	healthcheck.RegisterHealthServer(gs)

	log.Printf("server listening at %s", addr)
	return gs.Serve(listener)
}
