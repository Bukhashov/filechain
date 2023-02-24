package main

import (
	"fmt"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/Bukhashov/filechain/pkg/logging"
	"github.com/Bukhashov/filechain/configs"
	"github.com/Bukhashov/filechain/pkg/client/postgresql"
	"github.com/Bukhashov/filechain/api/v1/rest"
	"github.com/Bukhashov/filechain/api/v1/rpc"
)

func main() {
	logger := logging.GetLogger()
	cfg := configs.GetConfig()

	// POSTGRES деректер қорына қосылу
	postgresqlClinet, err := postgresql.NewClient(context.TODO(), cfg.Storage); if err != nil {
		fmt.Print(err)
	}
	// FACE API
	// Суреттегі адам бейнесін өндейтін серверге протокол gRPC мен қосылады
	FaceGFPCClinet, err := grpc.Dial(":5050", grpc.WithTransportCredentials(insecure.NewCredentials())); if err != nil {
		fmt.Print(err);
	}
	c := make(chan int)

	apiRest := rest.NewRest(postgresqlClinet, FaceGFPCClinet, cfg, &logger)
	apiGgpc := rpc.NewRpc(postgresqlClinet, FaceGFPCClinet, cfg, &logger)

	// REST және GRPC API жеке жеке екі поток та жұмыс жасайды
	go apiRest.Run()
	go apiGgpc.Run()
	
	<- c
}