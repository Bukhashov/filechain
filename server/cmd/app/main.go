package main

import (
	"context"
	"fmt"
	"github.com/Bukhashov/filechain/api/v1/rest"
	"github.com/Bukhashov/filechain/api/v1/rpc"
	"github.com/Bukhashov/filechain/internal/config"
	"github.com/Bukhashov/filechain/pkg/client/postgresql"
	"github.com/Bukhashov/filechain/pkg/logging"
	
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main(){
	logger := logging.GetLogger()
	cfg := config.GetConfig()
	
	// POSTGRES деректер қорына қосылады
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
