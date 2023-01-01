package rpc

import (
	"github.com/Bukhashov/filechain/pkg/logging"
	"github.com/Bukhashov/filechain/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

type Rpc interface {
	Run()
}

type rpc struct {
	client 			*pgxpool.Pool
	faceGrpcClient 	*grpc.ClientConn
	cfg 			*config.Config
	logger 			*logging.Logger
}

func NewRpc(client *pgxpool.Pool, faceGrpcClient *grpc.ClientConn, cfg *config.Config, logger *logging.Logger) Rpc{
	return &rpc{
		client: client,
		faceGrpcClient: faceGrpcClient,
		cfg: cfg,
		logger: logger,
	}
}