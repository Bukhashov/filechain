package rest

import (
	"fmt"
	"net/http"
	"github.com/Bukhashov/filechain/internal/handler/user"
	// "github.com/Bukhashov/filechain/internal/handler/file"
	"github.com/Bukhashov/filechain/internal/handler/folder"
)

func (r *rest) Run() {

	u := user.NewUser(r.client, r.faceGrpcClient, r.cfg, r.logger)
	
	http.HandleFunc(API_PATH+"/user/singup", u.Singup)
	http.HandleFunc(API_PATH+"/user/singin", u.Singin)
	http.HandleFunc(API_PATH+"/user/delete", u.Delete)
	http.HandleFunc(API_PATH+"/user/folder", u.Folder)

	fDer := folder.NewFolder(r.client, r.cfg, r.logger)
	
	http.HandleFunc(API_PATH+"/new/floder", fDer.New)
	http.HandleFunc(API_PATH+"/get/floder", fDer.New)
	http.HandleFunc(API_PATH+"/new/floder", fDer.New)

	// http.HandleFunc(API_PATH+"/new/block", b.New)
	// http.HandleFunc(API_PATH+"/get/block", b.Get)
	// http.HandleFunc(API_PATH+"/get/all/block", b.All)
	// http.HandleFunc(API_PATH+"/update/block", b.Update)

	err := http.ListenAndServe(fmt.Sprintf(":%s", r.cfg.Lesten.Port), nil); if err != nil {
		r.logger.Info(err)
		return
	}
	

}