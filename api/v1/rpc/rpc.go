package rpc

import (
	"google.golang.org/grpc"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/Bukhashov/filechain/pkg/logging"
	"github.com/Bukhashov/filechain/configs"
)

type rpc struct {
	client 			*pgxpool.Pool
	faceGrpcClient 	*grpc.ClientConn
	cfg 			*configs.Config
	logger 			*logging.Logger
}

func (r *rpc) Run() {

}

func NewRpc(client *pgxpool.Pool, faceGrpcClient *grpc.ClientConn, cfg *configs.Config, logger *logging.Logger) *rpc{
	return &rpc{
		client: client,
		faceGrpcClient: faceGrpcClient,
		cfg: cfg,
		logger: logger,
	}
}