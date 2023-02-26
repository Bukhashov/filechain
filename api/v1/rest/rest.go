package rest

import (
	"fmt"
	// "net/http"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Bukhashov/filechain/pkg/logging"
	"github.com/Bukhashov/filechain/configs"
	"github.com/Bukhashov/filechain/internal/handler/user"
	"github.com/Bukhashov/filechain/internal/handler/file"
	"github.com/Bukhashov/filechain/internal/handler/folder"
	
)

type rest struct {
	// router *gin.Engine
	client 			*pgxpool.Pool
	faceGrpcClient 	*grpc.ClientConn
	cfg 			*configs.Config
	logger 			*logging.Logger
}

func (res *rest) Run() {
	router := gin.Default();

	userHandler := user.NewUser(res.client, res.faceGrpcClient, res.cfg, res.logger)
	folderHandler := folder.NewFolder(res.client, res.cfg, res.logger)
	fileHandler := file.NewFile(res.client, *res.cfg, res.logger)

	v1 := router.Group("/api/v1"); {
		authPath := v1.Group("/auth"); {
			authPath.POST("/singup", userHandler.Singup)
			authPath.POST("/singin", userHandler.Singin)
			authPath.POST("/delete", userHandler.Delete)
		}
		folderPath := v1.Group("/folder"); {
			folderPath.POST("/new", folderHandler.New)
			folderPath.POST("/get", folderHandler.New) //###
		}
		filePath := v1.Group("/file"); {
			filePath.POST("/add", fileHandler.Add)
		}
	}

	
	// http.HandleFunc("/new/block", fubc)
	// http.HandleFunc(API_PATH+"/get/block", b.Get)
	// http.HandleFunc(API_PATH+"/get/all/block", b.All)
	// http.HandleFunc(API_PATH+"/update/block", b.Update)
	fmt.Printf("REST API RUN IP %s PORT %s", "127.0.0.1", res.cfg.Lesten.Port)
	
	err := router.Run(fmt.Sprintf(":%s", res.cfg.Lesten.Port)); if err != nil {
        panic("[Error] failed to start Gin server due to: " + err.Error())
        return
    }

	// err := http.ListenAndServe(fmt.Sprintf(":%s", res.cfg.Lesten.Port), nil); if err != nil {
	// 	res.logger.Info(err)
	// 	return
	// }

}

func NewRest(client *pgxpool.Pool, faceGrpcClient *grpc.ClientConn, cfg *configs.Config, logger *logging.Logger) *rest{
	return &rest{
		client: client,
		faceGrpcClient: faceGrpcClient,
		cfg: cfg,
		logger: logger,
	}
}

