package folder

import (
	"net/http"
	"time"

	"github.com/Bukhashov/filechain/configs"
	"github.com/Bukhashov/filechain/internal/model"
	"github.com/Bukhashov/filechain/pkg/logging"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Folder interface {
	New(w http.ResponseWriter, req *http.Request)
	GetAllAddress(w http.ResponseWriter, req *http.Request)
}
type folder struct {
	client	*pgxpool.Pool
	logger	*logging.Logger
	config	*configs.Config
	Model 	model.Folder
	Dto		Dto
}

type Requrest struct {
	Data	Data
	Addres	string
}
type BadRequrest struct {
	Data	Data
	Massage	string
}
type Data struct{
	Accepted 	time.Time
	GiveAway	time.Time
}

func NewFolder(client *pgxpool.Pool, config *configs.Config, logger *logging.Logger) Folder {
	return &folder{
		client: client,
		config: config,
		logger: logger,
	}
}