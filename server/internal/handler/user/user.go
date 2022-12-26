package user

import (
	"net/http"
	"github.com/Bukhashov/filechain/internal/config"
	"github.com/Bukhashov/filechain/pkg/logging"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User interface{
	Singup(w http.ResponseWriter, req *http.Request)
	Singin(w http.ResponseWriter, req *http.Request)
	Delete(w http.ResponseWriter, req *http.Request)
}

type user struct {
	client 	*pgxpool.Pool
	logger	*logging.Logger
	config	config.Config
	ID 		string `json:"id"`
	Name 	string `json:"name"`
	Email	string `json:"email"`
	Dto		Dto
}

func NewUser(client *pgxpool.Pool, config config.Config, logger *logging.Logger) User { 
	return &user{
		config: config,
		logger: logger,
		client: client,
	}
}