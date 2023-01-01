package rest

import (
	"github.com/Bukhashov/filechain/pkg/logging"
	"github.com/Bukhashov/filechain/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

const (
	API_PATH = "/api/v1"
)

type Rest interface {
	Run()
}

type rest struct {
	client 			*pgxpool.Pool
	faceGrpcClient 	*grpc.ClientConn
	cfg 			*config.Config
	logger 			*logging.Logger
}

func NewRest(client *pgxpool.Pool, faceGrpcClient *grpc.ClientConn, cfg *config.Config, logger *logging.Logger) Rest{
	return &rest{
		client: client,
		faceGrpcClient: faceGrpcClient,
		cfg: cfg,
		logger: logger,
	}
}