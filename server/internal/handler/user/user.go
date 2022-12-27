package user

import (
	"net/http"
	"github.com/Bukhashov/filechain/internal/config"
	"github.com/Bukhashov/filechain/pkg/logging"
	"github.com/Bukhashov/filechain/pkg/pb"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

type User interface{
	Singup(w http.ResponseWriter, req *http.Request)
	Singin(w http.ResponseWriter, req *http.Request)
	Delete(w http.ResponseWriter, req *http.Request)
}

type user struct {
	client 	*pgxpool.Pool
	service	pb.FaceClient
	logger	*logging.Logger
	config	config.Config
	ID 		string `json:"id"`
	Name 	string `json:"name"`
	Email	string `json:"email"`
	Dto		Dto
}

func NewUser(client *pgxpool.Pool, cc *grpc.ClientConn, config config.Config, logger *logging.Logger) User { 
	return &user{
		config: config,
		service: pb.NewFaceClient(cc),
		logger: logger,
		client: client,
	}
}