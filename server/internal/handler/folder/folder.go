package folder

import (
	"net/http"
	"time"
	"github.com/Bukhashov/filechain/internal/config"
	"github.com/Bukhashov/filechain/pkg/logging"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Folder interface {
	New(w http.ResponseWriter, req *http.Request)
	GetOne(w http.ResponseWriter, req *http.Request)
	GetAll(w http.ResponseWriter, req *http.Request)
	History(w http.ResponseWriter, req *http.Request)
}
type folder struct {
	client	*pgxpool.Pool
	logger	*logging.Logger
	config	*config.Config
}

type BadRequrest struct {
	Data	Data
	Massage	string
}
type Data struct{
	Accepted 	time.Time
	GiveAway	time.Time
}

func NewFolder(client *pgxpool.Pool, config *config.Config, logger *logging.Logger) Folder {
	return &folder{
		client: client,
		config: config,
		logger: logger,
	}
}