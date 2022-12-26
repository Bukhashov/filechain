package user

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Bukhashov/filechain/pkg/face"
	"google.golang.org/grpc"
)

func (u *user) Singup(w http.ResponseWriter, r *http.Request){
	var err error
	var file multipart.File
	format := []string {".png", ".jpg"}
	format_is := false
	
	err = r.ParseMultipartForm(32 << 20); if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	u.Dto.Email = r.FormValue("email")
	u.Dto.Name = r.FormValue("name")
	if u.Dto.Email == "" || u.Dto.Name == ""{
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	file, u.Dto.File, err = r.FormFile("img"); if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	
	
	fileFormat := filepath.Ext(u.Dto.File.Filename)
	
	for i := range(format) {
		if format[i] == fileFormat {
			format_is = true
		}
	}
	if !format_is{
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	
	tmpFile, err := os.Create("./assets/tmp/" + u.Dto.File.Filename); if err != nil {
		fmt.Print("tmp file err")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer tmpFile.Close();

	_, err = io.Copy(tmpFile, file); if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	conn, err := grpc.Dial(":5050", grpc.WithInsecure()); if err != nil {
		fmt.Print(err)
	}
	defer conn.Close()

	c := face.NewFaceClient(conn)


	c.Find(context.TODO(), &face.FindRequest{
		FindOneof: &face.FindRequest_Metadata{
			Metadata: &face.Metadata{ 
				Filename: "ff",
				Extension: ".png",
			}, 
		},
	})

}