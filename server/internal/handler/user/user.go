package user

import (
	"time"
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
	client	*pgxpool.Pool
	service	pb.FaceClient
	logger	*logging.Logger
	config	*config.Config
	Dto		Dto
	Model 	UserModel
	Token	Token
}

type Token struct {
	jwt 	string
}

type BadRequrest struct {
	Data	Data
	Massage	string
}
type ResponsData struct {
	Data	Data
	Massage	string
	Token	string
}

type Data struct{
	Accepted 	time.Time
	GiveAway	time.Time
}

const (
	TmpImagePath = "./assets/image/tmp/"
	FaceImagePath = "./assets/image/face/"
)

func NewUser(client *pgxpool.Pool, cc *grpc.ClientConn, config *config.Config, logger *logging.Logger) User { 
	return &user{
		config: config,
		service: pb.NewFaceClient(cc),
		logger: logger,
		client: client,
	}
}