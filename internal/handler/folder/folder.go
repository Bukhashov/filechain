package folder

import (
	"github.com/Bukhashov/filechain/configs"
	"github.com/Bukhashov/filechain/internal/model"
	"github.com/Bukhashov/filechain/internal/dto"
	"github.com/Bukhashov/filechain/pkg/logging"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/gin-gonic/gin"
)

type Folder interface {
	New(c *gin.Context)
	GetAllAddress(c *gin.Context)
}
type folder struct {
	client	*pgxpool.Pool
	logger	*logging.Logger
	config	*configs.Config
	Model 	model.Folder
	Dto		dto.Folder
}

func NewFolder(client *pgxpool.Pool, config *configs.Config, logger *logging.Logger) Folder {
	return &folder{
		client: client,
		config: config,
		logger: logger,
	}
}