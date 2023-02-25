package file

import (
	"github.com/Bukhashov/filechain/configs"
	"github.com/Bukhashov/filechain/pkg/logging"
	"github.com/Bukhashov/filechain/internal/model"
	"github.com/Bukhashov/filechain/internal/dto"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/gin-gonic/gin"
)

type Filechain interface {
	Add(c *gin.Context)
	Get(c *gin.Context)
	All(c *gin.Context)
}

type filechain struct {
	client 	*pgxpool.Pool
	logger	*logging.Logger
	config	configs.Config
	Dto		dto.File
	Model	model.File
}

const (
	TmpFilePath = "./assets/file/tmp/"
)

func NewFile(client *pgxpool.Pool, config configs.Config, logger *logging.Logger) Filechain{
	return &filechain {
		config: config,
		logger: logger,
		client: client,
	}
}