package user

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/Bukhashov/filechain/configs"
	"github.com/Bukhashov/filechain/internal/model"
	"github.com/Bukhashov/filechain/internal/dto"
	"github.com/Bukhashov/filechain/pkg/logging"
	"github.com/Bukhashov/filechain/pkg/pb"
)

// massageStatusInternalServerError := MassageStatusInternalServerError{}

// "An error occurred on the server, please try again later"

type User interface{
	Singup(c *gin.Context)
	Singin(c *gin.Context)
	Delete(c *gin.Context)
}

type user struct {
	client	*pgxpool.Pool
	service	pb.FaceClient
	logger	*logging.Logger
	config	*configs.Config
	Dto		dto.User
	Model 	model.User
	Token	Token
}

type Token struct {
	Jwt 	string
}

const (
	TmpImagePath = "./assets/image/tmp/"
	FaceImagePath = "./assets/image/face/"
)

func NewUser(client *pgxpool.Pool, cc *grpc.ClientConn, config *configs.Config, logger *logging.Logger) User { 
	return &user{
		config: config,
		service: pb.NewFaceClient(cc),
		logger: logger,
		client: client,
	}
}