package rest

import (
	"fmt"
	"net/http"
	"github.com/Bukhashov/filechain/internal/handler/user"
)

func (r *rest) Run() {

	u := user.NewUser(r.client, r.faceGrpcClient, r.cfg, r.logger)
	
	http.HandleFunc(API_PATH+"/user/singup", u.Singup)
	http.HandleFunc(API_PATH+"/user/singin", u.Singin)
	http.HandleFunc(API_PATH+"/user/delete", u.Delete)


	// b := block.NewBlock(client, cfg, logger)
	// http.HandleFunc(API_PATH+"/new/block", b.New)
	// http.HandleFunc(API_PATH+"/get/block", b.Get)
	// http.HandleFunc(API_PATH+"/get/all/block", b.All)
	// http.HandleFunc(API_PATH+"/update/block", b.Update)

	err := http.ListenAndServe(fmt.Sprintf(":%s", r.cfg.Lesten.Port), nil); if err != nil {
		r.logger.Info(err)
		return
	}
	

}