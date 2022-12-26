package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Bukhashov/filechain/internal/config"
	"github.com/Bukhashov/filechain/internal/handler/block"
	"github.com/Bukhashov/filechain/internal/handler/user"
	"github.com/Bukhashov/filechain/pkg/client/postgresql"
	"github.com/Bukhashov/filechain/pkg/logging"

	"github.com/jackc/pgx/v5/pgxpool"
	
)

const (
	API_PATH = "/api/v1"
)

func main(){
	logger := logging.GetLogger()
	cfg := config.GetConfig()
	postgresqlClinet, err := postgresql.NewClient(context.TODO(), cfg.Storage); if err != nil {
		fmt.Print(err)
	}
	
	c := make(chan int)

	go RunRESTServer(c, postgresqlClinet, *cfg, &logger)
	go RunGRPCServer(c)
	
	<- c
}

func RunRESTServer(status chan int, client *pgxpool.Pool, cfg config.Config, logger *logging.Logger) {
	u := user.NewUser(client, cfg, logger)
	http.HandleFunc(API_PATH+"/user/singup", u.Singup)
	http.HandleFunc(API_PATH+"/user/singin", u.Singin)
	http.HandleFunc(API_PATH+"/user/delete", u.Delete)
	
	b := block.NewBlock(client, cfg, logger)
	http.HandleFunc(API_PATH+"/new/block", b.New)
	http.HandleFunc(API_PATH+"/get/block", b.Get)
	http.HandleFunc(API_PATH+"/get/all/block", b.All)
	http.HandleFunc(API_PATH+"/update/block", b.Update)
	
	err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.Lesten.Port), nil); if err != nil {
		close(status)
		return
	}
}

func RunGRPCServer(status chan int) {
	// flag.Parse()
	// lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *))

	// s := grpc.NewServer()

}